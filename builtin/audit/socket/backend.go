package socket

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/parseutil"
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

	socketType, ok := conf.Config["socket_type"]
	if !ok {
		socketType = "tcp"
	}

	writeDeadline, ok := conf.Config["write_timeout"]
	if !ok {
		writeDeadline = "2s"
	}
	writeDuration, err := parseutil.ParseDurationSecond(writeDeadline)
	if err != nil {
		return nil, err
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

	conn, err := net.Dial(socketType, address)
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
		writeDuration: writeDuration,
		address:       address,
		socketType:    socketType,
	}

	switch format {
	case "json":
		b.formatter.AuditFormatWriter = &audit.JSONFormatWriter{
			Prefix: conf.Config["prefix"],
		}
	case "jsonx":
		b.formatter.AuditFormatWriter = &audit.JSONxFormatWriter{
			Prefix: conf.Config["prefix"],
		}
	}

	return b, nil
}

// Backend is the audit backend for the socket audit transport.
type Backend struct {
	connection net.Conn

	formatter    audit.AuditFormatter
	formatConfig audit.FormatterConfig

	writeDuration time.Duration
	address       string
	socketType    string

	sync.Mutex
}

func (b *Backend) GetHash(data string) string {
	return audit.HashString(b.formatConfig.Salt, data)
}

func (b *Backend) LogRequest(auth *logical.Auth, req *logical.Request, outerErr error) error {
	var buf bytes.Buffer
	if err := b.formatter.FormatRequest(&buf, b.formatConfig, auth, req, outerErr); err != nil {
		return err
	}

	b.Lock()
	defer b.Unlock()

	err := b.write(buf.Bytes())
	if err != nil {
		rErr := b.reconnect()
		if rErr != nil {
			err = multierror.Append(err, rErr)
		} else {
			// Try once more after reconnecting
			err = b.write(buf.Bytes())
		}
	}

	return err
}

func (b *Backend) LogResponse(auth *logical.Auth, req *logical.Request,
	resp *logical.Response, outerErr error) error {
	var buf bytes.Buffer
	if err := b.formatter.FormatResponse(&buf, b.formatConfig, auth, req, resp, outerErr); err != nil {
		return err
	}

	b.Lock()
	defer b.Unlock()

	err := b.write(buf.Bytes())
	if err != nil {
		rErr := b.reconnect()
		if rErr != nil {
			err = multierror.Append(err, rErr)
		} else {
			// Try once more after reconnecting
			err = b.write(buf.Bytes())
		}
	}

	return err
}

func (b *Backend) write(buf []byte) error {
	err := b.connection.SetWriteDeadline(time.Now().Add(b.writeDuration))
	if err != nil {
		return err
	}

	_, err = b.connection.Write(buf)
	if err != nil {
		return err
	}

	return err
}

func (b *Backend) reconnect() error {
	conn, err := net.Dial(b.socketType, b.address)
	if err != nil {
		return err
	}

	b.connection.Close()
	b.connection = conn

	return nil
}

func (b *Backend) Reload() error {
	b.Lock()
	defer b.Unlock()

	err := b.reconnect()

	return err
}
