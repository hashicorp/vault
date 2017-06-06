package dbplugin

import (
	"time"

	metrics "github.com/armon/go-metrics"
	log "github.com/mgutz/logxi/v1"
)

// ---- Tracing Middleware Domain ----

// databaseTracingMiddleware wraps a implementation of Database and executes
// trace logging on function call.
type databaseTracingMiddleware struct {
	next   Database
	logger log.Logger

	typeStr string
}

func (mw *databaseTracingMiddleware) Type() (string, error) {
	return mw.next.Type()
}

func (mw *databaseTracingMiddleware) CreateUser(statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error) {
	defer func(then time.Time) {
		mw.logger.Trace("database", "operation", "CreateUser", "status", "finished", "type", mw.typeStr, "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("database", "operation", "CreateUser", "status", "started", "type", mw.typeStr)
	return mw.next.CreateUser(statements, usernameConfig, expiration)
}

func (mw *databaseTracingMiddleware) RenewUser(statements Statements, username string, expiration time.Time) (err error) {
	defer func(then time.Time) {
		mw.logger.Trace("database", "operation", "RenewUser", "status", "finished", "type", mw.typeStr, "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("database", "operation", "RenewUser", "status", "started", mw.typeStr)
	return mw.next.RenewUser(statements, username, expiration)
}

func (mw *databaseTracingMiddleware) RevokeUser(statements Statements, username string) (err error) {
	defer func(then time.Time) {
		mw.logger.Trace("database", "operation", "RevokeUser", "status", "finished", "type", mw.typeStr, "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("database", "operation", "RevokeUser", "status", "started", "type", mw.typeStr)
	return mw.next.RevokeUser(statements, username)
}

func (mw *databaseTracingMiddleware) Initialize(conf map[string]interface{}, verifyConnection bool) (err error) {
	defer func(then time.Time) {
		mw.logger.Trace("database", "operation", "Initialize", "status", "finished", "type", mw.typeStr, "verify", verifyConnection, "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("database", "operation", "Initialize", "status", "started", "type", mw.typeStr)
	return mw.next.Initialize(conf, verifyConnection)
}

func (mw *databaseTracingMiddleware) Close() (err error) {
	defer func(then time.Time) {
		mw.logger.Trace("database", "operation", "Close", "status", "finished", "type", mw.typeStr, "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("database", "operation", "Close", "status", "started", "type", mw.typeStr)
	return mw.next.Close()
}

// ---- Metrics Middleware Domain ----

// databaseMetricsMiddleware wraps an implementation of Databases and on
// function call logs metrics about this instance.
type databaseMetricsMiddleware struct {
	next Database

	typeStr string
}

func (mw *databaseMetricsMiddleware) Type() (string, error) {
	return mw.next.Type()
}

func (mw *databaseMetricsMiddleware) CreateUser(statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error) {
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
	return mw.next.CreateUser(statements, usernameConfig, expiration)
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
