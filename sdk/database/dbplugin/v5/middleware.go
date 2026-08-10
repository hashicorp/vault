// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	metrics "github.com/hashicorp/go-metrics/compat"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/jackc/pgx/v5/pgconn"
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

	// namespace, mountPoint and connectionName identify the mount that this
	// instance serves. When set, they are attached to emitted metrics as
	// labels. All are empty unless the operator opted in, in which case emission
	// is unlabeled and therefore identical to the behavior before these fields
	// existed.
	namespace      string
	mountPoint     string
	connectionName string
}

// labels builds the telemetry labels identifying this plugin instance, or nil
// when the operator has not opted in to per-mount metrics.
//
// A fresh slice is returned on every call because go-metrics appends to the
// slice it is handed, so a shared slice could be mutated by an earlier
// emission.
func (mw databaseMetricsMiddleware) labels() []metrics.Label {
	if mw.namespace == "" && mw.mountPoint == "" && mw.connectionName == "" {
		return nil
	}

	labels := make([]metrics.Label, 0, 3)
	if mw.namespace != "" {
		labels = append(labels, metrics.Label{Name: "namespace", Value: mw.namespace})
	}
	if mw.mountPoint != "" {
		labels = append(labels, metrics.Label{Name: "mount_point", Value: mw.mountPoint})
	}
	if mw.connectionName != "" {
		labels = append(labels, metrics.Label{Name: "connection_name", Value: mw.connectionName})
	}
	return labels
}

// incrOperation counts an attempted operation, both across all database plugins
// and for this plugin type.
func (mw databaseMetricsMiddleware) incrOperation(operation string) {
	metrics.IncrCounterWithLabels([]string{"database", operation}, 1, mw.labels())
	metrics.IncrCounterWithLabels([]string{"database", mw.typeStr, operation}, 1, mw.labels())
}

// incrOperationError counts a failed operation.
func (mw databaseMetricsMiddleware) incrOperationError(operation string) {
	metrics.IncrCounterWithLabels([]string{"database", operation, "error"}, 1, mw.labels())
	metrics.IncrCounterWithLabels([]string{"database", mw.typeStr, operation, "error"}, 1, mw.labels())
}

// measureOperationSince records how long an operation took.
func (mw databaseMetricsMiddleware) measureOperationSince(operation string, start time.Time) {
	metrics.MeasureSinceWithLabels([]string{"database", operation}, start, mw.labels())
	metrics.MeasureSinceWithLabels([]string{"database", mw.typeStr, operation}, start, mw.labels())
}

func (mw databaseMetricsMiddleware) PluginVersion() logical.PluginVersion {
	defer func(start time.Time) {
		mw.measureOperationSince("PluginVersion", start)
	}(time.Now())

	mw.incrOperation("PluginVersion")

	if versioner, ok := mw.next.(logical.PluginVersioner); ok {
		return versioner.PluginVersion()
	}
	return logical.EmptyPluginVersion
}

func (mw databaseMetricsMiddleware) Initialize(ctx context.Context, req InitializeRequest) (resp InitializeResponse, err error) {
	defer func(start time.Time) {
		mw.measureOperationSince("Initialize", start)

		if err != nil {
			mw.incrOperationError("Initialize")
		}
	}(time.Now())

	mw.incrOperation("Initialize")
	return mw.next.Initialize(ctx, req)
}

func (mw databaseMetricsMiddleware) NewUser(ctx context.Context, req NewUserRequest) (resp NewUserResponse, err error) {
	defer func(start time.Time) {
		mw.measureOperationSince("NewUser", start)

		if err != nil {
			mw.incrOperationError("NewUser")
		}
	}(time.Now())

	mw.incrOperation("NewUser")
	return mw.next.NewUser(ctx, req)
}

func (mw databaseMetricsMiddleware) UpdateUser(ctx context.Context, req UpdateUserRequest) (resp UpdateUserResponse, err error) {
	defer func(start time.Time) {
		mw.measureOperationSince("UpdateUser", start)

		if err != nil {
			mw.incrOperationError("UpdateUser")
		}
	}(time.Now())

	mw.incrOperation("UpdateUser")
	return mw.next.UpdateUser(ctx, req)
}

func (mw databaseMetricsMiddleware) DeleteUser(ctx context.Context, req DeleteUserRequest) (resp DeleteUserResponse, err error) {
	defer func(start time.Time) {
		mw.measureOperationSince("DeleteUser", start)

		if err != nil {
			mw.incrOperationError("DeleteUser")
		}
	}(time.Now())

	mw.incrOperation("DeleteUser")
	return mw.next.DeleteUser(ctx, req)
}

func (mw databaseMetricsMiddleware) Type() (string, error) {
	return mw.next.Type()
}

func (mw databaseMetricsMiddleware) Close() (err error) {
	defer func(start time.Time) {
		mw.measureOperationSince("Close", start)

		if err != nil {
			mw.incrOperationError("Close")
		}
	}(time.Now())

	mw.incrOperation("Close")
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

	if errwrap.ContainsType(err, new(url.Error)) || errwrap.ContainsType(err, new(pgconn.ParseConfigError)) {
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
