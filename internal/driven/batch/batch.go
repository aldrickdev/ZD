package batch

import (
	"zd/internal/core/domain"
)

type Batch struct {
	data []*domain.UserEvent
}

func New() *Batch {
	return &Batch{}
}

func (b *Batch) Add(data *domain.UserEvent) {
	b.data = append(b.data, data)
}

func (b *Batch) Drain() []*domain.UserEvent {
	data := b.data
	b.data = []*domain.UserEvent{}

	return data
}
