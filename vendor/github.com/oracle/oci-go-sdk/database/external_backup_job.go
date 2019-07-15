// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ExternalBackupJob Provides all the details that apply to an external backup job.
type ExternalBackupJob struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the associated backup resource.
	BackupId *string `mandatory:"true" json:"backupId"`

	// An indicator for the provisioning state of the resource. If `TRUE`, the resource is still being provisioned.
	Provisioning *bool `mandatory:"true" json:"provisioning"`

	// The Swift path to use as a destination for the standalone backup.
	SwiftPath *string `mandatory:"true" json:"swiftPath"`

	// The name of the Swift compartment bucket where the backup should be stored.
	BucketName *string `mandatory:"true" json:"bucketName"`

	// The tag for RMAN to apply to the backup.
	Tag *string `mandatory:"true" json:"tag"`

	// The Swift user name to use for transferring the standalone backup to the designated Swift compartment bucket.
	UserName *string `mandatory:"true" json:"userName"`

	// The auth token to use for access to the Swift compartment bucket that will store the standalone backup.
	// For information about auth tokens, see Working with Auth Tokens (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingcredentials.htm#two).
	SwiftPassword *string `mandatory:"false" json:"swiftPassword"`
}

func (m ExternalBackupJob) String() string {
	return common.PointerString(m)
}
