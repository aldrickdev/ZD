package zendeskservice

import (
	"zd/internal/core/domain"
	"zd/internal/core/ports"
)

type service struct {
	// Driven
	q ports.UserEventQueue
	b ports.Batch

	// Core
	z domain.ZendeskMock
}

func New(q ports.UserEventQueue, b ports.Batch, userServiceLocation, eventPath, userPath string) service {
	z := domain.NewZendeskMock(
		userServiceLocation,
		eventPath,
		userPath,
	)

	return service{
		q: q,
		b: b,
		z: z,
	}
}

// Not sure if pointer is needed here
func (s service) BatchUserEvent() error {
	ue, err := s.z.GetUserEvent()
	if err != nil {
		return err
	}

	s.b.Add(ue)

	return nil
}

func (s service) PublishBatch() error {
	events := s.b.Drain()

	err := s.q.PublishBatch(events)
	if err != nil {
		return err
	}

	return nil
}

func (s service) GenerateUserEvent() error {
	_, err := s.z.GetUserEvent()
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetUserEvent() (*domain.UserEvent, error) {
	ue, err := s.z.GetUserEvent()
	if err != nil {
		return nil, err
	}

	err = s.q.Publish(*ue)
	if err != nil {
		return nil, err
	}
	return ue, nil
}
