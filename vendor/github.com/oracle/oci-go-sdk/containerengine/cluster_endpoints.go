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

// ClusterEndpoints The properties that define endpoints for a cluster.
type ClusterEndpoints struct {

	// The Kubernetes API server endpoint.
	Kubernetes *string `mandatory:"false" json:"kubernetes"`
}

func (m ClusterEndpoints) String() string {
	return common.PointerString(m)
}
