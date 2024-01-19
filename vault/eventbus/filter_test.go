// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package eventbus

import (
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/stretchr/testify/assert"
)

// TestFilters_AddRemoveMatchLocal checks that basic matching, adding, and removing of patterns all work.
func TestFilters_AddRemoveMatchLocal(t *testing.T) {
	f := NewFilters("self")
	ns := &namespace.Namespace{
		ID:   "ns1",
		Path: "ns1",
	}

	assert.False(t, f.localMatch(ns, "abc"))
	assert.False(t, f.anyMatch(ns, "abc"))
	f.addNsPattern("self", ns, "abc")
	assert.True(t, f.localMatch(ns, "abc"))
	assert.False(t, f.localMatch(ns, "abcd"))
	assert.True(t, f.anyMatch(ns, "abc"))
	assert.False(t, f.anyMatch(ns, "abcd"))
	f.removeNsPattern("self", ns, "abc")
	assert.False(t, f.localMatch(ns, "abc"))
	assert.False(t, f.anyMatch(ns, "abc"))
}
