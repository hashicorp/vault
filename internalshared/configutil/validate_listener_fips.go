// +build fips_140_3

package configutil

import "errors"

func IsValidListener(listener *Listener) error {
	return errors.New("invalid configuration: all HTTP listeners must have TLS enabled")
}
