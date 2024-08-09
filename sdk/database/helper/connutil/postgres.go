// Copyright (c) 2019-2021 Jack Christensen

// MIT License

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Copied from https://github.com/jackc/pgconn/blob/1860f4e57204614f40d05a5c76a43e8d80fde9da/config.go

package connutil

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

// openPostgres parses the connection string and opens a connection to the database.
//
// If sslinline is set, strips the connection string of all ssl settings and
// creates a TLS config based on the settings provided, then uses the
// RegisterConnConfig function to create a new connection. This is necessary
// because the pgx driver does not support the sslinline parameter and instead
// expects to source ssl material from the file system.
//
// Deprecated: openPostgres will be removed in a future version of the Vault SDK.
func openPostgres(driverName, connString string) (*sql.DB, error) {
	if ok, _ := strconv.ParseBool(os.Getenv(pluginutil.PluginUsePostgresSSLInline)); !ok {
		return nil, fmt.Errorf("failed to open postgres connection with deprecated funtion, set feature flag to enable")
	}

	var options pgconn.ParseConfigOptions

	settings := make(map[string]string)
	if connString != "" {
		var err error
		// connString may be a database URL or a DSN
		if strings.HasPrefix(connString, "postgres://") || strings.HasPrefix(connString, "postgresql://") {
			settings, err = parsePostgresURLSettings(connString)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as URL: %w", err)
			}
		} else {
			settings, err = parsePostgresDSNSettings(connString)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as DSN: %w", err)
			}
		}
	}

	// get the inline flag
	sslInline := settings["sslinline"] == "true"

	// if sslinline is not set, open a regular connection
	if !sslInline {
		return sql.Open(driverName, connString)
	}

	// generate a new DSN without the ssl settings
	newConnStr := []string{"sslmode=disable"}
	for k, v := range settings {
		switch k {
		case "sslinline", "sslcert", "sslkey", "sslrootcert", "sslmode":
			continue
		}

		newConnStr = append(newConnStr, fmt.Sprintf("%s='%s'", k, v))
	}

	// parse the updated config
	config, err := pgx.ParseConfig(strings.Join(newConnStr, " "))
	if err != nil {
		return nil, err
	}

	// create a TLS config
	fallbacks := []*pgconn.FallbackConfig{}

	hosts := strings.Split(settings["host"], ",")
	ports := strings.Split(settings["port"], ",")

	for i, host := range hosts {
		var portStr string
		if i < len(ports) {
			portStr = ports[i]
		} else {
			portStr = ports[0]
		}

		port, err := parsePort(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %w", err)
		}

		var tlsConfigs []*tls.Config

		// Ignore TLS settings if Unix domain socket like libpq
		if network, _ := pgconn.NetworkAddress(host, port); network == "unix" {
			tlsConfigs = append(tlsConfigs, nil)
		} else {
			var err error
			tlsConfigs, err = configPostgresTLS(settings, host, options)
			if err != nil {
				return nil, fmt.Errorf("failed to configure TLS: %w", err)
			}
		}

		for _, tlsConfig := range tlsConfigs {
			fallbacks = append(fallbacks, &pgconn.FallbackConfig{
				Host:      host,
				Port:      port,
				TLSConfig: tlsConfig,
			})
		}
	}

	config.Host = fallbacks[0].Host
	config.Port = fallbacks[0].Port
	config.TLSConfig = fallbacks[0].TLSConfig
	config.Fallbacks = fallbacks[1:]

	return sql.Open(driverName, stdlib.RegisterConnConfig(config))
}

// configPostgresTLS uses libpq's TLS parameters to construct  []*tls.Config. It is
// necessary to allow returning multiple TLS configs as sslmode "allow" and
// "prefer" allow fallback.
//
// Copied from https://github.com/jackc/pgconn/blob/1860f4e57204614f40d05a5c76a43e8d80fde9da/config.go
// and modified to read ssl material by value instead of file location.
func configPostgresTLS(settings map[string]string, thisHost string, parseConfigOptions pgconn.ParseConfigOptions) ([]*tls.Config, error) {
	host := thisHost
	sslmode := settings["sslmode"]
	sslrootcert := settings["sslrootcert"]
	sslcert := settings["sslcert"]
	sslkey := settings["sslkey"]
	sslpassword := settings["sslpassword"]
	sslsni := settings["sslsni"]

	// Match libpq default behavior
	if sslmode == "" {
		sslmode = "prefer"
	}
	if sslsni == "" {
		sslsni = "1"
	}

	tlsConfig := &tls.Config{}

	switch sslmode {
	case "disable":
		return []*tls.Config{nil}, nil
	case "allow", "prefer":
		tlsConfig.InsecureSkipVerify = true
	case "require":
		// According to PostgreSQL documentation, if a root CA file exists,
		// the behavior of sslmode=require should be the same as that of verify-ca
		//
		// See https://www.postgresql.org/docs/12/libpq-ssl.html
		if sslrootcert != "" {
			goto nextCase
		}
		tlsConfig.InsecureSkipVerify = true
		break
	nextCase:
		fallthrough
	case "verify-ca":
		// Don't perform the default certificate verification because it
		// will verify the hostname. Instead, verify the server's
		// certificate chain ourselves in VerifyPeerCertificate and
		// ignore the server name. This emulates libpq's verify-ca
		// behavior.
		//
		// See https://github.com/golang/go/issues/21971#issuecomment-332693931
		// and https://pkg.go.dev/crypto/tls?tab=doc#example-Config-VerifyPeerCertificate
		// for more info.
		tlsConfig.InsecureSkipVerify = true
		tlsConfig.VerifyPeerCertificate = func(certificates [][]byte, _ [][]*x509.Certificate) error {
			certs := make([]*x509.Certificate, len(certificates))
			for i, asn1Data := range certificates {
				cert, err := x509.ParseCertificate(asn1Data)
				if err != nil {
					return errors.New("failed to parse certificate from server: " + err.Error())
				}
				certs[i] = cert
			}

			// Leave DNSName empty to skip hostname verification.
			opts := x509.VerifyOptions{
				Roots:         tlsConfig.RootCAs,
				Intermediates: x509.NewCertPool(),
			}
			// Skip the first cert because it's the leaf. All others
			// are intermediates.
			for _, cert := range certs[1:] {
				opts.Intermediates.AddCert(cert)
			}
			_, err := certs[0].Verify(opts)
			return err
		}
	case "verify-full":
		tlsConfig.ServerName = host
	default:
		return nil, errors.New("sslmode is invalid")
	}

	if sslrootcert != "" {
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM([]byte(sslrootcert)) {
			return nil, errors.New("unable to add CA to cert pool")
		}

		tlsConfig.RootCAs = caCertPool
		tlsConfig.ClientCAs = caCertPool
	}

	if (sslcert != "" && sslkey == "") || (sslcert == "" && sslkey != "") {
		return nil, errors.New(`both "sslcert" and "sslkey" are required`)
	}

	if sslcert != "" && sslkey != "" {
		block, _ := pem.Decode([]byte(sslkey))
		var pemKey []byte
		var decryptedKey []byte
		var decryptedError error
		// If PEM is encrypted, attempt to decrypt using pass phrase
		if x509.IsEncryptedPEMBlock(block) {
			// Attempt decryption with pass phrase
			// NOTE: only supports RSA (PKCS#1)
			if sslpassword != "" {
				decryptedKey, decryptedError = x509.DecryptPEMBlock(block, []byte(sslpassword))
			}
			// if sslpassword not provided or has decryption error when use it
			// try to find sslpassword with callback function
			if sslpassword == "" || decryptedError != nil {
				if parseConfigOptions.GetSSLPassword != nil {
					sslpassword = parseConfigOptions.GetSSLPassword(context.Background())
				}
				if sslpassword == "" {
					return nil, fmt.Errorf("unable to find sslpassword")
				}
			}
			decryptedKey, decryptedError = x509.DecryptPEMBlock(block, []byte(sslpassword))
			// Should we also provide warning for PKCS#1 needed?
			if decryptedError != nil {
				return nil, fmt.Errorf("unable to decrypt key: %w", decryptedError)
			}

			pemBytes := pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: decryptedKey,
			}
			pemKey = pem.EncodeToMemory(&pemBytes)
		} else {
			pemKey = pem.EncodeToMemory(block)
		}

		cert, err := tls.X509KeyPair([]byte(sslcert), pemKey)
		if err != nil {
			return nil, fmt.Errorf("unable to load cert: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// Set Server Name Indication (SNI), if enabled by connection parameters.
	// Per RFC 6066, do not set it if the host is a literal IP address (IPv4
	// or IPv6).
	if sslsni == "1" && net.ParseIP(host) == nil {
		tlsConfig.ServerName = host
	}

	switch sslmode {
	case "allow":
		return []*tls.Config{nil, tlsConfig}, nil
	case "prefer":
		return []*tls.Config{tlsConfig, nil}, nil
	case "require", "verify-ca", "verify-full":
		return []*tls.Config{tlsConfig}, nil
	default:
		panic("BUG: bad sslmode should already have been caught")
	}
}

func parsePort(s string) (uint16, error) {
	port, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}
	if port < 1 || port > math.MaxUint16 {
		return 0, errors.New("outside range")
	}
	return uint16(port), nil
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func parsePostgresURLSettings(connString string) (map[string]string, error) {
	settings := make(map[string]string)

	url, err := url.Parse(connString)
	if err != nil {
		return nil, err
	}

	if url.User != nil {
		settings["user"] = url.User.Username()
		if password, present := url.User.Password(); present {
			settings["password"] = password
		}
	}

	// Handle multiple host:port's in url.Host by splitting them into host,host,host and port,port,port.
	var hosts []string
	var ports []string
	for _, host := range strings.Split(url.Host, ",") {
		if host == "" {
			continue
		}
		if isIPOnly(host) {
			hosts = append(hosts, strings.Trim(host, "[]"))
			continue
		}
		h, p, err := net.SplitHostPort(host)
		if err != nil {
			return nil, fmt.Errorf("failed to split host:port in '%s', err: %w", host, err)
		}
		if h != "" {
			hosts = append(hosts, h)
		}
		if p != "" {
			ports = append(ports, p)
		}
	}
	if len(hosts) > 0 {
		settings["host"] = strings.Join(hosts, ",")
	}
	if len(ports) > 0 {
		settings["port"] = strings.Join(ports, ",")
	}

	database := strings.TrimLeft(url.Path, "/")
	if database != "" {
		settings["database"] = database
	}

	nameMap := map[string]string{
		"dbname": "database",
	}

	for k, v := range url.Query() {
		if k2, present := nameMap[k]; present {
			k = k2
		}

		settings[k] = v[0]
	}

	return settings, nil
}

func parsePostgresDSNSettings(s string) (map[string]string, error) {
	settings := make(map[string]string)

	nameMap := map[string]string{
		"dbname": "database",
	}

	for len(s) > 0 {
		var key, val string
		eqIdx := strings.IndexRune(s, '=')
		if eqIdx < 0 {
			return nil, errors.New("invalid dsn")
		}

		key = strings.Trim(s[:eqIdx], " \t\n\r\v\f")
		s = strings.TrimLeft(s[eqIdx+1:], " \t\n\r\v\f")
		if len(s) == 0 {
		} else if s[0] != '\'' {
			end := 0
			for ; end < len(s); end++ {
				if asciiSpace[s[end]] == 1 {
					break
				}
				if s[end] == '\\' {
					end++
					if end == len(s) {
						return nil, errors.New("invalid backslash")
					}
				}
			}
			val = strings.Replace(strings.Replace(s[:end], "\\\\", "\\", -1), "\\'", "'", -1)
			if end == len(s) {
				s = ""
			} else {
				s = s[end+1:]
			}
		} else { // quoted string
			s = s[1:]
			end := 0
			for ; end < len(s); end++ {
				if s[end] == '\'' {
					break
				}
				if s[end] == '\\' {
					end++
				}
			}
			if end == len(s) {
				return nil, errors.New("unterminated quoted string in connection info string")
			}
			val = strings.Replace(strings.Replace(s[:end], "\\\\", "\\", -1), "\\'", "'", -1)
			if end == len(s) {
				s = ""
			} else {
				s = s[end+1:]
			}
		}

		if k, ok := nameMap[key]; ok {
			key = k
		}

		if key == "" {
			return nil, errors.New("invalid dsn")
		}

		settings[key] = val
	}

	return settings, nil
}

func isIPOnly(host string) bool {
	return net.ParseIP(strings.Trim(host, "[]")) != nil || !strings.Contains(host, ":")
}
