package transit

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"testing"

	"github.com/google/tink/go/kwp/subtle"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestTransit_Import(t *testing.T) {
	// Set up shared backend for subtests
	b, s := createBackendWithStorage(t)

	t.Run(
		"import into a key fails before wrapping key is generated",
		func(t *testing.T) {
			// Generate a fake wrapping key
			fakeWrappingKey, err := rsa.GenerateKey(rand.Reader, 4096)
			if err != nil {
				t.Fatalf("failed to generate fake wrapping key: %s", err)
			}
			// Roll an AES128 key and import
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}
			targetKey, err := uuid.GenerateRandomBytes(128)
			if err != nil {
				t.Fatalf("failed to generate target key: %s", err)
			}
			importBlob, err := wrapTargetKeyForImport(&fakeWrappingKey.PublicKey, targetKey, "aes128-gcm96")
			if err != nil {
				t.Fatalf("failed to wrap target key for import: %s", err)
			}
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"type":       "aes128-gcm96",
					"ciphertext": importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err == nil {
				t.Fatal("import prior to wrapping key generation incorrectly succeeded")
			}
		},
	)

	// Generate public wrapping key for import usage
	req := &logical.Request{
		Storage:   s,
		Operation: logical.ReadOperation,
		Path:      "wrapping_key",
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("error generating public wrapping key: %s", err)
	}
	if resp == nil || resp.Data == nil || resp.Data["public_key"] == nil {
		t.Fatal("expected non-nil public wrapping key response")
	}
	pubKeyPEM := resp.Data["public_key"]
	pubKeyBlock, _ := pem.Decode([]byte(pubKeyPEM.(string)))
	rawPubKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse public wrapping key: %s", err)
	}
	wrappingKey, ok := rawPubKey.(*rsa.PublicKey)
	if !ok || wrappingKey.Size() != 512 {
		t.Fatal("public wrapping key is not an RSA 4096 key")
	}

	t.Run(
		"import into existing key fails",
		func(t *testing.T) {
			// Generate a key ID
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate a key ID: %s", err)
			}

			// Create an AES128 key within Transit
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s", keyID),
				Data: map[string]interface{}{
					"type": "aes128-gcm96",
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("unexpected request error: %s", err)
			}

			// Roll an AES128 key and import
			targetKey, err := uuid.GenerateRandomBytes(128)
			if err != nil {
				t.Fatalf("failed to generate target key: %s", err)
			}
			importBlob, err := wrapTargetKeyForImport(wrappingKey, targetKey, "aes128-gcm96")
			if err != nil {
				t.Fatalf("failed to wrap target key for import: %s", err)
			}
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"type":       "aes128-gcm96",
					"ciphertext": importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err == nil {
				t.Fatal("import into an existing key incorrectly succeeded")
			}
		},
	)
}

func wrapTargetKeyForImport(wrappingKey *rsa.PublicKey, targetKey interface{}, targetKeyType string) (string, error) {
	// Generate ephemeral AES-256 key
	ephKey, err := uuid.GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}

	// Wrap ephemeral AES key with public wrapping key
	ephKeyWrapped, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, wrappingKey, ephKey, []byte{})
	if err != nil {
		return "", err
	}

	// Create KWP instance for wrapping target key
	kwp, err := subtle.NewKWP(ephKey)
	if err != nil {
		return "", err
	}

	// Format target key for wrapping
	var preppedTargetKey []byte
	var ok bool
	switch targetKeyType {
	case "aes128-gcm96", "aes256-gcm96", "chacha20-poly1305":
		preppedTargetKey, ok = targetKey.([]byte)
		if !ok {
			return "", errors.New("target key not provided in byte format")
		}
	default:
		preppedTargetKey, err = x509.MarshalPKCS8PrivateKey(targetKey)
		if err != nil {
			return "", err
		}
	}

	// Wrap target key with KWP
	targetKeyWrapped, err := kwp.Wrap(preppedTargetKey)
	if err != nil {
		return "", err
	}

	// Combine two wrapped keys into a single blob and base64 encode
	wrappedKeys := append(ephKeyWrapped, targetKeyWrapped...)
	return base64.RawURLEncoding.EncodeToString(wrappedKeys), nil
}
