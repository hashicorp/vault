package packngo

import "path"

type FabricServiceTokenType string

const (
	fabricServiceTokenBasePath                        = "/fabric-service-tokens"
	FabricServiceTokenASide    FabricServiceTokenType = "a_side"
	FabricServiceTokenZSide    FabricServiceTokenType = "z_side"
)

// FabricServiceToken represents an Equinix Metal metro
type FabricServiceToken struct {
	*Href            `json:",inline"`
	ExpiresAt        *Timestamp             `json:"expires_at,omitempty"`
	ID               string                 `json:"id"`
	MaxAllowedSpeed  uint64                 `json:"max_allowed_speed,omitempty"`
	Role             ConnectionPortRole     `json:"role,omitempty"`
	ServiceTokenType FabricServiceTokenType `json:"service_token_type,omitempty"`
	State            string                 `json:"state,omitempty"`
}

func (f FabricServiceToken) String() string {
	return Stringify(f)
}

// FabricServiceTokenServiceOp implements FabricServiceTokenService
type FabricServiceTokenServiceOp struct {
	client *Client
}

func (s *FabricServiceTokenServiceOp) Get(id string, opts *GetOptions) (*FabricServiceToken, *Response, error) {
	endpointPath := path.Join(fabricServiceTokenBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	fst := new(FabricServiceToken)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, fst)
	if err != nil {
		return nil, resp, err
	}
	return fst, resp, err
}
