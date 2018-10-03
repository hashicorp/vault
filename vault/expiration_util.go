// +build !enterprise

package vault

import (
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
)

func (m *ExpirationManager) leaseView(*namespace.Namespace) *BarrierView {
	return m.idView
}

func (m *ExpirationManager) tokenIndexView(*namespace.Namespace) *BarrierView {
	return m.tokenView
}

func (m *ExpirationManager) collectLeases() (map[*namespace.Namespace][]string, int, error) {
	leaseCount := 0
	existing := make(map[*namespace.Namespace][]string)
	keys, err := logical.CollectKeys(m.quitContext, m.leaseView(namespace.RootNamespace))
	if err != nil {
		return nil, 0, errwrap.Wrapf("failed to scan for leases: {{err}}", err)
	}
	existing[namespace.RootNamespace] = keys
	leaseCount += len(keys)
	return existing, leaseCount, nil
}
