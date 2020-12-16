package ioutils // import "github.com/ory/dockertest/v3/docker/pkg/ioutils"

import (
	"io/ioutil"

	"github.com/ory/dockertest/v3/docker/pkg/longpath"
)

// TempDir is the equivalent of ioutil.TempDir, except that the result is in Windows longpath format.
func TempDir(dir, prefix string) (string, error) {
	tempDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return "", err
	}
	return longpath.AddPrefix(tempDir), nil
}
