// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	"github.com/hashicorp/vault/sdk/helper/consts"
)

type certMethod struct {
	logger    hclog.Logger
	mountPath string
	name      string

	caCert     string
	clientCert string
	clientKey  string
	reload     bool

	// Client is the cached client to use if cert info was provided.
	client *api.Client
}

var _ auth.AuthMethodWithClient = &certMethod{}

func NewCertAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}

	// Not concerned if the conf.Config is empty as the 'name'
	// parameter is optional when using TLS Auth

	c := &certMethod{
		logger:    conf.Logger,
		mountPath: conf.MountPath,
	}

	if conf.Config != nil {
		nameRaw, ok := conf.Config["name"]
		if !ok {
			nameRaw = ""
		}
		c.name, ok = nameRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'name' config value to string")
		}

		caCertRaw, ok := conf.Config["ca_cert"]
		if ok {
			c.caCert, ok = caCertRaw.(string)
			if !ok {
				return nil, errors.New("could not convert 'ca_cert' config value to string")
			}
		}

		clientCertRaw, ok := conf.Config["client_cert"]
		if ok {
			c.clientCert, ok = clientCertRaw.(string)
			if !ok {
				return nil, errors.New("could not convert 'cert_file' config value to string")
			}
		}

		clientKeyRaw, ok := conf.Config["client_key"]
		if ok {
			c.clientKey, ok = clientKeyRaw.(string)
			if !ok {
				return nil, errors.New("could not convert 'cert_key' config value to string")
			}
		}

		reload, ok := conf.Config["reload"]
		if ok {
			c.reload, ok = reload.(bool)
			if !ok {
				return nil, errors.New("could not convert 'reload' config value to bool")
			}
		}
	}

	return c, nil
}

func (c *certMethod) Authenticate(_ context.Context, client *api.Client) (string, http.Header, map[string]interface{}, error) {
	c.logger.Trace("beginning authentication")

	authMap := map[string]interface{}{}

	if c.name != "" {
		authMap["name"] = c.name
	}

	return fmt.Sprintf("%s/login", c.mountPath), nil, authMap, nil
}

func (c *certMethod) NewCreds() chan struct{} {
	return nil
}

func (c *certMethod) CredSuccess() {}

func (c *certMethod) Shutdown() {}

// AuthClient uses the existing client's address and returns a new client with
// the auto-auth method's certificate information if that's provided in its
// config map.
func (c *certMethod) AuthClient(client *api.Client) (*api.Client, error) {
	c.logger.Trace("deriving auth client to use")

	clientToAuth := client

	if c.caCert != "" || (c.clientKey != "" && c.clientCert != "") {
		// Return cached client if present
		if c.client != nil && !c.reload {
			return c.client, nil
		}

		config := api.DefaultConfig()
		if config.Error != nil {
			return nil, config.Error
		}
		config.Address = client.Address()

		t := &api.TLSConfig{
			CACert:     c.caCert,
			ClientCert: c.clientCert,
			ClientKey:  c.clientKey,
		}

		// Setup TLS config
		if err := config.ConfigureTLS(t); err != nil {
			return nil, err
		}

		var err error
		clientToAuth, err = api.NewClient(config)
		if err != nil {
			return nil, err
		}
		if ns := client.Headers().Get(consts.NamespaceHeaderName); ns != "" {
			clientToAuth.SetNamespace(ns)
		}

		// Cache the client for future use
		c.client = clientToAuth
	}

	return clientToAuth, nil
}
