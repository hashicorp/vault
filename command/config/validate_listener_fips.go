// +build fips_140_3

package config

import (
	"errors"
	"github.com/hashicorp/vault/internalshared/configutil"
)

func IsValidListener(listener *configutil.Listener) error {
	return errors.New("invalid configuration: all HTTP listeners must have TLS enabled")
}
