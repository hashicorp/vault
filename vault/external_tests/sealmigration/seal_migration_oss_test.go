//go:build !enterprise

package sealmigration

import "testing"

func TestSealMigration_TransitToShamir_Recovery(t *testing.T) {
	t.Parallel()
	testVariousBackends(t, ParamTestSealMigrationTransitToShamir_Recovery, BasePort_TransitToShamir_Recovery, true)
}
