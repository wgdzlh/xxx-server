package task

import (
	"os"

	log "xxx-server/application/logger"
	repo "xxx-server/domain/repository"
	"xxx-server/infrastructure/config"

	"go.uber.org/zap"
)

const (
	HISTORY_RESET_TASK_EXP = "0 0 1 * *"

	HISTORY_RESET_TASK = "HistoryResetTask"
)

var (
	xxxTifDirRoot string
	workflowRepo  repo.WorkflowRepository
)

func SetupBackgroundTasks(wfr repo.WorkflowRepository) (cron *CronManager) {
	xxxTifDirRoot = config.C.Cron.XxxTifDirRoot
	if err := os.MkdirAll(xxxTifDirRoot, os.ModePerm); err != nil {
		log.Fatal("failed to make xxxTifDirRoot in cron", zap.Error(err))
	}
	workflowRepo = wfr
	cron = NewCronManger()

	// 定时历史XXX数据重置任务
	hrt := NewHistoryResetTask()
	cron.StartJob(hrt.Name, HISTORY_RESET_TASK_EXP, hrt.Start)
	return
}
