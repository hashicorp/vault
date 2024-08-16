// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"fmt"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
)

func (c *SQLConnectionProducer) getCloudSQLDriverType() (string, error) {
	var driverType string
	// using switch case for future extensibility
	switch c.Type {
	case dbTypePostgres:
		driverType = cloudSQLPostgres
	default:
		return "", fmt.Errorf("unsupported DB type for cloud IAM: %s", c.Type)
	}

	return driverType, nil
}

func (c *SQLConnectionProducer) registerDrivers(driverName string, credentials string, usePrivateIP bool) (func() error, error) {
	typ, err := c.getCloudSQLDriverType()
	if err != nil {
		return nil, err
	}

	opts, err := GetCloudSQLAuthOptions(credentials, usePrivateIP)
	if err != nil {
		return nil, err
	}

	// using switch case for future extensibility
	switch typ {
	case cloudSQLPostgres:
		return pgxv4.RegisterDriver(driverName, opts...)
	}

	return nil, fmt.Errorf("unrecognized cloudsql type encountered: %s", typ)
}

// GetCloudSQLAuthOptions takes a credentials JSON and returns
// a set of GCP CloudSQL options - always WithIAMAUthN, and then the appropriate file/JSON option.
func GetCloudSQLAuthOptions(credentials string, usePrivateIP bool) ([]cloudsqlconn.Option, error) {
	opts := []cloudsqlconn.Option{cloudsqlconn.WithIAMAuthN()}

	if credentials != "" {
		opts = append(opts, cloudsqlconn.WithCredentialsJSON([]byte(credentials)))
	}

	if usePrivateIP {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}

	return opts, nil
}
