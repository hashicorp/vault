// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Web Application Acceleration and Security Services API
//
// OCI Web Application Acceleration and Security Services
//

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ThreatFeedAction The action to take for a request that has been determined to be potentially malicious.
type ThreatFeedAction struct {

	// The unique key of the object for which the action applies.
	Key *string `mandatory:"true" json:"key"`

	// The selected action. If unspecified, defaults to `OFF`.
	Action ThreatFeedActionActionEnum `mandatory:"true" json:"action"`
}

func (m ThreatFeedAction) String() string {
	return common.PointerString(m)
}

// ThreatFeedActionActionEnum Enum with underlying type: string
type ThreatFeedActionActionEnum string

// Set of constants representing the allowable values for ThreatFeedActionActionEnum
const (
	ThreatFeedActionActionOff    ThreatFeedActionActionEnum = "OFF"
	ThreatFeedActionActionDetect ThreatFeedActionActionEnum = "DETECT"
	ThreatFeedActionActionBlock  ThreatFeedActionActionEnum = "BLOCK"
)

var mappingThreatFeedActionAction = map[string]ThreatFeedActionActionEnum{
	"OFF":    ThreatFeedActionActionOff,
	"DETECT": ThreatFeedActionActionDetect,
	"BLOCK":  ThreatFeedActionActionBlock,
}

// GetThreatFeedActionActionEnumValues Enumerates the set of values for ThreatFeedActionActionEnum
func GetThreatFeedActionActionEnumValues() []ThreatFeedActionActionEnum {
	values := make([]ThreatFeedActionActionEnum, 0)
	for _, v := range mappingThreatFeedActionAction {
		values = append(values, v)
	}
	return values
}
