// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package server

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func addExpectedEntConfig(c *Config, sentinelModules []string)                         {}
func addExpectedDefaultEntConfig(c *Config)                                            {}
func addExpectedEntSanitizedConfig(c map[string]interface{}, sentinelModules []string) {}
