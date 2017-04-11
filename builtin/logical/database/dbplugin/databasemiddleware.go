package dbplugin

import (
	"time"

	metrics "github.com/armon/go-metrics"
	log "github.com/mgutz/logxi/v1"
)

// ---- Tracing Middleware Domain ----

// databaseTracingMiddleware wraps a implementation of DatabaseType and executes
// trace logging on function call.
type databaseTracingMiddleware struct {
	next   DatabaseType
	logger log.Logger

	typeStr string
}

func (mw *databaseTracingMiddleware) Type() string {
	return mw.next.Type()
}

func (mw *databaseTracingMiddleware) CreateUser(statements Statements, usernamePrefix string, expiration time.Time) (username string, password string, err error) {
	if mw.logger.IsTrace() {
		defer func(then time.Time) {
			mw.logger.Trace("database/CreateUser: finished", "type", mw.typeStr, "err", err, "took", time.Since(then))
		}(time.Now())

		mw.logger.Trace("database/CreateUser: starting", "type", mw.typeStr)
	}
	return mw.next.CreateUser(statements, usernamePrefix, expiration)
}

func (mw *databaseTracingMiddleware) RenewUser(statements Statements, username string, expiration time.Time) (err error) {
	if mw.logger.IsTrace() {
		defer func(then time.Time) {
			mw.logger.Trace("database/RenewUser: finished", "type", mw.typeStr, "err", err, "took", time.Since(then))
		}(time.Now())

		mw.logger.Trace("database/RenewUser: starting", "type", mw.typeStr)
	}
	return mw.next.RenewUser(statements, username, expiration)
}

func (mw *databaseTracingMiddleware) RevokeUser(statements Statements, username string) (err error) {
	if mw.logger.IsTrace() {
		defer func(then time.Time) {
			mw.logger.Trace("database/RevokeUser: finished", "type", mw.typeStr, "err", err, "took", time.Since(then))
		}(time.Now())

		mw.logger.Trace("database/RevokeUser: starting", "type", mw.typeStr)
	}
	return mw.next.RevokeUser(statements, username)
}

func (mw *databaseTracingMiddleware) Initialize(conf map[string]interface{}, verifyConnection bool) (err error) {
	if mw.logger.IsTrace() {
		defer func(then time.Time) {
			mw.logger.Trace("database/Initialize: finished", "type", mw.typeStr, "verify", verifyConnection, "err", err, "took", time.Since(then))
		}(time.Now())

		mw.logger.Trace("database/Initialize: starting", "type", mw.typeStr)
	}
	return mw.next.Initialize(conf, verifyConnection)
}

func (mw *databaseTracingMiddleware) Close() (err error) {
	if mw.logger.IsTrace() {
		defer func(then time.Time) {
			mw.logger.Trace("database/Close: finished", "type", mw.typeStr, "err", err, "took", time.Since(then))
		}(time.Now())

		mw.logger.Trace("database/Close: starting", "type", mw.typeStr)
	}
	return mw.next.Close()
}

// ---- Metrics Middleware Domain ----

// databaseMetricsMiddleware wraps an implementation of DatabaseTypes and on
// function call logs metrics about this instance.
type databaseMetricsMiddleware struct {
	next DatabaseType

	typeStr string
}

func (mw *databaseMetricsMiddleware) Type() string {
	return mw.next.Type()
}

func (mw *databaseMetricsMiddleware) CreateUser(statements Statements, usernamePrefix string, expiration time.Time) (username string, password string, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "CreateUser"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "CreateUser"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "CreateUser", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "CreateUser", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "CreateUser"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "CreateUser"}, 1)
	return mw.next.CreateUser(statements, usernamePrefix, expiration)
}

func (mw *databaseMetricsMiddleware) RenewUser(statements Statements, username string, expiration time.Time) (err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "RenewUser"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "RenewUser"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "RenewUser", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "RenewUser", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "RenewUser"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "RenewUser"}, 1)
	return mw.next.RenewUser(statements, username, expiration)
}

func (mw *databaseMetricsMiddleware) RevokeUser(statements Statements, username string) (err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "RevokeUser"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "RevokeUser"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "RevokeUser", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "RevokeUser", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "RevokeUser"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "RevokeUser"}, 1)
	return mw.next.RevokeUser(statements, username)
}

func (mw *databaseMetricsMiddleware) Initialize(conf map[string]interface{}, verifyConnection bool) (err error) {
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
	return mw.next.Initialize(conf, verifyConnection)
}

func (mw *databaseMetricsMiddleware) Close() (err error) {
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
