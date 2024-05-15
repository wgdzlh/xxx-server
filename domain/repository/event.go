package repository

import "xxx-server/domain/entity"

type StatusEventBus interface {
	Subscribe(chan entity.StatusEvent) error
	UnSubscribe(chan entity.StatusEvent) error
	Publish(entity.StatusEvent) error
}
