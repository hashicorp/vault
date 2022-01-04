//go:build !enterprise

package pki

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/certutil"
)

func fetchManagedKey(_ context.Context, _ *backend, _ *certutil.ParsedCertBundle) error {
	// No-op
	return nil
}
