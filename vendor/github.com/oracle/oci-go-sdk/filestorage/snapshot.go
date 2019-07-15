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

// Snapshot A point-in-time snapshot of a specified file system.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type Snapshot struct {

	// The OCID of the file system from which the snapshot
	// was created.
	FileSystemId *string `mandatory:"true" json:"fileSystemId"`

	// The OCID of the snapshot.
	Id *string `mandatory:"true" json:"id"`

	// The current state of the snapshot.
	LifecycleState SnapshotLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Name of the snapshot. This value is immutable.
	// Avoid entering confidential information.
	// Example: `Sunday`
	Name *string `mandatory:"true" json:"name"`

	// The date and time the snapshot was created, expressed
	// in RFC 3339 (https://tools.ietf.org/rfc/rfc3339) timestamp format.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

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

func (m Snapshot) String() string {
	return common.PointerString(m)
}

// SnapshotLifecycleStateEnum Enum with underlying type: string
type SnapshotLifecycleStateEnum string

// Set of constants representing the allowable values for SnapshotLifecycleStateEnum
const (
	SnapshotLifecycleStateCreating SnapshotLifecycleStateEnum = "CREATING"
	SnapshotLifecycleStateActive   SnapshotLifecycleStateEnum = "ACTIVE"
	SnapshotLifecycleStateDeleting SnapshotLifecycleStateEnum = "DELETING"
	SnapshotLifecycleStateDeleted  SnapshotLifecycleStateEnum = "DELETED"
)

var mappingSnapshotLifecycleState = map[string]SnapshotLifecycleStateEnum{
	"CREATING": SnapshotLifecycleStateCreating,
	"ACTIVE":   SnapshotLifecycleStateActive,
	"DELETING": SnapshotLifecycleStateDeleting,
	"DELETED":  SnapshotLifecycleStateDeleted,
}

// GetSnapshotLifecycleStateEnumValues Enumerates the set of values for SnapshotLifecycleStateEnum
func GetSnapshotLifecycleStateEnumValues() []SnapshotLifecycleStateEnum {
	values := make([]SnapshotLifecycleStateEnum, 0)
	for _, v := range mappingSnapshotLifecycleState {
		values = append(values, v)
	}
	return values
}
