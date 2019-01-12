package oidc

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/hashicorp/vault/api"
)

const defaultPath = "jwt"

// CLIHandler struct
type CLIHandler struct{}

type loginResp struct {
	secret *api.Secret
	err    error
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	defer signal.Stop(ch)

	doneCh := make(chan loginResp)

	role := m["role"]
	if role == "" {
		return nil, errors.New("a 'role' must be specified")
	}

	path := m["path"]
	if path == "" {
		path = defaultPath
	}

	server := &http.Server{}

	http.HandleFunc(fmt.Sprintf("/v1/auth/%s/oidc/callback", path), func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		code := query.Get("code")
		state := query.Get("state")
		data := map[string][]string{
			"code":  []string{code},
			"state": []string{state},
		}

		secret, err := c.Logical().ReadWithData(fmt.Sprintf("auth/%s/oidc/callback", path), data)

		w.Write([]byte("Yay!"))
		doneCh <- loginResp{secret, err}
	})

	go func() {
		err := server.Serve(startListening())
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	secret, err := fetchAuthURL(c, role, path)
	if err != nil {
		return nil, err
	}

	authURL := secret.Data["auth_url"].(string)
	if authURL == "" {
		return nil, errors.New(fmt.Sprintf("invalid role '%s'", role))
	}

	fmt.Fprintf(os.Stderr, "Complete the login via your OIDC provider. Launching browser to:\n\n    %s\n\n\n", secret.Data["auth_url"].(string))
	exec.Command("open", secret.Data["auth_url"].(string)).Start()

	// Wait on either the callback to finish or the signal to come through
	select {
	case s := <-doneCh:
		return s.secret, s.err
	case <-ch:
		return nil, errors.New("interrupted")
	}
}

func fetchAuthURL(c *api.Client, role, path string) (*api.Secret, error) {
	data := map[string]interface{}{
		"role":         role,
		"redirect_uri": fmt.Sprintf("http://127.0.0.1:8300/v1/auth/%s/oidc/callback", path),
	}

	return c.Logical().Write(fmt.Sprintf("auth/%s/oidc/auth_url", path), data)
}

func startListening() net.Listener {
	listener, err := net.Listen("tcp", "127.0.0.1:8300")
	if err != nil {
		panic(err)
	}

	return listener
}

// Help method for okta cli
func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=okta [CONFIG K=V...]

  The Okta auth method allows users to authenticate using Okta.

  Authenticate as "sally":

      $ vault login -method=okta username=sally
      Password (will be hidden):

  Authenticate as "bob":

      $ vault login -method=okta username=bob password=password

Configuration:

  password=<string>
      Okta password to use for authentication. If not provided, the CLI will
      prompt for this on stdin.

  username=<string>
      Okta username to use for authentication.
`

	return strings.TrimSpace(help)
}
