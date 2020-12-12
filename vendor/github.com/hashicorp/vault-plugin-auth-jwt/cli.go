package jwtauth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/base62"
)

const defaultMount = "oidc"
const defaultListenAddress = "localhost"
const defaultPort = "8250"
const defaultCallbackHost = "localhost"
const defaultCallbackMethod = "http"

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

	listenAddress, ok := m["listenaddress"]
	if !ok {
		listenAddress = defaultListenAddress
	}

	port, ok := m["port"]
	if !ok {
		port = defaultPort
	}

	callbackHost, ok := m["callbackhost"]
	if !ok {
		callbackHost = defaultCallbackHost
	}

	callbackMethod, ok := m["callbackmethod"]
	if !ok {
		callbackMethod = defaultCallbackMethod
	}

	callbackPort, ok := m["callbackport"]
	if !ok {
		callbackPort = port
	}

	role := m["role"]

	authURL, clientNonce, err := fetchAuthURL(c, role, mount, callbackPort, callbackMethod, callbackHost)
	if err != nil {
		return nil, err
	}

	// Set up callback handler
	http.HandleFunc("/oidc/callback", callbackHandler(c, mount, clientNonce, doneCh))

	listener, err := net.Listen("tcp", listenAddress+":"+port)
	if err != nil {
		return nil, err
	}
	defer listener.Close()

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

	// Wait for either the callback to finish, SIGINT to be received or up to 2 minutes
	select {
	case s := <-doneCh:
		return s.secret, s.err
	case <-sigintCh:
		return nil, errors.New("Interrupted")
	case <-time.After(2 * time.Minute):
		return nil, errors.New("Timed out waiting for response from provider")
	}
}

func callbackHandler(c *api.Client, mount string, clientNonce string, doneCh chan<- loginResp) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var response string
		var secret *api.Secret
		var err error

		defer func() {
			w.Write([]byte(response))
			doneCh <- loginResp{secret, err}
		}()

		// Pull any parameters from either the body or query parameters.
		// FormValue prioritizes body values, if found.
		data := map[string][]string{
			"state":        {req.FormValue("state")},
			"code":         {req.FormValue("code")},
			"id_token":     {req.FormValue("id_token")},
			"client_nonce": {clientNonce},
		}

		// If this is a POST, then the form_post response_mode is being used and the flow
		// involves an extra step. First POST the data to Vault, and then issue a GET with
		// the same state/code to complete the auth as normal.
		if req.Method == http.MethodPost {
			url := c.Address() + path.Join("/v1/auth", mount, "oidc/callback")
			resp, err := http.PostForm(url, data)
			if err != nil {
				summary, detail := parseError(err)
				response = errorHTML(summary, detail)
				return
			}
			defer resp.Body.Close()

			// An id_token will never be part of a redirect GET, so remove it here too.
			delete(data, "id_token")
		}

		secret, err = c.Logical().ReadWithData(fmt.Sprintf("auth/%s/oidc/callback", mount), data)
		if err != nil {
			summary, detail := parseError(err)
			response = errorHTML(summary, detail)
		} else {
			response = successHTML
		}
	}
}

func fetchAuthURL(c *api.Client, role, mount, callbackport string, callbackMethod string, callbackHost string) (string, string, error) {
	var authURL string

	clientNonce, err := base62.Random(20)
	if err != nil {
		return "", "", err
	}

	redirectURI := fmt.Sprintf("%s://%s:%s/oidc/callback", callbackMethod, callbackHost, callbackport)
	data := map[string]interface{}{
		"role":         role,
		"redirect_uri": redirectURI,
		"client_nonce": clientNonce,
	}

	secret, err := c.Logical().Write(fmt.Sprintf("auth/%s/oidc/auth_url", mount), data)
	if err != nil {
		return "", "", err
	}

	if secret != nil {
		authURL = secret.Data["auth_url"].(string)
	}

	if authURL == "" {
		return "", "", fmt.Errorf("Unable to authorize role %q with redirect_uri %q. Check Vault logs for more information.", role, redirectURI)
	}

	return authURL, clientNonce, nil
}

// isWSL tests if the binary is being run in Windows Subsystem for Linux
func isWSL() bool {
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		return false
	}
	data, err := ioutil.ReadFile("/proc/version")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read /proc/version.\n")
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

// openURL opens the specified URL in the default browser of the user.
// Source: https://stackoverflow.com/a/39324149/453290
func openURL(url string) error {
	var cmd string
	var args []string

	switch {
	case "windows" == runtime.GOOS || isWSL():
		cmd = "cmd.exe"
		args = []string{"/c", "start"}
		url = strings.Replace(url, "&", "^&", -1)
	case "darwin" == runtime.GOOS:
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

  listenaddress=<string>
    Optional address to bind the OIDC callback listener to (default: localhost).

  port=<string>
    Optional localhost port to use for OIDC callback (default: 8250).

  callbackmethod=<string>
    Optional method to to use in OIDC redirect_uri (default: http).

  callbackhost=<string>
    Optional callback host address to use in OIDC redirect_uri (default: localhost).

  callbackport=<string>
      Optional port to to use in OIDC redirect_uri (default: the value set for port).
`

	return strings.TrimSpace(help)
}
