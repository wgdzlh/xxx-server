package http_handler

import (
	"os"
	"strings"
	"sync"

	"xxx-server/application/app"
	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"
	"xxx-server/interface/resp"

	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type TaskHandler struct {
	name         string
	repo         repo.CronTaskRepo
	restartMutex sync.Mutex
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		name: "后台任务",
		repo: app.Rs.CronTask,
	}
}

// @Summary 批量删除后台任务
// @Tags    后台任务
// @Param   ids query []uint64 true "ids"
// @Produce json
// @Success 200 {object} resp.Response
// @Router  /task [delete]
func (h *TaskHandler) Deletes(c *gin.Context) {
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

// @Summary 获取单个后台任务详情
// @Tags    后台任务
// @Param   id path uint64 true "id"
// @Produce json
// @Success 200 {object} resp.CronTaskResponse
// @Router  /task/{id} [get]
func (h *TaskHandler) QueryDetail(c *gin.Context) {
	var (
		tag  = h.name + "::获取详情"
		item *entity.CronTask
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

// @Summary 查询后台任务列表
// @Tags    后台任务
// @Param   id         query []uint64 false "ID，支持多值"
// @Param   genre      query string   false "任务大类型：enum(自动任务-cron,下载任务-download)"
// @Param   type       query string   false "任务子类型：enum(短临降雨-short,7天预报-mid,专题图-factor,反演报告-deduce)"
// @Param   adcode     query string   false "行政区划编码"
// @Param   factor     query string   false "要素类型"
// @Param   status     query string   false "任务状态：enum(NotStarted,InProc,Done,Failed)"
// @Param   created_at query []string false "创建时间（起止时间格式：2022-07-20 10:00:00,2022-07-21 10:00:00），只传一个表示不限制另一个"
// @Param   page       query uint32   false "页码，从1开始"
// @Param   size       query uint32   false "页面大小，默认10"
// @Produce json
// @Success 200 {object} resp.CronTasksResponse
// @Router  /task [get]
func (h *TaskHandler) Query(c *gin.Context) {
	var (
		tag   = h.name + "::查询列表"
		items []*entity.CronTask
	)
	filter := getQueryFilter(c)
	items, total, err := h.repo.QueryList(filter)
	if err != nil {
		errOpResp(c, tag, err)
		return
	}
	log.Debug(tag, zap.Any("filter", filter), zap.Int64("total", total))
	resp.JSON(c, resp.Success, resp.ListData{
		List:  items,
		Total: total,
	})
}

// @Summary 重启后台任务
// @Tags    后台任务
// @Param   id path uint64 true "id"
// @Produce json
// @Success 200 {object} resp.Response
// @Router  /task/{id} [put]
func (h *TaskHandler) Restart(c *gin.Context) {
	var (
		tag = h.name + "::重启"
		ct  *entity.CronTask
	)
	idStr := c.Param("id")
	id := utils.StrToUint64(idStr)
	if id == 0 {
		errIdResp(c, tag, idStr)
		return
	}
	h.restartMutex.Lock()
	defer h.restartMutex.Unlock()
	ct, err := h.repo.QueryDetail(id)
	if err != nil {
		errOpResp(c, tag, err)
		return
	}
	needRestart := ct.Genre == entity.TASK_GENRE_CRON && ct.Status == entity.TASK_FAILED
	if needRestart {
		if err = app.SchedulerSvr.RestartJob(ct); err != nil {
			errOpResp(c, tag, err)
			return
		}
	}
	log.Info(tag, zap.Uint64("id", id), zap.Bool("restarted", needRestart))
	resp.JSON(c, resp.Success)
}

// @Summary 申请下载任务数据
// @Tags    后台任务
// @Accept  json
// @Param   id             path uint64                true "id"
// @Param   factorThemeReq body entity.FactorThemeReq true "下载任务参数"
// @Produce json
// @Success 200 {object} resp.IdResponse
// @Router  /task/{id}/download [post]
func (h *TaskHandler) Download(c *gin.Context) {
	var (
		tag     = h.name + "::申请下载"
		idStr   = c.Param("id")
		req     entity.FactorThemeReq
		dTaskId uint64
	)
	id := utils.StrToUint64(idStr)
	if id == 0 {
		errIdResp(c, tag, idStr)
		return
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		errParamResp(c, tag, err)
		return
	}
	if req.DistrictLevel != "province" && req.DistrictLevel != "city" {
		errParamResp(c, tag, repo.ErrInvalidDistrictLevel)
		return
	}
	ct, err := h.repo.QueryDetail(id)
	if err != nil {
		errOpResp(c, tag, err)
		return
	}
	downloadable := ct.Genre == entity.TASK_GENRE_CRON && ct.Status == entity.TASK_DONE
	if downloadable {
		if dTaskId, err = app.SchedulerSvr.DownloadJob(ct, &req); err != nil {
			errOpResp(c, tag, err)
			return
		}
	}
	log.Info(tag, zap.Uint64("id", id), zap.Bool("started", downloadable))
	resp.JSON(c, resp.Success, resp.IdData{
		Id: dTaskId,
	})
}

// @Summary 下载任务数据
// @Tags    后台任务
// @Param   id path uint64 true "id"
// @Produce application/zip
// @Success 200 {file} name.zip
// @Failure 404
// @Failure 400
// @Failure 500
// @Router  /task/{id}/zip [get]
func (h *TaskHandler) GetZip(c *gin.Context) {
	var (
		tag   = h.name + "::下载数据"
		idStr = c.Param("id")
	)
	id := utils.StrToUint64(idStr)
	if id == 0 {
		errIdRespStatus(c, tag, idStr)
		return
	}
	ct, err := h.repo.QueryDetail(id)
	if err != nil {
		errOpRespStatus(c, tag, err)
		return
	}
	canDownload := ct.Genre == entity.TASK_GENRE_DOWN && ct.Status == entity.TASK_DONE
	if !canDownload {
		errOpRespStatus(c, tag, repo.ErrTaskNotReady)
		return
	}
	targetZip := ct.Dir + utils.ZIP_EXT
	if _, err = os.Stat(targetZip); err != nil {
		errNotFoundRespStatus(c, tag, repo.ErrZipNotExist)
		return
	}
	var tp entity.FactorThemeReq
	json.Unmarshal(ct.Params, &tp)
	factors := make([]string, len(tp.Themes))
	for i := range factors {
		factors[i] = tp.Themes[i].Factor
	}
	dZipName := utils.ConcatWithUnderscore(ct.Name, tp.District, strings.Join(factors, "-"), utils.RetainOnlyDigits(tp.ReportTime))
	c.FileAttachment(targetZip, dZipName+utils.ZIP_EXT)
	log.Info(tag, zap.Uint64("id", id))
}
