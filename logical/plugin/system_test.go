package plugin

import (
	"testing"

	"github.com/hashicorp/vault/logical"
)

func Test_impl(t *testing.T) {
	var _ logical.SystemView = new(SystemViewClient)
}
