package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/lukemassa/go-store"
)

func TestStoreBasic(t *testing.T) {
	s := store.New(func(ctx context.Context) (interface{}, error) {
		return 5, nil
	},
		store.Config{},
	)

	// If Get() is called immediately, the store will not have a chance
	// yet to set up the value
	res, err := s.Get()
	if res != nil {
		t.Errorf("Expected Get() to return nil, returned %v", res)
	}
	if err.Error() != "Value is not yet set" {
		t.Errorf("Expected Get() to have error, got %v", err)
	}
	err = s.Wait(1 * time.Second)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	res, err = s.Get()
	if res != 5 {
		t.Errorf("Expected Get() to return 5, returned %v", res)
	}
	if err != nil {
		t.Errorf("Expected Get() to have no error, got %v", err)
	}
}
