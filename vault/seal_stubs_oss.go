//go:build !enterprise

package vault

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

// isSealOldKeyError returns true if a value was decrypted using the
// old "unwrapSeal".
func isSealOldKeyError(err error) bool {
	return false
}
