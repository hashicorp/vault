// +build !enterprise

package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/namespace"
)

// getBarrierView is a convenience method to return a BarrierView scoped to the namespace and the key type.
func (r *ManagedKeyRegistry) getBarrierView(ctx context.Context, keyType ManagedKeyType) (*BarrierView, *namespace.Namespace, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	if keyType != ManagedKeyTypePkcs11 {
		return nil, nil, fmt.Errorf("unsupported key type: %s", keyType)
	}

	keyTypeSubPath := string(keyType) + "/"

	view := r.view.SubView(keyTypeSubPath)

	return view, ns, nil
}
