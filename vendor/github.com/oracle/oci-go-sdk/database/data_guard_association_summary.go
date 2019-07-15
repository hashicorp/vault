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

// DataGuardAssociationSummary The properties that define a Data Guard association.
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized, talk to an
// administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// For information about endpoints and signing API requests, see
// About the API (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm). For information about available SDKs and tools, see
// SDKS and Other Tools (https://docs.cloud.oracle.com/Content/API/Concepts/sdks.htm).
type DataGuardAssociationSummary struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Data Guard association.
	Id *string `mandatory:"true" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the reporting database.
	DatabaseId *string `mandatory:"true" json:"databaseId"`

	// The role of the reporting database in this Data Guard association.
	Role DataGuardAssociationSummaryRoleEnum `mandatory:"true" json:"role"`

	// The current state of the Data Guard association.
	LifecycleState DataGuardAssociationSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the DB system containing the associated
	// peer database.
	PeerDbSystemId *string `mandatory:"true" json:"peerDbSystemId"`

	// The role of the peer database in this Data Guard association.
	PeerRole DataGuardAssociationSummaryPeerRoleEnum `mandatory:"true" json:"peerRole"`

	// The protection mode of this Data Guard association. For more information, see
	// Oracle Data Guard Protection Modes (http://docs.oracle.com/database/122/SBYDB/oracle-data-guard-protection-modes.htm#SBYDB02000)
	// in the Oracle Data Guard documentation.
	ProtectionMode DataGuardAssociationSummaryProtectionModeEnum `mandatory:"true" json:"protectionMode"`

	// Additional information about the current lifecycleState, if available.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the database home containing the associated peer database.
	PeerDbHomeId *string `mandatory:"false" json:"peerDbHomeId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the associated peer database.
	PeerDatabaseId *string `mandatory:"false" json:"peerDatabaseId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the peer database's Data Guard association.
	PeerDataGuardAssociationId *string `mandatory:"false" json:"peerDataGuardAssociationId"`

	// The lag time between updates to the primary database and application of the redo data on the standby database,
	// as computed by the reporting database.
	// Example: `9 seconds`
	ApplyLag *string `mandatory:"false" json:"applyLag"`

	// The rate at which redo logs are synced between the associated databases.
	// Example: `180 Mb per second`
	ApplyRate *string `mandatory:"false" json:"applyRate"`

	// The redo transport type used by this Data Guard association.  For more information, see
	// Redo Transport Services (http://docs.oracle.com/database/122/SBYDB/oracle-data-guard-redo-transport-services.htm#SBYDB00400)
	// in the Oracle Data Guard documentation.
	TransportType DataGuardAssociationSummaryTransportTypeEnum `mandatory:"false" json:"transportType,omitempty"`

	// The date and time the Data Guard association was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m DataGuardAssociationSummary) String() string {
	return common.PointerString(m)
}

// DataGuardAssociationSummaryRoleEnum Enum with underlying type: string
type DataGuardAssociationSummaryRoleEnum string

// Set of constants representing the allowable values for DataGuardAssociationSummaryRoleEnum
const (
	DataGuardAssociationSummaryRolePrimary         DataGuardAssociationSummaryRoleEnum = "PRIMARY"
	DataGuardAssociationSummaryRoleStandby         DataGuardAssociationSummaryRoleEnum = "STANDBY"
	DataGuardAssociationSummaryRoleDisabledStandby DataGuardAssociationSummaryRoleEnum = "DISABLED_STANDBY"
)

var mappingDataGuardAssociationSummaryRole = map[string]DataGuardAssociationSummaryRoleEnum{
	"PRIMARY":          DataGuardAssociationSummaryRolePrimary,
	"STANDBY":          DataGuardAssociationSummaryRoleStandby,
	"DISABLED_STANDBY": DataGuardAssociationSummaryRoleDisabledStandby,
}

// GetDataGuardAssociationSummaryRoleEnumValues Enumerates the set of values for DataGuardAssociationSummaryRoleEnum
func GetDataGuardAssociationSummaryRoleEnumValues() []DataGuardAssociationSummaryRoleEnum {
	values := make([]DataGuardAssociationSummaryRoleEnum, 0)
	for _, v := range mappingDataGuardAssociationSummaryRole {
		values = append(values, v)
	}
	return values
}

// DataGuardAssociationSummaryLifecycleStateEnum Enum with underlying type: string
type DataGuardAssociationSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for DataGuardAssociationSummaryLifecycleStateEnum
const (
	DataGuardAssociationSummaryLifecycleStateProvisioning DataGuardAssociationSummaryLifecycleStateEnum = "PROVISIONING"
	DataGuardAssociationSummaryLifecycleStateAvailable    DataGuardAssociationSummaryLifecycleStateEnum = "AVAILABLE"
	DataGuardAssociationSummaryLifecycleStateUpdating     DataGuardAssociationSummaryLifecycleStateEnum = "UPDATING"
	DataGuardAssociationSummaryLifecycleStateTerminating  DataGuardAssociationSummaryLifecycleStateEnum = "TERMINATING"
	DataGuardAssociationSummaryLifecycleStateTerminated   DataGuardAssociationSummaryLifecycleStateEnum = "TERMINATED"
	DataGuardAssociationSummaryLifecycleStateFailed       DataGuardAssociationSummaryLifecycleStateEnum = "FAILED"
)

var mappingDataGuardAssociationSummaryLifecycleState = map[string]DataGuardAssociationSummaryLifecycleStateEnum{
	"PROVISIONING": DataGuardAssociationSummaryLifecycleStateProvisioning,
	"AVAILABLE":    DataGuardAssociationSummaryLifecycleStateAvailable,
	"UPDATING":     DataGuardAssociationSummaryLifecycleStateUpdating,
	"TERMINATING":  DataGuardAssociationSummaryLifecycleStateTerminating,
	"TERMINATED":   DataGuardAssociationSummaryLifecycleStateTerminated,
	"FAILED":       DataGuardAssociationSummaryLifecycleStateFailed,
}

// GetDataGuardAssociationSummaryLifecycleStateEnumValues Enumerates the set of values for DataGuardAssociationSummaryLifecycleStateEnum
func GetDataGuardAssociationSummaryLifecycleStateEnumValues() []DataGuardAssociationSummaryLifecycleStateEnum {
	values := make([]DataGuardAssociationSummaryLifecycleStateEnum, 0)
	for _, v := range mappingDataGuardAssociationSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}

// DataGuardAssociationSummaryPeerRoleEnum Enum with underlying type: string
type DataGuardAssociationSummaryPeerRoleEnum string

// Set of constants representing the allowable values for DataGuardAssociationSummaryPeerRoleEnum
const (
	DataGuardAssociationSummaryPeerRolePrimary         DataGuardAssociationSummaryPeerRoleEnum = "PRIMARY"
	DataGuardAssociationSummaryPeerRoleStandby         DataGuardAssociationSummaryPeerRoleEnum = "STANDBY"
	DataGuardAssociationSummaryPeerRoleDisabledStandby DataGuardAssociationSummaryPeerRoleEnum = "DISABLED_STANDBY"
)

var mappingDataGuardAssociationSummaryPeerRole = map[string]DataGuardAssociationSummaryPeerRoleEnum{
	"PRIMARY":          DataGuardAssociationSummaryPeerRolePrimary,
	"STANDBY":          DataGuardAssociationSummaryPeerRoleStandby,
	"DISABLED_STANDBY": DataGuardAssociationSummaryPeerRoleDisabledStandby,
}

// GetDataGuardAssociationSummaryPeerRoleEnumValues Enumerates the set of values for DataGuardAssociationSummaryPeerRoleEnum
func GetDataGuardAssociationSummaryPeerRoleEnumValues() []DataGuardAssociationSummaryPeerRoleEnum {
	values := make([]DataGuardAssociationSummaryPeerRoleEnum, 0)
	for _, v := range mappingDataGuardAssociationSummaryPeerRole {
		values = append(values, v)
	}
	return values
}

// DataGuardAssociationSummaryProtectionModeEnum Enum with underlying type: string
type DataGuardAssociationSummaryProtectionModeEnum string

// Set of constants representing the allowable values for DataGuardAssociationSummaryProtectionModeEnum
const (
	DataGuardAssociationSummaryProtectionModeAvailability DataGuardAssociationSummaryProtectionModeEnum = "MAXIMUM_AVAILABILITY"
	DataGuardAssociationSummaryProtectionModePerformance  DataGuardAssociationSummaryProtectionModeEnum = "MAXIMUM_PERFORMANCE"
	DataGuardAssociationSummaryProtectionModeProtection   DataGuardAssociationSummaryProtectionModeEnum = "MAXIMUM_PROTECTION"
)

var mappingDataGuardAssociationSummaryProtectionMode = map[string]DataGuardAssociationSummaryProtectionModeEnum{
	"MAXIMUM_AVAILABILITY": DataGuardAssociationSummaryProtectionModeAvailability,
	"MAXIMUM_PERFORMANCE":  DataGuardAssociationSummaryProtectionModePerformance,
	"MAXIMUM_PROTECTION":   DataGuardAssociationSummaryProtectionModeProtection,
}

// GetDataGuardAssociationSummaryProtectionModeEnumValues Enumerates the set of values for DataGuardAssociationSummaryProtectionModeEnum
func GetDataGuardAssociationSummaryProtectionModeEnumValues() []DataGuardAssociationSummaryProtectionModeEnum {
	values := make([]DataGuardAssociationSummaryProtectionModeEnum, 0)
	for _, v := range mappingDataGuardAssociationSummaryProtectionMode {
		values = append(values, v)
	}
	return values
}

// DataGuardAssociationSummaryTransportTypeEnum Enum with underlying type: string
type DataGuardAssociationSummaryTransportTypeEnum string

// Set of constants representing the allowable values for DataGuardAssociationSummaryTransportTypeEnum
const (
	DataGuardAssociationSummaryTransportTypeSync     DataGuardAssociationSummaryTransportTypeEnum = "SYNC"
	DataGuardAssociationSummaryTransportTypeAsync    DataGuardAssociationSummaryTransportTypeEnum = "ASYNC"
	DataGuardAssociationSummaryTransportTypeFastsync DataGuardAssociationSummaryTransportTypeEnum = "FASTSYNC"
)

var mappingDataGuardAssociationSummaryTransportType = map[string]DataGuardAssociationSummaryTransportTypeEnum{
	"SYNC":     DataGuardAssociationSummaryTransportTypeSync,
	"ASYNC":    DataGuardAssociationSummaryTransportTypeAsync,
	"FASTSYNC": DataGuardAssociationSummaryTransportTypeFastsync,
}

// GetDataGuardAssociationSummaryTransportTypeEnumValues Enumerates the set of values for DataGuardAssociationSummaryTransportTypeEnum
func GetDataGuardAssociationSummaryTransportTypeEnumValues() []DataGuardAssociationSummaryTransportTypeEnum {
	values := make([]DataGuardAssociationSummaryTransportTypeEnum, 0)
	for _, v := range mappingDataGuardAssociationSummaryTransportType {
		values = append(values, v)
	}
	return values
}
