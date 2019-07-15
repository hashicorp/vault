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

// InstanceConfigurationCreateVnicDetails Contains the properties of the VNIC for an instance configuration. See CreateVnicDetails
// and Instance Configurations (https://docs.cloud.oracle.com/Content/Compute/Concepts/instancemanagement.htm#config) for more information.
type InstanceConfigurationCreateVnicDetails struct {

	// Whether the VNIC should be assigned a public IP address. See the `assignPublicIp` attribute of CreateVnicDetails
	// for more information.
	AssignPublicIp *bool `mandatory:"false" json:"assignPublicIp"`

	// A user-friendly name for the VNIC. Does not have to be unique.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The hostname for the VNIC's primary private IP.
	// See the `hostnameLabel` attribute of CreateVnicDetails for more information.
	HostnameLabel *string `mandatory:"false" json:"hostnameLabel"`

	// A list of the OCIDs of the network security groups (NSGs) to add the VNIC to. For more
	// information about NSGs, see
	// NetworkSecurityGroup.
	NsgIds []string `mandatory:"false" json:"nsgIds"`

	// A private IP address of your choice to assign to the VNIC.
	// See the `privateIp` attribute of CreateVnicDetails for more information.
	PrivateIp *string `mandatory:"false" json:"privateIp"`

	// Whether the source/destination check is disabled on the VNIC.
	// See the `skipSourceDestCheck` attribute of CreateVnicDetails for more information.
	SkipSourceDestCheck *bool `mandatory:"false" json:"skipSourceDestCheck"`

	// The OCID of the subnet to create the VNIC in.
	// See the `subnetId` attribute of CreateVnicDetails for more information.
	SubnetId *string `mandatory:"false" json:"subnetId"`
}

func (m InstanceConfigurationCreateVnicDetails) String() string {
	return common.PointerString(m)
}
