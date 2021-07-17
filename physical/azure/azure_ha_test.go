package azure

import (
	"testing"

	"github.com/hashicorp/vault/sdk/physical"
)

func TestAzureHABackend(t *testing.T) {
	backend, cleanup := testFixture(t, withHA())
	defer cleanup()
	// point the two HA backends to the same backend azurite storage instance.
	// if you point the two HA backends to two diffferent azurite storage
	// instances, both HA backends will be able to aquire leadership.
	physical.ExerciseHABackend(t, backend, backend)
}
