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
		return "", fmt.Errorf("'code' var must be set, open a browser to: https://accounts.google.com/o/oauth2/v2/auth?response_type=code&client_id=158113233735-figmusvbkf0ui8g8u58am2tkumf9cnl8.apps.googleusercontent.com&redirect_uri=urn:ietf:wg:oauth:2.0:oob&scope=email")
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
//TODO nathang: client_id is something that is set on the config_path, not hardcode to my specific test app...
	help := `
The Google credential provider allows you to authenticate with Google.
To use it, specify the "code" parameter. The value should be a personal access
code for your Google account. You can generate a personal access token by browsing to google url
https://accounts.google.com/o/oauth2/v2/auth?response_type=code&client_id=158113233735-figmusvbkf0ui8g8u58am2tkumf9cnl8.apps.googleusercontent.com&redirect_uri=urn:ietf:wg:oauth:2.0:oob&scope=email

    Example: vault auth -method=google code=<code>

Key/Value Pairs:

    mount=google      The mountpoint for the Google credential provider.
                      Defaults to "google"

    code=<code>     The Google access code for authentication.
	`

	return strings.TrimSpace(help)
}
