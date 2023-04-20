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
	Forwarding          *ForwardingInfo
}

// ForwardingInfo should be used to describe the hosts a request is forwarded 'from' and 'to'.
type ForwardingInfo struct {
	From string
	To   string
}

type MarshalOptions struct {
	ValueHasher func(string) string
}

type OptMarshaler interface {
	MarshalJSONWithOptions(*MarshalOptions) ([]byte, error)
}

// IsPresent can be used to determine whether 'from' and/or 'to' forwarding information is present.
func (f *ForwardingInfo) IsPresent() bool {
	return len(f.From) > 0 || len(f.To) > 0
}
