// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package ocsp

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	httpHeaderContentType      = "Content-Type"
	httpHeaderAccept           = "accept"
	httpHeaderUserAgent        = "User-Agent"
	httpHeaderServiceName      = "X-Snowflake-Service"
	httpHeaderContentLength    = "Content-Length"
	httpHeaderHost             = "Host"
	httpHeaderValueOctetStream = "application/octet-stream"
	httpHeaderContentEncoding  = "Content-Encoding"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

const (
	// requestGUIDKey is attached to every request against Snowflake
	requestGUIDKey string = "request_guid"
	// retryCounterKey is attached to query-request from the second time
	retryCounterKey string = "retryCounter"
	// requestIDKey is attached to all requests to Snowflake
	requestIDKey string = "requestId"
)

// This class takes in an url during construction and replaces the value of
// request_guid every time replace() is called. If the url does not contain
// request_guid, just return the original url
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
generated UUID
*/
func (replacer *requestGUIDReplace) replace() *url.URL {
	replacer.urlValues.Del(requestGUIDKey)
	uuid, _ := uuid.GenerateUUID()
	replacer.urlValues.Add(requestGUIDKey, uuid)
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
	logger   hclog.Logger
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
	instance.logger = hclog.New(hclog.DefaultOptions)
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
	r.logger.Info("retryHTTP", "totalTimeout", totalTimeout)
	retryCounter := 0
	sleepTime := time.Duration(0)

	var rIDReplacer requestGUIDReplacer
	var rUpdater retryCounterUpdater

	for {
		r.logger.Debug("retry count", "retryCounter", retryCounter)
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
			r.logger.Warn(
				"failed http connection. no response is returned. retrying...", "err", err)
		} else {
			if res.StatusCode == http.StatusOK || r.raise4XX && res != nil && res.StatusCode >= 400 && res.StatusCode < 500 {
				// exit if success
				// or
				// abort connection if raise4XX flag is enabled and the range of HTTP status code are 4XX.
				// This is currently used for Snowflake login. The caller must generate an error object based on HTTP status.
				break
			}
			r.logger.Warn(
				"failed http connection. retrying...\n", "statusCode", res.StatusCode)
			res.Body.Close()
		}
		// uses decorrelated jitter backoff
		sleepTime = defaultWaitAlgo.decorr(retryCounter, sleepTime)

		if totalTimeout > 0 {
			r.logger.Info("to timeout: ", "totalTimeout", totalTimeout)
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
		r.logger.Info("sleeping to retry", "sleepTime", sleepTime, "totalTimeout", totalTimeout)

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
		if urlError.Err.Error() == "OCSP status revoked" {
			// Certificate Revoked
			return true, nil
		}
		if _, ok := urlError.Err.(x509.CertificateInvalidError); ok {
			// Certificate is invalid
			return true, err
		}
		if _, ok := urlError.Err.(x509.UnknownAuthorityError); ok {
			// Certificate is self-signed
			return true, err
		}
		errString := urlError.Err.Error()
		if runtime.GOOS == "darwin" && strings.HasPrefix(errString, "x509:") && strings.HasSuffix(errString, "certificate is expired") {
			// Certificate is expired
			return true, err
		}

	}
	return false, err
}

/*
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "{}"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
