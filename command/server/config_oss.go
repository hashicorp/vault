//go:build !enterprise

package server

func (c *Config) IsMultisealEnabled() bool {
	return false
}
