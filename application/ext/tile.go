package ext

import (
	"fmt"
	"strconv"
	"strings"

	"xxx-server/application/client"
	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"
	"xxx-server/infrastructure/config"
	"xxx-server/infrastructure/persistence"

	"go.uber.org/zap"
)

const (
	TILE_SERVICE_MIN_ZOOM = 2
	TILE_SERVICE_MAX_ZOOM = 20

	TS_SRID = 4490

	TILE_SERVICE_SOURCE = "data-service"
	TILE_SET_MARK       = "highres"
	TILE_SERVICE_PREFIX = "highres_"

	TILE_SRV_STATUS_SUCCEED = "succeeded"
	TILE_SRV_STATUS_FAILED  = "failed"
)

var (
	TileServiceSrids = []int32{3857, 4490}

	DefaultRealTimeScopes = []int32{100, 101}
)

type TileServiceDelegate struct {
	outSrv repo.OutSvrRepo
	client *client.HttpClient
	logTag string
}

func NewTileServiceRepository() repo.TileServiceRepository {
	return &TileServiceDelegate{
		outSrv: persistence.R.OutSvr,
		client: client.NewHttpClient("TileServiceClient", config.C.Ext.HttpTimeout),
		logTag: "TileServiceDelegate:",
	}
}

func (d *TileServiceDelegate) GetGridServiceStatus(srvId int64) (status string, err error) {
	svrUrl := config.C.Ext.TileServiceGet + "?data_class=grid&id=" + strconv.FormatInt(srvId, 10)
	resp := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Service struct {
				Status string `json:"status"`
			} `json:"service"`
		} `json:"data"`
	}{}
	if err = d.client.Get(svrUrl, nil, &resp); err != nil {
		log.Error(d.logTag+"check grid service failed", zap.Error(err))
		return
	}
	if resp.Code != 0 {
		log.Error(d.logTag+"check grid service err", zap.Int("code", resp.Code), zap.String("msg", resp.Message))
		return
	}
	status = resp.Data.Service.Status
	return
}

func (d *TileServiceDelegate) PublishGridService(metaIds []uint64, bandOrder, mask string, zoom [2]int32) (serviceId int64, err error) {
	tSrids := TileServiceSrids
	if len(config.C.Tile.ServiceSrids) > 0 {
		tSrids = config.C.Tile.ServiceSrids
	}
	param := entity.InitServiceParam{
		ServiceSource: TILE_SERVICE_SOURCE,
		CutClass:      "real_time",
		DataClass:     "grid",
		Srid:          tSrids[0],
		Srids:         tSrids,
		DataIds:       metaIds,
		MinZoom:       zoom[0],
		MaxZoom:       zoom[1],
		ServiceName:   TILE_SERVICE_PREFIX + utils.UInt64sToStr(metaIds, '-') + "_" + utils.GetNowTimeTag(),
	}
	param.GridService.GridSource = "result"
	param.GridService.MaskFile = mask
	param.GridService.Render = "rgba"
	param.GridService.Rgba = BandOrderToRgba(bandOrder)

	gridResult := entity.TileServiceResp{}
	if err = d.client.PostJson(config.C.Ext.TileServiceCreate, param, &gridResult); err != nil {
		return
	}
	if gridResult.Code != 0 {
		log.Error(d.logTag+"PublishGridService failed", zap.String("msg", gridResult.Message))
		err = repo.ErrTileServiceReqFailed
		return
	}
	// data, _ := json.Marshal(param)
	serviceId = gridResult.Data
	return
}

func BandOrderToRgba(bandOrder string) string {
	if (strings.Contains(bandOrder, ",R") || strings.Contains(bandOrder, "R,")) && strings.Contains(bandOrder, "G") && strings.Contains(bandOrder, "B") {
		orders := strings.Split(bandOrder, ",")
		bandIdxMap := map[string]int{}
		for i, o := range orders {
			bandIdxMap[o] = i
		}
		return fmt.Sprintf("%d,%d,%d,-1", bandIdxMap["R"], bandIdxMap["G"], bandIdxMap["B"])
	}
	return "2,1,0,-1"
}
