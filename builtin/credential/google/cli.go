package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "google"
	}

	code, ok := m["code"]
	if !ok {
		return "", fmt.Errorf("'code' var must be set, access read auth/%s/%s for a link to obtain the code from google", mount, PATH_CODE_URL)
	}

	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		"code": code,
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
The Google credential provider allows you to authenticate with Google.
To use it, specify the "code" parameter. The value should be a personal access
code for your Google account. You can generate a personal access token by browsing to a google url.
after configuring the backend with a google application secret and id to identify as, access auth/google/code_url to see the url.

    Example: vault auth -method=google code=<code>

Key/Value Pairs:

    mount=google      The mountpoint for the Google credential provider.
                      Defaults to "google"

    code=<code>     The Google access code for authentication.
	`

	return strings.TrimSpace(help)
}
