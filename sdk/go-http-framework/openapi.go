package openapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"github.com/fatih/structs"
)

type Extensions map[string]interface{}

func genericMarshal(i interface{}) ([]byte, error) {
	if i == nil {
		// Will print null but cleaner to fall back
		return json.Marshal(i)
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Struct {
		s := structs.Map(i)

		ef := v.FieldByName("Extensions")
		if !ef.IsZero() {
			ext := ef.Interface()
			if ext != nil {
				extMap, ok := ext.(map[string]interface{})
				if !ok {
					panic(fmt.Sprintf("Extensions field has incorrect type %T", ext))
				}
				for k, v := range extMap {
					s[k] = v
				}
			}
		}

		return json.Marshal(s)
	}

	return json.Marshal(i)
}

type OpenAPI struct {
	OpenAPI    string     `structs:"openapi"`
	Info       Info       `structs:"info"`
	Servers    []Server   `structs:"servers,omitempty"`
	Paths      Paths      `structs:"paths"`
	Extensions Extensions `structs:"-"`
}

func (o *OpenAPI) MarshalJSON() ([]byte, error) {
	if o == nil {
		return []byte("null"), nil
	}

	s := structs.Map(o)

	for k, v := range o.Extensions {
		s[k] = v
	}

	return json.Marshal(s)
}

type Info struct {
	Title          string   `structs:"title"`
	Description    *string  `structs:"description,omitempty"`
	TermsOfService *string  `structs:"termsOfService,omitempty"`
	License        *License `structs:"license,omitempty"`
	Version        string   `structs:"version"`
}

type Server struct {
	URL         string `structs:"url"`
	Description string `structs:"description,omitempty"`
}

type License struct {
	Name string  `structs:"name"`
	URL  *string `structs:"url,omitempty"`
}

type Paths map[string]OASPathItem

type OASPathItem struct {
	Description string         `structs:"description,omitempty"`
	Parameters  []OASParameter `structs:"parameters,omitempty"`
	Get         *OASOperation  `structs:"get,omitempty"`
	Post        *OASOperation  `structs:"post,omitempty"`
	Delete      *OASOperation  `structs:"delete,omitempty"`
}

// NewOASOperation creates an empty OpenAPI Operations object.
func NewOASOperation() *OASOperation {
	return &OASOperation{
		Responses: make(map[int]*OASResponse),
	}
}

type OASOperation struct {
	Summary     string               `structs:"summary,omitempty"`
	Description string               `structs:"description,omitempty"`
	OperationID string               `structs:"operationId,omitempty"`
	Tags        []string             `structs:"tags,omitempty"`
	Parameters  []OASParameter       `structs:"parameters,omitempty"`
	RequestBody *OASRequestBody      `structs:"requestBody,omitempty"`
	Responses   map[int]*OASResponse `structs:"responses"`
	Deprecated  bool                 `structs:"deprecated,omitempty"`
}

type OASParameter struct {
	Name        string     `structs:"name"`
	Description string     `structs:"description,omitempty"`
	In          string     `structs:"in"`
	Schema      *OASSchema `structs:"schema,omitempty"`
	Required    bool       `structs:"required,omitempty"`
	Deprecated  bool       `structs:"deprecated,omitempty"`
}

type OASRequestBody struct {
	Description string     `structs:"description,omitempty"`
	Content     OASContent `structs:"content,omitempty"`
}

type OASContent map[string]*OASMediaTypeObject

type OASMediaTypeObject struct {
	Schema *OASSchema `structs:"schema,omitempty"`
}

type OASSchema struct {
	Type        string                `structs:"type,omitempty"`
	Description string                `structs:"description,omitempty"`
	Properties  map[string]*OASSchema `structs:"properties,omitempty"`

	// Required is a list of keys in Properties that are required to be present. This is a different
	// approach than OASParameter (unfortunately), but is how JSONSchema handles 'required'.
	Required []string `structs:"required,omitempty"`

	Items      *OASSchema    `structs:"items,omitempty"`
	Format     string        `structs:"format,omitempty"`
	Pattern    string        `structs:"pattern,omitempty"`
	Enum       []interface{} `structs:"enum,omitempty"`
	Default    interface{}   `structs:"default,omitempty"`
	Example    interface{}   `structs:"example,omitempty"`
	Deprecated bool          `structs:"deprecated,omitempty"`
}

type OASResponse struct {
	Description string     `structs:"description"`
	Content     OASContent `structs:"content,omitempty"`
}

var OASStdRespOK = &OASResponse{
	Description: "OK",
}

var OASStdRespNoContent = &OASResponse{
	Description: "empty body",
}

// Regex for handling optional and named parameters in paths, and string cleanup.
// Predefined here to avoid substantial recompilation.

// Capture optional path elements in ungreedy (?U) fashion
// Both "(leases/)?renew" and "(/(?P<name>.+))?" formats are detected
var optRe = regexp.MustCompile(`(?U)\([^(]*\)\?|\(/\(\?P<[^(]*\)\)\?`)

var reqdRe = regexp.MustCompile(`\(?\?P<(\w+)>[^)]*\)?`)             // Capture required parameters, e.g. "(?P<name>regex)"
var altRe = regexp.MustCompile(`\((.*)\|(.*)\)`)                     // Capture alternation elements, e.g. "(raw/?$|raw/(?P<path>.+))"
var pathFieldsRe = regexp.MustCompile(`{(\w+)}`)                     // Capture OpenAPI-style named parameters, e.g. "lookup/{urltoken}",
var cleanCharsRe = regexp.MustCompile("[()^$?]")                     // Set of regex characters that will be stripped during cleaning
var cleanSuffixRe = regexp.MustCompile(`/\?\$?$`)                    // Path suffix patterns that will be stripped during cleaning
var wsRe = regexp.MustCompile(`\s+`)                                 // Match whitespace, to be compressed during cleaning
var altFieldsGroupRe = regexp.MustCompile(`\(\?P<\w+>\w+(\|\w+)+\)`) // Match named groups that limit options, e.g. "(?<foo>a|b|c)"
var altFieldsRe = regexp.MustCompile(`\w+(\|\w+)+`)                  // Match an options set, e.g. "a|b|c"
var nonWordRe = regexp.MustCompile(`[^\w]+`)                         // Match a sequence of non-word characters
