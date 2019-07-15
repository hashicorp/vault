// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Load Balancing API
//
// API for the Load Balancing service. Use this API to manage load balancers, backend sets, and related items. For more
// information, see Overview of Load Balancing (https://docs.cloud.oracle.com/iaas/Content/Balance/Concepts/balanceoverview.htm).
//

package loadbalancer

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// UpdateRuleSetDetails An updated set of rules that overwrites the existing set of rules.
type UpdateRuleSetDetails struct {

	// An array of rules that compose the rule set.
	Items []Rule `mandatory:"true" json:"items"`
}

func (m UpdateRuleSetDetails) String() string {
	return common.PointerString(m)
}

// UnmarshalJSON unmarshals from json
func (m *UpdateRuleSetDetails) UnmarshalJSON(data []byte) (e error) {
	model := struct {
		Items []rule `json:"items"`
	}{}

	e = json.Unmarshal(data, &model)
	if e != nil {
		return
	}
	m.Items = make([]Rule, len(model.Items))
	for i, n := range model.Items {
		nn, err := n.UnmarshalPolymorphicJSON(n.JsonData)
		if err != nil {
			return err
		}
		if nn != nil {
			m.Items[i] = nn.(Rule)
		} else {
			m.Items[i] = nil
		}
	}
	return
}
