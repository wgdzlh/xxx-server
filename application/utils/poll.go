package utils

import (
	"time"
)

type Poller struct {
	f        func()
	interval time.Duration
	ticker   *time.Ticker
	done     chan struct{}
}

func NewPoller(interval time.Duration, f func()) *Poller {
	return &Poller{
		f:        f,
		interval: interval,
	}
}

func (p *Poller) Start() {
	if p.ticker != nil {
		return
	}
	p.ticker = time.NewTicker(p.interval)
	p.done = make(chan struct{})
	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-p.ticker.C:
				p.f()
			}
		}
	}()
}

func (p *Poller) Stop() {
	if p.ticker == nil {
		return
	}
	p.ticker.Stop()
	p.ticker = nil
	close(p.done)
}
