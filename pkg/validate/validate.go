// Package validate provides util error validations.
package validate

import (
	"errors"
)

// IsEmpty returns if a string is empty
func IsEmpty(s string) bool {
	return s == ""
}

// Cause iterates through all the wrapped errors until the root
// error value is reached.
func Cause(err error) error {
	root := err
	for {
		if err = errors.Unwrap(root); err == nil {
			return root
		}
		root = err
	}
}
