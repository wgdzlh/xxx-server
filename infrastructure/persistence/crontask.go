package persistence

import (
	"fmt"
	"strings"

	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"
	"xxx-server/infrastructure/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CronTaskRepoImp struct {
	logTag string
}

const (
	CRON_TASK_REPO_LOG_TAG = "CronTaskRepo::"
)

func NewCronTaskRepo() repo.CronTaskRepo {
	return &CronTaskRepoImp{
		logTag: CRON_TASK_REPO_LOG_TAG,
	}
}

func (r *CronTaskRepoImp) getTable() *gorm.DB {
	return R.db.Table(entity.CronTaskTableName)
}

// 创建
func (r *CronTaskRepoImp) Create(item *entity.CronTask) (err error) {
	err = R.db.Create(item).Error
	if err != nil {
		log.Error(r.logTag+"create err", zap.Error(err))
		return
	}
	log.Info(r.logTag+"created", zap.Uint64("id", item.Id))
	return
}

// 更新
func (r *CronTaskRepoImp) Update(id uint64, fields any) (err error) {
	err = r.getTable().Where("id = ?", id).Updates(fields).Error
	if err != nil {
		log.Error(r.logTag+"update err", zap.Uint64("id", id), zap.Error(err))
		return TransSQLError(err)
	}
	log.Info(r.logTag+"updated", zap.Uint64("id", id), zap.Any("fields", fields))
	return
}

// 删除
func (r *CronTaskRepoImp) Delete(ids ...uint64) (err error) {
	err = r.getTable().Where("id IN ?", ids).Delete(nil).Error
	if err != nil {
		log.Error(r.logTag+"delete err", zap.Error(err))
		return
	}
	log.Info(r.logTag+"deleted", zap.Any("ids", ids))
	return
}

// 详情
func (r *CronTaskRepoImp) QueryDetail(id uint64, selected ...string) (ret *entity.CronTask, err error) {
	cur := r.getTable().Where("id = ?", id)
	if len(selected) > 0 {
		cur = cur.Select(selected)
	}
	if err = cur.Take(&ret).Error; err != nil {
		log.Error(r.logTag+"query detail err", zap.Error(err))
	}
	return
}

// 详情
func (r *CronTaskRepoImp) GetXxxReportTime(id uint64) (ret entity.JSONTime, err error) {
	if err = r.getTable().Where("id = ?", id).Select(`ext->>'report_time'`).Scan(&ret).Error; err != nil {
		log.Error(r.logTag+"get xxx report time err", zap.Error(err))
	}
	return
}

// 获取待重启任务
func (r *CronTaskRepoImp) GetUnended(typ string) (rets []*entity.CronTask, err error) {
	err = r.getTable().Where("type = ? AND end_at IS NULL AND (status = ? OR status = ?)", typ, entity.TASK_IN_PROC, entity.TASK_ST_PROC).Find(&rets).Error
	if err != nil {
		log.Error(r.logTag+"get unended task err", zap.Error(err))
	}
	return
}

// 获取数据准备阶段运行中任务个数
func (r *CronTaskRepoImp) GetInProcCount(typ string) (count int64, err error) {
	err = r.getTable().Where("type = ? AND end_at IS NULL AND status = ?", typ, entity.TASK_IN_PROC).Count(&count).Error
	if err != nil {
		log.Error(r.logTag+"get in proc count err", zap.Error(err))
	}
	return
}

// 获取服务发布阶段任务个数
func (r *CronTaskRepoImp) GetInStageProcCount(typ string) (count int64, err error) {
	err = r.getTable().Where("type = ? AND end_at IS NULL AND status = ?", typ, entity.TASK_ST_PROC).Count(&count).Error
	if err != nil {
		log.Error(r.logTag+"get in stage proc count err", zap.Error(err))
	}
	return
}

// 请求列表
func (r *CronTaskRepoImp) QueryList(filter map[string]string) (rets []*entity.CronTask, total int64, err error) {
	var (
		v  string
		ok bool
	)
	cur := r.getTable()
	if v, ok = filter[F_ID]; ok {
		cur = cur.Where("id IN ?", strings.Split(v, config.Q_SEP))
	}
	if v, ok = filter[F_GENRE]; ok {
		cur = cur.Where("genre = ?", v)
	}
	if v, ok = filter[F_TYPE]; ok {
		cur = cur.Where("type = ?", v)
	}
	if v, ok = filter[F_STATUS]; ok {
		cur = cur.Where("status = ?", v)
	}
	if v, ok = filter[F_ADCODE]; ok {
		codes, _ := R.OutSvr.GetDistrictAllCodes(v)
		if len(codes) > 0 {
			codesArr := Int32ToIntArray(codes)
			cur = cur.Where("(params->>'adcode')::int = ANY(?)", codesArr)
		}
	}
	if v, ok = filter[F_FACTOR]; ok {
		cur = cur.Where("params->'themes' @> ?", fmt.Sprintf(`[{"factor":"%s"}]`, v))
	}
	if v, ok = filter[F_REPORT_TIME]; ok {
		cur = cur.Where("ext->>'report_time' <= ?", v)
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
	if v, ok = filter[F_ORDER_BY]; ok {
		if filter[F_DESC] == V_TRUE {
			cur = cur.Order(v + " DESC")
		} else {
			cur = cur.Order(v)
		}
	}

	err = cur.Limit(size).Order("id DESC").Find(&rets).Error
	if err != nil {
		log.Error(r.logTag+"query list err", zap.Error(err))
	}
	return
}
