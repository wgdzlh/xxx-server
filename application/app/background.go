package app

import (
	"log"
	"os"
	"time"

	"xxx-server/infrastructure/config"

	"go.uber.org/zap"
)

const (
	SERVICE_CHECK_INTERVAL = 30 // in seconds

	SCHEDULE_LOG_TAG = "Schedule::"

	FINISHED_SHORT_SID_OPT = "xxx_short_fin_sid"
	FINISHED_MID_SID_OPT   = "xxx_mid_fin_sid"
	WIND_DIR_PNG_OPT       = "xxx_wind_png"
)

type Scheduler struct {
	tmpDir string
	logTag string
}

func setupScheduler() *Scheduler {
	s := &Scheduler{
		tmpDir: config.C.Server.TmpDir,
		logTag: SCHEDULE_LOG_TAG,
	}
	if err := os.MkdirAll(s.tmpDir, os.ModePerm); err != nil {
		log.Fatal("failed to make server tmp dir", zap.Error(err))
	}
	if !config.C.Cron.DisableRecover {
		s.recover()
		time.Sleep(time.Second) // make sure recover run first before mq
	}
	if !config.C.Cron.DisableLoops {
		s.runLoops()
	}
	config.AddCallback(s.updateCronTaskExps)
	return s
}

func (s *Scheduler) recover() {
}

func (s *Scheduler) runLoops() {
	intervalSecs := config.C.Cron.ServiceCheckInterval
	if intervalSecs <= 0 {
		intervalSecs = SERVICE_CHECK_INTERVAL
	}
	ticker := time.NewTicker(time.Second * time.Duration(intervalSecs))
	go func() {
		for range ticker.C {
		}
	}()
}

// 更新后台定时任务调度周期
func (s *Scheduler) updateCronTaskExps(sc *config.SelfConfig) {
	// Cron.ResetJob(task.SHORT_XXX_TASK, sc.Cron.ShortXxxExp)
	// Cron.ResetJob(task.MID_XXX_TASK, sc.Cron.MidXxxExp)
}
