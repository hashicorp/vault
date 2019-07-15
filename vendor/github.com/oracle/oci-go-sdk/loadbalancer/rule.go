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

// Rule An object that represents an action to apply to a listener.
type Rule interface {
}

type rule struct {
	JsonData []byte
	Action   string `json:"action"`
}

// UnmarshalJSON unmarshals json
func (m *rule) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalerrule rule
	s := struct {
		Model Unmarshalerrule
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.Action = s.Model.Action

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *rule) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.Action {
	case "ADD_HTTP_REQUEST_HEADER":
		mm := AddHttpRequestHeaderRule{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "REMOVE_HTTP_REQUEST_HEADER":
		mm := RemoveHttpRequestHeaderRule{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "EXTEND_HTTP_REQUEST_HEADER_VALUE":
		mm := ExtendHttpRequestHeaderValueRule{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "REMOVE_HTTP_RESPONSE_HEADER":
		mm := RemoveHttpResponseHeaderRule{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "ADD_HTTP_RESPONSE_HEADER":
		mm := AddHttpResponseHeaderRule{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "EXTEND_HTTP_RESPONSE_HEADER_VALUE":
		mm := ExtendHttpResponseHeaderValueRule{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

func (m rule) String() string {
	return common.PointerString(m)
}
