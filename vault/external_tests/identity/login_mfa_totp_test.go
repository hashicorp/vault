package identity

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/totp"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

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

	// Creating a user in the userpass auth mount
	_, err = client.Logical().Write("auth/userpass/users/testuser", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("failed to configure userpass backend: %v", err)
	}

	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatalf("bb")
	}
	var mountAccessor string
	if auths != nil && auths["userpass/"] != nil {
		mountAccessor = auths["userpass/"].Accessor
	}

	userClient, err := client.Clone()
	if err != nil {
		t.Fatalf("failed to clone the client")
	}
	userClient.SetToken(client.Token())

	var entityID string
	var groupID string
	{
		resp, err := userClient.Logical().Write("identity/entity", map[string]interface{}{
			"name": "test-entity",
			"metadata": map[string]string{
				"email":        "test@hashicorp.com",
				"phone_number": "123-456-7890",
			},
		})
		if err != nil {
			t.Fatalf("failed to create an entity")
		}
		entityID = resp.Data["id"].(string)

		// Create a group
		resp, err = client.Logical().Write("identity/group", map[string]interface{}{
			"name":              "engineering",
			"member_entity_ids": []string{entityID},
		})
		if err != nil {
			t.Fatalf("failed to create an identity group")
		}
		groupID = resp.Data["id"].(string)

		_, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
			"name":           "testuser",
			"canonical_id":   entityID,
			"mount_accessor": mountAccessor,
		})
		if err != nil {
			t.Fatalf("failed to create an entity alias")
		}

	}

	// configure TOTP secret engine
	var totpPasscode string
	var methodID string
	var userpassToken string
	// login MFA
	{
		// create a config
		resp1, err := client.Logical().Write("identity/mfa/method/totp", map[string]interface{}{
			"issuer":    "yCorp",
			"period":    5,
			"algorithm": "SHA1",
			"digits":    6,
			"skew":      1,
			"key_size":  10,
			"qr_size":   100,
		})

		if err != nil || (resp1 == nil) {
			t.Fatalf("bad: resp: %#v\n err: %v", resp1, err)
		}

		methodID = resp1.Data["method_id"].(string)
		if methodID == "" {
			t.Fatalf("method ID is empty")
		}

		secret, err := client.Logical().Write(fmt.Sprintf("identity/mfa/method/totp/admin-generate"), map[string]interface{}{
			"entity_id": entityID,
			"method_id": methodID,
		})
		if err != nil {
			t.Fatalf("failed to generate a TOTP secret on an entity: %v", err)
		}
		totpURL := secret.Data["url"].(string)

		_, err = client.Logical().Write("totp/keys/loginMFA", map[string]interface{}{
			"url": totpURL,
		})
		if err != nil {
			t.Fatalf("failed to register a TOTP URL: %v", err)
		}

		secret, err = client.Logical().Read("totp/code/loginMFA")
		if err != nil {
			t.Fatalf("failed to create totp passcode: %v", err)
		}
		totpPasscode = secret.Data["code"].(string)

		// creating MFAEnforcementConfig
		_, err = client.Logical().Write("identity/mfa/login-enforcement/randomName", map[string]interface{}{
			"auth_method_accessors": []string{mountAccessor},
			"auth_method_types":     []string{"userpass"},
			"identity_group_ids":    []string{groupID},
			"identity_entity_ids":   []string{entityID},
			"name":                  "randomName",
			"mfa_method_ids":        []string{methodID},
		})
		if err != nil {
			t.Fatalf("failed to configure MFAEnforcementConfig: %v", err)
		}

		// MFA single-phase login
		userClient.AddHeader("X-Vault-MFA", fmt.Sprintf("%s:%s", methodID, totpPasscode))
		secret, err = userClient.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
			"password": "testpassword",
		})
		if err != nil {
			t.Fatalf("MFA failed: %v", err)
		}

		userpassToken = secret.Auth.ClientToken

		userClient.SetToken(client.Token())
		secret, err = userClient.Logical().Write("auth/token/lookup", map[string]interface{}{
			"token": userpassToken,
		})
		if err != nil {
			t.Fatalf("failed to lookup userpass authenticated token: %v", err)
		}

		entityIDCheck := secret.Data["entity_id"].(string)
		if entityIDCheck != entityID {
			t.Fatalf("different entityID assigned")
		}

		// Two-phase login
		user2Client, err := client.Clone()
		if err != nil {
			t.Fatalf("failed to clone the client")
		}
		headers := user2Client.Headers()
		headers.Del("X-Vault-MFA")
		user2Client.SetHeaders(headers)
		secret, err = user2Client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
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
			t.Fatalf("")
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
		totpResp, err := client.Logical().Read("totp/code/loginMFA")
		if err != nil {
			t.Fatalf("failed to create totp passcode: %v", err)
		}
		totpPasscode = totpResp.Data["code"].(string)

		secret, err = user2Client.Logical().Write("sys/mfa/validate", map[string]interface{}{
			"mfa_request_id": secret.Auth.MFARequirement.MFARequestID,
			"mfa_payload": map[string][]string{
				methodID: {totpPasscode},
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
		secret, err = user2Client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
			"password": "testpassword",
		})
		if err != nil {
			t.Fatalf("MFA failed: %v", err)
		}

		if secret.Auth == nil || secret.Auth.MFARequirement == nil {
			t.Fatalf("two phase login returned nil MFARequirement")
		}

		_, err = user2Client.Logical().Write("sys/mfa/validate", map[string]interface{}{
			"mfa_request_id": secret.Auth.MFARequirement.MFARequestID,
			"mfa_payload": map[string][]string{
				methodID: {totpPasscode},
			},
		})
		if err == nil {
			t.Fatalf("MFA succeeded with an already used passcode")
		}
		if !strings.Contains(err.Error(), "code already used") {
			t.Fatalf("expected error message to mention code already used")
		}

		// Destroy the secret so that the token can self generate
		_, err = userClient.Logical().Write(fmt.Sprintf("identity/mfa/method/totp/admin-destroy"), map[string]interface{}{
			"entity_id": entityID,
			"method_id": methodID,
		})
		if err != nil {
			t.Fatalf("failed to destroy the MFA secret: %s", err)
		}
	}
}
