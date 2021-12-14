package vault

import (
	"context"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)


const (
	pkcs11     = ManagedKeyTypePkcs11
	theKeyName = "the-key-name"
)

func TestManagedKeyRegistry_crud(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	r := core.managedKeyRegistry

	ctx := namespace.RootContext(nil)
	var keys []string

	// The registry starts out empty
	keys = listManagedKeys(t, r, ctx, pkcs11)
	assert.Empty(t, keys, "The registry starts out empty")

	// Put a key on the registry
	keyConfig := &ManagedKeyConfiguration{
		Type: pkcs11,
		Name: theKeyName,
		RawParameters: map[string]interface{}{
			"slot": "123",
			"token_label": "the token label",
			"pin": "the pin",
			"key_id": "foo",
			"key_label": "bar",
			"mechanism": "0x1041",
		},
	}
	setManagedKey(t, r, ctx, keyConfig)

	// ListManagedKeys
	keys = listManagedKeys(t, r, ctx, pkcs11)
	assert.Contains(t, keys, theKeyName)
	assert.Len(t, keys, 1)

	// GetManagedKey
	actualConfig := getManagedKey(t, r, ctx, theKeyName, pkcs11)
	assert.Equal(t, keyConfig, actualConfig)

	// DeleteManagedKey
	deleteManagedKey(t, r, ctx, theKeyName, pkcs11)
	keys = listManagedKeys(t, r, ctx, pkcs11)
	assert.Nil(t, getManagedKey(t, r, ctx, theKeyName, pkcs11))
	assert.Empty(t, keys, "Deleting the key leaves the registry empty again")
}

func TestNewManagedKeyRegistry_supported_key_types(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	r := core.managedKeyRegistry

	ctx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)
	var err error

	// We support pkcs11
	_, err = r.ListManagedKeys(ctx, "pkcs11")
	assert.NoError(t, err)

	// And that is it
	_, err = r.ListManagedKeys(ctx, "awskms")
	assert.Error(t, err)
}

func TestManagedKeyRegistry_DeleteManagedKey(t *testing.T) {
	// Test that DeleteManagedKey pays attention to mount configuration

	core, _, _ := TestCoreUnsealed(t)
	r := core.managedKeyRegistry
	ctx := namespace.RootContext(nil)

	// Create a key

	keyConfig := &ManagedKeyConfiguration{
		Type: pkcs11,
		Name: theKeyName,
		RawParameters: map[string]interface{}{
			"key_id": "foo",
			"key_label": "bar",
		},
	}
	setManagedKey(t, r, ctx, keyConfig)

	// Create a mount and configure theKeyName to be allowed for use on it
	mountEntry, err := core.mounts.find(ctx, "secret/")
	if err != nil {
		t.Fatal(err)
	}
	mountEntry.Config.AllowedManagedKeys = []string{theKeyName}
	mountEntry.SyncCache()

	require.Error(t, r.DeleteManagedKey(ctx, theKeyName, pkcs11), "Key cannot be deleted since the mount's tune permits it use")

	mountEntry.Config.AllowedManagedKeys = []string{}
	mountEntry.SyncCache()
	require.NoError(t, r.DeleteManagedKey(ctx, theKeyName, pkcs11), "Key can be deleted since the mount's tune no longer permits it use")
}

func setManagedKey(t *testing.T, r *ManagedKeyRegistry, ctx context.Context, keyConfig *ManagedKeyConfiguration) {
	t.Helper()

	err := r.SetManagedKey(ctx, keyConfig)
	if err != nil {
		t.Fatal(err)
	}
}

func deleteManagedKey(t *testing.T, r *ManagedKeyRegistry, ctx context.Context, name string, keyType ManagedKeyType) {
	t.Helper()

	err := r.DeleteManagedKey(ctx, name, keyType)
	if err != nil {
		t.Fatal(err)
	}
}

func listManagedKeys(t *testing.T, r *ManagedKeyRegistry, ctx context.Context, keyType ManagedKeyType) []string {
	t.Helper()

	keys, err := r.ListManagedKeys(ctx, keyType)
	if err != nil {
		t.Fatal(err)
	}

	return keys
}

func getManagedKey(t *testing.T, r *ManagedKeyRegistry, ctx context.Context, name string, keyType ManagedKeyType) *ManagedKeyConfiguration{
	t.Helper()

	keyConfig, err := r.GetManagedKey(ctx, name, keyType)
	if err != nil {
		t.Fatal(err)
	}

	return keyConfig
}
