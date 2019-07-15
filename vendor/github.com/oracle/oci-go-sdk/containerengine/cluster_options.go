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

// ClusterOptions Options for creating or updating clusters.
type ClusterOptions struct {

	// Available Kubernetes versions.
	KubernetesVersions []string `mandatory:"false" json:"kubernetesVersions"`
}

func (m ClusterOptions) String() string {
	return common.PointerString(m)
}
