//go:build !fips_140_3

package config

import "github.com/hashicorp/vault/internalshared/configutil"

func IsValidListener(listener *configutil.Listener) error {
	return nil
}
