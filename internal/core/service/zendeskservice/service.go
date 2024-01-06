package zendeskservice

import (
	"fmt"
	"zd/internal/core/domain"
	"zd/internal/core/ports"
)

const (
	CallbackTypeImmediate = "CALLBACK_TYPE_IMMEDIATE"
	CallbackTypeLatest    = "CALLBACK_TYPE_LATEST"
)

type service struct {
	q                         ports.UserEventQueue
	b                         ports.Batch
	z                         domain.ZendeskMock
	publishingCallbacks       []func(domain.FullUserEvent) error
	latestPublishingCallbacks []func(domain.FullUserEvent) error
	latestFullUserEvent       domain.FullUserEvent
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
	ue, err := s.z.GetUserEvent()
	if err != nil {
		return err
	}

	s.b.Add(ue)

	err = s.q.Publish(*ue)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) PublishNewUserEvent() error {
	fmt.Println("New")
	ue, err := s.z.GetFullUserEvent()
	if err != nil {
		return err
	}

	s.latestFullUserEvent = *ue

	for _, callback := range s.publishingCallbacks {
		callback(*ue)
	}

	return nil
}

func (s *service) PublishLatestUserEvent() error {
	fmt.Println("Latest")
	for _, callback := range s.latestPublishingCallbacks {
		callback(s.latestFullUserEvent)
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

func (s *service) RegisterPublishingCallback(callback func(domain.FullUserEvent) error, callbackType string) {
	switch callbackType {
	case CallbackTypeImmediate:
		s.publishingCallbacks = append(s.publishingCallbacks, callback)

	case CallbackTypeLatest:
		s.latestPublishingCallbacks = append(s.latestPublishingCallbacks, callback)
	}
}
