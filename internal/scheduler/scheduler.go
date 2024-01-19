package scheduler

import (
	"fmt"
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
			err := s.fn()
			if err != nil {
				fmt.Printf("Failed to run scheduled function: %q\n", err)
			}
		}
	}()
}
