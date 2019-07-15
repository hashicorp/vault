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

// AutonomousDatabaseBackupSummary An Autonomous Database backup.
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized, talk to an administrator. If you're an administrator who needs to write policies to give users access, see Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type AutonomousDatabaseBackupSummary struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Autonomous Database backup.
	Id *string `mandatory:"true" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Autonomous Database.
	AutonomousDatabaseId *string `mandatory:"true" json:"autonomousDatabaseId"`

	// The user-friendly name for the backup. The name does not have to be unique.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The type of backup.
	Type AutonomousDatabaseBackupSummaryTypeEnum `mandatory:"true" json:"type"`

	// Indicates whether the backup is user-initiated or automatic.
	IsAutomatic *bool `mandatory:"true" json:"isAutomatic"`

	// The current state of the backup.
	LifecycleState AutonomousDatabaseBackupSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the backup started.
	TimeStarted *common.SDKTime `mandatory:"false" json:"timeStarted"`

	// The date and time the backup completed.
	TimeEnded *common.SDKTime `mandatory:"false" json:"timeEnded"`

	// Additional information about the current lifecycle state.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// The size of the database in terabytes at the time the backup was taken.
	DatabaseSizeInTBs *float32 `mandatory:"false" json:"databaseSizeInTBs"`
}

func (m AutonomousDatabaseBackupSummary) String() string {
	return common.PointerString(m)
}

// AutonomousDatabaseBackupSummaryTypeEnum Enum with underlying type: string
type AutonomousDatabaseBackupSummaryTypeEnum string

// Set of constants representing the allowable values for AutonomousDatabaseBackupSummaryTypeEnum
const (
	AutonomousDatabaseBackupSummaryTypeIncremental AutonomousDatabaseBackupSummaryTypeEnum = "INCREMENTAL"
	AutonomousDatabaseBackupSummaryTypeFull        AutonomousDatabaseBackupSummaryTypeEnum = "FULL"
)

var mappingAutonomousDatabaseBackupSummaryType = map[string]AutonomousDatabaseBackupSummaryTypeEnum{
	"INCREMENTAL": AutonomousDatabaseBackupSummaryTypeIncremental,
	"FULL":        AutonomousDatabaseBackupSummaryTypeFull,
}

// GetAutonomousDatabaseBackupSummaryTypeEnumValues Enumerates the set of values for AutonomousDatabaseBackupSummaryTypeEnum
func GetAutonomousDatabaseBackupSummaryTypeEnumValues() []AutonomousDatabaseBackupSummaryTypeEnum {
	values := make([]AutonomousDatabaseBackupSummaryTypeEnum, 0)
	for _, v := range mappingAutonomousDatabaseBackupSummaryType {
		values = append(values, v)
	}
	return values
}

// AutonomousDatabaseBackupSummaryLifecycleStateEnum Enum with underlying type: string
type AutonomousDatabaseBackupSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for AutonomousDatabaseBackupSummaryLifecycleStateEnum
const (
	AutonomousDatabaseBackupSummaryLifecycleStateCreating AutonomousDatabaseBackupSummaryLifecycleStateEnum = "CREATING"
	AutonomousDatabaseBackupSummaryLifecycleStateActive   AutonomousDatabaseBackupSummaryLifecycleStateEnum = "ACTIVE"
	AutonomousDatabaseBackupSummaryLifecycleStateDeleting AutonomousDatabaseBackupSummaryLifecycleStateEnum = "DELETING"
	AutonomousDatabaseBackupSummaryLifecycleStateDeleted  AutonomousDatabaseBackupSummaryLifecycleStateEnum = "DELETED"
	AutonomousDatabaseBackupSummaryLifecycleStateFailed   AutonomousDatabaseBackupSummaryLifecycleStateEnum = "FAILED"
)

var mappingAutonomousDatabaseBackupSummaryLifecycleState = map[string]AutonomousDatabaseBackupSummaryLifecycleStateEnum{
	"CREATING": AutonomousDatabaseBackupSummaryLifecycleStateCreating,
	"ACTIVE":   AutonomousDatabaseBackupSummaryLifecycleStateActive,
	"DELETING": AutonomousDatabaseBackupSummaryLifecycleStateDeleting,
	"DELETED":  AutonomousDatabaseBackupSummaryLifecycleStateDeleted,
	"FAILED":   AutonomousDatabaseBackupSummaryLifecycleStateFailed,
}

// GetAutonomousDatabaseBackupSummaryLifecycleStateEnumValues Enumerates the set of values for AutonomousDatabaseBackupSummaryLifecycleStateEnum
func GetAutonomousDatabaseBackupSummaryLifecycleStateEnumValues() []AutonomousDatabaseBackupSummaryLifecycleStateEnum {
	values := make([]AutonomousDatabaseBackupSummaryLifecycleStateEnum, 0)
	for _, v := range mappingAutonomousDatabaseBackupSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
