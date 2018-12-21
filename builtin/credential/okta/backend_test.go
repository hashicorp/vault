package okta

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/policyutil"

	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_Config(t *testing.T) {
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

	configData := map[string]interface{}{
		"organization": os.Getenv("OKTA_ORG"),
		"base_url":     "oktapreview.com",
	}

	updatedDuration := time.Hour * 1
	configDataToken := map[string]interface{}{
		"token": token,
		"ttl":   "1h",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testConfigCreate(t, configData),
			testLoginWrite(t, username, "wrong", "E0000004", 0, nil),
			testLoginWrite(t, username, password, "user is not a member of any authorized policy", 0, nil),
			testAccUserGroups(t, username, "local_grouP,lOcal_group2", []string{"user_policy"}),
			testAccGroups(t, "local_groUp", "loCal_group_policy"),
			testLoginWrite(t, username, password, "", defaultLeaseTTLVal, []string{"local_group_policy", "user_policy"}),
			testAccGroups(t, "everyoNe", "everyone_grouP_policy,eveRy_group_policy2"),
			testLoginWrite(t, username, password, "", defaultLeaseTTLVal, []string{"local_group_policy", "user_policy"}),
			testConfigUpdate(t, configDataToken),
			testConfigRead(t, token, configData),
			testLoginWrite(t, username, password, "", updatedDuration, []string{"everyone_group_policy", "every_group_policy2", "local_group_policy", "user_policy"}),
			testAccGroups(t, "locAl_group2", "testgroup_group_policy"),
			testLoginWrite(t, username, password, "", updatedDuration, []string{"everyone_group_policy", "every_group_policy2", "local_group_policy", "testgroup_group_policy", "user_policy"}),
		},
	})
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

			if resp.Data["organization"] != d["organization"] {
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
