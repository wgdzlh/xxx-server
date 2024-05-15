package repository

import "xxx-server/domain/entity"

type XxxDataRepo interface {
	Create(item *entity.XxxData) (err error)
	Delete(ids ...uint64) (err error)
	QueryDetail(id uint64) (ret entity.XxxData, err error)
	QueryList(page, size int) (rets []entity.XxxData, total int64, err error)
}
