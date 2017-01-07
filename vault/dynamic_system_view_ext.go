// +build vault,!ent

package vault

// IsPrimary checks if this is a primary Vault instance.
func (d dynamicSystemView) IsPrimary() bool {
	return true
}
