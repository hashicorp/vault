// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Container Engine for Kubernetes API
//
// API for the Container Engine for Kubernetes service. Use this API to build, deploy,
// and manage cloud-native applications. For more information, see
// Overview of Container Engine for Kubernetes (https://docs.cloud.oracle.com/iaas/Content/ContEng/Concepts/contengoverview.htm).
//

package containerengine

import (
	"github.com/oracle/oci-go-sdk/common"
)

// NodePoolSummary The properties that define a node pool summary.
type NodePoolSummary struct {

	// The OCID of the node pool.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the compartment in which the node pool exists.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The OCID of the cluster to which this node pool is attached.
	ClusterId *string `mandatory:"false" json:"clusterId"`

	// The name of the node pool.
	Name *string `mandatory:"false" json:"name"`

	// The version of Kubernetes running on the nodes in the node pool.
	KubernetesVersion *string `mandatory:"false" json:"kubernetesVersion"`

	// The OCID of the image running on the nodes in the node pool.
	NodeImageId *string `mandatory:"false" json:"nodeImageId"`

	// The name of the image running on the nodes in the node pool.
	NodeImageName *string `mandatory:"false" json:"nodeImageName"`

	// The name of the node shape of the nodes in the node pool.
	NodeShape *string `mandatory:"false" json:"nodeShape"`

	// A list of key/value pairs to add to nodes after they join the Kubernetes cluster.
	InitialNodeLabels []KeyValue `mandatory:"false" json:"initialNodeLabels"`

	// The SSH public key on each node in the node pool.
	SshPublicKey *string `mandatory:"false" json:"sshPublicKey"`

	// The number of nodes in each subnet.
	QuantityPerSubnet *int `mandatory:"false" json:"quantityPerSubnet"`

	// The OCIDs of the subnets in which to place nodes for this node pool.
	SubnetIds []string `mandatory:"false" json:"subnetIds"`
}

func (m NodePoolSummary) String() string {
	return common.PointerString(m)
}
