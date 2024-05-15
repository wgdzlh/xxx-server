package task

import (
	log "xxx-server/application/logger"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type CronManager struct {
	core      *cron.Cron
	nameToJob map[string]*Job
}

type Job struct {
	Id  cron.EntryID
	Exp string
	Fn  func()
}

func NewCronManger() *CronManager {
	cm := &CronManager{
		core:      cron.New(),
		nameToJob: map[string]*Job{},
	}
	cm.core.Start()
	// config.AddCallback(cm.syncExp)
	return cm
}

func (m *CronManager) StartJob(name, exp string, fn func()) (err error) {
	if _, ok := m.nameToJob[name]; ok {
		log.Warn("cron job already started", zap.String("name", name))
		return
	}
	var id cron.EntryID = -1
	if exp != "" {
		id, err = m.core.AddFunc(exp, fn)
		if err != nil {
			log.Fatal("start job failed", zap.String("name", name), zap.Error(err))
			return
		}
	}
	m.nameToJob[name] = &Job{id, exp, fn}
	log.Info("new cron job started", zap.String("name", name), zap.String("exp", exp))
	return
}

func (m *CronManager) ResetJob(name, exp string) (err error) {
	tJob, ok := m.nameToJob[name]
	if !ok {
		// log.Info("cron job not found", zap.String("name", name))
		return
	}
	if tJob.Exp == exp {
		return
	}
	var id cron.EntryID = -1
	if exp != "" {
		id, err = m.core.AddFunc(exp, tJob.Fn)
		if err != nil {
			log.Error("reset job failed", zap.Error(err))
			return
		}
	}
	if tJob.Id != -1 {
		m.core.Remove(tJob.Id)
	}
	tJob.Id = id
	tJob.Exp = exp
	log.Info("reset job succeed", zap.String("name", name), zap.String("exp", exp))
	return
}

// func (m *CronManager) syncExp(sc *config.SelfConfig) {
// if sc.Cron.TrimLogExp != config.C.Cron.TrimLogExp {
// 	m.ResetJob(TRIM_LOG_TASK, sc.Cron.TrimLogExp)
// }
// }
