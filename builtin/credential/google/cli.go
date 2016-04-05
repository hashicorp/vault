package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)


//CLIHandler for "vault auth -method=google ..."
type CLIHandler struct{}


//Auth logic for handling the google authentication code
func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = BackendName
	}

	code, ok := m[googleAuthCodeParameterName]
	if !ok {
		return "", fmt.Errorf("'%s' var must be set: %s", googleAuthCodeParameterName, readCodeUrlPathHelp)
	}

	path := fmt.Sprintf("auth/%s/%s", mount, loginPath)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		googleAuthCodeParameterName: code,
	})
	if err != nil {
		return "", err
	}
	if secret == nil {
		return "", fmt.Errorf("empty response from credential provider")
	}

	return secret.Auth.ClientToken, nil
}

//Help message on how to authenticate with google
func (h *CLIHandler) Help() string {

	return strings.TrimSpace(googleBackendHelp)
}
