package ext

import (
	"os"
	"path/filepath"
	"strconv"

	"xxx-server/application/client"
	log "xxx-server/application/logger"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"
	"xxx-server/infrastructure/config"

	json "github.com/json-iterator/go"
	"github.com/wgdzlh/gdalib"
	"go.uber.org/zap"
)

const (
	DATA_STORAGE_TYPE = "DATA_SERVICE"
)

var (
	xxxTifExtent = gdalib.PointsToWkt(8133511.7542897, 15081375.1438741, 1970599.25164441, 7188255.41580702)

	defaultLow  = 0.00
	defaultHigh = 1.00
)

type DataStoreDelegate struct {
	client *client.HttpClient
	logTag string
}

func NewDataStoreRepository() repo.DataStoreRepository {
	d := &DataStoreDelegate{
		client: client.NewHttpClient("DataStoreClient", config.C.Ext.HttpTimeout),
		logTag: "DataStoreDelegate:",
	}
	config.AddCallback(d.changeClientTimeout)
	return d
}

// 将矢量入库到综合服务平台
func (d *DataStoreDelegate) StoreShape(shp string, fields ...string) (vectorId int64, err error) {
	req := entity.ShapeEntry{
		ShpFilePath:    shp,
		StorageType:    DATA_STORAGE_TYPE,
		ServiceColumns: fields,
	}

	resp := entity.VectorStorageResp{}

	err = d.client.PostJson(config.C.Ext.VectorStorage, req, &resp)
	vectorId = resp.Data.VectorStorageId
	if err == nil && vectorId <= 0 {
		err = repo.ErrInvalidVectorId
	}
	log.Info(d.logTag+"store shape", zap.Any("req", req), zap.Any("resp", resp))
	return
}

// 将栅格入库到综合服务平台
func (d *DataStoreDelegate) StoreRaster(path, params string) (gridMetaId int64, err error) {
	var imgInfo entity.ParsedImgInfo
	if params == "" {
		imgInfo.Name = filepath.Base(path)
		imgInfo.Space = 3857
	} else {
		var bs []byte
		bs, err = os.ReadFile(params)
		if err != nil {
			log.Error(d.logTag+"read raster info json failed", zap.String("path", params), zap.Error(err))
			return
		}
		if err = json.Unmarshal(bs, &imgInfo); err != nil {
			return
		}
	}

	req := entity.GridStorageReq{
		FileName:     imgInfo.Name,
		Path:         path,
		StorageType:  DATA_STORAGE_TYPE,
		ProductLevel: imgInfo.ProductLevel,
		Space:        strconv.FormatInt(imgInfo.Space, 10),
		Proj:         imgInfo.Proj,
		Extent:       imgInfo.Extent,
		AdCode:       imgInfo.AdCode,
		Resolution:   imgInfo.Resolution,
		BandCount:    imgInfo.BandCount,
		BandOrder:    imgInfo.BandOrder,
	}

	resp := entity.GridStorageResp{}

	err = d.client.PostJson(config.C.Ext.GridStorage, req, &resp)
	gridMetaId = resp.Data.ResultGridMetaId
	if err == nil && gridMetaId <= 0 {
		err = repo.ErrInvalidGridId
	}
	log.Info(d.logTag+"store grid", zap.String("path", path), zap.Any("resp", resp))
	return
}

func (d *DataStoreDelegate) changeClientTimeout(sc *config.SelfConfig) {
	if sc.Ext.HttpTimeout != config.C.Ext.HttpTimeout {
		d.client.SetTimeout(sc.Ext.HttpTimeout)
	}
}
