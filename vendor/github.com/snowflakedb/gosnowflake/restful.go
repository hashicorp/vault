// Copyright (c) 2017-2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
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
	sessionRequestPath       = "/session"
	heartBeatPath            = "/session/heartbeat"
)

// FuncGetType httpclient GET method to return http.Response
type FuncGetType func(context.Context, *snowflakeRestful, *url.URL, map[string]string, time.Duration) (*http.Response, error)

// FuncPostType httpclient POST method to return http.Response
type FuncPostType func(context.Context, *snowflakeRestful, *url.URL, map[string]string, []byte, time.Duration, bool) (*http.Response, error)

type snowflakeRestful struct {
	Host           string
	Port           int
	Protocol       string
	LoginTimeout   time.Duration // Login timeout
	RequestTimeout time.Duration // request timeout

	Client        *http.Client
	TokenAccessor TokenAccessor
	HeartBeat     *heartbeat

	Connection *snowflakeConn

	FuncPostQuery       func(context.Context, *snowflakeRestful, *url.Values, map[string]string, []byte, time.Duration, uuid.UUID, *Config) (*execResponse, error)
	FuncPostQueryHelper func(context.Context, *snowflakeRestful, *url.Values, map[string]string, []byte, time.Duration, uuid.UUID, *Config) (*execResponse, error)
	FuncPost            FuncPostType
	FuncGet             FuncGetType
	FuncRenewSession    func(context.Context, *snowflakeRestful, time.Duration) error
	FuncPostAuth        func(context.Context, *snowflakeRestful, *url.Values, map[string]string, []byte, time.Duration) (*authResponse, error)
	FuncCloseSession    func(context.Context, *snowflakeRestful, time.Duration) error
	FuncCancelQuery     func(context.Context, *snowflakeRestful, uuid.UUID, time.Duration) error

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
	raise4XX bool) (
	*http.Response, error) {
	return newRetryHTTP(
		ctx, sr.Client, http.NewRequest, fullURL, headers, timeout).doPost().setBody(body).doRaise4XX(raise4XX).execute()
}

func getRestful(
	ctx context.Context,
	sr *snowflakeRestful,
	fullURL *url.URL,
	headers map[string]string,
	timeout time.Duration) (
	*http.Response, error) {
	return newRetryHTTP(
		ctx, sr.Client, http.NewRequest, fullURL, headers, timeout).execute()
}

func postRestfulQuery(
	ctx context.Context,
	sr *snowflakeRestful,
	params *url.Values,
	headers map[string]string,
	body []byte,
	timeout time.Duration,
	requestID uuid.UUID,
	cfg *Config) (
	data *execResponse, err error) {

	data, err = sr.FuncPostQueryHelper(ctx, sr, params, headers, body, timeout, requestID, cfg)

	// errors other than context timeout and cancel would be returned to upper layers
	if err != context.Canceled && err != context.DeadlineExceeded {
		return data, err
	}

	err = sr.FuncCancelQuery(context.TODO(), sr, requestID, timeout)
	if err != nil {
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
	requestID uuid.UUID,
	cfg *Config) (
	data *execResponse, err error) {
	logger.Infof("params: %v", params)
	params.Add(requestIDKey, requestID.String())
	params.Add("clientStartTime", strconv.FormatInt(time.Now().Unix(), 10))
	params.Add(requestGUIDKey, uuid.New().String())
	token, _, _ := sr.TokenAccessor.GetTokens()
	if token != "" {
		headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, token)
	}

	var resp *http.Response
	fullURL := sr.getFullURL(queryRequestPath, params)
	resp, err = sr.FuncPost(ctx, sr, fullURL, headers, body, timeout, true)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		logger.WithContext(ctx).Infof("postQuery: resp: %v", resp)
		var respd execResponse
		err = json.NewDecoder(resp.Body).Decode(&respd)
		if err != nil {
			logger.WithContext(ctx).Errorf("failed to decode JSON. err: %v", err)
			return nil, err
		}
		if respd.Code == sessionExpiredCode {
			err = sr.renewExpiredSessionToken(ctx, timeout, token)
			if err != nil {
				return nil, err
			}
			return sr.FuncPostQuery(ctx, sr, params, headers, body, timeout, requestID, cfg)
		}

		if queryIDChan := getQueryIDChan(ctx); queryIDChan != nil {
			queryIDChan <- respd.Data.QueryID
			close(queryIDChan)
			ctx = WithQueryIDChan(ctx, nil)
		}

		var resultURL string
		isSessionRenewed := false
		noResult := isAsyncMode(ctx)

		// if asynchronous query in progress, kick off retrieval but return object
		if respd.Code == queryInProgressAsyncCode && noResult {
			// placeholder object to return to user while retrieving results
			rows := new(snowflakeRows)
			res := new(snowflakeResult)
			switch resType := getResultType(ctx); resType {
			case execResultType:
				res.queryID = respd.Data.QueryID
				res.status = QueryStatusInProgress
				res.errChannel = make(chan error)
				respd.Data.AsyncResult = res
			case queryResultType:
				rows.queryID = respd.Data.QueryID
				rows.status = QueryStatusInProgress
				rows.errChannel = make(chan error)
				respd.Data.AsyncRows = rows
			default:
				return &respd, nil
			}

			// spawn goroutine to retrieve asynchronous results
			go getAsync(ctx, sr, headers, sr.getFullURL(respd.Data.GetResultURL, nil), timeout, res, rows, cfg)
			return &respd, nil
		}
		for isSessionRenewed || respd.Code == queryInProgressCode ||
			respd.Code == queryInProgressAsyncCode {
			if !isSessionRenewed {
				resultURL = respd.Data.GetResultURL
			}

			logger.Info("ping pong")
			token, _, _ := sr.TokenAccessor.GetTokens()
			headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, token)
			fullURL := sr.getFullURL(resultURL, nil)

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
				err = sr.renewExpiredSessionToken(ctx, timeout, token)
				if err != nil {
					return nil, err
				}
				isSessionRenewed = true
			} else {
				isSessionRenewed = false
			}
		}
		return &respd, nil
	}
	b, err := ioutil.ReadAll(resp.Body)
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
	params.Add("delete", "true")
	params.Add(requestIDKey, getOrGenerateRequestIDFromContext(ctx).String())
	params.Add(requestGUIDKey, uuid.New().String())
	fullURL := sr.getFullURL(sessionRequestPath, params)

	headers := getHeaders()
	token, _, _ := sr.TokenAccessor.GetTokens()
	headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, token)

	resp, err := sr.FuncPost(ctx, sr, fullURL, headers, nil, 5*time.Second, false)
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
	b, err := ioutil.ReadAll(resp.Body)
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
	params.Add(requestIDKey, getOrGenerateRequestIDFromContext(ctx).String())
	params.Add(requestGUIDKey, uuid.New().String())
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

	resp, err := sr.FuncPost(ctx, sr, fullURL, headers, reqBody, timeout, false)
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
	b, err := ioutil.ReadAll(resp.Body)
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
	cnt, _ := val.(int)
	return cnt
}

func cancelQuery(ctx context.Context, sr *snowflakeRestful, requestID uuid.UUID, timeout time.Duration) error {
	logger.WithContext(ctx).Info("cancel query")
	params := &url.Values{}
	params.Add(requestIDKey, getOrGenerateRequestIDFromContext(ctx).String())
	params.Add(requestGUIDKey, uuid.New().String())

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

	resp, err := sr.FuncPost(ctx, sr, fullURL, headers, reqByte, timeout, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var respd cancelQueryResponse
		err = json.NewDecoder(resp.Body).Decode(&respd)
		if err != nil {
			logger.WithContext(ctx).Errorf("failed to decode JSON. err: %v", err)
			return err
		}
		ctxRetry := getCancelRetry(ctx)
		if !respd.Success && respd.Code == sessionExpiredCode {
			err := sr.FuncRenewSession(ctx, sr, timeout)
			if err != nil {
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
	b, err := ioutil.ReadAll(resp.Body)
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
	c, _ := v.(chan<- string)
	return c
}
