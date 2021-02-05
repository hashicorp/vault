package packngo

const twoFactorAuthAppPath = "/user/otp/app"
const twoFactorAuthSmsPath = "/user/otp/sms"

// TwoFactorAuthService interface defines available two factor authentication functions
type TwoFactorAuthService interface {
	EnableApp(string) (*Response, error)
	DisableApp(string) (*Response, error)
	EnableSms(string) (*Response, error)
	DisableSms(string) (*Response, error)
	ReceiveSms() (*Response, error)
	SeedApp() (string, *Response, error)
}

// TwoFactorAuthServiceOp implements TwoFactorAuthService
type TwoFactorAuthServiceOp struct {
	client *Client
}

// EnableApp function enables two factor auth using authenticatior app
func (s *TwoFactorAuthServiceOp) EnableApp(token string) (resp *Response, err error) {
	headers := map[string]string{"x-otp-token": token}
	return s.client.DoRequestWithHeader("POST", headers, twoFactorAuthAppPath, nil, nil)
}

// EnableSms function enables two factor auth using sms
func (s *TwoFactorAuthServiceOp) EnableSms(token string) (resp *Response, err error) {
	headers := map[string]string{"x-otp-token": token}
	return s.client.DoRequestWithHeader("POST", headers, twoFactorAuthSmsPath, nil, nil)
}

// ReceiveSms orders the auth service to issue an SMS token
func (s *TwoFactorAuthServiceOp) ReceiveSms() (resp *Response, err error) {
	return s.client.DoRequest("POST", twoFactorAuthSmsPath+"/receive", nil, nil)
}

// DisableApp function disables two factor auth using
func (s *TwoFactorAuthServiceOp) DisableApp(token string) (resp *Response, err error) {
	headers := map[string]string{"x-otp-token": token}
	return s.client.DoRequestWithHeader("DELETE", headers, twoFactorAuthAppPath, nil, nil)
}

// DisableSms function disables two factor auth using
func (s *TwoFactorAuthServiceOp) DisableSms(token string) (resp *Response, err error) {
	headers := map[string]string{"x-otp-token": token}
	return s.client.DoRequestWithHeader("DELETE", headers, twoFactorAuthSmsPath, nil, nil)
}

// SeedApp orders the auth service to issue a token via google authenticator
func (s *TwoFactorAuthServiceOp) SeedApp() (otpURI string, resp *Response, err error) {
	ret := &map[string]string{}
	resp, err = s.client.DoRequest("POST", twoFactorAuthAppPath+"/receive", nil, ret)

	return (*ret)["otp_uri"], resp, err
}
