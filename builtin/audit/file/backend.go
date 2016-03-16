package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/copystructure"
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

	b := &Backend{
		path:         path,
		logRaw:       logRaw,
		hmacAccessor: hmacAccessor,
		salt:         conf.Salt,
	}

	// Ensure that the file can be successfully opened for writing;
	// otherwise it will be too late to catch later without problems
	// (ref: https://github.com/hashicorp/vault/issues/550)
	if err := b.open(); err != nil {
		return nil, fmt.Errorf("sanity check failed; unable to open %s for writing", path)
	}

	return b, nil
}

// Backend is the audit backend for the file-based audit store.
//
// NOTE: This audit backend is currently very simple: it appends to a file.
// It doesn't do anything more at the moment to assist with rotation
// or reset the write cursor, this should be done in the future.
type Backend struct {
	path         string
	logRaw       bool
	hmacAccessor bool
	salt         *salt.Salt

	once sync.Once
	f    *os.File
}

func (b *Backend) GetHash(data string) string {
	return audit.HashString(b.salt, data)
}

func (b *Backend) LogRequest(auth *logical.Auth, req *logical.Request, outerErr error) error {
	if err := b.open(); err != nil {
		return err
	}
	if !b.logRaw {
		// Before we copy the structure we must nil out some data
		// otherwise we will cause reflection to panic and die
		if req.Connection != nil && req.Connection.ConnState != nil {
			origReq := req
			origState := req.Connection.ConnState
			req.Connection.ConnState = nil
			defer func() {
				origReq.Connection.ConnState = origState
			}()
		}

		// Copy the structures
		cp, err := copystructure.Copy(auth)
		if err != nil {
			return err
		}
		auth = cp.(*logical.Auth)

		cp, err = copystructure.Copy(req)
		if err != nil {
			return err
		}
		req = cp.(*logical.Request)

		// Hash any sensitive information
		if err := audit.Hash(b.salt, auth); err != nil {
			return err
		}
		if err := audit.Hash(b.salt, req); err != nil {
			return err
		}

	}

	var format audit.FormatJSON
	return format.FormatRequest(b.f, auth, req, outerErr)
}

func (b *Backend) LogResponse(
	auth *logical.Auth,
	req *logical.Request,
	resp *logical.Response,
	err error) error {
	if err := b.open(); err != nil {
		return err
	}
	if !b.logRaw {
		// Before we copy the structure we must nil out some data
		// otherwise we will cause reflection to panic and die
		if req.Connection != nil && req.Connection.ConnState != nil {
			origReq := req
			origState := req.Connection.ConnState
			req.Connection.ConnState = nil
			defer func() {
				origReq.Connection.ConnState = origState
			}()
		}

		// Copy the structure
		cp, err := copystructure.Copy(auth)
		if err != nil {
			return err
		}
		auth = cp.(*logical.Auth)

		cp, err = copystructure.Copy(req)
		if err != nil {
			return err
		}
		req = cp.(*logical.Request)

		cp, err = copystructure.Copy(resp)
		if err != nil {
			return err
		}
		resp = cp.(*logical.Response)

		// Hash any sensitive information

		// Cache and restore accessor in the auth
		var accessor string
		if !b.hmacAccessor && auth != nil && auth.Accessor != "" {
			accessor = auth.Accessor
		}
		if err := audit.Hash(b.salt, auth); err != nil {
			return err
		}
		if accessor != "" {
			auth.Accessor = accessor
		}

		if err := audit.Hash(b.salt, req); err != nil {
			return err
		}

		// Cache and restore accessor in the response
		accessor = ""
		if !b.hmacAccessor && resp != nil && resp.Auth != nil && resp.Auth.Accessor != "" {
			accessor = resp.Auth.Accessor
		}
		if err := audit.Hash(b.salt, resp); err != nil {
			return err
		}
		if accessor != "" {
			resp.Auth.Accessor = accessor
		}
	}

	var format audit.FormatJSON
	return format.FormatResponse(b.f, auth, req, resp, err)
}

func (b *Backend) open() error {
	if b.f != nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(b.path), 0600); err != nil {
		return err
	}

	var err error
	b.f, err = os.OpenFile(b.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	return nil
}
