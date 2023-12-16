package schedule

import (
	"math/rand"
	"time"
	"zd/internal/core/ports"
)

type ScheduledFunc func() error

type Schedule struct {
	zendeskService ports.ZendeskService
	maxInterval    uint
	fn             ScheduledFunc
}

func New(zs ports.ZendeskService, mi uint, fn ScheduledFunc) *Schedule {
	return &Schedule{
		zendeskService: zs,
		maxInterval:    mi,
		fn:             fn,
	}
}

func (s Schedule) Run() {
	go func() {
		randomNumberGenerator := rand.New(rand.NewSource(time.Now().Unix()))
		for {
			randomInterval := randomNumberGenerator.Intn(int(s.maxInterval))
			time.Sleep(time.Second * time.Duration(randomInterval))

			s.fn()
		}
	}()
}
