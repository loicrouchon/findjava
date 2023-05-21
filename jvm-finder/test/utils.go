package test

import (
	"reflect"
	"strings"
	"testing"
)

func AssertEquals(t *testing.T, description string, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf(`Expecting
    %s
to be equal to:
    %#v
but got:
    %#v`,
			description, expected, actual)
	}
}

func AssertErrorEquals(t *testing.T, description string, expected error, actual error) {
	if actual != expected {
		t.Fatalf(`Expecting
    %s
to fail with:
    %v
but got:
    %v`,
			description, expected, actual)
	}
}

func AssertErrorContains(t *testing.T, description string, expected string, actual error) {
	if expected == "" && actual == nil {
		return
	}
	if !strings.Contains(actual.Error(), expected) {
		t.Fatalf(`Expecting error
    %s
to contain message:
    %v
but got:
    %v`,
			description, expected, actual)
	}
}

func AssertNoError(t *testing.T, description string, actual error) {
	if actual != nil {
		t.Fatalf(`Expecting
    %s
to return no error but got:
    %v`,
			description, actual)
	}
}
