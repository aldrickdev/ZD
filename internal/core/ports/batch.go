package ports

import "zd/internal/core/domain"

// Port that the drivers will use to make use of the core code
type Batch interface {
	Add(*domain.UserEvent)
	Drain() []*domain.UserEvent
}
