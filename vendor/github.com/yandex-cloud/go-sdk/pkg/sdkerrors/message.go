// Copyright (c) 2018 Yandex LLC. All rights reserved.
// Author: Dmitry Novikov <novikoff@yandex-team.ru>

package sdkerrors

import (
	"fmt"

	"google.golang.org/grpc/status"
)

func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return WithMessage(err, fmt.Sprintf(format, args...))
}

func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}

	withMessage := errWithMessage{err, message}
	if _, ok := err.(statusErr); ok {
		return &statusErrWithMessage{withMessage}
	}
	return &withMessage
}

type statusErr interface {
	GRPCStatus() *status.Status
}

type errWithMessage struct {
	err     error
	message string
}

type statusErrWithMessage struct {
	errWithMessage
}

func (e *errWithMessage) Error() string {
	return e.message + ": " + e.err.Error()
}

func (e *errWithMessage) Cause() error {
	return e.err
}

func (e *statusErrWithMessage) GRPCStatus() *status.Status {
	return status.Convert(e.err)
}
