// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// SetUpGoModSourceFromPath opens the file at path, reads its contents into
// source.Data, and sets source.Name to path. The file is closed before the
// function returns, with any close error joined to the return error.
func SetUpGoModSourceFromPath(path string, source *ModSource) (err error) {
	if source == nil {
		return errors.New("you must provide a mod source")
	}

	aPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	f, err := os.Open(aPath)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	source.Name = path
	source.Data, err = io.ReadAll(f)

	return err
}
