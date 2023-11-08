// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/command-server/agentproxyshared/sink"

	hclog "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
)

// fileSink is a Sink implementation that writes a token to a file
type fileSink struct {
	path   string
	mode   os.FileMode
	logger hclog.Logger
}

// NewFileSink creates a new file sink with the given configuration
func NewFileSink(conf *sink.SinkConfig) (sink.Sink, error) {
	if conf.Logger == nil {
		return nil, errors.New("nil logger provided")
	}

	conf.Logger.Info("creating file sink")

	f := &fileSink{
		logger: conf.Logger,
		mode:   0o640,
	}

	pathRaw, ok := conf.Config["path"]
	if !ok {
		return nil, errors.New("'path' not specified for file sink")
	}
	path, ok := pathRaw.(string)
	if !ok {
		return nil, errors.New("could not parse 'path' as string")
	}

	f.path = path

	if modeRaw, ok := conf.Config["mode"]; ok {
		f.logger.Debug("verifying override for default file sink mode")
		mode, typeOK := modeRaw.(int)
		if !typeOK {
			return nil, errors.New("could not parse 'mode' as integer")
		}

		if !os.FileMode(mode).IsRegular() {
			return nil, fmt.Errorf("file mode does not represent a regular file")
		}

		f.logger.Debug("overriding default file sink", "mode", mode)
		f.mode = os.FileMode(mode)
	}

	if err := f.WriteToken(""); err != nil {
		return nil, fmt.Errorf("error during write check: %w", err)
	}

	f.logger.Info("file sink configured", "path", f.path, "mode", f.mode)

	return f, nil
}

// WriteToken implements the Server interface and writes the token to a path on
// disk. It writes into the path's directory into a temp file and does an
// atomic rename to ensure consistency. If a blank token is passed in, it
// performs a write check but does not write a blank value to the final
// location.
func (f *fileSink) WriteToken(token string) error {
	f.logger.Trace("enter write_token", "path", f.path)
	defer f.logger.Trace("exit write_token", "path", f.path)

	u, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("error generating a uuid during write check: %w", err)
	}

	targetDir := filepath.Dir(f.path)
	fileName := filepath.Base(f.path)
	tmpSuffix := strings.Split(u, "-")[0]

	tmpFile, err := os.OpenFile(filepath.Join(targetDir, fmt.Sprintf("%s.tmp.%s", fileName, tmpSuffix)), os.O_WRONLY|os.O_CREATE, f.mode)
	if err != nil {
		return fmt.Errorf("error opening temp file in dir %s for writing: %w", targetDir, err)
	}

	valToWrite := token
	if token == "" {
		valToWrite = u
	}

	_, err = tmpFile.WriteString(valToWrite)
	if err != nil {
		// Attempt closing and deleting but ignore any error
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return fmt.Errorf("error writing to %s: %w", tmpFile.Name(), err)
	}

	err = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("error closing %s: %w", tmpFile.Name(), err)
	}

	// Now, if we were just doing a write check (blank token), remove the file
	// and exit; otherwise, atomically rename it
	if token == "" {
		err = os.Remove(tmpFile.Name())
		if err != nil {
			return fmt.Errorf("error removing temp file %s during write check: %w", tmpFile.Name(), err)
		}
		return nil
	}

	err = os.Rename(tmpFile.Name(), f.path)
	if err != nil {
		return fmt.Errorf("error renaming temp file %s to target file %s: %w", tmpFile.Name(), f.path, err)
	}

	f.logger.Info("token written", "path", f.path)
	return nil
}
