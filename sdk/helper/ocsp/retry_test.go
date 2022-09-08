// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package ocsp

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func fakeRequestFunc(_, _ string, _ io.Reader) (*http.Request, error) {
	return nil, nil
}

type fakeHTTPError struct {
	err     string
	timeout bool
}

func (e *fakeHTTPError) Error() string   { return e.err }
func (e *fakeHTTPError) Timeout() bool   { return e.timeout }
func (e *fakeHTTPError) Temporary() bool { return true }

type fakeResponseBody struct {
	body []byte
	cnt  int
}

func (b *fakeResponseBody) Read(p []byte) (n int, err error) {
	if b.cnt == 0 {
		copy(p, b.body)
		b.cnt = 1
		return len(b.body), nil
	}
	b.cnt = 0
	return 0, io.EOF
}

func (b *fakeResponseBody) Close() error {
	return nil
}

type fakeHTTPClient struct {
	cnt     int    // number of retry
	success bool   // return success after retry in cnt times
	timeout bool   // timeout
	body    []byte // return body
	logger  hclog.Logger
}

func (c *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.cnt--
	if c.cnt < 0 {
		c.cnt = 0
	}
	c.logger.Info("fakeHTTPClient", "cnt", c.cnt)

	var retcode int
	if c.success && c.cnt == 0 {
		retcode = 200
	} else {
		if c.timeout {
			// simulate timeout
			time.Sleep(time.Second * 1)
			return nil, &fakeHTTPError{
				err:     "Whatever reason (Client.Timeout exceeded while awaiting headers)",
				timeout: true,
			}
		}
		retcode = 0
	}

	ret := &http.Response{
		StatusCode: retcode,
		Body:       &fakeResponseBody{body: c.body},
	}
	return ret, nil
}

func TestRequestGUID(t *testing.T) {
	var ridReplacer requestGUIDReplacer
	var testURL *url.URL
	var actualURL *url.URL
	retryTime := 4

	// empty url
	testURL = &url.URL{}
	ridReplacer = newRequestGUIDReplace(testURL)
	for i := 0; i < retryTime; i++ {
		actualURL = ridReplacer.replace()
		if actualURL.String() != "" {
			t.Fatalf("empty url not replaced by an empty one, got %s", actualURL)
		}
	}

	// url with on retry id
	testURL = &url.URL{
		Path: "/" + requestIDKey + "=123-1923-9?param2=value",
	}
	ridReplacer = newRequestGUIDReplace(testURL)
	for i := 0; i < retryTime; i++ {
		actualURL = ridReplacer.replace()

		if actualURL != testURL {
			t.Fatalf("url without retry id not replaced by origin one, got %s", actualURL)
		}
	}

	// url with retry id
	// With both prefix and suffix
	prefix := "/" + requestIDKey + "=123-1923-9?" + requestGUIDKey + "="
	suffix := "?param2=value"
	testURL = &url.URL{
		Path: prefix + "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" + suffix,
	}
	ridReplacer = newRequestGUIDReplace(testURL)
	for i := 0; i < retryTime; i++ {
		actualURL = ridReplacer.replace()
		if (!strings.HasPrefix(actualURL.Path, prefix)) ||
			(!strings.HasSuffix(actualURL.Path, suffix)) ||
			len(testURL.Path) != len(actualURL.Path) {
			t.Fatalf("Retry url not replaced correctedly: \n origin: %s \n result: %s", testURL, actualURL)
		}
	}

	// With no suffix
	prefix = "/" + requestIDKey + "=123-1923-9?" + requestGUIDKey + "="
	suffix = ""
	testURL = &url.URL{
		Path: prefix + "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" + suffix,
	}
	ridReplacer = newRequestGUIDReplace(testURL)
	for i := 0; i < retryTime; i++ {
		actualURL = ridReplacer.replace()
		if (!strings.HasPrefix(actualURL.Path, prefix)) ||
			(!strings.HasSuffix(actualURL.Path, suffix)) ||
			len(testURL.Path) != len(actualURL.Path) {
			t.Fatalf("Retry url not replaced correctedly: \n origin: %s \n result: %s", testURL, actualURL)
		}

	}
	// With no prefix
	prefix = requestGUIDKey + "="
	suffix = "?param2=value"
	testURL = &url.URL{
		Path: prefix + "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" + suffix,
	}
	ridReplacer = newRequestGUIDReplace(testURL)
	for i := 0; i < retryTime; i++ {
		actualURL = ridReplacer.replace()
		if (!strings.HasPrefix(actualURL.Path, prefix)) ||
			(!strings.HasSuffix(actualURL.Path, suffix)) ||
			len(testURL.Path) != len(actualURL.Path) {
			t.Fatalf("Retry url not replaced correctedly: \n origin: %s \n result: %s", testURL, actualURL)
		}
	}
}

func TestRetryQuerySuccess(t *testing.T) {
	c := New(testLogFactory)
	c.Logger().Info("Retry N times and Success")
	client := &fakeHTTPClient{
		cnt:     3,
		success: true,
	}
	urlPtr, err := url.Parse("https://fakeaccountretrysuccess.snowflakecomputing.com:443/queries/v1/query-request?" + requestIDKey + "=testid&clientStartTime=123456")
	if err != nil {
		t.Fatal("failed to parse the test URL")
	}
	_, err = newRetryHTTP(context.TODO(),
		client,
		fakeRequestFunc, urlPtr, make(map[string]string), 60*time.Second).doPost().setBody([]byte{0}).execute()
	if err != nil {
		t.Fatal("failed to run retry")
	}
	var values url.Values
	values, err = url.ParseQuery(urlPtr.RawQuery)
	if err != nil {
		t.Fatal("failed to fail to parse the URL")
	}
	retry, err := strconv.Atoi(values.Get(retryCounterKey))
	if err != nil {
		t.Fatalf("failed to get retry counter: %v", err)
	}
	if retry < 2 {
		t.Fatalf("not enough retry counter: %v", retry)
	}
}
func TestRetryQueryFail(t *testing.T) {
	c := New(testLogFactory)
	c.Logger().Info("Retry N times and Fail")
	client := &fakeHTTPClient{
		cnt:     4,
		success: false,
	}
	urlPtr, err := url.Parse("https://fakeaccountretryfail.snowflakecomputing.com:443/queries/v1/query-request?" + requestIDKey + "=testid&clientStartTime=123456")
	if err != nil {
		t.Fatal("failed to parse the test URL")
	}
	_, err = newRetryHTTP(context.TODO(),
		client,
		fakeRequestFunc, urlPtr, make(map[string]string), 60*time.Second).doPost().setBody([]byte{0}).execute()
	if err == nil {
		t.Fatal("should fail to run retry")
	}
	var values url.Values
	values, err = url.ParseQuery(urlPtr.RawQuery)
	if err != nil {
		t.Fatalf("failed to fail to parse the URL: %v", err)
	}
	retry, err := strconv.Atoi(values.Get(retryCounterKey))
	if err != nil {
		t.Fatalf("failed to get retry counter: %v", err)
	}
	if retry < 2 {
		t.Fatalf("not enough retry counter: %v", retry)
	}
}
func TestRetryLoginRequest(t *testing.T) {
	client := &fakeHTTPClient{
		cnt:     3,
		success: true,
		timeout: true,
		logger:  hclog.New(hclog.DefaultOptions),
	}
	client.logger.Info("Retry N times for timeouts and Success")
	urlPtr, err := url.Parse("https://fakeaccountretrylogin.snowflakecomputing.com:443/login-request?request_id=testid")
	if err != nil {
		t.Fatal("failed to parse the test URL")
	}
	_, err = newRetryHTTP(context.TODO(),
		client,
		fakeRequestFunc, urlPtr, make(map[string]string), 60*time.Second).doPost().setBody([]byte{0}).execute()
	if err != nil {
		t.Fatal("failed to run retry")
	}
	var values url.Values
	values, err = url.ParseQuery(urlPtr.RawQuery)
	if err != nil {
		t.Fatalf("failed to fail to parse the URL: %v", err)
	}
	if values.Get(retryCounterKey) != "" {
		t.Fatalf("no retry counter should be attached: %v", retryCounterKey)
	}
	client.logger.Info("Retry N times for timeouts and Fail")
	client = &fakeHTTPClient{
		cnt:     10,
		success: false,
		timeout: true,
		logger:  hclog.New(hclog.DefaultOptions),
	}
	_, err = newRetryHTTP(context.TODO(),
		client,
		fakeRequestFunc, urlPtr, make(map[string]string), 10*time.Second).doPost().setBody([]byte{0}).execute()
	if err == nil {
		t.Fatal("should fail to run retry")
	}
	values, err = url.ParseQuery(urlPtr.RawQuery)
	if err != nil {
		t.Fatalf("failed to fail to parse the URL: %v", err)
	}
	if values.Get(retryCounterKey) != "" {
		t.Fatalf("no retry counter should be attached: %v", retryCounterKey)
	}
}
