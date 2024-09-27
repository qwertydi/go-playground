package util

import (
	"time"
)

type TimeServiceHandler interface {
	GetCurrentTime() time.Time
}

type TimeServiceHandlerImpl struct{}

func (TimeServiceHandlerImpl) GetCurrentTime() time.Time {
	return time.Now()
}

// TimeService initializes a new instance of Service
func TimeService() *TimeServiceHandlerImpl {
	return &TimeServiceHandlerImpl{}
}
