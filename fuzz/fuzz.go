// +build gofuzz

package fuzz

import (
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/random"
	"github.com/hashicorp/vault/vault"
)

func FuzzParseACLPolicy(data []byte) int {
	_, err := vault.ParseACLPolicy(namespace.RootNamespace, string(data))
	if err != nil {
		return 0
	}
	return 1
}

func FuzzParsePolicy(data []byte) int {
	_, err := random.ParsePolicy(string(data))
	if err != nil {
		return 0
	}
	return 1
}
