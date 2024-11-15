// Copyright (c) 2020-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"database/sql/driver"
)

// InternalSnowflakeDriver is the interface for an internal Snowflake driver
type InternalSnowflakeDriver interface {
	Open(dsn string) (driver.Conn, error)
	OpenWithConfig(ctx context.Context, config Config) (driver.Conn, error)
}

// Connector creates Driver with the specified Config
type Connector struct {
	driver InternalSnowflakeDriver
	cfg    Config
}

// NewConnector creates a new connector with the given SnowflakeDriver and Config.
func NewConnector(driver InternalSnowflakeDriver, config Config) Connector {
	return Connector{driver, config}
}

// Connect creates a new connection.
func (t Connector) Connect(ctx context.Context) (driver.Conn, error) {
	cfg := t.cfg
	err := fillMissingConfigParameters(&cfg)
	if err != nil {
		return nil, err
	}
	return t.driver.OpenWithConfig(ctx, cfg)
}

// Driver creates a new driver.
func (t Connector) Driver() driver.Driver {
	return t.driver
}
