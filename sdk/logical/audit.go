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

// LogInputBexpr is used for evaluating boolean expressions with go-bexpr.
type LogInputBexpr struct {
	MountPoint string `bexpr:"mount_point"`
	MountType  string `bexpr:"mount_type"`
	Namespace  string `bexpr:"namespace"`
	Operation  string `bexpr:"operation"`
	Path       string `bexpr:"path"`
}

// BexprDatum returns values from a LogInput formatted for use in evaluating go-bexpr boolean expressions.
// The namespace should be supplied from the current request's context.
func (l *LogInput) BexprDatum(namespace string) *LogInputBexpr {
	var mountPoint string
	var mountType string
	var operation string
	var path string

	if l.Request != nil {
		mountPoint = l.Request.MountPoint
		mountType = l.Request.MountType
		operation = string(l.Request.Operation)
		path = l.Request.Path
	}

	return &LogInputBexpr{
		MountPoint: mountPoint,
		MountType:  mountType,
		Namespace:  namespace,
		Operation:  operation,
		Path:       path,
	}
}

// Clone will attempt to create a deep copy (almost) of the LogInput.
// If the LogInput type or any of the subtypes referenced by LogInput fields are
// changed, then the Clone methods will need to be updated.
// NOTE: Does not deep clone the LogInput.OuterError field as it represents an
// error interface.
// NOTE: LogInput.Request.Connection (at the time of writing) is also not deep-copied
// and remains a pointer, see Request.Clone for more information.
func (l *LogInput) Clone() (*LogInput, error) {
	// Clone Auth
	auth, err := cloneAuth(l.Auth)
	if err != nil {
		return nil, err
	}

	// Clone Request
	var req *Request
	if l.Request != nil {
		req, err = l.Request.Clone()
		if err != nil {
			return nil, err
		}
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
