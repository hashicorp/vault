package google

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

const GOOGLE_CODE_ENV_VAR_NAME = "GOOGLE_CODE"
const GOOGLE_DOMAIN_ENV_VAR_NAME = "GOOGLE_DOMAIN"

func environmentVariable(name string) string {

	return os.Getenv(name)
}

func googleCode() string {
	return environmentVariable(GOOGLE_CODE_ENV_VAR_NAME)
}

func googleDomain() string {
	return environmentVariable(GOOGLE_DOMAIN_ENV_VAR_NAME)
}

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
		"code": googleCode(),
	}
	config_data1 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "",
		"max_ttl":      "",
	}
	expectedTTL1, _ := time.ParseDuration("24h0m0s")
	config_data2 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "1h",
		"max_ttl":      "2h",
	}
	expectedTTL2, _ := time.ParseDuration("1h0m0s")
	config_data3 := map[string]interface{}{
		"domain": googleDomain(),
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
	maxLeaseTTLVal := time.Hour * 24 * 30
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
			testAccStepConfig(t),
			testAccMap(t, "default", "root"),
			testAccMap(t, "oWnErs", "root"),
			testAccLogin(t, []string{"root"}),
		},
	})
}

func testAccPreCheck(t *testing.T) {
	requiredEnvVars := []string{ GOOGLE_CODE_ENV_VAR_NAME, GOOGLE_DOMAIN_ENV_VAR_NAME }
	for _, envVar := range requiredEnvVars {
		if value := environmentVariable(envVar); value == "" {
			t.Fatal(fmt.Sprintf("missing environment variable %s", envVar))
		}
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"domain": googleDomain(),
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

func testAccLogin(t *testing.T, keys []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"code": googleCode(),
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth(keys),
	}
}
