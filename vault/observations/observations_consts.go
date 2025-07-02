// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package observations

const (
	// lease
	ObservationTypeLeaseCreationAuth    = "lease/create/auth"
	ObservationTypeLeaseCreationNonAuth = "lease/create/non-auth"
	ObservationTypeLeaseRenewAuth       = "lease/renew/auth"
	ObservationTypeLeaseRenewNonAuth    = "lease/renew/non-auth"
	ObservationTypeLeaseLazyRevoke      = "lease/lazy-revoke"
	ObservationTypeLeaseRevocation      = "lease/revoke"

	// policy
	ObservationTypePolicyACLEvaluation = "policy/acl/evaluation"

	// mount
	ObservationTypeMountAuthEnable     = "mount/auth/enable"
	ObservationTypeMountAuthDisable    = "mount/auth/disable"
	ObservationTypeMountSecretsEnable  = "mount/secrets/enable"
	ObservationTypeMountSecretsDisable = "mount/secrets/disable"

	// namespace
	ObservationTypeNamespaceCreate = "namespace/create"
	ObservationTypeNamespacePatch  = "namespace/patch"
	ObservationTypeNamespaceDelete = "namespace/delete"
)
