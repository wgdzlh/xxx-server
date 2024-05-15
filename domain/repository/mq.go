package repository

import (
	"xxx-server/domain/entity"

	"github.com/wgdzlh/mqlib"
)

type MqApi interface {
	SendMsg(msg *mqlib.Message) error
	SendWorkflowReq(req *entity.WorkflowReq) (err error)
}
