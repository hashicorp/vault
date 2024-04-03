package audit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestErrors_External_Internal_Upstream checks error messages for internal/external
// consumption match the expectations, including when they have an upstream error configured.
func TestErrors_External_Internal_Upstream(t *testing.T) {
	err := NewError("magic.op", "the message", ErrInvalidParameter)

	// Check the internal/external messages look right.
	require.Equal(t, "magic.op: the message: invalid parameter", err.Internal().Error())
	require.Equal(t, "the message: invalid parameter", err.External().Error())

	// Configure an upstream error (we should have the same error both sides of this call)
	err2 := err.Wrap(NewError("op2", "this is an error", errors.New("upstream error")))
	require.Equal(t, err2, err)

	// Check the internal/external messages look right (the internal includes the upstream error).
	require.Equal(t, "magic.op: the message: invalid parameter: op2: this is an error: upstream error", err2.Internal().Error())
	require.Equal(t, "the message: invalid parameter: this is an error: upstream error", err.External().Error())
}

func TestErrors_Internal(t *testing.T) {
	err1 := NewError("op.Thing1", "error1's error msg", errors.New("a sad error"))
	err2 := NewError("op.Thing2", "error2's error msg", errors.New("a bad error"))
	err3 := NewError("op.Cat", "error3's error msg", errors.New("lowly error"))
	err1.Wrap(err2.Wrap(err3))

	res1 := err1.Internal()
	require.EqualError(t, res1, "op.Thing1: error1's error msg: a sad error: op.Thing2: error2's error msg: a bad error: op.Cat: error3's error msg: lowly error")
}

func TestErrors_External(t *testing.T) {
	err1 := NewError("op.Thing1", "error1's error msg", errors.New("a sad error"))
	err2 := NewError("op.Thing2", "error2's error msg", errors.New("a bad error"))
	err3 := NewError("op.Cat", "error3's error msg", errors.New("lowly error"))

	err1.Wrap(err2.Wrap(err3))

	res1 := err1.External()
	require.EqualError(t, res1, "error1's error msg: a sad error: error2's error msg: a bad error: error3's error msg: lowly error")
}
