package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/textproto"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/structs"
)

var (
	ErrDuplicateHeader             = errors.New("duplicate header encountered in header definition map")
	ErrInvalidHeaderTarget         = errors.New("target in headers map not header object or reference object")
	ErrInvalidExampleTarget        = errors.New("target in examples map not example object or reference object")
	ErrInvalidPath                 = errors.New("invalid path supplied in slugs in Paths object")
	ErrInvalidParameterType        = errors.New("entry in parameters slice has incorrect type")
	ErrInvalidRequestBodyType      = errors.New("request body has incorrect type")
	ErrInvalidStatusCode           = errors.New("status code is invalid")
	ErrInvalidResponsesDefaultType = errors.New(`"default" value in responses has incorrect type`)
	ErrInvalidSchemaType           = errors.New(`"schema" value in parameter or header object has incorrect type`)
	ErrPathRequiredValueIncorrect  = errors.New(`'required' not set and true when "in" is "path"`)
	ErrBadContentMapLength         = errors.New(`"content" map contains too many entries`)
	ErrSchemaAndContent            = errors.New(`both "schema" and "content" found in parameter or header object`)
)

func badParamInSomethingError(param, thing string) error {
	return fmt.Errorf(`%q parameter in %q object is invalid`, param, thing)
}

type Extensions map[string]interface{}

func addExtensionsToMap(v reflect.Value, m map[string]interface{}) {
	ef := v.FieldByName("Extensions")
	if ef.IsValid() && !ef.IsZero() {
		ext := ef.Interface()
		if ext != nil {
			extMap, ok := ext.(Extensions)
			if !ok {
				panic(fmt.Sprintf("Extensions field has incorrect type %T", ext))
			}
			for k, v := range extMap {
				m[k] = v
			}
		}
	}
}

func genericMarshal(i interface{}) ([]byte, error) {
	if i == nil {
		// Will print null but cleaner to fall back
		return json.Marshal(i)
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	if v.Kind() == reflect.Struct {
		s := structs.Map(i)
		addExtensionsToMap(v, s)

		return json.Marshal(s)
	}

	return json.Marshal(i)
}

type OpenAPI struct {
	OpenAPI string   `structs:"openapi"`
	Info    Info     `structs:"info"`
	Servers []Server `structs:"servers,omitempty"`
	Paths   Paths    `structs:"paths"`

	Extensions Extensions `structs:"-"`
}

func (o OpenAPI) MarshalJSON() ([]byte, error) {
	if o.OpenAPI == "" {
		return nil, badParamInSomethingError("openapi", "OpenAPI")
	}

	return genericMarshal(o)
}

type Info struct {
	Title          string   `structs:"title"`
	Description    *string  `structs:"description,omitempty"`
	TermsOfService *string  `structs:"termsOfService,omitempty"`
	License        *License `structs:"license,omitempty"`
	Version        string   `structs:"version"`

	Extensions Extensions `structs:"-"`
}

func (i Info) MarshalJSON() ([]byte, error) {
	switch {
	case i.Title == "":
		return nil, badParamInSomethingError("title", "Info")

	case i.Version == "":
		return nil, badParamInSomethingError("version", "Info")
	}

	return genericMarshal(i)
}

type License struct {
	Name string  `structs:"name"`
	URL  *string `structs:"url,omitempty"`
}

func (l License) MarshalJSON() ([]byte, error) {
	if l.Name == "" {
		return nil, badParamInSomethingError("name", "License")
	}

	return genericMarshal(l)
}

type Server struct {
	URL         string  `structs:"url"`
	Description *string `structs:"description,omitempty"`

	Extensions Extensions `structs:"-"`
}

func (s Server) MarshalJSON() ([]byte, error) {
	if s.URL == "" {
		return nil, badParamInSomethingError("url", "Server")
	}

	return genericMarshal(s)
}

type ExternalDocumentation struct {
	Description *string `structs:"description,omitempty"`
	URL         string  `structs:"url"`

	Extensions Extensions `structs:"-"`
}

func (e ExternalDocumentation) MarshalJSON() ([]byte, error) {
	if e.URL == "" {
		return nil, badParamInSomethingError("url", "External Documentation")
	}

	return genericMarshal(e)
}

type HeadersMap map[string]interface{}

func (h HeadersMap) Validate() (HeadersMap, error) {
	if len(h) == 0 {
		return h, nil
	}

	dupMap := make(HeadersMap, len(h))

	for k, v := range h {
		canonKey := textproto.CanonicalMIMEHeaderKey(k)
		if canonKey == "Content-Type" {
			continue
		}
		if _, ok := dupMap[canonKey]; ok {
			return h, ErrDuplicateHeader
		}
		switch v.(type) {
		case Reference, *Reference, Header, *Header:
		default:
			return h, ErrInvalidHeaderTarget
		}
		dupMap[canonKey] = v
	}

	return dupMap, nil
}

type Response struct {
	Description string               `structs:"description"`
	Headers     HeadersMap           `structs:"headers,omitempty"`
	Content     map[string]MediaType `structs:"content,omitempty"`

	Extensions Extensions `structs:"-"`
}

// MarshalJSON marshals a response. Note that it can perform transormations on
// the Response object as required by the spec. For instance, it will elide
// response headers of "Content-Type".
func (r *Response) MarshalJSON() ([]byte, error) {
	if r.Description == "" {
		return nil, badParamInSomethingError("description", "Response")
	}

	var err error
	r.Headers, err = r.Headers.Validate()
	if err != nil {
		return nil, err
	}

	return genericMarshal(r)
}

type Reference struct {
	Ref string `structs:"$ref"`
}

func (r *Reference) MarshalJSON() ([]byte, error) {
	if r.Ref == "" {
		return nil, badParamInSomethingError("$ref", "Reference")
	}

	return genericMarshal(r)
}

type MediaType struct {
	Schema   interface{}            `structs:"schema,omitempty"`
	Example  interface{}            `structs:"example,omitempty"`
	Examples map[string]interface{} `structs:"examples,omitempty"`
	Encoding map[string]Encoding    `structs:"encoding,omitempty"`

	Extensions Extensions `structs:"-"`
}

func (m MediaType) MarshalJSON() ([]byte, error) {
	for _, v := range m.Examples {
		switch v.(type) {
		case Example, *Example, Reference, *Reference:
		default:
			return nil, ErrInvalidExampleTarget
		}
	}

	return genericMarshal(m)
}

type Example struct {
	Summary       *string     `structs:"summary,omitempty"`
	Description   *string     `structs:"description,omitempty"`
	Value         interface{} `structs:"value,omitempty"`
	ExternalValue *string     `structs:"externalValue,omitempty"`

	Extensions Extensions `structs:"-"`
}

func (e Example) MarshalJSON() ([]byte, error) {
	return genericMarshal(e)
}

type Encoding struct {
	ContentType   *string    `structs:"contentType,omitempty"`
	Headers       HeadersMap `structs:"headers,omitempty"`
	Style         *string    `structs:"style,omitempty"`
	Explode       *bool      `structs:"explode,omitempty"`
	AllowReserved *bool      `structs:"allowReserved,omitempty"`

	Extensions Extensions `structs:"-"`
}

func (e *Encoding) MarshalJSON() ([]byte, error) {
	var err error
	e.Headers, err = e.Headers.Validate()
	if err != nil {
		return nil, err
	}

	return genericMarshal(e)
}

type Header struct {
	Description     *string                `structs:"description,omitempty"`
	Required        *bool                  `structs:"required,omitempty"`
	Deprecated      *bool                  `structs:"deprecated,omitempty"`
	AllowEmptyValue *bool                  `structs:"allowEmptyValue,omitempty"`
	Style           *string                `structs:"style,omitempty"`
	Explode         *bool                  `structs:"explode,omitempty"`
	AllowReserved   *bool                  `structs:"allowReserved,omitempty"`
	Schema          interface{}            `structs:"schema,omitempty"`
	Example         interface{}            `structs:"example,omitempty"`
	Examples        map[string]interface{} `structs:"examples,omitempty"`
	Content         map[string]MediaType   `structs:"content,omitempty"`

	Extensions Extensions `structs:"-"`
}

func (h Header) MarshalJSON() ([]byte, error) {
	if h.Schema != nil {
		switch h.Schema.(type) {
		// FIXME
		//case Schema, *Schema, Reference, *Reference:
		case Reference, *Reference:
		default:
			return nil, ErrInvalidSchemaType
		}
	}

	for _, v := range h.Examples {
		switch v.(type) {
		case Example, *Example, Reference, *Reference:
		default:
			return nil, ErrInvalidExampleTarget
		}
	}

	switch len(h.Content) {
	case 0:
	case 1:
		if h.Schema != nil {
			return nil, ErrSchemaAndContent
		}
	default:
		return nil, ErrBadContentMapLength
	}

	return genericMarshal(h)
}

// Paths corresponds to the Paths Object
//
// Note: because of the capability for the Paths object to have Extensions, it
// can't be implemented as a simple map. However serialization requires the
// items to be pulled into the top level. So although Slugs is a map, the items
// in Slugs (similar to Extensions) will be pulled into the top level at
// serialization time.
type Paths struct {
	Slugs map[string]PathItem `structs:"-"`

	Extensions Extensions `structs:"-"`
}

func (p Paths) MarshalJSON() ([]byte, error) {
	paths := make(map[string]interface{}, len(p.Slugs)+len(p.Extensions))
	for k, v := range p.Slugs {
		if !strings.HasPrefix(k, "/") {
			return nil, ErrInvalidPath
		}
		paths[k] = v
	}

	for k, v := range p.Extensions {
		paths[k] = v
	}

	return json.Marshal(paths)
}

type PathItem struct {
	Ref         *string       `structs:"$ref,omitempty"`
	Summary     *string       `structs:"summary,omitempty"`
	Description *string       `structs:"description,omitempty"`
	Get         *Operation    `structs:"get,omitempty"`
	Put         *Operation    `structs:"put,omitempty"`
	Post        *Operation    `structs:"post,omitempty"`
	Delete      *Operation    `structs:"delete,omitempty"`
	Options     *Operation    `structs:"options,omitempty"`
	Head        *Operation    `structs:"head,omitempty"`
	Patch       *Operation    `structs:"patch,omitempty"`
	Trace       *Operation    `structs:"trace,omitempty"`
	Servers     []Server      `structs:"servers,omitempty"`
	Parameters  []interface{} `structs:"parameters,omitempty"`

	Extensions Extensions `structs:"-"`
}

func (p PathItem) MarshalJSON() ([]byte, error) {
	for _, v := range p.Parameters {
		switch v.(type) {
		case Parameter, *Parameter, Reference, *Reference:
		default:
			return nil, ErrInvalidParameterType
		}
	}

	return genericMarshal(p)
}

type Parameter struct {
	Name            string                 `structs:"name"`
	In              string                 `structs:"in"`
	Description     *string                `structs:"description,omitempty"`
	Required        *bool                  `structs:"required,omitempty"`
	Deprecated      *bool                  `structs:"deprecated,omitempty"`
	AllowEmptyValue *bool                  `structs:"allowEmptyValue,omitempty"`
	Style           *string                `structs:"style,omitempty"`
	Explode         *bool                  `structs:"explode,omitempty"`
	AllowReserved   *bool                  `structs:"allowReserved,omitempty"`
	Schema          interface{}            `structs:"schema,omitempty"`
	Example         interface{}            `structs:"example,omitempty"`
	Examples        map[string]interface{} `structs:"examples,omitempty"`
	Content         map[string]MediaType   `structs:"content,omitempty"`

	Extensions Extensions `structs:"-"`
}

func (p Parameter) MarshalJSON() ([]byte, error) {
	switch {
	case p.Name == "":
		return nil, badParamInSomethingError("name", "Parameter")

	case p.In == "":
		return nil, badParamInSomethingError("in", "Parameter")
	}

	switch p.In {
	case "query", "header", "cookie":

	case "path":
		if p.Required == nil || !*p.Required {
			return nil, ErrPathRequiredValueIncorrect
		}

	default:
		return nil, badParamInSomethingError("in", "Parameter")
	}

	if p.Schema != nil {
		switch p.Schema.(type) {
		//FIXME
		//case Schema, *Schema, Reference, *Reference:
		case Reference, *Reference:
		default:
			return nil, ErrInvalidSchemaType
		}
	}

	for _, v := range p.Examples {
		switch v.(type) {
		case Example, *Example, Reference, *Reference:
		default:
			return nil, ErrInvalidExampleTarget
		}
	}

	switch len(p.Content) {
	case 0:
	case 1:
		if p.Schema != nil {
			return nil, ErrSchemaAndContent
		}
	default:
		return nil, ErrBadContentMapLength
	}

	return genericMarshal(p)
}

type Operation struct {
	Tags         []string               `structs:"tags,omitempty"`
	Summary      *string                `structs:"summary,omitempty"`
	Description  *string                `structs:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `structs:"externalDocs,omitempty"`
	OperationID  *string                `structs:"operationId,omitempty"`
	Parameters   []interface{}          `structs:"parameters,omitempty"`
	RequestBody  interface{}            `structs:"requestBody,omitempty"`
	Responses    Responses              `structs:"responses,omitempty"`
	Deprecated   *bool                  `structs:"deprecated,omitempty"`
	Servers      []Server               `structs:"servers,omitempty"`

	Extensions Extensions `structs:"-"`
}

func (o Operation) MarshalJSON() ([]byte, error) {
	for _, v := range o.Parameters {
		switch v.(type) {
		case Parameter, *Parameter, Reference, *Reference:
		default:
			return nil, ErrInvalidParameterType
		}
	}

	if o.RequestBody != nil {
		switch o.RequestBody.(type) {
		case RequestBody, *RequestBody, Reference, *Reference:
		default:
			return nil, ErrInvalidRequestBodyType
		}
	}

	return genericMarshal(o)
}

type Responses struct {
	Statuses map[string]interface{} `structs:"-"`
	Default  interface{}            `structs:"-"`

	Extensions Extensions `structs:"-"`
}

func (r Responses) MarshalJSON() ([]byte, error) {
	statuses := make(map[string]interface{}, len(r.Statuses)+len(r.Extensions)+1)

	if r.Default != nil {
		switch r.Default.(type) {
		case Response, *Response, Reference, *Reference:
		default:
			return nil, ErrInvalidResponsesDefaultType
		}
	}

	for k, v := range r.Statuses {
		switch v.(type) {
		case Response, *Response, Reference, *Reference:
		default:
			return nil, ErrInvalidResponsesDefaultType
		}

		switch k {
		case "1XX", "2XX", "3XX", "4XX", "5XX":
		default:
			intVal, err := strconv.ParseInt(k, 10, 16)
			if err != nil {
				return nil, ErrInvalidStatusCode
			}

			if intVal < 100 || intVal > 599 {
				return nil, ErrInvalidStatusCode
			}
		}

		statuses[k] = v
	}

	for k, v := range r.Extensions {
		statuses[k] = v
	}

	return json.Marshal(statuses)
}

type RequestBody struct {
	Description *string              `structs:"description,omitempty"`
	Content     map[string]MediaType `structs:"content"`
	Required    *bool                `structs:"required,omitempty"`
}

func (r RequestBody) MarshalJSON() ([]byte, error) {
	if len(r.Content) == 0 {
		return nil, badParamInSomethingError("content", "Request Body")
	}

	return genericMarshal(r)
}
