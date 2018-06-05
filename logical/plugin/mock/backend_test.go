package mock

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestBackend_impl(t *testing.T) {
	var _ logical.Backend = new(backend)
}
