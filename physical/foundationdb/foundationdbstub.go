//go:build !foundationdb

package foundationdb

import (
	"errors"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/sdk/physical"
)

func NewFDBBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	return nil, errors.New("FoundationDB backend not available in this Vault build")
}
