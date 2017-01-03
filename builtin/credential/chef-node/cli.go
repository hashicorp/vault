package chefNode

import (
	"fmt"
	"net/http"
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
	r := c.NewRequest("POST", "/v1/"+path)
	hr, err := r.ToHTTP()
	if err != nil {
		return "", err
	}

	conf := &config{
		ClientKey:  clientKey,
		ClientName: clientName,
	}

	// Use the vault logical path
	mungedURL := &url.URL{
		Scheme: hr.URL.Scheme,
		Host:   hr.URL.Host,
		Path:   path,
	}

	headers, err := authHeaders(conf, mungedURL, "POST")
	if err != nil {
		return "", err
	}

	hr.Header = headers
	client := &http.Client{}
	resp, err := client.Do(hr)
	if err != nil {
		return "", err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", nil
	}

	return secret.Auth.ClientToken, nil
}

func (h *CLIHandler) Help() string {
	return `Chef Node credential provider`
}
