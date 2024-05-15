package repository

import (
	"errors"
	"xxx-server/domain/entity"
)

type SubCmdRepo interface {
	Exec(input entity.AnyJson, args ...string) (out []byte, err error)
}

type PythonCmdRepo interface {
	SubCmdRepo
	Submit(input entity.AnyJson) (out entity.AnyJson, err error)
}

var (
	ErrPyCmdSubmitTimeout = errors.New("py cmd submit timeout")
)
