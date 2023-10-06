//go:build !enterprise

package builtinplugins

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func entAddExtPlugins(r *registry) {
}
