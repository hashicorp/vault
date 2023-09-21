// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !386 && !arm && !windows

package raft

const initialMmapSize = 100 * 1024 * 1024 * 1024 // 100GB
