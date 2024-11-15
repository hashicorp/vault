// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iamutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/errwrap"
	"google.golang.org/api/googleapi"
)

type ApiHandle struct {
	c         *http.Client
	userAgent string
}

func GetApiHandle(client *http.Client, userAgent string) *ApiHandle {
	return &ApiHandle{
		c:         client,
		userAgent: userAgent,
	}
}

func (h *ApiHandle) DoGetRequest(ctx context.Context, r Resource, out interface{}) (err error) {
	config := r.GetConfig()
	req, err := constructRequest(r, &config.GetMethod, nil)
	if err != nil {
		return errwrap.Wrapf("Unable to construct Get request: {{err}}", err)
	}
	return h.doRequest(ctx, req, out)
}

func (h *ApiHandle) DoSetRequest(ctx context.Context, r Resource, data io.Reader, out interface{}) error {
	config := r.GetConfig()
	req, err := constructRequest(r, &config.SetMethod, data)
	if err != nil {
		return errwrap.Wrapf("Unable to construct Set request: {{err}}", err)
	}
	return h.doRequest(ctx, req, out)
}

func (h *ApiHandle) doRequest(ctx context.Context, req *http.Request, out interface{}) error {
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	if h.userAgent != "" {
		req.Header.Set("User-Agent", h.userAgent)
	}

	resp, err := h.c.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer googleapi.CloseBody(resp)

	if err := googleapi.CheckResponse(resp); err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return errwrap.Wrapf("unable to decode JSON resp to output interface: {{err}}", err)
	}
	return nil
}

func constructRequest(r Resource, restMethod *RestMethod, data io.Reader) (*http.Request, error) {
	config := r.GetConfig()
	if data == nil && config != nil && config.Service == "cloudresourcemanager" {
		// In order to support Resource Manager policies with conditional bindings,
		// we need to request the policy version of 3. This request parameter is backwards compatible
		// and will return version 1 policies if they are not yet updated to version 3.
		requestPolicyVersion3 := `{"options": {"requestedPolicyVersion": 3}}`
		data = strings.NewReader(requestPolicyVersion3)
	}
	req, err := http.NewRequest(
		restMethod.HttpMethod,
		googleapi.ResolveRelative(restMethod.BaseURL, restMethod.Path),
		data)
	if err != nil {
		return nil, err
	}

	if req.Header == nil {
		req.Header = make(http.Header)
	}
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	relId := r.GetRelativeId()
	replacementMap := make(map[string]string)

	if strings.Contains(restMethod.Path, "{+resource}") {
		// +resource is used to represent full relative resource name
		if len(config.Parameters) == 1 && config.Parameters[0] == "resource" {
			relName := ""
			tkns := strings.Split(config.TypeKey, "/")
			for _, colId := range tkns {
				if colName, ok := relId.IdTuples[colId]; ok {
					relName += fmt.Sprintf("%s/%s/", colId, colName)
				}
			}
			replacementMap["resource"] = strings.Trim(relName, "/")
		}
	} else {
		for colId, resId := range relId.IdTuples {
			rId, ok := config.CollectionReplacementKeys[colId]
			if !ok {
				return nil, fmt.Errorf("expected value for collection id %s", colId)
			}
			replacementMap[rId] = resId
		}
	}

	googleapi.Expand(req.URL, replacementMap)
	return req, nil
}
