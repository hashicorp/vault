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

// AutonomousDatabaseBackup An Autonomous Database backup.
type AutonomousDatabaseBackup struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Autonomous Database backup.
	Id *string `mandatory:"true" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Autonomous Database.
	AutonomousDatabaseId *string `mandatory:"true" json:"autonomousDatabaseId"`

	// The user-friendly name for the backup. The name does not have to be unique.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The type of backup.
	Type AutonomousDatabaseBackupTypeEnum `mandatory:"true" json:"type"`

	// Indicates whether the backup is user-initiated or automatic.
	IsAutomatic *bool `mandatory:"true" json:"isAutomatic"`

	// The current state of the backup.
	LifecycleState AutonomousDatabaseBackupLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the backup started.
	TimeStarted *common.SDKTime `mandatory:"false" json:"timeStarted"`

	// The date and time the backup completed.
	TimeEnded *common.SDKTime `mandatory:"false" json:"timeEnded"`

	// Additional information about the current lifecycle state.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// The size of the database in terabytes at the time the backup was taken.
	DatabaseSizeInTBs *float32 `mandatory:"false" json:"databaseSizeInTBs"`
}

func (m AutonomousDatabaseBackup) String() string {
	return common.PointerString(m)
}

// AutonomousDatabaseBackupTypeEnum Enum with underlying type: string
type AutonomousDatabaseBackupTypeEnum string

// Set of constants representing the allowable values for AutonomousDatabaseBackupTypeEnum
const (
	AutonomousDatabaseBackupTypeIncremental AutonomousDatabaseBackupTypeEnum = "INCREMENTAL"
	AutonomousDatabaseBackupTypeFull        AutonomousDatabaseBackupTypeEnum = "FULL"
)

var mappingAutonomousDatabaseBackupType = map[string]AutonomousDatabaseBackupTypeEnum{
	"INCREMENTAL": AutonomousDatabaseBackupTypeIncremental,
	"FULL":        AutonomousDatabaseBackupTypeFull,
}

// GetAutonomousDatabaseBackupTypeEnumValues Enumerates the set of values for AutonomousDatabaseBackupTypeEnum
func GetAutonomousDatabaseBackupTypeEnumValues() []AutonomousDatabaseBackupTypeEnum {
	values := make([]AutonomousDatabaseBackupTypeEnum, 0)
	for _, v := range mappingAutonomousDatabaseBackupType {
		values = append(values, v)
	}
	return values
}

// AutonomousDatabaseBackupLifecycleStateEnum Enum with underlying type: string
type AutonomousDatabaseBackupLifecycleStateEnum string

// Set of constants representing the allowable values for AutonomousDatabaseBackupLifecycleStateEnum
const (
	AutonomousDatabaseBackupLifecycleStateCreating AutonomousDatabaseBackupLifecycleStateEnum = "CREATING"
	AutonomousDatabaseBackupLifecycleStateActive   AutonomousDatabaseBackupLifecycleStateEnum = "ACTIVE"
	AutonomousDatabaseBackupLifecycleStateDeleting AutonomousDatabaseBackupLifecycleStateEnum = "DELETING"
	AutonomousDatabaseBackupLifecycleStateDeleted  AutonomousDatabaseBackupLifecycleStateEnum = "DELETED"
	AutonomousDatabaseBackupLifecycleStateFailed   AutonomousDatabaseBackupLifecycleStateEnum = "FAILED"
)

var mappingAutonomousDatabaseBackupLifecycleState = map[string]AutonomousDatabaseBackupLifecycleStateEnum{
	"CREATING": AutonomousDatabaseBackupLifecycleStateCreating,
	"ACTIVE":   AutonomousDatabaseBackupLifecycleStateActive,
	"DELETING": AutonomousDatabaseBackupLifecycleStateDeleting,
	"DELETED":  AutonomousDatabaseBackupLifecycleStateDeleted,
	"FAILED":   AutonomousDatabaseBackupLifecycleStateFailed,
}

// GetAutonomousDatabaseBackupLifecycleStateEnumValues Enumerates the set of values for AutonomousDatabaseBackupLifecycleStateEnum
func GetAutonomousDatabaseBackupLifecycleStateEnumValues() []AutonomousDatabaseBackupLifecycleStateEnum {
	values := make([]AutonomousDatabaseBackupLifecycleStateEnum, 0)
	for _, v := range mappingAutonomousDatabaseBackupLifecycleState {
		values = append(values, v)
	}
	return values
}
