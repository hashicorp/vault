package jwtauth

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strings"

	"github.com/hashicorp/vault/api"
)

const defaultMount = "oidc"
const defaultPort = "8250"

var errorRegex = regexp.MustCompile(`(?s)Errors:.*\* *(.*)`)

type CLIHandler struct{}

type loginResp struct {
	secret *api.Secret
	err    error
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	// handle ctrl-c while waiting for the callback
	sigintCh := make(chan os.Signal, 1)
	signal.Notify(sigintCh, os.Interrupt)
	defer signal.Stop(sigintCh)

	doneCh := make(chan loginResp)

	mount, ok := m["mount"]
	if !ok {
		mount = defaultMount
	}

	port, ok := m["port"]
	if !ok {
		port = defaultPort
	}

	role := m["role"]

	authURL, err := fetchAuthURL(c, role, mount, port)
	if err != nil {
		return nil, err
	}

	// Set up callback handler
	http.HandleFunc("/oidc/callback", func(w http.ResponseWriter, req *http.Request) {
		var response string

		query := req.URL.Query()
		code := query.Get("code")
		state := query.Get("state")
		data := map[string][]string{
			"code":  {code},
			"state": {state},
		}

		secret, err := c.Logical().ReadWithData(fmt.Sprintf("auth/%s/oidc/callback", mount), data)
		if err != nil {
			summary, detail := parseError(err)
			response = errorHTML(summary, detail)
		} else {
			response = successHTML
		}

		w.Write([]byte(response))
		doneCh <- loginResp{secret, err}
	})

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	// Open the default browser to the callback URL.
	fmt.Fprintf(os.Stderr, "Complete the login via your OIDC provider. Launching browser to:\n\n    %s\n\n\n", authURL)
	if err := openURL(authURL); err != nil {
		fmt.Fprintf(os.Stderr, "Error attempting to automatically open browser: '%s'.\nPlease visit the authorization URL manually.", err)
	}

	// Start local server
	go func() {
		err := http.Serve(listener, nil)
		if err != nil && err != http.ErrServerClosed {
			doneCh <- loginResp{nil, err}
		}
	}()

	// Wait for either the callback to finish or SIGINT to be received
	select {
	case s := <-doneCh:
		return s.secret, s.err
	case <-sigintCh:
		return nil, errors.New("Interrupted")
	}
}

func fetchAuthURL(c *api.Client, role, mount, port string) (string, error) {
	data := map[string]interface{}{
		"role":         role,
		"redirect_uri": fmt.Sprintf("http://localhost:%s/oidc/callback", port),
	}

	secret, err := c.Logical().Write(fmt.Sprintf("auth/%s/oidc/auth_url", mount), data)
	if err != nil {
		return "", err
	}

	authURL := secret.Data["auth_url"].(string)
	if authURL == "" {
		return "", errors.New(fmt.Sprintf("Unable to authorize role %q. Check Vault logs for more information.", role))
	}

	return authURL, nil
}

// openURL opens the specified URL in the default browser of the user.
// Source: https://stackoverflow.com/a/39324149/453290
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// parseError converts error from the API into summary and detailed portions.
// This is used to present a nicer UI by splitting up *known* prefix sentences
// from the rest of the text. e.g.
//
//    "No response from provider. Gateway timeout from upstream proxy."
//
// becomes:
//
//    "No response from provider.", "Gateway timeout from upstream proxy."
func parseError(err error) (string, string) {
	headers := []string{errNoResponse, errLoginFailed, errTokenVerification}
	summary := "Login error"
	detail := ""

	errorParts := errorRegex.FindStringSubmatch(err.Error())
	switch len(errorParts) {
	case 0:
		summary = ""
	case 1:
		detail = errorParts[0]
	case 2:
		for _, h := range headers {
			if strings.HasPrefix(errorParts[1], h) {
				summary = h
				detail = strings.TrimSpace(errorParts[1][len(h):])
				break
			}
		}
		if detail == "" {
			detail = errorParts[1]
		}
	}

	return summary, detail
}

// Help method for OIDC cli
func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=oidc [CONFIG K=V...]

  The OIDC auth method allows users to authenticate using an OIDC provider.
  The provider must be configured as part of a role by the operator.

  Authenticate using role "engineering":

      $ vault login -method=oidc role=engineering
      Complete the login via your OIDC provider. Launching browser to:

          https://accounts.google.com/o/oauth2/v2/...

  The default browser will be opened for the user to complete the login. Alternatively,
  the user may visit the provided URL directly.

Configuration:

  role=<string>
      Vault role of type "OIDC" to use for authentication.

  port=<string>
      Optional localhost port to use for OIDC callback (default: 8250).
`

	return strings.TrimSpace(help)
}
