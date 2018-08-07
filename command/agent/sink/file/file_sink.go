package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/command/agent/sink"
)

// fileSink is a Sink implementation that writes a token to a file
type fileSink struct {
	path   string
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

	if err := f.WriteToken(""); err != nil {
		return nil, errwrap.Wrapf("error during write check: {{err}}", err)
	}

	f.logger.Info("file sink configured", "path", f.path)

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
		return errwrap.Wrapf("error generating a uuid during write check: {{err}}", err)
	}

	targetDir := filepath.Dir(f.path)
	fileName := filepath.Base(f.path)
	tmpSuffix := strings.Split(u, "-")[0]

	tmpFile, err := os.OpenFile(filepath.Join(targetDir, fmt.Sprintf("%s.tmp.%s", fileName, tmpSuffix)), os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error opening temp file in dir %s for writing: {{err}}", targetDir), err)
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
		return errwrap.Wrapf(fmt.Sprintf("error writing to %s: {{err}}", tmpFile.Name()), err)
	}

	err = tmpFile.Close()
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error closing %s: {{err}}", tmpFile.Name()), err)
	}

	// Now, if we were just doing a write check (blank token), remove the file
	// and exit; otherwise, atomically rename it
	if token == "" {
		err = os.Remove(tmpFile.Name())
		if err != nil {
			return errwrap.Wrapf(fmt.Sprintf("error removing temp file %s during write check: {{err}}", tmpFile.Name()), err)
		}
		return nil
	}

	err = os.Rename(tmpFile.Name(), f.path)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("error renaming temp file %s to target file %s: {{err}}", tmpFile.Name(), f.path), err)
	}

	f.logger.Info("token written", "path", f.path)
	return nil
}
