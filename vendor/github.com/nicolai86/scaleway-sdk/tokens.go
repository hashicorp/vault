package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Token represents a  Token
type Token struct {
	UserID             string `json:"user_id"`
	Description        string `json:"description,omitempty"`
	Roles              Role   `json:"roles"`
	Expires            string `json:"expires"`
	InheritsUsersPerms bool   `json:"inherits_user_perms"`
	ID                 string `json:"id"`
}

// Role represents a  Token UserId Role
type Role struct {
	Organization Organization `json:"organization,omitempty"`
	Role         string       `json:"role,omitempty"`
}

type getTokenResponse struct {
	Token Token `json:"token"`
}

type getTokensResponse struct {
	Tokens []Token `json:"tokens"`
}

func (s *API) GetTokens() ([]Token, error) {
	query := url.Values{}
	// TODO per_page=20&page=2
	resp, err := s.GetResponsePaginate(s.computeAPI, "tokens", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var token getTokensResponse

	if err = json.Unmarshal(body, &token); err != nil {
		return nil, err
	}
	return token.Tokens, nil
}

type CreateTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Expires  bool   `json:"expires"`
}

func (s *API) CreateToken(req *CreateTokenRequest) (*Token, error) {
	resp, err := s.PostResponse(AccountAPI, "tokens", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return nil, err
	}
	var data getTokenResponse

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data.Token, nil
}

type UpdateTokenRequest struct {
	Description string `json:"description,omitempty"`
	Expires     bool   `json:"expires"`
	ID          string `json:"-"`
}

func (s *API) UpdateToken(req *UpdateTokenRequest) (*Token, error) {
	resp, err := s.PatchResponse(AccountAPI, fmt.Sprintf("tokens/%s", req.ID), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var data getTokenResponse

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data.Token, nil
}

func (s *API) GetToken(id string) (*Token, error) {
	query := url.Values{}
	resp, err := s.GetResponsePaginate(AccountAPI, fmt.Sprintf("tokens/%s", id), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var data getTokenResponse

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data.Token, nil
}

func (s *API) DeleteToken(id string) error {
	resp, err := s.DeleteResponse(AccountAPI, fmt.Sprintf("tokens/%s", id))
	if err != nil {
		return err
	}

	if _, err = s.handleHTTPError([]int{http.StatusNoContent}, resp); err != nil {
		return err
	}
	return nil
}
