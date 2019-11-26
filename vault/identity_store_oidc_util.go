// +build !enterprise

package vault

import (
	"github.com/hashicorp/vault/helper/namespace"
)

func (i *IdentityStore) listNamespacePaths() []string {
	return []string{namespace.RootNamespace.Path}
}
