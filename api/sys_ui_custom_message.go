// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (c *Sys) ListUICustomMessages(req UICustomMessageListRequest) (*Secret, error) {
	return c.ListUICustomMessagesWithContext(context.Background(), req)
}

func (c *Sys) ListUICustomMessagesWithContext(ctx context.Context, req UICustomMessageListRequest) (*Secret, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest("LIST", "/v1/sys/config/ui/custom-messages/")
	if req.Active != nil {
		r.Params.Add("active", strconv.FormatBool(*req.Active))
	}
	if req.Authenticated != nil {
		r.Params.Add("authenticated", strconv.FormatBool(*req.Authenticated))
	}
	if req.Type != nil {
		r.Params.Add("type", *req.Type)
	}

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	return secret, nil
}

func (c *Sys) CreateUICustomMessage(req UICustomMessageRequest) (*Secret, error) {
	return c.CreateUICustomMessageWithContext(context.Background(), req)
}

func (c *Sys) CreateUICustomMessageWithContext(ctx context.Context, req UICustomMessageRequest) (*Secret, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodPost, "/v1/sys/config/ui/custom-messages")
	if err := r.SetJSONBody(&req); err != nil {
		return nil, fmt.Errorf("error encoding request body to json: %w", err)
	}

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("error sending request to server: %w", err)
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not parse secret from server response: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	return secret, nil
}

func (c *Sys) ReadUICustomMessage(id string) (*Secret, error) {
	return c.ReadUICustomMessageWithContext(context.Background(), id)
}

func (c *Sys) ReadUICustomMessageWithContext(ctx context.Context, id string) (*Secret, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodGet, fmt.Sprintf("/v1/sys/config/ui/custom-messages/%s", id))

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("error sending request to server: %w", err)
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not parse secret from server response: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	return secret, nil
}

func (c *Sys) UpdateUICustomMessage(id string, req UICustomMessageRequest) error {
	return c.UpdateUICustomMessageWithContext(context.Background(), id, req)
}

func (c *Sys) UpdateUICustomMessageWithContext(ctx context.Context, id string, req UICustomMessageRequest) error {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodPost, fmt.Sprintf("/v1/sys/config/ui/custom-messages/%s", id))
	if err := r.SetJSONBody(&req); err != nil {
		return fmt.Errorf("error encoding request body to json: %w", err)
	}

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err != nil {
		return fmt.Errorf("error sending request to server: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (c *Sys) DeleteUICustomMessage(id string) error {
	return c.DeletePolicyWithContext(context.Background(), id)
}

func (c *Sys) DeleteUICustomMessageWithContext(ctx context.Context, id string) error {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/sys/config/ui/custom-messages/%s", id))

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err != nil {
		return fmt.Errorf("error sending request to server: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

type UICustomMessageListRequest struct {
	Authenticated *bool
	Type          *string
	Active        *bool
}

func (r *UICustomMessageListRequest) WithAuthenticated(value bool) *UICustomMessageListRequest {
	r.Authenticated = &value

	return r
}

func (r *UICustomMessageListRequest) WithType(value string) *UICustomMessageListRequest {
	r.Type = &value

	return r
}

func (r *UICustomMessageListRequest) WithActive(value bool) *UICustomMessageListRequest {
	r.Active = &value

	return r
}

type UICustomMessageRequest struct {
	Title         string               `json:"title"`
	Message       string               `json:"message"`
	Authenticated bool                 `json:"authenticated"`
	Type          string               `json:"type"`
	StartTime     string               `json:"start_time"`
	EndTime       string               `json:"end_time,omitempty"`
	Link          *uiCustomMessageLink `json:"link,omitempty"`
	Options       map[string]any       `json:"options,omitempty"`
}

func (r *UICustomMessageRequest) WithLink(title, href string) *UICustomMessageRequest {
	r.Link = &uiCustomMessageLink{
		Title: title,
		Href:  href,
	}

	return r
}

type uiCustomMessageLink struct {
	Title string
	Href  string
}

func (l uiCustomMessageLink) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)

	m[l.Title] = l.Href

	return json.Marshal(m)
}

func (l *uiCustomMessageLink) UnmarshalJSON(b []byte) error {
	m := make(map[string]string)

	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	for k, v := range m {
		l.Title = k
		l.Href = v
		break
	}

	return nil
}
