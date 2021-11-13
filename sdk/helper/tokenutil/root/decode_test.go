package root

import (
	"context"
	hclog "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

func TestSimpleRootGeneration(t *testing.T) {
	core := vault.TestCore(t)
	core.SetClusterHandler(nil)
	defer core.Shutdown()

	barrierConfig := &vault.SealConfig{
		SecretShares:    3,
		SecretThreshold: 3,
		StoredShares:    1,
	}

	recoveryConfig := &vault.SealConfig{
		SecretShares:    3,
		SecretThreshold: 3,
	}

	initParams := &vault.InitParams{
		BarrierConfig:    barrierConfig,
		RecoveryConfig:   recoveryConfig,
		LegacyShamirSeal: true,
	}

	result, err := core.Initialize(context.Background(), initParams)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	unsealKeys := result.SecretShares
	for _, unsealKey := range unsealKeys {
		if _, err := core.Unseal(vault.TestKeyCopy(unsealKey)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	token, err := DecodeRootToken(result.RootToken, "", 0)

	if core.Sealed() {
		t.Fatal("should not be sealed")
	}

	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(token)

	id, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  fmt.Sprintf("root decoder-%s", id),
		Level: hclog.Trace,
	})
	logger.Info("raw token is", result.RootToken, "and client token is", token)

	secret, err := client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Data == nil || secret.Data["id"].(string) == "" {
		t.Fatalf("failed to perform lookup self through agent")
	}
}
