// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strconv"
	"sync"
	"testing"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"github.com/tink-crypto/tink-go/v2/kwp/subtle"
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
	"hmac",
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

const (
	nssFormattedEd25519Key = "MGcCAQAwFAYHKoZIzj0CAQYJKwYBBAHaRw8BBEwwSgIBAQQgfJm5R+LK4FMwGzOpemTBXksimEVOVCE8QeC+XBBfNU+hIwMhADaif7IhYx46IHcRTy1z8LeyhABep+UB8Da6olMZGx0i"
	rsaPSSFormattedKey     = "MIIEvAIBADALBgkqhkiG9w0BAQoEggSoMIIEpAIBAAKCAQEAiFXSBaicB534+2qMZTVzQHMjuhb4NM9hi5H4EAFiYHEBuvm2BAk58NdBK3wiMq/p7Ewu5NQI0gJ7GlcV1MBU94U6MEmWNd0ztmlz37esEDuaCDhmLEBHKRzs8Om0bY9vczcNwcnRIYusP2KMxon3Gv2C86M2Jahig70AIq0E9C7esfrlYxFnoxUfO09XyYfiHlZY59+/dhyULp/RDIvaQ0/DqSSnYmXw8vRQ1gp6DqIzxx3j8ikUrpE7MK6348keFQj1eb83Z5w8qgIdceHHH4wbIAW7qWCPJ/vIJp8Pe1NEanlef61pDut2YcljvN79ccjX/QyqwqYv6xX2uzSlpQIDAQABAoIBACtpBCAoIVJtkv9e3EhHniR55PjWYn7SP5GEz3MtNalWokHqS/H6DBhrOcWCV5NDHx1N3qqe9xYDkzX+X6Wn/gX4RmBkte79uX8OEca8wY1DpRaT+riBWQc2vh0xlPFDuC177KX1QGFJi3V9SCzZdjSCXyV7pPyVopSm4/mmlMq5ANfN8bcHAtcArP7vPzEdckJqurjwHyzsUZJa9sk3OL3rBkKy5bmoPebE1ZQ7C+9eA4u9MKSy95WpTiqMe3rRhvr6zj4bzEvzS9M4r2EdwgAn4FyDwtGdOqtfbtSLTikb73f4MSINnWbt3YPBfRC4PGjWXIN2sMG5XYC3KH+RKbsCgYEAu0HOFInH8OtWiUY0aqRKZuo7lrBczNa5gnce3ZYnNkfrPlu1Xp0SjUkEWukznBLO0N9lvG9j3ksUDTQlPoKarJb9uf/1H0tYHhHm6mP8mH87yfVn2bLb3VPeIQYb+MXnDrwNVCAtxhuHlpnXJPldeuVKeRigHUNIEs76UMiiLqMCgYEAumJxm5NrKk0LXUQmeZolLh0lM/shg8zW7Vi3Ksz5Pe4Pcmg+hTbHjZuJwK6HesljEA0JDNkS0+5hkqiS5UDnj94XfDbi08/kKbPYA12GPVSRNTJxL8q70rFnEUZuMBeL0SKMPhEfR2z5TDDZUBoO6HBUUwgJAij1EsXrBAb0BxcCgYBKS3eKKohLi/PPjy0oynpCjtiJlvuawe7kVoLGg9aW8L3jBdvV6Bf+OmQh9bhmSggIUzo4IzHKdptECdZlEMhxhY6xh14nxmr1s0Cc6oLDtmdwX4+OjioxjB7rl1Ltxwc/j1jycbn3ieCn3e3AW7e9FNARb7XHJnSoEbq65n+CZQKBgQChLPozYAL/HIrkR0fCRmM6gmemkNeFo0CFFP+oWoJ6ZIAlHjJafmmIcmVoI0TzEG3C9pLJ8nmOnYjxCyekakEUryi9+LSkGBWlXmlBV8H7DUNYrlskyfssEs8fKDmnCuWUn3yJO8NBv+HBWkjCNRaJOIIjH0KzBHoRludJnz2tVwKBgQCsQF5lvcXefNfQojbhF+9NfyhvAc7EsMTXQhP9HEj0wVqTuuqyGyu8meXEkcQPRl6yD/yZKuMREDNNck4KV2fdGekBsh8zBgpxdHQ2DcbfxZfNgv3yoX3f0grb/ApQNJb3DVW9FVRigue8XPzFOFX/demJmkUnTg3zGFnXLXjgxg=="
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

func TestTransit_ImportNSSEd25519Key(t *testing.T) {
	generateKeys(t)
	b, s, obsRecorder := createBackendWithObservationRecorder(t)

	wrappingKey, err := b.getWrappingKey(context.Background(), s)
	if err != nil || wrappingKey == nil {
		t.Fatalf("failed to retrieve public wrapping key: %s", err)
	}
	privWrappingKey := wrappingKey.Keys[strconv.Itoa(wrappingKey.LatestVersion)].RSAKey
	pubWrappingKey := &privWrappingKey.PublicKey

	rawPKCS8, err := base64.StdEncoding.DecodeString(nssFormattedEd25519Key)
	if err != nil {
		t.Fatalf("failed to parse nss base64: %v", err)
	}

	blob := wrapTargetPKCS8ForImport(t, pubWrappingKey, rawPKCS8, "SHA256")
	req := &logical.Request{
		Storage:   s,
		Operation: logical.UpdateOperation,
		Path:      "keys/nss-ed25519/import",
		Data: map[string]interface{}{
			"ciphertext": blob,
			"type":       "ed25519",
		},
	}

	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("failed to import NSS-formatted Ed25519 key: %v", err)
	}

	// Verify observation was recorded
	importObservations := obsRecorder.ObservationsByType(ObservationTypeTransitKeyImport)
	require.Len(t, importObservations, 1)
	require.Equal(t, "nss-ed25519", importObservations[0].Data["key_name"])
}

func TestTransit_ImportRSAPSS(t *testing.T) {
	generateKeys(t)
	b, s, obsRecorder := createBackendWithObservationRecorder(t)

	wrappingKey, err := b.getWrappingKey(context.Background(), s)
	if err != nil || wrappingKey == nil {
		t.Fatalf("failed to retrieve public wrapping key: %s", err)
	}
	privWrappingKey := wrappingKey.Keys[strconv.Itoa(wrappingKey.LatestVersion)].RSAKey
	pubWrappingKey := &privWrappingKey.PublicKey

	rawPKCS8, err := base64.StdEncoding.DecodeString(rsaPSSFormattedKey)
	if err != nil {
		t.Fatalf("failed to parse rsa-pss base64: %v", err)
	}

	blob := wrapTargetPKCS8ForImport(t, pubWrappingKey, rawPKCS8, "SHA256")
	req := &logical.Request{
		Storage:   s,
		Operation: logical.UpdateOperation,
		Path:      "keys/rsa-pss/import",
		Data: map[string]interface{}{
			"ciphertext": blob,
			"type":       "rsa-2048",
		},
	}

	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("failed to import RSA-PSS private key: %v", err)
	}

	importObservations := obsRecorder.ObservationsByType(ObservationTypeTransitKeyImport)
	require.Len(t, importObservations, 1)
	require.Equal(t, "rsa-pss", importObservations[0].Data["key_name"])
}

func TestTransit_Import(t *testing.T) {
	generateKeys(t)
	b, s, obsRecorder := createBackendWithObservationRecorder(t)
	checkImportObservation := func(t *testing.T, keyName string) {
		t.Helper()
		obs := obsRecorder.LastObservationOfType(ObservationTypeTransitKeyImport)
		require.NotNil(t, obs)
		require.Equal(t, keyName, obs.Data["key_name"])
	}
	t.Run(
		"import into a key fails before wrapping key is read",
		func(t *testing.T) {
			fakeWrappingKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 4096)
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
					checkImportObservation(t, keyID)
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

			checkImportObservation(t, keyID)
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

			obs := obsRecorder.LastObservationOfType(ObservationTypeTransitKeyRotateSuccess)
			require.NotNil(t, obs)
			require.Equal(t, obs.Data["key_name"], keyID)
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

			checkImportObservation(t, keyID)

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

	t.Run(
		"import public key ed25519",
		func(t *testing.T) {
			keyType := "ed25519"
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}

			// Get keys
			privateKey := getKey(t, keyType)
			publicKeyBytes, err := getPublicKey(privateKey, keyType)
			if err != nil {
				t.Fatal(err)
			}

			// Import key
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"public_key": publicKeyBytes,
					"type":       keyType,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to import ed25519 key: %v", err)
			}
			checkImportObservation(t, keyID)
		})

	t.Run(
		"import public key ecdsa",
		func(t *testing.T) {
			keyType := "ecdsa-p256"
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}

			// Get keys
			privateKey := getKey(t, keyType)
			publicKeyBytes, err := getPublicKey(privateKey, keyType)
			if err != nil {
				t.Fatal(err)
			}

			// Import key
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"public_key": publicKeyBytes,
					"type":       keyType,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to import public key: %s", err)
			}
			checkImportObservation(t, keyID)
		})
}

func TestTransit_ImportVersion(t *testing.T) {
	generateKeys(t)
	b, s, obsRecorder := createBackendWithObservationRecorder(t)

	t.Run(
		"import into a key version fails before wrapping key is read",
		func(t *testing.T) {
			fakeWrappingKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 4096)
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

	t.Run(
		"import rsa public key and update version with private counterpart",
		func(t *testing.T) {
			keyType := "rsa-2048"
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}

			// Get keys
			privateKey := getKey(t, keyType)
			importBlob := wrapTargetKeyForImport(t, pubWrappingKey, privateKey, keyType, "SHA256")
			publicKeyBytes, err := getPublicKey(privateKey, keyType)
			if err != nil {
				t.Fatal(err)
			}

			// Import RSA public key
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"public_key": publicKeyBytes,
					"type":       keyType,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to import public key: %s", err)
			}

			// Update version - import RSA private key
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import_version", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to update key: %s", err)
			}

			obs := obsRecorder.LastObservationOfType(ObservationTypeTransitKeyImport)
			require.NotNil(t, obs)
			require.Equal(t, keyID, obs.Data["key_name"])
			require.Equal(t, keyType, obs.Data["type"])
			require.NotContains(t, obs.Data, "import_version")
		},
	)
}

func TestTransit_ImportVersionWithPublicKeys(t *testing.T) {
	generateKeys(t)
	b, s, obsRecorder := createBackendWithObservationRecorder(t)

	// Retrieve public wrapping key
	wrappingKey, err := b.getWrappingKey(context.Background(), s)
	if err != nil || wrappingKey == nil {
		t.Fatalf("failed to retrieve public wrapping key: %s", err)
	}
	privWrappingKey := wrappingKey.Keys[strconv.Itoa(wrappingKey.LatestVersion)].RSAKey
	pubWrappingKey := &privWrappingKey.PublicKey

	// Import a public key then import private should give us one key
	t.Run(
		"import rsa public key and update version with private counterpart",
		func(t *testing.T) {
			keyType := "ecdsa-p256"
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}

			// Get keys
			privateKey := getKey(t, keyType)
			importBlob := wrapTargetKeyForImport(t, pubWrappingKey, privateKey, keyType, "SHA256")
			publicKeyBytes, err := getPublicKey(privateKey, keyType)
			if err != nil {
				t.Fatal(err)
			}

			// Import EC public key
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"public_key": publicKeyBytes,
					"type":       keyType,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to import public key: %s", err)
			}

			// Update version - import EC private key
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import_version", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to update key: %s", err)
			}

			// We should have one key on export
			req = &logical.Request{
				Storage:   s,
				Operation: logical.ReadOperation,
				Path:      fmt.Sprintf("export/public-key/%s", keyID),
			}
			resp, err := b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to export key: %s", err)
			}

			if len(resp.Data["keys"].(map[string]string)) != 1 {
				t.Fatalf("expected 1 key but got %v: %v", len(resp.Data["keys"].(map[string]string)), resp)
			}
		},
	)

	// Import a private and then public should give us two keys
	t.Run(
		"import ec private key and then its public counterpart",
		func(t *testing.T) {
			keyType := "ecdsa-p256"
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}

			// Get keys
			privateKey := getKey(t, keyType)
			importBlob := wrapTargetKeyForImport(t, pubWrappingKey, privateKey, keyType, "SHA256")
			publicKeyBytes, err := getPublicKey(privateKey, keyType)
			if err != nil {
				t.Fatal(err)
			}

			// Import EC private key
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob,
					"type":       keyType,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to update key: %s", err)
			}

			// Update version - Import EC public key
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import_version", keyID),
				Data: map[string]interface{}{
					"public_key": publicKeyBytes,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to import public key: %s", err)
			}

			// We should have two keys on export
			req = &logical.Request{
				Storage:   s,
				Operation: logical.ReadOperation,
				Path:      fmt.Sprintf("export/public-key/%s", keyID),
			}
			resp, err := b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to export key: %s", err)
			}

			if len(resp.Data["keys"].(map[string]string)) != 2 {
				t.Fatalf("expected 2 key but got %v: %v", len(resp.Data["keys"].(map[string]string)), resp)
			}
		},
	)

	// Import a public and another public should allow us to insert two private key.
	t.Run(
		"import two public keys and two private keys in reverse order",
		func(t *testing.T) {
			keyType := "ecdsa-p256"
			keyID, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed to generate key ID: %s", err)
			}

			// Get keys
			privateKey1 := getKey(t, keyType)
			importBlob1 := wrapTargetKeyForImport(t, pubWrappingKey, privateKey1, keyType, "SHA256")
			publicKeyBytes1, err := getPublicKey(privateKey1, keyType)
			if err != nil {
				t.Fatal(err)
			}

			privateKey2, err := generateKey(keyType)
			if err != nil {
				t.Fatal(err)
			}
			importBlob2 := wrapTargetKeyForImport(t, pubWrappingKey, privateKey2, keyType, "SHA256")
			publicKeyBytes2, err := getPublicKey(privateKey2, keyType)
			if err != nil {
				t.Fatal(err)
			}

			// Import EC public key
			req := &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import", keyID),
				Data: map[string]interface{}{
					"public_key": publicKeyBytes1,
					"type":       keyType,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to update key: %s", err)
			}

			// Update version - Import second EC public key
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import_version", keyID),
				Data: map[string]interface{}{
					"public_key": publicKeyBytes2,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to import public key: %s", err)
			}

			// We should have two keys on export
			req = &logical.Request{
				Storage:   s,
				Operation: logical.ReadOperation,
				Path:      fmt.Sprintf("export/public-key/%s", keyID),
			}
			resp, err := b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to export key: %s", err)
			}

			if len(resp.Data["keys"].(map[string]string)) != 2 {
				t.Fatalf("expected 2 key but got %v: %v", len(resp.Data["keys"].(map[string]string)), resp)
			}

			// Import second private key first, with no options.
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import_version", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob2,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to import private key: %s", err)
			}

			// Import first private key second, with a version
			req = &logical.Request{
				Storage:   s,
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("keys/%s/import_version", keyID),
				Data: map[string]interface{}{
					"ciphertext": importBlob1,
					"version":    1,
				},
			}
			_, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to import private key: %s", err)
			}

			obs := obsRecorder.LastObservationOfType(ObservationTypeTransitKeyImport)
			require.NotNil(t, obs)
			require.Equal(t, keyID, obs.Data["key_name"])
			require.Equal(t, 1, obs.Data["import_version"])

			// We should still have two keys on export
			req = &logical.Request{
				Storage:   s,
				Operation: logical.ReadOperation,
				Path:      fmt.Sprintf("export/public-key/%s", keyID),
			}
			resp, err = b.HandleRequest(context.Background(), req)
			if err != nil {
				t.Fatalf("failed to export key: %s", err)
			}

			if len(resp.Data["keys"].(map[string]string)) != 2 {
				t.Fatalf("expected 2 key but got %v: %v", len(resp.Data["keys"].(map[string]string)), resp)
			}
		},
	)
}

func wrapTargetKeyForImport(t *testing.T, wrappingKey *rsa.PublicKey, targetKey interface{}, targetKeyType string, hashFnName string) string {
	t.Helper()

	// Format target key for wrapping
	var preppedTargetKey []byte
	var ok bool
	var err error
	switch targetKeyType {
	case "aes128-gcm96", "aes256-gcm96", "chacha20-poly1305", "hmac":
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

	return wrapTargetPKCS8ForImport(t, wrappingKey, preppedTargetKey, hashFnName)
}

func wrapTargetPKCS8ForImport(t *testing.T, wrappingKey *rsa.PublicKey, preppedTargetKey []byte, hashFnName string) string {
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
	case "aes256-gcm96", "hmac":
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
		return cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	case "rsa-3072":
		return cryptoutil.GenerateRSAKey(rand.Reader, 3072)
	case "rsa-4096":
		return cryptoutil.GenerateRSAKey(rand.Reader, 4096)
	default:
		return nil, fmt.Errorf("failed to generate unsupported key type: %s", keyType)
	}
}

func getPublicKey(privateKey crypto.PrivateKey, keyType string) ([]byte, error) {
	var publicKey crypto.PublicKey
	var publicKeyBytes []byte
	switch keyType {
	case "rsa-2048", "rsa-3072", "rsa-4096":
		publicKey = privateKey.(*rsa.PrivateKey).Public()
	case "ecdsa-p256", "ecdsa-p384", "ecdsa-p521":
		publicKey = privateKey.(*ecdsa.PrivateKey).Public()
	case "ed25519":
		publicKey = privateKey.(ed25519.PrivateKey).Public()
	default:
		return publicKeyBytes, fmt.Errorf("failed to get public key from %s key", keyType)
	}

	publicKeyBytes, err := publicKeyToBytes(publicKey)
	if err != nil {
		return publicKeyBytes, err
	}

	return publicKeyBytes, nil
}

func publicKeyToBytes(publicKey crypto.PublicKey) ([]byte, error) {
	var publicKeyBytesPem []byte
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return publicKeyBytesPem, fmt.Errorf("failed to marshal public key: %s", err)
	}

	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	return pem.EncodeToMemory(pemBlock), nil
}
