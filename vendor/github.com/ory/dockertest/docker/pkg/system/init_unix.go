// +build !windows

package system // import "github.com/ory/dockertest/docker/pkg/system"

// InitLCOW does nothing since LCOW is a windows only feature
func InitLCOW(experimental bool) {
}
