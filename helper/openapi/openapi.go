// Package openapi provides structures and helpers that align with a subset of version 3 of the
// OpenAPI specification: https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md
package openapi

import (
	"github.com/hashicorp/vault/version"
)

const Version = "3.0.2"

// NewDocument returns an empty OpenAPI document.
func NewDocument() *Document {
	return &Document{
		OpenAPIVersion: Version,
		Info: Info{
			Title:       "HashiCorp Vault API",
			Description: "HTTP API that gives you full access to Vault. All API routes are prefixed with `/v1/`.",
			Version:     version.GetVersion().Version,
			License: License{
				Name: "Mozilla Public License 2.0",
				URL:  "https://www.mozilla.org/en-US/MPL/2.0",
			},
		},
		Paths: make(map[string]*PathItem),
	}
}

type Document struct {
	OpenAPIVersion string               `json:"openapi"`
	Info           Info                 `json:"info"`
	Paths          map[string]*PathItem `json:"paths"`
}

type Info struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Version     string  `json:"version"`
	License     License `json:"license"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PathItem struct {
	Description     string      `json:"description,omitempty"`
	Parameters      []Parameter `json:"parameters,omitempty"`
	Sudo            bool        `json:"x-vault-sudo,omitempty" mapstructure:"x-vault-sudo"`
	Unauthenticated bool        `json:"x-vault-unauthenticated,omitempty" mapstructure:"x-vault-unauthenticated"`
	Create          bool        `json:"x-vault-create,omitempty" mapstructure:"x-vault-create"`

	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
}

// NewOperation creates an empty OpenAPI Operations object.
func NewOperation() *Operation {
	return &Operation{
		Responses: make(map[string]*Response),
	}
}

type Operation struct {
	Summary     string               `json:"summary,omitempty"`
	Description string               `json:"description,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Parameters  []Parameter          `json:"parameters,omitempty"`
	RequestBody *RequestBody         `json:"requestBody,omitempty"`
	Responses   map[string]*Response `json:"responses"`
	Deprecated  bool                 `json:"deprecated,omitempty"`
}

type Parameter struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	In          string  `json:"in"`
	Schema      *Schema `json:"schema,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Deprecated  bool    `json:"deprecated,omitempty"`
}

type RequestBody struct {
	Description string   `json:"description,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Content     *Content `json:"content,omitempty"`
}

type Content map[string]*MediaTypeObject

type MediaTypeObject struct {
	Schema *Schema `json:"schema,omitempty"`
}

type Schema struct {
	Type        string             `json:"type,omitempty"`
	Description string             `json:"description,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
	Format      string             `json:"format,omitempty"`
	Example     interface{}        `json:"example,omitempty"`
	Deprecated  bool               `json:"deprecated,omitempty"`
}

type Response struct {
	Description string   `json:"description"`
	Content     *Content `json:"content,omitempty"`
}

var StdRespOK = &Response{
	Description: "OK",
}

var StdRespNoContent = &Response{
	Description: "empty body",
}
