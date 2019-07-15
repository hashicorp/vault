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

// ExportSummary Summary information for an export.
type ExportSummary struct {

	// The OCID of this export's export set.
	ExportSetId *string `mandatory:"true" json:"exportSetId"`

	// The OCID of this export's file system.
	FileSystemId *string `mandatory:"true" json:"fileSystemId"`

	// The OCID of this export.
	Id *string `mandatory:"true" json:"id"`

	// The current state of this export.
	LifecycleState ExportSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Path used to access the associated file system.
	// Avoid entering confidential information.
	// Example: `/mediafiles`
	Path *string `mandatory:"true" json:"path"`

	// The date and time the export was created, expressed
	// in RFC 3339 (https://tools.ietf.org/rfc/rfc3339) timestamp format.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`
}

func (m ExportSummary) String() string {
	return common.PointerString(m)
}

// ExportSummaryLifecycleStateEnum Enum with underlying type: string
type ExportSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for ExportSummaryLifecycleStateEnum
const (
	ExportSummaryLifecycleStateCreating ExportSummaryLifecycleStateEnum = "CREATING"
	ExportSummaryLifecycleStateActive   ExportSummaryLifecycleStateEnum = "ACTIVE"
	ExportSummaryLifecycleStateDeleting ExportSummaryLifecycleStateEnum = "DELETING"
	ExportSummaryLifecycleStateDeleted  ExportSummaryLifecycleStateEnum = "DELETED"
)

var mappingExportSummaryLifecycleState = map[string]ExportSummaryLifecycleStateEnum{
	"CREATING": ExportSummaryLifecycleStateCreating,
	"ACTIVE":   ExportSummaryLifecycleStateActive,
	"DELETING": ExportSummaryLifecycleStateDeleting,
	"DELETED":  ExportSummaryLifecycleStateDeleted,
}

// GetExportSummaryLifecycleStateEnumValues Enumerates the set of values for ExportSummaryLifecycleStateEnum
func GetExportSummaryLifecycleStateEnumValues() []ExportSummaryLifecycleStateEnum {
	values := make([]ExportSummaryLifecycleStateEnum, 0)
	for _, v := range mappingExportSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
