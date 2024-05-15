package http_handler

import (
	"os"
	"path/filepath"

	"xxx-server/application/app"
	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"
	"xxx-server/infrastructure/config"
	"xxx-server/interface/resp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SettingHandler struct {
	name string
	repo repo.SettingRepo
}

func NewSettingHandler() *SettingHandler {
	return &SettingHandler{
		name: "设置",
		repo: app.Rs.Setting,
	}
}

// @Summary 新增设置项
// @Tags    设置
// @Accept  json
// @Param   section path string         true "设置所属的分栏"
// @Param   value   body entity.AnyJson true "设置的值"
// @Produce json
// @Success 200 {object} resp.IdResponse
// @Router  /setting/{section} [post]
func (h *SettingHandler) Create(c *gin.Context) {
	var (
		tag  = h.name + "::新增"
		item entity.Setting
	)
	err := c.ShouldBindJSON(&item.Value)
	if err != nil {
		errParamResp(c, tag, err)
		return
	}
	item.Section = c.Param("section")
	if err = h.repo.Create(&item); err != nil {
		errOpResp(c, tag, err)
		return
	}
	resp.JSON(c, resp.Success, resp.IdData{
		Id: item.Id,
	})
}

// @Summary 批量删除设置项
// @Tags    设置
// @Param   ids query []uint64 true "ids"
// @Produce json
// @Success 200 {object} resp.Response
// @Router  /setting [delete]
func (h *SettingHandler) Deletes(c *gin.Context) {
	var (
		tag    = h.name + "::批量删除"
		idsStr = c.Query("ids")
	)
	ids := utils.IdsToUint64s(idsStr)
	if len(ids) == 0 {
		errIdResp(c, tag, idsStr)
		return
	}
	if err := h.repo.Delete(ids...); err != nil {
		errOpResp(c, tag, err)
		return
	}
	log.Info(tag, zap.Any("ids", ids))
	resp.JSON(c, resp.Success)
}

// @Summary 获取单个设置项详情
// @Tags    设置
// @Param   id path uint64 true "id"
// @Produce json
// @Success 200 {object} resp.SettingResponse
// @Router  /setting/{id} [get]
func (h *SettingHandler) QueryDetail(c *gin.Context) {
	var (
		tag = h.name + "::获取详情"
	)
	idStr := c.Param("id")
	id := utils.StrToUint64(idStr)
	if id == 0 {
		errIdResp(c, tag, idStr)
		return
	}
	item, err := h.repo.QueryDetail(id)
	if err != nil {
		errOpResp(c, tag, err)
		return
	}
	log.Debug(tag, zap.Uint64("id", id))
	resp.JSON(c, resp.Success, item)
}

// @Summary 查询设置列表
// @Tags    设置
// @Param   section query string false "设置所属的分栏，enum(预警-alert)"
// @Produce json
// @Success 200 {object} resp.SettingsResponse
// @Router  /setting [get]
func (h *SettingHandler) Query(c *gin.Context) {
	var (
		tag = h.name + "::查询列表"
	)
	items, err := h.repo.Query(c.Query("section"))
	if err != nil {
		errOpResp(c, tag, err)
		return
	}
	resp.JSON(c, resp.Success, items)
}

// @Summary 查询特定地区预警设置
// @Tags    设置
// @Param   adcode path string true "adcode"
// @Produce json
// @Success 200 {object} resp.SettingResponse
// @Router  /setting/alert/{adcode} [get]
func (h *SettingHandler) QueryAlertForAdcode(c *gin.Context) {
	var (
		tag = h.name + "::查询地区预警"
	)
	item, err := h.repo.QueryAdcodeDetail(c.Param("adcode"))
	if err != nil {
		errOpResp(c, tag, err)
		return
	}
	resp.JSON(c, resp.Success, item)
}

// @Summary     将包含shp的zip压缩包转为WKT(srid=4326)
// @Description 使用multipart form格式上传zip：curl -X POST http://xxx -F "file=@/home/test/test.zip" -H "Content-Type: multipart/form-data"
// @Tags        设置
// @Accept      multipart/form-data
// @Param       file formData file true "upload target"
// @Produce     json
// @Success     200 {object} resp.Response
// @Router      /setting/shapefile [post]
func (h *SettingHandler) ShpToWkt(c *gin.Context) {
	var (
		tag    = h.name + "::shp转WKT"
		dstDir = filepath.Join(config.C.Server.TmpDir, uuid.NewString())
		upZip  = dstDir + utils.ZIP_EXT
		wkt    string
	)
	file, err := c.FormFile(MULTIPART_FILE_KEY)
	if err != nil {
		errParamResp(c, tag, err)
		return
	}
	defer func() {
		os.Remove(upZip)
		os.RemoveAll(dstDir)
	}()
	log.Info(tag+":uploading file", zap.String("name", file.Filename))
	err = c.SaveUploadedFile(file, upZip)
	if err != nil {
		errOpResp(c, tag+":SaveUploadedFile", err)
		return
	}
	shp, err := utils.GetShpFromZip(upZip, dstDir)
	if err != nil {
		errOpResp(c, tag+":GetShpFromZip", err)
		return
	}
	if shp != "" {
		wkt, err = app.GdalRepo.GetWktFromShp(shp)
		if err != nil {
			errOpResp(c, tag+":GetWktFromShp", err)
			return
		}
	}
	resp.JSON(c, resp.Success, wkt)
}
