//go:build ui

package http

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
)

// content is our static web server content.
//
//go:embed web_ui/*
var content embed.FS

// assetFS is a http Filesystem that serves the generated web UI from the
// "ember-dist" make step or the given directory.
func assetFS(uiDir string) http.FileSystem {
	var f fs.FS
	if uiDir == "" {
		// sub out to web_ui, where the generated content lives
		var err error
		f, err = fs.Sub(content, "web_ui")
		if err != nil {
			panic(err)
		}
	} else {
		f = os.DirFS(uiDir)
	}
	return http.FS(f)
}
