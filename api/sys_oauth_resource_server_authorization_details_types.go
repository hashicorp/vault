// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

const oAuthResourceServerAuthorizationDetailsTypesPath = "/v1/sys/oauth-resource-server/authorization-details-types"

// OAuthResourceServerAuthorizationDetailsTypeSummary contains catalog summary
// metadata for a published authorization_details type.
type OAuthResourceServerAuthorizationDetailsTypeSummary struct {
	Type          string `json:"type" mapstructure:"type"`
	Status        string `json:"status" mapstructure:"status"`
	SchemaVersion string `json:"schema_version" mapstructure:"schema_version"`
	Title         string `json:"title" mapstructure:"title"`
	Description   string `json:"description" mapstructure:"description"`
	DocsURL       string `json:"docs_url" mapstructure:"docs_url"`
	IntroducedIn  string `json:"introduced_in" mapstructure:"introduced_in"`
	DeprecatedIn  string `json:"deprecated_in,omitempty" mapstructure:"deprecated_in,omitempty"`
	SchemaPath    string `json:"schema_path" mapstructure:"schema_path"`
}

// OAuthResourceServerAuthorizationDetailsTypesCatalog contains discovery
// metadata for all supported authorization_details types.
type OAuthResourceServerAuthorizationDetailsTypesCatalog struct {
	APIVersion                 string                                               `json:"api_version" mapstructure:"api_version"`
	DefaultUnknownTypeBehavior string                                               `json:"default_unknown_type_behavior" mapstructure:"default_unknown_type_behavior"`
	Types                      []OAuthResourceServerAuthorizationDetailsTypeSummary `json:"types" mapstructure:"types"`
}

// OAuthResourceServerAuthorizationDetailsTypeDetail contains detailed
// metadata for a single published authorization_details type.
type OAuthResourceServerAuthorizationDetailsTypeDetail struct {
	Type              string         `json:"type" mapstructure:"type"`
	Status            string         `json:"status" mapstructure:"status"`
	SchemaVersion     string         `json:"schema_version" mapstructure:"schema_version"`
	Title             string         `json:"title" mapstructure:"title"`
	Description       string         `json:"description" mapstructure:"description"`
	DocsURL           string         `json:"docs_url" mapstructure:"docs_url"`
	IntroducedIn      string         `json:"introduced_in" mapstructure:"introduced_in"`
	DeprecatedIn      string         `json:"deprecated_in,omitempty" mapstructure:"deprecated_in,omitempty"`
	RequiredFields    []string       `json:"required_fields" mapstructure:"required_fields"`
	OptionalFields    []string       `json:"optional_fields" mapstructure:"optional_fields"`
	MatchingSemantics map[string]any `json:"matching_semantics" mapstructure:"matching_semantics"`
	Examples          []any          `json:"examples" mapstructure:"examples"`
	SchemaPath        string         `json:"schema_path" mapstructure:"schema_path"`
}

// OAuthResourceServerAuthorizationDetailsTypeSchema contains the published
// machine-readable JSON schema for a single authorization_details type.
type OAuthResourceServerAuthorizationDetailsTypeSchema struct {
	Type            string         `json:"type" mapstructure:"type"`
	SchemaVersion   string         `json:"schema_version" mapstructure:"schema_version"`
	JSONSchemaDraft string         `json:"json_schema_draft" mapstructure:"json_schema_draft"`
	Schema          map[string]any `json:"schema" mapstructure:"schema"`
}

// ListOAuthResourceServerAuthorizationDetailsTypes returns all published
// authorization_details type identifiers.
func (c *Sys) ListOAuthResourceServerAuthorizationDetailsTypes() ([]string, error) {
	return c.ListOAuthResourceServerAuthorizationDetailsTypesWithContext(context.Background())
}

// ListOAuthResourceServerAuthorizationDetailsTypesWithContext returns all
// published authorization_details type identifiers.
func (c *Sys) ListOAuthResourceServerAuthorizationDetailsTypesWithContext(ctx context.Context) ([]string, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	req := c.c.NewRequest(http.MethodGet, oAuthResourceServerAuthorizationDetailsTypesPath)
	req.Params.Set("list", "true")

	secret, err := c.readSecretWithContext(ctx, req)
	if err != nil {
		return nil, err
	}

	keysRaw, ok := secret.Data["keys"]
	if !ok {
		return nil, errors.New("data from server response does not contain keys")
	}

	var keys []string
	if err := mapstructure.Decode(keysRaw, &keys); err != nil {
		return nil, err
	}

	return keys, nil
}

// ReadOAuthResourceServerAuthorizationDetailsTypesCatalog returns catalog
// metadata for all published authorization_details types.
func (c *Sys) ReadOAuthResourceServerAuthorizationDetailsTypesCatalog() (*OAuthResourceServerAuthorizationDetailsTypesCatalog, error) {
	return c.ReadOAuthResourceServerAuthorizationDetailsTypesCatalogWithContext(context.Background())
}

// ReadOAuthResourceServerAuthorizationDetailsTypesCatalogWithContext returns
// catalog metadata for all published authorization_details types.
func (c *Sys) ReadOAuthResourceServerAuthorizationDetailsTypesCatalogWithContext(ctx context.Context) (*OAuthResourceServerAuthorizationDetailsTypesCatalog, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	req := c.c.NewRequest(http.MethodGet, oAuthResourceServerAuthorizationDetailsTypesPath)
	secret, err := c.readSecretWithContext(ctx, req)
	if err != nil {
		return nil, err
	}

	var result OAuthResourceServerAuthorizationDetailsTypesCatalog
	if err := mapstructure.Decode(secret.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ReadOAuthResourceServerAuthorizationDetailsType returns detailed metadata for
// a single published authorization_details type.
func (c *Sys) ReadOAuthResourceServerAuthorizationDetailsType(typeName string) (*OAuthResourceServerAuthorizationDetailsTypeDetail, error) {
	return c.ReadOAuthResourceServerAuthorizationDetailsTypeWithContext(context.Background(), typeName)
}

// ReadOAuthResourceServerAuthorizationDetailsTypeWithContext returns detailed
// metadata for a single published authorization_details type.
func (c *Sys) ReadOAuthResourceServerAuthorizationDetailsTypeWithContext(ctx context.Context, typeName string) (*OAuthResourceServerAuthorizationDetailsTypeDetail, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	req := c.c.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", oAuthResourceServerAuthorizationDetailsTypesPath, typeName))
	secret, err := c.readSecretWithContext(ctx, req)
	if err != nil {
		return nil, err
	}

	var result OAuthResourceServerAuthorizationDetailsTypeDetail
	if err := mapstructure.Decode(secret.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ReadOAuthResourceServerAuthorizationDetailsTypeSchema returns the
// machine-readable JSON schema for a single published authorization_details
// type.
func (c *Sys) ReadOAuthResourceServerAuthorizationDetailsTypeSchema(typeName string) (*OAuthResourceServerAuthorizationDetailsTypeSchema, error) {
	return c.ReadOAuthResourceServerAuthorizationDetailsTypeSchemaWithContext(context.Background(), typeName)
}

// ReadOAuthResourceServerAuthorizationDetailsTypeSchemaWithContext returns the
// machine-readable JSON schema for a single published authorization_details
// type.
func (c *Sys) ReadOAuthResourceServerAuthorizationDetailsTypeSchemaWithContext(ctx context.Context, typeName string) (*OAuthResourceServerAuthorizationDetailsTypeSchema, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	req := c.c.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/schema", oAuthResourceServerAuthorizationDetailsTypesPath, typeName))
	secret, err := c.readSecretWithContext(ctx, req)
	if err != nil {
		return nil, err
	}

	var result OAuthResourceServerAuthorizationDetailsTypeSchema
	if err := mapstructure.Decode(secret.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// readSecretWithContext issues the request and returns a parsed Secret with
// non-empty Data.
func (c *Sys) readSecretWithContext(ctx context.Context, req *Request) (*Secret, error) {
	resp, err := c.c.rawRequestWithContext(ctx, req)
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
