// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"fmt"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
)

var configurableAuthTypes = []string{
	AuthTypeGCPIAM,
}

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

func (c *SQLConnectionProducer) registerDrivers(driverName string, credentials string) (func() error, error) {
	typ, err := c.getCloudSQLDriverType()
	if err != nil {
		return nil, err
	}

	opts, err := GetCloudSQLAuthOptions(credentials)
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
func GetCloudSQLAuthOptions(credentials string) ([]cloudsqlconn.Option, error) {
	opts := []cloudsqlconn.Option{cloudsqlconn.WithIAMAuthN()}

	if credentials != "" {
		opts = append(opts, cloudsqlconn.WithCredentialsJSON([]byte(credentials)))
	}

	return opts, nil
}

func ValidateAuthType(authType string) bool {
	var valid bool
	for _, typ := range configurableAuthTypes {
		if authType == typ {
			valid = true
			break
		}
	}

	return valid
}
