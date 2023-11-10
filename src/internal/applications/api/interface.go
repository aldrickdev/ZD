package api

import "zd/internal/applications/core/zendesk"

type EventProducer interface {
	GetUserEvent() (*zendesk.UserEvent, error)
}
