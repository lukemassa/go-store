package store_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/lukemassa/go-store"
)

const testString = "This content is used for tests do NOT edit.\n"

// Actually a const but compiler won't allow it
var testBytes = []byte(testString)

func TestURLStore(t *testing.T) {

	s := store.NewHTTP("https://raw.githubusercontent.com/lukemassa/go-store/master/test_values/test.txt", store.Config{})
	s.Wait(5 * time.Second)
	res, err := s.Get()

	if bytes.Compare(testBytes, res) != 0 {
		t.Errorf("Actual value %s does not match expected %s", string(res), testString)
	}
	if err != nil {
		t.Errorf("Expected nil error, found %v", err)
	}
}
