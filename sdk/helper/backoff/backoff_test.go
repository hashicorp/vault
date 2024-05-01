// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestBackoff_Basic tests that basic exponential backoff works as expected up to a max of 3 times.
func TestBackoff_Basic(t *testing.T) {
	for i := 0; i < 100; i++ {
		b := NewBackoff(3, 1*time.Millisecond, 10*time.Millisecond)
		x, err := b.Next()
		assert.Nil(t, err)
		assert.LessOrEqual(t, x, 1*time.Millisecond)
		assert.GreaterOrEqual(t, x, 750*time.Microsecond)

		x2, err := b.Next()
		assert.Nil(t, err)
		assert.LessOrEqual(t, x2, x*2)
		assert.GreaterOrEqual(t, x2, x*3/4)

		x3, err := b.Next()
		assert.Nil(t, err)
		assert.LessOrEqual(t, x3, x2*2)
		assert.GreaterOrEqual(t, x3, x2*3/4)

		_, err = b.Next()
		assert.NotNil(t, err)
	}
}

// TestBackoff_ZeroRetriesAlwaysFails checks that if retries is set to zero, then an error is returned immediately.
func TestBackoff_ZeroRetriesAlwaysFails(t *testing.T) {
	b := NewBackoff(0, 1*time.Millisecond, 10*time.Millisecond)
	_, err := b.Next()
	assert.NotNil(t, err)
}

// TestBackoff_MaxIsEnforced checks that the maximum backoff is enforced.
func TestBackoff_MaxIsEnforced(t *testing.T) {
	b := NewBackoff(1001, 1*time.Millisecond, 2*time.Millisecond)
	for i := 0; i < 1000; i++ {
		x, err := b.Next()
		assert.LessOrEqual(t, x, 2*time.Millisecond)
		assert.Nil(t, err)
	}
}
