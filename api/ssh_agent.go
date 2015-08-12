package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/hcl"
	"github.com/mitchellh/mapstructure"
)

// Default path at which SSH backend will be mounted
const SSHAgentDefaultMountPoint = "ssh"

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
	Message  string `mapstructure:"message"`
	Username string `mapstructure:"username"`
	IP       string `mapstructure:"ip"`
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

// Returns a HTTP client that uses TLS verification (TLS 1.2) with the given
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
	return &client
}

// Loads agent's configuration from the file and populates the corresponding
// in memory structure.
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

// Verifies if the key provided by user is present in Vault server. If yes,
// the response will contain the IP address and username associated with the
// key.
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
