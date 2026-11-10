// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package ibmcloud

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

// CLIHandler implements the Vault CLI auth handler interface
type CLIHandler struct{}

// Auth authenticates to Vault using IBM Cloud credentials. The credential is
// exchanged with IBM Cloud IAM locally; only the resulting IAM access token is
// sent to Vault.
func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "ibmcloud"
	}

	role, ok := m["role"]
	if !ok {
		return nil, fmt.Errorf("'role' must be specified")
	}

	token, err := resolveToken(m)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"role":  role,
		"token": token,
	}

	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, data)
	if err != nil {
		return nil, err
	}

	if secret == nil {
		return nil, fmt.Errorf("empty response from credential provider")
	}

	return secret, nil
}

// resolveToken returns the IAM access token to send to Vault, exchanging an
// API key or a compute resource token with IAM, or querying the VPC metadata
// service when no token is supplied.
func resolveToken(m map[string]string) (string, error) {
	if token := m["token"]; token != "" {
		return token, nil
	}

	iamEndpoint := m["iam_endpoint"]
	if iamEndpoint == "" {
		iamEndpoint = defaultIAMEndpoint
	}

	client, err := NewIBMCloudIAMClient(iamEndpoint, nil)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Explicit compute resource token provided
	if crToken := m["cr_token"]; crToken != "" {
		resp, err := client.AuthenticateWithCRToken(ctx, crToken, m["trusted_profile_id"])
		if err != nil {
			return "", fmt.Errorf("failed to exchange CR token: %w", err)
		}
		return resp.AccessToken, nil
	}

	// 2. Explicit API key provided in map or environment variable
	apiKey := m["apikey"]
	if apiKey == "" && m["type"] != "trusted_profile" && m["trusted_profile_id"] == "" {
		apiKey = os.Getenv("IBMCLOUD_API_KEY")
	}

	if apiKey != "" && m["type"] != "trusted_profile" && m["trusted_profile_id"] == "" {
		resp, err := client.AuthenticateWithAPIKey(ctx, apiKey)
		if err != nil {
			return "", fmt.Errorf("failed to exchange API key: %w", err)
		}
		return resp.AccessToken, nil
	}

	// 3. Trusted profile on VPC instance (explicit type=trusted_profile, or trusted_profile_id set, or fallback when no apiKey).
	// This branch is always reached when control gets here: if apiKey != "" the
	// step-2 block above returned already, so apiKey == "" is always true here,
	// making the condition always true. The outer if is kept for clarity.
	resp, err := client.ExchangeVPCMetadata(ctx, m["metadata_endpoint"], m["trusted_profile_id"])
	if err != nil {
		if m["type"] == "" && m["trusted_profile_id"] == "" {
			// No explicit trusted-profile intent: the caller provided nothing usable.
			return "", fmt.Errorf("one of 'token', 'apikey' (or IBMCLOUD_API_KEY), 'cr_token', or running on a VPC instance with a trusted profile is required (metadata exchange failed: %w)", err)
		}
		return "", fmt.Errorf("failed to exchange VPC metadata token: %w", err)
	}
	return resp.AccessToken, nil
}

// Help returns the help text for the IBM Cloud auth method
func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=ibmcloud [CONFIG K=V...]

  The IBM Cloud auth method authenticates with an IBM Cloud IAM access token.
  The CLI obtains the token locally by exchanging your API key, compute resource
  token, or by querying the VPC metadata service. Only the resulting IAM access
  token is sent to Vault.

  Authenticate using an existing IAM access token:

      $ vault login -method=ibmcloud \
          role=my-role \
          token=your-iam-access-token

  Authenticate using an API key (user or service ID):

      $ vault login -method=ibmcloud \
          role=my-role \
          apikey=your-api-key

      # or using the environment variable:
      $ export IBMCLOUD_API_KEY=your-api-key
      $ vault login -method=ibmcloud role=my-role

  Authenticate using a VPC virtual server instance (auto-discovered trusted profile):

      $ vault login -method=ibmcloud \
          role=my-role \
          type=trusted_profile

  Authenticate using a specific trusted profile on a VPC instance:

      $ vault login -method=ibmcloud \
          role=my-role \
          type=trusted_profile \
          trusted_profile_id=Profile-xxxx

  Authenticate using an explicit compute resource token:

      $ vault login -method=ibmcloud \
          role=my-role \
          cr_token=your-cr-token \
          trusted_profile_id=Profile-xxxx

Configuration:

  mount=<string>
      Path where the IBM Cloud auth method is mounted. This is usually
      provided via the -path flag in the "vault login" command, but it can
      be specified here as well. If specified here, it takes precedence
      over the value for -path. The default value is "ibmcloud".

  role=<string>
      Name of the role to authenticate against. This is required.

  token=<string>
      IBM Cloud IAM access token, used as is.

  apikey=<string>
      IBM Cloud API key to exchange for an IAM access token. Can also be
      specified via the IBMCLOUD_API_KEY environment variable.

  type=<string>
      Authentication type: "apikey" or "trusted_profile".

  cr_token=<string>
      Compute resource token to exchange for an IAM access token.

  trusted_profile_id=<string>
      Trusted profile ID. When authenticating from a VPC VSI without an explicit
      CR token, this is optional; omitting it uses the instance's default
      trusted profile.

  metadata_endpoint=<string>
      VPC metadata service endpoint used to fetch instance identity tokens.
      The default value is "http://169.254.169.254".

  iam_endpoint=<string>
      IBM Cloud IAM endpoint used for the exchange. The default value is
      "https://iam.cloud.ibm.com".
`

	return strings.TrimSpace(help)
}

// Made with Bob
