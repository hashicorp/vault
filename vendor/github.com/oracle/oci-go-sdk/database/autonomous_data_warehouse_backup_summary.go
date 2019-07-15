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

// AutonomousDataWarehouseBackupSummary **Deprecated.** See AutonomousDataWarehouseBackupSummary for reference information about Autonomous Data Warehouse backups.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type AutonomousDataWarehouseBackupSummary struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Autonomous Data Warehouse backup.
	Id *string `mandatory:"true" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Autonomous Data Warehouse.
	AutonomousDataWarehouseId *string `mandatory:"true" json:"autonomousDataWarehouseId"`

	// The user-friendly name for the backup. The name does not have to be unique.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The type of backup.
	Type AutonomousDataWarehouseBackupSummaryTypeEnum `mandatory:"true" json:"type"`

	// Indicates whether the backup is user-initiated or automatic.
	IsAutomatic *bool `mandatory:"true" json:"isAutomatic"`

	// The current state of the backup.
	LifecycleState AutonomousDataWarehouseBackupSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the backup started.
	TimeStarted *common.SDKTime `mandatory:"false" json:"timeStarted"`

	// The date and time the backup completed.
	TimeEnded *common.SDKTime `mandatory:"false" json:"timeEnded"`

	// Additional information about the current lifecycle state.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`
}

func (m AutonomousDataWarehouseBackupSummary) String() string {
	return common.PointerString(m)
}

// AutonomousDataWarehouseBackupSummaryTypeEnum Enum with underlying type: string
type AutonomousDataWarehouseBackupSummaryTypeEnum string

// Set of constants representing the allowable values for AutonomousDataWarehouseBackupSummaryTypeEnum
const (
	AutonomousDataWarehouseBackupSummaryTypeIncremental AutonomousDataWarehouseBackupSummaryTypeEnum = "INCREMENTAL"
	AutonomousDataWarehouseBackupSummaryTypeFull        AutonomousDataWarehouseBackupSummaryTypeEnum = "FULL"
)

var mappingAutonomousDataWarehouseBackupSummaryType = map[string]AutonomousDataWarehouseBackupSummaryTypeEnum{
	"INCREMENTAL": AutonomousDataWarehouseBackupSummaryTypeIncremental,
	"FULL":        AutonomousDataWarehouseBackupSummaryTypeFull,
}

// GetAutonomousDataWarehouseBackupSummaryTypeEnumValues Enumerates the set of values for AutonomousDataWarehouseBackupSummaryTypeEnum
func GetAutonomousDataWarehouseBackupSummaryTypeEnumValues() []AutonomousDataWarehouseBackupSummaryTypeEnum {
	values := make([]AutonomousDataWarehouseBackupSummaryTypeEnum, 0)
	for _, v := range mappingAutonomousDataWarehouseBackupSummaryType {
		values = append(values, v)
	}
	return values
}

// AutonomousDataWarehouseBackupSummaryLifecycleStateEnum Enum with underlying type: string
type AutonomousDataWarehouseBackupSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for AutonomousDataWarehouseBackupSummaryLifecycleStateEnum
const (
	AutonomousDataWarehouseBackupSummaryLifecycleStateCreating AutonomousDataWarehouseBackupSummaryLifecycleStateEnum = "CREATING"
	AutonomousDataWarehouseBackupSummaryLifecycleStateActive   AutonomousDataWarehouseBackupSummaryLifecycleStateEnum = "ACTIVE"
	AutonomousDataWarehouseBackupSummaryLifecycleStateDeleting AutonomousDataWarehouseBackupSummaryLifecycleStateEnum = "DELETING"
	AutonomousDataWarehouseBackupSummaryLifecycleStateDeleted  AutonomousDataWarehouseBackupSummaryLifecycleStateEnum = "DELETED"
	AutonomousDataWarehouseBackupSummaryLifecycleStateFailed   AutonomousDataWarehouseBackupSummaryLifecycleStateEnum = "FAILED"
)

var mappingAutonomousDataWarehouseBackupSummaryLifecycleState = map[string]AutonomousDataWarehouseBackupSummaryLifecycleStateEnum{
	"CREATING": AutonomousDataWarehouseBackupSummaryLifecycleStateCreating,
	"ACTIVE":   AutonomousDataWarehouseBackupSummaryLifecycleStateActive,
	"DELETING": AutonomousDataWarehouseBackupSummaryLifecycleStateDeleting,
	"DELETED":  AutonomousDataWarehouseBackupSummaryLifecycleStateDeleted,
	"FAILED":   AutonomousDataWarehouseBackupSummaryLifecycleStateFailed,
}

// GetAutonomousDataWarehouseBackupSummaryLifecycleStateEnumValues Enumerates the set of values for AutonomousDataWarehouseBackupSummaryLifecycleStateEnum
func GetAutonomousDataWarehouseBackupSummaryLifecycleStateEnumValues() []AutonomousDataWarehouseBackupSummaryLifecycleStateEnum {
	values := make([]AutonomousDataWarehouseBackupSummaryLifecycleStateEnum, 0)
	for _, v := range mappingAutonomousDataWarehouseBackupSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
