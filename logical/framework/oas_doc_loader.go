package framework

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/oas"
	"github.com/hashicorp/vault/logical"
)

// Regex for handling optional and named parameters, and string cleanup
var optRe = regexp.MustCompile(`(?U)\(.*\)\?`)
var altRe = regexp.MustCompile(`\((.*)\|(.*)\)`)
var reqdRe = regexp.MustCompile(`\(?\?P<(\w+)>[^)]*\)?`)
var cleanRe = regexp.MustCompile("[()^$?]")
var wsRe = regexp.MustCompile(`\s+`)

// DocumentPaths parses all paths in a framework.Backend into OAS paths.
func DocumentPaths(backend *Backend, doc *oas.OASDoc) {
	var rootPaths []string

	sp := backend.SpecialPaths()
	if sp != nil {
		rootPaths = sp.Root
	}

	for _, p := range backend.Paths {
		documentPath(p, rootPaths, doc)
	}
}

// documentPath parses a framework.Path into one or more oas.Paths,
// methods, etc. It will infer path and body fields from when they're
// provided, and use existing path help if available.
func documentPath(p *Path, rootPaths []string, doc *oas.OASDoc) {
	var httpMethod string

	paths := expandPattern(p.Pattern)

	for _, path := range paths {
		pm := oas.PathMethods{}

		// Test for exact or prefix match of root paths
		for _, root := range rootPaths {
			if root == path ||
				(strings.HasSuffix(root, "*") && strings.HasPrefix(path, root[0:len(root)-1])) {
				pm.Root = true
				break
			}
		}

		// Add details for every registered operation
		for opType := range p.Callbacks {
			m := oas.NewMethodDetail()

			m.Summary = cleanString(p.HelpSynopsis)
			m.Description = cleanString(p.HelpDescription)

			switch opType {
			case logical.CreateOperation, logical.UpdateOperation:
				httpMethod = "POST"
				m.Responses[200] = oas.StdRespOK
			case logical.DeleteOperation:
				httpMethod = "DELETE"
				m.Responses[204] = oas.StdRespNoContent
			case logical.ReadOperation, logical.ListOperation:
				httpMethod = "GET"
				m.Responses[200] = oas.StdRespOK
			default:
				log.L().Warn("unknown operation", "type", opType)
				httpMethod = "GET"
				m.Responses[200] = oas.StdRespOK
			}

			// Extract path fields into parameters
			fieldSet := make(map[string]bool)
			params := oas.PathFields(path)

			for _, param := range params {
				fieldSet[param] = true
				oasType := convertType(p.Fields[param].Type)
				p := oas.Parameter{
					Name:        param,
					Description: cleanString(p.Fields[param].Description),
					In:          "path",
					Type:        oasType.typ,
					Required:    true,
				}
				m.Parameters = append(m.Parameters, p)
			}

			// LIST is exported as a GET with a `list` query parameter
			if opType == logical.ListOperation {
				m.Parameters = append(m.Parameters, oas.Parameter{
					Name:        "list",
					Description: "Return a list if `true`",
					In:          "query",
					Type:        "string",
				})
			}

			// Add any fields not present in the path as body parameters for POST
			if httpMethod == "POST" {
				s := oas.NewSchema()
				for name, field := range p.Fields {
					if !fieldSet[name] {
						oasField := convertType(field.Type)
						p := oas.Property{
							Type:        oasField.typ,
							Description: cleanString(field.Description),
							Format:      oasField.format,
							Attrs:       field.Attrs,
						}
						if oasField.typ == "array" {
							p.Items = &oas.Property{
								Type: oasField.items,
							}
						}
						s.Properties[name] = &p
					}
				}

				m.Parameters = append(m.Parameters,
					oas.Parameter{
						Name:     "body",
						In:       "body",
						Schema:   s,
						Required: true,
					},
				)
			}

			// Add explicitly specified reponse details
			if respSchema, ok := p.Responses[opType]; ok {
				for code, resp := range respSchema {
					m.Responses[code] = oas.Response{
						Description: resp.Description,
						Example:     resp.Example,
					}
				}
			}

			switch httpMethod {
			case "POST":
				pm.Post = m
			case "GET":
				pm.Get = m
			case "DELETE":
				pm.Delete = m
			}
		}

		doc.Paths["/"+path] = &pm
	}
}

// expandPattern expands a regex pattern by generating permutations of any optional parameters
// and changing named parameters into their {openapi} equivalents.
func expandPattern(pattern string) []string {
	var paths []string

	// This construct is added by GenericNameRegex and is much easier to remove now
	// than to compensate for in the other regexes.
	pattern = strings.Replace(pattern, `\w(([\w-.]+)?\w)?`, "", -1)

	// Initialize paths with the original pattern or the halves of an
	// alternation, which is also present in some patterns.
	matches := altRe.FindAllStringSubmatch(pattern, -1)
	if len(matches) > 0 {
		paths = []string{matches[0][1], matches[0][2]}
	} else {
		paths = []string{pattern}
	}

	// expand all optional regex elements into two paths. This approach is really only useful up to 2 optional
	// groups, but we probably don't want to deal with the exponential increase beyond that anyway.
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		match := optRe.FindStringIndex(p)
		if match != nil {
			paths[i] = p[0:match[0]] + p[match[0]+1:match[1]-2] + p[match[1]:]
			paths = append(paths, p[0:match[0]]+p[match[1]:])
			i--
		}
	}

	// replace named parameters (?P<foo>) with {foo}
	replacedPaths := make([]string, 0)
	for _, path := range paths {
		result := reqdRe.FindAllStringSubmatch(path, -1)
		if result != nil {
			for _, p := range result {
				par := p[1]
				path = strings.Replace(path, p[0], fmt.Sprintf("{%s}", par), 1)
			}
		}
		path = cleanRe.ReplaceAllString(path, "")
		replacedPaths = append(replacedPaths, path)
	}
	return replacedPaths
}

type oasType struct {
	typ    string
	items  string
	format string
}

// convertType translates a FieldType into an OpenAPI type.
// In the case of arrays, a subtype is returns as well.
func convertType(t FieldType) oasType {
	ret := oasType{}

	switch t {
	case TypeString, TypeNameString:
		ret.typ = "string"
	case TypeInt:
		ret.typ = "number"
	case TypeDurationSecond:
		ret.typ = "number"
		ret.format = "seconds"
	case TypeBool:
		ret.typ = "boolean"
	case TypeMap:
		ret.typ = "object"
		ret.format = "map"
	case TypeKVPairs:
		ret.typ = "object"
		ret.format = "kvpairs"
	case TypeSlice, TypeStringSlice, TypeCommaStringSlice:
		ret.typ = "array"
		ret.items = "string"
	case TypeCommaIntSlice:
		ret.typ = "array"
		ret.items = "number"
	default:
		log.L().Warn("error parsing field type", "type", t)
		ret.format = "unknown"
	}

	return ret
}

// cleanString prepares s for inclusion in the output
func cleanString(s string) string {
	s = strings.TrimSpace(s)
	return wsRe.ReplaceAllString(s, " ")
}
