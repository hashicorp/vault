// +build !enterprise

package server

func addExpectedEntConfig(c *Config, sentinelModules []string)                         {}
func addExpectedEntSanitizedConfig(c map[string]interface{}, sentinelModules []string) {}
