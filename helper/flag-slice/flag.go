// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sliceflag

import "strings"

// StringFlag implements the flag.Value interface and allows multiple
// calls to the same variable to append a list.
type StringFlag []string

func (s *StringFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *StringFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}
