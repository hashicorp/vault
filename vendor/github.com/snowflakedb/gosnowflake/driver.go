// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"os"
	"strings"
	"sync"
)

var paramsMutex *sync.Mutex

// SnowflakeDriver is a context of Go Driver
type SnowflakeDriver struct{}

// Open creates a new connection.
func (d SnowflakeDriver) Open(dsn string) (driver.Conn, error) {
	var cfg *Config
	var err error
	logger.Info("Open")
	ctx := context.Background()
	if dsn == "autoConfig" {
		cfg, err = loadConnectionConfig()
	} else {
		cfg, err = ParseDSN(dsn)
	}
	if err != nil {
		return nil, err
	}
	return d.OpenWithConfig(ctx, *cfg)
}

// OpenWithConfig creates a new connection with the given Config.
func (d SnowflakeDriver) OpenWithConfig(ctx context.Context, config Config) (driver.Conn, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	if config.Tracing != "" {
		logger.SetLogLevel(config.Tracing)
	}
	logger.WithContext(ctx).Info("OpenWithConfig")
	sc, err := buildSnowflakeConn(ctx, config)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(strings.ToLower(config.Host), cnDomain) {
		logger.WithContext(ctx).Info("Connecting to CHINA Snowflake domain")
	} else {
		logger.WithContext(ctx).Info("Connecting to GLOBAL Snowflake domain")
	}

	if err = authenticateWithConfig(sc); err != nil {
		return nil, err
	}
	sc.connectionTelemetry(&config)

	sc.startHeartBeat()
	sc.internal = &httpClient{sr: sc.rest}
	return sc, nil
}

func runningOnGithubAction() bool {
	return os.Getenv("GITHUB_ACTIONS") != ""
}

// GOSNOWFLAKE_SKIP_REGISTERATION is an environment variable which can be set client side to
// bypass dbSql driver registration. This should not be used if sql.Open() is used as the method
// to connect to the server, as sql.Open will require registration so it can map the driver name
// to the driver type, which in this case is "snowflake" and SnowflakeDriver{}. If you wish to call
// into multiple versions of the driver from one client, this is needed because calling register
// twice with the same name on init will cause the driver to panic.
func skipRegisteration() bool {
	return os.Getenv("GOSNOWFLAKE_SKIP_REGISTERATION") != ""
}

var logger = CreateDefaultLogger()

func init() {
	if !skipRegisteration() {
		sql.Register("snowflake", &SnowflakeDriver{})
	}
	logger.SetLogLevel("error")
	if runningOnGithubAction() {
		logger.SetLogLevel("fatal")
	}
	paramsMutex = &sync.Mutex{}
}
