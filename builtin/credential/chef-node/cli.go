package chefnode

import (
	"fmt"
	"net/url"

	"strings"

	"github.com/hashicorp/vault/api"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "chef-node"
	}

	clientName, ok := m["client_name"]
	if !ok {
		return "", fmt.Errorf("Chef client name must be as value for 'client_name' token")
	}

	clientKey, ok := m["client_key"]
	if !ok {
		return "", fmt.Errorf("Chef client key must be provided as value for 'client_key token")
	}
	path := fmt.Sprintf("auth/%s/login", mount)

	conf := &config{
		ClientKey:  clientKey,
		ClientName: clientName,
	}

	// Use the vault logical path
	vaultURL, err := url.Parse(c.Address())
	if err != nil {
		return "", err
	}
	vaultURL.Path = "/v1/" + path

	headers, err := authHeaders(conf, vaultURL, "POST", false)
	if err != nil {
		return "", err
	}

	sigVer := headers.Get("X-Ops-Sign")
	sig := headers.Get("X-Ops-Authorization")
	ts := headers.Get("X-Ops-Timestamp")

	secret, err := c.Logical().Write(path, map[string]interface{}{
		"signature_version": sigVer,
		"signature":         sig,
		"client_name":       clientName,
		"timestamp":         ts,
	})
	if err != nil {
		return "", err
	}
	if secret == nil {
		return "", fmt.Errorf("empty response from credential provider")
	}

	return secret.Auth.ClientToken, nil
}

func (h *CLIHandler) Help() string {
	help := `
The chef-node credential provider allows a chef node to authenticate against a Chef server.
Before using it you must configure a client for Vault to connect to the Chef server.  The
client information and the Chef server API endpoint must then be configured through the 'config'
endpoint.

The vault CLI can also be used to authenticate by giving a client name and client key in pem
format to use for authentication.

Example: vault auth -method=chef-node client_name=test_node client_key=@test_node.pem
`
	return strings.TrimSpace(help)
}
