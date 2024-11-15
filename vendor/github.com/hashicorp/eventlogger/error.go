// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventlogger

import "errors"

var (
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrNodeNotFound     = errors.New("node not found")
)
