package ext

import (
	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"
	"xxx-server/infrastructure/config"

	"go.uber.org/zap"
)

const (
	XxxTaskType          = "xxx"
	MidForecastTaskType  = XxxTaskType + "-mid"
	ShortPrecastTaskType = XxxTaskType + "-short"
	ThemeTaskType        = "theme"
	AlertTaskType        = "alert"
	MidAlertTaskType     = AlertTaskType + "-mid"
	ShortAlertTaskType   = AlertTaskType + "-short"

	SepInId = "_"

	WORKFLOW_DELEGATE_LOG_TAG = "WorkflowDelegate:"
)

type WorkflowDelegate struct {
	mq     repo.MqApi
	logTag string
}

func NewWorkflowRepository(mq repo.MqApi) repo.WorkflowRepository {
	d := &WorkflowDelegate{
		mq:     mq,
		logTag: WORKFLOW_DELEGATE_LOG_TAG,
	}
	return d
}

// XXX专题图
func (d *WorkflowDelegate) XxxTheme(taskId string, input *entity.XxxThemeInput) (err error) {
	cfg := config.C.Workflow.XxxTheme
	return d.sendReq(taskId, ThemeTaskType, taskId, cfg, input)
}

func (d *WorkflowDelegate) sendReq(taskId, taskType, tName string, cfg config.WorkflowReq, inputs ...any) (err error) {
	hasOps := len(cfg.OpId) > 0
	if hasOps && len(cfg.OpId) != len(inputs) {
		log.Error(d.logTag + "miss matched config and input count")
		return
	}
	var (
		wfi entity.WorkflowInputs
		wfp entity.WorkflowParams
	)
	if hasOps {
		params := make(map[string]any, len(inputs))
		for i, input := range inputs {
			params[cfg.OpId[i]] = input
		}
		wfi = entity.WorkflowInputs{params}
	} else {
		wfp = inputs
	}
	req := entity.WorkflowReq{
		TaskId:            taskId,
		TaskType:          taskType,
		WorkflowServiceId: cfg.ServiceId,
		Name:              cfg.Name,
		QueueId:           cfg.QueueId,
		GroupId:           cfg.GroupId,
		Priority:          cfg.Priority,
		Input:             wfi,
		Param:             wfp,
	}
	if tName != "" {
		req.Name += "_" + tName
	}
	req.Name += "_" + utils.GetNowTimeTag()
	err = d.mq.SendWorkflowReq(&req)
	log.Info(d.logTag+"workflow via mq", zap.Bool("succeed", err == nil))
	return
}
