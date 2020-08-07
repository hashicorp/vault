// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/address"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

// AuthenticatorFactory constructs an authenticator.
type AuthenticatorFactory func(cred *Cred) (Authenticator, error)

var authFactories = make(map[string]AuthenticatorFactory)

func init() {
	RegisterAuthenticatorFactory("", newDefaultAuthenticator)
	RegisterAuthenticatorFactory(SCRAMSHA1, newScramSHA1Authenticator)
	RegisterAuthenticatorFactory(SCRAMSHA256, newScramSHA256Authenticator)
	RegisterAuthenticatorFactory(MONGODBCR, newMongoDBCRAuthenticator)
	RegisterAuthenticatorFactory(PLAIN, newPlainAuthenticator)
	RegisterAuthenticatorFactory(GSSAPI, newGSSAPIAuthenticator)
	RegisterAuthenticatorFactory(MongoDBX509, newMongoDBX509Authenticator)
}

// CreateAuthenticator creates an authenticator.
func CreateAuthenticator(name string, cred *Cred) (Authenticator, error) {
	if f, ok := authFactories[name]; ok {
		return f(cred)
	}

	return nil, newAuthError(fmt.Sprintf("unknown authenticator: %s", name), nil)
}

// RegisterAuthenticatorFactory registers the authenticator factory.
func RegisterAuthenticatorFactory(name string, factory AuthenticatorFactory) {
	authFactories[name] = factory
}

// HandshakeOptions packages options that can be passed to the Handshaker()
// function.  DBUser is optional but must be of the form <dbname.username>;
// if non-empty, then the connection will do SASL mechanism negotiation.
type HandshakeOptions struct {
	AppName               string
	Authenticator         Authenticator
	Compressors           []string
	DBUser                string
	PerformAuthentication func(description.Server) bool
}

type authHandshaker struct {
	wrapped driver.Handshaker
	options *HandshakeOptions
}

// GetDescription performs an isMaster to retrieve the initial description for conn.
func (ah *authHandshaker) GetDescription(ctx context.Context, addr address.Address, conn driver.Connection) (description.Server, error) {
	if ah.wrapped != nil {
		return ah.wrapped.GetDescription(ctx, addr, conn)
	}

	desc, err := operation.NewIsMaster().
		AppName(ah.options.AppName).
		Compressors(ah.options.Compressors).
		SASLSupportedMechs(ah.options.DBUser).
		GetDescription(ctx, addr, conn)
	if err != nil {
		return description.Server{}, newAuthError("handshake failure", err)
	}
	return desc, nil
}

// FinishHandshake performs authentication for conn if necessary.
func (ah *authHandshaker) FinishHandshake(ctx context.Context, conn driver.Connection) error {
	performAuth := ah.options.PerformAuthentication
	if performAuth == nil {
		performAuth = func(serv description.Server) bool {
			// Authentication is possible against all server types except arbiters
			return serv.Kind != description.RSArbiter
		}
	}
	desc := conn.Description()
	if performAuth(desc) && ah.options.Authenticator != nil {
		err := ah.options.Authenticator.Auth(ctx, desc, conn)
		if err != nil {
			return newAuthError("auth error", err)
		}
	}

	if ah.wrapped == nil {
		return nil
	}
	return ah.wrapped.FinishHandshake(ctx, conn)
}

// Handshaker creates a connection handshaker for the given authenticator.
func Handshaker(h driver.Handshaker, options *HandshakeOptions) driver.Handshaker {
	return &authHandshaker{
		wrapped: h,
		options: options,
	}
}

// Authenticator handles authenticating a connection.
type Authenticator interface {
	// Auth authenticates the connection.
	Auth(context.Context, description.Server, driver.Connection) error
}

func newAuthError(msg string, inner error) error {
	return &Error{
		message: msg,
		inner:   inner,
	}
}

func newError(err error, mech string) error {
	return &Error{
		message: fmt.Sprintf("unable to authenticate using mechanism \"%s\"", mech),
		inner:   err,
	}
}

// Error is an error that occurred during authentication.
type Error struct {
	message string
	inner   error
}

func (e *Error) Error() string {
	if e.inner == nil {
		return e.message
	}
	return fmt.Sprintf("%s: %s", e.message, e.inner)
}

// Inner returns the wrapped error.
func (e *Error) Inner() error {
	return e.inner
}

// Message returns the message.
func (e *Error) Message() string {
	return e.message
}
