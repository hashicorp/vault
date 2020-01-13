package transit

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/hashicorp/go-hclog"
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
	Encrypt(plaintext []byte) (ciphertext []byte, err error)
	Decrypt(ciphertext []byte) (plaintext []byte, err error)
}

type TransitClient struct {
	client          *api.Client
	lifetimeWatcher *api.Renewer

	mountPath string
	keyName   string
}

func newTransitClient(logger hclog.Logger, config map[string]string) (*TransitClient, map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	var mountPath, keyName string
	switch {
	case os.Getenv(EnvTransitWrapperMountPath) != "":
		mountPath = os.Getenv(EnvTransitWrapperMountPath)
	case os.Getenv(EnvVaultTransitSealMountPath) != "":
		mountPath = os.Getenv(EnvVaultTransitSealMountPath)
	case config["mount_path"] != "":
		mountPath = config["mount_path"]
	default:
		return nil, nil, fmt.Errorf("mount_path is required")
	}

	switch {
	case os.Getenv(EnvTransitWrapperKeyName) != "":
		keyName = os.Getenv(EnvTransitWrapperKeyName)
	case os.Getenv(EnvVaultTransitSealKeyName) != "":
		keyName = os.Getenv(EnvVaultTransitSealKeyName)
	case config["key_name"] != "":
		keyName = config["key_name"]
	default:
		return nil, nil, fmt.Errorf("key_name is required")
	}

	var disableRenewal bool
	var disableRenewalRaw string
	switch {
	case os.Getenv(EnvTransitWrapperDisableRenewal) != "":
		disableRenewalRaw = os.Getenv(EnvTransitWrapperDisableRenewal)
	case os.Getenv(EnvVaultTransitSealDisableRenewal) != "":
		disableRenewalRaw = os.Getenv(EnvVaultTransitSealDisableRenewal)
	case config["disable_renewal"] != "":
		disableRenewalRaw = config["disable_renewal"]
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
	case os.Getenv("VAULT_NAMESPACE") != "":
		namespace = os.Getenv("VAULT_NAMESPACE")
	case config["namespace"] != "":
		namespace = config["namespace"]
	}

	apiConfig := api.DefaultConfig()
	if config["address"] != "" {
		apiConfig.Address = config["address"]
	}
	if config["tls_ca_cert"] != "" || config["tls_ca_path"] != "" || config["tls_client_cert"] != "" || config["tls_client_key"] != "" ||
		config["tls_server_name"] != "" || config["tls_skip_verify"] != "" {
		var tlsSkipVerify bool
		if config["tls_skip_verify"] != "" {
			var err error
			tlsSkipVerify, err = strconv.ParseBool(config["tls_skip_verify"])
			if err != nil {
				return nil, nil, err
			}
		}

		tlsConfig := &api.TLSConfig{
			CACert:        config["tls_ca_cert"],
			CAPath:        config["tls_ca_path"],
			ClientCert:    config["tls_client_cert"],
			ClientKey:     config["tls_client_key"],
			TLSServerName: config["tls_server_name"],
			Insecure:      tlsSkipVerify,
		}
		if err := apiConfig.ConfigureTLS(tlsConfig); err != nil {
			return nil, nil, err
		}
	}

	apiClient, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, nil, err
	}
	if config["token"] != "" {
		apiClient.SetToken(config["token"])
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

	sealInfo := make(map[string]string)
	sealInfo["address"] = apiClient.Address()
	sealInfo["mount_path"] = mountPath
	sealInfo["key_name"] = keyName
	if namespace != "" {
		sealInfo["namespace"] = namespace
	}

	return client, sealInfo, nil
}

func (c *TransitClient) Close() {
	if c.lifetimeWatcher != nil {
		c.lifetimeWatcher.Stop()
	}
}

func (c *TransitClient) Encrypt(plaintext []byte) ([]byte, error) {
	encPlaintext := base64.StdEncoding.EncodeToString(plaintext)
	path := path.Join(c.mountPath, "encrypt", c.keyName)
	secret, err := c.client.Logical().Write(path, map[string]interface{}{
		"plaintext": encPlaintext,
	})
	if err != nil {
		return nil, err
	}

	return []byte(secret.Data["ciphertext"].(string)), nil
}

func (c *TransitClient) Decrypt(ciphertext []byte) ([]byte, error) {
	path := path.Join(c.mountPath, "decrypt", c.keyName)
	secret, err := c.client.Logical().Write(path, map[string]interface{}{
		"ciphertext": string(ciphertext),
	})
	if err != nil {
		return nil, err
	}

	plaintext, err := base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string))
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func (c *TransitClient) GetMountPath() string {
	return c.mountPath
}

func (c *TransitClient) GetApiClient() *api.Client {
	return c.client
}
