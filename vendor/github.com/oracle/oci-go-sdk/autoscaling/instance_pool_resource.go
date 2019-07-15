// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Autoscaling API
//
// APIs for dynamically scaling Compute resources to meet application requirements.
// For information about the Compute service, see Overview of the Compute Service (https://docs.cloud.oracle.com/Content/Compute/Concepts/computeoverview.htm).
//

package autoscaling

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// InstancePoolResource A Compute instance pool.
type InstancePoolResource struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the resource that is managed by the autoscaling configuration.
	Id *string `mandatory:"true" json:"id"`
}

//GetId returns Id
func (m InstancePoolResource) GetId() *string {
	return m.Id
}

func (m InstancePoolResource) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m InstancePoolResource) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeInstancePoolResource InstancePoolResource
	s := struct {
		DiscriminatorParam string `json:"type"`
		MarshalTypeInstancePoolResource
	}{
		"instancePool",
		(MarshalTypeInstancePoolResource)(m),
	}

	return json.Marshal(&s)
}
