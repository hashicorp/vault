// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//+build gssapi
//+build windows linux darwin

package auth

import (
	"context"
	"fmt"
	"net"

	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth/internal/gssapi"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
)

// GSSAPI is the mechanism name for GSSAPI.
const GSSAPI = "GSSAPI"

func newGSSAPIAuthenticator(cred *Cred) (Authenticator, error) {
	if cred.Source != "" && cred.Source != "$external" {
		return nil, newAuthError("GSSAPI source must be empty or $external", nil)
	}

	return &GSSAPIAuthenticator{
		Username:    cred.Username,
		Password:    cred.Password,
		PasswordSet: cred.PasswordSet,
		Props:       cred.Props,
	}, nil
}

// GSSAPIAuthenticator uses the GSSAPI algorithm over SASL to authenticate a connection.
type GSSAPIAuthenticator struct {
	Username    string
	Password    string
	PasswordSet bool
	Props       map[string]string
}

// Auth authenticates the connection.
func (a *GSSAPIAuthenticator) Auth(ctx context.Context, desc description.Server, conn driver.Connection) error {
	target := desc.Addr.String()
	hostname, _, err := net.SplitHostPort(target)
	if err != nil {
		return newAuthError(fmt.Sprintf("invalid endpoint (%s) specified: %s", target, err), nil)
	}

	client, err := gssapi.New(hostname, a.Username, a.Password, a.PasswordSet, a.Props)

	if err != nil {
		return newAuthError("error creating gssapi", err)
	}
	return ConductSaslConversation(ctx, conn, "$external", client)
}
