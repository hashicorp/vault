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

// DataGuardAssociation The representation of DataGuardAssociation
type DataGuardAssociation struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Data Guard association.
	Id *string `mandatory:"true" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the reporting database.
	DatabaseId *string `mandatory:"true" json:"databaseId"`

	// The role of the reporting database in this Data Guard association.
	Role DataGuardAssociationRoleEnum `mandatory:"true" json:"role"`

	// The current state of the Data Guard association.
	LifecycleState DataGuardAssociationLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the DB system containing the associated
	// peer database.
	PeerDbSystemId *string `mandatory:"true" json:"peerDbSystemId"`

	// The role of the peer database in this Data Guard association.
	PeerRole DataGuardAssociationPeerRoleEnum `mandatory:"true" json:"peerRole"`

	// The protection mode of this Data Guard association. For more information, see
	// Oracle Data Guard Protection Modes (http://docs.oracle.com/database/122/SBYDB/oracle-data-guard-protection-modes.htm#SBYDB02000)
	// in the Oracle Data Guard documentation.
	ProtectionMode DataGuardAssociationProtectionModeEnum `mandatory:"true" json:"protectionMode"`

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
	TransportType DataGuardAssociationTransportTypeEnum `mandatory:"false" json:"transportType,omitempty"`

	// The date and time the Data Guard association was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m DataGuardAssociation) String() string {
	return common.PointerString(m)
}

// DataGuardAssociationRoleEnum Enum with underlying type: string
type DataGuardAssociationRoleEnum string

// Set of constants representing the allowable values for DataGuardAssociationRoleEnum
const (
	DataGuardAssociationRolePrimary         DataGuardAssociationRoleEnum = "PRIMARY"
	DataGuardAssociationRoleStandby         DataGuardAssociationRoleEnum = "STANDBY"
	DataGuardAssociationRoleDisabledStandby DataGuardAssociationRoleEnum = "DISABLED_STANDBY"
)

var mappingDataGuardAssociationRole = map[string]DataGuardAssociationRoleEnum{
	"PRIMARY":          DataGuardAssociationRolePrimary,
	"STANDBY":          DataGuardAssociationRoleStandby,
	"DISABLED_STANDBY": DataGuardAssociationRoleDisabledStandby,
}

// GetDataGuardAssociationRoleEnumValues Enumerates the set of values for DataGuardAssociationRoleEnum
func GetDataGuardAssociationRoleEnumValues() []DataGuardAssociationRoleEnum {
	values := make([]DataGuardAssociationRoleEnum, 0)
	for _, v := range mappingDataGuardAssociationRole {
		values = append(values, v)
	}
	return values
}

// DataGuardAssociationLifecycleStateEnum Enum with underlying type: string
type DataGuardAssociationLifecycleStateEnum string

// Set of constants representing the allowable values for DataGuardAssociationLifecycleStateEnum
const (
	DataGuardAssociationLifecycleStateProvisioning DataGuardAssociationLifecycleStateEnum = "PROVISIONING"
	DataGuardAssociationLifecycleStateAvailable    DataGuardAssociationLifecycleStateEnum = "AVAILABLE"
	DataGuardAssociationLifecycleStateUpdating     DataGuardAssociationLifecycleStateEnum = "UPDATING"
	DataGuardAssociationLifecycleStateTerminating  DataGuardAssociationLifecycleStateEnum = "TERMINATING"
	DataGuardAssociationLifecycleStateTerminated   DataGuardAssociationLifecycleStateEnum = "TERMINATED"
	DataGuardAssociationLifecycleStateFailed       DataGuardAssociationLifecycleStateEnum = "FAILED"
)

var mappingDataGuardAssociationLifecycleState = map[string]DataGuardAssociationLifecycleStateEnum{
	"PROVISIONING": DataGuardAssociationLifecycleStateProvisioning,
	"AVAILABLE":    DataGuardAssociationLifecycleStateAvailable,
	"UPDATING":     DataGuardAssociationLifecycleStateUpdating,
	"TERMINATING":  DataGuardAssociationLifecycleStateTerminating,
	"TERMINATED":   DataGuardAssociationLifecycleStateTerminated,
	"FAILED":       DataGuardAssociationLifecycleStateFailed,
}

// GetDataGuardAssociationLifecycleStateEnumValues Enumerates the set of values for DataGuardAssociationLifecycleStateEnum
func GetDataGuardAssociationLifecycleStateEnumValues() []DataGuardAssociationLifecycleStateEnum {
	values := make([]DataGuardAssociationLifecycleStateEnum, 0)
	for _, v := range mappingDataGuardAssociationLifecycleState {
		values = append(values, v)
	}
	return values
}

// DataGuardAssociationPeerRoleEnum Enum with underlying type: string
type DataGuardAssociationPeerRoleEnum string

// Set of constants representing the allowable values for DataGuardAssociationPeerRoleEnum
const (
	DataGuardAssociationPeerRolePrimary         DataGuardAssociationPeerRoleEnum = "PRIMARY"
	DataGuardAssociationPeerRoleStandby         DataGuardAssociationPeerRoleEnum = "STANDBY"
	DataGuardAssociationPeerRoleDisabledStandby DataGuardAssociationPeerRoleEnum = "DISABLED_STANDBY"
)

var mappingDataGuardAssociationPeerRole = map[string]DataGuardAssociationPeerRoleEnum{
	"PRIMARY":          DataGuardAssociationPeerRolePrimary,
	"STANDBY":          DataGuardAssociationPeerRoleStandby,
	"DISABLED_STANDBY": DataGuardAssociationPeerRoleDisabledStandby,
}

// GetDataGuardAssociationPeerRoleEnumValues Enumerates the set of values for DataGuardAssociationPeerRoleEnum
func GetDataGuardAssociationPeerRoleEnumValues() []DataGuardAssociationPeerRoleEnum {
	values := make([]DataGuardAssociationPeerRoleEnum, 0)
	for _, v := range mappingDataGuardAssociationPeerRole {
		values = append(values, v)
	}
	return values
}

// DataGuardAssociationProtectionModeEnum Enum with underlying type: string
type DataGuardAssociationProtectionModeEnum string

// Set of constants representing the allowable values for DataGuardAssociationProtectionModeEnum
const (
	DataGuardAssociationProtectionModeAvailability DataGuardAssociationProtectionModeEnum = "MAXIMUM_AVAILABILITY"
	DataGuardAssociationProtectionModePerformance  DataGuardAssociationProtectionModeEnum = "MAXIMUM_PERFORMANCE"
	DataGuardAssociationProtectionModeProtection   DataGuardAssociationProtectionModeEnum = "MAXIMUM_PROTECTION"
)

var mappingDataGuardAssociationProtectionMode = map[string]DataGuardAssociationProtectionModeEnum{
	"MAXIMUM_AVAILABILITY": DataGuardAssociationProtectionModeAvailability,
	"MAXIMUM_PERFORMANCE":  DataGuardAssociationProtectionModePerformance,
	"MAXIMUM_PROTECTION":   DataGuardAssociationProtectionModeProtection,
}

// GetDataGuardAssociationProtectionModeEnumValues Enumerates the set of values for DataGuardAssociationProtectionModeEnum
func GetDataGuardAssociationProtectionModeEnumValues() []DataGuardAssociationProtectionModeEnum {
	values := make([]DataGuardAssociationProtectionModeEnum, 0)
	for _, v := range mappingDataGuardAssociationProtectionMode {
		values = append(values, v)
	}
	return values
}

// DataGuardAssociationTransportTypeEnum Enum with underlying type: string
type DataGuardAssociationTransportTypeEnum string

// Set of constants representing the allowable values for DataGuardAssociationTransportTypeEnum
const (
	DataGuardAssociationTransportTypeSync     DataGuardAssociationTransportTypeEnum = "SYNC"
	DataGuardAssociationTransportTypeAsync    DataGuardAssociationTransportTypeEnum = "ASYNC"
	DataGuardAssociationTransportTypeFastsync DataGuardAssociationTransportTypeEnum = "FASTSYNC"
)

var mappingDataGuardAssociationTransportType = map[string]DataGuardAssociationTransportTypeEnum{
	"SYNC":     DataGuardAssociationTransportTypeSync,
	"ASYNC":    DataGuardAssociationTransportTypeAsync,
	"FASTSYNC": DataGuardAssociationTransportTypeFastsync,
}

// GetDataGuardAssociationTransportTypeEnumValues Enumerates the set of values for DataGuardAssociationTransportTypeEnum
func GetDataGuardAssociationTransportTypeEnumValues() []DataGuardAssociationTransportTypeEnum {
	values := make([]DataGuardAssociationTransportTypeEnum, 0)
	for _, v := range mappingDataGuardAssociationTransportType {
		values = append(values, v)
	}
	return values
}
