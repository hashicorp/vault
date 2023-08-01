// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"fmt"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
)

func (c *SQLConnectionProducer) getCloudSQLDBType() (string, error) {
	var dbType string
	switch c.Type {
	case "mssql":
		dbType = cloudSQLMSSQL
	case dbTypePostgres:
		dbType = cloudSQLPostgres
	default:
		return "", fmt.Errorf("unrecognized DB type: %s", c.Type)
	}

	return dbType, nil
}

func (c *SQLConnectionProducer) registerDrivers(filename, credentials interface{}) (func() error, error) {
	typ, err := c.getCloudSQLDBType()
	if err != nil {
		return nil, err
	}

	if cacheGet(typ) != nil {
		// drivers have already been registered
		// return
		fmt.Printf("drivers have already been registered, returning\n")
		return nil, nil
	}

	opts, err := getAuthOptions(filename, credentials)
	if err != nil {
		return nil, err
	}
	// @TODO add support for other drivers
	switch typ {
	case cloudSQLMSSQL:
		// return mssql.RegisterDriver(cloudSQLMSSQL, cloudsqlconn.WithCredentialsFile("key.json"))
	case cloudSQLPostgres:
		return registerDriverPostgres(opts)
	}

	return nil, fmt.Errorf("unrecognized cloudsql type encountered: %s", typ)
}

func registerDriverPostgres(opts cloudsqlconn.Option) (func() error, error) {
	return pgxv4.RegisterDriver(cloudSQLPostgres, opts)
}

func getAuthOptions(filename, credentials interface{}) (cloudsqlconn.Option, error) {
	if filename != nil {
		v, ok := filename.(string)
		if !ok {
			return nil, fmt.Errorf("error converting file name to string")
		}

		fmt.Printf("registering driver with credential file\n")
		return cloudsqlconn.WithCredentialsFile(v), nil
	}

	if credentials != nil {
		v, ok := credentials.([]byte)
		if !ok {
			return nil, fmt.Errorf("error converting JSON data to bytes")
		}

		fmt.Printf("registering driver with credential json\n")
		return cloudsqlconn.WithCredentialsJSON(v), nil

	}

	return cloudsqlconn.WithIAMAuthN(), nil
}

func cacheCleanup(typ string, f func() error) {
	basicCleanupCache[typ] = f
}

func cacheGet(typ string) func() error {
	return basicCleanupCache[typ]
}
