package batch

import (
	"testing"
	"zd/internal/core/domain"
)

func TestBatch(t *testing.T) {
	t.Run("Adding User Event", func(t *testing.T) {
		batcher := New()

		userEvent := domain.UserEvent{
			UserID:  1,
			EventID: 1,
		}

		batcher.Add(&userEvent)
		currentBatchLength := len(batcher.data)
		if currentBatchLength != 1 {
			t.Errorf("expected %v, got %v", 1, currentBatchLength)
		}
	})

	t.Run("Draining the batch", func(t *testing.T) {
		batcher := New()

		userEvent := domain.UserEvent{
			UserID:  1,
			EventID: 1,
		}

		batcher.Add(&userEvent)
		currentBatchLength := len(batcher.data)
		if currentBatchLength != 1 {
			t.Errorf("expected %v, got %v", 1, currentBatchLength)
		}

		data := batcher.Drain()
		returnedLength := len(data)

		if returnedLength != 1 {
			t.Errorf("got %v, expected %v", returnedLength, 1)
		}
	})
}
