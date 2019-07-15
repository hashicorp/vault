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

// LaunchAutonomousExadataInfrastructureDetails Describes the input parameters to launch a new Autonomous Exadata Infrastructure.
type LaunchAutonomousExadataInfrastructureDetails struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment the Autonomous Exadata Infrastructure belongs in.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The availability domain where the Autonomous Exadata Infrastructure is located.
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the subnet the Autonomous Exadata Infrastructure is associated with.
	// **Subnet Restrictions:**
	// - For Autonomous Exadata Infrastructures, do not use a subnet that overlaps with 192.168.128.0/20
	// These subnets are used by the Oracle Clusterware private interconnect on the database instance.
	// Specifying an overlapping subnet will cause the private interconnect to malfunction.
	// This restriction applies to both the client subnet and backup subnet.
	SubnetId *string `mandatory:"true" json:"subnetId"`

	// The shape of the Autonomous Exadata Infrastructure. The shape determines resources allocated to the Autonomous Exadata Infrastructure (CPU cores, memory and storage). To get a list of shapes, use the ListDbSystemShapes operation.
	Shape *string `mandatory:"true" json:"shape"`

	// The user-friendly name for the Autonomous Exadata Infrastructure. It does not have to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// A domain name used for the Autonomous Exadata Infrastructure. If the Oracle-provided Internet and VCN
	// Resolver is enabled for the specified subnet, the domain name for the subnet is used
	// (don't provide one). Otherwise, provide a valid DNS domain name. Hyphens (-) are not permitted.
	Domain *string `mandatory:"false" json:"domain"`

	// The Oracle license model that applies to all the databases in the Autonomous Exadata Infrastructure. The default is BRING_YOUR_OWN_LICENSE.
	LicenseModel LaunchAutonomousExadataInfrastructureDetailsLicenseModelEnum `mandatory:"false" json:"licenseModel,omitempty"`

	MaintenanceWindowDetails *MaintenanceWindow `mandatory:"false" json:"maintenanceWindowDetails"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m LaunchAutonomousExadataInfrastructureDetails) String() string {
	return common.PointerString(m)
}

// LaunchAutonomousExadataInfrastructureDetailsLicenseModelEnum Enum with underlying type: string
type LaunchAutonomousExadataInfrastructureDetailsLicenseModelEnum string

// Set of constants representing the allowable values for LaunchAutonomousExadataInfrastructureDetailsLicenseModelEnum
const (
	LaunchAutonomousExadataInfrastructureDetailsLicenseModelLicenseIncluded     LaunchAutonomousExadataInfrastructureDetailsLicenseModelEnum = "LICENSE_INCLUDED"
	LaunchAutonomousExadataInfrastructureDetailsLicenseModelBringYourOwnLicense LaunchAutonomousExadataInfrastructureDetailsLicenseModelEnum = "BRING_YOUR_OWN_LICENSE"
)

var mappingLaunchAutonomousExadataInfrastructureDetailsLicenseModel = map[string]LaunchAutonomousExadataInfrastructureDetailsLicenseModelEnum{
	"LICENSE_INCLUDED":       LaunchAutonomousExadataInfrastructureDetailsLicenseModelLicenseIncluded,
	"BRING_YOUR_OWN_LICENSE": LaunchAutonomousExadataInfrastructureDetailsLicenseModelBringYourOwnLicense,
}

// GetLaunchAutonomousExadataInfrastructureDetailsLicenseModelEnumValues Enumerates the set of values for LaunchAutonomousExadataInfrastructureDetailsLicenseModelEnum
func GetLaunchAutonomousExadataInfrastructureDetailsLicenseModelEnumValues() []LaunchAutonomousExadataInfrastructureDetailsLicenseModelEnum {
	values := make([]LaunchAutonomousExadataInfrastructureDetailsLicenseModelEnum, 0)
	for _, v := range mappingLaunchAutonomousExadataInfrastructureDetailsLicenseModel {
		values = append(values, v)
	}
	return values
}
