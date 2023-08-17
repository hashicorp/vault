// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	upAuth "github.com/hashicorp/vault/api/auth/userpass"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
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
	t.Setenv("VAULT_AUDIT_DISABLE_EVENTLOGGER", "true")

	noop := corehelpers.TestNoopAudit(t, nil)

	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"totp": totp.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"noop": func(ctx context.Context, config *audit.BackendConfig, _ bool, _ audit.HeaderFormatter) (audit.Backend, error) {
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

	// Get a random set of chars to seed our entity and alias names
	userseed := base64.StdEncoding.EncodeToString([]byte("couple of test users"))
	entity1 := userseed[0:3]
	testuser1 := userseed[3:6]
	entity2 := userseed[6:9]
	testuser2 := userseed[9:12]

	// Creating two users in the userpass auth mount
	userClient1, entityID1, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, entity1, testuser1)
	userClient2, entityID2, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, entity2, testuser2)
	waitPeriod := 5
	totpConfig := map[string]interface{}{
		"issuer":                  "yCorp",
		"period":                  waitPeriod,
		"algorithm":               "SHA1",
		"digits":                  6,
		"skew":                    1,
		"key_size":                10,
		"qr_size":                 100,
		"max_validation_attempts": 3,
		"method_name":             "foo",
	}

	methodID := testhelpers.SetupTOTPMethod(t, client, totpConfig)

	// registering EntityIDs in the TOTP secret Engine for MethodID
	enginePath1 := testhelpers.RegisterEntityInTOTPEngine(t, client, entityID1, methodID)
	enginePath2 := testhelpers.RegisterEntityInTOTPEngine(t, client, entityID2, methodID)

	// Configure a default login enforcement
	enforcementConfig := map[string]interface{}{
		"auth_method_types": []string{"userpass"},
		"name":              methodID[0:4],
		"mfa_method_ids":    []string{methodID},
	}

	testhelpers.SetupMFALoginEnforcement(t, client, enforcementConfig)

	userpassPath := fmt.Sprintf("auth/userpass/login/%s", testuser1)

	// MFA single-phase login
	verifyLoginRequest := func(secret *api.Secret) {
		userpassToken := secret.Auth.ClientToken
		userClient1.SetToken(client.Token())
		secret, err := userClient1.Logical().WriteWithContext(context.Background(), "auth/token/lookup", map[string]interface{}{
			"token": userpassToken,
		})
		if err != nil {
			t.Fatalf("failed to lookup userpass authenticated token: %v", err)
		}

		entityIDCheck := secret.Data["entity_id"].(string)
		if entityIDCheck != entityID1 {
			t.Fatalf("different entityID assigned")
		}
	}

	// helper function to clear the MFA request header
	clearMFARequestHeaders := func(c *api.Client) {
		headers := c.Headers()
		headers.Del("X-Vault-MFA")
		c.SetHeaders(headers)
	}

	var secret *api.Secret
	var err error
	var methodIdentifier string

	singlePhaseLoginFunc := func() error {
		totpPasscode := testhelpers.GetTOTPCodeFromEngine(t, client, enginePath1)
		userClient1.AddHeader("X-Vault-MFA", fmt.Sprintf("%s:%s", methodIdentifier, totpPasscode))
		defer clearMFARequestHeaders(userClient1)
		secret, err = userClient1.Logical().WriteWithContext(context.Background(), userpassPath, map[string]interface{}{
			"password": "testpassword",
		})
		if err != nil {
			return fmt.Errorf("MFA failed for identifier %s: %v", methodIdentifier, err)
		}
		return nil
	}

	// single phase login for both method name and method ID
	methodIdentifier = totpConfig["method_name"].(string)
	testhelpers.RetryUntilAtCadence(t, 20*time.Second, 100*time.Millisecond, singlePhaseLoginFunc)
	verifyLoginRequest(secret)

	methodIdentifier = methodID
	// Need to wait a bit longer to avoid hitting maximum allowed consecutive
	// failed TOTP validation
	testhelpers.RetryUntilAtCadence(t, 20*time.Second, time.Duration(waitPeriod)*time.Second, singlePhaseLoginFunc)
	verifyLoginRequest(secret)

	// Two-phase login
	secret, err = userClient1.Logical().WriteWithContext(context.Background(), userpassPath, map[string]interface{}{
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
	mfaConstraints, ok := secret.Auth.MFARequirement.MFAConstraints[methodID[0:4]]
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
	var mfaReqID string
	var totpPasscode1 string
	mfaValidateFunc := func() error {
		totpPasscode1 = testhelpers.GetTOTPCodeFromEngine(t, client, enginePath1)
		secret, err = userClient1.Logical().WriteWithContext(context.Background(), "sys/mfa/validate", map[string]interface{}{
			"mfa_request_id": mfaReqID,
			"mfa_payload": map[string][]string{
				methodIdentifier: {totpPasscode1},
			},
		})
		if err != nil {
			return fmt.Errorf("MFA failed: %v", err)
		}
		if secret.Auth == nil || secret.Auth.ClientToken == "" {
			t.Fatalf("successful mfa validation did not return a client token")
		}

		return nil
	}

	methodIdentifier = methodID
	mfaReqID = secret.Auth.MFARequirement.MFARequestID
	testhelpers.RetryUntilAtCadence(t, 20*time.Second, time.Duration(waitPeriod)*time.Second, mfaValidateFunc)

	// two phase login with method name
	secret, err = userClient1.Logical().WriteWithContext(context.Background(), userpassPath, map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("MFA failed: %v", err)
	}

	methodIdentifier = totpConfig["method_name"].(string)
	mfaReqID = secret.Auth.MFARequirement.MFARequestID
	testhelpers.RetryUntilAtCadence(t, 20*time.Second, time.Duration(waitPeriod)*time.Second, mfaValidateFunc)

	// checking audit log
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
	secret, err = userClient1.Logical().WriteWithContext(context.Background(), userpassPath, map[string]interface{}{
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
		t.Fatalf("got: %+v, expected: code already used", err.Error())
	}

	// check for reaching max failed validation requests
	secret, err = userClient1.Logical().WriteWithContext(context.Background(), userpassPath, map[string]interface{}{
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
	doTwoPhaseLogin(t, userClient2, enginePath2, methodID, testuser2)

	// let's see if user1 is able to login after 5 seconds
	time.Sleep(5 * time.Second)
	doTwoPhaseLogin(t, userClient1, enginePath1, methodID, testuser1)

	// Destroy the secret so that the token can self generate
	_, err = client.Logical().WriteWithContext(context.Background(), fmt.Sprintf("identity/mfa/method/totp/admin-destroy"), map[string]interface{}{
		"entity_id": entityID1,
		"method_id": methodID,
	})
	if err != nil {
		t.Fatalf("failed to destroy the MFA secret: %s", err)
	}
}
