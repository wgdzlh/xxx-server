package repository

import (
	"errors"

	"xxx-server/domain/entity"
)

type WorkflowRepository interface {
	XxxTheme(taskId string, input *entity.XxxThemeInput) (err error)
}

type DataStoreRepository interface {
	StoreShape(shp string, fields ...string) (vectorId int64, err error)
	StoreRaster(path, params string) (gridMetaId int64, err error)
}

type TileServiceRepository interface {
	GetGridServiceStatus(srvId int64) (status string, err error)
	PublishGridService(metaIds []uint64, bandOrder, mask string, zoom [2]int32) (serviceId int64, err error)
}

var (
	ErrInvalidVectorId = errors.New("invalid vector id")
	ErrInvalidGridId   = errors.New("invalid grid id")

	ErrInvalidDataSetId    = errors.New("invalid data set id")
	ErrInvalidServiceSetId = errors.New("invalid service set id")

	ErrUnmatchedImgSrvIds = errors.New("unmatched img and service ids")

	ErrEmptyShp             = errors.New("empty shp")
	ErrEmptyTif             = errors.New("empty tif")
	ErrNotSupportedFileType = errors.New("file type not supported")
	ErrNotCached            = errors.New("id not found in cache")

	ErrTileServiceReqFailed = errors.New("tile service req failed")
)
