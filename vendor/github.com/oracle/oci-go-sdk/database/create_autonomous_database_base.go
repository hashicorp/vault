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

// CreateAutonomousDatabaseBase Details to create an Oracle Autonomous Database.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type CreateAutonomousDatabaseBase interface {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment of the autonomous database.
	GetCompartmentId() *string

	// The database name. The name must begin with an alphabetic character and can contain a maximum of 14 alphanumeric characters. Special characters are not permitted. The database name must be unique in the tenancy.
	GetDbName() *string

	// The number of CPU Cores to be made available to the database.
	GetCpuCoreCount() *int

	// The size, in terabytes, of the data volume that will be created and attached to the database. This storage can later be scaled up if needed.
	GetDataStorageSizeInTBs() *int

	// The password must be between 12 and 30 characters long, and must contain at least 1 uppercase, 1 lowercase, and 1 numeric character. It cannot contain the double quote symbol (") or the username "admin", regardless of casing.
	GetAdminPassword() *string

	// The autonomous database workload type. OLTP indicates an Autonomous Transaction Processing database and DW indicates an Autonomous Data Warehouse. The default is OLTP.
	GetDbWorkload() CreateAutonomousDatabaseBaseDbWorkloadEnum

	// The user-friendly name for the Autonomous Database. The name does not have to be unique.
	GetDisplayName() *string

	// The Oracle license model that applies to the Oracle Autonomous Database. The default is BRING_YOUR_OWN_LICENSE.
	GetLicenseModel() CreateAutonomousDatabaseBaseLicenseModelEnum

	// If set to true, indicates that an Autonomous Database preview version is being provisioned, and that the preview version's terms of service have been accepted.
	GetIsPreviewVersionWithServiceTermsAccepted() *bool

	// Indicates if auto scaling is enabled for the Autonomous Database CPU core count. The default value is false.
	GetIsAutoScalingEnabled() *bool

	// True if it is dedicated database.
	GetIsDedicated() *bool

	// The Autonomous Container Database OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	GetAutonomousContainerDatabaseId() *string

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	GetFreeformTags() map[string]string

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	GetDefinedTags() map[string]map[string]interface{}
}

type createautonomousdatabasebase struct {
	JsonData                                 []byte
	CompartmentId                            *string                                      `mandatory:"true" json:"compartmentId"`
	DbName                                   *string                                      `mandatory:"true" json:"dbName"`
	CpuCoreCount                             *int                                         `mandatory:"true" json:"cpuCoreCount"`
	DataStorageSizeInTBs                     *int                                         `mandatory:"true" json:"dataStorageSizeInTBs"`
	AdminPassword                            *string                                      `mandatory:"true" json:"adminPassword"`
	DbWorkload                               CreateAutonomousDatabaseBaseDbWorkloadEnum   `mandatory:"false" json:"dbWorkload,omitempty"`
	DisplayName                              *string                                      `mandatory:"false" json:"displayName"`
	LicenseModel                             CreateAutonomousDatabaseBaseLicenseModelEnum `mandatory:"false" json:"licenseModel,omitempty"`
	IsPreviewVersionWithServiceTermsAccepted *bool                                        `mandatory:"false" json:"isPreviewVersionWithServiceTermsAccepted"`
	IsAutoScalingEnabled                     *bool                                        `mandatory:"false" json:"isAutoScalingEnabled"`
	IsDedicated                              *bool                                        `mandatory:"false" json:"isDedicated"`
	AutonomousContainerDatabaseId            *string                                      `mandatory:"false" json:"autonomousContainerDatabaseId"`
	FreeformTags                             map[string]string                            `mandatory:"false" json:"freeformTags"`
	DefinedTags                              map[string]map[string]interface{}            `mandatory:"false" json:"definedTags"`
	Source                                   string                                       `json:"source"`
}

// UnmarshalJSON unmarshals json
func (m *createautonomousdatabasebase) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalercreateautonomousdatabasebase createautonomousdatabasebase
	s := struct {
		Model Unmarshalercreateautonomousdatabasebase
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.CompartmentId = s.Model.CompartmentId
	m.DbName = s.Model.DbName
	m.CpuCoreCount = s.Model.CpuCoreCount
	m.DataStorageSizeInTBs = s.Model.DataStorageSizeInTBs
	m.AdminPassword = s.Model.AdminPassword
	m.DbWorkload = s.Model.DbWorkload
	m.DisplayName = s.Model.DisplayName
	m.LicenseModel = s.Model.LicenseModel
	m.IsPreviewVersionWithServiceTermsAccepted = s.Model.IsPreviewVersionWithServiceTermsAccepted
	m.IsAutoScalingEnabled = s.Model.IsAutoScalingEnabled
	m.IsDedicated = s.Model.IsDedicated
	m.AutonomousContainerDatabaseId = s.Model.AutonomousContainerDatabaseId
	m.FreeformTags = s.Model.FreeformTags
	m.DefinedTags = s.Model.DefinedTags
	m.Source = s.Model.Source

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *createautonomousdatabasebase) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.Source {
	case "DATABASE":
		mm := CreateAutonomousDatabaseCloneDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "NONE":
		mm := CreateAutonomousDatabaseDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

//GetCompartmentId returns CompartmentId
func (m createautonomousdatabasebase) GetCompartmentId() *string {
	return m.CompartmentId
}

//GetDbName returns DbName
func (m createautonomousdatabasebase) GetDbName() *string {
	return m.DbName
}

//GetCpuCoreCount returns CpuCoreCount
func (m createautonomousdatabasebase) GetCpuCoreCount() *int {
	return m.CpuCoreCount
}

//GetDataStorageSizeInTBs returns DataStorageSizeInTBs
func (m createautonomousdatabasebase) GetDataStorageSizeInTBs() *int {
	return m.DataStorageSizeInTBs
}

//GetAdminPassword returns AdminPassword
func (m createautonomousdatabasebase) GetAdminPassword() *string {
	return m.AdminPassword
}

//GetDbWorkload returns DbWorkload
func (m createautonomousdatabasebase) GetDbWorkload() CreateAutonomousDatabaseBaseDbWorkloadEnum {
	return m.DbWorkload
}

//GetDisplayName returns DisplayName
func (m createautonomousdatabasebase) GetDisplayName() *string {
	return m.DisplayName
}

//GetLicenseModel returns LicenseModel
func (m createautonomousdatabasebase) GetLicenseModel() CreateAutonomousDatabaseBaseLicenseModelEnum {
	return m.LicenseModel
}

//GetIsPreviewVersionWithServiceTermsAccepted returns IsPreviewVersionWithServiceTermsAccepted
func (m createautonomousdatabasebase) GetIsPreviewVersionWithServiceTermsAccepted() *bool {
	return m.IsPreviewVersionWithServiceTermsAccepted
}

//GetIsAutoScalingEnabled returns IsAutoScalingEnabled
func (m createautonomousdatabasebase) GetIsAutoScalingEnabled() *bool {
	return m.IsAutoScalingEnabled
}

//GetIsDedicated returns IsDedicated
func (m createautonomousdatabasebase) GetIsDedicated() *bool {
	return m.IsDedicated
}

//GetAutonomousContainerDatabaseId returns AutonomousContainerDatabaseId
func (m createautonomousdatabasebase) GetAutonomousContainerDatabaseId() *string {
	return m.AutonomousContainerDatabaseId
}

//GetFreeformTags returns FreeformTags
func (m createautonomousdatabasebase) GetFreeformTags() map[string]string {
	return m.FreeformTags
}

//GetDefinedTags returns DefinedTags
func (m createautonomousdatabasebase) GetDefinedTags() map[string]map[string]interface{} {
	return m.DefinedTags
}

func (m createautonomousdatabasebase) String() string {
	return common.PointerString(m)
}

// CreateAutonomousDatabaseBaseDbWorkloadEnum Enum with underlying type: string
type CreateAutonomousDatabaseBaseDbWorkloadEnum string

// Set of constants representing the allowable values for CreateAutonomousDatabaseBaseDbWorkloadEnum
const (
	CreateAutonomousDatabaseBaseDbWorkloadOltp CreateAutonomousDatabaseBaseDbWorkloadEnum = "OLTP"
	CreateAutonomousDatabaseBaseDbWorkloadDw   CreateAutonomousDatabaseBaseDbWorkloadEnum = "DW"
)

var mappingCreateAutonomousDatabaseBaseDbWorkload = map[string]CreateAutonomousDatabaseBaseDbWorkloadEnum{
	"OLTP": CreateAutonomousDatabaseBaseDbWorkloadOltp,
	"DW":   CreateAutonomousDatabaseBaseDbWorkloadDw,
}

// GetCreateAutonomousDatabaseBaseDbWorkloadEnumValues Enumerates the set of values for CreateAutonomousDatabaseBaseDbWorkloadEnum
func GetCreateAutonomousDatabaseBaseDbWorkloadEnumValues() []CreateAutonomousDatabaseBaseDbWorkloadEnum {
	values := make([]CreateAutonomousDatabaseBaseDbWorkloadEnum, 0)
	for _, v := range mappingCreateAutonomousDatabaseBaseDbWorkload {
		values = append(values, v)
	}
	return values
}

// CreateAutonomousDatabaseBaseLicenseModelEnum Enum with underlying type: string
type CreateAutonomousDatabaseBaseLicenseModelEnum string

// Set of constants representing the allowable values for CreateAutonomousDatabaseBaseLicenseModelEnum
const (
	CreateAutonomousDatabaseBaseLicenseModelLicenseIncluded     CreateAutonomousDatabaseBaseLicenseModelEnum = "LICENSE_INCLUDED"
	CreateAutonomousDatabaseBaseLicenseModelBringYourOwnLicense CreateAutonomousDatabaseBaseLicenseModelEnum = "BRING_YOUR_OWN_LICENSE"
)

var mappingCreateAutonomousDatabaseBaseLicenseModel = map[string]CreateAutonomousDatabaseBaseLicenseModelEnum{
	"LICENSE_INCLUDED":       CreateAutonomousDatabaseBaseLicenseModelLicenseIncluded,
	"BRING_YOUR_OWN_LICENSE": CreateAutonomousDatabaseBaseLicenseModelBringYourOwnLicense,
}

// GetCreateAutonomousDatabaseBaseLicenseModelEnumValues Enumerates the set of values for CreateAutonomousDatabaseBaseLicenseModelEnum
func GetCreateAutonomousDatabaseBaseLicenseModelEnumValues() []CreateAutonomousDatabaseBaseLicenseModelEnum {
	values := make([]CreateAutonomousDatabaseBaseLicenseModelEnum, 0)
	for _, v := range mappingCreateAutonomousDatabaseBaseLicenseModel {
		values = append(values, v)
	}
	return values
}
