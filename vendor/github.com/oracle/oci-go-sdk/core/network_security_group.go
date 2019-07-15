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

// NetworkSecurityGroup A *network security group* (NSG) provides virtual firewall rules for a specific set of
// Vnic in a VCN. Compare NSGs with SecurityList,
// which provide virtual firewall rules to all the VNICs in a *subnet*.
// A network security group consists of two items:
//   * The set of Vnic that all have the same security rule needs (for
//     example, a group of Compute instances all running the same application)
//   * A set of NSG SecurityRule that apply to the VNICs in the group
// After creating an NSG, you can add VNICs and security rules to it. For example, when you create
// an instance, you can specify one or more NSGs to add the instance to (see
// CreateVnicDetails). Or you can add an existing
// instance to an NSG with UpdateVnic.
// To add security rules to an NSG, see
// AddNetworkSecurityGroupSecurityRules.
// To list the VNICs in an NSG, see
// ListNetworkSecurityGroupVnics.
// To list the security rules in an NSG, see
// ListNetworkSecurityGroupSecurityRules.
// For more information about network security groups, see
// Network Security Groups (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/networksecuritygroups.htm).
// **Important:** Oracle Cloud Infrastructure Compute service images automatically include firewall rules (for example,
// Linux iptables, Windows firewall). If there are issues with some type of access to an instance,
// make sure all of the following are set correctly:
//   * Any security rules in any NSGs the instance's VNIC belongs to
//   * Any SecurityList associated with the instance's subnet
//   * The instance's OS firewall rules
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type NetworkSecurityGroup struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment the network security group is in.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the network security group.
	Id *string `mandatory:"true" json:"id"`

	// The network security group's current state.
	LifecycleState NetworkSecurityGroupLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the network security group was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the network security group's VCN.
	VcnId *string `mandatory:"true" json:"vcnId"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// A user-friendly name. Does not have to be unique.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m NetworkSecurityGroup) String() string {
	return common.PointerString(m)
}

// NetworkSecurityGroupLifecycleStateEnum Enum with underlying type: string
type NetworkSecurityGroupLifecycleStateEnum string

// Set of constants representing the allowable values for NetworkSecurityGroupLifecycleStateEnum
const (
	NetworkSecurityGroupLifecycleStateProvisioning NetworkSecurityGroupLifecycleStateEnum = "PROVISIONING"
	NetworkSecurityGroupLifecycleStateAvailable    NetworkSecurityGroupLifecycleStateEnum = "AVAILABLE"
	NetworkSecurityGroupLifecycleStateTerminating  NetworkSecurityGroupLifecycleStateEnum = "TERMINATING"
	NetworkSecurityGroupLifecycleStateTerminated   NetworkSecurityGroupLifecycleStateEnum = "TERMINATED"
)

var mappingNetworkSecurityGroupLifecycleState = map[string]NetworkSecurityGroupLifecycleStateEnum{
	"PROVISIONING": NetworkSecurityGroupLifecycleStateProvisioning,
	"AVAILABLE":    NetworkSecurityGroupLifecycleStateAvailable,
	"TERMINATING":  NetworkSecurityGroupLifecycleStateTerminating,
	"TERMINATED":   NetworkSecurityGroupLifecycleStateTerminated,
}

// GetNetworkSecurityGroupLifecycleStateEnumValues Enumerates the set of values for NetworkSecurityGroupLifecycleStateEnum
func GetNetworkSecurityGroupLifecycleStateEnumValues() []NetworkSecurityGroupLifecycleStateEnum {
	values := make([]NetworkSecurityGroupLifecycleStateEnum, 0)
	for _, v := range mappingNetworkSecurityGroupLifecycleState {
		values = append(values, v)
	}
	return values
}
