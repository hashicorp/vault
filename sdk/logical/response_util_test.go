// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/consts"
)

func TestResponseUtil_RespondErrorCommon_basic(t *testing.T) {
	testCases := []struct {
		title          string
		req            *Request
		resp           *Response
		respErr        error
		expectedStatus int
		expectedErr    error
	}{
		{
			title:          "Throttled, no error",
			respErr:        ErrUpstreamRateLimited,
			resp:           &Response{},
			expectedStatus: 502,
		},
		{
			title:   "Throttled, with error",
			respErr: ErrUpstreamRateLimited,
			resp: &Response{
				Data: map[string]interface{}{
					"error": "rate limited",
				},
			},
			expectedStatus: 502,
		},
		{
			title: "Read not found",
			req: &Request{
				Operation: ReadOperation,
			},
			respErr:        nil,
			expectedStatus: 404,
		},
		{
			title: "Header not found",
			req: &Request{
				Operation: HeaderOperation,
			},
			respErr:        nil,
			expectedStatus: 404,
		},
		{
			title: "List with response and no keys",
			req: &Request{
				Operation: ListOperation,
			},
			resp:           &Response{},
			respErr:        nil,
			expectedStatus: 404,
		},
		{
			title: "List with response and keys",
			req: &Request{
				Operation: ListOperation,
			},
			resp: &Response{
				Data: map[string]interface{}{
					"keys": []string{"some", "things", "here"},
				},
			},
			respErr:        nil,
			expectedStatus: 0,
		},
		{
			title:   "Invalid Credentials error ",
			respErr: ErrInvalidCredentials,
			resp: &Response{
				Data: map[string]interface{}{
					"error": "error due to wrong credentials",
				},
			},
			expectedErr:    errors.New("error due to wrong credentials"),
			expectedStatus: 400,
		},
		{
			title:   "Overloaded error",
			respErr: consts.ErrOverloaded,
			resp: &Response{
				Data: map[string]interface{}{
					"error": "overloaded, try again later",
				},
			},
			expectedErr:    consts.ErrOverloaded,
			expectedStatus: 503,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var status int
			var err, respErr error
			if tc.respErr != nil {
				respErr = tc.respErr
			}
			status, err = RespondErrorCommon(tc.req, tc.resp, respErr)
			if status != tc.expectedStatus {
				t.Fatalf("Expected (%d) status code, got (%d)", tc.expectedStatus, status)
			}
			if tc.expectedErr != nil {
				if !strings.Contains(tc.expectedErr.Error(), err.Error()) {
					t.Fatalf("Expected error to contain:\n%s\n\ngot:\n%s\n", tc.expectedErr, err)
				}
			}
		})
	}
}
