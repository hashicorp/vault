package failovercluster

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"

	"github.com/KnicKnic/go-windows/pkg/cluster"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
)

// FailoverClusterSeal is an auto-seal that uses ClusterEncrypt
// for crypto operations. See https://docs.microsoft.com/en-us/windows/win32/api/resapi/nf-resapi-clusterencrypt
// for more info
type FailoverClusterSeal struct {
	resourceName string
	currentKeyID *atomic.Value
	logger       log.Logger
}

// Ensure that we are implementing AutoSealAccess
var _ seal.Access = (*FailoverClusterSeal)(nil)

func NewSeal(logger log.Logger) *FailoverClusterSeal {
	v := &FailoverClusterSeal{
		logger:       logger,
		currentKeyID: new(atomic.Value),
	}
	v.currentKeyID.Store("")
	return v
}

// SetConfig sets the fields on the FailoverClusterSeal object based on
// values from the config parameter.
//
// Order of precedence:
// * Environment variable
// * Value from Vault configuration file
func (v *FailoverClusterSeal) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	switch {
	case os.Getenv("FAILOVERCLUSTER_RESOURCE_NAME") != "":
		v.resourceName = os.Getenv("FAILOVERCLUSTER_RESOURCE_NAME")
	case config["resource_name"] != "":
		v.resourceName = config["resource_name"]
	default:
		return nil, errors.New("resource name is required")
	}

	// Map that holds non-sensitive configuration info
	sealInfo := make(map[string]string)
	sealInfo["resource_name"] = v.resourceName

	return sealInfo, nil
}

// Init is called during core.Initialize.  This is a no-op.
func (v *FailoverClusterSeal) Init(context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op.
func (v *FailoverClusterSeal) Finalize(context.Context) error {
	return nil
}

// SealType returns the seal type for this particular seal implementation.
func (v *FailoverClusterSeal) SealType() string {
	return seal.FailoverCluster
}

// KeyID returns the last known key id.
func (v *FailoverClusterSeal) KeyID() string {
	return v.currentKeyID.Load().(string)
}

// Encrypt is used to encrypt using FailoverCluster Crypto Apis.
// This returns the ciphertext, and/or any errors from this
// call.
func (v *FailoverClusterSeal) Encrypt(ctx context.Context, plaintext []byte) (blob *physical.EncryptedBlobInfo, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", seal.FailoverCluster, "encrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", seal.FailoverCluster, "encrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
	metrics.IncrCounter([]string{"seal", seal.FailoverCluster, "encrypt"}, 1)

	if plaintext == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	env, err := seal.NewEnvelope().Encrypt(plaintext)
	if err != nil {
		return nil, errwrap.Wrapf("error wrapping dat: {{err}}", err)
	}

	encryptedData, err := v.clusterEncryptData(env.Key)
	if err != nil {
		return nil, err
	}

	keyVersion := "1"

	v.currentKeyID.Store(keyVersion)

	ret := &physical.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &physical.SealKeyInfo{
			KeyID:      keyVersion,
			WrappedKey: []byte(encryptedData),
		},
	}

	return ret, nil
}

func (v *FailoverClusterSeal) clusterEncryptData(plainData []byte) ([]byte, error) {

	handle, err := cluster.OpenClusterCryptProvider(v.resourceName,
		cluster.MS_ENH_RSA_AES_PROV,
		cluster.PROV_RSA_AES,
		cluster.CLUS_CREATE_CRYPT_CONTAINER_NOT_FOUND)
	defer handle.CloseClusterCryptProvider()

	if err != nil {
		return nil, err
	}

	encrypted, err := handle.ClusterEncrypt(plainData)
	return encrypted, err
}

func (v *FailoverClusterSeal) clusterDecryptData(encryptedData []byte) ([]byte, error) {

	handle, err := cluster.OpenClusterCryptProvider(v.resourceName,
		cluster.MS_ENH_RSA_AES_PROV,
		cluster.PROV_RSA_AES,
		cluster.CLUS_CREATE_CRYPT_CONTAINER_NOT_FOUND)
	defer handle.CloseClusterCryptProvider()

	if err != nil {
		return nil, err
	}
	decrypted, err := handle.ClusterDecrypt(encryptedData)
	return decrypted, err
}

// Decrypt is used to decrypt the ciphertext.
func (v *FailoverClusterSeal) Decrypt(ctx context.Context, in *physical.EncryptedBlobInfo) (pt []byte, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", seal.FailoverCluster, "decrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", seal.FailoverCluster, "decrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", seal.FailoverCluster, "decrypt"}, 1)

	if in == nil {
		return nil, errors.New("given input for decryption is nil")
	}

	if in.KeyInfo == nil {
		return nil, errors.New("key info is nil")
	}

	if in.KeyInfo.KeyID != "1" {
		return nil, errors.New("Invalid KeyID for FailoverCluster")
	}

	decryptedData, err := v.clusterDecryptData(in.KeyInfo.WrappedKey)
	if err != nil {
		return nil, err
	}

	envInfo := &seal.EnvelopeInfo{
		Key:        decryptedData,
		IV:         in.IV,
		Ciphertext: in.Ciphertext,
	}
	return seal.NewEnvelope().Decrypt(envInfo)
}
