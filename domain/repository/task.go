package repository

import (
	"errors"

	"xxx-server/domain/entity"
)

type BgTask interface {
	DoTaskStep(taskId uint64, restart ...bool) (err error)
	AwaitStep(taskId uint64, stepName string, isFailed bool, ids ...string)
}

type BackgroundTaskRepo interface {
	Create(*entity.BackgroundTask) error
	Update(id uint64, fields any) error
	UpdateProcedure(id uint64, step entity.TaskStep, done ...bool) error
	Delete(ids ...uint64) error
	QueryDetail(id uint64, selected ...string) (*entity.BackgroundTask, error)
	QueryList(filter map[string]string) ([]*entity.BackgroundTask, int64, error)
	GetUnended(typ string) ([]*entity.BackgroundTask, error)
	GetStaged(typ string) ([]*entity.BackgroundTask, error)
	GetInProcCount(typ string) (count int64, err error)
	GetInStageProcCount(typ string) (count int64, err error)
}

type CronTaskRepo interface {
	Create(*entity.CronTask) error
	Update(id uint64, fields any) error
	Delete(ids ...uint64) error
	QueryDetail(id uint64, selected ...string) (*entity.CronTask, error)
	QueryList(filter map[string]string) ([]*entity.CronTask, int64, error)
	GetUnended(typ string) ([]*entity.CronTask, error)
	GetInProcCount(typ string) (count int64, err error)
	GetInStageProcCount(typ string) (count int64, err error)
	GetXxxReportTime(id uint64) (ret entity.JSONTime, err error)
}

var (
	ErrTaskNotReady         = errors.New("task not ready")
	ErrZipNotExist          = errors.New("zip not exist")
	ErrInvalidFactors       = errors.New("invalid factors")
	ErrInvalidDistrictLevel = errors.New("invalid district level")
)
