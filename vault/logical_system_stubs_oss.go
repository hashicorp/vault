//go:build !enterprise

package vault

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func entUnauthenticatedPaths() []string {
	return []string{}
}
