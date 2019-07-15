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

// FileSystem An NFS file system. To allow access to a file system, add it
// to an export set and associate the export set with a mount
// target. The same file system can be in multiple export sets and
// associated with multiple mount targets.
// To use any of the API operations, you must be authorized in an
// IAM policy. If you're not authorized, talk to an
// administrator. If you're an administrator who needs to write
// policies to give users access, see Getting Started with
// Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type FileSystem struct {

	// The number of bytes consumed by the file system, including
	// any snapshots. This number reflects the metered size of the file
	// system and is updated asynchronously with respect to
	// updates to the file system.
	MeteredBytes *int64 `mandatory:"true" json:"meteredBytes"`

	// The OCID of the compartment that contains the file system.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	// Example: `My file system`
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID of the file system.
	Id *string `mandatory:"true" json:"id"`

	// The current state of the file system.
	LifecycleState FileSystemLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the file system was created, expressed in
	// RFC 3339 (https://tools.ietf.org/rfc/rfc3339) timestamp format.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The availability domain the file system is in. May be unset
	// as a blank or NULL value.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"false" json:"availabilityDomain"`

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

func (m FileSystem) String() string {
	return common.PointerString(m)
}

// FileSystemLifecycleStateEnum Enum with underlying type: string
type FileSystemLifecycleStateEnum string

// Set of constants representing the allowable values for FileSystemLifecycleStateEnum
const (
	FileSystemLifecycleStateCreating FileSystemLifecycleStateEnum = "CREATING"
	FileSystemLifecycleStateActive   FileSystemLifecycleStateEnum = "ACTIVE"
	FileSystemLifecycleStateDeleting FileSystemLifecycleStateEnum = "DELETING"
	FileSystemLifecycleStateDeleted  FileSystemLifecycleStateEnum = "DELETED"
)

var mappingFileSystemLifecycleState = map[string]FileSystemLifecycleStateEnum{
	"CREATING": FileSystemLifecycleStateCreating,
	"ACTIVE":   FileSystemLifecycleStateActive,
	"DELETING": FileSystemLifecycleStateDeleting,
	"DELETED":  FileSystemLifecycleStateDeleted,
}

// GetFileSystemLifecycleStateEnumValues Enumerates the set of values for FileSystemLifecycleStateEnum
func GetFileSystemLifecycleStateEnumValues() []FileSystemLifecycleStateEnum {
	values := make([]FileSystemLifecycleStateEnum, 0)
	for _, v := range mappingFileSystemLifecycleState {
		values = append(values, v)
	}
	return values
}
