// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package observations

const (
	ObservationTypeLeaseCreationAuth    = "lease/create/auth"
	ObservationTypeLeaseCreationNonAuth = "lease/create/non-auth"
	ObservationTypeLeaseRenewAuth       = "lease/renew/auth"
	ObservationTypeLeaseRenewNonAuth    = "lease/renew/non-auth"
	ObservationTypeLeaseLazyRevoke      = "lease/lazy-revoke"
	ObservationTypeLeaseRevocation      = "lease/revoke"
	ObservationTypePolicyACLEvaluation  = "policy/acl/evaluation"
)
