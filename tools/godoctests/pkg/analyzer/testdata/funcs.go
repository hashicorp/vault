// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testdata

import "testing"

// Test_GoDocOK is a test that has a go doc
func Test_GoDocOK(t *testing.T) {}

func Test_NoGoDocFails(t *testing.T) {} // want "Test Test_NoGoDocFails is missing a go doc"

// This test does not have a go doc beginning with the function name
func Test_BadGoDocFails(t *testing.T) {} // want "Test Test_BadGoDocFails must have a go doc beginning with the function name"

func test_TestHelperNoGoDocOK(t *testing.T) {}

func Test_DifferentSignatureNoGoDocOK() {}

func Test_DifferentSignature2NoGoDocOK(t *testing.T, a int) {}
