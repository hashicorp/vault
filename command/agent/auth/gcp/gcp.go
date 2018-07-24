package gcp

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
)

const (
	typeGCE          = "gce"
	typeIAM          = "iam"
	identityEndpoint = "http://metadata/computeMetadata/v1/instance/service-accounts/%s/identity"
)

type gcpMethod struct {
	logger         hclog.Logger
	authType       string
	mountPath      string
	role           string
	credentials    string
	serviceAccount string
}

func NewGCPAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

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

	return g, nil
}

func (g *gcpMethod) Authenticate(ctx context.Context, client *api.Client) (*api.Secret, error) {
	g.logger.Trace("beginning authentication")

	data := make(map[string]interface{})

	switch g.authType {
	case typeGCE:
		httpClient := cleanhttp.DefaultClient()

		// Fetch token
		{
			req, err := http.NewRequest("GET", fmt.Sprintf(identityEndpoint, g.serviceAccount), nil)
			if err != nil {
				return nil, errwrap.Wrapf("error creating request: {{err}}", err)
			}
			req = req.WithContext(ctx)
			req.Header.Add("Metadata-Flavor", "Google")
			q := req.URL.Query()
			q.Add("audience", fmt.Sprintf("%s/vault/%s", client.Address(), g.role))
			q.Add("format", "full")
			req.URL.RawQuery = q.Encode()
			resp, err := httpClient.Do(req)
			if err != nil {
				return nil, errwrap.Wrapf("error fetching instance token: {{err}}", err)
			}
			if resp == nil {
				return nil, errors.New("empty response fetching instance toke")
			}
			defer resp.Body.Close()
			token, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, errwrap.Wrapf("error reading instance token response body: {{err}}", err)
			}
			data["jwt"] = string(token)
		}

	default:
		/*
			var err error
			data, err = gcpauth.GenerateLoginData(g.accessKey, g.secretKey, g.sessionToken, g.headerValue)
			if err != nil {
				return nil, errwrap.Wrapf("error creating login value: {{err}}", err)
			}
		*/
	}

	data["role"] = g.role

	secret, err := client.Logical().Write(fmt.Sprintf("%s/login", g.mountPath), data)
	if err != nil {
		return nil, errwrap.Wrapf("error logging in: {{err}}", err)
	}

	return secret, nil
}

func (g *gcpMethod) NewCreds() chan struct{} {
	return nil
}

func (g *gcpMethod) Shutdown() {
}
