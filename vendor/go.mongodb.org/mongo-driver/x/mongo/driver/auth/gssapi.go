// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build gssapi && (windows || linux || darwin)
// +build gssapi
// +build windows linux darwin

package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth/internal/gssapi"
)

// GSSAPI is the mechanism name for GSSAPI.
const GSSAPI = "GSSAPI"

func newGSSAPIAuthenticator(cred *Cred, _ *http.Client) (Authenticator, error) {
	if cred.Source != "" && cred.Source != sourceExternal {
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
func (a *GSSAPIAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	target := cfg.Description.Addr.String()
	hostname, _, err := net.SplitHostPort(target)
	if err != nil {
		return newAuthError(fmt.Sprintf("invalid endpoint (%s) specified: %s", target, err), nil)
	}

	client, err := gssapi.New(hostname, a.Username, a.Password, a.PasswordSet, a.Props)

	if err != nil {
		return newAuthError("error creating gssapi", err)
	}
	return ConductSaslConversation(ctx, cfg, sourceExternal, client)
}

// Reauth reauthenticates the connection.
func (a *GSSAPIAuthenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("GSSAPI does not support reauthentication", nil)
}
