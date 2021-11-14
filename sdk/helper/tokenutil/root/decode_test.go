package root

import (
	"context"
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

	// Check if the seal configuration is valid
	if err := barrierConfig.Validate(); err != nil {
		t.Fatal("invalid seal configuration", "error", err)
	}

	initResult, err := core.Initialize(context.Background(), initParams)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	unsealKeys := initResult.SecretShares

	err = core.GenerateRootInit("", "", vault.GenerateStandardRootTokenStrategy)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var encodedRoot string
	for _, unsealKey := range unsealKeys {
		updateRes, err := core.GenerateRootUpdate(context.TODO(), unsealKey, "", vault.GenerateStandardRootTokenStrategy)
		if err != nil {
			t.Fatalf("unseal err: %s", err)
		} else if updateRes.Progress == recoveryConfig.SecretShares {
			encodedRoot = updateRes.EncodedToken
		}
	}

	token, err := DecodeRootToken(encodedRoot, "", 0)
	if err != nil {
		t.Fatalf("unseal err: %s", err)
	}

	if token != initResult.RootToken {
		t.Fatal("Decoded root token is different than the original token (", token, initResult.RootToken)
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

	secret, err := client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || secret.Data == nil || secret.Data["id"].(string) == "" {
		t.Fatalf("failed to perform lookup self through agent")
	}
}
