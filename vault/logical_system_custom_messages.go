// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	uicustommessages "github.com/hashicorp/vault/vault/ui_custom_messages"
)

// uiCustomMessagePaths returns a slice of *framework.Path elements that are
// added to the receiver SystemBackend.
func (b *SystemBackend) uiCustomMessagePaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "config/ui/custom-messages/$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "ui-config",
				OperationVerb:   "list",
				OperationSuffix: "custom-messages",
			},
			Fields: map[string]*framework.FieldSchema{
				"authenticated": {
					Type:     framework.TypeBool,
					Required: false,
				},
				"active": {
					Type:     framework.TypeBool,
					Required: false,
				},
				"type": {
					Type:     framework.TypeString,
					Required: false,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleListCustomMessages,
					Summary:  "Lists custom messages",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
						}},
					},
				},
			},
		},
		{
			Pattern: "config/ui/custom-messages$",

			ExistenceCheck: b.handleCustomMessageExistenceCheck,

			Fields: map[string]*framework.FieldSchema{
				"title": {
					Type:     framework.TypeString,
					Required: true,
				},
				"type": {
					Type:     framework.TypeString,
					Required: false,
					Default:  uicustommessages.BannerMessageType,
				},
				"authenticated": {
					Type:     framework.TypeBool,
					Required: false,
					Default:  true,
				},
				"message": {
					Type:     framework.TypeString,
					Required: true,
				},
				"start_time": {
					Type:     framework.TypeTime,
					Required: true,
				},
				"end_time": {
					Type:     framework.TypeTime,
					Required: false,
				},
				"link": {
					Type:     framework.TypeMap,
					Required: false,
				},
				"options": {
					Type:     framework.TypeMap,
					Required: false,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleCreateCustomMessages,
					Summary:  "Create custom message",
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
						}},
					},
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "create",
						OperationSuffix: "custom-message",
					},
				},
			},
		},
		{
			Pattern: "config/ui/custom-messages/" + framework.MatchAllRegex("id"),
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "The unique identifier for the custom message",
				},
				"title": {
					Type:     framework.TypeString,
					Required: true,
				},
				"type": {
					Type:     framework.TypeString,
					Required: false,
					Default:  uicustommessages.BannerMessageType,
				},
				"authenticated": {
					Type:     framework.TypeBool,
					Required: false,
					Default:  true,
				},
				"message": {
					Type:     framework.TypeString,
					Required: true,
				},
				"start_time": {
					Type:     framework.TypeTime,
					Required: true,
				},
				"end_time": {
					Type:     framework.TypeTime,
					Required: false,
				},
				"link": {
					Type:     framework.TypeMap,
					Required: false,
				},
				"options": {
					Type:     framework.TypeMap,
					Required: false,
				},
			},
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "ui-config",
				OperationSuffix: "custom-message",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleReadCustomMessage,
					Summary:  "Read custom message",

					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"id": {
									Type:     framework.TypeString,
									Required: true,
								},
							},
						}},
					},

					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "read",
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDeleteCustomMessage,
					Summary:  "Delete custom message",

					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"id": {
									Type:     framework.TypeString,
									Required: true,
								},
							},
						}},
					},

					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "delete",
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleUpdateCustomMessage,
					Summary:  "Update custom message",

					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"id": {
									Type:     framework.TypeString,
									Required: true,
								},
								"title": {
									Type:     framework.TypeString,
									Required: true,
								},
								"authenticated": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"type": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"message": {
									Type:     framework.TypeString,
									Required: true,
								},
								"link": {
									Type:     framework.TypeMap,
									Required: false,
								},
								"start_time": {
									Type:     framework.TypeTime,
									Required: true,
								},
								"end_time": {
									Type:     framework.TypeTime,
									Required: true,
								},
								"options": {
									Type:     framework.TypeMap,
									Required: false,
								},
							},
						}},
					},

					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "update",
					},
				},
			},
		},
	}
}

// handleListCustomMessages is the operation callback for the LIST operation of
// the custom messages endpoint.
func (b *SystemBackend) handleListCustomMessages(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var filters uicustommessages.FindFilter

	err := parameterValidateAndUse[bool]("authenticated", func(v bool) error {
		filters.Authenticated(v)
		return nil
	}, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	err = parameterValidateAndUse[string]("type", filters.Type, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	err = parameterValidateAndUse[bool]("active", func(v bool) error {
		filters.Active(v)
		return nil
	}, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	messages, err := b.Core.customMessageManager.FindMessages(ctx, filters)
	if err != nil {
		return logical.ErrorResponse("failed to retrieve list of custom messages: %s", err), nil
	}

	keyInfos := make(map[string]any)
	keys := make([]string, len(messages))

	for i, message := range messages {
		m := make(map[string]any)
		m["title"] = message.Title
		m["type"] = message.Type
		m["authenticated"] = message.Authenticated
		m["start_time"] = message.StartTime
		m["end_time"] = message.EndTime
		m["active"] = message.Active()

		keyInfos[message.ID] = m
		keys[i] = message.ID
	}

	return logical.ListResponseWithInfo(keys, keyInfos), nil
}

// parameterValidateAndUse is a helper that retrieves a parameter from the
// provided framework.FieldData if it exists and is valid then calls the
// provided setter method (filterSetter) using that parameter value as the
// argument. If the parameter contains an invalid value, an error is returned.
// If the parameter does not have a value, nothing happens.
func parameterValidateAndUse[T bool | string](parameterName string, filterSetter func(T) error, d *framework.FieldData) error {
	value, ok, err := d.GetOkErr(parameterName)
	if err != nil {
		return fmt.Errorf("invalid %s parameter value: %s", parameterName, err)
	}

	if ok {
		filterSetter(value.(T))
	}

	return nil
}

// parameterValidateOrReportMissing is a helper that retrieves a parameter from
// the provided framework.FieldData if it exists and is valid. If the parameter
// contains an invalid value, an error is returned. If the parameter does not
// have a value, an error is returned.
func parameterValidateOrReportMissing[T string | bool | time.Time](parameterName string, d *framework.FieldData) (T, error) {
	var empty T

	value, ok, err := d.GetOkErr(parameterName)
	if err != nil {
		return empty, fmt.Errorf("invalid %s parameter value: %s", parameterName, err)
	}

	if !ok {
		return empty, fmt.Errorf("missing %s parameter value", parameterName)
	}

	return value.(T), nil
}

// parameterValidateOrUseDefault is a helper that retrieves a parameter from
// the provided framework.FieldData if it exists and is valid. If the parameter
// contains an invalid value, an error is returned. If the parameter does not
// have a value, its default value is returned.
func parameterValidateOrUseDefault[T string | bool](parameterName string, d *framework.FieldData) (T, error) {
	var empty T

	value, ok, err := d.GetOkErr(parameterName)
	if err != nil {
		return empty, fmt.Errorf("invalid %s parameter value: %s", parameterName, err)
	}

	if !ok {
		value = d.GetDefaultOrZero(parameterName)
	}

	return value.(T), nil
}

// parameterValidateMap is a helper that retrieves a parameter from the provided
// framework.FieldData if it exists and is valid. If the parameter contains an
// invalid value, an error is returned. If the parameter does not have a value,
// nothing happens.
func parameterValidateMap(parameterName string, d *framework.FieldData) (map[string]any, error) {
	value, ok, err := d.GetOkErr(parameterName)
	if err != nil {
		return nil, fmt.Errorf("invalid %s parameter value: %s", parameterName, err)
	}

	if ok {
		return value.(map[string]any), nil
	}

	return nil, nil
}

// handleCreateCustomMessages is the operation callback for the CREATE operation
// of the custom messages endpoint.
func (b *SystemBackend) handleCreateCustomMessages(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	title, err := parameterValidateOrReportMissing[string]("title", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	authenticated, err := parameterValidateOrUseDefault[bool]("authenticated", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	messageType, err := parameterValidateOrUseDefault[string]("type", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	messageValue, err := parameterValidateOrReportMissing[string]("message", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	_, err = base64.StdEncoding.DecodeString(messageValue)
	if err != nil {
		return logical.ErrorResponse("invalid message parameter value, must be base64 encoded"), nil
	}

	startTime, err := parameterValidateOrReportMissing[time.Time]("start_time", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	var endTime *time.Time
	endTimeValue, ok, err := d.GetOkErr("end_time")
	if err != nil {
		return logical.ErrorResponse("invalid end_time parameter value: %s", err), nil
	}
	if ok {
		value := endTimeValue.(time.Time)
		endTime = &value
	}

	linkMap, err := parameterValidateMap("link", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	link, resp := validateLinkMap(linkMap)
	if resp != nil {
		return resp, nil
	}

	options, err := parameterValidateMap("options", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	message := &uicustommessages.Message{
		Title:         title,
		Authenticated: authenticated,
		Type:          messageType,
		Message:       messageValue,
		StartTime:     startTime,
		EndTime:       endTime,
		Link:          link,
		Options:       options,
	}

	message, err = b.Core.customMessageManager.AddMessage(ctx, *message)
	if err != nil {
		return logical.ErrorResponse("failed to create custom message: %s", err), nil
	}

	var endTimeResponse any
	if message.EndTime != nil {
		endTimeResponse = message.EndTime.Format(time.RFC3339Nano)
	}

	return &logical.Response{
		Data: map[string]any{
			"id":            message.ID,
			"authenticated": message.Authenticated,
			"type":          message.Type,
			"message":       message.Message,
			"start_time":    message.StartTime.Format(time.RFC3339Nano),
			"end_time":      endTimeResponse,
			"link":          message.Link,
			"options":       message.Options,
			"active":        message.Active(),
		},
	}, nil
}

// validateLinkMap takes care of detecting either incomplete or invalid link
// parameter value.
// An invalid link parameter value is one where there's either
// an empty string for the title key or the href value or the href value not
// being a string. A linkMap that is invalid, results in only a logical.Response
// containing an error response being returned.
// An incomplete link parameter value is one where the linkMap is either nil or
// empty, where the title key is an empty string and the href value is an empty
// string. A linkMap that is incomplete, results in neither a MessageLink nor a
// Response being returned.
// If the linkMap is neither invalid nor incomplete, a MessageLink is returned.
func validateLinkMap(linkMap map[string]any) (*uicustommessages.MessageLink, *logical.Response) {
	if len(linkMap) > 1 {
		return nil, logical.ErrorResponse("invalid number of elements in link parameter value; only a single element can be provided")
	}

	for k, v := range linkMap {
		href, ok := v.(string)

		switch {
		case !ok:
			// href value is not a string, so return an error
			return nil, logical.ErrorResponse(fmt.Sprintf("invalid url for %q key in link parameter value", k))
		case len(k) == 0 && len(href) > 0:
			// no title key, but there's an href value, so return an error
			return nil, logical.ErrorResponse("missing title key in link parameter value")
		case len(k) > 0 && len(href) == 0:
			// no href value, but there's a title key, so return an error
			return nil, logical.ErrorResponse(fmt.Sprintf("missing url for %q key in link parameter value", k))
		case len(k) != 0 && len(href) != 0:
			// when title key and href value are not empty, return a MessageLink
			// pointer
			return &uicustommessages.MessageLink{
				Title: k,
				Href:  href,
			}, nil
		}

	}

	// no title key and no href value, treat it as if no link was specified
	return nil, nil
}

// handleReadCustomMessage is the operation callback for the READ operation of
// the custom messages endpoint.
func (b *SystemBackend) handleReadCustomMessage(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	id, err := parameterValidateOrReportMissing[string]("id", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	message, err := b.Core.customMessageManager.ReadMessage(ctx, id)
	switch {
	case errors.Is(err, logical.ErrNotFound):
		return nil, err
	case err != nil:
		return logical.ErrorResponse("failed to retrieve custom message: %s", err), nil
	}

	var endTimeResponse any
	if message.EndTime != nil {
		endTimeResponse = message.EndTime.Format(time.RFC3339Nano)
	}

	var linkResponse map[string]string = nil
	if message.Link != nil {
		linkResponse = make(map[string]string)

		linkResponse[message.Link.Title] = message.Link.Href
	}

	return &logical.Response{
		Data: map[string]any{
			"id":            id,
			"authenticated": message.Authenticated,
			"type":          message.Type,
			"message":       message.Message,
			"start_time":    message.StartTime.Format(time.RFC3339Nano),
			"end_time":      endTimeResponse,
			"link":          linkResponse,
			"options":       message.Options,
			"active":        message.Active(),
			"title":         message.Title,
		},
	}, nil
}

// handleUpdateCustomMessage is the operation callback for the UPDATE operation
// of the custom messages endpoint.
func (b *SystemBackend) handleUpdateCustomMessage(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	messageID, err := parameterValidateOrReportMissing[string]("id", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	title, err := parameterValidateOrReportMissing[string]("title", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	authenticated, err := parameterValidateOrUseDefault[bool]("authenticated", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	messageType, err := parameterValidateOrUseDefault[string]("type", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	messageValue, err := parameterValidateOrReportMissing[string]("message", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	_, err = base64.StdEncoding.DecodeString(messageValue)
	if err != nil {
		return logical.ErrorResponse("invalid message parameter value, must be base64 encoded"), nil
	}

	linkMap, err := parameterValidateMap("link", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	link, resp := validateLinkMap(linkMap)
	if resp != nil {
		return resp, nil
	}

	options, err := parameterValidateMap("options", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	startTime, err := parameterValidateOrReportMissing[time.Time]("start_time", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	var endTime *time.Time
	endTimeValue, ok, err := d.GetOkErr("end_time")
	if err != nil {
		return logical.ErrorResponse("invalid end_time parameter value: %s", err), nil
	}
	if ok {
		value := endTimeValue.(time.Time)
		endTime = &value
	}

	message := &uicustommessages.Message{
		ID:            messageID,
		Title:         title,
		Authenticated: authenticated,
		Type:          messageType,
		Message:       messageValue,
		Link:          link,
		Options:       options,
		StartTime:     startTime,
		EndTime:       endTime,
	}

	message, err = b.Core.customMessageManager.UpdateMessage(ctx, *message)
	switch {
	case errors.Is(err, logical.ErrNotFound):
		return nil, err
	case err != nil:
		return logical.ErrorResponse("failed to update custom message: %s", err), nil
	}

	var endTimeResponse any
	if message.EndTime != nil {
		endTimeResponse = message.EndTime.Format(time.RFC3339Nano)
	}

	return &logical.Response{
		Data: map[string]any{
			"id":            message.ID,
			"active":        message.Active(),
			"start_time":    message.StartTime.Format(time.RFC3339Nano),
			"end_time":      endTimeResponse,
			"type":          message.Type,
			"authenticated": message.Authenticated,
		},
	}, nil
}

// handleDeleteCustomMessage is the operation callback for the DELETE operation
// of the custom messages endpoint.
func (b *SystemBackend) handleDeleteCustomMessage(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	id, err := parameterValidateOrReportMissing[string]("id", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if err = b.Core.customMessageManager.DeleteMessage(ctx, id); err != nil {
		return logical.ErrorResponse("failed to delete custom message: %s", err), nil
	}

	return nil, nil
}

// handleCustomMessageExistenceCheck is the function that fills the
// framework.Path ExistenceCheck role for custom messages.
func (b *SystemBackend) handleCustomMessageExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	_, ok := d.Schema["id"]
	return ok, nil
}
