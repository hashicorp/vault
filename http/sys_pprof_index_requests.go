// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/hashicorp/vault/vault"
)

// handlePprofNamedIndexRequest returns an http.Handler that serves the pprof
// profile corresponding to "name".  The list of valid profiles includes:
// "cmdline", "profile", "symbol", "trace", "allocs", "block", "goroutine",
// "heap", "mutex", "threadcreate".
//
// Where applicable, we apply runtime settings to ensure the profiles produce
// output.
func handlePprofNamedIndexRequest(core *vault.Core, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		core.Logger().Debug("handling pprof", "name", name)

		// Blocking profiles need to be provided a sampling rate.
		switch name {
		case "block":
			runtime.SetBlockProfileRate(1)
			defer runtime.SetBlockProfileRate(0)
		case "mutex":
			runtime.SetMutexProfileFraction(1)
			defer runtime.SetMutexProfileFraction(0)
		}

		pprof.Handler(name).ServeHTTP(w, r)
	})
}
