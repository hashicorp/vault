// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package okta

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/base62"
	pwd "github.com/hashicorp/go-secure-stdlib/password"
	"github.com/hashicorp/vault/api"
)

// CLIHandler struct
type CLIHandler struct{}

// Auth cli method
func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "okta"
	}

	username, ok := m["username"]
	if !ok {
		return nil, fmt.Errorf("'username' var must be set")
	}
	password, ok := m["password"]
	if !ok {
		fmt.Fprintf(os.Stderr, "Password (will be hidden): ")
		var err error
		password, err = pwd.Read(os.Stdin)
		fmt.Fprintf(os.Stderr, "\n")
		if err != nil {
			return nil, err
		}
	}

	data := map[string]interface{}{
		"password": password,
	}

	// Okta or Google totp code
	if totp, ok := m["totp"]; ok {
		data["totp"] = totp
	}

	// provider is an optional parameter
	if provider, ok := m["provider"]; ok {
		data["provider"] = provider
	}

	nonce := base62.MustRandom(20)
	data["nonce"] = nonce

	// Create a done channel to signal termination of the login so that we can
	// clean up the goroutine
	doneCh := make(chan struct{})
	defer close(doneCh)

	go func() {
		for {
			timer := time.NewTimer(time.Second)
			select {
			case <-doneCh:
				timer.Stop()
				return
			case <-timer.C:
			}

			resp, _ := c.Logical().Read(fmt.Sprintf("auth/%s/verify/%s", mount, nonce))
			if resp != nil {
				fmt.Fprintf(os.Stderr, "In Okta Verify, tap the number %q\n", resp.Data["correct_answer"].(json.Number))
				return
			}
		}
	}()

	path := fmt.Sprintf("auth/%s/login/%s", mount, username)
	secret, err := c.Logical().Write(path, data)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("empty response from credential provider")
	}

	return secret, nil
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
