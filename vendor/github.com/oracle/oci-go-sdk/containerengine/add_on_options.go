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

// AddOnOptions The properties that define options for supported add-ons.
type AddOnOptions struct {

	// Whether or not to enable the Kubernetes Dashboard add-on.
	IsKubernetesDashboardEnabled *bool `mandatory:"false" json:"isKubernetesDashboardEnabled"`

	// Whether or not to enable the Tiller add-on.
	IsTillerEnabled *bool `mandatory:"false" json:"isTillerEnabled"`
}

func (m AddOnOptions) String() string {
	return common.PointerString(m)
}
