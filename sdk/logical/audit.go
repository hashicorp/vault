// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

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
