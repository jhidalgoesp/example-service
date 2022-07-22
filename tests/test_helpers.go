// Package tests contains util test assertions and third party mocks.
package tests

import (
	"errors"
	"reflect"
	"testing"
)

// AssertNilError asserts an error is nil
func AssertNilError(t testing.TB, val interface{}) {
	t.Helper()

	if val != nil {
		t.Errorf("Expected nil, got: %s.", val)
	}
}

// AssertStrings asserts two strings are equal
func AssertStrings(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

// AssertSameInt asserts two int are equal
func AssertSameInt(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

// AssertSameErrors asserts two errors are equal
func AssertSameErrors(t testing.TB, got, want error) {
	if !errors.Is(got, want) {
		t.Errorf("Expected %v, got %v", got, want)
	}
}

// AssertNotNilError asserts a type is not nil
func AssertNotNilError(t testing.TB, val interface{}) {
	t.Helper()
	if val == nil {
		t.Errorf("Expected error to be not nil.")
	}
}

// AssertKindIsStruct asserts a type is kind struct
func AssertKindIsStruct(t testing.TB, val interface{}) {
	t.Helper()
	valType := reflect.TypeOf(val).Kind()

	if valType != reflect.Struct {
		t.Errorf("Expected a struct, got %v", valType)
	}
}

// AssertKindIsPointer asserts a type is kind struct
func AssertKindIsPointer(t testing.TB, val interface{}) {
	t.Helper()
	valType := reflect.TypeOf(val).Kind()

	if valType != reflect.Ptr {
		t.Errorf("Expected a pointer, got %v", valType)
	}
}

// AssertSliceLength asserts a slice has a defined length
func AssertSliceLength(t testing.TB, got, want int) {
	if got != want {
		t.Errorf("Expected %d items, got %d items.", got, want)
	}
}
