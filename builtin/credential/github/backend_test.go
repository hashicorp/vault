package github

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_basic(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccMap(t, "default", "foo"),
			testAccMap(t, "oWnErs", "bar"),
			testAccLogin(t, []string{"bar", "foo"}),
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GITHUB_TOKEN"); v == "" {
		t.Fatal("GITHUB_USER must be set for acceptance tests")
	}

	if v := os.Getenv("GITHUB_ORG"); v == "" {
		t.Fatal("GITHUB_ORG must be set for acceptance tests")
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"organization": os.Getenv("GITHUB_ORG"),
		},
	}
}

func testAccMap(t *testing.T, k string, v string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "map/teams/" + k,
		Data: map[string]interface{}{
			"value": v,
		},
	}
}

func testAccLogin(t *testing.T, keys []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"token": os.Getenv("GITHUB_TOKEN"),
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth(keys),
	}
}
