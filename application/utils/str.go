package utils

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"xxx-server/infrastructure/config"

	json "github.com/json-iterator/go"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	dateStrFormat = "20060102150405"

	dayDateTime   = "000000"
	monthDateTime = "01" + dayDateTime
	yearDateTime  = "01" + monthDateTime

	firstQuarter  = "01" + monthDateTime
	secondQuarter = "04" + monthDateTime
	thirdQuarter  = "07" + monthDateTime
	forthQuarter  = "10" + monthDateTime
)

var (
	ErrDateStrMisformed = errors.New("date string misformed")
)

func StrToInt32(s string) int32 {
	i, _ := strconv.Atoi(s)
	return int32(i)
}

func StrToUint32(s string) uint32 {
	i, _ := strconv.Atoi(s)
	return uint32(i)
}

func StrToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func StrToUint64(s string) uint64 {
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}

func StrToInt(s string) int {
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

func StrsToInt64s(ss []string) []int64 {
	ret := make([]int64, len(ss))
	for i, s := range ss {
		ret[i], _ = strconv.ParseInt(s, 10, 64)
	}
	return ret
}

func StrsToFloat64s(ss []string) []float64 {
	ret := make([]float64, len(ss))
	for i, s := range ss {
		ret[i], _ = strconv.ParseFloat(s, 64)
	}
	return ret
}

func Int32sToStr(ids []int32, sep byte) string {
	var ret strings.Builder
	for i, id := range ids {
		if i > 0 {
			ret.WriteByte(sep)
		}
		ret.WriteString(strconv.FormatInt(int64(id), 10))
	}
	return ret.String()
}

func Int64sToStr(ids []int64, sep byte) string {
	var ret strings.Builder
	for i, id := range ids {
		if i > 0 {
			ret.WriteByte(sep)
		}
		ret.WriteString(strconv.FormatInt(id, 10))
	}
	return ret.String()
}

func UInt64sToStr(ids []uint64, sep byte) string {
	var ret strings.Builder
	for i, id := range ids {
		if i > 0 {
			ret.WriteByte(sep)
		}
		ret.WriteString(strconv.FormatUint(id, 10))
	}
	return ret.String()
}

func IdsToUint64s(s string) []uint64 {
	ids := strings.Split(s, config.Q_SEP)
	rets := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if i, e := strconv.ParseUint(id, 10, 64); e == nil {
			rets = append(rets, i)
		}
	}
	return rets
}

func IdsToInt64s(s string) []int64 {
	ids := strings.Split(s, config.Q_SEP)
	rets := make([]int64, 0, len(ids))
	for _, id := range ids {
		if i, e := strconv.ParseInt(id, 10, 64); e == nil {
			rets = append(rets, i)
		}
	}
	return rets
}

func IdsToInt32s(s string) []int32 {
	ids := strings.Split(s, config.Q_SEP)
	rets := make([]int32, 0, len(ids))
	for _, id := range ids {
		if i, e := strconv.ParseInt(id, 10, 32); e == nil {
			rets = append(rets, int32(i))
		}
	}
	return rets
}

func ConcatWithUnderscore(ss ...string) string {
	return strings.Join(ss, "_")
}

func ConcatColumns(cols []string) string {
	return `"` + strings.Join(cols, `","`) + `"`
}

func CstStrToTime(s, format string) (time.Time, error) {
	return time.ParseInLocation(format, s, config.TZ)
}

func StrToTime(s string) (time.Time, error) {
	return time.ParseInLocation(config.Q_TIME_FORMAT, s, config.TZ)
}

func TimeToStr(t time.Time) string {
	return t.Format(config.Q_TIME_FORMAT)
}

func DateStrToTime(s string) (t time.Time, err error) {
	var suffix string
	if len(s) > 14 {
		s = s[:14]
	}
	switch len(s) {
	case 4:
		suffix = yearDateTime
	case 6:
		switch s[4] {
		case 'Q':
			switch s[5] {
			case '1':
				suffix = firstQuarter
			case '2':
				suffix = secondQuarter
			case '3':
				suffix = thirdQuarter
			case '4':
				suffix = forthQuarter
			default:
				err = ErrDateStrMisformed
				return
			}
			s = s[:4]
		case '0', '1':
			suffix = monthDateTime
		default:
			err = ErrDateStrMisformed
			return
		}
	case 8:
		suffix = dayDateTime
	case 14:
	default:
		err = ErrDateStrMisformed
		return
	}
	return time.ParseInLocation(dateStrFormat, s+suffix, config.TZ)
}

func FindLastYear(input string, refs []string) (lastYearIdx int) {
	lastYearIdx = -1
	if len(input) < 4 {
		return
	}
	inputInt, err := strconv.Atoi(input[:4])
	if err != nil {
		return
	}
	lastYear := strconv.Itoa(inputInt - 1)
	if len(lastYear) < 4 {
		return
	}
	lastYear += input[4:]
	for i := len(refs) - 1; i >= 0; i-- {
		if refs[i] == lastYear {
			lastYearIdx = i
			break
		}
	}
	return
}

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func S2B(s string) []byte {
	const MaxInt32 = 1<<31 - 1
	return (*[MaxInt32]byte)(unsafe.Pointer((*reflect.StringHeader)(
		unsafe.Pointer(&s)).Data))[: len(s)&MaxInt32 : len(s)&MaxInt32]
}

func TrimTailComma(s string) string {
	return strings.TrimRight(s, ",")
}

func ReplaceUnderline(s string) string {
	return strings.ReplaceAll(s, "_", "-")
}

func GetNowTimeTag() string {
	const tf = "20060102150405.000"
	t := time.Now().Format(tf)
	return t[:len(tf)-4] + t[len(tf)-3:]
}

func GetJsonBody(data any) (body io.Reader, err error) {
	bs, err := json.Marshal(data)
	if err != nil {
		return
	}
	body = bytes.NewReader(bs)
	return
}

func StrToHex(s string) string {
	src := S2B(s)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return B2S(dst)
}

func BsToHex(bs []byte) string {
	dst := make([]byte, hex.EncodedLen(len(bs)))
	hex.Encode(dst, bs)
	return B2S(dst)
}

func FindInStrs(ss []string, target string) int {
	for i, s := range ss {
		if s == target {
			return i
		}
	}
	return -1
}

func FindsInStrs(ss []string, targets []string) []int {
	ret := make([]int, len(targets))
OUT:
	for i, t := range targets {
		for j, s := range ss {
			if s == t {
				ret[i] = j
				continue OUT
			}
		}
		ret[i] = -1
	}
	return ret
}

func Contains(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}

func ContainsAll(group, sub []string) bool {
OUT:
	for _, s := range sub {
		for _, a := range group {
			if a == s {
				continue OUT
			}
		}
		return false
	}
	return true
}

func ContainsAny(group, sub []string) bool {
	for _, s := range sub {
		for _, a := range group {
			if a == s {
				return true
			}
		}
	}
	return false
}

func RemoveFromSs(ss []string, target string) []string {
	ret := make([]string, 0, len(ss))
	for _, s := range ss {
		if s == target {
			continue
		}
		ret = append(ret, s)
	}
	return ret
}

func ContainsId(ids []uint64, target uint64) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}

func RemoveFromIds(ids []uint64, target uint64) []uint64 {
	ret := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == target {
			continue
		}
		ret = append(ret, id)
	}
	return ret
}

// GBK 转 UTF-8
func GbkToUtf8(s []byte) (d []byte, e error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e = io.ReadAll(reader)
	return
}

// UTF-8 转 GBK
func Utf8ToGbk(s []byte) (d []byte, e error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e = io.ReadAll(reader)
	return
}

// GBK string 转 UTF-8
func GbkStrToUtf8(s string) (d string, e error) {
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	t, e := io.ReadAll(reader)
	if e != nil {
		return
	}
	d = B2S(t)
	return
}

// UTF-8 string 转 GBK
func Utf8StrToGbk(s string) (d string, e error) {
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	t, e := io.ReadAll(reader)
	if e != nil {
		return
	}
	d = B2S(t)
	return
}

func PurifyForUtf8(s string) string {
	return strings.ToValidUTF8(strings.ReplaceAll(s, "\x00", ""), "")
}

func YearToPeriod(year string) []string {
	return []string{year + "-01-01", year + "-12-31"}
}

func SimplifyCronExp(cron string) string {
	ss := strings.Split(cron, " ")
	if len(ss) >= 6 {
		return strings.Join(ss[1:6], " ")
	} else if len(ss) == 5 {
		return cron
	}
	return ""
}
