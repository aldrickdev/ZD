package api

import "zd/internal/applications/core/zendesk"

type Application struct {
	EventProducer EventProducer
}

func NewApplication(ep EventProducer) *Application {
	return &Application{
		EventProducer: ep,
	}
}

func (a Application) GetUserEvent() (*zendesk.UserEvent, error) {
	return a.EventProducer.GetUserEvent()
}
