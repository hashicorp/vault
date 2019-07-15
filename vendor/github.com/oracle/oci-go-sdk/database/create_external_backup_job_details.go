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

// CreateExternalBackupJobDetails The representation of CreateExternalBackupJobDetails
type CreateExternalBackupJobDetails struct {

	// The targeted availability domain for the backup.
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment where this backup should be created.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name for the backup. This name does not have to be unique.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// A valid Oracle Database version.
	DbVersion *string `mandatory:"true" json:"dbVersion"`

	// The name of the database from which the backup is being taken.
	DbName *string `mandatory:"true" json:"dbName"`

	// The `DBID` of the Oracle Database being backed up.
	ExternalDatabaseIdentifier *int64 `mandatory:"true" json:"externalDatabaseIdentifier"`

	// The character set for the database.
	CharacterSet *string `mandatory:"true" json:"characterSet"`

	// The national character set for the database.
	NcharacterSet *string `mandatory:"true" json:"ncharacterSet"`

	// The mode (single instance or RAC) of the database being backed up.
	DatabaseMode CreateExternalBackupJobDetailsDatabaseModeEnum `mandatory:"true" json:"databaseMode"`

	// The Oracle Database edition to use for creating a database from this standalone backup.
	// Note that 2-node RAC DB systems require Enterprise Edition - Extreme Performance.
	DatabaseEdition CreateExternalBackupJobDetailsDatabaseEditionEnum `mandatory:"true" json:"databaseEdition"`

	// The `DB_UNIQUE_NAME` of the Oracle Database being backed up.
	DbUniqueName *string `mandatory:"false" json:"dbUniqueName"`

	// The pluggable database name.
	PdbName *string `mandatory:"false" json:"pdbName"`
}

func (m CreateExternalBackupJobDetails) String() string {
	return common.PointerString(m)
}

// CreateExternalBackupJobDetailsDatabaseModeEnum Enum with underlying type: string
type CreateExternalBackupJobDetailsDatabaseModeEnum string

// Set of constants representing the allowable values for CreateExternalBackupJobDetailsDatabaseModeEnum
const (
	CreateExternalBackupJobDetailsDatabaseModeSi  CreateExternalBackupJobDetailsDatabaseModeEnum = "SI"
	CreateExternalBackupJobDetailsDatabaseModeRac CreateExternalBackupJobDetailsDatabaseModeEnum = "RAC"
)

var mappingCreateExternalBackupJobDetailsDatabaseMode = map[string]CreateExternalBackupJobDetailsDatabaseModeEnum{
	"SI":  CreateExternalBackupJobDetailsDatabaseModeSi,
	"RAC": CreateExternalBackupJobDetailsDatabaseModeRac,
}

// GetCreateExternalBackupJobDetailsDatabaseModeEnumValues Enumerates the set of values for CreateExternalBackupJobDetailsDatabaseModeEnum
func GetCreateExternalBackupJobDetailsDatabaseModeEnumValues() []CreateExternalBackupJobDetailsDatabaseModeEnum {
	values := make([]CreateExternalBackupJobDetailsDatabaseModeEnum, 0)
	for _, v := range mappingCreateExternalBackupJobDetailsDatabaseMode {
		values = append(values, v)
	}
	return values
}

// CreateExternalBackupJobDetailsDatabaseEditionEnum Enum with underlying type: string
type CreateExternalBackupJobDetailsDatabaseEditionEnum string

// Set of constants representing the allowable values for CreateExternalBackupJobDetailsDatabaseEditionEnum
const (
	CreateExternalBackupJobDetailsDatabaseEditionStandardEdition                     CreateExternalBackupJobDetailsDatabaseEditionEnum = "STANDARD_EDITION"
	CreateExternalBackupJobDetailsDatabaseEditionEnterpriseEdition                   CreateExternalBackupJobDetailsDatabaseEditionEnum = "ENTERPRISE_EDITION"
	CreateExternalBackupJobDetailsDatabaseEditionEnterpriseEditionHighPerformance    CreateExternalBackupJobDetailsDatabaseEditionEnum = "ENTERPRISE_EDITION_HIGH_PERFORMANCE"
	CreateExternalBackupJobDetailsDatabaseEditionEnterpriseEditionExtremePerformance CreateExternalBackupJobDetailsDatabaseEditionEnum = "ENTERPRISE_EDITION_EXTREME_PERFORMANCE"
)

var mappingCreateExternalBackupJobDetailsDatabaseEdition = map[string]CreateExternalBackupJobDetailsDatabaseEditionEnum{
	"STANDARD_EDITION":                       CreateExternalBackupJobDetailsDatabaseEditionStandardEdition,
	"ENTERPRISE_EDITION":                     CreateExternalBackupJobDetailsDatabaseEditionEnterpriseEdition,
	"ENTERPRISE_EDITION_HIGH_PERFORMANCE":    CreateExternalBackupJobDetailsDatabaseEditionEnterpriseEditionHighPerformance,
	"ENTERPRISE_EDITION_EXTREME_PERFORMANCE": CreateExternalBackupJobDetailsDatabaseEditionEnterpriseEditionExtremePerformance,
}

// GetCreateExternalBackupJobDetailsDatabaseEditionEnumValues Enumerates the set of values for CreateExternalBackupJobDetailsDatabaseEditionEnum
func GetCreateExternalBackupJobDetailsDatabaseEditionEnumValues() []CreateExternalBackupJobDetailsDatabaseEditionEnum {
	values := make([]CreateExternalBackupJobDetailsDatabaseEditionEnum, 0)
	for _, v := range mappingCreateExternalBackupJobDetailsDatabaseEdition {
		values = append(values, v)
	}
	return values
}
