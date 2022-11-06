//go:build !enterprise

package server

func addExpectedEntConfig(c *Config, sentinelModules []string)                 {}
func addExpectedEntSanitizedConfig(c map[string]any, sentinelModules []string) {}
