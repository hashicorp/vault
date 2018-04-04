package iamutil

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/errwrap"
	"google.golang.org/api/gensupport"
	"google.golang.org/api/googleapi"
)

type IamHandle struct {
	c         *http.Client
	userAgent string
}

func GetIamHandle(client *http.Client, userAgent string) *IamHandle {
	return &IamHandle{
		c:         client,
		userAgent: userAgent,
	}
}

func (h *IamHandle) GetIamPolicy(ctx context.Context, r IamResource) (*Policy, error) {
	req, err := r.GetIamPolicyRequest()
	if err != nil {
		return nil, errwrap.Wrapf("unable to construct GetIamPolicy request: {{err}}", err)
	}
	var p Policy
	if err := h.doRequest(ctx, req, &p); err != nil {
		return nil, errwrap.Wrapf("unable to get policy: {{err}}", err)
	}
	return &p, nil
}

func (h *IamHandle) SetIamPolicy(ctx context.Context, r IamResource, p *Policy) (*Policy, error) {
	req, err := r.SetIamPolicyRequest(p)
	if err != nil {
		return nil, errwrap.Wrapf("unable to construct SetIamPolicy request: {{err}}", err)
	}
	var out Policy
	if err := h.doRequest(ctx, req, &out); err != nil {
		return nil, errwrap.Wrapf("unable to set policy: {{err}}", err)
	}
	return &out, nil
}

func (h *IamHandle) doRequest(ctx context.Context, req *http.Request, out interface{}) error {
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	if h.userAgent != "" {
		req.Header.Set("User-Agent", h.userAgent)
	}

	resp, err := gensupport.SendRequest(ctx, h.c, req)
	defer googleapi.CloseBody(resp)

	if resp != nil && resp.StatusCode == http.StatusNotModified {
		return &googleapi.Error{
			Code:   resp.StatusCode,
			Header: resp.Header,
		}
	}
	if err != nil {
		return err
	}

	if err := googleapi.CheckResponse(resp); err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return errwrap.Wrapf("unable to decode JSON resp to output interface: {{err}}", err)
	}
	return nil
}
