package scheduler

import (
	"math/rand"
	"time"
)

type ScheduledFunc func() error

type Schedule struct {
	maxInterval int
	random      bool
	fn          ScheduledFunc
}

func New(mi int, isRandom bool, fn ScheduledFunc) *Schedule {
	return &Schedule{
		maxInterval: mi,
		random:      isRandom,
		fn:          fn,
	}
}

func (s Schedule) Run() {
	go func() {
		randomNumberGenerator := rand.New(rand.NewSource(time.Now().Unix()))
		for {
			wait := s.maxInterval
			if s.random {
				wait = randomNumberGenerator.Intn(int(s.maxInterval))
			}
			time.Sleep(time.Second * time.Duration(wait))
			s.fn()
		}
	}()
}
