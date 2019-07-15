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

// VirtualCircuitPublicPrefix A public IP prefix and its details. With a public virtual circuit, the customer
// specifies the customer-owned public IP prefixes to advertise across the connection.
// For more information, see FastConnect Overview (https://docs.cloud.oracle.com/Content/Network/Concepts/fastconnect.htm).
type VirtualCircuitPublicPrefix struct {

	// Publix IP prefix (CIDR) that the customer specified.
	CidrBlock *string `mandatory:"true" json:"cidrBlock"`

	// Oracle must verify that the customer owns the public IP prefix before traffic
	// for that prefix can flow across the virtual circuit. Verification can take a
	// few business days. `IN_PROGRESS` means Oracle is verifying the prefix. `COMPLETED`
	// means verification succeeded. `FAILED` means verification failed and traffic for
	// this prefix will not flow across the connection.
	VerificationState VirtualCircuitPublicPrefixVerificationStateEnum `mandatory:"true" json:"verificationState"`
}

func (m VirtualCircuitPublicPrefix) String() string {
	return common.PointerString(m)
}

// VirtualCircuitPublicPrefixVerificationStateEnum Enum with underlying type: string
type VirtualCircuitPublicPrefixVerificationStateEnum string

// Set of constants representing the allowable values for VirtualCircuitPublicPrefixVerificationStateEnum
const (
	VirtualCircuitPublicPrefixVerificationStateInProgress VirtualCircuitPublicPrefixVerificationStateEnum = "IN_PROGRESS"
	VirtualCircuitPublicPrefixVerificationStateCompleted  VirtualCircuitPublicPrefixVerificationStateEnum = "COMPLETED"
	VirtualCircuitPublicPrefixVerificationStateFailed     VirtualCircuitPublicPrefixVerificationStateEnum = "FAILED"
)

var mappingVirtualCircuitPublicPrefixVerificationState = map[string]VirtualCircuitPublicPrefixVerificationStateEnum{
	"IN_PROGRESS": VirtualCircuitPublicPrefixVerificationStateInProgress,
	"COMPLETED":   VirtualCircuitPublicPrefixVerificationStateCompleted,
	"FAILED":      VirtualCircuitPublicPrefixVerificationStateFailed,
}

// GetVirtualCircuitPublicPrefixVerificationStateEnumValues Enumerates the set of values for VirtualCircuitPublicPrefixVerificationStateEnum
func GetVirtualCircuitPublicPrefixVerificationStateEnumValues() []VirtualCircuitPublicPrefixVerificationStateEnum {
	values := make([]VirtualCircuitPublicPrefixVerificationStateEnum, 0)
	for _, v := range mappingVirtualCircuitPublicPrefixVerificationState {
		values = append(values, v)
	}
	return values
}
