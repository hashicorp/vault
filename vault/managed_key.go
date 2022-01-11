package vault

import "context"

type ManagedKey interface {
	Name() string
}

type ManagedKeySystemView interface {
	GetManagedKey(context.Context, string) (ManagedKey, error)
}
