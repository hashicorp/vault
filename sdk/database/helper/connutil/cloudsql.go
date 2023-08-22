// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"fmt"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
	"cloud.google.com/go/cloudsqlconn/sqlserver/mssql"
)

func (c *SQLConnectionProducer) getCloudSQLDriverName() (string, error) {
	var driverName string
	switch c.Type {
	case dbTypeMSSQL:
		driverName = cloudSQLMSSQL
	case dbTypePostgres:
		driverName = cloudSQLPostgres
	default:
		return "", fmt.Errorf("unrecognized DB type: %s", c.Type)
	}

	return driverName, nil
}

func (c *SQLConnectionProducer) registerDrivers(driverName string, credentials string) (func() error, error) {
	typ, err := c.getCloudSQLDriverName()
	if err != nil {
		return nil, err
	}

	opts, err := GetCloudSQLAuthOptions(credentials)
	if err != nil {
		return nil, err
	}

	switch typ {
	case cloudSQLMSSQL:
		return mssql.RegisterDriver(driverName, opts...)
	case cloudSQLPostgres:
		return pgxv4.RegisterDriver(driverName, opts...)
	}

	return nil, fmt.Errorf("unrecognized cloudsql type encountered: %s", typ)
}

// GetCloudSQLAuthOptions takes a credentials (file) or a credentialsJSON (the actual data) and returns
// a set of GCP CloudSQL options - always WithIAMAUthN, and then the appropriate file/JSON option.

func GetCloudSQLAuthOptions(credentials string) ([]cloudsqlconn.Option, error) {
	opts := []cloudsqlconn.Option{cloudsqlconn.WithIAMAuthN()}

	if credentials != "" {
		opts = append(opts, cloudsqlconn.WithCredentialsJSON([]byte(credentials)))
	}

	return opts, nil
}

func cacheDrivers(driverName string, f cloudSQLCleanup) {
	driversMu.Lock()
	defer driversMu.Unlock()

	drivers[driverName] = f
}

func cachePop(typ string) cloudSQLCleanup {
	driversMu.Lock()
	defer driversMu.Unlock()

	var cleanup cloudSQLCleanup
	if f, ok := drivers[typ]; ok {
		cleanup = f
		delete(drivers, typ)
	}
	return cleanup
}
