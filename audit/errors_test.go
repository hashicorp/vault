// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestErrors_ConvertToExternalError is used to check that we 'mute' errors which
// have an internal error in their tree.
func TestErrors_ConvertToExternalError(t *testing.T) {
	t.Parallel()

	err := fmt.Errorf("wrap this error: %w", ErrInternal)
	res := ConvertToExternalError(err)
	require.EqualError(t, res, "audit system internal error")

	err = fmt.Errorf("test: %w", errors.New("this is just an error"))
	res = ConvertToExternalError(err)
	require.Equal(t, "test: this is just an error", res.Error())
}
