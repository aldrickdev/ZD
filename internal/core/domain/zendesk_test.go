package domain

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type requesterMock struct{}

func (r requesterMock) Get(url string) ([]byte, error) {
	urlSplit := strings.Split(url, "/")
	urlSplitLen := len(urlSplit)
	requestedResource := urlSplit[urlSplitLen-1]

	switch requestedResource {
	case "event":
		testResource := []Event{
			{
				ID:     10,
				Name:   "Test Event",
				Points: 10,
			},
		}
		testData, err := json.Marshal(testResource)
		if err != nil {
			return nil, err
		}
		return testData, nil

	case "user":
		testResource := []User{
			{
				ID:   5,
				Name: "Test User",
				Pod:  "The Best",
			},
		}
		testData, err := json.Marshal(testResource)
		if err != nil {
			return nil, err
		}
		return testData, nil

	default:
		return nil, fmt.Errorf("an unsupported supported resource (%s) has been requested", requestedResource)
	}
}

func TestZendesk(t *testing.T) {
	requester := requesterMock{}
	z := NewZendesk(requester, "", "/event", "/user")

	got, err := z.GetUserEvent()
	if err != nil {
		t.Errorf("failed to generate a User Event: %s", err)
	}
	want := &UserEvent{
		UserID:  5,
		EventID: 10,
	}

	if *got != *want {
		t.Errorf("failed to generate the expected User Event: got %v, want %v", got, want)
	}
}
