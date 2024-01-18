package service

import (
	"zd/internal/core"
	service_dependencies "zd/internal/service/dependencies"
)

const (
	CallbackTypeImmediate = "CALLBACK_TYPE_IMMEDIATE"
	CallbackTypeLatest    = "CALLBACK_TYPE_LATEST"
)

type service struct {
	// Dependencies
	queue service_dependencies.QueueBroker
	core  service_dependencies.Core

	// Callback Store
	publishingCallbacks       []func(*core.FullUserEvent) error
	latestPublishingCallbacks []func(*core.FullUserEvent) error
	latestFullUserEvent       *core.FullUserEvent
}

func New(queueBroker service_dependencies.QueueBroker, core service_dependencies.Core) service {
	return service{
		queue: queueBroker,
		core:  core,
	}
}

func (s *service) PublishNewUserEvent() error {
	ue, err := s.core.GetFullUserEvent()
	if err != nil {
		return err
	}

	s.latestFullUserEvent = ue

	for _, callback := range s.publishingCallbacks {
		err := callback(ue)
		return err
	}

	return nil
}

func (s *service) PublishLatestUserEvent() error {
	for _, callback := range s.latestPublishingCallbacks {
		err := callback(s.latestFullUserEvent)
		return err
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
