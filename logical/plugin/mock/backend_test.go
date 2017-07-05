package mock

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestMockBackend_impl(t *testing.T) {
	var _ logical.Backend = new(backend)
}
