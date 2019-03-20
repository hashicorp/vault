// +build !enterprise

package vault

import (
	"github.com/hashicorp/vault/helper/namespace"
)

func (ts *TokenStore) baseView(ns *namespace.Namespace) *BarrierView {
	return ts.baseBarrierView
}

func (ts *TokenStore) idView(ns *namespace.Namespace) *BarrierView {
	return ts.idBarrierView
}

func (ts *TokenStore) accessorView(ns *namespace.Namespace) *BarrierView {
	return ts.accessorBarrierView
}

func (ts *TokenStore) parentView(ns *namespace.Namespace) *BarrierView {
	return ts.parentBarrierView
}

func (ts *TokenStore) rolesView(ns *namespace.Namespace) *BarrierView {
	return ts.rolesBarrierView
}
