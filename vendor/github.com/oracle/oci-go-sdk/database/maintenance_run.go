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

// MaintenanceRun Details of a Maintenance Run.
type MaintenanceRun struct {

	// The OCID of the Maintenance Run.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The user-friendly name for the Maintenance Run.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The current state of the Maintenance Run.
	LifecycleState MaintenanceRunLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

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
	TargetResourceType MaintenanceRunTargetResourceTypeEnum `mandatory:"false" json:"targetResourceType,omitempty"`

	// The ID of the target resource on which the Maintenance Run occurs.
	TargetResourceId *string `mandatory:"false" json:"targetResourceId"`

	// Maintenance type.
	MaintenanceType MaintenanceRunMaintenanceTypeEnum `mandatory:"false" json:"maintenanceType,omitempty"`

	// Maintenance sub-type.
	MaintenanceSubtype MaintenanceRunMaintenanceSubtypeEnum `mandatory:"false" json:"maintenanceSubtype,omitempty"`
}

func (m MaintenanceRun) String() string {
	return common.PointerString(m)
}

// MaintenanceRunLifecycleStateEnum Enum with underlying type: string
type MaintenanceRunLifecycleStateEnum string

// Set of constants representing the allowable values for MaintenanceRunLifecycleStateEnum
const (
	MaintenanceRunLifecycleStateScheduled  MaintenanceRunLifecycleStateEnum = "SCHEDULED"
	MaintenanceRunLifecycleStateInProgress MaintenanceRunLifecycleStateEnum = "IN_PROGRESS"
	MaintenanceRunLifecycleStateSucceeded  MaintenanceRunLifecycleStateEnum = "SUCCEEDED"
	MaintenanceRunLifecycleStateSkipped    MaintenanceRunLifecycleStateEnum = "SKIPPED"
	MaintenanceRunLifecycleStateFailed     MaintenanceRunLifecycleStateEnum = "FAILED"
)

var mappingMaintenanceRunLifecycleState = map[string]MaintenanceRunLifecycleStateEnum{
	"SCHEDULED":   MaintenanceRunLifecycleStateScheduled,
	"IN_PROGRESS": MaintenanceRunLifecycleStateInProgress,
	"SUCCEEDED":   MaintenanceRunLifecycleStateSucceeded,
	"SKIPPED":     MaintenanceRunLifecycleStateSkipped,
	"FAILED":      MaintenanceRunLifecycleStateFailed,
}

// GetMaintenanceRunLifecycleStateEnumValues Enumerates the set of values for MaintenanceRunLifecycleStateEnum
func GetMaintenanceRunLifecycleStateEnumValues() []MaintenanceRunLifecycleStateEnum {
	values := make([]MaintenanceRunLifecycleStateEnum, 0)
	for _, v := range mappingMaintenanceRunLifecycleState {
		values = append(values, v)
	}
	return values
}

// MaintenanceRunTargetResourceTypeEnum Enum with underlying type: string
type MaintenanceRunTargetResourceTypeEnum string

// Set of constants representing the allowable values for MaintenanceRunTargetResourceTypeEnum
const (
	MaintenanceRunTargetResourceTypeExadataInfrastructure MaintenanceRunTargetResourceTypeEnum = "AUTONOMOUS_EXADATA_INFRASTRUCTURE"
	MaintenanceRunTargetResourceTypeContainerDatabase     MaintenanceRunTargetResourceTypeEnum = "AUTONOMOUS_CONTAINER_DATABASE"
)

var mappingMaintenanceRunTargetResourceType = map[string]MaintenanceRunTargetResourceTypeEnum{
	"AUTONOMOUS_EXADATA_INFRASTRUCTURE": MaintenanceRunTargetResourceTypeExadataInfrastructure,
	"AUTONOMOUS_CONTAINER_DATABASE":     MaintenanceRunTargetResourceTypeContainerDatabase,
}

// GetMaintenanceRunTargetResourceTypeEnumValues Enumerates the set of values for MaintenanceRunTargetResourceTypeEnum
func GetMaintenanceRunTargetResourceTypeEnumValues() []MaintenanceRunTargetResourceTypeEnum {
	values := make([]MaintenanceRunTargetResourceTypeEnum, 0)
	for _, v := range mappingMaintenanceRunTargetResourceType {
		values = append(values, v)
	}
	return values
}

// MaintenanceRunMaintenanceTypeEnum Enum with underlying type: string
type MaintenanceRunMaintenanceTypeEnum string

// Set of constants representing the allowable values for MaintenanceRunMaintenanceTypeEnum
const (
	MaintenanceRunMaintenanceTypePlanned   MaintenanceRunMaintenanceTypeEnum = "PLANNED"
	MaintenanceRunMaintenanceTypeUnplanned MaintenanceRunMaintenanceTypeEnum = "UNPLANNED"
)

var mappingMaintenanceRunMaintenanceType = map[string]MaintenanceRunMaintenanceTypeEnum{
	"PLANNED":   MaintenanceRunMaintenanceTypePlanned,
	"UNPLANNED": MaintenanceRunMaintenanceTypeUnplanned,
}

// GetMaintenanceRunMaintenanceTypeEnumValues Enumerates the set of values for MaintenanceRunMaintenanceTypeEnum
func GetMaintenanceRunMaintenanceTypeEnumValues() []MaintenanceRunMaintenanceTypeEnum {
	values := make([]MaintenanceRunMaintenanceTypeEnum, 0)
	for _, v := range mappingMaintenanceRunMaintenanceType {
		values = append(values, v)
	}
	return values
}

// MaintenanceRunMaintenanceSubtypeEnum Enum with underlying type: string
type MaintenanceRunMaintenanceSubtypeEnum string

// Set of constants representing the allowable values for MaintenanceRunMaintenanceSubtypeEnum
const (
	MaintenanceRunMaintenanceSubtypeQuarterly MaintenanceRunMaintenanceSubtypeEnum = "QUARTERLY"
	MaintenanceRunMaintenanceSubtypeHardware  MaintenanceRunMaintenanceSubtypeEnum = "HARDWARE"
	MaintenanceRunMaintenanceSubtypeCritical  MaintenanceRunMaintenanceSubtypeEnum = "CRITICAL"
)

var mappingMaintenanceRunMaintenanceSubtype = map[string]MaintenanceRunMaintenanceSubtypeEnum{
	"QUARTERLY": MaintenanceRunMaintenanceSubtypeQuarterly,
	"HARDWARE":  MaintenanceRunMaintenanceSubtypeHardware,
	"CRITICAL":  MaintenanceRunMaintenanceSubtypeCritical,
}

// GetMaintenanceRunMaintenanceSubtypeEnumValues Enumerates the set of values for MaintenanceRunMaintenanceSubtypeEnum
func GetMaintenanceRunMaintenanceSubtypeEnumValues() []MaintenanceRunMaintenanceSubtypeEnum {
	values := make([]MaintenanceRunMaintenanceSubtypeEnum, 0)
	for _, v := range mappingMaintenanceRunMaintenanceSubtype {
		values = append(values, v)
	}
	return values
}
