package service

import (
	"zd/internal/core"
	"zd/internal/dependencies"
)

const (
	CallbackTypeImmediate = "CALLBACK_TYPE_IMMEDIATE"
	CallbackTypeLatest    = "CALLBACK_TYPE_LATEST"
)

type service struct {
	// Dependencies
	queue dependencies.QueueBroker
	core  dependencies.Core

	// Callback Store
	publishingCallbacks       []func(*core.FullUserEvent) error
	latestPublishingCallbacks []func(*core.FullUserEvent) error
	latestFullUserEvent       *core.FullUserEvent
}

func New(queueBroker dependencies.QueueBroker, userServiceLocation, eventPath, userPath string) service {
	z := core.NewZendeskMock(
		userServiceLocation,
		eventPath,
		userPath,
	)

	return service{
		queue: queueBroker,
		core:  z,
	}
}

func (s *service) PublishNewUserEvent() error {
	ue, err := s.core.GetFullUserEvent()
	if err != nil {
		return err
	}

	s.latestFullUserEvent = ue

	for _, callback := range s.publishingCallbacks {
		callback(ue)
	}

	return nil
}

func (s *service) PublishLatestUserEvent() error {
	for _, callback := range s.latestPublishingCallbacks {
		callback(s.latestFullUserEvent)
	}

	return nil
}

func (s *service) RegisterPublishingCallback(callback func(*core.FullUserEvent) error, callbackType string) {
	switch callbackType {
	case CallbackTypeImmediate:
		s.publishingCallbacks = append(s.publishingCallbacks, callback)

	case CallbackTypeLatest:
		s.latestPublishingCallbacks = append(s.latestPublishingCallbacks, callback)
	}
}
