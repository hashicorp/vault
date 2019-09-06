// Copyright © 2019, Oracle and/or its affiliates.
package ocikms

import (
	"os"
	"reflect"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"golang.org/x/net/context"
)

/*
* To run these tests, ensure you setup:
* 1. OCI SDK with your credentials. Refer to here:
*		https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdkconfig.htm
* 2. Go to ocikms folder: vault/vault/seal/ocikms
*		VAULT_OCIKMS_SEAL_KEY_ID="your-kms-key" VAULT_OCIKMS_CRYPTO_ENDPOINT="your-kms-crypto-endpoint" go test
 */

func TestOCIKMSSeal(t *testing.T) {
	initSeal(t)
}

func TestOCIKMSSeal_LifeCycle(t *testing.T) {
	s := initSeal(t)

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

func initSeal(t *testing.T) *OCIKMSSeal {
	// Skip tests if we are not running acceptance tests
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	s := NewSeal(logging.NewVaultLogger(log.Trace))
	_, err := s.SetConfig(nil)
	if err == nil {
		t.Fatal("expected error when OCIKMSSeal required values are not provided")
	}

	mockConfig := map[string]string{

		"auth_type_api_key": "true",
	}

	_, err = s.SetConfig(mockConfig)
	if err != nil {
		t.Fatalf("error setting seal config: %v", err)
	}

	return s
}
