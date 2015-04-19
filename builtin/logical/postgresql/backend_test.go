package postgresql

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
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PG_URL"); v == "" {
		t.Fatal("PG_URL must be set for acceptance tests")
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "config/connection",
		Data: map[string]interface{}{
			"value": os.Getenv("PG_URL"),
		},
	}
}
