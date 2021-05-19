package test

import (
	"regexp"
	"testing"
)

// AssertErrorMatch will compare the error with the desired string, using regexp
// to match on the message
func AssertErrorMatch(t *testing.T, msg string, testErr error) {
	t.Helper()
	if msg == "" && testErr == nil {
		return
	}
	if msg != "" && testErr == nil {
		t.Fatalf("wanted error matching %s, got nil", msg)
	}
	match, err := regexp.MatchString(msg, testErr.Error())
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("failed to match error %s against %s", testErr, msg)
	}
}
