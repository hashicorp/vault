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

// AutonomousDataWarehouseBackup **Deprecated.** See AutonomousDatabaseBackup for reference information about Autonomous Data Warehouse backups.
type AutonomousDataWarehouseBackup struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Autonomous Data Warehouse backup.
	Id *string `mandatory:"true" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Autonomous Data Warehouse.
	AutonomousDataWarehouseId *string `mandatory:"true" json:"autonomousDataWarehouseId"`

	// The user-friendly name for the backup. The name does not have to be unique.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The type of backup.
	Type AutonomousDataWarehouseBackupTypeEnum `mandatory:"true" json:"type"`

	// Indicates whether the backup is user-initiated or automatic.
	IsAutomatic *bool `mandatory:"true" json:"isAutomatic"`

	// The current state of the backup.
	LifecycleState AutonomousDataWarehouseBackupLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the backup started.
	TimeStarted *common.SDKTime `mandatory:"false" json:"timeStarted"`

	// The date and time the backup completed.
	TimeEnded *common.SDKTime `mandatory:"false" json:"timeEnded"`

	// Additional information about the current lifecycle state.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`
}

func (m AutonomousDataWarehouseBackup) String() string {
	return common.PointerString(m)
}

// AutonomousDataWarehouseBackupTypeEnum Enum with underlying type: string
type AutonomousDataWarehouseBackupTypeEnum string

// Set of constants representing the allowable values for AutonomousDataWarehouseBackupTypeEnum
const (
	AutonomousDataWarehouseBackupTypeIncremental AutonomousDataWarehouseBackupTypeEnum = "INCREMENTAL"
	AutonomousDataWarehouseBackupTypeFull        AutonomousDataWarehouseBackupTypeEnum = "FULL"
)

var mappingAutonomousDataWarehouseBackupType = map[string]AutonomousDataWarehouseBackupTypeEnum{
	"INCREMENTAL": AutonomousDataWarehouseBackupTypeIncremental,
	"FULL":        AutonomousDataWarehouseBackupTypeFull,
}

// GetAutonomousDataWarehouseBackupTypeEnumValues Enumerates the set of values for AutonomousDataWarehouseBackupTypeEnum
func GetAutonomousDataWarehouseBackupTypeEnumValues() []AutonomousDataWarehouseBackupTypeEnum {
	values := make([]AutonomousDataWarehouseBackupTypeEnum, 0)
	for _, v := range mappingAutonomousDataWarehouseBackupType {
		values = append(values, v)
	}
	return values
}

// AutonomousDataWarehouseBackupLifecycleStateEnum Enum with underlying type: string
type AutonomousDataWarehouseBackupLifecycleStateEnum string

// Set of constants representing the allowable values for AutonomousDataWarehouseBackupLifecycleStateEnum
const (
	AutonomousDataWarehouseBackupLifecycleStateCreating AutonomousDataWarehouseBackupLifecycleStateEnum = "CREATING"
	AutonomousDataWarehouseBackupLifecycleStateActive   AutonomousDataWarehouseBackupLifecycleStateEnum = "ACTIVE"
	AutonomousDataWarehouseBackupLifecycleStateDeleting AutonomousDataWarehouseBackupLifecycleStateEnum = "DELETING"
	AutonomousDataWarehouseBackupLifecycleStateDeleted  AutonomousDataWarehouseBackupLifecycleStateEnum = "DELETED"
	AutonomousDataWarehouseBackupLifecycleStateFailed   AutonomousDataWarehouseBackupLifecycleStateEnum = "FAILED"
)

var mappingAutonomousDataWarehouseBackupLifecycleState = map[string]AutonomousDataWarehouseBackupLifecycleStateEnum{
	"CREATING": AutonomousDataWarehouseBackupLifecycleStateCreating,
	"ACTIVE":   AutonomousDataWarehouseBackupLifecycleStateActive,
	"DELETING": AutonomousDataWarehouseBackupLifecycleStateDeleting,
	"DELETED":  AutonomousDataWarehouseBackupLifecycleStateDeleted,
	"FAILED":   AutonomousDataWarehouseBackupLifecycleStateFailed,
}

// GetAutonomousDataWarehouseBackupLifecycleStateEnumValues Enumerates the set of values for AutonomousDataWarehouseBackupLifecycleStateEnum
func GetAutonomousDataWarehouseBackupLifecycleStateEnumValues() []AutonomousDataWarehouseBackupLifecycleStateEnum {
	values := make([]AutonomousDataWarehouseBackupLifecycleStateEnum, 0)
	for _, v := range mappingAutonomousDataWarehouseBackupLifecycleState {
		values = append(values, v)
	}
	return values
}
