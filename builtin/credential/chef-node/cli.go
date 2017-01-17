package chefnode

import (
	"fmt"
	"net/url"

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
	mungedURL, err := url.Parse(c.Address())
	if err != nil {
		return "", err
	}
	mungedURL.Path = path

	headers, err := authHeaders(conf, mungedURL, "POST", false)
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
	return `Chef Node credential provider`
}
