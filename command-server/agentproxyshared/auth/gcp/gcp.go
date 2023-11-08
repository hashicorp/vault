// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package gcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/vault/command-server/agentproxyshared/auth"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-gcp-common/gcputil"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/api"
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
		logger:         conf.Logger,
		mountPath:      conf.MountPath,
		serviceAccount: "default",
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
			jwtBytes, err := ioutil.ReadAll(resp.Body)
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
		if g.serviceAccount == "" && credentials != nil {
			serviceAccount = credentials.ClientEmail
		} else {
			serviceAccount = g.serviceAccount
		}
		if serviceAccount == "" {
			retErr = errors.New("could not obtain service account from credentials (possibly Application Default Credentials are being used); a service account to authenticate as must be provided")
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
