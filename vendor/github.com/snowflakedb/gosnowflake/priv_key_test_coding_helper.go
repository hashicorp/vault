// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"crypto/rsa"
	"crypto/x509"
)

func parsePKCS8PrivateKey(block []byte) (*rsa.PrivateKey, error) {
	privKey, err := x509.ParsePKCS8PrivateKey(block)
	if err != nil {
		return nil, &SnowflakeError{
			Number:  ErrCodePrivateKeyParseError,
			Message: "Error decoding private key using PKCS8.",
		}
	}
	return privKey.(*rsa.PrivateKey), nil
}

func marshalPKCS8PrivateKey(key *rsa.PrivateKey) ([]byte, error) {
	keyInBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, &SnowflakeError{
			Number:  ErrCodePrivateKeyParseError,
			Message: "Error encoding private key using PKCS8."}
	}
	return keyInBytes, nil

}
