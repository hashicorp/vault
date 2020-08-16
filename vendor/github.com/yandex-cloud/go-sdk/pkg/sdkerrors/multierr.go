// Copyright (c) 2018 Yandex LLC. All rights reserved.
// Author: Dmitry Novikov <novikoff@yandex-team.ru>

package sdkerrors

import (
	"fmt"
	"strings"
)

type multerr struct {
	errs []error
}

func (e *multerr) Errors() []error {
	return e.errs
}

func (e *multerr) Error() string {
	lines := make([]string, len(e.errs))
	for k, v := range e.errs {
		lines[k] = v.Error()
	}
	return strings.Join(lines, "\n")
}

func Errors(err error) []error {
	if err == nil {
		return nil
	}
	switch err := err.(type) {
	case interface {
		Errors() []error
	}:
		// go.uber.org/multierr
		return err.Errors()
	case interface {
		WrappedErrors() []error
	}:
		// github.com/hashicorp/go-multierror
		return err.WrappedErrors()
	default:
	}
	return []error{err}
}

func Append(lhs, rhs error) error {
	if lhs == nil {
		return rhs
	} else if rhs == nil {
		return lhs
	}
	var result []error
	result = append(result, Errors(lhs)...)
	result = append(result, Errors(rhs)...)
	return &multerr{result}
}

func CombineGoroutines(funcs ...func() error) error {
	errChan := make(chan error, len(funcs))
	for _, f := range funcs {
		go func(f func() error) {
			var err error
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("Panic recovered: %v", r)
				} else {
					errChan <- err
				}
			}()
			err = f()
		}(f)
	}
	var errs error
	for i := 0; i < cap(errChan); i++ {
		errs = Append(errs, <-errChan)
	}
	return errs
}
