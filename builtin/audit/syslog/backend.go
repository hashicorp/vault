package file

import (
	"bytes"
	"strconv"

	"github.com/hashicorp/go-syslog"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/copystructure"
)

func Factory(conf map[string]string) (audit.Backend, error) {
	// Get facility or default to AUTH
	facility, ok := conf["facility"]
	if !ok {
		facility = "AUTH"
	}

	// Get tag or default to 'vault'
	tag, ok := conf["tag"]
	if !ok {
		tag = "vault"
	}

	// Check if raw logging is enabled
	logRaw := false
	if raw, ok := conf["log_raw"]; ok {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, err
		}
		logRaw = b
	}

	// Get the logger
	logger, err := gsyslog.NewLogger(gsyslog.LOG_INFO, facility, tag)
	if err != nil {
		return nil, err
	}

	b := &Backend{
		logger: logger,
		logRaw: logRaw,
	}
	return b, nil
}

// Backend is the audit backend for the syslog-based audit store.
type Backend struct {
	logger gsyslog.Syslogger
	logRaw bool
}

func (b *Backend) LogRequest(auth *logical.Auth, req *logical.Request, outerErr error) error {
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
		if err := audit.Hash(auth); err != nil {
			return err
		}
		if err := audit.Hash(req); err != nil {
			return err
		}
	}

	// Encode the entry as JSON
	var buf bytes.Buffer
	var format audit.FormatJSON
	if err := format.FormatRequest(&buf, auth, req, outerErr); err != nil {
		return err
	}

	// Write out to syslog
	_, err := b.logger.Write(buf.Bytes())
	return err
}

func (b *Backend) LogResponse(auth *logical.Auth, req *logical.Request,
	resp *logical.Response, err error) error {
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
		if err := audit.Hash(auth); err != nil {
			return err
		}
		if err := audit.Hash(req); err != nil {
			return err
		}
		if err := audit.Hash(resp); err != nil {
			return err
		}
	}

	// Encode the entry as JSON
	var buf bytes.Buffer
	var format audit.FormatJSON
	if err := format.FormatResponse(&buf, auth, req, resp, err); err != nil {
		return err
	}

	// Write otu to syslog
	_, err = b.logger.Write(buf.Bytes())
	return err
}
