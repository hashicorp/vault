package retry

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/linode/linodego/pkg/errors"
)

const retryAfterHeaderName = "Retry-After"

// ConditionFunc represents a function that evaluates whether a request
// should be retried, given the response and error.
type ConditionFunc resty.RetryConditionFunc

// LinodeBusyRetryCondition is a ConditionFunc that retries requests which have
// resulted in "Linode busy." errors.
func LinodeBusyRetryCondition(r *resty.Response, _ error) bool {
	apiError, ok := r.Error().(*errors.APIError)
	linodeBusy := ok && apiError.Error() == "Linode busy."
	retry := r.StatusCode() == http.StatusBadRequest && linodeBusy
	return retry
}

// TooManyRequestsRetryConditon is a ConditionFunc that retries requests which
// have resulted in an HTTP status code 429.
func TooManyRequestsRetryCondition(r *resty.Response, _ error) bool {
	return r.StatusCode() == http.StatusTooManyRequests
}

// ServiceUnavailableRetryCondition is a ConditionFunc that retries requests
// have resulted in an HTTP status code 503.
func ServiceUnavailableRetryCondition(r *resty.Response, _ error) bool {
	return r.StatusCode() == http.StatusServiceUnavailable
}

// RespectRetryAfter configures the resty client to respect Retry-After
// headers.
func RespectRetryAfter(client *resty.Client, resp *resty.Response) (time.Duration, error) {
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
