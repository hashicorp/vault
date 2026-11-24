// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	upAuth "github.com/hashicorp/vault/api/auth/userpass"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoginMFA_PastLifetime verifies, end to end, that a pending MFA login
// cannot be completed once it is past its lifetime, even if background cleanup
// has not removed it yet. Older expired entries are queued ahead of the
// target, so the request ID is still in the queue when validation is attempted.
// It also checks that a new login with the updated password still succeeds.
func TestLoginMFA_PastLifetime(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{"userpass": userpass.Factory},
		LogicalBackends:    map[string]logical.Factory{"totp": totp.Factory},
	})
	client := cluster.Cores[0].Client
	core := cluster.Cores[0].Core

	testhelpers.SetupTOTPMount(t, client)
	mountAccessor := testhelpers.SetupUserpassMountAccessor(t, client)
	_, entityID, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "ent-lifetime", "user-lifetime")
	methodID := testhelpers.SetupTOTPMethod(t, client, map[string]interface{}{
		"issuer": "vault", "period": 30, "algorithm": "SHA1", "digits": 6,
		"skew": 1, "key_size": 20, "qr_size": 200, "max_validation_attempts": 5,
	})
	// RegisterEntityInTOTPEngine also creates the login enforcement for this
	// entity, so no separate enforcement is configured here. A second
	// enforcement would require the same TOTP code to be used twice.
	totpPath := testhelpers.RegisterEntityInTOTPEngine(t, client, entityID, methodID)

	// Non-root client for the user flow.
	userClient, err := client.Clone()
	require.NoError(t, err)
	userClient.ClearToken()

	oldPw, err := upAuth.NewUserpassAuth("user-lifetime", &upAuth.Password{FromString: "testpassword"})
	require.NoError(t, err)
	mfaSecret, err := userClient.Auth().MFALogin(context.Background(), oldPw)
	require.NoError(t, err)
	reqID := mfaSecret.Auth.MFARequirement.MFARequestID

	// Backdate the real pending entry past lifetime.
	entry, err := core.PopMFAResponseAuthByID(reqID)
	require.NoError(t, err)
	require.NotNil(t, entry)
	entry.TimeOfStorage = time.Now().Add(-6 * time.Minute)
	require.NoError(t, core.SaveMFAResponseAuth(entry))

	for i := 0; i < 32; i++ {
		require.NoError(t, core.SaveMFAResponseAuth(&vault.MFACachedAuthResponse{
			RequestID:     fmt.Sprintf("queued-%d", i),
			TimeOfStorage: time.Now().Add(-time.Hour - time.Duration(i)*time.Second),
		}))
	}

	_, err = client.Logical().Write("auth/userpass/users/user-lifetime/password", map[string]interface{}{
		"password": "updated-password",
	})
	require.NoError(t, err)
	_, err = userClient.Auth().MFALogin(context.Background(), oldPw)
	require.Error(t, err)

	code := testhelpers.GetTOTPCodeFromEngine(t, client, totpPath)
	secret, err := userClient.Auth().MFAValidate(context.Background(), mfaSecret,
		map[string]interface{}{methodID: []string{code}})
	require.ErrorContains(t, err, "invalid request ID")
	require.Nil(t, secret)

	newPw, err := upAuth.NewUserpassAuth("user-lifetime", &upAuth.Password{FromString: "updated-password"})
	require.NoError(t, err)
	fresh, err := userClient.Auth().MFALogin(context.Background(), newPw)
	require.NoError(t, err)
	secret, err = userClient.Auth().MFAValidate(context.Background(), fresh,
		map[string]interface{}{methodID: []string{code}})
	require.NoError(t, err)
	require.NotNil(t, secret)
	require.NotEmpty(t, secret.Auth.ClientToken)
}

// TestLoginMFA_Redemption verifies that when several clients validate the same
// live MFA request at the same time with a valid code, exactly one token is
// issued.
func TestLoginMFA_Redemption(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{"userpass": userpass.Factory},
		LogicalBackends:    map[string]logical.Factory{"totp": totp.Factory},
	})
	client := cluster.Cores[0].Client

	testhelpers.SetupTOTPMount(t, client)
	mountAccessor := testhelpers.SetupUserpassMountAccessor(t, client)
	_, entityID, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "ent-concurrent", "user-concurrent")
	methodID := testhelpers.SetupTOTPMethod(t, client, map[string]interface{}{
		"issuer": "vault", "period": 30, "algorithm": "SHA1", "digits": 6,
		"skew": 1, "key_size": 20, "qr_size": 200, "max_validation_attempts": 5,
	})
	// RegisterEntityInTOTPEngine also creates the login enforcement for this
	// entity, so no separate enforcement is configured here. A second
	// enforcement would require the same TOTP code to be used twice.
	totpPath := testhelpers.RegisterEntityInTOTPEngine(t, client, entityID, methodID)

	userClient, err := client.Clone()
	require.NoError(t, err)
	userClient.ClearToken()

	pw, err := upAuth.NewUserpassAuth("user-concurrent", &upAuth.Password{FromString: "testpassword"})
	require.NoError(t, err)
	mfaSecret, err := userClient.Auth().MFALogin(context.Background(), pw)
	require.NoError(t, err)
	code := testhelpers.GetTOTPCodeFromEngine(t, client, totpPath)

	var wg sync.WaitGroup
	var tokens atomic.Int64
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, err := userClient.Clone()
			if !assert.NoError(t, err) {
				return
			}
			c.ClearToken()
			s, err := c.Auth().MFAValidate(context.Background(), mfaSecret,
				map[string]interface{}{methodID: []string{code}})
			if err == nil && s != nil && s.Auth != nil && s.Auth.ClientToken != "" {
				tokens.Add(1)
			}
		}()
	}
	wg.Wait()
	require.EqualValues(t, 1, tokens.Load())
}

// TestLoginMFA_DroppedOnSeal verifies that pending MFA state is in-memory only.
func TestLoginMFA_DroppedOnSeal(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{"userpass": userpass.Factory},
		LogicalBackends:    map[string]logical.Factory{"totp": totp.Factory},
	})
	client := cluster.Cores[0].Client

	testhelpers.SetupTOTPMount(t, client)
	mountAccessor := testhelpers.SetupUserpassMountAccessor(t, client)
	_, entityID, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "ent-seal", "user-seal")
	methodID := testhelpers.SetupTOTPMethod(t, client, map[string]interface{}{
		"issuer": "vault", "period": 30, "algorithm": "SHA1", "digits": 6,
		"skew": 1, "key_size": 20, "qr_size": 200, "max_validation_attempts": 5,
	})
	// RegisterEntityInTOTPEngine also creates the login enforcement for this
	// entity, so no separate enforcement is configured here. A second
	// enforcement would require the same TOTP code to be used twice.
	totpPath := testhelpers.RegisterEntityInTOTPEngine(t, client, entityID, methodID)

	userClient, err := client.Clone()
	require.NoError(t, err)
	userClient.ClearToken()

	pw, err := upAuth.NewUserpassAuth("user-seal", &upAuth.Password{FromString: "testpassword"})
	require.NoError(t, err)
	mfaSecret, err := userClient.Auth().MFALogin(context.Background(), pw)
	require.NoError(t, err)

	require.NoError(t, client.Sys().Seal())
	cluster.UnsealCore(t, cluster.Cores[0])
	testhelpers.WaitForActiveNode(t, cluster)

	code := testhelpers.GetTOTPCodeFromEngine(t, client, totpPath)
	_, err = userClient.Auth().MFAValidate(context.Background(), mfaSecret,
		map[string]interface{}{methodID: []string{code}})
	require.ErrorContains(t, err, "invalid request ID")
}

// TestLoginMFA_PendingEntityLimit verifies, end to end, that an entity cannot
// hold more than 100 live pending MFA logins at once: the 101st first-factor
// login is rejected, while a different entity can still start an MFA login.
// This guards against one authenticated first-factor client building an
// unbounded backlog of pending MFA state.
func TestLoginMFA_PendingEntityLimit(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{"userpass": userpass.Factory},
		LogicalBackends:    map[string]logical.Factory{"totp": totp.Factory},
	})
	client := cluster.Cores[0].Client

	testhelpers.SetupTOTPMount(t, client)
	mountAccessor := testhelpers.SetupUserpassMountAccessor(t, client)
	_, entityID, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "ent-limit", "user-limit")
	_, otherEntityID, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "ent-other", "user-other")
	methodID := testhelpers.SetupTOTPMethod(t, client, map[string]interface{}{
		"issuer": "vault", "period": 30, "algorithm": "SHA1", "digits": 6,
		"skew": 1, "key_size": 20, "qr_size": 200, "max_validation_attempts": 5,
	})
	testhelpers.RegisterEntityInTOTPEngine(t, client, entityID, methodID)
	// The helper keys the enforcement on the method ID, so rewrite it to cover
	// both entities rather than registering the second one separately.
	_, err := client.Logical().Write(fmt.Sprintf("identity/mfa/login-enforcement/%s", methodID[0:4]), map[string]interface{}{
		"name":                methodID[0:4],
		"identity_entity_ids": []string{entityID, otherEntityID},
		"mfa_method_ids":      []string{methodID},
	})
	require.NoError(t, err)

	userClient, err := client.Clone()
	require.NoError(t, err)
	userClient.ClearToken()

	pw, err := upAuth.NewUserpassAuth("user-limit", &upAuth.Password{FromString: "testpassword"})
	require.NoError(t, err)
	for i := 0; i < 100; i++ {
		secret, err := userClient.Auth().MFALogin(context.Background(), pw)
		require.NoError(t, err, "login %d should be accepted", i+1)
		require.NotNil(t, secret.Auth.MFARequirement)
	}

	_, err = userClient.Auth().MFALogin(context.Background(), pw)
	require.ErrorContains(t, err, "too many pending MFA requests")

	otherPw, err := upAuth.NewUserpassAuth("user-other", &upAuth.Password{FromString: "testpassword"})
	require.NoError(t, err)
	secret, err := userClient.Auth().MFALogin(context.Background(), otherPw)
	require.NoError(t, err)
	require.NotNil(t, secret.Auth.MFARequirement)
}

// TestLoginMFA_FailedValidationDoesNotLockUser verifies that wrong MFA codes
// sent to mfa/validate do not count toward Vault's user lockout, which tracks
// failed first-factor logins on userpass, ldap and approle mounts. The user
// sends more wrong codes than the lockout threshold and can still log in with
// the correct password. As a control, the same number of wrong passwords does
// lock the user out, which shows lockout is active on the mount.
func TestLoginMFA_FailedValidationDoesNotLockUser(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{"userpass": userpass.Factory},
		LogicalBackends:    map[string]logical.Factory{"totp": totp.Factory},
	})
	client := cluster.Cores[0].Client

	testhelpers.SetupTOTPMount(t, client)
	mountAccessor := testhelpers.SetupUserpassMountAccessor(t, client)
	_, err := client.Logical().Write("sys/auth/userpass/tune", map[string]interface{}{
		"user_lockout_config": map[string]interface{}{
			"lockout_threshold": "3",
			"lockout_duration":  "15m",
		},
	})
	require.NoError(t, err)

	_, entityID, _ := testhelpers.CreateEntityAndAlias(t, client, mountAccessor, "ent-lockout", "user-lockout")
	methodID := testhelpers.SetupTOTPMethod(t, client, map[string]interface{}{
		"issuer": "vault", "period": 30, "algorithm": "SHA1", "digits": 6,
		"skew": 1, "key_size": 20, "qr_size": 200, "max_validation_attempts": 5,
	})
	totpPath := testhelpers.RegisterEntityInTOTPEngine(t, client, entityID, methodID)

	userClient, err := client.Clone()
	require.NoError(t, err)
	userClient.ClearToken()

	pw, err := upAuth.NewUserpassAuth("user-lockout", &upAuth.Password{FromString: "testpassword"})
	require.NoError(t, err)
	mfaSecret, err := userClient.Auth().MFALogin(context.Background(), pw)
	require.NoError(t, err)

	// Derive a code that is guaranteed to be wrong by changing the last digit
	// of the current valid code.
	code := testhelpers.GetTOTPCodeFromEngine(t, client, totpPath)
	wrongCode := code[:5] + string(rune('0'+(code[5]-'0'+1)%10))

	// One more wrong code than the lockout threshold, and fewer than the TOTP
	// method's max_validation_attempts, so only user lockout could block the
	// next login.
	for i := 0; i < 4; i++ {
		_, err = userClient.Auth().MFAValidate(context.Background(), mfaSecret,
			map[string]interface{}{methodID: []string{wrongCode}})
		require.Error(t, err, "wrong code %d should fail validation", i+1)
	}

	secret, err := userClient.Auth().MFALogin(context.Background(), pw)
	require.NoError(t, err, "wrong MFA codes must not lock the user out")
	require.NotNil(t, secret.Auth.MFARequirement)

	// Control: wrong passwords up to the threshold do lock the user out.
	badPw, err := upAuth.NewUserpassAuth("user-lockout", &upAuth.Password{FromString: "wrong-password"})
	require.NoError(t, err)
	for i := 0; i < 3; i++ {
		_, err = userClient.Auth().MFALogin(context.Background(), badPw)
		require.Error(t, err)
	}
	_, err = userClient.Auth().MFALogin(context.Background(), pw)
	require.ErrorContains(t, err, "permission denied")
}
