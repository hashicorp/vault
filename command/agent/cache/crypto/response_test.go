package crypto

import (
	"context"
	"fmt"
	"testing"

	log "github.com/hashicorp/go-hclog"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestCrypto_ResponseWrappingNewKey(t *testing.T) {
	var err error
	coreConfig := &vault.CoreConfig{
		DisableMlock:       true,
		DisableCache:       true,
		Logger:             log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client

	responseWrappedKey, err := NewResponseEncrypter([]byte{}, client, ResponseWrappedTokenTTL)
	if err != nil {
		t.Fatalf(fmt.Sprintf("unexpected error: %s", err))
	}

	key := responseWrappedKey.GetKey()
	if key == nil {
		t.Fatalf(fmt.Sprintf("key is nil, it shouldn't be: %s", key))
	}

	plaintextInput := []byte("test")
	aad := []byte("")

	ciphertext, err := responseWrappedKey.Encrypt(nil, plaintextInput, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if ciphertext == nil {
		t.Fatalf("ciphertext nil, it shouldn't be")
	}

	plaintext, err := responseWrappedKey.Decrypt(nil, ciphertext, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if string(plaintext) != string(plaintextInput) {
		t.Fatalf("expected %s, got %s", plaintextInput, plaintext)
	}

	token, err := responseWrappedKey.GetPersistentKey()
	if err != nil {
		t.Fatalf("unxpected error: %s", err)
	}

	if token == nil {
		t.Fatalf("persistent token nil, it shouldn't be")
	}
}

func TestCrypto_ResponseWrappingExistingKey(t *testing.T) {
	var err error

	coreConfig := &vault.CoreConfig{
		DisableMlock:       true,
		DisableCache:       true,
		Logger:             log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client

	responseWrappedKey, err := NewResponseEncrypter([]byte{}, client, ResponseWrappedTokenTTL)
	if err != nil {
		t.Fatalf(fmt.Sprintf("unexpected error: %s", err))
	}

	key := responseWrappedKey.GetKey()
	if key == nil {
		t.Fatalf(fmt.Sprintf("key is nil, it shouldn't be: %s", key))
	}

	plaintextInput := []byte("test")
	aad := []byte("")

	ciphertext, err := responseWrappedKey.Encrypt(nil, plaintextInput, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if ciphertext == nil {
		t.Fatalf("ciphertext nil, it shouldn't be")
	}

	token, err := responseWrappedKey.GetPersistentKey()
	if err != nil {
		t.Fatalf("unxpected error: %s", err)
	}

	if token == nil {
		t.Fatalf("persistent token nil, it shouldn't be")
	}

	responseWrappedKeyExisting, err := NewResponseEncrypter(token, client, ResponseWrappedTokenTTL)
	if err != nil {
		t.Fatalf(fmt.Sprintf("unexpected error: %s", err))
	}

	existingKey := responseWrappedKeyExisting.GetKey()
	if existingKey == nil {
		t.Fatalf(fmt.Sprintf("key is nil, it shouldn't be: %s", existingKey))
	}

	if string(key) != string(existingKey) {
		t.Fatalf("keys don't match, they should: old: %s, new: %s", key, existingKey)
	}

	plaintext, err := responseWrappedKeyExisting.Decrypt(nil, ciphertext, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if string(plaintext) != string(plaintextInput) {
		t.Fatalf("expected %s, got %s", plaintextInput, plaintext)
	}
}

func TestCrypto_ResponseWrappingRenewer(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		DisableMlock:       true,
		DisableCache:       true,
		Logger:             log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client

	responseWrappedKey, err := NewResponseEncrypter([]byte{}, client, "5")
	if err != nil {
		t.Fatalf(fmt.Sprintf("unexpected error: %s", err))
	}

	key := responseWrappedKey.GetKey()
	if key == nil {
		t.Fatalf(fmt.Sprintf("key is nil, it shouldn't be: %s", key))
	}

	plaintextInput := []byte("test")
	aad := []byte("")

	ciphertext, err := responseWrappedKey.Encrypt(nil, plaintextInput, aad)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if ciphertext == nil {
		t.Fatalf("ciphertext nil, it shouldn't be")
	}

	token, err := responseWrappedKey.GetPersistentKey()
	if err != nil {
		t.Fatalf("unxpected error: %s", err)
	}

	if token == nil {
		t.Fatalf("persistent token nil, it shouldn't be")
	}

	if !responseWrappedKey.Renewable() {
		t.Fatalf("response wrapped key isn't renewable, it should be")
	}

	ctx := context.Background()
	go responseWrappedKey.Renewer(ctx)

	firstToken := responseWrappedKey.token
	if firstToken == nil {
		t.Fatalf("first wrapped token is nil, it shouldn't be")
	}

	<-responseWrappedKey.Notify

	secondToken := responseWrappedKey.token
	if secondToken == nil {
		t.Fatalf("second wrapped token is nil, it shouldn't be")
	}

	if string(firstToken) == string(secondToken) {
		t.Fatalf("first token and rewrapped token mathch, they shouldn't")
	}
}
