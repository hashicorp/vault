package linodego

import (
	"context"
)

// SendPhoneNumberVerificationCodeOptions fields are those accepted by SendPhoneNumberVerificationCode
type SendPhoneNumberVerificationCodeOptions struct {
	ISOCode     string `json:"iso_code"`
	PhoneNumber string `json:"phone_number"`
}

// VerifyPhoneNumberOptions fields are those accepted by VerifyPhoneNumber
type VerifyPhoneNumberOptions struct {
	OTPCode string `json:"otp_code"`
}

// SendPhoneNumberVerificationCode sends a one-time verification code via SMS message to the submitted phone number.
func (c *Client) SendPhoneNumberVerificationCode(ctx context.Context, opts SendPhoneNumberVerificationCodeOptions) error {
	e := "profile/phone-number"
	_, err := doPOSTRequest[any](ctx, c, e, opts)

	return err
}

// DeletePhoneNumber deletes the verified phone number for the User making this request.
func (c *Client) DeletePhoneNumber(ctx context.Context) error {
	e := "profile/phone-number"
	err := doDELETERequest(ctx, c, e)
	return err
}

// VerifyPhoneNumber verifies a phone number by confirming the one-time code received via SMS message after accessing the Phone Verification Code Send command.
func (c *Client) VerifyPhoneNumber(ctx context.Context, opts VerifyPhoneNumberOptions) error {
	e := "profile/phone-number/verify"
	_, err := doPOSTRequest[any](ctx, c, e, opts)

	return err
}
