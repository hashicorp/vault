// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package audit

// formatterConfigEnt provides extensions to a formatterConfig which behave differently
// for Enterprise and community edition.
// NOTE: Use newFormatterConfigEnt to initialize the formatterConfigEnt struct.
type formatterConfigEnt struct{}

// newFormatterConfigEnt should be used to create formatterConfigEnt.
func newFormatterConfigEnt(config map[string]string) (formatterConfigEnt, error) {
	return formatterConfigEnt{}, nil
}
