package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/logical"
)

func Factory(conf *audit.BackendConfig) (audit.Backend, error) {
	if conf.Salt == nil {
		return nil, fmt.Errorf("nil salt")
	}

	path, ok := conf.Config["file_path"]
	if !ok {
		path, ok = conf.Config["path"]
		if !ok {
			return nil, fmt.Errorf("file_path is required")
		}
	}

	format, ok := conf.Config["format"]
	if !ok {
		format = "json"
	}
	switch format {
	case "json", "jsonx":
	default:
		return nil, fmt.Errorf("unknown format type %s", format)
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
		mode = os.FileMode(m)
	}

	b := &Backend{
		path: path,
		mode: mode,
		formatConfig: audit.FormatterConfig{
			Raw:          logRaw,
			Salt:         conf.Salt,
			HMACAccessor: hmacAccessor,
		},
	}

	switch format {
	case "json":
		b.formatter.AuditFormatWriter = &audit.JSONFormatWriter{}
	case "jsonx":
		b.formatter.AuditFormatWriter = &audit.JSONxFormatWriter{}
	}

	// Ensure that the file can be successfully opened for writing;
	// otherwise it will be too late to catch later without problems
	// (ref: https://github.com/hashicorp/vault/issues/550)
	if err := b.open(); err != nil {
		return nil, fmt.Errorf("sanity check failed; unable to open %s for writing: %v", path, err)
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

	fileLock sync.RWMutex
	f        *os.File
	mode     os.FileMode
}

func (b *Backend) GetHash(data string) string {
	return audit.HashString(b.formatConfig.Salt, data)
}

func (b *Backend) LogRequest(auth *logical.Auth, req *logical.Request, outerErr error) error {
	b.fileLock.Lock()
	defer b.fileLock.Unlock()

	if err := b.open(); err != nil {
		return err
	}

	return b.formatter.FormatRequest(b.f, b.formatConfig, auth, req, outerErr)
}

func (b *Backend) LogResponse(
	auth *logical.Auth,
	req *logical.Request,
	resp *logical.Response,
	err error) error {

	b.fileLock.Lock()
	defer b.fileLock.Unlock()

	if err := b.open(); err != nil {
		return err
	}

	return b.formatter.FormatResponse(b.f, b.formatConfig, auth, req, resp, err)
}

// The file lock must be held before calling this
func (b *Backend) open() error {
	if b.f != nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(b.path), b.mode); err != nil {
		return err
	}

	var err error
	b.f, err = os.OpenFile(b.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, b.mode)
	if err != nil {
		return err
	}

	// Change the file mode in case the log file already existed
	err = os.Chmod(b.path, b.mode)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) Reload() error {
	b.fileLock.Lock()
	defer b.fileLock.Unlock()

	if b.f == nil {
		return b.open()
	}

	err := b.f.Close()
	// Set to nil here so that even if we error out, on the next access open()
	// will be tried
	b.f = nil
	if err != nil {
		return err
	}

	return b.open()
}
