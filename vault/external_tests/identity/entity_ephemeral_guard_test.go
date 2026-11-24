// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"net/http"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestEntityEphemeralGuard_MergeRejectsEphemeralToEntity verifies that an ephemeral destination is rejected before a persistent source entity can be changed.
func TestEntityEphemeralGuard_MergeRejectsEphemeralToEntity(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	entityResp, err := client.Logical().Write("identity/entity", map[string]interface{}{"name": "entity-a"})
	require.NoError(t, err, "create source entity")
	require.NotNil(t, entityResp, "entity creation should return a response")
	entityID, ok := entityResp.Data["id"].(string)
	require.True(t, ok, "created entity should have an ID")
	require.NotEmpty(t, entityID, "created entity ID should not be empty")

	_, err = client.Logical().Write("identity/entity/merge", map[string]interface{}{
		"to_entity_id":    "eid.abc123",
		"from_entity_ids": []string{entityID},
	})
	require.Error(t, err, "merge to an ephemeral entity must be rejected")
	require.ErrorContains(t, err, "ephemeral", "rejection should explain that the ID is ephemeral")
	var responseError *api.ResponseError
	require.ErrorAs(t, err, &responseError, "rejection should be an HTTP response error")
	require.Equal(t, http.StatusBadRequest, responseError.StatusCode, "ephemeral merge should return HTTP 400")

	unchanged, err := client.Logical().Read("identity/entity/id/" + entityID)
	require.NoError(t, err, "read source entity after rejected merge")
	require.NotNil(t, unchanged, "source entity should remain after rejected merge")
	require.Equal(t, entityID, unchanged.Data["id"], "source entity ID should remain unchanged")
	require.Equal(t, "entity-a", unchanged.Data["name"], "source entity name should remain unchanged")
	require.Empty(t, unchanged.Data["merged_entity_ids"], "source entity should not record a merge")
}

// TestEntityEphemeralGuard_MergeRejectsEphemeralFromEntity verifies that validating every source ID up front prevents a partial multi-entity merge.
func TestEntityEphemeralGuard_MergeRejectsEphemeralFromEntity(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	toResp, err := client.Logical().Write("identity/entity", map[string]interface{}{"name": "entity-b"})
	require.NoError(t, err, "create destination entity")
	require.NotNil(t, toResp, "destination creation should return a response")
	toID, ok := toResp.Data["id"].(string)
	require.True(t, ok, "destination entity should have an ID")
	require.NotEmpty(t, toID, "destination entity ID should not be empty")

	fromResp, err := client.Logical().Write("identity/entity", map[string]interface{}{"name": "entity-c"})
	require.NoError(t, err, "create source entity")
	require.NotNil(t, fromResp, "source creation should return a response")
	fromID, ok := fromResp.Data["id"].(string)
	require.True(t, ok, "source entity should have an ID")
	require.NotEmpty(t, fromID, "source entity ID should not be empty")

	_, err = client.Logical().Write("identity/entity/merge", map[string]interface{}{
		"to_entity_id":    toID,
		"from_entity_ids": []string{fromID, "eid.abc123"},
	})
	require.Error(t, err, "merge with an ephemeral source must be rejected")
	require.ErrorContains(t, err, "ephemeral", "rejection should explain that the ID is ephemeral")
	var responseError *api.ResponseError
	require.ErrorAs(t, err, &responseError, "rejection should be an HTTP response error")
	require.Equal(t, http.StatusBadRequest, responseError.StatusCode, "ephemeral merge should return HTTP 400")

	unchangedFrom, err := client.Logical().Read("identity/entity/id/" + fromID)
	require.NoError(t, err, "read source entity after rejected merge")
	require.NotNil(t, unchangedFrom, "source entity should remain after rejected merge")
	require.Equal(t, "entity-c", unchangedFrom.Data["name"], "source entity name should remain unchanged")

	unchangedTo, err := client.Logical().Read("identity/entity/id/" + toID)
	require.NoError(t, err, "read destination entity after rejected merge")
	require.NotNil(t, unchangedTo, "destination entity should remain after rejected merge")
	require.Empty(t, unchangedTo.Data["merged_entity_ids"], "destination entity should not record a partial merge")
}
