// Copyright (c) HashiCorp, Inc.
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
)
