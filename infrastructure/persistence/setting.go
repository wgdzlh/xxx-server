package persistence

import (
	log "xxx-server/application/logger"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SettingRepoImp struct {
	logTag string
}

const (
	SETTING_REPO_LOG_TAG = "SettingRepo::"
)

func NewSettingRepo() repo.SettingRepo {
	r := &SettingRepoImp{
		logTag: SETTING_REPO_LOG_TAG,
	}
	return r
}

func (r *SettingRepoImp) getTable() *gorm.DB {
	return R.db.Table(entity.SettingTableName)
}

// 创建
func (r *SettingRepoImp) Create(item *entity.Setting) (err error) {
	if item == nil || item.Section == "" || len(item.Value) == 0 {
		err = repo.ErrInvalidSetting
		return
	}
	err = R.db.FirstOrCreate(item, *item).Error
	if err != nil {
		log.Error(r.logTag+"create Setting err", zap.Error(err))
		return
	}
	log.Info(r.logTag+"created", zap.Any("item", item))
	return
}

// 删除
func (r *SettingRepoImp) Delete(ids ...uint64) (err error) {
	err = r.getTable().Where("id IN ?", ids).Delete(nil).Error
	if err != nil {
		log.Error(r.logTag+"delete err", zap.Error(err))
		return
	}
	log.Info(r.logTag+"deleted", zap.Any("ids", ids))
	return
}

// 检查是否有alert
func (r *SettingRepoImp) CheckAlert(ids ...uint64) (hasAlert bool, err error) {
	err = R.db.Raw(`SELECT true FROM "`+entity.SettingTableName+`" WHERE id IN ? AND section = 'alert' LIMIT 1`,
		ids).Scan(&hasAlert).Error
	if err != nil {
		log.Error(r.logTag+"check alert err", zap.Error(err))
	}
	return
}

// 获取详情
func (r *SettingRepoImp) QueryDetail(id uint64) (ret entity.Setting, err error) {
	if err = r.getTable().Where("id = ?", id).Take(&ret).Error; err != nil {
		log.Error(r.logTag+"query detail err", zap.Error(err))
	}
	return
}

// 获取地区预警详情
func (r *SettingRepoImp) QueryAdcodeDetail(adcode string) (ret entity.Setting, err error) {
	if err = r.getTable().Where(`section = 'alert' AND value->>'lat_lon' IS NULL AND value->>'adcode' = ?`, adcode).Limit(1).Scan(&ret).Error; err != nil {
		log.Error(r.logTag+"query adcode detail err", zap.Error(err))
	}
	return
}

// 获取列表
func (r *SettingRepoImp) Query(section string) (rets []entity.Setting, err error) {
	cur := r.getTable()
	if section != "" {
		cur = cur.Where("section = ?", section)
	}
	err = cur.Order("id DESC").Find(&rets).Error
	if err != nil {
		log.Error(r.logTag+"query err", zap.Error(err))
	}
	return
}
