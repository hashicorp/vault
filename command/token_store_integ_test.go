package command

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical/file"
	"github.com/hashicorp/vault/vault"

	vaulthttp "github.com/hashicorp/vault/http"
	log "github.com/mgutz/logxi/v1"
)

func TestTokenStore_Integ_TokenCreation(t *testing.T) {
	/*
		if os.Getenv(logicaltest.TestEnvVar) == "" {
			t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
			return
		}
	*/
	filePath, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(filePath)

	logger := logformat.NewVaultLogger(log.LevelTrace)

	config := map[string]string{
		"path": filePath,
	}

	underlying, err := file.NewFileBackend(config, logger)
	if err != nil {
		t.Fatal(err)
	}

	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NullLog,
		Physical:     underlying,
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	client.SetToken(cluster.RootToken)

	count := 1000
	for i := 1; i <= count; i++ {
		id := strconv.Itoa(i)
		tcr := &api.TokenCreateRequest{
			ID:          id,
			Policies:    []string{"default"},
			TTL:         "1m",
			DisplayName: "test-" + id,
		}

		secret, err := client.Auth().Token().Create(tcr)
		if err != nil {
			t.Fatal(err)
		}

		if secret.Auth.ClientToken != id {
			t.Fatalf("failed to create the token in iteration %q", id)
		}
	}
}
