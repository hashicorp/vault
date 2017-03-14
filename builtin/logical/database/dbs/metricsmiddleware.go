package dbs

import (
	"time"

	metrics "github.com/armon/go-metrics"
)

type databaseMetricsMiddleware struct {
	next DatabaseType

	typeStr string
}

func (mw *databaseMetricsMiddleware) Type() string {
	return mw.next.Type()
}

func (mw *databaseMetricsMiddleware) CreateUser(statements Statements, username, password, expiration string) (err error) {
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
	return mw.next.CreateUser(statements, username, password, expiration)
}

func (mw *databaseMetricsMiddleware) RenewUser(statements Statements, username, expiration string) (err error) {
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

func (mw *databaseMetricsMiddleware) Initialize(conf map[string]interface{}) (err error) {
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
	return mw.next.Initialize(conf)
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

func (mw *databaseMetricsMiddleware) GenerateUsername(displayName string) (_ string, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "GenerateUsername"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "GenerateUsername"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "GenerateUsername", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "GenerateUsername", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "GenerateUsername"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "GenerateUsername"}, 1)
	return mw.next.GenerateUsername(displayName)
}

func (mw *databaseMetricsMiddleware) GeneratePassword() (_ string, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "GeneratePassword"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "GeneratePassword"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "GeneratePassword", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "GeneratePassword", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "GeneratePassword"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "GeneratePassword"}, 1)
	return mw.next.GeneratePassword()
}

func (mw *databaseMetricsMiddleware) GenerateExpiration(duration time.Duration) (_ string, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "GenerateExpiration"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "GenerateExpiration"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "GenerateExpiration", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "GenerateExpiration", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "GenerateExpiration"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "GenerateExpiration"}, 1)
	return mw.next.GenerateExpiration(duration)
}
