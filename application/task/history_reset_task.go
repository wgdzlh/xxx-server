package task

import (
	log "xxx-server/application/logger"
)

type HistoryResetTask struct {
	Name string
}

func NewHistoryResetTask() *HistoryResetTask {
	return &HistoryResetTask{
		Name: HISTORY_RESET_TASK,
	}
}

func (t *HistoryResetTask) Start() {
	log.Info(t.Name + " begin")
	// if err := persistence.R.HistoryXxx.ResetCurMonth(); err != nil {
	// 	log.Error(t.Name+" failed to reset cur month history", zap.Error(err))
	// }
	log.Info(t.Name + " done")
}
