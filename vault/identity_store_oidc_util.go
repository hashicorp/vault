// +build !enterprise

package vault

import (
	"github.com/hashicorp/vault/helper/namespace"
)

func (i *IdentityStore) listNamespaces() []*namespace.Namespace {
	return []*namespace.Namespace{namespace.RootNamespace}
}
