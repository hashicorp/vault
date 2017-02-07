package socket

import (
	"bytes"
	"fmt"
	"net"
	"strconv"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/logical"
)

func Factory(conf *audit.BackendConfig) (audit.Backend, error) {
	if conf.Salt == nil {
		return nil, fmt.Errorf("nil salt passed in")
	}

	address, ok := conf.Config["address"]
	if !ok {
		return nil, fmt.Errorf("address is required")
	}

	socket_type, ok := conf.Config["socket_type"]
	if !ok {
		socket_type = "tcp"
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

	conn, err := net.Dial(socket_type, address)
	if err != nil {
		return nil, err
	}

	b := &Backend{
		connection: conn,
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

	return b, nil
}

// Backend is the audit backend for the socket audit transport.
type Backend struct {
	connection net.Conn

	formatter    audit.AuditFormatter
	formatConfig audit.FormatterConfig
}

func (b *Backend) GetHash(data string) string {
	return audit.HashString(b.formatConfig.Salt, data)
}

func (b *Backend) LogRequest(auth *logical.Auth, req *logical.Request, outerErr error) error {
	var buf bytes.Buffer
	if err := b.formatter.FormatRequest(&buf, b.formatConfig, auth, req, outerErr); err != nil {
		return err
	}

	b.connection.Write(buf.Bytes())
	return nil
}

func (b *Backend) LogResponse(auth *logical.Auth, req *logical.Request,
	resp *logical.Response, err error) error {
	var buf bytes.Buffer
	if err := b.formatter.FormatResponse(&buf, b.formatConfig, auth, req, resp, err); err != nil {
		return err
	}
	b.connection.Write(buf.Bytes())
	return nil
}

func (b *Backend) Reload() error {
	return nil
}
