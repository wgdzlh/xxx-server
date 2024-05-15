package utils

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
)

const (
	FILE_EXT_SHP  = ".shp"
	FILE_EXT_CPG  = ".cpg"
	FILE_EXT_PNG  = ".png"
	FILE_EXT_JSON = ".json"

	FILE_EXT_TXT = ".txt"

	UTF8  = "UTF8"
	UTF_8 = "UTF-8"

	degToRad = math.Pi / 180
)

var (
	ErrNoShpInZip = errors.New("no shp in zip")
	ErrDimNotMet  = errors.New("dimension of input data not met")
)

func GetUniqSubDir(parentPath string) (path string, err error) {
	path = filepath.Join(parentPath, uuid.NewString())
	err = os.Mkdir(path, os.ModePerm)
	return
}

func GetDateSubDir(parentPath, date string) (path string, err error) {
	path = filepath.Join(parentPath, date)
	err = os.MkdirAll(path, os.ModePerm)
	return
}

func GetFilenameWithoutExt(path string) (name string) {
	name = filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(path))
	return
}

func GetShpInZip(zipFile, dstDir string) (path string, utf8 bool, err error) {
	shpFiles, err := Unzip(zipFile, dstDir)
	if err != nil {
		return
	}
	os.Remove(zipFile)
	for _, file := range shpFiles {
		if strings.HasSuffix(file, FILE_EXT_SHP) {
			path = file
			continue
		}
		if strings.HasSuffix(file, FILE_EXT_CPG) {
			enc, e := os.ReadFile(file)
			if e == nil && len(enc) > 0 {
				encStr := strings.ToUpper(string(enc))
				utf8 = encStr == UTF_8 || encStr == UTF8
			}
		}
	}
	if path == "" {
		err = ErrNoShpInZip
	}
	return
}

func GetDistrictInShpName(shp string) (district string) {
	district = strings.TrimSuffix(filepath.Base(shp), FILE_EXT_SHP)
	return
}

func GetBasicBandIdx(bandOrder string) (idx [3]string, invalid bool) {
	bands := strings.Split(bandOrder, ",")
	for i, b := range bands {
		switch b {
		case "R":
			idx[0] = strconv.Itoa(i + 1)
		case "G":
			idx[1] = strconv.Itoa(i + 1)
		case "B":
			idx[2] = strconv.Itoa(i + 1)
		}
	}
	for _, b := range idx {
		if b == "" {
			invalid = true
			break
		}
	}
	return
}

func RetainOnlyDigits(s string) string {
	var sb strings.Builder
	for _, b := range s {
		if b >= '0' && b <= '9' {
			sb.WriteRune(b)
		}
	}
	return sb.String()
}

func GetShpFromZip(zipFile, dstDir string) (shp string, err error) {
	shpFiles, err := Unzip(zipFile, dstDir)
	if err != nil {
		return
	}
	for _, file := range shpFiles {
		if strings.HasSuffix(file, FILE_EXT_SHP) {
			shp = file
			break
		}
	}
	return
}

func OutputWindPng(dst string, width, height int, wd, ws []int16) (meta []byte, err error) {
	pngF, err := os.Create(dst)
	if err != nil {
		return
	}
	defer pngF.Close()
	// jsonF, err := os.Create(dst + FILE_EXT_JSON)
	// if err != nil {
	// 	return
	// }
	// defer jsonF.Close()
	return UploadWindPng(pngF, width, height, wd, ws)
}

func UploadWindPng(pngW io.Writer, width, height int, wd, ws []int16) (meta []byte, err error) {
	n := width * height
	if len(wd) != n || len(ws) != n {
		err = ErrDimNotMet
		return
	}
	var (
		// Create a colored image of the given width and height.
		img              = image.NewNRGBA(image.Rect(0, 0, width, height))
		uu               = make([]float32, n)
		vv               = make([]float32, n)
		rad              float64
		min_u, min_v     float32 = math.MaxFloat32, math.MaxFloat32
		max_u, max_v             = -min_u, -min_v
		ratio_u, ratio_v float32
	)
	for i := range wd {
		rad = float64(wd[i]) / 10.0 * degToRad
		uu[i] = -float32(float64(ws[i]) * math.Sin(rad))
		vv[i] = -float32(float64(ws[i]) * math.Cos(rad))
		if uu[i] < min_u {
			min_u = uu[i]
		}
		if uu[i] > max_u {
			max_u = uu[i]
		}
		if vv[i] < min_v {
			min_v = vv[i]
		}
		if vv[i] > max_v {
			max_v = vv[i]
		}
	}
	if max_u > min_u {
		ratio_u = 255.0 / (max_u - min_u)
	}
	if max_v > min_v {
		ratio_v = 255.0 / (max_v - min_v)
	}
	meta, err = json.Marshal(struct {
		MaxU float32 `json:"max_u"`
		MinU float32 `json:"min_u"`
		MaxV float32 `json:"max_v"`
		MinV float32 `json:"min_v"`
	}{max_u, min_u, max_v, min_v})
	if err != nil {
		return
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := x + y*width
			img.Set(x, y, color.NRGBA{
				R: uint8(ratio_u * (uu[i] - min_u)),
				G: uint8(ratio_v * (vv[i] - min_v)),
				B: 0,
				A: 255,
			})
		}
	}
	err = png.Encode(pngW, img)
	return
}
