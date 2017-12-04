package approle

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func createBackendWithStorage(t *testing.T) (*backend, logical.Storage) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	if b == nil {
		t.Fatalf("failed to create backend")
	}
	err = b.Backend.Setup(config)
	if err != nil {
		t.Fatal(err)
	}
	return b, config.StorageView
}
