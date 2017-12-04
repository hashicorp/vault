package github

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_Config(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 2
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	login_data := map[string]interface{}{
		// This token has to be replaced with a working token for the test to work.
		"token": os.Getenv("GITHUB_TOKEN"),
	}
	config_data1 := map[string]interface{}{
		"organization": os.Getenv("GITHUB_ORG"),
		"ttl":          "",
		"max_ttl":      "",
	}
	expectedTTL1, _ := time.ParseDuration("24h0m0s")
	config_data2 := map[string]interface{}{
		"organization": os.Getenv("GITHUB_ORG"),
		"ttl":          "1h",
		"max_ttl":      "2h",
	}
	expectedTTL2, _ := time.ParseDuration("1h0m0s")
	config_data3 := map[string]interface{}{
		"organization": os.Getenv("GITHUB_ORG"),
		"ttl":          "50h",
		"max_ttl":      "50h",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, config_data1),
			testLoginWrite(t, login_data, expectedTTL1.Nanoseconds(), false),
			testConfigWrite(t, config_data2),
			testLoginWrite(t, login_data, expectedTTL2.Nanoseconds(), false),
			testConfigWrite(t, config_data3),
			testLoginWrite(t, login_data, 0, true),
		},
	})
}

func testLoginWrite(t *testing.T, d map[string]interface{}, expectedTTL int64, expectFail bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		ErrorOk:   true,
		Data:      d,
		Check: func(resp *logical.Response) error {
			if resp.IsError() && expectFail {
				return nil
			}
			var actualTTL int64
			actualTTL = resp.Auth.LeaseOptions.TTL.Nanoseconds()
			if actualTTL != expectedTTL {
				return fmt.Errorf("TTL mismatched. Expected: %d Actual: %d", expectedTTL, resp.Auth.LeaseOptions.TTL.Nanoseconds())
			}
			return nil
		},
	}
}

func testConfigWrite(t *testing.T, d map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      d,
	}
}

func TestBackend_basic(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, false),
			testAccMap(t, "default", "fakepol"),
			testAccMap(t, "oWnErs", "fakepol"),
			testAccLogin(t, []string{"default", "fakepol"}),
			testAccStepConfig(t, true),
			testAccMap(t, "default", "fakepol"),
			testAccMap(t, "oWnErs", "fakepol"),
			testAccLogin(t, []string{"default", "fakepol"}),
			testAccStepConfigWithBaseURL(t),
			testAccMap(t, "default", "fakepol"),
			testAccMap(t, "oWnErs", "fakepol"),
			testAccLogin(t, []string{"default", "fakepol"}),
			testAccMap(t, "default", "fakepol"),
			testAccStepConfig(t, true),
			mapUserToPolicy(t, os.Getenv("GITHUB_USER"), "userpolicy"),
			testAccLogin(t, []string{"default", "fakepol", "userpolicy"}),
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GITHUB_TOKEN"); v == "" {
		t.Skip("GITHUB_TOKEN must be set for acceptance tests")
	}

	if v := os.Getenv("GITHUB_ORG"); v == "" {
		t.Skip("GITHUB_ORG must be set for acceptance tests")
	}

	if v := os.Getenv("GITHUB_BASEURL"); v == "" {
		t.Skip("GITHUB_BASEURL must be set for acceptance tests (use 'https://api.github.com' if you don't know what you're doing)")
	}
}

func testAccStepConfig(t *testing.T, upper bool) logicaltest.TestStep {
	ts := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"organization": os.Getenv("GITHUB_ORG"),
		},
	}
	if upper {
		ts.Data["organization"] = strings.ToUpper(os.Getenv("GITHUB_ORG"))
	}
	return ts
}

func testAccStepConfigWithBaseURL(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"organization": os.Getenv("GITHUB_ORG"),
			"base_url":     os.Getenv("GITHUB_BASEURL"),
		},
	}
}

func testAccMap(t *testing.T, k string, v string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "map/teams/" + k,
		Data: map[string]interface{}{
			"value": v,
		},
	}
}

func mapUserToPolicy(t *testing.T, k string, v string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "map/users/" + k,
		Data: map[string]interface{}{
			"value": v,
		},
	}
}

func testAccLogin(t *testing.T, policies []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"token": os.Getenv("GITHUB_TOKEN"),
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth(policies),
	}
}
