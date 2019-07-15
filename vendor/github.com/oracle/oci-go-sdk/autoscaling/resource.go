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

// Resource A resource that is managed by an autoscaling configuration. The only supported type is "instancePool."
// Each instance pool can have one autoscaling configuration.
type Resource interface {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the resource that is managed by the autoscaling configuration.
	GetId() *string
}

type resource struct {
	JsonData []byte
	Id       *string `mandatory:"true" json:"id"`
	Type     string  `json:"type"`
}

// UnmarshalJSON unmarshals json
func (m *resource) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalerresource resource
	s := struct {
		Model Unmarshalerresource
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.Id = s.Model.Id
	m.Type = s.Model.Type

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *resource) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.Type {
	case "instancePool":
		mm := InstancePoolResource{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

//GetId returns Id
func (m resource) GetId() *string {
	return m.Id
}

func (m resource) String() string {
	return common.PointerString(m)
}
