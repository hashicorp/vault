package crypto

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	mathRand "math/rand"
	"time"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-kms-wrapping/wrappers/aead"
	"github.com/hashicorp/vault/api"
)

var _ KeyManager = (*ResponseEncrypter)(nil)

const ResponseWrappedTokenTTL = "60"

// ResponseEncrypter ...
type ResponseEncrypter struct {
	renewable bool
	wrapper   *aead.Wrapper
	token     []byte
	ttl       string
	Client    *api.Client
	Notify    chan struct{}
	Stop      chan struct{}
	Done      chan error
}

// NewResponseEncrypter ..
func NewResponseEncrypter(existingToken []byte, client *api.Client, ttl string) (*ResponseEncrypter, error) {
	if ttl == "" {
		ttl = ResponseWrappedTokenTTL
	}

	r := &ResponseEncrypter{
		renewable: true,
		wrapper:   aead.NewWrapper(nil),
		Client:    client,
		Notify:    make(chan struct{}, 1),
		Stop:      make(chan struct{}, 1),
		Done:      make(chan error, 1),
		ttl:       ttl,
	}
	r.wrapper.SetConfig(map[string]string{"key_id": KeyID})

	var rootKey []byte
	switch tokenLength := len(existingToken); {
	case tokenLength == 0:
		newKey := make([]byte, 32)
		_, err := rand.Read(newKey)
		if err != nil {
			return r, err
		}
		rootKey = newKey
	case tokenLength > 0:
		r.token = existingToken
		key, err := r.unwrap()
		if err != nil {
			return r, err
		}
		rootKey = key
	default:
		return r, fmt.Errorf("unknown error")
	}

	if err := r.wrapper.SetAESGCMKeyBytes(rootKey); err != nil {
		return r, err
	}

	return r, nil
}

// GetKey ...
func (r *ResponseEncrypter) GetKey() []byte {
	return r.wrapper.GetKeyBytes()
}

// GetPersistentKey ...
func (r *ResponseEncrypter) GetPersistentKey() ([]byte, error) {
	if r.token == nil {
		if r.Client.Token() == "" {
			return nil, fmt.Errorf("response wrapping requires a token set on client")
		}

		if err := r.wrapForStorage(); err != nil {
			return nil, err
		}
	}
	return r.token, nil
}

// Renewable ...
func (r *ResponseEncrypter) Renewable() bool {
	return r.renewable
}

// Renewer ...
func (r *ResponseEncrypter) Renewer(ctx context.Context) error {
	for {
		token, ttl, err := r.rewrap()
		if err != nil {
			r.Done <- err
			return err
		}
		r.token = []byte(token)
		r.Notify <- struct{}{}

		sleep := float64(time.Duration(ttl) * time.Second)
		sleep = sleep * (.60 + mathRand.Float64()*0.1)
		sleepDuration := time.Duration(sleep)

		select {
		case <-time.After(sleepDuration):
		case <-ctx.Done():
			return nil
		case <-r.Stop:
			// Should we try to rewrap before stopping?
			r.Done <- nil
		}
	}
}

// Encrypt ...
func (r *ResponseEncrypter) Encrypt(ctx context.Context, plaintext, aad []byte) ([]byte, error) {
	blob, err := r.wrapper.Encrypt(ctx, plaintext, aad)
	if err != nil {
		return nil, err
	}
	return blob.Ciphertext, nil
}

// Decrypt ...
func (r *ResponseEncrypter) Decrypt(ctx context.Context, ciphertext, aad []byte) ([]byte, error) {
	blob := &wrapping.EncryptedBlobInfo{
		Ciphertext: ciphertext,
		KeyInfo: &wrapping.KeyInfo{
			KeyID: KeyID,
		},
	}
	return r.wrapper.Decrypt(ctx, blob, aad)
}

func (r *ResponseEncrypter) wrapForStorage() error {
	ttl := ResponseWrappedTokenTTL
	if r.ttl != "" {
		ttl = r.ttl
	}

	r.Client.AddHeader("X-Vault-Wrap-TTL", ttl)
	token, err := r.wrap()
	if err != nil {
		return err
	}

	r.token = []byte(token)
	return nil
}

func (r *ResponseEncrypter) wrap() (string, error) {
	b64Key := base64.StdEncoding.EncodeToString(r.wrapper.GetKeyBytes())
	data := map[string]interface{}{"key": b64Key}

	secret, err := r.Client.Logical().Write("/sys/wrapping/wrap", data)
	if err != nil {
		return "", err
	}
	return secret.WrapInfo.Token, nil
}

func (r *ResponseEncrypter) unwrap() ([]byte, error) {
	// Clear any previous headers set else Vault might
	// treat this as a wrap request (ie "X-Vault-Wrap-TTL")
	r.Client.SetHeaders(nil)
	secret, err := r.Client.Logical().Unwrap(string(r.token))
	if err != nil {
		return nil, err
	}

	key, ok := secret.Data["key"]
	if !ok {
		return nil, fmt.Errorf("key not found in unwrap response")
	}
	return base64.StdEncoding.DecodeString(key.(string))
}

func (r *ResponseEncrypter) rewrap() (string, int, error) {
	data := map[string]interface{}{"token": string(r.token)}
	secret, err := r.Client.Logical().Write("/sys/wrapping/rewrap", data)
	if err != nil {
		return "", -1, err
	}
	return secret.WrapInfo.Token, secret.WrapInfo.TTL, nil
}
