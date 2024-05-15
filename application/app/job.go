package app

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"xxx-server/application/ext"
	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"
	"xxx-server/infrastructure/config"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
)

func (s *Scheduler) ProcWorkflowMqRet(_ string, keys []string, body []byte) (err error) {
	var ret entity.WorkflowMqResp
	e := json.Unmarshal(body, &ret)
	if e != nil {
		log.Error(s.logTag+"workflow ret Unmarshal failed", zap.Error(e))
		return
	}
	if ret.Code == 7 {
		log.Info(s.logTag+"workflow paused", zap.Int("code", ret.Code), zap.Any("keys", keys))
		return
		// case 5:
		// 	log.Info(s.logTag+"workflow started", zap.Any("keys", keys))
		// 	return
	}
	if ret.CustomMsgId == "" || len(keys) == 0 {
		log.Error(s.logTag+"workflow invalid", zap.Any("keys", keys))
		return
	}
	var job *entity.CronTask
	if strings.Contains(keys[0], ext.XxxTaskType) {
		if job, err = s.getJob(ret); job != nil && err == nil {
			go s.procXxxJob(job) // 耗时任务
		}
	} else if strings.Contains(keys[0], ext.ThemeTaskType) {
		if job, err = s.getJob(ret); job != nil && err == nil {
			err = s.procThemeJob(job)
		}
	} else {
		log.Error(s.logTag+"unknown workflow keys", zap.Any("keys", keys))
	}
	return
}

func (s *Scheduler) getJob(ret entity.WorkflowMqResp) (job *entity.CronTask, err error) {
	taskId := utils.StrToUint64(ret.CustomMsgId)
	if ret.Code == 5 || ret.Code == 4 {
		var trs []entity.TaskResult
		if e := json.Unmarshal(ret.Data, &trs); e != nil || len(trs) == 0 {
			log.Error(s.logTag+"unmarshal job task results failed", zap.Error(e), zap.Any("res", ret.Data))
			return
		}
		if ret.Code == 5 {
			err = Rs.CronTask.Update(taskId, struct {
				ExtId int64
			}{trs[0].Id})
		} else {
			if trs[0].Progress == 0 {
				err = Rs.CronTask.Update(taskId, struct {
					StartAt time.Time
					Status  string
				}{time.Now(), entity.TASK_IN_PROC})
			} else {
				err = Rs.CronTask.Update(taskId, struct {
					Progress float32
				}{trs[0].Progress})
			}
		}
		return
	}

	if ret.Code != 0 {
		log.Error(s.logTag+"xxx job failed", zap.Uint64("id", taskId))
		err = Rs.CronTask.Update(taskId, struct {
			EndAt  time.Time
			Status string
			ErrLog string
		}{time.Now(), entity.TASK_FAILED, utils.PurifyForUtf8(ret.Msg)})
		return
	}
	job, err = Rs.CronTask.QueryDetail(taskId)
	return
}

func (s *Scheduler) finishJob(job *entity.CronTask, e error) error {
	status := entity.TASK_DONE
	errLog := ""
	if e != nil {
		status = entity.TASK_FAILED
		errLog = utils.PurifyForUtf8(e.Error())
	}
	return Rs.CronTask.Update(job.Id, entity.MapJson{
		"end_at":   time.Now(),
		"progress": 1.0,
		"status":   status,
		"ext":      job.Ext,
		"err_log":  errLog,
	})
}

func (s *Scheduler) procXxxJob(job *entity.CronTask) {
	var (
		err error
		mte entity.XxxTaskExt
	)
	switch job.Type {
	case entity.TASK_TYPE_SHORT:
	case entity.TASK_TYPE_MID:
	}
	job.Ext, _ = json.Marshal(mte)
	s.finishJob(job, err)
}

func (s *Scheduler) procThemeJob(job *entity.CronTask) error {
	retDir := job.Dir
	pDir := filepath.Dir(retDir)
	e := os.Chdir(pDir)
	if e == nil {
		retZip := retDir + utils.ZIP_EXT
		retDir = filepath.Base(retDir)
		if e = utils.CompressZip(retZip, retDir); e != nil {
			log.Error(s.logTag+"CompressZip failed", zap.Error(e))
		} else {
			os.RemoveAll(retDir)
		}
	}
	return s.finishJob(job, e)
}

func (s *Scheduler) RestartJob(ct *entity.CronTask) (err error) {
	if ct.Progress == 1.0 { // already finished all ops, no need to restart workflow
		if err = Rs.CronTask.Update(ct.Id, entity.MapJson{
			"end_at":  nil,
			"err_log": "",
			"status":  entity.TASK_IN_PROC,
		}); err != nil {
			return
		}
		go s.procXxxJob(ct)
		return
	}
	// taskId := strconv.FormatUint(ct.Id, 10)
	// ts := utils.GetNowTimeTag()
	switch ct.Type {
	case entity.TASK_TYPE_SHORT:
	case entity.TASK_TYPE_MID:
	}
	if err != nil {
		log.Error(s.logTag+"trigger task failed", zap.Error(err))
		return
	}
	err = Rs.CronTask.Update(ct.Id, entity.MapJson{
		"start_at": nil,
		"end_at":   nil,
		"err_log":  "",
		"dir":      ct.Dir,
		"ext_id":   0,
		"progress": 0,
		"status":   entity.TASK_INIT,
	})
	return
}

func (s *Scheduler) DownloadJob(ct *entity.CronTask, req *entity.FactorThemeReq) (tid uint64, err error) {
	var (
		task = entity.CronTask{
			Name:  "XXX专题图生成",
			Genre: entity.TASK_GENRE_DOWN,
			Type:  entity.TASK_TYPE_FACTOR,
			SrcId: ct.Id,
		}
		fcTypes, titles = transFactors(ct.Type == entity.TASK_TYPE_MID, req)
	)
	if len(fcTypes) == 0 {
		err = repo.ErrInvalidFactors
		return
	}
	reportTime, err := ct.GetReportTime()
	if err != nil {
		return
	}
	req.ReportTime = reportTime.Format(entity.THEME_TF)
	task.Dir, err = utils.GetDateSubDir(config.C.Cron.ThemeDirRoot, "factor_"+uuid.NewString())
	if err != nil {
		log.Error(s.logTag+"failed to make sub dir", zap.Error(err))
		return
	}
	task.Params, _ = json.Marshal(req)
	err = Rs.CronTask.Create(&task)
	if err != nil {
		log.Error(s.logTag+"failed to create task", zap.Error(err))
		os.Remove(task.Dir)
		return
	}
	tid = task.Id
	taskId := strconv.FormatUint(task.Id, 10)
	input := &entity.XxxThemeInput{
		ForecastTime: req.ReportTime,
		TiffDir:      ct.Dir,
		Adcode:       req.Adcode,
		FcType:       fcTypes,
		Title:        titles,
		FigureDir:    task.Dir,
	}
	err = WorkflowRepo.XxxTheme(taskId, input)
	if err != nil {
		log.Error(s.logTag+"failed to trigger task", zap.Error(err))
		Rs.CronTask.Update(task.Id, struct {
			EndAt  time.Time
			Status string
			ErrLog string
		}{time.Now(), entity.TASK_FAILED, utils.PurifyForUtf8(err.Error())})
	}
	return
}

func transFactors(isMid bool, req *entity.FactorThemeReq) (factors, titles []string) {
	var inputFac string
	for _, t := range req.Themes {
		switch t.Factor {
		case "降雨":
			if isMid {
				inputFac = "Precip6h"
			} else {
				inputFac = "Short_Precip"
			}
		case "气压":
			inputFac = "Pressure"
		case "温度":
			inputFac = "T2"
		case "湿度":
			inputFac = "Rh2"
		case "风场":
			inputFac = "WS10"
		default:
			continue
		}
		factors = append(factors, inputFac)
		titles = append(titles, t.Title)
	}
	return
}

// func transDistLevel(level string) string {
// 	if level == "district" {
// 		return "county"
// 	}
// 	return level
// }
