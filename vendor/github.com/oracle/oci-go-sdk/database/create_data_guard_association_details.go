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

// CreateDataGuardAssociationDetails The configuration details for creating a Data Guard association between databases.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type CreateDataGuardAssociationDetails interface {

	// A strong password for the `SYS`, `SYSTEM`, and `PDB Admin` users to apply during standby creation.
	// The password must contain no fewer than nine characters and include:
	// * At least two uppercase characters.
	// * At least two lowercase characters.
	// * At least two numeric characters.
	// * At least two special characters. Valid special characters include "_", "#", and "-" only.
	// **The password MUST be the same as the primary admin password.**
	GetDatabaseAdminPassword() *string

	// The protection mode to set up between the primary and standby databases. For more information, see
	// Oracle Data Guard Protection Modes (http://docs.oracle.com/database/122/SBYDB/oracle-data-guard-protection-modes.htm#SBYDB02000)
	// in the Oracle Data Guard documentation.
	// **IMPORTANT** - The only protection mode currently supported by the Database service is MAXIMUM_PERFORMANCE.
	GetProtectionMode() CreateDataGuardAssociationDetailsProtectionModeEnum

	// The redo transport type to use for this Data Guard association.  Valid values depend on the specified `protectionMode`:
	// * MAXIMUM_AVAILABILITY - SYNC or FASTSYNC
	// * MAXIMUM_PERFORMANCE - ASYNC
	// * MAXIMUM_PROTECTION - SYNC
	// For more information, see
	// Redo Transport Services (http://docs.oracle.com/database/122/SBYDB/oracle-data-guard-redo-transport-services.htm#SBYDB00400)
	// in the Oracle Data Guard documentation.
	// **IMPORTANT** - The only transport type currently supported by the Database service is ASYNC.
	GetTransportType() CreateDataGuardAssociationDetailsTransportTypeEnum
}

type createdataguardassociationdetails struct {
	JsonData              []byte
	DatabaseAdminPassword *string                                             `mandatory:"true" json:"databaseAdminPassword"`
	ProtectionMode        CreateDataGuardAssociationDetailsProtectionModeEnum `mandatory:"true" json:"protectionMode"`
	TransportType         CreateDataGuardAssociationDetailsTransportTypeEnum  `mandatory:"true" json:"transportType"`
	CreationType          string                                              `json:"creationType"`
}

// UnmarshalJSON unmarshals json
func (m *createdataguardassociationdetails) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalercreatedataguardassociationdetails createdataguardassociationdetails
	s := struct {
		Model Unmarshalercreatedataguardassociationdetails
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.DatabaseAdminPassword = s.Model.DatabaseAdminPassword
	m.ProtectionMode = s.Model.ProtectionMode
	m.TransportType = s.Model.TransportType
	m.CreationType = s.Model.CreationType

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *createdataguardassociationdetails) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.CreationType {
	case "NewDbSystem":
		mm := CreateDataGuardAssociationWithNewDbSystemDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "ExistingDbSystem":
		mm := CreateDataGuardAssociationToExistingDbSystemDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

//GetDatabaseAdminPassword returns DatabaseAdminPassword
func (m createdataguardassociationdetails) GetDatabaseAdminPassword() *string {
	return m.DatabaseAdminPassword
}

//GetProtectionMode returns ProtectionMode
func (m createdataguardassociationdetails) GetProtectionMode() CreateDataGuardAssociationDetailsProtectionModeEnum {
	return m.ProtectionMode
}

//GetTransportType returns TransportType
func (m createdataguardassociationdetails) GetTransportType() CreateDataGuardAssociationDetailsTransportTypeEnum {
	return m.TransportType
}

func (m createdataguardassociationdetails) String() string {
	return common.PointerString(m)
}

// CreateDataGuardAssociationDetailsProtectionModeEnum Enum with underlying type: string
type CreateDataGuardAssociationDetailsProtectionModeEnum string

// Set of constants representing the allowable values for CreateDataGuardAssociationDetailsProtectionModeEnum
const (
	CreateDataGuardAssociationDetailsProtectionModeAvailability CreateDataGuardAssociationDetailsProtectionModeEnum = "MAXIMUM_AVAILABILITY"
	CreateDataGuardAssociationDetailsProtectionModePerformance  CreateDataGuardAssociationDetailsProtectionModeEnum = "MAXIMUM_PERFORMANCE"
	CreateDataGuardAssociationDetailsProtectionModeProtection   CreateDataGuardAssociationDetailsProtectionModeEnum = "MAXIMUM_PROTECTION"
)

var mappingCreateDataGuardAssociationDetailsProtectionMode = map[string]CreateDataGuardAssociationDetailsProtectionModeEnum{
	"MAXIMUM_AVAILABILITY": CreateDataGuardAssociationDetailsProtectionModeAvailability,
	"MAXIMUM_PERFORMANCE":  CreateDataGuardAssociationDetailsProtectionModePerformance,
	"MAXIMUM_PROTECTION":   CreateDataGuardAssociationDetailsProtectionModeProtection,
}

// GetCreateDataGuardAssociationDetailsProtectionModeEnumValues Enumerates the set of values for CreateDataGuardAssociationDetailsProtectionModeEnum
func GetCreateDataGuardAssociationDetailsProtectionModeEnumValues() []CreateDataGuardAssociationDetailsProtectionModeEnum {
	values := make([]CreateDataGuardAssociationDetailsProtectionModeEnum, 0)
	for _, v := range mappingCreateDataGuardAssociationDetailsProtectionMode {
		values = append(values, v)
	}
	return values
}

// CreateDataGuardAssociationDetailsTransportTypeEnum Enum with underlying type: string
type CreateDataGuardAssociationDetailsTransportTypeEnum string

// Set of constants representing the allowable values for CreateDataGuardAssociationDetailsTransportTypeEnum
const (
	CreateDataGuardAssociationDetailsTransportTypeSync     CreateDataGuardAssociationDetailsTransportTypeEnum = "SYNC"
	CreateDataGuardAssociationDetailsTransportTypeAsync    CreateDataGuardAssociationDetailsTransportTypeEnum = "ASYNC"
	CreateDataGuardAssociationDetailsTransportTypeFastsync CreateDataGuardAssociationDetailsTransportTypeEnum = "FASTSYNC"
)

var mappingCreateDataGuardAssociationDetailsTransportType = map[string]CreateDataGuardAssociationDetailsTransportTypeEnum{
	"SYNC":     CreateDataGuardAssociationDetailsTransportTypeSync,
	"ASYNC":    CreateDataGuardAssociationDetailsTransportTypeAsync,
	"FASTSYNC": CreateDataGuardAssociationDetailsTransportTypeFastsync,
}

// GetCreateDataGuardAssociationDetailsTransportTypeEnumValues Enumerates the set of values for CreateDataGuardAssociationDetailsTransportTypeEnum
func GetCreateDataGuardAssociationDetailsTransportTypeEnumValues() []CreateDataGuardAssociationDetailsTransportTypeEnum {
	values := make([]CreateDataGuardAssociationDetailsTransportTypeEnum, 0)
	for _, v := range mappingCreateDataGuardAssociationDetailsTransportType {
		values = append(values, v)
	}
	return values
}
