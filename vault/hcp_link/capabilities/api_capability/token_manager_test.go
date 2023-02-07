package api_capability

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	scada "github.com/hashicorp/hcp-scada-provider"
	sdkResource "github.com/hashicorp/hcp-sdk-go/resource"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/hcp_link/internal"
)

func getHCPConfig(t *testing.T, clientID, clientSecret string) *configutil.HCPLinkConfig {
	resourceIDRaw, ok := os.LookupEnv("HCP_RESOURCE_ID")
	if !ok {
		t.Skip("failed to find the HCP resource ID")
	}
	res, err := sdkResource.FromString(resourceIDRaw)
	if err != nil {
		t.Fatalf("failed to parse the resource ID, %v", err.Error())
	}
	return &configutil.HCPLinkConfig{
		ResourceIDRaw: resourceIDRaw,
		Resource:      &res,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
	}
}

func getTestScadaConfig(t *testing.T) *scada.Config {
	clientID, ok := os.LookupEnv("HCP_CLIENT_ID")
	if !ok {
		t.Skip("HCP client ID not found in env")
	}
	clientSecret, ok := os.LookupEnv("HCP_CLIENT_SECRET")
	if !ok {
		t.Skip("HCP client secret not found in env")
	}

	if _, ok := os.LookupEnv("HCP_API_ADDRESS"); !ok {
		t.Skip("failed to find HCP_API_ADDRESS in the environment")
	}
	if _, ok := os.LookupEnv("HCP_SCADA_ADDRESS"); !ok {
		t.Skip("failed to find HCP_SCADA_ADDRESS in the environment")
	}
	if _, ok := os.LookupEnv("HCP_AUTH_URL"); !ok {
		t.Skip("failed to find HCP_AUTH_URL in the environment")
	}

	hcpConfig := getHCPConfig(t, clientID, clientSecret)

	scadaConfig, err := internal.NewScadaConfig(hcpConfig, hclog.New(nil))
	if err != nil {
		t.Fatalf("failed to initialize Scada config")
	}
	return scadaConfig
}

func TestCreateTokenCoreSealedUnSealed(t *testing.T) {
	t.Parallel()
	core := vault.TestCore(t)
	logger := hclog.New(nil)

	scadaConfig := getTestScadaConfig(t)

	tm, err := NewHCPLinkTokenManager(scadaConfig, core, logger)
	if err != nil {
		t.Fatalf("failed to instantiate token manager")
	}

	if tm.latestToken != "" {
		t.Fatalf("unexpected latest token")
	}

	ctx := context.Background()
	tm.createToken(ctx)

	if tm.latestToken != "" {
		t.Fatalf("unexpected latest token while core is sealed")
	}

	// unsealing core
	vault.TestInitUnsealCore(t, core)

	tm.createToken(ctx)
	if tm.latestToken == "" {
		t.Fatalf("latestToken should not be empty")
	}
	latestTokenOld := tm.latestToken

	// running update token again should not change the token
	tm.createToken(ctx)

	if tm.latestToken == latestTokenOld {
		t.Fatalf("latestToken should have been refreshed")
	}
}

func TestShutdownTokenManagerForgetsTokenPolicy(t *testing.T) {
	t.Parallel()
	core := vault.TestCore(t)
	logger := hclog.New(nil)

	scadaConfig := getTestScadaConfig(t)

	tm, err := NewHCPLinkTokenManager(scadaConfig, core, logger)
	if err != nil {
		t.Fatalf("failed to instantiate token manager")
	}

	// unsealing core
	vault.TestInitUnsealCore(t, core)

	tm.HandleTokenPolicy(context.Background(), true)
	if tm.GetLatestToken() == "" {
		t.Fatalf("token manager did not update both token and policy")
	}

	tm.Shutdown()

	if tm.GetLatestToken() != "" || tm.policy != "" {
		t.Fatalf("shutting down TM did not forget both token and policy")
	}
}

func TestSealVaultTokenManagerForgetsTokenPolicy(t *testing.T) {
	t.Parallel()
	core := vault.TestCore(t)
	logger := hclog.New(nil)

	scadaConfig := getTestScadaConfig(t)

	tm, err := NewHCPLinkTokenManager(scadaConfig, core, logger)
	if err != nil {
		t.Fatalf("failed to instantiate token manager")
	}

	// unsealing core
	vault.TestInitUnsealCore(t, core)

	ctx := context.Background()
	tm.HandleTokenPolicy(ctx, true)

	if tm.GetLatestToken() == "" {
		t.Fatalf("token manager did not update both token and policy")
	}

	// seal core
	err = vault.TestCoreSeal(core)
	if err != nil {
		t.Fatalf("failed to seal core")
	}

	tm.HandleTokenPolicy(ctx, true)
	if tm.GetLatestToken() != "" || tm.GetPolicy() != "" {
		t.Fatalf("vault is seal, TM did not forget both token and policy")
	}
}

func TestCreateTokenWithTTL(t *testing.T) {
	t.Parallel()
	core := vault.TestCore(t)
	logger := hclog.New(nil)

	scadaConfig := getTestScadaConfig(t)

	tm, err := NewHCPLinkTokenManager(scadaConfig, core, logger)
	if err != nil {
		t.Fatalf("failed to instantiate token manager")
	}

	// unsealing core
	vault.TestInitUnsealCore(t, core)

	ttl := 7 * time.Second
	err = tm.SetTokenTTL(ttl)
	if err != nil {
		t.Fatalf("failed to set token TTL")
	}

	ctx := context.Background()
	tm.createToken(ctx)

	latestToken := tm.GetLatestToken()
	if latestToken == "" {
		t.Fatalf("latestToken should not be empty")
	}

	te, err := core.LookupToken(namespace.RootContext(nil), latestToken)
	if err != nil {
		t.Fatalf("failed to look up token")
	}
	if te.TTL != ttl {
		t.Fatalf("ttl is not as expected")
	}

	// sleep until the token is expired
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		te, err = core.LookupToken(namespace.RootContext(nil), latestToken)
		if err == nil && te == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil && te != nil {
		t.Fatalf("token did not expire as expected")
	}
}

func TestHandleTokenPolicyWithTokenTTL(t *testing.T) {
	t.Parallel()
	core := vault.TestCore(t)
	logger := hclog.New(nil)

	scadaConfig := getTestScadaConfig(t)

	tm, err := NewHCPLinkTokenManager(scadaConfig, core, logger)
	if err != nil {
		t.Fatalf("failed to instantiate token manager")
	}

	// unsealing core
	vault.TestInitUnsealCore(t, core)

	ttl := 7 * time.Second
	err = tm.SetTokenTTL(ttl)
	if err != nil {
		t.Fatalf("failed to set token TTL")
	}

	ctx := context.Background()
	tm.createToken(ctx)

	latestToken := tm.GetLatestToken()
	if latestToken == "" {
		t.Fatalf("latestToken should not be empty")
	}

	// waiting until TTL - 5 seconds is past
	time.Sleep(ttl - tokenExpiryOffset + 1)

	newToken := tm.HandleTokenPolicy(ctx, true)
	if newToken == "" || newToken == latestToken {
		t.Fatalf("token did not refreshed after expiry")
	}
}

func TestHandleTokenPolicySealedUnsealed(t *testing.T) {
	t.Parallel()
	core := vault.TestCore(t)
	logger := hclog.New(nil)

	scadaConfig := getTestScadaConfig(t)

	tm, err := NewHCPLinkTokenManager(scadaConfig, core, logger)
	if err != nil {
		t.Fatalf("failed to instantiate token manager")
	}

	tm.latestToken = "non empty"
	tm.policy = "non empty"

	ctx := context.Background()
	tm.HandleTokenPolicy(ctx, true)
	if tm.GetLatestToken() != "" || tm.GetPolicy() != "" {
		t.Fatalf("on sealed Vault, token and policy was not deleted")
	}

	tm.latestToken = "non empty"
	tm.policy = "non empty"

	// unsealing core
	vault.TestInitUnsealCore(t, core)

	tm.HandleTokenPolicy(ctx, true)

	if tm.GetLatestToken() == "non empty" {
		t.Fatalf("latestToken and policy should have been updated")
	}
}

func TestForgetTokenPolicySealedUnsealed(t *testing.T) {
	t.Parallel()
	core := vault.TestCore(t)
	logger := hclog.New(nil)

	scadaConfig := getTestScadaConfig(t)

	tm, err := NewHCPLinkTokenManager(scadaConfig, core, logger)
	if err != nil {
		t.Fatalf("failed to instantiate token manager")
	}

	// unsealing core
	vault.TestInitUnsealCore(t, core)

	ctx := context.Background()
	tm.HandleTokenPolicy(ctx, true)

	if tm.GetLatestToken() == "" {
		t.Fatalf("on sealed Vault, token and policy was not refreshed")
	}

	tm.ForgetTokenPolicy()

	if tm.GetLatestToken() != "" || tm.GetPolicy() != "" {
		t.Fatalf("on sealed Vault, token and policy was not deleted")
	}
}
