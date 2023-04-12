//go:build !enterprise

package server

func addExpectedEntConfig(c *Config, sentinelModules []string)                         {}
func addExpectedDefaultEntConfig(c *Config)                                            {}
func addExpectedEntSanitizedConfig(c map[string]interface{}, sentinelModules []string) {}
