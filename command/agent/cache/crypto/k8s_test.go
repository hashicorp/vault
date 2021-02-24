package crypto

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestCrypto_KubernetesNewKey(t *testing.T) {
	k8sKey, err := NewK8s([]byte{})
	if err != nil {
		t.Fatalf(fmt.Sprintf("unexpected error: %s", err))
	}

	key := k8sKey.GetKey()
	if key == nil {
		t.Fatalf(fmt.Sprintf("key is nil, it shouldn't be: %s", key))
	}

	persistentKey, _ := k8sKey.GetPersistentKey()
	if persistentKey == nil {
		t.Fatalf(fmt.Sprintf("key is nil, it shouldn't be: %s", persistentKey))
	}

	if string(key) != string(persistentKey) {
		t.Fatalf("keys don't match, they should: key: %s, persistentKey: %s", key, persistentKey)
	}

	plaintextInput := []byte("test")
	aad := []byte("kubernetes")

	ciphertext, err := k8sKey.Encrypt(nil, plaintextInput, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if ciphertext == nil {
		t.Fatalf("ciphertext nil, it shouldn't be")
	}

	plaintext, err := k8sKey.Decrypt(nil, ciphertext, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if string(plaintext) != string(plaintextInput) {
		t.Fatalf("expected %s, got %s", plaintextInput, plaintext)
	}
}

func TestCrypto_KubernetesExistingKey(t *testing.T) {
	rootKey := make([]byte, 32)
	n, err := rand.Read(rootKey)
	if err != nil {
		t.Fatal(err)
	}
	if n != 32 {
		t.Fatal(n)
	}

	k8sKey, err := NewK8s(rootKey)
	if err != nil {
		t.Fatalf(fmt.Sprintf("unexpected error: %s", err))
	}

	key := k8sKey.GetKey()
	if key == nil {
		t.Fatalf(fmt.Sprintf("key is nil, it shouldn't be: %s", key))
	}

	if string(key) != string(rootKey) {
		t.Fatalf(fmt.Sprintf("expected keys to be the same, they weren't: expected: %s, got: %s", rootKey, key))
	}

	persistentKey, _ := k8sKey.GetPersistentKey()
	if persistentKey == nil {
		t.Fatalf("key is nil, it shouldn't be")
	}

	if string(persistentKey) != string(rootKey) {
		t.Fatalf(fmt.Sprintf("expected keys to be the same, they weren't: expected: %s, got: %s", rootKey, persistentKey))
	}

	if string(key) != string(persistentKey) {
		t.Fatalf(fmt.Sprintf("expected keys to be the same, they weren't: %s %s", rootKey, persistentKey))
	}

	plaintextInput := []byte("test")
	aad := []byte("kubernetes")

	ciphertext, err := k8sKey.Encrypt(nil, plaintextInput, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if ciphertext == nil {
		t.Fatalf("ciphertext nil, it shouldn't be")
	}

	plaintext, err := k8sKey.Decrypt(nil, ciphertext, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if string(plaintext) != string(plaintextInput) {
		t.Fatalf("expected %s, got %s", plaintextInput, plaintext)
	}
}

func TestCrypto_KubernetesPassGeneratedKey(t *testing.T) {
	k8sFirstKey, err := NewK8s([]byte{})
	if err != nil {
		t.Fatalf(fmt.Sprintf("unexpected error: %s", err))
	}

	firstPersistentKey := k8sFirstKey.GetKey()
	if firstPersistentKey == nil {
		t.Fatalf(fmt.Sprintf("key is nil, it shouldn't be: %s", firstPersistentKey))
	}

	plaintextInput := []byte("test")
	aad := []byte("kubernetes")

	ciphertext, err := k8sFirstKey.Encrypt(nil, plaintextInput, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if ciphertext == nil {
		t.Fatalf("ciphertext nil, it shouldn't be")
	}

	k8sLoadedKey, err := NewK8s(firstPersistentKey)
	if err != nil {
		t.Fatalf(fmt.Sprintf("unexpected error: %s", err))
	}

	loadedKey, _ := k8sLoadedKey.GetPersistentKey()
	if loadedKey == nil {
		t.Fatalf(fmt.Sprintf("key is nil, it shouldn't be: %s", loadedKey))
	}

	if string(loadedKey) != string(firstPersistentKey) {
		t.Fatalf(fmt.Sprintf("keys do not match"))
	}

	plaintext, err := k8sLoadedKey.Decrypt(nil, ciphertext, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if string(plaintext) != string(plaintextInput) {
		t.Fatalf("expected %s, got %s", plaintextInput, plaintext)
	}
}
