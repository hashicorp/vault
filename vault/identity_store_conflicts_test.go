// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/stretchr/testify/require"
)

// TestDuplicateReportingErrorResolver tests that the reporting error resolver
// correctly records, identifies, sorts and outputs information about duplicate
// entities.
func TestDuplicateReportingErrorResolver(t *testing.T) {
	t.Parallel()

	entities := [][]string{
		// Some unduplicated entities in a different namespaces
		{"root", "foo"},
		{"root", "bar"},
		{"admin", "foo"},
		{"developers", "BAR"},

		// Some exact-match duplicates in different namespaces
		{"root", "exact-dupe-1"},
		{"root", "exact-dupe-1"},
		{"admin", "exact-dupe-1"},
		{"admin", "exact-dupe-1"},

		// Some different-case duplicates in different namespaces
		{"root", "different-case-dupe-1"},
		{"root", "DIFFERENT-CASE-DUPE-1"},
		{"admin", "different-case-dupe-1"},
		{"admin", "DIFFERENT-case-DUPE-1"},
		{"admin", "different-case-DUPE-1"},
	}

	// Note that `local-` prefix here is used to define a mount as local as well
	// as used in it's name.
	aliases := [][]string{
		// Unduplicated aliases on different mounts
		{"mount1", "alias1"},
		{"mount2", "alias1"},
		{"mount2", "alias2"},
		{"local-mount", "alias1"},

		// We don't bother testing exact-match aliases since they will have been
		// merged by the time they are reported)

		// Different-case aliases on different mounts (some local)
		{"mount1", "different-case-alias-1"},
		{"mount1", "DIFFERENT-CASE-ALIAS-1"},
		{"mount2", "different-case-alias-1"},
		{"mount2", "DIFFERENT-CASE-ALIAS-1"},
		{"mount2", "different-CASE-ALIAS-1"},
		{"local-mount", "DIFFERENT-CASE-ALIAS-1"},
		{"local-mount", "different-CASE-ALIAS-1"},
	}

	expectReport := `
DUPLICATES DETECTED, see following logs for details and refer to https://developer.hashicorp.com/vault/docs/upgrading/identity-deduplication for resolution.:
1 different-case local entity alias duplicates found (potential security risk):
local entity-alias "DIFFERENT-CASE-ALIAS-1" with mount accessor "local-mount" duplicates 1 others: id="00000000-0000-0000-0000-000000000009" canonical_id="11111111-0000-0000-0000-000000000009" force_deduplication="would merge others into this entity"
local entity-alias "different-CASE-ALIAS-1" with mount accessor "local-mount" duplicates 1 others: id="00000000-0000-0000-0000-000000000010" canonical_id="11111111-0000-0000-0000-000000000010" force_deduplication="would merge into entity 11111111-0000-0000-0000-000000000009"
end of different-case local entity-alias duplicates:
2 different-case entity alias duplicates found (potential security risk):
entity-alias "different-case-alias-1" with mount accessor "mount1" duplicates 1 others: id="00000000-0000-0000-0000-000000000004" canonical_id="11111111-0000-0000-0000-000000000004" force_deduplication="would merge others into this entity"
entity-alias "DIFFERENT-CASE-ALIAS-1" with mount accessor "mount1" duplicates 1 others: id="00000000-0000-0000-0000-000000000005" canonical_id="11111111-0000-0000-0000-000000000005" force_deduplication="would merge into entity 11111111-0000-0000-0000-000000000004"
entity-alias "different-case-alias-1" with mount accessor "mount2" duplicates 2 others: id="00000000-0000-0000-0000-000000000006" canonical_id="11111111-0000-0000-0000-000000000006" force_deduplication="would merge others into this entity"
entity-alias "DIFFERENT-CASE-ALIAS-1" with mount accessor "mount2" duplicates 2 others: id="00000000-0000-0000-0000-000000000007" canonical_id="11111111-0000-0000-0000-000000000007" force_deduplication="would merge into entity 11111111-0000-0000-0000-000000000006"
entity-alias "different-CASE-ALIAS-1" with mount accessor "mount2" duplicates 2 others: id="00000000-0000-0000-0000-000000000008" canonical_id="11111111-0000-0000-0000-000000000008" force_deduplication="would merge into entity 11111111-0000-0000-0000-000000000006"
end of different-case entity-alias duplicates:
4 entity duplicates found:
entity "different-case-dupe-1" with namespace ID "admin" duplicates 2 others: id="00000000-0000-0000-0000-000000000010" force_deduplication="would not rename"
entity "DIFFERENT-case-DUPE-1" with namespace ID "admin" duplicates 2 others: id="00000000-0000-0000-0000-000000000011" force_deduplication="would rename to DIFFERENT-case-DUPE-1-00000000-0000-0000-0000-000000000011"
entity "different-case-DUPE-1" with namespace ID "admin" duplicates 2 others: id="00000000-0000-0000-0000-000000000012" force_deduplication="would rename to different-case-DUPE-1-00000000-0000-0000-0000-000000000012"
entity "exact-dupe-1" with namespace ID "admin" duplicates 1 others: id="00000000-0000-0000-0000-000000000006" force_deduplication="would not rename"
entity "exact-dupe-1" with namespace ID "admin" duplicates 1 others: id="00000000-0000-0000-0000-000000000007" force_deduplication="would rename to exact-dupe-1-00000000-0000-0000-0000-000000000007"
entity "different-case-dupe-1" with namespace ID "root" duplicates 1 others: id="00000000-0000-0000-0000-000000000008" force_deduplication="would not rename"
entity "DIFFERENT-CASE-DUPE-1" with namespace ID "root" duplicates 1 others: id="00000000-0000-0000-0000-000000000009" force_deduplication="would rename to DIFFERENT-CASE-DUPE-1-00000000-0000-0000-0000-000000000009"
entity "exact-dupe-1" with namespace ID "root" duplicates 1 others: id="00000000-0000-0000-0000-000000000004" force_deduplication="would not rename"
entity "exact-dupe-1" with namespace ID "root" duplicates 1 others: id="00000000-0000-0000-0000-000000000005" force_deduplication="would rename to exact-dupe-1-00000000-0000-0000-0000-000000000005"
end of entity duplicates:
4 group duplicates found:
group "different-case-dupe-1" with namespace ID "admin" duplicates 2 others: id="00000000-0000-0000-0000-000000000010" force_deduplication="would not rename"
group "DIFFERENT-case-DUPE-1" with namespace ID "admin" duplicates 2 others: id="00000000-0000-0000-0000-000000000011" force_deduplication="would rename to DIFFERENT-case-DUPE-1-00000000-0000-0000-0000-000000000011"
group "different-case-DUPE-1" with namespace ID "admin" duplicates 2 others: id="00000000-0000-0000-0000-000000000012" force_deduplication="would rename to different-case-DUPE-1-00000000-0000-0000-0000-000000000012"
group "exact-dupe-1" with namespace ID "admin" duplicates 1 others: id="00000000-0000-0000-0000-000000000006" force_deduplication="would not rename"
group "exact-dupe-1" with namespace ID "admin" duplicates 1 others: id="00000000-0000-0000-0000-000000000007" force_deduplication="would rename to exact-dupe-1-00000000-0000-0000-0000-000000000007"
group "different-case-dupe-1" with namespace ID "root" duplicates 1 others: id="00000000-0000-0000-0000-000000000008" force_deduplication="would not rename"
group "DIFFERENT-CASE-DUPE-1" with namespace ID "root" duplicates 1 others: id="00000000-0000-0000-0000-000000000009" force_deduplication="would rename to DIFFERENT-CASE-DUPE-1-00000000-0000-0000-0000-000000000009"
group "exact-dupe-1" with namespace ID "root" duplicates 1 others: id="00000000-0000-0000-0000-000000000004" force_deduplication="would not rename"
group "exact-dupe-1" with namespace ID "root" duplicates 1 others: id="00000000-0000-0000-0000-000000000005" force_deduplication="would rename to exact-dupe-1-00000000-0000-0000-0000-000000000005"
end of group duplicates:
end of identity duplicate report, refer to https://developer.hashicorp.com/vault/docs/upgrading/identity-deduplication for resolution.:
`

	// Create a new errorResolver
	r := newDuplicateReportingErrorResolver(log.NewNullLogger())

	for i, pair := range entities {
		// Create a fake UUID based on the index this makes sure sort order is
		// preserved when eyeballing the expected report.
		id := fmt.Sprintf("00000000-0000-0000-0000-%012d", i)
		// Create a new entity with the pair
		entity := &identity.Entity{
			ID:          id,
			Name:        pair[1],
			NamespaceID: pair[0],
		}

		// Call ResolveEntities, assume existing is nil for now. In real life we
		// should be passed the existing entity for the exact match dupes but we
		// don't depend on that so it's fine to omit.
		_ = r.ResolveEntities(context.Background(), nil, entity)
		// Don't care about the actual error here since it would be ignored in
		// case-sensitive mode anyway.

		// Also, since the data model is the same, pretend these are groups too
		group := &identity.Group{
			ID:          id,
			Name:        pair[1],
			NamespaceID: pair[0],
		}
		_ = r.ResolveGroups(context.Background(), nil, group)
	}

	// Load aliases second because that is realistic and yet we want to report on
	// them first.
	for i, pair := range aliases {
		entity := &identity.Entity{
			ID:          fmt.Sprintf("11111111-0000-0000-0000-%012d", i),
			Name:        pair[1] + "-entity",
			NamespaceID: "root",
		}
		alias := &identity.Alias{
			ID:            fmt.Sprintf("00000000-0000-0000-0000-%012d", i),
			CanonicalID:   entity.ID,
			Name:          pair[1],
			MountAccessor: pair[0],
			// Parse our hacky DSL to define some alias mounts as local
			Local: strings.HasPrefix(pair[0], "local-"),
		}
		_ = r.ResolveAliases(context.Background(), entity, nil, alias)
	}

	// "log" the report and check it matches expected report below.
	var testLog identityTestWarnLogger
	r.LogReport(&testLog)

	// Dump the raw report to make it easier to copy paste/read
	t.Log("\n\n" + testLog.buf.String())

	require.Equal(t,
		strings.TrimSpace(expectReport),
		strings.TrimSpace(testLog.buf.String()),
	)
}

type identityTestWarnLogger struct {
	buf bytes.Buffer
}

func (l *identityTestWarnLogger) Warn(msg string, args ...interface{}) {
	l.buf.WriteString(msg + ":")
	if len(args)%2 != 0 {
		panic("args must be key-value pairs")
	}
	for i := 0; i < len(args); i += 2 {
		l.buf.WriteString(fmt.Sprintf(" %s=%q", args[i], args[i+1]))
	}
	l.buf.WriteString("\n")
}

// TestDuplicateRenameResolver tests that the rename resolver
// correctly renames pre-existing entities and groups.
func TestDuplicateRenameResolver(t *testing.T) {
	t.Parallel()

	entities := map[string][]string{
		// 2 non-duplicates and 2 duplicates
		"root": {
			"foo",
			"bar",
			"exact-dupe-1", "exact-dupe-1",
			"different-case-dupe-1", "DIFFERENT-CASE-DUPE-1",
		},
		// 1 non-duplicate and 3 duplicates
		"admin": {
			"foo",
			"exact-dupe-1", "exact-dupe-1",
			"different-case-dupe-1", "DIFFERENT-case-DUPE-1", "different-case-DUPE-1",
		},
		"developers": {"BAR"},
	}

	// Create a new errorResolver
	r := &renameResolver{log.NewNullLogger()}

	seenEntities := make(map[string]*identity.Entity)
	seenGroups := make(map[string]*identity.Group)

	for ns, entityList := range entities {
		for i, name := range entityList {

			id := fmt.Sprintf("00000000-0000-0000-0000-%012d", i)
			// Create a new entity with the name/ns pair
			entity := &identity.Entity{
				ID:          id,
				Name:        name,
				NamespaceID: ns,
			}

			// Simulate a MemDB lookup
			existingEntity := seenEntities[name]
			err := r.ResolveEntities(context.Background(), existingEntity, entity)
			require.NoError(t, err)

			if existingEntity != nil {
				require.Equal(t, name+"-"+id, entity.Name)
				require.Equal(t, existingEntity.ID, entity.Metadata["duplicate_of_canonical_id"])
			} else {
				seenEntities[name] = entity
			}

			// Also, since the data model is the same, pretend these are groups too
			group := &identity.Group{
				ID:          id,
				Name:        name,
				NamespaceID: ns,
			}

			// More MemDB mocking
			existingGroup := seenGroups[name]
			err = r.ResolveGroups(context.Background(), existingGroup, group)
			require.NoError(t, err)

			if existingGroup != nil {
				require.Equal(t, name+"-"+id, group.Name)
				require.Equal(t, existingGroup.ID, group.Metadata["duplicate_of_canonical_id"])
			} else {
				seenGroups[name] = group
			}
		}
	}

	// No need to test entity alias merges here, since that's handled separately.
}
