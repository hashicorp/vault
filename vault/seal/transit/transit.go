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

	disableRenewalRaw := "false"
	switch {
	case os.Getenv("VAULT_TRANSIT_SEAL_DISABLE_RENEWAL") != "":
		disableRenewalRaw = os.Getenv("VAULT_TRANSIT_SEAL_DISABLE_RENEWAL")
	case config["key_name"] != "":
		disableRenewalRaw = config["disable_renewal"]
	}
	disableRenewal, err := strconv.ParseBool(disableRenewalRaw)
	if err != nil {
		return nil, err
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
		s.client = client

		if !disableRenewal {
			tokenInfo, err := client.Auth().Token().LookupSelf()
			if err != nil {
				return nil, err
			}

			// Only set up renewer if the token can be renewed
			if tokenInfo.Renewable {
				// Build a api.SecretAuth block for the renewer
				secretAuth := &api.SecretAuth{
					ClientToken:   s.client.Token(),
					LeaseDuration: tokenInfo.LeaseDuration,
					Renewable:     tokenInfo.Renewable,
				}
				renewer, err := s.client.NewRenewer(&api.RenewerInput{
					Secret: &api.Secret{
						Auth: secretAuth,
					},
				})
				if err != nil {
					return nil, err
				}
				s.renewer = renewer

				go func() {
					err := <-renewer.DoneCh()
					s.logger.Info("renewer done channel triggered")
					if err != nil {
						s.logger.Error("error renewing token", "error", err)
					}
				}()
				go s.renewer.Renew()
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
func (s *Seal) Encrypt(_ context.Context, plaintext []byte) (*physical.EncryptedBlobInfo, error) {
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
func (s *Seal) Decrypt(_ context.Context, in *physical.EncryptedBlobInfo) ([]byte, error) {
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
