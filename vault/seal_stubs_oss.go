//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/physical"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

// isSealOldKeyError returns true if a value was decrypted using the
// old "unwrapSeal".
func isSealOldKeyError(err error) bool {
	return false
}

func startPartialSealRewrapping(c *Core) {
	// nothing to do
}

func GetPartiallySealWrappedPaths(ctx context.Context, backend physical.Backend) ([]string, error) {
	return nil, nil
}
