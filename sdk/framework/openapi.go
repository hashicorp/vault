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
func documentPaths(backend *Backend, defaultMountPath string, doc *OASDocument) error {
	for _, p := range backend.Paths {
		if err := documentPath(p, backend.SpecialPaths(), defaultMountPath, backend.BackendType, doc); err != nil {
			return err
		}
	}

	return nil
}

// documentPath parses a framework.Path into one or more OpenAPI paths.
func documentPath(p *Path, specialPaths *logical.Paths, defaultMountPath string, backendType logical.BackendType, doc *OASDocument) error {
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

		if defaultMountPath != "sys" && defaultMountPath != "identity" {
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
			op.OperationID = constructRequestResponseIdentifier(opType, defaultMountPath, path, "")

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
					requestName := constructRequestResponseIdentifier(opType, defaultMountPath, path, "request")
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
						responseName := constructRequestResponseIdentifier(opType, defaultMountPath, path, "response")
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
// we use custom path mappings instead to generate request/response names.
type knownPathKey struct {
	mount string
	path  string
}

type knownPathValue struct {
	name             string
	operationImplied bool // if true, operation will not be appended to the resulting name
}

var knownPathMappings = map[knownPathKey]knownPathValue{
	{mount: "auth/alicloud", path: "login"}:                                            {name: "AliCloudLogin", operationImplied: true},
	{mount: "auth/alicloud", path: "role"}:                                             {name: "AliCloudAuthRoles"},
	{mount: "auth/alicloud", path: "role/{role}"}:                                      {name: "AliCloudAuthRole"},
	{mount: "auth/alicloud", path: "roles"}:                                            {name: "AliCloudAuthRoles2"},
	{mount: "auth/approle", path: "login"}:                                             {name: "AppRoleLogin", operationImplied: true},
	{mount: "auth/approle", path: "role"}:                                              {name: "AppRoleRoles"},
	{mount: "auth/approle", path: "role/{role_name}"}:                                  {name: "AppRoleRole"},
	{mount: "auth/approle", path: "role/{role_name}/bind-secret-id"}:                   {name: "AppRoleBindSecretID"},
	{mount: "auth/approle", path: "role/{role_name}/bound-cidr-list"}:                  {name: "AppRoleBoundCIDRList"},
	{mount: "auth/approle", path: "role/{role_name}/custom-secret-id"}:                 {name: "AppRoleCustomSecretID"},
	{mount: "auth/approle", path: "role/{role_name}/local-secret-ids"}:                 {name: "AppRoleLocalSecretIDs"},
	{mount: "auth/approle", path: "role/{role_name}/period"}:                           {name: "AppRolePeriod"},
	{mount: "auth/approle", path: "role/{role_name}/policies"}:                         {name: "AppRolePolicies"},
	{mount: "auth/approle", path: "role/{role_name}/role-id"}:                          {name: "AppRoleRoleID"},
	{mount: "auth/approle", path: "role/{role_name}/secret-id"}:                        {name: "AppRoleSecretID"},
	{mount: "auth/approle", path: "role/{role_name}/secret-id-accessor/destroy"}:       {name: "AppRoleSecretIDAccessorDestroy"},
	{mount: "auth/approle", path: "role/{role_name}/secret-id-accessor/lookup"}:        {name: "AppRoleSecretIDAccessorLookup"},
	{mount: "auth/approle", path: "role/{role_name}/secret-id-bound-cidrs"}:            {name: "AppRoleSecretIDBoundCIDRs"},
	{mount: "auth/approle", path: "role/{role_name}/secret-id-num-uses"}:               {name: "AppRoleSecretIDNumberOfUses"},
	{mount: "auth/approle", path: "role/{role_name}/secret-id-ttl"}:                    {name: "AppRoleSecretIDTTL"},
	{mount: "auth/approle", path: "role/{role_name}/secret-id/destroy"}:                {name: "AppRoleSecretIDDestroy"},
	{mount: "auth/approle", path: "role/{role_name}/secret-id/lookup"}:                 {name: "AppRoleSecretIDLookup"},
	{mount: "auth/approle", path: "role/{role_name}/token-bound-cidrs"}:                {name: "AppRoleTokenBoundCIDRs"},
	{mount: "auth/approle", path: "role/{role_name}/token-max-ttl"}:                    {name: "AppRoleTokenMaxTTL"},
	{mount: "auth/approle", path: "role/{role_name}/token-num-uses"}:                   {name: "AppRoleTokenNumberOfUses"},
	{mount: "auth/approle", path: "role/{role_name}/token-ttl"}:                        {name: "AppRoleTokenTTL"},
	{mount: "auth/approle", path: "tidy/secret-id"}:                                    {name: "AppRoleTidySecretID"},
	{mount: "auth/aws", path: "config/certificate/{cert_name}"}:                        {name: "AWSConfigCertificate"},
	{mount: "auth/aws", path: "config/certificates"}:                                   {name: "AWSConfigCertificates"},
	{mount: "auth/aws", path: "config/client"}:                                         {name: "AWSConfigClient"},
	{mount: "auth/aws", path: "config/identity"}:                                       {name: "AWSConfigIdentity"},
	{mount: "auth/aws", path: "config/rotate-root"}:                                    {name: "AWSConfigRotateRoot"},
	{mount: "auth/aws", path: "config/sts"}:                                            {name: "AWSConfigSecurityTokenService"},
	{mount: "auth/aws", path: "config/sts/{account_id}"}:                               {name: "AWSConfigSecurityTokenServiceAccount"},
	{mount: "auth/aws", path: "config/tidy/identity-accesslist"}:                       {name: "AWSConfigTidyIdentityAccesslist"},
	{mount: "auth/aws", path: "config/tidy/identity-whitelist"}:                        {name: "AWSConfigTidyIdentityWhitelist"},
	{mount: "auth/aws", path: "config/tidy/roletag-blacklist"}:                         {name: "AWSConfigTidyRoleTagBlacklist"},
	{mount: "auth/aws", path: "config/tidy/roletag-denylist"}:                          {name: "AWSConfigTidyRoleTagDenylist"},
	{mount: "auth/aws", path: "identity-accesslist"}:                                   {name: "AWSIdentityAccesslist"},
	{mount: "auth/aws", path: "identity-accesslist/{instance_id}"}:                     {name: "AWSIdentityAccesslistFor"},
	{mount: "auth/aws", path: "identity-whitelist"}:                                    {name: "AWSIdentityWhitelist"},
	{mount: "auth/aws", path: "identity-whitelist/{instance_id}"}:                      {name: "AWSIdentityWhitelistFor"},
	{mount: "auth/aws", path: "login"}:                                                 {name: "AWSLogin", operationImplied: true},
	{mount: "auth/aws", path: "role"}:                                                  {name: "AWSAuthRoles"},
	{mount: "auth/aws", path: "role/{role}"}:                                           {name: "AWSAuthRole"},
	{mount: "auth/aws", path: "role/{role}/tag"}:                                       {name: "AWSAuthRoleTag"},
	{mount: "auth/aws", path: "roles"}:                                                 {name: "AWSAuthRoles2"},
	{mount: "auth/aws", path: "roletag-blacklist"}:                                     {name: "AWSRoleTagBlacklist"},
	{mount: "auth/aws", path: "roletag-blacklist/{role_tag}"}:                          {name: "AWSRoleTagBlacklistFor"},
	{mount: "auth/aws", path: "roletag-denylist"}:                                      {name: "AWSRoleTagDenylist"},
	{mount: "auth/aws", path: "roletag-denylist/{role_tag}"}:                           {name: "AWSRoleTagDenylistFor"},
	{mount: "auth/aws", path: "tidy/identity-accesslist"}:                              {name: "AWSTidyIdentityAccesslist"},
	{mount: "auth/aws", path: "tidy/identity-whitelist"}:                               {name: "AWSTidyIdentityWhitelist"},
	{mount: "auth/aws", path: "tidy/roletag-blacklist"}:                                {name: "AWSTidyRoleTagBlacklist"},
	{mount: "auth/aws", path: "tidy/roletag-denylist"}:                                 {name: "AWSTidyRoleTagDenylist"},
	{mount: "auth/azure", path: "config"}:                                              {name: "AzureConfig"},
	{mount: "auth/azure", path: "login"}:                                               {name: "AzureLogin", operationImplied: true},
	{mount: "auth/azure", path: "role"}:                                                {name: "AzureRoles"},
	{mount: "auth/azure", path: "role/{name}"}:                                         {name: "AzureRole"},
	{mount: "auth/centrify", path: "config"}:                                           {name: "CentrifyConfig"},
	{mount: "auth/centrify", path: "login"}:                                            {name: "CentrifyLogin", operationImplied: true},
	{mount: "auth/cert", path: "certs"}:                                                {name: "Certificates"},
	{mount: "auth/cert", path: "certs/{name}"}:                                         {name: "Certificate"},
	{mount: "auth/cert", path: "config"}:                                               {name: "CertificateConfig"},
	{mount: "auth/cert", path: "crls"}:                                                 {name: "CertificateCRLs"},
	{mount: "auth/cert", path: "crls/{name}"}:                                          {name: "CertificateCRL"},
	{mount: "auth/cert", path: "login"}:                                                {name: "CertificateLogin", operationImplied: true},
	{mount: "auth/cf", path: "config"}:                                                 {name: "CloudFoundryConfig"},
	{mount: "auth/cf", path: "login"}:                                                  {name: "CloudFoundryLogin", operationImplied: true},
	{mount: "auth/cf", path: "roles"}:                                                  {name: "CloudFoundryRoles"},
	{mount: "auth/cf", path: "roles/{role}"}:                                           {name: "CloudFoundryRole"},
	{mount: "auth/gcp", path: "config"}:                                                {name: "GoogleCloudConfig"},
	{mount: "auth/gcp", path: "login"}:                                                 {name: "GoogleCloudLogin", operationImplied: true},
	{mount: "auth/gcp", path: "role"}:                                                  {name: "GoogleCloudRoles"},
	{mount: "auth/gcp", path: "role/{name}"}:                                           {name: "GoogleCloudRole"},
	{mount: "auth/gcp", path: "role/{name}/labels"}:                                    {name: "GoogleCloudRoleLabels"},
	{mount: "auth/gcp", path: "role/{name}/service-accounts"}:                          {name: "GoogleCloudRoleServiceAccounts"},
	{mount: "auth/gcp", path: "roles"}:                                                 {name: "GoogleCloudRoles2"},
	{mount: "auth/github", path: "config"}:                                             {name: "GitHubConfig"},
	{mount: "auth/github", path: "login"}:                                              {name: "GitHubLogin", operationImplied: true},
	{mount: "auth/github", path: "map/teams"}:                                          {name: "GitHubMapTeams"},
	{mount: "auth/github", path: "map/teams/{key}"}:                                    {name: "GitHubMapTeam"},
	{mount: "auth/github", path: "map/users"}:                                          {name: "GitHubMapUsers"},
	{mount: "auth/github", path: "map/users/{key}"}:                                    {name: "GitHubMapUser"},
	{mount: "auth/jwt", path: "config"}:                                                {name: "JWTConfig"},
	{mount: "auth/jwt", path: "login"}:                                                 {name: "JWTLogin", operationImplied: true},
	{mount: "auth/jwt", path: "oidc/auth_url"}:                                         {name: "JWTOIDCAuthURL"},
	{mount: "auth/jwt", path: "oidc/callback"}:                                         {name: "JWTOIDCCallback"},
	{mount: "auth/jwt", path: "role"}:                                                  {name: "JWTRoles"},
	{mount: "auth/jwt", path: "role/{name}"}:                                           {name: "JWTRole"},
	{mount: "auth/kerberos", path: "config"}:                                           {name: "KerberosConfig"},
	{mount: "auth/kerberos", path: "config/ldap"}:                                      {name: "KerberosConfigLDAP"},
	{mount: "auth/kerberos", path: "groups"}:                                           {name: "KerberosGroups"},
	{mount: "auth/kerberos", path: "groups/{name}"}:                                    {name: "KerberosGroup"},
	{mount: "auth/kerberos", path: "login"}:                                            {name: "KerberosLogin", operationImplied: true},
	{mount: "auth/kubernetes", path: "config"}:                                         {name: "KubernetesConfig"},
	{mount: "auth/kubernetes", path: "login"}:                                          {name: "KubernetesLogin", operationImplied: true},
	{mount: "auth/kubernetes", path: "role"}:                                           {name: "KubernetesRoles"},
	{mount: "auth/kubernetes", path: "role/{name}"}:                                    {name: "KubernetesRole"},
	{mount: "auth/ldap", path: "config"}:                                               {name: "LDAPConfig"},
	{mount: "auth/ldap", path: "groups"}:                                               {name: "LDAPGroups"},
	{mount: "auth/ldap", path: "groups/{name}"}:                                        {name: "LDAPGroup"},
	{mount: "auth/ldap", path: "login/{username}"}:                                     {name: "LDAPLogin", operationImplied: true},
	{mount: "auth/ldap", path: "users"}:                                                {name: "LDAPUsers"},
	{mount: "auth/ldap", path: "users/{name}"}:                                         {name: "LDAPUser"},
	{mount: "auth/oci", path: "config"}:                                                {name: "OCIConfig"},
	{mount: "auth/oci", path: "login"}:                                                 {name: "OCILogin", operationImplied: true},
	{mount: "auth/oci", path: "login/{role}"}:                                          {name: "OCILoginWithRole", operationImplied: true},
	{mount: "auth/oci", path: "role"}:                                                  {name: "OCIRoles"},
	{mount: "auth/oci", path: "role/{role}"}:                                           {name: "OCIRole"},
	{mount: "auth/oidc", path: "config"}:                                               {name: "OIDCConfig"},
	{mount: "auth/oidc", path: "login"}:                                                {name: "OIDCLogin", operationImplied: true},
	{mount: "auth/oidc", path: "oidc/auth_url"}:                                        {name: "OIDCAuthURL"},
	{mount: "auth/oidc", path: "oidc/callback"}:                                        {name: "OIDCCallback"},
	{mount: "auth/oidc", path: "role"}:                                                 {name: "OIDCRoles"},
	{mount: "auth/oidc", path: "role/{name}"}:                                          {name: "OIDCRole"},
	{mount: "auth/okta", path: "config"}:                                               {name: "OktaConfig"},
	{mount: "auth/okta", path: "groups"}:                                               {name: "OktaGroups"},
	{mount: "auth/okta", path: "groups/{name}"}:                                        {name: "OktaGroup"},
	{mount: "auth/okta", path: "login/{username}"}:                                     {name: "OktaLogin", operationImplied: true},
	{mount: "auth/okta", path: "users"}:                                                {name: "OktaUsers"},
	{mount: "auth/okta", path: "users/{name}"}:                                         {name: "OktaUser"},
	{mount: "auth/okta", path: "verify/{nonce}"}:                                       {name: "OktaVerify"},
	{mount: "auth/radius", path: "config"}:                                             {name: "RadiusConfig"},
	{mount: "auth/radius", path: "login"}:                                              {name: "RadiusLogin", operationImplied: true},
	{mount: "auth/radius", path: "login/{urlusername}"}:                                {name: "RadiusLoginWithUsername", operationImplied: true},
	{mount: "auth/radius", path: "users"}:                                              {name: "RadiusUsers"},
	{mount: "auth/radius", path: "users/{name}"}:                                       {name: "RadiusUser"},
	{mount: "auth/token", path: "accessors/"}:                                          {name: "TokenAccessors"},
	{mount: "auth/token", path: "create"}:                                              {name: "TokenCreate"},
	{mount: "auth/token", path: "create-orphan"}:                                       {name: "TokenCreateOrphan"},
	{mount: "auth/token", path: "create/{role_name}"}:                                  {name: "TokenCreateWithRole"},
	{mount: "auth/token", path: "lookup"}:                                              {name: "TokenLookup"},
	{mount: "auth/token", path: "lookup-accessor"}:                                     {name: "TokenLookupAccessor"},
	{mount: "auth/token", path: "lookup-self"}:                                         {name: "TokenLookupSelf"},
	{mount: "auth/token", path: "renew"}:                                               {name: "TokenRenew"},
	{mount: "auth/token", path: "renew-accessor"}:                                      {name: "TokenRenewAccessor"},
	{mount: "auth/token", path: "renew-self"}:                                          {name: "TokenRenewSelf"},
	{mount: "auth/token", path: "revoke"}:                                              {name: "TokenRevoke"},
	{mount: "auth/token", path: "revoke-accessor"}:                                     {name: "TokenRevokeAccessor"},
	{mount: "auth/token", path: "revoke-orphan"}:                                       {name: "TokenRevokeOrphan"},
	{mount: "auth/token", path: "revoke-self"}:                                         {name: "TokenRevokeSelf"},
	{mount: "auth/token", path: "roles"}:                                               {name: "TokenRoles"},
	{mount: "auth/token", path: "roles/{role_name}"}:                                   {name: "TokenRole"},
	{mount: "auth/token", path: "tidy"}:                                                {name: "TokenTidy"},
	{mount: "auth/userpass", path: "login/{username}"}:                                 {name: "UserpassLogin", operationImplied: true},
	{mount: "auth/userpass", path: "users"}:                                            {name: "UserpassUsers"},
	{mount: "auth/userpass", path: "users/{username}"}:                                 {name: "UserpassUser"},
	{mount: "auth/userpass", path: "users/{username}/password"}:                        {name: "UserpassUserPassword"},
	{mount: "auth/userpass", path: "users/{username}/policies"}:                        {name: "UserpassUserPolicies"},
	{mount: "ad", path: "config"}:                                                      {name: "ActiveDirectoryConfig"},
	{mount: "ad", path: "creds/{name}"}:                                                {name: "ActiveDirectoryCredentials"},
	{mount: "ad", path: "library"}:                                                     {name: "ActiveDirectoryLibraries"},
	{mount: "ad", path: "library/manage/{name}/check-in"}:                              {name: "ActiveDirectoryLibraryManageCheckIn"},
	{mount: "ad", path: "library/{name}"}:                                              {name: "ActiveDirectoryLibrary"},
	{mount: "ad", path: "library/{name}/check-in"}:                                     {name: "ActiveDirectoryLibraryCheckIn"},
	{mount: "ad", path: "library/{name}/check-out"}:                                    {name: "ActiveDirectoryLibraryCheckOut"},
	{mount: "ad", path: "library/{name}/status"}:                                       {name: "ActiveDirectoryLibraryStatus"},
	{mount: "ad", path: "roles"}:                                                       {name: "ActiveDirectoryRoles"},
	{mount: "ad", path: "roles/{name}"}:                                                {name: "ActiveDirectoryRole"},
	{mount: "ad", path: "rotate-role/{name}"}:                                          {name: "ActiveDirectoryRotateRole"},
	{mount: "ad", path: "rotate-root"}:                                                 {name: "ActiveDirectoryRotateRoot"},
	{mount: "alicloud", path: "config"}:                                                {name: "AliCloudConfig"},
	{mount: "alicloud", path: "creds/{name}"}:                                          {name: "AliCloudCredentials"},
	{mount: "alicloud", path: "role"}:                                                  {name: "AliCloudRoles"},
	{mount: "alicloud", path: "role/{name}"}:                                           {name: "AliCloudRole"},
	{mount: "aws", path: "config/lease"}:                                               {name: "AWSConfigLease"},
	{mount: "aws", path: "config/root"}:                                                {name: "AWSConfigRoot"},
	{mount: "aws", path: "config/rotate-root"}:                                         {name: "AWSConfigRotateRoot"},
	{mount: "aws", path: "creds"}:                                                      {name: "AWSCredentials"},
	{mount: "aws", path: "roles"}:                                                      {name: "AWSRoles"},
	{mount: "aws", path: "roles/{name}"}:                                               {name: "AWSRole"},
	{mount: "aws", path: "sts/{name}"}:                                                 {name: "AWSSecurityTokenService"},
	{mount: "azure", path: "config"}:                                                   {name: "AzureConfig"},
	{mount: "azure", path: "creds/{role}"}:                                             {name: "AzureCredentials"},
	{mount: "azure", path: "roles"}:                                                    {name: "AzureRoles"},
	{mount: "azure", path: "roles/{name}"}:                                             {name: "AzureRole"},
	{mount: "azure", path: "rotate-root"}:                                              {name: "AzureRotateRoot"},
	{mount: "consul", path: "config/access"}:                                           {name: "ConsulConfigAccess"},
	{mount: "consul", path: "creds/{role}"}:                                            {name: "ConsulCredentials"},
	{mount: "consul", path: "roles"}:                                                   {name: "ConsulRoles"},
	{mount: "consul", path: "roles/{name}"}:                                            {name: "ConsulRole"},
	{mount: "cubbyhole", path: "{path}"}:                                               {name: "Cubbyhole"},
	{mount: "gcp", path: "config"}:                                                     {name: "GoogleCloudConfig"},
	{mount: "gcp", path: "config/rotate-root"}:                                         {name: "GoogleCloudConfigRotateRoot"},
	{mount: "gcp", path: "key/{roleset}"}:                                              {name: "GoogleCloudKey"},
	{mount: "gcp", path: "roleset/{name}"}:                                             {name: "GoogleCloudRoleset"},
	{mount: "gcp", path: "roleset/{name}/rotate"}:                                      {name: "GoogleCloudRolesetRotate"},
	{mount: "gcp", path: "roleset/{name}/rotate-key"}:                                  {name: "GoogleCloudRolesetRotateKey"},
	{mount: "gcp", path: "roleset/{roleset}/key"}:                                      {name: "GoogleCloudRolesetKey"},
	{mount: "gcp", path: "roleset/{roleset}/token"}:                                    {name: "GoogleCloudRolesetToken"},
	{mount: "gcp", path: "rolesets"}:                                                   {name: "GoogleCloudRolesets"},
	{mount: "gcp", path: "static-account/{name}"}:                                      {name: "GoogleCloudStaticAccount"},
	{mount: "gcp", path: "static-account/{name}/key"}:                                  {name: "GoogleCloudStaticAccountKey"},
	{mount: "gcp", path: "static-account/{name}/rotate-key"}:                           {name: "GoogleCloudStaticAccountRotateKey"},
	{mount: "gcp", path: "static-account/{name}/token"}:                                {name: "GoogleCloudStaticAccountToken"},
	{mount: "gcp", path: "static-accounts"}:                                            {name: "GoogleCloudStaticAccounts"},
	{mount: "gcp", path: "token/{roleset}"}:                                            {name: "GoogleCloudToken"},
	{mount: "gcpkms", path: "config"}:                                                  {name: "GoogleCloudKMSConfig"},
	{mount: "gcpkms", path: "decrypt/{key}"}:                                           {name: "GoogleCloudKMSDecrypt"},
	{mount: "gcpkms", path: "encrypt/{key}"}:                                           {name: "GoogleCloudKMSEncrypt"},
	{mount: "gcpkms", path: "keys"}:                                                    {name: "GoogleCloudKMSKeys"},
	{mount: "gcpkms", path: "keys/config/{key}"}:                                       {name: "GoogleCloudKMSKeysConfig"},
	{mount: "gcpkms", path: "keys/deregister/{key}"}:                                   {name: "GoogleCloudKMSKeysDeregister"},
	{mount: "gcpkms", path: "keys/register/{key}"}:                                     {name: "GoogleCloudKMSKeysRegister"},
	{mount: "gcpkms", path: "keys/rotate/{key}"}:                                       {name: "GoogleCloudKMSKeysRotate"},
	{mount: "gcpkms", path: "keys/trim/{key}"}:                                         {name: "GoogleCloudKMSKeysTrim"},
	{mount: "gcpkms", path: "keys/{key}"}:                                              {name: "GoogleCloudKMSKey"},
	{mount: "gcpkms", path: "pubkey/{key}"}:                                            {name: "GoogleCloudKMSPubkey"},
	{mount: "gcpkms", path: "reencrypt/{key}"}:                                         {name: "GoogleCloudKMSReencrypt"},
	{mount: "gcpkms", path: "sign/{key}"}:                                              {name: "GoogleCloudKMSSign"},
	{mount: "gcpkms", path: "verify/{key}"}:                                            {name: "GoogleCloudKMSVerify"},
	{mount: "identity", path: "alias"}:                                                 {name: "IdentityAlias"},
	{mount: "identity", path: "alias/id"}:                                              {name: "IdentityAliasesByID"},
	{mount: "identity", path: "alias/id/{id}"}:                                         {name: "IdentityAliasByID"},
	{mount: "identity", path: "entity"}:                                                {name: "IdentityEntity"},
	{mount: "identity", path: "entity-alias"}:                                          {name: "IdentityEntityAlias"},
	{mount: "identity", path: "entity-alias/id"}:                                       {name: "IdentityEntityAliasesByID"},
	{mount: "identity", path: "entity-alias/id/{id}"}:                                  {name: "IdentityEntityAliasByID"},
	{mount: "identity", path: "entity/batch-delete"}:                                   {name: "IdentityEntityBatchDelete"},
	{mount: "identity", path: "entity/id"}:                                             {name: "IdentityEntitiesByID"},
	{mount: "identity", path: "entity/id/{id}"}:                                        {name: "IdentityEntityByID"},
	{mount: "identity", path: "entity/merge"}:                                          {name: "IdentityEntityMerge"},
	{mount: "identity", path: "entity/name"}:                                           {name: "IdentityEntitiesByName"},
	{mount: "identity", path: "entity/name/{name}"}:                                    {name: "IdentityEntityByName"},
	{mount: "identity", path: "group"}:                                                 {name: "IdentityGroup"},
	{mount: "identity", path: "group-alias"}:                                           {name: "IdentityGroupAlias"},
	{mount: "identity", path: "group-alias/id"}:                                        {name: "IdentityGroupAliasesByID"},
	{mount: "identity", path: "group-alias/id/{id}"}:                                   {name: "IdentityGroupAliasByID"},
	{mount: "identity", path: "group/id"}:                                              {name: "IdentityGroupsByID"},
	{mount: "identity", path: "group/id/{id}"}:                                         {name: "IdentityGroupByID"},
	{mount: "identity", path: "group/name"}:                                            {name: "IdentityGroupsByName"},
	{mount: "identity", path: "group/name/{name}"}:                                     {name: "IdentityGroupByName"},
	{mount: "identity", path: "lookup/entity"}:                                         {name: "IdentityLookupEntity"},
	{mount: "identity", path: "lookup/group"}:                                          {name: "IdentityLookupGroup"},
	{mount: "identity", path: "mfa/login-enforcement"}:                                 {name: "IdentityMFALoginEnforcements"},
	{mount: "identity", path: "mfa/login-enforcement/{name}"}:                          {name: "IdentityMFALoginEnforcement"},
	{mount: "identity", path: "mfa/method"}:                                            {name: "IdentityMFAMethods"},
	{mount: "identity", path: "mfa/method/duo"}:                                        {name: "IdentityMFAMethodsDuo"},
	{mount: "identity", path: "mfa/method/duo/{method_id}"}:                            {name: "IdentityMFAMethodDuo"},
	{mount: "identity", path: "mfa/method/okta"}:                                       {name: "IdentityMFAMethodsOkta"},
	{mount: "identity", path: "mfa/method/okta/{method_id}"}:                           {name: "IdentityMFAMethodOkta"},
	{mount: "identity", path: "mfa/method/pingid"}:                                     {name: "IdentityMFAMethodsPingID"},
	{mount: "identity", path: "mfa/method/pingid/{method_id}"}:                         {name: "IdentityMFAMethodPingID"},
	{mount: "identity", path: "mfa/method/totp"}:                                       {name: "IdentityMFAMethodsTOTP"},
	{mount: "identity", path: "mfa/method/totp/admin-destroy"}:                         {name: "IdentityMFAMethodTOTPAdminDestroy"},
	{mount: "identity", path: "mfa/method/totp/admin-generate"}:                        {name: "IdentityMFAMethodTOTPAdminGenerate"},
	{mount: "identity", path: "mfa/method/totp/generate"}:                              {name: "IdentityMFAMethodTOTPGenerate"},
	{mount: "identity", path: "mfa/method/totp/{method_id}"}:                           {name: "IdentityMFAMethodTOTP"},
	{mount: "identity", path: "mfa/method/{method_id}"}:                                {name: "IdentityMFAMethod"},
	{mount: "identity", path: "oidc/.well-known/keys"}:                                 {name: "IdentityOIDCWellKnownKeys"},
	{mount: "identity", path: "oidc/.well-known/openid-configuration"}:                 {name: "IdentityOIDCWellKnownOpenIDConfiguration"},
	{mount: "identity", path: "oidc/assignment"}:                                       {name: "IdentityOIDCAssignments"},
	{mount: "identity", path: "oidc/assignment/{name}"}:                                {name: "IdentityOIDCAssignment"},
	{mount: "identity", path: "oidc/client"}:                                           {name: "IdentityOIDCClients"},
	{mount: "identity", path: "oidc/client/{name}"}:                                    {name: "IdentityOIDCClient"},
	{mount: "identity", path: "oidc/config"}:                                           {name: "IdentityOIDCConfig"},
	{mount: "identity", path: "oidc/introspect"}:                                       {name: "IdentityOIDCIntrospect"},
	{mount: "identity", path: "oidc/key"}:                                              {name: "IdentityOIDCKeys"},
	{mount: "identity", path: "oidc/key/{name}"}:                                       {name: "IdentityOIDCKey"},
	{mount: "identity", path: "oidc/key/{name}/rotate"}:                                {name: "IdentityOIDCKeyRotate"},
	{mount: "identity", path: "oidc/provider"}:                                         {name: "IdentityOIDCProviders"},
	{mount: "identity", path: "oidc/provider/{name}"}:                                  {name: "IdentityOIDCProvider"},
	{mount: "identity", path: "oidc/provider/{name}/.well-known/keys"}:                 {name: "IdentityOIDCProviderWellKnownKeys"},
	{mount: "identity", path: "oidc/provider/{name}/.well-known/openid-configuration"}: {name: "IdentityOIDCProviderWellKnownOpenIDConfiguration"},
	{mount: "identity", path: "oidc/provider/{name}/authorize"}:                        {name: "IdentityOIDCProviderAuthorize"},
	{mount: "identity", path: "oidc/provider/{name}/token"}:                            {name: "IdentityOIDCProviderToken"},
	{mount: "identity", path: "oidc/provider/{name}/userinfo"}:                         {name: "IdentityOIDCProviderUserInfo"},
	{mount: "identity", path: "oidc/role"}:                                             {name: "IdentityOIDCRoles"},
	{mount: "identity", path: "oidc/role/{name}"}:                                      {name: "IdentityOIDCRole"},
	{mount: "identity", path: "oidc/scope"}:                                            {name: "IdentityOIDCScopes"},
	{mount: "identity", path: "oidc/scope/{name}"}:                                     {name: "IdentityOIDCScope"},
	{mount: "identity", path: "oidc/token/{name}"}:                                     {name: "IdentityOIDCToken"},
	{mount: "identity", path: "persona"}:                                               {name: "IdentityPersona"},
	{mount: "identity", path: "persona/id"}:                                            {name: "IdentityPersonaIDs"},
	{mount: "identity", path: "persona/id/{id}"}:                                       {name: "IdentityPersonaID"},
	{mount: "kubernetes", path: "config"}:                                              {name: "KubernetesConfig"},
	{mount: "kubernetes", path: "creds/{name}"}:                                        {name: "KubernetesCredentials"},
	{mount: "kubernetes", path: "roles"}:                                               {name: "KubernetesRoles"},
	{mount: "kubernetes", path: "roles/{name}"}:                                        {name: "KubernetesRole"},
	{mount: "kv", path: "{path}"}:                                                      {name: "KVv1Secret"},
	{mount: "ldap", path: "config"}:                                                    {name: "LDAPConfig"},
	{mount: "ldap", path: "creds/{name}"}:                                              {name: "LDAPCredentials"},
	{mount: "ldap", path: "library"}:                                                   {name: "LDAPLibraries"},
	{mount: "ldap", path: "library/manage/{name}/check-in"}:                            {name: "LDAPLibraryManageCheckIn"},
	{mount: "ldap", path: "library/{name}"}:                                            {name: "LDAPLibrary"},
	{mount: "ldap", path: "library/{name}/check-in"}:                                   {name: "LDAPLibraryCheckIn"},
	{mount: "ldap", path: "library/{name}/check-out"}:                                  {name: "LDAPLibraryCheckOut"},
	{mount: "ldap", path: "library/{name}/status"}:                                     {name: "LDAPLibraryStatus"},
	{mount: "ldap", path: "role"}:                                                      {name: "LDAPRoles"},
	{mount: "ldap", path: "role/{name}"}:                                               {name: "LDAPRole"},
	{mount: "ldap", path: "rotate-role/{name}"}:                                        {name: "LDAPRotateRole"},
	{mount: "ldap", path: "rotate-root"}:                                               {name: "LDAPRotateRoot"},
	{mount: "ldap", path: "static-cred/{name}"}:                                        {name: "LDAPStaticCredentials"},
	{mount: "ldap", path: "static-role"}:                                               {name: "LDAPStaticRoles"},
	{mount: "ldap", path: "static-role/{name}"}:                                        {name: "LDAPStaticRole"},
	{mount: "mongodbatlas", path: "config"}:                                            {name: "MongoDBAtlasConfig"},
	{mount: "mongodbatlas", path: "creds/{name}"}:                                      {name: "MongoDBAtlasCredentials"},
	{mount: "mongodbatlas", path: "roles"}:                                             {name: "MongoDBAtlasRoles"},
	{mount: "mongodbatlas", path: "roles/{name}"}:                                      {name: "MongoDBAtlasRole"},
	{mount: "nomad", path: "config/access"}:                                            {name: "NomadConfigAccess"},
	{mount: "nomad", path: "config/lease"}:                                             {name: "NomadConfigLease"},
	{mount: "nomad", path: "creds/{name}"}:                                             {name: "NomadCredentials"},
	{mount: "nomad", path: "role"}:                                                     {name: "NomadRoles"},
	{mount: "nomad", path: "role/{name}"}:                                              {name: "NomadRole"},
	{mount: "openldap", path: "config"}:                                                {name: "OpenLDAPConfig"},
	{mount: "openldap", path: "creds/{name}"}:                                          {name: "OpenLDAPCredentials"},
	{mount: "openldap", path: "library"}:                                               {name: "OpenLDAPLibraries"},
	{mount: "openldap", path: "library/manage/{name}/check-in"}:                        {name: "OpenLDAPLibraryManageCheckIn"},
	{mount: "openldap", path: "library/{name}"}:                                        {name: "OpenLDAPLibrary"},
	{mount: "openldap", path: "library/{name}/check-in"}:                               {name: "OpenLDAPLibraryCheckIn"},
	{mount: "openldap", path: "library/{name}/check-out"}:                              {name: "OpenLDAPLibraryCheckOut"},
	{mount: "openldap", path: "library/{name}/status"}:                                 {name: "OpenLDAPLibraryStatus"},
	{mount: "openldap", path: "role"}:                                                  {name: "OpenLDAPRoles"},
	{mount: "openldap", path: "role/{name}"}:                                           {name: "OpenLDAPRole"},
	{mount: "openldap", path: "rotate-role/{name}"}:                                    {name: "OpenLDAPRotateRole"},
	{mount: "openldap", path: "rotate-root"}:                                           {name: "OpenLDAPRotateRoot"},
	{mount: "openldap", path: "static-cred/{name}"}:                                    {name: "OpenLDAPStaticCredentials"},
	{mount: "openldap", path: "static-role"}:                                           {name: "OpenLDAPStaticRoles"},
	{mount: "openldap", path: "static-role/{name}"}:                                    {name: "OpenLDAPStaticRole"},
	{mount: "pki", path: "bundle"}:                                                     {name: "PKIBundle"},
	{mount: "pki", path: "ca"}:                                                         {name: "PKICA"},
	{mount: "pki", path: "ca/pem"}:                                                     {name: "PKICAPem"},
	{mount: "pki", path: "ca_chain"}:                                                   {name: "PKICaChain"},
	{mount: "pki", path: "cert"}:                                                       {name: "PKICerts"},
	{mount: "pki", path: "cert/ca_chain"}:                                              {name: "PKICertCaChain"},
	{mount: "pki", path: "cert/{serial}"}:                                              {name: "PKICert"},
	{mount: "pki", path: "cert/{serial}/raw"}:                                          {name: "PKICertRaw"},
	{mount: "pki", path: "cert/{serial}/raw/pem"}:                                      {name: "PKICertRawPem"},
	{mount: "pki", path: "certs"}:                                                      {name: "PKICerts"},
	{mount: "pki", path: "certs/revoked"}:                                              {name: "PKICertsRevoked"},
	{mount: "pki", path: "config/auto-tidy"}:                                           {name: "PKIConfigAutoTidy"},
	{mount: "pki", path: "config/ca"}:                                                  {name: "PKIConfigCa"},
	{mount: "pki", path: "config/cluster"}:                                             {name: "PKIConfigCluster"},
	{mount: "pki", path: "config/crl"}:                                                 {name: "PKIConfigCRL"},
	{mount: "pki", path: "config/issuers"}:                                             {name: "PKIConfigIssuers"},
	{mount: "pki", path: "config/keys"}:                                                {name: "PKIConfigKeys"},
	{mount: "pki", path: "config/urls"}:                                                {name: "PKIConfigURLs"},
	{mount: "pki", path: "crl"}:                                                        {name: "PKICRL"},
	{mount: "pki", path: "crl/rotate"}:                                                 {name: "PKICRLRotate"},
	{mount: "pki", path: "crl/rotate-delta"}:                                           {name: "PKICRLRotateDelta"},
	{mount: "pki", path: "delta-crl"}:                                                  {name: "PKIDeltaCRL"},
	{mount: "pki", path: "intermediate/cross-sign"}:                                    {name: "PKIIntermediateCrossSign"},
	{mount: "pki", path: "intermediate/generate/{exported}"}:                           {name: "PKIIntermediateGenerate"},
	{mount: "pki", path: "intermediate/set-signed"}:                                    {name: "PKIIntermediateSetSigned"},
	{mount: "pki", path: "internal|exported"}:                                          {name: "PKIInternalExported"},
	{mount: "pki", path: "issue/{role}"}:                                               {name: "PKIIssueRole"},
	{mount: "pki", path: "issuer/{issuer_ref}/issue/{role}"}:                           {name: "PKIIssuerIssueRole"},
	{mount: "pki", path: "issuer/{issuer_ref}/resign-crls"}:                            {name: "PKIIssuerResignCRLs"},
	{mount: "pki", path: "issuer/{issuer_ref}/revoke"}:                                 {name: "PKIIssuerRevoke"},
	{mount: "pki", path: "issuer/{issuer_ref}/sign-intermediate"}:                      {name: "PKIIssuerSignIntermediate"},
	{mount: "pki", path: "issuer/{issuer_ref}/sign-revocation-list"}:                   {name: "PKIIssuerSignRevocationList"},
	{mount: "pki", path: "issuer/{issuer_ref}/sign-self-issued"}:                       {name: "PKIIssuerSignSelfIssued"},
	{mount: "pki", path: "issuer/{issuer_ref}/sign-verbatim"}:                          {name: "PKIIssuerSignVerbatim"},
	{mount: "pki", path: "issuer/{issuer_ref}/sign-verbatim/{role}"}:                   {name: "PKIIssuerSignVerbatimRole"},
	{mount: "pki", path: "issuer/{issuer_ref}/sign/{role}"}:                            {name: "PKIIssuerSignRole"},
	{mount: "pki", path: "issuers"}:                                                    {name: "PKIIssuers"},
	{mount: "pki", path: "issuers/generate/intermediate/{exported}"}:                   {name: "PKIIssuersGenerateIntermediateExported"},
	{mount: "pki", path: "issuers/generate/root/{exported}"}:                           {name: "PKIIssuersGenerateRootExported"},
	{mount: "pki", path: "key/{key_ref}"}:                                              {name: "PKIKey"},
	{mount: "pki", path: "keys"}:                                                       {name: "PKIKeys"},
	{mount: "pki", path: "keys/import"}:                                                {name: "PKIKeysImport"},
	{mount: "pki", path: "kms"}:                                                        {name: "PKIKMS"},
	{mount: "pki", path: "ocsp"}:                                                       {name: "PKIOCSP"},
	{mount: "pki", path: "ocsp/{req}"}:                                                 {name: "PKIOCSPReq"},
	{mount: "pki", path: "revoke"}:                                                     {name: "PKIRevoke"},
	{mount: "pki", path: "revoke-with-key"}:                                            {name: "PKIRevokeWithKey"},
	{mount: "pki", path: "roles"}:                                                      {name: "PKIRoles"},
	{mount: "pki", path: "roles/{name}"}:                                               {name: "PKIRole"},
	{mount: "pki", path: "root"}:                                                       {name: "PKIRoot"},
	{mount: "pki", path: "root/generate/{exported}"}:                                   {name: "PKIRootGenerate"},
	{mount: "pki", path: "root/replace"}:                                               {name: "PKIRootReplace"},
	{mount: "pki", path: "root/rotate/{exported}"}:                                     {name: "PKIRootRotate"},
	{mount: "pki", path: "root/sign-intermediate"}:                                     {name: "PKIRootSignIntermediate"},
	{mount: "pki", path: "root/sign-self-issued"}:                                      {name: "PKIRootSignSelfIssued"},
	{mount: "pki", path: "sign-verbatim"}:                                              {name: "PKISignVerbatim"},
	{mount: "pki", path: "sign-verbatim/{role}"}:                                       {name: "PKISignVerbatimRole"},
	{mount: "pki", path: "sign/{role}"}:                                                {name: "PKISignRole"},
	{mount: "pki", path: "tidy"}:                                                       {name: "PKITidy"},
	{mount: "pki", path: "tidy-cancel"}:                                                {name: "PKITidyCancel"},
	{mount: "pki", path: "tidy-status"}:                                                {name: "PKITidyStatus"},
	{mount: "rabbitmq", path: "config/connection"}:                                     {name: "RabbitMQConfigConnection"},
	{mount: "rabbitmq", path: "config/lease"}:                                          {name: "RabbitMQConfigLease"},
	{mount: "rabbitmq", path: "creds/{name}"}:                                          {name: "RabbitMQCredentials"},
	{mount: "rabbitmq", path: "roles"}:                                                 {name: "RabbitMQRoles"},
	{mount: "rabbitmq", path: "roles/{name}"}:                                          {name: "RabbitMQRole"},
	{mount: "secret", path: "config"}:                                                  {name: "SecretConfig"},
	{mount: "secret", path: "data/{path}"}:                                             {name: "Secret"},
	{mount: "secret", path: "delete/{path}"}:                                           {name: "DeleteSecretVersions", operationImplied: true},
	{mount: "secret", path: "destroy/{path}"}:                                          {name: "DestroySecretVersions", operationImplied: true},
	{mount: "secret", path: "metadata/{path}"}:                                         {name: "SecretMetadata"},
	{mount: "secret", path: "subkeys/{path}"}:                                          {name: "SecretSubkeys"},
	{mount: "secret", path: "undelete/{path}"}:                                         {name: "UndeleteSecretVersions", operationImplied: true},
	{mount: "ssh", path: "config/ca"}:                                                  {name: "SSHConfigCA"},
	{mount: "ssh", path: "config/zeroaddress"}:                                         {name: "SSHConfigZeroAddress"},
	{mount: "ssh", path: "creds/{role}"}:                                               {name: "SSHCredentials"},
	{mount: "ssh", path: "issue/{role}"}:                                               {name: "SSHIssue", operationImplied: true},
	{mount: "ssh", path: "keys/{key_name}"}:                                            {name: "SSHKeys"},
	{mount: "ssh", path: "lookup"}:                                                     {name: "SSHLookup", operationImplied: true},
	{mount: "ssh", path: "public_key"}:                                                 {name: "SSHPublicKey"},
	{mount: "ssh", path: "roles"}:                                                      {name: "SSHRoles"},
	{mount: "ssh", path: "roles/{role}"}:                                               {name: "SSHRole"},
	{mount: "ssh", path: "sign/{role}"}:                                                {name: "SSHSign", operationImplied: true},
	{mount: "ssh", path: "verify"}:                                                     {name: "SSHVerify", operationImplied: true},
	{mount: "sys", path: "audit"}:                                                      {name: "AuditDevices"},
	{mount: "sys", path: "audit/{path}"}:                                               {name: "AuditDevicesWith"},
	{mount: "sys", path: "audit-hash/{path}"}:                                          {name: "AuditHash"},
	{mount: "sys", path: "auth"}:                                                       {name: "AuthMethods"},
	{mount: "sys", path: "auth/{path}"}:                                                {name: "AuthMethodsWith"},
	{mount: "sys", path: "auth/{path}/tune"}:                                           {name: "AuthMethodTune"},
	{mount: "sys", path: "capabilities"}:                                               {name: "Capabilities"},
	{mount: "sys", path: "capabilities-accessor"}:                                      {name: "CapabilitiesAccessor"},
	{mount: "sys", path: "capabilities-self"}:                                          {name: "CapabilitiesSelf"},
	{mount: "sys", path: "config/auditing/request-headers"}:                            {name: "ConfigAuditingRequestHeaders"},
	{mount: "sys", path: "config/auditing/request-headers/{header}"}:                   {name: "ConfigAuditingRequestHeader"},
	{mount: "sys", path: "config/cors"}:                                                {name: "ConfigCORS"},
	{mount: "sys", path: "config/reload/{subsystem}"}:                                  {name: "ConfigReloadSubsystem"},
	{mount: "sys", path: "config/state/sanitized"}:                                     {name: "ConfigStateSanitized"},
	{mount: "sys", path: "config/ui/headers/"}:                                         {name: "ConfigUIHeaders"},
	{mount: "sys", path: "config/ui/headers/{header}"}:                                 {name: "ConfigUIHeader"},
	{mount: "sys", path: "generate-root"}:                                              {name: "GenerateRoot"},
	{mount: "sys", path: "generate-root/attempt"}:                                      {name: "GenerateRootAttempt"},
	{mount: "sys", path: "generate-root/update"}:                                       {name: "GenerateRootUpdate"},
	{mount: "sys", path: "ha-status"}:                                                  {name: "HAStatus"},
	{mount: "sys", path: "health"}:                                                     {name: "Health"},
	{mount: "sys", path: "host-info"}:                                                  {name: "HostInfo"},
	{mount: "sys", path: "in-flight-req"}:                                              {name: "InFlightRequests"},
	{mount: "sys", path: "init"}:                                                       {name: "Init", operationImplied: true},
	{mount: "sys", path: "internal/counters/activity"}:                                 {name: "InternalCountersActivity"},
	{mount: "sys", path: "internal/counters/activity/export"}:                          {name: "InternalCountersActivityExport"},
	{mount: "sys", path: "internal/counters/activity/monthly"}:                         {name: "InternalCountersActivityMonthly"},
	{mount: "sys", path: "internal/counters/config"}:                                   {name: "InternalCountersConfig"},
	{mount: "sys", path: "internal/counters/entities"}:                                 {name: "InternalCountersEntities"},
	{mount: "sys", path: "internal/counters/requests"}:                                 {name: "InternalCountersRequests"},
	{mount: "sys", path: "internal/counters/tokens"}:                                   {name: "InternalCountersTokens"},
	{mount: "sys", path: "internal/inspect/router/{tag}"}:                              {name: "InternalInspectRouter"},
	{mount: "sys", path: "internal/specs/openapi"}:                                     {name: "InternalSpecsOpenAPI"},
	{mount: "sys", path: "internal/ui/feature-flags"}:                                  {name: "InternalUIFeatureFlags"},
	{mount: "sys", path: "internal/ui/mounts"}:                                         {name: "InternalUIMounts"},
	{mount: "sys", path: "internal/ui/mounts/{path}"}:                                  {name: "InternalUIMount"},
	{mount: "sys", path: "internal/ui/namespaces"}:                                     {name: "InternalUINamespaces"},
	{mount: "sys", path: "internal/ui/resultant-acl"}:                                  {name: "InternalUIResultantACL"},
	{mount: "sys", path: "key-status"}:                                                 {name: "KeyStatus"},
	{mount: "sys", path: "leader"}:                                                     {name: "Leader"},
	{mount: "sys", path: "leases"}:                                                     {name: "Leases"},
	{mount: "sys", path: "leases/count"}:                                               {name: "LeasesCount"},
	{mount: "sys", path: "leases/lookup"}:                                              {name: "LeasesLookup"},
	{mount: "sys", path: "leases/lookup/{prefix}"}:                                     {name: "LeasesLookupPrefix"},
	{mount: "sys", path: "leases/renew"}:                                               {name: "LeasesRenew"},
	{mount: "sys", path: "leases/renew/{url_lease_id}"}:                                {name: "LeasesRenew2"},
	{mount: "sys", path: "leases/revoke"}:                                              {name: "LeasesRevoke"},
	{mount: "sys", path: "leases/revoke-force/{prefix}"}:                               {name: "LeasesRevokeForce"},
	{mount: "sys", path: "leases/revoke-prefix/{prefix}"}:                              {name: "LeasesRevokePrefix"},
	{mount: "sys", path: "leases/revoke/{url_lease_id}"}:                               {name: "LeasesRevoke2"},
	{mount: "sys", path: "leases/tidy"}:                                                {name: "LeasesTidy"},
	{mount: "sys", path: "loggers"}:                                                    {name: "Loggers"},
	{mount: "sys", path: "loggers/{name}"}:                                             {name: "Logger"},
	{mount: "sys", path: "metrics"}:                                                    {name: "Metrics"},
	{mount: "sys", path: "mfa/validate"}:                                               {name: "MFAValidate"},
	{mount: "sys", path: "monitor"}:                                                    {name: "Monitor", operationImplied: true},
	{mount: "sys", path: "mounts"}:                                                     {name: "Mounts"},
	{mount: "sys", path: "mounts/{path}"}:                                              {name: "MountsWith"},
	{mount: "sys", path: "mounts/{path}/tune"}:                                         {name: "MountsTune"},
	{mount: "sys", path: "plugins/catalog"}:                                            {name: "PluginsCatalog"},
	{mount: "sys", path: "plugins/catalog/{type}"}:                                     {name: "PluginsCatalogByType"},
	{mount: "sys", path: "plugins/catalog/{type}/{name}"}:                              {name: "PluginsCatalogByTypeByName"},
	{mount: "sys", path: "plugins/reload/backend"}:                                     {name: "PluginsReloadBackend"},
	{mount: "sys", path: "policies/acl"}:                                               {name: "PoliciesACL"},
	{mount: "sys", path: "policies/acl/{name}"}:                                        {name: "PoliciesACL"},
	{mount: "sys", path: "policies/password"}:                                          {name: "PoliciesPassword"},
	{mount: "sys", path: "policies/password/{name}"}:                                   {name: "PoliciesPassword"},
	{mount: "sys", path: "policies/password/{name}/generate"}:                          {name: "PoliciesPasswordGenerate"},
	{mount: "sys", path: "policy"}:                                                     {name: "Policies"},
	{mount: "sys", path: "policy/{name}"}:                                              {name: "Policy"},
	{mount: "sys", path: "pprof/"}:                                                     {name: "Pprof"},
	{mount: "sys", path: "pprof/allocs"}:                                               {name: "PprofAllocs"},
	{mount: "sys", path: "pprof/block"}:                                                {name: "PprofBlock"},
	{mount: "sys", path: "pprof/cmdline"}:                                              {name: "PprofCmdline"},
	{mount: "sys", path: "pprof/goroutine"}:                                            {name: "PprofGoroutine"},
	{mount: "sys", path: "pprof/heap"}:                                                 {name: "PprofHeap"},
	{mount: "sys", path: "pprof/mutex"}:                                                {name: "PprofMutex"},
	{mount: "sys", path: "pprof/profile"}:                                              {name: "PprofProfile"},
	{mount: "sys", path: "pprof/symbol"}:                                               {name: "PprofSymbol"},
	{mount: "sys", path: "pprof/threadcreate"}:                                         {name: "PprofThreadcreate"},
	{mount: "sys", path: "pprof/trace"}:                                                {name: "PprofTrace"},
	{mount: "sys", path: "quotas/config"}:                                              {name: "QuotasConfig"},
	{mount: "sys", path: "quotas/rate-limit"}:                                          {name: "QuotasRateLimits"},
	{mount: "sys", path: "quotas/rate-limit/{name}"}:                                   {name: "QuotasRateLimit"},
	{mount: "sys", path: "raw"}:                                                        {name: "Raw"},
	{mount: "sys", path: "raw/{path}"}:                                                 {name: "RawPath"},
	{mount: "sys", path: "rekey/backup"}:                                               {name: "RekeyBackup", operationImplied: true},
	{mount: "sys", path: "rekey/init"}:                                                 {name: "RekeyInit", operationImplied: true},
	{mount: "sys", path: "rekey/recovery-key-backup"}:                                  {name: "RekeyRecoveryKeyBackup", operationImplied: true},
	{mount: "sys", path: "rekey/update"}:                                               {name: "RekeyUpdate", operationImplied: true},
	{mount: "sys", path: "rekey/verify"}:                                               {name: "RekeyVerify", operationImplied: true},
	{mount: "sys", path: "remount"}:                                                    {name: "Remount", operationImplied: true},
	{mount: "sys", path: "remount/status/{migration_id}"}:                              {name: "RemountStatus"},
	{mount: "sys", path: "renew"}:                                                      {name: "Renew", operationImplied: true},
	{mount: "sys", path: "renew/{url_lease_id}"}:                                       {name: "Renew", operationImplied: true},
	{mount: "sys", path: "replication/status"}:                                         {name: "ReplicationStatus"},
	{mount: "sys", path: "revoke"}:                                                     {name: "Revoke", operationImplied: true},
	{mount: "sys", path: "revoke-force/{prefix}"}:                                      {name: "RevokeForce", operationImplied: true},
	{mount: "sys", path: "revoke-prefix/{prefix}"}:                                     {name: "RevokePrefix", operationImplied: true},
	{mount: "sys", path: "revoke/{url_lease_id}"}:                                      {name: "Revoke, operationImplied: true"},
	{mount: "sys", path: "rotate"}:                                                     {name: "Rotate", operationImplied: true},
	{mount: "sys", path: "rotate/config"}:                                              {name: "RotateConfig"},
	{mount: "sys", path: "seal"}:                                                       {name: "Seal", operationImplied: true},
	{mount: "sys", path: "seal-status"}:                                                {name: "SealStatus"},
	{mount: "sys", path: "step-down"}:                                                  {name: "StepDown", operationImplied: true},
	{mount: "sys", path: "tools/hash"}:                                                 {name: "ToolsHashes"},
	{mount: "sys", path: "tools/hash/{urlalgorithm}"}:                                  {name: "ToolsHash"},
	{mount: "sys", path: "tools/random"}:                                               {name: "ToolsRandom"},
	{mount: "sys", path: "tools/random/{source}"}:                                      {name: "ToolsRandomSource"},
	{mount: "sys", path: "tools/random/{source}/{urlbytes}"}:                           {name: "ToolsRandomSourceBytes"},
	{mount: "sys", path: "unseal"}:                                                     {name: "Unseal", operationImplied: true},
	{mount: "sys", path: "version-history/"}:                                           {name: "VersionHistory"},
	{mount: "sys", path: "wrapping/lookup"}:                                            {name: "WrappingLookup", operationImplied: true},
	{mount: "sys", path: "wrapping/rewrap"}:                                            {name: "WrappingRewrap", operationImplied: true},
	{mount: "sys", path: "wrapping/unwrap"}:                                            {name: "WrappingUnwrap", operationImplied: true},
	{mount: "sys", path: "wrapping/wrap"}:                                              {name: "WrappingWrap", operationImplied: true},
	{mount: "terraform", path: "config"}:                                               {name: "TerraformConfig"},
	{mount: "terraform", path: "creds/{name}"}:                                         {name: "TerraformCredentials"},
	{mount: "terraform", path: "role"}:                                                 {name: "TerraformRoles"},
	{mount: "terraform", path: "role/{name}"}:                                          {name: "TerraformRole"},
	{mount: "terraform", path: "rotate-role/{name}"}:                                   {name: "TerraformRotateRole"},
	{mount: "totp", path: "code/{name}"}:                                               {name: "TOTPCode"},
	{mount: "totp", path: "keys"}:                                                      {name: "TOTPKeys"},
	{mount: "totp", path: "keys/{name}"}:                                               {name: "TOTPKey"},
	{mount: "transit", path: "backup/{name}"}:                                          {name: "TransitBackup"},
	{mount: "transit", path: "cache-config"}:                                           {name: "TransitCacheConfig"},
	{mount: "transit", path: "datakey/{plaintext}/{name}"}:                             {name: "TransitDatakey"},
	{mount: "transit", path: "decrypt/{name}"}:                                         {name: "TransitDecrypt"},
	{mount: "transit", path: "encrypt/{name}"}:                                         {name: "TransitEncrypt"},
	{mount: "transit", path: "export/{type}/{name}"}:                                   {name: "TransitExport"},
	{mount: "transit", path: "export/{type}/{name}/{version}"}:                         {name: "TransitExportVersion"},
	{mount: "transit", path: "hash"}:                                                   {name: "TransitHash"},
	{mount: "transit", path: "hash/{urlalgorithm}"}:                                    {name: "TransitHash"},
	{mount: "transit", path: "hmac/{name}"}:                                            {name: "TransitHMAC"},
	{mount: "transit", path: "hmac/{name}/{urlalgorithm}"}:                             {name: "TransitHMACAlgorithm"},
	{mount: "transit", path: "keys"}:                                                   {name: "TransitKeys"},
	{mount: "transit", path: "keys/{name}"}:                                            {name: "TransitKey"},
	{mount: "transit", path: "keys/{name}/config"}:                                     {name: "TransitKeyConfig"},
	{mount: "transit", path: "keys/{name}/import"}:                                     {name: "TransitKeyImport"},
	{mount: "transit", path: "keys/{name}/import_version"}:                             {name: "TransitKeyImportVersion"},
	{mount: "transit", path: "keys/{name}/rotate"}:                                     {name: "TransitKeyRotate"},
	{mount: "transit", path: "keys/{name}/trim"}:                                       {name: "TransitKeyTrim"},
	{mount: "transit", path: "random"}:                                                 {name: "TransitRandom"},
	{mount: "transit", path: "random/{source}"}:                                        {name: "TransitRandomSource"},
	{mount: "transit", path: "random/{source}/{urlbytes}"}:                             {name: "TransitRandomSourceBytes"},
	{mount: "transit", path: "restore"}:                                                {name: "TransitRestore"},
	{mount: "transit", path: "restore/{name}"}:                                         {name: "TransitRestoreKey"},
	{mount: "transit", path: "rewrap/{name}"}:                                          {name: "TransitRewrap"},
	{mount: "transit", path: "sign/{name}"}:                                            {name: "TransitSign"},
	{mount: "transit", path: "sign/{name}/{urlalgorithm}"}:                             {name: "TransitSignAlgorithm"},
	{mount: "transit", path: "verify/{name}"}:                                          {name: "TransitVerify"},
	{mount: "transit", path: "verify/{name}/{urlalgorithm}"}:                           {name: "TransitVerifyAlgorithm"},
	{mount: "transit", path: "wrapping_key"}:                                           {name: "TransitWrappingKey"},
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
	// For most request names "write" seems to be more applicable in place of "update"
	var operationStr string
	if operation == logical.UpdateOperation {
		operationStr = "write"
	} else {
		operationStr = string(operation)
	}

	var parts []string

	// Append either the known mapping or operation + mount + path split by non-word characters
	if mapping, ok := knownPathMappings[knownPathKey{mount: mount, path: path}]; ok {
		// Certain names have operations implied in the name, e.g. Seal/Unseal
		if !mapping.operationImplied {
			parts = append(parts, nonWordRe.Split(strings.ToLower(operationStr), -1)...)
		}
		parts = append(parts, mapping.name)
	} else {
		parts = append(parts, nonWordRe.Split(strings.ToLower(operationStr), -1)...)
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
