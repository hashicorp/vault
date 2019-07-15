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

// ThreatFeed The settings of the threat intelligence feed. You can block requests from IP addresses based on their reputations with various commercial and open source threat feeds.
type ThreatFeed struct {

	// The unique key of the threat intelligence feed.
	Key *string `mandatory:"false" json:"key"`

	// The name of the threat intelligence feed.
	Name *string `mandatory:"false" json:"name"`

	// The action to take when traffic is flagged as malicious by data from the threat intelligence feed. If unspecified, defaults to `OFF`.
	Action ThreatFeedActionEnum `mandatory:"false" json:"action,omitempty"`

	// The description of the threat intelligence feed.
	Description *string `mandatory:"false" json:"description"`
}

func (m ThreatFeed) String() string {
	return common.PointerString(m)
}

// ThreatFeedActionEnum Enum with underlying type: string
type ThreatFeedActionEnum string

// Set of constants representing the allowable values for ThreatFeedActionEnum
const (
	ThreatFeedActionOff    ThreatFeedActionEnum = "OFF"
	ThreatFeedActionDetect ThreatFeedActionEnum = "DETECT"
	ThreatFeedActionBlock  ThreatFeedActionEnum = "BLOCK"
)

var mappingThreatFeedAction = map[string]ThreatFeedActionEnum{
	"OFF":    ThreatFeedActionOff,
	"DETECT": ThreatFeedActionDetect,
	"BLOCK":  ThreatFeedActionBlock,
}

// GetThreatFeedActionEnumValues Enumerates the set of values for ThreatFeedActionEnum
func GetThreatFeedActionEnumValues() []ThreatFeedActionEnum {
	values := make([]ThreatFeedActionEnum, 0)
	for _, v := range mappingThreatFeedAction {
		values = append(values, v)
	}
	return values
}
