// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFilters_AddRemoveMatchLocal checks that basic matching, adding, and removing of patterns all work.
func TestFilters_AddRemoveMatchLocal(t *testing.T) {
	f := NewFilters("self")

	assert.False(t, f.localMatch("ns1", "abc"))
	assert.False(t, f.anyMatch("ns1", "abc"))
	f.addPattern("self", []string{"ns1"}, "abc")
	assert.True(t, f.localMatch("ns1", "abc"))
	assert.False(t, f.localMatch("ns1", "abcd"))
	assert.True(t, f.anyMatch("ns1", "abc"))
	assert.False(t, f.anyMatch("ns1", "abcd"))
	f.removePattern("self", []string{"ns1"}, "abc")
	assert.False(t, f.localMatch("ns1", "abc"))
	assert.False(t, f.anyMatch("ns1", "abc"))
}

// TestFilters_ParallelAnyMatch checks that anyMatch works with parallel set to true.
func TestFilters_ParallelAnyMatch(t *testing.T) {
	f := NewFilters("self")
	f.parallel = true

	f.addPattern("self", []string{"ns1"}, "abc")
	assert.True(t, f.anyMatch("ns1", "abc"))
	assert.False(t, f.anyMatch("ns1", "abcd"))
}
