// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	golangjwt "github.com/golang-jwt/jwt/v5"
	jwtauth "github.com/hashicorp/vault-plugin-auth-jwt"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestIdentityStore_ExternalGroupMemberDuplication_DuplicateGroupsClaim
// verifies that a JWT/OIDC login whose groups claim contains duplicate entries
// does not append the entity ID to member_entity_ids more than once (VAULT-48623).
func TestIdentityStore_ExternalGroupMemberDuplication_DuplicateGroupsClaim(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"jwt": jwtauth.Factory,
		},
	})
	client := cluster.Cores[0].Client

	// Generate a local RSA signing key; no external IdP is needed.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)
	pubKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyBytes}))

	require.NoError(t, client.Sys().EnableAuthWithOptions("jwt-dup", &api.EnableAuthOptions{
		Type: "jwt",
	}))

	_, err = client.Logical().Write("auth/jwt-dup/config", map[string]interface{}{
		"jwt_validation_pubkeys": []string{pubKeyPEM},
	})
	require.NoError(t, err)

	_, err = client.Logical().Write("auth/jwt-dup/role/dup-role", map[string]interface{}{
		"role_type":       "jwt",
		"bound_audiences": "vault-dup-aud",
		"user_claim":      "sub",
		"groups_claim":    "groups",
		"token_policies":  "default",
		"ttl":             "1h",
	})
	require.NoError(t, err)

	authList, err := client.Sys().ListAuth()
	require.NoError(t, err)
	mountAccessor := authList["jwt-dup/"].Accessor
	require.NotEmpty(t, mountAccessor)

	signJWT := func(sub string, groups []string) string {
		claims := golangjwt.MapClaims{
			"sub":    sub,
			"aud":    "vault-dup-aud",
			"iat":    time.Now().Unix(),
			"exp":    time.Now().Add(time.Hour).Unix(),
			"groups": groups,
		}
		token := golangjwt.NewWithClaims(golangjwt.SigningMethodRS256, claims)
		signed, signErr := token.SignedString(privateKey)
		require.NoError(t, signErr)
		return signed
	}

	// Seed the entity with no group claim so it is fully established before
	// the external group and its alias mapping are created.
	const subject = "dup-user"
	seedResp, err := client.Logical().Write("auth/jwt-dup/login", map[string]interface{}{
		"role": "dup-role",
		"jwt":  signJWT(subject, []string{}),
	})
	require.NoError(t, err)
	require.NotNil(t, seedResp, "seed login response must not be nil")
	require.NotNil(t, seedResp.Auth, "seed login must return auth")
	entityID := seedResp.Auth.EntityID
	require.NotEmpty(t, entityID)

	groupResp, err := client.Logical().Write("identity/group", map[string]interface{}{
		"name": "dup-test-engineering",
		"type": "external",
	})
	require.NoError(t, err)
	require.NotNil(t, groupResp, "group creation must return a response")
	groupID, ok := groupResp.Data["id"].(string)
	require.True(t, ok, "group response must contain a string id")
	require.NotEmpty(t, groupID)

	const claimName = "engineering"
	_, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           claimName,
		"mount_accessor": mountAccessor,
		"canonical_id":   groupID,
	})
	require.NoError(t, err)

	readMembers := func(label string) []string {
		t.Helper()
		secret, readErr := client.Logical().Read("identity/group/id/" + groupID)
		require.NoError(t, readErr)
		require.NotNil(t, secret)
		raw, _ := secret.Data["member_entity_ids"].([]interface{})
		ids := make([]string, 0, len(raw))
		for _, v := range raw {
			id, ok := v.(string)
			require.Truef(t, ok, "%s: member_entity_ids entry %v is not a string", label, v)
			ids = append(ids, id)
		}
		t.Logf("%s: member_entity_ids count=%d ids=%v", label, len(ids), ids)
		return ids
	}

	// requireSoleMember asserts the group's membership is exactly the seeded
	// entity — verifying the identity, not just the list length. This catches
	// both duplication ([entityID, entityID]) and any wrong/extra entity ID.
	requireSoleMember := func(label string) {
		t.Helper()
		require.Equalf(t, []string{entityID}, readMembers(label),
			"%s: member_entity_ids must contain exactly the seeded entity %s once", label, entityID)
	}

	require.Empty(t, readMembers("initial"), "group must start with zero members")

	// Login 1: single occurrence — normal path, entity must be added once.
	_, err = client.Logical().Write("auth/jwt-dup/login", map[string]interface{}{
		"role": "dup-role",
		"jwt":  signJWT(subject, []string{claimName}),
	})
	require.NoError(t, err, "login 1 (single group claim) must succeed")
	requireSoleMember("login 1")

	// Login 2: duplicate group name — triggers the bug without the fix.
	// diffGroups misclassifies the second copy of the resolved group as new,
	// causing an unguarded append even though the entity is already a member.
	_, err = client.Logical().Write("auth/jwt-dup/login", map[string]interface{}{
		"role": "dup-role",
		"jwt":  signJWT(subject, []string{claimName, claimName}),
	})
	require.NoError(t, err, "login 2 (duplicate group claim) must succeed")
	requireSoleMember("login 2")

	// Login 3: repeat to confirm no unbounded growth on subsequent logins.
	_, err = client.Logical().Write("auth/jwt-dup/login", map[string]interface{}{
		"role": "dup-role",
		"jwt":  signJWT(subject, []string{claimName, claimName}),
	})
	require.NoError(t, err, "login 3 (duplicate group claim) must succeed")
	requireSoleMember("login 3")
}
