package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Quota represents a map of quota (name, value)
type Quotas map[string]int

// GetQuotas represents the response of GET /organizations/{orga_id}/quotas
type GetQuotas struct {
	Quotas Quotas `json:"quotas"`
}

// GetQuotas returns a GetQuotas
func (s *API) GetQuotas() (Quotas, error) {
	resp, err := s.GetResponsePaginate(AccountAPI, fmt.Sprintf("organizations/%s/quotas", s.Organization), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var quotas GetQuotas

	if err = json.Unmarshal(body, &quotas); err != nil {
		return nil, err
	}
	return quotas.Quotas, nil
}
