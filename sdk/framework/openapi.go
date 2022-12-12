package framework

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// Regex for handling optional and named parameters in paths, and string cleanup.
// Predefined here to avoid substantial recompilation.

// Capture optional path elements in ungreedy (?U) fashion
// Both "(leases/)?renew" and "(/(?P<name>.+))?" formats are detected
var optRe = regexp.MustCompile(`(?U)\([^(]*\)\?|\(/\(\?P<[^(]*\)\)\?`)

var (
	altFieldsGroupRe = regexp.MustCompile(`\(\?P<\w+>\w+(\|\w+)+\)`)              // Match named groups that limit options, e.g. "(?<foo>a|b|c)"
	altFieldsRe      = regexp.MustCompile(`\w+(\|\w+)+`)                          // Match an options set, e.g. "a|b|c"
	altRe            = regexp.MustCompile(`\((.*)\|(.*)\)`)                       // Capture alternation elements, e.g. "(raw/?$|raw/(?P<path>.+))"
	altRootsRe       = regexp.MustCompile(`^\(([\w\-_]+(?:\|[\w\-_]+)+)\)(/.*)$`) // Pattern starting with alts, e.g. "(root1|root2)/(?P<name>regex)"
	cleanCharsRe     = regexp.MustCompile("[()^$?]")                              // Set of regex characters that will be stripped during cleaning
	cleanSuffixRe    = regexp.MustCompile(`/\?\$?$`)                              // Path suffix patterns that will be stripped during cleaning
	nonWordRe        = regexp.MustCompile(`[^a-zA-Z0-9]+`)                        // Match a sequence of non-word characters
	pathFieldsRe     = regexp.MustCompile(`{(\w+)}`)                              // Capture OpenAPI-style named parameters, e.g. "lookup/{urltoken}",
	reqdRe           = regexp.MustCompile(`\(?\?P<(\w+)>[^)]*\)?`)                // Capture required parameters, e.g. "(?P<name>regex)"
	wsRe             = regexp.MustCompile(`\s+`)                                  // Match whitespace, to be compressed during cleaning
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
	paths := expandPattern(p.Pattern)

	for _, path := range paths {

		log.L().Warn(
			fmt.Sprintf(
				`"{prefix: "%s", path: "%s"}: "%s",`,
				requestResponsePrefix,
				path,
				constructRequestIdentifier(logical.Operation(""), path, requestResponsePrefix, "")))

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

		defaultMountPath := requestResponsePrefix
		if requestResponsePrefix == "kv" {
			defaultMountPath = "secret"
		}

		if defaultMountPath != "system" && defaultMountPath != "identity" {
			p := OASParameter{
				Name:        fmt.Sprintf("%s_mount_path", defaultMountPath),
				Description: "Path where the backend was mounted; the endpoint path will be offset by the mount path",
				In:          "path",
				Schema: &OASSchema{
					Type:    "string",
					Default: defaultMountPath,
				},
				Required: false,
			}

			pi.Parameters = append(pi.Parameters, p)
		}

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
			if props.Unpublished {
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

			op.Summary = props.Summary
			op.Description = props.Description
			op.Deprecated = props.Deprecated

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
					requestName := constructRequestIdentifier(opType, path, requestResponsePrefix, "request")
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
						responseName := constructRequestIdentifier(opType, path, requestResponsePrefix, "response")
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

// When constructing a request or response name from a prefix + path, certain
// paths result in very long names or names with duplicate parts. For such paths,
// we use custom path mappings instead.
type knownPathKey struct {
	prefix string
	path   string
}

var knownPathMappings = map[knownPathKey]string{
	{prefix: "ad", path: "config"}:                                                      "AdConfig",
	{prefix: "ad", path: "creds/{name}"}:                                                "AdCredsName",
	{prefix: "ad", path: "library"}:                                                     "AdLibrary",
	{prefix: "ad", path: "library/manage/{name}/check-in"}:                              "AdLibraryManageNameCheckIn",
	{prefix: "ad", path: "library/{name}"}:                                              "AdLibraryName",
	{prefix: "ad", path: "library/{name}/check-in"}:                                     "AdLibraryNameCheckIn",
	{prefix: "ad", path: "library/{name}/check-out"}:                                    "AdLibraryNameCheckOut",
	{prefix: "ad", path: "library/{name}/status"}:                                       "AdLibraryNameStatus",
	{prefix: "ad", path: "roles"}:                                                       "AdRoles",
	{prefix: "ad", path: "roles/{name}"}:                                                "AdRolesName",
	{prefix: "ad", path: "rotate-role/{name}"}:                                          "AdRotateRoleName",
	{prefix: "ad", path: "rotate-root"}:                                                 "AdRotateRoot",
	{prefix: "alicloud", path: "config"}:                                                "AlicloudConfig",
	{prefix: "alicloud", path: "creds/{name}"}:                                          "AlicloudCredsName",
	{prefix: "alicloud", path: "role"}:                                                  "AlicloudRole",
	{prefix: "alicloud", path: "role/{name}"}:                                           "AlicloudRoleName",
	{prefix: "auth/alicloud", path: "login"}:                                            "AuthAlicloudLogin",
	{prefix: "auth/alicloud", path: "role"}:                                             "AuthAlicloudRole",
	{prefix: "auth/alicloud", path: "role/{role}"}:                                      "AuthAlicloudRoleRole",
	{prefix: "auth/alicloud", path: "roles"}:                                            "AuthAlicloudRoles",
	{prefix: "auth/approle", path: "login"}:                                             "AuthApproleLogin",
	{prefix: "auth/approle", path: "role"}:                                              "AuthApproleRole",
	{prefix: "auth/approle", path: "role/{role_name}"}:                                  "AuthApproleRoleRoleName",
	{prefix: "auth/approle", path: "role/{role_name}/bind-secret-id"}:                   "AuthApproleRoleRoleNameBindSecretId",
	{prefix: "auth/approle", path: "role/{role_name}/bound-cidr-list"}:                  "AuthApproleRoleRoleNameBoundCidrList",
	{prefix: "auth/approle", path: "role/{role_name}/custom-secret-id"}:                 "AuthApproleRoleRoleNameCustomSecretId",
	{prefix: "auth/approle", path: "role/{role_name}/local-secret-ids"}:                 "AuthApproleRoleRoleNameLocalSecretIds",
	{prefix: "auth/approle", path: "role/{role_name}/period"}:                           "AuthApproleRoleRoleNamePeriod",
	{prefix: "auth/approle", path: "role/{role_name}/policies"}:                         "AuthApproleRoleRoleNamePolicies",
	{prefix: "auth/approle", path: "role/{role_name}/role-id"}:                          "AuthApproleRoleRoleNameRoleId",
	{prefix: "auth/approle", path: "role/{role_name}/secret-id"}:                        "AuthApproleRoleRoleNameSecretId",
	{prefix: "auth/approle", path: "role/{role_name}/secret-id-accessor/destroy"}:       "AuthApproleRoleRoleNameSecretIdAccessorDestroy",
	{prefix: "auth/approle", path: "role/{role_name}/secret-id-accessor/lookup"}:        "AuthApproleRoleRoleNameSecretIdAccessorLookup",
	{prefix: "auth/approle", path: "role/{role_name}/secret-id-bound-cidrs"}:            "AuthApproleRoleRoleNameSecretIdBoundCidrs",
	{prefix: "auth/approle", path: "role/{role_name}/secret-id-num-uses"}:               "AuthApproleRoleRoleNameSecretIdNumUses",
	{prefix: "auth/approle", path: "role/{role_name}/secret-id-ttl"}:                    "AuthApproleRoleRoleNameSecretIdTtl",
	{prefix: "auth/approle", path: "role/{role_name}/secret-id/destroy"}:                "AuthApproleRoleRoleNameSecretIdDestroy",
	{prefix: "auth/approle", path: "role/{role_name}/secret-id/lookup"}:                 "AuthApproleRoleRoleNameSecretIdLookup",
	{prefix: "auth/approle", path: "role/{role_name}/token-bound-cidrs"}:                "AuthApproleRoleRoleNameTokenBoundCidrs",
	{prefix: "auth/approle", path: "role/{role_name}/token-max-ttl"}:                    "AuthApproleRoleRoleNameTokenMaxTtl",
	{prefix: "auth/approle", path: "role/{role_name}/token-num-uses"}:                   "AuthApproleRoleRoleNameTokenNumUses",
	{prefix: "auth/approle", path: "role/{role_name}/token-ttl"}:                        "AuthApproleRoleRoleNameTokenTtl",
	{prefix: "auth/approle", path: "tidy/secret-id"}:                                    "AuthApproleTidySecretId",
	{prefix: "auth/aws", path: "config/certificate/{cert_name}"}:                        "AuthAwsConfigCertificateCertName",
	{prefix: "auth/aws", path: "config/certificates"}:                                   "AuthAwsConfigCertificates",
	{prefix: "auth/aws", path: "config/client"}:                                         "AuthAwsConfigClient",
	{prefix: "auth/aws", path: "config/identity"}:                                       "AuthAwsConfigIdentity",
	{prefix: "auth/aws", path: "config/rotate-root"}:                                    "AuthAwsConfigRotateRoot",
	{prefix: "auth/aws", path: "config/sts"}:                                            "AuthAwsConfigSts",
	{prefix: "auth/aws", path: "config/sts/{account_id}"}:                               "AuthAwsConfigStsAccountId",
	{prefix: "auth/aws", path: "config/tidy/identity-accesslist"}:                       "AuthAwsConfigTidyIdentityAccesslist",
	{prefix: "auth/aws", path: "config/tidy/identity-whitelist"}:                        "AuthAwsConfigTidyIdentityWhitelist",
	{prefix: "auth/aws", path: "config/tidy/roletag-blacklist"}:                         "AuthAwsConfigTidyRoletagBlacklist",
	{prefix: "auth/aws", path: "config/tidy/roletag-denylist"}:                          "AuthAwsConfigTidyRoletagDenylist",
	{prefix: "auth/aws", path: "identity-accesslist"}:                                   "AuthAwsIdentityAccesslist",
	{prefix: "auth/aws", path: "identity-accesslist/{instance_id}"}:                     "AuthAwsIdentityAccesslistInstanceId",
	{prefix: "auth/aws", path: "identity-whitelist"}:                                    "AuthAwsIdentityWhitelist",
	{prefix: "auth/aws", path: "identity-whitelist/{instance_id}"}:                      "AuthAwsIdentityWhitelistInstanceId",
	{prefix: "auth/aws", path: "login"}:                                                 "AuthAwsLogin",
	{prefix: "auth/aws", path: "role"}:                                                  "AuthAwsRole",
	{prefix: "auth/aws", path: "role/{role}"}:                                           "AuthAwsRoleRole",
	{prefix: "auth/aws", path: "role/{role}/tag"}:                                       "AuthAwsRoleRoleTag",
	{prefix: "auth/aws", path: "roles"}:                                                 "AuthAwsRoles",
	{prefix: "auth/aws", path: "roletag-blacklist"}:                                     "AuthAwsRoletagBlacklist",
	{prefix: "auth/aws", path: "roletag-blacklist/{role_tag}"}:                          "AuthAwsRoletagBlacklistRoleTag",
	{prefix: "auth/aws", path: "roletag-denylist"}:                                      "AuthAwsRoletagDenylist",
	{prefix: "auth/aws", path: "roletag-denylist/{role_tag}"}:                           "AuthAwsRoletagDenylistRoleTag",
	{prefix: "auth/aws", path: "tidy/identity-accesslist"}:                              "AuthAwsTidyIdentityAccesslist",
	{prefix: "auth/aws", path: "tidy/identity-whitelist"}:                               "AuthAwsTidyIdentityWhitelist",
	{prefix: "auth/aws", path: "tidy/roletag-blacklist"}:                                "AuthAwsTidyRoletagBlacklist",
	{prefix: "auth/aws", path: "tidy/roletag-denylist"}:                                 "AuthAwsTidyRoletagDenylist",
	{prefix: "auth/azure", path: "config"}:                                              "AuthAzureConfig",
	{prefix: "auth/azure", path: "login"}:                                               "AuthAzureLogin",
	{prefix: "auth/azure", path: "role"}:                                                "AuthAzureRole",
	{prefix: "auth/azure", path: "role/{name}"}:                                         "AuthAzureRoleName",
	{prefix: "auth/centrify", path: "config"}:                                           "AuthCentrifyConfig",
	{prefix: "auth/centrify", path: "login"}:                                            "AuthCentrifyLogin",
	{prefix: "auth/cert", path: "certs"}:                                                "AuthCertCerts",
	{prefix: "auth/cert", path: "certs/{name}"}:                                         "AuthCertCertsName",
	{prefix: "auth/cert", path: "config"}:                                               "AuthCertConfig",
	{prefix: "auth/cert", path: "crls"}:                                                 "AuthCertCrls",
	{prefix: "auth/cert", path: "crls/{name}"}:                                          "AuthCertCrlsName",
	{prefix: "auth/cert", path: "login"}:                                                "AuthCertLogin",
	{prefix: "auth/cf", path: "config"}:                                                 "AuthCfConfig",
	{prefix: "auth/cf", path: "login"}:                                                  "AuthCfLogin",
	{prefix: "auth/cf", path: "roles"}:                                                  "AuthCfRoles",
	{prefix: "auth/cf", path: "roles/{role}"}:                                           "AuthCfRolesRole",
	{prefix: "auth/gcp", path: "config"}:                                                "AuthGcpConfig",
	{prefix: "auth/gcp", path: "login"}:                                                 "AuthGcpLogin",
	{prefix: "auth/gcp", path: "role"}:                                                  "AuthGcpRole",
	{prefix: "auth/gcp", path: "role/{name}"}:                                           "AuthGcpRoleName",
	{prefix: "auth/gcp", path: "role/{name}/labels"}:                                    "AuthGcpRoleNameLabels",
	{prefix: "auth/gcp", path: "role/{name}/service-accounts"}:                          "AuthGcpRoleNameServiceAccounts",
	{prefix: "auth/gcp", path: "roles"}:                                                 "AuthGcpRoles",
	{prefix: "auth/github", path: "config"}:                                             "AuthGithubConfig",
	{prefix: "auth/github", path: "login"}:                                              "AuthGithubLogin",
	{prefix: "auth/github", path: "map/teams"}:                                          "AuthGithubMapTeams",
	{prefix: "auth/github", path: "map/teams/{key}"}:                                    "AuthGithubMapTeamsKey",
	{prefix: "auth/github", path: "map/users"}:                                          "AuthGithubMapUsers",
	{prefix: "auth/github", path: "map/users/{key}"}:                                    "AuthGithubMapUsersKey",
	{prefix: "auth/jwt", path: "config"}:                                                "AuthJwtConfig",
	{prefix: "auth/jwt", path: "login"}:                                                 "AuthJwtLogin",
	{prefix: "auth/jwt", path: "oidc/auth_url"}:                                         "AuthJwtOidcAuthUrl",
	{prefix: "auth/jwt", path: "oidc/callback"}:                                         "AuthJwtOidcCallback",
	{prefix: "auth/jwt", path: "role"}:                                                  "AuthJwtRole",
	{prefix: "auth/jwt", path: "role/{name}"}:                                           "AuthJwtRoleName",
	{prefix: "auth/kerberos", path: "config"}:                                           "AuthKerberosConfig",
	{prefix: "auth/kerberos", path: "config/ldap"}:                                      "AuthKerberosConfigLdap",
	{prefix: "auth/kerberos", path: "groups"}:                                           "AuthKerberosGroups",
	{prefix: "auth/kerberos", path: "groups/{name}"}:                                    "AuthKerberosGroupsName",
	{prefix: "auth/kerberos", path: "login"}:                                            "AuthKerberosLogin",
	{prefix: "auth/kubernetes", path: "config"}:                                         "AuthKubernetesConfig",
	{prefix: "auth/kubernetes", path: "login"}:                                          "AuthKubernetesLogin",
	{prefix: "auth/kubernetes", path: "role"}:                                           "AuthKubernetesRole",
	{prefix: "auth/kubernetes", path: "role/{name}"}:                                    "AuthKubernetesRoleName",
	{prefix: "auth/ldap", path: "config"}:                                               "AuthLdapConfig",
	{prefix: "auth/ldap", path: "groups"}:                                               "AuthLdapGroups",
	{prefix: "auth/ldap", path: "groups/{name}"}:                                        "AuthLdapGroupsName",
	{prefix: "auth/ldap", path: "login/{username}"}:                                     "AuthLdapLoginUsername",
	{prefix: "auth/ldap", path: "users"}:                                                "AuthLdapUsers",
	{prefix: "auth/ldap", path: "users/{name}"}:                                         "AuthLdapUsersName",
	{prefix: "auth/oci", path: "config"}:                                                "AuthOciConfig",
	{prefix: "auth/oci", path: "login"}:                                                 "AuthOciLogin",
	{prefix: "auth/oci", path: "login/{role}"}:                                          "AuthOciLoginRole",
	{prefix: "auth/oci", path: "role"}:                                                  "AuthOciRole",
	{prefix: "auth/oci", path: "role/{role}"}:                                           "AuthOciRoleRole",
	{prefix: "auth/oidc", path: "config"}:                                               "AuthOidcConfig",
	{prefix: "auth/oidc", path: "login"}:                                                "AuthOidcLogin",
	{prefix: "auth/oidc", path: "oidc/auth_url"}:                                        "AuthOidcOidcAuthUrl",
	{prefix: "auth/oidc", path: "oidc/callback"}:                                        "AuthOidcOidcCallback",
	{prefix: "auth/oidc", path: "role"}:                                                 "AuthOidcRole",
	{prefix: "auth/oidc", path: "role/{name}"}:                                          "AuthOidcRoleName",
	{prefix: "auth/okta", path: "config"}:                                               "AuthOktaConfig",
	{prefix: "auth/okta", path: "groups"}:                                               "AuthOktaGroups",
	{prefix: "auth/okta", path: "groups/{name}"}:                                        "AuthOktaGroupsName",
	{prefix: "auth/okta", path: "login/{username}"}:                                     "AuthOktaLoginUsername",
	{prefix: "auth/okta", path: "users"}:                                                "AuthOktaUsers",
	{prefix: "auth/okta", path: "users/{name}"}:                                         "AuthOktaUsersName",
	{prefix: "auth/okta", path: "verify/{nonce}"}:                                       "AuthOktaVerifyNonce",
	{prefix: "auth/radius", path: "config"}:                                             "AuthRadiusConfig",
	{prefix: "auth/radius", path: "login"}:                                              "AuthRadiusLogin",
	{prefix: "auth/radius", path: "login/{urlusername}"}:                                "AuthRadiusLoginUrlusername",
	{prefix: "auth/radius", path: "users"}:                                              "AuthRadiusUsers",
	{prefix: "auth/radius", path: "users/{name}"}:                                       "AuthRadiusUsersName",
	{prefix: "auth/token", path: "accessors/"}:                                          "AuthTokenAccessors",
	{prefix: "auth/token", path: "create"}:                                              "AuthTokenCreate",
	{prefix: "auth/token", path: "create-orphan"}:                                       "AuthTokenCreateOrphan",
	{prefix: "auth/token", path: "create/{role_name}"}:                                  "AuthTokenCreateRoleName",
	{prefix: "auth/token", path: "lookup"}:                                              "AuthTokenLookup",
	{prefix: "auth/token", path: "lookup-accessor"}:                                     "AuthTokenLookupAccessor",
	{prefix: "auth/token", path: "lookup-self"}:                                         "AuthTokenLookupSelf",
	{prefix: "auth/token", path: "renew"}:                                               "AuthTokenRenew",
	{prefix: "auth/token", path: "renew-accessor"}:                                      "AuthTokenRenewAccessor",
	{prefix: "auth/token", path: "renew-self"}:                                          "AuthTokenRenewSelf",
	{prefix: "auth/token", path: "revoke"}:                                              "AuthTokenRevoke",
	{prefix: "auth/token", path: "revoke-accessor"}:                                     "AuthTokenRevokeAccessor",
	{prefix: "auth/token", path: "revoke-orphan"}:                                       "AuthTokenRevokeOrphan",
	{prefix: "auth/token", path: "revoke-self"}:                                         "AuthTokenRevokeSelf",
	{prefix: "auth/token", path: "roles"}:                                               "AuthTokenRoles",
	{prefix: "auth/token", path: "roles/{role_name}"}:                                   "AuthTokenRolesRoleName",
	{prefix: "auth/token", path: "tidy"}:                                                "AuthTokenTidy",
	{prefix: "auth/userpass", path: "login/{username}"}:                                 "AuthUserpassLoginUsername",
	{prefix: "auth/userpass", path: "users"}:                                            "AuthUserpassUsers",
	{prefix: "auth/userpass", path: "users/{username}"}:                                 "AuthUserpassUsersUsername",
	{prefix: "auth/userpass", path: "users/{username}/password"}:                        "AuthUserpassUsersUsernamePassword",
	{prefix: "auth/userpass", path: "users/{username}/policies"}:                        "AuthUserpassUsersUsernamePolicies",
	{prefix: "aws", path: "config/lease"}:                                               "AwsConfigLease",
	{prefix: "aws", path: "config/root"}:                                                "AwsConfigRoot",
	{prefix: "aws", path: "config/rotate-root"}:                                         "AwsConfigRotateRoot",
	{prefix: "aws", path: "creds"}:                                                      "AwsCreds",
	{prefix: "aws", path: "roles"}:                                                      "AwsRoles",
	{prefix: "aws", path: "roles/{name}"}:                                               "AwsRolesName",
	{prefix: "aws", path: "sts/{name}"}:                                                 "AwsStsName",
	{prefix: "azure", path: "config"}:                                                   "AzureConfig",
	{prefix: "azure", path: "creds/{role}"}:                                             "AzureCredsRole",
	{prefix: "azure", path: "roles"}:                                                    "AzureRoles",
	{prefix: "azure", path: "roles/{name}"}:                                             "AzureRolesName",
	{prefix: "azure", path: "rotate-root"}:                                              "AzureRotateRoot",
	{prefix: "consul", path: "config/access"}:                                           "ConsulConfigAccess",
	{prefix: "consul", path: "creds/{role}"}:                                            "ConsulCredsRole",
	{prefix: "consul", path: "roles"}:                                                   "ConsulRoles",
	{prefix: "consul", path: "roles/{name}"}:                                            "ConsulRolesName",
	{prefix: "cubbyhole", path: "{path}"}:                                               "CubbyholePath",
	{prefix: "gcp", path: "config"}:                                                     "GcpConfig",
	{prefix: "gcp", path: "config/rotate-root"}:                                         "GcpConfigRotateRoot",
	{prefix: "gcp", path: "key/{roleset}"}:                                              "GcpKeyRoleset",
	{prefix: "gcp", path: "roleset/{name}"}:                                             "GcpRolesetName",
	{prefix: "gcp", path: "roleset/{name}/rotate"}:                                      "GcpRolesetNameRotate",
	{prefix: "gcp", path: "roleset/{name}/rotate-key"}:                                  "GcpRolesetNameRotateKey",
	{prefix: "gcp", path: "roleset/{roleset}/key"}:                                      "GcpRolesetRolesetKey",
	{prefix: "gcp", path: "roleset/{roleset}/token"}:                                    "GcpRolesetRolesetToken",
	{prefix: "gcp", path: "rolesets"}:                                                   "GcpRolesets",
	{prefix: "gcp", path: "static-account/{name}"}:                                      "GcpStaticAccountName",
	{prefix: "gcp", path: "static-account/{name}/key"}:                                  "GcpStaticAccountNameKey",
	{prefix: "gcp", path: "static-account/{name}/rotate-key"}:                           "GcpStaticAccountNameRotateKey",
	{prefix: "gcp", path: "static-account/{name}/token"}:                                "GcpStaticAccountNameToken",
	{prefix: "gcp", path: "static-accounts"}:                                            "GcpStaticAccounts",
	{prefix: "gcp", path: "token/{roleset}"}:                                            "GcpTokenRoleset",
	{prefix: "gcpkms", path: "config"}:                                                  "GcpkmsConfig",
	{prefix: "gcpkms", path: "decrypt/{key}"}:                                           "GcpkmsDecryptKey",
	{prefix: "gcpkms", path: "encrypt/{key}"}:                                           "GcpkmsEncryptKey",
	{prefix: "gcpkms", path: "keys"}:                                                    "GcpkmsKeys",
	{prefix: "gcpkms", path: "keys/config/{key}"}:                                       "GcpkmsKeysConfigKey",
	{prefix: "gcpkms", path: "keys/deregister/{key}"}:                                   "GcpkmsKeysDeregisterKey",
	{prefix: "gcpkms", path: "keys/register/{key}"}:                                     "GcpkmsKeysRegisterKey",
	{prefix: "gcpkms", path: "keys/rotate/{key}"}:                                       "GcpkmsKeysRotateKey",
	{prefix: "gcpkms", path: "keys/trim/{key}"}:                                         "GcpkmsKeysTrimKey",
	{prefix: "gcpkms", path: "keys/{key}"}:                                              "GcpkmsKeysKey",
	{prefix: "gcpkms", path: "pubkey/{key}"}:                                            "GcpkmsPubkeyKey",
	{prefix: "gcpkms", path: "reencrypt/{key}"}:                                         "GcpkmsReencryptKey",
	{prefix: "gcpkms", path: "sign/{key}"}:                                              "GcpkmsSignKey",
	{prefix: "gcpkms", path: "verify/{key}"}:                                            "GcpkmsVerifyKey",
	{prefix: "identity", path: "alias"}:                                                 "IdentityAlias",
	{prefix: "identity", path: "alias/id"}:                                              "IdentityAliasId",
	{prefix: "identity", path: "alias/id/{id}"}:                                         "IdentityAliasIdId",
	{prefix: "identity", path: "entity"}:                                                "IdentityEntity",
	{prefix: "identity", path: "entity-alias"}:                                          "IdentityEntityAlias",
	{prefix: "identity", path: "entity-alias/id"}:                                       "IdentityEntityAliasId",
	{prefix: "identity", path: "entity-alias/id/{id}"}:                                  "IdentityEntityAliasIdId",
	{prefix: "identity", path: "entity/batch-delete"}:                                   "IdentityEntityBatchDelete",
	{prefix: "identity", path: "entity/id"}:                                             "IdentityEntityId",
	{prefix: "identity", path: "entity/id/{id}"}:                                        "IdentityEntityIdId",
	{prefix: "identity", path: "entity/merge"}:                                          "IdentityEntityMerge",
	{prefix: "identity", path: "entity/name"}:                                           "IdentityEntityName",
	{prefix: "identity", path: "entity/name/{name}"}:                                    "IdentityEntityNameName",
	{prefix: "identity", path: "group"}:                                                 "IdentityGroup",
	{prefix: "identity", path: "group-alias"}:                                           "IdentityGroupAlias",
	{prefix: "identity", path: "group-alias/id"}:                                        "IdentityGroupAliasId",
	{prefix: "identity", path: "group-alias/id/{id}"}:                                   "IdentityGroupAliasIdId",
	{prefix: "identity", path: "group/id"}:                                              "IdentityGroupId",
	{prefix: "identity", path: "group/id/{id}"}:                                         "IdentityGroupIdId",
	{prefix: "identity", path: "group/name"}:                                            "IdentityGroupName",
	{prefix: "identity", path: "group/name/{name}"}:                                     "IdentityGroupNameName",
	{prefix: "identity", path: "lookup/entity"}:                                         "IdentityLookupEntity",
	{prefix: "identity", path: "lookup/group"}:                                          "IdentityLookupGroup",
	{prefix: "identity", path: "mfa/login-enforcement"}:                                 "IdentityMfaLoginEnforcement",
	{prefix: "identity", path: "mfa/login-enforcement/{name}"}:                          "IdentityMfaLoginEnforcementName",
	{prefix: "identity", path: "mfa/method"}:                                            "IdentityMfaMethod",
	{prefix: "identity", path: "mfa/method/duo"}:                                        "IdentityMfaMethodDuo",
	{prefix: "identity", path: "mfa/method/duo/{method_id}"}:                            "IdentityMfaMethodDuoMethodId",
	{prefix: "identity", path: "mfa/method/okta"}:                                       "IdentityMfaMethodOkta",
	{prefix: "identity", path: "mfa/method/okta/{method_id}"}:                           "IdentityMfaMethodOktaMethodId",
	{prefix: "identity", path: "mfa/method/pingid"}:                                     "IdentityMfaMethodPingid",
	{prefix: "identity", path: "mfa/method/pingid/{method_id}"}:                         "IdentityMfaMethodPingidMethodId",
	{prefix: "identity", path: "mfa/method/totp"}:                                       "IdentityMfaMethodTotp",
	{prefix: "identity", path: "mfa/method/totp/admin-destroy"}:                         "IdentityMfaMethodTotpAdminDestroy",
	{prefix: "identity", path: "mfa/method/totp/admin-generate"}:                        "IdentityMfaMethodTotpAdminGenerate",
	{prefix: "identity", path: "mfa/method/totp/generate"}:                              "IdentityMfaMethodTotpGenerate",
	{prefix: "identity", path: "mfa/method/totp/{method_id}"}:                           "IdentityMfaMethodTotpMethodId",
	{prefix: "identity", path: "mfa/method/{method_id}"}:                                "IdentityMfaMethodMethodId",
	{prefix: "identity", path: "oidc/.well-known/keys"}:                                 "IdentityOidcWellKnownKeys",
	{prefix: "identity", path: "oidc/.well-known/openid-configuration"}:                 "IdentityOidcWellKnownOpenidConfiguration",
	{prefix: "identity", path: "oidc/assignment"}:                                       "IdentityOidcAssignment",
	{prefix: "identity", path: "oidc/assignment/{name}"}:                                "IdentityOidcAssignmentName",
	{prefix: "identity", path: "oidc/client"}:                                           "IdentityOidcClient",
	{prefix: "identity", path: "oidc/client/{name}"}:                                    "IdentityOidcClientName",
	{prefix: "identity", path: "oidc/config"}:                                           "IdentityOidcConfig",
	{prefix: "identity", path: "oidc/introspect"}:                                       "IdentityOidcIntrospect",
	{prefix: "identity", path: "oidc/key"}:                                              "IdentityOidcKey",
	{prefix: "identity", path: "oidc/key/{name}"}:                                       "IdentityOidcKeyName",
	{prefix: "identity", path: "oidc/key/{name}/rotate"}:                                "IdentityOidcKeyNameRotate",
	{prefix: "identity", path: "oidc/provider"}:                                         "IdentityOidcProvider",
	{prefix: "identity", path: "oidc/provider/{name}"}:                                  "IdentityOidcProviderName",
	{prefix: "identity", path: "oidc/provider/{name}/.well-known/keys"}:                 "IdentityOidcProviderNameWellKnownKeys",
	{prefix: "identity", path: "oidc/provider/{name}/.well-known/openid-configuration"}: "IdentityOidcProviderNameWellKnownOpenidConfiguration",
	{prefix: "identity", path: "oidc/provider/{name}/authorize"}:                        "IdentityOidcProviderNameAuthorize",
	{prefix: "identity", path: "oidc/provider/{name}/token"}:                            "IdentityOidcProviderNameToken",
	{prefix: "identity", path: "oidc/provider/{name}/userinfo"}:                         "IdentityOidcProviderNameUserinfo",
	{prefix: "identity", path: "oidc/role"}:                                             "IdentityOidcRole",
	{prefix: "identity", path: "oidc/role/{name}"}:                                      "IdentityOidcRoleName",
	{prefix: "identity", path: "oidc/scope"}:                                            "IdentityOidcScope",
	{prefix: "identity", path: "oidc/scope/{name}"}:                                     "IdentityOidcScopeName",
	{prefix: "identity", path: "oidc/token/{name}"}:                                     "IdentityOidcTokenName",
	{prefix: "identity", path: "persona"}:                                               "IdentityPersona",
	{prefix: "identity", path: "persona/id"}:                                            "IdentityPersonaId",
	{prefix: "identity", path: "persona/id/{id}"}:                                       "IdentityPersonaIdId",
	{prefix: "kubernetes", path: "config"}:                                              "KubernetesConfig",
	{prefix: "kubernetes", path: "creds/{name}"}:                                        "KubernetesCredsName",
	{prefix: "kubernetes", path: "roles"}:                                               "KubernetesRoles",
	{prefix: "kubernetes", path: "roles/{name}"}:                                        "KubernetesRolesName",
	{prefix: "kv", path: "{path}"}:                                                      "KvPath",
	{prefix: "ldap", path: "config"}:                                                    "LdapConfig",
	{prefix: "ldap", path: "creds/{name}"}:                                              "LdapCredsName",
	{prefix: "ldap", path: "library"}:                                                   "LdapLibrary",
	{prefix: "ldap", path: "library/manage/{name}/check-in"}:                            "LdapLibraryManageNameCheckIn",
	{prefix: "ldap", path: "library/{name}"}:                                            "LdapLibraryName",
	{prefix: "ldap", path: "library/{name}/check-in"}:                                   "LdapLibraryNameCheckIn",
	{prefix: "ldap", path: "library/{name}/check-out"}:                                  "LdapLibraryNameCheckOut",
	{prefix: "ldap", path: "library/{name}/status"}:                                     "LdapLibraryNameStatus",
	{prefix: "ldap", path: "role"}:                                                      "LdapRole",
	{prefix: "ldap", path: "role/{name}"}:                                               "LdapRoleName",
	{prefix: "ldap", path: "rotate-role/{name}"}:                                        "LdapRotateRoleName",
	{prefix: "ldap", path: "rotate-root"}:                                               "LdapRotateRoot",
	{prefix: "ldap", path: "static-cred/{name}"}:                                        "LdapStaticCredName",
	{prefix: "ldap", path: "static-role"}:                                               "LdapStaticRole",
	{prefix: "ldap", path: "static-role/{name}"}:                                        "LdapStaticRoleName",
	{prefix: "mongodbatlas", path: "config"}:                                            "MongodbatlasConfig",
	{prefix: "mongodbatlas", path: "creds/{name}"}:                                      "MongodbatlasCredsName",
	{prefix: "mongodbatlas", path: "roles"}:                                             "MongodbatlasRoles",
	{prefix: "mongodbatlas", path: "roles/{name}"}:                                      "MongodbatlasRolesName",
	{prefix: "nomad", path: "config/access"}:                                            "NomadConfigAccess",
	{prefix: "nomad", path: "config/lease"}:                                             "NomadConfigLease",
	{prefix: "nomad", path: "creds/{name}"}:                                             "NomadCredsName",
	{prefix: "nomad", path: "role"}:                                                     "NomadRole",
	{prefix: "nomad", path: "role/{name}"}:                                              "NomadRoleName",
	{prefix: "openldap", path: "config"}:                                                "OpenldapConfig",
	{prefix: "openldap", path: "creds/{name}"}:                                          "OpenldapCredsName",
	{prefix: "openldap", path: "library"}:                                               "OpenldapLibrary",
	{prefix: "openldap", path: "library/manage/{name}/check-in"}:                        "OpenldapLibraryManageNameCheckIn",
	{prefix: "openldap", path: "library/{name}"}:                                        "OpenldapLibraryName",
	{prefix: "openldap", path: "library/{name}/check-in"}:                               "OpenldapLibraryNameCheckIn",
	{prefix: "openldap", path: "library/{name}/check-out"}:                              "OpenldapLibraryNameCheckOut",
	{prefix: "openldap", path: "library/{name}/status"}:                                 "OpenldapLibraryNameStatus",
	{prefix: "openldap", path: "role"}:                                                  "OpenldapRole",
	{prefix: "openldap", path: "role/{name}"}:                                           "OpenldapRoleName",
	{prefix: "openldap", path: "rotate-role/{name}"}:                                    "OpenldapRotateRoleName",
	{prefix: "openldap", path: "rotate-root"}:                                           "OpenldapRotateRoot",
	{prefix: "openldap", path: "static-cred/{name}"}:                                    "OpenldapStaticCredName",
	{prefix: "openldap", path: "static-role"}:                                           "OpenldapStaticRole",
	{prefix: "openldap", path: "static-role/{name}"}:                                    "OpenldapStaticRoleName",
	{prefix: "pki", path: "/delta"}:                                                     "PkiDelta",
	{prefix: "pki", path: "/delta/pem"}:                                                 "PkiDeltaPem",
	{prefix: "pki", path: "/der"}:                                                       "PkiDer",
	{prefix: "pki", path: "/json"}:                                                      "PkiJson",
	{prefix: "pki", path: "/pem"}:                                                       "PkiPem",
	{prefix: "pki", path: "bundle"}:                                                     "PkiBundle",
	{prefix: "pki", path: "ca"}:                                                         "PkiCa",
	{prefix: "pki", path: "ca/pem"}:                                                     "PkiCaPem",
	{prefix: "pki", path: "ca_chain"}:                                                   "PkiCaChain",
	{prefix: "pki", path: "cert"}:                                                       "PkiCert",
	{prefix: "pki", path: "cert/ca_chain"}:                                              "PkiCertCaChain",
	{prefix: "pki", path: "cert/{serial}"}:                                              "PkiCertSerial",
	{prefix: "pki", path: "cert/{serial}/raw"}:                                          "PkiCertSerialRaw",
	{prefix: "pki", path: "cert/{serial}/raw/pem"}:                                      "PkiCertSerialRawPem",
	{prefix: "pki", path: "certs"}:                                                      "PkiCerts",
	{prefix: "pki", path: "certs/revoked"}:                                              "PkiCertsRevoked",
	{prefix: "pki", path: "config/auto-tidy"}:                                           "PkiConfigAutoTidy",
	{prefix: "pki", path: "config/ca"}:                                                  "PkiConfigCa",
	{prefix: "pki", path: "config/cluster"}:                                             "PkiConfigCluster",
	{prefix: "pki", path: "config/crl"}:                                                 "PkiConfigCrl",
	{prefix: "pki", path: "config/issuers"}:                                             "PkiConfigIssuers",
	{prefix: "pki", path: "config/keys"}:                                                "PkiConfigKeys",
	{prefix: "pki", path: "config/urls"}:                                                "PkiConfigUrls",
	{prefix: "pki", path: "crl"}:                                                        "PkiCrl",
	{prefix: "pki", path: "crl/rotate"}:                                                 "PkiCrlRotate",
	{prefix: "pki", path: "crl/rotate-delta"}:                                           "PkiCrlRotateDelta",
	{prefix: "pki", path: "delta-crl"}:                                                  "PkiDeltaCrl",
	{prefix: "pki", path: "intermediate/cross-sign"}:                                    "PkiIntermediateCrossSign",
	{prefix: "pki", path: "intermediate/generate/{exported}"}:                           "PkiIntermediateGenerateExported",
	{prefix: "pki", path: "intermediate/set-signed"}:                                    "PkiIntermediateSetSigned",
	{prefix: "pki", path: "internal|exported"}:                                          "PkiInternalExported",
	{prefix: "pki", path: "issue/{role}"}:                                               "PkiIssueRole",
	{prefix: "pki", path: "issuer/{issuer_ref}/issue/{role}"}:                           "PkiIssuerIssuerRefIssueRole",
	{prefix: "pki", path: "issuer/{issuer_ref}/resign-crls"}:                            "PkiIssuerIssuerRefResignCrls",
	{prefix: "pki", path: "issuer/{issuer_ref}/revoke"}:                                 "PkiIssuerIssuerRefRevoke",
	{prefix: "pki", path: "issuer/{issuer_ref}/sign-intermediate"}:                      "PkiIssuerIssuerRefSignIntermediate",
	{prefix: "pki", path: "issuer/{issuer_ref}/sign-revocation-list"}:                   "PkiIssuerIssuerRefSignRevocationList",
	{prefix: "pki", path: "issuer/{issuer_ref}/sign-self-issued"}:                       "PkiIssuerIssuerRefSignSelfIssued",
	{prefix: "pki", path: "issuer/{issuer_ref}/sign-verbatim"}:                          "PkiIssuerIssuerRefSignVerbatim",
	{prefix: "pki", path: "issuer/{issuer_ref}/sign-verbatim/{role}"}:                   "PkiIssuerIssuerRefSignVerbatimRole",
	{prefix: "pki", path: "issuer/{issuer_ref}/sign/{role}"}:                            "PkiIssuerIssuerRefSignRole",
	{prefix: "pki", path: "issuers"}:                                                    "PkiIssuers",
	{prefix: "pki", path: "issuers/generate/intermediate/{exported}"}:                   "PkiIssuersGenerateIntermediateExported",
	{prefix: "pki", path: "issuers/generate/root/{exported}"}:                           "PkiIssuersGenerateRootExported",
	{prefix: "pki", path: "key/{key_ref}"}:                                              "PkiKeyKeyRef",
	{prefix: "pki", path: "keys"}:                                                       "PkiKeys",
	{prefix: "pki", path: "keys/import"}:                                                "PkiKeysImport",
	{prefix: "pki", path: "kms"}:                                                        "PkiKms",
	{prefix: "pki", path: "ocsp"}:                                                       "PkiOcsp",
	{prefix: "pki", path: "ocsp/{req}"}:                                                 "PkiOcspReq",
	{prefix: "pki", path: "revoke"}:                                                     "PkiRevoke",
	{prefix: "pki", path: "revoke-with-key"}:                                            "PkiRevokeWithKey",
	{prefix: "pki", path: "roles"}:                                                      "PkiRoles",
	{prefix: "pki", path: "roles/{name}"}:                                               "PkiRolesName",
	{prefix: "pki", path: "root"}:                                                       "PkiRoot",
	{prefix: "pki", path: "root/generate/{exported}"}:                                   "PkiRootGenerateExported",
	{prefix: "pki", path: "root/replace"}:                                               "PkiRootReplace",
	{prefix: "pki", path: "root/rotate/{exported}"}:                                     "PkiRootRotateExported",
	{prefix: "pki", path: "root/sign-intermediate"}:                                     "PkiRootSignIntermediate",
	{prefix: "pki", path: "root/sign-self-issued"}:                                      "PkiRootSignSelfIssued",
	{prefix: "pki", path: "sign-verbatim"}:                                              "PkiSignVerbatim",
	{prefix: "pki", path: "sign-verbatim/{role}"}:                                       "PkiSignVerbatimRole",
	{prefix: "pki", path: "sign/{role}"}:                                                "PkiSignRole",
	{prefix: "pki", path: "tidy"}:                                                       "PkiTidy",
	{prefix: "pki", path: "tidy-cancel"}:                                                "PkiTidyCancel",
	{prefix: "pki", path: "tidy-status"}:                                                "PkiTidyStatus",
	{prefix: "pki", path: "{issuer_ref}/crl/pem|/der|/delta/pem"}:                       "PkiIssuerRefCrlPemDerDeltaPem",
	{prefix: "pki", path: "{issuer_ref}/der|/pem"}:                                      "PkiIssuerRefDerPem",
	{prefix: "rabbitmq", path: "config/connection"}:                                     "RabbitmqConfigConnection",
	{prefix: "rabbitmq", path: "config/lease"}:                                          "RabbitmqConfigLease",
	{prefix: "rabbitmq", path: "creds/{name}"}:                                          "RabbitmqCredsName",
	{prefix: "rabbitmq", path: "roles"}:                                                 "RabbitmqRoles",
	{prefix: "rabbitmq", path: "roles/{name}"}:                                          "RabbitmqRolesName",
	{prefix: "secret", path: ".*"}:                                                      "Secret",
	{prefix: "secret", path: "config"}:                                                  "SecretConfig",
	{prefix: "secret", path: "data/{path}"}:                                             "SecretDataPath",
	{prefix: "secret", path: "delete/{path}"}:                                           "SecretDeletePath",
	{prefix: "secret", path: "destroy/{path}"}:                                          "SecretDestroyPath",
	{prefix: "secret", path: "metadata/{path}"}:                                         "SecretMetadataPath",
	{prefix: "secret", path: "subkeys/{path}"}:                                          "SecretSubkeysPath",
	{prefix: "secret", path: "undelete/{path}"}:                                         "SecretUndeletePath",
	{prefix: "ssh", path: "config/ca"}:                                                  "SshConfigCa",
	{prefix: "ssh", path: "config/zeroaddress"}:                                         "SshConfigZeroaddress",
	{prefix: "ssh", path: "creds/{role}"}:                                               "SshCredsRole",
	{prefix: "ssh", path: "issue/{role}"}:                                               "SshIssueRole",
	{prefix: "ssh", path: "keys/{key_name}"}:                                            "SshKeysKeyName",
	{prefix: "ssh", path: "lookup"}:                                                     "SshLookup",
	{prefix: "ssh", path: "public_key"}:                                                 "SshPublicKey",
	{prefix: "ssh", path: "roles"}:                                                      "SshRoles",
	{prefix: "ssh", path: "roles/{role}"}:                                               "SshRolesRole",
	{prefix: "ssh", path: "sign/{role}"}:                                                "SshSignRole",
	{prefix: "ssh", path: "verify"}:                                                     "SshVerify",
	{prefix: "sys", path: "audit"}:                                                      "SysAudit",
	{prefix: "sys", path: "audit-hash/{path}"}:                                          "SysAuditHashPath",
	{prefix: "sys", path: "audit/{path}"}:                                               "SysAuditPath",
	{prefix: "sys", path: "auth"}:                                                       "SysAuth",
	{prefix: "sys", path: "auth/{path}"}:                                                "SysAuthPath",
	{prefix: "sys", path: "auth/{path}/tune"}:                                           "SysAuthPathTune",
	{prefix: "sys", path: "capabilities"}:                                               "SysCapabilities",
	{prefix: "sys", path: "capabilities-accessor"}:                                      "SysCapabilitiesAccessor",
	{prefix: "sys", path: "capabilities-self"}:                                          "SysCapabilitiesSelf",
	{prefix: "sys", path: "config/auditing/request-headers"}:                            "SysConfigAuditingRequestHeaders",
	{prefix: "sys", path: "config/auditing/request-headers/{header}"}:                   "SysConfigAuditingRequestHeadersHeader",
	{prefix: "sys", path: "config/cors"}:                                                "SysConfigCors",
	{prefix: "sys", path: "config/reload/{subsystem}"}:                                  "SysConfigReloadSubsystem",
	{prefix: "sys", path: "config/state/sanitized"}:                                     "SysConfigStateSanitized",
	{prefix: "sys", path: "config/ui/headers/"}:                                         "SysConfigUiHeaders",
	{prefix: "sys", path: "config/ui/headers/{header}"}:                                 "SysConfigUiHeadersHeader",
	{prefix: "sys", path: "generate-root"}:                                              "SysGenerateRoot",
	{prefix: "sys", path: "generate-root/attempt"}:                                      "SysGenerateRootAttempt",
	{prefix: "sys", path: "generate-root/update"}:                                       "SysGenerateRootUpdate",
	{prefix: "sys", path: "ha-status"}:                                                  "SysHaStatus",
	{prefix: "sys", path: "health"}:                                                     "SysHealth",
	{prefix: "sys", path: "host-info"}:                                                  "SysHostInfo",
	{prefix: "sys", path: "in-flight-req"}:                                              "SysInFlightReq",
	{prefix: "sys", path: "init"}:                                                       "SysInit",
	{prefix: "sys", path: "internal/counters/activity"}:                                 "SysInternalCountersActivity",
	{prefix: "sys", path: "internal/counters/activity/export"}:                          "SysInternalCountersActivityExport",
	{prefix: "sys", path: "internal/counters/activity/monthly"}:                         "SysInternalCountersActivityMonthly",
	{prefix: "sys", path: "internal/counters/config"}:                                   "SysInternalCountersConfig",
	{prefix: "sys", path: "internal/counters/entities"}:                                 "SysInternalCountersEntities",
	{prefix: "sys", path: "internal/counters/requests"}:                                 "SysInternalCountersRequests",
	{prefix: "sys", path: "internal/counters/tokens"}:                                   "SysInternalCountersTokens",
	{prefix: "sys", path: "internal/inspect/router/{tag}"}:                              "SysInternalInspectRouterTag",
	{prefix: "sys", path: "internal/specs/openapi"}:                                     "SysInternalSpecsOpenapi",
	{prefix: "sys", path: "internal/ui/feature-flags"}:                                  "SysInternalUiFeatureFlags",
	{prefix: "sys", path: "internal/ui/mounts"}:                                         "SysInternalUiMounts",
	{prefix: "sys", path: "internal/ui/mounts/{path}"}:                                  "SysInternalUiMountsPath",
	{prefix: "sys", path: "internal/ui/namespaces"}:                                     "SysInternalUiNamespaces",
	{prefix: "sys", path: "internal/ui/resultant-acl"}:                                  "SysInternalUiResultantAcl",
	{prefix: "sys", path: "key-status"}:                                                 "SysKeyStatus",
	{prefix: "sys", path: "leader"}:                                                     "SysLeader",
	{prefix: "sys", path: "leases"}:                                                     "SysLeases",
	{prefix: "sys", path: "leases/count"}:                                               "SysLeasesCount",
	{prefix: "sys", path: "leases/lookup"}:                                              "SysLeasesLookup",
	{prefix: "sys", path: "leases/lookup/"}:                                             "SysLeasesLookup",
	{prefix: "sys", path: "leases/lookup/{prefix}"}:                                     "SysLeasesLookupPrefix",
	{prefix: "sys", path: "leases/renew"}:                                               "SysLeasesRenew",
	{prefix: "sys", path: "leases/renew/{url_lease_id}"}:                                "SysLeasesRenewUrlLeaseId",
	{prefix: "sys", path: "leases/revoke"}:                                              "SysLeasesRevoke",
	{prefix: "sys", path: "leases/revoke-force/{prefix}"}:                               "SysLeasesRevokeForcePrefix",
	{prefix: "sys", path: "leases/revoke-prefix/{prefix}"}:                              "SysLeasesRevokePrefixPrefix",
	{prefix: "sys", path: "leases/revoke/{url_lease_id}"}:                               "SysLeasesRevokeUrlLeaseId",
	{prefix: "sys", path: "leases/tidy"}:                                                "SysLeasesTidy",
	{prefix: "sys", path: "loggers"}:                                                    "SysLoggers",
	{prefix: "sys", path: "loggers/{name}"}:                                             "SysLoggersName",
	{prefix: "sys", path: "metrics"}:                                                    "SysMetrics",
	{prefix: "sys", path: "mfa/validate"}:                                               "SysMfaValidate",
	{prefix: "sys", path: "monitor"}:                                                    "SysMonitor",
	{prefix: "sys", path: "mounts"}:                                                     "SysMounts",
	{prefix: "sys", path: "mounts/{path}"}:                                              "SysMountsPath",
	{prefix: "sys", path: "mounts/{path}/tune"}:                                         "SysMountsPathTune",
	{prefix: "sys", path: "plugins/catalog"}:                                            "SysPluginsCatalog",
	{prefix: "sys", path: "plugins/catalog/{name}"}:                                     "SysPluginsCatalogName",
	{prefix: "sys", path: "plugins/catalog/{type}"}:                                     "SysPluginsCatalogType",
	{prefix: "sys", path: "plugins/catalog/{type}/{name}"}:                              "SysPluginsCatalogTypeName",
	{prefix: "sys", path: "plugins/reload/backend"}:                                     "SysPluginsReloadBackend",
	{prefix: "sys", path: "policies/acl"}:                                               "SysPoliciesAcl",
	{prefix: "sys", path: "policies/acl/{name}"}:                                        "SysPoliciesAclName",
	{prefix: "sys", path: "policies/password"}:                                          "SysPoliciesPassword",
	{prefix: "sys", path: "policies/password/{name}"}:                                   "SysPoliciesPasswordName",
	{prefix: "sys", path: "policies/password/{name}/generate"}:                          "SysPoliciesPasswordNameGenerate",
	{prefix: "sys", path: "policy"}:                                                     "SysPolicy",
	{prefix: "sys", path: "policy/{name}"}:                                              "SysPolicyName",
	{prefix: "sys", path: "pprof/"}:                                                     "SysPprof",
	{prefix: "sys", path: "pprof/allocs"}:                                               "SysPprofAllocs",
	{prefix: "sys", path: "pprof/block"}:                                                "SysPprofBlock",
	{prefix: "sys", path: "pprof/cmdline"}:                                              "SysPprofCmdline",
	{prefix: "sys", path: "pprof/goroutine"}:                                            "SysPprofGoroutine",
	{prefix: "sys", path: "pprof/heap"}:                                                 "SysPprofHeap",
	{prefix: "sys", path: "pprof/mutex"}:                                                "SysPprofMutex",
	{prefix: "sys", path: "pprof/profile"}:                                              "SysPprofProfile",
	{prefix: "sys", path: "pprof/symbol"}:                                               "SysPprofSymbol",
	{prefix: "sys", path: "pprof/threadcreate"}:                                         "SysPprofThreadcreate",
	{prefix: "sys", path: "pprof/trace"}:                                                "SysPprofTrace",
	{prefix: "sys", path: "quotas/config"}:                                              "SysQuotasConfig",
	{prefix: "sys", path: "quotas/rate-limit"}:                                          "SysQuotasRateLimit",
	{prefix: "sys", path: "quotas/rate-limit/{name}"}:                                   "SysQuotasRateLimitName",
	{prefix: "sys", path: "raw"}:                                                        "SysRaw",
	{prefix: "sys", path: "raw/{path}"}:                                                 "SysRawPath",
	{prefix: "sys", path: "rekey/backup"}:                                               "SysRekeyBackup",
	{prefix: "sys", path: "rekey/init"}:                                                 "SysRekeyInit",
	{prefix: "sys", path: "rekey/recovery-key-backup"}:                                  "SysRekeyRecoveryKeyBackup",
	{prefix: "sys", path: "rekey/update"}:                                               "SysRekeyUpdate",
	{prefix: "sys", path: "rekey/verify"}:                                               "SysRekeyVerify",
	{prefix: "sys", path: "remount"}:                                                    "SysRemount",
	{prefix: "sys", path: "remount/status/{migration_id}"}:                              "SysRemountStatusMigrationId",
	{prefix: "sys", path: "renew"}:                                                      "SysRenew",
	{prefix: "sys", path: "renew/{url_lease_id}"}:                                       "SysRenewUrlLeaseId",
	{prefix: "sys", path: "replication/status"}:                                         "SysReplicationStatus",
	{prefix: "sys", path: "revoke"}:                                                     "SysRevoke",
	{prefix: "sys", path: "revoke-force/{prefix}"}:                                      "SysRevokeForcePrefix",
	{prefix: "sys", path: "revoke-prefix/{prefix}"}:                                     "SysRevokePrefixPrefix",
	{prefix: "sys", path: "revoke/{url_lease_id}"}:                                      "SysRevokeUrlLeaseId",
	{prefix: "sys", path: "rotate"}:                                                     "SysRotate",
	{prefix: "sys", path: "rotate/config"}:                                              "SysRotateConfig",
	{prefix: "sys", path: "seal"}:                                                       "SysSeal",
	{prefix: "sys", path: "seal-status"}:                                                "SysSealStatus",
	{prefix: "sys", path: "step-down"}:                                                  "SysStepDown",
	{prefix: "sys", path: "tools/hash"}:                                                 "SysToolsHash",
	{prefix: "sys", path: "tools/hash/{urlalgorithm}"}:                                  "SysToolsHashUrlalgorithm",
	{prefix: "sys", path: "tools/random"}:                                               "SysToolsRandom",
	{prefix: "sys", path: "tools/random/{source}"}:                                      "SysToolsRandomSource",
	{prefix: "sys", path: "tools/random/{source}/{urlbytes}"}:                           "SysToolsRandomSourceUrlbytes",
	{prefix: "sys", path: "tools/random/{urlbytes}"}:                                    "SysToolsRandomUrlbytes",
	{prefix: "sys", path: "unseal"}:                                                     "SysUnseal",
	{prefix: "sys", path: "version-history/"}:                                           "SysVersionHistory",
	{prefix: "sys", path: "wrapping/lookup"}:                                            "SysWrappingLookup",
	{prefix: "sys", path: "wrapping/rewrap"}:                                            "SysWrappingRewrap",
	{prefix: "sys", path: "wrapping/unwrap"}:                                            "SysWrappingUnwrap",
	{prefix: "sys", path: "wrapping/wrap"}:                                              "SysWrappingWrap",
	{prefix: "terraform", path: "config"}:                                               "TerraformConfig",
	{prefix: "terraform", path: "creds/{name}"}:                                         "TerraformCredsName",
	{prefix: "terraform", path: "role"}:                                                 "TerraformRole",
	{prefix: "terraform", path: "role/{name}"}:                                          "TerraformRoleName",
	{prefix: "terraform", path: "rotate-role/{name}"}:                                   "TerraformRotateRoleName",
	{prefix: "totp", path: "code/{name}"}:                                               "TotpCodeName",
	{prefix: "totp", path: "keys"}:                                                      "TotpKeys",
	{prefix: "totp", path: "keys/{name}"}:                                               "TotpKeysName",
	{prefix: "transit", path: "backup/{name}"}:                                          "TransitBackupName",
	{prefix: "transit", path: "cache-config"}:                                           "TransitCacheConfig",
	{prefix: "transit", path: "datakey/{plaintext}/{name}"}:                             "TransitDatakeyPlaintextName",
	{prefix: "transit", path: "decrypt/{name}"}:                                         "TransitDecryptName",
	{prefix: "transit", path: "encrypt/{name}"}:                                         "TransitEncryptName",
	{prefix: "transit", path: "export/{type}/{name}"}:                                   "TransitExportTypeName",
	{prefix: "transit", path: "export/{type}/{name}/{version}"}:                         "TransitExportTypeNameVersion",
	{prefix: "transit", path: "hash"}:                                                   "TransitHash",
	{prefix: "transit", path: "hash/{urlalgorithm}"}:                                    "TransitHashUrlalgorithm",
	{prefix: "transit", path: "hmac/{name}"}:                                            "TransitHmacName",
	{prefix: "transit", path: "hmac/{name}/{urlalgorithm}"}:                             "TransitHmacNameUrlalgorithm",
	{prefix: "transit", path: "keys"}:                                                   "TransitKeys",
	{prefix: "transit", path: "keys/{name}"}:                                            "TransitKeysName",
	{prefix: "transit", path: "keys/{name}/config"}:                                     "TransitKeysNameConfig",
	{prefix: "transit", path: "keys/{name}/import"}:                                     "TransitKeysNameImport",
	{prefix: "transit", path: "keys/{name}/import_version"}:                             "TransitKeysNameImportVersion",
	{prefix: "transit", path: "keys/{name}/rotate"}:                                     "TransitKeysNameRotate",
	{prefix: "transit", path: "keys/{name}/trim"}:                                       "TransitKeysNameTrim",
	{prefix: "transit", path: "random"}:                                                 "TransitRandom",
	{prefix: "transit", path: "random/{source}"}:                                        "TransitRandomSource",
	{prefix: "transit", path: "random/{source}/{urlbytes}"}:                             "TransitRandomSourceUrlbytes",
	{prefix: "transit", path: "random/{urlbytes}"}:                                      "TransitRandomUrlbytes",
	{prefix: "transit", path: "restore"}:                                                "TransitRestore",
	{prefix: "transit", path: "restore/{name}"}:                                         "TransitRestoreName",
	{prefix: "transit", path: "rewrap/{name}"}:                                          "TransitRewrapName",
	{prefix: "transit", path: "sign/{name}"}:                                            "TransitSignName",
	{prefix: "transit", path: "sign/{name}/{urlalgorithm}"}:                             "TransitSignNameUrlalgorithm",
	{prefix: "transit", path: "verify/{name}"}:                                          "TransitVerifyName",
	{prefix: "transit", path: "verify/{name}/{urlalgorithm}"}:                           "TransitVerifyNameUrlalgorithm",
	{prefix: "transit", path: "wrapping_key"}:                                           "TransitWrappingKey",
}

// constructRequestIdentifier joins the given inputs into a title case string,
// e.g. 'UpdateSecretConfigLeaseRequest'. This function is used to generate:
//
//   - operation id
//   - request name
//   - response name
//
// For certain prefix + path combinations, which would otherwise result in an
// ugly string, the function uses a custom lookup table to construct part of
// the string instead.
func constructRequestIdentifier(operation logical.Operation, path, prefix, suffix string) string {
	var parts []string

	// Split the operation by non-word characters
	parts = append(parts, nonWordRe.Split(strings.ToLower(string(operation)), -1)...)

	// Append either the known mapping or prefix + path split by non-word characters
	if mapping, ok := knownPathMappings[knownPathKey{prefix: prefix, path: path}]; ok {
		parts = append(parts, mapping)
	} else {
		parts = append(parts, nonWordRe.Split(strings.ToLower(prefix), -1)...)
		parts = append(parts, nonWordRe.Split(strings.ToLower(path), -1)...)
	}

	parts = append(parts, suffix)

	// Title case everything & join the result into a string
	title := cases.Title(language.English)

	for i, s := range parts {
		parts[i] = title.String(s)
	}

	return strings.Join(parts, "")
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

// expandPattern expands a regex pattern by generating permutations of any optional parameters
// and changing named parameters into their {openapi} equivalents.
func expandPattern(pattern string) []string {
	var paths []string

	// Determine if the pattern starts with an alternation for multiple roots
	// example (root1|root2)/(?P<name>regex) -> match['(root1|root2)/(?P<name>regex)','root1|root2','/(?P<name>regex)']
	match := altRootsRe.FindStringSubmatch(pattern)
	if len(match) == 3 {
		var expandedRoots []string
		for _, root := range strings.Split(match[1], "|") {
			expandedRoots = append(expandedRoots, expandPattern(root+match[2])...)
		}
		return expandedRoots
	}

	// GenericNameRegex adds a regex that complicates our parsing. It is much easier to
	// detect and remove it now than to compensate for in the other regexes.
	//
	// example: (?P<foo>\\w(([\\w-.]+)?\\w)?) -> (?P<foo>)
	base := GenericNameRegex("")
	start := strings.Index(base, ">")
	end := strings.LastIndex(base, ")")
	regexToRemove := ""
	if start != -1 && end != -1 && end > start {
		regexToRemove = base[start+1 : end]
	}
	pattern = strings.ReplaceAll(pattern, regexToRemove, "")

	// Simplify named fields that have limited options, e.g. (?P<foo>a|b|c) -> (<P<foo>.+)
	pattern = altFieldsGroupRe.ReplaceAllStringFunc(pattern, func(s string) string {
		return altFieldsRe.ReplaceAllString(s, ".+")
	})

	// Initialize paths with the original pattern or the halves of an
	// alternation, which is also present in some patterns.
	matches := altRe.FindAllStringSubmatch(pattern, -1)
	if len(matches) > 0 {
		paths = []string{matches[0][1], matches[0][2]}
	} else {
		paths = []string{pattern}
	}

	// Expand all optional regex elements into two paths. This approach is really only useful up to 2 optional
	// groups, but we probably don't want to deal with the exponential increase beyond that anyway.
	for i := 0; i < len(paths); i++ {
		p := paths[i]

		// match is a 2-element slice that will have a start and end index
		// for the left-most match of a regex of form: (lease/)?
		match := optRe.FindStringIndex(p)

		if match != nil {
			// create a path that includes the optional element but without
			// parenthesis or the '?' character.
			paths[i] = p[:match[0]] + p[match[0]+1:match[1]-2] + p[match[1]:]

			// create a path that excludes the optional element.
			paths = append(paths, p[:match[0]]+p[match[1]:])
			i--
		}
	}

	// Replace named parameters (?P<foo>) with {foo}
	var replacedPaths []string

	for _, path := range paths {
		result := reqdRe.FindAllStringSubmatch(path, -1)
		if result != nil {
			for _, p := range result {
				par := p[1]
				path = strings.Replace(path, p[0], fmt.Sprintf("{%s}", par), 1)
			}
		}
		// Final cleanup
		path = cleanSuffixRe.ReplaceAllString(path, "")
		path = cleanCharsRe.ReplaceAllString(path, "")
		replacedPaths = append(replacedPaths, path)
	}

	return replacedPaths
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
