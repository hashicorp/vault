package msdsn

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type (
	Encryption int
	Log        uint64
)

const (
	EncryptionOff      = 0
	EncryptionRequired = 1
	EncryptionDisabled = 3
)

const (
	LogErrors      Log = 1
	LogMessages    Log = 2
	LogRows        Log = 4
	LogSQL         Log = 8
	LogParams      Log = 16
	LogTransaction Log = 32
	LogDebug       Log = 64
	LogRetries     Log = 128
)

type Config struct {
	Port       uint64
	Host       string
	Instance   string
	Database   string
	User       string
	Password   string
	Encryption Encryption
	TLSConfig  *tls.Config

	FailOverPartner string
	FailOverPort    uint64

	// If true the TLSConfig servername should use the routed server.
	HostInCertificateProvided bool

	// Read Only intent for application database.
	// NOTE: This does not make queries to most databases read-only.
	ReadOnlyIntent bool

	LogFlags Log

	ServerSPN   string
	Workstation string
	AppName     string

	// If true disables database/sql's automatic retry of queries
	// that start on bad connections.
	DisableRetry bool

	// Do not use the following.

	DialTimeout time.Duration // DialTimeout defaults to 15s. Set negative to disable.
	ConnTimeout time.Duration // Use context for timeouts.
	KeepAlive   time.Duration // Leave at default.
	PacketSize  uint16
}

func SetupTLS(certificate string, insecureSkipVerify bool, hostInCertificate string) (*tls.Config, error) {
	config := tls.Config{
		ServerName:         hostInCertificate,
		InsecureSkipVerify: insecureSkipVerify,

		// fix for https://github.com/denisenkom/go-mssqldb/issues/166
		// Go implementation of TLS payload size heuristic algorithm splits single TDS package to multiple TCP segments,
		// while SQL Server seems to expect one TCP segment per encrypted TDS package.
		// Setting DynamicRecordSizingDisabled to true disables that algorithm and uses 16384 bytes per TLS package
		DynamicRecordSizingDisabled: true,
	}
	if len(certificate) == 0 {
		return &config, nil
	}
	pem, err := ioutil.ReadFile(certificate)
	if err != nil {
		return nil, fmt.Errorf("cannot read certificate %q: %w", certificate, err)
	}
	if strings.Contains(config.ServerName, ":") && !insecureSkipVerify {
		err := setupTLSCommonName(&config, pem)
		if err != skipSetup {
			return &config, err
		}
	}
	certs := x509.NewCertPool()
	certs.AppendCertsFromPEM(pem)
	config.RootCAs = certs
	return &config, nil
}

var skipSetup = errors.New("skip setting up TLS")

func Parse(dsn string) (Config, map[string]string, error) {
	p := Config{}

	var params map[string]string
	if strings.HasPrefix(dsn, "odbc:") {
		parameters, err := splitConnectionStringOdbc(dsn[len("odbc:"):])
		if err != nil {
			return p, params, err
		}
		params = parameters
	} else if strings.HasPrefix(dsn, "sqlserver://") {
		parameters, err := splitConnectionStringURL(dsn)
		if err != nil {
			return p, params, err
		}
		params = parameters
	} else {
		params = splitConnectionString(dsn)
	}

	strlog, ok := params["log"]
	if ok {
		flags, err := strconv.ParseUint(strlog, 10, 64)
		if err != nil {
			return p, params, fmt.Errorf("invalid log parameter '%s': %s", strlog, err.Error())
		}
		p.LogFlags = Log(flags)
	}
	server := params["server"]
	parts := strings.SplitN(server, `\`, 2)
	p.Host = parts[0]
	if p.Host == "." || strings.ToUpper(p.Host) == "(LOCAL)" || p.Host == "" {
		p.Host = "localhost"
	}
	if len(parts) > 1 {
		p.Instance = parts[1]
	}
	p.Database = params["database"]
	p.User = params["user id"]
	p.Password = params["password"]

	p.Port = 0
	strport, ok := params["port"]
	if ok {
		var err error
		p.Port, err = strconv.ParseUint(strport, 10, 16)
		if err != nil {
			f := "invalid tcp port '%v': %v"
			return p, params, fmt.Errorf(f, strport, err.Error())
		}
	}

	// https://docs.microsoft.com/en-us/sql/database-engine/configure-windows/configure-the-network-packet-size-server-configuration-option\
	strpsize, ok := params["packet size"]
	if ok {
		var err error
		psize, err := strconv.ParseUint(strpsize, 0, 16)
		if err != nil {
			f := "invalid packet size '%v': %v"
			return p, params, fmt.Errorf(f, strpsize, err.Error())
		}

		// Ensure packet size falls within the TDS protocol range of 512 to 32767 bytes
		// NOTE: Encrypted connections have a maximum size of 16383 bytes.  If you request
		// a higher packet size, the server will respond with an ENVCHANGE request to
		// alter the packet size to 16383 bytes.
		p.PacketSize = uint16(psize)
		if p.PacketSize < 512 {
			p.PacketSize = 512
		} else if p.PacketSize > 32767 {
			p.PacketSize = 32767
		}
	}

	// https://msdn.microsoft.com/en-us/library/dd341108.aspx
	//
	// Do not set a connection timeout. Use Context to manage such things.
	// Default to zero, but still allow it to be set.
	if strconntimeout, ok := params["connection timeout"]; ok {
		timeout, err := strconv.ParseUint(strconntimeout, 10, 64)
		if err != nil {
			f := "invalid connection timeout '%v': %v"
			return p, params, fmt.Errorf(f, strconntimeout, err.Error())
		}
		p.ConnTimeout = time.Duration(timeout) * time.Second
	}
	p.DialTimeout = 15 * time.Second
	if strdialtimeout, ok := params["dial timeout"]; ok {
		timeout, err := strconv.ParseUint(strdialtimeout, 10, 64)
		if err != nil {
			f := "invalid dial timeout '%v': %v"
			return p, params, fmt.Errorf(f, strdialtimeout, err.Error())
		}
		p.DialTimeout = time.Duration(timeout) * time.Second
	}

	// default keep alive should be 30 seconds according to spec:
	// https://msdn.microsoft.com/en-us/library/dd341108.aspx
	p.KeepAlive = 30 * time.Second
	if keepAlive, ok := params["keepalive"]; ok {
		timeout, err := strconv.ParseUint(keepAlive, 10, 64)
		if err != nil {
			f := "invalid keepAlive value '%s': %s"
			return p, params, fmt.Errorf(f, keepAlive, err.Error())
		}
		p.KeepAlive = time.Duration(timeout) * time.Second
	}

	var (
		trustServerCert   = false
		certificate       = ""
		hostInCertificate = ""
	)
	encrypt, ok := params["encrypt"]
	if ok {
		if strings.EqualFold(encrypt, "DISABLE") {
			p.Encryption = EncryptionDisabled
		} else {
			e, err := strconv.ParseBool(encrypt)
			if err != nil {
				f := "invalid encrypt '%s': %s"
				return p, params, fmt.Errorf(f, encrypt, err.Error())
			}
			if e {
				p.Encryption = EncryptionRequired
			}
		}
	} else {
		trustServerCert = true
	}
	trust, ok := params["trustservercertificate"]
	if ok {
		var err error
		trustServerCert, err = strconv.ParseBool(trust)
		if err != nil {
			f := "invalid trust server certificate '%s': %s"
			return p, params, fmt.Errorf(f, trust, err.Error())
		}
	}
	certificate = params["certificate"]
	hostInCertificate, ok = params["hostnameincertificate"]
	if ok {
		p.HostInCertificateProvided = true
	} else {
		hostInCertificate = p.Host
		p.HostInCertificateProvided = false
	}

	if p.Encryption != EncryptionDisabled {
		var err error
		p.TLSConfig, err = SetupTLS(certificate, trustServerCert, hostInCertificate)
		if err != nil {
			return p, params, fmt.Errorf("failed to setup TLS: %w", err)
		}
	}

	serverSPN, ok := params["serverspn"]
	if ok {
		p.ServerSPN = serverSPN
	} else {
		p.ServerSPN = generateSpn(p.Host, p.Port)
	}

	workstation, ok := params["workstation id"]
	if ok {
		p.Workstation = workstation
	} else {
		workstation, err := os.Hostname()
		if err == nil {
			p.Workstation = workstation
		}
	}

	appname, ok := params["app name"]
	if !ok {
		appname = "go-mssqldb"
	}
	p.AppName = appname

	appintent, ok := params["applicationintent"]
	if ok {
		if appintent == "ReadOnly" {
			if p.Database == "" {
				return p, params, fmt.Errorf("database must be specified when ApplicationIntent is ReadOnly")
			}
			p.ReadOnlyIntent = true
		}
	}

	failOverPartner, ok := params["failoverpartner"]
	if ok {
		p.FailOverPartner = failOverPartner
	}

	failOverPort, ok := params["failoverport"]
	if ok {
		var err error
		p.FailOverPort, err = strconv.ParseUint(failOverPort, 0, 16)
		if err != nil {
			f := "invalid failover port '%v': %v"
			return p, params, fmt.Errorf(f, failOverPort, err.Error())
		}
	}

	disableRetry, ok := params["disableretry"]
	if ok {
		var err error
		p.DisableRetry, err = strconv.ParseBool(disableRetry)
		if err != nil {
			f := "invalid disableRetry '%s': %s"
			return p, params, fmt.Errorf(f, disableRetry, err.Error())
		}
	} else {
		p.DisableRetry = disableRetryDefault
	}

	return p, params, nil
}

// convert connectionParams to url style connection string
// used mostly for testing
func (p Config) URL() *url.URL {
	q := url.Values{}
	if p.Database != "" {
		q.Add("database", p.Database)
	}
	if p.LogFlags != 0 {
		q.Add("log", strconv.FormatUint(uint64(p.LogFlags), 10))
	}
	host := p.Host
	if p.Port > 0 {
		host = fmt.Sprintf("%s:%d", p.Host, p.Port)
	}
	q.Add("disableRetry", fmt.Sprintf("%t", p.DisableRetry))
	res := url.URL{
		Scheme: "sqlserver",
		Host:   host,
		User:   url.UserPassword(p.User, p.Password),
	}
	if p.Instance != "" {
		res.Path = p.Instance
	}
	if len(q) > 0 {
		res.RawQuery = q.Encode()
	}
	return &res
}

var adoSynonyms = map[string]string{
	"application name": "app name",
	"data source":      "server",
	"address":          "server",
	"network address":  "server",
	"addr":             "server",
	"user":             "user id",
	"uid":              "user id",
	"initial catalog":  "database",
}

func splitConnectionString(dsn string) (res map[string]string) {
	res = map[string]string{}
	parts := strings.Split(dsn, ";")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		lst := strings.SplitN(part, "=", 2)
		name := strings.TrimSpace(strings.ToLower(lst[0]))
		if len(name) == 0 {
			continue
		}
		var value string = ""
		if len(lst) > 1 {
			value = strings.TrimSpace(lst[1])
		}
		synonym, hasSynonym := adoSynonyms[name]
		if hasSynonym {
			name = synonym
		}
		// "server" in ADO can include a protocol and a port.
		// We only support tcp protocol
		if name == "server" {
			value = strings.TrimPrefix(value, "tcp:")
			serverParts := strings.Split(value, ",")
			if len(serverParts) == 2 && len(serverParts[1]) > 0 {
				value = serverParts[0]
				res["port"] = serverParts[1]
			}
		}
		res[name] = value
	}
	return res
}

// Splits a URL of the form sqlserver://username:password@host/instance?param1=value&param2=value
func splitConnectionStringURL(dsn string) (map[string]string, error) {
	res := map[string]string{}

	u, err := url.Parse(dsn)
	if err != nil {
		return res, err
	}

	if u.Scheme != "sqlserver" {
		return res, fmt.Errorf("scheme %s is not recognized", u.Scheme)
	}

	if u.User != nil {
		res["user id"] = u.User.Username()
		p, exists := u.User.Password()
		if exists {
			res["password"] = p
		}
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}

	if len(u.Path) > 0 {
		res["server"] = host + "\\" + u.Path[1:]
	} else {
		res["server"] = host
	}

	if len(port) > 0 {
		res["port"] = port
	}

	query := u.Query()
	for k, v := range query {
		if len(v) > 1 {
			return res, fmt.Errorf("key %s provided more than once", k)
		}
		res[strings.ToLower(k)] = v[0]
	}

	return res, nil
}

// Splits a URL in the ODBC format
func splitConnectionStringOdbc(dsn string) (map[string]string, error) {
	res := map[string]string{}

	type parserState int
	const (
		// Before the start of a key
		parserStateBeforeKey parserState = iota

		// Inside a key
		parserStateKey

		// Beginning of a value. May be bare or braced
		parserStateBeginValue

		// Inside a bare value
		parserStateBareValue

		// Inside a braced value
		parserStateBracedValue

		// A closing brace inside a braced value.
		// May be the end of the value or an escaped closing brace, depending on the next character
		parserStateBracedValueClosingBrace

		// After a value. Next character should be a semicolon or whitespace.
		parserStateEndValue
	)

	var state = parserStateBeforeKey

	var key string
	var value string

	for i, c := range dsn {
		switch state {
		case parserStateBeforeKey:
			switch {
			case c == '=':
				return res, fmt.Errorf("unexpected character = at index %d. Expected start of key or semi-colon or whitespace", i)
			case !unicode.IsSpace(c) && c != ';':
				state = parserStateKey
				key += string(c)
			}

		case parserStateKey:
			switch c {
			case '=':
				key = normalizeOdbcKey(key)
				state = parserStateBeginValue

			case ';':
				// Key without value
				key = normalizeOdbcKey(key)
				res[key] = value
				key = ""
				value = ""
				state = parserStateBeforeKey

			default:
				key += string(c)
			}

		case parserStateBeginValue:
			switch {
			case c == '{':
				state = parserStateBracedValue
			case c == ';':
				// Empty value
				res[key] = value
				key = ""
				state = parserStateBeforeKey
			case unicode.IsSpace(c):
				// Ignore whitespace
			default:
				state = parserStateBareValue
				value += string(c)
			}

		case parserStateBareValue:
			if c == ';' {
				res[key] = strings.TrimRightFunc(value, unicode.IsSpace)
				key = ""
				value = ""
				state = parserStateBeforeKey
			} else {
				value += string(c)
			}

		case parserStateBracedValue:
			if c == '}' {
				state = parserStateBracedValueClosingBrace
			} else {
				value += string(c)
			}

		case parserStateBracedValueClosingBrace:
			if c == '}' {
				// Escaped closing brace
				value += string(c)
				state = parserStateBracedValue
				continue
			}

			// End of braced value
			res[key] = value
			key = ""
			value = ""

			// This character is the first character past the end,
			// so it needs to be parsed like the parserStateEndValue state.
			state = parserStateEndValue
			switch {
			case c == ';':
				state = parserStateBeforeKey
			case unicode.IsSpace(c):
				// Ignore whitespace
			default:
				return res, fmt.Errorf("unexpected character %c at index %d. Expected semi-colon or whitespace", c, i)
			}

		case parserStateEndValue:
			switch {
			case c == ';':
				state = parserStateBeforeKey
			case unicode.IsSpace(c):
				// Ignore whitespace
			default:
				return res, fmt.Errorf("unexpected character %c at index %d. Expected semi-colon or whitespace", c, i)
			}
		}
	}

	switch state {
	case parserStateBeforeKey: // Okay
	case parserStateKey: // Unfinished key. Treat as key without value.
		key = normalizeOdbcKey(key)
		res[key] = value
	case parserStateBeginValue: // Empty value
		res[key] = value
	case parserStateBareValue:
		res[key] = strings.TrimRightFunc(value, unicode.IsSpace)
	case parserStateBracedValue:
		return res, fmt.Errorf("unexpected end of braced value at index %d", len(dsn))
	case parserStateBracedValueClosingBrace: // End of braced value
		res[key] = value
	case parserStateEndValue: // Okay
	}

	return res, nil
}

// Normalizes the given string as an ODBC-format key
func normalizeOdbcKey(s string) string {
	return strings.ToLower(strings.TrimRightFunc(s, unicode.IsSpace))
}

func generateSpn(host string, port uint64) string {
	return fmt.Sprintf("MSSQLSvc/%s:%d", host, port)
}
