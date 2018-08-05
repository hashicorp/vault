// +build !foundationdb

package foundationdb

import (
	"fmt"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/physical"
)

func NewFDBBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	return nil, fmt.Errorf("FoundationDB backend not available in this Vault build")
}
