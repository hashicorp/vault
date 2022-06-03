package transit

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/google/tink/go/kwp/subtle"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
)

var keyTypes = []string{
	"aes256-gcm96",
	"aes128-gcm96",
	"chacha20-poly1305",
	"ed25519",
	"ecdsa-p256",
	"ecdsa-p384",
	"ecdsa-p521",
	"rsa-2048",
	"rsa-3072",
	"rsa-4096",
}

var hashFns = []string{
	"SHA256",
	"SHA1",
	"SHA224",
	"SHA384",
	"SHA512",
}

var (
	keysLock sync.RWMutex
	keys     = map[string]interface{}{}
)

func generateKeys(t *testing.T) {
	t.Helper()

	keysLock.Lock()
	defer keysLock.Unlock()

	if len(keys) > 0 {
		return
	}

	for _, keyType := range keyTypes {
		key, err := generateKey(keyType)
		if err != nil {
			t.Fatalf("failed to generate %s key: %s", keyType, err)
		}
		keys[keyType] = key
	}
}

func getKey(t *testing.T, keyType string) interface{} {
	t.Helper()

	keysLock.RLock()
	defer keysLock.RUnlock()

	key, ok := keys[keyType]
	if !ok {
		t.Fatalf("no pre-generated key of type: %s", keyType)
	}

	return key
}

func TestTransit_Import(t *testing.T) {
	generateKeys(t)
	b, s := createBackendWithStorage(t)

	t.Run(
		"import into a key fails before wrapping key is read",
		func(t *testing.T) {
			fakeWrappingKey, err := rsa.GenerateKey(rand.Reader, 4096)
			if err != nil {
				t.Fatalf("failed to generate fake wrapping key: %s", err)
			}
			// Roll an AES256 key and import
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}
			targetKey := getKey(t, "aes256-gcm96")
			importBlob := wrapTargetKeyForImport(t, &fakeWrappingKey.PublicKey, targetKey, "aes256-gcm96", "SHA256")
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err == nil {
				t.Fatal("import prior to wrapping key generation incorrectly succeeded")
			}
		},
	)

	// Retrieve public wrapping key
	wrappingKey, err := b.getWrappingKey(context.Background(), s)
	if err != nil || wrappingKey == nil {
		t.Fatalf("failed to retrieve public wrapping key: %s", err)
	}
	privWrappingKey := wrappingKey.Keys[strconv.Itoa(wrappingKey.LatestVersion)].RSAKey
	pubWrappingKey := &privWrappingKey.PublicKey

	t.Run(
		"import into an existing key fails",
		func(t *testing.T) {
			// Generate a key ID
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate a key ID: %s", err)
			}

			// Create an AES256 key within Transit
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s", keyID),
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("unexpected error creating key: %s", err)
			}

			targetKey := getKey(t, "aes256-gcm96")
			importBlob := wrapTargetKeyForImport(t, pubWrappingKey, targetKey, "aes256-gcm96", "SHA256")
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err == nil {
				t.Fatal("import into an existing key incorrectly succeeded")
			}
		},
	)

	for _, keyType := range keyTypes {
		priv := getKey(t, keyType)
		for _, hashFn := range hashFns {
			t.Run(
				fmt.Sprintf("%s/%s", keyType, hashFn),
				func(t *testing.T) {
					keyID, err := uuid.GenerateUUID()
					if err != nil {
						t.Fatalf("failed to generate key ID: %s", err)
					}
					importBlob := wrapTargetKeyForImport(t, pubWrappingKey, priv, keyType, hashFn)
					req := &logical.Request{
						Storage:   s,
						Operation: logical.UpdateOperation,
						Path:      fmt.Sprintf("keys/%s/import", keyID),
						Data: map[string]interface{}{
							"type":          keyType,
							"hash_function": hashFn,
							"ciphertext":    importBlob,
						},
					}
					_, err = b.HandleRequest(context.Background(), req)
					if err != nil {
						t.Fatalf("failed to import valid key: %s", err)
					}
				},
			)

			// Shouldn't need to test every combination of key and hash function
			if keyType != "aes256-gcm96" {
				break
			}
		}
	}

	failures := []struct {
		name       string
		ciphertext interface{}
		keyType    interface{}
		hashFn     interface{}
	}{
		{
			name: "nil ciphertext",
		},
		{
			name:       "empty string ciphertext",
			ciphertext: "",
		},
		{
			name:       "ciphertext not base64",
			ciphertext: "this isn't correct",
		},
		{
			name:       "ciphertext too short",
			ciphertext: "ZmFrZSBjaXBoZXJ0ZXh0Cg",
		},
		{
			name:    "invalid key type",
			keyType: "fake-key-type",
		},
		{
			name:   "invalid hash function",
			hashFn: "fake-hash-fn",
		},
	}
	for _, tt := range failures {
		t.Run(
			tt.name,
			func(t *testing.T) {
				keyID, err := uuid.GenerateUUID()
				if err != nil {
					t.Fatalf("failed to generate key ID: %s", err)
				}
				req := &logical.Request{
					Storage:   s,
					Operation: logical.UpdateOperation,
					Path:      fmt.Sprintf("keys/%s/import", keyID),
					Data:      map[string]interface{}{},
				}
				if tt.ciphertext != nil {
					req.Data["ciphertext"] = tt.ciphertext
				}
				if tt.keyType != nil {
					req.Data["type"] = tt.keyType
				}
				if tt.hashFn != nil {
					req.Data["hash_function"] = tt.hashFn
				}
				_, err = b.HandleRequest(context.Background(), req)
				if err == nil {
					t.Fatal("invalid import request incorrectly succeeded")
				}
			},
		)
	}

	t.Run(
		"disallow import of convergent keys",
		func(t *testing.T) {
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}
			targetKey := getKey(t, "aes256-gcm96")
			importBlob := wrapTargetKeyForImport(t, pubWrappingKey, targetKey, "aes256-gcm96", "SHA256")
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"convergent_encryption": true,
					"ciphertext":            importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err == nil {
				t.Fatal("import of convergent key incorrectly succeeded")
			}
		},
	)

	t.Run(
		"allow_rotation=true enables rotation within vault",
		func(t *testing.T) {
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}
			targetKey := getKey(t, "aes256-gcm96")

			// Import key
			importBlob := wrapTargetKeyForImport(t, pubWrappingKey, targetKey, "aes256-gcm96", "SHA256")
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"allow_rotation": true,
					"ciphertext":     importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to import key: %s", err)
			}

			// Rotate key
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/rotate", keyID),
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to rotate key: %s", err)
			}
		},
	)

	t.Run(
		"allow_rotation=false disables rotation within vault",
		func(t *testing.T) {
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}
			targetKey := getKey(t, "aes256-gcm96")

			// Import key
			importBlob := wrapTargetKeyForImport(t, pubWrappingKey, targetKey, "aes256-gcm96", "SHA256")
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"allow_rotation": false,
					"ciphertext":     importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to import key: %s", err)
			}

			// Rotate key
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/rotate", keyID),
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err == nil {
				t.Fatal("rotation of key with allow_rotation incorrectly succeeded")
			}
		},
	)
}

func TestTransit_ImportVersion(t *testing.T) {
	generateKeys(t)
	b, s := createBackendWithStorage(t)

	t.Run(
		"import into a key version fails before wrapping key is read",
		func(t *testing.T) {
			fakeWrappingKey, err := rsa.GenerateKey(rand.Reader, 4096)
			if err != nil {
				t.Fatalf("failed to generate fake wrapping key: %s", err)
			}
			// Roll an AES256 key and import
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}
			targetKey := getKey(t, "aes256-gcm96")
			importBlob := wrapTargetKeyForImport(t, &fakeWrappingKey.PublicKey, targetKey, "aes256-gcm96", "SHA256")
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import_version", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err == nil {
				t.Fatal("import_version prior to wrapping key generation incorrectly succeeded")
			}
		},
	)

	// Retrieve public wrapping key
	wrappingKey, err := b.getWrappingKey(context.Background(), s)
	if err != nil || wrappingKey == nil {
		t.Fatalf("failed to retrieve public wrapping key: %s", err)
	}
	privWrappingKey := wrappingKey.Keys[strconv.Itoa(wrappingKey.LatestVersion)].RSAKey
	pubWrappingKey := &privWrappingKey.PublicKey

	t.Run(
		"import into a non-existent key fails",
		func(t *testing.T) {
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}
			targetKey := getKey(t, "aes256-gcm96")
			importBlob := wrapTargetKeyForImport(t, pubWrappingKey, targetKey, "aes256-gcm96", "SHA256")
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import_version", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err == nil {
				t.Fatal("import_version into a non-existent key incorrectly succeeded")
			}
		},
	)

	t.Run(
		"import into an internally-generated key fails",
		func(t *testing.T) {
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}

			// Roll a key within Transit
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s", keyID),
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to generate a key within transit: %s", err)
			}

			// Attempt to import into newly generated key
			targetKey := getKey(t, "aes256-gcm96")
			importBlob := wrapTargetKeyForImport(t, pubWrappingKey, targetKey, "aes256-gcm96", "SHA256")
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import_version", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err == nil {
				t.Fatal("import_version into an internally-generated key incorrectly succeeded")
			}
		},
	)

	t.Run(
		"imported key version type must match existing key type",
		func(t *testing.T) {
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}

			// Import an RSA key
			targetKey := getKey(t, "rsa-2048")
			importBlob := wrapTargetKeyForImport(t, pubWrappingKey, targetKey, "rsa-2048", "SHA256")
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob,
					"type":       "rsa-2048",
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to generate a key within transit: %s", err)
			}

			// Attempt to import an AES key version into existing RSA key
			targetKey = getKey(t, "aes256-gcm96")
			importBlob = wrapTargetKeyForImport(t, pubWrappingKey, targetKey, "aes256-gcm96", "SHA256")
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import_version", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err == nil {
				t.Fatal("import_version into a key of a different type incorrectly succeeded")
			}
		},
	)
}

func wrapTargetKeyForImport(t *testing.T, wrappingKey *rsa.PublicKey, targetKey interface{}, targetKeyType string, hashFnName string) string {
	t.Helper()

	// Generate an ephemeral AES-256 key
	ephKey, err := uuid.GenerateRandomBytes(32)
	if err != nil {
		t.Fatalf("failed to wrap target key for import: %s", err)
	}

	// Parse the hash function name into an actual function
	hashFn, err := parseHashFn(hashFnName)
	if err != nil {
		t.Fatalf("failed to wrap target key for import: %s", err)
	}

	// Wrap ephemeral AES key with public wrapping key
	ephKeyWrapped, err := rsa.EncryptOAEP(hashFn, rand.Reader, wrappingKey, ephKey, []byte{})
	if err != nil {
		t.Fatalf("failed to wrap target key for import: %s", err)
	}

	// Create KWP instance for wrapping target key
	kwp, err := subtle.NewKWP(ephKey)
	if err != nil {
		t.Fatalf("failed to wrap target key for import: %s", err)
	}

	// Format target key for wrapping
	var preppedTargetKey []byte
	var ok bool
	switch targetKeyType {
	case "aes128-gcm96", "aes256-gcm96", "chacha20-poly1305":
		preppedTargetKey, ok = targetKey.([]byte)
		if !ok {
			t.Fatal("failed to wrap target key for import: symmetric key not provided in byte format")
		}
	default:
		preppedTargetKey, err = x509.MarshalPKCS8PrivateKey(targetKey)
		if err != nil {
			t.Fatalf("failed to wrap target key for import: %s", err)
		}
	}

	// Wrap target key with KWP
	targetKeyWrapped, err := kwp.Wrap(preppedTargetKey)
	if err != nil {
		t.Fatalf("failed to wrap target key for import: %s", err)
	}

	// Combined wrapped keys into a single blob and base64 encode
	wrappedKeys := append(ephKeyWrapped, targetKeyWrapped...)
	return base64.StdEncoding.EncodeToString(wrappedKeys)
}

func generateKey(keyType string) (interface{}, error) {
	switch keyType {
	case "aes128-gcm96":
		return uuid.GenerateRandomBytes(16)
	case "aes256-gcm96":
		return uuid.GenerateRandomBytes(32)
	case "chacha20-poly1305":
		return uuid.GenerateRandomBytes(32)
	case "ed25519":
		_, priv, err := ed25519.GenerateKey(rand.Reader)
		return priv, err
	case "ecdsa-p256":
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "ecdsa-p384":
		return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "ecdsa-p521":
		return ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	case "rsa-2048":
		return rsa.GenerateKey(rand.Reader, 2048)
	case "rsa-3072":
		return rsa.GenerateKey(rand.Reader, 3072)
	case "rsa-4096":
		return rsa.GenerateKey(rand.Reader, 4096)
	default:
		return nil, fmt.Errorf("failed to generate unsupported key type: %s", keyType)
	}
}
