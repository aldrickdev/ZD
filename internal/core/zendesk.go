package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"zd/internal/utils"
)

type ZendeskMock struct {
	userServiceLocation string
	eventPath           string
	userPath            string
}

func NewZendeskMock(userServiceLocation, eventPath, userPath string) ZendeskMock {
	return ZendeskMock{
		userServiceLocation: userServiceLocation,
		eventPath:           eventPath,
		userPath:            userPath,
	}
}
func (z ZendeskMock) GetUserEvent() (*UserEvent, error) {
	users, err := z.getAvailableUsers()
	if err != nil {
		return nil, fmt.Errorf("error while getting all available users: %s", err)
	}
	if len(users) == 0 {
		return nil, nil
	}

	events, err := z.getAvailableEvents()
	if err != nil {
		return nil, fmt.Errorf("error while getting all available events: %s", err)
	}
	if len(events) == 0 {
		return nil, nil
	}

	randomUser := randomSelection(users)
	randomEvent := randomSelection(events)

	return &UserEvent{
		UserID:  randomUser.Id,
		EventID: randomEvent.ID,
	}, nil
}
func (z ZendeskMock) GetFullUserEvent() (*FullUserEvent, error) {
	users, err := z.getAvailableUsers()
	if err != nil {
		return nil, fmt.Errorf("error while getting all available users: %s", err)
	}
	if len(users) == 0 {
		return nil, nil
	}

	events, err := z.getAvailableEvents()
	if err != nil {
		return nil, fmt.Errorf("error while getting all available events: %s", err)
	}
	if len(events) == 0 {
		return nil, nil
	}

	randomUser := randomSelection(users)
	randomEvent := randomSelection(events)

	return &FullUserEvent{
		User:  randomUser,
		Event: randomEvent,
	}, nil
}
func (z ZendeskMock) getAvailableEvents() ([]Event, error) {

	requestURL := fmt.Sprintf(
		"http://%s%s",
		z.userServiceLocation,
		z.eventPath,
	)
	data, err := utils.GetRequest(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error while getting available events: %s", err)
	}

	events := []Event{}
	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %s", err)
	}

	return events, nil
}
func (z ZendeskMock) getAvailableUsers() ([]User, error) {
	requestURL := fmt.Sprintf(
		"http://%s%s",
		z.userServiceLocation,
		z.userPath,
	)
	data, err := utils.GetRequest(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error while getting available users: %s", err)
	}

	users := []User{}
	err = json.Unmarshal(data, &users)
	if err != nil {
		fmt.Printf("Error Here: %q\n", err)
		return nil, fmt.Errorf("error unmarshaling json: %s", err)
	}

	return users, nil
}
func randomSelection[O User | Event](obj []O) O {
	randomNumberGenerator := rand.New(rand.NewSource(time.Now().Unix()))
	lastIndex := len(obj) - 1
	randomNumber := randomNumberGenerator.Intn(lastIndex)

	return obj[randomNumber]
}
