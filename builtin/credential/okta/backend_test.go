// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package okta

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	stepwise "github.com/hashicorp/vault-testing-stepwise"
	dockerEnvironment "github.com/hashicorp/vault-testing-stepwise/environments/docker"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	thstepwise "github.com/hashicorp/vault/sdk/helper/testhelpers/stepwise"
	thutils "github.com/hashicorp/vault/sdk/helper/testhelpers/utils"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/stretchr/testify/require"
)

// To run this test, set the following env variables:
// VAULT_ACC=1
// OKTA_ORG=dev-219337
// OKTA_PREVIEW=1 to use the okta preview URL (optional)
// OKTA_API_TOKEN=<generate via web UI, see Confluence for login details>
// OKTA_USERNAME=test3@example.com
// OKTA_PASSWORD=<find in 1password>
//
// Using free integrater accounts to test incurs rate limits; so you will need
// to set the following to skip creating/deleting test groups
// SKIP_GROUPS_CREATE_DELETE=<non-empty to skip> (optional)
//
// You will need to install the Okta client app on your mobile device and
// setup MFA in order to use the Okta web UI.  This test does not exercise
// MFA however (which is an enterprise feature), and therefore the test
// user in OKTA_USERNAME should not be configured with it.  Currently
// test3@example.com is not a member of testgroup, which is the group with
// the profile that requires MFA.

const (
	envOktaOrg                = "OKTA_ORG"
	envOktaUsername           = "OKTA_USERNAME"
	envOktaPassword           = "OKTA_PASSWORD"
	envOktaAPIToken           = "OKTA_API_TOKEN"
	envOktaPreview            = "OKTA_PREVIEW"
	envSkipGroupsCreateDelete = "SKIP_GROUPS_CREATE_DELETE"
)

func TestBackend_Config(t *testing.T) {
	// Ensure required environment variables are set.
	requiredEnvs := []string{
		"VAULT_ACC",
		envOktaOrg,
		envOktaUsername,
		envOktaPassword,
		envOktaAPIToken,
	}
	thutils.SkipUnlessEnvVarsSet(t, requiredEnvs)

	defaultLeaseTTL := time.Hour * 12
	maxLeaseTTL := time.Hour * 24

	username := os.Getenv(envOktaUsername)
	password := os.Getenv(envOktaPassword)
	token := os.Getenv(envOktaAPIToken)

	if os.Getenv(envSkipGroupsCreateDelete) != "" {
		t.Logf("Skipping creation/deletion of Okta test groups")
	} else {
		groupIDs := createOktaGroups(t, username, token, os.Getenv(envOktaOrg))
		defer deleteOktaGroups(t, token, os.Getenv(envOktaOrg), groupIDs)
	}

	configData := map[string]interface{}{
		"org_name": os.Getenv(envOktaOrg),
		"base_url": defaultBaseURL,
	}

	if os.Getenv(envOktaPreview) != "" {
		configData["base_url"] = previewBaseURL
	}

	updatedDuration := time.Hour * 1
	configDataToken := map[string]interface{}{
		"api_token": token,
		"token_ttl": "1h",
	}

	envOptions := &stepwise.MountOptions{
		RegistryName:    "okta-auth",
		PluginType:      api.PluginTypeCredential,
		PluginName:      "okta",
		MountPathPrefix: "okta-auth",
		MountConfigInput: api.MountConfigInput{
			DefaultLeaseTTL: fmt.Sprintf("%d", int(defaultLeaseTTL.Seconds())),
			MaxLeaseTTL:     fmt.Sprintf("%d", int(maxLeaseTTL.Seconds())),
		},
	}
	stepwise.Run(t, stepwise.Case{
		Precheck:    func() { testAccPreCheck(t) },
		Environment: dockerEnvironment.NewEnvironment("okta", envOptions),
		Steps: []stepwise.Step{
			testConfigCreate(configData),
			// 2. Login with bad password, expect failure (E0000004=okta auth failure).
			testLoginWrite(username, "wrong", "E0000004", 0, nil),
			// 3. Make our user belong to two groups and have one user-specific policy.
			testAccUserGroups(username, "local_grouP,lOcal_group2", []string{"user_policy"}),
			// 4. Create the group local_group, assign it a single policy.
			testAccGroups("local_groUp", "loCal_group_policy"),
			// 5. Login with good password, expect user to have their user-specific
			// policy and the policy of the one valid group they belong to.
			testLoginWrite(username, password, "", defaultLeaseTTL, []string{"local_group_policy", "user_policy"}),
			// 6. Create the group everyone, assign it two policies.  This is a
			// magic group name in okta that always exists and which every
			// user automatically belongs to.
			testAccGroups("everyoNe", "everyone_grouP_policy,eveRy_group_policy2"),
			// 7. Login as before, expect same result
			testLoginWrite(username, password, "", defaultLeaseTTL, []string{"local_group_policy", "user_policy"}),
			// 8. Add API token so we can lookup groups
			testConfigUpdate(configDataToken),
			testConfigRead(token, configData),
			// 10. Login should now lookup okta groups; since all okta users are
			// in the "everyone" group, that should be returned; since we
			// defined policies attached to the everyone group, we should now
			// see those policies attached to returned vault token.
			testLoginWrite(username, password, "", updatedDuration,
				[]string{"everyone_group_policy", "every_group_policy2", "local_group_policy", "user_policy"}),
			testAccGroups("locAl_group2", "testgroup_group_policy"),
			testLoginWrite(username, password, "", updatedDuration,
				[]string{"everyone_group_policy", "every_group_policy2", "local_group_policy", "testgroup_group_policy", "user_policy"}),
			testAccLogin(username, password,
				[]string{"default", "everyone_group_policy", "every_group_policy2", "local_group_policy", "testgroup_group_policy", "user_policy"}),
		},
	})
}

func createOktaGroups(t *testing.T, username string, token string, org string) []string {
	t.Helper()
	orgURL := "https://" + org + "."
	if os.Getenv(envOktaPreview) != "" {
		orgURL += previewBaseURL
	} else {
		orgURL += defaultBaseURL
	}

	cfg, err := okta.NewConfiguration(okta.WithOrgUrl(orgURL), okta.WithToken(token))
	require.Nil(t, err)
	client := okta.NewAPIClient(cfg)
	ctx := context.Background()

	users, _, err := client.UserAPI.ListUsers(ctx).Q(username).Execute()
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
		groups, _, err := client.GroupAPI.ListGroups(ctx).Q(name).Execute()
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
	t.Helper()
	orgURL := "https://" + org + "." + previewBaseURL
	cfg, err := okta.NewConfiguration(okta.WithOrgUrl(orgURL), okta.WithToken(token))
	require.Nil(t, err)
	client := okta.NewAPIClient(cfg)

	for _, groupID := range groupIDs {
		_, err := client.GroupAPI.DeleteGroup(context.Background(), groupID).Execute()
		require.Nil(t, err)
	}
}

func testLoginWrite(username, password, reason string, expectedTTL time.Duration, policies []string) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "login/" + username,
		Data: map[string]interface{}{
			"password": password,
		},
		Assert: func(resp *api.Secret, err error) error {
			if reason != "" {
				if !strings.Contains(err.Error(), reason) {
					return fmt.Errorf("expected error containing %q, got no error", reason)
				}
				return nil
			}

			if err != nil {
				return err
			}

			if resp == nil {
				return fmt.Errorf("expected non-nil response")
			}

			if resp.Auth != nil {
				if !policyutil.EquivalentPolicies(resp.Auth.Policies, policies) {
					return fmt.Errorf("policy mismatch expected %v but got %v", policies, resp.Auth.Policies)
				}

				actualTTL := resp.Auth.LeaseDuration
				if time.Duration(actualTTL)*time.Second != expectedTTL {
					return fmt.Errorf("TTL mismatch expected %v but got %v", expectedTTL, actualTTL)
				}
			}

			return nil
		},
	}
}

func testConfigCreate(d map[string]interface{}) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.WriteOperation,
		Path:      "config",
		Data:      d,
	}
}

func testConfigUpdate(d map[string]interface{}) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "config",
		Data:      d,
	}
}

func testConfigRead(token string, d map[string]interface{}) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.ReadOperation,
		Path:      "config",
		Assert: func(resp *api.Secret, err error) error {
			if err != nil {
				return err
			}

			if resp.Data["error"] != nil {
				return fmt.Errorf("error reading config: %v", resp.Data["error"])
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
	t.Helper()
	if v := os.Getenv(envOktaUsername); v == "" {
		t.Fatalf("%s must be set for acceptance tests", envOktaUsername)
	}

	if v := os.Getenv(envOktaPassword); v == "" {
		t.Fatalf("%s must be set for acceptance tests", envOktaPassword)
	}

	if v := os.Getenv(envOktaOrg); v == "" {
		t.Fatalf("%s must be set for acceptance tests", envOktaOrg)
	}

	if v := os.Getenv(envOktaAPIToken); v == "" {
		t.Fatalf("%s must be set for acceptance tests", envOktaAPIToken)
	}

	if v := os.Getenv(envSkipGroupsCreateDelete); v == "" {
		t.Fatalf("%s must be set for acceptance tests", envSkipGroupsCreateDelete)
	}
}

func testAccUserGroups(user string, groups interface{}, policies interface{}) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "users/" + user,
		Data: map[string]interface{}{
			"groups":   groups,
			"policies": policies,
		},
	}
}

func testAccGroups(group string, policies interface{}) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "groups/" + group,
		Data: map[string]interface{}{
			"policies": policies,
		},
	}
}

func testAccLogin(user, password string, keys []string) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "login/" + user,
		Data: map[string]interface{}{
			"password": password,
		},
		Unauthenticated: true,
		Assert:          thstepwise.NewAssertAuthPoliciesFunc(keys),
	}
}
