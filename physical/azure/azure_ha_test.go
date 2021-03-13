package azure

import (
	"testing"

	"github.com/hashicorp/vault/sdk/physical"
)

func TestAzureHABackend(t *testing.T) {
	checkTestPreReqs(t)
	backend1, cleanup1 := testFixture(t)
	defer cleanup1()
	backend2, cleanup2 := testFixture(t)
	defer cleanup2()

	physical.ExerciseHABackend(t, backend1, backend2)
}
