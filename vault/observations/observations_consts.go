// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package observations

const (
	// ObservationTypeLeaseCreationAuth is emitted when a lease connected to auth is created
	ObservationTypeLeaseCreationAuth = "lease/create/auth"
	// ObservationTypeLeaseCreationNonAuth is emitted when a lease NOT connected to auth is created
	ObservationTypeLeaseCreationNonAuth = "lease/create/non-auth"
	// ObservationTypeLeaseRenewAuth is emitted when a lease connected to auth is renewed
	ObservationTypeLeaseRenewAuth = "lease/renew/auth"
	// ObservationTypeLeaseRenewNonAuth is emitted when a lease NOT connected to auth is renewed
	ObservationTypeLeaseRenewNonAuth = "lease/renew/non-auth"
	// ObservationTypeLeaseLazyRevoke is emitted when a lease is lazy-revoked
	ObservationTypeLeaseLazyRevoke = "lease/lazy-revoke"
	// ObservationTypeLeaseRevocation is emitted when a lease is revoked
	ObservationTypeLeaseRevocation = "lease/revoke"

	// ObservationTypePolicyACLEvaluation is emitted when an ACL policy is evaluated
	ObservationTypePolicyACLEvaluation = "policy/acl/evaluation"

	// ObservationTypeMountAuthEnable is emitted when an auth mount is enabled
	ObservationTypeMountAuthEnable = "mount/auth/enable"
	// ObservationTypeMountAuthDisable is emitted when an auth mount is disabled
	ObservationTypeMountAuthDisable = "mount/auth/disable"
	// ObservationTypeMountSecretsEnable is emitted when a secret mount is enabled
	ObservationTypeMountSecretsEnable = "mount/secrets/enable"
	// ObservationTypeMountSecretsDisable is emitted when a secret mount is disabled
	ObservationTypeMountSecretsDisable = "mount/secrets/disable"

	// ObservationTypeNamespaceCreate is emitted when a namespace is created
	ObservationTypeNamespaceCreate = "namespace/create"
	// ObservationTypeNamespacePatch is emitted when a namespace is patched
	ObservationTypeNamespacePatch = "namespace/patch"
	// ObservationTypeNamespaceDelete is emitted when a namespace is deleted
	ObservationTypeNamespaceDelete = "namespace/delete"

	// ObservationTypeEntityUpsert is emitted when an entity is upserted
	ObservationTypeEntityUpsert = "identity/entity/upsert"
	// ObservationTypeEntityDelete is emitted when an entity is deleted
	ObservationTypeEntityDelete = "identity/entity/delete"

	// ObservationTypeAliasUpsert is emitted when an alias is upserted.
	// NOTE: Currently we don't allow by-factors modification of group aliases the
	// way we do with entities. Instead, the group itself is updated, not the alias.
	ObservationTypeAliasUpsert = "identity/alias/upsert"
	// ObservationTypeAliasDelete is emitted when an alias is deleted
	ObservationTypeAliasDelete = "identity/alias/delete"

	// ObservationTypeGroupUpsert is emitted when a group is upserted
	ObservationTypeGroupUpsert = "identity/group/upsert"
	// ObservationTypeGroupDelete is emitted when a group is deleted
	ObservationTypeGroupDelete = "identity/group/delete"
)
