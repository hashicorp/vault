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

// Cluster A Kubernetes cluster. Avoid entering confidential information.
type Cluster struct {

	// The OCID of the cluster.
	Id *string `mandatory:"false" json:"id"`

	// The name of the cluster.
	Name *string `mandatory:"false" json:"name"`

	// The OCID of the compartment in which the cluster exists.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The OCID of the virtual cloud network (VCN) in which the cluster exists.
	VcnId *string `mandatory:"false" json:"vcnId"`

	// The version of Kubernetes running on the cluster masters.
	KubernetesVersion *string `mandatory:"false" json:"kubernetesVersion"`

	// Optional attributes for the cluster.
	Options *ClusterCreateOptions `mandatory:"false" json:"options"`

	// Metadata about the cluster.
	Metadata *ClusterMetadata `mandatory:"false" json:"metadata"`

	// The state of the cluster masters.
	LifecycleState ClusterLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// Details about the state of the cluster masters.
	LifecycleDetails *string `mandatory:"false" json:"lifecycleDetails"`

	// Endpoints served up by the cluster masters.
	Endpoints *ClusterEndpoints `mandatory:"false" json:"endpoints"`

	// Available Kubernetes versions to which the clusters masters may be upgraded.
	AvailableKubernetesUpgrades []string `mandatory:"false" json:"availableKubernetesUpgrades"`
}

func (m Cluster) String() string {
	return common.PointerString(m)
}

// ClusterLifecycleStateEnum Enum with underlying type: string
type ClusterLifecycleStateEnum string

// Set of constants representing the allowable values for ClusterLifecycleStateEnum
const (
	ClusterLifecycleStateCreating ClusterLifecycleStateEnum = "CREATING"
	ClusterLifecycleStateActive   ClusterLifecycleStateEnum = "ACTIVE"
	ClusterLifecycleStateFailed   ClusterLifecycleStateEnum = "FAILED"
	ClusterLifecycleStateDeleting ClusterLifecycleStateEnum = "DELETING"
	ClusterLifecycleStateDeleted  ClusterLifecycleStateEnum = "DELETED"
	ClusterLifecycleStateUpdating ClusterLifecycleStateEnum = "UPDATING"
)

var mappingClusterLifecycleState = map[string]ClusterLifecycleStateEnum{
	"CREATING": ClusterLifecycleStateCreating,
	"ACTIVE":   ClusterLifecycleStateActive,
	"FAILED":   ClusterLifecycleStateFailed,
	"DELETING": ClusterLifecycleStateDeleting,
	"DELETED":  ClusterLifecycleStateDeleted,
	"UPDATING": ClusterLifecycleStateUpdating,
}

// GetClusterLifecycleStateEnumValues Enumerates the set of values for ClusterLifecycleStateEnum
func GetClusterLifecycleStateEnumValues() []ClusterLifecycleStateEnum {
	values := make([]ClusterLifecycleStateEnum, 0)
	for _, v := range mappingClusterLifecycleState {
		values = append(values, v)
	}
	return values
}
