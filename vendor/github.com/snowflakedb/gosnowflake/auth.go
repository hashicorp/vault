// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	clientType = "Go"
)

const (
	idToken                        = "ID_TOKEN"
	mfaToken                       = "MFATOKEN"
	clientStoreTemporaryCredential = "CLIENT_STORE_TEMPORARY_CREDENTIAL"
	clientRequestMfaToken          = "CLIENT_REQUEST_MFA_TOKEN"
	idTokenAuthenticator           = "ID_TOKEN"
)

// AuthType indicates the type of authentication in Snowflake
type AuthType int

const (
	// AuthTypeSnowflake is the general username password authentication
	AuthTypeSnowflake AuthType = iota
	// AuthTypeOAuth is the OAuth authentication
	AuthTypeOAuth
	// AuthTypeExternalBrowser is to use a browser to access an Fed and perform SSO authentication
	AuthTypeExternalBrowser
	// AuthTypeOkta is to use a native okta URL to perform SSO authentication on Okta
	AuthTypeOkta
	// AuthTypeJwt is to use Jwt to perform authentication
	AuthTypeJwt
	// AuthTypeTokenAccessor is to use the provided token accessor and bypass authentication
	AuthTypeTokenAccessor
	// AuthTypeUsernamePasswordMFA is to use username and password with mfa
	AuthTypeUsernamePasswordMFA
)

func determineAuthenticatorType(cfg *Config, value string) error {
	upperCaseValue := strings.ToUpper(value)
	lowerCaseValue := strings.ToLower(value)
	if strings.Trim(value, " ") == "" || upperCaseValue == AuthTypeSnowflake.String() {
		cfg.Authenticator = AuthTypeSnowflake
		return nil
	} else if upperCaseValue == AuthTypeOAuth.String() {
		cfg.Authenticator = AuthTypeOAuth
		return nil
	} else if upperCaseValue == AuthTypeJwt.String() {
		cfg.Authenticator = AuthTypeJwt
		return nil
	} else if upperCaseValue == AuthTypeExternalBrowser.String() {
		cfg.Authenticator = AuthTypeExternalBrowser
		return nil
	} else if upperCaseValue == AuthTypeUsernamePasswordMFA.String() {
		cfg.Authenticator = AuthTypeUsernamePasswordMFA
		return nil
	} else if upperCaseValue == AuthTypeTokenAccessor.String() {
		cfg.Authenticator = AuthTypeTokenAccessor
		return nil
	} else {
		// possibly Okta case
		oktaURLString, err := url.QueryUnescape(lowerCaseValue)
		if err != nil {
			return &SnowflakeError{
				Number:      ErrCodeFailedToParseAuthenticator,
				Message:     errMsgFailedToParseAuthenticator,
				MessageArgs: []interface{}{lowerCaseValue},
			}
		}

		oktaURL, err := url.Parse(oktaURLString)
		if err != nil {
			return &SnowflakeError{
				Number:      ErrCodeFailedToParseAuthenticator,
				Message:     errMsgFailedToParseAuthenticator,
				MessageArgs: []interface{}{oktaURLString},
			}
		}

		if oktaURL.Scheme != "https" {
			return &SnowflakeError{
				Number:      ErrCodeFailedToParseAuthenticator,
				Message:     errMsgFailedToParseAuthenticator,
				MessageArgs: []interface{}{oktaURLString},
			}
		}
		cfg.OktaURL = oktaURL
		cfg.Authenticator = AuthTypeOkta
	}
	return nil
}

func (authType AuthType) String() string {
	switch authType {
	case AuthTypeSnowflake:
		return "SNOWFLAKE"
	case AuthTypeOAuth:
		return "OAUTH"
	case AuthTypeExternalBrowser:
		return "EXTERNALBROWSER"
	case AuthTypeOkta:
		return "OKTA"
	case AuthTypeJwt:
		return "SNOWFLAKE_JWT"
	case AuthTypeTokenAccessor:
		return "TOKENACCESSOR"
	case AuthTypeUsernamePasswordMFA:
		return "USERNAME_PASSWORD_MFA"
	default:
		return "UNKNOWN"
	}
}

// platform consists of compiler and architecture type in string
var platform = fmt.Sprintf("%v-%v", runtime.Compiler, runtime.GOARCH)

// operatingSystem is the runtime operating system.
var operatingSystem = runtime.GOOS

// userAgent shows up in User-Agent HTTP header
var userAgent = fmt.Sprintf("%v/%v (%v-%v) %v/%v",
	clientType,
	SnowflakeGoDriverVersion,
	operatingSystem,
	runtime.GOARCH,
	runtime.Compiler,
	runtime.Version())

type authRequestClientEnvironment struct {
	Application string `json:"APPLICATION"`
	Os          string `json:"OS"`
	OsVersion   string `json:"OS_VERSION"`
	OCSPMode    string `json:"OCSP_MODE"`
	GoVersion   string `json:"GO_VERSION"`
}

type authRequestData struct {
	ClientAppID             string                       `json:"CLIENT_APP_ID"`
	ClientAppVersion        string                       `json:"CLIENT_APP_VERSION"`
	SvnRevision             string                       `json:"SVN_REVISION"`
	AccountName             string                       `json:"ACCOUNT_NAME"`
	LoginName               string                       `json:"LOGIN_NAME,omitempty"`
	Password                string                       `json:"PASSWORD,omitempty"`
	RawSAMLResponse         string                       `json:"RAW_SAML_RESPONSE,omitempty"`
	ExtAuthnDuoMethod       string                       `json:"EXT_AUTHN_DUO_METHOD,omitempty"`
	Passcode                string                       `json:"PASSCODE,omitempty"`
	Authenticator           string                       `json:"AUTHENTICATOR,omitempty"`
	SessionParameters       map[string]interface{}       `json:"SESSION_PARAMETERS,omitempty"`
	ClientEnvironment       authRequestClientEnvironment `json:"CLIENT_ENVIRONMENT"`
	BrowserModeRedirectPort string                       `json:"BROWSER_MODE_REDIRECT_PORT,omitempty"`
	ProofKey                string                       `json:"PROOF_KEY,omitempty"`
	Token                   string                       `json:"TOKEN,omitempty"`
}
type authRequest struct {
	Data authRequestData `json:"data"`
}

type nameValueParameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type authResponseSessionInfo struct {
	DatabaseName  string `json:"databaseName"`
	SchemaName    string `json:"schemaName"`
	WarehouseName string `json:"warehouseName"`
	RoleName      string `json:"roleName"`
}

type authResponseMain struct {
	Token               string                  `json:"token,omitempty"`
	Validity            time.Duration           `json:"validityInSeconds,omitempty"`
	MasterToken         string                  `json:"masterToken,omitempty"`
	MasterValidity      time.Duration           `json:"masterValidityInSeconds"`
	MfaToken            string                  `json:"mfaToken,omitempty"`
	MfaTokenValidity    time.Duration           `json:"mfaTokenValidityInSeconds"`
	IDToken             string                  `json:"idToken,omitempty"`
	IDTokenValidity     time.Duration           `json:"idTokenValidityInSeconds"`
	DisplayUserName     string                  `json:"displayUserName"`
	ServerVersion       string                  `json:"serverVersion"`
	FirstLogin          bool                    `json:"firstLogin"`
	RemMeToken          string                  `json:"remMeToken"`
	RemMeValidity       time.Duration           `json:"remMeValidityInSeconds"`
	HealthCheckInterval time.Duration           `json:"healthCheckInterval"`
	NewClientForUpgrade string                  `json:"newClientForUpgrade"`
	SessionID           int64                   `json:"sessionId"`
	Parameters          []nameValueParameter    `json:"parameters"`
	SessionInfo         authResponseSessionInfo `json:"sessionInfo"`
	TokenURL            string                  `json:"tokenUrl,omitempty"`
	SSOURL              string                  `json:"ssoUrl,omitempty"`
	ProofKey            string                  `json:"proofKey,omitempty"`
}

type authResponse struct {
	Data    authResponseMain `json:"data"`
	Message string           `json:"message"`
	Code    string           `json:"code"`
	Success bool             `json:"success"`
}

func postAuth(
	ctx context.Context,
	sr *snowflakeRestful,
	client *http.Client,
	params *url.Values,
	headers map[string]string,
	bodyCreator bodyCreatorType,
	timeout time.Duration) (
	data *authResponse, err error) {
	params.Set(requestIDKey, getOrGenerateRequestIDFromContext(ctx).String())
	params.Set(requestGUIDKey, NewUUID().String())

	fullURL := sr.getFullURL(loginRequestPath, params)
	logger.WithContext(ctx).Infof("full URL: %v", fullURL)
	resp, err := sr.FuncAuthPost(ctx, client, fullURL, headers, bodyCreator, timeout, sr.MaxRetryCount)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var respd authResponse
		err = json.NewDecoder(resp.Body).Decode(&respd)
		if err != nil {
			logger.WithContext(ctx).Errorf("failed to decode JSON. err: %v", err)
			return nil, err
		}
		return &respd, nil
	}
	switch resp.StatusCode {
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		// service availability or connectivity issue. Most likely server side issue.
		return nil, &SnowflakeError{
			Number:      ErrCodeServiceUnavailable,
			SQLState:    SQLStateConnectionWasNotEstablished,
			Message:     errMsgServiceUnavailable,
			MessageArgs: []interface{}{resp.StatusCode, fullURL},
		}
	case http.StatusUnauthorized, http.StatusForbidden:
		// failed to connect to db. account name may be wrong
		return nil, &SnowflakeError{
			Number:      ErrCodeFailedToConnect,
			SQLState:    SQLStateConnectionRejected,
			Message:     errMsgFailedToConnect,
			MessageArgs: []interface{}{resp.StatusCode, fullURL},
		}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to extract HTTP response body. err: %v", err)
		return nil, err
	}
	logger.WithContext(ctx).Infof("HTTP: %v, URL: %v, Body: %v", resp.StatusCode, fullURL, b)
	logger.WithContext(ctx).Infof("Header: %v", resp.Header)
	return nil, &SnowflakeError{
		Number:      ErrFailedToAuth,
		SQLState:    SQLStateConnectionRejected,
		Message:     errMsgFailedToAuth,
		MessageArgs: []interface{}{resp.StatusCode, fullURL},
	}
}

// Generates a map of headers needed to authenticate
// with Snowflake.
func getHeaders() map[string]string {
	headers := make(map[string]string)
	headers[httpHeaderContentType] = headerContentTypeApplicationJSON
	headers[httpHeaderAccept] = headerAcceptTypeApplicationSnowflake
	headers[httpClientAppID] = clientType
	headers[httpClientAppVersion] = SnowflakeGoDriverVersion
	headers[httpHeaderUserAgent] = userAgent
	return headers
}

// Used to authenticate the user with Snowflake.
func authenticate(
	ctx context.Context,
	sc *snowflakeConn,
	samlResponse []byte,
	proofKey []byte,
) (resp *authResponseMain, err error) {
	if sc.cfg.Authenticator == AuthTypeTokenAccessor {
		logger.WithContext(ctx).Info("Bypass authentication using existing token from token accessor")
		sessionInfo := authResponseSessionInfo{
			DatabaseName:  sc.cfg.Database,
			SchemaName:    sc.cfg.Schema,
			WarehouseName: sc.cfg.Warehouse,
			RoleName:      sc.cfg.Role,
		}
		token, masterToken, sessionID := sc.cfg.TokenAccessor.GetTokens()
		return &authResponseMain{
			Token:       token,
			MasterToken: masterToken,
			SessionID:   sessionID,
			SessionInfo: sessionInfo,
		}, nil
	}

	headers := getHeaders()
	clientEnvironment := authRequestClientEnvironment{
		Application: sc.cfg.Application,
		Os:          operatingSystem,
		OsVersion:   platform,
		OCSPMode:    sc.cfg.ocspMode(),
		GoVersion:   runtime.Version(),
	}

	sessionParameters := make(map[string]interface{})
	paramsMutex.Lock()
	for k, v := range sc.cfg.Params {
		// upper casing to normalize keys
		sessionParameters[strings.ToUpper(k)] = *v
	}
	paramsMutex.Unlock()

	sessionParameters[sessionClientValidateDefaultParameters] = sc.cfg.ValidateDefaultParameters != ConfigBoolFalse
	if sc.cfg.ClientRequestMfaToken == ConfigBoolTrue {
		sessionParameters[clientRequestMfaToken] = true
	}
	if sc.cfg.ClientStoreTemporaryCredential == ConfigBoolTrue {
		sessionParameters[clientStoreTemporaryCredential] = true
	}
	bodyCreator := func() ([]byte, error) {
		return createRequestBody(sc, sessionParameters, clientEnvironment, proofKey, samlResponse)
	}

	params := &url.Values{}
	if sc.cfg.Database != "" {
		params.Add("databaseName", sc.cfg.Database)
	}
	if sc.cfg.Schema != "" {
		params.Add("schemaName", sc.cfg.Schema)
	}
	if sc.cfg.Warehouse != "" {
		params.Add("warehouse", sc.cfg.Warehouse)
	}
	if sc.cfg.Role != "" {
		params.Add("roleName", sc.cfg.Role)
	}

	logger.WithContext(ctx).WithContext(sc.ctx).Infof("PARAMS for Auth: %v, %v, %v, %v, %v, %v",
		params, sc.rest.Protocol, sc.rest.Host, sc.rest.Port, sc.rest.LoginTimeout, sc.cfg.Authenticator.String())

	respd, err := sc.rest.FuncPostAuth(ctx, sc.rest, sc.rest.getClientFor(sc.cfg.Authenticator), params, headers, bodyCreator, sc.rest.LoginTimeout)
	if err != nil {
		return nil, err
	}
	if !respd.Success {
		logger.WithContext(ctx).Errorln("Authentication FAILED")
		sc.rest.TokenAccessor.SetTokens("", "", -1)
		if sessionParameters[clientRequestMfaToken] == true {
			deleteCredential(sc, mfaToken)
		}
		if sessionParameters[clientStoreTemporaryCredential] == true {
			deleteCredential(sc, idToken)
		}
		code, err := strconv.Atoi(respd.Code)
		if err != nil {
			code = -1
			return nil, err
		}
		return nil, (&SnowflakeError{
			Number:   code,
			SQLState: SQLStateConnectionRejected,
			Message:  respd.Message,
		}).exceptionTelemetry(sc)
	}
	logger.WithContext(ctx).Info("Authentication SUCCESS")
	sc.rest.TokenAccessor.SetTokens(respd.Data.Token, respd.Data.MasterToken, respd.Data.SessionID)
	if sessionParameters[clientRequestMfaToken] == true {
		token := respd.Data.MfaToken
		setCredential(sc, mfaToken, token)
	}
	if sessionParameters[clientStoreTemporaryCredential] == true {
		token := respd.Data.IDToken
		setCredential(sc, idToken, token)
	}
	return &respd.Data, nil
}

func createRequestBody(sc *snowflakeConn, sessionParameters map[string]interface{},
	clientEnvironment authRequestClientEnvironment, proofKey []byte, samlResponse []byte,
) ([]byte, error) {
	requestMain := authRequestData{
		ClientAppID:       clientType,
		ClientAppVersion:  SnowflakeGoDriverVersion,
		AccountName:       sc.cfg.Account,
		SessionParameters: sessionParameters,
		ClientEnvironment: clientEnvironment,
	}

	switch sc.cfg.Authenticator {
	case AuthTypeExternalBrowser:
		if sc.cfg.IDToken != "" {
			requestMain.Authenticator = idTokenAuthenticator
			requestMain.Token = sc.cfg.IDToken
			requestMain.LoginName = sc.cfg.User
		} else {
			requestMain.ProofKey = string(proofKey)
			requestMain.Token = string(samlResponse)
			requestMain.LoginName = sc.cfg.User
			requestMain.Authenticator = AuthTypeExternalBrowser.String()
		}
	case AuthTypeOAuth:
		requestMain.LoginName = sc.cfg.User
		requestMain.Authenticator = AuthTypeOAuth.String()
		requestMain.Token = sc.cfg.Token
	case AuthTypeOkta:
		samlResponse, err := authenticateBySAML(
			sc.ctx,
			sc.rest,
			sc.cfg.OktaURL,
			sc.cfg.Application,
			sc.cfg.Account,
			sc.cfg.User,
			sc.cfg.Password,
			sc.cfg.DisableSamlURLCheck)
		if err != nil {
			return nil, err
		}
		requestMain.RawSAMLResponse = string(samlResponse)
	case AuthTypeJwt:
		requestMain.Authenticator = AuthTypeJwt.String()

		jwtTokenString, err := prepareJWTToken(sc.cfg)
		if err != nil {
			return nil, err
		}
		requestMain.Token = jwtTokenString
	case AuthTypeSnowflake:
		logger.WithContext(sc.ctx).Info("Username and password")
		requestMain.LoginName = sc.cfg.User
		requestMain.Password = sc.cfg.Password
		switch {
		case sc.cfg.PasscodeInPassword:
			requestMain.ExtAuthnDuoMethod = "passcode"
		case sc.cfg.Passcode != "":
			requestMain.Passcode = sc.cfg.Passcode
			requestMain.ExtAuthnDuoMethod = "passcode"
		}
	case AuthTypeUsernamePasswordMFA:
		logger.WithContext(sc.ctx).Info("Username and password MFA")
		requestMain.LoginName = sc.cfg.User
		requestMain.Password = sc.cfg.Password
		switch {
		case sc.cfg.MfaToken != "":
			requestMain.Token = sc.cfg.MfaToken
		case sc.cfg.PasscodeInPassword:
			requestMain.ExtAuthnDuoMethod = "passcode"
		case sc.cfg.Passcode != "":
			requestMain.Passcode = sc.cfg.Passcode
			requestMain.ExtAuthnDuoMethod = "passcode"
		}
	}

	authRequest := authRequest{
		Data: requestMain,
	}
	jsonBody, err := json.Marshal(authRequest)
	if err != nil {
		return nil, err
	}
	return jsonBody, nil
}

// Generate a JWT token in string given the configuration
func prepareJWTToken(config *Config) (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(config.PrivateKey.Public())
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(pubBytes)

	accountName := extractAccountName(config.Account)
	userName := strings.ToUpper(config.User)

	issueAtTime := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": fmt.Sprintf("%s.%s.%s", accountName, userName, "SHA256:"+base64.StdEncoding.EncodeToString(hash[:])),
		"sub": fmt.Sprintf("%s.%s", accountName, userName),
		"iat": issueAtTime.Unix(),
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp": issueAtTime.Add(config.JWTExpireTimeout).Unix(),
	})

	tokenString, err := token.SignedString(config.PrivateKey)

	if err != nil {
		return "", err
	}

	return tokenString, err
}

// Authenticate with sc.cfg
func authenticateWithConfig(sc *snowflakeConn) error {
	var authData *authResponseMain
	var samlResponse []byte
	var proofKey []byte
	var err error
	//var consentCacheIdToken = true

	if sc.cfg.Authenticator == AuthTypeExternalBrowser {
		if (runtime.GOOS == "windows" || runtime.GOOS == "darwin") && sc.cfg.ClientStoreTemporaryCredential == configBoolNotSet {
			sc.cfg.ClientStoreTemporaryCredential = ConfigBoolTrue
		}
		if sc.cfg.ClientStoreTemporaryCredential == ConfigBoolTrue {
			fillCachedIDToken(sc)
		}
		// Disable console login by default
		if sc.cfg.DisableConsoleLogin == configBoolNotSet {
			sc.cfg.DisableConsoleLogin = ConfigBoolTrue
		}
	}

	if sc.cfg.Authenticator == AuthTypeUsernamePasswordMFA {
		if (runtime.GOOS == "windows" || runtime.GOOS == "darwin") && sc.cfg.ClientRequestMfaToken == configBoolNotSet {
			sc.cfg.ClientRequestMfaToken = ConfigBoolTrue
		}
		if sc.cfg.ClientRequestMfaToken == ConfigBoolTrue {
			fillCachedMfaToken(sc)
		}
	}

	logger.WithContext(sc.ctx).Infof("Authenticating via %v", sc.cfg.Authenticator.String())
	switch sc.cfg.Authenticator {
	case AuthTypeExternalBrowser:
		if sc.cfg.IDToken == "" {
			samlResponse, proofKey, err = authenticateByExternalBrowser(
				sc.ctx,
				sc.rest,
				sc.cfg.Authenticator.String(),
				sc.cfg.Application,
				sc.cfg.Account,
				sc.cfg.User,
				sc.cfg.Password,
				sc.cfg.ExternalBrowserTimeout,
				sc.cfg.DisableConsoleLogin)
			if err != nil {
				sc.cleanup()
				return err
			}
		}
	}
	authData, err = authenticate(
		sc.ctx,
		sc,
		samlResponse,
		proofKey)
	if err != nil {
		sc.cleanup()
		return err
	}
	sc.populateSessionParameters(authData.Parameters)
	sc.ctx = context.WithValue(sc.ctx, SFSessionIDKey, authData.SessionID)
	return nil
}

func fillCachedIDToken(sc *snowflakeConn) {
	getCredential(sc, idToken)
}

func fillCachedMfaToken(sc *snowflakeConn) {
	getCredential(sc, mfaToken)
}
