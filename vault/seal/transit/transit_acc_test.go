package transit

import (
	"context"
	"fmt"
	"path"
	"reflect"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/ory/dockertest"
)

func TestTransitSeal_Lifecycle(t *testing.T) {
	cleanup, retAddress, token, mountPath, keyName, _ := prepareTestContainer(t)
	defer cleanup()

	sealConfig := map[string]string{
		"address":    retAddress,
		"token":      token,
		"mount_path": mountPath,
		"key_name":   keyName,
	}
	s := NewSeal(logging.NewVaultLogger(log.Trace))
	_, err := s.SetConfig(sealConfig)
	if err != nil {
		t.Fatalf("error setting seal config: %v", err)
	}

	// Test Encrypt and Decrypt calls
	input := []byte("foo")
	swi, err := s.Encrypt(context.Background(), input)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	pt, err := s.Decrypt(context.Background(), swi)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	if !reflect.DeepEqual(input, pt) {
		t.Fatalf("expected %s, got %s", input, pt)
	}
}

func TestTransitSeal_TokenRenewal(t *testing.T) {
	cleanup, retAddress, token, mountPath, keyName, tlsConfig := prepareTestContainer(t)
	defer cleanup()

	clientConfig := &api.Config{
		Address: retAddress,
	}
	if err := clientConfig.ConfigureTLS(tlsConfig); err != nil {
		t.Fatalf("err: %s", err)
	}

	remoteClient, err := api.NewClient(clientConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	remoteClient.SetToken(token)

	req := &api.TokenCreateRequest{
		Period: "5s",
	}
	rsp, err := remoteClient.Auth().Token().Create(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	sealConfig := map[string]string{
		"address":    retAddress,
		"token":      rsp.Auth.ClientToken,
		"mount_path": mountPath,
		"key_name":   keyName,
	}
	s := NewSeal(logging.NewVaultLogger(log.Trace))
	_, err = s.SetConfig(sealConfig)
	if err != nil {
		t.Fatalf("error setting seal config: %v", err)
	}

	time.Sleep(7 * time.Second)

	// Test Encrypt and Decrypt calls
	input := []byte("foo")
	swi, err := s.Encrypt(context.Background(), input)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	pt, err := s.Decrypt(context.Background(), swi)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	if !reflect.DeepEqual(input, pt) {
		t.Fatalf("expected %s, got %s", input, pt)
	}
}

func prepareTestContainer(t *testing.T) (cleanup func(), retAddress, token, mountPath, keyName string, tlsConfig *api.TLSConfig) {
	testToken, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testMountPath, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testKeyName, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	dockerOptions := &dockertest.RunOptions{
		Repository: "vault",
		Tag:        "latest",
		Cmd: []string{"server", "-log-level=trace", "-dev", fmt.Sprintf("-dev-root-token-id=%s", testToken),
			"-dev-listen-address=0.0.0.0:8200"},
	}
	resource, err := pool.RunWithOptions(dockerOptions)
	if err != nil {
		t.Fatalf("Could not start local Vault docker container: %s", err)
	}

	cleanup = func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	retAddress = fmt.Sprintf("http://127.0.0.1:%s", resource.GetPort("8200/tcp"))
	tlsConfig = &api.TLSConfig{
		Insecure: true,
	}

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		vaultConfig := api.DefaultConfig()
		vaultConfig.Address = retAddress
		if err := vaultConfig.ConfigureTLS(tlsConfig); err != nil {
			return err
		}
		vault, err := api.NewClient(vaultConfig)
		if err != nil {
			return err
		}
		vault.SetToken(testToken)

		// Set up transit
		if err := vault.Sys().Mount(testMountPath, &api.MountInput{
			Type: "transit",
		}); err != nil {
			return err
		}

		// Create default aesgcm key
		if _, err := vault.Logical().Write(path.Join(testMountPath, "keys", testKeyName), map[string]interface{}{}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to vault: %s", err)
	}
	return cleanup, retAddress, testToken, testMountPath, testKeyName, tlsConfig
}
