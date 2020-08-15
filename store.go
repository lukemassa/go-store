package store

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

const unsetErrorString = "Value is not yet set"

// Config configure the store
type Config struct {
	// How frequently the store is updated
	// Defaults to 60 seconds
	UpdateInterval time.Duration

	// The max amount of time for the update function to run
	// Must be less than the update interval
	// Defaults to 50 seconds
	UpdateTimeout time.Duration

	// How long the updated value is valid for. After this time,
	// Get() will return an error. Should be several multiples of
	// the UpdateInterval so updates are allowed to fail.
	//
	// If set to 0 (the default), entries are always valid.
	// However, in most cases this is inappropriate and should be set
	ResultValidity time.Duration
}

// Store a store to hold the value
type Store struct {
	config     Config
	cache      *cache.Cache
	updateFunc func(context.Context) (interface{}, error)
}

type cacheValue struct {
	value interface{}
	err   error
}

// New create a new store
// Function called to get the current value of the store
// It will be passed a new context that deadlines at
// updateTimeout, which it should respect
func New(updateFunc func(context.Context) (interface{}, error), config Config) Store {
	store := Store{
		updateFunc: updateFunc,
	}

	if config.UpdateInterval == 0 {
		config.UpdateInterval = 60 * time.Second
	}
	if config.UpdateTimeout == 0 {
		config.UpdateTimeout = 50 * time.Second
	}
	if config.UpdateTimeout > config.UpdateInterval {
		panic("Update timeout cannot be longer than update interval")
	}
	store.config = config
	store.cache = cache.New(config.ResultValidity, 0)

	// Seed the cache with an empty value so we can differentiate
	// between never set up with expired
	store.cache.Set("value", cacheValue{
		value: nil,
		err:   errors.New(unsetErrorString),
	}, 0)
	store.start()

	return store
}

func (s *Store) start() {

	ticker := time.NewTicker(s.config.UpdateInterval)
	go func() {
		for ; true; <-ticker.C {
			s.update()
		}
	}()
}
func (s *Store) update() {
	ctx, cancel := context.WithTimeout(context.TODO(), s.config.UpdateTimeout)
	defer cancel()

	val, err := s.updateFunc(ctx)
	s.cache.Set("value", cacheValue{
		value: val,
		err:   err,
	}, s.config.ResultValidity)
}

// Get the value from the store
func (s *Store) Get() (interface{}, error) {
	val, ok := s.cache.Get("value")
	if !ok {
		return nil, errors.New("Value has expired")
	}
	// This casting is safe because I'm the only one who puts anything into the value
	return val.(cacheValue).value, val.(cacheValue).err
}

// Wait until the cache contains a successful run
func (s *Store) Wait(maxWait time.Duration) error {
	waitPeriod := 1 * time.Millisecond
	maxWaitPeriod := 1 * time.Second
	end := time.Now().Add(maxWait)

	for !time.Now().After(end) {
		_, err := s.Get()
		if err == nil || err.Error() != unsetErrorString {
			return nil
		}
		time.Sleep(waitPeriod)
		waitPeriod *= 2
		if waitPeriod > maxWaitPeriod {
			waitPeriod = maxWaitPeriod
		}
	}
	return errors.New("Still unset")
}
