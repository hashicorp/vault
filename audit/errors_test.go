package audit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestErrors_External_Internal_Upstream checks error messages for internal/external
// consumption match the expectations, including when they have an upstream error configured.
func TestErrors_External_Internal_Upstream(t *testing.T) {
	err := NewAuditError("magic.op", "the message", ErrInvalidParameter)

	// Check the internal/external messages look right.
	require.Equal(t, "magic.op: the message: invalid parameter", err.Internal().Error())
	require.Equal(t, "the message: invalid parameter", err.External().Error())

	// Configure an upstream error (we should have the same error both sides of this call)
	err2 := err.SetUpstream(errors.New("upstream error"))
	require.Equal(t, err2, err)

	// Check the internal/external messages look right (the internal includes the upstream error).
	require.Equal(t, "magic.op: the message: upstream error", err2.Internal().Error())
	require.Equal(t, "the message: invalid parameter", err.External().Error())
}
