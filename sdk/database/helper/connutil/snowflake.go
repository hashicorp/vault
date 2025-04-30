// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/snowflakedb/gosnowflake"
)

// Open the DB connection to Snowflake or return an error.
func openSnowflake(connectionURL, username, providedPrivateKey string) (*sql.DB, error) {
	// Parse thee connection_url for required fields. Should be of
	// the form <account_name>.snowflakecomputing.com/<db_name>
	accountName, dbName, err := parseSnowflakeFieldsFromURL(connectionURL)
	if err != nil {
		return nil, err
	}

	privateKey, err := getPrivateKey(providedPrivateKey)
	if err != nil {
		return nil, err
	}

	snowflakeConfig := &gosnowflake.Config{
		Account:       accountName,
		Database:      dbName,
		User:          username,
		Authenticator: gosnowflake.AuthTypeJwt,
		PrivateKey:    privateKey,
	}
	connector := gosnowflake.NewConnector(gosnowflake.SnowflakeDriver{}, *snowflakeConfig)

	return sql.OpenDB(connector), nil
}

// Parse thee connection_url for required fields.
func parseSnowflakeFieldsFromURL(connectionURL string) (string, string, error) {
	nameDelim := strings.Index(connectionURL, ".snowflakecomputing")
	domainDelim := strings.Index(connectionURL, ".com/")
	if nameDelim == -1 || domainDelim == -1 {
		return "", "", fmt.Errorf("invalid connection URL, missing snowflakecomputing.com")
	}

	accountName := connectionURL[:nameDelim]
	dbName := connectionURL[domainDelim:]

	return accountName, dbName, nil
}

// Open and decode the private key file
func getPrivateKey(providedPrivateKey string) (*rsa.PrivateKey, error) {
	var block *pem.Block

	// If the provided data was the key itself, use it directly.
	if strings.HasPrefix(providedPrivateKey, "-----BEGIN PRIVATE KEY-----") {
		block, _ = pem.Decode([]byte(providedPrivateKey))
	} else {
		keyFile, err := os.ReadFile(providedPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}

		block, _ = pem.Decode(keyFile)
	}

	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode the the private key value")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key to PKCS8: %w", err)
	}

	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key was parsed into an unexpected type")
	}

	return privateKey, nil
}
