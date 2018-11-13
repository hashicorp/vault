package alicloudkms

import (
	"context"
	"os"
	"reflect"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
)

// This test executes real calls. The calls themselves should be free,
// but the KMS key used is generally not free. Alibaba doesn't publish
// the price but it can be assumed to be around $1/month because that's
// what AWS charges for the same.
//
// To run this test, the following env variables need to be set:
//   - VAULT_ALICLOUDKMS_SEAL_KEY_ID
//   - ALICLOUD_REGION
//   - ALICLOUD_ACCESS_KEY
//   - ALICLOUD_SECRET_KEY
func TestAccAliCloudKMSSeal_Lifecycle(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}

	s := NewSeal(logging.NewVaultLogger(log.Trace))
	_, err := s.SetConfig(nil)
	if err != nil {
		t.Fatalf("err : %s", err)
	}

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
