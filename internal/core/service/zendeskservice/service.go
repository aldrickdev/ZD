package zendeskservice

import (
	"zd/internal/core/domain"
	"zd/internal/core/ports"
)

type service struct {
	// Driven
	q ports.UserEventQueue

	// Core
	z domain.Zendesk
}

func New(q ports.UserEventQueue, userServiceLocation, eventPath, userPath string) service {
	z := domain.NewZendesk(
		userServiceLocation,
		eventPath,
		userPath,
	)

	return service{
		q: q,
		z: z,
	}
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
