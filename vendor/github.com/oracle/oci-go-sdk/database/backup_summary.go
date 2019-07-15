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

// BackupSummary A database backup.
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized, talk to an administrator. If you're an administrator who needs to write policies to give users access, see Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type BackupSummary struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the backup.
	Id *string `mandatory:"false" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the database.
	DatabaseId *string `mandatory:"false" json:"databaseId"`

	// The user-friendly name for the backup. The name does not have to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The type of backup.
	Type BackupSummaryTypeEnum `mandatory:"false" json:"type,omitempty"`

	// The date and time the backup started.
	TimeStarted *common.SDKTime `mandatory:"false" json:"timeStarted"`

	// The date and time the backup was completed.
	TimeEnded *common.SDKTime `mandatory:"false" json:"timeEnded"`

	// Additional information about the current lifecycleState.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// The name of the availability domain where the database backup is stored.
	AvailabilityDomain *string `mandatory:"false" json:"availabilityDomain"`

	// The current state of the backup.
	LifecycleState BackupSummaryLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The Oracle Database edition of the DB system from which the database backup was taken.
	DatabaseEdition BackupSummaryDatabaseEditionEnum `mandatory:"false" json:"databaseEdition,omitempty"`

	// The size of the database in gigabytes at the time the backup was taken.
	DatabaseSizeInGBs *float64 `mandatory:"false" json:"databaseSizeInGBs"`
}

func (m BackupSummary) String() string {
	return common.PointerString(m)
}

// BackupSummaryTypeEnum Enum with underlying type: string
type BackupSummaryTypeEnum string

// Set of constants representing the allowable values for BackupSummaryTypeEnum
const (
	BackupSummaryTypeIncremental BackupSummaryTypeEnum = "INCREMENTAL"
	BackupSummaryTypeFull        BackupSummaryTypeEnum = "FULL"
	BackupSummaryTypeVirtualFull BackupSummaryTypeEnum = "VIRTUAL_FULL"
)

var mappingBackupSummaryType = map[string]BackupSummaryTypeEnum{
	"INCREMENTAL":  BackupSummaryTypeIncremental,
	"FULL":         BackupSummaryTypeFull,
	"VIRTUAL_FULL": BackupSummaryTypeVirtualFull,
}

// GetBackupSummaryTypeEnumValues Enumerates the set of values for BackupSummaryTypeEnum
func GetBackupSummaryTypeEnumValues() []BackupSummaryTypeEnum {
	values := make([]BackupSummaryTypeEnum, 0)
	for _, v := range mappingBackupSummaryType {
		values = append(values, v)
	}
	return values
}

// BackupSummaryLifecycleStateEnum Enum with underlying type: string
type BackupSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for BackupSummaryLifecycleStateEnum
const (
	BackupSummaryLifecycleStateCreating  BackupSummaryLifecycleStateEnum = "CREATING"
	BackupSummaryLifecycleStateActive    BackupSummaryLifecycleStateEnum = "ACTIVE"
	BackupSummaryLifecycleStateDeleting  BackupSummaryLifecycleStateEnum = "DELETING"
	BackupSummaryLifecycleStateDeleted   BackupSummaryLifecycleStateEnum = "DELETED"
	BackupSummaryLifecycleStateFailed    BackupSummaryLifecycleStateEnum = "FAILED"
	BackupSummaryLifecycleStateRestoring BackupSummaryLifecycleStateEnum = "RESTORING"
)

var mappingBackupSummaryLifecycleState = map[string]BackupSummaryLifecycleStateEnum{
	"CREATING":  BackupSummaryLifecycleStateCreating,
	"ACTIVE":    BackupSummaryLifecycleStateActive,
	"DELETING":  BackupSummaryLifecycleStateDeleting,
	"DELETED":   BackupSummaryLifecycleStateDeleted,
	"FAILED":    BackupSummaryLifecycleStateFailed,
	"RESTORING": BackupSummaryLifecycleStateRestoring,
}

// GetBackupSummaryLifecycleStateEnumValues Enumerates the set of values for BackupSummaryLifecycleStateEnum
func GetBackupSummaryLifecycleStateEnumValues() []BackupSummaryLifecycleStateEnum {
	values := make([]BackupSummaryLifecycleStateEnum, 0)
	for _, v := range mappingBackupSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}

// BackupSummaryDatabaseEditionEnum Enum with underlying type: string
type BackupSummaryDatabaseEditionEnum string

// Set of constants representing the allowable values for BackupSummaryDatabaseEditionEnum
const (
	BackupSummaryDatabaseEditionStandardEdition                     BackupSummaryDatabaseEditionEnum = "STANDARD_EDITION"
	BackupSummaryDatabaseEditionEnterpriseEdition                   BackupSummaryDatabaseEditionEnum = "ENTERPRISE_EDITION"
	BackupSummaryDatabaseEditionEnterpriseEditionHighPerformance    BackupSummaryDatabaseEditionEnum = "ENTERPRISE_EDITION_HIGH_PERFORMANCE"
	BackupSummaryDatabaseEditionEnterpriseEditionExtremePerformance BackupSummaryDatabaseEditionEnum = "ENTERPRISE_EDITION_EXTREME_PERFORMANCE"
)

var mappingBackupSummaryDatabaseEdition = map[string]BackupSummaryDatabaseEditionEnum{
	"STANDARD_EDITION":                       BackupSummaryDatabaseEditionStandardEdition,
	"ENTERPRISE_EDITION":                     BackupSummaryDatabaseEditionEnterpriseEdition,
	"ENTERPRISE_EDITION_HIGH_PERFORMANCE":    BackupSummaryDatabaseEditionEnterpriseEditionHighPerformance,
	"ENTERPRISE_EDITION_EXTREME_PERFORMANCE": BackupSummaryDatabaseEditionEnterpriseEditionExtremePerformance,
}

// GetBackupSummaryDatabaseEditionEnumValues Enumerates the set of values for BackupSummaryDatabaseEditionEnum
func GetBackupSummaryDatabaseEditionEnumValues() []BackupSummaryDatabaseEditionEnum {
	values := make([]BackupSummaryDatabaseEditionEnum, 0)
	for _, v := range mappingBackupSummaryDatabaseEdition {
		values = append(values, v)
	}
	return values
}
