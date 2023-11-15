package ports

import "zd/internal/core/domain"

type UserEventQueue interface {
	Publish(domain.UserEvent) error
}
