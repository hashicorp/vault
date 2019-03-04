package token

import (
	"encoding/base64"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/api"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/jsonutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestTokenStore_CreateOrphanResponse(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	secret, err := client.Auth().Token().CreateOrphan(&api.TokenCreateRequest{
		Policies: []string{"default"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !secret.Auth.Orphan {
		t.Fatalf("failed to set orphan as true, got: %#v", secret.Auth)
	}
}

func TestTokenStore_TokenInvalidEntityID(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
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

	// Enable userpass auth
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Add a user to userpass backend
	_, err = client.Logical().Write("auth/userpass/users/testuser", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}

	secret, err := client.Logical().Write("auth/userpass/login/testuser", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	clientToken := secret.Auth.ClientToken

	secret, err = client.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": clientToken,
	})
	if err != nil {
		t.Fatal(err)
	}

	entityID := secret.Data["entity_id"].(string)

	_, err = client.Logical().Delete("identity/entity/id/" + entityID)
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken(clientToken)

	secret, err = client.Logical().Write("auth/token/lookup-self", nil)
	if err == nil {
		t.Fatalf("expected error due to token being invalid when its entity is invalid")
	}
}

func TestTokenStore_IdentityPolicies(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"ldap": credLdap.Factory,
		},
		EnableRaw: true,
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

	// Create user in LDAP auth. We add two groups, but we should filter out
	// the ones that don't match aliases later (we will check for this)
	_, err = client.Logical().Write("auth/ldap/users/tesla", map[string]interface{}{
		"policies": "default",
		"groups":   "testgroup1,testgroup2",
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

	expectedPolicies := []string{
		"default",
		"testgroup1-policy",
	}
	if !reflect.DeepEqual(expectedPolicies, secret.Auth.Policies) {
		t.Fatalf("bad: identity policies; expected: %#v\nactual: %#v", expectedPolicies, secret.Auth.Policies)
	}

	// At this point there shouldn't be any identity policy on the token
	secret, err = client.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": ldapClientToken,
	})
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
	secret, err = client.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": ldapClientToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	identityPolicies := secret.Data["identity_policies"].([]interface{})
	var actualPolicies []string
	for _, item := range identityPolicies {
		actualPolicies = append(actualPolicies, item.(string))
	}
	sort.Strings(actualPolicies)

	expectedPolicies = []string{
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
	secret, err = client.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": ldapClientToken,
	})
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
	secret, err = client.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": ldapClientToken,
	})
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

	// Log in and get a new token, then renew it. See issue #4829. The logic is
	// continued after the next block.
	secret, err = client.Logical().Write("auth/ldap/login/tesla", map[string]interface{}{
		"password": "password",
	})
	if err != nil {
		t.Fatal(err)
	}
	token4829 := secret.Auth.ClientToken

	// Check that the lease for the token contains only the single group; this
	// should be true for both as one was fresh and the other was a renew
	// (which is why we do the renew check on the 4839 token after this block)
	secret, err = client.Logical().List("sys/raw/sys/expire/id/auth/ldap/login/tesla/")
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range secret.Data["keys"].([]interface{}) {
		secret, err := client.Logical().Read("sys/raw/sys/expire/id/auth/ldap/login/tesla/" + key.(string))
		if err != nil {
			t.Fatal(err)
		}
		//t.Logf("%#v", *secret)
		var resp logical.Response
		if err := jsonutil.DecodeJSON([]byte(secret.Data["value"].(string)), &resp); err != nil {
			t.Fatal(err)
		}
		if len(resp.Auth.GroupAliases) != 1 || resp.Auth.GroupAliases[0].Name != "testgroup1" {
			t.Fatalf("bad: %#v", resp.Auth.GroupAliases)
		}
	}

	secret, err = client.Auth().Token().Renew(token4829, 10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTokenStore_CIDRBlocks(t *testing.T) {
	testPolicy := `
path "auth/token/create" {
	capabilities = ["update"]
}
`

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

	_, err = client.Logical().Write("sys/policies/acl/test", map[string]interface{}{
		"policy": testPolicy,
	})
	if err != nil {
		t.Fatal(err)
	}

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
		"bound_cidrs":      []string{"127.0.0.1/32", "1.2.3.4/8", "5.6.7.8/24"},
		"allowed_policies": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	secret, err = client.Auth().Token().CreateWithRole(&api.TokenCreateRequest{
		Policies: []string{"test", "default"},
	}, "testrole")
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(secret.Auth.ClientToken)
	_, err = client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}

	// Before moving on, validate that a child token created from this token
	// inherits the bound cidr blocks
	client.SetToken(secret.Auth.ClientToken)
	childSecret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"default"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(childSecret.Auth.ClientToken)
	childInfo, err := client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(childInfo.Data["bound_cidrs"], []interface{}{"127.0.0.1", "1.2.3.4/8", "5.6.7.8/24"}); diff != nil {
		t.Fatal(diff)
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
		"bound_cidrs":      []string{"1.2.3.4/8", "5.6.7.8/24"},
		"allowed_policies": "",
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

func TestTokenStore_RevocationOnStartup(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client
	rootToken := client.Token()

	type leaseEntry struct {
		LeaseID         string                 `json:"lease_id"`
		ClientToken     string                 `json:"client_token"`
		Path            string                 `json:"path"`
		Data            map[string]interface{} `json:"data"`
		Secret          *logical.Secret        `json:"secret"`
		Auth            *logical.Auth          `json:"auth"`
		IssueTime       time.Time              `json:"issue_time"`
		ExpireTime      time.Time              `json:"expire_time"`
		LastRenewalTime time.Time              `json:"last_renewal_time"`
	}

	var secret *api.Secret
	var err error
	var tokens []string
	// Create tokens
	for i := 0; i < 500; i++ {
		secret, err = client.Auth().Token().Create(&api.TokenCreateRequest{
			Policies: []string{"default"},
		})
		if err != nil {
			t.Fatal(err)
		}
		tokens = append(tokens, secret.Auth.ClientToken)
	}

	const tokenPath string = "sys/raw/sys/token/id/"
	secret, err = client.Logical().List(tokenPath)
	if err != nil {
		t.Fatal(err)
	}
	totalTokens := len(secret.Data["keys"].([]interface{}))

	// Get the list of leases
	const leasePath string = "sys/raw/sys/expire/id/auth/token/create/"
	secret, err = client.Logical().List(leasePath)
	if err != nil {
		t.Fatal(err)
	}
	leases := secret.Data["keys"].([]interface{})
	if len(leases) != 500 {
		t.Fatalf("unexpected number of leases: %d", len(leases))
	}

	// Holds non-root leases
	var validLeases []string
	// Fake times in the past
	for _, lease := range leases {
		secret, err = client.Logical().Read(leasePath + lease.(string))
		var entry leaseEntry
		if err := jsonutil.DecodeJSON([]byte(secret.Data["value"].(string)), &entry); err != nil {
			t.Fatal(err)
		}
		if entry.ExpireTime.IsZero() {
			continue
		}
		validLeases = append(validLeases, lease.(string))
		entry.IssueTime = entry.IssueTime.Add(-1 * time.Hour * 24 * 365)
		entry.ExpireTime = entry.ExpireTime.Add(-1 * time.Hour * 24 * 365)
		jsonEntry, err := jsonutil.EncodeJSON(&entry)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := client.Logical().Write(leasePath+lease.(string), map[string]interface{}{
			"value": string(jsonEntry),
		}); err != nil {
			t.Fatal(err)
		}
	}

	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	var status *api.SealStatusResponse
	for i := 0; i < len(cluster.BarrierKeys); i++ {
		status, err = client.Sys().Unseal(string(base64.StdEncoding.EncodeToString(cluster.BarrierKeys[i])))
		if err != nil {
			t.Fatal(err)
		}
		if !status.Sealed {
			break
		}
	}
	if status.Sealed {
		t.Fatal("did not unseal properly")
	}

	// Give lease loading some time to process
	time.Sleep(5 * time.Second)

	for i, token := range tokens {
		client.SetToken(token)
		_, err := client.Logical().Write("cubbyhole/foo", map[string]interface{}{
			"value": "bar",
		})
		if err == nil {
			t.Errorf("expected error but did not get one, token num %d", i)
		}
	}

	expectedLeases := len(leases) - len(validLeases)

	client.SetToken(rootToken)
	secret, err = client.Logical().List(leasePath)
	if err != nil {
		t.Fatal(err)
	}

	switch {
	case secret == nil:
		if expectedLeases != 0 {
			t.Fatalf("nil secret back but expected %d leases", expectedLeases)
		}

	case secret.Data == nil:
		if expectedLeases != 0 {
			t.Fatalf("nil secret data back but expected %d leases, secret is %#v", expectedLeases, *secret)
		}

	default:
		leasesLeft := len(secret.Data["keys"].([]interface{}))
		if leasesLeft != expectedLeases {
			t.Fatalf("found %d leases left, expected %d", leasesLeft, expectedLeases)
		}
	}

	expectedTokens := totalTokens - len(validLeases)
	secret, err = client.Logical().List(tokenPath)
	if err != nil {
		t.Fatal(err)
	}
	tokensLeft := len(secret.Data["keys"].([]interface{}))
	if tokensLeft != expectedTokens {
		t.Fatalf("found %d tokens left, expected %d", tokensLeft, expectedTokens)
	}
}
