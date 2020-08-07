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

// PLAIN is the mechanism name for PLAIN.
const PLAIN = "PLAIN"

func newPlainAuthenticator(cred *Cred) (Authenticator, error) {
	return &PlainAuthenticator{
		Username: cred.Username,
		Password: cred.Password,
	}, nil
}

// PlainAuthenticator uses the PLAIN algorithm over SASL to authenticate a connection.
type PlainAuthenticator struct {
	Username string
	Password string
}

// Auth authenticates the connection.
func (a *PlainAuthenticator) Auth(ctx context.Context, _ description.Server, conn driver.Connection) error {
	return ConductSaslConversation(ctx, conn, "$external", &plainSaslClient{
		username: a.Username,
		password: a.Password,
	})
}

type plainSaslClient struct {
	username string
	password string
}

func (c *plainSaslClient) Start() (string, []byte, error) {
	b := []byte("\x00" + c.username + "\x00" + c.password)
	return PLAIN, b, nil
}

func (c *plainSaslClient) Next(challenge []byte) ([]byte, error) {
	return nil, newAuthError("unexpected server challenge", nil)
}

func (c *plainSaslClient) Completed() bool {
	return true
}
