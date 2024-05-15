package repository

import (
	"errors"
)

const (
	GPDistrictTableName      = "district"
	OriginalImgInfoTableName = "t_original_grid_meta_image_info"
	OriginalPkgInfoTableName = "t_original_grid_meta_package_info"

	ResultImgInfoTableName = "t_result_grid_meta_image_info"
	VectorInfoTableName    = "t_vector_meta_info"
	GridDataTableName      = "t_grid_storage"
	VectorDataTableName    = "t_vector_storage"

	FlightDomInfoTableName = "t_result_flight_meta_image_info"

	FileDescTableName = "t_file_describe_backup"

	LabelTableName = "t_label_info"

	DataSetRelationTableName = "t_data_set_file_storage_relation"
)

type OutSvrRepo interface {
	GetImgPath(id uint64, isFlight bool) (path string, err error)
	GetVecPath(id uint64) (path string, err error)
	GetImgPathWithIntersects(id uint64, wkt string) (path, bands, cropped string, err error)
	GetTileSvrId(dsId uint64) (svrId uint64, err error)
	GetDistrictCodes(adcode string) (codes []string, err error)
	GetDistrictName(adcode string) (name string, err error)
	GetDistrictGeoJson(adcode string) (extent string, err error)
	GetDistrictWkt(adcode string) (extent string, err error)
	GetDistrictSimpWkt(adcode string) (extent string, err error)
	GetDistrictHexEwkb(adcode string) (extent string, err error)
	GetResultImageDataId(metaId int64) (dataId int64, err error)
	GetGridMetaId(dataId int64) (metaId int64, err error)
	GetRetImgBandOrder(id uint64) (band string, err error)
	GetImgExtent(id uint64) (wkt string, err error)
	GetImgExtentHex(id uint64) (hex string, err error)
	GetTileServiceId(gridId int64) (tileId int64, err error)

	GetTileServiceIds(gridIds []int64) (tileIds []int64, err error)
	GetDistrictCodeOfSpan(span []float64) (code string, err error)
	GetDistrictSubs(adcode string) (subs []int32, err error)
	GetDistrictAllCodes(adcode string) (subs []int32, err error)
	GetDistrictPathName(adcode string) (pn string, err error)
}

var (
	ErrInvalidDistrictRoute  = errors.New("invalid district route")
	ErrInvalidDistrictName   = errors.New("invalid district name")
	ErrInvalidDistrictExtent = errors.New("invalid district extent")
	ErrUnknownDistrict       = errors.New("unknown district")
	ErrInvalidImgExtent      = errors.New("invalid img extent")
	ErrFailedParsing         = errors.New("failed parsing")
	ErrFailedTileService     = errors.New("failed tile service")
	ErrRasterNotFound        = errors.New("raster not found")
	ErrEmptyServiceId        = errors.New("empty service id")
	ErrEmptyGridDataId       = errors.New("empty grid data id")
	ErrEmptyGridMetaId       = errors.New("empty grid meta id")
	ErrEmptyVectorPath       = errors.New("empty vector path")
	ErrEmptyRasterPath       = errors.New("empty raster path")
)
