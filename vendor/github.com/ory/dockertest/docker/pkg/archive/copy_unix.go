// +build !windows

package archive // import "github.com/ory/dockertest/docker/pkg/archive"

import (
	"path/filepath"
)

func normalizePath(path string) string {
	return filepath.ToSlash(path)
}
