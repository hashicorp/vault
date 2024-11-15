// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"strconv"
	"strings"
)

// formatFloat converts the floating-point number f to a string,
// after rounding it to the passed unit.
//
// Uses 'f' format (-ddd.dddddd, no exponent), and uses at most
// maxPrec digits after the decimal point.
func formatFloat(f float64, maxPrec int) string {
	v := strconv.FormatFloat(f, 'f', -1, 64)

	idx := strings.LastIndex(v, ".")
	if idx == -1 {
		return v
	}

	sublen := idx + maxPrec + 1
	if sublen > len(v) {
		sublen = len(v)
	}

	return v[:sublen]
}

// pointerOf returns a pointer to a.
func pointerOf[A any](a A) *A {
	return &a
}

// pointerCopy returns a new pointer to a.
func pointerCopy[A any](a *A) *A {
	if a == nil {
		return nil
	}
	na := *a
	return &na
}
