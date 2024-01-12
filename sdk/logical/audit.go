// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"fmt"

	"github.com/mitchellh/copystructure"
)

// LogInput is used as the input to the audit system on which audit entries are based.
type LogInput struct {
	Type                string
	Auth                *Auth
	Request             *Request
	Response            *Response
	OuterErr            error
	NonHMACReqDataKeys  []string
	NonHMACRespDataKeys []string
}

type MarshalOptions struct {
	ValueHasher func(string) string
}

type OptMarshaler interface {
	MarshalJSONWithOptions(*MarshalOptions) ([]byte, error)
}

// Clone will attempt to create a deep copy of the LogInput.
// This is preferred over attempting to use other libraries such as copystructure.
// Audit formatting methods which parse LogInput into an audit request/response
// entry, call receivers on the LogInput struct to get their value. These values
// would be lost using copystructure as it cannot copy unexported fields.
// If the LogInput type or any of the subtypes referenced by LogInput fields are
// changed, then the Clone methods will need to be updated.
// NOTE: Does not deep clone the LogInput.OuterError field as it represents an
// error interface.
func (l *LogInput) Clone() (*LogInput, error) {
	// Clone Auth
	auth, err := cloneAuth(l.Auth)
	if err != nil {
		return nil, err
	}

	// Clone Request
	req, err := cloneRequest(l.Request)
	if err != nil {
		return nil, err
	}

	// Clone Response
	resp, err := cloneResponse(l.Response)
	if err != nil {
		return nil, err
	}

	// Copy HMAC keys
	reqDataKeys := make([]string, len(l.NonHMACReqDataKeys))
	copy(reqDataKeys, l.NonHMACReqDataKeys)
	respDataKeys := make([]string, len(l.NonHMACRespDataKeys))
	copy(respDataKeys, l.NonHMACRespDataKeys)

	// OuterErr is just linked in a non-deep way as it's an interface, and we
	// don't know for sure which type this might actually be.
	// At the time of writing this code, OuterErr isn't modified by anything,
	// so we shouldn't get any race issues.
	cloned := &LogInput{
		Type:                l.Type,
		Auth:                auth,
		Request:             req,
		Response:            resp,
		OuterErr:            l.OuterErr,
		NonHMACReqDataKeys:  reqDataKeys,
		NonHMACRespDataKeys: respDataKeys,
	}

	return cloned, nil
}

// clone will deep-copy the supplied struct.
// However, it cannot copy unexported fields or evaluate methods.
func clone[V any](s V) (V, error) {
	var result V

	data, err := copystructure.Copy(s)
	if err != nil {
		return result, err
	}

	result = data.(V)

	return result, err
}

// cloneAuth deep copies an Auth struct.
func cloneAuth(auth *Auth) (*Auth, error) {
	// If auth is nil, there's nothing to clone.
	if auth == nil {
		return nil, nil
	}

	auth, err := clone[*Auth](auth)
	if err != nil {
		return nil, fmt.Errorf("unable to clone auth: %w", err)
	}

	return auth, nil
}

// cloneRequest deep copies a Request struct.
// It will set unexported fields which were only previously accessible outside
// the package via receiver methods.
func cloneRequest(request *Request) (*Request, error) {
	// If request is nil, there's nothing to clone.
	if request == nil {
		return nil, nil
	}

	req, err := clone[*Request](request)
	if err != nil {
		return nil, fmt.Errorf("unable to clone request: %w", err)
	}

	// Add the unexported values that were only retrievable via receivers.
	req.mountClass = request.MountClass()
	req.mountRunningVersion = request.MountRunningVersion()
	req.mountRunningSha256 = request.MountRunningSha256()
	req.mountIsExternalPlugin = request.MountIsExternalPlugin()
	// This needs to be overwritten as the internal connection state is not cloned properly
	// mainly the big.Int serial numbers within the x509.Certificate objects get mangled.
	req.Connection = request.Connection

	return req, nil
}

// cloneResponse deep copies a Response struct.
func cloneResponse(response *Response) (*Response, error) {
	// If response is nil, there's nothing to clone.
	if response == nil {
		return nil, nil
	}

	resp, err := clone[*Response](response)
	if err != nil {
		return nil, fmt.Errorf("unable to clone response: %w", err)
	}

	return resp, nil
}
