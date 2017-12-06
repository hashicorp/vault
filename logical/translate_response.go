package logical

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// This logic was pulled from the http package so that it can be used for
// encoding wrapped responses as well. It simply translates the logical
// response to an http response, with the values we want and omitting the
// values we don't.
func LogicalResponseToHTTPResponse(input *Response) *HTTPResponse {
	httpResp := &HTTPResponse{
		Data:     input.Data,
		Warnings: input.Warnings,
	}

	if input.Secret != nil {
		httpResp.LeaseID = input.Secret.LeaseID
		httpResp.Renewable = input.Secret.Renewable
		httpResp.LeaseDuration = int(input.Secret.TTL.Seconds())
	}

	// If we have authentication information, then
	// set up the result structure.
	if input.Auth != nil {
		httpResp.Auth = &HTTPAuth{
			ClientToken:   input.Auth.ClientToken,
			Accessor:      input.Auth.Accessor,
			Policies:      input.Auth.Policies,
			Metadata:      input.Auth.Metadata,
			LeaseDuration: int(input.Auth.TTL.Seconds()),
			Renewable:     input.Auth.Renewable,
			EntityID:      input.Auth.EntityID,
		}
	}

	return httpResp
}

func HTTPResponseToLogicalResponse(input *HTTPResponse) *Response {
	logicalResp := &Response{
		Data:     input.Data,
		Warnings: input.Warnings,
	}

	if input.LeaseID != "" {
		logicalResp.Secret = &Secret{
			LeaseID: input.LeaseID,
		}
		logicalResp.Secret.Renewable = input.Renewable
		logicalResp.Secret.TTL = time.Second * time.Duration(input.LeaseDuration)
	}

	if input.Auth != nil {
		logicalResp.Auth = &Auth{
			ClientToken: input.Auth.ClientToken,
			Accessor:    input.Auth.Accessor,
			Policies:    input.Auth.Policies,
			Metadata:    input.Auth.Metadata,
			EntityID:    input.Auth.EntityID,
		}
		logicalResp.Auth.Renewable = input.Auth.Renewable
		logicalResp.Auth.TTL = time.Second * time.Duration(input.Auth.LeaseDuration)
	}

	return logicalResp
}

type HTTPResponse struct {
	RequestID     string                 `json:"request_id"`
	LeaseID       string                 `json:"lease_id"`
	Renewable     bool                   `json:"renewable"`
	LeaseDuration int                    `json:"lease_duration"`
	Data          map[string]interface{} `json:"data"`
	WrapInfo      *HTTPWrapInfo          `json:"wrap_info"`
	Warnings      []string               `json:"warnings"`
	Auth          *HTTPAuth              `json:"auth"`
}

type HTTPAuth struct {
	ClientToken   string            `json:"client_token"`
	Accessor      string            `json:"accessor"`
	Policies      []string          `json:"policies"`
	Metadata      map[string]string `json:"metadata"`
	LeaseDuration int               `json:"lease_duration"`
	Renewable     bool              `json:"renewable"`
	EntityID      string            `json:"entity_id"`
}

type HTTPWrapInfo struct {
	Token           string `json:"token"`
	Accessor        string `json:"accessor"`
	TTL             int    `json:"ttl"`
	CreationTime    string `json:"creation_time"`
	CreationPath    string `json:"creation_path"`
	WrappedAccessor string `json:"wrapped_accessor,omitempty"`
}

type HTTPSysInjector struct {
	Response *HTTPResponse
}

func (h HTTPSysInjector) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(h.Response)
	if err != nil {
		return nil, err
	}

	// Fast path no data or empty data
	if h.Response.Data == nil || len(h.Response.Data) == 0 {
		return j, nil
	}

	// Marshaling a response will always be a JSON object, meaning it will
	// always start with '{', so we hijack this to prepend necessary values

	// Make a guess at the capacity, and write the object opener
	buf := bytes.NewBuffer(make([]byte, 0, len(j)*2))
	buf.WriteRune('{')

	for k, v := range h.Response.Data {
		// Marshal each key/value individually
		mk, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		mv, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		// Write into the final buffer. We'll never have a valid response
		// without any fields so we can unconditionally add a comma after each.
		buf.WriteString(fmt.Sprintf("%s: %s, ", mk, mv))
	}

	// Add the rest, without the first '{'
	buf.Write(j[1:])

	return buf.Bytes(), nil
}
