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

func requireEphemeralBadRequest(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err, "ephemeral entity IDs must be rejected")
	require.ErrorContains(t, err, "ephemeral", "rejection should explain that the ID is ephemeral")
	var responseError *api.ResponseError
	require.ErrorAs(t, err, &responseError, "rejection should be an HTTP response error")
	require.Equal(t, http.StatusBadRequest, responseError.StatusCode, "ephemeral entity ID should return HTTP 400")
}

// TestEntityEphemeralGuard_GroupMembersRejectEphemeralIDs verifies invalid member IDs are rejected on both create and update without leaving a partially changed group.
func TestEntityEphemeralGuard_GroupMembersRejectEphemeralIDs(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	firstEntity, err := client.Logical().Write("identity/entity", map[string]interface{}{"name": "group-member-first"})
	require.NoError(t, err, "create first persistent entity")
	firstEntityID := firstEntity.Data["id"].(string)
	secondEntity, err := client.Logical().Write("identity/entity", map[string]interface{}{"name": "group-member-second"})
	require.NoError(t, err, "create second persistent entity")
	secondEntityID := secondEntity.Data["id"].(string)

	_, err = client.Logical().Write("identity/group", map[string]interface{}{
		"name":              "group-create-ephemeral",
		"member_entity_ids": []string{"eid.x", firstEntityID},
	})
	requireEphemeralBadRequest(t, err)
	missingGroup, err := client.Logical().Read("identity/group/name/group-create-ephemeral")
	require.NoError(t, err, "read rejected group creation")
	require.Nil(t, missingGroup, "rejected group should not be created")

	group, err := client.Logical().Write("identity/group", map[string]interface{}{
		"name":              "group-update-ephemeral",
		"member_entity_ids": []string{firstEntityID},
	})
	require.NoError(t, err, "create persistent group")
	groupID := group.Data["id"].(string)

	_, err = client.Logical().Write("identity/group/id/"+groupID, map[string]interface{}{
		"member_entity_ids": []string{"eid.x", secondEntityID},
	})
	requireEphemeralBadRequest(t, err)

	unchanged, err := client.Logical().Read("identity/group/id/" + groupID)
	require.NoError(t, err, "read group after rejected update")
	require.NotNil(t, unchanged, "group should remain after rejected update")
	require.Equal(t, []interface{}{firstEntityID}, unchanged.Data["member_entity_ids"], "rejected update must not change group membership")
}

// TestEntityEphemeralGuard_AliasCanonicalIDsRejectEphemeralIDs verifies alias creation and reassignment reject ephemeral canonical IDs without creating or changing an alias.
func TestEntityEphemeralGuard_AliasCanonicalIDsRejectEphemeralIDs(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{Type: "userpass"})
	require.NoError(t, err, "enable userpass auth mount")
	authResponse, err := client.Logical().Read("sys/auth")
	require.NoError(t, err, "read auth mounts")
	userpassMount, ok := authResponse.Data["userpass/"].(map[string]interface{})
	require.True(t, ok, "userpass auth mount should be present")
	mountAccessor, ok := userpassMount["accessor"].(string)
	require.True(t, ok, "userpass mount should expose an accessor")

	entity, err := client.Logical().Write("identity/entity", map[string]interface{}{"name": "alias-canonical-entity"})
	require.NoError(t, err, "create persistent entity")
	entityID := entity.Data["id"].(string)

	_, err = client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "alias-create-ephemeral",
		"mount_accessor": mountAccessor,
		"canonical_id":   "eid.x",
	})
	requireEphemeralBadRequest(t, err)
	entityAfterCreate, err := client.Logical().Read("identity/entity/id/" + entityID)
	require.NoError(t, err, "read entity after rejected alias creation")
	require.Empty(t, entityAfterCreate.Data["aliases"], "rejected alias creation must not attach an alias")

	alias, err := client.Logical().Write("identity/entity-alias", map[string]interface{}{
		"name":           "alias-update-persistent",
		"mount_accessor": mountAccessor,
		"canonical_id":   entityID,
	})
	require.NoError(t, err, "create persistent alias")
	aliasID := alias.Data["id"].(string)

	_, err = client.Logical().Write("identity/entity-alias/id/"+aliasID, map[string]interface{}{"canonical_id": "eid.x"})
	requireEphemeralBadRequest(t, err)
	unchangedAlias, err := client.Logical().Read("identity/entity-alias/id/" + aliasID)
	require.NoError(t, err, "read alias after rejected update")
	require.NotNil(t, unchangedAlias, "alias should remain after rejected update")
	require.Equal(t, entityID, unchangedAlias.Data["canonical_id"], "rejected update must preserve alias canonical ID")
}

// TestEntityEphemeralGuard_EntityBuilderRejectsEphemeralIDs verifies entity create and update requests reject reserved IDs before lookup and do not create named entities.
func TestEntityEphemeralGuard_EntityBuilderRejectsEphemeralIDs(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	_, err := client.Logical().Write("identity/entity", map[string]interface{}{
		"id":   "eid.x",
		"name": "entity-create-ephemeral",
	})
	requireEphemeralBadRequest(t, err)
	missingEntity, err := client.Logical().Read("identity/entity/name/entity-create-ephemeral")
	require.NoError(t, err, "read rejected entity creation")
	require.Nil(t, missingEntity, "rejected entity create should not create a named entity")

	_, err = client.Logical().Write("identity/entity/id/eid.update", map[string]interface{}{"name": "entity-update-ephemeral"})
	requireEphemeralBadRequest(t, err)
	missingUpdatedEntity, err := client.Logical().Read("identity/entity/name/entity-update-ephemeral")
	require.NoError(t, err, "read rejected entity update")
	require.Nil(t, missingUpdatedEntity, "rejected entity update should not create an entity")
}

// TestEntityEphemeralGuard_DeleteRejectsEphemeralIDs verifies individual and batch deletion reject ephemeral IDs before deleting any persistent entity.
func TestEntityEphemeralGuard_DeleteRejectsEphemeralIDs(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	entity, err := client.Logical().Write("identity/entity", map[string]interface{}{"name": "batch-delete-survivor"})
	require.NoError(t, err, "create persistent entity")
	entityID := entity.Data["id"].(string)

	_, err = client.Logical().Delete("identity/entity/id/eid.x")
	requireEphemeralBadRequest(t, err)

	_, err = client.Logical().Write("identity/entity/batch-delete", map[string]interface{}{
		"entity_ids": []string{"eid.x", entityID},
	})
	requireEphemeralBadRequest(t, err)
	remainingEntity, err := client.Logical().Read("identity/entity/id/" + entityID)
	require.NoError(t, err, "read entity after rejected batch delete")
	require.NotNil(t, remainingEntity, "batch delete must not delete persistent entities when any ID is ephemeral")
}

// TestEntityEphemeralGuard_MissingEphemeralRecordStillRejected verifies a cleaned-up ephemeral ID is rejected by its reserved prefix rather than reported as an ordinary missing entity.
func TestEntityEphemeralGuard_MissingEphemeralRecordStillRejected(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	_, err := client.Logical().Write("identity/entity/id/eid.cleaned-up", map[string]interface{}{"name": "cleaned-up-ephemeral"})
	requireEphemeralBadRequest(t, err)
}
