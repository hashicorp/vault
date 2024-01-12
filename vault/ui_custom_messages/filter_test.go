// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package uicustommessages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFindFilterAuthenticated is a trivial test that verifies that the
// FindFilter struct is correctly mutated when calling the Authenticated setter
// method.
func TestFindFilterAuthenticated(t *testing.T) {
	var filter FindFilter

	// Check that initially the authenticated field is nil
	assert.Nil(t, filter.authenticated)

	// Check that after calling Authenticated with false, that the authenticated
	// field points to a false value.
	filter.Authenticated(false)

	assert.NotNil(t, filter.authenticated)
	assert.False(t, *filter.authenticated)

	// Check that after calling Authenticated with true, that the authenticated
	// field points to a true value.
	filter.Authenticated(true)

	assert.NotNil(t, filter.authenticated)
	assert.True(t, *filter.authenticated)
}

// TestFindFilterActive is a trivial test that verifies that the FindFilter
// struct is correctly mutated when calling the Active setter method.
func TestFindFilterActive(t *testing.T) {
	var filter FindFilter

	// Check that initially the active field is nil
	assert.Nil(t, filter.active)

	// Check that after calling Active with false, that the active field points
	// to a false value.
	filter.Active(false)

	assert.NotNil(t, filter.active)
	assert.False(t, *filter.active)

	// Check that after calling Active with true, that the active field points
	// to a true value.
	filter.Active(true)

	assert.NotNil(t, filter.active)
	assert.True(t, *filter.active)
}

// TestFindFilterType is a trivial test that verifies that the FindFilter
// struct is correctly mutated when calling the Type setter method.
func TestFindFilterType(t *testing.T) {
	var filter FindFilter

	// Check that initially the messageType field is empty
	assert.Empty(t, filter.messageType)

	// Check that allowed message types can be used to set the field.
	for _, el := range allowedMessageTypes {
		filter.messageType = ""

		assert.NoError(t, filter.Type(el))
		assert.Equal(t, el, filter.messageType)
	}

	// Check that other values have no effect.
	filter.messageType = ""

	err := filter.Type("watermark")
	assert.Error(t, err)
	assert.Empty(t, filter.messageType)
}
