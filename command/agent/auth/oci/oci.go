package oci

import (
	"context"
	"errors"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
)

const (
	typeAPIKey            = "apiKey"
	typeInstancePrincipal = "instancePrincipal"
)

type ociMethod struct {
	logger   hclog.Logger
	creds    chan struct{}
	authType string
	authRole string
}

func NewOCIAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	o := &ociMethod{
		logger: conf.Logger,
		creds:  make(chan struct{}),
	}

	typeRaw, ok := conf.Config["type"]
	if !ok {
		return nil, errors.New("missing 'type' value")
	}
	o.authType, ok = typeRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'type' config value to string")
	}

	roleRaw, ok := conf.Config["role"]
	if !ok {
		return nil, errors.New("missing 'role' value")
	}
	o.authRole, ok = roleRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role' config value to string")
	}

	switch {
	case o.authRole == "":
		return nil, errors.New("'role' value is empty")
	case o.authType == "":
		return nil, errors.New("'type' value is empty")
	case o.authType != typeAPIKey && o.authType != typeInstancePrincipal:
		return nil, errors.New("'type' value is invalid")
	}

	// next we have to find either the API key or the token from instance metadata. most of the code
	// in the OCI auth method leverages generated library methods for this. TBD

	return o, nil
}

func (o *ociMethod) Authenticate(ctx context.Context, client *api.Client) (retToken string, retData map[string]interface{}, retErr error) {
	return "", make(map[string]interface{}), nil
}

func (o *ociMethod) NewCreds() chan struct{} {
	return o.creds
}

func (o *ociMethod) CredSuccess() {}

func (o *ociMethod) Shutdown() {}
