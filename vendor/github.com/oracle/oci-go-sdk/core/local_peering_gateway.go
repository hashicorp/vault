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

// LocalPeeringGateway A local peering gateway (LPG) is an object on a VCN that lets that VCN peer
// with another VCN in the same region. *Peering* means that the two VCNs can
// communicate using private IP addresses, but without the traffic traversing the
// internet or routing through your on-premises network. For more information,
// see VCN Peering (https://docs.cloud.oracle.com/Content/Network/Tasks/VCNpeering.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type LocalPeeringGateway struct {

	// The OCID of the compartment containing the LPG.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid
	// entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The LPG's Oracle ID (OCID).
	Id *string `mandatory:"true" json:"id"`

	// Whether the VCN at the other end of the peering is in a different tenancy.
	// Example: `false`
	IsCrossTenancyPeering *bool `mandatory:"true" json:"isCrossTenancyPeering"`

	// The LPG's current lifecycle state.
	LifecycleState LocalPeeringGatewayLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Whether the LPG is peered with another LPG. `NEW` means the LPG has not yet been
	// peered. `PENDING` means the peering is being established. `REVOKED` means the
	// LPG at the other end of the peering has been deleted.
	PeeringStatus LocalPeeringGatewayPeeringStatusEnum `mandatory:"true" json:"peeringStatus"`

	// The date and time the LPG was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the VCN the LPG belongs to.
	VcnId *string `mandatory:"true" json:"vcnId"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The smallest aggregate CIDR that contains all the CIDR routes advertised by the VCN
	// at the other end of the peering from this LPG. See `peerAdvertisedCidrDetails` for
	// the individual CIDRs. The value is `null` if the LPG is not peered.
	// Example: `192.168.0.0/16`, or if aggregated with `172.16.0.0/24` then `128.0.0.0/1`
	PeerAdvertisedCidr *string `mandatory:"false" json:"peerAdvertisedCidr"`

	// The specific ranges of IP addresses available on or via the VCN at the other
	// end of the peering from this LPG. The value is `null` if the LPG is not peered.
	// You can use these as destination CIDRs for route rules to route a subnet's
	// traffic to this LPG.
	// Example: [`192.168.0.0/16`, `172.16.0.0/24`]
	PeerAdvertisedCidrDetails []string `mandatory:"false" json:"peerAdvertisedCidrDetails"`

	// Additional information regarding the peering status, if applicable.
	PeeringStatusDetails *string `mandatory:"false" json:"peeringStatusDetails"`

	// The OCID of the route table the LPG is using. For information about why you
	// would associate a route table with an LPG, see
	// Advanced Scenario: Transit Routing (https://docs.cloud.oracle.com/Content/Network/Tasks/transitrouting.htm).
	RouteTableId *string `mandatory:"false" json:"routeTableId"`
}

func (m LocalPeeringGateway) String() string {
	return common.PointerString(m)
}

// LocalPeeringGatewayLifecycleStateEnum Enum with underlying type: string
type LocalPeeringGatewayLifecycleStateEnum string

// Set of constants representing the allowable values for LocalPeeringGatewayLifecycleStateEnum
const (
	LocalPeeringGatewayLifecycleStateProvisioning LocalPeeringGatewayLifecycleStateEnum = "PROVISIONING"
	LocalPeeringGatewayLifecycleStateAvailable    LocalPeeringGatewayLifecycleStateEnum = "AVAILABLE"
	LocalPeeringGatewayLifecycleStateTerminating  LocalPeeringGatewayLifecycleStateEnum = "TERMINATING"
	LocalPeeringGatewayLifecycleStateTerminated   LocalPeeringGatewayLifecycleStateEnum = "TERMINATED"
)

var mappingLocalPeeringGatewayLifecycleState = map[string]LocalPeeringGatewayLifecycleStateEnum{
	"PROVISIONING": LocalPeeringGatewayLifecycleStateProvisioning,
	"AVAILABLE":    LocalPeeringGatewayLifecycleStateAvailable,
	"TERMINATING":  LocalPeeringGatewayLifecycleStateTerminating,
	"TERMINATED":   LocalPeeringGatewayLifecycleStateTerminated,
}

// GetLocalPeeringGatewayLifecycleStateEnumValues Enumerates the set of values for LocalPeeringGatewayLifecycleStateEnum
func GetLocalPeeringGatewayLifecycleStateEnumValues() []LocalPeeringGatewayLifecycleStateEnum {
	values := make([]LocalPeeringGatewayLifecycleStateEnum, 0)
	for _, v := range mappingLocalPeeringGatewayLifecycleState {
		values = append(values, v)
	}
	return values
}

// LocalPeeringGatewayPeeringStatusEnum Enum with underlying type: string
type LocalPeeringGatewayPeeringStatusEnum string

// Set of constants representing the allowable values for LocalPeeringGatewayPeeringStatusEnum
const (
	LocalPeeringGatewayPeeringStatusInvalid LocalPeeringGatewayPeeringStatusEnum = "INVALID"
	LocalPeeringGatewayPeeringStatusNew     LocalPeeringGatewayPeeringStatusEnum = "NEW"
	LocalPeeringGatewayPeeringStatusPeered  LocalPeeringGatewayPeeringStatusEnum = "PEERED"
	LocalPeeringGatewayPeeringStatusPending LocalPeeringGatewayPeeringStatusEnum = "PENDING"
	LocalPeeringGatewayPeeringStatusRevoked LocalPeeringGatewayPeeringStatusEnum = "REVOKED"
)

var mappingLocalPeeringGatewayPeeringStatus = map[string]LocalPeeringGatewayPeeringStatusEnum{
	"INVALID": LocalPeeringGatewayPeeringStatusInvalid,
	"NEW":     LocalPeeringGatewayPeeringStatusNew,
	"PEERED":  LocalPeeringGatewayPeeringStatusPeered,
	"PENDING": LocalPeeringGatewayPeeringStatusPending,
	"REVOKED": LocalPeeringGatewayPeeringStatusRevoked,
}

// GetLocalPeeringGatewayPeeringStatusEnumValues Enumerates the set of values for LocalPeeringGatewayPeeringStatusEnum
func GetLocalPeeringGatewayPeeringStatusEnumValues() []LocalPeeringGatewayPeeringStatusEnum {
	values := make([]LocalPeeringGatewayPeeringStatusEnum, 0)
	for _, v := range mappingLocalPeeringGatewayPeeringStatus {
		values = append(values, v)
	}
	return values
}
