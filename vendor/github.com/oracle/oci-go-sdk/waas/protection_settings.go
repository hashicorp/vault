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

// ProtectionSettings The settings used for protection rules.
type ProtectionSettings struct {

	// If `action` is set to `BLOCK`, this specifies how the traffic is blocked when detected as malicious by a protection rule. If unspecified, defaults to `SET_RESPONSE_CODE`.
	BlockAction ProtectionSettingsBlockActionEnum `mandatory:"false" json:"blockAction,omitempty"`

	// The response code returned when `action` is set to `BLOCK`, `blockAction` is set to `SET_RESPONSE_CODE`, and the traffic is detected as malicious by a protection rule. If unspecified, defaults to `403`.
	BlockResponseCode *int `mandatory:"false" json:"blockResponseCode"`

	// The message to show on the error page when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_ERROR_PAGE`, and the traffic is detected as malicious by a protection rule. If unspecified, defaults to 'Access to the website is blocked.'
	BlockErrorPageMessage *string `mandatory:"false" json:"blockErrorPageMessage"`

	// The error code to show on the error page when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_ERROR_PAGE`, and the traffic is detected as malicious by a protection rule. If unspecified, defaults to `403`.
	BlockErrorPageCode *string `mandatory:"false" json:"blockErrorPageCode"`

	// The description text to show on the error page when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_ERROR_PAGE`, and the traffic is detected as malicious by a protection rule. If unspecified, defaults to `Access blocked by website owner. Please contact support.`
	BlockErrorPageDescription *string `mandatory:"false" json:"blockErrorPageDescription"`

	// The maximum number of arguments allowed to be passed to your application before an action is taken. If unspecified, defaults to `255`.
	MaxArgumentCount *int `mandatory:"false" json:"maxArgumentCount"`

	// The maximum length allowed for each argument name, in characters. If unspecified, defaults to `400`.
	MaxNameLengthPerArgument *int `mandatory:"false" json:"maxNameLengthPerArgument"`

	// The maximum length allowed for the sum of all argument names, in characters. If unspecified, defaults to `64000`.
	MaxTotalNameLengthOfArguments *int `mandatory:"false" json:"maxTotalNameLengthOfArguments"`

	// The length of time to analyze traffic traffic, in days. After the analysis period, `WafRecommendations` will be populated. If unspecified, defaults to `10`.
	// Use `GET /waasPolicies/{waasPolicyId}/wafRecommendations` to view WAF recommendations.
	RecommendationsPeriodInDays *int `mandatory:"false" json:"recommendationsPeriodInDays"`

	// Inspects the response body of origin responses. Can be used to detect leakage of sensitive data. If unspecified, defaults to `false`.
	// **Note:** Only origin responses with a Content-Type matching a value in `mediaTypes` will be inspected.
	IsResponseInspected *bool `mandatory:"false" json:"isResponseInspected"`

	// The maximum response size to be fully inspected, in binary kilobytes (KiB). Anything over this limit will be partially inspected. If unspecified, defaults to `1024`.
	MaxResponseSizeInKiB *int `mandatory:"false" json:"maxResponseSizeInKiB"`

	// The list of allowed HTTP methods. If unspecified, default to `[OPTIONS, GET, HEAD, POST]`.
	AllowedHttpMethods []ProtectionSettingsAllowedHttpMethodsEnum `mandatory:"false" json:"allowedHttpMethods,omitempty"`

	// The list of media types to allow for inspection, if `isResponseInspected` is enabled. Only responses with MIME types in this list will be inspected. If unspecified, defaults to `[`text/html`, `text/plain`, `text/xml`]`.
	//     Supported MIME types include:
	//     - text/html
	//     - text/plain
	//     - text/asp
	//     - text/css
	//     - text/x-script
	//     - application/json
	//     - text/webviewhtml
	//     - text/x-java-source
	//     - application/x-javascript
	//     - application/javascript
	//     - application/ecmascript
	//     - text/javascript
	//     - text/ecmascript
	//     - text/x-script.perl
	//     - text/x-script.phyton
	//     - application/plain
	//     - application/xml
	//     - text/xml
	MediaTypes []string `mandatory:"false" json:"mediaTypes"`
}

func (m ProtectionSettings) String() string {
	return common.PointerString(m)
}

// ProtectionSettingsBlockActionEnum Enum with underlying type: string
type ProtectionSettingsBlockActionEnum string

// Set of constants representing the allowable values for ProtectionSettingsBlockActionEnum
const (
	ProtectionSettingsBlockActionShowErrorPage   ProtectionSettingsBlockActionEnum = "SHOW_ERROR_PAGE"
	ProtectionSettingsBlockActionSetResponseCode ProtectionSettingsBlockActionEnum = "SET_RESPONSE_CODE"
)

var mappingProtectionSettingsBlockAction = map[string]ProtectionSettingsBlockActionEnum{
	"SHOW_ERROR_PAGE":   ProtectionSettingsBlockActionShowErrorPage,
	"SET_RESPONSE_CODE": ProtectionSettingsBlockActionSetResponseCode,
}

// GetProtectionSettingsBlockActionEnumValues Enumerates the set of values for ProtectionSettingsBlockActionEnum
func GetProtectionSettingsBlockActionEnumValues() []ProtectionSettingsBlockActionEnum {
	values := make([]ProtectionSettingsBlockActionEnum, 0)
	for _, v := range mappingProtectionSettingsBlockAction {
		values = append(values, v)
	}
	return values
}

// ProtectionSettingsAllowedHttpMethodsEnum Enum with underlying type: string
type ProtectionSettingsAllowedHttpMethodsEnum string

// Set of constants representing the allowable values for ProtectionSettingsAllowedHttpMethodsEnum
const (
	ProtectionSettingsAllowedHttpMethodsOptions  ProtectionSettingsAllowedHttpMethodsEnum = "OPTIONS"
	ProtectionSettingsAllowedHttpMethodsGet      ProtectionSettingsAllowedHttpMethodsEnum = "GET"
	ProtectionSettingsAllowedHttpMethodsHead     ProtectionSettingsAllowedHttpMethodsEnum = "HEAD"
	ProtectionSettingsAllowedHttpMethodsPost     ProtectionSettingsAllowedHttpMethodsEnum = "POST"
	ProtectionSettingsAllowedHttpMethodsPut      ProtectionSettingsAllowedHttpMethodsEnum = "PUT"
	ProtectionSettingsAllowedHttpMethodsDelete   ProtectionSettingsAllowedHttpMethodsEnum = "DELETE"
	ProtectionSettingsAllowedHttpMethodsTrace    ProtectionSettingsAllowedHttpMethodsEnum = "TRACE"
	ProtectionSettingsAllowedHttpMethodsConnect  ProtectionSettingsAllowedHttpMethodsEnum = "CONNECT"
	ProtectionSettingsAllowedHttpMethodsPatch    ProtectionSettingsAllowedHttpMethodsEnum = "PATCH"
	ProtectionSettingsAllowedHttpMethodsPropfind ProtectionSettingsAllowedHttpMethodsEnum = "PROPFIND"
)

var mappingProtectionSettingsAllowedHttpMethods = map[string]ProtectionSettingsAllowedHttpMethodsEnum{
	"OPTIONS":  ProtectionSettingsAllowedHttpMethodsOptions,
	"GET":      ProtectionSettingsAllowedHttpMethodsGet,
	"HEAD":     ProtectionSettingsAllowedHttpMethodsHead,
	"POST":     ProtectionSettingsAllowedHttpMethodsPost,
	"PUT":      ProtectionSettingsAllowedHttpMethodsPut,
	"DELETE":   ProtectionSettingsAllowedHttpMethodsDelete,
	"TRACE":    ProtectionSettingsAllowedHttpMethodsTrace,
	"CONNECT":  ProtectionSettingsAllowedHttpMethodsConnect,
	"PATCH":    ProtectionSettingsAllowedHttpMethodsPatch,
	"PROPFIND": ProtectionSettingsAllowedHttpMethodsPropfind,
}

// GetProtectionSettingsAllowedHttpMethodsEnumValues Enumerates the set of values for ProtectionSettingsAllowedHttpMethodsEnum
func GetProtectionSettingsAllowedHttpMethodsEnumValues() []ProtectionSettingsAllowedHttpMethodsEnum {
	values := make([]ProtectionSettingsAllowedHttpMethodsEnum, 0)
	for _, v := range mappingProtectionSettingsAllowedHttpMethods {
		values = append(values, v)
	}
	return values
}
