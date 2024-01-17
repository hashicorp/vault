// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package http

import (
	"net/http"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func adjustRequest(c *vault.Core, listener *configutil.Listener, r *http.Request) (*http.Request, int, error) {
	return r, 0, nil
}

func handleEntPaths(nsPath string, core *vault.Core, r *http.Request) http.Handler {
	return nil
}
