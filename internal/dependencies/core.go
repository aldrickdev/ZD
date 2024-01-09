package dependencies

import "zd/internal/core"

type Core interface {
	GetFullUserEvent() (*core.FullUserEvent, error)
}
