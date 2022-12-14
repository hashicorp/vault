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
			op.OperationID = constructRequestResponseIdentifier(opType, path, requestResponsePrefix, "")

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
					requestName := constructRequestResponseIdentifier(opType, path, requestResponsePrefix, "request")
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
						responseName := constructRequestResponseIdentifier(opType, path, requestResponsePrefix, "response")
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
	mount string
	path  string
}

var knownPathMappings = map[knownPathKey]string{
	{mount: "auth/alicloud", path: "login"}:                                            "AliCloudLogin",
	{mount: "auth/alicloud", path: "role"}:                                             "AliCloudAuthRoles",
	{mount: "auth/alicloud", path: "role/{role}"}:                                      "AliCloudAuthRole",
	{mount: "auth/alicloud", path: "roles"}:                                            "AliCloudAuthRoles2",
	{mount: "auth/approle", path: "login"}:                                             "AppRoleLogin",
	{mount: "auth/approle", path: "role"}:                                              "AppRoleRoles",
	{mount: "auth/approle", path: "role/{role_name}"}:                                  "AppRoleRole",
	{mount: "auth/approle", path: "role/{role_name}/bind-secret-id"}:                   "AppRoleBindSecretID",
	{mount: "auth/approle", path: "role/{role_name}/bound-cidr-list"}:                  "AppRoleBoundCIDRList",
	{mount: "auth/approle", path: "role/{role_name}/custom-secret-id"}:                 "AppRoleCustomSecretID",
	{mount: "auth/approle", path: "role/{role_name}/local-secret-ids"}:                 "AppRoleLocalSecretIDs",
	{mount: "auth/approle", path: "role/{role_name}/period"}:                           "AppRolePeriod",
	{mount: "auth/approle", path: "role/{role_name}/policies"}:                         "AppRolePolicies",
	{mount: "auth/approle", path: "role/{role_name}/role-id"}:                          "AppRoleRoleID",
	{mount: "auth/approle", path: "role/{role_name}/secret-id"}:                        "AppRoleSecretID",
	{mount: "auth/approle", path: "role/{role_name}/secret-id-accessor/destroy"}:       "AppRoleSecretIDAccessorDestroy",
	{mount: "auth/approle", path: "role/{role_name}/secret-id-accessor/lookup"}:        "AppRoleSecretIDAccessorLookup",
	{mount: "auth/approle", path: "role/{role_name}/secret-id-bound-cidrs"}:            "AppRoleSecretIDBoundCIDRs",
	{mount: "auth/approle", path: "role/{role_name}/secret-id-num-uses"}:               "AppRoleSecretIDNumberOfUses",
	{mount: "auth/approle", path: "role/{role_name}/secret-id-ttl"}:                    "AppRoleSecretIDTTL",
	{mount: "auth/approle", path: "role/{role_name}/secret-id/destroy"}:                "AppRoleSecretIDDestroy",
	{mount: "auth/approle", path: "role/{role_name}/secret-id/lookup"}:                 "AppRoleSecretIDLookup",
	{mount: "auth/approle", path: "role/{role_name}/token-bound-cidrs"}:                "AppRoleTokenBoundCIDRs",
	{mount: "auth/approle", path: "role/{role_name}/token-max-ttl"}:                    "AppRoleTokenMaxTTL",
	{mount: "auth/approle", path: "role/{role_name}/token-num-uses"}:                   "AppRoleTokenNumberOfUses",
	{mount: "auth/approle", path: "role/{role_name}/token-ttl"}:                        "AppRoleTokenTTL",
	{mount: "auth/approle", path: "tidy/secret-id"}:                                    "AppRoleTidySecretID",
	{mount: "auth/aws", path: "config/certificate/{cert_name}"}:                        "AWSConfigCertificate",
	{mount: "auth/aws", path: "config/certificates"}:                                   "AWSConfigCertificates",
	{mount: "auth/aws", path: "config/client"}:                                         "AWSConfigClient",
	{mount: "auth/aws", path: "config/identity"}:                                       "AWSConfigIdentity",
	{mount: "auth/aws", path: "config/rotate-root"}:                                    "AWSConfigRotateRoot",
	{mount: "auth/aws", path: "config/sts"}:                                            "AWSConfigSecurityTokenService",
	{mount: "auth/aws", path: "config/sts/{account_id}"}:                               "AWSConfigSecurityTokenServiceAccount",
	{mount: "auth/aws", path: "config/tidy/identity-accesslist"}:                       "AWSConfigTidyIdentityAccesslist",
	{mount: "auth/aws", path: "config/tidy/identity-whitelist"}:                        "AWSConfigTidyIdentityWhitelist",
	{mount: "auth/aws", path: "config/tidy/roletag-blacklist"}:                         "AWSConfigTidyRoleTagBlacklist",
	{mount: "auth/aws", path: "config/tidy/roletag-denylist"}:                          "AWSConfigTidyRoleTagDenylist",
	{mount: "auth/aws", path: "identity-accesslist"}:                                   "AWSIdentityAccesslist",
	{mount: "auth/aws", path: "identity-accesslist/{instance_id}"}:                     "AWSIdentityAccesslistFor",
	{mount: "auth/aws", path: "identity-whitelist"}:                                    "AWSIdentityWhitelist",
	{mount: "auth/aws", path: "identity-whitelist/{instance_id}"}:                      "AWSIdentityWhitelistFor",
	{mount: "auth/aws", path: "login"}:                                                 "AWSLogin",
	{mount: "auth/aws", path: "role"}:                                                  "AWSAuthRoles",
	{mount: "auth/aws", path: "role/{role}"}:                                           "AWSAuthRole",
	{mount: "auth/aws", path: "role/{role}/tag"}:                                       "AWSAuthRoleTag",
	{mount: "auth/aws", path: "roles"}:                                                 "AWSAuthRoles2",
	{mount: "auth/aws", path: "roletag-blacklist"}:                                     "AWSRoleTagBlacklist",
	{mount: "auth/aws", path: "roletag-blacklist/{role_tag}"}:                          "AWSRoleTagBlacklistFor",
	{mount: "auth/aws", path: "roletag-denylist"}:                                      "AWSRoleTagDenylist",
	{mount: "auth/aws", path: "roletag-denylist/{role_tag}"}:                           "AWSRoleTagDenylistFor",
	{mount: "auth/aws", path: "tidy/identity-accesslist"}:                              "AWSTidyIdentityAccesslist",
	{mount: "auth/aws", path: "tidy/identity-whitelist"}:                               "AWSTidyIdentityWhitelist",
	{mount: "auth/aws", path: "tidy/roletag-blacklist"}:                                "AWSTidyRoleTagBlacklist",
	{mount: "auth/aws", path: "tidy/roletag-denylist"}:                                 "AWSTidyRoleTagDenylist",
	{mount: "auth/azure", path: "config"}:                                              "AzureConfig",
	{mount: "auth/azure", path: "login"}:                                               "AzureLogin",
	{mount: "auth/azure", path: "role"}:                                                "AzureRoles",
	{mount: "auth/azure", path: "role/{name}"}:                                         "AzureRole",
	{mount: "auth/centrify", path: "config"}:                                           "CentrifyConfig",
	{mount: "auth/centrify", path: "login"}:                                            "CentrifyLogin",
	{mount: "auth/cert", path: "certs"}:                                                "Certificates",
	{mount: "auth/cert", path: "certs/{name}"}:                                         "Certificate",
	{mount: "auth/cert", path: "config"}:                                               "CertificateConfig",
	{mount: "auth/cert", path: "crls"}:                                                 "CertificateCRLs",
	{mount: "auth/cert", path: "crls/{name}"}:                                          "CertificateCRL",
	{mount: "auth/cert", path: "login"}:                                                "CertificateLogin",
	{mount: "auth/cf", path: "config"}:                                                 "CloudFoundryConfig",
	{mount: "auth/cf", path: "login"}:                                                  "CloudFoundryLogin",
	{mount: "auth/cf", path: "roles"}:                                                  "CloudFoundryRoles",
	{mount: "auth/cf", path: "roles/{role}"}:                                           "CloudFoundryRole",
	{mount: "auth/gcp", path: "config"}:                                                "GoogleCloudConfig",
	{mount: "auth/gcp", path: "login"}:                                                 "GoogleCloudLogin",
	{mount: "auth/gcp", path: "role"}:                                                  "GoogleCloudRoles",
	{mount: "auth/gcp", path: "role/{name}"}:                                           "GoogleCloudRole",
	{mount: "auth/gcp", path: "role/{name}/labels"}:                                    "GoogleCloudRoleLabels",
	{mount: "auth/gcp", path: "role/{name}/service-accounts"}:                          "GoogleCloudRoleServiceAccounts",
	{mount: "auth/gcp", path: "roles"}:                                                 "GoogleCloudRoles2",
	{mount: "auth/github", path: "config"}:                                             "GitHubConfig",
	{mount: "auth/github", path: "login"}:                                              "GitHubLogin",
	{mount: "auth/github", path: "map/teams"}:                                          "GitHubMapTeams",
	{mount: "auth/github", path: "map/teams/{key}"}:                                    "GitHubMapTeam",
	{mount: "auth/github", path: "map/users"}:                                          "GitHubMapUsers",
	{mount: "auth/github", path: "map/users/{key}"}:                                    "GitHubMapUser",
	{mount: "auth/jwt", path: "config"}:                                                "JWTConfig",
	{mount: "auth/jwt", path: "login"}:                                                 "JWTLogin",
	{mount: "auth/jwt", path: "oidc/auth_url"}:                                         "JWTOIDCAuthURL",
	{mount: "auth/jwt", path: "oidc/callback"}:                                         "JWTOIDCCallback",
	{mount: "auth/jwt", path: "role"}:                                                  "JWTRoles",
	{mount: "auth/jwt", path: "role/{name}"}:                                           "JWTRole",
	{mount: "auth/kerberos", path: "config"}:                                           "KerberosConfig",
	{mount: "auth/kerberos", path: "config/ldap"}:                                      "KerberosConfigLDAP",
	{mount: "auth/kerberos", path: "groups"}:                                           "KerberosGroups",
	{mount: "auth/kerberos", path: "groups/{name}"}:                                    "KerberosGroup",
	{mount: "auth/kerberos", path: "login"}:                                            "KerberosLogin",
	{mount: "auth/kubernetes", path: "config"}:                                         "KubernetesConfig",
	{mount: "auth/kubernetes", path: "login"}:                                          "KubernetesLogin",
	{mount: "auth/kubernetes", path: "role"}:                                           "KubernetesRoles",
	{mount: "auth/kubernetes", path: "role/{name}"}:                                    "KubernetesRole",
	{mount: "auth/ldap", path: "config"}:                                               "LDAPConfig",
	{mount: "auth/ldap", path: "groups"}:                                               "LDAPGroups",
	{mount: "auth/ldap", path: "groups/{name}"}:                                        "LDAPGroup",
	{mount: "auth/ldap", path: "login/{username}"}:                                     "LDAPLogin",
	{mount: "auth/ldap", path: "users"}:                                                "LDAPUsers",
	{mount: "auth/ldap", path: "users/{name}"}:                                         "LDAPUser",
	{mount: "auth/oci", path: "config"}:                                                "OCIConfig",
	{mount: "auth/oci", path: "login"}:                                                 "OCILogin",
	{mount: "auth/oci", path: "login/{role}"}:                                          "OCILoginWithRole",
	{mount: "auth/oci", path: "role"}:                                                  "OCIRoles",
	{mount: "auth/oci", path: "role/{role}"}:                                           "OCIRole",
	{mount: "auth/oidc", path: "config"}:                                               "OIDCConfig",
	{mount: "auth/oidc", path: "login"}:                                                "OIDCLogin",
	{mount: "auth/oidc", path: "oidc/auth_url"}:                                        "OIDCAuthURL",
	{mount: "auth/oidc", path: "oidc/callback"}:                                        "OIDCCallback",
	{mount: "auth/oidc", path: "role"}:                                                 "OIDCRoles",
	{mount: "auth/oidc", path: "role/{name}"}:                                          "OIDCRole",
	{mount: "auth/okta", path: "config"}:                                               "OktaConfig",
	{mount: "auth/okta", path: "groups"}:                                               "OktaGroups",
	{mount: "auth/okta", path: "groups/{name}"}:                                        "OktaGroup",
	{mount: "auth/okta", path: "login/{username}"}:                                     "OktaLogin",
	{mount: "auth/okta", path: "users"}:                                                "OktaUsers",
	{mount: "auth/okta", path: "users/{name}"}:                                         "OktaUser",
	{mount: "auth/okta", path: "verify/{nonce}"}:                                       "OktaVerify",
	{mount: "auth/radius", path: "config"}:                                             "RadiusConfig",
	{mount: "auth/radius", path: "login"}:                                              "RadiusLogin",
	{mount: "auth/radius", path: "login/{urlusername}"}:                                "RadiusLoginWithUsername",
	{mount: "auth/radius", path: "users"}:                                              "RadiusUsers",
	{mount: "auth/radius", path: "users/{name}"}:                                       "RadiusUser",
	{mount: "auth/token", path: "accessors/"}:                                          "TokenAccessors",
	{mount: "auth/token", path: "create"}:                                              "TokenCreate",
	{mount: "auth/token", path: "create-orphan"}:                                       "TokenCreateOrphan",
	{mount: "auth/token", path: "create/{role_name}"}:                                  "TokenCreateWithRole",
	{mount: "auth/token", path: "lookup"}:                                              "TokenLookup",
	{mount: "auth/token", path: "lookup-accessor"}:                                     "TokenLookupAccessor",
	{mount: "auth/token", path: "lookup-self"}:                                         "TokenLookupSelf",
	{mount: "auth/token", path: "renew"}:                                               "TokenRenew",
	{mount: "auth/token", path: "renew-accessor"}:                                      "TokenRenewAccessor",
	{mount: "auth/token", path: "renew-self"}:                                          "TokenRenewSelf",
	{mount: "auth/token", path: "revoke"}:                                              "TokenRevoke",
	{mount: "auth/token", path: "revoke-accessor"}:                                     "TokenRevokeAccessor",
	{mount: "auth/token", path: "revoke-orphan"}:                                       "TokenRevokeOrphan",
	{mount: "auth/token", path: "revoke-self"}:                                         "TokenRevokeSelf",
	{mount: "auth/token", path: "roles"}:                                               "TokenRoles",
	{mount: "auth/token", path: "roles/{role_name}"}:                                   "TokenRole",
	{mount: "auth/token", path: "tidy"}:                                                "TokenTidy",
	{mount: "auth/userpass", path: "login/{username}"}:                                 "UserpassLogin",
	{mount: "auth/userpass", path: "users"}:                                            "UserpassUsers",
	{mount: "auth/userpass", path: "users/{username}"}:                                 "UserpassUser",
	{mount: "auth/userpass", path: "users/{username}/password"}:                        "UserpassUserPassword",
	{mount: "auth/userpass", path: "users/{username}/policies"}:                        "UserpassUserPolicies",
	{mount: "ad", path: "config"}:                                                      "ActiveDirectoryConfig",
	{mount: "ad", path: "creds/{name}"}:                                                "ActiveDirectoryCredentials",
	{mount: "ad", path: "library"}:                                                     "ActiveDirectoryLibraries",
	{mount: "ad", path: "library/manage/{name}/check-in"}:                              "ActiveDirectoryLibraryManageCheckIn",
	{mount: "ad", path: "library/{name}"}:                                              "ActiveDirectoryLibrary",
	{mount: "ad", path: "library/{name}/check-in"}:                                     "ActiveDirectoryLibraryCheckIn",
	{mount: "ad", path: "library/{name}/check-out"}:                                    "ActiveDirectoryLibraryCheckOut",
	{mount: "ad", path: "library/{name}/status"}:                                       "ActiveDirectoryLibraryStatus",
	{mount: "ad", path: "roles"}:                                                       "ActiveDirectoryRoles",
	{mount: "ad", path: "roles/{name}"}:                                                "ActiveDirectoryRole",
	{mount: "ad", path: "rotate-role/{name}"}:                                          "ActiveDirectoryRotateRole",
	{mount: "ad", path: "rotate-root"}:                                                 "ActiveDirectoryRotateRoot",
	{mount: "alicloud", path: "config"}:                                                "AliCloudConfig",
	{mount: "alicloud", path: "creds/{name}"}:                                          "AliCloudCredentials",
	{mount: "alicloud", path: "role"}:                                                  "AliCloudRoles",
	{mount: "alicloud", path: "role/{name}"}:                                           "AliCloudRole",
	{mount: "aws", path: "config/lease"}:                                               "AWSConfigLease",
	{mount: "aws", path: "config/root"}:                                                "AWSConfigRoot",
	{mount: "aws", path: "config/rotate-root"}:                                         "AWSConfigRotateRoot",
	{mount: "aws", path: "creds"}:                                                      "AWSCredentials",
	{mount: "aws", path: "roles"}:                                                      "AWSRoles",
	{mount: "aws", path: "roles/{name}"}:                                               "AWSRole",
	{mount: "aws", path: "sts/{name}"}:                                                 "AWSSecurityTokenService",
	{mount: "azure", path: "config"}:                                                   "AzureConfig",
	{mount: "azure", path: "creds/{role}"}:                                             "AzureCredentials",
	{mount: "azure", path: "roles"}:                                                    "AzureRoles",
	{mount: "azure", path: "roles/{name}"}:                                             "AzureRole",
	{mount: "azure", path: "rotate-root"}:                                              "AzureRotateRoot",
	{mount: "consul", path: "config/access"}:                                           "ConsulConfigAccess",
	{mount: "consul", path: "creds/{role}"}:                                            "ConsulCredentials",
	{mount: "consul", path: "roles"}:                                                   "ConsulRoles",
	{mount: "consul", path: "roles/{name}"}:                                            "ConsulRole",
	{mount: "cubbyhole", path: "{path}"}:                                               "Cubbyhole",
	{mount: "gcp", path: "config"}:                                                     "GoogleCloudConfig",
	{mount: "gcp", path: "config/rotate-root"}:                                         "GoogleCloudConfigRotateRoot",
	{mount: "gcp", path: "key/{roleset}"}:                                              "GoogleCloudKey",
	{mount: "gcp", path: "roleset/{name}"}:                                             "GoogleCloudRoleset",
	{mount: "gcp", path: "roleset/{name}/rotate"}:                                      "GoogleCloudRolesetRotate",
	{mount: "gcp", path: "roleset/{name}/rotate-key"}:                                  "GoogleCloudRolesetRotateKey",
	{mount: "gcp", path: "roleset/{roleset}/key"}:                                      "GoogleCloudRolesetKey",
	{mount: "gcp", path: "roleset/{roleset}/token"}:                                    "GoogleCloudRolesetToken",
	{mount: "gcp", path: "rolesets"}:                                                   "GoogleCloudRolesets",
	{mount: "gcp", path: "static-account/{name}"}:                                      "GoogleCloudStaticAccount",
	{mount: "gcp", path: "static-account/{name}/key"}:                                  "GoogleCloudStaticAccountKey",
	{mount: "gcp", path: "static-account/{name}/rotate-key"}:                           "GoogleCloudStaticAccountRotateKey",
	{mount: "gcp", path: "static-account/{name}/token"}:                                "GoogleCloudStaticAccountToken",
	{mount: "gcp", path: "static-accounts"}:                                            "GoogleCloudStaticAccounts",
	{mount: "gcp", path: "token/{roleset}"}:                                            "GoogleCloudToken",
	{mount: "gcpkms", path: "config"}:                                                  "GoogleCloudKMSConfig",
	{mount: "gcpkms", path: "decrypt/{key}"}:                                           "GoogleCloudKMSDecrypt",
	{mount: "gcpkms", path: "encrypt/{key}"}:                                           "GoogleCloudKMSEncrypt",
	{mount: "gcpkms", path: "keys"}:                                                    "GoogleCloudKMSKeys",
	{mount: "gcpkms", path: "keys/config/{key}"}:                                       "GoogleCloudKMSKeysConfig",
	{mount: "gcpkms", path: "keys/deregister/{key}"}:                                   "GoogleCloudKMSKeysDeregister",
	{mount: "gcpkms", path: "keys/register/{key}"}:                                     "GoogleCloudKMSKeysRegister",
	{mount: "gcpkms", path: "keys/rotate/{key}"}:                                       "GoogleCloudKMSKeysRotate",
	{mount: "gcpkms", path: "keys/trim/{key}"}:                                         "GoogleCloudKMSKeysTrim",
	{mount: "gcpkms", path: "keys/{key}"}:                                              "GoogleCloudKMSKey",
	{mount: "gcpkms", path: "pubkey/{key}"}:                                            "GoogleCloudKMSPubkey",
	{mount: "gcpkms", path: "reencrypt/{key}"}:                                         "GoogleCloudKMSReencrypt",
	{mount: "gcpkms", path: "sign/{key}"}:                                              "GoogleCloudKMSSign",
	{mount: "gcpkms", path: "verify/{key}"}:                                            "GoogleCloudKMSVerify",
	{mount: "identity", path: "alias"}:                                                 "IdentityAlias",
	{mount: "identity", path: "alias/id"}:                                              "IdentityAliasesByID",
	{mount: "identity", path: "alias/id/{id}"}:                                         "IdentityAliasByID",
	{mount: "identity", path: "entity"}:                                                "IdentityEntity",
	{mount: "identity", path: "entity-alias"}:                                          "IdentityEntityAlias",
	{mount: "identity", path: "entity-alias/id"}:                                       "IdentityEntityAliasesByID",
	{mount: "identity", path: "entity-alias/id/{id}"}:                                  "IdentityEntityAliasByID",
	{mount: "identity", path: "entity/batch-delete"}:                                   "IdentityEntityBatchDelete",
	{mount: "identity", path: "entity/id"}:                                             "IdentityEntitiesByID",
	{mount: "identity", path: "entity/id/{id}"}:                                        "IdentityEntityByID",
	{mount: "identity", path: "entity/merge"}:                                          "IdentityEntityMerge",
	{mount: "identity", path: "entity/name"}:                                           "IdentityEntitiesByName",
	{mount: "identity", path: "entity/name/{name}"}:                                    "IdentityEntityByName",
	{mount: "identity", path: "group"}:                                                 "IdentityGroup",
	{mount: "identity", path: "group-alias"}:                                           "IdentityGroupAlias",
	{mount: "identity", path: "group-alias/id"}:                                        "IdentityGroupAliasesByID",
	{mount: "identity", path: "group-alias/id/{id}"}:                                   "IdentityGroupAliasByID",
	{mount: "identity", path: "group/id"}:                                              "IdentityGroupsByID",
	{mount: "identity", path: "group/id/{id}"}:                                         "IdentityGroupByID",
	{mount: "identity", path: "group/name"}:                                            "IdentityGroupsByName",
	{mount: "identity", path: "group/name/{name}"}:                                     "IdentityGroupByName",
	{mount: "identity", path: "lookup/entity"}:                                         "IdentityLookupEntity",
	{mount: "identity", path: "lookup/group"}:                                          "IdentityLookupGroup",
	{mount: "identity", path: "mfa/login-enforcement"}:                                 "IdentityMFALoginEnforcements",
	{mount: "identity", path: "mfa/login-enforcement/{name}"}:                          "IdentityMFALoginEnforcement",
	{mount: "identity", path: "mfa/method"}:                                            "IdentityMFAMethods",
	{mount: "identity", path: "mfa/method/duo"}:                                        "IdentityMFAMethodsDuo",
	{mount: "identity", path: "mfa/method/duo/{method_id}"}:                            "IdentityMFAMethodDuo",
	{mount: "identity", path: "mfa/method/okta"}:                                       "IdentityMFAMethodsOkta",
	{mount: "identity", path: "mfa/method/okta/{method_id}"}:                           "IdentityMFAMethodOkta",
	{mount: "identity", path: "mfa/method/pingid"}:                                     "IdentityMFAMethodsPingID",
	{mount: "identity", path: "mfa/method/pingid/{method_id}"}:                         "IdentityMFAMethodPingID",
	{mount: "identity", path: "mfa/method/totp"}:                                       "IdentityMFAMethodsTOTP",
	{mount: "identity", path: "mfa/method/totp/admin-destroy"}:                         "IdentityMFAMethodTOTPAdminDestroy",
	{mount: "identity", path: "mfa/method/totp/admin-generate"}:                        "IdentityMFAMethodTOTPAdminGenerate",
	{mount: "identity", path: "mfa/method/totp/generate"}:                              "IdentityMFAMethodTOTPGenerate",
	{mount: "identity", path: "mfa/method/totp/{method_id}"}:                           "IdentityMFAMethodTOTP",
	{mount: "identity", path: "mfa/method/{method_id}"}:                                "IdentityMFAMethod",
	{mount: "identity", path: "oidc/.well-known/keys"}:                                 "IdentityOIDCWellKnownKeys",
	{mount: "identity", path: "oidc/.well-known/openid-configuration"}:                 "IdentityOIDCWellKnownOpenIDConfiguration",
	{mount: "identity", path: "oidc/assignment"}:                                       "IdentityOIDCAssignments",
	{mount: "identity", path: "oidc/assignment/{name}"}:                                "IdentityOIDCAssignment",
	{mount: "identity", path: "oidc/client"}:                                           "IdentityOIDCClients",
	{mount: "identity", path: "oidc/client/{name}"}:                                    "IdentityOIDCClient",
	{mount: "identity", path: "oidc/config"}:                                           "IdentityOIDCConfig",
	{mount: "identity", path: "oidc/introspect"}:                                       "IdentityOIDCIntrospect",
	{mount: "identity", path: "oidc/key"}:                                              "IdentityOIDCKeys",
	{mount: "identity", path: "oidc/key/{name}"}:                                       "IdentityOIDCKey",
	{mount: "identity", path: "oidc/key/{name}/rotate"}:                                "IdentityOIDCKeyRotate",
	{mount: "identity", path: "oidc/provider"}:                                         "IdentityOIDCProviders",
	{mount: "identity", path: "oidc/provider/{name}"}:                                  "IdentityOIDCProvider",
	{mount: "identity", path: "oidc/provider/{name}/.well-known/keys"}:                 "IdentityOIDCProviderWellKnownKeys",
	{mount: "identity", path: "oidc/provider/{name}/.well-known/openid-configuration"}: "IdentityOIDCProviderWellKnownOpenIDConfiguration",
	{mount: "identity", path: "oidc/provider/{name}/authorize"}:                        "IdentityOIDCProviderAuthorize",
	{mount: "identity", path: "oidc/provider/{name}/token"}:                            "IdentityOIDCProviderToken",
	{mount: "identity", path: "oidc/provider/{name}/userinfo"}:                         "IdentityOIDCProviderUserInfo",
	{mount: "identity", path: "oidc/role"}:                                             "IdentityOIDCRoles",
	{mount: "identity", path: "oidc/role/{name}"}:                                      "IdentityOIDCRole",
	{mount: "identity", path: "oidc/scope"}:                                            "IdentityOIDCScopes",
	{mount: "identity", path: "oidc/scope/{name}"}:                                     "IdentityOIDCScope",
	{mount: "identity", path: "oidc/token/{name}"}:                                     "IdentityOIDCToken",
	{mount: "identity", path: "persona"}:                                               "IdentityPersona",
	{mount: "identity", path: "persona/id"}:                                            "IdentityPersonaIDs",
	{mount: "identity", path: "persona/id/{id}"}:                                       "IdentityPersonaID",
	{mount: "kubernetes", path: "config"}:                                              "KubernetesConfig",
	{mount: "kubernetes", path: "creds/{name}"}:                                        "KubernetesCredentials",
	{mount: "kubernetes", path: "roles"}:                                               "KubernetesRoles",
	{mount: "kubernetes", path: "roles/{name}"}:                                        "KubernetesRole",
	{mount: "kv", path: "{path}"}:                                                      "KVv1Secret",
	{mount: "ldap", path: "config"}:                                                    "LDAPConfig",
	{mount: "ldap", path: "creds/{name}"}:                                              "LDAPCredentials",
	{mount: "ldap", path: "library"}:                                                   "LDAPLibraries",
	{mount: "ldap", path: "library/manage/{name}/check-in"}:                            "LDAPLibraryManageCheckIn",
	{mount: "ldap", path: "library/{name}"}:                                            "LDAPLibrary",
	{mount: "ldap", path: "library/{name}/check-in"}:                                   "LDAPLibraryCheckIn",
	{mount: "ldap", path: "library/{name}/check-out"}:                                  "LDAPLibraryCheckOut",
	{mount: "ldap", path: "library/{name}/status"}:                                     "LDAPLibraryStatus",
	{mount: "ldap", path: "role"}:                                                      "LDAPRoles",
	{mount: "ldap", path: "role/{name}"}:                                               "LDAPRole",
	{mount: "ldap", path: "rotate-role/{name}"}:                                        "LDAPRotateRole",
	{mount: "ldap", path: "rotate-root"}:                                               "LDAPRotateRoot",
	{mount: "ldap", path: "static-cred/{name}"}:                                        "LDAPStaticCredentials",
	{mount: "ldap", path: "static-role"}:                                               "LDAPStaticRoles",
	{mount: "ldap", path: "static-role/{name}"}:                                        "LDAPStaticRole",
	{mount: "mongodbatlas", path: "config"}:                                            "MongoDBAtlasConfig",
	{mount: "mongodbatlas", path: "creds/{name}"}:                                      "MongoDBAtlasCredentials",
	{mount: "mongodbatlas", path: "roles"}:                                             "MongoDBAtlasRoles",
	{mount: "mongodbatlas", path: "roles/{name}"}:                                      "MongoDBAtlasRole",
	{mount: "nomad", path: "config/access"}:                                            "NomadConfigAccess",
	{mount: "nomad", path: "config/lease"}:                                             "NomadConfigLease",
	{mount: "nomad", path: "creds/{name}"}:                                             "NomadCredentials",
	{mount: "nomad", path: "role"}:                                                     "NomadRoles",
	{mount: "nomad", path: "role/{name}"}:                                              "NomadRole",
	{mount: "openldap", path: "config"}:                                                "OpenLDAPConfig",
	{mount: "openldap", path: "creds/{name}"}:                                          "OpenLDAPCredentials",
	{mount: "openldap", path: "library"}:                                               "OpenLDAPLibraries",
	{mount: "openldap", path: "library/manage/{name}/check-in"}:                        "OpenLDAPLibraryManageCheckIn",
	{mount: "openldap", path: "library/{name}"}:                                        "OpenLDAPLibrary",
	{mount: "openldap", path: "library/{name}/check-in"}:                               "OpenLDAPLibraryCheckIn",
	{mount: "openldap", path: "library/{name}/check-out"}:                              "OpenLDAPLibraryCheckOut",
	{mount: "openldap", path: "library/{name}/status"}:                                 "OpenLDAPLibraryStatus",
	{mount: "openldap", path: "role"}:                                                  "OpenLDAPRoles",
	{mount: "openldap", path: "role/{name}"}:                                           "OpenLDAPRole",
	{mount: "openldap", path: "rotate-role/{name}"}:                                    "OpenLDAPRotateRole",
	{mount: "openldap", path: "rotate-root"}:                                           "OpenLDAPRotateRoot",
	{mount: "openldap", path: "static-cred/{name}"}:                                    "OpenLDAPStaticCredentials",
	{mount: "openldap", path: "static-role"}:                                           "OpenLDAPStaticRoles",
	{mount: "openldap", path: "static-role/{name}"}:                                    "OpenLDAPStaticRole",
	{mount: "pki", path: "bundle"}:                                                     "PKIBundle",
	{mount: "pki", path: "ca"}:                                                         "PKICa",
	{mount: "pki", path: "ca/pem"}:                                                     "PKICaPem",
	{mount: "pki", path: "ca_chain"}:                                                   "PKICaChain",
	{mount: "pki", path: "cert"}:                                                       "PKICerts",
	{mount: "pki", path: "cert/ca_chain"}:                                              "PKICertCaChain",
	{mount: "pki", path: "cert/{serial}"}:                                              "PKICert",
	{mount: "pki", path: "cert/{serial}/raw"}:                                          "PKICertRaw",
	{mount: "pki", path: "cert/{serial}/raw/pem"}:                                      "PKICertRawPem",
	{mount: "pki", path: "certs"}:                                                      "PKICerts",
	{mount: "pki", path: "certs/revoked"}:                                              "PKICertsRevoked",
	{mount: "pki", path: "config/auto-tidy"}:                                           "PKIConfigAutoTidy",
	{mount: "pki", path: "config/ca"}:                                                  "PKIConfigCa",
	{mount: "pki", path: "config/cluster"}:                                             "PKIConfigCluster",
	{mount: "pki", path: "config/crl"}:                                                 "PKIConfigCRL",
	{mount: "pki", path: "config/issuers"}:                                             "PKIConfigIssuers",
	{mount: "pki", path: "config/keys"}:                                                "PKIConfigKeys",
	{mount: "pki", path: "config/urls"}:                                                "PKIConfigURLs",
	{mount: "pki", path: "crl"}:                                                        "PKICRL",
	{mount: "pki", path: "crl/rotate"}:                                                 "PKICRLRotate",
	{mount: "pki", path: "crl/rotate-delta"}:                                           "PKICRLRotateDelta",
	{mount: "pki", path: "delta-crl"}:                                                  "PKIDeltaCRL",
	{mount: "pki", path: "intermediate/cross-sign"}:                                    "PKIIntermediateCrossSign",
	{mount: "pki", path: "intermediate/generate/{exported}"}:                           "PKIIntermediateGenerate",
	{mount: "pki", path: "intermediate/set-signed"}:                                    "PKIIntermediateSetSigned",
	{mount: "pki", path: "internal|exported"}:                                          "PKIInternalExported",
	{mount: "pki", path: "issue/{role}"}:                                               "PKIIssueRole",
	{mount: "pki", path: "issuer/{issuer_ref}/issue/{role}"}:                           "PKIIssuerIssueRole",
	{mount: "pki", path: "issuer/{issuer_ref}/resign-crls"}:                            "PKIIssuerResignCRLs",
	{mount: "pki", path: "issuer/{issuer_ref}/revoke"}:                                 "PKIIssuerRevoke",
	{mount: "pki", path: "issuer/{issuer_ref}/sign-intermediate"}:                      "PKIIssuerSignIntermediate",
	{mount: "pki", path: "issuer/{issuer_ref}/sign-revocation-list"}:                   "PKIIssuerSignRevocationList",
	{mount: "pki", path: "issuer/{issuer_ref}/sign-self-issued"}:                       "PKIIssuerSignSelfIssued",
	{mount: "pki", path: "issuer/{issuer_ref}/sign-verbatim"}:                          "PKIIssuerSignVerbatim",
	{mount: "pki", path: "issuer/{issuer_ref}/sign-verbatim/{role}"}:                   "PKIIssuerSignVerbatimRole",
	{mount: "pki", path: "issuer/{issuer_ref}/sign/{role}"}:                            "PKIIssuerSignRole",
	{mount: "pki", path: "issuers"}:                                                    "PKIIssuers",
	{mount: "pki", path: "issuers/generate/intermediate/{exported}"}:                   "PKIIssuersGenerateIntermediateExported",
	{mount: "pki", path: "issuers/generate/root/{exported}"}:                           "PKIIssuersGenerateRootExported",
	{mount: "pki", path: "key/{key_ref}"}:                                              "PKIKey",
	{mount: "pki", path: "keys"}:                                                       "PKIKeys",
	{mount: "pki", path: "keys/import"}:                                                "PKIKeysImport",
	{mount: "pki", path: "kms"}:                                                        "PKIKMS",
	{mount: "pki", path: "ocsp"}:                                                       "PKIOCSP",
	{mount: "pki", path: "ocsp/{req}"}:                                                 "PKIOCSPReq",
	{mount: "pki", path: "revoke"}:                                                     "PKIRevoke",
	{mount: "pki", path: "revoke-with-key"}:                                            "PKIRevokeWithKey",
	{mount: "pki", path: "roles"}:                                                      "PKIRoles",
	{mount: "pki", path: "roles/{name}"}:                                               "PKIRole",
	{mount: "pki", path: "root"}:                                                       "PKIRoot",
	{mount: "pki", path: "root/generate/{exported}"}:                                   "PKIRootGenerate",
	{mount: "pki", path: "root/replace"}:                                               "PKIRootReplace",
	{mount: "pki", path: "root/rotate/{exported}"}:                                     "PKIRootRotate",
	{mount: "pki", path: "root/sign-intermediate"}:                                     "PKIRootSignIntermediate",
	{mount: "pki", path: "root/sign-self-issued"}:                                      "PKIRootSignSelfIssued",
	{mount: "pki", path: "sign-verbatim"}:                                              "PKISignVerbatim",
	{mount: "pki", path: "sign-verbatim/{role}"}:                                       "PKISignVerbatimRole",
	{mount: "pki", path: "sign/{role}"}:                                                "PKISignRole",
	{mount: "pki", path: "tidy"}:                                                       "PKITidy",
	{mount: "pki", path: "tidy-cancel"}:                                                "PKITidyCancel",
	{mount: "pki", path: "tidy-status"}:                                                "PKITidyStatus",
	{mount: "rabbitmq", path: "config/connection"}:                                     "RabbitMQConfigConnection",
	{mount: "rabbitmq", path: "config/lease"}:                                          "RabbitMQConfigLease",
	{mount: "rabbitmq", path: "creds/{name}"}:                                          "RabbitMQCredentials",
	{mount: "rabbitmq", path: "roles"}:                                                 "RabbitMQRoles",
	{mount: "rabbitmq", path: "roles/{name}"}:                                          "RabbitMQRole",
	{mount: "secret", path: "config"}:                                                  "SecretConfig",
	{mount: "secret", path: "data/{path}"}:                                             "Secret",
	{mount: "secret", path: "delete/{path}"}:                                           "SecretVersions",
	{mount: "secret", path: "destroy/{path}"}:                                          "SecretVersions",
	{mount: "secret", path: "metadata/{path}"}:                                         "SecretMetadata",
	{mount: "secret", path: "subkeys/{path}"}:                                          "SecretSubkeys",
	{mount: "secret", path: "undelete/{path}"}:                                         "SecretVersions",
	{mount: "ssh", path: "config/ca"}:                                                  "SSHConfigCA",
	{mount: "ssh", path: "config/zeroaddress"}:                                         "SSHConfigZeroAddress",
	{mount: "ssh", path: "creds/{role}"}:                                               "SSHCredentials",
	{mount: "ssh", path: "issue/{role}"}:                                               "SSHIssue",
	{mount: "ssh", path: "keys/{key_name}"}:                                            "SSHKeys",
	{mount: "ssh", path: "lookup"}:                                                     "SSHLookup",
	{mount: "ssh", path: "public_key"}:                                                 "SSHPublicKey",
	{mount: "ssh", path: "roles"}:                                                      "SSHRoles",
	{mount: "ssh", path: "roles/{role}"}:                                               "SSHRole",
	{mount: "ssh", path: "sign/{role}"}:                                                "SSHSign",
	{mount: "ssh", path: "verify"}:                                                     "SSHVerify",
	{mount: "sys", path: "audit"}:                                                      "AuditDevices",
	{mount: "sys", path: "audit/{path}"}:                                               "AuditDevicesWith",
	{mount: "sys", path: "audit-hash/{path}"}:                                          "AuditHash",
	{mount: "sys", path: "auth"}:                                                       "AuthMethods",
	{mount: "sys", path: "auth/{path}"}:                                                "AuthMethodsWith",
	{mount: "sys", path: "auth/{path}/tune"}:                                           "AuthMethodTune",
	{mount: "sys", path: "capabilities"}:                                               "Capabilities",
	{mount: "sys", path: "capabilities-accessor"}:                                      "CapabilitiesAccessor",
	{mount: "sys", path: "capabilities-self"}:                                          "CapabilitiesSelf",
	{mount: "sys", path: "config/auditing/request-headers"}:                            "ConfigAuditingRequestHeaders",
	{mount: "sys", path: "config/auditing/request-headers/{header}"}:                   "ConfigAuditingRequestHeader",
	{mount: "sys", path: "config/cors"}:                                                "ConfigCORS",
	{mount: "sys", path: "config/reload/{subsystem}"}:                                  "ConfigReloadSubsystem",
	{mount: "sys", path: "config/state/sanitized"}:                                     "ConfigStateSanitized",
	{mount: "sys", path: "config/ui/headers/"}:                                         "ConfigUIHeaders",
	{mount: "sys", path: "config/ui/headers/{header}"}:                                 "ConfigUIHeader",
	{mount: "sys", path: "generate-root"}:                                              "GenerateRoot",
	{mount: "sys", path: "generate-root/attempt"}:                                      "GenerateRootAttempt",
	{mount: "sys", path: "generate-root/update"}:                                       "GenerateRootUpdate",
	{mount: "sys", path: "ha-status"}:                                                  "HAStatus",
	{mount: "sys", path: "health"}:                                                     "Health",
	{mount: "sys", path: "host-info"}:                                                  "HostInfo",
	{mount: "sys", path: "in-flight-req"}:                                              "InFlightRequests",
	{mount: "sys", path: "init"}:                                                       "Init",
	{mount: "sys", path: "internal/counters/activity"}:                                 "InternalCountersActivity",
	{mount: "sys", path: "internal/counters/activity/export"}:                          "InternalCountersActivityExport",
	{mount: "sys", path: "internal/counters/activity/monthly"}:                         "InternalCountersActivityMonthly",
	{mount: "sys", path: "internal/counters/config"}:                                   "InternalCountersConfig",
	{mount: "sys", path: "internal/counters/entities"}:                                 "InternalCountersEntities",
	{mount: "sys", path: "internal/counters/requests"}:                                 "InternalCountersRequests",
	{mount: "sys", path: "internal/counters/tokens"}:                                   "InternalCountersTokens",
	{mount: "sys", path: "internal/inspect/router/{tag}"}:                              "InternalInspectRouter",
	{mount: "sys", path: "internal/specs/openapi"}:                                     "InternalSpecsOpenAPI",
	{mount: "sys", path: "internal/ui/feature-flags"}:                                  "InternalUIFeatureFlags",
	{mount: "sys", path: "internal/ui/mounts"}:                                         "InternalUIMounts",
	{mount: "sys", path: "internal/ui/mounts/{path}"}:                                  "InternalUIMount",
	{mount: "sys", path: "internal/ui/namespaces"}:                                     "InternalUINamespaces",
	{mount: "sys", path: "internal/ui/resultant-acl"}:                                  "InternalUIResultantACL",
	{mount: "sys", path: "key-status"}:                                                 "KeyStatus",
	{mount: "sys", path: "leader"}:                                                     "Leader",
	{mount: "sys", path: "leases"}:                                                     "Leases",
	{mount: "sys", path: "leases/count"}:                                               "LeasesCount",
	{mount: "sys", path: "leases/lookup"}:                                              "LeasesLookup",
	{mount: "sys", path: "leases/lookup/{prefix}"}:                                     "LeasesLookupPrefix",
	{mount: "sys", path: "leases/renew"}:                                               "LeasesRenew",
	{mount: "sys", path: "leases/renew/{url_lease_id}"}:                                "LeasesRenew2",
	{mount: "sys", path: "leases/revoke"}:                                              "LeasesRevoke",
	{mount: "sys", path: "leases/revoke-force/{prefix}"}:                               "LeasesRevokeForce",
	{mount: "sys", path: "leases/revoke-prefix/{prefix}"}:                              "LeasesRevokePrefix",
	{mount: "sys", path: "leases/revoke/{url_lease_id}"}:                               "LeasesRevoke2",
	{mount: "sys", path: "leases/tidy"}:                                                "LeasesTidy",
	{mount: "sys", path: "loggers"}:                                                    "Loggers",
	{mount: "sys", path: "loggers/{name}"}:                                             "Logger",
	{mount: "sys", path: "metrics"}:                                                    "Metrics",
	{mount: "sys", path: "mfa/validate"}:                                               "MFAValidate",
	{mount: "sys", path: "monitor"}:                                                    "Monitor",
	{mount: "sys", path: "mounts"}:                                                     "Mounts",
	{mount: "sys", path: "mounts/{path}"}:                                              "MountsWith",
	{mount: "sys", path: "mounts/{path}/tune"}:                                         "MountsTune",
	{mount: "sys", path: "plugins/catalog"}:                                            "PluginsCatalog",
	{mount: "sys", path: "plugins/catalog/{type}"}:                                     "PluginsCatalogByType",
	{mount: "sys", path: "plugins/catalog/{type}/{name}"}:                              "PluginsCatalogByTypeByName",
	{mount: "sys", path: "plugins/reload/backend"}:                                     "PluginsReloadBackend",
	{mount: "sys", path: "policies/acl"}:                                               "PoliciesACL",
	{mount: "sys", path: "policies/acl/{name}"}:                                        "PoliciesACL",
	{mount: "sys", path: "policies/password"}:                                          "PoliciesPassword",
	{mount: "sys", path: "policies/password/{name}"}:                                   "PoliciesPassword",
	{mount: "sys", path: "policies/password/{name}/generate"}:                          "PoliciesPasswordGenerate",
	{mount: "sys", path: "policy"}:                                                     "Policies",
	{mount: "sys", path: "policy/{name}"}:                                              "Policy",
	{mount: "sys", path: "pprof/"}:                                                     "Pprof",
	{mount: "sys", path: "pprof/allocs"}:                                               "PprofAllocs",
	{mount: "sys", path: "pprof/block"}:                                                "PprofBlock",
	{mount: "sys", path: "pprof/cmdline"}:                                              "PprofCmdline",
	{mount: "sys", path: "pprof/goroutine"}:                                            "PprofGoroutine",
	{mount: "sys", path: "pprof/heap"}:                                                 "PprofHeap",
	{mount: "sys", path: "pprof/mutex"}:                                                "PprofMutex",
	{mount: "sys", path: "pprof/profile"}:                                              "PprofProfile",
	{mount: "sys", path: "pprof/symbol"}:                                               "PprofSymbol",
	{mount: "sys", path: "pprof/threadcreate"}:                                         "PprofThreadcreate",
	{mount: "sys", path: "pprof/trace"}:                                                "PprofTrace",
	{mount: "sys", path: "quotas/config"}:                                              "QuotasConfig",
	{mount: "sys", path: "quotas/rate-limit"}:                                          "QuotasRateLimits",
	{mount: "sys", path: "quotas/rate-limit/{name}"}:                                   "QuotasRateLimit",
	{mount: "sys", path: "raw"}:                                                        "Raw",
	{mount: "sys", path: "raw/{path}"}:                                                 "RawPath",
	{mount: "sys", path: "rekey/backup"}:                                               "RekeyBackup",
	{mount: "sys", path: "rekey/init"}:                                                 "RekeyInit",
	{mount: "sys", path: "rekey/recovery-key-backup"}:                                  "RekeyRecoveryKeyBackup",
	{mount: "sys", path: "rekey/update"}:                                               "RekeyUpdate",
	{mount: "sys", path: "rekey/verify"}:                                               "RekeyVerify",
	{mount: "sys", path: "remount"}:                                                    "Remount",
	{mount: "sys", path: "remount/status/{migration_id}"}:                              "RemountStatus",
	{mount: "sys", path: "renew"}:                                                      "Renew",
	{mount: "sys", path: "renew/{url_lease_id}"}:                                       "Renew",
	{mount: "sys", path: "replication/status"}:                                         "ReplicationStatus",
	{mount: "sys", path: "revoke"}:                                                     "Revoke",
	{mount: "sys", path: "revoke-force/{prefix}"}:                                      "RevokeForce",
	{mount: "sys", path: "revoke-prefix/{prefix}"}:                                     "RevokePrefix",
	{mount: "sys", path: "revoke/{url_lease_id}"}:                                      "Revoke",
	{mount: "sys", path: "rotate"}:                                                     "Rotate",
	{mount: "sys", path: "rotate/config"}:                                              "RotateConfig",
	{mount: "sys", path: "seal"}:                                                       "Seal",
	{mount: "sys", path: "seal-status"}:                                                "SealStatus",
	{mount: "sys", path: "step-down"}:                                                  "StepDown",
	{mount: "sys", path: "tools/hash"}:                                                 "ToolsHashes",
	{mount: "sys", path: "tools/hash/{urlalgorithm}"}:                                  "ToolsHash",
	{mount: "sys", path: "tools/random"}:                                               "ToolsRandom",
	{mount: "sys", path: "tools/random/{source}"}:                                      "ToolsRandomSource",
	{mount: "sys", path: "tools/random/{source}/{urlbytes}"}:                           "ToolsRandomSourceBytes",
	{mount: "sys", path: "unseal"}:                                                     "Unseal",
	{mount: "sys", path: "version-history/"}:                                           "VersionHistory",
	{mount: "sys", path: "wrapping/lookup"}:                                            "WrappingLookup",
	{mount: "sys", path: "wrapping/rewrap"}:                                            "WrappingRewrap",
	{mount: "sys", path: "wrapping/unwrap"}:                                            "WrappingUnwrap",
	{mount: "sys", path: "wrapping/wrap"}:                                              "WrappingWrap",
	{mount: "terraform", path: "config"}:                                               "TerraformConfig",
	{mount: "terraform", path: "creds/{name}"}:                                         "TerraformCredentials",
	{mount: "terraform", path: "role"}:                                                 "TerraformRoles",
	{mount: "terraform", path: "role/{name}"}:                                          "TerraformRole",
	{mount: "terraform", path: "rotate-role/{name}"}:                                   "TerraformRotateRole",
	{mount: "totp", path: "code/{name}"}:                                               "TOTPCode",
	{mount: "totp", path: "keys"}:                                                      "TOTPKeys",
	{mount: "totp", path: "keys/{name}"}:                                               "TOTPKey",
	{mount: "transit", path: "backup/{name}"}:                                          "TransitBackup",
	{mount: "transit", path: "cache-config"}:                                           "TransitCacheConfig",
	{mount: "transit", path: "datakey/{plaintext}/{name}"}:                             "TransitDatakey",
	{mount: "transit", path: "decrypt/{name}"}:                                         "TransitDecrypt",
	{mount: "transit", path: "encrypt/{name}"}:                                         "TransitEncrypt",
	{mount: "transit", path: "export/{type}/{name}"}:                                   "TransitExport",
	{mount: "transit", path: "export/{type}/{name}/{version}"}:                         "TransitExportVersion",
	{mount: "transit", path: "hash"}:                                                   "TransitHash",
	{mount: "transit", path: "hash/{urlalgorithm}"}:                                    "TransitHash",
	{mount: "transit", path: "hmac/{name}"}:                                            "TransitHMAC",
	{mount: "transit", path: "hmac/{name}/{urlalgorithm}"}:                             "TransitHMACAlgorithm",
	{mount: "transit", path: "keys"}:                                                   "TransitKeys",
	{mount: "transit", path: "keys/{name}"}:                                            "TransitKey",
	{mount: "transit", path: "keys/{name}/config"}:                                     "TransitKeyConfig",
	{mount: "transit", path: "keys/{name}/import"}:                                     "TransitKeyImport",
	{mount: "transit", path: "keys/{name}/import_version"}:                             "TransitKeyImportVersion",
	{mount: "transit", path: "keys/{name}/rotate"}:                                     "TransitKeyRotate",
	{mount: "transit", path: "keys/{name}/trim"}:                                       "TransitKeyTrim",
	{mount: "transit", path: "random"}:                                                 "TransitRandom",
	{mount: "transit", path: "random/{source}"}:                                        "TransitRandomSource",
	{mount: "transit", path: "random/{source}/{urlbytes}"}:                             "TransitRandomSourceBytes",
	{mount: "transit", path: "restore"}:                                                "TransitRestore",
	{mount: "transit", path: "restore/{name}"}:                                         "TransitRestoreKey",
	{mount: "transit", path: "rewrap/{name}"}:                                          "TransitRewrap",
	{mount: "transit", path: "sign/{name}"}:                                            "TransitSign",
	{mount: "transit", path: "sign/{name}/{urlalgorithm}"}:                             "TransitSignAlgorithm",
	{mount: "transit", path: "verify/{name}"}:                                          "TransitVerify",
	{mount: "transit", path: "verify/{name}/{urlalgorithm}"}:                           "TransitVerifyAlgorithm",
	{mount: "transit", path: "wrapping_key"}:                                           "TransitWrappingKey",
}

// constructRequestResponseIdentifier joins the given inputs into a title case
// string, e.g. 'UpdateNomadConfigLeaseRequest'. This function is used to
// generate:
//
//   - OpenAPI operation ID
//   - OpenAPI request names
//   - OpenAPI response names
//
// For certain prefix + path combinations, which would otherwise result in an
// ugly string, the function uses a custom lookup table to construct part of
// the string instead.
func constructRequestResponseIdentifier(operation logical.Operation, mount, path, suffix string) string {
	operationStr := string(operation)

	// Replace the operation prefix (usually update/POST) with the actual operation
	if mount == "secret" {
		for _, prefix := range []string{
			"delete",
			"destroy",
			"undelete",
		} {
			if strings.HasPrefix(path, prefix) {
				operationStr = prefix
			}
		}
	}

	// Remove the operation prefix (usually update/POST) from the login requests
	if strings.HasPrefix(mount, "auth/") && strings.HasPrefix(path, "login") {
		operationStr = ""
	}

	// Split the operation by non-word characters (if any)
	parts := nonWordRe.Split(strings.ToLower(operationStr), -1)

	// Append either the known mapping or mount + path split by non-word characters
	if mapping, ok := knownPathMappings[knownPathKey{mount: mount, path: path}]; ok {
		parts = append(parts, mapping)
	} else {
		parts = append(parts, nonWordRe.Split(strings.ToLower(mount), -1)...)
		parts = append(parts, nonWordRe.Split(strings.ToLower(path), -1)...)
	}

	parts = append(parts, suffix)

	// Title case everything & join the result into a string
	title := cases.Title(language.English, cases.NoLower)

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
