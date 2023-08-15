// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"fmt"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
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

	fmt.Printf("returning driver type: %s for dbType: %s", driverName, c.Type)
	return driverName, nil
}

func (c *SQLConnectionProducer) registerDrivers(filename, credentials interface{}) (func() error, error) {
	typ, err := c.getCloudSQLDriverName()
	if err != nil {
		return nil, err
	}

	if isDriverRegistered(typ) {
		// driver already registered
		fmt.Printf("drivers have already been registered, returning\n")
		return nil, nil
	}

	opts, err := GetCloudSQLAuthOptions(filename, credentials)
	if err != nil {
		return nil, err
	}

	switch typ {
	case cloudSQLMSSQL:
		// return registerDriverMSSQL(opts)
	case cloudSQLPostgres:
		return registerDriverPostgres(opts...)
	}

	return nil, fmt.Errorf("unrecognized cloudsql type encountered: %s", typ)
}

func registerDriverPostgres(opts ...cloudsqlconn.Option) (func() error, error) {
	return pgxv4.RegisterDriver(cloudSQLPostgres, opts...)
}

//func registerDriverMSSQL(opts cloudsqlconn.Option) (func() error, error) {
//	return mssql.RegisterDriver(cloudSQLMSSQL, opts)
//}

func GetCloudSQLAuthOptions(filename, credentials interface{}) ([]cloudsqlconn.Option, error) {
	opts := []cloudsqlconn.Option{cloudsqlconn.WithIAMAuthN()}

	if filename != nil {
		v, ok := filename.(string)
		if !ok {
			return nil, fmt.Errorf("error converting file name to string")
		}

		fmt.Printf("registering driver with credential file\n")
		opts = append(opts, cloudsqlconn.WithCredentialsFile(v))
	}

	if credentials != nil {
		fmt.Printf("registering driver with credential json\n")
		switch v := credentials.(type) {
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
