// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !ui

package http

import (
	"net/http"
	"os"
)

func init() {
	// set uiBuiltIn to false to indicate the ui is not built in. See
	// http/handler.go
	uiBuiltIn = false
}

// assetFS serves the UI from the given directory or defaults to a stub when
// Vault is built without a UI.
func assetFS(uiDir string) http.FileSystem {
	if uiDir != "" {
		return http.FS(os.DirFS(uiDir))
	}

	return nil
}
