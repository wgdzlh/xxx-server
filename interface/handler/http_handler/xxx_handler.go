package http_handler

import (
	"xxx-server/application/app"
	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	"xxx-server/domain/entity"
	"xxx-server/domain/repository"
	"xxx-server/interface/dto"
	"xxx-server/interface/resp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type XxxDataHandler struct {
	name string
	repo repository.XxxDataRepo
}

const (
	MULTIPART_DISTRICT_KEY = "district"
)

func NewXxxDataHandler() *XxxDataHandler {
	return &XxxDataHandler{
		name: "XXX数据",
		repo: app.Rs.XxxData,
	}
}

// @Summary 新增XXX数据
// @Tags    XXX数据
// @Accept  json
// @Param   value body dto.NewXxxData true "数据的值"
// @Produce json
// @Success 200 {object} resp.IdResponse
// @Router  /xxx [post]
func (h *XxxDataHandler) Create(c *gin.Context) {
	var (
		tag = h.name + "::新增"
		in  dto.NewXxxData
	)
	err := c.ShouldBindJSON(&in)
	if err != nil {
		errParamResp(c, tag, err)
		return
	}
	item := entity.XxxData{
		District: in.District,
		Path:     in.Path,
	}
	if err = h.repo.Create(&item); err != nil {
		errOpResp(c, tag, err)
		return
	}
	resp.JSON(c, resp.Success, resp.IdData{
		Id: item.Id,
	})
}

// @Summary 批量删除XXX数据
// @Tags    XXX数据
// @Param   ids query []uint64 true "ids"
// @Produce json
// @Success 200 {object} resp.Response
// @Router  /xxx [delete]
func (h *XxxDataHandler) Deletes(c *gin.Context) {
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

// @Summary 获取单个XXX数据详情
// @Tags    XXX数据
// @Param   id path uint64 true "id"
// @Produce json
// @Success 200 {object} resp.Response
// @Router  /xxx/{id} [get]
func (h *XxxDataHandler) QueryDetail(c *gin.Context) {
	var (
		tag  = h.name + "::获取详情"
		item entity.XxxData
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
	log.Info(tag, zap.Uint64("id", id))
	resp.JSON(c, resp.Success, item)
}

// @Summary 查询XXX数据列表
// @Tags    XXX数据
// @Param   page query uint32 false "页码，从1开始"
// @Param   size query uint32 false "页面大小，默认10"
// @Produce json
// @Success 200 {object} resp.Response
// @Router  /xxx [get]
func (h *XxxDataHandler) QueryList(c *gin.Context) {
	var (
		tag   = h.name + "::查询列表"
		items []entity.XxxData
	)
	page := utils.StrToInt(c.Query("page"))
	size := utils.StrToInt(c.Query("size"))
	items, total, err := h.repo.QueryList(page, size)
	if err != nil {
		errOpResp(c, tag, err)
		return
	}
	log.Debug(tag, zap.String("query", c.Request.URL.RawQuery), zap.Int64("total", total))
	resp.JSON(c, resp.Success, resp.ListData{
		List:  items,
		Total: total,
	})
}
