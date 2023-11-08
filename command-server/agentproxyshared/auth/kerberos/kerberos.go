// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package kerberos

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/command-server/agentproxyshared/auth"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	kerberos "github.com/hashicorp/vault-plugin-auth-kerberos"
	"github.com/hashicorp/vault/api"
	"github.com/jcmturner/gokrb5/v8/spnego"
)

type kerberosMethod struct {
	logger    hclog.Logger
	mountPath string
	loginCfg  *kerberos.LoginCfg
}

func NewKerberosAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}
	username, err := read("username", conf.Config)
	if err != nil {
		return nil, err
	}
	service, err := read("service", conf.Config)
	if err != nil {
		return nil, err
	}
	realm, err := read("realm", conf.Config)
	if err != nil {
		return nil, err
	}
	keytabPath, err := read("keytab_path", conf.Config)
	if err != nil {
		return nil, err
	}
	krb5ConfPath, err := read("krb5conf_path", conf.Config)
	if err != nil {
		return nil, err
	}

	disableFast := false
	disableFastRaw, ok := conf.Config["disable_fast_negotiation"]
	if ok {
		disableFast, err = parseutil.ParseBool(disableFastRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'disable_fast_negotiation': %s", err)
		}
	}

	return &kerberosMethod{
		logger:    conf.Logger,
		mountPath: conf.MountPath,
		loginCfg: &kerberos.LoginCfg{
			Username:               username,
			Service:                service,
			Realm:                  realm,
			KeytabPath:             keytabPath,
			Krb5ConfPath:           krb5ConfPath,
			DisableFASTNegotiation: disableFast,
		},
	}, nil
}

func (k *kerberosMethod) Authenticate(context.Context, *api.Client) (string, http.Header, map[string]interface{}, error) {
	k.logger.Trace("beginning authentication")
	authHeaderVal, err := kerberos.GetAuthHeaderVal(k.loginCfg)
	if err != nil {
		return "", nil, nil, err
	}
	var header http.Header
	header = make(map[string][]string)
	header.Set(spnego.HTTPHeaderAuthRequest, authHeaderVal)
	return k.mountPath + "/login", header, make(map[string]interface{}), nil
}

// These functions are implemented to meet the AuthHandler interface,
// but we don't need to take advantage of them.
func (k *kerberosMethod) NewCreds() chan struct{} { return nil }
func (k *kerberosMethod) CredSuccess()            {}
func (k *kerberosMethod) Shutdown()               {}

// read reads a key from a map and convert its value to a string.
func read(key string, m map[string]interface{}) (string, error) {
	raw, ok := m[key]
	if !ok {
		return "", fmt.Errorf("%q is required", key)
	}
	v, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%q must be a string", key)
	}
	return v, nil
}
