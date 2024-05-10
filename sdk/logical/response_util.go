// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/errwrap"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/helper/consts"
)

// RespondErrorCommon pulls most of the functionality from http's
// respondErrorCommon and some of http's handleLogical and makes it available
// to both the http package and elsewhere.
func RespondErrorCommon(req *Request, resp *Response, err error) (int, error) {
	if err == nil && (resp == nil || !resp.IsError()) {
		switch {
		case req.Operation == ReadOperation || req.Operation == HeaderOperation:
			if resp == nil {
				return http.StatusNotFound, nil
			}

		// Basically: if we have empty "keys" or no keys at all, 404. This
		// provides consistency with GET.
		case req.Operation == ListOperation && (resp == nil || resp.WrapInfo == nil):
			if resp == nil {
				return http.StatusNotFound, nil
			}
			if len(resp.Data) == 0 {
				if len(resp.Warnings) > 0 {
					return 0, nil
				}
				return http.StatusNotFound, nil
			}
			keysRaw, ok := resp.Data["keys"]
			if !ok || keysRaw == nil {
				// If we don't have keys but have other data, return as-is
				if len(resp.Data) > 0 || len(resp.Warnings) > 0 {
					return 0, nil
				}
				return http.StatusNotFound, nil
			}

			var keys []string
			switch keysRaw.(type) {
			case []interface{}:
				keys = make([]string, len(keysRaw.([]interface{})))
				for i, el := range keysRaw.([]interface{}) {
					s, ok := el.(string)
					if !ok {
						return http.StatusInternalServerError, nil
					}
					keys[i] = s
				}

			case []string:
				keys = keysRaw.([]string)
			default:
				return http.StatusInternalServerError, nil
			}

			if len(keys) == 0 {
				return http.StatusNotFound, nil
			}
		}

		return 0, nil
	}

	if errwrap.ContainsType(err, new(ReplicationCodedError)) {
		var allErrors error
		var codedErr *ReplicationCodedError
		errwrap.Walk(err, func(inErr error) {
			// The Walk function does not just traverse leaves, and execute the
			// callback function on the entire error first. So, if the error is
			// of type multierror.Error, we may want to skip storing the entire
			// error first to avoid adding duplicate errors when walking down
			// the leaf errors
			if _, ok := inErr.(*multierror.Error); ok {
				return
			}
			newErr, ok := inErr.(*ReplicationCodedError)
			if ok {
				codedErr = newErr
			} else {
				// if the error is of type fmt.wrapError which is typically
				// made by calling fmt.Errorf("... %w", err), allErrors will
				// contain duplicated error messages
				allErrors = multierror.Append(allErrors, inErr)
			}
		})
		if allErrors != nil {
			return codedErr.Code, multierror.Append(fmt.Errorf("errors from both primary and secondary; primary error was %v; secondary errors follow", codedErr.Msg), allErrors)
		}
		return codedErr.Code, errors.New(codedErr.Msg)
	}

	// Start out with internal server error since in most of these cases there
	// won't be a response so this won't be overridden
	statusCode := http.StatusInternalServerError
	// If we actually have a response, start out with bad request
	if resp != nil {
		statusCode = http.StatusBadRequest
	}

	// Now, check the error itself; if it has a specific logical error, set the
	// appropriate code
	if err != nil {
		switch {
		case errwrap.Contains(err, consts.ErrOverloaded.Error()):
			statusCode = http.StatusServiceUnavailable
		case errwrap.ContainsType(err, new(StatusBadRequest)):
			statusCode = http.StatusBadRequest
		case errwrap.Contains(err, ErrPermissionDenied.Error()):
			statusCode = http.StatusForbidden
		case errwrap.Contains(err, consts.ErrInvalidWrappingToken.Error()):
			statusCode = http.StatusBadRequest
		case errwrap.Contains(err, ErrUnsupportedOperation.Error()):
			statusCode = http.StatusMethodNotAllowed
		case errwrap.Contains(err, ErrUnsupportedPath.Error()):
			statusCode = http.StatusNotFound
		case errwrap.Contains(err, ErrInvalidRequest.Error()):
			statusCode = http.StatusBadRequest
		case errwrap.Contains(err, ErrUpstreamRateLimited.Error()):
			statusCode = http.StatusBadGateway
		case errwrap.Contains(err, ErrRateLimitQuotaExceeded.Error()):
			statusCode = http.StatusTooManyRequests
		case errwrap.Contains(err, ErrLeaseCountQuotaExceeded.Error()):
			statusCode = http.StatusTooManyRequests
		case errwrap.Contains(err, ErrMissingRequiredState.Error()):
			statusCode = http.StatusPreconditionFailed
		case errwrap.Contains(err, ErrPathFunctionalityRemoved.Error()):
			statusCode = http.StatusNotFound
		case errwrap.Contains(err, ErrRelativePath.Error()):
			statusCode = http.StatusBadRequest
		case errwrap.Contains(err, ErrInvalidCredentials.Error()):
			statusCode = http.StatusBadRequest
		case errors.Is(err, ErrNotFound):
			statusCode = http.StatusNotFound
		}
	}

	if respErr := resp.Error(); respErr != nil {
		err = fmt.Errorf("%s", respErr.Error())

		// Don't let other error codes override the overloaded status code
		if strings.Contains(respErr.Error(), consts.ErrOverloaded.Error()) {
			statusCode = http.StatusServiceUnavailable
		}
	}

	return statusCode, err
}

// AdjustErrorStatusCode adjusts the status that will be sent in error
// conditions in a way that can be shared across http's respondError and other
// locations.
func AdjustErrorStatusCode(status *int, err error) {
	// Handle nested errors
	if t, ok := err.(*multierror.Error); ok {
		for _, e := range t.Errors {
			AdjustErrorStatusCode(status, e)
		}
	}

	// Adjust status code when overloaded
	if errwrap.Contains(err, consts.ErrOverloaded.Error()) {
		*status = http.StatusServiceUnavailable
	}

	// Adjust status code when sealed
	if errwrap.Contains(err, consts.ErrSealed.Error()) {
		*status = http.StatusServiceUnavailable
	}

	if errwrap.Contains(err, consts.ErrAPILocked.Error()) {
		*status = http.StatusServiceUnavailable
	}

	// Adjust status code on
	if errwrap.Contains(err, "http: request body too large") {
		*status = http.StatusRequestEntityTooLarge
	}

	// Allow HTTPCoded error passthrough to specify a code
	if t, ok := err.(HTTPCodedError); ok {
		*status = t.Code()
	}
}

func RespondError(w http.ResponseWriter, status int, err error) {
	AdjustErrorStatusCode(&status, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	type ErrorResponse struct {
		Errors []string `json:"errors"`
	}
	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func RespondErrorAndData(w http.ResponseWriter, status int, data interface{}, err error) {
	AdjustErrorStatusCode(&status, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	type ErrorAndDataResponse struct {
		Errors []string    `json:"errors"`
		Data   interface{} `json:"data"`
	}
	resp := &ErrorAndDataResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}
	resp.Data = data

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}
