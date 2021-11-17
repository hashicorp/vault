// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/description"
)

func newDefaultAuthenticator(cred *Cred) (Authenticator, error) {
	scram, err := newScramSHA256Authenticator(cred)
	if err != nil {
		return nil, newAuthError("failed to create internal authenticator", err)
	}
	speculative, ok := scram.(SpeculativeAuthenticator)
	if !ok {
		typeErr := fmt.Errorf("expected SCRAM authenticator to be SpeculativeAuthenticator but got %T", scram)
		return nil, newAuthError("failed to create internal authenticator", typeErr)
	}

	return &DefaultAuthenticator{
		Cred:                     cred,
		speculativeAuthenticator: speculative,
	}, nil
}

// DefaultAuthenticator uses SCRAM-SHA-1 or MONGODB-CR depending
// on the server version.
type DefaultAuthenticator struct {
	Cred *Cred

	// The authenticator to use for speculative authentication. Because the correct auth mechanism is unknown when doing
	// the initial isMaster, SCRAM-SHA-256 is used for the speculative attempt.
	speculativeAuthenticator SpeculativeAuthenticator
}

var _ SpeculativeAuthenticator = (*DefaultAuthenticator)(nil)

// CreateSpeculativeConversation creates a speculative conversation for SCRAM authentication.
func (a *DefaultAuthenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) {
	return a.speculativeAuthenticator.CreateSpeculativeConversation()
}

// Auth authenticates the connection.
func (a *DefaultAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	var actual Authenticator
	var err error

	switch chooseAuthMechanism(cfg) {
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

	return actual.Auth(ctx, cfg)
}

// If a server provides a list of supported mechanisms, we choose
// SCRAM-SHA-256 if it exists or else MUST use SCRAM-SHA-1.
// Otherwise, we decide based on what is supported.
func chooseAuthMechanism(cfg *Config) string {
	if saslSupportedMechs := cfg.HandshakeInfo.SaslSupportedMechs; saslSupportedMechs != nil {
		for _, v := range saslSupportedMechs {
			if v == SCRAMSHA256 {
				return v
			}
		}
		return SCRAMSHA1
	}

	if err := scramSHA1Supported(cfg.HandshakeInfo.Description.WireVersion); err == nil {
		return SCRAMSHA1
	}

	return MONGODBCR
}

// scramSHA1Supported returns an error if the given server version does not support scram-sha-1.
func scramSHA1Supported(wireVersion *description.VersionRange) error {
	if wireVersion != nil && wireVersion.Max < 3 {
		return fmt.Errorf("SCRAM-SHA-1 is only supported for servers 3.0 or newer")
	}

	return nil
}
