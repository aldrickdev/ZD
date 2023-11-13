package ports

import "zd/internal/applications/core/zendesk"

type APIPort interface {
	GetUserEvent() (*zendesk.UserEvent, error)
}
