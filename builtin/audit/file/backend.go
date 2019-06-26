package file

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *audit.BackendConfig) (audit.Backend, error) {
	if conf.SaltConfig == nil {
		return nil, fmt.Errorf("nil salt config")
	}
	if conf.SaltView == nil {
		return nil, fmt.Errorf("nil salt view")
	}

	path, ok := conf.Config["file_path"]
	if !ok {
		path, ok = conf.Config["path"]
		if !ok {
			return nil, fmt.Errorf("file_path is required")
		}
	}

	// normalize path if configured for stdout
	if strings.EqualFold(path, "stdout") {
		path = "stdout"
	}
	if strings.EqualFold(path, "discard") {
		path = "discard"
	}

	format, ok := conf.Config["format"]
	if !ok {
		format = "json"
	}
	switch format {
	case "json", "jsonx":
	default:
		return nil, fmt.Errorf("unknown format type %q", format)
	}

	// Check if hashing of accessor is disabled
	hmacAccessor := true
	if hmacAccessorRaw, ok := conf.Config["hmac_accessor"]; ok {
		value, err := strconv.ParseBool(hmacAccessorRaw)
		if err != nil {
			return nil, err
		}
		hmacAccessor = value
	}

	// Check if raw logging is enabled
	logRaw := false
	if raw, ok := conf.Config["log_raw"]; ok {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, err
		}
		logRaw = b
	}

	// Check if mode is provided
	mode := os.FileMode(0600)
	if modeRaw, ok := conf.Config["mode"]; ok {
		m, err := strconv.ParseUint(modeRaw, 8, 32)
		if err != nil {
			return nil, err
		}
		if m != 0 {
			mode = os.FileMode(m)
		}
	}

	b := &Backend{
		path:       path,
		mode:       mode,
		saltConfig: conf.SaltConfig,
		saltView:   conf.SaltView,
		formatConfig: audit.FormatterConfig{
			Raw:          logRaw,
			HMACAccessor: hmacAccessor,
		},
	}

	switch format {
	case "json":
		b.formatter.AuditFormatWriter = &audit.JSONFormatWriter{
			Prefix:   conf.Config["prefix"],
			SaltFunc: b.Salt,
		}
	case "jsonx":
		b.formatter.AuditFormatWriter = &audit.JSONxFormatWriter{
			Prefix:   conf.Config["prefix"],
			SaltFunc: b.Salt,
		}
	}

	switch path {
	case "stdout", "discard":
		// no need to test opening file if outputting to stdout or discarding
	default:
		// Ensure that the file can be successfully opened for writing;
		// otherwise it will be too late to catch later without problems
		// (ref: https://github.com/hashicorp/vault/issues/550)
		b.fileLock.RLock()
		defer b.fileLock.RUnlock()
		if err := b.open(); err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("sanity check failed; unable to open %q for writing: {{err}}", path), err)
		}
	}

	return b, nil
}

// Backend is the audit backend for the file-based audit store.
//
// NOTE: This audit backend is currently very simple: it appends to a file.
// It doesn't do anything more at the moment to assist with rotation
// or reset the write cursor, this should be done in the future.
type Backend struct {
	path string

	formatter    audit.AuditFormatter
	formatConfig audit.FormatterConfig

	fileLock      sync.RWMutex // locks file assignment; open/close operations.
	fileWriteLock sync.Mutex   // locks writes to the file
	f             *os.File
	mode          os.FileMode

	saltMutex  sync.RWMutex
	salt       *salt.Salt
	saltConfig *salt.Config
	saltView   logical.Storage
}

type lockedFileWriter struct {
	fileLock *sync.Mutex
	f        *os.File
}

func (lf lockedFileWriter) Write(p []byte) (int, error) {
	lf.fileLock.Lock()
	defer lf.fileLock.Unlock()
	return lf.f.Write(p)
}

var _ audit.Backend = (*Backend)(nil)
var _ io.Writer = lockedFileWriter{}

func (b *Backend) Salt(ctx context.Context) (*salt.Salt, error) {
	b.saltMutex.RLock()
	if b.salt != nil {
		defer b.saltMutex.RUnlock()
		return b.salt, nil
	}
	b.saltMutex.RUnlock()
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	if b.salt != nil {
		return b.salt, nil
	}
	salt, err := salt.NewSalt(ctx, b.saltView, b.saltConfig)
	if err != nil {
		return nil, err
	}
	b.salt = salt
	return salt, nil
}

func (b *Backend) GetHash(ctx context.Context, data string) (string, error) {
	salt, err := b.Salt(ctx)
	if err != nil {
		return "", err
	}
	return audit.HashString(salt, data), nil
}

func (b *Backend) LogRequest(ctx context.Context, in *logical.LogInput) error {
	return b.log(ctx, in, b.formatter.FormatRequest)
}

func (b *Backend) LogResponse(ctx context.Context, in *logical.LogInput) error {
	return b.log(ctx, in, b.formatter.FormatResponse)
}

type formatter func(ctx context.Context, w io.Writer, config audit.FormatterConfig, in *logical.LogInput) error

func (b *Backend) log(ctx context.Context, in *logical.LogInput, format formatter) error {
	b.fileLock.RLock()
	defer b.fileLock.RUnlock()

	switch b.path {
	case "stdout":
		return format(ctx, os.Stdout, b.formatConfig, in)
	case "discard":
		return format(ctx, ioutil.Discard, b.formatConfig, in)
	}

	if err := b.open(); err != nil {
		return err
	}

	writer := lockedFileWriter{f: b.f, fileLock: &b.fileWriteLock}
	if err := format(ctx, writer, b.formatConfig, in); err == nil {
		return nil
	}

	// Opportunistically try to re-open the FD, once per call
	b.close()
	if err := b.open(); err != nil {
		return err
	}

	writer.f = b.f
	return b.formatter.FormatRequest(ctx, writer, b.formatConfig, in)
}

// upgradeFileLock exchanges a held read lock for a write lock. It returns a
// function which can be deferred to exchange the write lock back for a
// read lock.
func (b *Backend) upgradeFileLock() func() {
	// Upgrade to a write lock
	b.fileLock.RUnlock()
	b.fileLock.Lock()
	return func() {
		// Downgrade to read lock
		b.fileLock.Unlock()
		b.fileLock.RLock()
	}
}

// The file lock must be held for reading before calling this
func (b *Backend) close() (err error) {
	defer b.upgradeFileLock()()
	err = b.f.Close()
	b.f = nil
	return
}

// The file lock must be held for reading before calling this
func (b *Backend) open() error {
	if b.f != nil {
		return nil
	}

	defer b.upgradeFileLock()()
	if err := os.MkdirAll(filepath.Dir(b.path), b.mode); err != nil {
		return err
	}

	var err error
	b.f, err = os.OpenFile(b.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, b.mode)
	if err != nil {
		return err
	}

	// Change the file mode in case the log file already existed. We special
	// case /dev/null since we can't chmod it and bypass if the mode is zero
	switch b.path {
	case "/dev/null":
	default:
		if b.mode != 0 {
			err = os.Chmod(b.path, b.mode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Backend) Reload(_ context.Context) error {
	switch b.path {
	case "stdout", "discard":
		return nil
	}

	b.fileLock.RLock()
	defer b.fileLock.RUnlock()

	if b.f == nil {
		return b.open()
	}

	// b.f is set to nil inside of close.
	// Even if we error out, on the next access open() will be tried
	if err := b.close(); err != nil {
		return err
	}

	return b.open()
}

func (b *Backend) Invalidate(_ context.Context) {
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	b.salt = nil
}
