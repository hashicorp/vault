// +build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/helper/namespace"
)

func (i *IdentityStore) listNamespacePaths(ctx context.Context) []string {
	return []string{namespace.RootNamespace.Path}
}
