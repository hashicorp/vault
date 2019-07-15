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

// WafBlockedRequest The representation of WafBlockedRequest
type WafBlockedRequest struct {

	// The date and time the blocked requests were observed, expressed in RFC 3339 timestamp format.
	TimeObserved *common.SDKTime `mandatory:"false" json:"timeObserved"`

	// The number of seconds the data covers.
	TimeRangeInSeconds *int `mandatory:"false" json:"timeRangeInSeconds"`

	// The specific Web Application Firewall feature that blocked the requests, such as JavaScript Challenge or Access Control.
	WafFeature WafBlockedRequestWafFeatureEnum `mandatory:"false" json:"wafFeature,omitempty"`

	// The count of blocked requests.
	Count *int `mandatory:"false" json:"count"`
}

func (m WafBlockedRequest) String() string {
	return common.PointerString(m)
}

// WafBlockedRequestWafFeatureEnum Enum with underlying type: string
type WafBlockedRequestWafFeatureEnum string

// Set of constants representing the allowable values for WafBlockedRequestWafFeatureEnum
const (
	WafBlockedRequestWafFeatureProtectionRules            WafBlockedRequestWafFeatureEnum = "PROTECTION_RULES"
	WafBlockedRequestWafFeatureJsChallenge                WafBlockedRequestWafFeatureEnum = "JS_CHALLENGE"
	WafBlockedRequestWafFeatureAccessRules                WafBlockedRequestWafFeatureEnum = "ACCESS_RULES"
	WafBlockedRequestWafFeatureThreatFeeds                WafBlockedRequestWafFeatureEnum = "THREAT_FEEDS"
	WafBlockedRequestWafFeatureHumanInteractionChallenge  WafBlockedRequestWafFeatureEnum = "HUMAN_INTERACTION_CHALLENGE"
	WafBlockedRequestWafFeatureDeviceFingerprintChallenge WafBlockedRequestWafFeatureEnum = "DEVICE_FINGERPRINT_CHALLENGE"
	WafBlockedRequestWafFeatureCaptcha                    WafBlockedRequestWafFeatureEnum = "CAPTCHA"
	WafBlockedRequestWafFeatureAddressRateLimiting        WafBlockedRequestWafFeatureEnum = "ADDRESS_RATE_LIMITING"
)

var mappingWafBlockedRequestWafFeature = map[string]WafBlockedRequestWafFeatureEnum{
	"PROTECTION_RULES":             WafBlockedRequestWafFeatureProtectionRules,
	"JS_CHALLENGE":                 WafBlockedRequestWafFeatureJsChallenge,
	"ACCESS_RULES":                 WafBlockedRequestWafFeatureAccessRules,
	"THREAT_FEEDS":                 WafBlockedRequestWafFeatureThreatFeeds,
	"HUMAN_INTERACTION_CHALLENGE":  WafBlockedRequestWafFeatureHumanInteractionChallenge,
	"DEVICE_FINGERPRINT_CHALLENGE": WafBlockedRequestWafFeatureDeviceFingerprintChallenge,
	"CAPTCHA":                      WafBlockedRequestWafFeatureCaptcha,
	"ADDRESS_RATE_LIMITING":        WafBlockedRequestWafFeatureAddressRateLimiting,
}

// GetWafBlockedRequestWafFeatureEnumValues Enumerates the set of values for WafBlockedRequestWafFeatureEnum
func GetWafBlockedRequestWafFeatureEnumValues() []WafBlockedRequestWafFeatureEnum {
	values := make([]WafBlockedRequestWafFeatureEnum, 0)
	for _, v := range mappingWafBlockedRequestWafFeature {
		values = append(values, v)
	}
	return values
}
