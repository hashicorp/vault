package stdout

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/logical"
)

func Factory(conf *audit.BackendConfig) (audit.Backend, error) {
	if conf.Salt == nil {
		return nil, fmt.Errorf("Nil salt passed in")
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

	b := &Backend{
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

type Backend struct {
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

	_, err := os.Stdout.Write(buf.Bytes())
	return err
}

func (b *Backend) LogResponse(auth *logical.Auth, req *logical.Request,
	resp *logical.Response, err error) error {
	var buf bytes.Buffer
	if err := b.formatter.FormatResponse(&buf, b.formatConfig, auth, req, resp, err); err != nil {
		return err
	}

	_, err = os.Stdout.Write(buf.Bytes())
	return err
}

func (b *Backend) Reload() error {
	return nil
}
