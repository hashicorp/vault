package plugin

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestStorage_impl(t *testing.T) {
	var _ logical.Storage = new(StorageClient)
}
