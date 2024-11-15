// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build gssapi && !windows && !linux && !darwin
// +build gssapi,!windows,!linux,!darwin

package auth

import (
	"fmt"
	"net/http"
	"runtime"
)

// GSSAPI is the mechanism name for GSSAPI.
const GSSAPI = "GSSAPI"

func newGSSAPIAuthenticator(*Cred, *http.Client) (Authenticator, error) {
	return nil, newAuthError(fmt.Sprintf("GSSAPI is not supported on %s", runtime.GOOS), nil)
}
