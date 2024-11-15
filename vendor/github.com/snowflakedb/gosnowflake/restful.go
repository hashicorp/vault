// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// HTTP headers
const (
	headerSnowflakeToken   = "Snowflake Token=\"%v\""
	headerAuthorizationKey = "Authorization"

	headerContentTypeApplicationJSON     = "application/json"
	headerAcceptTypeApplicationSnowflake = "application/snowflake"
)

// Snowflake Server Error code
const (
	queryInProgressCode      = "333333"
	queryInProgressAsyncCode = "333334"
	sessionExpiredCode       = "390112"
	queryNotExecuting        = "000605"
)

// Snowflake Server Endpoints
const (
	loginRequestPath         = "/session/v1/login-request"
	queryRequestPath         = "/queries/v1/query-request"
	tokenRequestPath         = "/session/token-request"
	abortRequestPath         = "/queries/v1/abort-request"
	authenticatorRequestPath = "/session/authenticator-request"
	monitoringQueriesPath    = "/monitoring/queries"
	sessionRequestPath       = "/session"
	heartBeatPath            = "/session/heartbeat"
	consoleLoginRequestPath  = "/console/login"
)

type (
	funcGetType      func(context.Context, *snowflakeRestful, *url.URL, map[string]string, time.Duration) (*http.Response, error)
	funcPostType     func(context.Context, *snowflakeRestful, *url.URL, map[string]string, []byte, time.Duration, currentTimeProvider, *Config) (*http.Response, error)
	funcAuthPostType func(context.Context, *http.Client, *url.URL, map[string]string, bodyCreatorType, time.Duration, int) (*http.Response, error)
	bodyCreatorType  func() ([]byte, error)
)

var emptyBodyCreator = func() ([]byte, error) {
	return []byte{}, nil
}

type snowflakeRestful struct {
	Host           string
	Port           int
	Protocol       string
	LoginTimeout   time.Duration // Login timeout
	RequestTimeout time.Duration // request timeout
	MaxRetryCount  int

	Client        *http.Client
	JWTClient     *http.Client
	TokenAccessor TokenAccessor
	HeartBeat     *heartbeat

	Connection *snowflakeConn

	FuncPostQuery       func(context.Context, *snowflakeRestful, *url.Values, map[string]string, []byte, time.Duration, UUID, *Config) (*execResponse, error)
	FuncPostQueryHelper func(context.Context, *snowflakeRestful, *url.Values, map[string]string, []byte, time.Duration, UUID, *Config) (*execResponse, error)
	FuncPost            funcPostType
	FuncGet             funcGetType
	FuncAuthPost        funcAuthPostType
	FuncRenewSession    func(context.Context, *snowflakeRestful, time.Duration) error
	FuncCloseSession    func(context.Context, *snowflakeRestful, time.Duration) error
	FuncCancelQuery     func(context.Context, *snowflakeRestful, UUID, time.Duration) error

	FuncPostAuth     func(context.Context, *snowflakeRestful, *http.Client, *url.Values, map[string]string, bodyCreatorType, time.Duration) (*authResponse, error)
	FuncPostAuthSAML func(context.Context, *snowflakeRestful, map[string]string, []byte, time.Duration) (*authResponse, error)
	FuncPostAuthOKTA func(context.Context, *snowflakeRestful, map[string]string, []byte, string, time.Duration) (*authOKTAResponse, error)
	FuncGetSSO       func(context.Context, *snowflakeRestful, *url.Values, map[string]string, string, time.Duration) ([]byte, error)
}

func (sr *snowflakeRestful) getURL() *url.URL {
	return &url.URL{
		Scheme: sr.Protocol,
		Host:   sr.Host + ":" + strconv.Itoa(sr.Port),
	}
}

func (sr *snowflakeRestful) getFullURL(path string, params *url.Values) *url.URL {
	ret := &url.URL{
		Scheme: sr.Protocol,
		Host:   sr.Host + ":" + strconv.Itoa(sr.Port),
		Path:   path,
	}
	if params != nil {
		ret.RawQuery = params.Encode()
	}
	return ret
}

// We need separate client for JWT, because if token processing takes too long, token may be already expired.
func (sr *snowflakeRestful) getClientFor(authType AuthType) *http.Client {
	switch authType {
	case AuthTypeJwt:
		return sr.JWTClient
	default:
		return sr.Client
	}
}

// Renew the snowflake session if the current token is still the stale token specified
func (sr *snowflakeRestful) renewExpiredSessionToken(ctx context.Context, timeout time.Duration, expiredToken string) error {
	err := sr.TokenAccessor.Lock()
	if err != nil {
		return err
	}
	defer sr.TokenAccessor.Unlock()
	currentToken, _, _ := sr.TokenAccessor.GetTokens()
	if expiredToken == currentToken || currentToken == "" {
		// Only renew the session if the current token is still the expired token or current token is empty
		return sr.FuncRenewSession(ctx, sr, timeout)
	}
	return nil
}

type renewSessionResponse struct {
	Data    renewSessionResponseMain `json:"data"`
	Message string                   `json:"message"`
	Code    string                   `json:"code"`
	Success bool                     `json:"success"`
}

type renewSessionResponseMain struct {
	SessionToken        string        `json:"sessionToken"`
	ValidityInSecondsST time.Duration `json:"validityInSecondsST"`
	MasterToken         string        `json:"masterToken"`
	ValidityInSecondsMT time.Duration `json:"validityInSecondsMT"`
	SessionID           int64         `json:"sessionId"`
}

type cancelQueryResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Success bool        `json:"success"`
}

type telemetryResponse struct {
	Data    interface{}       `json:"data,omitempty"`
	Message string            `json:"message"`
	Code    string            `json:"code"`
	Success bool              `json:"success"`
	Headers map[string]string `json:"headers,omitempty"`
}

func postRestful(
	ctx context.Context,
	sr *snowflakeRestful,
	fullURL *url.URL,
	headers map[string]string,
	body []byte,
	timeout time.Duration,
	currentTimeProvider currentTimeProvider,
	cfg *Config) (
	*http.Response, error) {
	return newRetryHTTP(ctx, sr.Client, http.NewRequest, fullURL, headers, timeout, sr.MaxRetryCount, currentTimeProvider, cfg).
		doPost().
		setBody(body).
		execute()
}

func getRestful(
	ctx context.Context,
	sr *snowflakeRestful,
	fullURL *url.URL,
	headers map[string]string,
	timeout time.Duration) (
	*http.Response, error) {
	return newRetryHTTP(ctx, sr.Client, http.NewRequest, fullURL, headers, timeout, sr.MaxRetryCount, defaultTimeProvider, nil).execute()
}

func postAuthRestful(
	ctx context.Context,
	client *http.Client,
	fullURL *url.URL,
	headers map[string]string,
	bodyCreator bodyCreatorType,
	timeout time.Duration,
	maxRetryCount int) (
	*http.Response, error) {
	return newRetryHTTP(ctx, client, http.NewRequest, fullURL, headers, timeout, maxRetryCount, defaultTimeProvider, nil).
		doPost().
		setBodyCreator(bodyCreator).
		execute()
}

func postRestfulQuery(
	ctx context.Context,
	sr *snowflakeRestful,
	params *url.Values,
	headers map[string]string,
	body []byte,
	timeout time.Duration,
	requestID UUID,
	cfg *Config) (
	data *execResponse, err error) {

	data, err = sr.FuncPostQueryHelper(ctx, sr, params, headers, body, timeout, requestID, cfg)

	// errors other than context timeout and cancel would be returned to upper layers
	if err != context.Canceled && err != context.DeadlineExceeded {
		return data, err
	}

	if err = sr.FuncCancelQuery(context.Background(), sr, requestID, timeout); err != nil {
		return nil, err
	}
	return nil, ctx.Err()
}

func postRestfulQueryHelper(
	ctx context.Context,
	sr *snowflakeRestful,
	params *url.Values,
	headers map[string]string,
	body []byte,
	timeout time.Duration,
	requestID UUID,
	cfg *Config) (
	data *execResponse, err error) {
	logger.WithContext(ctx).Infof("params: %v", params)
	params.Set(requestIDKey, requestID.String())
	params.Set(requestGUIDKey, NewUUID().String())
	token, _, _ := sr.TokenAccessor.GetTokens()
	if token != "" {
		headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, token)
	}

	var resp *http.Response
	fullURL := sr.getFullURL(queryRequestPath, params)
	resp, err = sr.FuncPost(ctx, sr, fullURL, headers, body, timeout, defaultTimeProvider, cfg)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		logger.WithContext(ctx).Infof("postQuery: resp: %v", resp)
		var respd execResponse
		if err = json.NewDecoder(resp.Body).Decode(&respd); err != nil {
			logger.WithContext(ctx).Errorf("failed to decode JSON. err: %v", err)
			return nil, err
		}
		if respd.Code == sessionExpiredCode {
			if err = sr.renewExpiredSessionToken(ctx, timeout, token); err != nil {
				return nil, err
			}
			return sr.FuncPostQuery(ctx, sr, params, headers, body, timeout, requestID, cfg)
		}

		if queryIDChan := getQueryIDChan(ctx); queryIDChan != nil {
			queryIDChan <- respd.Data.QueryID
			close(queryIDChan)
			ctx = WithQueryIDChan(ctx, nil)
		}

		isSessionRenewed := false

		// if asynchronous query in progress, kick off retrieval but return object
		if respd.Code == queryInProgressAsyncCode && isAsyncMode(ctx) {
			return sr.processAsync(ctx, &respd, headers, timeout, cfg)
		}
		for isSessionRenewed || respd.Code == queryInProgressCode ||
			respd.Code == queryInProgressAsyncCode {
			if !isSessionRenewed {
				fullURL = sr.getFullURL(respd.Data.GetResultURL, nil)
			}

			logger.WithContext(ctx).Info("ping pong")
			token, _, _ = sr.TokenAccessor.GetTokens()
			headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, token)

			resp, err = sr.FuncGet(ctx, sr, fullURL, headers, timeout)
			if err != nil {
				logger.WithContext(ctx).Errorf("failed to get response. err: %v", err)
				return nil, err
			}
			respd = execResponse{} // reset the response
			err = json.NewDecoder(resp.Body).Decode(&respd)
			resp.Body.Close()
			if err != nil {
				logger.WithContext(ctx).Errorf("failed to decode JSON. err: %v", err)
				return nil, err
			}
			if respd.Code == sessionExpiredCode {
				if err = sr.renewExpiredSessionToken(ctx, timeout, token); err != nil {
					return nil, err
				}
				isSessionRenewed = true
			} else {
				isSessionRenewed = false
			}
		}
		return &respd, nil
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to extract HTTP response body. err: %v", err)
		return nil, err
	}
	logger.WithContext(ctx).Infof("HTTP: %v, URL: %v, Body: %v", resp.StatusCode, fullURL, b)
	logger.WithContext(ctx).Infof("Header: %v", resp.Header)
	return nil, &SnowflakeError{
		Number:      ErrFailedToPostQuery,
		SQLState:    SQLStateConnectionFailure,
		Message:     errMsgFailedToPostQuery,
		MessageArgs: []interface{}{resp.StatusCode, fullURL},
	}
}

func closeSession(ctx context.Context, sr *snowflakeRestful, timeout time.Duration) error {
	logger.WithContext(ctx).Info("close session")
	params := &url.Values{}
	params.Set("delete", "true")
	params.Set(requestIDKey, getOrGenerateRequestIDFromContext(ctx).String())
	params.Set(requestGUIDKey, NewUUID().String())
	fullURL := sr.getFullURL(sessionRequestPath, params)

	headers := getHeaders()
	token, _, _ := sr.TokenAccessor.GetTokens()
	headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, token)

	resp, err := sr.FuncPost(ctx, sr, fullURL, headers, nil, 5*time.Second, defaultTimeProvider, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var respd renewSessionResponse
		if err = json.NewDecoder(resp.Body).Decode(&respd); err != nil {
			logger.WithContext(ctx).Errorf("failed to decode JSON. err: %v", err)
			return err
		}
		if !respd.Success && respd.Code != sessionExpiredCode {
			c, err := strconv.Atoi(respd.Code)
			if err != nil {
				return err
			}
			return &SnowflakeError{
				Number:  c,
				Message: respd.Message,
			}
		}
		return nil
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to extract HTTP response body. err: %v", err)
		return err
	}
	logger.WithContext(ctx).Infof("HTTP: %v, URL: %v, Body: %v", resp.StatusCode, fullURL, b)
	logger.WithContext(ctx).Infof("Header: %v", resp.Header)
	return &SnowflakeError{
		Number:      ErrFailedToCloseSession,
		SQLState:    SQLStateConnectionFailure,
		Message:     errMsgFailedToCloseSession,
		MessageArgs: []interface{}{resp.StatusCode, fullURL},
	}
}

func renewRestfulSession(ctx context.Context, sr *snowflakeRestful, timeout time.Duration) error {
	logger.WithContext(ctx).Info("start renew session")
	params := &url.Values{}
	params.Set(requestIDKey, getOrGenerateRequestIDFromContext(ctx).String())
	params.Set(requestGUIDKey, NewUUID().String())
	fullURL := sr.getFullURL(tokenRequestPath, params)

	token, masterToken, _ := sr.TokenAccessor.GetTokens()
	headers := getHeaders()
	headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, masterToken)

	body := make(map[string]string)
	body["oldSessionToken"] = token
	body["requestType"] = "RENEW"

	var reqBody []byte
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := sr.FuncPost(ctx, sr, fullURL, headers, reqBody, timeout, defaultTimeProvider, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var respd renewSessionResponse
		err = json.NewDecoder(resp.Body).Decode(&respd)
		if err != nil {
			logger.WithContext(ctx).Errorf("failed to decode JSON. err: %v", err)
			return err
		}
		if !respd.Success {
			c, err := strconv.Atoi(respd.Code)
			if err != nil {
				return err
			}
			return &SnowflakeError{
				Number:  c,
				Message: respd.Message,
			}
		}
		sr.TokenAccessor.SetTokens(respd.Data.SessionToken, respd.Data.MasterToken, respd.Data.SessionID)
		return nil
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to extract HTTP response body. err: %v", err)
		return err
	}
	logger.WithContext(ctx).Infof("HTTP: %v, URL: %v, Body: %v", resp.StatusCode, fullURL, b)
	logger.WithContext(ctx).Infof("Header: %v", resp.Header)
	return &SnowflakeError{
		Number:      ErrFailedToRenewSession,
		SQLState:    SQLStateConnectionFailure,
		Message:     errMsgFailedToRenew,
		MessageArgs: []interface{}{resp.StatusCode, fullURL},
	}
}

func getCancelRetry(ctx context.Context) int {
	val := ctx.Value(cancelRetry)
	if val == nil {
		return 5
	}
	cnt, ok := val.(int)
	if !ok {
		return -1
	}
	return cnt
}

func cancelQuery(ctx context.Context, sr *snowflakeRestful, requestID UUID, timeout time.Duration) error {
	logger.WithContext(ctx).Info("cancel query")
	params := &url.Values{}
	params.Set(requestIDKey, getOrGenerateRequestIDFromContext(ctx).String())
	params.Set(requestGUIDKey, NewUUID().String())

	fullURL := sr.getFullURL(abortRequestPath, params)

	headers := getHeaders()
	token, _, _ := sr.TokenAccessor.GetTokens()
	headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, token)

	req := make(map[string]string)
	req[requestIDKey] = requestID.String()

	reqByte, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := sr.FuncPost(ctx, sr, fullURL, headers, reqByte, timeout, defaultTimeProvider, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var respd cancelQueryResponse
		if err = json.NewDecoder(resp.Body).Decode(&respd); err != nil {
			logger.WithContext(ctx).Errorf("failed to decode JSON. err: %v", err)
			return err
		}
		ctxRetry := getCancelRetry(ctx)
		if !respd.Success && respd.Code == sessionExpiredCode {
			if err = sr.FuncRenewSession(ctx, sr, timeout); err != nil {
				return err
			}
			return sr.FuncCancelQuery(ctx, sr, requestID, timeout)
		} else if !respd.Success && respd.Code == queryNotExecuting && ctxRetry != 0 {
			return sr.FuncCancelQuery(context.WithValue(ctx, cancelRetry, ctxRetry-1), sr, requestID, timeout)
		} else if respd.Success {
			return nil
		} else {
			c, err := strconv.Atoi(respd.Code)
			if err != nil {
				return err
			}
			return &SnowflakeError{
				Number:  c,
				Message: respd.Message,
			}
		}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to extract HTTP response body. err: %v", err)
		return err
	}
	logger.WithContext(ctx).Infof("HTTP: %v, URL: %v, Body: %v", resp.StatusCode, fullURL, b)
	logger.WithContext(ctx).Infof("Header: %v", resp.Header)
	return &SnowflakeError{
		Number:      ErrFailedToCancelQuery,
		SQLState:    SQLStateConnectionFailure,
		Message:     errMsgFailedToCancelQuery,
		MessageArgs: []interface{}{resp.StatusCode, fullURL},
	}
}

func getQueryIDChan(ctx context.Context) chan<- string {
	v := ctx.Value(queryIDChannel)
	if v == nil {
		return nil
	}
	c, ok := v.(chan<- string)
	if !ok {
		return nil
	}
	return c
}
