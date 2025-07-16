// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package influxdb

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/mitchellh/mapstructure"
)

// influxdbConnectionProducer implements ConnectionProducer and provides an
// interface for influxdb databases to make connections.
type influxdbConnectionProducer struct {
	Host              string      `json:"host" structs:"host" mapstructure:"host"`
	Username          string      `json:"username" structs:"username" mapstructure:"username"`
	Password          string      `json:"password" structs:"password" mapstructure:"password"`
	Port              string      `json:"port" structs:"port" mapstructure:"port"` // default to 8086
	TLS               bool        `json:"tls" structs:"tls" mapstructure:"tls"`
	InsecureTLS       bool        `json:"insecure_tls" structs:"insecure_tls" mapstructure:"insecure_tls"`
	ConnectTimeoutRaw interface{} `json:"connect_timeout" structs:"connect_timeout" mapstructure:"connect_timeout"`
	TLSMinVersion     string      `json:"tls_min_version" structs:"tls_min_version" mapstructure:"tls_min_version"`
	PemBundle         string      `json:"pem_bundle" structs:"pem_bundle" mapstructure:"pem_bundle"`
	PemJSON           string      `json:"pem_json" structs:"pem_json" mapstructure:"pem_json"`

	connectTimeout time.Duration
	certificate    string
	privateKey     string
	issuingCA      string
	rawConfig      map[string]interface{}

	Initialized bool
	Type        string
	client      influx.Client
	sync.Mutex
}

func (i *influxdbConnectionProducer) Initialize(ctx context.Context, req dbplugin.InitializeRequest) (dbplugin.InitializeResponse, error) {
	i.Lock()
	defer i.Unlock()

	i.rawConfig = req.Config

	err := mapstructure.WeakDecode(req.Config, i)
	if err != nil {
		return dbplugin.InitializeResponse{}, err
	}

	if i.ConnectTimeoutRaw == nil {
		i.ConnectTimeoutRaw = "5s"
	}
	if i.Port == "" {
		i.Port = "8086"
	}
	i.connectTimeout, err = parseutil.ParseDurationSecond(i.ConnectTimeoutRaw)
	if err != nil {
		return dbplugin.InitializeResponse{}, fmt.Errorf("invalid connect_timeout: %w", err)
	}

	switch {
	case len(i.Host) == 0:
		return dbplugin.InitializeResponse{}, fmt.Errorf("host cannot be empty")
	case len(i.Username) == 0:
		return dbplugin.InitializeResponse{}, fmt.Errorf("username cannot be empty")
	case len(i.Password) == 0:
		return dbplugin.InitializeResponse{}, fmt.Errorf("password cannot be empty")
	}

	var certBundle *certutil.CertBundle
	var parsedCertBundle *certutil.ParsedCertBundle
	switch {
	case len(i.PemJSON) != 0:
		parsedCertBundle, err = certutil.ParsePKIJSON([]byte(i.PemJSON))
		if err != nil {
			return dbplugin.InitializeResponse{}, fmt.Errorf("could not parse given JSON; it must be in the format of the output of the PKI backend certificate issuing command: %w", err)
		}
		certBundle, err = parsedCertBundle.ToCertBundle()
		if err != nil {
			return dbplugin.InitializeResponse{}, fmt.Errorf("Error marshaling PEM information: %w", err)
		}
		i.certificate = certBundle.Certificate
		i.privateKey = certBundle.PrivateKey
		i.issuingCA = certBundle.IssuingCA
		i.TLS = true

	case len(i.PemBundle) != 0:
		parsedCertBundle, err = certutil.ParsePEMBundle(i.PemBundle)
		if err != nil {
			return dbplugin.InitializeResponse{}, fmt.Errorf("Error parsing the given PEM information: %w", err)
		}
		certBundle, err = parsedCertBundle.ToCertBundle()
		if err != nil {
			return dbplugin.InitializeResponse{}, fmt.Errorf("Error marshaling PEM information: %w", err)
		}
		i.certificate = certBundle.Certificate
		i.privateKey = certBundle.PrivateKey
		i.issuingCA = certBundle.IssuingCA
		i.TLS = true
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	i.Initialized = true

	if req.VerifyConnection {
		if _, err := i.Connection(ctx); err != nil {
			return dbplugin.InitializeResponse{}, fmt.Errorf("error verifying connection: %w", err)
		}
	}

	resp := dbplugin.InitializeResponse{
		Config: req.Config,
	}

	return resp, nil
}

func (i *influxdbConnectionProducer) Connection(_ context.Context) (interface{}, error) {
	if !i.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	// If we already have a DB, return it
	if i.client != nil {
		return i.client, nil
	}

	cli, err := i.createClient()
	if err != nil {
		return nil, err
	}

	//  Store the session in backend for reuse
	i.client = cli

	return cli, nil
}

func (i *influxdbConnectionProducer) Close() error {
	// Grab the write lock
	i.Lock()
	defer i.Unlock()

	if i.client != nil {
		i.client.Close()
	}

	i.client = nil

	return nil
}

func (i *influxdbConnectionProducer) createClient() (influx.Client, error) {
	clientConfig := influx.HTTPConfig{
		Addr:      fmt.Sprintf("http://%s:%s", i.Host, i.Port),
		Username:  i.Username,
		Password:  i.Password,
		UserAgent: "vault-influxdb-plugin",
		Timeout:   i.connectTimeout,
	}

	if i.TLS {
		tlsConfig := &tls.Config{}
		if len(i.certificate) > 0 || len(i.issuingCA) > 0 {
			if len(i.certificate) > 0 && len(i.privateKey) == 0 {
				return nil, fmt.Errorf("found certificate for TLS authentication but no private key")
			}

			certBundle := &certutil.CertBundle{}
			if len(i.certificate) > 0 {
				certBundle.Certificate = i.certificate
				certBundle.PrivateKey = i.privateKey
			}
			if len(i.issuingCA) > 0 {
				certBundle.IssuingCA = i.issuingCA
			}

			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate bundle: %w", err)
			}

			tlsConfig, err = parsedCertBundle.GetTLSConfig(certutil.TLSClient)
			if err != nil || tlsConfig == nil {
				return nil, fmt.Errorf("failed to get TLS configuration: tlsConfig:%#v err:%w", tlsConfig, err)
			}
		}

		tlsConfig.InsecureSkipVerify = i.InsecureTLS

		if i.TLSMinVersion != "" {
			var ok bool
			tlsConfig.MinVersion, ok = tlsutil.TLSLookup[i.TLSMinVersion]
			if !ok {
				return nil, fmt.Errorf("invalid 'tls_min_version' in config")
			}
		} else {
			// MinVersion was not being set earlier. Reset it to
			// zero to gracefully handle upgrades.
			tlsConfig.MinVersion = 0
		}

		clientConfig.TLSConfig = tlsConfig
		clientConfig.Addr = fmt.Sprintf("https://%s:%s", i.Host, i.Port)
	}

	cli, err := influx.NewHTTPClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	// Checking server status
	_, _, err = cli.Ping(i.connectTimeout)
	if err != nil {
		return nil, fmt.Errorf("error checking cluster status: %w", err)
	}

	// verifying infos about the connection
	isAdmin, err := isUserAdmin(cli, i.Username)
	if err != nil {
		return nil, fmt.Errorf("error getting if provided username is admin: %w", err)
	}
	if !isAdmin {
		return nil, fmt.Errorf("the provided user is not an admin of the influxDB server")
	}

	return cli, nil
}

func (i *influxdbConnectionProducer) secretValues() map[string]string {
	return map[string]string{
		i.Password:  "[password]",
		i.PemBundle: "[pem_bundle]",
		i.PemJSON:   "[pem_json]",
	}
}

func isUserAdmin(cli influx.Client, user string) (bool, error) {
	q := influx.NewQuery("SHOW USERS", "", "")
	response, err := cli.Query(q)
	if err != nil {
		return false, err
	}
	if response == nil {
		return false, fmt.Errorf("empty response")
	}
	if response.Error() != nil {
		return false, response.Error()
	}
	for _, res := range response.Results {
		for _, serie := range res.Series {
			for _, val := range serie.Values {
				if val[0].(string) == user && val[1].(bool) {
					return true, nil
				}
			}
		}
	}
	return false, fmt.Errorf("the provided username is not a valid user in the influxdb")
}
