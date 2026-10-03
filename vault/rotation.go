// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

const (
	// rotationLocalSubPath is the sub-path used for the rotation manager's
	// Local (non-replicated) view. This is nested under the system view.
	rotationLocalSubPath = "rotation-local/"

	// orphanLocalSubPath is the sub-path used for the rotation manager's
	// Local (non-replicated) view for orphaned rotations. This is nested under the system view.
	orphanLocalSubPath = "orphaned-rotation-local/"
)
