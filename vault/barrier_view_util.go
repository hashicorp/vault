// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

func runICheck(v *BarrierView, expandedKey string, roErr error) bool { return true }
