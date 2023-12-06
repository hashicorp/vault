package vault

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	uicustommessages "github.com/hashicorp/vault/vault/ui_custom_messages"
)

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
					Required: true,
				},
				"authenticated": {
					Type:     framework.TypeBool,
					Required: true,
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
					Required: true,
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
					Required: true,
				},
				"authenticated": {
					Type:     framework.TypeBool,
					Required: true,
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
					Required: true,
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

	err := parameterValidateAndUse[bool]("authenticated", filters.Authenticated, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	err = parameterValidateAndUse[string]("type", filters.Type, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	err = parameterValidateAndUse[bool]("active", filters.Active, d)
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
// argument.
func parameterValidateAndUse[T bool | string](parameterName string, filterSetter func(T), d *framework.FieldData) error {
	value, ok, err := d.GetOkErr(parameterName)
	if err != nil {
		return fmt.Errorf("invalid %s parameter value: %s", parameterName, err)
	}

	if ok {
		filterSetter(value.(T))
	}

	return nil
}

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

func parameterValidate[T map[string]any](parameterName string, variable *T, d *framework.FieldData) error {
	value, ok, err := d.GetOkErr(parameterName)
	if err != nil {
		return fmt.Errorf("invalid %s parameter value: %s", parameterName, err)
	}

	if ok {
		*variable = value.(T)
	}

	return nil
}

// handleCreateCustomMessages is the operation callback for the CREATE operation
// of the custom messages endpoint.
func (b *SystemBackend) handleCreateCustomMessages(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	title, err := parameterValidateOrReportMissing[string]("title", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	authenticated, err := parameterValidateOrReportMissing[bool]("authenticated", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	messageType, err := parameterValidateOrReportMissing[string]("type", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	messageValue, err := parameterValidateOrReportMissing[string]("message", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	startTime, err := parameterValidateOrReportMissing[time.Time]("start_time", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	endTime, err := parameterValidateOrReportMissing[time.Time]("end_time", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	link := make(map[string]any)
	err = parameterValidate[map[string]any]("link", &link, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	options := make(map[string]any)
	err = parameterValidate[map[string]any]("options", &options, d)
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

	message, err = b.Core.customMessageManager.CreateMessage(ctx, *message)
	if err != nil {
		return logical.ErrorResponse("failed to create custom message: %s", err), nil
	}

	return &logical.Response{
		Data: map[string]any{
			"id": message.ID,
			"data": map[string]any{
				"authenticated": message.Authenticated,
				"type":          message.Type,
				"message":       message.Message,
				"start_time":    message.StartTime.Format(time.RFC3339Nano),
				"end_time":      message.EndTime.Format(time.RFC3339Nano),
				"link":          message.Link,
				"options":       message.Options,
				"active":        message.Active(),
			},
		},
	}, nil
}

func (b *SystemBackend) handleReadCustomMessage(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	id, err := parameterValidateOrReportMissing[string]("id", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	message, err := b.Core.customMessageManager.ReadMessage(ctx, id)
	if err != nil {
		return logical.ErrorResponse("failed to retrieve custom message: %s", err), nil
	}

	if message == nil {
		return nil, logical.ErrCustomMessageNotFound
	}

	return &logical.Response{
		Data: map[string]any{
			"id": id,
			"data": map[string]any{
				"authenticated": message.Authenticated,
				"type":          message.Type,
				"message":       message.Message,
				"start_time":    message.StartTime.Format(time.RFC3339Nano),
				"end_time":      message.EndTime.Format(time.RFC3339Nano),
				"link":          message.Link,
				"options":       message.Options,
				"active":        message.Active(),
				"title":         message.Title,
			},
		},
	}, nil
}

func (b *SystemBackend) handleUpdateCustomMessage(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	messageID, err := parameterValidateOrReportMissing[string]("id", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	title, err := parameterValidateOrReportMissing[string]("title", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	authenticated, err := parameterValidateOrReportMissing[bool]("authenticated", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	messageType, err := parameterValidateOrReportMissing[string]("type", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	messageValue, err := parameterValidateOrReportMissing[string]("message", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	var link map[string]any
	err = parameterValidate[map[string]any]("link", &link, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	var options map[string]any
	err = parameterValidate[map[string]any]("options", &options, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	startTime, err := parameterValidateOrReportMissing[time.Time]("start_time", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	endTime, err := parameterValidateOrReportMissing[time.Time]("end_time", d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
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
	if err != nil {
		return logical.ErrorResponse("failed to update custom message: %s", err), nil
	}

	if message == nil {
		return nil, logical.ErrCustomMessageNotFound
	}

	return &logical.Response{
		Data: map[string]any{
			"id": message.ID,
			"data": map[string]any{
				"active":        message.Active(),
				"start_time":    message.StartTime.Format(time.RFC3339Nano),
				"end_time":      message.EndTime.Format(time.RFC3339Nano),
				"type":          message.Type,
				"authenticated": message.Authenticated,
			},
		},
	}, nil
}

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

func (b *SystemBackend) handleCustomMessageExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	_, ok := d.Schema["id"]
	return ok, nil
}
