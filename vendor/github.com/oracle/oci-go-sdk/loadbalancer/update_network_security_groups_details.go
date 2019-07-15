// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Load Balancing API
//
// API for the Load Balancing service. Use this API to manage load balancers, backend sets, and related items. For more
// information, see Overview of Load Balancing (https://docs.cloud.oracle.com/iaas/Content/Balance/Concepts/balanceoverview.htm).
//

package loadbalancer

import (
	"github.com/oracle/oci-go-sdk/common"
)

// UpdateNetworkSecurityGroupsDetails An object representing an updated list of NSGs that overwrites the existing list of NSGs. In particular, if the load balancer had no prior NSGs configured, these with be the new NSGs to be used by the load balancer. If the load balancer used to have a list of NSGs configured, and this list contains no entries, then the load balancer will contain no NSGs after this call completes.
type UpdateNetworkSecurityGroupsDetails struct {

	// The array of NSG OCIDs (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) to be used by this Load Balancer.
	NetworkSecurityGroupIds []string `mandatory:"false" json:"networkSecurityGroupIds"`
}

func (m UpdateNetworkSecurityGroupsDetails) String() string {
	return common.PointerString(m)
}
