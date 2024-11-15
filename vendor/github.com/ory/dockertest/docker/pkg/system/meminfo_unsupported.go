// +build !linux,!windows

package system // import "github.com/ory/dockertest/docker/pkg/system"

// ReadMemInfo is not supported on platforms other than linux and windows.
func ReadMemInfo() (*MemInfo, error) {
	return nil, ErrNotSupportedPlatform
}
