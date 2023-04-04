// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sealmigration

import (
	"testing"
)

// TestSealMigration_TransitToShamir_Pre14 tests transit-to-shamir seal
// migration, using the pre-1.4 method of bring down the whole cluster to do
// the migration.
func TestSealMigration_TransitToShamir_Pre14(t *testing.T) {
	t.Parallel()
	// Note that we do not test integrated raft storage since this is
	// a pre-1.4 test.
	testVariousBackends(t, ParamTestSealMigrationTransitToShamir_Pre14, BasePort_TransitToShamir_Pre14, false)
}

// TestSealMigration_ShamirToTransit_Pre14 tests shamir-to-transit seal
// migration, using the pre-1.4 method of bring down the whole cluster to do
// the migration.
func TestSealMigration_ShamirToTransit_Pre14(t *testing.T) {
	t.Parallel()
	// Note that we do not test integrated raft storage since this is
	// a pre-1.4 test.
	testVariousBackends(t, ParamTestSealMigrationShamirToTransit_Pre14, BasePort_ShamirToTransit_Pre14, false)
}
