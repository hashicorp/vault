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

// AutonomousExadataInfrastructure The representation of AutonomousExadataInfrastructure
type AutonomousExadataInfrastructure struct {

	// The OCID of the Autonomous Exadata Infrastructure.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The user-friendly name for the Autonomous Exadata Infrastructure.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The name of the availability domain that the Autonomous Exadata Infrastructure is located in.
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID of the subnet the Autonomous Exadata Infrastructure is associated with.
	// **Subnet Restrictions:**
	// - For Autonomous Databases with Autonomous Exadata Infrastructure, do not use a subnet that overlaps with 192.168.128.0/20
	// These subnets are used by the Oracle Clusterware private interconnect on the database instance.
	// Specifying an overlapping subnet will cause the private interconnect to malfunction.
	// This restriction applies to both the client subnet and backup subnet.
	SubnetId *string `mandatory:"true" json:"subnetId"`

	// The shape of the Autonomous Exadata Infrastructure. The shape determines resources to allocate to the Autonomous Exadata Infrastructure (CPU cores, memory and storage).
	Shape *string `mandatory:"true" json:"shape"`

	// The host name for the Autonomous Exadata Infrastructure node.
	Hostname *string `mandatory:"true" json:"hostname"`

	// The domain name for the Autonomous Exadata Infrastructure.
	Domain *string `mandatory:"true" json:"domain"`

	// The current lifecycle state of the Autonomous Exadata Infrastructure.
	LifecycleState AutonomousExadataInfrastructureLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	MaintenanceWindow *MaintenanceWindow `mandatory:"true" json:"maintenanceWindow"`

	// Additional information about the current lifecycle state of the Autonomous Exadata Infrastructure.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// The Oracle license model that applies to all databases in the Autonomous Exadata Infrastructure. The default is BRING_YOUR_OWN_LICENSE.
	LicenseModel AutonomousExadataInfrastructureLicenseModelEnum `mandatory:"false" json:"licenseModel,omitempty"`

	// The date and time the Autonomous Exadata Infrastructure was created.
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
}

func (m AutonomousExadataInfrastructure) String() string {
	return common.PointerString(m)
}

// AutonomousExadataInfrastructureLifecycleStateEnum Enum with underlying type: string
type AutonomousExadataInfrastructureLifecycleStateEnum string

// Set of constants representing the allowable values for AutonomousExadataInfrastructureLifecycleStateEnum
const (
	AutonomousExadataInfrastructureLifecycleStateProvisioning          AutonomousExadataInfrastructureLifecycleStateEnum = "PROVISIONING"
	AutonomousExadataInfrastructureLifecycleStateAvailable             AutonomousExadataInfrastructureLifecycleStateEnum = "AVAILABLE"
	AutonomousExadataInfrastructureLifecycleStateUpdating              AutonomousExadataInfrastructureLifecycleStateEnum = "UPDATING"
	AutonomousExadataInfrastructureLifecycleStateTerminating           AutonomousExadataInfrastructureLifecycleStateEnum = "TERMINATING"
	AutonomousExadataInfrastructureLifecycleStateTerminated            AutonomousExadataInfrastructureLifecycleStateEnum = "TERMINATED"
	AutonomousExadataInfrastructureLifecycleStateFailed                AutonomousExadataInfrastructureLifecycleStateEnum = "FAILED"
	AutonomousExadataInfrastructureLifecycleStateMaintenanceInProgress AutonomousExadataInfrastructureLifecycleStateEnum = "MAINTENANCE_IN_PROGRESS"
)

var mappingAutonomousExadataInfrastructureLifecycleState = map[string]AutonomousExadataInfrastructureLifecycleStateEnum{
	"PROVISIONING":            AutonomousExadataInfrastructureLifecycleStateProvisioning,
	"AVAILABLE":               AutonomousExadataInfrastructureLifecycleStateAvailable,
	"UPDATING":                AutonomousExadataInfrastructureLifecycleStateUpdating,
	"TERMINATING":             AutonomousExadataInfrastructureLifecycleStateTerminating,
	"TERMINATED":              AutonomousExadataInfrastructureLifecycleStateTerminated,
	"FAILED":                  AutonomousExadataInfrastructureLifecycleStateFailed,
	"MAINTENANCE_IN_PROGRESS": AutonomousExadataInfrastructureLifecycleStateMaintenanceInProgress,
}

// GetAutonomousExadataInfrastructureLifecycleStateEnumValues Enumerates the set of values for AutonomousExadataInfrastructureLifecycleStateEnum
func GetAutonomousExadataInfrastructureLifecycleStateEnumValues() []AutonomousExadataInfrastructureLifecycleStateEnum {
	values := make([]AutonomousExadataInfrastructureLifecycleStateEnum, 0)
	for _, v := range mappingAutonomousExadataInfrastructureLifecycleState {
		values = append(values, v)
	}
	return values
}

// AutonomousExadataInfrastructureLicenseModelEnum Enum with underlying type: string
type AutonomousExadataInfrastructureLicenseModelEnum string

// Set of constants representing the allowable values for AutonomousExadataInfrastructureLicenseModelEnum
const (
	AutonomousExadataInfrastructureLicenseModelLicenseIncluded     AutonomousExadataInfrastructureLicenseModelEnum = "LICENSE_INCLUDED"
	AutonomousExadataInfrastructureLicenseModelBringYourOwnLicense AutonomousExadataInfrastructureLicenseModelEnum = "BRING_YOUR_OWN_LICENSE"
)

var mappingAutonomousExadataInfrastructureLicenseModel = map[string]AutonomousExadataInfrastructureLicenseModelEnum{
	"LICENSE_INCLUDED":       AutonomousExadataInfrastructureLicenseModelLicenseIncluded,
	"BRING_YOUR_OWN_LICENSE": AutonomousExadataInfrastructureLicenseModelBringYourOwnLicense,
}

// GetAutonomousExadataInfrastructureLicenseModelEnumValues Enumerates the set of values for AutonomousExadataInfrastructureLicenseModelEnum
func GetAutonomousExadataInfrastructureLicenseModelEnumValues() []AutonomousExadataInfrastructureLicenseModelEnum {
	values := make([]AutonomousExadataInfrastructureLicenseModelEnum, 0)
	for _, v := range mappingAutonomousExadataInfrastructureLicenseModel {
		values = append(values, v)
	}
	return values
}
