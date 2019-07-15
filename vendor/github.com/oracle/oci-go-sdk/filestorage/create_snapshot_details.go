// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// File Storage Service API
//
// The API for the File Storage Service.
//

package filestorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// CreateSnapshotDetails Details for creating the snapshot.
type CreateSnapshotDetails struct {

	// The OCID of the file system to take a snapshot of.
	FileSystemId *string `mandatory:"true" json:"fileSystemId"`

	// Name of the snapshot. This value is immutable. It must also be unique with respect
	// to all other non-DELETED snapshots on the associated file
	// system.
	// Avoid entering confidential information.
	// Example: `Sunday`
	Name *string `mandatory:"true" json:"name"`

	// Free-form tags for this resource. Each tag is a simple key-value pair
	//  with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m CreateSnapshotDetails) String() string {
	return common.PointerString(m)
}
