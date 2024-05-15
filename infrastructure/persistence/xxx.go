package persistence

import (
	log "xxx-server/application/logger"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type XxxDataRepoImp struct {
	logTag string
}

const (
	XXX_DATA_REPO_LOG_TAG = "XxxDataRepo::"
)

func NewXxxDataRepo() repo.XxxDataRepo {
	r := &XxxDataRepoImp{
		logTag: XXX_DATA_REPO_LOG_TAG,
	}
	return r
}

func (r *XxxDataRepoImp) getTable() *gorm.DB {
	return R.db.Table(entity.XxxDataTableName)
}

// 创建
func (r *XxxDataRepoImp) Create(item *entity.XxxData) (err error) {
	err = R.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "district"}},
		UpdateAll: true,
	}).Create(item).Error
	if err != nil {
		log.Error(r.logTag+"create xxx data err", zap.Error(err))
	}
	return
}

// 删除
func (r *XxxDataRepoImp) Delete(ids ...uint64) (err error) {
	err = r.getTable().Where("id IN ?", ids).Delete(nil).Error
	if err != nil {
		log.Error(r.logTag+"delete err", zap.Error(err))
		return
	}
	log.Info(r.logTag+"deleted", zap.Any("ids", ids))
	return
}

// 获取详情
func (r *XxxDataRepoImp) QueryDetail(id uint64) (ret entity.XxxData, err error) {
	cur := r.getTable().Where("id = ?", id)
	if err = cur.Take(&ret).Error; err != nil {
		log.Error(r.logTag+"query detail err", zap.Error(err))
	}
	return
}

// 获取列表
func (r *XxxDataRepoImp) QueryList(page, size int) (rets []entity.XxxData, total int64, err error) {
	cur := r.getTable()
	if err = cur.Count(&total).Error; err != nil {
		log.Error(r.logTag+"query total err", zap.Error(err))
		return
	}
	if size <= 0 {
		size = 10
	}
	if page > 1 {
		cur = cur.Offset(size * (page - 1))
	}
	err = cur.Order("id DESC").Limit(size).Find(&rets).Error
	if err != nil {
		log.Error(r.logTag+"query list err", zap.Error(err))
	}
	return
}
