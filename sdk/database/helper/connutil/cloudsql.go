// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"fmt"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/mysql/mysql"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
)

func (c *SQLConnectionProducer) getCloudSQLDBType() (string, error) {
	var dbType string
	switch c.Type {
	case "mssql":
		dbType = cloudSQLMSSQL
	case dbTypePostgres:
		dbType = cloudSQLPostgres
	case "mysql":
		dbType = cloudSQLMySQL
	default:
		return "", fmt.Errorf("unrecognized DB type: %s", c.Type)
	}

	return dbType, nil
}

func (c *SQLConnectionProducer) registerDrivers() (func() error, error) {
	typ, err := c.getCloudSQLDBType()
	if err != nil {
		return nil, err
	}
	// @TODO add support for other drivers
	switch typ {
	case cloudSQLMSSQL:
		// return mssql.RegisterDriver(cloudSQLMSSQL, cloudsqlconn.WithCredentialsFile("key.json"))
	case cloudSQLPostgres:
		return registerDriverPostgres()
	case cloudSQLMySQL:
		return registerDriverMySQL()
	}

	return nil, fmt.Errorf("unrecognized cloudsql type encountered: %s", typ)
}

// @TODO add support for credentials file
func registerDriverPostgres() (func() error, error) {
	fmt.Printf("registering driver for %s\n", cloudSQLPostgres)
	return pgxv4.RegisterDriver(cloudSQLPostgres, cloudsqlconn.WithIAMAuthN())
}

func registerDriverMySQL() (func() error, error) {
	fmt.Printf("registering driver for %s\n", cloudSQLMySQL)
	return mysql.RegisterDriver(cloudSQLMySQL, cloudsqlconn.WithIAMAuthN())
}
