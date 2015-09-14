package api

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/hcl"
	"github.com/mitchellh/mapstructure"
)

const (
	// Default path at which SSH backend will be mounted in Vault server
	SSHAgentDefaultMountPoint = "ssh"

	// Echo request message sent as OTP by the agent
	VerifyEchoRequest = "verify-echo-request"

	// Echo response message sent as a response to OTP matching echo request
	VerifyEchoResponse = "verify-echo-response"
)

// This is a structure representing an SSH agent which can talk to vault server
// in order to verify the OTP entered by the user. It contains the path at which
// SSH backend is mounted at the server.
type SSHAgent struct {
	c          *Client
	MountPoint string
}

// SSHVerifyResp is a structure representing the fields in Vault server's
// response.
type SSHVerifyResponse struct {
	// Usually empty. If the request OTP is echo request message, this will
	// be set to the corresponding echo response message.
	Message string `mapstructure:"message"`

	// Username associated with the OTP
	Username string `mapstructure:"username"`

	// IP associated with the OTP
	IP string `mapstructure:"ip"`
}

// Structure which represents the entries from the agent's configuration file.
type SSHAgentConfig struct {
	VaultAddr       string `hcl:"vault_addr"`
	SSHMountPoint   string `hcl:"ssh_mount_point"`
	CACert          string `hcl:"ca_cert"`
	CAPath          string `hcl:"ca_path"`
	TLSSkipVerify   bool   `hcl:"tls_skip_verify"`
	AllowedCidrList string `hcl:"allowed_cidr_list"`
}

// Returns a HTTP client that uses TLS verification (TLS 1.2) for a given
// certificate pool.
func (c *SSHAgentConfig) TLSClient(certPool *x509.CertPool) *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.TLSSkipVerify,
		MinVersion:         tls.VersionTLS12,
		RootCAs:            certPool,
	}

	client := *http.DefaultClient
	client.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSClientConfig:     tlsConfig,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	// From https://github.com/michiwend/gomusicbrainz/pull/4/files
	defaultRedirectLimit := 30

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > defaultRedirectLimit {
			return fmt.Errorf("%d consecutive requests(redirects)", len(via))
		}
		if len(via) == 0 {
			// No redirects
			return nil
		}
		// mutate the subsequent redirect requests with the first Header
		if token := via[0].Header.Get("X-Vault-Token"); len(token) != 0 {
			req.Header.Set("X-Vault-Token", token)
		}
		return nil
	}

	return &client
}

// Returns a new client for the configuration. This client will be used by the
// SSH agent to communicate with Vault server and verify the OTP entered by user.
// If the configuration supplies Vault SSL certificates, then the client will
// have TLS configured in its transport.
func (c *SSHAgentConfig) NewClient() (*Client, error) {
	// Creating a default client configuration for communicating with vault server.
	clientConfig := DefaultConfig()

	// Pointing the client to the actual address of vault server.
	clientConfig.Address = c.VaultAddr

	// Check if certificates are provided via config file.
	if c.CACert != "" || c.CAPath != "" || c.TLSSkipVerify {
		var certPool *x509.CertPool
		var err error
		if c.CACert != "" {
			certPool, err = loadCACert(c.CACert)
		} else if c.CAPath != "" {
			certPool, err = loadCAPath(c.CAPath)
		}
		if err != nil {
			return nil, err
		}

		// Change the configuration to have an HTTP client with TLS enabled.
		clientConfig.HttpClient = c.TLSClient(certPool)
	}

	// Creating the client object for the given configuration
	client, err := NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Load agent's configuration from the file and populate the corresponding
// in-memory structure.
//
// Vault address is a required parameter.
// Mount point defaults to "ssh".
func LoadSSHAgentConfig(path string) (*SSHAgentConfig, error) {
	var config SSHAgentConfig
	contents, err := ioutil.ReadFile(path)
	if !os.IsNotExist(err) {
		obj, err := hcl.Parse(string(contents))
		if err != nil {
			return nil, err
		}

		if err := hcl.DecodeObject(&config, obj); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	if config.VaultAddr == "" {
		return nil, fmt.Errorf("config missing vault_addr")
	}
	if config.SSHMountPoint == "" {
		config.SSHMountPoint = SSHAgentDefaultMountPoint
	}

	return &config, nil
}

// Creates an SSHAgent object which can talk to Vault server with SSH backend
// mounted at default path ("ssh").
func (c *Client) SSHAgent() *SSHAgent {
	return c.SSHAgentWithMountPoint(SSHAgentDefaultMountPoint)
}

// Creates an SSHAgent object which can talk to Vault server with SSH backend
// mounted at a specific mount point.
func (c *Client) SSHAgentWithMountPoint(mountPoint string) *SSHAgent {
	return &SSHAgent{
		c:          c,
		MountPoint: mountPoint,
	}
}

// Verifies if the key provided by user is present in Vault server. The response
// will contain the IP address and username associated with the OTP. In case the
// OTP matches the echo request message, instead of searching an entry for the OTP,
// an echo response message is returned. This feature is used by agent to verify if
// its configured correctly.
func (c *SSHAgent) Verify(otp string) (*SSHVerifyResponse, error) {
	data := map[string]interface{}{
		"otp": otp,
	}
	verifyPath := fmt.Sprintf("/v1/%s/verify", c.MountPoint)
	r := c.c.NewRequest("PUT", verifyPath)
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}

	if secret.Data == nil {
		return nil, nil
	}

	var verifyResp SSHVerifyResponse
	err = mapstructure.Decode(secret.Data, &verifyResp)
	if err != nil {
		return nil, err
	}
	return &verifyResp, nil
}

// Loads the certificate from given path and creates a certificate pool from it.
func loadCACert(path string) (*x509.CertPool, error) {
	certs, err := loadCertFromPEM(path)
	if err != nil {
		return nil, err
	}

	result := x509.NewCertPool()
	for _, cert := range certs {
		result.AddCert(cert)
	}

	return result, nil
}

// Loads the certificates present in the given directory and creates a
// certificate pool from it.
func loadCAPath(path string) (*x509.CertPool, error) {
	result := x509.NewCertPool()
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		certs, err := loadCertFromPEM(path)
		if err != nil {
			return err
		}

		for _, cert := range certs {
			result.AddCert(cert)
		}
		return nil
	}

	return result, filepath.Walk(path, fn)
}

// Creates a certificate from the given path
func loadCertFromPEM(path string) ([]*x509.Certificate, error) {
	pemCerts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	certs := make([]*x509.Certificate, 0, 5)
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		certs = append(certs, cert)
	}

	return certs, nil
}
