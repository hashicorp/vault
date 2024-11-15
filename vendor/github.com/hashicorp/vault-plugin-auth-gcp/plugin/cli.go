// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/api"
	"golang.org/x/oauth2"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
)

type CLIHandler struct{}

func getSignedJwt(role string, m map[string]string) (string, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cleanhttp.DefaultClient())

	credentials, tokenSource, err := gcputil.FindCredentials(m["credentials"], ctx, iamcredentials.CloudPlatformScope)
	if err != nil {
		return "", fmt.Errorf("could not obtain credentials: %v", err)
	}

	httpClient := oauth2.NewClient(ctx, tokenSource)

	serviceAccount, ok := m["service_account"]
	if !ok && credentials != nil {
		serviceAccount = credentials.ClientEmail
	}
	if serviceAccount == "" {
		// Check if the metadata server is available.
		if !metadata.OnGCE() {
			return "", errors.New("could not obtain service account from credentials (are you using Application Default Credentials?). You must provide a service account to authenticate as")
		}
		metadataClient := metadata.NewClient(cleanhttp.DefaultClient())
		v := url.Values{}
		v.Set("audience", fmt.Sprintf("http://vault/%s", role))
		v.Set("format", "full")
		path := "instance/service-accounts/default/identity?" + v.Encode()
		instanceJwt, err := metadataClient.Get(path)
		if err != nil {
			return "", fmt.Errorf("unable to read the identity token: %w", err)
		}
		return instanceJwt, nil

	} else {
		ttl := time.Duration(defaultIamMaxJwtExpMinutes) * time.Minute
		jwtExpStr, ok := m["jwt_exp"]
		if ok {
			ttl, err = parseutil.ParseDurationSecond(jwtExpStr)
			if err != nil {
				return "", fmt.Errorf("could not parse jwt_exp '%s' into integer value", jwtExpStr)
			}
		}

		jwtPayload := map[string]interface{}{
			"aud": fmt.Sprintf("http://vault/%s", role),
			"sub": serviceAccount,
			"exp": time.Now().Add(ttl).Unix(),
		}
		payloadBytes, err := json.Marshal(jwtPayload)
		if err != nil {
			return "", fmt.Errorf("could not convert JWT payload to JSON string: %v", err)
		}

		jwtReq := &iamcredentials.SignJwtRequest{
			Payload: string(payloadBytes),
		}

		iamClient, err := iamcredentials.NewService(ctx, option.WithHTTPClient(httpClient))
		if err != nil {
			return "", fmt.Errorf("could not create IAM client: %v", err)
		}

		resourceName := fmt.Sprintf(gcputil.ServiceAccountCredentialsTemplate, serviceAccount)
		resp, err := iamClient.Projects.ServiceAccounts.SignJwt(resourceName, jwtReq).Do()
		if err != nil {
			return "", fmt.Errorf("unable to sign JWT for %s using given Vault credentials: %v", resourceName, err)
		}

		return resp.SignedJwt, nil
	}
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	role, ok := m["role"]
	if !ok {
		return nil, errors.New("role is required")
	}

	mount, ok := m["mount"]
	if !ok {
		mount = "gcp"
	}

	var loginToken string
	var err error
	if v, ok := m["jwt"]; ok {
		loginToken = v
	} else {
		loginToken, err = getSignedJwt(role, m)
		if err != nil {
			return nil, err
		}
	}

	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(
		path,
		map[string]interface{}{
			"role": role,
			"jwt":  loginToken,
		})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("empty response from credential provider")
	}

	return secret, nil
}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=gcp [CONFIG K=V...]

The GCP credential provider allows you to authenticate with
a GCP IAM service account key. To use it, you provide a valid GCP IAM
credential JSON either explicitly on the command line (not recommended),
or through any GCP Application Default Credentials.

Authenticate using Application Default Credentials:

  Example: vault login -method=gcp role=my-iam-role

Authenticate using explicitly passed-in credentials:
  
	Example:
	vault login -method=gcp role=my-iam-role -credentials=@path/to/creds role=my-iam-role

This tool generates a signed JWT signed using the given credentials.

Configuration:

  role                                				
	Required. The name of the role you're requesting a token for.
  
  mount=gcp                           				
	This is usually provided via the -path flag in the "vault login" command,
	but it can be specified here as well. If specified here, it takes 
	precedence over the value for -path.
	
  credentials=<string>			
	Explicitly specified GCP credentials in JSON string format (not recommended)

  jwt_exp=<minutes>
	Time until the generated JWT expires in minutes.
	The given IAM role will have a max_jwt_exp field, the
	time in minutes that all valid authentication JWTs
	must expire within (from time of authentication).
	Defaults to 15 minutes, the default max_jwt_exp for a role.
	Must be less than an hour. 

  service_account=<string>	
	Service account to generate a JWT for. Defaults to credentials 
	"client_email" if "credentials" specified and this value is not. 
	The actual credential must have the "iam.serviceAccounts.signJWT" 
	permissions on this service account. 
`

	return strings.TrimSpace(help)
}
