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

// SecurityList A set of virtual firewall rules for your VCN. Security lists are configured at the subnet
// level, but the rules are applied to the ingress and egress traffic for the individual instances
// in the subnet. The rules can be stateful or stateless. For more information, see
// Security Lists (https://docs.cloud.oracle.com/Content/Network/Concepts/securitylists.htm).
// **Note:** Compare security lists to NetworkSecurityGroups,
// which let you apply a set of security rules to a *specific set of VNICs* instead of an entire
// subnet. Oracle recommends using network security groups instead of security lists, although you
// can use either or both together.
// **Important:** Oracle Cloud Infrastructure Compute service images automatically include firewall rules (for example,
// Linux iptables, Windows firewall). If there are issues with some type of access to an instance,
// make sure both the security lists associated with the instance's subnet and the instance's
// firewall rules are set correctly.
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type SecurityList struct {

	// The OCID of the compartment containing the security list.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// Rules for allowing egress IP packets.
	EgressSecurityRules []EgressSecurityRule `mandatory:"true" json:"egressSecurityRules"`

	// The security list's Oracle Cloud ID (OCID).
	Id *string `mandatory:"true" json:"id"`

	// Rules for allowing ingress IP packets.
	IngressSecurityRules []IngressSecurityRule `mandatory:"true" json:"ingressSecurityRules"`

	// The security list's current state.
	LifecycleState SecurityListLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the security list was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the VCN the security list belongs to.
	VcnId *string `mandatory:"true" json:"vcnId"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m SecurityList) String() string {
	return common.PointerString(m)
}

// SecurityListLifecycleStateEnum Enum with underlying type: string
type SecurityListLifecycleStateEnum string

// Set of constants representing the allowable values for SecurityListLifecycleStateEnum
const (
	SecurityListLifecycleStateProvisioning SecurityListLifecycleStateEnum = "PROVISIONING"
	SecurityListLifecycleStateAvailable    SecurityListLifecycleStateEnum = "AVAILABLE"
	SecurityListLifecycleStateTerminating  SecurityListLifecycleStateEnum = "TERMINATING"
	SecurityListLifecycleStateTerminated   SecurityListLifecycleStateEnum = "TERMINATED"
)

var mappingSecurityListLifecycleState = map[string]SecurityListLifecycleStateEnum{
	"PROVISIONING": SecurityListLifecycleStateProvisioning,
	"AVAILABLE":    SecurityListLifecycleStateAvailable,
	"TERMINATING":  SecurityListLifecycleStateTerminating,
	"TERMINATED":   SecurityListLifecycleStateTerminated,
}

// GetSecurityListLifecycleStateEnumValues Enumerates the set of values for SecurityListLifecycleStateEnum
func GetSecurityListLifecycleStateEnumValues() []SecurityListLifecycleStateEnum {
	values := make([]SecurityListLifecycleStateEnum, 0)
	for _, v := range mappingSecurityListLifecycleState {
		values = append(values, v)
	}
	return values
}
