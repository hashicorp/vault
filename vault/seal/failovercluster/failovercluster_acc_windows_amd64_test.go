package failovercluster

import (
	"context"
	"os"
	"reflect"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func TestFailoverClusterVault_SetConfig(t *testing.T) {
	if os.Getenv("FAILOVERCLUSTER_RESOURCE_NAME") == "" {
		t.SkipNow()
	}

	seal := NewSeal(logging.NewVaultLogger(log.Trace))

	resourceName := os.Getenv("FAILOVERCLUSTER_RESOURCE_NAME")
	os.Unsetenv("FAILOVERCLUSTER_RESOURCE_NAME")

	// Attempt to set config, expect failure due to missing config
	_, err := seal.SetConfig(nil)
	if err == nil {
		t.Fatal("expected error when FailoverCluster config values are not provided")
	}

	os.Setenv("FAILOVERCLUSTER_RESOURCE_NAME", resourceName)

	_, err = seal.SetConfig(nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFailoverClusterVault_Lifecycle(t *testing.T) {
	if os.Getenv("FAILOVERCLUSTER_RESOURCE_NAME") == "" {
		t.SkipNow()
	}

	s := NewSeal(logging.NewVaultLogger(log.Trace))
	_, err := s.SetConfig(nil)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
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
