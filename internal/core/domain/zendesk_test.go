package domain

import (
	"testing"
)

func TestZendesk(t *testing.T) {
	z := NewZendesk("localhost:4001", "/api/v1/event", "/api/v1/user")
	got, err := z.GetUserEvent()
	if err != nil {
		t.Errorf("failed to generate a User Event: %s", err)
	}
	want := &UserEvent{
		UserID:  50,
		EventID: 25,
	}

	if got != want {
		t.Errorf("failed to generate the expected User Event: got %v, want %v", got, want)
	}
}
