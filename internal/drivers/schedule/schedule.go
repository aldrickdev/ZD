package schedule

import (
	"math/rand"
	"time"
	"zd/internal/core/ports"
)

type Schedule struct {
	zendeskService ports.ZendeskService
	maxInterval    uint
}

func New(zs ports.ZendeskService, mi uint) *Schedule {
	return &Schedule{
		zendeskService: zs,
		maxInterval:    mi,
	}
}

func (s Schedule) Run() {
	go func() {
		randomNumberGenerator := rand.New(rand.NewSource(time.Now().Unix()))
		for {
			randomInterval := randomNumberGenerator.Intn(int(s.maxInterval))
			time.Sleep(time.Second * time.Duration(randomInterval))

			s.zendeskService.GetUserEvent()
		}
	}()
}
