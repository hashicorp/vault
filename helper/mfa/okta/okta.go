// Package okta provides an Okta MFA handler to authenticate users
// with Okta Push or OTP. This handler is registered as the "okta" type in
// mfa_config.
package okta

import (
	"github.com/chrismalek/oktasdk-go/okta"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"time"
)

// OktaPaths returns path functions to configure Okta.
func OktaPaths() []*framework.Path {
	return []*framework.Path{
		pathOktaConfig(),
	}
}

// OktaRootPaths returns the paths that are used to configure Okta.
func OktaRootPaths() []string {
	return []string{
		"okta/config",
		"okta/debug",
	}
}

type userMFAFactor struct {
	ID          string    `json:"id,omitempty"`
	FactorType  string    `json:"factorType,omitempty"`
	Provider    string    `json:"provider,omitempty"`
	VendorName  string    `json:"vendorName,omitempty"`
	Status      string    `json:"status,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	Links       struct {
		Self   FactorLink `json:"self"`
		Verify FactorLink `json:"verify"`
	} `json:"_links"`
	Profile struct {
		CredentialID string `json:"credentialId,omitempty"`
	} `json:"profile,omitempty"`
}

// OktaHandler interacts with the Okta Auth API to authenticate a user
// login request. If successful, the original response from the login
// backend is returned.
func OktaHandler(req *logical.Request, d *framework.FieldData, resp *logical.Response) (
	*logical.Response, error) {
	response := &oktaAuthRequest{}
	response.successResp = resp
	response.username, _ = resp.Auth.Metadata["username"]
	response.method = d.Get("method").(string)
	response.passcode = d.Get("passcode").(string)
	response.ipAddr = req.Connection.RemoteAddr

	config, err := GetOktaConfig(req)
	if err != nil || config == nil {
		return logical.ErrorResponse("Could not load Okta MFA configuration"), nil
	}

	client := okta.NewClientWithBaseURL(
		nil,
		config.BaseURL,
		config.ApiToken,
	)

	var userFactors []userMFAFactor

	// Don't need this, just get MFAs with uid.
	request, err := client.NewRequest("GET",
		"/api/v1/users/"+resp.Data["uid"].(string)+"/factors",
		nil)
	if err != nil {
		println(err.Error())
		return logical.ErrorResponse("Error creating request for Okta."), nil
	}

	_, err = client.Do(request, &userFactors)
	if err != nil {
		println(err.Error())
		return logical.ErrorResponse("Error retrieving user from Okta."), nil
	}

	var factorToUse (userMFAFactor)
	var authFactor AuthFactor
	factorSet := false

	for _, factor := range userFactors {
		if factor.Status == "ACTIVE" {
			switch {
			case !factorSet || (factor.FactorType == "push" && response.passcode == ""):
				pushFactor := PushFactor{
					Mfa: factor,
				}
				authFactor = &pushFactor
				factorToUse = factor
				factorSet = true
			case !factorSet || response.passcode != "" || ((factor.FactorType == "token:software:totp") &&
				(factorToUse.FactorType != "push")):
				otpFactor := OTPFactor{
					Mfa:      factor,
					Passcode: response.passcode,
				}
				authFactor = &otpFactor
				factorToUse = factor
				factorSet = true
			}
		}
	}

	if !factorSet {
		return logical.ErrorResponse("No valid MFA available from Okta."), nil
	}

	return authFactor.DoAuth(client, response), nil
}

type AuthFactor interface {
	DoAuth(client *okta.Client, response *oktaAuthRequest) *logical.Response
}

type PushFactor struct {
	Mfa userMFAFactor
}

// Handles Okta Verify with Push MFA attempts
func (pf *PushFactor) DoAuth(client *okta.Client, response *oktaAuthRequest) *logical.Response {
	request, _ := client.NewRequest(
		"POST",
		pf.Mfa.Links.Self.URL+"/verify",
		nil,
	)

	var factorReply mfaResponse
	_, err := client.Do(request, &factorReply)
	if err != nil {
		println(err.Error())
		return logical.ErrorResponse("Error in request")
	}

	err = blockUntilMFAVerified(client, factorReply)
	if err != nil {
		return logical.ErrorResponse(err.Error())
	}

	return response.successResp
}

type OTPFactor struct {
	Mfa      userMFAFactor
	Passcode string
}

type passCode struct {
	PassCode string `json:"passCode"`
}

// Handles Okta Verify with One Time Password (OTP) MFA attempts
func (otpf *OTPFactor) DoAuth(client *okta.Client, response *oktaAuthRequest) *logical.Response {

	pass := passCode{
		PassCode: otpf.Passcode,
	}

	request, err := client.NewRequest(
		"POST",
		otpf.Mfa.Links.Verify.URL,
		pass,
	)
	if err != nil {
		println(err.Error())
		return logical.ErrorResponse("Error in request")
	}

	var factorReply mfaResponse
	resp, err := client.Do(request, &factorReply)
	if resp.Status == "403 Forbidden" {
		return logical.ErrorResponse("Invalid MFA passcode.")
	}
	if err != nil {
		println(err.Error())
		return logical.ErrorResponse("Error in request")
	}

	if factorReply.FactorResult == "SUCCESS" {
		return response.successResp
	} else {
		return logical.ErrorResponse("Error in request")
	}
}

// Blocks thread until MFA is verified, or error is returned.
func blockUntilMFAVerified(client *okta.Client, response mfaResponse) error {
	request, _ := client.NewRequest("GET",
		response.Links.Poll.URL,
		nil)

	var mfaValidationResponse mfaResponse
	_, err := client.Do(request, &mfaValidationResponse)
	if err != nil {
		return err
	}

	// Waits for up to 5 minutes
	for i := 0; i < 300; i++ {
		client.Do(request, &mfaValidationResponse)
		switch mfaValidationResponse.FactorResult {
		case "SUCCESS":
			return nil
		case "REJECTED":
			var err mfaError
			err.error = "Could not authenticate Okta user."
			return error(&err)
		case "TIMEOUT":
			var err mfaError
			err.error = "Could not authenticate Okta user."
			return error(&err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	var unknownError mfaError
	unknownError.error = "Could not authenticate Okta user in time."
	return error(&unknownError)
}

type mfaResponse struct {
	FactorResult string   `json:"factorResult"`
	Links        MFALinks `json:"_links"`
}

type mfaError struct {
	error string
}

func (mfa *mfaError) Error() string {
	return mfa.error
}

type FactorLink struct {
	URL   string `json:"href"`
	Hints struct {
		Allow []string `json:"allow"`
	} `json:"hints"`
}

type MFALinks struct {
	Poll   FactorLink `json:"poll,omitempty"`
	Cancel FactorLink `json:"cancel,omitempty"`
	Verify FactorLink `json:"verify,omitempty"`
	Factor FactorLink `json:"factor,omitempty"`
}

type FactorStatus struct {
	Expires string   `json:"expiresAt"`
	Result  string   `json:"factorResult"`
	Links   MFALinks `json:"_links"`
}

type oktaAuthRequest struct {
	successResp *logical.Response
	username    string
	method      string
	passcode    string
	ipAddr      string
}
