package store

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPStore a store to hold the value
type HTTPStore struct {
	store Store
}

// NewHTTP store the result of an http get call
func NewHTTP(url string, config Config) HTTPStore {
	store := New(func(ctx context.Context) (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}, config)

	return HTTPStore{
		store: store,
	}
}

// Get gets the result of the http call
func (s *HTTPStore) Get() ([]byte, error) {
	res, err := s.store.Get()
	if err != nil {
		return nil, err
	}
	return res.([]byte), err
}

// Wait until value is available
func (s *HTTPStore) Wait(maxWait time.Duration) error {
	return s.store.Wait(maxWait)
}
