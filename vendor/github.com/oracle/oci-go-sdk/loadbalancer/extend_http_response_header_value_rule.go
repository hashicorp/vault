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

// ExtendHttpResponseHeaderValueRule An object that represents the action of modifying a response header value. This rule applies only to HTTP listeners.
// This rule adds a prefix, a suffix, or both to the header value.
// **NOTES:**
// *  This rule requires a value for a prefix, suffix, or both.
// *  The system does not support this rule for headers with multiple values.
// *  The system does not distinquish between underscore and dash characters in headers. That is, it treats
//    `example_header_name` and `example-header-name` as identical.  If two such headers appear in a request, the system
//    applies the action to the first header it finds. The affected header cannot be determined in advance. Oracle
//    recommends that you do not rely on underscore or dash characters to uniquely distinguish header names.
type ExtendHttpResponseHeaderValueRule struct {

	// A header name that conforms to RFC 7230.
	// Example: `example_header_name`
	Header *string `mandatory:"true" json:"header"`

	// A string to prepend to the header value. The resulting header value must still conform to RFC 7230.
	// Example: `example_prefix_value`
	Prefix *string `mandatory:"false" json:"prefix"`

	// A string to append to the header value. The resulting header value must still conform to RFC 7230.
	// Example: `example_suffix_value`
	Suffix *string `mandatory:"false" json:"suffix"`
}

func (m ExtendHttpResponseHeaderValueRule) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m ExtendHttpResponseHeaderValueRule) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeExtendHttpResponseHeaderValueRule ExtendHttpResponseHeaderValueRule
	s := struct {
		DiscriminatorParam string `json:"action"`
		MarshalTypeExtendHttpResponseHeaderValueRule
	}{
		"EXTEND_HTTP_RESPONSE_HEADER_VALUE",
		(MarshalTypeExtendHttpResponseHeaderValueRule)(m),
	}

	return json.Marshal(&s)
}
