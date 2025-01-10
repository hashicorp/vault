// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/identity"
)

var errDuplicateIdentityName = errors.New("duplicate identity name")

// ConflictResolver defines the interface for resolving conflicts between
// entities, groups, and aliases. All methods should implement a check for
// existing=nil. This is an intentional design choice to allow the caller to
// search for extra information if necessary.
type ConflictResolver interface {
	ResolveEntities(ctx context.Context, existing, duplicate *identity.Entity) error
	ResolveGroups(ctx context.Context, existing, duplicate *identity.Group) error
	ResolveAliases(ctx context.Context, parent *identity.Entity, existing, duplicate *identity.Alias) error
}

// The errorResolver is a ConflictResolver that logs a warning message when a
// pre-existing Identity artifact is found with the same factors as a new one.
type errorResolver struct {
	logger log.Logger
}

// ResolveEntities logs a warning message when a pre-existing Entity is found
// and returns a duplicate name error, which should be handled by the caller by
// putting the system in case-sensitive mode.
func (r *errorResolver) ResolveEntities(ctx context.Context, existing, duplicate *identity.Entity) error {
	if existing == nil {
		return nil
	}

	r.logger.Warn(errDuplicateIdentityName.Error(),
		"entity_name", duplicate.Name,
		"duplicate_of_name", existing.Name,
		"duplicate_of_id", existing.ID,
		"action", "merge the duplicate entities into one")

	return errDuplicateIdentityName
}

// ResolveGroups logs a warning message when a pre-existing Group is found and
// returns a duplicate name error, which should be handled by the caller by
// putting the system in case-sensitive mode.
func (r *errorResolver) ResolveGroups(ctx context.Context, existing, duplicate *identity.Group) error {
	if existing == nil {
		return nil
	}

	r.logger.Warn(errDuplicateIdentityName.Error(),
		"group_name", duplicate.Name,
		"duplicate_of_name", existing.Name,
		"duplicate_of_id", existing.ID,
		"action", "merge the contents of duplicated groups into one and delete the other")

	return errDuplicateIdentityName
}

// ResolveAliases logs a warning message when a pre-existing Alias is found and
// returns a duplicate name error, which should be handled by the caller by
// putting the system in case-sensitive mode.
func (r *errorResolver) ResolveAliases(ctx context.Context, parent *identity.Entity, existing, duplicate *identity.Alias) error {
	if existing == nil {
		return nil
	}

	r.logger.Warn(errDuplicateIdentityName.Error(),
		"alias_name", duplicate.Name,
		"mount_accessor", duplicate.MountAccessor,
		"local", duplicate.Local,
		"entity_name", parent.Name,
		"alias_canonical_id", duplicate.CanonicalID,
		"duplicate_of_name", existing.Name,
		"duplicate_of_canonical_id", existing.CanonicalID,
		"action", "merge the canonical entity IDs into one")

	return errDuplicateIdentityName
}

// duplicateReportingErrorResolver collects duplicate information and optionally
// logs a report on all the duplicates. We don't embed an errorResolver here
// because we _don't_ want it's side effect of warning on just some duplicates
// as we go as that's confusing when we have a more comprehensive report. The
// only other behavior it has is to return a constant error which we can just do
// ourselves.
type duplicateReportingErrorResolver struct {
	// seen* track the unique factors for each identity artifact, so
	// that we can report on any duplication including different-case duplicates
	// when in case-sensitive mode.
	//
	// Since this is only ever called from `load*` methods on IdentityStore during
	// an unseal we can assume that it's all from a single goroutine and does'nt
	// need locking.
	seenEntities     map[string][]*identity.Entity
	seenGroups       map[string][]*identity.Group
	seenAliases      map[string][]*identity.Alias
	seenLocalAliases map[string][]*identity.Alias
	logger           hclog.Logger
}

func newDuplicateReportingErrorResolver(logger hclog.Logger) *duplicateReportingErrorResolver {
	return &duplicateReportingErrorResolver{
		seenEntities:     make(map[string][]*identity.Entity),
		seenGroups:       make(map[string][]*identity.Group),
		seenAliases:      make(map[string][]*identity.Alias),
		seenLocalAliases: make(map[string][]*identity.Alias),
		logger:           logger,
	}
}

func (r *duplicateReportingErrorResolver) ResolveEntities(ctx context.Context, existing, duplicate *identity.Entity) error {
	entityKey := fmt.Sprintf("%s/%s", duplicate.NamespaceID, strings.ToLower(duplicate.Name))
	r.seenEntities[entityKey] = append(r.seenEntities[entityKey], duplicate)
	return errDuplicateIdentityName
}

func (r *duplicateReportingErrorResolver) ResolveGroups(ctx context.Context, existing, duplicate *identity.Group) error {
	groupKey := fmt.Sprintf("%s/%s", duplicate.NamespaceID, strings.ToLower(duplicate.Name))
	r.seenGroups[groupKey] = append(r.seenGroups[groupKey], duplicate)
	return errDuplicateIdentityName
}

func (r *duplicateReportingErrorResolver) ResolveAliases(ctx context.Context, parent *identity.Entity, existing, duplicate *identity.Alias) error {
	aliasKey := fmt.Sprintf("%s/%s", duplicate.MountAccessor, strings.ToLower(duplicate.Name))
	if duplicate.Local {
		r.seenLocalAliases[aliasKey] = append(r.seenLocalAliases[aliasKey], duplicate)
	} else {
		r.seenAliases[aliasKey] = append(r.seenAliases[aliasKey], duplicate)
	}
	return errDuplicateIdentityName
}

type identityDuplicateReportEntry struct {
	artifactType   string
	scope          string
	name           string
	id             string
	canonicalID    string
	resolutionHint string
	index          int // we care about preserving load order in reporting
	numOthers      int
}

type identityDuplicateReport struct {
	entities                []identityDuplicateReportEntry
	groups                  []identityDuplicateReportEntry
	aliases                 []identityDuplicateReportEntry
	localAliases            []identityDuplicateReportEntry
	numEntityDuplicates     int
	numGroupDuplicates      int
	numAliasDuplicates      int
	numLocalAliasDuplicates int
}

func (r *identityDuplicateReportEntry) Description() string {
	scopeField := "namespace ID"
	if r.artifactType == "entity-alias" || r.artifactType == "local entity-alias" {
		scopeField = "mount accessor"
	}
	return fmt.Sprintf("%s %q with %s %q duplicates %d others",
		r.artifactType, r.name, scopeField, r.scope, r.numOthers)
}

// Labels returns metadata pairs suitable for passing to a logger each slice
// element corresponds alternately to a key and then a value.
func (r *identityDuplicateReportEntry) Labels() []interface{} {
	args := []interface{}{"id", r.id}
	if r.canonicalID != "" {
		args = append(args, "canonical_id")
		args = append(args, r.canonicalID)
	}
	if r.resolutionHint != "" {
		args = append(args, "force_deduplication")
		args = append(args, r.resolutionHint)
	}
	return args
}

func (r *duplicateReportingErrorResolver) Report() identityDuplicateReport {
	var report identityDuplicateReport

	for _, entities := range r.seenEntities {
		if len(entities) <= 1 {
			// Fast path, skip non-duplicates
			continue
		}
		report.numEntityDuplicates++
		// We don't care if it's an exact match or not for entities since we'll
		// rename in either case when we force a de-dupe.
		for idx, entity := range entities {
			r := identityDuplicateReportEntry{
				artifactType: "entity",
				scope:        entity.NamespaceID,
				name:         entity.Name,
				id:           entity.ID,
				index:        idx,
				numOthers:    len(entities) - 1,
			}
			if idx < len(entities)-1 {
				r.resolutionHint = fmt.Sprintf("would rename to %s-%s", entity.Name, entity.ID)
			} else {
				r.resolutionHint = "would not rename"
			}
			report.entities = append(report.entities, r)
		}
	}
	sortReportEntries(report.entities)

	for _, groups := range r.seenGroups {
		if len(groups) <= 1 {
			// Fast path, skip non-duplicates
			continue
		}
		report.numGroupDuplicates++
		// We don't care if it's an exact match or not for groups since we'll
		// rename in either case when we force a de-dupe.
		for idx, group := range groups {
			r := identityDuplicateReportEntry{
				artifactType: "group",
				scope:        group.NamespaceID,
				name:         group.Name,
				id:           group.ID,
				index:        idx,
				numOthers:    len(groups) - 1,
			}
			if idx < len(groups)-1 {
				r.resolutionHint = fmt.Sprintf("would rename to %s-%s", group.Name, group.ID)
			} else {
				r.resolutionHint = "would not rename"
			}
			report.groups = append(report.groups, r)
		}
	}
	sortReportEntries(report.groups)

	reportAliases(&report, r.seenAliases, false)
	reportAliases(&report, r.seenLocalAliases, true)

	return report
}

func reportAliases(report *identityDuplicateReport, seen map[string][]*identity.Alias, local bool) {
	artType := "entity-alias"
	if local {
		artType = "local entity-alias"
	}
	for _, aliases := range seen {
		if len(aliases) <= 1 {
			// Fast path, skip non-duplicates
			continue
		}
		if local {
			report.numLocalAliasDuplicates++
		} else {
			report.numAliasDuplicates++
		}
		// We can't have exact match duplicated for aliases at this point because
		// the would have been merged during load. These are different-case
		// duplicates that must be handled.
		for idx, alias := range aliases {
			r := identityDuplicateReportEntry{
				artifactType: artType,
				scope:        alias.MountAccessor,
				name:         alias.Name,
				id:           alias.ID,
				canonicalID:  alias.CanonicalID,
				index:        idx,
				numOthers:    len(aliases) - 1,
			}
			if idx < len(aliases)-1 {
				r.resolutionHint = fmt.Sprintf("would merge into entity %s", aliases[len(aliases)-1].CanonicalID)
			} else {
				r.resolutionHint = "would merge others into this entity"
			}
			if local {
				report.localAliases = append(report.localAliases, r)
			} else {
				report.aliases = append(report.aliases, r)
			}
		}
	}
	sortReportEntries(report.aliases)
}

func sortReportEntries(es []identityDuplicateReportEntry) {
	sort.Slice(es, func(i, j int) bool {
		a, b := es[i], es[j]
		if a.scope != b.scope {
			return a.scope < b.scope
		}
		aName, bName := strings.ToLower(a.name), strings.ToLower(b.name)
		if aName != bName {
			return aName < bName
		}
		return a.index < b.index
	})
}

// Warner is a subset of hclog.Logger that only has the Warn method to make
// testing simpler.
type Warner interface {
	Warn(msg string, args ...interface{})
}

// TODO set this correctly.
const identityDuplicateReportUrl = "https://developer.hashicorp.com/vault/docs/upgrading/identity-deduplication"

func (r *duplicateReportingErrorResolver) LogReport(log Warner) {
	report := r.Report()

	if report.numEntityDuplicates == 0 && report.numGroupDuplicates == 0 && report.numAliasDuplicates == 0 {
		return
	}

	log.Warn("DUPLICATES DETECTED, see following logs for details and refer to " +
		identityDuplicateReportUrl + " for resolution.")

	// Aliases first since they are most critical to resolve. Local first because
	// all the rest can be ignored on a perf secondary.
	if len(report.localAliases) > 0 {
		log.Warn(fmt.Sprintf("%d different-case local entity alias duplicates found (potential security risk)", report.numLocalAliasDuplicates))
		for _, e := range report.localAliases {
			log.Warn(e.Description(), e.Labels()...)
		}
		log.Warn("end of different-case local entity-alias duplicates")
	}
	if len(report.aliases) > 0 {
		log.Warn(fmt.Sprintf("%d different-case entity alias duplicates found (potential security risk)", report.numAliasDuplicates))
		for _, e := range report.aliases {
			log.Warn(e.Description(), e.Labels()...)
		}
		log.Warn("end of different-case entity-alias duplicates")
	}

	if len(report.entities) > 0 {
		log.Warn(fmt.Sprintf("%d entity duplicates found", report.numEntityDuplicates))
		for _, e := range report.entities {
			log.Warn(e.Description(), e.Labels()...)
		}
		log.Warn("end of entity duplicates")
	}

	if len(report.groups) > 0 {
		log.Warn(fmt.Sprintf("%d group duplicates found", report.numGroupDuplicates))
		for _, e := range report.groups {
			log.Warn(e.Description(), e.Labels()...)
		}
		log.Warn("end of group duplicates")
	}
	log.Warn("end of identity duplicate report, refer to " +
		identityDuplicateReportUrl + " for resolution.")
}
