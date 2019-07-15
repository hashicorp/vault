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

// AddHttpRequestHeaderRule An object that represents the action of adding a header to a request.
// This rule applies only to HTTP listeners.
// **NOTES:**
// *  If a matching header already exists in the request, the system removes all of its occurrences, and then adds the
//    new header.
// *  The system does not distinquish between underscore and dash characters in headers. That is, it treats
//   `example_header_name` and `example-header-name` as identical. Oracle recommends that you do not rely on underscore
//   or dash characters to uniquely distinguish header names.
type AddHttpRequestHeaderRule struct {

	// A header name that conforms to RFC 7230.
	// Example: `example_header_name`
	Header *string `mandatory:"true" json:"header"`

	// A header value that conforms to RFC 7230.
	// Example: `example_value`
	Value *string `mandatory:"true" json:"value"`
}

func (m AddHttpRequestHeaderRule) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m AddHttpRequestHeaderRule) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeAddHttpRequestHeaderRule AddHttpRequestHeaderRule
	s := struct {
		DiscriminatorParam string `json:"action"`
		MarshalTypeAddHttpRequestHeaderRule
	}{
		"ADD_HTTP_REQUEST_HEADER",
		(MarshalTypeAddHttpRequestHeaderRule)(m),
	}

	return json.Marshal(&s)
}
