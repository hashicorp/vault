package file

import (
	"bytes"

	"github.com/hashicorp/go-syslog"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/logical"
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

	// Get the logger
	logger, err := gsyslog.NewLogger(gsyslog.LOG_INFO, facility, tag)
	if err != nil {
		return nil, err
	}

	b := &Backend{
		logger: logger,
	}
	return b, nil
}

// Backend is the audit backend for the syslog-based audit store.
type Backend struct {
	logger gsyslog.Syslogger
}

func (b *Backend) LogRequest(auth *logical.Auth, req *logical.Request) error {
	var buf bytes.Buffer
	var format audit.FormatJSON
	if err := format.FormatRequest(&buf, auth, req); err != nil {
		return err
	}
	_, err := b.logger.Write(buf.Bytes())
	return err
}

func (b *Backend) LogResponse(auth *logical.Auth, req *logical.Request,
	resp *logical.Response, err error) error {
	var buf bytes.Buffer
	var format audit.FormatJSON
	if err := format.FormatResponse(&buf, auth, req, resp, err); err != nil {
		return err
	}
	_, err = b.logger.Write(buf.Bytes())
	return err
}
