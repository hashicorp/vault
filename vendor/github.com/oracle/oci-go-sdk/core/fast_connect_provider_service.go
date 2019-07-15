// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"github.com/oracle/oci-go-sdk/common"
)

// FastConnectProviderService A service offering from a supported provider. For more information,
// see FastConnect Overview (https://docs.cloud.oracle.com/Content/Network/Concepts/fastconnect.htm).
type FastConnectProviderService struct {

	// The OCID of the service offered by the provider.
	Id *string `mandatory:"true" json:"id"`

	// Who is responsible for managing the private peering BGP information.
	PrivatePeeringBgpManagement FastConnectProviderServicePrivatePeeringBgpManagementEnum `mandatory:"true" json:"privatePeeringBgpManagement"`

	// The name of the provider.
	ProviderName *string `mandatory:"true" json:"providerName"`

	// The name of the service offered by the provider.
	ProviderServiceName *string `mandatory:"true" json:"providerServiceName"`

	// Who is responsible for managing the public peering BGP information.
	PublicPeeringBgpManagement FastConnectProviderServicePublicPeeringBgpManagementEnum `mandatory:"true" json:"publicPeeringBgpManagement"`

	// Who is responsible for managing the ASN information for the network at the other end
	// of the connection from Oracle.
	CustomerAsnManagement FastConnectProviderServiceCustomerAsnManagementEnum `mandatory:"true" json:"customerAsnManagement"`

	// Who is responsible for managing the provider service key.
	ProviderServiceKeyManagement FastConnectProviderServiceProviderServiceKeyManagementEnum `mandatory:"true" json:"providerServiceKeyManagement"`

	// Who is responsible for managing the virtual circuit bandwidth.
	BandwithShapeManagement FastConnectProviderServiceBandwithShapeManagementEnum `mandatory:"true" json:"bandwithShapeManagement"`

	// Total number of cross-connect or cross-connect groups required for the virtual circuit.
	RequiredTotalCrossConnects *int `mandatory:"true" json:"requiredTotalCrossConnects"`

	// Provider service type.
	Type FastConnectProviderServiceTypeEnum `mandatory:"true" json:"type"`

	// The location of the provider's website or portal. This portal is where you can get information
	// about the provider service, create a virtual circuit connection from the provider to Oracle
	// Cloud Infrastructure, and retrieve your provider service key for that virtual circuit connection.
	// Example: `https://example.com`
	Description *string `mandatory:"false" json:"description"`

	// An array of virtual circuit types supported by this service.
	SupportedVirtualCircuitTypes []FastConnectProviderServiceSupportedVirtualCircuitTypesEnum `mandatory:"false" json:"supportedVirtualCircuitTypes,omitempty"`
}

func (m FastConnectProviderService) String() string {
	return common.PointerString(m)
}

// FastConnectProviderServicePrivatePeeringBgpManagementEnum Enum with underlying type: string
type FastConnectProviderServicePrivatePeeringBgpManagementEnum string

// Set of constants representing the allowable values for FastConnectProviderServicePrivatePeeringBgpManagementEnum
const (
	FastConnectProviderServicePrivatePeeringBgpManagementCustomerManaged FastConnectProviderServicePrivatePeeringBgpManagementEnum = "CUSTOMER_MANAGED"
	FastConnectProviderServicePrivatePeeringBgpManagementProviderManaged FastConnectProviderServicePrivatePeeringBgpManagementEnum = "PROVIDER_MANAGED"
	FastConnectProviderServicePrivatePeeringBgpManagementOracleManaged   FastConnectProviderServicePrivatePeeringBgpManagementEnum = "ORACLE_MANAGED"
)

var mappingFastConnectProviderServicePrivatePeeringBgpManagement = map[string]FastConnectProviderServicePrivatePeeringBgpManagementEnum{
	"CUSTOMER_MANAGED": FastConnectProviderServicePrivatePeeringBgpManagementCustomerManaged,
	"PROVIDER_MANAGED": FastConnectProviderServicePrivatePeeringBgpManagementProviderManaged,
	"ORACLE_MANAGED":   FastConnectProviderServicePrivatePeeringBgpManagementOracleManaged,
}

// GetFastConnectProviderServicePrivatePeeringBgpManagementEnumValues Enumerates the set of values for FastConnectProviderServicePrivatePeeringBgpManagementEnum
func GetFastConnectProviderServicePrivatePeeringBgpManagementEnumValues() []FastConnectProviderServicePrivatePeeringBgpManagementEnum {
	values := make([]FastConnectProviderServicePrivatePeeringBgpManagementEnum, 0)
	for _, v := range mappingFastConnectProviderServicePrivatePeeringBgpManagement {
		values = append(values, v)
	}
	return values
}

// FastConnectProviderServicePublicPeeringBgpManagementEnum Enum with underlying type: string
type FastConnectProviderServicePublicPeeringBgpManagementEnum string

// Set of constants representing the allowable values for FastConnectProviderServicePublicPeeringBgpManagementEnum
const (
	FastConnectProviderServicePublicPeeringBgpManagementCustomerManaged FastConnectProviderServicePublicPeeringBgpManagementEnum = "CUSTOMER_MANAGED"
	FastConnectProviderServicePublicPeeringBgpManagementProviderManaged FastConnectProviderServicePublicPeeringBgpManagementEnum = "PROVIDER_MANAGED"
	FastConnectProviderServicePublicPeeringBgpManagementOracleManaged   FastConnectProviderServicePublicPeeringBgpManagementEnum = "ORACLE_MANAGED"
)

var mappingFastConnectProviderServicePublicPeeringBgpManagement = map[string]FastConnectProviderServicePublicPeeringBgpManagementEnum{
	"CUSTOMER_MANAGED": FastConnectProviderServicePublicPeeringBgpManagementCustomerManaged,
	"PROVIDER_MANAGED": FastConnectProviderServicePublicPeeringBgpManagementProviderManaged,
	"ORACLE_MANAGED":   FastConnectProviderServicePublicPeeringBgpManagementOracleManaged,
}

// GetFastConnectProviderServicePublicPeeringBgpManagementEnumValues Enumerates the set of values for FastConnectProviderServicePublicPeeringBgpManagementEnum
func GetFastConnectProviderServicePublicPeeringBgpManagementEnumValues() []FastConnectProviderServicePublicPeeringBgpManagementEnum {
	values := make([]FastConnectProviderServicePublicPeeringBgpManagementEnum, 0)
	for _, v := range mappingFastConnectProviderServicePublicPeeringBgpManagement {
		values = append(values, v)
	}
	return values
}

// FastConnectProviderServiceSupportedVirtualCircuitTypesEnum Enum with underlying type: string
type FastConnectProviderServiceSupportedVirtualCircuitTypesEnum string

// Set of constants representing the allowable values for FastConnectProviderServiceSupportedVirtualCircuitTypesEnum
const (
	FastConnectProviderServiceSupportedVirtualCircuitTypesPublic  FastConnectProviderServiceSupportedVirtualCircuitTypesEnum = "PUBLIC"
	FastConnectProviderServiceSupportedVirtualCircuitTypesPrivate FastConnectProviderServiceSupportedVirtualCircuitTypesEnum = "PRIVATE"
)

var mappingFastConnectProviderServiceSupportedVirtualCircuitTypes = map[string]FastConnectProviderServiceSupportedVirtualCircuitTypesEnum{
	"PUBLIC":  FastConnectProviderServiceSupportedVirtualCircuitTypesPublic,
	"PRIVATE": FastConnectProviderServiceSupportedVirtualCircuitTypesPrivate,
}

// GetFastConnectProviderServiceSupportedVirtualCircuitTypesEnumValues Enumerates the set of values for FastConnectProviderServiceSupportedVirtualCircuitTypesEnum
func GetFastConnectProviderServiceSupportedVirtualCircuitTypesEnumValues() []FastConnectProviderServiceSupportedVirtualCircuitTypesEnum {
	values := make([]FastConnectProviderServiceSupportedVirtualCircuitTypesEnum, 0)
	for _, v := range mappingFastConnectProviderServiceSupportedVirtualCircuitTypes {
		values = append(values, v)
	}
	return values
}

// FastConnectProviderServiceCustomerAsnManagementEnum Enum with underlying type: string
type FastConnectProviderServiceCustomerAsnManagementEnum string

// Set of constants representing the allowable values for FastConnectProviderServiceCustomerAsnManagementEnum
const (
	FastConnectProviderServiceCustomerAsnManagementCustomerManaged FastConnectProviderServiceCustomerAsnManagementEnum = "CUSTOMER_MANAGED"
	FastConnectProviderServiceCustomerAsnManagementProviderManaged FastConnectProviderServiceCustomerAsnManagementEnum = "PROVIDER_MANAGED"
	FastConnectProviderServiceCustomerAsnManagementOracleManaged   FastConnectProviderServiceCustomerAsnManagementEnum = "ORACLE_MANAGED"
)

var mappingFastConnectProviderServiceCustomerAsnManagement = map[string]FastConnectProviderServiceCustomerAsnManagementEnum{
	"CUSTOMER_MANAGED": FastConnectProviderServiceCustomerAsnManagementCustomerManaged,
	"PROVIDER_MANAGED": FastConnectProviderServiceCustomerAsnManagementProviderManaged,
	"ORACLE_MANAGED":   FastConnectProviderServiceCustomerAsnManagementOracleManaged,
}

// GetFastConnectProviderServiceCustomerAsnManagementEnumValues Enumerates the set of values for FastConnectProviderServiceCustomerAsnManagementEnum
func GetFastConnectProviderServiceCustomerAsnManagementEnumValues() []FastConnectProviderServiceCustomerAsnManagementEnum {
	values := make([]FastConnectProviderServiceCustomerAsnManagementEnum, 0)
	for _, v := range mappingFastConnectProviderServiceCustomerAsnManagement {
		values = append(values, v)
	}
	return values
}

// FastConnectProviderServiceProviderServiceKeyManagementEnum Enum with underlying type: string
type FastConnectProviderServiceProviderServiceKeyManagementEnum string

// Set of constants representing the allowable values for FastConnectProviderServiceProviderServiceKeyManagementEnum
const (
	FastConnectProviderServiceProviderServiceKeyManagementCustomerManaged FastConnectProviderServiceProviderServiceKeyManagementEnum = "CUSTOMER_MANAGED"
	FastConnectProviderServiceProviderServiceKeyManagementProviderManaged FastConnectProviderServiceProviderServiceKeyManagementEnum = "PROVIDER_MANAGED"
	FastConnectProviderServiceProviderServiceKeyManagementOracleManaged   FastConnectProviderServiceProviderServiceKeyManagementEnum = "ORACLE_MANAGED"
)

var mappingFastConnectProviderServiceProviderServiceKeyManagement = map[string]FastConnectProviderServiceProviderServiceKeyManagementEnum{
	"CUSTOMER_MANAGED": FastConnectProviderServiceProviderServiceKeyManagementCustomerManaged,
	"PROVIDER_MANAGED": FastConnectProviderServiceProviderServiceKeyManagementProviderManaged,
	"ORACLE_MANAGED":   FastConnectProviderServiceProviderServiceKeyManagementOracleManaged,
}

// GetFastConnectProviderServiceProviderServiceKeyManagementEnumValues Enumerates the set of values for FastConnectProviderServiceProviderServiceKeyManagementEnum
func GetFastConnectProviderServiceProviderServiceKeyManagementEnumValues() []FastConnectProviderServiceProviderServiceKeyManagementEnum {
	values := make([]FastConnectProviderServiceProviderServiceKeyManagementEnum, 0)
	for _, v := range mappingFastConnectProviderServiceProviderServiceKeyManagement {
		values = append(values, v)
	}
	return values
}

// FastConnectProviderServiceBandwithShapeManagementEnum Enum with underlying type: string
type FastConnectProviderServiceBandwithShapeManagementEnum string

// Set of constants representing the allowable values for FastConnectProviderServiceBandwithShapeManagementEnum
const (
	FastConnectProviderServiceBandwithShapeManagementCustomerManaged FastConnectProviderServiceBandwithShapeManagementEnum = "CUSTOMER_MANAGED"
	FastConnectProviderServiceBandwithShapeManagementProviderManaged FastConnectProviderServiceBandwithShapeManagementEnum = "PROVIDER_MANAGED"
	FastConnectProviderServiceBandwithShapeManagementOracleManaged   FastConnectProviderServiceBandwithShapeManagementEnum = "ORACLE_MANAGED"
)

var mappingFastConnectProviderServiceBandwithShapeManagement = map[string]FastConnectProviderServiceBandwithShapeManagementEnum{
	"CUSTOMER_MANAGED": FastConnectProviderServiceBandwithShapeManagementCustomerManaged,
	"PROVIDER_MANAGED": FastConnectProviderServiceBandwithShapeManagementProviderManaged,
	"ORACLE_MANAGED":   FastConnectProviderServiceBandwithShapeManagementOracleManaged,
}

// GetFastConnectProviderServiceBandwithShapeManagementEnumValues Enumerates the set of values for FastConnectProviderServiceBandwithShapeManagementEnum
func GetFastConnectProviderServiceBandwithShapeManagementEnumValues() []FastConnectProviderServiceBandwithShapeManagementEnum {
	values := make([]FastConnectProviderServiceBandwithShapeManagementEnum, 0)
	for _, v := range mappingFastConnectProviderServiceBandwithShapeManagement {
		values = append(values, v)
	}
	return values
}

// FastConnectProviderServiceTypeEnum Enum with underlying type: string
type FastConnectProviderServiceTypeEnum string

// Set of constants representing the allowable values for FastConnectProviderServiceTypeEnum
const (
	FastConnectProviderServiceTypeLayer2 FastConnectProviderServiceTypeEnum = "LAYER2"
	FastConnectProviderServiceTypeLayer3 FastConnectProviderServiceTypeEnum = "LAYER3"
)

var mappingFastConnectProviderServiceType = map[string]FastConnectProviderServiceTypeEnum{
	"LAYER2": FastConnectProviderServiceTypeLayer2,
	"LAYER3": FastConnectProviderServiceTypeLayer3,
}

// GetFastConnectProviderServiceTypeEnumValues Enumerates the set of values for FastConnectProviderServiceTypeEnum
func GetFastConnectProviderServiceTypeEnumValues() []FastConnectProviderServiceTypeEnum {
	values := make([]FastConnectProviderServiceTypeEnum, 0)
	for _, v := range mappingFastConnectProviderServiceType {
		values = append(values, v)
	}
	return values
}
