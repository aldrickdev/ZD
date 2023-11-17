package domain

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"zd/internal/utils"
)

type Zendesk struct {
	userServiceLocation string
	eventPath           string
	userPath            string
}

func NewZendesk(userServiceLocation, eventPath, userPath string) Zendesk {
	return Zendesk{
		userServiceLocation: userServiceLocation,
		eventPath:           eventPath,
		userPath:            userPath,
	}
}

func (z Zendesk) GetUserEvent() (*UserEvent, error) {
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
		UserID:  randomUser.ID,
		EventID: randomEvent.ID,
	}, nil
}
func (z Zendesk) getAvailableEvents() ([]Event, error) {
	// "http://localhost:4001/api/v1/event",
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
func (z Zendesk) getAvailableUsers() ([]User, error) {
	// "http://localhost:4001/api/v1/user"
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
