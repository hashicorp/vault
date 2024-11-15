package linodego

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/http2"
)

const (
	// nolint:unused
	httpRetryAfterHeaderName = "Retry-After"
	// nolint:unused
	httpMaintenanceModeHeaderName = "X-Maintenance-Mode"

	// nolint:unused
	httpDefaultRetryCount = 1000
)

// RetryConditional is a type alias for a function that determines if a request should be retried based on the response and error.
// nolint:unused
type httpRetryConditional func(*http.Response, error) bool

// RetryAfter is a type alias for a function that determines the duration to wait before retrying based on the response.
// nolint:unused
type httpRetryAfter func(*http.Response) (time.Duration, error)

// Configures http.Client to lock until enough time has passed to retry the request as determined by the Retry-After response header.
// If the Retry-After header is not set, we fall back to the value of SetPollDelay.
// nolint:unused
func httpConfigureRetries(c *httpClient) {
	c.retryConditionals = append(c.retryConditionals, httpcheckRetryConditionals(c))
	c.retryAfter = httpRespectRetryAfter
}

// nolint:unused
func httpcheckRetryConditionals(c *httpClient) httpRetryConditional {
	return func(resp *http.Response, err error) bool {
		for _, retryConditional := range c.retryConditionals {
			retry := retryConditional(resp, err)
			if retry {
				log.Printf("[INFO] Received error %v - Retrying", err)
				return true
			}
		}
		return false
	}
}

// nolint:unused
func httpRespectRetryAfter(resp *http.Response) (time.Duration, error) {
	retryAfterStr := resp.Header.Get(retryAfterHeaderName)
	if retryAfterStr == "" {
		return 0, nil
	}

	retryAfter, err := strconv.Atoi(retryAfterStr)
	if err != nil {
		return 0, err
	}

	duration := time.Duration(retryAfter) * time.Second
	log.Printf("[INFO] Respecting Retry-After Header of %d (%s)", retryAfter, duration)
	return duration, nil
}

// Retry conditions

// nolint:unused
func httpLinodeBusyRetryCondition(resp *http.Response, _ error) bool {
	apiError, ok := getAPIError(resp)
	linodeBusy := ok && apiError.Error() == "Linode busy."
	retry := resp.StatusCode == http.StatusBadRequest && linodeBusy
	return retry
}

// nolint:unused
func httpTooManyRequestsRetryCondition(resp *http.Response, _ error) bool {
	return resp.StatusCode == http.StatusTooManyRequests
}

// nolint:unused
func httpServiceUnavailableRetryCondition(resp *http.Response, _ error) bool {
	serviceUnavailable := resp.StatusCode == http.StatusServiceUnavailable

	// During maintenance events, the API will return a 503 and add
	// an `X-MAINTENANCE-MODE` header. Don't retry during maintenance
	// events, only for legitimate 503s.
	if serviceUnavailable && resp.Header.Get(maintenanceModeHeaderName) != "" {
		log.Printf("[INFO] Linode API is under maintenance, request will not be retried - please see status.linode.com for more information")
		return false
	}

	return serviceUnavailable
}

// nolint:unused
func httpRequestTimeoutRetryCondition(resp *http.Response, _ error) bool {
	return resp.StatusCode == http.StatusRequestTimeout
}

// nolint:unused
func httpRequestGOAWAYRetryCondition(_ *http.Response, err error) bool {
	return errors.As(err, &http2.GoAwayError{})
}

// nolint:unused
func httpRequestNGINXRetryCondition(resp *http.Response, _ error) bool {
	return resp.StatusCode == http.StatusBadRequest &&
		resp.Header.Get("Server") == "nginx" &&
		resp.Header.Get("Content-Type") == "text/html"
}

// Helper function to extract APIError from response
// nolint:unused
func getAPIError(resp *http.Response) (*APIError, bool) {
	var apiError APIError
	err := json.NewDecoder(resp.Body).Decode(&apiError)
	if err != nil {
		return nil, false
	}
	return &apiError, true
}
