// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cassandra

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/mitchellh/mapstructure"
)

// cassandraConnectionProducer implements ConnectionProducer and provides an
// interface for cassandra databases to make connections.
type cassandraConnectionProducer struct {
	Hosts              string      `json:"hosts" structs:"hosts" mapstructure:"hosts"`
	Port               int         `json:"port" structs:"port" mapstructure:"port"`
	Username           string      `json:"username" structs:"username" mapstructure:"username"`
	Password           string      `json:"password" structs:"password" mapstructure:"password"`
	TLS                bool        `json:"tls" structs:"tls" mapstructure:"tls"`
	InsecureTLS        bool        `json:"insecure_tls" structs:"insecure_tls" mapstructure:"insecure_tls"`
	TLSServerName      string      `json:"tls_server_name" structs:"tls_server_name" mapstructure:"tls_server_name"`
	ProtocolVersion    int         `json:"protocol_version" structs:"protocol_version" mapstructure:"protocol_version"`
	ConnectTimeoutRaw  interface{} `json:"connect_timeout" structs:"connect_timeout" mapstructure:"connect_timeout"`
	SocketKeepAliveRaw interface{} `json:"socket_keep_alive" structs:"socket_keep_alive" mapstructure:"socket_keep_alive"`
	TLSMinVersion      string      `json:"tls_min_version" structs:"tls_min_version" mapstructure:"tls_min_version"`
	Consistency        string      `json:"consistency" structs:"consistency" mapstructure:"consistency"`
	LocalDatacenter    string      `json:"local_datacenter" structs:"local_datacenter" mapstructure:"local_datacenter"`
	PemBundle          string      `json:"pem_bundle" structs:"pem_bundle" mapstructure:"pem_bundle"`
	PemJSON            string      `json:"pem_json" structs:"pem_json" mapstructure:"pem_json"`
	SkipVerification   bool        `json:"skip_verification" structs:"skip_verification" mapstructure:"skip_verification"`

	connectTimeout  time.Duration
	socketKeepAlive time.Duration
	sslOpts         *gocql.SslOptions
	rawConfig       map[string]interface{}

	Initialized bool
	Type        string
	session     *gocql.Session
	sync.Mutex
}

func (c *cassandraConnectionProducer) Initialize(ctx context.Context, req dbplugin.InitializeRequest) error {
	c.Lock()
	defer c.Unlock()

	c.rawConfig = req.Config

	err := mapstructure.WeakDecode(req.Config, c)
	if err != nil {
		return err
	}

	if c.ConnectTimeoutRaw == nil {
		c.ConnectTimeoutRaw = "5s"
	}
	c.connectTimeout, err = parseutil.ParseDurationSecond(c.ConnectTimeoutRaw)
	if err != nil {
		return fmt.Errorf("invalid connect_timeout: %w", err)
	}

	if c.SocketKeepAliveRaw == nil {
		c.SocketKeepAliveRaw = "0s"
	}
	c.socketKeepAlive, err = parseutil.ParseDurationSecond(c.SocketKeepAliveRaw)
	if err != nil {
		return fmt.Errorf("invalid socket_keep_alive: %w", err)
	}

	switch {
	case len(c.Hosts) == 0:
		return fmt.Errorf("hosts cannot be empty")
	case len(c.Username) == 0:
		return fmt.Errorf("username cannot be empty")
	case len(c.Password) == 0:
		return fmt.Errorf("password cannot be empty")
	case len(c.PemJSON) > 0 && len(c.PemBundle) > 0:
		return fmt.Errorf("cannot specify both pem_json and pem_bundle")
	}

	var tlsMinVersion uint16 = tls.VersionTLS12
	if c.TLSMinVersion != "" {
		ver, exists := tlsutil.TLSLookup[c.TLSMinVersion]
		if !exists {
			return fmt.Errorf("unrecognized TLS version [%s]", c.TLSMinVersion)
		}
		tlsMinVersion = ver
	}

	switch {
	case len(c.PemJSON) != 0:
		cfg, err := jsonBundleToTLSConfig(c.PemJSON, tlsMinVersion, c.TLSServerName, c.InsecureTLS)
		if err != nil {
			return fmt.Errorf("failed to parse pem_json: %w", err)
		}
		c.sslOpts = &gocql.SslOptions{
			Config:                 cfg,
			EnableHostVerification: !cfg.InsecureSkipVerify,
		}
		c.TLS = true

	case len(c.PemBundle) != 0:
		cfg, err := pemBundleToTLSConfig(c.PemBundle, tlsMinVersion, c.TLSServerName, c.InsecureTLS)
		if err != nil {
			return fmt.Errorf("failed to parse pem_bundle: %w", err)
		}
		c.sslOpts = &gocql.SslOptions{
			Config:                 cfg,
			EnableHostVerification: !cfg.InsecureSkipVerify,
		}
		c.TLS = true

	case c.InsecureTLS:
		c.sslOpts = &gocql.SslOptions{
			EnableHostVerification: !c.InsecureTLS,
		}
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	if req.VerifyConnection {
		if _, err := c.Connection(ctx); err != nil {
			return fmt.Errorf("error verifying connection: %w", err)
		}
	}

	return nil
}

func (c *cassandraConnectionProducer) Connection(ctx context.Context) (interface{}, error) {
	if !c.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	// If we already have a DB, return it
	if c.session != nil && !c.session.Closed() {
		return c.session, nil
	}

	session, err := c.createSession(ctx)
	if err != nil {
		return nil, err
	}

	//  Store the session in backend for reuse
	c.session = session

	return session, nil
}

func (c *cassandraConnectionProducer) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.session != nil {
		c.session.Close()
	}

	c.session = nil

	return nil
}

func (c *cassandraConnectionProducer) createSession(ctx context.Context) (*gocql.Session, error) {
	hosts := strings.Split(c.Hosts, ",")
	clusterConfig := gocql.NewCluster(hosts...)
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Password,
	}

	if c.Port != 0 {
		clusterConfig.Port = c.Port
	}

	clusterConfig.ProtoVersion = c.ProtocolVersion
	if clusterConfig.ProtoVersion == 0 {
		clusterConfig.ProtoVersion = 2
	}

	clusterConfig.Timeout = c.connectTimeout
	clusterConfig.ConnectTimeout = c.connectTimeout
	clusterConfig.SocketKeepalive = c.socketKeepAlive
	clusterConfig.SslOpts = c.sslOpts

	if c.LocalDatacenter != "" {
		clusterConfig.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy(c.LocalDatacenter)
	}

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	if c.Consistency != "" {
		consistencyValue, err := gocql.ParseConsistencyWrapper(c.Consistency)
		if err != nil {
			session.Close()
			return nil, err
		}

		session.SetConsistency(consistencyValue)
	}

	if !c.SkipVerification {
		err = session.Query(`LIST ALL`).WithContext(ctx).Exec()
		if err != nil && len(c.Username) != 0 && strings.Contains(err.Error(), "not authorized") {
			rowNum := session.Query(dbutil.QueryHelper(`LIST CREATE ON ALL ROLES OF '{{username}}';`, map[string]string{
				"username": c.Username,
			})).Iter().NumRows()

			if rowNum < 1 {
				session.Close()
				return nil, fmt.Errorf("error validating connection info: No role create permissions found, previous error: %w", err)
			}
		} else if err != nil {
			session.Close()
			return nil, fmt.Errorf("error validating connection info: %w", err)
		}
	}

	return session, nil
}

func (c *cassandraConnectionProducer) secretValues() map[string]string {
	return map[string]string{
		c.Password:  "[password]",
		c.PemBundle: "[pem_bundle]",
		c.PemJSON:   "[pem_json]",
	}
}
