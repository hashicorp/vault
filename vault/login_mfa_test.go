// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/identity/mfa"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	otplib "github.com/pquerna/otp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFactors(t *testing.T) {
	testcases := []struct {
		name                string
		invalidMFAHeaderVal []string
		expectedError       string
	}{
		{
			"two headers with passcode",
			[]string{"passcode", "foo"},
			"found multiple passcodes for the same MFA method",
		},
		{
			"single header with passcode=",
			[]string{"passcode="},
			"invalid passcode",
		},
		{
			"single invalid header",
			[]string{"foo="},
			"found an invalid MFA cred",
		},
		{
			"single header equal char",
			[]string{"=="},
			"found an invalid MFA cred",
		},
		{
			"two headers with passcode=",
			[]string{"passcode=foo", "foo"},
			"found multiple passcodes for the same MFA method",
		},
		{
			"two headers invalid name",
			[]string{"passcode=foo", "passcode=bar"},
			"found multiple passcodes for the same MFA method",
		},
		{
			"two headers, two invalid",
			[]string{"foo", "bar"},
			"found multiple passcodes for the same MFA method",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseMfaFactors(tc.invalidMFAHeaderVal)
			if err == nil {
				t.Fatal("nil error returned")
			}
			if !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected %s, got %v", tc.expectedError, err)
			}
		})
	}
}

// TestMFAConfigToMap verifies that different MFA method configurations (TOTP,
// Okta, Duo, PingID) are correctly converted to a map format for API
// responses. The test covers both login MFA and policy MFA scenarios, since
// these share objects/protobufs, but can have slightly different responses.
func TestMFAConfigToMap(t *testing.T) {
	testCases := map[string]struct {
		config         *mfa.Config
		expectedResult map[string]any
		isLoginMFA     bool
	}{
		"totp-with-login-mfa": {
			config: &mfa.Config{
				Type: mfaMethodTypeTOTP,
				Config: &mfa.Config_TOTPConfig{
					TOTPConfig: &mfa.TOTPConfig{
						Issuer:                "TestIssuer",
						Period:                30,
						Digits:                6,
						Skew:                  1,
						KeySize:               20,
						QRSize:                200,
						Algorithm:             int32(otplib.AlgorithmSHA1),
						MaxValidationAttempts: 5,
						EnableSelfEnrollment:  true,
					},
				},
			},
			expectedResult: map[string]interface{}{
				"algorithm":               "SHA1",
				"digits":                  int32(6),
				"enable_self_enrollment":  true,
				"id":                      "",
				"issuer":                  "TestIssuer",
				"key_size":                uint32(20),
				"max_validation_attempts": uint32(5),
				"name":                    "",
				"namespace_id":            "",
				"namespace_path":          "/",
				"period":                  uint32(30),
				"qr_size":                 int32(200),
				"skew":                    uint32(1),
				"type":                    "totp",
			},
			isLoginMFA: true,
		},
		"totp-ent-step-up-mfa-self-enrollment-false": {
			config: &mfa.Config{
				Type: mfaMethodTypeTOTP,
				Config: &mfa.Config_TOTPConfig{
					TOTPConfig: &mfa.TOTPConfig{
						Issuer:                "TestIssuer",
						Period:                30,
						Digits:                6,
						Skew:                  1,
						KeySize:               20,
						QRSize:                200,
						Algorithm:             int32(otplib.AlgorithmSHA1),
						MaxValidationAttempts: 5,
						EnableSelfEnrollment:  false,
					},
				},
			},
			expectedResult: map[string]interface{}{
				"algorithm":               "SHA1",
				"digits":                  int32(6),
				"id":                      "",
				"issuer":                  "TestIssuer",
				"key_size":                uint32(20),
				"max_validation_attempts": uint32(5),
				"name":                    "",
				"namespace_id":            "",
				"namespace_path":          "/",
				"period":                  uint32(30),
				"qr_size":                 int32(200),
				"skew":                    uint32(1),
				"type":                    "totp",
			},
			isLoginMFA: false,
		},
		"totp-ent-step-up-mfa-self-enrollment-true": {
			config: &mfa.Config{
				Type: mfaMethodTypeTOTP,
				Config: &mfa.Config_TOTPConfig{
					TOTPConfig: &mfa.TOTPConfig{
						Issuer:                "TestIssuer",
						Period:                30,
						Digits:                6,
						Skew:                  1,
						KeySize:               20,
						QRSize:                200,
						Algorithm:             int32(otplib.AlgorithmSHA1),
						MaxValidationAttempts: 5,
						EnableSelfEnrollment:  true,
					},
				},
			},
			expectedResult: map[string]interface{}{
				"algorithm":               "SHA1",
				"digits":                  int32(6),
				"id":                      "",
				"issuer":                  "TestIssuer",
				"key_size":                uint32(20),
				"max_validation_attempts": uint32(5),
				"name":                    "",
				"namespace_id":            "",
				"namespace_path":          "/",
				"period":                  uint32(30),
				"qr_size":                 int32(200),
				"skew":                    uint32(1),
				"type":                    "totp",
			},
			isLoginMFA: false,
		},
		"okta-prod": {
			config: &mfa.Config{
				Type: mfaMethodTypeOkta,
				Config: &mfa.Config_OktaConfig{
					OktaConfig: &mfa.OktaConfig{
						OrgName:      "TestOrg",
						APIToken:     "test-token",
						PrimaryEmail: true,
						Production:   true,
					},
				},
			},
			expectedResult: map[string]interface{}{
				"production":      true,
				"id":              "",
				"mount_accessor":  "",
				"name":            "",
				"namespace_id":    "",
				"namespace_path":  "/",
				"org_name":        "TestOrg",
				"type":            "okta",
				"username_format": "",
			},
		},
		"okta-non-prod": {
			config: &mfa.Config{
				Type: mfaMethodTypeOkta,
				Config: &mfa.Config_OktaConfig{
					OktaConfig: &mfa.OktaConfig{
						OrgName:      "TestOrg",
						APIToken:     "test-token",
						BaseURL:      "https://test.okta.com",
						PrimaryEmail: true,
					},
				},
			},
			expectedResult: map[string]interface{}{
				"base_url":        "https://test.okta.com",
				"id":              "",
				"mount_accessor":  "",
				"name":            "",
				"namespace_id":    "",
				"namespace_path":  "/",
				"org_name":        "TestOrg",
				"type":            "okta",
				"username_format": "",
			},
		},
		"duo": {
			config: &mfa.Config{
				Type: mfaMethodTypeDuo,
				Config: &mfa.Config_DuoConfig{
					DuoConfig: &mfa.DuoConfig{
						APIHostname:    "api.duo.com",
						IntegrationKey: "integration-key",
						SecretKey:      "secret-key",
						PushInfo:       "push-info",
						UsePasscode:    true,
					},
				},
			},
			expectedResult: map[string]interface{}{
				"api_hostname":    "api.duo.com",
				"pushinfo":        "push-info",
				"mount_accessor":  "",
				"username_format": "",
				"use_passcode":    true,
				"type":            "duo",
				"id":              "",
				"name":            "",
				"namespace_id":    "",
				"namespace_path":  "/",
			},
			isLoginMFA: false,
		},
		"pingid": {
			config: &mfa.Config{
				Type: mfaMethodTypePingID,
				Config: &mfa.Config_PingIDConfig{
					PingIDConfig: &mfa.PingIDConfig{
						UseSignature:     true,
						IDPURL:           "https://idp.pingid.com",
						OrgAlias:         "org-alias",
						AdminURL:         "https://admin.pingid.com",
						AuthenticatorURL: "https://authenticator.pingid.com",
					},
				},
			},
			expectedResult: map[string]interface{}{
				"use_signature":     true,
				"idp_url":           "https://idp.pingid.com",
				"org_alias":         "org-alias",
				"admin_url":         "https://admin.pingid.com",
				"authenticator_url": "https://authenticator.pingid.com",
				"type":              "pingid",
				"id":                "",
				"name":              "",
				"namespace_id":      "",
				"namespace_path":    "/",
			},
			isLoginMFA: false,
		},
	}
	backend := &MFABackend{}
	backend.namespacer = &MockNamespacer{}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualResult, err := backend.mfaConfigToMap(tc.config, tc.isLoginMFA)
			require.NoError(t, err)
			require.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

// TestMfaConfigToMap_InvalidType validates that mfaConfigToMap returns an error
// when the method type is not recognized.
func TestMfaConfigToMap_InvalidType(t *testing.T) {
	backend := &MFABackend{}
	mConfig := &mfa.Config{
		Type:   "invalid-type",
		Config: nil,
	}

	_, err := backend.mfaConfigToMap(mConfig, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid method type")
}

// TestParseTOTPConfig validates that the parseTOTPConfig function returns the
// expected self-enrollment value given various combinations of isEnterprise and
// isLoginMFA inputs.
func TestParseTOTPConfig(t *testing.T) {
	testCases := map[string]struct {
		isEnterprise   bool
		isLoginMFA     bool
		expectedResult *mfa.TOTPConfig
		wantErr        bool
	}{
		"login-mfa-ent": {
			isEnterprise: true,
			isLoginMFA:   true,
			expectedResult: &mfa.TOTPConfig{
				Issuer:                "TestIssuer",
				Period:                30,
				Digits:                6,
				Skew:                  1,
				KeySize:               20,
				QRSize:                200,
				Algorithm:             int32(otplib.AlgorithmSHA1),
				MaxValidationAttempts: 5,
				EnableSelfEnrollment:  true,
			},
		},
		"login-mfa-ce": {
			isEnterprise:   false,
			isLoginMFA:     true,
			expectedResult: &mfa.TOTPConfig{},
			wantErr:        true,
		},
		"step-up-mfa-ent": {
			isEnterprise: true,
			isLoginMFA:   false,
			expectedResult: &mfa.TOTPConfig{
				Issuer:                "TestIssuer",
				Period:                30,
				Digits:                6,
				Skew:                  1,
				KeySize:               20,
				QRSize:                200,
				Algorithm:             int32(otplib.AlgorithmSHA1),
				MaxValidationAttempts: 5,
				EnableSelfEnrollment:  false,
			},
		},
		"step-up-mfa-ce": {
			isEnterprise: false,
			isLoginMFA:   false,
			expectedResult: &mfa.TOTPConfig{
				Issuer:                "TestIssuer",
				Period:                30,
				Digits:                6,
				Skew:                  1,
				KeySize:               20,
				QRSize:                200,
				Algorithm:             int32(otplib.AlgorithmSHA1),
				MaxValidationAttempts: 5,
				EnableSelfEnrollment:  false,
			},
			wantErr: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			rawData := map[string]interface{}{
				"algorithm":               "SHA1",
				"digits":                  6,
				"period":                  30,
				"skew":                    1,
				"key_size":                20,
				"issuer":                  "TestIssuer",
				"max_validation_attempts": 5,
				"enable_self_enrollment":  true,
				"qr_size":                 200,
			}
			config := &mfa.Config{}
			schema := mfaTOTPPaths(&IdentityStore{})[0].Fields
			data := &framework.FieldData{
				Raw:    rawData,
				Schema: schema,
			}
			err := parseTOTPConfig(config, data, tc.isEnterprise, tc.isLoginMFA)
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, config.Config)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedResult, config.Config.(*mfa.Config_TOTPConfig).TOTPConfig)
			}
		})
	}
}

// TestParseTOTPConfig_UnhappyPaths tests various error scenarios for the
// parseTOTPConfig function.
func TestParseTOTPConfig_UnhappyPaths(t *testing.T) {
	testCases := map[string]struct {
		config         *mfa.Config
		data           *framework.FieldData
		expectedErrMsg string
	}{
		"nil-config": {
			config:         nil,
			data:           nil,
			expectedErrMsg: "config is nil",
		},
		"nil-field-data": {
			config:         &mfa.Config{},
			data:           nil,
			expectedErrMsg: "field data is nil",
		},
		"invalid-algorithm": {
			config: &mfa.Config{},
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"algorithm":               "SHA007",
					"digits":                  6,
					"period":                  30,
					"skew":                    1,
					"key_size":                20,
					"issuer":                  "TestIssuer",
					"max_validation_attempts": 5,
					"enable_self_enrollment":  true,
					"qr_size":                 200,
				},
				Schema: mfaTOTPPaths(&IdentityStore{})[0].Fields,
			},
			expectedErrMsg: "unrecognized algorithm",
		},
		"invalid-digits": {
			config: &mfa.Config{},
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"algorithm":               "SHA1",
					"digits":                  7, // Invalid value, should be 6 or 8
					"period":                  30,
					"skew":                    1,
					"key_size":                20,
					"issuer":                  "TestIssuer",
					"max_validation_attempts": 5,
					"enable_self_enrollment":  true,
					"qr_size":                 200,
				},
				Schema: mfaTOTPPaths(&IdentityStore{})[0].Fields,
			},
			expectedErrMsg: "digits can only be 6 or 8",
		},
		"invalid-period": {
			config: &mfa.Config{},
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"algorithm":               "SHA1",
					"digits":                  6,
					"period":                  0, // Invalid value, should be positive
					"skew":                    1,
					"key_size":                20,
					"issuer":                  "TestIssuer",
					"max_validation_attempts": 5,
					"enable_self_enrollment":  true,
					"qr_size":                 200,
				},
				Schema: mfaTOTPPaths(&IdentityStore{})[0].Fields,
			},
			expectedErrMsg: "period must be greater than zero",
		},
		"invalid-skew": {
			config: &mfa.Config{},
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"algorithm":               "SHA1",
					"digits":                  6,
					"period":                  30,
					"skew":                    2, // Invalid value, should be 0 or 1
					"key_size":                20,
					"issuer":                  "TestIssuer",
					"max_validation_attempts": 5,
					"enable_self_enrollment":  true,
					"qr_size":                 200,
				},
				Schema: mfaTOTPPaths(&IdentityStore{})[0].Fields,
			},
			expectedErrMsg: "skew must be 0 or 1",
		},
		"invalid-key-size": {
			config: &mfa.Config{},
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"algorithm":               "SHA1",
					"digits":                  6,
					"period":                  30,
					"skew":                    1,
					"key_size":                0, // Invalid value, should be positive
					"issuer":                  "TestIssuer",
					"max_validation_attempts": 5,
					"enable_self_enrollment":  true,
					"qr_size":                 200,
				},
				Schema: mfaTOTPPaths(&IdentityStore{})[0].Fields,
			},
			expectedErrMsg: "key_size must be greater than zero",
		},
		"invalid-issuer": {
			config: &mfa.Config{},
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"algorithm":               "SHA1",
					"digits":                  6,
					"period":                  30,
					"skew":                    1,
					"key_size":                20,
					"issuer":                  "", // Invalid value, should not be empty
					"max_validation_attempts": 5,
					"enable_self_enrollment":  true,
					"qr_size":                 200,
				},
				Schema: mfaTOTPPaths(&IdentityStore{})[0].Fields,
			},
			expectedErrMsg: "issuer must be set",
		},
		"invalid-max-validation-attempts": {
			config: &mfa.Config{},
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"algorithm":               "SHA1",
					"digits":                  6,
					"period":                  30,
					"skew":                    1,
					"key_size":                20,
					"issuer":                  "TestIssuer",
					"max_validation_attempts": -1, // Invalid value, should be greater than zero
					"enable_self_enrollment":  true,
					"qr_size":                 200,
				},
				Schema: mfaTOTPPaths(&IdentityStore{})[0].Fields,
			},
			expectedErrMsg: "max_validation_attempts must be greater than zero",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := parseTOTPConfig(tc.config, tc.data, true, true)
			require.Error(t, err)
			require.ErrorContains(t, err, tc.expectedErrMsg)
		})
	}
}

// MockNamespacer is a mock implementation of the Namespacer interface for testing purposes.
type MockNamespacer struct {
	Namespacer
}

// NamespaceByID returns a namespace with the same ID as the passed argument and with a path of `/`.
func (m *MockNamespacer) NamespaceByID(ctx context.Context, id string) (*namespace.Namespace, error) {
	return &namespace.Namespace{
		ID:   id,
		Path: "/",
	}, nil
}
