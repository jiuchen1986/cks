package test

import (
	"fmt"

	"github.com/pkg/errors"
)

// ReturnError returns an error for testing
func ReturnError() error {
	return errors.New("error occured in test/test.go:ReturnError")
}

// ReturnNestedError returns an error from ReturnError for testing
func ReturnNestedError() error {
	return ReturnError()
}

// ReturnWrappedError returns an wrapped error from ReturnError for testing
func ReturnWrappedError() error {
	return errors.Wrap(ReturnError(), "wrapped error in test/test.go:ReturnWrappeddError")
}

// ReturnFmtError returns an wrapped error from ReturnError with fmt.Errorf for testing
func ReturnFmtError() error {
	return fmt.Errorf("error in test/test.go:ReturnFmtError: %v", ReturnError())
}
