package persistence

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	repo "xxx-server/domain/repository"

	"github.com/lib/pq"
	"github.com/wgdzlh/gdalib"
	"go.uber.org/zap"
)

const (
	OUT_SRV_LOG_TAG = "OutSvr::"
	RESULT_IMG_TYPE = "RESULTS"

	TILE_ID_RETRIES      = 3000
	TILE_ID_RETRY_PERIOD = time.Second * 20
)

type OutSvrRepoImp struct {
	logTag string
}

type Label struct {
	Type   *int64
	Date   *string
	Region *string
}

func NewOutSvrRepo() repo.OutSvrRepo {
	return &OutSvrRepoImp{
		logTag: OUT_SRV_LOG_TAG,
	}
}

// 获取影像路径信息
func (r *OutSvrRepoImp) GetImgPath(id uint64, isFlight bool) (path string, err error) {
	var pathInfo struct {
		Directory string
		Name      string
	}
	cur := R.im
	if isFlight {
		cur = cur.Raw(`SELECT directory, file_name "name" FROM "` + repo.FlightDomInfoTableName + `" WHERE id = ?`)
	} else {
		cur = cur.Raw(`SELECT directory, name FROM "`+repo.ResultImgInfoTableName+`" WHERE id = ?`, id)
	}
	if err = cur.Take(&pathInfo).Error; err != nil {
		log.Error(r.logTag+"failed to get img path", zap.Error(err))
		return
	}
	path = filepath.Join(pathInfo.Directory, pathInfo.Name)
	return
}

// 获取矢量存储路径
func (r *OutSvrRepoImp) GetVecPath(id uint64) (path string, err error) {
	err = R.ds.Raw(`SELECT vector_file FROM "`+repo.VectorDataTableName+`" WHERE id = ? AND vector_file IS NOT NULL`, id).Scan(&path).Error
	if err == nil && path == "" {
		err = repo.ErrEmptyVectorPath
	}
	if err != nil {
		log.Error(r.logTag+"failed to get vec path", zap.Error(err))
	}
	return
}

// 获取与特定范围交叉的成果影像路径信息
func (r *OutSvrRepoImp) GetImgPathWithIntersects(id uint64, wkt string) (path, bands, cropped string, err error) {
	var pathInfo struct {
		Directory string
		Name      string
		BandOrder string
		Inter     string
	}
	cur := R.im
	if wkt != "" {
		cur = cur.Raw(`SELECT directory, name, band_order, ST_AsText(ST_Intersection(ST_GeomFromText(?, 4326), extent)) "inter" FROM "`+repo.ResultImgInfoTableName+`" WHERE id = ? AND ST_Intersects(ST_GeomFromText(?, 4326), extent)`, wkt, id, wkt)
	} else {
		cur = cur.Raw(`SELECT directory, name, band_order FROM "`+repo.ResultImgInfoTableName+`" WHERE id = ?`, id)
	}
	if err = cur.Scan(&pathInfo).Error; err != nil {
		log.Error(r.logTag+"failed to get img path", zap.Error(err))
		return
	}
	if pathInfo.Directory == "" || pathInfo.Name == "" {
		log.Error(r.logTag + "empty img path")
		err = repo.ErrEmptyRasterPath
		return
	}
	path = filepath.Join(pathInfo.Directory, pathInfo.Name)
	bands = pathInfo.BandOrder
	cropped = pathInfo.Inter
	return
}

// 获取切片服务id
func (r *OutSvrRepoImp) GetTileSvrId(dsId uint64) (svrId uint64, err error) {
	var sId *uint64
	err = R.ds.Raw(`SELECT tile_service_id FROM "`+repo.VectorDataTableName+`" WHERE id = ?`, dsId).Scan(&sId).Error
	if err != nil {
		log.Error(r.logTag+"get tile_service_id failed", zap.Error(err))
		return
	}
	if sId != nil {
		svrId = *sId
	}
	return
}

// 获取地区完整Adcode列表
func (r *OutSvrRepoImp) GetDistrictCodes(adcode string) (codes []string, err error) {
	var ret struct {
		Acroutes pq.Int32Array `gorm:"type:int4[]"`
	}
	err = R.im.Raw(`SELECT acroutes FROM district WHERE adcode = ? LIMIT 1`, adcode).Scan(&ret).Error
	if err == nil && len(ret.Acroutes) == 0 {
		err = repo.ErrInvalidDistrictRoute
	}
	if err != nil {
		log.Error(r.logTag+"failed to get district codes", zap.Error(err))
		return
	}
	codes = make([]string, 0, len(ret.Acroutes))
	for _, code := range ret.Acroutes[1:] {
		codes = append(codes, strconv.FormatInt(int64(code), 10))
	}
	codes = append(codes, adcode)
	return
}

// 获取地区名称
func (r *OutSvrRepoImp) GetDistrictName(adcode string) (name string, err error) {
	err = R.im.Raw(`SELECT name FROM district WHERE adcode = ? LIMIT 1`, adcode).Scan(&name).Error
	if err == nil && name == "" {
		err = repo.ErrInvalidDistrictName
	}
	if err != nil {
		log.Error(r.logTag+"failed to get district name", zap.Error(err))
	}
	return
}

// 获取地区完整名称
func (r *OutSvrRepoImp) GetDistrictPathName(adcode string) (pn string, err error) {
	var ret struct {
		Name     string
		Acroutes pq.Int32Array `gorm:"type:int4[]"`
	}
	if err = R.im.Raw(`SELECT name, acroutes FROM district WHERE adcode = ? LIMIT 1`, adcode).Scan(&ret).Error; err != nil {
		return
	}
	if ret.Name == "" {
		err = repo.ErrInvalidDistrictName
		return
	}
	if len(ret.Acroutes) <= 1 {
		pn = ret.Name
		return
	}
	codeArr := Int32ToIntArray(ret.Acroutes[1:])
	var names []string
	if err = R.im.Raw(`SELECT name FROM district WHERE adcode::int = ANY(?) ORDER BY array_position(?, adcode::int)`,
		codeArr, codeArr).Scan(&names).Error; err != nil {
		return
	}
	pn = strings.Join(append(names, ret.Name), ",")
	return
}

// 获取地区code与名称
func (r *OutSvrRepoImp) GetDistrictCodeOfSpan(span []float64) (code string, err error) {
	if len(span) != 4 {
		return
	}
	var (
		targetWkt = gdalib.PointsToWkt(span[0], span[1], span[2], span[3])
		ret       struct {
			Adcode string
			// Name   string
		}
	)
	err = R.im.Raw(`SELECT adcode, ST_Area(ST_Intersection(geom, ST_GeomFromText(?, 4326))::geography) "area" FROM district WHERE ST_Intersects(geom, ST_GeomFromText(?, 4326)) ORDER BY area DESC, cardinality(acroutes) DESC LIMIT 1`, targetWkt, targetWkt).Scan(&ret).Error
	if err == nil && ret.Adcode == "" {
		err = repo.ErrUnknownDistrict
	}
	if err != nil {
		log.Error(r.logTag+"failed to get district code and name", zap.Error(err))
		return
	}
	code = ret.Adcode
	// name = ret.Name
	return
}

// 获取地区末级子区划编码
func (r *OutSvrRepoImp) GetDistrictSubs(adcode string) (subs []int32, err error) {
	// var tmp []struct {
	// 	Adcode int32
	// 	Parent int32
	// }
	err = R.im.Raw(`SELECT adcode::int4 FROM district WHERE children_num = 0 AND ? = ANY(acroutes) ORDER BY 1`, utils.StrToInt(adcode)).Scan(&subs).Error
	if err != nil {
		log.Error(r.logTag+"failed to get district subs", zap.Error(err))
		// return
	}
	// nt := len(tmp)
	// if nt == 0 {
	// 	return
	// }
	// pMap := map[int32]struct{}{}
	// for _, t := range tmp {
	// 	pMap[t.Parent] = struct{}{}
	// }
	// subs = make([]int32, 0, nt)
	// for _, t := range tmp {
	// 	if _, ok := pMap[t.Adcode]; !ok {
	// 		subs = append(subs, t.Adcode)
	// 	}
	// }
	return
}

func (r *OutSvrRepoImp) GetDistrictAllCodes(adcode string) (subs []int32, err error) {
	code := utils.StrToInt32(adcode)
	err = R.im.Raw(`SELECT adcode::int4 FROM district WHERE ? = ANY(acroutes)`, code).Scan(&subs).Error
	if err != nil {
		log.Error(r.logTag+"failed to get district subs", zap.Error(err))
		return
	}
	subs = append(subs, code)
	return
}

// 获取地区范围geojson
func (r *OutSvrRepoImp) GetDistrictGeoJson(adcode string) (extent string, err error) {
	err = R.im.Raw(`SELECT ST_AsGeoJson(geom) FROM district WHERE adcode = ? LIMIT 1`, adcode).Scan(&extent).Error
	if err == nil && extent == "" {
		err = repo.ErrInvalidDistrictExtent
	}
	if err != nil {
		log.Error(r.logTag+"failed to get district geojson", zap.Error(err))
	}
	return
}

// 获取地区范围wkt
func (r *OutSvrRepoImp) GetDistrictWkt(adcode string) (extent string, err error) {
	err = R.im.Raw(`SELECT ST_AsText(geom) FROM district WHERE adcode = ? LIMIT 1`, adcode).Scan(&extent).Error
	if err == nil && extent == "" {
		err = repo.ErrInvalidDistrictExtent
	}
	if err != nil {
		log.Error(r.logTag+"failed to get district wkt", zap.Error(err))
	}
	return
}

// 获取地区范围简化wkt
func (r *OutSvrRepoImp) GetDistrictSimpWkt(adcode string) (extent string, err error) {
	err = R.im.Raw(`SELECT ST_AsText(simplified_geom) FROM district WHERE adcode = ? LIMIT 1`, adcode).Scan(&extent).Error
	if err == nil && extent == "" {
		err = repo.ErrInvalidDistrictExtent
	}
	if err != nil {
		log.Error(r.logTag+"failed to get district simp wkt", zap.Error(err))
	}
	return
}

// 获取地区范围HexEwkb
func (r *OutSvrRepoImp) GetDistrictHexEwkb(adcode string) (extent string, err error) {
	err = R.im.Raw(`SELECT geom FROM district WHERE adcode = ? LIMIT 1`, adcode).Scan(&extent).Error
	if err == nil && extent == "" {
		err = repo.ErrInvalidDistrictExtent
	}
	if err != nil {
		log.Error(r.logTag+"failed to get district hex ewkb", zap.Error(err))
	}
	return
}

// 获取成果影像DataId
func (r *OutSvrRepoImp) GetResultImageDataId(metaId int64) (dataId int64, err error) {
	err = R.ds.Raw(`SELECT id FROM "`+repo.GridDataTableName+`" WHERE meta_id = ? AND grid_type = 'RESULTS' LIMIT 1`, metaId).Scan(&dataId).Error
	if err == nil && dataId == 0 {
		err = repo.ErrEmptyGridDataId
	}
	return
}

// 获取影像MetaId
func (r *OutSvrRepoImp) GetGridMetaId(dataId int64) (metaId int64, err error) {
	err = R.ds.Raw(`SELECT meta_id FROM "`+repo.GridDataTableName+`" WHERE id = ?`, dataId).Scan(&metaId).Error
	if err == nil && metaId == 0 {
		err = repo.ErrEmptyGridMetaId
	}
	return
}

// 获取成果影像band_order
func (r *OutSvrRepoImp) GetRetImgBandOrder(id uint64) (band string, err error) {
	if err = R.im.Raw(`SELECT band_order FROM "`+repo.ResultImgInfoTableName+`" WHERE id = ?`, id).Scan(&band).Error; err != nil {
		log.Error(r.logTag+"get band order of ret img failed", zap.Uint64("id", id), zap.Error(err))
	}
	return
}

// 获取影像范围wkt
func (r *OutSvrRepoImp) GetImgExtent(id uint64) (wkt string, err error) {
	err = R.im.Raw(`SELECT ST_AsText(extent) FROM "`+repo.ResultImgInfoTableName+`" WHERE id = ?`, id).Scan(&wkt).Error
	if err == nil && wkt == "" {
		err = repo.ErrInvalidImgExtent
	}
	if err != nil {
		log.Error(r.logTag+"failed to get img wkt", zap.Error(err))
	}
	return
}

// 获取影像范围HexEwkb
func (r *OutSvrRepoImp) GetImgExtentHex(id uint64) (hex string, err error) {
	err = R.im.Raw(`SELECT extent FROM "`+repo.ResultImgInfoTableName+`" WHERE id = ?`, id).Scan(&hex).Error
	if err == nil && hex == "" {
		err = repo.ErrInvalidImgExtent
	}
	if err != nil {
		log.Error(r.logTag+"failed to get img hex ewkb", zap.Error(err))
	}
	return
}

// 获取影像免切片服务id
func (r *OutSvrRepoImp) GetTileServiceId(gridId int64) (tileId int64, err error) {
	const (
		FAILED  = "FAIL"
		SUCCESS = "SUCCESS"
	)
	var q = struct {
		TileServiceId     *int64
		ParseStatus       *string
		TileServiceStatus *string
	}{}
	for i := 0; i < TILE_ID_RETRIES; i++ {
		if err = R.ds.Raw(`SELECT tile_service_id, parse_status, tile_service_status FROM "`+repo.GridDataTableName+`" WHERE id = ?`, gridId).Scan(&q).Error; err != nil {
			return
		}
		if q.ParseStatus != nil && *q.ParseStatus == FAILED {
			err = repo.ErrFailedParsing
			return
		}
		if q.TileServiceStatus != nil {
			if *q.TileServiceStatus == FAILED {
				err = repo.ErrFailedTileService
				return
			}
			if *q.TileServiceStatus == SUCCESS && q.TileServiceId != nil && *q.TileServiceId != 0 {
				tileId = *q.TileServiceId
				return
			}
		}
		time.Sleep(TILE_ID_RETRY_PERIOD)
	}
	err = repo.ErrEmptyServiceId
	return
}

// 获取多个影像免切片服务id
func (r *OutSvrRepoImp) GetTileServiceIds(gridIds []int64) (tileIds []int64, err error) {
	if len(gridIds) == 0 {
		return
	}
	var idPairs []struct {
		Id            int64
		TileServiceId int64
	}
	if err = R.ds.Raw(`SELECT id, tile_service_id FROM "`+repo.GridDataTableName+`" WHERE tile_service_id > 0 AND tile_service_status = 'SUCCESS' AND id IN ?`,
		gridIds).Scan(&idPairs).Error; err != nil {
		return
	}
	idMap := make(map[int64]int64, len(idPairs))
	for _, p := range idPairs {
		idMap[p.Id] = p.TileServiceId
	}
	tileIds = make([]int64, len(gridIds))
	for i, gid := range gridIds {
		tileIds[i] = idMap[gid]
	}
	// n := len(finTileIds)
	// if n == len(gridIds) {
	// 	tileIds = finTileIds
	// 	return
	// }
	// if n < entity.MID_OUT_NUM*entity.MID_SRV_FAC {
	// 	err = repo.ErrEmptyServiceId
	// 	return
	// }
	// i := 0
	// for _, gid := range gridIds {
	// 	if gid == 0 {
	// 		tileIds = append(tileIds, 0)
	// 	} else if i < n {
	// 		tileIds = append(tileIds, finTileIds[i])
	// 		i++
	// 	}
	// }
	// if len(tileIds) != len(gridIds) {
	// 	err = repo.ErrEmptyServiceId
	// }
	return
}
