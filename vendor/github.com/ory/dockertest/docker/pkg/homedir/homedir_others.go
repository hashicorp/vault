// +build !linux

package homedir // import "github.com/ory/dockertest/docker/pkg/homedir"

import (
	"errors"
)

// GetStatic is not needed for non-linux systems.
// (Precisely, it is needed only for glibc-based linux systems.)
func GetStatic() (string, error) {
	return "", errors.New("homedir.GetStatic() is not supported on this system")
}
