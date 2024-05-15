package disk

import (
	"io"
	"mime/multipart"
	"path/filepath"

	"xxx-server/application/client"
	log "xxx-server/application/logger"
	"xxx-server/domain/entity"
	"xxx-server/infrastructure/config"

	"go.uber.org/zap"
)

const (
	SEAWEED_TIMEOUT = 300 // Secs

	SEAWEED_LOG_TAG = "SeaweedRepo::"
)

type SeaweedRepo struct {
	c      *client.HttpClient
	logTag string
}

type MultipartFormFile struct {
	pw *io.PipeWriter
	mw *multipart.Writer
	fw io.Writer
}

func NewSwClient() *SeaweedRepo {
	r := &SeaweedRepo{
		c:      client.NewHttpClient("SeaweedClient", SEAWEED_TIMEOUT),
		logTag: SEAWEED_LOG_TAG,
	}
	return r
}

func (m *MultipartFormFile) Write(p []byte) (n int, err error) {
	return m.fw.Write(p)
}

func (m *MultipartFormFile) Close() error {
	m.mw.Close()
	return m.pw.Close()
}

func (r *SeaweedRepo) DeleteIfExist(dstDir string) (err error) {
	url := config.C.Ext.SeaweedAddr + dstDir + "?recursive=true"
	resp := entity.AnyJson{}
	err = r.c.Delete(url, nil, &resp)
	log.Info(r.logTag+"DeleteIfExist", zap.Bool("exist", len(resp) == 0))
	return
}

func (r *SeaweedRepo) GetUploadWriter(dstPath string) (w io.WriteCloser, err error) {
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	// must start pipe reader end first, to avoid blocking of write in mw.CreateFormFile
	go func() {
		url := config.C.Ext.SeaweedAddr + dstPath
		log.Info(r.logTag+"PostFile start", zap.String("path", dstPath))
		resp := entity.AnyJson{}
		if e := r.c.PostFile(url, mw.FormDataContentType(), pr, &resp); e != nil {
			log.Error(r.logTag+"PostFile failed", zap.Error(e))
			return
		}
		log.Info(r.logTag+"PostFile end", zap.Any("resp", resp))
	}()
	mf := &MultipartFormFile{
		pw: pw,
		mw: mw,
	}
	if mf.fw, err = mw.CreateFormFile("file", filepath.Base(dstPath)); err != nil {
		log.Error(r.logTag+"failed to create form file field", zap.Error(err))
		return
	}
	w = mf
	return
}
