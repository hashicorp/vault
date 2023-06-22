// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const respErrEmptyName = "empty name attribute given"

// SubscriptionConfig is used by the Factory function to configure a Subscription
// object.
type SubscriptionConfig struct {
	PluginName    string `json:"plugin_name" structs:"plugin_name" mapstructure:"plugin_name"`
	PluginVersion string `json:"plugin_version" structs:"plugin_version" mapstructure:"plugin_version"`
	// Settings stores the subscription-specific settings needed by each event plugin type.
	Settings map[string]interface{} `json:"settings" structs:"settings" mapstructure:"settings"`
}

// pathResetSubscription configures a path to reset a plugin.
func pathResetSubscription(b *eventBackend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("reset/%s", framework.GenericNameRegex("name")),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixEvent,
			OperationVerb:   "reset",
			OperationSuffix: operationSuffixSubscription,
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of this subscription",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathSubscriptionReset(),
			},
		},

		HelpSynopsis:    pathResetSubscriptionHelpSyn,
		HelpDescription: pathResetSubscriptionHelpDesc,
	}
}

func (b *eventBackend) pathSubscriptionReset() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		// Close plugin and delete the entry in the subscriptions cache.
		if err := b.ClearSubscription(name); err != nil {
			return nil, err
		}

		// Execute plugin again, we don't need the object so throw away.
		if _, err := b.GetSubscription(ctx, req.Storage, name); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// pathConfigureSubscription returns a configured framework.Path setup to
// operate on plugins.
func pathConfigureSubscription(b *eventBackend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("%s%s", eventSubscriptionPath, framework.GenericNameRegex("name")),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixEvent,
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of this subscription",
			},

			"plugin_name": {
				Type: framework.TypeString,
				Description: `The name of a builtin or previously registered
				plugin known to vault. This endpoint will create an instance of
				that plugin type.`,
			},

			"plugin_version": {
				Type:        framework.TypeString,
				Description: `The version of the plugin to use.`,
			},

			"verify": {
				Type:    framework.TypeBool,
				Default: true,
				Description: `If true, the subscription settings are verified by
				actually connecting to the subscription. Defaults to true.`,
			},
		},

		ExistenceCheck: b.subscriptionExistenceCheck(),

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.subscriptionWriteHandler(),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: operationSuffixSubscription,
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.subscriptionWriteHandler(),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: operationSuffixSubscription,
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.subscriptionReadHandler(),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: operationSuffixSubscription,
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.subscriptionDeleteHandler(),
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "delete",
					OperationSuffix: operationSuffixSubscription,
				},
			},
		},

		HelpSynopsis:    pathSubscriptionHelpSyn,
		HelpDescription: pathSubscriptionHelpDesc,
	}
}

func (b *eventBackend) subscriptionExistenceCheck() framework.ExistenceFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
		name := data.Get("name").(string)
		if name == "" {
			return false, errors.New(`missing "name" parameter`)
		}

		entry, err := req.Storage.Get(ctx, fmt.Sprintf("%s%s", eventSubscriptionPath, name))
		if err != nil {
			return false, fmt.Errorf("failed to read subscription configuration: %w", err)
		}

		return entry != nil, nil
	}
}

func pathListSubscription(b *eventBackend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf(eventSubscriptionPath + "?$"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixEvent,
			OperationSuffix: operationSuffixSubscription,
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.subscriptionListHandler(),
			},
		},

		HelpSynopsis:    pathSubscriptionHelpSyn,
		HelpDescription: pathSubscriptionHelpDesc,
	}
}

func (b *eventBackend) subscriptionListHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		entries, err := req.Storage.List(ctx, eventSubscriptionPath)
		if err != nil {
			return nil, err
		}

		return logical.ListResponse(entries), nil
	}
}

// subscriptionReadHandler reads out the subscription configuration
func (b *eventBackend) subscriptionReadHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		entry, err := req.Storage.Get(ctx, fmt.Sprintf("%s%s", eventSubscriptionPath, name))
		if err != nil {
			return nil, fmt.Errorf("failed to read subscription configuration: %w", err)
		}
		if entry == nil {
			return nil, nil
		}

		var config SubscriptionConfig
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, err
		}

		if versions.IsBuiltinVersion(config.PluginVersion) {
			// This gets treated as though it's empty when mounting, and will get
			// overwritten to be empty when the config is next written. See #18051.
			config.PluginVersion = ""
		}

		return &logical.Response{
			Data: structs.New(config).Map(),
		}, nil
	}
}

// subscriptionDeleteHandler deletes the subscription configuration
func (b *eventBackend) subscriptionDeleteHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		err := req.Storage.Delete(ctx, fmt.Sprintf("%s%s", eventSubscriptionPath, name))
		if err != nil {
			return nil, fmt.Errorf("failed to delete subscription configuration: %w", err)
		}

		if err := b.ClearSubscription(name); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// subscriptionWriteHandler returns a handler function for creating and updating
// both builtin and event plugin types.
func (b *eventBackend) subscriptionWriteHandler() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		return nil, fmt.Errorf("not implemented yet")
	}
}

const pathSubscriptionHelpSyn = `
Configure subscription details for an event plugin.
`

const pathSubscriptionHelpDesc = ``

const pathResetSubscriptionHelpSyn = `
Resets the subscription.
`

const pathResetSubscriptionHelpDesc = `
This path resets the subscription by closing the existing event plugin
instance and running a new one.
`
