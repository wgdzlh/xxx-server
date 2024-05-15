package persistence

import (
	"fmt"
	"strings"
	"time"

	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"
	"xxx-server/infrastructure/config"

	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BackgroundTaskRepoImp struct {
	logTag string
}

const (
	TIF_TASK_REPO_LOG_TAG = "HighResTifTaskRepo::"
)

func NewBackgroundTaskRepo() repo.BackgroundTaskRepo {
	return &BackgroundTaskRepoImp{
		logTag: TIF_TASK_REPO_LOG_TAG,
	}
}

func (r *BackgroundTaskRepoImp) getTable() *gorm.DB {
	return R.db.Table(entity.BackgroundTaskTableName)
}

// 创建
func (r *BackgroundTaskRepoImp) Create(item *entity.BackgroundTask) (err error) {
	err = R.db.Create(item).Error
	if err != nil {
		log.Error(r.logTag+"create err", zap.Error(err))
		return
	}
	log.Info(r.logTag+"created", zap.Uint64("id", item.Id))
	return
}

// 更新
func (r *BackgroundTaskRepoImp) Update(id uint64, fields any) (err error) {
	err = r.getTable().Where("id = ?", id).Updates(fields).Error
	if err != nil {
		log.Error(r.logTag+"update err", zap.Uint64("id", id), zap.Error(err))
		return TransSQLError(err)
	}
	log.Info(r.logTag+"updated", zap.Uint64("id", id), zap.Any("fields", fields))
	return
}

// 更新流程
func (r *BackgroundTaskRepoImp) UpdateProcedure(id uint64, p entity.TaskStep, done ...bool) (err error) {
	if p.Idx < 0 {
		log.Error(r.logTag + "step index invalid")
		return
	}
	pjs, err := json.Marshal(p)
	if err != nil {
		return
	}
	// if p.TaskId != 0 {
	newValue := gorm.Expr(fmt.Sprintf(`jsonb_set(procedure, '{%d}', procedure->%d || '%s', true)`, p.Idx, p.Idx, pjs))
	// } else if p.TifNum > 0 {
	// 	newValue = gorm.Expr(fmt.Sprintf(`jsonb_set(procedure, '{%d}', procedure->%d || '{"status": "%s", "ratio": %f}', true)`, p.Idx, p.Idx, p.Status, p.Ratio))
	// } else {
	// 	newValue = gorm.Expr(fmt.Sprintf(`jsonb_set(procedure, '{%d, status}', '"%s"', true)`, p.Idx, p.Status))
	// }
	err = r.getTable().Where("id = ?", id).Update("procedure", newValue).Error
	if err != nil {
		log.Error(r.logTag+"update procedure err", zap.Uint64("id", id), zap.Error(err))
		return
	}
	if p.Status == entity.TASK_FAILED {
		r.Update(id, struct {
			EndAt  time.Time
			Status string
		}{time.Now(), entity.TASK_FAILED})
	} else if len(done) > 0 {
		if done[0] {
			r.Update(id, struct {
				EndAt  time.Time
				Status string
			}{time.Now(), entity.TASK_DONE})
		} else {
			r.Update(id, struct {
				StageEndAt time.Time
				Status     string
			}{time.Now(), entity.TASK_STAGED})
		}
	}
	log.Info(r.logTag+"procedure updated", zap.Uint64("id", id), zap.Int("pIdx", p.Idx))
	return
}

// 删除
func (r *BackgroundTaskRepoImp) Delete(ids ...uint64) (err error) {
	err = r.getTable().Where("id IN ?", ids).Delete(nil).Error
	if err != nil {
		log.Error(r.logTag+"delete err", zap.Error(err))
		return
	}
	log.Info(r.logTag+"deleted", zap.Any("ids", ids))
	return
}

// 详情
func (r *BackgroundTaskRepoImp) QueryDetail(id uint64, selected ...string) (ret *entity.BackgroundTask, err error) {
	cur := r.getTable().Where("id = ?", id)
	if len(selected) > 0 {
		cur = cur.Select(selected)
	}
	if err = cur.Take(&ret).Error; err != nil {
		log.Error(r.logTag+"query detail err", zap.Error(err))
	}
	return
}

// 获取待重启任务
func (r *BackgroundTaskRepoImp) GetUnended(typ string) (rets []*entity.BackgroundTask, err error) {
	err = r.getTable().Where("type = ? AND end_at IS NULL AND (status = ? OR status = ?)", typ, entity.TASK_IN_PROC, entity.TASK_ST_PROC).Find(&rets).Error
	if err != nil {
		log.Error(r.logTag+"get unended task err", zap.Error(err))
	}
	return
}

// 获取阶段完结任务
func (r *BackgroundTaskRepoImp) GetStaged(typ string) (rets []*entity.BackgroundTask, err error) {
	err = r.getTable().Where("type = ? AND end_at IS NULL AND status = ?", typ, entity.TASK_STAGED).Order("id").Find(&rets).Error
	if err != nil {
		log.Error(r.logTag+"get unended staged task err", zap.Error(err))
	}
	return
}

// 获取数据准备阶段运行中任务个数
func (r *BackgroundTaskRepoImp) GetInProcCount(typ string) (count int64, err error) {
	err = r.getTable().Where("type = ? AND end_at IS NULL AND status = ?", typ, entity.TASK_IN_PROC).Count(&count).Error
	if err != nil {
		log.Error(r.logTag+"get in proc count err", zap.Error(err))
	}
	return
}

// 获取服务发布阶段任务个数
func (r *BackgroundTaskRepoImp) GetInStageProcCount(typ string) (count int64, err error) {
	err = r.getTable().Where("type = ? AND end_at IS NULL AND status = ?", typ, entity.TASK_ST_PROC).Count(&count).Error
	if err != nil {
		log.Error(r.logTag+"get in stage proc count err", zap.Error(err))
	}
	return
}

// 请求列表
func (r *BackgroundTaskRepoImp) QueryList(filter map[string]string) (rets []*entity.BackgroundTask, total int64, err error) {
	var (
		v  string
		ok bool
	)
	cur := r.getTable()
	if v, ok = filter[F_ID]; ok {
		cur = cur.Where("id IN ?", strings.Split(v, config.Q_SEP))
	}
	if v, ok = filter[F_TYPE]; ok {
		cur = cur.Where("type = ?", v)
	}
	if v, ok = filter[F_CREATED_AT]; ok {
		pair := strings.SplitN(v, config.Q_SEP, 2)
		if len(pair) == 0 {
			goto next
		}
		startTime, e := utils.StrToTime(pair[0])
		validStart := e == nil
		if len(pair) > 1 {
			endTime, e := utils.StrToTime(pair[1])
			if e == nil {
				if validStart {
					cur = cur.Where("created_at BETWEEN ? AND ?", startTime, endTime)
				} else {
					cur = cur.Where("created_at <= ?", endTime)
				}
				goto next
			}
		}
		if validStart {
			cur = cur.Where("created_at >= ?", startTime)
		} else {
			log.Error(r.logTag+"parse time err", zap.Error(e))
		}
	}

next:
	if err = cur.Count(&total).Error; err != nil {
		log.Error(r.logTag+"query count err", zap.Error(err))
		return
	}

	size := 10
	if v, ok = filter[F_PAGESIZE]; ok {
		if newSize := utils.StrToInt(v); newSize > 0 {
			size = newSize
		}
	}
	if v, ok = filter[F_PAGE]; ok {
		if page := utils.StrToInt(v); page > 1 {
			cur = cur.Offset(size * (page - 1))
		}
	}

	err = cur.Limit(size).Omit("process_task").Order("id DESC").Find(&rets).Error
	if err != nil {
		log.Error(r.logTag+"query list err", zap.Error(err))
	}
	return
}
