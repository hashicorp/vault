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

func (c *SQLConnectionProducer) registerDrivers(driverName string, credentials, credentialsJSON interface{}) (func() error, error) {
	typ, err := c.getCloudSQLDriverName()
	if err != nil {
		return nil, err
	}

	if isDriverRegistered(typ) {
		// driver already registered
		fmt.Printf("drivers have already been registered, returning\n")
		return nil, nil
	}

	opts, err := GetCloudSQLAuthOptions(credentials, credentialsJSON)
	if err != nil {
		return nil, err
	}

	switch typ {
	case cloudSQLMSSQL:
		return pgxv4.RegisterDriver(driverName, opts...)
	case cloudSQLPostgres:
		return mssql.RegisterDriver(driverName, opts...)
	}

	return nil, fmt.Errorf("unrecognized cloudsql type encountered: %s", typ)
}

// GetCloudSQLAuthOptions takes a credentials (file) or a credentialsJSON (the actual data) and returns
// a set of GCP CloudSQL options - always WithIAMAUthN, and then the appropriate file/JSON option.
func GetCloudSQLAuthOptions(credentials, credentialsJSON interface{}) ([]cloudsqlconn.Option, error) {
	opts := []cloudsqlconn.Option{cloudsqlconn.WithIAMAuthN()}
	if credentials != nil {
		v, ok := credentials.(string)
		if !ok {
			return nil, fmt.Errorf("error converting file name to string")
		}

		fmt.Printf("registering driver with credential file\n")
		opts = append(opts, cloudsqlconn.WithCredentialsFile(v))
	}

	if credentialsJSON != nil {
		fmt.Printf("registering driver with credential json\n")
		switch v := credentialsJSON.(type) {
		case string:
			opts = append(opts, cloudsqlconn.WithCredentialsJSON([]byte(v)))
		case []byte:
			opts = append(opts, cloudsqlconn.WithCredentialsJSON(v))
		default:
			return nil, fmt.Errorf("error converting credentials of type %T to []byte", credentials)
		}
	}

	return opts, nil
}

func cacheDrivers(typ string, f cloudSQLCleanup) {
	driversMu.Lock()
	defer driversMu.Unlock()

	drivers[typ] = f
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

func isDriverRegistered(typ string) bool {
	driversMu.Lock()
	defer driversMu.Unlock()

	return drivers[typ] != nil
}
