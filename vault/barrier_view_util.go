// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

func runICheck(v *BarrierView, expandedKey string, roErr error) bool { return true }
