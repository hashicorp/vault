package local

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
)

const (
	// EnvLocalSealKeyGlob is an environment variable name whose value
	// is expected to contain a shell glob pattern.  The files that match
	// this pattern will be used as LocalSeal's encryption and decryption
	// keys.  These filenames must contain a numerical index as one of
	// their period delimited components which identifies the key version
	// (e.g. foo.5.key is considered to have index 5)
	EnvLocalSealKeyGlob = "VAULT_LOCAL_SEAL_KEY_GLOB"

	// KeyLen is the length of encryption keys used by this seal
	KeyLen = 32

	// NonceLen is the length of encryption nonces used by this seal
	NonceLen = 12
)

// LocalSeal provides a seal.Access implementation that uses on-disk secrets to
// support auto seal functionality.
type LocalSeal struct {
	keyGlob      string
	logger       log.Logger
	currentKeyID *atomic.Value
}

// ensure LocalSeal implements the seal.Access interface
var _ seal.Access = (*LocalSeal)(nil)

// NewSeal creates a new LocalSeal with the provided logger
func NewSeal(logger log.Logger) *LocalSeal {
	l := &LocalSeal{
		logger:       logger,
		currentKeyID: new(atomic.Value),
	}
	l.currentKeyID.Store("")
	return l
}

// Init is called during core.Initialize
func (l *LocalSeal) Init(context.Context) error {
	return nil
}

// Finalize is called during shutdown
func (l *LocalSeal) Finalize(context.Context) error {
	return nil
}

// SealType returns the seal type for this particular seal implementation.
func (l *LocalSeal) SealType() string {
	return seal.Local
}

// SetConfig sets the fields on the LocalSeal object based on
// values from the config parameter.
func (l *LocalSeal) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	switch {
	case os.Getenv(EnvLocalSealKeyGlob) != "":
		l.keyGlob = os.Getenv(EnvLocalSealKeyGlob)
	case config["key_glob"] != "":
		l.keyGlob = config["key_glob"]
	default:
		return nil, fmt.Errorf("'key_glob' not found for local seal configuration")
	}

	if _, err := filepath.Glob(l.keyGlob); err != nil {
		return nil, errwrap.Wrapf("'key_glob' is invalid for local seal configuration: {{err}}", err)
	}

	// Map that holds non-sensitive configuration info to return
	sealInfo := make(map[string]string)
	sealInfo["key_glob"] = l.keyGlob
	return sealInfo, nil
}

// KeyID returns the last known key id
func (l *LocalSeal) KeyID() string {
	return l.currentKeyID.Load().(string)
}

// Encrypt is used to encrypt the master key using a local on-disk secret.
// Returns the ciphertext, and/or any errors from this call.
func (l *LocalSeal) Encrypt(ctx context.Context, plaintext []byte) (*physical.EncryptedBlobInfo, error) {
	if plaintext == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	f, keyIdx, err := l.getCurrentFile()
	if err != nil {
		return nil, err
	}

	key, err := l.readKey(f)
	if err != nil {
		return nil, err
	}

	aead, err := l.aeadEncrypter(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, NonceLen)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %v", err)
	}

	enc := &physical.EncryptedBlobInfo{
		Ciphertext: aead.Seal(nil, nonce, plaintext, nil),
		IV:         nonce,
		KeyInfo: &physical.SealKeyInfo{
			KeyID: strconv.Itoa(keyIdx),
		},
	}
	l.currentKeyID.Store(enc.KeyInfo.KeyID)
	return enc, nil
}

// Decrypt is used to decrypt the ciphertext.
func (l *LocalSeal) Decrypt(ctx context.Context, in *physical.EncryptedBlobInfo) ([]byte, error) {
	if in == nil {
		return nil, errors.New("given ciphertext for decryption is nil")
	}

	if in.KeyInfo == nil {
		return nil, errors.New("key info is nil")
	}

	f, err := l.getKeyFile(in.KeyInfo.KeyID)
	if err != nil {
		return nil, err
	}

	key, err := l.readKey(f)
	if err != nil {
		return nil, err
	}

	aead, err := l.aeadEncrypter(key)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, in.IV, in.Ciphertext, nil)
}

// getCurrentFile finds the key with the highest index that matches l.keyGlob.
// On success returns the full path of the file, its index and a nil error.
func (l *LocalSeal) getCurrentFile() (string, int, error) {
	files, err := filepath.Glob(l.keyGlob)
	if err != nil {
		return "", 0, err
	}

	maxIdx := -1
	maxFile := ""
	for _, f := range files {
		fIdx, err := l.getIndex(f)
		if err == nil && fIdx > maxIdx {
			maxIdx = fIdx
			maxFile = f
		}
	}

	if maxFile == "" {
		return "", 0, fmt.Errorf("unable to find current seal secret")
	}

	return maxFile, maxIdx, nil
}

// getKeyFile finds the full path key that has the given index
func (l *LocalSeal) getKeyFile(keyIndex string) (string, error) {
	idx, err := strconv.Atoi(keyIndex)
	if err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("error parsing key index '%s': {{err}}", keyIndex), err)
	}

	files, err := filepath.Glob(l.keyGlob)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		fIdx, err := l.getIndex(f)
		if err == nil && fIdx == idx {
			return f, nil
		}
	}

	return "", fmt.Errorf("unable to find seal secret for index '%s'", keyIndex)
}

// getIndex returns the key index of the given filename
func (l *LocalSeal) getIndex(filename string) (int, error) {
	parts := strings.Split(path.Base(filename), ".")
	for i := len(parts) - 1; i >= 0; i-- {
		if v, err := strconv.Atoi(parts[i]); err == nil {
			return v, nil
		}
	}

	return 0, fmt.Errorf("could not determine index of '%s'", filename)
}

func (l *LocalSeal) readKey(filename string) ([]byte, error) {
	key, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("error reading key from '%s': {{err}}", filename), err)
	}

	if len(key) < KeyLen {
		return nil, fmt.Errorf("key '%s' is too small", filename)
	}

	return key[:KeyLen], nil
}

func (l *LocalSeal) aeadEncrypter(key []byte) (cipher.AEAD, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create cipher: {{err}}", err)
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, errors.New("failed to initialize GCM mode")
	}

	return gcm, nil
}
