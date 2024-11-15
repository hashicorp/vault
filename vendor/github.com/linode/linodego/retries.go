package linodego

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/http2"
)

const (
	retryAfterHeaderName      = "Retry-After"
	maintenanceModeHeaderName = "X-Maintenance-Mode"

	defaultRetryCount = 1000
)

// type RetryConditional func(r *resty.Response) (shouldRetry bool)
type RetryConditional resty.RetryConditionFunc

// type RetryAfter func(c *resty.Client, r *resty.Response) (time.Duration, error)
type RetryAfter resty.RetryAfterFunc

// Configures resty to
// lock until enough time has passed to retry the request as determined by the Retry-After response header.
// If the Retry-After header is not set, we fall back to value of SetPollDelay.
func configureRetries(c *Client) {
	c.resty.
		SetRetryCount(defaultRetryCount).
		AddRetryCondition(checkRetryConditionals(c)).
		SetRetryAfter(respectRetryAfter)
}

func checkRetryConditionals(c *Client) func(*resty.Response, error) bool {
	return func(r *resty.Response, err error) bool {
		for _, retryConditional := range c.retryConditionals {
			retry := retryConditional(r, err)
			if retry {
				log.Printf("[INFO] Received error %s - Retrying", r.Error())
				return true
			}
		}
		return false
	}
}

// SetLinodeBusyRetry configures resty to retry specifically on "Linode busy." errors
// The retry wait time is configured in SetPollDelay
func linodeBusyRetryCondition(r *resty.Response, _ error) bool {
	apiError, ok := r.Error().(*APIError)
	linodeBusy := ok && apiError.Error() == "Linode busy."
	retry := r.StatusCode() == http.StatusBadRequest && linodeBusy
	return retry
}

func tooManyRequestsRetryCondition(r *resty.Response, _ error) bool {
	return r.StatusCode() == http.StatusTooManyRequests
}

func serviceUnavailableRetryCondition(r *resty.Response, _ error) bool {
	serviceUnavailable := r.StatusCode() == http.StatusServiceUnavailable

	// During maintenance events, the API will return a 503 and add
	// an `X-MAINTENANCE-MODE` header. Don't retry during maintenance
	// events, only for legitimate 503s.
	if serviceUnavailable && r.Header().Get(maintenanceModeHeaderName) != "" {
		log.Printf("[INFO] Linode API is under maintenance, request will not be retried - please see status.linode.com for more information")
		return false
	}

	return serviceUnavailable
}

func requestTimeoutRetryCondition(r *resty.Response, _ error) bool {
	return r.StatusCode() == http.StatusRequestTimeout
}

func requestGOAWAYRetryCondition(_ *resty.Response, e error) bool {
	return errors.As(e, &http2.GoAwayError{})
}

func requestNGINXRetryCondition(r *resty.Response, _ error) bool {
	return r.StatusCode() == http.StatusBadRequest &&
		r.Header().Get("Server") == "nginx" &&
		r.Header().Get("Content-Type") == "text/html"
}

func respectRetryAfter(client *resty.Client, resp *resty.Response) (time.Duration, error) {
	retryAfterStr := resp.Header().Get(retryAfterHeaderName)
	if retryAfterStr == "" {
		return 0, nil
	}

	retryAfter, err := strconv.Atoi(retryAfterStr)
	if err != nil {
		return 0, err
	}

	duration := time.Duration(retryAfter) * time.Second
	log.Printf("[INFO] Respecting Retry-After Header of %d (%s) (max %s)", retryAfter, duration, client.RetryMaxWaitTime)
	return duration, nil
}
