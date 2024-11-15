// Copyright (c) 2016, 2018, 2022, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.

package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/sony/gobreaker"
)

// ServiceError models all potential errors generated the service call
type ServiceError interface {
	// The http status code of the error
	GetHTTPStatusCode() int

	// The human-readable error string as sent by the service
	GetMessage() string

	// A short error code that defines the error, meant for programmatic parsing.
	// See https://docs.cloud.oracle.com/Content/API/References/apierrors.htm
	GetCode() string

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	GetOpcRequestID() string
}

type servicefailure struct {
	StatusCode   int
	Code         string `json:"code,omitempty"`
	Message      string `json:"message,omitempty"`
	OpcRequestID string `json:"opc-request-id"`
}

func newServiceFailureFromResponse(response *http.Response) error {
	var err error

	se := servicefailure{
		StatusCode:   response.StatusCode,
		Code:         "BadErrorResponse",
		OpcRequestID: response.Header.Get("opc-request-id")}

	//If there is an error consume the body, entirely
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		se.Message = fmt.Sprintf("The body of the response was not readable, due to :%s", err.Error())
		return se
	}

	err = json.Unmarshal(body, &se)
	if err != nil {
		Debugf("Error response could not be parsed due to: %s", err.Error())
		se.Message = fmt.Sprintf("Failed to parse json from response body due to: %s. With response body %s.", err.Error(), string(body[:]))
		return se
	}
	return se
}

func (se servicefailure) Error() string {
	return fmt.Sprintf("Service error:%s. %s. http status code: %d. Opc request id: %s",
		se.Code, se.Message, se.StatusCode, se.OpcRequestID)
}

func (se servicefailure) GetHTTPStatusCode() int {
	return se.StatusCode

}

func (se servicefailure) GetMessage() string {
	return se.Message
}

func (se servicefailure) GetCode() string {
	return se.Code
}

func (se servicefailure) GetOpcRequestID() string {
	return se.OpcRequestID
}

// IsServiceError returns false if the error is not service side, otherwise true
// additionally it returns an interface representing the ServiceError
func IsServiceError(err error) (failure ServiceError, ok bool) {
	failure, ok = err.(servicefailure)
	return
}

type deadlineExceededByBackoffError struct{}

func (deadlineExceededByBackoffError) Error() string {
	return "now() + computed backoff duration exceeds request deadline"
}

// DeadlineExceededByBackoff is the error returned by Call() when GetNextDuration() returns a time.Duration that would
// force the user to wait past the request deadline before re-issuing a request. This enables us to exit early, since
// we cannot succeed based on the configured retry policy.
var DeadlineExceededByBackoff error = deadlineExceededByBackoffError{}

// NonSeekableRequestRetryFailure is the error returned when the request is with binary request body, and is configured
// retry, but the request body is not retryable
type NonSeekableRequestRetryFailure struct {
	err error
}

func (ne NonSeekableRequestRetryFailure) Error() string {
	if ne.err == nil {
		return fmt.Sprintf("Unable to perform Retry on this request body type, which did not implement seek() interface")
	}
	return fmt.Sprintf("%s. Unable to perform Retry on this request body type, which did not implement seek() interface", ne.err.Error())
}

// IsNetworkError validates if an error is a net.Error and check if it's temporary or timeout
func IsNetworkError(err error) bool {
	if r, ok := err.(net.Error); ok && (r.Temporary() || r.Timeout()) {
		return true
	}
	return false
}

// IsCircuitBreakerError validates if an error's text is Open state ErrOpenState or HalfOpen state ErrTooManyRequests
func IsCircuitBreakerError(err error) bool {
	if err.Error() == gobreaker.ErrOpenState.Error() || err.Error() == gobreaker.ErrTooManyRequests.Error() {
		return true
	}
	return false
}

func getCircuitBreakerError(request *http.Request, err error, cbr *OciCircuitBreaker) error {
	cbErr := fmt.Errorf(" %s. This request was not sent to the service. Look for earlier errors to determine why the circuit breaker was opened.\n An open circuit breaker means %s service failed too many times in the recent past. "+
		"Because the circuit breaker has been opened, requests within the openStateWindow of %.2f seconds since the circuit breaker was opened will not be sent to the service.\n"+
		"For more information on the exact errors with which the %s service responds to your requests, enable Info-level logs and rerun your code.\n "+
		"URL which circuit breaker prevented request to - %s \n Circuit Breaker Info \n Name - %s \n State - %s \n  Number of requests - %d \n Number of success - %d \n Number of failures - %d ",
		err, cbr.Cbst.serviceName, cbr.Cbst.openStateWindow.Seconds(), cbr.Cbst.serviceName, request.URL.Host+request.URL.Path, cbr.Cbst.name, cbr.Cb.State().String(), cbr.Cb.Counts().Requests, cbr.Cb.Counts().TotalSuccesses, cbr.Cb.Counts().TotalFailures)
	return cbErr
}

// StatErrCode is a type which wraps error's statusCode and errorCode from service end
type StatErrCode struct {
	statusCode int
	errorCode  string
}
