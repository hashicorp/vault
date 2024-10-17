// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package okta

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/stretchr/testify/require"
)

// To run this test, set the following env variables:
// VAULT_ACC=1
// OKTA_ORG=dev-219337
// OKTA_API_TOKEN=<generate via web UI, see Confluence for login details>
// OKTA_USERNAME=test3@example.com
// OKTA_PASSWORD=<find in 1password>
//
// You will need to install the Okta client app on your mobile device and
// setup MFA in order to use the Okta web UI.  This test does not exercise
// MFA however (which is an enterprise feature), and therefore the test
// user in OKTA_USERNAME should not be configured with it.  Currently
// test3@example.com is not a member of testgroup, which is the group with
// the profile that requires MFA.
func TestBackend_Config(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}

	// Ensure each cred is populated.
	credNames := []string{
		"OKTA_USERNAME",
		"OKTA_PASSWORD",
		"OKTA_API_TOKEN",
	}
	testhelpers.SkipUnlessEnvVarsSet(t, credNames)

	defaultLeaseTTLVal := time.Hour * 12
	maxLeaseTTLVal := time.Hour * 24
	b, err := Factory(context.Background(), &logical.BackendConfig{
		Logger: logging.NewVaultLogger(log.Trace),
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	username := os.Getenv("OKTA_USERNAME")
	password := os.Getenv("OKTA_PASSWORD")
	token := os.Getenv("OKTA_API_TOKEN")
	groupIDs := createOktaGroups(t, username, token, os.Getenv("OKTA_ORG"))
	defer deleteOktaGroups(t, token, os.Getenv("OKTA_ORG"), groupIDs)

	configData := map[string]interface{}{
		"org_name": os.Getenv("OKTA_ORG"),
		"base_url": "oktapreview.com",
	}

	updatedDuration := time.Hour * 1
	configDataToken := map[string]interface{}{
		"api_token": token,
		"token_ttl": "1h",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest:    true,
		PreCheck:          func() { testAccPreCheck(t) },
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			testConfigCreate(t, configData),
			// 2. Login with bad password, expect failure (E0000004=okta auth failure).
			testLoginWrite(t, username, "wrong", "E0000004", 0, nil),
			// 3. Make our user belong to two groups and have one user-specific policy.
			testAccUserGroups(t, username, "local_grouP,lOcal_group2", []string{"user_policy"}),
			// 4. Create the group local_group, assign it a single policy.
			testAccGroups(t, "local_groUp", "loCal_group_policy"),
			// 5. Login with good password, expect user to have their user-specific
			// policy and the policy of the one valid group they belong to.
			testLoginWrite(t, username, password, "", defaultLeaseTTLVal, []string{"local_group_policy", "user_policy"}),
			// 6. Create the group everyone, assign it two policies.  This is a
			// magic group name in okta that always exists and which every
			// user automatically belongs to.
			testAccGroups(t, "everyoNe", "everyone_grouP_policy,eveRy_group_policy2"),
			// 7. Login as before, expect same result
			testLoginWrite(t, username, password, "", defaultLeaseTTLVal, []string{"local_group_policy", "user_policy"}),
			// 8. Add API token so we can lookup groups
			testConfigUpdate(t, configDataToken),
			testConfigRead(t, token, configData),
			// 10. Login should now lookup okta groups; since all okta users are
			// in the "everyone" group, that should be returned; since we
			// defined policies attached to the everyone group, we should now
			// see those policies attached to returned vault token.
			testLoginWrite(t, username, password, "", updatedDuration, []string{"everyone_group_policy", "every_group_policy2", "local_group_policy", "user_policy"}),
			testAccGroups(t, "locAl_group2", "testgroup_group_policy"),
			testLoginWrite(t, username, password, "", updatedDuration, []string{"everyone_group_policy", "every_group_policy2", "local_group_policy", "testgroup_group_policy", "user_policy"}),
		},
	})
}

func createOktaGroups(t *testing.T, username string, token string, org string) []string {
	orgURL := "https://" + org + "." + previewBaseURL
	cfg, err := okta.NewConfiguration(okta.WithOrgUrl(orgURL), okta.WithToken(token))
	require.Nil(t, err)
	client := okta.NewAPIClient(cfg)
	ctx := context.Background()

	users, _, err := client.UserAPI.ListUsers(ctx).Filter(username).Execute()
	require.Nil(t, err)
	require.Len(t, users, 1)
	userID := users[0].GetId()
	var groupIDs []string

	// Verify that login's call to list the groups of the user logging in will page
	// through multiple result sets; note here
	// https://developer.okta.com/docs/reference/api/groups/#list-groups-with-defaults
	// that "If you don't specify a value for limit and don't specify a query,
	// only 200 results are returned for most orgs."
	for i := 0; i < 201; i++ {
		name := fmt.Sprintf("TestGroup%d", i)
		groups, _, err := client.GroupAPI.ListGroups(ctx).Filter(name).Execute()
		require.Nil(t, err)

		var groupID string
		if len(groups) == 0 {
			group, _, err := client.GroupAPI.CreateGroup(ctx).Group(okta.Group{
				Profile: &okta.GroupProfile{
					Name: okta.PtrString(fmt.Sprintf("TestGroup%d", i)),
				},
			}).Execute()
			require.Nil(t, err)
			groupID = group.GetId()
		} else {
			groupID = groups[0].GetId()
		}
		groupIDs = append(groupIDs, groupID)

		_, err = client.GroupAPI.AssignUserToGroup(ctx, groupID, userID).Execute()
		require.Nil(t, err)
	}
	return groupIDs
}

func deleteOktaGroups(t *testing.T, token string, org string, groupIDs []string) {
	orgURL := "https://" + org + "." + previewBaseURL
	cfg, err := okta.NewConfiguration(okta.WithOrgUrl(orgURL), okta.WithToken(token))
	require.Nil(t, err)
	client := okta.NewAPIClient(cfg)

	for _, groupID := range groupIDs {
		_, err := client.GroupAPI.DeleteGroup(context.Background(), groupID).Execute()
		require.Nil(t, err)
	}
}

func testLoginWrite(t *testing.T, username, password, reason string, expectedTTL time.Duration, policies []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + username,
		ErrorOk:   true,
		Data: map[string]interface{}{
			"password": password,
		},
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				if reason == "" || !strings.Contains(resp.Error().Error(), reason) {
					return resp.Error()
				}
			} else if reason != "" {
				return fmt.Errorf("expected error containing %q, got no error", reason)
			}

			if resp.Auth != nil {
				if !policyutil.EquivalentPolicies(resp.Auth.Policies, policies) {
					return fmt.Errorf("policy mismatch expected %v but got %v", policies, resp.Auth.Policies)
				}

				actualTTL := resp.Auth.LeaseOptions.TTL
				if actualTTL != expectedTTL {
					return fmt.Errorf("TTL mismatch expected %v but got %v", expectedTTL, actualTTL)
				}
			}

			return nil
		},
	}
}

func testConfigCreate(t *testing.T, d map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.CreateOperation,
		Path:      "config",
		Data:      d,
	}
}

func testConfigUpdate(t *testing.T, d map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      d,
	}
}

func testConfigRead(t *testing.T, token string, d map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "config",
		Check: func(resp *logical.Response) error {
			if resp.IsError() {
				return resp.Error()
			}

			if resp.Data["org_name"] != d["org_name"] {
				return fmt.Errorf("org mismatch expected %s but got %s", d["organization"], resp.Data["Org"])
			}

			if resp.Data["base_url"] != d["base_url"] {
				return fmt.Errorf("BaseURL mismatch expected %s but got %s", d["base_url"], resp.Data["BaseURL"])
			}

			for _, value := range resp.Data {
				if value == token {
					return fmt.Errorf("token should not be returned on a read request")
				}
			}

			return nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OKTA_USERNAME"); v == "" {
		t.Fatal("OKTA_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("OKTA_PASSWORD"); v == "" {
		t.Fatal("OKTA_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("OKTA_ORG"); v == "" {
		t.Fatal("OKTA_ORG must be set for acceptance tests")
	}

	if v := os.Getenv("OKTA_API_TOKEN"); v == "" {
		t.Fatal("OKTA_API_TOKEN must be set for acceptance tests")
	}
}

func testAccUserGroups(t *testing.T, user string, groups interface{}, policies interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "users/" + user,
		Data: map[string]interface{}{
			"groups":   groups,
			"policies": policies,
		},
	}
}

func testAccGroups(t *testing.T, group string, policies interface{}) logicaltest.TestStep {
	t.Logf("[testAccGroups] - Registering group %s, policy %s", group, policies)
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "groups/" + group,
		Data: map[string]interface{}{
			"policies": policies,
		},
	}
}

func testAccLogin(t *testing.T, user, password string, keys []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login/" + user,
		Data: map[string]interface{}{
			"password": password,
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth(keys),
	}
}
