package framework

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"regexp/syntax"
	"sort"
	"strconv"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

// OpenAPI specification (OAS): https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md
const OASVersion = "3.0.2"

// NewOASDocument returns an empty OpenAPI document.
func NewOASDocument(version string) *OASDocument {
	return &OASDocument{
		Version: OASVersion,
		Info: OASInfo{
			Title:       "HashiCorp Vault API",
			Description: "HTTP API that gives you full access to Vault. All API routes are prefixed with `/v1/`.",
			Version:     version,
			License: OASLicense{
				Name: "Mozilla Public License 2.0",
				URL:  "https://www.mozilla.org/en-US/MPL/2.0",
			},
		},
		Paths: make(map[string]*OASPathItem),
		Components: OASComponents{
			Schemas: make(map[string]*OASSchema),
		},
	}
}

// NewOASDocumentFromMap builds an OASDocument from an existing map version of a document.
// If a document has been decoded from JSON or received from a plugin, it will be as a map[string]interface{}
// and needs special handling beyond the default mapstructure decoding.
func NewOASDocumentFromMap(input map[string]interface{}) (*OASDocument, error) {
	// The Responses map uses integer keys (the response code), but once translated into JSON
	// (e.g. during the plugin transport) these become strings. mapstructure will not coerce these back
	// to integers without a custom decode hook.
	decodeHook := func(src reflect.Type, tgt reflect.Type, inputRaw interface{}) (interface{}, error) {
		// Only alter data if:
		//  1. going from string to int
		//  2. string represent an int in status code range (100-599)
		if src.Kind() == reflect.String && tgt.Kind() == reflect.Int {
			if input, ok := inputRaw.(string); ok {
				if intval, err := strconv.Atoi(input); err == nil {
					if intval >= 100 && intval < 600 {
						return intval, nil
					}
				}
			}
		}
		return inputRaw, nil
	}

	doc := new(OASDocument)

	config := &mapstructure.DecoderConfig{
		DecodeHook: decodeHook,
		Result:     doc,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(input); err != nil {
		return nil, err
	}

	return doc, nil
}

type OASDocument struct {
	Version    string                  `json:"openapi" mapstructure:"openapi"`
	Info       OASInfo                 `json:"info"`
	Paths      map[string]*OASPathItem `json:"paths"`
	Components OASComponents           `json:"components"`
}

type OASComponents struct {
	Schemas map[string]*OASSchema `json:"schemas"`
}

type OASInfo struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
	License     OASLicense `json:"license"`
}

type OASLicense struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type OASPathItem struct {
	Description       string             `json:"description,omitempty"`
	Parameters        []OASParameter     `json:"parameters,omitempty"`
	Sudo              bool               `json:"x-vault-sudo,omitempty" mapstructure:"x-vault-sudo"`
	Unauthenticated   bool               `json:"x-vault-unauthenticated,omitempty" mapstructure:"x-vault-unauthenticated"`
	CreateSupported   bool               `json:"x-vault-createSupported,omitempty" mapstructure:"x-vault-createSupported"`
	DisplayNavigation bool               `json:"x-vault-displayNavigation,omitempty" mapstructure:"x-vault-displayNavigation"`
	DisplayAttrs      *DisplayAttributes `json:"x-vault-displayAttrs,omitempty" mapstructure:"x-vault-displayAttrs"`

	Get    *OASOperation `json:"get,omitempty"`
	Post   *OASOperation `json:"post,omitempty"`
	Delete *OASOperation `json:"delete,omitempty"`
}

// NewOASOperation creates an empty OpenAPI Operations object.
func NewOASOperation() *OASOperation {
	return &OASOperation{
		Responses: make(map[int]*OASResponse),
	}
}

type OASOperation struct {
	Summary     string               `json:"summary,omitempty"`
	Description string               `json:"description,omitempty"`
	OperationID string               `json:"operationId,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Parameters  []OASParameter       `json:"parameters,omitempty"`
	RequestBody *OASRequestBody      `json:"requestBody,omitempty"`
	Responses   map[int]*OASResponse `json:"responses"`
	Deprecated  bool                 `json:"deprecated,omitempty"`
}

type OASParameter struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	In          string     `json:"in"`
	Schema      *OASSchema `json:"schema,omitempty"`
	Required    bool       `json:"required,omitempty"`
	Deprecated  bool       `json:"deprecated,omitempty"`
}

type OASRequestBody struct {
	Description string     `json:"description,omitempty"`
	Required    bool       `json:"required,omitempty"`
	Content     OASContent `json:"content,omitempty"`
}

type OASContent map[string]*OASMediaTypeObject

type OASMediaTypeObject struct {
	Schema *OASSchema `json:"schema,omitempty"`
}

type OASSchema struct {
	Ref         string                `json:"$ref,omitempty"`
	Type        string                `json:"type,omitempty"`
	Description string                `json:"description,omitempty"`
	Properties  map[string]*OASSchema `json:"properties,omitempty"`

	// Required is a list of keys in Properties that are required to be present. This is a different
	// approach than OASParameter (unfortunately), but is how JSONSchema handles 'required'.
	Required []string `json:"required,omitempty"`

	Items      *OASSchema    `json:"items,omitempty"`
	Format     string        `json:"format,omitempty"`
	Pattern    string        `json:"pattern,omitempty"`
	Enum       []interface{} `json:"enum,omitempty"`
	Default    interface{}   `json:"default,omitempty"`
	Example    interface{}   `json:"example,omitempty"`
	Deprecated bool          `json:"deprecated,omitempty"`
	// DisplayName      string             `json:"x-vault-displayName,omitempty" mapstructure:"x-vault-displayName,omitempty"`
	DisplayValue     interface{}        `json:"x-vault-displayValue,omitempty" mapstructure:"x-vault-displayValue,omitempty"`
	DisplaySensitive bool               `json:"x-vault-displaySensitive,omitempty" mapstructure:"x-vault-displaySensitive,omitempty"`
	DisplayGroup     string             `json:"x-vault-displayGroup,omitempty" mapstructure:"x-vault-displayGroup,omitempty"`
	DisplayAttrs     *DisplayAttributes `json:"x-vault-displayAttrs,omitempty" mapstructure:"x-vault-displayAttrs,omitempty"`
}

type OASResponse struct {
	Description string     `json:"description"`
	Content     OASContent `json:"content,omitempty"`
}

var OASStdRespOK = &OASResponse{
	Description: "OK",
}

var OASStdRespNoContent = &OASResponse{
	Description: "empty body",
}

// Regex for handling fields in paths, and string cleanup.
// Predefined here to avoid substantial recompilation.

var (
	nonWordRe    = regexp.MustCompile(`[^\w]+`)  // Match a sequence of non-word characters
	pathFieldsRe = regexp.MustCompile(`{(\w+)}`) // Capture OpenAPI-style named parameters, e.g. "lookup/{urltoken}",
	wsRe         = regexp.MustCompile(`\s+`)     // Match whitespace, to be compressed during cleaning
)

// documentPaths parses all paths in a framework.Backend into OpenAPI paths.
func documentPaths(backend *Backend, requestResponsePrefix string, doc *OASDocument) error {
	for _, p := range backend.Paths {
		if err := documentPath(p, backend.SpecialPaths(), requestResponsePrefix, backend.BackendType, doc); err != nil {
			return err
		}
	}

	return nil
}

// documentPath parses a framework.Path into one or more OpenAPI paths.
func documentPath(p *Path, specialPaths *logical.Paths, requestResponsePrefix string, backendType logical.BackendType, doc *OASDocument) error {
	var sudoPaths []string
	var unauthPaths []string

	if specialPaths != nil {
		sudoPaths = specialPaths.Root
		unauthPaths = specialPaths.Unauthenticated
	}

	// Convert optional parameters into distinct patterns to be processed independently.
	forceUnpublished := false
	paths, err := expandPattern(p.Pattern)
	if err != nil {
		if errors.Is(err, errUnsupportableRegexpOperationForOpenAPI) {
			// Pattern cannot be transformed into sensible OpenAPI paths. In this case, we override the later
			// processing to use the regexp, as is, as the path, and behave as if Unpublished was set on every
			// operation (meaning the operations will not be represented in the OpenAPI document).
			//
			// This allows a human reading the OpenAPI document to notice that, yes, a path handler does exist,
			// even though it was not able to contribute actual OpenAPI operations.
			forceUnpublished = true
			paths = []string{p.Pattern}
		} else {
			return err
		}
	}

	for pathIndex, path := range paths {
		// Construct a top level PathItem which will be populated as the path is processed.
		pi := OASPathItem{
			Description: cleanString(p.HelpSynopsis),
		}

		pi.Sudo = specialPathMatch(path, sudoPaths)
		pi.Unauthenticated = specialPathMatch(path, unauthPaths)
		pi.DisplayAttrs = p.DisplayAttrs

		// If the newer style Operations map isn't defined, create one from the legacy fields.
		operations := p.Operations
		if operations == nil {
			operations = make(map[logical.Operation]OperationHandler)

			for opType, cb := range p.Callbacks {
				operations[opType] = &PathOperation{
					Callback: cb,
					Summary:  p.HelpSynopsis,
				}
			}
		}

		// Process path and header parameters, which are common to all operations.
		// Body fields will be added to individual operations.
		pathFields, bodyFields := splitFields(p.Fields, path)

		for name, field := range pathFields {
			location := "path"
			required := true

			if field == nil {
				continue
			}

			if field.Query {
				location = "query"
				required = false
			}

			t := convertType(field.Type)
			p := OASParameter{
				Name:        name,
				Description: cleanString(field.Description),
				In:          location,
				Schema: &OASSchema{
					Type:         t.baseType,
					Pattern:      t.pattern,
					Enum:         field.AllowedValues,
					Default:      field.Default,
					DisplayAttrs: field.DisplayAttrs,
				},
				Required:   required,
				Deprecated: field.Deprecated,
			}
			pi.Parameters = append(pi.Parameters, p)
		}

		// Sort parameters for a stable output
		sort.Slice(pi.Parameters, func(i, j int) bool {
			return strings.ToLower(pi.Parameters[i].Name) < strings.ToLower(pi.Parameters[j].Name)
		})

		// Process each supported operation by building up an Operation object
		// with descriptions, properties and examples from the framework.Path data.
		for opType, opHandler := range operations {
			props := opHandler.Properties()
			if props.Unpublished || forceUnpublished {
				continue
			}

			if opType == logical.CreateOperation {
				pi.CreateSupported = true

				// If both Create and Update are defined, only process Update.
				if operations[logical.UpdateOperation] != nil {
					continue
				}
			}

			// If both List and Read are defined, only process Read.
			if opType == logical.ListOperation && operations[logical.ReadOperation] != nil {
				continue
			}

			op := NewOASOperation()

			operationID := constructOperationID(
				path,
				pathIndex,
				p.DisplayAttrs,
				opType,
				props.DisplayAttrs,
				requestResponsePrefix,
			)

			op.Summary = props.Summary
			op.Description = props.Description
			op.Deprecated = props.Deprecated
			op.OperationID = operationID

			// Add any fields not present in the path as body parameters for POST.
			if opType == logical.CreateOperation || opType == logical.UpdateOperation {
				s := &OASSchema{
					Type:       "object",
					Properties: make(map[string]*OASSchema),
					Required:   make([]string, 0),
				}

				for name, field := range bodyFields {
					// Removing this field from the spec as it is deprecated in favor of using "sha256"
					// The duplicate sha_256 and sha256 in these paths cause issues with codegen
					if name == "sha_256" && strings.Contains(path, "plugins/catalog/") {
						continue
					}

					openapiField := convertType(field.Type)
					if field.Required {
						s.Required = append(s.Required, name)
					}

					p := OASSchema{
						Type:         openapiField.baseType,
						Description:  cleanString(field.Description),
						Format:       openapiField.format,
						Pattern:      openapiField.pattern,
						Enum:         field.AllowedValues,
						Default:      field.Default,
						Deprecated:   field.Deprecated,
						DisplayAttrs: field.DisplayAttrs,
					}
					if openapiField.baseType == "array" {
						p.Items = &OASSchema{
							Type: openapiField.items,
						}
					}
					s.Properties[name] = &p
				}

				// If examples were given, use the first one as the sample
				// of this schema.
				if len(props.Examples) > 0 {
					s.Example = props.Examples[0].Data
				}

				// Set the final request body. Only JSON request data is supported.
				if len(s.Properties) > 0 || s.Example != nil {
					requestName := operationID + "-request"
					doc.Components.Schemas[requestName] = s
					op.RequestBody = &OASRequestBody{
						Required: true,
						Content: OASContent{
							"application/json": &OASMediaTypeObject{
								Schema: &OASSchema{Ref: fmt.Sprintf("#/components/schemas/%s", requestName)},
							},
						},
					}
				}
			}

			// LIST is represented as GET with a `list` query parameter.
			if opType == logical.ListOperation {
				// Only accepts List (due to the above skipping of ListOperations that also have ReadOperations)
				op.Parameters = append(op.Parameters, OASParameter{
					Name:        "list",
					Description: "Must be set to `true`",
					Required:    true,
					In:          "query",
					Schema:      &OASSchema{Type: "string", Enum: []interface{}{"true"}},
				})
			} else if opType == logical.ReadOperation && operations[logical.ListOperation] != nil {
				// Accepts both Read and List
				op.Parameters = append(op.Parameters, OASParameter{
					Name:        "list",
					Description: "Return a list if `true`",
					In:          "query",
					Schema:      &OASSchema{Type: "string"},
				})
			}

			// Add tags based on backend type
			var tags []string
			switch backendType {
			case logical.TypeLogical:
				tags = []string{"secrets"}
			case logical.TypeCredential:
				tags = []string{"auth"}
			}

			op.Tags = append(op.Tags, tags...)

			// Set default responses.
			if len(props.Responses) == 0 {
				if opType == logical.DeleteOperation {
					op.Responses[204] = OASStdRespNoContent
				} else {
					op.Responses[200] = OASStdRespOK
				}
			}

			// Add any defined response details.
			for code, responses := range props.Responses {
				var description string
				content := make(OASContent)

				for i, resp := range responses {
					if i == 0 {
						description = resp.Description
					}
					if resp.Example != nil {
						mediaType := resp.MediaType
						if mediaType == "" {
							mediaType = "application/json"
						}

						// create a version of the response that will not emit null items
						cr := cleanResponse(resp.Example)

						// Only one example per media type is allowed, so first one wins
						if _, ok := content[mediaType]; !ok {
							content[mediaType] = &OASMediaTypeObject{
								Schema: &OASSchema{
									Example: cr,
								},
							}
						}
					}

					responseSchema := &OASSchema{
						Type:       "object",
						Properties: make(map[string]*OASSchema),
					}

					for name, field := range resp.Fields {
						openapiField := convertType(field.Type)
						p := OASSchema{
							Type:         openapiField.baseType,
							Description:  cleanString(field.Description),
							Format:       openapiField.format,
							Pattern:      openapiField.pattern,
							Enum:         field.AllowedValues,
							Default:      field.Default,
							Deprecated:   field.Deprecated,
							DisplayAttrs: field.DisplayAttrs,
						}
						if openapiField.baseType == "array" {
							p.Items = &OASSchema{
								Type: openapiField.items,
							}
						}
						responseSchema.Properties[name] = &p
					}

					if len(resp.Fields) != 0 {
						responseName := operationID + "-response"
						doc.Components.Schemas[responseName] = responseSchema
						content = OASContent{
							"application/json": &OASMediaTypeObject{
								Schema: &OASSchema{Ref: fmt.Sprintf("#/components/schemas/%s", responseName)},
							},
						}
					}
				}

				op.Responses[code] = &OASResponse{
					Description: description,
					Content:     content,
				}
			}

			switch opType {
			case logical.CreateOperation, logical.UpdateOperation:
				pi.Post = op
			case logical.ReadOperation, logical.ListOperation:
				pi.Get = op
			case logical.DeleteOperation:
				pi.Delete = op
			}
		}

		doc.Paths["/"+path] = &pi
	}

	return nil
}

func specialPathMatch(path string, specialPaths []string) bool {
	// Test for exact or prefix match of special paths.
	for _, sp := range specialPaths {
		if sp == path ||
			(strings.HasSuffix(sp, "*") && strings.HasPrefix(path, sp[0:len(sp)-1])) {
			return true
		}
	}
	return false
}

// constructOperationID joins the given inputs into a hyphen-separated
// lower-case operation id, which is also used as a prefix for request and
// response names.
//
// The OperationPrefix / -Verb / -Suffix found in display attributes will be
// used, if provided. Otherwise, the function falls back to using the path and
// the operation.
//
// Examples of generated operation identifiers:
//   - kvv2-write
//   - kvv2-read
//   - google-cloud-login
//   - google-cloud-write-role
func constructOperationID(
	path string,
	pathIndex int,
	pathAttributes *DisplayAttributes,
	operation logical.Operation,
	operationAttributes *DisplayAttributes,
	defaultPrefix string,
) string {
	var (
		prefix string
		verb   string
		suffix string
	)

	if operationAttributes != nil {
		prefix = operationAttributes.OperationPrefix
		verb = operationAttributes.OperationVerb
		suffix = operationAttributes.OperationSuffix
	}

	if pathAttributes != nil {
		if prefix == "" {
			prefix = pathAttributes.OperationPrefix
		}
		if verb == "" {
			verb = pathAttributes.OperationVerb
		}
		if suffix == "" {
			suffix = pathAttributes.OperationSuffix
		}
	}

	// A single suffix string can contain multiple pipe-delimited strings. To
	// determine the actual suffix, we attempt to match it by the index of the
	// paths returned from `expandPattern(...)`. For example:
	//
	//  pki/
	//  	Pattern: "keys/generate/(internal|exported|kms)",
	//      DisplayAttrs: {
	//          ...
	//          OperationSuffix: "internal-key|exported-key|kms-key",
	//      },
	//
	//  will expand into three paths and corresponding suffixes:
	//
	//      path 0: "keys/generate/internal"  suffix: internal-key
	//      path 1: "keys/generate/exported"  suffix: exported-key
	//      path 2: "keys/generate/kms"       suffix: kms-key
	//
	if suffixes := strings.Split(suffix, "|"); len(suffixes) > 1 || pathIndex > 0 {
		// if the index is out of bounds, fall back to the old logic
		if pathIndex >= len(suffixes) {
			suffix = ""
		} else {
			suffix = suffixes[pathIndex]
		}
	}

	// hyphenate is a helper that hyphenates the given slice except the empty elements
	hyphenate := func(parts []string) string {
		filtered := make([]string, 0, len(parts))
		for _, e := range parts {
			if e != "" {
				filtered = append(filtered, e)
			}
		}
		return strings.ToLower(strings.Join(filtered, "-"))
	}

	// fall back to using the path + operation to construct the operation id
	var (
		needPrefix = prefix == "" && verb == ""
		needVerb   = verb == ""
		needSuffix = suffix == "" && (verb == "" || pathIndex > 0)
	)

	if needPrefix {
		prefix = defaultPrefix
	}

	if needVerb {
		if operation == logical.UpdateOperation {
			verb = "write"
		} else {
			verb = string(operation)
		}
	}

	if needSuffix {
		suffix = hyphenate(nonWordRe.Split(strings.ToLower(path), -1))
	}

	return hyphenate([]string{prefix, verb, suffix})
}

// expandPattern expands a regex pattern by generating permutations of any optional parameters
// and changing named parameters into their {openapi} equivalents.
func expandPattern(pattern string) ([]string, error) {
	// Happily, the Go regexp library exposes its underlying "parse to AST" functionality, so we can rely on that to do
	// the hard work of interpreting the regexp syntax.
	rx, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		// This should be impossible to reach, since regexps have previously been compiled with MustCompile in
		// Backend.init.
		panic(err)
	}

	paths, err := collectPathsFromRegexpAST(rx)
	if err != nil {
		return nil, err
	}

	return paths, nil
}

type pathCollector struct {
	strings.Builder
	conditionalSlashAppendedAtLength int
}

// collectPathsFromRegexpAST performs a depth-first recursive walk through a regexp AST, collecting an OpenAPI-style
// path as it goes.
//
// Each time it encounters alternation (a|b) or an optional part (a?), it forks its processing to produce additional
// results, to account for each possibility. Note: This does mean that an input pattern with lots of these regexp
// features can produce a lot of different OpenAPI endpoints. At the time of writing, the most complex known example is
//
//	"issuer/" + framework.GenericNameRegex(issuerRefParam) + "/crl(/pem|/der|/delta(/pem|/der)?)?"
//
// in the PKI secrets engine which expands to 6 separate paths.
//
// Each named capture group - i.e. (?P<name>something here) - is replaced with an OpenAPI parameter - i.e. {name} - and
// the subtree of regexp AST inside the parameter is completely skipped.
func collectPathsFromRegexpAST(rx *syntax.Regexp) ([]string, error) {
	pathCollectors, err := collectPathsFromRegexpASTInternal(rx, []*pathCollector{{}})
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(pathCollectors))
	for _, collector := range pathCollectors {
		if collector.conditionalSlashAppendedAtLength != collector.Len() {
			paths = append(paths, collector.String())
		}
	}
	return paths, nil
}

var errUnsupportableRegexpOperationForOpenAPI = errors.New("path regexp uses an operation that cannot be translated to an OpenAPI pattern")

func collectPathsFromRegexpASTInternal(rx *syntax.Regexp, appendingTo []*pathCollector) ([]*pathCollector, error) {
	var err error

	// Depending on the type of this regexp AST node (its Op, i.e. operation), figure out whether it contributes any
	// characters to the URL path, and whether we need to recurse through child AST nodes.
	//
	// Each element of the appendingTo slice tracks a separate path, defined by the alternatives chosen when traversing
	// the | and ? conditional regexp features, and new elements are added as each of these features are traversed.
	//
	// To share this slice across multiple recursive calls of this function, it is passed down as a parameter to each
	// recursive call, potentially modified throughout this switch block, and passed back up as a return value at the
	// end of this function - the parent call uses the return value to update its own local variable.
	switch rx.Op {

	// These AST operations are leaf nodes (no children), that match zero characters, so require no processing at all
	case syntax.OpEmptyMatch: // e.g. (?:)
	case syntax.OpBeginLine: // i.e. ^ when (?m)
	case syntax.OpEndLine: // i.e. $ when (?m)
	case syntax.OpBeginText: // i.e. \A, or ^ when (?-m)
	case syntax.OpEndText: // i.e. \z, or $ when (?-m)
	case syntax.OpWordBoundary: // i.e. \b
	case syntax.OpNoWordBoundary: // i.e. \B

	// OpConcat simply represents multiple parts of the pattern appearing one after the other, so just recurse through
	// those pieces.
	case syntax.OpConcat:
		for _, child := range rx.Sub {
			appendingTo, err = collectPathsFromRegexpASTInternal(child, appendingTo)
			if err != nil {
				return nil, err
			}
		}

	// OpLiteral is a literal string in the pattern - append it to the paths we are building.
	case syntax.OpLiteral:
		for _, collector := range appendingTo {
			collector.WriteString(string(rx.Rune))
		}

	// OpAlternate, i.e. a|b, means we clone all of the pathCollector instances we are currently accumulating paths
	// into, and independently recurse through each alternate option.
	case syntax.OpAlternate: // i.e |
		var totalAppendingTo []*pathCollector
		lastIndex := len(rx.Sub) - 1
		for index, child := range rx.Sub {
			var childAppendingTo []*pathCollector
			if index == lastIndex {
				// Optimization: last time through this loop, we can simply re-use the existing set of pathCollector
				// instances, as we no longer need to preserve them unmodified to make further copies of.
				childAppendingTo = appendingTo
			} else {
				for _, collector := range appendingTo {
					newCollector := new(pathCollector)
					newCollector.WriteString(collector.String())
					newCollector.conditionalSlashAppendedAtLength = collector.conditionalSlashAppendedAtLength
					childAppendingTo = append(childAppendingTo, newCollector)
				}
			}
			childAppendingTo, err = collectPathsFromRegexpASTInternal(child, childAppendingTo)
			if err != nil {
				return nil, err
			}
			totalAppendingTo = append(totalAppendingTo, childAppendingTo...)
		}
		appendingTo = totalAppendingTo

	// OpQuest, i.e. a?, is much like an alternation between exactly two options, one of which is the empty string.
	case syntax.OpQuest:
		child := rx.Sub[0]
		var childAppendingTo []*pathCollector
		for _, collector := range appendingTo {
			newCollector := new(pathCollector)
			newCollector.WriteString(collector.String())
			newCollector.conditionalSlashAppendedAtLength = collector.conditionalSlashAppendedAtLength
			childAppendingTo = append(childAppendingTo, newCollector)
		}
		childAppendingTo, err = collectPathsFromRegexpASTInternal(child, childAppendingTo)
		if err != nil {
			return nil, err
		}
		appendingTo = append(appendingTo, childAppendingTo...)

		// Many Vault path patterns end with `/?` to accept paths that end with or without a slash. Our current
		// convention for generating the OpenAPI is to strip away these slashes. To do that, this very special case
		// detects when we just appended a single conditional slash, and records the length of the path at this point,
		// so we can later discard this path variant, if nothing else is appended to it later.
		if child.Op == syntax.OpLiteral && string(child.Rune) == "/" {
			for _, collector := range childAppendingTo {
				collector.conditionalSlashAppendedAtLength = collector.Len()
			}
		}

	// OpCapture, i.e. ( ) or (?P<name> ), a capturing group
	case syntax.OpCapture:
		if rx.Name == "" {
			// In Vault, an unnamed capturing group is not actually used for capturing.
			// We treat it exactly the same as OpConcat.
			for _, child := range rx.Sub {
				appendingTo, err = collectPathsFromRegexpASTInternal(child, appendingTo)
				if err != nil {
					return nil, err
				}
			}
		} else {
			// A named capturing group is replaced with the OpenAPI parameter syntax, and the regexp inside the group
			// is NOT added to the OpenAPI path.
			for _, builder := range appendingTo {
				builder.WriteRune('{')
				builder.WriteString(rx.Name)
				builder.WriteRune('}')
			}
		}

	// Any other kind of operation is a problem, and will trigger an error, resulting in the pattern being left out of
	// the OpenAPI entirely - that's better than generating a path which is incorrect.
	//
	// The Op types we expect to hit the default condition are:
	//
	//     OpCharClass    - i.e. [something]
	//     OpAnyCharNotNL - i.e. .
	//     OpAnyChar      - i.e. (?s:.)
	//     OpStar         - i.e. *
	//     OpPlus         - i.e. +
	//     OpRepeat       - i.e. {N}, {N,M}, etc.
	//
	// In any of these conditions, there is no sensible translation of the path to OpenAPI syntax. (Note, this only
	// applies to these appearing outside of a named capture group, otherwise they are handled in the previous case.)
	//
	// At the time of writing, the only pattern in the builtin Vault plugins that hits this codepath is the ".*"
	// pattern in the KVv2 secrets engine, which is not a valid path, but rather, is a catch-all used to implement
	// custom error handling behaviour to guide users who attempt to treat a KVv2 as a KVv1. It is already marked as
	// Unpublished, so is withheld from the OpenAPI anyway.
	//
	// For completeness, one other Op type exists, OpNoMatch, which is never generated by syntax.Parse - only by
	// subsequent Simplify in preparation to Compile, which is not used here.
	default:
		return nil, errUnsupportableRegexpOperationForOpenAPI
	}

	return appendingTo, nil
}

// schemaType is a subset of the JSON Schema elements used as a target
// for conversions from Vault's standard FieldTypes.
type schemaType struct {
	baseType string
	items    string
	format   string
	pattern  string
}

// convertType translates a FieldType into an OpenAPI type.
// In the case of arrays, a subtype is returned as well.
func convertType(t FieldType) schemaType {
	ret := schemaType{}

	switch t {
	case TypeString, TypeHeader:
		ret.baseType = "string"
	case TypeNameString:
		ret.baseType = "string"
		ret.pattern = `\w([\w-.]*\w)?`
	case TypeLowerCaseString:
		ret.baseType = "string"
		ret.format = "lowercase"
	case TypeInt:
		ret.baseType = "integer"
	case TypeInt64:
		ret.baseType = "integer"
		ret.format = "int64"
	case TypeDurationSecond, TypeSignedDurationSecond:
		ret.baseType = "integer"
		ret.format = "seconds"
	case TypeBool:
		ret.baseType = "boolean"
	case TypeMap:
		ret.baseType = "object"
		ret.format = "map"
	case TypeKVPairs:
		ret.baseType = "object"
		ret.format = "kvpairs"
	case TypeSlice:
		ret.baseType = "array"
		ret.items = "object"
	case TypeStringSlice, TypeCommaStringSlice:
		ret.baseType = "array"
		ret.items = "string"
	case TypeCommaIntSlice:
		ret.baseType = "array"
		ret.items = "integer"
	case TypeTime:
		ret.baseType = "string"
		ret.format = "date-time"
	case TypeFloat:
		ret.baseType = "number"
		ret.format = "float"
	default:
		log.L().Warn("error parsing field type", "type", t)
		ret.format = "unknown"
	}

	return ret
}

// cleanString prepares s for inclusion in the output
func cleanString(s string) string {
	// clean leading/trailing whitespace, and replace whitespace runs into a single space
	s = strings.TrimSpace(s)
	s = wsRe.ReplaceAllString(s, " ")
	return s
}

// splitFields partitions fields into path and body groups
// The input pattern is expected to have been run through expandPattern,
// with paths parameters denotes in {braces}.
func splitFields(allFields map[string]*FieldSchema, pattern string) (pathFields, bodyFields map[string]*FieldSchema) {
	pathFields = make(map[string]*FieldSchema)
	bodyFields = make(map[string]*FieldSchema)

	for _, match := range pathFieldsRe.FindAllStringSubmatch(pattern, -1) {
		name := match[1]
		pathFields[name] = allFields[name]
	}

	for name, field := range allFields {
		if _, ok := pathFields[name]; !ok {
			if field.Query {
				pathFields[name] = field
			} else {
				bodyFields[name] = field
			}
		}
	}

	return pathFields, bodyFields
}

// cleanedResponse is identical to logical.Response but with nulls
// removed from from JSON encoding
type cleanedResponse struct {
	Secret   *logical.Secret            `json:"secret,omitempty"`
	Auth     *logical.Auth              `json:"auth,omitempty"`
	Data     map[string]interface{}     `json:"data,omitempty"`
	Redirect string                     `json:"redirect,omitempty"`
	Warnings []string                   `json:"warnings,omitempty"`
	WrapInfo *wrapping.ResponseWrapInfo `json:"wrap_info,omitempty"`
	Headers  map[string][]string        `json:"headers,omitempty"`
}

func cleanResponse(resp *logical.Response) *cleanedResponse {
	return &cleanedResponse{
		Secret:   resp.Secret,
		Auth:     resp.Auth,
		Data:     resp.Data,
		Redirect: resp.Redirect,
		Warnings: resp.Warnings,
		WrapInfo: resp.WrapInfo,
		Headers:  resp.Headers,
	}
}

// CreateOperationIDs generates unique operationIds for all paths/methods.
// The transform will convert path/method into camelcase. e.g.:
//
// /sys/tools/random/{urlbytes} -> postSysToolsRandomUrlbytes
//
// In the unlikely case of a duplicate ids, a numeric suffix is added:
//
//	postSysToolsRandomUrlbytes_2
//
// An optional user-provided suffix ("context") may also be appended.
//
// Deprecated: operationID's are now populated using `constructOperationID`.
// This function is here for backwards compatibility with older plugins.
func (d *OASDocument) CreateOperationIDs(context string) {
	opIDCount := make(map[string]int)
	var paths []string

	// traverse paths in a stable order to ensure stable output
	for path := range d.Paths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		pi := d.Paths[path]
		for _, method := range []string{"get", "post", "delete"} {
			var oasOperation *OASOperation
			switch method {
			case "get":
				oasOperation = pi.Get
			case "post":
				oasOperation = pi.Post
			case "delete":
				oasOperation = pi.Delete
			}

			if oasOperation == nil {
				continue
			}

			if oasOperation.OperationID != "" {
				continue
			}

			// Discard "_mount_path" from any {thing_mount_path} parameters
			path = strings.Replace(path, "_mount_path", "", 1)

			// Space-split on non-words, title case everything, recombine
			opID := nonWordRe.ReplaceAllString(strings.ToLower(path), " ")
			opID = strings.Title(opID)
			opID = method + strings.ReplaceAll(opID, " ", "")

			// deduplicate operationIds. This is a safeguard, since generated IDs should
			// already be unique given our current path naming conventions.
			opIDCount[opID]++
			if opIDCount[opID] > 1 {
				opID = fmt.Sprintf("%s_%d", opID, opIDCount[opID])
			}

			if context != "" {
				opID += "_" + context
			}

			oasOperation.OperationID = opID
		}
	}
}
