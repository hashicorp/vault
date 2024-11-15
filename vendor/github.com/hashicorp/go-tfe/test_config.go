// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

type TestConfig struct {
	TestsEnabled bool `jsonapi:"attr,tests-enabled"`
}
