// Package oas provides structures and helpers that align with a subset of version 2 of the
// OpenAPI specification: https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md
package oas

import (
	"regexp"

	"github.com/hashicorp/vault/version"
)

type OASDoc struct {
	Swagger string     `json:"swagger"`
	Info    Info       `json:"info"`
	Paths   PathTopMap `json:"paths"`
}

func NewOASDoc() OASDoc {
	return OASDoc{
		Swagger: "2.0",
		Info: Info{
			Title:   "HashiCorp Vault API",
			Version: version.GetVersion().Version,
		},
		Paths: make(PathTopMap),
	}
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type PathTopMap map[string]*PathMethods

type MethodDetail struct {
	Summary     string           `json:"summary"`
	Description string           `json:"description"`
	Produces    []string         `json:"produces"`
	Parameters  []Parameter      `json:"parameters,omitempty"`
	Responses   map[int]Response `json:"responses"`
}

func NewMethodDetail() *MethodDetail {
	return &MethodDetail{
		Responses: make(map[int]Response),
		Produces:  []string{"application/json"},
	}
}

type PathMethods struct {
	Get    *MethodDetail `json:"get,omitempty"`
	Post   *MethodDetail `json:"post,omitempty"`
	Delete *MethodDetail `json:"delete,omitempty"`
	Root   bool          `json:"x-vault-root"`
}

type Parameter struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	In          string            `json:"in"`
	Type        string            `json:"type,omitempty"`
	Items       *Property         `json:"items,omitempty"`
	Schema      *Schema           `json:"schema,omitempty"`
	Required    bool              `json:"required,omitempty"`
	Attrs       map[string]string `json:"x-attrs,omitempty"`
}

type Property struct {
	Type        string            `json:"type"`
	Description string            `json:"description,omitempty"`
	Items       *Property         `json:"items,omitempty"`
	Format      string            `json:"format,omitempty"`
	Attrs       map[string]string `json:"x-attrs,omitempty"`
}

type Schema struct {
	Type       string               `json:"type"`
	Properties map[string]*Property `json:"properties"`
}

type Response struct {
	Description string `json:"description"`
	Example     string `json:"example,omitempty"`
}

var StdRespOK = Response{
	Description: "OK",
}

var StdRespNoContent = Response{
	Description: "empty body",
}

func NewSchema() *Schema {
	return &Schema{
		Type:       "object",
		Properties: make(map[string]*Property),
	}
}

var pathFieldsRe = regexp.MustCompile(`{(\w+)}`)

// PathFields extracts named parameters from an OAS path
func PathFields(pattern string) []string {
	r := pathFieldsRe.FindAllStringSubmatch(pattern, -1)
	ret := make([]string, 0, len(r))
	for _, t := range r {
		ret = append(ret, t[1])
	}
	return ret
}
