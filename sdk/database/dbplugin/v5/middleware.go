// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/grpc/status"
)

// ///////////////////////////////////////////////////
// Tracing Middleware
// ///////////////////////////////////////////////////

var (
	_ Database                = databaseTracingMiddleware{}
	_ logical.PluginVersioner = databaseTracingMiddleware{}
)

// databaseTracingMiddleware wraps a implementation of Database and executes
// trace logging on function call.
type databaseTracingMiddleware struct {
	next   Database
	logger log.Logger
}

func (mw databaseTracingMiddleware) PluginVersion() (resp logical.PluginVersion) {
	defer func(then time.Time) {
		mw.logger.Trace("version",
			"status", "finished",
			"version", resp,
			"took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("version", "status", "started")
	if versioner, ok := mw.next.(logical.PluginVersioner); ok {
		return versioner.PluginVersion()
	}
	return logical.EmptyPluginVersion
}

func (mw databaseTracingMiddleware) Initialize(ctx context.Context, req InitializeRequest) (resp InitializeResponse, err error) {
	defer func(then time.Time) {
		mw.logger.Trace("initialize",
			"status", "finished",
			"verify", req.VerifyConnection,
			"err", err,
			"took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("initialize", "status", "started")
	return mw.next.Initialize(ctx, req)
}

func (mw databaseTracingMiddleware) NewUser(ctx context.Context, req NewUserRequest) (resp NewUserResponse, err error) {
	defer func(then time.Time) {
		mw.logger.Trace("create user",
			"status", "finished",
			"err", err,
			"took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("create user",
		"status", "started")
	return mw.next.NewUser(ctx, req)
}

func (mw databaseTracingMiddleware) UpdateUser(ctx context.Context, req UpdateUserRequest) (resp UpdateUserResponse, err error) {
	defer func(then time.Time) {
		mw.logger.Trace("update user",
			"status", "finished",
			"err", err,
			"took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("update user", "status", "started")
	return mw.next.UpdateUser(ctx, req)
}

func (mw databaseTracingMiddleware) DeleteUser(ctx context.Context, req DeleteUserRequest) (resp DeleteUserResponse, err error) {
	defer func(then time.Time) {
		mw.logger.Trace("delete user",
			"status", "finished",
			"err", err,
			"took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("delete user",
		"status", "started")
	return mw.next.DeleteUser(ctx, req)
}

func (mw databaseTracingMiddleware) Type() (string, error) {
	return mw.next.Type()
}

func (mw databaseTracingMiddleware) Close() (err error) {
	defer func(then time.Time) {
		mw.logger.Trace("close",
			"status", "finished",
			"err", err,
			"took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("close",
		"status", "started")
	return mw.next.Close()
}

// ///////////////////////////////////////////////////
// Metrics Middleware Domain
// ///////////////////////////////////////////////////

var (
	_ Database                = databaseMetricsMiddleware{}
	_ logical.PluginVersioner = databaseMetricsMiddleware{}
)

// databaseMetricsMiddleware wraps an implementation of Databases and on
// function call logs metrics about this instance.
type databaseMetricsMiddleware struct {
	next Database

	typeStr string
}

func (mw databaseMetricsMiddleware) PluginVersion() logical.PluginVersion {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "PluginVersion"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "PluginVersion"}, now)
	}(time.Now())

	metrics.IncrCounter([]string{"database", "PluginVersion"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "PluginVersion"}, 1)

	if versioner, ok := mw.next.(logical.PluginVersioner); ok {
		return versioner.PluginVersion()
	}
	return logical.EmptyPluginVersion
}

func (mw databaseMetricsMiddleware) Initialize(ctx context.Context, req InitializeRequest) (resp InitializeResponse, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "Initialize"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "Initialize"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "Initialize", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "Initialize", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "Initialize"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "Initialize"}, 1)
	return mw.next.Initialize(ctx, req)
}

func (mw databaseMetricsMiddleware) NewUser(ctx context.Context, req NewUserRequest) (resp NewUserResponse, err error) {
	defer func(start time.Time) {
		metrics.MeasureSince([]string{"database", "NewUser"}, start)
		metrics.MeasureSince([]string{"database", mw.typeStr, "NewUser"}, start)

		if err != nil {
			metrics.IncrCounter([]string{"database", "NewUser", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "NewUser", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "NewUser"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "NewUser"}, 1)
	return mw.next.NewUser(ctx, req)
}

func (mw databaseMetricsMiddleware) UpdateUser(ctx context.Context, req UpdateUserRequest) (resp UpdateUserResponse, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "UpdateUser"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "UpdateUser"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "UpdateUser", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "UpdateUser", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "UpdateUser"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "UpdateUser"}, 1)
	return mw.next.UpdateUser(ctx, req)
}

func (mw databaseMetricsMiddleware) DeleteUser(ctx context.Context, req DeleteUserRequest) (resp DeleteUserResponse, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "DeleteUser"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "DeleteUser"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "DeleteUser", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "DeleteUser", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "DeleteUser"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "DeleteUser"}, 1)
	return mw.next.DeleteUser(ctx, req)
}

func (mw databaseMetricsMiddleware) Type() (string, error) {
	return mw.next.Type()
}

func (mw databaseMetricsMiddleware) Close() (err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "Close"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "Close"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "Close", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "Close", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "Close"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "Close"}, 1)
	return mw.next.Close()
}

// ///////////////////////////////////////////////////
// Error Sanitizer Middleware Domain
// ///////////////////////////////////////////////////

var (
	_ Database                = (*DatabaseErrorSanitizerMiddleware)(nil)
	_ logical.PluginVersioner = (*DatabaseErrorSanitizerMiddleware)(nil)
)

// DatabaseErrorSanitizerMiddleware wraps an implementation of Databases and
// sanitizes returned error messages
type DatabaseErrorSanitizerMiddleware struct {
	next      Database
	secretsFn secretsFn
}

type secretsFn func() map[string]string

func NewDatabaseErrorSanitizerMiddleware(next Database, secrets secretsFn) DatabaseErrorSanitizerMiddleware {
	return DatabaseErrorSanitizerMiddleware{
		next:      next,
		secretsFn: secrets,
	}
}

func (mw DatabaseErrorSanitizerMiddleware) Initialize(ctx context.Context, req InitializeRequest) (resp InitializeResponse, err error) {
	resp, err = mw.next.Initialize(ctx, req)
	return resp, mw.sanitize(err)
}

func (mw DatabaseErrorSanitizerMiddleware) NewUser(ctx context.Context, req NewUserRequest) (resp NewUserResponse, err error) {
	resp, err = mw.next.NewUser(ctx, req)
	return resp, mw.sanitize(err)
}

func (mw DatabaseErrorSanitizerMiddleware) UpdateUser(ctx context.Context, req UpdateUserRequest) (UpdateUserResponse, error) {
	resp, err := mw.next.UpdateUser(ctx, req)
	return resp, mw.sanitize(err)
}

func (mw DatabaseErrorSanitizerMiddleware) DeleteUser(ctx context.Context, req DeleteUserRequest) (DeleteUserResponse, error) {
	resp, err := mw.next.DeleteUser(ctx, req)
	return resp, mw.sanitize(err)
}

func (mw DatabaseErrorSanitizerMiddleware) Type() (string, error) {
	dbType, err := mw.next.Type()
	return dbType, mw.sanitize(err)
}

func (mw DatabaseErrorSanitizerMiddleware) Close() (err error) {
	return mw.sanitize(mw.next.Close())
}

func (mw DatabaseErrorSanitizerMiddleware) PluginVersion() logical.PluginVersion {
	if versioner, ok := mw.next.(logical.PluginVersioner); ok {
		return versioner.PluginVersion()
	}
	return logical.EmptyPluginVersion
}

// sanitize errors by removing any sensitive strings within their messages. This uses
// the secretsFn to determine what fields should be sanitized.
func (mw DatabaseErrorSanitizerMiddleware) sanitize(err error) error {
	if err == nil {
		return nil
	}
	if errwrap.ContainsType(err, new(url.Error)) {
		return errors.New("unable to parse connection url")
	}
	if mw.secretsFn == nil {
		return err
	}
	for find, replace := range mw.secretsFn() {
		if find == "" {
			continue
		}

		// Attempt to keep the status code attached to the
		// error while changing the actual error message
		s, ok := status.FromError(err)
		if ok {
			err = status.Error(s.Code(), strings.ReplaceAll(s.Message(), find, replace))
			continue
		}

		err = errors.New(strings.ReplaceAll(err.Error(), find, replace))
	}
	return err
}
