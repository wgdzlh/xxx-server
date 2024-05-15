package config

import (
	"sync"
	"time"
)

type Callback func(*SelfConfig)

const (
	Q_SEP         = ","
	Q_TIME_FORMAT = "2006-01-02 15:04:05"
	Q_DATE_FORMAT = "2006-01-02"

	APP        = "xxx-server"
	APP_PREFIX = "/" + APP
)

var (
	C         = &SelfConfig{}
	cMutex    sync.Mutex
	callbacks []Callback
	TZ        = time.FixedZone("CST", 60*60*8)
)

func setGlobalConfig(in *SelfConfig) {
	cMutex.Lock()
	defer cMutex.Unlock()
	for _, c := range callbacks {
		c(in)
	}
	C = in
}

func AddCallback(f Callback) {
	callbacks = append(callbacks, f)
}
