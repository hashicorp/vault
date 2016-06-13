package logical

import "time"

// This logic was pulled from the http package so that it can be used for
// encoding wrapped responses as well. It simply translates the logical request
// to an http response, with the values we want and omitting the values we
// don't.
func SanitizeResponse(input *Response) *HTTPResponse {
	logicalResp := &HTTPResponse{
		Data:     input.Data,
		Warnings: input.Warnings(),
	}

	if input.Secret != nil {
		logicalResp.LeaseID = input.Secret.LeaseID
		logicalResp.Renewable = input.Secret.Renewable
		logicalResp.LeaseDuration = int(input.Secret.TTL.Seconds())
	}

	// If we have authentication information, then
	// set up the result structure.
	if input.Auth != nil {
		logicalResp.Auth = &HTTPAuth{
			ClientToken:   input.Auth.ClientToken,
			Accessor:      input.Auth.Accessor,
			Policies:      input.Auth.Policies,
			Metadata:      input.Auth.Metadata,
			LeaseDuration: int(input.Auth.TTL.Seconds()),
			Renewable:     input.Auth.Renewable,
		}
	}

	return logicalResp
}

type HTTPResponse struct {
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
}

type HTTPWrapInfo struct {
	Token           string    `json:"token"`
	TTL             int       `json:"ttl"`
	CreationTime    time.Time `json:"creation_time"`
	WrappedAccessor string    `json:"wrapped_accessor,omitempty"`
}
