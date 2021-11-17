// Copyright (c) 2017-2019 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"context"

	"sync"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// requestGUIDKey is attached to every request against Snowflake
const requestGUIDKey string = "request_guid"

// retryCounterKey is attached to query-request from the second time
const retryCounterKey string = "retryCounter"

// requestIDKey is attached to all requests to Snowflake
const requestIDKey string = "requestId"

// This class takes in an url during construction and replace the
// value of request_guid every time the replace() is called
// When the url does not contain request_guid, just return the original
// url
type requestGUIDReplacer interface {
	// replace the url with new ID
	replace() *url.URL
}

// Make requestGUIDReplacer given a url string
func newRequestGUIDReplace(urlPtr *url.URL) requestGUIDReplacer {
	values, err := url.ParseQuery(urlPtr.RawQuery)
	if err != nil {
		// nop if invalid query parameters
		return &transientReplace{urlPtr}
	}
	if len(values.Get(requestGUIDKey)) == 0 {
		// nop if no request_guid is included.
		return &transientReplace{urlPtr}
	}

	return &requestGUIDReplace{urlPtr, values}
}

// this replacer does nothing but replace the url
type transientReplace struct {
	urlPtr *url.URL
}

func (replacer *transientReplace) replace() *url.URL {
	return replacer.urlPtr
}

/*
requestGUIDReplacer is a one-shot object that is created out of the retry loop and
called with replace to change the retry_guid's value upon every retry
*/
type requestGUIDReplace struct {
	urlPtr    *url.URL
	urlValues url.Values
}

/**
This function would replace they value of the requestGUIDKey in a url with a newly
generated uuid
*/
func (replacer *requestGUIDReplace) replace() *url.URL {
	replacer.urlValues.Del(requestGUIDKey)
	replacer.urlValues.Add(requestGUIDKey, uuid.New().String())
	replacer.urlPtr.RawQuery = replacer.urlValues.Encode()
	return replacer.urlPtr
}

type retryCounterUpdater interface {
	replaceOrAdd(retry int) *url.URL
}

type retryCounterUpdate struct {
	urlPtr    *url.URL
	urlValues url.Values
}

// this replacer does nothing but replace the url
type transientReplaceOrAdd struct {
	urlPtr *url.URL
}

func (replaceOrAdder *transientReplaceOrAdd) replaceOrAdd(retry int) *url.URL {
	return replaceOrAdder.urlPtr
}

func (replacer *retryCounterUpdate) replaceOrAdd(retry int) *url.URL {
	replacer.urlValues.Del(retryCounterKey)
	replacer.urlValues.Add(retryCounterKey, strconv.Itoa(retry))
	replacer.urlPtr.RawQuery = replacer.urlValues.Encode()
	return replacer.urlPtr
}

func newRetryUpdate(urlPtr *url.URL) retryCounterUpdater {
	if !strings.HasPrefix(urlPtr.Path, queryRequestPath) {
		// nop if not query-request
		return &transientReplaceOrAdd{urlPtr}
	}
	values, err := url.ParseQuery(urlPtr.RawQuery)
	if err != nil {
		// nop if the URL is not valid
		return &transientReplaceOrAdd{urlPtr}
	}
	return &retryCounterUpdate{urlPtr, values}
}

type waitAlgo struct {
	mutex *sync.Mutex   // required for random.Int63n
	base  time.Duration // base wait time
	cap   time.Duration // maximum wait time
}

func randSecondDuration(n time.Duration) time.Duration {
	return time.Duration(random.Int63n(int64(n/time.Second))) * time.Second
}

// decorrelated jitter backoff
func (w *waitAlgo) decorr(attempt int, sleep time.Duration) time.Duration {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	t := 3*sleep - w.base
	switch {
	case t > 0:
		return durationMin(w.cap, randSecondDuration(t)+w.base)
	case t < 0:
		return durationMin(w.cap, randSecondDuration(-t)+3*sleep)
	}
	return w.base
}

var defaultWaitAlgo = &waitAlgo{
	mutex: &sync.Mutex{},
	base:  5 * time.Second,
	cap:   160 * time.Second,
}

type requestFunc func(method, urlStr string, body io.Reader) (*http.Request, error)

type clientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type retryHTTP struct {
	ctx      context.Context
	client   clientInterface
	req      requestFunc
	method   string
	fullURL  *url.URL
	headers  map[string]string
	body     []byte
	timeout  time.Duration
	raise4XX bool
}

func newRetryHTTP(ctx context.Context,
	client clientInterface,
	req requestFunc,
	fullURL *url.URL,
	headers map[string]string,
	timeout time.Duration) *retryHTTP {
	instance := retryHTTP{}
	instance.ctx = ctx
	instance.client = client
	instance.req = req
	instance.method = "GET"
	instance.fullURL = fullURL
	instance.headers = headers
	instance.body = nil
	instance.timeout = timeout
	instance.raise4XX = false
	return &instance
}

func (r *retryHTTP) doRaise4XX(raise4XX bool) *retryHTTP {
	r.raise4XX = raise4XX
	return r
}

func (r *retryHTTP) doPost() *retryHTTP {
	r.method = "POST"
	return r
}

func (r *retryHTTP) setBody(body []byte) *retryHTTP {
	r.body = body
	return r
}

func (r *retryHTTP) execute() (res *http.Response, err error) {
	totalTimeout := r.timeout
	logger.WithContext(r.ctx).Infof("retryHTTP.totalTimeout: %v", totalTimeout)
	retryCounter := 0
	sleepTime := time.Duration(0)

	var rIDReplacer requestGUIDReplacer
	var rUpdater retryCounterUpdater

	for {
		req, err := r.req(r.method, r.fullURL.String(), bytes.NewReader(r.body))
		if err != nil {
			return nil, err
		}
		if req != nil {
			// req can be nil in tests
			req = req.WithContext(r.ctx)
		}
		for k, v := range r.headers {
			req.Header.Set(k, v)
		}
		res, err = r.client.Do(req)
		if err != nil {
			// check if it can retry.
			doExit, err := r.isRetryableError(err)
			if doExit {
				return res, err
			}
			// cannot just return 4xx and 5xx status as the error can be sporadic. run often helps.
			logger.WithContext(r.ctx).Warningf(
				"failed http connection. no response is returned. err: %v. retrying...\n", err)
		} else {
			if res.StatusCode == http.StatusOK || r.raise4XX && res != nil && res.StatusCode >= 400 && res.StatusCode < 500 {
				// exit if success
				// or
				// abort connection if raise4XX flag is enabled and the range of HTTP status code are 4XX.
				// This is currently used for Snowflake login. The caller must generate an error object based on HTTP status.
				break
			}
			logger.WithContext(r.ctx).Warningf(
				"failed http connection. HTTP Status: %v. retrying...\n", res.StatusCode)
			res.Body.Close()
		}
		// uses decorrelated jitter backoff
		sleepTime = defaultWaitAlgo.decorr(retryCounter, sleepTime)

		if totalTimeout > 0 {
			logger.WithContext(r.ctx).Infof("to timeout: %v", totalTimeout)
			// if any timeout is set
			totalTimeout -= sleepTime
			if totalTimeout <= 0 {
				if err != nil {
					return nil, err
				}
				if res != nil {
					return nil, fmt.Errorf("timeout after %s. HTTP Status: %v. Hanging?", r.timeout, res.StatusCode)
				}
				return nil, fmt.Errorf("timeout after %s. Hanging?", r.timeout)
			}
		}
		retryCounter++
		if rIDReplacer == nil {
			rIDReplacer = newRequestGUIDReplace(r.fullURL)
		}
		r.fullURL = rIDReplacer.replace()
		if rUpdater == nil {
			rUpdater = newRetryUpdate(r.fullURL)
		}
		r.fullURL = rUpdater.replaceOrAdd(retryCounter)
		logger.WithContext(r.ctx).Infof("sleeping %v. to timeout: %v. retrying", sleepTime, totalTimeout)

		await := time.NewTimer(sleepTime)
		select {
		case <-await.C:
			// retry the request
		case <-r.ctx.Done():
			await.Stop()
			return res, r.ctx.Err()
		}
	}
	return res, err
}

func (r *retryHTTP) isRetryableError(err error) (bool, error) {
	urlError, isURLError := err.(*url.Error)
	if isURLError {
		// context cancel or timeout
		if urlError.Err == context.DeadlineExceeded || urlError.Err == context.Canceled {
			return true, urlError.Err
		}
		if driverError, ok := urlError.Err.(*SnowflakeError); ok {
			// Certificate Revoked
			if driverError.Number == ErrOCSPStatusRevoked {
				return true, err
			}
		}
		if _, ok := urlError.Err.(x509.CertificateInvalidError); ok {
			// Certificate is invalid
			return true, err
		}
		if _, ok := urlError.Err.(x509.UnknownAuthorityError); ok {
			// Certificate is self-signed
			return true, err
		}

	}
	return false, err
}
