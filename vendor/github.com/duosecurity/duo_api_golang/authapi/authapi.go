package authapi

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/duosecurity/duo_api_golang"
)

type AuthApi struct {
	duoapi.DuoApi
}

// Build a new Duo Auth API object.
// api is a duoapi.DuoApi object used to make the Duo Rest API calls.
// Example: authapi.NewAuthApi(*duoapi.NewDuoApi(ikey,skey,host,userAgent,duoapi.SetTimeout(10*time.Second)))
func NewAuthApi(api duoapi.DuoApi) *AuthApi {
	return &AuthApi{api}
}

// API calls will return a StatResult object.  On success, Stat is 'OK'.
// On error, Stat is 'FAIL', and Code, Message, and Message_Detail
// contain error information.
type StatResult struct {
	Stat           string
	Code           *int32
	Message        *string
	Message_Detail *string
}

// Return object for the 'Ping' API call.
type PingResult struct {
	StatResult
	Response struct {
		Time int64
	}
}

// Duo's Ping method. https://www.duosecurity.com/docs/authapi#/ping
// This is an unsigned Duo Rest API call which returns the Duo system's time.
// Use this method to determine whether your system time is in sync with Duo's.
func (api *AuthApi) Ping() (*PingResult, error) {
	_, body, err := api.Call("GET", "/auth/v2/ping", nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}
	ret := &PingResult{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Return object for the 'Check' API call.
type CheckResult struct {
	StatResult
	Response struct {
		Time int64
	}
}

// Call Duo's Check method. https://www.duosecurity.com/docs/authapi#/check
// Check is a signed Duo API call, which returns the Duo system's time.
// Use this method to determine whether your ikey, skey and host are correct,
// and whether your system time is in sync with Duo's.
func (api *AuthApi) Check() (*CheckResult, error) {
	_, body, err := api.SignedCall("GET", "/auth/v2/check", nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}
	ret := &CheckResult{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Return object for the 'Logo' API call.
type LogoResult struct {
	StatResult
	png *[]byte
}

// Duo's Logo method. https://www.duosecurity.com/docs/authapi#/logo
// If the API call is successful, the configured logo png is returned.  Othwerwise,
// error information is returned in the LogoResult return value.
func (api *AuthApi) Logo() (*LogoResult, error) {
	resp, body, err := api.SignedCall("GET", "/auth/v2/logo", nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		ret := &LogoResult{StatResult: StatResult{Stat: "OK"},
			png: &body}
		return ret, nil
	}
	ret := &LogoResult{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Optional parameter for the Enroll method.
func EnrollUsername(username string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("username", username)
	}
}

// Optional parameter for the Enroll method.
func EnrollValidSeconds(secs uint64) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("valid_secs", strconv.FormatUint(secs, 10))
	}
}

// Enroll return type.
type EnrollResult struct {
	StatResult
	Response struct {
		Activation_Barcode string
		Activation_Code    string
		Expiration         int64
		User_Id            string
		Username           string
	}
}

// Duo's Enroll method. https://www.duosecurity.com/docs/authapi#/enroll
// Use EnrollUsername() to include the optional username parameter.
// Use EnrollValidSeconds() to change the default validation time limit that the
// user has to complete enrollment.
func (api *AuthApi) Enroll(options ...func(*url.Values)) (*EnrollResult, error) {
	opts := url.Values{}
	for _, o := range options {
		o(&opts)
	}

	_, body, err := api.SignedCall("POST", "/auth/v2/enroll", opts, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}
	ret := &EnrollResult{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Response is "success", "invalid" or "waiting".
type EnrollStatusResult struct {
	StatResult
	Response string
}

// Duo's EnrollStatus method. https://www.duosecurity.com/docs/authapi#/enroll_status
// Return the status of an outstanding Enrollment.
func (api *AuthApi) EnrollStatus(userid string,
	activationCode string) (*EnrollStatusResult, error) {
	queryArgs := url.Values{}
	queryArgs.Set("user_id", userid)
	queryArgs.Set("activation_code", activationCode)

	_, body, err := api.SignedCall("POST",
		"/auth/v2/enroll_status",
		queryArgs,
		duoapi.UseTimeout)

	if err != nil {
		return nil, err
	}
	ret := &EnrollStatusResult{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Preauth return type.
type PreauthResult struct {
	StatResult
	Response struct {
		Result            string
		Status_Msg        string
		Enroll_Portal_Url string
		Devices           []struct {
			Device       string
			Type         string
			Name         string
			Number       string
			Capabilities []string
		}
	}
}

func PreauthUserId(userid string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("user_id", userid)
	}
}

func PreauthUsername(username string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("username", username)
	}
}

func PreauthIpAddr(ip string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("ipaddr", ip)
	}
}

func PreauthTrustedToken(trustedtoken string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("trusted_device_token", trustedtoken)
	}
}

// Duo's Preauth method. https://www.duosecurity.com/docs/authapi#/preauth
// options Optional values to include in the preauth call.
// Use PreauthUserId to specify the user_id parameter.
// Use PreauthUsername to specify the username parameter.  You must
// specify PreauthUserId or PreauthUsername, but not both.
// Use PreauthIpAddr to include the ipaddr parameter, the ip address
// of the client attempting authroization.
// Use PreauthTrustedToken to specify the trusted_device_token parameter.
func (api *AuthApi) Preauth(options ...func(*url.Values)) (*PreauthResult, error) {
	opts := url.Values{}
	for _, o := range options {
		o(&opts)
	}
	_, body, err := api.SignedCall("POST", "/auth/v2/preauth", opts, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}
	ret := &PreauthResult{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func AuthUserId(userid string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("user_id", userid)
	}
}

func AuthUsername(username string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("username", username)
	}
}

func AuthIpAddr(ip string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("ipaddr", ip)
	}
}

func AuthAsync() func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("async", "1")
	}
}

func AuthDevice(device string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("device", device)
	}
}

func AuthType(type_ string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("type", type_)
	}
}

func AuthDisplayUsername(username string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("display_username", username)
	}
}

func AuthPushinfo(pushinfo string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("pushinfo", pushinfo)
	}
}

func AuthPasscode(passcode string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("passcode", passcode)
	}
}

// Auth return type.
type AuthResult struct {
	StatResult
	Response struct {
		// Synchronous
		Result               string
		Status               string
		Status_Msg           string
		Trusted_Device_Token string
		// Asynchronous
		Txid string
	}
}

// Duo's Auth method. https://www.duosecurity.com/docs/authapi#/auth
// Factor must be one of 'auto', 'push', 'passcode', 'sms' or 'phone'.
// Use AuthUserId to specify the user_id.
// Use AuthUsername to speicy the username.  You must specify either AuthUserId
// or AuthUsername, but not both.
// Use AuthIpAddr to include the client's IP address.
// Use AuthAsync to toggle whether the call blocks for the user's response or not.
// If used asynchronously, get the auth status with the AuthStatus method.
// When using factor 'push', use AuthDevice to specify the device ID to push to.
// When using factor 'push', use AuthType to display some extra auth text to the user.
// When using factor 'push', use AuthDisplayUsername to display some extra text
// to the user.
// When using factor 'push', use AuthPushInfo to include some URL-encoded key/value
// pairs to display to the user.
// When using factor 'passcode', use AuthPasscode to specify the passcode entered
// by the user.
// When using factor 'sms' or 'phone', use AuthDevice to specify which device
// should receive the SMS or phone call.
func (api *AuthApi) Auth(factor string, options ...func(*url.Values)) (*AuthResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}
	params.Set("factor", factor)

	var apiOps []duoapi.DuoApiOption
	if _, ok := params["async"]; ok == true {
		apiOps = append(apiOps, duoapi.UseTimeout)
	}

	_, body, err := api.SignedCall("POST", "/auth/v2/auth", params, apiOps...)
	if err != nil {
		return nil, err
	}
	ret := &AuthResult{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// AuthStatus return type.
type AuthStatusResult struct {
	StatResult
	Response struct {
		Result               string
		Status               string
		Status_Msg           string
		Trusted_Device_Token string
	}
}

// Duo's auth_status method.  https://www.duosecurity.com/docs/authapi#/auth_status
// When using the Auth call in async mode, use this method to retrieve the
// result of the authentication attempt.
// txid is returned by the Auth call.
func (api *AuthApi) AuthStatus(txid string) (*AuthStatusResult, error) {
	opts := url.Values{}
	opts.Set("txid", txid)
	_, body, err := api.SignedCall("GET", "/auth/v2/auth_status", opts)
	if err != nil {
		return nil, err
	}
	ret := &AuthStatusResult{}
	if err = json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
