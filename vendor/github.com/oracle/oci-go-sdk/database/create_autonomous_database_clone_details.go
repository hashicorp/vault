// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// CreateAutonomousDatabaseCloneDetails Details to create an Oracle Autonomous Database by cloning an existing Autonomous Database.
type CreateAutonomousDatabaseCloneDetails struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment of the autonomous database.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The database name. The name must begin with an alphabetic character and can contain a maximum of 14 alphanumeric characters. Special characters are not permitted. The database name must be unique in the tenancy.
	DbName *string `mandatory:"true" json:"dbName"`

	// The number of CPU Cores to be made available to the database.
	CpuCoreCount *int `mandatory:"true" json:"cpuCoreCount"`

	// The size, in terabytes, of the data volume that will be created and attached to the database. This storage can later be scaled up if needed.
	DataStorageSizeInTBs *int `mandatory:"true" json:"dataStorageSizeInTBs"`

	// The password must be between 12 and 30 characters long, and must contain at least 1 uppercase, 1 lowercase, and 1 numeric character. It cannot contain the double quote symbol (") or the username "admin", regardless of casing.
	AdminPassword *string `mandatory:"true" json:"adminPassword"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the source Autonomous Database that you will clone to create a new Autonomous Database.
	SourceId *string `mandatory:"true" json:"sourceId"`

	// The user-friendly name for the Autonomous Database. The name does not have to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// If set to true, indicates that an Autonomous Database preview version is being provisioned, and that the preview version's terms of service have been accepted.
	IsPreviewVersionWithServiceTermsAccepted *bool `mandatory:"false" json:"isPreviewVersionWithServiceTermsAccepted"`

	// Indicates if auto scaling is enabled for the Autonomous Database CPU core count. The default value is false.
	IsAutoScalingEnabled *bool `mandatory:"false" json:"isAutoScalingEnabled"`

	// True if it is dedicated database.
	IsDedicated *bool `mandatory:"false" json:"isDedicated"`

	// The Autonomous Container Database OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	AutonomousContainerDatabaseId *string `mandatory:"false" json:"autonomousContainerDatabaseId"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The clone type.
	CloneType CreateAutonomousDatabaseCloneDetailsCloneTypeEnum `mandatory:"true" json:"cloneType"`

	// The autonomous database workload type. OLTP indicates an Autonomous Transaction Processing database and DW indicates an Autonomous Data Warehouse. The default is OLTP.
	DbWorkload CreateAutonomousDatabaseBaseDbWorkloadEnum `mandatory:"false" json:"dbWorkload,omitempty"`

	// The Oracle license model that applies to the Oracle Autonomous Database. The default is BRING_YOUR_OWN_LICENSE.
	LicenseModel CreateAutonomousDatabaseBaseLicenseModelEnum `mandatory:"false" json:"licenseModel,omitempty"`
}

//GetCompartmentId returns CompartmentId
func (m CreateAutonomousDatabaseCloneDetails) GetCompartmentId() *string {
	return m.CompartmentId
}

//GetDbName returns DbName
func (m CreateAutonomousDatabaseCloneDetails) GetDbName() *string {
	return m.DbName
}

//GetCpuCoreCount returns CpuCoreCount
func (m CreateAutonomousDatabaseCloneDetails) GetCpuCoreCount() *int {
	return m.CpuCoreCount
}

//GetDbWorkload returns DbWorkload
func (m CreateAutonomousDatabaseCloneDetails) GetDbWorkload() CreateAutonomousDatabaseBaseDbWorkloadEnum {
	return m.DbWorkload
}

//GetDataStorageSizeInTBs returns DataStorageSizeInTBs
func (m CreateAutonomousDatabaseCloneDetails) GetDataStorageSizeInTBs() *int {
	return m.DataStorageSizeInTBs
}

//GetAdminPassword returns AdminPassword
func (m CreateAutonomousDatabaseCloneDetails) GetAdminPassword() *string {
	return m.AdminPassword
}

//GetDisplayName returns DisplayName
func (m CreateAutonomousDatabaseCloneDetails) GetDisplayName() *string {
	return m.DisplayName
}

//GetLicenseModel returns LicenseModel
func (m CreateAutonomousDatabaseCloneDetails) GetLicenseModel() CreateAutonomousDatabaseBaseLicenseModelEnum {
	return m.LicenseModel
}

//GetIsPreviewVersionWithServiceTermsAccepted returns IsPreviewVersionWithServiceTermsAccepted
func (m CreateAutonomousDatabaseCloneDetails) GetIsPreviewVersionWithServiceTermsAccepted() *bool {
	return m.IsPreviewVersionWithServiceTermsAccepted
}

//GetIsAutoScalingEnabled returns IsAutoScalingEnabled
func (m CreateAutonomousDatabaseCloneDetails) GetIsAutoScalingEnabled() *bool {
	return m.IsAutoScalingEnabled
}

//GetIsDedicated returns IsDedicated
func (m CreateAutonomousDatabaseCloneDetails) GetIsDedicated() *bool {
	return m.IsDedicated
}

//GetAutonomousContainerDatabaseId returns AutonomousContainerDatabaseId
func (m CreateAutonomousDatabaseCloneDetails) GetAutonomousContainerDatabaseId() *string {
	return m.AutonomousContainerDatabaseId
}

//GetFreeformTags returns FreeformTags
func (m CreateAutonomousDatabaseCloneDetails) GetFreeformTags() map[string]string {
	return m.FreeformTags
}

//GetDefinedTags returns DefinedTags
func (m CreateAutonomousDatabaseCloneDetails) GetDefinedTags() map[string]map[string]interface{} {
	return m.DefinedTags
}

func (m CreateAutonomousDatabaseCloneDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m CreateAutonomousDatabaseCloneDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeCreateAutonomousDatabaseCloneDetails CreateAutonomousDatabaseCloneDetails
	s := struct {
		DiscriminatorParam string `json:"source"`
		MarshalTypeCreateAutonomousDatabaseCloneDetails
	}{
		"DATABASE",
		(MarshalTypeCreateAutonomousDatabaseCloneDetails)(m),
	}

	return json.Marshal(&s)
}

// CreateAutonomousDatabaseCloneDetailsCloneTypeEnum Enum with underlying type: string
type CreateAutonomousDatabaseCloneDetailsCloneTypeEnum string

// Set of constants representing the allowable values for CreateAutonomousDatabaseCloneDetailsCloneTypeEnum
const (
	CreateAutonomousDatabaseCloneDetailsCloneTypeFull     CreateAutonomousDatabaseCloneDetailsCloneTypeEnum = "FULL"
	CreateAutonomousDatabaseCloneDetailsCloneTypeMetadata CreateAutonomousDatabaseCloneDetailsCloneTypeEnum = "METADATA"
)

var mappingCreateAutonomousDatabaseCloneDetailsCloneType = map[string]CreateAutonomousDatabaseCloneDetailsCloneTypeEnum{
	"FULL":     CreateAutonomousDatabaseCloneDetailsCloneTypeFull,
	"METADATA": CreateAutonomousDatabaseCloneDetailsCloneTypeMetadata,
}

// GetCreateAutonomousDatabaseCloneDetailsCloneTypeEnumValues Enumerates the set of values for CreateAutonomousDatabaseCloneDetailsCloneTypeEnum
func GetCreateAutonomousDatabaseCloneDetailsCloneTypeEnumValues() []CreateAutonomousDatabaseCloneDetailsCloneTypeEnum {
	values := make([]CreateAutonomousDatabaseCloneDetailsCloneTypeEnum, 0)
	for _, v := range mappingCreateAutonomousDatabaseCloneDetailsCloneType {
		values = append(values, v)
	}
	return values
}
