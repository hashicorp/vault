package dbplugin

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"

	metrics "github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
)

// ---- Tracing Middleware Domain ----

// databaseTracingMiddleware wraps a implementation of Database and executes
// trace logging on function call.
type databaseTracingMiddleware struct {
	next   Database
	logger log.Logger
}

func (mw *databaseTracingMiddleware) Type() (string, error) {
	return mw.next.Type()
}

func (mw *databaseTracingMiddleware) CreateUser(ctx context.Context, statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error) {
	defer func(then time.Time) {
		mw.logger.Trace("create user", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("create user", "status", "started")
	return mw.next.CreateUser(ctx, statements, usernameConfig, expiration)
}

func (mw *databaseTracingMiddleware) RenewUser(ctx context.Context, statements Statements, username string, expiration time.Time) (err error) {
	defer func(then time.Time) {
		mw.logger.Trace("renew user", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("renew user", "status", "started")
	return mw.next.RenewUser(ctx, statements, username, expiration)
}

func (mw *databaseTracingMiddleware) RevokeUser(ctx context.Context, statements Statements, username string) (err error) {
	defer func(then time.Time) {
		mw.logger.Trace("revoke user", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("revoke user", "status", "started")
	return mw.next.RevokeUser(ctx, statements, username)
}

func (mw *databaseTracingMiddleware) RotateRootCredentials(ctx context.Context, statements []string) (conf map[string]interface{}, err error) {
	defer func(then time.Time) {
		mw.logger.Trace("rotate root credentials", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("rotate root credentials", "status", "started")
	return mw.next.RotateRootCredentials(ctx, statements)
}

func (mw *databaseTracingMiddleware) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := mw.Init(ctx, conf, verifyConnection)
	return err
}

func (mw *databaseTracingMiddleware) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (saveConf map[string]interface{}, err error) {
	defer func(then time.Time) {
		mw.logger.Trace("initialize", "status", "finished", "verify", verifyConnection, "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("initialize", "status", "started")
	return mw.next.Init(ctx, conf, verifyConnection)
}

func (mw *databaseTracingMiddleware) Close() (err error) {
	defer func(then time.Time) {
		mw.logger.Trace("close", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("close", "status", "started")
	return mw.next.Close()
}

func (mw *databaseTracingMiddleware) GenerateCredentials(ctx context.Context) (password string, err error) {
	defer func(then time.Time) {
		mw.logger.Trace("generate credentials", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("generate credentials", "status", "started")
	return mw.next.GenerateCredentials(ctx)
}

func (mw *databaseTracingMiddleware) SetCredentials(ctx context.Context, statements Statements, staticConfig StaticUserConfig) (username, password string, err error) {
	defer func(then time.Time) {
		mw.logger.Trace("set credentials", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	mw.logger.Trace("set credentials", "status", "started")
	return mw.next.SetCredentials(ctx, statements, staticConfig)
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

func (mw *databaseMetricsMiddleware) CreateUser(ctx context.Context, statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error) {
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
	return mw.next.CreateUser(ctx, statements, usernameConfig, expiration)
}

func (mw *databaseMetricsMiddleware) RenewUser(ctx context.Context, statements Statements, username string, expiration time.Time) (err error) {
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
	return mw.next.RenewUser(ctx, statements, username, expiration)
}

func (mw *databaseMetricsMiddleware) RevokeUser(ctx context.Context, statements Statements, username string) (err error) {
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
	return mw.next.RevokeUser(ctx, statements, username)
}

func (mw *databaseMetricsMiddleware) RotateRootCredentials(ctx context.Context, statements []string) (conf map[string]interface{}, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "RotateRootCredentials"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "RotateRootCredentials"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "RotateRootCredentials", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "RotateRootCredentials", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "RotateRootCredentials"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "RotateRootCredentials"}, 1)
	return mw.next.RotateRootCredentials(ctx, statements)
}

func (mw *databaseMetricsMiddleware) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := mw.Init(ctx, conf, verifyConnection)
	return err
}

func (mw *databaseMetricsMiddleware) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (saveConf map[string]interface{}, err error) {
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
	return mw.next.Init(ctx, conf, verifyConnection)
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

func (mw *databaseMetricsMiddleware) GenerateCredentials(ctx context.Context) (password string, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "GenerateCredentials"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "GenerateCredentials"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "GenerateCredentials", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "GenerateCredentials", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "GenerateCredentials"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "GenerateCredentials"}, 1)
	return mw.next.GenerateCredentials(ctx)
}

func (mw *databaseMetricsMiddleware) SetCredentials(ctx context.Context, statements Statements, staticConfig StaticUserConfig) (username, password string, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"database", "SetCredentials"}, now)
		metrics.MeasureSince([]string{"database", mw.typeStr, "SetCredentials"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"database", "SetCredentials", "error"}, 1)
			metrics.IncrCounter([]string{"database", mw.typeStr, "SetCredentials", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"database", "SetCredentials"}, 1)
	metrics.IncrCounter([]string{"database", mw.typeStr, "SetCredentials"}, 1)
	return mw.next.SetCredentials(ctx, statements, staticConfig)
}

// ---- Error Sanitizer Middleware Domain ----

// DatabaseErrorSanitizerMiddleware wraps an implementation of Databases and
// sanitizes returned error messages
type DatabaseErrorSanitizerMiddleware struct {
	l         sync.RWMutex
	next      Database
	secretsFn func() map[string]interface{}
}

func NewDatabaseErrorSanitizerMiddleware(next Database, secretsFn func() map[string]interface{}) *DatabaseErrorSanitizerMiddleware {
	return &DatabaseErrorSanitizerMiddleware{
		next:      next,
		secretsFn: secretsFn,
	}
}

func (mw *DatabaseErrorSanitizerMiddleware) Type() (string, error) {
	dbType, err := mw.next.Type()
	return dbType, mw.sanitize(err)
}

func (mw *DatabaseErrorSanitizerMiddleware) CreateUser(ctx context.Context, statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error) {
	username, password, err = mw.next.CreateUser(ctx, statements, usernameConfig, expiration)
	return username, password, mw.sanitize(err)
}

func (mw *DatabaseErrorSanitizerMiddleware) RenewUser(ctx context.Context, statements Statements, username string, expiration time.Time) (err error) {
	return mw.sanitize(mw.next.RenewUser(ctx, statements, username, expiration))
}

func (mw *DatabaseErrorSanitizerMiddleware) RevokeUser(ctx context.Context, statements Statements, username string) (err error) {
	return mw.sanitize(mw.next.RevokeUser(ctx, statements, username))
}

func (mw *DatabaseErrorSanitizerMiddleware) RotateRootCredentials(ctx context.Context, statements []string) (conf map[string]interface{}, err error) {
	conf, err = mw.next.RotateRootCredentials(ctx, statements)
	return conf, mw.sanitize(err)
}

func (mw *DatabaseErrorSanitizerMiddleware) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := mw.Init(ctx, conf, verifyConnection)
	return err
}

func (mw *DatabaseErrorSanitizerMiddleware) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (saveConf map[string]interface{}, err error) {
	saveConf, err = mw.next.Init(ctx, conf, verifyConnection)
	return saveConf, mw.sanitize(err)
}

func (mw *DatabaseErrorSanitizerMiddleware) Close() (err error) {
	return mw.sanitize(mw.next.Close())
}

// sanitize
func (mw *DatabaseErrorSanitizerMiddleware) sanitize(err error) error {
	if err == nil {
		return nil
	}
	if errwrap.ContainsType(err, new(url.Error)) {
		return errors.New("unable to parse connection url")
	}
	if mw.secretsFn != nil {
		for k, v := range mw.secretsFn() {
			if k == "" {
				continue
			}
			err = errors.New(strings.Replace(err.Error(), k, v.(string), -1))
		}
	}
	return err
}

func (mw *DatabaseErrorSanitizerMiddleware) GenerateCredentials(ctx context.Context) (password string, err error) {
	password, err = mw.next.GenerateCredentials(ctx)
	return password, mw.sanitize(err)
}

func (mw *DatabaseErrorSanitizerMiddleware) SetCredentials(ctx context.Context, statements Statements, staticConfig StaticUserConfig) (username, password string, err error) {
	username, password, err = mw.next.SetCredentials(ctx, statements, staticConfig)
	return username, password, mw.sanitize(err)
}
