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

// DeviceFingerprintChallenge The device fingerprint challenge settings. The device fingerprint challenge generates hashed signatures of both virtual and real browsers to identify and block malicious bots.
type DeviceFingerprintChallenge struct {

	// Enables or disables the device fingerprint challenge Web Application Firewall feature.
	IsEnabled *bool `mandatory:"true" json:"isEnabled"`

	// The action to take on requests from detected bots. If unspecified, defaults to `DETECT`.
	Action DeviceFingerprintChallengeActionEnum `mandatory:"false" json:"action,omitempty"`

	// The number of failed requests allowed before taking action. If unspecified, defaults to `10`.
	FailureThreshold *int `mandatory:"false" json:"failureThreshold"`

	// The number of seconds between challenges for the same IP address. If unspecified, defaults to `60`.
	ActionExpirationInSeconds *int `mandatory:"false" json:"actionExpirationInSeconds"`

	// The number of seconds before the failure threshold resets. If unspecified, defaults to `60`.
	FailureThresholdExpirationInSeconds *int `mandatory:"false" json:"failureThresholdExpirationInSeconds"`

	// The maximum number of IP addresses permitted with the same device fingerprint. If unspecified, defaults to `20`.
	MaxAddressCount *int `mandatory:"false" json:"maxAddressCount"`

	// The number of seconds before the maximum addresses count resets. If unspecified, defaults to `60`.
	MaxAddressCountExpirationInSeconds *int `mandatory:"false" json:"maxAddressCountExpirationInSeconds"`

	ChallengeSettings *BlockChallengeSettings `mandatory:"false" json:"challengeSettings"`
}

func (m DeviceFingerprintChallenge) String() string {
	return common.PointerString(m)
}

// DeviceFingerprintChallengeActionEnum Enum with underlying type: string
type DeviceFingerprintChallengeActionEnum string

// Set of constants representing the allowable values for DeviceFingerprintChallengeActionEnum
const (
	DeviceFingerprintChallengeActionDetect DeviceFingerprintChallengeActionEnum = "DETECT"
	DeviceFingerprintChallengeActionBlock  DeviceFingerprintChallengeActionEnum = "BLOCK"
)

var mappingDeviceFingerprintChallengeAction = map[string]DeviceFingerprintChallengeActionEnum{
	"DETECT": DeviceFingerprintChallengeActionDetect,
	"BLOCK":  DeviceFingerprintChallengeActionBlock,
}

// GetDeviceFingerprintChallengeActionEnumValues Enumerates the set of values for DeviceFingerprintChallengeActionEnum
func GetDeviceFingerprintChallengeActionEnumValues() []DeviceFingerprintChallengeActionEnum {
	values := make([]DeviceFingerprintChallengeActionEnum, 0)
	for _, v := range mappingDeviceFingerprintChallengeAction {
		values = append(values, v)
	}
	return values
}
