package transit

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strconv"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
)

type transitClientEncryptor interface {
	Close()
	Encrypt(plaintext []byte) (ciphertext []byte, err error)
	Decrypt(ciphertext []byte) (plaintext []byte, err error)
}

type transitClient struct {
	client  *api.Client
	renewer *api.Renewer

	mountPath string
	keyName   string
}

func newTransitClient(logger log.Logger, config map[string]string) (*transitClient, map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	var mountPath, keyName string
	switch {
	case os.Getenv("VAULT_TRANSIT_SEAL_MOUNT_PATH") != "":
		mountPath = os.Getenv("VAULT_TRANSIT_SEAL_MOUNT_PATH")
	case config["mount_path"] != "":
		mountPath = config["mount_path"]
	default:
		return nil, nil, fmt.Errorf("mount_path is required")
	}

	switch {
	case os.Getenv("VAULT_TRANSIT_SEAL_KEY_NAME") != "":
		keyName = os.Getenv("VAULT_TRANSIT_SEAL_KEY_NAME")
	case config["key_name"] != "":
		keyName = config["key_name"]
	default:
		return nil, nil, fmt.Errorf("key_name is required")
	}

	var disableRenewal bool
	var disableRenewalRaw string
	switch {
	case os.Getenv("VAULT_TRANSIT_SEAL_DISABLE_RENEWAL") != "":
		disableRenewalRaw = os.Getenv("VAULT_TRANSIT_SEAL_DISABLE_RENEWAL")
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
		logger.Info("no token provided to transit auto-seal")
	}

	client := &transitClient{
		client:    apiClient,
		mountPath: mountPath,
		keyName:   keyName,
	}

	if !disableRenewal && apiClient.Token() != "" {
		// Renew the token immediately to get a secret to pass to renewer
		secret, err := apiClient.Auth().Token().RenewTokenAsSelf(apiClient.Token(), 0)
		// If we don't get an error renewing, set up a renewer.  The token may not be renewable or not have
		// permission to renew-self.
		if err == nil {
			renewer, err := apiClient.NewRenewer(&api.RenewerInput{
				Secret: secret,
			})
			if err != nil {
				return nil, nil, err
			}
			client.renewer = renewer

			go func() {
				for {
					select {
					case err := <-renewer.DoneCh():
						logger.Info("shutting down token renewal")
						if err != nil {
							logger.Error("error renewing token", "error", err)
						}
						return
					case <-renewer.RenewCh():
						logger.Trace("successfully renewed token")
					}
				}
			}()
			go renewer.Renew()
		} else {
			logger.Info("unable to renew token, disabling renewal", "err", err)
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

func (c *transitClient) Close() {
	if c.renewer != nil {
		c.renewer.Stop()
	}
}

func (c *transitClient) Encrypt(plaintext []byte) ([]byte, error) {
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

func (c *transitClient) Decrypt(ciphertext []byte) ([]byte, error) {
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
