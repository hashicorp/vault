package identity

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	upAuth "github.com/hashicorp/vault/api/auth/userpass"
	"github.com/hashicorp/vault/helper/testhelpers"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/totp"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func doTwoPhaseLogin(t *testing.T, client *api.Client, totpCodePath, methodID, username string) {
	t.Helper()
	totpPasscode := testhelpers.GetTOTPCodeFromEngine(t, client, totpCodePath)

	upMethod, err := upAuth.NewUserpassAuth(username, &upAuth.Password{FromString: "testpassword"})

	mfaSecret, err := client.Auth().MFALogin(context.Background(), upMethod)
	if err != nil {
		t.Fatalf("failed to login with userpass auth method: %v", err)
	}

	secret, err := client.Auth().MFAValidate(
		context.Background(),
		mfaSecret,
		map[string]interface{}{
			methodID: []string{totpPasscode},
		},
	)
	if err != nil {
		t.Fatalf("MFA validation failed: %v", err)
	}

	if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
		t.Fatalf("MFA validation failed to return a ClientToken in secret: %v", secret)
	}
}

func TestLoginMfaGenerateTOTPTestAuditIncluded(t *testing.T) {
	var noop *vault.NoopAudit

	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"totp": totp.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"noop": func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
				noop = &vault.NoopAudit{
					Config: config,
				}
				return noop, nil
			},
		},
	},
		&vault.TestClusterOptions{
			HandlerFunc: vaulthttp.Handler,
		})

	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	// Enable the audit backend
	if err := client.Sys().EnableAuditWithOptions("noop", &api.EnableAuditOptions{Type: "noop"}); err != nil {
		t.Fatal(err)
	}

	testhelpers.SetupTOTPMount(t, client)
	mountAccessor := testhelpers.SetupUserpassMountAccessor(t, client)

	// Creating two users in the userpass auth mount
	userClient1, entityID1, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "entity1", "testuser1")
	userClient2, entityID2, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "entity2", "testuser2")

	totpConfig := map[string]interface{}{
		"issuer":                  "yCorp",
		"period":                  10,
		"algorithm":               "SHA512",
		"digits":                  6,
		"skew":                    0,
		"key_size":                20,
		"qr_size":                 200,
		"max_validation_attempts": 5,
	}

	methodID := testhelpers.SetupTOTPMethod(t, client, totpConfig)

	// registering EntityIDs in the TOTP secret Engine for MethodID
	enginePath1 := testhelpers.RegisterEntityInTOTPEngine(t, client, entityID1, methodID)
	enginePath2 := testhelpers.RegisterEntityInTOTPEngine(t, client, entityID2, methodID)

	// Configure a default login enforcement
	enforcementConfig := map[string]interface{}{
		"auth_method_types": []string{"userpass"},
		"name":              "randomName",
		"mfa_method_ids":    []string{methodID},
	}

	testhelpers.SetupMFALoginEnforcement(t, client, enforcementConfig)

	// MFA single-phase login
	totpPasscode1 := testhelpers.GetTOTPCodeFromEngine(t, userClient1, enginePath1)

	userClient1.AddHeader("X-Vault-MFA", fmt.Sprintf("%s:%s", methodID, totpPasscode1))
	secret, err := userClient1.Logical().WriteWithContext(context.Background(), "auth/userpass/login/testuser1", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("MFA failed: %v", err)
	}

	userpassToken := secret.Auth.ClientToken

	userClient1.SetToken(client.Token())
	secret, err = userClient1.Logical().WriteWithContext(context.Background(), "auth/token/lookup", map[string]interface{}{
		"token": userpassToken,
	})
	if err != nil {
		t.Fatalf("failed to lookup userpass authenticated token: %v", err)
	}

	entityIDCheck := secret.Data["entity_id"].(string)
	if entityIDCheck != entityID1 {
		t.Fatalf("different entityID assigned")
	}

	// Two-phase login
	headers := userClient1.Headers()
	headers.Del("X-Vault-MFA")
	userClient1.SetHeaders(headers)
	secret, err = userClient1.Logical().WriteWithContext(context.Background(), "auth/userpass/login/testuser1", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("MFA failed: %v", err)
	}

	if len(secret.Warnings) == 0 || !strings.Contains(strings.Join(secret.Warnings, ""), "A login request was issued that is subject to MFA validation") {
		t.Fatalf("first phase of login did not have a warning")
	}

	if secret.Auth == nil || secret.Auth.MFARequirement == nil {
		t.Fatalf("two phase login returned nil MFARequirement")
	}
	if secret.Auth.MFARequirement.MFARequestID == "" {
		t.Fatalf("MFARequirement contains empty MFARequestID")
	}
	if secret.Auth.MFARequirement.MFAConstraints == nil || len(secret.Auth.MFARequirement.MFAConstraints) == 0 {
		t.Fatalf("MFAConstraints is nil or empty")
	}
	mfaConstraints, ok := secret.Auth.MFARequirement.MFAConstraints["randomName"]
	if !ok {
		t.Fatalf("failed to find the mfaConstrains")
	}
	if mfaConstraints.Any == nil || len(mfaConstraints.Any) == 0 {
		t.Fatalf("expected to see the methodID is enforced in MFAConstaint.Any")
	}
	for _, mfaAny := range mfaConstraints.Any {
		if mfaAny.ID != methodID || mfaAny.Type != "totp" || !mfaAny.UsesPasscode {
			t.Fatalf("Invalid mfa constraints")
		}
	}

	// validation
	// waiting for 5 seconds so that a fresh code could be generated
	time.Sleep(5 * time.Second)
	totpPasscode1 = testhelpers.GetTOTPCodeFromEngine(t, client, enginePath1)

	secret, err = userClient1.Logical().WriteWithContext(context.Background(), "sys/mfa/validate", map[string]interface{}{
		"mfa_request_id": secret.Auth.MFARequirement.MFARequestID,
		"mfa_payload": map[string][]string{
			methodID: {totpPasscode1},
		},
	})
	if err != nil {
		t.Fatalf("MFA failed: %v", err)
	}

	if secret.Auth == nil || secret.Auth.ClientToken == "" {
		t.Fatalf("successful mfa validation did not return a client token")
	}

	if noop.Req == nil {
		t.Fatalf("no request was logged in audit log")
	}
	var found bool
	for _, req := range noop.Req {
		if req.Path == "sys/mfa/validate" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("mfa/validate was not logged in audit log")
	}

	// check for login request expiration
	secret, err = userClient1.Logical().WriteWithContext(context.Background(), "auth/userpass/login/testuser1", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("MFA failed: %v", err)
	}

	if secret.Auth == nil || secret.Auth.MFARequirement == nil {
		t.Fatalf("two phase login returned nil MFARequirement")
	}

	_, err = userClient1.Logical().WriteWithContext(context.Background(), "sys/mfa/validate", map[string]interface{}{
		"mfa_request_id": secret.Auth.MFARequirement.MFARequestID,
		"mfa_payload": map[string][]string{
			methodID: {totpPasscode1},
		},
	})
	if err == nil {
		t.Fatalf("MFA succeeded with an already used passcode")
	}
	if !strings.Contains(err.Error(), "code already used") {
		t.Fatalf("expected error message to mention code already used")
	}

	// check for reaching max failed validation requests
	secret, err = userClient1.Logical().WriteWithContext(context.Background(), "auth/userpass/login/testuser1", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("MFA failed: %v", err)
	}

	var maxErr error
	maxAttempts := 6
	i := 0
	for i = 0; i < maxAttempts; i++ {
		_, maxErr = userClient1.Logical().WriteWithContext(context.Background(), "sys/mfa/validate", map[string]interface{}{
			"mfa_request_id": secret.Auth.MFARequirement.MFARequestID,
			"mfa_payload": map[string][]string{
				methodID: {fmt.Sprintf("%d", i)},
			},
		})
		if maxErr == nil {
			t.Fatalf("MFA succeeded with an invalid passcode")
		}
	}
	if !strings.Contains(maxErr.Error(), "maximum TOTP validation attempts") {
		t.Fatalf("unexpected error message when exceeding max failed validation attempts: %s", maxErr.Error())
	}

	// let's make sure the configID is not blocked for other users
	doTwoPhaseLogin(t, userClient2, enginePath2, methodID, "testuser2")

	// let's see if user1 is able to login after 5 seconds
	time.Sleep(5 * time.Second)
	doTwoPhaseLogin(t, userClient1, enginePath1, methodID, "testuser1")

	// Destroy the secret so that the token can self generate
	_, err = client.Logical().WriteWithContext(context.Background(), fmt.Sprintf("identity/mfa/method/totp/admin-destroy"), map[string]interface{}{
		"entity_id": entityID1,
		"method_id": methodID,
	})
	if err != nil {
		t.Fatalf("failed to destroy the MFA secret: %s", err)
	}
}
