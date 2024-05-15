package persistence

import (
	"strings"

	log "xxx-server/application/logger"
	"xxx-server/domain/entity"
	repo "xxx-server/domain/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OptionRepoImp struct {
	logTag string
}

const (
	F_ID          = "id"
	F_STATUS      = "status"
	F_NAME        = "name"
	F_TYPE        = "type"
	F_GENRE       = "genre"
	F_ADCODE      = "adcode"
	F_FACTOR      = "factor"
	F_CENTER      = "center"
	F_REGION      = "region"
	F_RES_TIME    = "resource_time"
	F_CREATED_AT  = "created_at"
	F_REPORT_TIME = "report_time"
	F_LABEL       = "label"
	F_LABEL_ALL   = "label_all"
	F_ORDER_BY    = "order_by"
	F_DESC        = "desc"
	F_PAGE        = "page"
	F_PAGESIZE    = "size"

	V_TRUE = "1"

	OPTION_REPO_LOG_TAG = "OptionRepo::"
)

func NewOptionRepo() repo.OptionRepo {
	r := &OptionRepoImp{
		logTag: OPTION_REPO_LOG_TAG,
	}
	return r
}

func (r *OptionRepoImp) getTable() *gorm.DB {
	return R.db.Table(entity.OptionTableName)
}

// 创建
func (r *OptionRepoImp) Create(item *entity.Option) (err error) {
	err = R.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		UpdateAll: true,
	}).Create(item).Error
	if err != nil {
		log.Error(r.logTag+"create option err", zap.String("key", item.Key), zap.Error(err))
		return
	}
	log.Debug(r.logTag+"created", zap.String("key", item.Key))
	return
}

// 删除
func (r *OptionRepoImp) Delete(keys ...string) (err error) {
	err = r.getTable().Where("key IN ?", keys).Delete(nil).Error
	if err != nil {
		log.Error(r.logTag+"delete err", zap.Any("keys", keys), zap.Error(err))
		return
	}
	log.Info(r.logTag+"deleted", zap.Any("keys", keys))
	return
}

// 获取详情
func (r *OptionRepoImp) QueryDetail(key string) (ret entity.Option, err error) {
	if err = r.getTable().Where("key = ?", key).Scan(&ret).Error; err != nil {
		log.Error(r.logTag+"query detail err", zap.Error(err))
	}
	return
}

// 获取值
func (r *OptionRepoImp) GetValue(key string) (value string, err error) {
	ret, err := r.QueryDetail(key)
	value = ret.Value
	return
}

// 获取多值
func (r *OptionRepoImp) GetValues(keys ...string) (values []string, err error) {
	var rets []entity.Option
	if err = r.getTable().Where("key IN ?", keys).Find(&rets).Error; err != nil {
		log.Error(r.logTag+"query details err", zap.Error(err))
		return
	}
	retMap := make(map[string]string, len(rets))
	for _, r := range rets {
		retMap[r.Key] = r.Value
	}
	values = make([]string, len(keys))
	for i, k := range keys {
		values[i] = retMap[k]
	}
	// log.Debug(r.logTag+"got values", zap.Any("rets", rets))
	return
}

// 获取多值
func (r *OptionRepoImp) GetOpts(keys ...string) (rets []entity.Option, err error) {
	keyArr := ToStringArray(keys)
	if err = r.getTable().Where("key = ANY(?)", keyArr).Clauses(clause.OrderBy{Expression: gorm.Expr("array_position(?, key)", keyArr)}).
		Scan(&rets).Error; err != nil {
		log.Error(r.logTag+"query details err", zap.Error(err))
	}
	return
}

// 获取Grid Service详情
func (r *OptionRepoImp) QueryGridServiceDetail(keyPre, sid string) (ret entity.Option, err error) {
	if err = r.getTable().Where("key LIKE ? AND value LIKE ?", keyPre+"%", "%,"+sid).Scan(&ret).Error; err != nil {
		log.Error(r.logTag+"query grid service detail err", zap.Error(err))
	}
	return
}

// 获取列表
func (r *OptionRepoImp) QueryList() (rets []entity.Option, err error) {
	err = r.getTable().Order("key").Find(&rets).Error
	if err != nil {
		log.Error(r.logTag+"query list err", zap.Error(err))
	}
	return
}

// 获取以某前缀开头的Keys
func (r *OptionRepoImp) QueryPrefix(pre string) (rets []string, err error) {
	if err = r.getTable().Select("key").Where("key LIKE ?", pre+"%").Scan(&rets).Error; err != nil {
		log.Error(r.logTag+"query prefix err", zap.Error(err))
	}
	for i, r := range rets {
		rets[i] = strings.TrimPrefix(r, pre)
	}
	return
}

func (r *OptionRepoImp) GetInProcDates() (dates []string, err error) {
	err = r.getTable().Select("key").Where(`key LIKE '20%' AND value = ''`).Scan(&dates).Error
	if err != nil {
		log.Error(r.logTag+"query in proc task date err", zap.Error(err))
	}
	return
}
