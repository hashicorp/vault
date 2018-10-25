package gcpckms

import (
	"os"
	"reflect"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	context "golang.org/x/net/context"
)

const (
	// These values need to match the values from the hc-value-testing project
	gcpckmsTestProjectID  = "hc-vault-testing"
	gcpckmsTestLocationID = "global"
	gcpckmsTestKeyRing    = "vault-test-keyring"
	gcpckmsTestCryptoKey  = "vault-test-key"
)

func TestGCPCKMSSeal(t *testing.T) {
	// Do an error check before env vars are set
	s := NewSeal(logging.NewVaultLogger(log.Trace))
	_, err := s.SetConfig(nil)
	if err == nil {
		t.Fatal("expected error when GCPCKMSSeal required values are not provided")
	}

	// Now test for cases where CKMS values are provided
	checkAndSetEnvVars(t)

	configCases := map[string]map[string]string{
		"env_var": nil,
		"config": map[string]string{
			"credentials": os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		},
	}

	for name, config := range configCases {
		t.Run(name, func(t *testing.T) {
			s := NewSeal(logging.NewVaultLogger(log.Trace))
			_, err := s.SetConfig(config)
			if err != nil {
				t.Fatalf("error setting seal config: %v", err)
			}
		})
	}
}

func TestGCPCKMSSeal_Lifecycle(t *testing.T) {
	checkAndSetEnvVars(t)

	s := NewSeal(logging.NewVaultLogger(log.Trace))
	_, err := s.SetConfig(nil)
	if err != nil {
		t.Fatalf("error setting seal config: %v", err)
	}

	// Test Encrypt and Decrypt calls
	input := []byte("foo")
	swi, err := s.Encrypt(context.Background(), input)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	pt, err := s.Decrypt(context.Background(), swi)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	if !reflect.DeepEqual(input, pt) {
		t.Fatalf("expected %s, got %s", input, pt)
	}
}

// checkAndSetEnvVars check and sets the required env vars. It will skip tests that are
// not ran as acceptance tests since they require calling to external APIs.
func checkAndSetEnvVars(t *testing.T) {
	t.Helper()

	// Skip tests if we are not running acceptance tests
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" && os.Getenv(EnvGCPCKMSSealCredsPath) == "" {
		t.Fatal("unable to get GCP credentials via environment variables")
	}

	if os.Getenv(EnvGCPCKMSSealProject) == "" {
		os.Setenv(EnvGCPCKMSSealProject, gcpckmsTestProjectID)
	}

	if os.Getenv(EnvGCPCKMSSealLocation) == "" {
		os.Setenv(EnvGCPCKMSSealLocation, gcpckmsTestLocationID)
	}

	if os.Getenv(EnvGCPCKMSSealKeyRing) == "" {
		os.Setenv(EnvGCPCKMSSealKeyRing, gcpckmsTestKeyRing)
	}

	if os.Getenv(EnvGCPCKMSSealCryptoKey) == "" {
		os.Setenv(EnvGCPCKMSSealCryptoKey, gcpckmsTestCryptoKey)
	}
}
