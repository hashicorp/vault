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

// BlockChallengeSettings The challenge settings if `action` is set to `BLOCK`.
type BlockChallengeSettings struct {

	// The method used to block requests that fail the challenge, if `action` is set to `BLOCK`. If unspecified, defaults to `SHOW_ERROR_PAGE`.
	BlockAction BlockChallengeSettingsBlockActionEnum `mandatory:"false" json:"blockAction,omitempty"`

	// The response status code to return when `action` is set to `BLOCK`, `blockAction` is set to `SET_RESPONSE_CODE` or `SHOW_ERROR_PAGE`, and the request is blocked. If unspecified, defaults to `403`.
	BlockResponseCode *int `mandatory:"false" json:"blockResponseCode"`

	// The message to show on the error page when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_ERROR_PAGE`, and the request is blocked. If unspecified, defaults to `Access to the website is blocked`.
	BlockErrorPageMessage *string `mandatory:"false" json:"blockErrorPageMessage"`

	// The description text to show on the error page when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_ERROR_PAGE`, and the request is blocked. If unspecified, defaults to `Access blocked by website owner. Please contact support.`
	BlockErrorPageDescription *string `mandatory:"false" json:"blockErrorPageDescription"`

	// The error code to show on the error page when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_ERROR_PAGE` and the request is blocked. If unspecified, defaults to `403`.
	BlockErrorPageCode *string `mandatory:"false" json:"blockErrorPageCode"`

	// The title used when showing a CAPTCHA challenge when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_CAPTCHA`, and the request is blocked. If unspecified, defaults to `Are you human?`
	CaptchaTitle *string `mandatory:"false" json:"captchaTitle"`

	// The text to show in the header when showing a CAPTCHA challenge when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_CAPTCHA`, and the request is blocked. If unspecified, defaults to `We have detected an increased number of attempts to access this webapp. To help us keep this webapp secure, please let us know that you are not a robot by entering the text from captcha below.`
	CaptchaHeader *string `mandatory:"false" json:"captchaHeader"`

	// The text to show in the footer when showing a CAPTCHA challenge when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_CAPTCHA`, and the request is blocked. If unspecified, default to `Enter the letters and numbers as they are shown in image above`.
	CaptchaFooter *string `mandatory:"false" json:"captchaFooter"`

	// The text to show on the label of the CAPTCHA challenge submit button when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_CAPTCHA`, and the request is blocked. If unspecified, defaults to `Yes, I am human`.
	CaptchaSubmitLabel *string `mandatory:"false" json:"captchaSubmitLabel"`
}

func (m BlockChallengeSettings) String() string {
	return common.PointerString(m)
}

// BlockChallengeSettingsBlockActionEnum Enum with underlying type: string
type BlockChallengeSettingsBlockActionEnum string

// Set of constants representing the allowable values for BlockChallengeSettingsBlockActionEnum
const (
	BlockChallengeSettingsBlockActionSetResponseCode BlockChallengeSettingsBlockActionEnum = "SET_RESPONSE_CODE"
	BlockChallengeSettingsBlockActionShowErrorPage   BlockChallengeSettingsBlockActionEnum = "SHOW_ERROR_PAGE"
	BlockChallengeSettingsBlockActionShowCaptcha     BlockChallengeSettingsBlockActionEnum = "SHOW_CAPTCHA"
)

var mappingBlockChallengeSettingsBlockAction = map[string]BlockChallengeSettingsBlockActionEnum{
	"SET_RESPONSE_CODE": BlockChallengeSettingsBlockActionSetResponseCode,
	"SHOW_ERROR_PAGE":   BlockChallengeSettingsBlockActionShowErrorPage,
	"SHOW_CAPTCHA":      BlockChallengeSettingsBlockActionShowCaptcha,
}

// GetBlockChallengeSettingsBlockActionEnumValues Enumerates the set of values for BlockChallengeSettingsBlockActionEnum
func GetBlockChallengeSettingsBlockActionEnumValues() []BlockChallengeSettingsBlockActionEnum {
	values := make([]BlockChallengeSettingsBlockActionEnum, 0)
	for _, v := range mappingBlockChallengeSettingsBlockAction {
		values = append(values, v)
	}
	return values
}
