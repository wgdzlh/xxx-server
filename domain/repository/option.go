package repository

import "xxx-server/domain/entity"

type OptionRepo interface {
	Create(item *entity.Option) (err error)
	Delete(keys ...string) (err error)
	QueryDetail(key string) (ret entity.Option, err error)
	GetValue(key string) (value string, err error)
	GetValues(keys ...string) (values []string, err error)
	GetOpts(keys ...string) (values []entity.Option, err error)
	QueryGridServiceDetail(keyPre, sid string) (ret entity.Option, err error)
	QueryList() (rets []entity.Option, err error)
	QueryPrefix(pre string) (rets []string, err error)
	GetInProcDates() (dates []string, err error)
}
