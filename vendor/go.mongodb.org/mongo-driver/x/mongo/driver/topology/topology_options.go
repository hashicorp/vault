// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
)

// Option is a configuration option for a topology.
type Option func(*config) error

type config struct {
	mode                   MonitorMode
	replicaSetName         string
	seedList               []string
	serverOpts             []ServerOption
	cs                     connstring.ConnString // This must not be used for any logic in topology.Topology.
	uri                    string
	serverSelectionTimeout time.Duration
}

func newConfig(opts ...Option) (*config, error) {
	cfg := &config{
		seedList:               []string{"localhost:27017"},
		serverSelectionTimeout: 30 * time.Second,
	}

	for _, opt := range opts {
		err := opt(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// WithConnString configures the topology using the connection string.
func WithConnString(fn func(connstring.ConnString) connstring.ConnString) Option {
	return func(c *config) error {
		cs := fn(c.cs)
		c.cs = cs

		if cs.ServerSelectionTimeoutSet {
			c.serverSelectionTimeout = cs.ServerSelectionTimeout
		}

		var connOpts []ConnectionOption

		if cs.AppName != "" {
			c.serverOpts = append(c.serverOpts, WithServerAppName(func(string) string { return cs.AppName }))
		}

		if cs.Connect == connstring.SingleConnect || (cs.DirectConnectionSet && cs.DirectConnection) {
			c.mode = SingleMode
		}

		c.seedList = cs.Hosts

		if cs.ConnectTimeout > 0 {
			c.serverOpts = append(c.serverOpts, WithHeartbeatTimeout(func(time.Duration) time.Duration { return cs.ConnectTimeout }))
			connOpts = append(connOpts, WithConnectTimeout(func(time.Duration) time.Duration { return cs.ConnectTimeout }))
		}

		if cs.SocketTimeoutSet {
			connOpts = append(
				connOpts,
				WithReadTimeout(func(time.Duration) time.Duration { return cs.SocketTimeout }),
				WithWriteTimeout(func(time.Duration) time.Duration { return cs.SocketTimeout }),
			)
		}

		if cs.HeartbeatInterval > 0 {
			c.serverOpts = append(c.serverOpts, WithHeartbeatInterval(func(time.Duration) time.Duration { return cs.HeartbeatInterval }))
		}

		if cs.MaxConnIdleTime > 0 {
			connOpts = append(connOpts, WithIdleTimeout(func(time.Duration) time.Duration { return cs.MaxConnIdleTime }))
		}

		if cs.MaxPoolSizeSet {
			c.serverOpts = append(c.serverOpts, WithMaxConnections(func(uint64) uint64 { return cs.MaxPoolSize }))
		}

		if cs.MinPoolSizeSet {
			c.serverOpts = append(c.serverOpts, WithMinConnections(func(u uint64) uint64 { return cs.MinPoolSize }))
		}

		if cs.ReplicaSet != "" {
			c.replicaSetName = cs.ReplicaSet
		}

		var x509Username string
		if cs.SSL {
			tlsConfig := new(tls.Config)

			if cs.SSLCaFileSet {
				err := addCACertFromFile(tlsConfig, cs.SSLCaFile)
				if err != nil {
					return err
				}
			}

			if cs.SSLInsecure {
				tlsConfig.InsecureSkipVerify = true
			}

			if cs.SSLClientCertificateKeyFileSet {
				var keyPasswd string
				if cs.SSLClientCertificateKeyPasswordSet && cs.SSLClientCertificateKeyPassword != nil {
					keyPasswd = cs.SSLClientCertificateKeyPassword()
				}
				s, err := addClientCertFromFile(tlsConfig, cs.SSLClientCertificateKeyFile, keyPasswd)
				if err != nil {
					return err
				}

				// The Go x509 package gives the subject with the pairs in reverse order that we want.
				pairs := strings.Split(s, ",")
				b := bytes.NewBufferString("")

				for i := len(pairs) - 1; i >= 0; i-- {
					b.WriteString(pairs[i])

					if i > 0 {
						b.WriteString(",")
					}
				}

				x509Username = b.String()
			}

			connOpts = append(connOpts, WithTLSConfig(func(*tls.Config) *tls.Config { return tlsConfig }))
		}

		if cs.Username != "" || cs.AuthMechanism == auth.MongoDBX509 || cs.AuthMechanism == auth.GSSAPI {
			cred := &auth.Cred{
				Source:      "admin",
				Username:    cs.Username,
				Password:    cs.Password,
				PasswordSet: cs.PasswordSet,
				Props:       cs.AuthMechanismProperties,
			}

			if cs.AuthSource != "" {
				cred.Source = cs.AuthSource
			} else {
				switch cs.AuthMechanism {
				case auth.MongoDBX509:
					if cred.Username == "" {
						cred.Username = x509Username
					}
					fallthrough
				case auth.GSSAPI, auth.PLAIN:
					cred.Source = "$external"
				default:
					cred.Source = cs.Database
				}
			}

			authenticator, err := auth.CreateAuthenticator(cs.AuthMechanism, cred)
			if err != nil {
				return err
			}

			connOpts = append(connOpts, WithHandshaker(func(h Handshaker) Handshaker {
				options := &auth.HandshakeOptions{
					AppName:       cs.AppName,
					Authenticator: authenticator,
					Compressors:   cs.Compressors,
				}
				if cs.AuthMechanism == "" {
					// Required for SASL mechanism negotiation during handshake
					options.DBUser = cred.Source + "." + cred.Username
				}
				return auth.Handshaker(h, options)
			}))
		} else {
			// We need to add a non-auth Handshaker to the connection options
			connOpts = append(connOpts, WithHandshaker(func(h driver.Handshaker) driver.Handshaker {
				return operation.NewIsMaster().AppName(cs.AppName).Compressors(cs.Compressors)
			}))
		}

		if len(cs.Compressors) > 0 {
			connOpts = append(connOpts, WithCompressors(func(compressors []string) []string {
				return append(compressors, cs.Compressors...)
			}))

			for _, comp := range cs.Compressors {
				switch comp {
				case "zlib":
					connOpts = append(connOpts, WithZlibLevel(func(level *int) *int {
						return &cs.ZlibLevel
					}))
				case "zstd":
					connOpts = append(connOpts, WithZstdLevel(func(level *int) *int {
						return &cs.ZstdLevel
					}))
				}
			}

			c.serverOpts = append(c.serverOpts, WithCompressionOptions(func(opts ...string) []string {
				return append(opts, cs.Compressors...)
			}))
		}

		if len(connOpts) > 0 {
			c.serverOpts = append(c.serverOpts, WithConnectionOptions(func(opts ...ConnectionOption) []ConnectionOption {
				return append(opts, connOpts...)
			}))
		}

		return nil
	}
}

// WithMode configures the topology's monitor mode.
func WithMode(fn func(MonitorMode) MonitorMode) Option {
	return func(cfg *config) error {
		cfg.mode = fn(cfg.mode)
		return nil
	}
}

// WithReplicaSetName configures the topology's default replica set name.
func WithReplicaSetName(fn func(string) string) Option {
	return func(cfg *config) error {
		cfg.replicaSetName = fn(cfg.replicaSetName)
		return nil
	}
}

// WithSeedList configures a topology's seed list.
func WithSeedList(fn func(...string) []string) Option {
	return func(cfg *config) error {
		cfg.seedList = fn(cfg.seedList...)
		return nil
	}
}

// WithServerOptions configures a topology's server options for when a new server
// needs to be created.
func WithServerOptions(fn func(...ServerOption) []ServerOption) Option {
	return func(cfg *config) error {
		cfg.serverOpts = fn(cfg.serverOpts...)
		return nil
	}
}

// WithServerSelectionTimeout configures a topology's server selection timeout.
// A server selection timeout of 0 means there is no timeout for server selection.
func WithServerSelectionTimeout(fn func(time.Duration) time.Duration) Option {
	return func(cfg *config) error {
		cfg.serverSelectionTimeout = fn(cfg.serverSelectionTimeout)
		return nil
	}
}

// WithURI specifies the URI that was used to create the topology.
func WithURI(fn func(string) string) Option {
	return func(cfg *config) error {
		cfg.uri = fn(cfg.uri)
		return nil
	}
}

// addCACertFromFile adds a root CA certificate to the configuration given a path
// to the containing file.
func addCACertFromFile(cfg *tls.Config, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	certBytes, err := loadCert(data)
	if err != nil {
		return err
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return err
	}

	if cfg.RootCAs == nil {
		cfg.RootCAs = x509.NewCertPool()
	}

	cfg.RootCAs.AddCert(cert)

	return nil
}

func loadCert(data []byte) ([]byte, error) {
	var certBlock *pem.Block

	for certBlock == nil {
		if data == nil || len(data) == 0 {
			return nil, errors.New(".pem file must have both a CERTIFICATE and an RSA PRIVATE KEY section")
		}

		block, rest := pem.Decode(data)
		if block == nil {
			return nil, errors.New("invalid .pem file")
		}

		switch block.Type {
		case "CERTIFICATE":
			if certBlock != nil {
				return nil, errors.New("multiple CERTIFICATE sections in .pem file")
			}

			certBlock = block
		}

		data = rest
	}

	return certBlock.Bytes, nil
}

// addClientCertFromFile adds a client certificate to the configuration given a path to the
// containing file and returns the certificate's subject name.
func addClientCertFromFile(cfg *tls.Config, clientFile, keyPasswd string) (string, error) {
	data, err := ioutil.ReadFile(clientFile)
	if err != nil {
		return "", err
	}

	var currentBlock *pem.Block
	var certBlock, certDecodedBlock, keyBlock []byte

	remaining := data
	start := 0
	for {
		currentBlock, remaining = pem.Decode(remaining)
		if currentBlock == nil {
			break
		}

		if currentBlock.Type == "CERTIFICATE" {
			certBlock = data[start : len(data)-len(remaining)]
			certDecodedBlock = currentBlock.Bytes
			start += len(certBlock)
		} else if strings.HasSuffix(currentBlock.Type, "PRIVATE KEY") {
			if keyPasswd != "" && x509.IsEncryptedPEMBlock(currentBlock) {
				var encoded bytes.Buffer
				buf, err := x509.DecryptPEMBlock(currentBlock, []byte(keyPasswd))
				if err != nil {
					return "", err
				}

				pem.Encode(&encoded, &pem.Block{Type: currentBlock.Type, Bytes: buf})
				keyBlock = encoded.Bytes()
				start = len(data) - len(remaining)
			} else {
				keyBlock = data[start : len(data)-len(remaining)]
				start += len(keyBlock)
			}
		}
	}
	if len(certBlock) == 0 {
		return "", fmt.Errorf("failed to find CERTIFICATE")
	}
	if len(keyBlock) == 0 {
		return "", fmt.Errorf("failed to find PRIVATE KEY")
	}

	cert, err := tls.X509KeyPair(certBlock, keyBlock)
	if err != nil {
		return "", err
	}

	cfg.Certificates = append(cfg.Certificates, cert)

	// The documentation for the tls.X509KeyPair indicates that the Leaf certificate is not
	// retained.
	crt, err := x509.ParseCertificate(certDecodedBlock)
	if err != nil {
		return "", err
	}

	return x509CertSubject(crt), nil
}
