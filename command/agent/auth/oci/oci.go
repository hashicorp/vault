// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oci

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	"github.com/oracle/oci-go-sdk/common"
	ociAuth "github.com/oracle/oci-go-sdk/common/auth"
)

const (
	typeAPIKey   = "apikey"
	typeInstance = "instance"

	/*

		IAM creds can be inferred from instance metadata or the container
		identity service, and those creds expire at varying intervals with
		new creds becoming available at likewise varying intervals. Let's
		default to polling once a minute so all changes can be picked up
		rather quickly. This is configurable, however.

	*/
	defaultCredCheckFreqSeconds = 60 * time.Second

	defaultConfigFileName    = "config"
	defaultConfigDirName     = ".oci"
	configFilePathEnvVarName = "OCI_CONFIG_FILE"
	secondaryConfigDirName   = ".oraclebmc"
)

func NewOCIAuthMethod(conf *auth.AuthConfig, vaultAddress string) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	a := &ociMethod{
		logger:       conf.Logger,
		vaultAddress: vaultAddress,
		mountPath:    conf.MountPath,
		credsFound:   make(chan struct{}),
		stopCh:       make(chan struct{}),
	}

	typeRaw, ok := conf.Config["type"]
	if !ok {
		return nil, errors.New("missing 'type' value")
	}
	authType, ok := typeRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'type' config value to string")
	}

	roleRaw, ok := conf.Config["role"]
	if !ok {
		return nil, errors.New("missing 'role' value")
	}
	a.role, ok = roleRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role' config value to string")
	}

	// Check for an optional custom frequency at which we should poll for creds.
	credCheckFreqSec := defaultCredCheckFreqSeconds
	if checkFreqRaw, ok := conf.Config["credential_poll_interval"]; ok {
		checkFreq, err := parseutil.ParseDurationSecond(checkFreqRaw)
		if err != nil {
			return nil, fmt.Errorf("could not parse credential_poll_interval: %v", err)
		}
		credCheckFreqSec = checkFreq
	}

	switch {
	case a.role == "":
		return nil, errors.New("'role' value is empty")
	case authType == "":
		return nil, errors.New("'type' value is empty")
	case authType != typeAPIKey && authType != typeInstance:
		return nil, errors.New("'type' value is invalid")
	case authType == typeAPIKey:
		defaultConfigFile := getDefaultConfigFilePath()
		homeFolder := getHomeFolder()
		secondaryConfigFile := path.Join(homeFolder, secondaryConfigDirName, defaultConfigFileName)

		environmentProvider := common.ConfigurationProviderEnvironmentVariables("OCI", "")
		defaultFileProvider, _ := common.ConfigurationProviderFromFile(defaultConfigFile, "")
		secondaryFileProvider, _ := common.ConfigurationProviderFromFile(secondaryConfigFile, "")

		provider, _ := common.ComposingConfigurationProvider([]common.ConfigurationProvider{environmentProvider, defaultFileProvider, secondaryFileProvider})
		a.configurationProvider = provider
	case authType == typeInstance:
		configurationProvider, err := ociAuth.InstancePrincipalConfigurationProvider()
		if err != nil {
			return nil, fmt.Errorf("failed to create instance principal configuration provider: %v", err)
		}
		a.configurationProvider = configurationProvider
	}

	// Do an initial population of the creds because we want to err right away if we can't
	// even get a first set.
	creds, err := a.configurationProvider.KeyID()
	if err != nil {
		return nil, err
	}
	a.lastCreds = creds

	go a.pollForCreds(credCheckFreqSec)

	return a, nil
}

type ociMethod struct {
	logger       hclog.Logger
	vaultAddress string
	mountPath    string

	configurationProvider common.ConfigurationProvider
	role                  string

	// These are used to share the latest creds safely across goroutines.
	credLock  sync.Mutex
	lastCreds string

	// Notifies the outer environment that it should call Authenticate again.
	credsFound chan struct{}

	// Detects that the outer environment is closing.
	stopCh chan struct{}
}

func (a *ociMethod) Authenticate(context.Context, *api.Client) (string, http.Header, map[string]interface{}, error) {
	a.credLock.Lock()
	defer a.credLock.Unlock()

	a.logger.Trace("beginning authentication")

	requestPath := fmt.Sprintf("/v1/%s/login/%s", a.mountPath, a.role)
	requestURL := fmt.Sprintf("%s%s", a.vaultAddress, requestPath)

	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", nil, nil, fmt.Errorf("error creating authentication request: %w", err)
	}

	request.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	signer := common.DefaultRequestSigner(a.configurationProvider)

	err = signer.Sign(request)

	if err != nil {
		return "", nil, nil, fmt.Errorf("error signing authentication request: %w", err)
	}

	parsedVaultAddress, err := url.Parse(a.vaultAddress)
	if err != nil {
		return "", nil, nil, fmt.Errorf("unable to parse vault address: %w", err)
	}

	request.Header.Set("Host", parsedVaultAddress.Host)
	request.Header.Set("(request-target)", fmt.Sprintf("%s %s", "get", requestPath))

	data := map[string]interface{}{
		"request_headers": request.Header,
	}

	return fmt.Sprintf("%s/login/%s", a.mountPath, a.role), nil, data, nil
}

func (a *ociMethod) NewCreds() chan struct{} {
	return a.credsFound
}

func (a *ociMethod) CredSuccess() {}

func (a *ociMethod) Shutdown() {
	close(a.credsFound)
	close(a.stopCh)
}

func (a *ociMethod) pollForCreds(frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()
	for {
		select {
		case <-a.stopCh:
			a.logger.Trace("shutdown triggered, stopping OCI auth handler")
			return
		case <-ticker.C:
			if err := a.checkCreds(); err != nil {
				a.logger.Warn("unable to retrieve current creds, retaining last creds", "error", err)
			}
		}
	}
}

func (a *ociMethod) checkCreds() error {
	a.credLock.Lock()
	defer a.credLock.Unlock()

	a.logger.Trace("checking for new credentials")
	currentCreds, err := a.configurationProvider.KeyID()
	if err != nil {
		return err
	}
	// These will always have different pointers regardless of whether their
	// values are identical, hence the use of DeepEqual.
	if currentCreds == a.lastCreds {
		a.logger.Trace("credentials are unchanged")
		return nil
	}
	a.lastCreds = currentCreds
	a.logger.Trace("new credentials detected, triggering Authenticate")
	a.credsFound <- struct{}{}
	return nil
}

func getHomeFolder() string {
	current, e := user.Current()
	if e != nil {
		// Give up and try to return something sensible
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return home
	}
	return current.HomeDir
}

func getDefaultConfigFilePath() string {
	homeFolder := getHomeFolder()
	defaultConfigFile := path.Join(homeFolder, defaultConfigDirName, defaultConfigFileName)
	if _, err := os.Stat(defaultConfigFile); err == nil {
		return defaultConfigFile
	}

	// Read configuration file path from OCI_CONFIG_FILE env var
	fallbackConfigFile, existed := os.LookupEnv(configFilePathEnvVarName)
	if !existed {
		return defaultConfigFile
	}
	if _, err := os.Stat(fallbackConfigFile); os.IsNotExist(err) {
		return defaultConfigFile
	}
	return fallbackConfigFile
}
