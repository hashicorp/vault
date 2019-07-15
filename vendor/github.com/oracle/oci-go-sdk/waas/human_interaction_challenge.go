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

// HumanInteractionChallenge The human interaction challenge settings. The human interaction challenge checks various event listeners in the user's browser to determine if there is a human user making a request.
type HumanInteractionChallenge struct {

	// Enables or disables the human interaction challenge Web Application Firewall feature.
	IsEnabled *bool `mandatory:"true" json:"isEnabled"`

	// The action to take against requests from detected bots. If unspecified, defaults to `DETECT`.
	Action HumanInteractionChallengeActionEnum `mandatory:"false" json:"action,omitempty"`

	// The number of failed requests before taking action. If unspecified, defaults to `10`.
	FailureThreshold *int `mandatory:"false" json:"failureThreshold"`

	// The number of seconds between challenges for the same IP address. If unspecified, defaults to `60`.
	ActionExpirationInSeconds *int `mandatory:"false" json:"actionExpirationInSeconds"`

	// The number of seconds before the failure threshold resets. If unspecified, defaults to  `60`.
	FailureThresholdExpirationInSeconds *int `mandatory:"false" json:"failureThresholdExpirationInSeconds"`

	// The number of interactions required to pass the challenge. If unspecified, defaults to `3`.
	InteractionThreshold *int `mandatory:"false" json:"interactionThreshold"`

	// The number of seconds to record the interactions from the user. If unspecified, defaults to `15`.
	RecordingPeriodInSeconds *int `mandatory:"false" json:"recordingPeriodInSeconds"`

	// Adds an additional HTTP header to requests that fail the challenge before being passed to the origin. Only applicable when the `action` is set to `DETECT`.
	SetHttpHeader *Header `mandatory:"false" json:"setHttpHeader"`

	ChallengeSettings *BlockChallengeSettings `mandatory:"false" json:"challengeSettings"`
}

func (m HumanInteractionChallenge) String() string {
	return common.PointerString(m)
}

// HumanInteractionChallengeActionEnum Enum with underlying type: string
type HumanInteractionChallengeActionEnum string

// Set of constants representing the allowable values for HumanInteractionChallengeActionEnum
const (
	HumanInteractionChallengeActionDetect HumanInteractionChallengeActionEnum = "DETECT"
	HumanInteractionChallengeActionBlock  HumanInteractionChallengeActionEnum = "BLOCK"
)

var mappingHumanInteractionChallengeAction = map[string]HumanInteractionChallengeActionEnum{
	"DETECT": HumanInteractionChallengeActionDetect,
	"BLOCK":  HumanInteractionChallengeActionBlock,
}

// GetHumanInteractionChallengeActionEnumValues Enumerates the set of values for HumanInteractionChallengeActionEnum
func GetHumanInteractionChallengeActionEnumValues() []HumanInteractionChallengeActionEnum {
	values := make([]HumanInteractionChallengeActionEnum, 0)
	for _, v := range mappingHumanInteractionChallengeAction {
		values = append(values, v)
	}
	return values
}
