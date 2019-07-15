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

// AutonomousContainerDatabaseSummary An Autonomous Container Database is a container database service that enables the customer to host one or more databases within the container database. A basic container database runs on a single Autonomous Exadata Infrastructure from an availability domain without the Extreme Availability features enabled.
type AutonomousContainerDatabaseSummary struct {

	// The OCID of the Autonomous Container Database.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The user-provided name for the Autonomous Container Database.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The service level agreement type of the container database. The default is STANDARD.
	ServiceLevelAgreementType AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnum `mandatory:"true" json:"serviceLevelAgreementType"`

	// The OCID of the Autonomous Exadata Infrastructure.
	AutonomousExadataInfrastructureId *string `mandatory:"true" json:"autonomousExadataInfrastructureId"`

	// The current state of the Autonomous Container Database.
	LifecycleState AutonomousContainerDatabaseSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Database Patch model preference.
	PatchModel AutonomousContainerDatabaseSummaryPatchModelEnum `mandatory:"true" json:"patchModel"`

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

func (m AutonomousContainerDatabaseSummary) String() string {
	return common.PointerString(m)
}

// AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnum Enum with underlying type: string
type AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnum string

// Set of constants representing the allowable values for AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnum
const (
	AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeStandard        AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnum = "STANDARD"
	AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeMissionCritical AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnum = "MISSION_CRITICAL"
)

var mappingAutonomousContainerDatabaseSummaryServiceLevelAgreementType = map[string]AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnum{
	"STANDARD":         AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeStandard,
	"MISSION_CRITICAL": AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeMissionCritical,
}

// GetAutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnumValues Enumerates the set of values for AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnum
func GetAutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnumValues() []AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnum {
	values := make([]AutonomousContainerDatabaseSummaryServiceLevelAgreementTypeEnum, 0)
	for _, v := range mappingAutonomousContainerDatabaseSummaryServiceLevelAgreementType {
		values = append(values, v)
	}
	return values
}

// AutonomousContainerDatabaseSummaryLifecycleStateEnum Enum with underlying type: string
type AutonomousContainerDatabaseSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for AutonomousContainerDatabaseSummaryLifecycleStateEnum
const (
	AutonomousContainerDatabaseSummaryLifecycleStateProvisioning          AutonomousContainerDatabaseSummaryLifecycleStateEnum = "PROVISIONING"
	AutonomousContainerDatabaseSummaryLifecycleStateAvailable             AutonomousContainerDatabaseSummaryLifecycleStateEnum = "AVAILABLE"
	AutonomousContainerDatabaseSummaryLifecycleStateUpdating              AutonomousContainerDatabaseSummaryLifecycleStateEnum = "UPDATING"
	AutonomousContainerDatabaseSummaryLifecycleStateTerminating           AutonomousContainerDatabaseSummaryLifecycleStateEnum = "TERMINATING"
	AutonomousContainerDatabaseSummaryLifecycleStateTerminated            AutonomousContainerDatabaseSummaryLifecycleStateEnum = "TERMINATED"
	AutonomousContainerDatabaseSummaryLifecycleStateFailed                AutonomousContainerDatabaseSummaryLifecycleStateEnum = "FAILED"
	AutonomousContainerDatabaseSummaryLifecycleStateBackupInProgress      AutonomousContainerDatabaseSummaryLifecycleStateEnum = "BACKUP_IN_PROGRESS"
	AutonomousContainerDatabaseSummaryLifecycleStateRestoring             AutonomousContainerDatabaseSummaryLifecycleStateEnum = "RESTORING"
	AutonomousContainerDatabaseSummaryLifecycleStateRestoreFailed         AutonomousContainerDatabaseSummaryLifecycleStateEnum = "RESTORE_FAILED"
	AutonomousContainerDatabaseSummaryLifecycleStateRestarting            AutonomousContainerDatabaseSummaryLifecycleStateEnum = "RESTARTING"
	AutonomousContainerDatabaseSummaryLifecycleStateMaintenanceInProgress AutonomousContainerDatabaseSummaryLifecycleStateEnum = "MAINTENANCE_IN_PROGRESS"
)

var mappingAutonomousContainerDatabaseSummaryLifecycleState = map[string]AutonomousContainerDatabaseSummaryLifecycleStateEnum{
	"PROVISIONING":            AutonomousContainerDatabaseSummaryLifecycleStateProvisioning,
	"AVAILABLE":               AutonomousContainerDatabaseSummaryLifecycleStateAvailable,
	"UPDATING":                AutonomousContainerDatabaseSummaryLifecycleStateUpdating,
	"TERMINATING":             AutonomousContainerDatabaseSummaryLifecycleStateTerminating,
	"TERMINATED":              AutonomousContainerDatabaseSummaryLifecycleStateTerminated,
	"FAILED":                  AutonomousContainerDatabaseSummaryLifecycleStateFailed,
	"BACKUP_IN_PROGRESS":      AutonomousContainerDatabaseSummaryLifecycleStateBackupInProgress,
	"RESTORING":               AutonomousContainerDatabaseSummaryLifecycleStateRestoring,
	"RESTORE_FAILED":          AutonomousContainerDatabaseSummaryLifecycleStateRestoreFailed,
	"RESTARTING":              AutonomousContainerDatabaseSummaryLifecycleStateRestarting,
	"MAINTENANCE_IN_PROGRESS": AutonomousContainerDatabaseSummaryLifecycleStateMaintenanceInProgress,
}

// GetAutonomousContainerDatabaseSummaryLifecycleStateEnumValues Enumerates the set of values for AutonomousContainerDatabaseSummaryLifecycleStateEnum
func GetAutonomousContainerDatabaseSummaryLifecycleStateEnumValues() []AutonomousContainerDatabaseSummaryLifecycleStateEnum {
	values := make([]AutonomousContainerDatabaseSummaryLifecycleStateEnum, 0)
	for _, v := range mappingAutonomousContainerDatabaseSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}

// AutonomousContainerDatabaseSummaryPatchModelEnum Enum with underlying type: string
type AutonomousContainerDatabaseSummaryPatchModelEnum string

// Set of constants representing the allowable values for AutonomousContainerDatabaseSummaryPatchModelEnum
const (
	AutonomousContainerDatabaseSummaryPatchModelUpdates         AutonomousContainerDatabaseSummaryPatchModelEnum = "RELEASE_UPDATES"
	AutonomousContainerDatabaseSummaryPatchModelUpdateRevisions AutonomousContainerDatabaseSummaryPatchModelEnum = "RELEASE_UPDATE_REVISIONS"
)

var mappingAutonomousContainerDatabaseSummaryPatchModel = map[string]AutonomousContainerDatabaseSummaryPatchModelEnum{
	"RELEASE_UPDATES":          AutonomousContainerDatabaseSummaryPatchModelUpdates,
	"RELEASE_UPDATE_REVISIONS": AutonomousContainerDatabaseSummaryPatchModelUpdateRevisions,
}

// GetAutonomousContainerDatabaseSummaryPatchModelEnumValues Enumerates the set of values for AutonomousContainerDatabaseSummaryPatchModelEnum
func GetAutonomousContainerDatabaseSummaryPatchModelEnumValues() []AutonomousContainerDatabaseSummaryPatchModelEnum {
	values := make([]AutonomousContainerDatabaseSummaryPatchModelEnum, 0)
	for _, v := range mappingAutonomousContainerDatabaseSummaryPatchModel {
		values = append(values, v)
	}
	return values
}
