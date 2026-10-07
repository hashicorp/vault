// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package spiffe

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
)

const operationPrefixSpiffe = "spiffe"

// envSVIDToken is the environment variable that may hold a JWT-SVID, as an
// alternative to passing the sensitive value.
const envSVIDToken = "VAULT_SPIFFE_SVID_TOKEN"

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount := operationPrefixSpiffe
	if v, ok := m["mount"]; ok && v != "" {
		mount = v
	}

	role := m["role"]
	typeVal := m["type"]

	switch typeVal {
	case "", "auto", "cert", "jwt":
		// valid case
	default:
		return nil, fmt.Errorf("invalid type %q (must be of: auto, cert, jwt)", typeVal)
	}

	// Resolve the JWT-SVID in order of precedence:
	//   1. svid_token_file=<path>
	//   2. VAULT_SPIFFE_SVID_TOKEN
	var token string
	if path, ok := m["svid_token_file"]; ok && path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading svid_token_file %q: %w", path, err)
		}
		token = strings.TrimSpace(string(raw))
	}
	if token == "" {
		token = os.Getenv(envSVIDToken)
	}

	switch typeVal {
	case "", "auto":
		if token != "" {
			c.AddHeader("Authorization", "Bearer "+token)
		}
	case "jwt":
		if token != "" {
			c.AddHeader("Authorization", "Bearer "+token)
		} else {
			return nil, fmt.Errorf("svid token is required")
		}
	}

	data := map[string]interface{}{}
	if role != "" {
		data["role"] = role
	}

	data["type"] = typeVal

	path := "auth/" + mount + "/login"

	return c.Logical().Write(path, data)
}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=spiffe [CONFIG K=V...]

  The SPIFFE auth method allows authentication using a SVID.

  Authenticate using a local x.509 SVID. The -client-cert and -client-key
  flags are included with the "vault login" command, NOT as configuration to the
  auth method.

      $ vault login \
          -client-cert=/path/to/svid.pem \
          -client-key=/path/to/svid.key \
          -method=spiffe type=cert

  Authenticate using a JWT-SVID stored in a file:

      $ vault login -method=spiffe type=jwt svid_token_file=/run/secrets/svid.jwt

  Authenticate using a JWT-SVID stored in an environment variable:

      $ VAULT_SPIFFE_SVID_TOKEN=<jwt-svid> vault login -method=spiffe type=jwt

Configuration:

  mount=<string>
    Path of the SPIFFE auth method. Defaults to "spiffe".

  role=<string>
    Name of the role to authenticate against. If unset, Vault iterates
    over all configured roles.

  type=<string>
    SVID source for authentication. Defaults to "auto". Must be one of:
      auto  - use JWT SVIDs if provided, otherwise fall back to the peer
              x.509 certificate.
      cert  - only use SVIDs from the peer x.509 certificate.
      jwt   - only use the JWT-SVID supplied via svid_token_file
              or VAULT_SPIFFE_SVID_TOKEN.

  svid_token_file=<path>
    Path to a file containing the JWT-SVID string. Used when type=jwt.
    Takes precedence over VAULT_SPIFFE_SVID_TOKEN.

  VAULT_SPIFFE_SVID_TOKEN
    Environment variable containing the JWT-SVID string. Used when type=jwt.
    Used only when svid_token_file is not set.
`
	return strings.TrimSpace(help)
}
