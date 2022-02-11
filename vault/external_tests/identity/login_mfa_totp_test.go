package identity

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

var loginMFACoreConfig = &vault.CoreConfig{
	CredentialBackends: map[string]logical.Factory{
		"userpass": userpass.Factory,
	},
	LogicalBackends: map[string]logical.Factory{
		"totp": totp.Factory,
	},
}

type totpCode struct {
	name          string
	methodID      string
	namespacePath string // this is tied to the entityID or the mount accessor
	entityID      string
}

func getNamespaceSpecificMountAccessor(namespace string, client *api.Client, t *testing.T) string {
	client.SetNamespace(namespace)
	auths, err := client.Sys().ListAuth()
	if err != nil || auths == nil || auths["userpass/"] == nil {
		t.Fatalf("failed to get the list of auths")
	}
	return auths["userpass/"].Accessor
}

func TestLoginMfaGenerateTOTPRoleTest(t *testing.T) {
	cluster := vault.NewTestCluster(t, loginMFACoreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	// Mount the TOTP backend
	mountInfo := &api.MountInput{
		Type: "totp",
	}
	err := client.Sys().Mount("totp", mountInfo)
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
		resp1, err := client.Logical().Write("identity/mfa/method-id/totp", map[string]interface{}{
			"issuer":    "yCorp",
			"period":    10000,
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

		secret, err := client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-generate"), map[string]interface{}{
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
		secret, err = user2Client.Logical().Write("sys/mfa/validate", map[string]interface{}{
			"mfa_request_id": secret.Auth.MFARequirement.MFARequestID,
			"mfa_payload": map[string][]string{
				methodID: {totpPasscode},
			},
		})
		if err != nil {
			t.Fatalf("MFA failed: %v", err)
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
		// give it enough time to make sure the request has expired
		time.Sleep(605 * time.Second)
		_, err = user2Client.Logical().Write("sys/mfa/validate", map[string]interface{}{
			"mfa_request_id": secret.Auth.MFARequirement.MFARequestID,
			"mfa_payload": map[string][]string{
				methodID: {totpPasscode},
			},
		})
		if err == nil {
			t.Fatalf("MFA succeeded: %v", err)
		}

		// Destroy the secret so that the token can self generate
		_, err = userClient.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-destroy"), map[string]interface{}{
			"entity_id": entityID,
			"method_id": methodID,
		})
		if err != nil {
			t.Fatalf("failed to destroy the MFA secret: %s", err)
		}
	}
}

//- an enforcement can be defined in any NS and applies to that NS and its children
//- a methodid can be defined in any NS and may be referenced by an enforcement in
//  that NS or its children
//- an entity may configure TOTP keys for methods defined in the entity's NS or
//  its parents
func TestNamespaceLoginMfaTotp(t *testing.T) {
	cluster := vault.NewTestCluster(t, loginMFACoreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	client := cluster.Cores[0].Client

	// Mount the TOTP backend
	mountInfo := &api.MountInput{
		Type: "totp",
	}
	err := client.Sys().Mount("totp", mountInfo)
	if err != nil {
		t.Fatalf("failed to mount totp backend: %v", err)
	}

	// Mount Userpass for the root namespace
	err = client.Sys().EnableAuthWithOptions("userpassRoot", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass auth: %v", err)
	}

	// Setup Namespaces and create Userpass mounts in each
	// creating namespaces ns1, ns1/ns2, and ns1/ns2/ns3

	// NS1
	ns1Path := testhelpers.CreateNamespace(t, client, "ns1/", "")

	// Enable Userpass authentication for the NS
	client.SetNamespace(ns1Path)
	err = client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass auth: %v", err)
	}
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatalf("bb")
	}
	mountAccessor1 := auths["userpass/"].Accessor

	resp, err := client.Logical().Write("identity/entity", map[string]interface{}{
		"name": "test-entity1",
		"metadata": map[string]string{
			"email":        "test@hashicorp.com",
			"phone_number": "123-456-7890",
		},
	})
	if err != nil {
		t.Fatalf("failed to creat an entity: %v", err)
	}
	entityID1 := resp.Data["id"].(string)

	_, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "testuser1",
		"canonical_id":   entityID1,
		"mount_accessor": getNamespaceSpecificMountAccessor(ns1Path, client, t),
	})
	if err != nil {
		t.Fatalf("failed to create an entity-alias: %v", err)
	}

	// Creating a user in the userpass auth mount
	_, err = client.Logical().Write("auth/userpass/users/testuser1", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("failed to configure userpass backend: %v", err)
	}

	// NS2
	// Create second namespace ns2
	ns2Path := testhelpers.CreateNamespace(t, client, "ns2/", "ns1/")
	client.SetNamespace(ns2Path)
	err = client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass auth: %v", err)
	}
	auths, err = client.Sys().ListAuth()
	if err != nil {
		t.Fatalf("bb")
	}
	mountAccessor2 := auths["userpass/"].Accessor

	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{
		"name": "test-entity2",
		"metadata": map[string]string{
			"email":        "test@hashicorp.com",
			"phone_number": "123-456-7890",
		},
	})
	if err != nil {
		t.Fatalf("failed to creat an entity: %v", err)
	}
	entityID2 := resp.Data["id"].(string)

	_, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "testuser2",
		"canonical_id":   entityID2,
		"mount_accessor": getNamespaceSpecificMountAccessor(ns2Path, client, t),
	})
	if err != nil {
		t.Fatalf("failed to create an entity-alias: %v", err)
	}
	// Creating a user in the userpass auth mount
	_, err = client.Logical().Write("auth/userpass/users/testuser2", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("failed to configure userpass backend: %v", err)
	}

	// NS3
	// Create third namespace within ns1
	ns3Path := testhelpers.CreateNamespace(t, client, "ns3/", "ns1/ns2/")
	client.SetNamespace(ns3Path)
	err = client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass auth: %v", err)
	}
	auths, err = client.Sys().ListAuth()
	if err != nil {
		t.Fatalf("bb")
	}
	mountAccessor3 := auths["userpass/"].Accessor
	// mountAccessors := []string{mountAccessor1, mountAccessor2, mountAccessor3}

	resp, err = client.Logical().Write("identity/entity", map[string]interface{}{
		"name": "test-entity3",
		"metadata": map[string]string{
			"email":        "test@hashicorp.com",
			"phone_number": "123-456-7890",
		},
	})
	if err != nil {
		t.Fatalf("failed to creat an entity: %v", err)
	}
	entityID3 := resp.Data["id"].(string)

	_, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "testuser3",
		"canonical_id":   entityID3,
		"mount_accessor": getNamespaceSpecificMountAccessor(ns3Path, client, t),
	})
	if err != nil {
		t.Fatalf("failed to create an entity-alias: %v", err)
	}
	// Creating a user in the userpass auth mount
	_, err = client.Logical().Write("auth/userpass/users/testuser3", map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("failed to configure userpass backend: %v", err)
	}

	namespaceEntityIDs := []string{entityID1, entityID2, entityID3}

	// Creating a group for all entities in all namespaces
	client.SetNamespace("")
	// Create a group
	resp, err = client.Logical().Write("identity/group", map[string]interface{}{
		"name":              "engineering",
		"member_entity_ids": []string{entityID1, entityID2, entityID3},
	})
	if err != nil {
		t.Fatalf("failed to create a group: %v", err)
	}

	// groupID := resp.Data["id"].(string)

	namespacePaths := []string{ns1Path, ns2Path, ns3Path}
	var namespaceMethodIDs []string
	for _, nsPath := range namespacePaths {
		// MFA TOTP Method for various NS
		client.SetNamespace(nsPath)
		// create a config
		resp, err = client.Logical().Write("identity/mfa/method-id/totp", map[string]interface{}{
			"issuer":    "yCorp",
			"period":    10000,
			"algorithm": "SHA1",
			"digits":    6,
			"skew":      1,
			"key_size":  10,
			"qr_size":   100,
		})

		if err != nil || (resp == nil) {
			t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
		}

		methodID := resp.Data["method_id"].(string)
		if methodID == "" {
			t.Fatalf("method ID is empty")
		}
		namespaceMethodIDs = append(namespaceMethodIDs, methodID)
	}

	methodIDTotpCodeNameMap := checkGenerateTotp(client, namespacePaths, namespaceMethodIDs, namespaceEntityIDs, t)

	// Creating Login enforcement in NS1
	client.SetNamespace(ns1Path)
	// creating MFAEnforcementConfig
	_, err = client.Logical().Write("identity/mfa/login-enforcement/LE11", map[string]interface{}{
		"auth_method_accessors": []string{mountAccessor1},
		"identity_entity_ids":   []string{entityID1},
		"name":                  "LE11",
		"mfa_method_ids":        []string{namespaceMethodIDs[0]},
	})
	if err != nil {
		t.Fatalf("failed to configure MFAEnforcementConfig: %v", err)
	}

	// creating MFAEnforcementConfig
	client.SetNamespace(ns2Path)
	_, err = client.Logical().Write("identity/mfa/login-enforcement/LE21", map[string]interface{}{
		"auth_method_accessors": []string{mountAccessor2},
		"identity_entity_ids":   []string{entityID2},
		"name":                  "LE21",
		"mfa_method_ids":        []string{namespaceMethodIDs[0]},
	})
	if err != nil {
		t.Fatalf("failed to configure MFAEnforcementConfig: %v", err)
	}

	_, err = client.Logical().Write("identity/mfa/login-enforcement/LE22", map[string]interface{}{
		"auth_method_accessors": []string{mountAccessor2},
		"identity_entity_ids":   []string{entityID2},
		"name":                  "LE22",
		"mfa_method_ids":        []string{namespaceMethodIDs[1]},
	})
	if err != nil {
		t.Fatalf("failed to configure MFAEnforcementConfig: %v", err)
	}

	client.SetNamespace(ns3Path)
	// creating MFAEnforcementConfig
	_, err = client.Logical().Write("identity/mfa/login-enforcement/LE33", map[string]interface{}{
		"auth_method_accessors": []string{mountAccessor3},
		"identity_entity_ids":   []string{entityID3},
		"name":                  "LE33",
		"mfa_method_ids":        []string{namespaceMethodIDs[2]},
	})
	if err != nil {
		t.Fatalf("failed to configure MFAEnforcementConfig: %v", err)
	}

	singlePhaseLogin(client, t, "testuser1", entityID1, methodIDTotpCodeNameMap)
	singlePhaseLogin(client, t, "testuser2", entityID2, methodIDTotpCodeNameMap)
	singlePhaseLogin(client, t, "testuser3", entityID3, methodIDTotpCodeNameMap)
	twoPhaseLogin(client, t, "testuser1", entityID1, methodIDTotpCodeNameMap)
	twoPhaseLogin(client, t, "testuser2", entityID2, methodIDTotpCodeNameMap)
	twoPhaseLogin(client, t, "testuser3", entityID3, methodIDTotpCodeNameMap)
	twoPhaseLoginDifferentNamespace(client, t, "testuser3", entityID3, methodIDTotpCodeNameMap)

	checkDestroyTotp(client, namespacePaths, namespaceMethodIDs, namespaceEntityIDs, t)
}

func singlePhaseLogin(client *api.Client, t *testing.T, username, entityID string, totpCodeMap map[string][]*totpCode) {
	headers := client.Headers()
	headers.Del("X-Vault-MFA")
	client.SetHeaders(headers)

	// getting the passcode
	client.SetNamespace("")

	totpCodeStruct := totpCodeMap[entityID]
	for _, codeStruct := range totpCodeStruct {
		secret, err := client.Logical().Read(fmt.Sprintf("totp/code/%s", codeStruct.name))
		if err != nil {
			t.Fatalf("failed to create totp passcode: %v", err)
		}
		passCode, ok := secret.Data["code"].(string)
		if !ok && passCode == "" {
			t.Fatalf("failed to generate a totp passcode")
		}

		// MFA single-phase login
		client.AddHeader("X-Vault-MFA", fmt.Sprintf("%s:%s", codeStruct.methodID, passCode))
	}

	// namespace is the same for the same entityID
	client.SetNamespace(totpCodeStruct[0].namespacePath)
	secret, err := client.Logical().Write(fmt.Sprintf("auth/userpass/login/%s", username), map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("%s MFA failed: %v", username, err)
	}

	userpassToken := secret.Auth.ClientToken

	secret, err = client.Logical().Write("auth/token/lookup", map[string]interface{}{
		"token": userpassToken,
	})
	if err != nil {
		t.Fatalf("failed to lookup userpass authenticated token: %v", err)
	}

	entityIDCheck := secret.Data["entity_id"].(string)
	if entityIDCheck != entityID {
		t.Fatalf("different entityID assigned")
	}
}

func twoPhaseLogin(client *api.Client, t *testing.T, username, entityID string, totpCodeMap map[string][]*totpCode) {
	// Two-phase login
	headers := client.Headers()
	headers.Del("X-Vault-MFA")
	client.SetHeaders(headers)

	client.SetNamespace("")
	totpCodeStruct := totpCodeMap[entityID]

	methodIDPasscodeMap := make(map[string][]string, 0)
	for _, codeStruct := range totpCodeStruct {
		secret, err := client.Logical().Read(fmt.Sprintf("totp/code/%s", codeStruct.name))
		if err != nil {
			t.Fatalf("failed to create totp passcode: %v", err)
		}
		passCode, ok := secret.Data["code"].(string)
		if !ok && passCode == "" {
			t.Fatalf("failed to generate a totp passcode")
		}
		methodIDPasscodeMap[codeStruct.methodID] = []string{passCode}
	}

	// namespace is the same for the same entityID
	client.SetNamespace(totpCodeStruct[0].namespacePath)

	secret, err := client.Logical().Write(fmt.Sprintf("auth/userpass/login/%s", username), map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("MFA failed: %v", err)
	}

	if secret.Auth == nil || secret.Auth.MFARequirement == nil {
		t.Fatalf("two phase login returned nil MFARequirement for username %s", username)
	}
	if secret.Auth.MFARequirement.MFARequestID == "" {
		t.Fatalf("MFARequirement contains empty MFARequestID")
	}
	if secret.Auth.MFARequirement.MFAConstraints == nil || len(secret.Auth.MFARequirement.MFAConstraints) == 0 {
		t.Fatalf("MFAConstraints is nil or empty")
	}

	// validation
	secretValidated, err := client.Logical().Write("sys/mfa/validate", map[string]interface{}{
		"mfa_request_id": secret.Auth.MFARequirement.MFARequestID,
		"mfa_payload":    methodIDPasscodeMap,
	})
	if err != nil || secretValidated == nil {
		t.Fatalf("MFA failed: %v", err)
	}

	// validate the same request the second time should fail
	secret, err = client.Logical().Write("sys/mfa/validate", map[string]interface{}{
		"mfa_request_id": secret.Auth.MFARequirement.MFARequestID,
		"mfa_payload":    methodIDPasscodeMap,
	})
	if err == nil {
		t.Fatalf("MFA validate did not fail as expected")
	}
	if !strings.Contains(err.Error(), "invalid request ID") {
		t.Fatalf("expected error invalid request ID, got %s", err.Error())
	}
}

func twoPhaseLoginDifferentNamespace(client *api.Client, t *testing.T, username, entityID string, totpCodeMap map[string][]*totpCode) {
	// Two-phase login
	headers := client.Headers()
	headers.Del("X-Vault-MFA")
	client.SetHeaders(headers)

	client.SetNamespace("")
	totpCodeStruct := totpCodeMap[entityID]

	methodIDPasscodeMap := make(map[string][]string, 0)
	for _, codeStruct := range totpCodeStruct {
		secret, err := client.Logical().Read(fmt.Sprintf("totp/code/%s", codeStruct.name))
		if err != nil {
			t.Fatalf("failed to create totp passcode: %v", err)
		}
		passCode, ok := secret.Data["code"].(string)
		if !ok && passCode == "" {
			t.Fatalf("failed to generate a totp passcode")
		}
		methodIDPasscodeMap[codeStruct.methodID] = []string{passCode}
	}

	// namespace is the same for the same entityID
	client.SetNamespace(totpCodeStruct[0].namespacePath)

	secret, err := client.Logical().Write(fmt.Sprintf("auth/userpass/login/%s", username), map[string]interface{}{
		"password": "testpassword",
	})
	if err != nil {
		t.Fatalf("MFA failed: %v", err)
	}

	if secret.Auth == nil || secret.Auth.MFARequirement == nil {
		t.Fatalf("two phase login returned nil MFARequirement for username %s", username)
	}
	if secret.Auth.MFARequirement.MFARequestID == "" {
		t.Fatalf("MFARequirement contains empty MFARequestID")
	}
	if secret.Auth.MFARequirement.MFAConstraints == nil || len(secret.Auth.MFARequirement.MFAConstraints) == 0 {
		t.Fatalf("MFAConstraints is nil or empty")
	}

	// validation in a different namespace
	client.SetNamespace("")
	_, err = client.Logical().Write("sys/mfa/validate", map[string]interface{}{
		"mfa_request_id": secret.Auth.MFARequirement.MFARequestID,
		"mfa_payload":    methodIDPasscodeMap,
	})
	if err == nil {
		t.Fatalf("expected MFA validate to fail: %v", err)
	}
	if !strings.Contains(err.Error(), "original request was issued in a different namesapce") {
		t.Fatalf("unexpected error returned")
	}
}

func checkGenerateTotp(client *api.Client, namespacePaths, namespaceMethodIDs, namespaceEntityIDs []string, t *testing.T) map[string][]*totpCode {
	codeNameStructMap := make(map[string][]*totpCode, 0)

	// generating totp in different namespace than entity namespace should fail
	client.SetNamespace("")
	secret, err := client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-generate"), map[string]interface{}{
		"entity_id": namespaceEntityIDs[0],
		"method_id": namespaceMethodIDs[1],
	})
	if err == nil {
		t.Fatalf("1failed to generate a TOTP secret on an entity: %v", err)
	}

	// non-root namespace
	client.SetNamespace(namespacePaths[2])
	secret, err = client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-generate"), map[string]interface{}{
		"entity_id": namespaceEntityIDs[2],
		"method_id": namespaceMethodIDs[2],
	})
	if err != nil {
		t.Fatalf("2failed to generate a TOTP secret on an entity: %v", err)
	}
	totpURL := secret.Data["url"].(string)

	name := namespaceMethodIDs[2] + namespaceEntityIDs[2]
	registerTotpUrl(client, totpURL, name, t)
	codeNameStructMap[namespaceEntityIDs[2]] = append(codeNameStructMap[namespaceEntityIDs[2]], &totpCode{
		name:          name,
		methodID:      namespaceMethodIDs[2],
		namespacePath: namespacePaths[2],
		entityID:      namespaceEntityIDs[2],
	})

	client.SetNamespace(namespacePaths[0])
	secret, err = client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-generate"), map[string]interface{}{
		"entity_id": namespaceEntityIDs[0],
		"method_id": namespaceMethodIDs[0],
	})
	if err != nil {
		t.Fatalf("3failed to generate a TOTP secret on an entity: %v", err)
	}
	totpURL = secret.Data["url"].(string)

	name = namespaceMethodIDs[0] + namespaceEntityIDs[0]
	registerTotpUrl(client, totpURL, name, t)
	codeNameStructMap[namespaceEntityIDs[0]] = append(codeNameStructMap[namespaceEntityIDs[0]], &totpCode{
		name:          name,
		methodID:      namespaceMethodIDs[0],
		namespacePath: namespacePaths[0],
		entityID:      namespaceEntityIDs[0],
	})

	client.SetNamespace(namespacePaths[1])
	secret, err = client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-generate"), map[string]interface{}{
		"entity_id": namespaceEntityIDs[1],
		"method_id": namespaceMethodIDs[0],
	})
	if err != nil {
		t.Fatalf("failed to generate a TOTP secret on an entity: %v", err)
	}
	totpURL = secret.Data["url"].(string)

	name = namespaceMethodIDs[0] + namespaceEntityIDs[1]
	registerTotpUrl(client, totpURL, name, t)
	codeNameStructMap[namespaceEntityIDs[1]] = append(codeNameStructMap[namespaceEntityIDs[1]], &totpCode{
		name:          name,
		methodID:      namespaceMethodIDs[0],
		namespacePath: namespacePaths[1],
		entityID:      namespaceEntityIDs[1],
	})

	client.SetNamespace(namespacePaths[1])
	secret, err = client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-generate"), map[string]interface{}{
		"entity_id": namespaceEntityIDs[1],
		"method_id": namespaceMethodIDs[1],
	})
	if err != nil {
		t.Fatalf("5failed to generate a TOTP secret on an entity: %v", err)
	}
	totpURL = secret.Data["url"].(string)

	name = namespaceMethodIDs[1] + namespaceEntityIDs[1]
	registerTotpUrl(client, totpURL, name, t)
	codeNameStructMap[namespaceEntityIDs[1]] = append(codeNameStructMap[namespaceEntityIDs[1]], &totpCode{
		name:          name,
		methodID:      namespaceMethodIDs[1],
		namespacePath: namespacePaths[1],
		entityID:      namespaceEntityIDs[1],
	})

	return codeNameStructMap
}

func registerTotpUrl(client *api.Client, totpURL, codeName string, t *testing.T) {
	client.SetNamespace("")
	_, err := client.Logical().Write(fmt.Sprintf("totp/keys/%s", codeName), map[string]interface{}{
		"url": totpURL,
	})
	if err != nil {
		t.Fatalf("failed to register a TOTP URL: %v", err)
	}
}

func checkDestroyTotp(client *api.Client, namespacePaths, namespaceMethodIDs, namespaceEntityIDs []string, t *testing.T) {
	client.SetNamespace("")

	_, err := client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-destroy"), map[string]interface{}{
		"entity_id": namespaceEntityIDs[0],
		"method_id": namespaceMethodIDs[1],
	})
	if err == nil {
		t.Fatalf("failed to destroy the MFA secret: %s", err)
	}

	// non-root namespace
	client.SetNamespace(namespacePaths[2])
	_, err = client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-destroy"), map[string]interface{}{
		"entity_id": namespaceEntityIDs[2],
		"method_id": namespaceMethodIDs[2],
	})
	if err != nil {
		t.Fatalf("failed to destroy a TOTP secret on an entity: %v", err)
	}

	client.SetNamespace(namespacePaths[0])
	_, err = client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-destroy"), map[string]interface{}{
		"entity_id": namespaceEntityIDs[0],
		"method_id": namespaceMethodIDs[0],
	})
	if err != nil {
		t.Fatalf("failed to destroy a TOTP secret on an entity: %v", err)
	}

	client.SetNamespace(namespacePaths[1])
	_, err = client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-destroy"), map[string]interface{}{
		"entity_id": namespaceEntityIDs[1],
		"method_id": namespaceMethodIDs[0],
	})
	if err != nil {
		t.Fatalf("failed to destroy a TOTP secret on an entity: %v", err)
	}

	client.SetNamespace(namespacePaths[1])
	_, err = client.Logical().Write(fmt.Sprintf("identity/mfa/method-id/totp/admin-destroy"), map[string]interface{}{
		"entity_id": namespaceEntityIDs[1],
		"method_id": namespaceMethodIDs[1],
	})
	if err != nil {
		t.Fatalf("failed to destroy a TOTP secret on an entity: %v", err)
	}
}
