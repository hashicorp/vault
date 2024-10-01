// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build ui

package http

import (
	"embed"
	"io/fs"
	"net/http"
)

// content is our static web server content.
//
//go:embed web_ui/*
var content embed.FS

// assetFS is a http Filesystem that serves the generated web UI from the
// "ember-dist" make step
func assetFS() http.FileSystem {
	// sub out to web_ui, where the generated content lives
	f, err := fs.Sub(content, "web_ui")
	if err != nil {
		panic(err)
	}
	return http.FS(f)
}
