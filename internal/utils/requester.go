package utils

import (
	"fmt"
	"io"
	"net/http"
)

type Requester interface {
	Get(string) ([]byte, error)
}

type Request struct{}

func NewRequester() Requester {
	return Request{}
}

func (r Request) Get(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error make request: %s", err)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the body: %s", err)
	}

	return data, nil
}
