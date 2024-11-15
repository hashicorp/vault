// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build !gssapi
// +build !gssapi

package auth

import "net/http"

// GSSAPI is the mechanism name for GSSAPI.
const GSSAPI = "GSSAPI"

func newGSSAPIAuthenticator(*Cred, *http.Client) (Authenticator, error) {
	return nil, newAuthError("GSSAPI support not enabled during build (-tags gssapi)", nil)
}
