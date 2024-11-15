// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/vault/api"
)

const (
	EnvTransitWrapperMountPath   = "TRANSIT_WRAPPER_MOUNT_PATH"
	EnvVaultTransitSealMountPath = "VAULT_TRANSIT_SEAL_MOUNT_PATH"

	EnvTransitWrapperKeyName   = "TRANSIT_WRAPPER_KEY_NAME"
	EnvVaultTransitSealKeyName = "VAULT_TRANSIT_SEAL_KEY_NAME"

	EnvTransitWrapperDisableRenewal   = "TRANSIT_WRAPPER_DISABLE_RENEWAL"
	EnvVaultTransitSealDisableRenewal = "VAULT_TRANSIT_SEAL_DISABLE_RENEWAL"
)

type transitClientEncryptor interface {
	Close()
	Encrypt(ctx context.Context, plaintext []byte) (ciphertext []byte, err error)
	Decrypt(ctx context.Context, ciphertext []byte) (plaintext []byte, err error)
}

type TransitClient struct {
	client          *api.Client
	lifetimeWatcher *api.Renewer

	mountPath string
	keyName   string
}

func newTransitClient(logger hclog.Logger, opts *options) (*TransitClient, *wrapping.WrapperConfig, error) {
	var mountPath, keyName string
	switch {
	case os.Getenv(EnvTransitWrapperMountPath) != "" && !opts.Options.WithDisallowEnvVars:
		mountPath = os.Getenv(EnvTransitWrapperMountPath)
	case os.Getenv(EnvVaultTransitSealMountPath) != "" && !opts.Options.WithDisallowEnvVars:
		mountPath = os.Getenv(EnvVaultTransitSealMountPath)
	case opts.withMountPath != "":
		mountPath = opts.withMountPath
	default:
		return nil, nil, fmt.Errorf("mount_path is required")
	}

	switch {
	case os.Getenv(EnvTransitWrapperKeyName) != "" && !opts.Options.WithDisallowEnvVars:
		keyName = os.Getenv(EnvTransitWrapperKeyName)
	case os.Getenv(EnvVaultTransitSealKeyName) != "" && !opts.Options.WithDisallowEnvVars:
		keyName = os.Getenv(EnvVaultTransitSealKeyName)
	case opts.withKeyName != "":
		keyName = opts.withKeyName
	default:
		return nil, nil, fmt.Errorf("key_name is required")
	}

	var disableRenewal bool
	var disableRenewalRaw string
	switch {
	case os.Getenv(EnvTransitWrapperDisableRenewal) != "" && !opts.Options.WithDisallowEnvVars:
		disableRenewalRaw = os.Getenv(EnvTransitWrapperDisableRenewal)
	case os.Getenv(EnvVaultTransitSealDisableRenewal) != "" && !opts.Options.WithDisallowEnvVars:
		disableRenewalRaw = os.Getenv(EnvVaultTransitSealDisableRenewal)
	case opts.withDisableRenewal != "":
		disableRenewalRaw = opts.withDisableRenewal
	}
	if disableRenewalRaw != "" {
		var err error
		disableRenewal, err = strconv.ParseBool(disableRenewalRaw)
		if err != nil {
			return nil, nil, err
		}
	}

	var namespace string
	switch {
	case os.Getenv("VAULT_NAMESPACE") != "" && !opts.Options.WithDisallowEnvVars:
		namespace = os.Getenv("VAULT_NAMESPACE")
	case opts.withNamespace != "":
		namespace = opts.withNamespace
	}

	apiConfig := api.DefaultConfig()
	if opts.withAddress != "" {
		apiConfig.Address = opts.withAddress
	}
	if opts.withTlsCaCert != "" ||
		opts.withTlsCaPath != "" ||
		opts.withTlsClientCert != "" ||
		opts.withTlsClientKey != "" ||
		opts.withTlsServerName != "" ||
		opts.withTlsSkipVerify {

		tlsConfig := &api.TLSConfig{
			CACert:        opts.withTlsCaCert,
			CAPath:        opts.withTlsCaPath,
			ClientCert:    opts.withTlsClientCert,
			ClientKey:     opts.withTlsClientKey,
			TLSServerName: opts.withTlsServerName,
			Insecure:      opts.withTlsSkipVerify,
		}
		if err := apiConfig.ConfigureTLS(tlsConfig); err != nil {
			return nil, nil, err
		}
	}

	apiClient, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, nil, err
	}
	if opts.withToken != "" {
		apiClient.SetToken(opts.withToken)
	}
	if namespace != "" {
		apiClient.SetNamespace(namespace)
	}
	if apiClient.Token() == "" {
		if logger != nil {
			logger.Info("no token provided to transit auto-seal")
		}
	}

	client := &TransitClient{
		client:    apiClient,
		mountPath: mountPath,
		keyName:   keyName,
	}

	if !disableRenewal && apiClient.Token() != "" {
		// Renew the token immediately to get a secret to pass to lifetime watcher
		secret, err := apiClient.Auth().Token().RenewTokenAsSelf(apiClient.Token(), 0)
		// If we don't get an error renewing, set up a lifetime watcher.  The token may not be renewable or not have
		// permission to renew-self.
		if err == nil {
			lifetimeWatcher, err := apiClient.NewLifetimeWatcher(&api.LifetimeWatcherInput{
				Secret: secret,
			})
			if err != nil {
				return nil, nil, err
			}
			client.lifetimeWatcher = lifetimeWatcher

			go func() {
				for {
					select {
					case err := <-lifetimeWatcher.DoneCh():
						if logger != nil {
							logger.Info("shutting down token renewal")
						}
						if err != nil {
							if logger != nil {
								logger.Error("error renewing token", "error", err)
							}
						}
						return
					case <-lifetimeWatcher.RenewCh():
						if logger != nil {
							logger.Trace("successfully renewed token")
						}
					}
				}
			}()
			go lifetimeWatcher.Start()
		} else {
			if logger != nil {
				logger.Info("unable to renew token, disabling renewal", "err", err)
			}
		}
	}

	wrapConfig := new(wrapping.WrapperConfig)
	wrapConfig.Metadata = make(map[string]string)
	wrapConfig.Metadata["address"] = apiClient.Address()
	wrapConfig.Metadata["mount_path"] = mountPath
	wrapConfig.Metadata["key_name"] = keyName
	if namespace != "" {
		wrapConfig.Metadata["namespace"] = namespace
	}

	return client, wrapConfig, nil
}

func (c *TransitClient) Close() {
	if c.lifetimeWatcher != nil {
		c.lifetimeWatcher.Stop()
	}
}

func (c *TransitClient) Encrypt(ctx context.Context, plaintext []byte) ([]byte, error) {
	encPlaintext := base64.StdEncoding.EncodeToString(plaintext)
	path := path.Join(c.mountPath, "encrypt", c.keyName)
	secret, err := c.client.Logical().WriteWithContext(ctx, path, map[string]interface{}{
		"plaintext": encPlaintext,
	})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, errors.New("after encrypt operation the returned secret from vault is nil")
	}
	if secret.Data == nil {
		return nil, errors.New("after encrypt operation no data was found in returned secret from vault")
	}
	ct := secret.Data["ciphertext"]
	if ct == nil {
		return nil, errors.New("after encrypt operation ciphertext was not found in data returned from vault")
	}
	ctStr, ok := ct.(string)
	if !ok {
		return nil, errors.New("after encrypt operation ciphertext in data returned from vault is not a string")
	}

	return []byte(ctStr), nil
}

func (c *TransitClient) Decrypt(ctx context.Context, ciphertext []byte) ([]byte, error) {
	path := path.Join(c.mountPath, "decrypt", c.keyName)
	secret, err := c.client.Logical().WriteWithContext(ctx, path, map[string]interface{}{
		"ciphertext": string(ciphertext),
	})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, errors.New("after decrypt operation the returned secret from vault is nil")
	}
	if secret.Data == nil {
		return nil, errors.New("after decrypt operation no data was found in returned secret from vault")
	}
	pt := secret.Data["plaintext"]
	if pt == nil {
		return nil, errors.New("after decrypt operation plaintext was not found in data returned from vault")
	}
	ptStr, ok := pt.(string)
	if !ok {
		return nil, errors.New("after decrypt operation plaintext in data returned from vault is not a string")
	}

	plaintext, err := base64.StdEncoding.DecodeString(ptStr)
	if err != nil {
		return nil, fmt.Errorf("error base64-decoding plaintext: %w", err)
	}
	return plaintext, nil
}

func (c *TransitClient) GetMountPath() string {
	return c.mountPath
}

func (c *TransitClient) GetApiClient() *api.Client {
	return c.client
}
