// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAuth_ActorEntityFields verifies that ActorEntityID and ActorEntityName
// round-trip correctly through JSON marshalling and unmarshalling, and that
// both fields appear with their expected JSON keys when set.
func TestAuth_ActorEntityFields(t *testing.T) {
	t.Parallel()

	auth := &Auth{
		EntityID:        "subject-entity-123",
		DisplayName:     "subject-user",
		ActorEntityID:   "actor-entity-456",
		ActorEntityName: "actor-service",
	}

	data, err := json.Marshal(auth)
	require.NoError(t, err)
	require.Contains(t, string(data), `"actor_entity_id":"actor-entity-456"`)
	require.Contains(t, string(data), `"actor_entity_name":"actor-service"`)

	var auth2 Auth
	err = json.Unmarshal(data, &auth2)
	require.NoError(t, err)
	require.Equal(t, "actor-entity-456", auth2.ActorEntityID)
	require.Equal(t, "actor-service", auth2.ActorEntityName)
}

// TestAuth_ActorEntityFields_OmitEmpty verifies that ActorEntityID and
// ActorEntityName are omitted from the JSON output when not set, preventing
// empty actor fields from appearing in audit entries for requests without
// an actor entity.
func TestAuth_ActorEntityFields_OmitEmpty(t *testing.T) {
	t.Parallel()

	auth := &Auth{
		EntityID:    "subject-entity-123",
		DisplayName: "subject-user",
		// ActorEntityID and ActorEntityName intentionally not set
	}

	data, err := json.Marshal(auth)
	require.NoError(t, err)
	require.NotContains(t, string(data), "actor_entity_id")
	require.NotContains(t, string(data), "actor_entity_name")
}
