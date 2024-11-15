// +build !linux

package archive // import "github.com/ory/dockertest/docker/pkg/archive"

func getWhiteoutConverter(format WhiteoutFormat) tarWhiteoutConverter {
	return nil
}
