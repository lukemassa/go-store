package store_test

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/lukemassa/go-store"
)

func TestGitStore(t *testing.T) {

	s := store.NewGitRepo("https://github.com/lukemassa/go-store.git", store.Config{})
	err := s.Wait(5 * time.Second)
	if err != nil {
		t.Errorf("Unexpected error after wait: %v", err)
		return
	}
	res, err := s.Get()
	if err != nil {
		t.Errorf("Unpexected error from get: %v", err)
		return
	}
	file, err := res.Filesystem.Open("test_values/test.txt")
	if err != nil {
		t.Errorf("Unexpected error opening test file: %v", err)
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Errorf("Unexpected error reading file: %v", err)
		return
	}
	if bytes.Compare(testBytes, content) != 0 {
		t.Errorf("Actual value %s does not match expected %s", string(content), testString)
	}

}
