// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package server

import (
	"testing"
)

func TestLoadConfigFile_topLevel(t *testing.T) {
	testLoadConfigFile_topLevel(t, nil)
}

func TestLoadConfigFile_json2(t *testing.T) {
	testLoadConfigFile_json2(t, nil)
}

func TestParseEntropy(t *testing.T) {
	testParseEntropy(t, true)
}
