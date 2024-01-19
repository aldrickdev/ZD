package core

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	core_dependencies "zd/internal/core/dependencies"
	"zd/internal/utils"
)

type ZendeskMock struct {
	cache               core_dependencies.Cache
	userServiceLocation string
	eventPath           string
	userPath            string
}

func NewZendeskMock(cache core_dependencies.Cache, userServiceLocation, eventPath, userPath string) ZendeskMock {
	return ZendeskMock{
		cache:               cache,
		userServiceLocation: userServiceLocation,
		eventPath:           eventPath,
		userPath:            userPath,
	}
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

	// Check to see if the Event Data has already been cached
	result, err := z.cache.CheckCache(context.Background(), requestURL)
	if err == nil {
		// Cache hit
		var eventCache []Event
		err = json.Unmarshal([]byte(result), &eventCache)
		if err == nil {
			return eventCache, nil

		} else {
			fmt.Printf("Failed to unmarshal the event data cache: %q\n", err)
		}
	}

	// Request up-to-date event data from the User Service
	data, err := utils.GetRequest(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error while getting available events: %s", err)
	}

	events := []Event{}
	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %s", err)
	}

	// Cache the latest Event Data
	_, err = z.cache.CacheData(context.Background(), requestURL, string(data), 5*time.Minute)
	if err != nil {
		fmt.Printf("Failed to cache the Event Data: %q\n", err)
	} else {
		fmt.Println("Event Data has been cached")
	}

	return events, nil
}
func (z ZendeskMock) getAvailableUsers() ([]User, error) {
	requestURL := fmt.Sprintf(
		"http://%s%s",
		z.userServiceLocation,
		z.userPath,
	)

	// Check to see if the user data has already been cached
	result, err := z.cache.CheckCache(context.Background(), requestURL)
	if err == nil {
		// Cache hit
		var userCache []User
		err = json.Unmarshal([]byte(result), &userCache)
		if err == nil {
			return userCache, nil
		} else {
			fmt.Printf("Failed to unmarshal the user data cache: %q\n", err)
		}
	}

	// Request up-to-date user data from the User Service
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

	// Cache the latest User Data
	_, err = z.cache.CacheData(context.Background(), requestURL, string(data), 1*time.Minute)
	if err != nil {
		fmt.Printf("Failed to cache the user data: %q\n", err)
	} else {
		fmt.Println("User data has been cached")
	}

	return users, nil
}
func randomSelection[O User | Event](obj []O) O {
	randomNumberGenerator := rand.New(rand.NewSource(time.Now().Unix()))
	lastIndex := len(obj) - 1
	randomNumber := randomNumberGenerator.Intn(lastIndex)

	return obj[randomNumber]
}
