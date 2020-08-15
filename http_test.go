package store_test

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/lukemassa/go-store"
)

func TestURLStore(t *testing.T) {
	// TODO: mock http somehow
	s := store.NewHTTP("http://icanhazip.com", store.Config{})
	s.Wait(5 * time.Second)
	res, err := s.Get()
	parsedRes := strings.Trim(string(res), "\n")
	ipaddr := net.ParseIP(parsedRes)

	if ipaddr == nil {
		t.Errorf("Actual value does not look like an IP address: %s", parsedRes)
	}
	if err != nil {
		t.Errorf("Expected nil error, found %v", err)
	}
}
