// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package token_file

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/vault/command-server/agentproxyshared/auth"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
)

type tokenFileMethod struct {
	logger    hclog.Logger
	mountPath string

	cachedToken   string
	tokenFilePath string
}

func NewTokenFileAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	a := &tokenFileMethod{
		logger:    conf.Logger,
		mountPath: "auth/token",
	}

	tokenFilePathRaw, ok := conf.Config["token_file_path"]
	if !ok {
		return nil, errors.New("missing 'token_file_path' value")
	}
	a.tokenFilePath, ok = tokenFilePathRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'token_file_path' config value to string")
	}
	if a.tokenFilePath == "" {
		return nil, errors.New("'token_file_path' value is empty")
	}

	return a, nil
}

func (a *tokenFileMethod) Authenticate(ctx context.Context, client *api.Client) (string, http.Header, map[string]interface{}, error) {
	token, err := os.ReadFile(a.tokenFilePath)
	if err != nil {
		if a.cachedToken == "" {
			return "", nil, nil, fmt.Errorf("error reading token file and no cached token known: %w", err)
		}
		a.logger.Warn("error reading token file", "error", err)
	}
	if len(token) == 0 {
		if a.cachedToken == "" {
			return "", nil, nil, errors.New("token file empty and no cached token known")
		}
		a.logger.Warn("token file exists but read empty value, re-using cached value")
	} else {
		a.cachedToken = strings.TrimSpace(string(token))
	}

	// i.e. auth/token/lookup-self
	return fmt.Sprintf("%s/lookup-self", a.mountPath), nil, map[string]interface{}{
		"token": a.cachedToken,
	}, nil
}

func (a *tokenFileMethod) NewCreds() chan struct{} {
	return nil
}

func (a *tokenFileMethod) CredSuccess() {
}

func (a *tokenFileMethod) Shutdown() {
}
