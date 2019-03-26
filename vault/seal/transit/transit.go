package transit

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault/seal"
)

// Seal is a seal that leverages Vault's Transit secret
// engine
type Seal struct {
	logger  log.Logger
	client  *api.Client
	renewer *api.Renewer

	mountPath string
	keyName   string

	currentKeyID *atomic.Value
}

var _ seal.Access = (*Seal)(nil)

// NewSeal creates a new transit seal
func NewSeal(logger log.Logger) *Seal {
	s := &Seal{
		logger:       logger.ResetNamed("seal-transit"),
		currentKeyID: new(atomic.Value),
	}
	s.currentKeyID.Store("")
	return s
}

// SetConfig processes the config info from the server config
func (s *Seal) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	switch {
	case os.Getenv("VAULT_TRANSIT_SEAL_MOUNT_PATH") != "":
		s.mountPath = os.Getenv("VAULT_TRANSIT_SEAL_MOUNT_PATH")
	case config["mount_path"] != "":
		s.mountPath = config["mount_path"]
	default:
		return nil, fmt.Errorf("mount_path is required")
	}

	switch {
	case os.Getenv("VAULT_TRANSIT_SEAL_KEY_NAME") != "":
		s.keyName = os.Getenv("VAULT_TRANSIT_SEAL_KEY_NAME")
	case config["key_name"] != "":
		s.keyName = config["key_name"]
	default:
		return nil, fmt.Errorf("key_name is required")
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
			return nil, err
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
				return nil, err
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
			return nil, err
		}
	}

	if s.client == nil {
		client, err := api.NewClient(apiConfig)
		if err != nil {
			return nil, err
		}
		if config["token"] != "" {
			client.SetToken(config["token"])
		}
		if namespace != "" {
			client.SetNamespace(namespace)
		}
		if client.Token() == "" {
			return nil, errors.New("missing token")
		}
		s.client = client

		// Send a value to test the seal and to set the current key id
		if _, err := s.Encrypt(context.Background(), []byte("a")); err != nil {
			return nil, err
		}

		if !disableRenewal {
			// Renew the token immediately to get a secret to pass to renewer
			secret, err := client.Auth().Token().RenewTokenAsSelf(s.client.Token(), 0)
			// If we don't get an error renewing, set up a renewer.  The token may not be renewable or not have
			// permission to renew-self.
			if err == nil {
				renewer, err := s.client.NewRenewer(&api.RenewerInput{
					Secret: secret,
				})
				if err != nil {
					return nil, err
				}
				s.renewer = renewer

				go func() {
					for {
						select {
						case err := <-renewer.DoneCh():
							s.logger.Info("shutting down token renewal")
							if err != nil {
								s.logger.Error("error renewing token", "error", err)
							}
							return
						case <-renewer.RenewCh():
							s.logger.Trace("successfully renewed token")
						}
					}
				}()
				go s.renewer.Renew()
			} else {
				s.logger.Info("unable to renew token, disabling renewal", "err", err)
			}
		}
	}

	sealInfo := make(map[string]string)
	sealInfo["address"] = s.client.Address()
	sealInfo["mount_path"] = s.mountPath
	sealInfo["key_name"] = s.keyName
	if namespace != "" {
		sealInfo["namespace"] = namespace
	}

	return sealInfo, nil
}

// Init is called during core.Initialize
func (s *Seal) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown
func (s *Seal) Finalize(_ context.Context) error {
	if s.renewer != nil {
		s.renewer.Stop()
	}

	return nil
}

// SealType returns the seal type for this particular seal implementation.
func (s *Seal) SealType() string {
	return seal.Transit
}

// KeyID returns the last known key id.
func (s *Seal) KeyID() string {
	return s.currentKeyID.Load().(string)
}

// Encrypt is used to encrypt using Vaults Transit engine
func (s *Seal) Encrypt(_ context.Context, plaintext []byte) (blob *physical.EncryptedBlobInfo, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "transit", "encrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "transit", "encrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "transit", "encrypt"}, 1)

	encPlaintext := base64.StdEncoding.EncodeToString(plaintext)
	path := path.Join(s.mountPath, "encrypt", s.keyName)
	secret, err := s.client.Logical().Write(path, map[string]interface{}{
		"plaintext": encPlaintext,
	})
	if err != nil {
		return nil, err
	}

	ciphertext := secret.Data["ciphertext"].(string)
	splitKey := strings.Split(ciphertext, ":")
	if len(splitKey) != 3 {
		return nil, errors.New("invalid ciphertext returned")
	}
	keyID := splitKey[1]
	s.currentKeyID.Store(keyID)

	ret := &physical.EncryptedBlobInfo{
		Ciphertext: []byte(ciphertext),
		KeyInfo: &physical.SealKeyInfo{
			KeyID: keyID,
		},
	}
	return ret, nil
}

// Decrypt is used to decrypt the ciphertext
func (s *Seal) Decrypt(_ context.Context, in *physical.EncryptedBlobInfo) (pt []byte, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "transit", "decrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "transit", "decrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "transit", "decrypt"}, 1)

	path := path.Join(s.mountPath, "decrypt", s.keyName)
	secret, err := s.client.Logical().Write(path, map[string]interface{}{
		"ciphertext": string(in.Ciphertext),
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
