// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package aws

const (

	// Root config observations
	// These observations have rotation_period, rotation_schedule,
	// rotation_window, and disable_automatic_rotation metadata

	ObservationTypeAWSRootConfigWrite  = "aws/config/root/write"
	ObservationTypeAWSRootConfigRead   = "aws/config/root/read"
	ObservationTypeAWSRootConfigRotate = "aws/config/root/rotate"

	// Lease config observations
	// These observations have lease and lease_max durations

	ObservationTypeAWSLeaseConfigWrite = "aws/config/lease/write"
	ObservationTypeAWSLeaseConfigRead  = "aws/config/lease/read"

	// Role related observations
	// These observations have role_name and credentials_type metadata

	ObservationTypeAWSRoleWrite = "aws/role/write"
	ObservationTypeAWSRoleRead  = "aws/role/read"

	// Role delete observation
	// This observation only has role_name metadata to avoid doing a storage read

	ObservationTypeAWSRoleDelete = "aws/role/delete"

	// Static role related observations
	// These observations have role_name and rotation_period metadata

	ObservationTypeAWSStaticRoleWrite  = "aws/static-role/write"
	ObservationTypeAWSStaticRoleRead   = "aws/static-role/read"
	ObservationTypeAWSStaticRoleDelete = "aws/static-role/delete"

	// Credential related observations
	// These observations have role_name, credentials_type, ttl, max_ttl, and is_sts metadata

	ObservationTypeAWSCredentialCreateSuccess = "aws/credential/create/success"
	ObservationTypeAWSCredentialCreateFail    = "aws/credential/create/fail"

	// Secret lifecycle observations
	// These observations don't have access to the role_name
	// They only have is_sts and credentials_type metadata

	ObservationTypeAWSCredentialRenew  = "aws/credential/renew"
	ObservationTypeAWSCredentialRevoke = "aws/credential/revoke"

	// Static credential related observations
	// These observations have role_name

	ObservationTypeAWSStaticCredentialRead   = "aws/static-credential/read"
	ObservationTypeAWSStaticCredentialRotate = "aws/static-credential/rotate"
)
