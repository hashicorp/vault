// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package gcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-gcp-common/gcputil"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	"golang.org/x/oauth2"
	"google.golang.org/api/iamcredentials/v1"
)

const (
	typeGCE                    = "gce"
	typeIAM                    = "iam"
	identityEndpoint           = "http://metadata/computeMetadata/v1/instance/service-accounts/%s/identity"
	defaultIamMaxJwtExpMinutes = 15
)

type gcpMethod struct {
	logger         hclog.Logger
	authType       string
	mountPath      string
	role           string
	credentials    string
	serviceAccount string
	project        string
	jwtExp         int64
}

func NewGCPAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	var err error

	g := &gcpMethod{
		logger:    conf.Logger,
		mountPath: conf.MountPath,
	}

	typeRaw, ok := conf.Config["type"]
	if !ok {
		return nil, errors.New("missing 'type' value")
	}
	g.authType, ok = typeRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'type' config value to string")
	}

	roleRaw, ok := conf.Config["role"]
	if !ok {
		return nil, errors.New("missing 'role' value")
	}
	g.role, ok = roleRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role' config value to string")
	}

	switch {
	case g.role == "":
		return nil, errors.New("'role' value is empty")
	case g.authType == "":
		return nil, errors.New("'type' value is empty")
	case g.authType != typeGCE && g.authType != typeIAM:
		return nil, errors.New("'type' value is invalid")
	}

	credentialsRaw, ok := conf.Config["credentials"]
	if ok {
		g.credentials, ok = credentialsRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'credentials' value into string")
		}
	}

	serviceAccountRaw, ok := conf.Config["service_account"]
	if ok {
		g.serviceAccount, ok = serviceAccountRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'service_account' value into string")
		}
	}

	projectRaw, ok := conf.Config["project"]
	if ok {
		g.project, ok = projectRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'project' value into string")
		}
	}

	jwtExpRaw, ok := conf.Config["jwt_exp"]
	if ok {
		g.jwtExp, err = parseutil.ParseInt(jwtExpRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'jwt_raw' into integer: %w", err)
		}
	}

	return g, nil
}

func (g *gcpMethod) Authenticate(ctx context.Context, client *api.Client) (retPath string, header http.Header, retData map[string]interface{}, retErr error) {
	g.logger.Trace("beginning authentication")

	data := make(map[string]interface{})
	var jwt string

	switch g.authType {
	case typeGCE:
		httpClient := cleanhttp.DefaultClient()

		// For GCE, "default" is a valid metadata-server alias when no
		// service_account is explicitly configured.
		serviceAccountGCE := g.serviceAccount
		if serviceAccountGCE == "" {
			serviceAccountGCE = "default"
		}

		// Fetch token
		{
			req, err := http.NewRequest("GET", fmt.Sprintf(identityEndpoint, g.serviceAccount), nil)
			if err != nil {
				retErr = fmt.Errorf("error creating request: %w", err)
				return
			}
			req = req.WithContext(ctx)
			req.Header.Add("Metadata-Flavor", "Google")
			q := req.URL.Query()
			q.Add("audience", fmt.Sprintf("%s/vault/%s", client.Address(), g.role))
			q.Add("format", "full")
			req.URL.RawQuery = q.Encode()
			resp, err := httpClient.Do(req)
			if err != nil {
				retErr = fmt.Errorf("error fetching instance token: %w", err)
				return
			}
			if resp == nil {
				retErr = errors.New("empty response fetching instance toke")
				return
			}
			defer resp.Body.Close()
			jwtBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				retErr = fmt.Errorf("error reading instance token response body: %w", err)
				return
			}

			jwt = string(jwtBytes)
		}

	default:
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cleanhttp.DefaultClient())

		credentials, tokenSource, err := gcputil.FindCredentials(g.credentials, ctx, iamcredentials.CloudPlatformScope)
		if err != nil {
			retErr = fmt.Errorf("could not obtain credentials: %w", err)
			return
		}

		httpClient := oauth2.NewClient(ctx, tokenSource)

		var serviceAccount string
		switch {
		case g.serviceAccount != "":
			// Explicitly configured — use it as-is.
			serviceAccount = g.serviceAccount
		case credentials != nil && credentials.ClientEmail != "":
			// Derive from GOOGLE_APPLICATION_CREDENTIALS / service account key.
			serviceAccount = credentials.ClientEmail
		default:
			retErr = errors.New("could not obtain service account: no 'service_account' configured and credentials do not contain a client email (Application Default Credentials may be in use); set 'service_account' explicitly in the auto_auth config")
			return
		}

		ttlMin := int64(defaultIamMaxJwtExpMinutes)
		if g.jwtExp != 0 {
			ttlMin = g.jwtExp
		}
		ttl := time.Minute * time.Duration(ttlMin)

		jwtPayload := map[string]interface{}{
			"aud": fmt.Sprintf("http://vault/%s", g.role),
			"sub": serviceAccount,
			"exp": time.Now().Add(ttl).Unix(),
		}
		payloadBytes, err := json.Marshal(jwtPayload)
		if err != nil {
			retErr = fmt.Errorf("could not convert JWT payload to JSON string: %w", err)
			return
		}

		jwtReq := &iamcredentials.SignJwtRequest{
			Payload: string(payloadBytes),
		}

		iamClient, err := iamcredentials.New(httpClient)
		if err != nil {
			retErr = fmt.Errorf("could not create IAM client: %w", err)
			return
		}

		resourceName := fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccount)
		resp, err := iamClient.Projects.ServiceAccounts.SignJwt(resourceName, jwtReq).Do()
		if err != nil {
			retErr = fmt.Errorf("unable to sign JWT for %s using given Vault credentials: %w", resourceName, err)
			return
		}

		jwt = resp.SignedJwt
	}

	data["role"] = g.role
	data["jwt"] = jwt

	return fmt.Sprintf("%s/login", g.mountPath), nil, data, nil
}

func (g *gcpMethod) NewCreds() chan struct{} {
	return nil
}

func (g *gcpMethod) CredSuccess() {
}

func (g *gcpMethod) Shutdown() {
}
