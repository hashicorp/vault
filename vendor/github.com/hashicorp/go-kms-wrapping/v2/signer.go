// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package wrapping

import (
	"context"
)

// SigInfoSigner defines common capabilities for creating a msg signature in the
// form of a SigInfo
type SigInfoSigner interface {
	// Sign creates a msg signature in the form of a SigInfo
	Sign(ctx context.Context, msg []byte, opt ...Option) (*SigInfo, error)
}

// SigInfoVerifier defines common capabilities for verifying a msg signature in
// the form of a SigInfo
type SigInfoVerifier interface {
	// Verify a msg signature in the form of a SigInfo
	Verify(ctx context.Context, msg []byte, sig *SigInfo, opt ...Option) (bool, error)
}
