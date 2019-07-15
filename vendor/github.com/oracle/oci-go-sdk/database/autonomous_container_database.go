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

// AutonomousContainerDatabase The representation of AutonomousContainerDatabase
type AutonomousContainerDatabase struct {

	// The OCID of the Autonomous Container Database.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The user-provided name for the Autonomous Container Database.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The service level agreement type of the container database. The default is STANDARD.
	ServiceLevelAgreementType AutonomousContainerDatabaseServiceLevelAgreementTypeEnum `mandatory:"true" json:"serviceLevelAgreementType"`

	// The OCID of the Autonomous Exadata Infrastructure.
	AutonomousExadataInfrastructureId *string `mandatory:"true" json:"autonomousExadataInfrastructureId"`

	// The current state of the Autonomous Container Database.
	LifecycleState AutonomousContainerDatabaseLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Database Patch model preference.
	PatchModel AutonomousContainerDatabasePatchModelEnum `mandatory:"true" json:"patchModel"`

	// Additional information about the current lifecycleState.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// The date and time the Autonomous was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the last maintenance run.
	LastMaintenanceRunId *string `mandatory:"false" json:"lastMaintenanceRunId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the next maintenance run.
	NextMaintenanceRunId *string `mandatory:"false" json:"nextMaintenanceRunId"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The availability domain of the Autonomous Container Database.
	AvailabilityDomain *string `mandatory:"false" json:"availabilityDomain"`

	BackupConfig *AutonomousContainerDatabaseBackupConfig `mandatory:"false" json:"backupConfig"`
}

func (m AutonomousContainerDatabase) String() string {
	return common.PointerString(m)
}

// AutonomousContainerDatabaseServiceLevelAgreementTypeEnum Enum with underlying type: string
type AutonomousContainerDatabaseServiceLevelAgreementTypeEnum string

// Set of constants representing the allowable values for AutonomousContainerDatabaseServiceLevelAgreementTypeEnum
const (
	AutonomousContainerDatabaseServiceLevelAgreementTypeStandard        AutonomousContainerDatabaseServiceLevelAgreementTypeEnum = "STANDARD"
	AutonomousContainerDatabaseServiceLevelAgreementTypeMissionCritical AutonomousContainerDatabaseServiceLevelAgreementTypeEnum = "MISSION_CRITICAL"
)

var mappingAutonomousContainerDatabaseServiceLevelAgreementType = map[string]AutonomousContainerDatabaseServiceLevelAgreementTypeEnum{
	"STANDARD":         AutonomousContainerDatabaseServiceLevelAgreementTypeStandard,
	"MISSION_CRITICAL": AutonomousContainerDatabaseServiceLevelAgreementTypeMissionCritical,
}

// GetAutonomousContainerDatabaseServiceLevelAgreementTypeEnumValues Enumerates the set of values for AutonomousContainerDatabaseServiceLevelAgreementTypeEnum
func GetAutonomousContainerDatabaseServiceLevelAgreementTypeEnumValues() []AutonomousContainerDatabaseServiceLevelAgreementTypeEnum {
	values := make([]AutonomousContainerDatabaseServiceLevelAgreementTypeEnum, 0)
	for _, v := range mappingAutonomousContainerDatabaseServiceLevelAgreementType {
		values = append(values, v)
	}
	return values
}

// AutonomousContainerDatabaseLifecycleStateEnum Enum with underlying type: string
type AutonomousContainerDatabaseLifecycleStateEnum string

// Set of constants representing the allowable values for AutonomousContainerDatabaseLifecycleStateEnum
const (
	AutonomousContainerDatabaseLifecycleStateProvisioning          AutonomousContainerDatabaseLifecycleStateEnum = "PROVISIONING"
	AutonomousContainerDatabaseLifecycleStateAvailable             AutonomousContainerDatabaseLifecycleStateEnum = "AVAILABLE"
	AutonomousContainerDatabaseLifecycleStateUpdating              AutonomousContainerDatabaseLifecycleStateEnum = "UPDATING"
	AutonomousContainerDatabaseLifecycleStateTerminating           AutonomousContainerDatabaseLifecycleStateEnum = "TERMINATING"
	AutonomousContainerDatabaseLifecycleStateTerminated            AutonomousContainerDatabaseLifecycleStateEnum = "TERMINATED"
	AutonomousContainerDatabaseLifecycleStateFailed                AutonomousContainerDatabaseLifecycleStateEnum = "FAILED"
	AutonomousContainerDatabaseLifecycleStateBackupInProgress      AutonomousContainerDatabaseLifecycleStateEnum = "BACKUP_IN_PROGRESS"
	AutonomousContainerDatabaseLifecycleStateRestoring             AutonomousContainerDatabaseLifecycleStateEnum = "RESTORING"
	AutonomousContainerDatabaseLifecycleStateRestoreFailed         AutonomousContainerDatabaseLifecycleStateEnum = "RESTORE_FAILED"
	AutonomousContainerDatabaseLifecycleStateRestarting            AutonomousContainerDatabaseLifecycleStateEnum = "RESTARTING"
	AutonomousContainerDatabaseLifecycleStateMaintenanceInProgress AutonomousContainerDatabaseLifecycleStateEnum = "MAINTENANCE_IN_PROGRESS"
)

var mappingAutonomousContainerDatabaseLifecycleState = map[string]AutonomousContainerDatabaseLifecycleStateEnum{
	"PROVISIONING":            AutonomousContainerDatabaseLifecycleStateProvisioning,
	"AVAILABLE":               AutonomousContainerDatabaseLifecycleStateAvailable,
	"UPDATING":                AutonomousContainerDatabaseLifecycleStateUpdating,
	"TERMINATING":             AutonomousContainerDatabaseLifecycleStateTerminating,
	"TERMINATED":              AutonomousContainerDatabaseLifecycleStateTerminated,
	"FAILED":                  AutonomousContainerDatabaseLifecycleStateFailed,
	"BACKUP_IN_PROGRESS":      AutonomousContainerDatabaseLifecycleStateBackupInProgress,
	"RESTORING":               AutonomousContainerDatabaseLifecycleStateRestoring,
	"RESTORE_FAILED":          AutonomousContainerDatabaseLifecycleStateRestoreFailed,
	"RESTARTING":              AutonomousContainerDatabaseLifecycleStateRestarting,
	"MAINTENANCE_IN_PROGRESS": AutonomousContainerDatabaseLifecycleStateMaintenanceInProgress,
}

// GetAutonomousContainerDatabaseLifecycleStateEnumValues Enumerates the set of values for AutonomousContainerDatabaseLifecycleStateEnum
func GetAutonomousContainerDatabaseLifecycleStateEnumValues() []AutonomousContainerDatabaseLifecycleStateEnum {
	values := make([]AutonomousContainerDatabaseLifecycleStateEnum, 0)
	for _, v := range mappingAutonomousContainerDatabaseLifecycleState {
		values = append(values, v)
	}
	return values
}

// AutonomousContainerDatabasePatchModelEnum Enum with underlying type: string
type AutonomousContainerDatabasePatchModelEnum string

// Set of constants representing the allowable values for AutonomousContainerDatabasePatchModelEnum
const (
	AutonomousContainerDatabasePatchModelUpdates         AutonomousContainerDatabasePatchModelEnum = "RELEASE_UPDATES"
	AutonomousContainerDatabasePatchModelUpdateRevisions AutonomousContainerDatabasePatchModelEnum = "RELEASE_UPDATE_REVISIONS"
)

var mappingAutonomousContainerDatabasePatchModel = map[string]AutonomousContainerDatabasePatchModelEnum{
	"RELEASE_UPDATES":          AutonomousContainerDatabasePatchModelUpdates,
	"RELEASE_UPDATE_REVISIONS": AutonomousContainerDatabasePatchModelUpdateRevisions,
}

// GetAutonomousContainerDatabasePatchModelEnumValues Enumerates the set of values for AutonomousContainerDatabasePatchModelEnum
func GetAutonomousContainerDatabasePatchModelEnumValues() []AutonomousContainerDatabasePatchModelEnum {
	values := make([]AutonomousContainerDatabasePatchModelEnum, 0)
	for _, v := range mappingAutonomousContainerDatabasePatchModel {
		values = append(values, v)
	}
	return values
}
