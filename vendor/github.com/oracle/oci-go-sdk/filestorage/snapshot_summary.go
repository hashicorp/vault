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

// SnapshotSummary Summary information for a snapshot.
type SnapshotSummary struct {

	// The OCID of the file system from which the
	// snapshot was created.
	FileSystemId *string `mandatory:"true" json:"fileSystemId"`

	// The OCID of the snapshot.
	Id *string `mandatory:"true" json:"id"`

	// The current state of the snapshot.
	LifecycleState SnapshotSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

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

func (m SnapshotSummary) String() string {
	return common.PointerString(m)
}

// SnapshotSummaryLifecycleStateEnum Enum with underlying type: string
type SnapshotSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for SnapshotSummaryLifecycleStateEnum
const (
	SnapshotSummaryLifecycleStateCreating SnapshotSummaryLifecycleStateEnum = "CREATING"
	SnapshotSummaryLifecycleStateActive   SnapshotSummaryLifecycleStateEnum = "ACTIVE"
	SnapshotSummaryLifecycleStateDeleting SnapshotSummaryLifecycleStateEnum = "DELETING"
	SnapshotSummaryLifecycleStateDeleted  SnapshotSummaryLifecycleStateEnum = "DELETED"
)

var mappingSnapshotSummaryLifecycleState = map[string]SnapshotSummaryLifecycleStateEnum{
	"CREATING": SnapshotSummaryLifecycleStateCreating,
	"ACTIVE":   SnapshotSummaryLifecycleStateActive,
	"DELETING": SnapshotSummaryLifecycleStateDeleting,
	"DELETED":  SnapshotSummaryLifecycleStateDeleted,
}

// GetSnapshotSummaryLifecycleStateEnumValues Enumerates the set of values for SnapshotSummaryLifecycleStateEnum
func GetSnapshotSummaryLifecycleStateEnumValues() []SnapshotSummaryLifecycleStateEnum {
	values := make([]SnapshotSummaryLifecycleStateEnum, 0)
	for _, v := range mappingSnapshotSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
