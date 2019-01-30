package gcpauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/parseutil"
	"golang.org/x/oauth2"
	"google.golang.org/api/iam/v1"
	"strings"
	"time"
)

type CLIHandler struct{}

func getSignedJwt(role string, m map[string]string) (string, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cleanhttp.DefaultClient())

	credentials, tokenSource, err := gcputil.FindCredentials(m["credentials"], ctx, iam.CloudPlatformScope)
	if err != nil {
		return "", fmt.Errorf("could not obtain credentials: %v", err)
	}

	httpClient := oauth2.NewClient(ctx, tokenSource)

	serviceAccount, ok := m["service_account"]
	if !ok && credentials != nil {
		serviceAccount = credentials.ClientEmail
	}
	if serviceAccount == "" {
		return "", errors.New("could not obtain service account from credentials (are you using Application Default Credentials?). You must provide a service account to authenticate as")
	}

	project, ok := m["project"]
	if !ok {
		if credentials != nil {
			project = credentials.ProjectId
		} else {
			project = "-"
		}
	}

	var ttl = time.Duration(defaultIamMaxJwtExpMinutes) * time.Minute
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

	jwtReq := &iam.SignJwtRequest{
		Payload: string(payloadBytes),
	}

	iamClient, err := iam.New(httpClient)
	if err != nil {
		return "", fmt.Errorf("could not create IAM client: %v", err)
	}

	resourceName := fmt.Sprintf("projects/%s/serviceAccounts/%s", project, serviceAccount)
	resp, err := iamClient.Projects.ServiceAccounts.SignJwt(resourceName, jwtReq).Do()
	if err != nil {
		return "", fmt.Errorf("unable to sign JWT for %s using given Vault credentials: %v", resourceName, err)
	}

	return resp.SignedJwt, nil
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

	loginToken, err := getSignedJwt(role, m)
	if err != nil {
		return nil, err
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
  
  project=<string>                                
	Project for the service account who will be authenticating to Vault.
    Defaults to the credential's "project_id" (if credentials are specified)."
`

	return strings.TrimSpace(help)
}
