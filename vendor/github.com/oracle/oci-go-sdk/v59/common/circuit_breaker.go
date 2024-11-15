// Copyright (c) 2016, 2018, 2022, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.

package common

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

const (
	// CircuitBreakerDefaultFailureRateThreshold is the requests failure rate which calculates in at most 120 seconds, once reaches to this rate, the circuit breaker state changes from closed to open
	CircuitBreakerDefaultFailureRateThreshold float64 = 0.80
	// CircuitBreakerDefaultClosedWindow is the default value of closeStateWindow, which is the cyclic period of the closed state
	CircuitBreakerDefaultClosedWindow time.Duration = 120 * time.Second
	// CircuitBreakerDefaultResetTimeout is the default value of openStateWindow, which is the wait time before setting the breaker to halfOpen state from open state
	CircuitBreakerDefaultResetTimeout time.Duration = 30 * time.Second
	// CircuitBreakerDefaultVolumeThreshold is the default value of minimumRequests in closed status
	CircuitBreakerDefaultVolumeThreshold uint32 = 10
	// DefaultCircuitBreakerName is the name of the circuit breaker
	DefaultCircuitBreakerName string = "DefaultCircuitBreaker"
	// DefaultCircuitBreakerServiceName is the servicename of the circuit breaker
	DefaultCircuitBreakerServiceName string = ""
)

// CircuitBreakerSetting wraps all exposed configurable params of circuit breaker
type CircuitBreakerSetting struct {
	// Name is the Circuit Breaker's identifier
	name string
	// isEnabled is the switch of the circuit breaker, used for disable circuit breaker
	isEnabled bool
	// closeStateWindow is the cyclic period of the closed state, the default value is 120 seconds
	closeStateWindow time.Duration
	// openStateWindow is the wait time before setting the breaker to halfOpen state from open state, the default value is 30 seconds
	openStateWindow time.Duration
	// failureRateThreshold is the failure rate which calculates in at most closeStateWindow seconds, once reaches to this rate, the circuit breaker state changes from closed to open
	// the circuit will transition from closed to open, the default value is 80%
	failureRateThreshold float64
	// minimumRequests is the minimum number of counted requests in closed state, the default value is 10 requests
	minimumRequests uint32
	// successStatCodeMap is the error(s) of StatusCode returned from service, which should be considered as the success or failure accounted by circuit breaker
	// successStatCodeMap and successStatErrCodeMap are combined to use, if both StatusCode and ErrorCode are required, no need to add it to successStatCodeMap,
	// the default value is [429, 500, 502, 503, 504]
	successStatCodeMap map[int]bool
	// successStatErrCodeMap is the error(s) of StatusCode and ErrorCode returned from service, which should be considered
	// as the success or failure accounted by circuit breaker
	// the default value is {409, "IncorrectState"}
	successStatErrCodeMap map[StatErrCode]bool
	// serviceName is the name of the service which can be set using withServiceName option for NewCircuitBreaker.
	// The default value is empty string
	serviceName string
}

// Convert CircuitBreakerSetting to human-readable string representation
func (cbst CircuitBreakerSetting) String() string {
	return fmt.Sprintf("{name=%v, isEnabled=%v, closeStateWindow=%v, openStateWindow=%v, failureRateThreshold=%v, minimumRequests=%v, successStatCodeMap=%v, successStatErrCodeMap=%v, serviceName=%v}",
		cbst.name, cbst.isEnabled, cbst.closeStateWindow, cbst.openStateWindow, cbst.failureRateThreshold, cbst.minimumRequests, cbst.successStatCodeMap, cbst.successStatErrCodeMap, cbst.serviceName)
}

// OciCircuitBreaker wraps all exposed configurable params of circuit breaker and 3P gobreaker CircuirBreaker
type OciCircuitBreaker struct {
	Cbst *CircuitBreakerSetting
	Cb   *gobreaker.CircuitBreaker
}

// NewOciCircuitBreaker is used for initializing specified oci circuit breaker configuration with circuit breaker settings
func NewOciCircuitBreaker(cbst *CircuitBreakerSetting, gbcb *gobreaker.CircuitBreaker) *OciCircuitBreaker {
	ocb := new(OciCircuitBreaker)
	ocb.Cbst = cbst
	ocb.Cb = gbcb

	return ocb
}

// CircuitBreakerOption is the type of the options for NewCircuitBreakerWithOptions.
type CircuitBreakerOption func(cbst *CircuitBreakerSetting)

// NewGoCircuitBreaker is a function to initialize a CircuitBreaker object with the specified configuration
// Add the interface, to allow the user directly use the 3P gobreaker.Setting's params.
func NewGoCircuitBreaker(st gobreaker.Settings) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(st)
}

// DefaultCircuitBreakerSetting is used for set circuit breaker with default config
func DefaultCircuitBreakerSetting() *CircuitBreakerSetting {
	successStatErrCodeMap := map[StatErrCode]bool{
		{409, "IncorrectState"}: false,
	}
	successStatCodeMap := map[int]bool{
		429: false,
		500: false,
		502: false,
		503: false,
		504: false,
	}
	return newCircuitBreakerSetting(
		WithName(DefaultCircuitBreakerName),
		WithIsEnabled(true),
		WithCloseStateWindow(CircuitBreakerDefaultClosedWindow),
		WithOpenStateWindow(CircuitBreakerDefaultResetTimeout),
		WithFailureRateThreshold(CircuitBreakerDefaultFailureRateThreshold),
		WithMinimumRequests(CircuitBreakerDefaultVolumeThreshold),
		WithSuccessStatErrCodeMap(successStatErrCodeMap),
		WithSuccessStatCodeMap(successStatCodeMap))
}

// DefaultCircuitBreakerSettingWithServiceName is used for set circuit breaker with default config
func DefaultCircuitBreakerSettingWithServiceName() *CircuitBreakerSetting {
	successStatErrCodeMap := map[StatErrCode]bool{
		{409, "IncorrectState"}: false,
	}
	successStatCodeMap := map[int]bool{
		429: false,
		500: false,
		502: false,
		503: false,
		504: false,
	}
	return newCircuitBreakerSetting(
		WithName(DefaultCircuitBreakerName),
		WithIsEnabled(true),
		WithCloseStateWindow(CircuitBreakerDefaultClosedWindow),
		WithOpenStateWindow(CircuitBreakerDefaultResetTimeout),
		WithFailureRateThreshold(CircuitBreakerDefaultFailureRateThreshold),
		WithMinimumRequests(CircuitBreakerDefaultVolumeThreshold),
		WithSuccessStatErrCodeMap(successStatErrCodeMap),
		WithSuccessStatCodeMap(successStatCodeMap),
		WithServiceName(DefaultCircuitBreakerServiceName))
}

// NoCircuitBreakerSetting is used for disable Circuit Breaker
func NoCircuitBreakerSetting() *CircuitBreakerSetting {
	return NewCircuitBreakerSettingWithOptions(WithIsEnabled(false))
}

// NewCircuitBreakerSettingWithOptions is a helper method to assemble a CircuitBreakerSetting object.
// It starts out with the values returned by defaultCircuitBreakerSetting().
func NewCircuitBreakerSettingWithOptions(opts ...CircuitBreakerOption) *CircuitBreakerSetting {
	cbst := DefaultCircuitBreakerSettingWithServiceName()
	// allow changing values
	for _, opt := range opts {
		opt(cbst)
	}
	if defaultLogger.LogLevel() == verboseLogging {
		Debugf("Circuit Breaker setting: %s\n", cbst.String())
	}

	return cbst
}

// NewCircuitBreaker is used for initialing specified circuit breaker configuration with base client
func NewCircuitBreaker(cbst *CircuitBreakerSetting) *OciCircuitBreaker {
	if !cbst.isEnabled {
		return nil
	}

	st := gobreaker.Settings{}
	customizeGoBreakerSetting(&st, cbst)
	gbcb := gobreaker.NewCircuitBreaker(st)

	return NewOciCircuitBreaker(cbst, gbcb)
}

func newCircuitBreakerSetting(opts ...CircuitBreakerOption) *CircuitBreakerSetting {
	cbSetting := CircuitBreakerSetting{}

	// allow changing values
	for _, opt := range opts {
		opt(&cbSetting)
	}
	return &cbSetting
}

// customizeGoBreakerSetting is used for converting CircuitBreakerSetting to 3P gobreaker's setting type
func customizeGoBreakerSetting(st *gobreaker.Settings, cbst *CircuitBreakerSetting) {
	st.Name = cbst.name
	st.Timeout = cbst.openStateWindow
	st.Interval = cbst.closeStateWindow
	st.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		if to == gobreaker.StateOpen {
			Debugf("Circuit Breaker %s is now in Open State\n", name)
		}
	}
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= cbst.minimumRequests && failureRatio >= cbst.failureRateThreshold
	}
	st.IsSuccessful = func(err error) bool {
		if serviceErr, ok := IsServiceError(err); ok {
			if isSuccessful, ok := cbst.successStatCodeMap[serviceErr.GetHTTPStatusCode()]; ok {
				return isSuccessful
			}
			if isSuccessful, ok := cbst.successStatErrCodeMap[StatErrCode{serviceErr.GetHTTPStatusCode(), serviceErr.GetCode()}]; ok {
				return isSuccessful
			}
		}
		return true
	}
}

// WithName is the option for NewCircuitBreaker that sets the Name.
func WithName(name string) CircuitBreakerOption {
	// this is the CircuitBreakerOption function type
	return func(cbst *CircuitBreakerSetting) {
		cbst.name = name
	}
}

// WithIsEnabled is the option for NewCircuitBreaker that sets the isEnabled.
func WithIsEnabled(isEnabled bool) CircuitBreakerOption {
	// this is the CircuitBreakerOption function type
	return func(cbst *CircuitBreakerSetting) {
		cbst.isEnabled = isEnabled
	}
}

// WithCloseStateWindow is the option for NewCircuitBreaker that sets the closeStateWindow.
func WithCloseStateWindow(window time.Duration) CircuitBreakerOption {
	// this is the CircuitBreakerOption function type
	return func(cbst *CircuitBreakerSetting) {
		cbst.closeStateWindow = window
	}
}

// WithOpenStateWindow is the option for NewCircuitBreaker that sets the openStateWindow.
func WithOpenStateWindow(window time.Duration) CircuitBreakerOption {
	// this is the CircuitBreakerOption function type
	return func(cbst *CircuitBreakerSetting) {
		cbst.openStateWindow = window
	}
}

// WithFailureRateThreshold is the option for NewCircuitBreaker that sets the failureRateThreshold.
func WithFailureRateThreshold(threshold float64) CircuitBreakerOption {
	// this is the CircuitBreakerOption function type
	return func(cbst *CircuitBreakerSetting) {
		cbst.failureRateThreshold = threshold
	}
}

// WithMinimumRequests is the option for NewCircuitBreaker that sets the minimumRequests.
func WithMinimumRequests(num uint32) CircuitBreakerOption {
	// this is the CircuitBreakerOption function type
	return func(cbst *CircuitBreakerSetting) {
		cbst.minimumRequests = num
	}
}

// WithSuccessStatCodeMap is the option for NewCircuitBreaker that sets the successStatCodeMap.
func WithSuccessStatCodeMap(successStatCodeMap map[int]bool) CircuitBreakerOption {
	// this is the CircuitBreakerOption function type
	return func(cbst *CircuitBreakerSetting) {
		cbst.successStatCodeMap = successStatCodeMap
	}
}

// WithSuccessStatErrCodeMap is the option for NewCircuitBreaker that sets the successStatErrCodeMap.
func WithSuccessStatErrCodeMap(successStatErrCodeMap map[StatErrCode]bool) CircuitBreakerOption {
	// this is the CircuitBreakerOption function type
	return func(cbst *CircuitBreakerSetting) {
		cbst.successStatErrCodeMap = successStatErrCodeMap
	}
}

// WithServiceName is the option for NewCircuitBreaker that sets the ServiceName.
func WithServiceName(serviceName string) CircuitBreakerOption {
	// this is the CircuitBreakerOption function type
	return func(cbst *CircuitBreakerSetting) {
		cbst.serviceName = serviceName
	}
}

// GlobalCircuitBreakerSetting is global level circuit breaker setting, it would impact all services, the precedence is lower
// than client level circuit breaker
var GlobalCircuitBreakerSetting *CircuitBreakerSetting = nil

// ConfigCircuitBreakerFromEnvVar is used for checking the circuit breaker environment variable setting, default value is nil
func ConfigCircuitBreakerFromEnvVar(baseClient *BaseClient) {
	if IsEnvVarTrue(isDefaultCircuitBreakerEnabled) {
		baseClient.Configuration.CircuitBreaker = NewCircuitBreaker(DefaultCircuitBreakerSetting())
		return
	}
	if IsEnvVarFalse(isDefaultCircuitBreakerEnabled) {
		baseClient.Configuration.CircuitBreaker = nil
	}
}

// ConfigCircuitBreakerFromGlobalVar is used for checking if global circuitBreakerSetting is configured, the priority is higher than cb env var
func ConfigCircuitBreakerFromGlobalVar(baseClient *BaseClient) {
	if GlobalCircuitBreakerSetting != nil {
		baseClient.Configuration.CircuitBreaker = NewCircuitBreaker(GlobalCircuitBreakerSetting)
	}
}
