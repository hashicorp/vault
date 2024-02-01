// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	uicustommessages "github.com/hashicorp/vault/vault/ui_custom_messages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHandleListCustomMessages verifies the proper functioning of the
// (*SystemBackend).handleListCustomMessages method. This test focuses on proper
// parsing of the request parameters (framework.FieldData) and when errors
// occur in the underlying logical.Storage.
func TestHandleListCustomMessages(t *testing.T) {
	startTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339Nano)
	endTime := time.Now().Add(time.Hour).Format(time.RFC3339Nano)

	storageEntry := &logical.StorageEntry{
		Key:   "sys/config/ui/custom-messages",
		Value: []byte(fmt.Sprintf(`{"messages":{"000":{"id":"000","title":"title","message":"message","type":"%s","authenticated":true,"start_time":"%s","end_time":"%s"}}}`, uicustommessages.BannerMessageType, startTime, endTime)),
	}

	storage := &logical.InmemStorage{}
	storage.Put(context.Background(), storageEntry)

	backend := &SystemBackend{
		Core: &Core{
			customMessageManager: uicustommessages.NewManager(storage),
		},
	}

	nsCtx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	fieldSchemas := map[string]*framework.FieldSchema{
		"authenticated": {
			Type: framework.TypeBool,
		},
		"type": {
			Type: framework.TypeString,
		},
		"active": {
			Type: framework.TypeBool,
		},
	}

	testcases := []struct {
		name              string
		fieldRaw          map[string]any
		expectKeysInData  bool
		expectErrorInData bool
	}{
		{
			name:             "no-filter-parameters",
			fieldRaw:         map[string]any{},
			expectKeysInData: true,
		},
		{
			name: "authenticated-false",
			fieldRaw: map[string]any{
				"authenticated": "false",
			},
		},
		{
			name: "authenticated-true",
			fieldRaw: map[string]any{
				"authenticated": "true",
			},
			expectKeysInData: true,
		},
		{
			name: "authenticated-invalid-value",
			fieldRaw: map[string]any{
				"authenticated": "fred",
			},
			expectErrorInData: true,
		},
		{
			name: "type-banner",
			fieldRaw: map[string]any{
				"type": uicustommessages.BannerMessageType,
			},
			expectKeysInData: true,
		},
		{
			name: "type-modal",
			fieldRaw: map[string]any{
				"type": uicustommessages.ModalMessageType,
			},
		},
		{
			name: "type-unrecognized-value",
			fieldRaw: map[string]any{
				"type": "fred",
			},
			expectKeysInData: true,
		},
		{
			name: "type-invalid-value",
			fieldRaw: map[string]any{
				"type": []int{0},
			},
			expectErrorInData: true,
		},
		{
			name: "active-false",
			fieldRaw: map[string]any{
				"active": "false",
			},
		},
		{
			name: "active-true",
			fieldRaw: map[string]any{
				"active": "true",
			},
			expectKeysInData: true,
		},
		{
			name: "active-invalid-value",
			fieldRaw: map[string]any{
				"active": "fred",
			},
			expectErrorInData: true,
		},
	}

	for _, testcase := range testcases {
		resp, err := backend.handleListCustomMessages(nsCtx, &logical.Request{}, &framework.FieldData{Schema: fieldSchemas, Raw: testcase.fieldRaw})

		assert.NoError(t, err, testcase.name)
		assert.NotNil(t, resp, testcase.name)
		if testcase.expectKeysInData {
			assert.Contains(t, resp.Data, "keys", testcase.name)
			assert.Equal(t, 1, len(resp.Data["keys"].([]string)), testcase.name)
			assert.Contains(t, resp.Data, "key_info", testcase.name)
			assert.IsType(t, map[string]any{}, resp.Data["key_info"], testcase.name)
			assert.Contains(t, resp.Data["key_info"].(map[string]any), "000", testcase.name)
			assert.Contains(t, resp.Data["key_info"].(map[string]any)["000"], "title", testcase.name)
			assert.Contains(t, resp.Data["key_info"].(map[string]any)["000"], "type", testcase.name)
			assert.Contains(t, resp.Data["key_info"].(map[string]any)["000"], "authenticated", testcase.name)
			assert.Contains(t, resp.Data["key_info"].(map[string]any)["000"], "start_time", testcase.name)
			assert.Contains(t, resp.Data["key_info"].(map[string]any)["000"], "end_time", testcase.name)
			assert.Contains(t, resp.Data["key_info"].(map[string]any)["000"], "active", testcase.name)
		} else {
			assert.NotContains(t, resp.Data, "keys", testcase.name)
		}

		if testcase.expectErrorInData {
			assert.Contains(t, resp.Data, "error", testcase.name)
		} else {
			assert.NotContains(t, resp.Data, "error", testcase.name)
		}
	}

	// Finally, test when the underlying storage returns an error
	backend.Core.customMessageManager = uicustommessages.NewManager(&testingStorage{
		getFails: true,
	})

	resp, err := backend.handleListCustomMessages(nsCtx, &logical.Request{}, &framework.FieldData{
		Schema: fieldSchemas,
		Raw:    map[string]any{},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "error")
}

// TestHandleCreateCustomMessage verifies the proper functioning of the
// (*SystemBackend).handleCreateCustomMessage method. The test focuses on
// missing and invalid request parameter (framework.FieldData) and errors in the
// underlying logical.Storage.
func TestHandleCreateCustomMessage(t *testing.T) {
	// Setup a system backend
	backend := &SystemBackend{
		Core: &Core{
			customMessageManager: uicustommessages.NewManager(&logical.InmemStorage{}),
		},
	}

	nsCtx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	fieldSchemas := map[string]*framework.FieldSchema{
		"title": {
			Type:     framework.TypeString,
			Required: true,
		},
		"message": {
			Type:     framework.TypeString,
			Required: true,
		},
		"authenticated": {
			Type:    framework.TypeBool,
			Default: true,
		},
		"type": {
			Type:    framework.TypeString,
			Default: uicustommessages.BannerMessageType,
		},
		"start_time": {
			Type:     framework.TypeTime,
			Required: true,
		},
		"end_time": {
			Type: framework.TypeTime,
		},
		"link": {
			Type: framework.TypeMap,
		},
		"options": {
			Type: framework.TypeMap,
		},
	}

	// the standard map of request parameters containing all required parameters
	// with valid values. The test cases will make a copy of this map and modify
	// it to test different conditions.
	fieldRaw := map[string]any{
		"title":      "title",
		"message":    base64.StdEncoding.EncodeToString([]byte("message")),
		"start_time": "2023-01-01T00:00:00Z",
	}

	testcases := []struct {
		name           string
		fieldRawDelete []string
		fieldRawUpdate map[string]any
		errorExpected  bool
	}{
		{
			name:           "title-parameter-missing",
			fieldRawDelete: []string{"title"},
			errorExpected:  true,
		},
		{
			name: "title-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"title": []bool{},
			},
			errorExpected: true,
		},
		{
			name: "authenticated-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"authenticated": "abc",
			},
			errorExpected: true,
		},
		{
			name: "type-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"type": []int{},
			},
			errorExpected: true,
		},
		{
			name:           "message-parameter-missing",
			fieldRawDelete: []string{"message"},
			errorExpected:  true,
		},
		{
			name: "message-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"message": map[int]string{},
			},
			errorExpected: true,
		},
		{
			name: "message-parameter-invalid-not-base64",
			fieldRawUpdate: map[string]any{
				"message": "The message not base64 encoded.",
			},
			errorExpected: true,
		},
		{
			name:           "start_time-parameter-missing",
			fieldRawDelete: []string{"start_time"},
			errorExpected:  true,
		},
		{
			name: "start_time-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"start_time": "friday",
			},
			errorExpected: true,
		},
		{
			name: "end_time-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"end_time": []int{},
			},
			errorExpected: true,
		},
		{
			name: "link-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"link": "not-a-map",
			},
			errorExpected: true,
		},
		{
			name: "link-parameter-href-invalid",
			fieldRawUpdate: map[string]any{
				"link": map[string]any{
					"click here": []int{},
				},
			},
			errorExpected: true,
		},
		{
			name: "link-parameter-multiple-links",
			fieldRawUpdate: map[string]any{
				"link": map[string]any{
					"click here":   "http://example.org",
					"click here 2": "http://ping.net",
				},
			},
			errorExpected: true,
		},
		{
			name: "options-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"options": "not-a-map",
			},
			errorExpected: true,
		},
		{
			name: "happy-path with defaults",
		},
		{
			name: "happy-path with non-defaults",
			fieldRawUpdate: map[string]any{
				"authenticated": false,
				"type":          uicustommessages.ModalMessageType,
				"end_time":      "2100-01-01T00:00:00Z",
				"options": map[string]any{
					"color": "red",
				},
				"link": map[string]string{
					"Details": "https://server.com/details",
				},
			},
		},
	}

	for _, testcase := range testcases {
		raw := map[string]any{}
		for k, v := range fieldRaw {
			raw[k] = v
		}

		for _, d := range testcase.fieldRawDelete {
			delete(raw, d)
		}

		for k, v := range testcase.fieldRawUpdate {
			raw[k] = v
		}

		resp, err := backend.handleCreateCustomMessages(nsCtx, &logical.Request{}, &framework.FieldData{
			Schema: fieldSchemas,
			Raw:    raw,
		})
		assert.NoError(t, err, testcase.name)
		assert.NotNil(t, resp, testcase.name)
		assert.NotNil(t, resp.Data, testcase.name)

		if testcase.errorExpected {
			assert.Contains(t, resp.Data, "error", testcase.name)

			if len(testcase.fieldRawDelete) > 0 {
				assert.Contains(t, resp.Data["error"], "missing", testcase.name)
				assert.Contains(t, resp.Data["error"], testcase.fieldRawDelete[0], testcase.name)
			}

			if len(testcase.fieldRawUpdate) > 0 {
				assert.Contains(t, resp.Data["error"], "invalid", testcase.name)

				for k := range testcase.fieldRawUpdate {
					assert.Contains(t, resp.Data["error"], k, testcase.name)
				}
			}
		} else {
			assert.Contains(t, resp.Data, "active", testcase.name)
			assert.Contains(t, resp.Data, "authenticated", testcase.name)
			assert.Contains(t, resp.Data, "type", testcase.name)
			assert.Contains(t, resp.Data, "start_time", testcase.name)
			assert.Contains(t, resp.Data, "end_time", testcase.name)
			assert.Contains(t, resp.Data, "id", testcase.name)
			assert.Contains(t, resp.Data, "options", testcase.name)
			assert.Contains(t, resp.Data, "link", testcase.name)
			_, ok := testcase.fieldRawUpdate["authenticated"]
			assert.Equal(t, !ok, resp.Data["authenticated"].(bool), testcase.name)
			if _, ok := testcase.fieldRawUpdate["type"]; !ok {
				assert.Equal(t, resp.Data["type"], uicustommessages.BannerMessageType, testcase.name)
			} else {
				assert.Equal(t, resp.Data["type"], uicustommessages.ModalMessageType, testcase.name)
			}
			if _, ok := testcase.fieldRawUpdate["end_time"]; !ok {
				assert.Nil(t, resp.Data["end_time"], testcase.name)
			} else {
				assert.NotNil(t, resp.Data["end_time"], testcase.name)
			}
			if _, ok := testcase.fieldRawUpdate["link"]; !ok {
				assert.Nil(t, resp.Data["link"], testcase.name)
			} else {
				assert.NotNil(t, resp.Data["link"], testcase.name)
			}
		}
	}

	// Finally, test when the underlying storage returns an error
	backend.Core.customMessageManager = uicustommessages.NewManager(&testingStorage{
		putFails:         true,
		getResponseValue: `{"messages":{}}`,
	})

	resp, err := backend.handleCreateCustomMessages(context.Background(), &logical.Request{}, &framework.FieldData{
		Schema: fieldSchemas,
		Raw:    fieldRaw,
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "error")
}

// TestHandleCreateAndUpdateCustomMessageLinkProperty verifies some particular
// cases in the handleCreateCustomMessage method surrounding the link property
// of the custom message. Because the link property is modeled as a JSON object,
// the possibility of an empty or partially constructed object is provided, so
// this test function verifies that each of these cases is appropriately
// handled.
func TestHandleCreateAndUpdateCustomMessageLinkProperty(t *testing.T) {
	testBackend := &SystemBackend{
		Core: &Core{
			customMessageManager: uicustommessages.NewManager(&logical.InmemStorage{}),
		},
	}

	fieldData := &framework.FieldData{
		Schema: map[string]*framework.FieldSchema{
			"title": {
				Type:     framework.TypeString,
				Required: true,
			},
			"message": {
				Type:     framework.TypeString,
				Required: true,
			},
			"authenticated": {
				Type:    framework.TypeBool,
				Default: true,
			},
			"type": {
				Type:    framework.TypeString,
				Default: "banner",
			},
			"start_time": {
				Type:     framework.TypeTime,
				Required: true,
			},
			"end_time": {
				Type: framework.TypeTime,
			},
			"link": {
				Type: framework.TypeMap,
			},
			"options": {
				Type: framework.TypeMap,
			},
		},
		Raw: map[string]any{
			"title":      "The title",
			"message":    "VGhlIG1lc3NhZ2UK",
			"start_time": "2024-01-01T00:00:00.000Z",
		},
	}

	for _, testcase := range []struct {
		name                  string
		linkField             any
		responseErrorExpected bool
	}{
		// First group of testcases, setup an incomplete link parameter,
		// which result in a CustomMessage with a Link field set to nil
		{
			name:      "nil-link",
			linkField: nil,
		},
		{
			name:      "empty-link-map",
			linkField: map[string]any{},
		},
		{
			name: "empty-key-and-value-link-map",
			linkField: map[string]any{
				"": "",
			},
		},
		{
			name: "empty-value-link-map",
			linkField: map[string]any{
				"Click me": "",
			},
		},
		{
			name: "empty-key-nil-value-link-map",
			linkField: map[string]any{
				"": nil,
			},
		},
		// Second group of testcases, setup an invalid link parameter,
		// which result in an error being returned and the message not being
		// created or updated.
		{
			name: "nil-value-link-map",
			linkField: map[string]any{
				"Click me": nil,
			},
			responseErrorExpected: true,
		},
		{
			name: "invalid-value-link-map",
			linkField: map[string]any{
				"Click me": []int{},
			},
			responseErrorExpected: true,
		},
		{
			name: "empty-key-link-map",
			linkField: map[string]any{
				"": "https://www.example.org",
			},
			responseErrorExpected: true,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			delete(fieldData.Raw, "link")

			// Create a custom message with no link to be used to test the
			// handleUpdateCustomMessage case.
			resp, err := testBackend.handleCreateCustomMessages(namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace), &logical.Request{}, fieldData)
			require.NoError(t, err)
			require.Contains(t, resp.Data, "id")
			updatableMessageId := resp.Data["id"]

			// Update the fieldData to include the testcase
			fieldData.Raw["link"] = testcase.linkField
			fieldData.Schema["id"] = &framework.FieldSchema{
				Type:     framework.TypeString,
				Required: true,
			}
			fieldData.Raw["id"] = updatableMessageId

			// Test the Update operation
			resp, err = testBackend.handleUpdateCustomMessage(namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace), &logical.Request{}, fieldData)
			assert.NoError(t, err, "UpdateCustomMessage")

			if testcase.responseErrorExpected {
				assert.Contains(t, resp.Data, "error", "UpdateCustomMessage")
			} else {
				assert.Nil(t, resp.Data["link"], "UpdateCustomMessage")
			}

			// Reset for testing the Create operation
			delete(fieldData.Schema, "id")
			delete(fieldData.Raw, "id")

			resp, err = testBackend.handleCreateCustomMessages(namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace), &logical.Request{}, fieldData)
			assert.NoError(t, err, "CreateCustomMessage")

			if testcase.responseErrorExpected {
				assert.Contains(t, resp.Data, "error", "CreateCustomMessage")
			} else {
				assert.Nil(t, resp.Data["link"], "CreateCustomMessage")
			}
		})
	}
}

// TestHandleReadCustomMessage verifies the proper functioning of the
// (*SystemBackend).handleReadCustomMessage method. The tests focus on missing
// or invalid request parameters as well as reading existing and non-
// existing custom messages and errors in the underlying storage.
func TestHandleReadCustomMessage(t *testing.T) {
	// Setup backend for storage and a sample custom message.
	storage := &logical.InmemStorage{}
	backend := &SystemBackend{
		Core: &Core{
			customMessageManager: uicustommessages.NewManager(storage),
		},
	}

	nsCtx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	now := time.Now()
	later := now.Add(time.Hour)
	earlier := now.Add(-1 * time.Hour)

	message := &uicustommessages.Message{
		Title:         "title",
		Message:       "message",
		Authenticated: false,
		Type:          uicustommessages.ModalMessageType,
		StartTime:     earlier,
		EndTime:       &later,
		Options:       make(map[string]any),
		Link: &uicustommessages.MessageLink{
			Title: "Click Here",
			Href:  "www.example.com",
		},
	}

	message, err := backend.Core.customMessageManager.AddMessage(nsCtx, *message)
	require.NoError(t, err)
	require.NotNil(t, message)

	fieldData := &framework.FieldData{
		Schema: map[string]*framework.FieldSchema{
			"id": {
				Type: framework.TypeString,
			},
		},
		Raw: map[string]any{
			"id": message.ID,
		},
	}

	// Check that reading the sample custom message succeeds.
	resp, err := backend.handleReadCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "id")
	assert.Equal(t, resp.Data["id"], message.ID)
	assert.Contains(t, resp.Data, "active")
	assert.Equal(t, resp.Data["active"], true)
	assert.Contains(t, resp.Data, "end_time")
	assert.NotNil(t, resp.Data["end_time"])
	assert.Contains(t, resp.Data, "link")
	assert.Equal(t, 1, len(resp.Data["link"].(map[string]string)))

	// Change the message so that it doesn't have an end time.
	message.EndTime = nil
	message.Link = nil
	message, err = backend.Core.customMessageManager.UpdateMessage(nsCtx, *message)
	require.NoError(t, err)
	require.NotNil(t, message)

	resp, err = backend.handleReadCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "id")
	assert.Equal(t, resp.Data["id"], message.ID)
	assert.Contains(t, resp.Data, "active")
	assert.Equal(t, resp.Data["active"], true)
	assert.Contains(t, resp.Data, "end_time")
	assert.Nil(t, resp.Data["end_time"])
	assert.Contains(t, resp.Data, "link")
	assert.Nil(t, resp.Data["link"])

	// Check that there's an error when trying to read a non-existant custom
	// message.
	fieldData.Raw["id"] = "def"

	resp, err = backend.handleReadCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.Error(t, err)
	assert.ErrorIs(t, err, logical.ErrNotFound)
	assert.Nil(t, resp)

	// Check that there's an error when the id parameter is invalid.
	fieldData.Raw["id"] = []bool{}

	resp, err = backend.handleReadCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "error")
	assert.Contains(t, resp.Data["error"], "invalid")

	// Check that there's an error response when the id parameter is missing.
	delete(fieldData.Raw, "id")

	resp, err = backend.handleReadCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "error")
	assert.Contains(t, resp.Data["error"], "missing")

	// Check that there's an error response when there's an error occurring in
	// the underlying storage.
	backend.Core.customMessageManager = uicustommessages.NewManager(&testingStorage{
		getFails: true,
	})

	fieldData.Raw["id"] = message.ID

	resp, err = backend.handleReadCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.NotContains(t, resp.Data, "id")
	assert.Contains(t, resp.Data, "error")
	assert.Contains(t, resp.Data["error"], "failed")
}

// TestHandleUpdateCustomMessage verifies the proper functioning of the
// (*SystemBackend).handleUpdateCustomMessage method. The tests focus on
// missing or invalid request parameters.
func TestHandleUpdateCustomMessage(t *testing.T) {
	storage := &logical.InmemStorage{}

	backend := &SystemBackend{
		Core: &Core{
			customMessageManager: uicustommessages.NewManager(storage),
		},
	}

	now := time.Now()
	startTime := now.Add(-1 * time.Hour).Format(time.RFC3339Nano)
	endTime := now.Add(time.Hour).Format(time.RFC3339Nano)
	startTime2 := now.UTC().Add(-2 * time.Hour).Format(time.RFC3339Nano)

	storageEntryValue := fmt.Sprintf(`{"messages":{"xyz":{"id":"xyz","title":"title","message":"message","authenticated":true,"type":"%s","start_time":"%s","end_time":"%s","link":null,"options":null}}}`, uicustommessages.ModalMessageType, startTime, endTime)

	storageEntry := &logical.StorageEntry{
		Key:   "sys/config/ui/custom-messages",
		Value: []byte(storageEntryValue),
	}

	nsCtx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	err := storage.Put(nsCtx, storageEntry)
	require.NoError(t, err)

	fieldData := &framework.FieldData{
		Schema: map[string]*framework.FieldSchema{
			"id": {
				Type:     framework.TypeString,
				Required: true,
			},
			"title": {
				Type:     framework.TypeString,
				Required: true,
			},
			"message": {
				Type:     framework.TypeString,
				Required: true,
			},
			"authenticated": {
				Type:    framework.TypeBool,
				Default: true,
			},
			"type": {
				Type:    framework.TypeString,
				Default: uicustommessages.BannerMessageType,
			},
			"start_time": {
				Type:     framework.TypeTime,
				Required: true,
			},
			"end_time": {
				Type: framework.TypeTime,
			},
			"link": {
				Type: framework.TypeMap,
			},
			"options": {
				Type: framework.TypeMap,
			},
		},
		Raw: map[string]any{
			"id":            "abc",
			"title":         "title",
			"message":       base64.StdEncoding.EncodeToString([]byte("message")),
			"authenticated": "true",
			"type":          uicustommessages.ModalMessageType,
			"start_time":    startTime,
			"end_time":      endTime,
			"link": map[string]any{
				"link-title": "http://link.url.com",
			},
			"options": map[string]any{},
		},
	}

	// Try to update non-existant custom message
	resp, err := backend.handleUpdateCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.Error(t, err)
	assert.ErrorIs(t, err, logical.ErrNotFound)
	assert.Nil(t, resp)

	// Try to update an existing custom message
	fieldData.Raw["id"] = "xyz"
	fieldData.Raw["start_time"] = startTime2

	resp, err = backend.handleUpdateCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, startTime2, resp.Data["start_time"].(string))

	// Update existing custom message with request containing only required
	// parameters
	delete(fieldData.Raw, "authenticated")
	delete(fieldData.Raw, "type")
	delete(fieldData.Raw, "end_time")

	resp, err = backend.handleUpdateCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.True(t, resp.Data["authenticated"].(bool))
	assert.Equal(t, uicustommessages.BannerMessageType, resp.Data["type"])
	assert.Nil(t, resp.Data["end_time"])

	testcases := []struct {
		name string
		// The testing logic below only works reliably when a single element
		// appears in the fieldRawDelete or fieldRawUpdate for each test case
		// No elements will also work fine for the happy path test.
		fieldRawDelete []string
		fieldRawUpdate map[string]any
	}{
		{
			name:           "id-parameter-missing",
			fieldRawDelete: []string{"id"},
		},
		{
			name: "id-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"id": []int{},
			},
		},
		{
			name:           "title-parameter-missing",
			fieldRawDelete: []string{"title"},
		},
		{
			name: "title-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"title": struct{}{},
			},
		},
		{
			name:           "message-parameter-missing",
			fieldRawDelete: []string{"message"},
		},
		{
			name: "message-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"message": []bool{},
			},
		},
		{
			name: "message-parameter-invalid-not-base64",
			fieldRawUpdate: map[string]any{
				"message": "The message not base64 encoded.",
			},
		},
		{
			name: "authenticated-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"authenticated": "fred",
			},
		},
		{
			name: "type-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"type": []int{1},
			},
		},
		{
			name:           "start_time-parameter-missing",
			fieldRawDelete: []string{"start_time"},
		},
		{
			name: "start_time-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"start_time": "tomorrow",
			},
		},
		{
			name: "end_time-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"end_time": "yesterday",
			},
		},
		{
			name: "link-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"link": "link",
			},
		},
		{
			name: "link-parameter-url-invalid",
			fieldRawUpdate: map[string]any{
				"link": map[string]any{
					"my-link": []int{},
				},
			},
		},
		{
			name: "link-parameter-multiple-links",
			fieldRawUpdate: map[string]any{
				"link": map[string]any{
					"click here":   "http://example.org",
					"click here 2": "http://ping.net",
				},
			},
		},
		{
			name: "options-parameter-invalid",
			fieldRawUpdate: map[string]any{
				"options": "options",
			},
		},
	}

	for _, testcase := range testcases {
		raw := map[string]any{}
		for k, v := range fieldData.Raw {
			raw[k] = v
		}

		for _, d := range testcase.fieldRawDelete {
			delete(raw, d)
		}

		for k, v := range testcase.fieldRawUpdate {
			raw[k] = v
		}

		resp, err := backend.handleUpdateCustomMessage(nsCtx, &logical.Request{}, &framework.FieldData{
			Schema: fieldData.Schema,
			Raw:    raw,
		})
		assert.NoError(t, err, testcase.name)
		assert.NotNil(t, resp, testcase.name)
		assert.NotNil(t, resp.Data, testcase.name)

		if len(testcase.fieldRawDelete) > 0 {
			assert.Contains(t, resp.Data, "error", testcase.name)
			assert.Contains(t, resp.Data["error"], "missing", testcase.name)
			assert.Contains(t, resp.Data["error"], testcase.fieldRawDelete[0], testcase.name)
		}

		if len(testcase.fieldRawUpdate) > 0 {
			var keyName string
			for k := range testcase.fieldRawUpdate {
				keyName = k
				break
			}
			assert.Contains(t, resp.Data, "error", testcase.name)
			assert.Contains(t, resp.Data["error"], "invalid", testcase.name)
			assert.Contains(t, resp.Data["error"], keyName, testcase.name)
		}

		if len(testcase.fieldRawDelete)+len(testcase.fieldRawUpdate) == 0 {
			assert.Contains(t, resp.Data, "active", testcase.name)
			assert.Contains(t, resp.Data, "id", testcase.name)
		}
	}

	// Check that there's an error response if an error occurred in the
	// underlying storage.
	backend.Core.customMessageManager = uicustommessages.NewManager(&testingStorage{
		putFails:         true,
		getResponseValue: storageEntryValue,
	})

	resp, err = backend.handleUpdateCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "error")
	assert.Contains(t, resp.Data["error"], "failed")
}

// TestHandleDeleteCustomMessage verifies the proper functioning of the
// (*SystemBackend).handleDeleteCustomMessage method. The test focuses on
// missing and invalid request parameters as well as errors occurring in the
// underlying storage.
func TestHandleDeleteCustomMessage(t *testing.T) {
	// Setup backend for testing
	storage := &logical.InmemStorage{}

	backend := &SystemBackend{
		Core: &Core{
			customMessageManager: uicustommessages.NewManager(storage),
		},
	}

	startTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339Nano)
	endTime := time.Now().Add(time.Hour).Format(time.RFC3339Nano)

	storageEntryValue := fmt.Sprintf(`{"messages":{"xyz":{"id":"xyz","title":"title","message":"message","authenticated":true,"type":"%s","start_time":"%s","end_time":"%s","link":{},"options":{}}}}`, uicustommessages.ModalMessageType, startTime, endTime)
	storageEntry := &logical.StorageEntry{
		Key:   "root",
		Value: []byte(storageEntryValue),
	}

	nsCtx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	err := storage.Put(nsCtx, storageEntry)
	require.NoError(t, err)

	fieldData := &framework.FieldData{
		Schema: map[string]*framework.FieldSchema{
			"id": {
				Type: framework.TypeString,
			},
		},
		Raw: map[string]any{
			"id": "abc",
		},
	}

	// Check if deleting a non-existant custom message is ok.
	resp, err := backend.handleDeleteCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.Nil(t, resp)

	// Check if the id parameter is invalid.
	fieldData.Raw["id"] = []int{}

	resp, err = backend.handleDeleteCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "error")
	assert.Contains(t, resp.Data["error"], "invalid")
	assert.Contains(t, resp.Data["error"], "id")

	// Check if the id parameter is missing.
	delete(fieldData.Raw, "id")

	resp, err = backend.handleDeleteCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "error")
	assert.Contains(t, resp.Data["error"], "missing")
	assert.Contains(t, resp.Data["error"], "id")

	// Check that deleting an existing message works.
	fieldData.Raw["id"] = "xyz"

	resp, err = backend.handleDeleteCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.Nil(t, resp)

	// Check that errors in the underlying storage result in an error response.
	backend.Core.customMessageManager = uicustommessages.NewManager(&testingStorage{
		getFails: true,
	})

	resp, err = backend.handleDeleteCustomMessage(nsCtx, &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Data)
	assert.Contains(t, resp.Data, "error")
	assert.Contains(t, resp.Data["error"], "failed")
}

// TestHandleCustomMessageExistenceCheck verifies the proper functioning of the
// (*SystemBackend).handleCustomMessageExistenceCheck method.
func TestHandleCustomMessageExistenceCheck(t *testing.T) {
	fieldData := &framework.FieldData{
		Schema: map[string]*framework.FieldSchema{
			"id": {
				Type: framework.TypeString,
			},
		},
	}

	backend := &SystemBackend{}

	// Check that when id is provided and valid, true and no error are returned.
	found, err := backend.handleCustomMessageExistenceCheck(context.Background(), &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.True(t, found)

	// Check that when id is provided but invalid that there's an error
	// returned.
	delete(fieldData.Schema, "id")

	found, err = backend.handleCustomMessageExistenceCheck(context.Background(), &logical.Request{}, fieldData)
	assert.NoError(t, err)
	assert.False(t, found)
}

// testingStorage is a struct with methods that satisfy the logical.Storage
// interface. It can be programmed to unconditionally return errors for any
// of its methods.
type testingStorage struct {
	listFails   bool
	getFails    bool
	deleteFails bool
	putFails    bool

	listResponseValue []string
	getResponseValue  string
}

// List returns a single key, dummy, unless the receiver has been configured to
// return an error, in which case it returns nil and a made up error.
func (s *testingStorage) List(_ context.Context, _ string) ([]string, error) {
	if s.listFails {
		return nil, errors.New("failure")
	}

	return append(make([]string, 0), s.listResponseValue...), nil
}

// Get returns a dummy logical.StorageEntry unless the receiver has been
// configured to return an error, in which case it returns nil and a made up
// error.
func (s *testingStorage) Get(_ context.Context, key string) (*logical.StorageEntry, error) {
	if s.getFails {
		return nil, errors.New("failure")
	}

	return &logical.StorageEntry{
		Key:   key,
		Value: []byte(s.getResponseValue),
	}, nil
}

// Delete returns nothing, unless the receiver has been configured to return an
// error, in which case it returns a made up error.
func (s *testingStorage) Delete(_ context.Context, _ string) error {
	if s.deleteFails {
		return errors.New("failure")
	}

	return nil
}

// Put returns nothing, unless the receiver has been configured to return an
// error, in which case it returns a made up error.
func (s *testingStorage) Put(_ context.Context, _ *logical.StorageEntry) error {
	if s.putFails {
		return errors.New("failure")
	}

	return nil
}
