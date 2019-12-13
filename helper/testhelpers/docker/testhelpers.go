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

type TestVaultContainer struct {
	Cleanup    func()
	RetAddress string
	Token      string
	MountPath  string
	KeyName    string
	TLSConfig  *api.TLSConfig
	Pool       *dockertest.Pool
}

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

func PrepareTestVaultContainer(t *testing.T) *TestVaultContainer {
	ret := &TestVaultContainer{}
	var err error
	ret.Token, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	ret.MountPath, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	ret.KeyName, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	ret.Pool, err = dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	dockerOptions := &dockertest.RunOptions{
		Repository: "vault",
		Tag:        "latest",
		Cmd: []string{"server", "-log-level=trace", "-dev", fmt.Sprintf("-dev-root-token-id=%s", ret.Token),
			"-dev-listen-address=0.0.0.0:8200"},
	}
	resource, err := ret.Pool.RunWithOptions(dockerOptions)
	if err != nil {
		t.Fatalf("Could not start local Vault docker container: %s", err)
	}

	ret.Cleanup = func() {
		CleanupResource(*t, ret.Pool, resource)
	}

	ret.RetAddress = fmt.Sprintf("http://127.0.0.1:%s", resource.GetPort("8200/tcp"))
	ret.TLSConfig = &api.TLSConfig{
		Insecure: true,
	}

	return ret
}

func (tv *TestVaultContainer) MountTransit(t *testing.T) {
	// exponential backoff-retry
	if err := tv.Pool.Retry(func() error {
		vaultConfig := api.DefaultConfig()
		vaultConfig.Address = tv.RetAddress
		if err := vaultConfig.ConfigureTLS(tv.TLSConfig); err != nil {
			return err
		}
		vault, err := api.NewClient(vaultConfig)
		if err != nil {
			return err
		}
		vault.SetToken(tv.Token)

		// Set up mount
		if err := vault.Sys().Mount(tv.MountPath, &api.MountInput{
			Type: "transit",
		}); err != nil {
			return err
		}
		// Create default aesgcm key
		if _, err := vault.Logical().Write(path.Join(tv.MountPath, "keys", tv.KeyName), map[string]interface{}{}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		tv.Cleanup()
		t.Fatalf("Could not connect to vault: %s", err)
	}
}
