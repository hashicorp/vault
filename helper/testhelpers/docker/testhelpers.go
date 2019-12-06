package docker

import (
	"fmt"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/ory/dockertest"
)

func CleanupResource(t testing.T, pool *dockertest.Pool, resource *dockertest.Resource) {
	var err error
	for i := 0; i < 10; i++ {
		err = pool.Purge(resource)
		if err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}

	if strings.Contains(err.Error(), "No such container") {
		return
	}
	t.Fatalf("Failed to cleanup local container: %s", err)
}

func PrepareTestContainer(t *testing.T) (cleanup func(), retAddress, token, mountPath, keyName string, tlsConfig *api.TLSConfig) {
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
		CleanupResource(*t, pool, resource)
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
