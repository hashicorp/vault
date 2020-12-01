// +build !linux

package archive // import "github.com/ory/dockertest/v3/docker/pkg/archive"

func getWhiteoutConverter(format WhiteoutFormat) tarWhiteoutConverter {
	return nil
}
