// +build !fips_140_3

package configutil

func IsValidListener(listener *Listener) error {
	return nil
}
