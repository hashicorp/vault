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

// WafLog A list of Web Application Firewall log entries. Each entry is a JSON object whose fields vary based on log type. Logs record what rules and countermeasures are triggered by requests and are used as a basis to move request handling into block mode.
type WafLog struct {

	// The action taken on the request.
	Action *string `mandatory:"false" json:"action"`

	// The CAPTCHA action taken on the request.
	CaptchaAction *string `mandatory:"false" json:"captchaAction"`

	// The CAPTCHA challenge answer that was expected.
	CaptchaExpected *string `mandatory:"false" json:"captchaExpected"`

	// The CAPTCHA challenge answer that was received.
	CaptchaReceived *string `mandatory:"false" json:"captchaReceived"`

	// The number of times the CAPTCHA challenge was failed.
	CaptchaFailCount *string `mandatory:"false" json:"captchaFailCount"`

	// The IPv4 address of the requesting client.
	ClientAddress *string `mandatory:"false" json:"clientAddress"`

	// The name of the country where the request was made.
	CountryName *string `mandatory:"false" json:"countryName"`

	// The `User-Agent` header value of the request.
	UserAgent *string `mandatory:"false" json:"userAgent"`

	// The domain where the request was sent.
	Domain *string `mandatory:"false" json:"domain"`

	// A map of protection rule keys to detection message details.
	ProtectionRuleDetections map[string]string `mandatory:"false" json:"protectionRuleDetections"`

	// The HTTP method of the request.
	HttpMethod *string `mandatory:"false" json:"httpMethod"`

	// The path and query string of the request.
	RequestUrl *string `mandatory:"false" json:"requestUrl"`

	// The map of header names to values of the request sent to the origin.
	HttpHeaders map[string]string `mandatory:"false" json:"httpHeaders"`

	// The `Referrer` header value of the request.
	Referrer *string `mandatory:"false" json:"referrer"`

	// The status code of the response.
	ResponseCode *int `mandatory:"false" json:"responseCode"`

	// The size in bytes of the response.
	ResponseSize *int `mandatory:"false" json:"responseSize"`

	// The incident key that matched the request.
	IncidentKey *string `mandatory:"false" json:"incidentKey"`

	// TODO: what is this? MD5 hash of the request? SHA1?
	Fingerprint *string `mandatory:"false" json:"fingerprint"`

	// The type of device that the request was made from.
	Device *string `mandatory:"false" json:"device"`

	// The ISO 3166-1 country code of the request.
	CountryCode *string `mandatory:"false" json:"countryCode"`

	// A map of header names to values of the original request.
	RequestHeaders map[string]string `mandatory:"false" json:"requestHeaders"`

	// The `ThreatFeed` key that matched the request.
	ThreatFeedKey *string `mandatory:"false" json:"threatFeedKey"`

	// The `AccessRule` key that matched the request.
	AccessRuleKey *string `mandatory:"false" json:"accessRuleKey"`

	// The `AddressRateLimiting` key that matched the request.
	AddressRateLimitingKey *string `mandatory:"false" json:"addressRateLimitingKey"`

	// The `Date` header value of the request.
	Timestamp *string `mandatory:"false" json:"timestamp"`

	// The type of log of the request.
	LogType *string `mandatory:"false" json:"logType"`

	// The address of the origin server where the request was sent.
	OriginAddress *string `mandatory:"false" json:"originAddress"`

	// The amount of time it took the origin server to respond to the request.
	// TODO: determine unit of time and example
	OriginResponseTime *string `mandatory:"false" json:"originResponseTime"`
}

func (m WafLog) String() string {
	return common.PointerString(m)
}
