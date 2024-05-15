package app

import (
	"xxx-server/application/ext"
	"xxx-server/application/mq"
	"xxx-server/application/task"
	repo "xxx-server/domain/repository"
	"xxx-server/infrastructure/disk"
	"xxx-server/infrastructure/persistence"
	"xxx-server/infrastructure/shell"

	"github.com/wgdzlh/gdalib"
)

var (
	Rs *persistence.Repositories

	DataStoreRepo   repo.DataStoreRepository
	TileServiceRepo repo.TileServiceRepository

	SchedulerSvr *Scheduler
	GdalRepo     *gdalib.GdalToolbox

	MqRepo       repo.MqApi
	WorkflowRepo repo.WorkflowRepository

	Cron *task.CronManager

	SeaweedRepo *disk.SeaweedRepo
)

func Init(rs *persistence.Repositories, cs *shell.Repositories) {
	Rs = rs

	DataStoreRepo = ext.NewDataStoreRepository()
	TileServiceRepo = ext.NewTileServiceRepository()
	GdalRepo = gdalib.NewGdalToolbox()
	SeaweedRepo = disk.NewSwClient()

	SchedulerSvr = setupScheduler()

	MqRepo = mq.NewMqService(mq.SubHandler{
		Topic: mq.TopicWorkflowResp,
		Tags:  []string{mq.TagWorkflowResp},
		F:     SchedulerSvr.ProcWorkflowMqRet,
	})
	WorkflowRepo = ext.NewWorkflowRepository(MqRepo)

	Cron = task.SetupBackgroundTasks(WorkflowRepo) // 启动定时任务
}
