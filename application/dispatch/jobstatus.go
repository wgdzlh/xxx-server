package dispatch

import (
	"errors"
	"time"
	log "xxx-server/application/logger"
	"xxx-server/domain/entity"
	"xxx-server/domain/repository"

	"go.uber.org/zap"
)

type StatusChan = chan entity.StatusEvent

type JobStatusEventBus struct {
	// Events are pushed to this channel
	events StatusChan
	// New client connections
	newClients chan StatusChan
	// Closed client connections
	closedClients chan StatusChan
	// Total client connections
	totalClients map[StatusChan]struct{}
	logTag       string
}

const (
	MSG_BUFF_SIZE = 10
	SUB_TIMEOUT   = time.Second * 5
	PUB_TIMEOUT   = time.Second * 5

	SSE_WAIT_PERIOD = time.Second * 10
)

var (
	emptyStatusEvent entity.StatusEvent
)

func NewJobStatusEventBus() repository.StatusEventBus {
	bus := &JobStatusEventBus{
		events:        make(StatusChan, MSG_BUFF_SIZE),
		newClients:    make(chan StatusChan, MSG_BUFF_SIZE),
		closedClients: make(chan StatusChan, MSG_BUFF_SIZE),
		totalClients:  map[StatusChan]struct{}{},
		logTag:        "JobStatusEventBus:",
	}
	bus.loop()
	return bus
}

// 发布消息
func (b *JobStatusEventBus) Publish(event entity.StatusEvent) (err error) {
	select {
	case b.events <- event:
	case <-time.After(PUB_TIMEOUT):
		eMsg := b.logTag + "pub que timeout"
		log.Error(eMsg)
		err = errors.New(eMsg)
	}
	return
}

// 订阅
func (b *JobStatusEventBus) Subscribe(clientChan chan entity.StatusEvent) (err error) {
	select {
	case b.newClients <- clientChan:
	case <-time.After(SUB_TIMEOUT):
		eMsg := b.logTag + "new sub que timeout"
		log.Error(eMsg)
		err = errors.New(eMsg)
	}
	return
}

// 取消订阅
func (b *JobStatusEventBus) UnSubscribe(clientChan chan entity.StatusEvent) (err error) {
	select {
	case b.closedClients <- clientChan:
	case <-time.After(SUB_TIMEOUT):
		eMsg := b.logTag + "close sub que timeout"
		log.Error(eMsg)
		err = errors.New(eMsg)
	}
	return
}

func (b *JobStatusEventBus) loop() {
	go func() {
		var (
			client StatusChan
			event  entity.StatusEvent
			ticker = time.NewTicker(SSE_WAIT_PERIOD)
		)
		for {
			select {
			// Add new available client
			case client = <-b.newClients:
				b.totalClients[client] = struct{}{}
				log.Info(b.logTag+"client added", zap.Int("total", len(b.totalClients)))
			// Remove closed client
			case client = <-b.closedClients:
				delete(b.totalClients, client)
				close(client)
				log.Info(b.logTag+"client removed", zap.Int("total", len(b.totalClients)))
			// Broadcast message to clients
			case event = <-b.events:
				b.fanOut(event)
			case <-ticker.C:
				b.fanOut(emptyStatusEvent)
			}
		}
	}()
}

func (b *JobStatusEventBus) fanOut(event entity.StatusEvent) {
	for client := range b.totalClients {
		select {
		case client <- event:
		default:
		}
	}
}
