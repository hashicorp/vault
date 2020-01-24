// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
)

func newDefaultAuthenticator(cred *Cred) (Authenticator, error) {
	return &DefaultAuthenticator{
		Cred: cred,
	}, nil
}

// DefaultAuthenticator uses SCRAM-SHA-1 or MONGODB-CR depending
// on the server version.
type DefaultAuthenticator struct {
	Cred *Cred
}

// Auth authenticates the connection.
func (a *DefaultAuthenticator) Auth(ctx context.Context, desc description.Server, conn driver.Connection) error {
	var actual Authenticator
	var err error

	switch chooseAuthMechanism(desc) {
	case SCRAMSHA256:
		actual, err = newScramSHA256Authenticator(a.Cred)
	case SCRAMSHA1:
		actual, err = newScramSHA1Authenticator(a.Cred)
	default:
		actual, err = newMongoDBCRAuthenticator(a.Cred)
	}

	if err != nil {
		return newAuthError("error creating authenticator", err)
	}

	return actual.Auth(ctx, desc, conn)
}

// If a server provides a list of supported mechanisms, we choose
// SCRAM-SHA-256 if it exists or else MUST use SCRAM-SHA-1.
// Otherwise, we decide based on what is supported.
func chooseAuthMechanism(desc description.Server) string {
	if desc.SaslSupportedMechs != nil {
		for _, v := range desc.SaslSupportedMechs {
			if v == SCRAMSHA256 {
				return v
			}
		}
		return SCRAMSHA1
	}

	if err := description.ScramSHA1Supported(desc.WireVersion); err == nil {
		return SCRAMSHA1
	}

	return MONGODBCR
}
