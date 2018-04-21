package vault_test

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestTokenStore_IdentityPolicies(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"ldap": credLdap.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// Enable LDAP auth
	err := client.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{
		Type: "ldap",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Configure LDAP auth
	_, err = client.Logical().Write("auth/ldap/config", map[string]interface{}{
		"url":      "ldap://ldap.forumsys.com",
		"userattr": "uid",
		"userdn":   "dc=example,dc=com",
		"groupdn":  "dc=example,dc=com",
		"binddn":   "cn=read-only-admin,dc=example,dc=com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create group in LDAP auth
	_, err = client.Logical().Write("auth/ldap/groups/testgroup1", map[string]interface{}{
		"policies": "testgroup1-policy",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create user in LDAP auth
	_, err = client.Logical().Write("auth/ldap/users/tesla", map[string]interface{}{
		"policies": "default",
		"groups":   "testgroup1",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Login using LDAP
	secret, err := client.Logical().Write("auth/ldap/login/tesla", map[string]interface{}{
		"password": "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	ldapClientToken := secret.Auth.ClientToken

	// At this point there shouldn't be any identity policy on the token
	secret, err = client.Logical().Read("auth/token/lookup/" + ldapClientToken)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := secret.Data["identity_policies"]
	if ok {
		t.Fatalf("identity_policies should not have been set")
	}

	// Extract the entity ID of the token and set some policies on the entity
	entityID := secret.Data["entity_id"].(string)
	_, err = client.Logical().Write("identity/entity/id/"+entityID, map[string]interface{}{
		"policies": []string{
			"entity_policy_1",
			"entity_policy_2",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Lookup the token and expect entity policies on the token
	secret, err = client.Logical().Read("auth/token/lookup/" + ldapClientToken)
	if err != nil {
		t.Fatal(err)
	}
	identityPolicies := secret.Data["identity_policies"].([]interface{})
	var actualPolicies []string
	for _, item := range identityPolicies {
		actualPolicies = append(actualPolicies, item.(string))
	}
	sort.Strings(actualPolicies)

	expectedPolicies := []string{
		"entity_policy_1",
		"entity_policy_2",
	}
	sort.Strings(expectedPolicies)
	if !reflect.DeepEqual(expectedPolicies, actualPolicies) {
		t.Fatalf("bad: identity policies; expected: %#v\nactual: %#v", expectedPolicies, actualPolicies)
	}

	// Create identity group and add entity as its member
	secret, err = client.Logical().Write("identity/group", map[string]interface{}{
		"policies": []string{
			"group_policy_1",
			"group_policy_2",
		},
		"member_entity_ids": []string{
			entityID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Lookup token and expect both entity and group policies on the token
	secret, err = client.Logical().Read("auth/token/lookup/" + ldapClientToken)
	if err != nil {
		t.Fatal(err)
	}
	identityPolicies = secret.Data["identity_policies"].([]interface{})
	actualPolicies = nil
	for _, item := range identityPolicies {
		actualPolicies = append(actualPolicies, item.(string))
	}
	sort.Strings(actualPolicies)

	expectedPolicies = []string{
		"entity_policy_1",
		"entity_policy_2",
		"group_policy_1",
		"group_policy_2",
	}
	sort.Strings(expectedPolicies)
	if !reflect.DeepEqual(expectedPolicies, actualPolicies) {
		t.Fatalf("bad: identity policies; expected: %#v\nactual: %#v", expectedPolicies, actualPolicies)
	}

	// Create an external group and renew the token. This should add external
	// group policies to the token.
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	ldapMountAccessor1 := auths["ldap/"].Accessor

	// Create an external group
	secret, err = client.Logical().Write("identity/group", map[string]interface{}{
		"type": "external",
		"policies": []string{
			"external_group_policy_1",
			"external_group_policy_2",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	ldapExtGroupID1 := secret.Data["id"].(string)

	// Associate a group from LDAP auth as a group-alias in the external group
	_, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           "testgroup1",
		"mount_accessor": ldapMountAccessor1,
		"canonical_id":   ldapExtGroupID1,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Renew token to refresh external group memberships
	secret, err = client.Auth().Token().Renew(ldapClientToken, 10)
	if err != nil {
		t.Fatal(err)
	}

	// Lookup token and expect entity, group and external group policies on the
	// token
	secret, err = client.Logical().Read("auth/token/lookup/" + ldapClientToken)
	if err != nil {
		t.Fatal(err)
	}
	identityPolicies = secret.Data["identity_policies"].([]interface{})
	actualPolicies = nil
	for _, item := range identityPolicies {
		actualPolicies = append(actualPolicies, item.(string))
	}
	sort.Strings(actualPolicies)

	expectedPolicies = []string{
		"entity_policy_1",
		"entity_policy_2",
		"group_policy_1",
		"group_policy_2",
		"external_group_policy_1",
		"external_group_policy_2",
	}
	sort.Strings(expectedPolicies)
	if !reflect.DeepEqual(expectedPolicies, actualPolicies) {
		t.Fatalf("bad: identity policies; expected: %#v\nactual: %#v", expectedPolicies, actualPolicies)
	}
}

func TestTokenStore_CIDRBlocks(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client
	rootToken := client.Token()

	var err error
	var secret *api.Secret

	// Test normally
	_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
		"bound_cidrs": []string{},
	})
	if err != nil {
		t.Fatal(err)
	}
	secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
		Policies: []string{"default"},
	}, "testrole")
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(secret.Auth.ClientToken)
	_, err = client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}

	// CIDR blocks, containing localhost
	client.SetToken(rootToken)
	_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
		"bound_cidrs": []string{"127.0.0.1/32", "1.2.3.4/8", "5.6.7.8/24"},
	})
	if err != nil {
		t.Fatal(err)
	}
	secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
		Policies: []string{"default"},
	}, "testrole")
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(secret.Auth.ClientToken)
	_, err = client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}

	// CIDR blocks, not containing localhost (should fail)
	client.SetToken(rootToken)
	_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
		"bound_cidrs": []string{"1.2.3.4/8", "5.6.7.8/24"},
	})
	if err != nil {
		t.Fatal(err)
	}
	secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
		Policies: []string{"default"},
	}, "testrole")
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(secret.Auth.ClientToken)
	_, err = client.Auth().Token().LookupSelf()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "permission denied") {
		t.Fatalf("unexpected error: %v", err)
	}

	// Root token, no ttl, should work
	client.SetToken(rootToken)
	_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
		"bound_cidrs": []string{"1.2.3.4/8", "5.6.7.8/24"},
	})
	if err != nil {
		t.Fatal(err)
	}
	secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{}, "testrole")
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(secret.Auth.ClientToken)
	_, err = client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}

	// Root token, ttl, should not work
	client.SetToken(rootToken)
	_, err = client.Logical().Write("auth/token/roles/testrole", map[string]interface{}{
		"bound_cidrs": []string{"1.2.3.4/8", "5.6.7.8/24"},
		"period":      3600,
	})
	if err != nil {
		t.Fatal(err)
	}
	secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{}, "testrole")
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(secret.Auth.ClientToken)
	_, err = client.Auth().Token().LookupSelf()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "permission denied") {
		t.Fatalf("unexpected error: %v", err)
	}
}
