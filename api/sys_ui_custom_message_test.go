// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var messageBase64 string = base64.StdEncoding.EncodeToString([]byte("message"))

// TestUICustomMessageJsonMarshalling verifies that json marshalling (struct to
// json) works with the uiCustomMessageRequest type.
func TestUICustomMessageJsonMarshalling(t *testing.T) {
	for _, testcase := range []struct {
		name         string
		request      UICustomMessageRequest
		expectedJSON string
	}{
		{
			name: "no-link-no-options",
			request: UICustomMessageRequest{
				Title:         "title",
				Message:       messageBase64,
				StartTime:     "2024-01-01T00:00:00.000Z",
				EndTime:       "",
				Type:          "banner",
				Authenticated: true,
			},
			expectedJSON: fmt.Sprintf(`{"title":"title","message":"%s","authenticated":true,"type":"banner","start_time":"2024-01-01T00:00:00.000Z"}`, messageBase64),
		},
		{
			name: "link-no-options",
			request: UICustomMessageRequest{
				Title:         "title",
				Message:       messageBase64,
				StartTime:     "2024-01-01T00:00:00.000Z",
				EndTime:       "",
				Type:          "modal",
				Authenticated: false,
				Link: &uiCustomMessageLink{
					Title: "Click here",
					Href:  "https://www.example.org",
				},
			},
			expectedJSON: fmt.Sprintf(`{"title":"title","message":"%s","authenticated":false,"type":"modal","start_time":"2024-01-01T00:00:00.000Z","link":{"Click here":"https://www.example.org"}}`, messageBase64),
		},
		{
			name: "no-link-options",
			request: UICustomMessageRequest{
				Title:         "title",
				Message:       messageBase64,
				StartTime:     "2024-01-01T00:00:00.000Z",
				EndTime:       "",
				Authenticated: true,
				Type:          "banner",
				Options: map[string]any{
					"key": "value",
				},
			},
			expectedJSON: fmt.Sprintf(`{"title":"title","message":"%s","authenticated":true,"type":"banner","start_time":"2024-01-01T00:00:00.000Z","options":{"key":"value"}}`, messageBase64),
		},
		{
			name: "link-and-options",
			request: UICustomMessageRequest{
				Title:         "title",
				Message:       messageBase64,
				StartTime:     "2024-01-01T00:00:00.000Z",
				EndTime:       "",
				Authenticated: true,
				Type:          "banner",
				Link: &uiCustomMessageLink{
					Title: "Click here",
					Href:  "https://www.example.org",
				},
				Options: map[string]any{
					"key": "value",
				},
			},
			expectedJSON: fmt.Sprintf(`{"title":"title","message":"%s","authenticated":true,"type":"banner","start_time":"2024-01-01T00:00:00.000Z","link":{"Click here":"https://www.example.org"},"options":{"key":"value"}}`, messageBase64),
		},
	} {
		tc := testcase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bytes, err := json.Marshal(&tc.request)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedJSON, string(bytes))
		})
	}
}

// TestUICustomMessageJsonUnmarshal verifies that json unmarshalling (json to
// struct) works with the uiCustomMessageRequest type.
func TestUICustomMessageJsonUnmarshal(t *testing.T) {
	for _, testcase := range []struct {
		name             string
		encodedBytes     string
		linkAssertion    func(assert.TestingT, any, ...any) bool
		checkLink        bool
		optionsAssertion func(assert.TestingT, any, ...any) bool
		checkOptions     bool
	}{
		{
			name:             "no-link-no-options",
			encodedBytes:     fmt.Sprintf(`{"title":"title","message":"%s","authenticated":false,"type":"modal","start_time":"2024-01-01T00:00:00.000Z"}`, messageBase64),
			linkAssertion:    assert.Nil,
			optionsAssertion: assert.Nil,
		},
		{
			name:             "link-no-options",
			encodedBytes:     fmt.Sprintf(`{"title":"title","message":"%s","authenticated":false,"type":"modal","start_time":"2024-01-01T00:00:00.000Z","link":{"Click here":"https://www.example.org"}}`, messageBase64),
			linkAssertion:    assert.NotNil,
			checkLink:        true,
			optionsAssertion: assert.Nil,
		},
		{
			name:             "no-link-options",
			encodedBytes:     fmt.Sprintf(`{"title":"title","message":"%s","authenticated":false,"type":"modal","start_time":"2024-01-01T00:00:00.000Z","options":{"key":"value"}}`, messageBase64),
			linkAssertion:    assert.Nil,
			optionsAssertion: assert.NotNil,
			checkOptions:     true,
		},
		{
			name:             "link-and-options",
			encodedBytes:     fmt.Sprintf(`{"title":"title","message":"%s","authenticated":false,"type":"modal","start_time":"2024-01-01T00:00:00.000Z","link":{"Click here":"https://www.example.org"},"options":{"key":"value"}}`, messageBase64),
			linkAssertion:    assert.NotNil,
			checkLink:        true,
			optionsAssertion: assert.NotNil,
			checkOptions:     true,
		},
	} {
		tc := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			var request UICustomMessageRequest

			err := json.Unmarshal([]byte(tc.encodedBytes), &request)
			assert.NoError(t, err)
			tc.linkAssertion(t, request.Link)
			tc.optionsAssertion(t, request.Options)

			if tc.checkLink {
				assert.Equal(t, "Click here", request.Link.Title)
				assert.Equal(t, "https://www.example.org", request.Link.Href)
			}

			if tc.checkOptions {
				assert.Contains(t, request.Options, "key")
			}
		})
	}
}

// TestUICustomMessageListRequestOptions verifies the correct behaviour of all
// of the With... methods of the UICustomMessageListRequest.
func TestUICustomMessageListRequestOptions(t *testing.T) {
	request := &UICustomMessageListRequest{}
	assert.Nil(t, request.Active)
	assert.Nil(t, request.Authenticated)
	assert.Nil(t, request.Type)

	request = (&UICustomMessageListRequest{}).WithActive(true)
	assert.NotNil(t, request.Active)
	assert.True(t, *request.Active)

	request = (&UICustomMessageListRequest{}).WithActive(false)
	assert.NotNil(t, request.Active)
	assert.False(t, *request.Active)

	request = (&UICustomMessageListRequest{}).WithAuthenticated(true)
	assert.NotNil(t, request.Authenticated)
	assert.True(t, *request.Authenticated)

	request = (&UICustomMessageListRequest{}).WithAuthenticated(false)
	assert.NotNil(t, request.Authenticated)
	assert.False(t, *request.Authenticated)

	request = (&UICustomMessageListRequest{}).WithType("banner")
	assert.NotNil(t, request.Type)
	assert.Equal(t, "banner", *request.Type)

	request = (&UICustomMessageListRequest{}).WithType("modal")
	assert.NotNil(t, request.Type)
	assert.Equal(t, "modal", *request.Type)
}
