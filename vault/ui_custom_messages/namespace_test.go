// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package uicustommessages

import (
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/stretchr/testify/assert"
)

// TestCommunityEditionNamespaceManagerGetParentNamespace verifies that the
// (*CommunityEditionNamespaceManager).GetParentNamespace behaves as intended,
// which is to always return namespace.RootNamespace, regardless of the input.
func TestCommunityEditionNamespaceManagerGetParentNamespace(t *testing.T) {
	testNsManager := &CommunityEditionNamespaceManager{}

	// Verify root namespace
	assert.Equal(t, namespace.RootNamespace, testNsManager.GetParentNamespace(namespace.RootNamespace.Path))

	// Verify a different namespace
	testNamespace := namespace.Namespace{
		ID:   "abc123",
		Path: "test/",
	}
	assert.Equal(t, namespace.RootNamespace, testNsManager.GetParentNamespace(testNamespace.Path))

	// Verify that even a random string results in the root namespace
	assert.Equal(t, namespace.RootNamespace, testNsManager.GetParentNamespace("blah"))
}
