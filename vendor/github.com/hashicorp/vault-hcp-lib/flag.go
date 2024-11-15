// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vaulthcplib

import "strconv"

type boolValue struct {
	target *bool
}

func (v *boolValue) String() string {
	return strconv.FormatBool(*v.target)
}

func (v *boolValue) Set(s string) error {
	value, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*v.target = value
	return nil
}

type stringValue struct {
	target *string
}

func (v *stringValue) String() string {
	return *v.target
}

func (v *stringValue) Set(s string) error {
	*v.target = s
	return nil
}
