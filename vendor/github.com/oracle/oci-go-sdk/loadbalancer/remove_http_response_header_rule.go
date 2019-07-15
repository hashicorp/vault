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

// RemoveHttpResponseHeaderRule An object that represents the action of removing a header from a response. This rule applies only to HTTP listeners.
// If the same header appears more than once in the response, the load balancer removes all occurances of the specified header.
// **Note:** The system does not distinquish between underscore and dash characters in headers. That is, it treats
// `example_header_name` and `example-header-name` as identical. Oracle recommends that you do not rely on underscore
// or dash characters to uniquely distinguish header names.
type RemoveHttpResponseHeaderRule struct {

	// A header name that conforms to RFC 7230.
	// Example: `example_header_name`
	Header *string `mandatory:"true" json:"header"`
}

func (m RemoveHttpResponseHeaderRule) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m RemoveHttpResponseHeaderRule) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeRemoveHttpResponseHeaderRule RemoveHttpResponseHeaderRule
	s := struct {
		DiscriminatorParam string `json:"action"`
		MarshalTypeRemoveHttpResponseHeaderRule
	}{
		"REMOVE_HTTP_RESPONSE_HEADER",
		(MarshalTypeRemoveHttpResponseHeaderRule)(m),
	}

	return json.Marshal(&s)
}
