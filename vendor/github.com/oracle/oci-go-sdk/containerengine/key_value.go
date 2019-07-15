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

// KeyValue The properties that define a key value pair.
type KeyValue struct {

	// The key of the pair.
	Key *string `mandatory:"false" json:"key"`

	// The value of the pair.
	Value *string `mandatory:"false" json:"value"`
}

func (m KeyValue) String() string {
	return common.PointerString(m)
}
