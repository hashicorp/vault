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

// MaintenanceRunSummary Details of a Maintenance Run.
type MaintenanceRunSummary struct {

	// The OCID of the Maintenance Run.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The user-friendly name for the Maintenance Run.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The current state of the Maintenance Run.
	LifecycleState MaintenanceRunSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the Maintenance Run is scheduled for.
	TimeScheduled *common.SDKTime `mandatory:"true" json:"timeScheduled"`

	// The text describing this Maintenance Run.
	Description *string `mandatory:"false" json:"description"`

	// Additional information about the current lifecycleState.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// The date and time the Maintenance Run starts.
	TimeStarted *common.SDKTime `mandatory:"false" json:"timeStarted"`

	// The date and time the Maintenance Run was completed.
	TimeEnded *common.SDKTime `mandatory:"false" json:"timeEnded"`

	// The type of the target resource on which the Maintenance Run occurs.
	TargetResourceType MaintenanceRunSummaryTargetResourceTypeEnum `mandatory:"false" json:"targetResourceType,omitempty"`

	// The ID of the target resource on which the Maintenance Run occurs.
	TargetResourceId *string `mandatory:"false" json:"targetResourceId"`

	// Maintenance type.
	MaintenanceType MaintenanceRunSummaryMaintenanceTypeEnum `mandatory:"false" json:"maintenanceType,omitempty"`

	// Maintenance sub-type.
	MaintenanceSubtype MaintenanceRunSummaryMaintenanceSubtypeEnum `mandatory:"false" json:"maintenanceSubtype,omitempty"`
}

func (m MaintenanceRunSummary) String() string {
	return common.PointerString(m)
}

// MaintenanceRunSummaryLifecycleStateEnum Enum with underlying type: string
type MaintenanceRunSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for MaintenanceRunSummaryLifecycleStateEnum
const (
	MaintenanceRunSummaryLifecycleStateScheduled  MaintenanceRunSummaryLifecycleStateEnum = "SCHEDULED"
	MaintenanceRunSummaryLifecycleStateInProgress MaintenanceRunSummaryLifecycleStateEnum = "IN_PROGRESS"
	MaintenanceRunSummaryLifecycleStateSucceeded  MaintenanceRunSummaryLifecycleStateEnum = "SUCCEEDED"
	MaintenanceRunSummaryLifecycleStateSkipped    MaintenanceRunSummaryLifecycleStateEnum = "SKIPPED"
	MaintenanceRunSummaryLifecycleStateFailed     MaintenanceRunSummaryLifecycleStateEnum = "FAILED"
)

var mappingMaintenanceRunSummaryLifecycleState = map[string]MaintenanceRunSummaryLifecycleStateEnum{
	"SCHEDULED":   MaintenanceRunSummaryLifecycleStateScheduled,
	"IN_PROGRESS": MaintenanceRunSummaryLifecycleStateInProgress,
	"SUCCEEDED":   MaintenanceRunSummaryLifecycleStateSucceeded,
	"SKIPPED":     MaintenanceRunSummaryLifecycleStateSkipped,
	"FAILED":      MaintenanceRunSummaryLifecycleStateFailed,
}

// GetMaintenanceRunSummaryLifecycleStateEnumValues Enumerates the set of values for MaintenanceRunSummaryLifecycleStateEnum
func GetMaintenanceRunSummaryLifecycleStateEnumValues() []MaintenanceRunSummaryLifecycleStateEnum {
	values := make([]MaintenanceRunSummaryLifecycleStateEnum, 0)
	for _, v := range mappingMaintenanceRunSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}

// MaintenanceRunSummaryTargetResourceTypeEnum Enum with underlying type: string
type MaintenanceRunSummaryTargetResourceTypeEnum string

// Set of constants representing the allowable values for MaintenanceRunSummaryTargetResourceTypeEnum
const (
	MaintenanceRunSummaryTargetResourceTypeExadataInfrastructure MaintenanceRunSummaryTargetResourceTypeEnum = "AUTONOMOUS_EXADATA_INFRASTRUCTURE"
	MaintenanceRunSummaryTargetResourceTypeContainerDatabase     MaintenanceRunSummaryTargetResourceTypeEnum = "AUTONOMOUS_CONTAINER_DATABASE"
)

var mappingMaintenanceRunSummaryTargetResourceType = map[string]MaintenanceRunSummaryTargetResourceTypeEnum{
	"AUTONOMOUS_EXADATA_INFRASTRUCTURE": MaintenanceRunSummaryTargetResourceTypeExadataInfrastructure,
	"AUTONOMOUS_CONTAINER_DATABASE":     MaintenanceRunSummaryTargetResourceTypeContainerDatabase,
}

// GetMaintenanceRunSummaryTargetResourceTypeEnumValues Enumerates the set of values for MaintenanceRunSummaryTargetResourceTypeEnum
func GetMaintenanceRunSummaryTargetResourceTypeEnumValues() []MaintenanceRunSummaryTargetResourceTypeEnum {
	values := make([]MaintenanceRunSummaryTargetResourceTypeEnum, 0)
	for _, v := range mappingMaintenanceRunSummaryTargetResourceType {
		values = append(values, v)
	}
	return values
}

// MaintenanceRunSummaryMaintenanceTypeEnum Enum with underlying type: string
type MaintenanceRunSummaryMaintenanceTypeEnum string

// Set of constants representing the allowable values for MaintenanceRunSummaryMaintenanceTypeEnum
const (
	MaintenanceRunSummaryMaintenanceTypePlanned   MaintenanceRunSummaryMaintenanceTypeEnum = "PLANNED"
	MaintenanceRunSummaryMaintenanceTypeUnplanned MaintenanceRunSummaryMaintenanceTypeEnum = "UNPLANNED"
)

var mappingMaintenanceRunSummaryMaintenanceType = map[string]MaintenanceRunSummaryMaintenanceTypeEnum{
	"PLANNED":   MaintenanceRunSummaryMaintenanceTypePlanned,
	"UNPLANNED": MaintenanceRunSummaryMaintenanceTypeUnplanned,
}

// GetMaintenanceRunSummaryMaintenanceTypeEnumValues Enumerates the set of values for MaintenanceRunSummaryMaintenanceTypeEnum
func GetMaintenanceRunSummaryMaintenanceTypeEnumValues() []MaintenanceRunSummaryMaintenanceTypeEnum {
	values := make([]MaintenanceRunSummaryMaintenanceTypeEnum, 0)
	for _, v := range mappingMaintenanceRunSummaryMaintenanceType {
		values = append(values, v)
	}
	return values
}

// MaintenanceRunSummaryMaintenanceSubtypeEnum Enum with underlying type: string
type MaintenanceRunSummaryMaintenanceSubtypeEnum string

// Set of constants representing the allowable values for MaintenanceRunSummaryMaintenanceSubtypeEnum
const (
	MaintenanceRunSummaryMaintenanceSubtypeQuarterly MaintenanceRunSummaryMaintenanceSubtypeEnum = "QUARTERLY"
	MaintenanceRunSummaryMaintenanceSubtypeHardware  MaintenanceRunSummaryMaintenanceSubtypeEnum = "HARDWARE"
	MaintenanceRunSummaryMaintenanceSubtypeCritical  MaintenanceRunSummaryMaintenanceSubtypeEnum = "CRITICAL"
)

var mappingMaintenanceRunSummaryMaintenanceSubtype = map[string]MaintenanceRunSummaryMaintenanceSubtypeEnum{
	"QUARTERLY": MaintenanceRunSummaryMaintenanceSubtypeQuarterly,
	"HARDWARE":  MaintenanceRunSummaryMaintenanceSubtypeHardware,
	"CRITICAL":  MaintenanceRunSummaryMaintenanceSubtypeCritical,
}

// GetMaintenanceRunSummaryMaintenanceSubtypeEnumValues Enumerates the set of values for MaintenanceRunSummaryMaintenanceSubtypeEnum
func GetMaintenanceRunSummaryMaintenanceSubtypeEnumValues() []MaintenanceRunSummaryMaintenanceSubtypeEnum {
	values := make([]MaintenanceRunSummaryMaintenanceSubtypeEnum, 0)
	for _, v := range mappingMaintenanceRunSummaryMaintenanceSubtype {
		values = append(values, v)
	}
	return values
}
