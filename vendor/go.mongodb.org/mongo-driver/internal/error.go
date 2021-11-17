// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package internal

import (
	"fmt"
)

// WrappedError represents an error that contains another error.
type WrappedError interface {
	// Message gets the basic message of the error.
	Message() string
	// Inner gets the inner error if one exists.
	Inner() error
}

// RolledUpErrorMessage gets a flattened error message.
func RolledUpErrorMessage(err error) string {
	if wrappedErr, ok := err.(WrappedError); ok {
		inner := wrappedErr.Inner()
		if inner != nil {
			return fmt.Sprintf("%s: %s", wrappedErr.Message(), RolledUpErrorMessage(inner))
		}

		return wrappedErr.Message()
	}

	return err.Error()
}

//UnwrapError attempts to unwrap the error down to its root cause.
func UnwrapError(err error) error {

	switch tErr := err.(type) {
	case WrappedError:
		return UnwrapError(tErr.Inner())
	case *multiError:
		return UnwrapError(tErr.errors[0])
	}

	return err
}

// WrapError wraps an error with a message.
func WrapError(inner error, message string) error {
	return &wrappedError{message, inner}
}

// WrapErrorf wraps an error with a message.
func WrapErrorf(inner error, format string, args ...interface{}) error {
	return &wrappedError{fmt.Sprintf(format, args...), inner}
}

// MultiError combines multiple errors into a single error. If there are no errors,
// nil is returned. If there is 1 error, it is returned. Otherwise, they are combined.
func MultiError(errors ...error) error {

	// remove nils from the error list
	var nonNils []error
	for _, e := range errors {
		if e != nil {
			nonNils = append(nonNils, e)
		}
	}

	switch len(nonNils) {
	case 0:
		return nil
	case 1:
		return nonNils[0]
	default:
		return &multiError{
			message: "multiple errors encountered",
			errors:  nonNils,
		}
	}
}

type multiError struct {
	message string
	errors  []error
}

func (e *multiError) Message() string {
	return e.message
}

func (e *multiError) Error() string {
	result := e.message
	for _, e := range e.errors {
		result += fmt.Sprintf("\n  %s", e)
	}
	return result
}

func (e *multiError) Errors() []error {
	return e.errors
}

type wrappedError struct {
	message string
	inner   error
}

func (e *wrappedError) Message() string {
	return e.message
}

func (e *wrappedError) Error() string {
	return RolledUpErrorMessage(e)
}

func (e *wrappedError) Inner() error {
	return e.inner
}
