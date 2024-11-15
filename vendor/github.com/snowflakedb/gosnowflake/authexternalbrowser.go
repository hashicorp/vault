// Copyright (c) 2019-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/browser"
)

const (
	successHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"/>
<title>SAML Response for Snowflake</title></head>
<body>
Your identity was confirmed and propagated to Snowflake %v.
You can close this window now and go back where you started from.
</body></html>`
)

const (
	bufSize = 8192
)

// Builds a response to show to the user after successfully
// getting a response from Snowflake.
func buildResponse(application string) bytes.Buffer {
	body := fmt.Sprintf(successHTML, application)
	t := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       nil,
		Header:        make(http.Header),
	}
	var b bytes.Buffer
	t.Write(&b)
	return b
}

// This opens a socket that listens on all available unicast
// and any anycast IP addresses locally. By specifying "0", we are
// able to bind to a free port.
func createLocalTCPListener() (*net.TCPListener, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}

	tcpListener, ok := l.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("failed to assert type as *net.TCPListener")
	}

	return tcpListener, nil
}

// Opens a browser window (or new tab) with the configured login Url.
// This can / will fail if running inside a shell with no display, ie
// ssh'ing into a box attempting to authenticate via external browser.
func openBrowser(loginURL string) error {
	err := browser.OpenURL(loginURL)
	if err != nil {
		logger.Infof("failed to open a browser. err: %v", err)
		return err
	}
	return nil
}

// Gets the IDP Url and Proof Key from Snowflake.
// Note: FuncPostAuthSaml will return a fully qualified error if
// there is something wrong getting data from Snowflake.
func getIdpURLProofKey(
	ctx context.Context,
	sr *snowflakeRestful,
	authenticator string,
	application string,
	account string,
	user string,
	callbackPort int) (string, string, error) {

	headers := make(map[string]string)
	headers[httpHeaderContentType] = headerContentTypeApplicationJSON
	headers[httpHeaderAccept] = headerContentTypeApplicationJSON
	headers[httpHeaderUserAgent] = userAgent

	clientEnvironment := authRequestClientEnvironment{
		Application: application,
		Os:          operatingSystem,
		OsVersion:   platform,
	}

	requestMain := authRequestData{
		ClientAppID:             clientType,
		ClientAppVersion:        SnowflakeGoDriverVersion,
		AccountName:             account,
		LoginName:               user,
		ClientEnvironment:       clientEnvironment,
		Authenticator:           authenticator,
		BrowserModeRedirectPort: strconv.Itoa(callbackPort),
	}

	authRequest := authRequest{
		Data: requestMain,
	}

	jsonBody, err := json.Marshal(authRequest)
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to serialize json. err: %v", err)
		return "", "", err
	}

	respd, err := sr.FuncPostAuthSAML(ctx, sr, headers, jsonBody, sr.LoginTimeout)
	if err != nil {
		return "", "", err
	}
	if !respd.Success {
		logger.WithContext(ctx).Errorln("Authentication FAILED")
		sr.TokenAccessor.SetTokens("", "", -1)
		code, err := strconv.Atoi(respd.Code)
		if err != nil {
			code = -1
			return "", "", err
		}
		return "", "", &SnowflakeError{
			Number:   code,
			SQLState: SQLStateConnectionRejected,
			Message:  respd.Message,
		}
	}
	return respd.Data.SSOURL, respd.Data.ProofKey, nil
}

// Gets the login URL for multiple SAML
func getLoginURL(sr *snowflakeRestful, user string, callbackPort int) (string, string, error) {
	proofKey := generateProofKey()

	params := &url.Values{}
	params.Add("login_name", user)
	params.Add("browser_mode_redirect_port", strconv.Itoa(callbackPort))
	params.Add("proof_key", proofKey)
	url := sr.getFullURL(consoleLoginRequestPath, params)

	return url.String(), proofKey, nil
}

func generateProofKey() string {
	randomness := getSecureRandom(32)
	return base64.StdEncoding.WithPadding(base64.StdPadding).EncodeToString(randomness)
}

// The response returned from Snowflake looks like so:
// GET /?token=encodedSamlToken
// Host: localhost:54001
// Connection: keep-alive
// Upgrade-Insecure-Requests: 1
// User-Agent: userAgentStr
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8
// Referer: https://myaccount.snowflakecomputing.com/fed/login
// Accept-Encoding: gzip, deflate, br
// Accept-Language: en-US,en;q=0.9
// This extracts the token portion of the response.
func getTokenFromResponse(response string) (string, error) {
	start := "GET /?token="
	arr := strings.Split(response, "\r\n")
	if !strings.HasPrefix(arr[0], start) {
		logger.Errorf("response is malformed. ")
		return "", &SnowflakeError{
			Number:      ErrFailedToParseResponse,
			SQLState:    SQLStateConnectionRejected,
			Message:     errMsgFailedToParseResponse,
			MessageArgs: []interface{}{response},
		}
	}
	token := strings.TrimPrefix(arr[0], start)
	token = strings.Split(token, " ")[0]
	return token, nil
}

type authenticateByExternalBrowserResult struct {
	escapedSamlResponse []byte
	proofKey            []byte
	err                 error
}

func authenticateByExternalBrowser(
	ctx context.Context,
	sr *snowflakeRestful,
	authenticator string,
	application string,
	account string,
	user string,
	password string,
	externalBrowserTimeout time.Duration,
	disableConsoleLogin ConfigBool,
) ([]byte, []byte, error) {
	resultChan := make(chan authenticateByExternalBrowserResult, 1)
	go GoroutineWrapper(
		ctx,
		func() {
			resultChan <- doAuthenticateByExternalBrowser(ctx, sr, authenticator, application, account, user, password, disableConsoleLogin)
		},
	)
	select {
	case <-time.After(externalBrowserTimeout):
		return nil, nil, errors.New("authentication timed out")
	case result := <-resultChan:
		return result.escapedSamlResponse, result.proofKey, result.err
	}
}

// Authentication by an external browser takes place via the following:
//   - the golang snowflake driver communicates to Snowflake that the user wishes to
//     authenticate via external browser
//   - snowflake sends back the IDP Url configured at the Snowflake side for the
//     provided account, or use the multiple SAML way via console login
//   - the default browser is opened to that URL
//   - user authenticates at the IDP, and is redirected to Snowflake
//   - Snowflake directs the user back to the driver
//   - authenticate is complete!
func doAuthenticateByExternalBrowser(
	ctx context.Context,
	sr *snowflakeRestful,
	authenticator string,
	application string,
	account string,
	user string,
	password string,
	disableConsoleLogin ConfigBool,
) authenticateByExternalBrowserResult {
	l, err := createLocalTCPListener()
	if err != nil {
		return authenticateByExternalBrowserResult{nil, nil, err}
	}
	defer l.Close()

	callbackPort := l.Addr().(*net.TCPAddr).Port

	var loginURL string
	var proofKey string
	if disableConsoleLogin == ConfigBoolTrue {
		// Gets the IDP URL and Proof Key from Snowflake
		loginURL, proofKey, err = getIdpURLProofKey(ctx, sr, authenticator, application, account, user, callbackPort)
	} else {
		// Multiple SAML way to do authentication via console login
		loginURL, proofKey, err = getLoginURL(sr, user, callbackPort)
	}

	if err != nil {
		return authenticateByExternalBrowserResult{nil, nil, err}
	}

	if err = openBrowser(loginURL); err != nil {
		return authenticateByExternalBrowserResult{nil, nil, err}
	}

	encodedSamlResponseChan := make(chan string)
	errChan := make(chan error)

	var encodedSamlResponse string
	var errFromGoroutine error
	conn, err := l.Accept()
	if err != nil {
		logger.WithContext(ctx).Errorf("unable to accept connection. err: %v", err)
		log.Fatal(err)
	}
	go func(c net.Conn) {
		var buf bytes.Buffer
		total := 0
		encodedSamlResponse := ""
		var errAccept error
		for {
			b := make([]byte, bufSize)
			n, err := c.Read(b)
			if err != nil {
				if err != io.EOF {
					logger.WithContext(ctx).Infof("error reading from socket. err: %v", err)
					errAccept = &SnowflakeError{
						Number:      ErrFailedToGetExternalBrowserResponse,
						SQLState:    SQLStateConnectionRejected,
						Message:     errMsgFailedToGetExternalBrowserResponse,
						MessageArgs: []interface{}{err},
					}
				}
				break
			}
			total += n
			buf.Write(b)
			if n < bufSize {
				// We successfully read all data
				s := string(buf.Bytes()[:total])
				encodedSamlResponse, errAccept = getTokenFromResponse(s)
				break
			}
			buf.Grow(bufSize)
		}
		if encodedSamlResponse != "" {
			httpResponse := buildResponse(application)
			c.Write(httpResponse.Bytes())
		}
		c.Close()
		encodedSamlResponseChan <- encodedSamlResponse
		errChan <- errAccept
	}(conn)

	encodedSamlResponse = <-encodedSamlResponseChan
	errFromGoroutine = <-errChan

	if errFromGoroutine != nil {
		return authenticateByExternalBrowserResult{nil, nil, errFromGoroutine}
	}

	escapedSamlResponse, err := url.QueryUnescape(encodedSamlResponse)
	if err != nil {
		logger.WithContext(ctx).Errorf("unable to unescape saml response. err: %v", err)
		return authenticateByExternalBrowserResult{nil, nil, err}
	}
	return authenticateByExternalBrowserResult{[]byte(escapedSamlResponse), []byte(proofKey), nil}
}
