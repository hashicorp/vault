// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dependency

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	rootcerts "github.com/hashicorp/go-rootcerts"
	nomadapi "github.com/hashicorp/nomad/api"
	vaultapi "github.com/hashicorp/vault/api"
	vaultkubernetesauth "github.com/hashicorp/vault/api/auth/kubernetes"
)

// ClientSet is a collection of clients that dependencies use to communicate
// with remote services like Consul or Vault.
type ClientSet struct {
	sync.RWMutex

	vault  *vaultClient
	consul *consulClient
	nomad  *nomadClient
}

// consulClient is a wrapper around a real Consul API client.
type consulClient struct {
	client    *consulapi.Client
	transport *http.Transport
}

// vaultClient is a wrapper around a real Vault API client.
type vaultClient struct {
	client     *vaultapi.Client
	httpClient *http.Client
}

// nomadClient is a wrapper around a real Nomad API client.
type nomadClient struct {
	client     *nomadapi.Client
	httpClient *http.Client
}

// TransportDialer is an interface that allows passing a custom dialer function
// to an HTTP client's transport config
type TransportDialer interface {
	// Dial is intended to match https://pkg.go.dev/net#Dialer.Dial
	Dial(network, address string) (net.Conn, error)

	// DialContext is intended to match https://pkg.go.dev/net#Dialer.DialContext
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// CreateConsulClientInput is used as input to the CreateConsulClient function.
type CreateConsulClientInput struct {
	Address      string
	Namespace    string
	Token        string
	TokenFile    string
	AuthEnabled  bool
	AuthUsername string
	AuthPassword string
	SSLEnabled   bool
	SSLVerify    bool
	SSLCert      string
	SSLKey       string
	SSLCACert    string
	SSLCAPath    string
	ServerName   string

	TransportDialKeepAlive       time.Duration
	TransportDialTimeout         time.Duration
	TransportDisableKeepAlives   bool
	TransportIdleConnTimeout     time.Duration
	TransportMaxIdleConns        int
	TransportMaxIdleConnsPerHost int
	TransportTLSHandshakeTimeout time.Duration
}

// CreateVaultClientInput is used as input to the CreateVaultClient function.
type CreateVaultClientInput struct {
	Address         string
	Namespace       string
	Token           string
	UnwrapToken     bool
	SSLEnabled      bool
	SSLVerify       bool
	SSLCert         string
	SSLKey          string
	SSLCACert       string
	SSLCACertBytes  string
	SSLCAPath       string
	ServerName      string
	ClientUserAgent string

	K8SAuthRoleName            string
	K8SServiceAccountTokenPath string
	K8SServiceAccountToken     string
	K8SServiceMountPath        string

	TransportCustomDialer        TransportDialer
	TransportDialKeepAlive       time.Duration
	TransportDialTimeout         time.Duration
	TransportDisableKeepAlives   bool
	TransportIdleConnTimeout     time.Duration
	TransportMaxIdleConns        int
	TransportMaxIdleConnsPerHost int
	TransportMaxConnsPerHost     int
	TransportTLSHandshakeTimeout time.Duration
}

// CreateNomadClientInput is used as input to the CreateNomadClient function.
type CreateNomadClientInput struct {
	Address      string
	Namespace    string
	Token        string
	AuthUsername string
	AuthPassword string
	SSLEnabled   bool
	SSLVerify    bool
	SSLCert      string
	SSLKey       string
	SSLCACert    string
	SSLCAPath    string
	ServerName   string

	TransportCustomDialer        TransportDialer
	TransportDialKeepAlive       time.Duration
	TransportDialTimeout         time.Duration
	TransportDisableKeepAlives   bool
	TransportIdleConnTimeout     time.Duration
	TransportMaxIdleConns        int
	TransportMaxIdleConnsPerHost int
	TransportTLSHandshakeTimeout time.Duration
}

// NewClientSet creates a new client set that is ready to accept clients.
func NewClientSet() *ClientSet {
	return &ClientSet{}
}

// CreateConsulClient creates a new Consul API client from the given input.
func (c *ClientSet) CreateConsulClient(i *CreateConsulClientInput) error {
	consulConfig := consulapi.DefaultConfig()

	if i.Address != "" {
		consulConfig.Address = i.Address
	}

	if i.Namespace != "" {
		consulConfig.Namespace = i.Namespace
	}

	if i.Token != "" {
		consulConfig.Token = i.Token
	}

	if i.TokenFile != "" {
		consulConfig.TokenFile = i.TokenFile
	}

	if i.AuthEnabled {
		consulConfig.HttpAuth = &consulapi.HttpBasicAuth{
			Username: i.AuthUsername,
			Password: i.AuthPassword,
		}
	}

	// This transport will attempt to keep connections open to the Consul server.
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   i.TransportDialTimeout,
			KeepAlive: i.TransportDialKeepAlive,
		}).Dial,
		DisableKeepAlives:   i.TransportDisableKeepAlives,
		MaxIdleConns:        i.TransportMaxIdleConns,
		IdleConnTimeout:     i.TransportIdleConnTimeout,
		MaxIdleConnsPerHost: i.TransportMaxIdleConnsPerHost,
		TLSHandshakeTimeout: i.TransportTLSHandshakeTimeout,
	}

	// Configure SSL
	if i.SSLEnabled {
		consulConfig.Scheme = "https"

		var tlsConfig tls.Config

		// Custom certificate or certificate and key
		if i.SSLCert != "" && i.SSLKey != "" {
			cert, err := tls.LoadX509KeyPair(i.SSLCert, i.SSLKey)
			if err != nil {
				return fmt.Errorf("client set: consul: %s", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		} else if i.SSLCert != "" {
			cert, err := tls.LoadX509KeyPair(i.SSLCert, i.SSLCert)
			if err != nil {
				return fmt.Errorf("client set: consul: %s", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		// Custom CA certificate
		if i.SSLCACert != "" || i.SSLCAPath != "" {
			rootConfig := &rootcerts.Config{
				CAFile: i.SSLCACert,
				CAPath: i.SSLCAPath,
			}
			if err := rootcerts.ConfigureTLS(&tlsConfig, rootConfig); err != nil {
				return fmt.Errorf("client set: consul configuring TLS failed: %s", err)
			}
		}

		// Construct all the certificates now
		tlsConfig.BuildNameToCertificate()

		// SSL verification
		if i.ServerName != "" {
			tlsConfig.ServerName = i.ServerName
			tlsConfig.InsecureSkipVerify = false
		}
		if !i.SSLVerify {
			log.Printf("[WARN] (clients) disabling consul SSL verification")
			tlsConfig.InsecureSkipVerify = true
		}

		// Save the TLS config on our transport
		transport.TLSClientConfig = &tlsConfig
	}

	// Setup the new transport
	consulConfig.Transport = transport

	// Create the API client
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		return fmt.Errorf("client set: consul: %s", err)
	}

	// Save the data on ourselves
	c.Lock()
	c.consul = &consulClient{
		client:    client,
		transport: transport,
	}
	c.Unlock()

	return nil
}

func (c *ClientSet) CreateVaultClient(i *CreateVaultClientInput) error {
	vaultConfig := vaultapi.DefaultConfig()

	if i.Address != "" {
		vaultConfig.Address = i.Address
	}

	// This transport will attempt to keep connections open to the Vault server.
	var dialer TransportDialer
	dialer = &net.Dialer{
		Timeout:   i.TransportDialTimeout,
		KeepAlive: i.TransportDialKeepAlive,
	}

	if i.TransportCustomDialer != nil {
		dialer = i.TransportCustomDialer
	}

	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		Dial:                dialer.Dial,
		DisableKeepAlives:   i.TransportDisableKeepAlives,
		MaxIdleConns:        i.TransportMaxIdleConns,
		IdleConnTimeout:     i.TransportIdleConnTimeout,
		MaxIdleConnsPerHost: i.TransportMaxIdleConnsPerHost,
		MaxConnsPerHost:     i.TransportMaxConnsPerHost,
		TLSHandshakeTimeout: i.TransportTLSHandshakeTimeout,
	}

	// Configure SSL
	if i.SSLEnabled {
		var tlsConfig tls.Config

		// Custom certificate or certificate and key
		if i.SSLCert != "" && i.SSLKey != "" {
			cert, err := tls.LoadX509KeyPair(i.SSLCert, i.SSLKey)
			if err != nil {
				return fmt.Errorf("client set: vault: %s", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		} else if i.SSLCert != "" {
			cert, err := tls.LoadX509KeyPair(i.SSLCert, i.SSLCert)
			if err != nil {
				return fmt.Errorf("client set: vault: %s", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		// Custom CA certificate
		if i.SSLCACert != "" || i.SSLCAPath != "" || i.SSLCACertBytes != "" {
			rootConfig := &rootcerts.Config{
				CAFile:        i.SSLCACert,
				CACertificate: []byte(i.SSLCACertBytes),
				CAPath:        i.SSLCAPath,
			}
			if err := rootcerts.ConfigureTLS(&tlsConfig, rootConfig); err != nil {
				return fmt.Errorf("client set: vault configuring TLS failed: %s", err)
			}
		}

		// Construct all the certificates now
		tlsConfig.BuildNameToCertificate()

		// SSL verification
		if i.ServerName != "" {
			tlsConfig.ServerName = i.ServerName
			tlsConfig.InsecureSkipVerify = false
		}
		if !i.SSLVerify {
			log.Printf("[WARN] (clients) disabling vault SSL verification")
			tlsConfig.InsecureSkipVerify = true
		}

		// Save the TLS config on our transport
		transport.TLSClientConfig = &tlsConfig
	}

	// Setup the new transport
	vaultConfig.HttpClient.Transport = transport

	// Create the client
	client, err := vaultapi.NewClient(vaultConfig)
	if err != nil {
		return fmt.Errorf("client set: vault: %s", err)
	}

	if i.ClientUserAgent != "" {
		client.SetCloneHeaders(true)
		client.AddHeader("User-Agent", i.ClientUserAgent)
	}

	// Set the namespace if given.
	if i.Namespace != "" {
		client.SetNamespace(i.Namespace)
	}

	// Set token using k8s auth method.
	if i.K8SAuthRoleName != "" && i.Token == "" {
		err = prepareK8SServiceTokenAuth(i, client)
		if err != nil {
			return fmt.Errorf("client set: vault: %w", err)
		}
	}

	if i.Token != "" {
		client.SetToken(i.Token)
	}

	// Save the data on ourselves
	c.Lock()
	c.vault = &vaultClient{
		client:     client,
		httpClient: vaultConfig.HttpClient,
	}
	c.Unlock()

	return nil
}

// CreateNomadClient creates a new Nomad API client from the given input.
func (c *ClientSet) CreateNomadClient(i *CreateNomadClientInput) error {
	conf := nomadapi.DefaultConfig()

	if i.Address != "" {
		conf.Address = i.Address
	}

	if i.Namespace != "" {
		conf.Namespace = i.Namespace
	}

	if i.Token != "" {
		conf.SecretID = i.Token
	}

	if i.AuthUsername != "" || i.AuthPassword != "" {
		conf.HttpAuth = &nomadapi.HttpBasicAuth{
			Username: i.AuthUsername,
			Password: i.AuthPassword,
		}
	}

	// This transport will attempt to keep connections open to the Nomad server.
	var dialer TransportDialer
	dialer = &net.Dialer{
		Timeout:   i.TransportDialTimeout,
		KeepAlive: i.TransportDialKeepAlive,
	}

	if i.TransportCustomDialer != nil {
		dialer = i.TransportCustomDialer
	}

	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		Dial:                dialer.Dial,
		DisableKeepAlives:   i.TransportDisableKeepAlives,
		MaxIdleConns:        i.TransportMaxIdleConns,
		IdleConnTimeout:     i.TransportIdleConnTimeout,
		MaxIdleConnsPerHost: i.TransportMaxIdleConnsPerHost,
		TLSHandshakeTimeout: i.TransportTLSHandshakeTimeout,
	}

	// Configure SSL
	if i.SSLEnabled {
		var tlsConfig tls.Config

		// Custom certificate or certificate and key
		if i.SSLCert != "" && i.SSLKey != "" {
			cert, err := tls.LoadX509KeyPair(i.SSLCert, i.SSLKey)
			if err != nil {
				return fmt.Errorf("client set: nomad: %s", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		} else if i.SSLCert != "" {
			cert, err := tls.LoadX509KeyPair(i.SSLCert, i.SSLCert)
			if err != nil {
				return fmt.Errorf("client set: nomad: %s", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		// Custom CA certificate
		if i.SSLCACert != "" || i.SSLCAPath != "" {
			rootConfig := &rootcerts.Config{
				CAFile: i.SSLCACert,
				CAPath: i.SSLCAPath,
			}
			if err := rootcerts.ConfigureTLS(&tlsConfig, rootConfig); err != nil {
				return fmt.Errorf("client set: nomad configuring TLS failed: %s", err)
			}
		}

		// Construct all the certificates now
		tlsConfig.BuildNameToCertificate()

		// SSL verification
		if i.ServerName != "" {
			tlsConfig.ServerName = i.ServerName
			tlsConfig.InsecureSkipVerify = false
		}
		if !i.SSLVerify {
			log.Printf("[WARN] (clients) disabling nomad SSL verification")
			tlsConfig.InsecureSkipVerify = true
		}

		// Save the TLS config on our transport
		transport.TLSClientConfig = &tlsConfig
	}

	conf.HttpClient = &http.Client{
		Transport: transport,
	}

	// Create the API client
	client, err := nomadapi.NewClient(conf)
	if err != nil {
		return fmt.Errorf("client set: nomad: %w", err)
	}

	// Save the data on ourselves
	c.Lock()
	c.nomad = &nomadClient{
		client:     client,
		httpClient: conf.HttpClient,
	}
	c.Unlock()

	return nil
}

// Consul returns the Consul client for this set.
func (c *ClientSet) Consul() *consulapi.Client {
	c.RLock()
	defer c.RUnlock()
	return c.consul.client
}

// Vault returns the Vault client for this set.
func (c *ClientSet) Vault() *vaultapi.Client {
	c.RLock()
	defer c.RUnlock()
	return c.vault.client
}

// Nomad returns the Nomad client for this set.
func (c *ClientSet) Nomad() *nomadapi.Client {
	c.RLock()
	defer c.RUnlock()
	return c.nomad.client
}

// Stop closes all idle connections for any attached clients.
func (c *ClientSet) Stop() {
	c.Lock()
	defer c.Unlock()

	if c.consul != nil {
		c.consul.transport.CloseIdleConnections()
	}

	if c.vault != nil {
		c.vault.httpClient.Transport.(*http.Transport).CloseIdleConnections()
	}

	if c.nomad != nil {
		c.nomad.httpClient.Transport.(*http.Transport).CloseIdleConnections()
	}
}

func prepareK8SServiceTokenAuth(
	i *CreateVaultClientInput,
	client *vaultapi.Client,
) (err error) {
	opts := make([]vaultkubernetesauth.LoginOption, 0, 2)

	switch {
	case i.K8SServiceAccountToken != "":
		opts = append(opts, vaultkubernetesauth.WithServiceAccountToken(
			i.K8SServiceAccountToken,
		))
	case i.K8SServiceAccountTokenPath != "":
		opts = append(opts, vaultkubernetesauth.WithServiceAccountTokenPath(
			i.K8SServiceAccountTokenPath,
		))
	default:
		// The Kubernetes service account token JWT will be retrieved
		// from /run/secrets/kubernetes.io/serviceaccount/token.
	}

	if i.K8SServiceMountPath != "" {
		opts = append(opts, vaultkubernetesauth.WithMountPath(
			i.K8SServiceMountPath,
		))
	}

	k8sAuth, err := vaultkubernetesauth.NewKubernetesAuth(i.K8SAuthRoleName, opts...)
	if err != nil {
		return fmt.Errorf("k8s auth: new kubernetes auth: %w", err)
	}

	ctx := context.TODO()
	sec, err := client.Auth().Login(ctx, k8sAuth)
	if err != nil {
		return fmt.Errorf("k8s auth: login: %w", err)
	}

	i.Token = sec.Auth.ClientToken

	return nil
}
