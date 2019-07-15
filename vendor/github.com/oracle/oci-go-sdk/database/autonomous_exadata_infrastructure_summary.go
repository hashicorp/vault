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

// AutonomousExadataInfrastructureSummary Infrastructure that enables the running of multiple Autonomous Databases within a dedicated DB system.
// For more information about Autonomous Exadata Infrastructure, see
// Overview of Autonomous Database (https://docs.cloud.oracle.com/iaas/Content/Database/Concepts/adboverview.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized, talk to an administrator. If you're an administrator who needs to write policies to give users access, see Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// For information about access control and compartments, see
// Overview of the Identity Service (https://docs.cloud.oracle.com/Content/Identity/Concepts/overview.htm).
// For information about availability domains, see
// Regions and Availability Domains (https://docs.cloud.oracle.com/Content/General/Concepts/regions.htm).
// To get a list of availability domains, use the ListAvailabilityDomains operation
// in the Identity service API.
type AutonomousExadataInfrastructureSummary struct {

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
	LifecycleState AutonomousExadataInfrastructureSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	MaintenanceWindow *MaintenanceWindow `mandatory:"true" json:"maintenanceWindow"`

	// Additional information about the current lifecycle state of the Autonomous Exadata Infrastructure.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// The Oracle license model that applies to all databases in the Autonomous Exadata Infrastructure. The default is BRING_YOUR_OWN_LICENSE.
	LicenseModel AutonomousExadataInfrastructureSummaryLicenseModelEnum `mandatory:"false" json:"licenseModel,omitempty"`

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

func (m AutonomousExadataInfrastructureSummary) String() string {
	return common.PointerString(m)
}

// AutonomousExadataInfrastructureSummaryLifecycleStateEnum Enum with underlying type: string
type AutonomousExadataInfrastructureSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for AutonomousExadataInfrastructureSummaryLifecycleStateEnum
const (
	AutonomousExadataInfrastructureSummaryLifecycleStateProvisioning          AutonomousExadataInfrastructureSummaryLifecycleStateEnum = "PROVISIONING"
	AutonomousExadataInfrastructureSummaryLifecycleStateAvailable             AutonomousExadataInfrastructureSummaryLifecycleStateEnum = "AVAILABLE"
	AutonomousExadataInfrastructureSummaryLifecycleStateUpdating              AutonomousExadataInfrastructureSummaryLifecycleStateEnum = "UPDATING"
	AutonomousExadataInfrastructureSummaryLifecycleStateTerminating           AutonomousExadataInfrastructureSummaryLifecycleStateEnum = "TERMINATING"
	AutonomousExadataInfrastructureSummaryLifecycleStateTerminated            AutonomousExadataInfrastructureSummaryLifecycleStateEnum = "TERMINATED"
	AutonomousExadataInfrastructureSummaryLifecycleStateFailed                AutonomousExadataInfrastructureSummaryLifecycleStateEnum = "FAILED"
	AutonomousExadataInfrastructureSummaryLifecycleStateMaintenanceInProgress AutonomousExadataInfrastructureSummaryLifecycleStateEnum = "MAINTENANCE_IN_PROGRESS"
)

var mappingAutonomousExadataInfrastructureSummaryLifecycleState = map[string]AutonomousExadataInfrastructureSummaryLifecycleStateEnum{
	"PROVISIONING":            AutonomousExadataInfrastructureSummaryLifecycleStateProvisioning,
	"AVAILABLE":               AutonomousExadataInfrastructureSummaryLifecycleStateAvailable,
	"UPDATING":                AutonomousExadataInfrastructureSummaryLifecycleStateUpdating,
	"TERMINATING":             AutonomousExadataInfrastructureSummaryLifecycleStateTerminating,
	"TERMINATED":              AutonomousExadataInfrastructureSummaryLifecycleStateTerminated,
	"FAILED":                  AutonomousExadataInfrastructureSummaryLifecycleStateFailed,
	"MAINTENANCE_IN_PROGRESS": AutonomousExadataInfrastructureSummaryLifecycleStateMaintenanceInProgress,
}

// GetAutonomousExadataInfrastructureSummaryLifecycleStateEnumValues Enumerates the set of values for AutonomousExadataInfrastructureSummaryLifecycleStateEnum
func GetAutonomousExadataInfrastructureSummaryLifecycleStateEnumValues() []AutonomousExadataInfrastructureSummaryLifecycleStateEnum {
	values := make([]AutonomousExadataInfrastructureSummaryLifecycleStateEnum, 0)
	for _, v := range mappingAutonomousExadataInfrastructureSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}

// AutonomousExadataInfrastructureSummaryLicenseModelEnum Enum with underlying type: string
type AutonomousExadataInfrastructureSummaryLicenseModelEnum string

// Set of constants representing the allowable values for AutonomousExadataInfrastructureSummaryLicenseModelEnum
const (
	AutonomousExadataInfrastructureSummaryLicenseModelLicenseIncluded     AutonomousExadataInfrastructureSummaryLicenseModelEnum = "LICENSE_INCLUDED"
	AutonomousExadataInfrastructureSummaryLicenseModelBringYourOwnLicense AutonomousExadataInfrastructureSummaryLicenseModelEnum = "BRING_YOUR_OWN_LICENSE"
)

var mappingAutonomousExadataInfrastructureSummaryLicenseModel = map[string]AutonomousExadataInfrastructureSummaryLicenseModelEnum{
	"LICENSE_INCLUDED":       AutonomousExadataInfrastructureSummaryLicenseModelLicenseIncluded,
	"BRING_YOUR_OWN_LICENSE": AutonomousExadataInfrastructureSummaryLicenseModelBringYourOwnLicense,
}

// GetAutonomousExadataInfrastructureSummaryLicenseModelEnumValues Enumerates the set of values for AutonomousExadataInfrastructureSummaryLicenseModelEnum
func GetAutonomousExadataInfrastructureSummaryLicenseModelEnumValues() []AutonomousExadataInfrastructureSummaryLicenseModelEnum {
	values := make([]AutonomousExadataInfrastructureSummaryLicenseModelEnum, 0)
	for _, v := range mappingAutonomousExadataInfrastructureSummaryLicenseModel {
		values = append(values, v)
	}
	return values
}
