package identity

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	upAuth "github.com/hashicorp/vault/api/auth/userpass"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/totp"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func createEntityAndAlias(client *api.Client, mountAccessor, entityName, aliasName string, t *testing.T) (*api.Client, string) {
	_, err := client.Logical().WriteWithContext(context.Background(), fmt.Sprintf("auth/userpass/users/%s", aliasName), map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("failed to configure userpass backend: %v", err)
	}

	userClient, err := client.Clone()
	if err != nil {
		t.Fatalf("failed to clone the client:%v", err)
	}
	userClient.SetToken(client.Token())

	resp, err := client.Logical().WriteWithContext(context.Background(), "identity/entity", map[string]interface{}{
		"name": entityName,
	})
	if err != nil {
		t.Fatalf("failed to create an entity:%v", err)
	}
	entityID := resp.Data["id"].(string)

	_, err = client.Logical().WriteWithContext(context.Background(), "identity/entity-alias", map[string]interface{}{
		"name":           aliasName,
		"canonical_id":   entityID,
		"mount_accessor": mountAccessor,
	})
	if err != nil {
		t.Fatalf("failed to create an entity alias:%v", err)
	}
	return userClient, entityID
}

func registerEntityInTOTPEngine(client *api.Client, entityID, methodID string, t *testing.T) string {
	totpGenName := fmt.Sprintf("%s-%s", entityID, methodID)
	secret, err := client.Logical().WriteWithContext(context.Background(), fmt.Sprintf("identity/mfa/method/totp/admin-generate"), map[string]interface{}{
		"entity_id": entityID,
		"method_id": methodID,
	})
	if err != nil {
		t.Fatalf("failed to generate a TOTP secret on an entity: %v", err)
	}
	totpURL := secret.Data["url"].(string)

	_, err = client.Logical().WriteWithContext(context.Background(), fmt.Sprintf("totp/keys/%s", totpGenName), map[string]interface{}{
		"url": totpURL,
	})
	if err != nil {
		t.Fatalf("failed to register a TOTP URL: %v", err)
	}
	return totpGenName
}

func doTwoPhaseLogin(client *api.Client, totpCodePath, methodID, username string, t *testing.T) {
	totpResp, err := client.Logical().ReadWithContext(context.Background(), totpCodePath)
	if err != nil {
		t.Fatalf("failed to create totp passcode: %v", err)
	}
	totpPasscode := totpResp.Data["code"].(string)

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
	err := client.Sys().EnableAuditWithOptions("noop", &api.EnableAuditOptions{Type: "noop"})
	if err != nil {
		t.Fatal(err)
	}

	// Mount the TOTP backend
	mountInfo := &api.MountInput{
		Type: "totp",
	}
	err = client.Sys().Mount("totp", mountInfo)
	if err != nil {
		t.Fatalf("failed to mount totp backend: %v", err)
	}

	// Enable Userpass authentication
	err = client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass auth: %v", err)
	}

	auths, err := client.Sys().ListAuthWithContext(context.Background())
	if err != nil {
		t.Fatalf("bb")
	}
	var mountAccessor string
	if auths != nil && auths["userpass/"] != nil {
		mountAccessor = auths["userpass/"].Accessor
	}

	// Creating two users in the userpass auth mount
	userClient1, entityID1 := createEntityAndAlias(client, mountAccessor, "entity1", "testuser1", t)
	userClient2, entityID2 := createEntityAndAlias(client, mountAccessor, "entity2", "testuser2", t)

	// configure TOTP secret engine
	var methodID string
	// login MFA
	{
		// create a config
		resp1, err := client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
			"issuer":                  "yCorp",
			"period":                  5,
			"algorithm":               "SHA1",
			"digits":                  6,
			"skew":                    1,
			"key_size":                10,
			"qr_size":                 100,
			"max_validation_attempts": 3,
		})

		if err != nil || (resp1 == nil) {
			t.Fatalf("bad: resp: %#v\n err: %v", resp1, err)
		}

		methodID = resp1.Data["method_id"].(string)
		if methodID == "" {
			t.Fatalf("method ID is empty")
		}

		// creating MFAEnforcementConfig
		_, err = client.Logical().WriteWithContext(context.Background(), "identity/mfa/login-enforcement/randomName", map[string]interface{}{
			"auth_method_types": []string{"userpass"},
			"name":              "randomName",
			"mfa_method_ids":    []string{methodID},
		})
		if err != nil {
			t.Fatalf("failed to configure MFAEnforcementConfig: %v", err)
		}
	}

	// registering EntityIDs in the TOTP secret Engine for MethodID
	totpEngineConfigName1 := registerEntityInTOTPEngine(client, entityID1, methodID, t)
	totpEngineConfigName2 := registerEntityInTOTPEngine(client, entityID2, methodID, t)

	// MFA single-phase login
	totpCodePath1 := fmt.Sprintf("totp/code/%s", totpEngineConfigName1)
	secret, err := client.Logical().ReadWithContext(context.Background(), totpCodePath1)
	if err != nil {
		t.Fatalf("failed to create totp passcode: %v", err)
	}
	totpPasscode1 := secret.Data["code"].(string)

	userClient1.AddHeader("X-Vault-MFA", fmt.Sprintf("%s:%s", methodID, totpPasscode1))
	secret, err = userClient1.Logical().WriteWithContext(context.Background(), "auth/userpass/login/testuser1", map[string]interface{}{
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
	// getting a fresh totp passcode for the validation step
	totpResp, err := client.Logical().ReadWithContext(context.Background(), totpCodePath1)
	if err != nil {
		t.Fatalf("failed to create totp passcode: %v", err)
	}
	totpPasscode1 = totpResp.Data["code"].(string)

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
	for i := 0; i < 4; i++ {
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
	if !strings.Contains(maxErr.Error(), "maximum TOTP validation attempts 4 exceeded the allowed attempts 3") {
		t.Fatalf("unexpected error message when exceeding max failed validation attempts")
	}

	// let's make sure the configID is not blocked for other users
	totpCodePath2 := fmt.Sprintf("totp/code/%s", totpEngineConfigName2)
	doTwoPhaseLogin(userClient2, totpCodePath2, methodID, "testuser2", t)

	// let's see if user1 is able to login after 5 seconds
	time.Sleep(5 * time.Second)
	// getting a fresh totp passcode for the validation step
	doTwoPhaseLogin(userClient1, totpCodePath1, methodID, "testuser1", t)

	// Destroy the secret so that the token can self generate
	_, err = client.Logical().WriteWithContext(context.Background(), fmt.Sprintf("identity/mfa/method/totp/admin-destroy"), map[string]interface{}{
		"entity_id": entityID1,
		"method_id": methodID,
	})
	if err != nil {
		t.Fatalf("failed to destroy the MFA secret: %s", err)
	}
}
