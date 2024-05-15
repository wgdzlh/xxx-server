package repository

import (
	"errors"
	"xxx-server/domain/entity"
)

type SettingRepo interface {
	Create(item *entity.Setting) (err error)
	Delete(ids ...uint64) (err error)
	QueryDetail(id uint64) (ret entity.Setting, err error)
	Query(section string) (ret []entity.Setting, err error)
	QueryAdcodeDetail(adcode string) (ret entity.Setting, err error)
	CheckAlert(ids ...uint64) (hasAlert bool, err error)
}

var (
	ErrInvalidSetting = errors.New("invalid setting")
)
