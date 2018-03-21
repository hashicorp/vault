package framework

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/logical"
)

// LoadBackend parse paths in a framework.Backend into apidoc paths,
// methods, etc. It will infer path and body fields from when they're
// provided, and use existing path help if available.
func LoadBackend(backend *Backend, doc *Top) {
	// TODO: refactor!
	for _, p := range backend.Paths {
		procFrameworkPath(p, []string{}, doc)
		//doc.Mounts[prefix] = append(doc.Mounts[prefix], paths...)
	}
}

// procFrameworkPath parses a framework.Path into one or more apidoc.Paths.
func procFrameworkPath(p *Path, rootPaths []string, top *Top) {
	var httpMethod string

	paths := expandPattern(p.Pattern)

	for _, path := range paths {
		pm := PathMethods{}

		for _, root := range rootPaths {
			if root == path ||
				(strings.HasSuffix(root, "*") && strings.HasPrefix(path, root[0:len(root)-1])) {
				pm.Root = true
				break
			}
		}

		for opType := range p.Callbacks {
			m := NewMethodDetail()
			m.Summary = cleanString(p.HelpSynopsis)
			m.Description = cleanString(p.HelpDescription)

			switch opType {
			case logical.CreateOperation, logical.UpdateOperation:
				httpMethod = "POST"
				m.Responses[200] = StdRespOK2
			case logical.DeleteOperation:
				httpMethod = "DELETE"
				m.Responses[204] = StdRespNoContent2
			case logical.ReadOperation, logical.ListOperation:
				httpMethod = "GET"
				m.Responses[200] = StdRespOK2
			default:
				panic(fmt.Sprintf("unknown operation type %v", opType))
			}

			fieldSet := make(map[string]bool)
			params := pathFields(path)

			// Extract path fields
			for _, param := range params {
				fieldSet[param] = true
				oasType := convertType(p.Fields[param].Type)
				p := Parameter{
					Name:        param,
					Description: cleanString(p.Fields[param].Description),
					In:          "path",
					Type:        oasType.typ,
					Required:    true,
					//SubType:     sub,
				}
				m.Parameters = append(m.Parameters, p)
				//m.PathFields = append(m.PathFields, Property{
				//	Name:        param,
				//	Type:        typ,
				//	SubType:     sub,
				//	Description: cleanString(p.Fields[param].Description),
				//})
			}

			if opType == logical.ListOperation {
				m.Parameters = append(m.Parameters, Parameter{
					Name:        "list",
					Description: "Return a list if `true`",
					In:          "query",
					Type:        "string",
				})
			}

			// It's assumed that any fields not present in the path can be part of
			// the body for POST/PUT methods.
			if httpMethod == "POST" || httpMethod == "PUT" {
				s := NewSchema()
				for name, field := range p.Fields {
					if !fieldSet[name] {
						oasType := convertType(field.Type)
						p := Property2{
							Type:        oasType.typ,
							Description: cleanString(field.Description),
							Format:      oasType.format,
							//SubType:     sub,
						}
						if oasType.typ == "array" {
							p.Items = &Property2{
								Type: oasType.items,
								//Format: oasType.format,
							}
						}
						s.Properties[name] = p
						//m.Parameters = append(m.Parameters, p)

						//m.BodyFields = append(m.BodyFields, Property{
						//	Name:        name,
						//	Description: cleanString(field.Description),
						//	Type:        typ,
						//	SubType:     sub,
						//})
					}
				}
				m.Parameters = append(m.Parameters,
					Parameter{
						Name:   "body",
						In:     "body",
						Schema: s,
					},
				)

			}

			switch httpMethod {
			case "POST":
				pm.Post = m
			case "GET":
				pm.Get = m
			}

			//methods[httpMethod] = &m
		}
		top.Paths["/"+path] = &pm

		//if len(methods) > 0 {
		//	newPath := Path2{
		//		Pattern: path,
		//		Methods: methods,
		//	}
		//	docPaths = append(docPaths, newPath)
		//}
	}
}

// Regexen for handling optional and named parameters
var optRe = regexp.MustCompile(`(?U)\(.*\)\?`)
var reqdRe = regexp.MustCompile(`\(\?P<(\w+)>[^)]*\)`)
var cleanRe = regexp.MustCompile("[()^$?]")

// expandPattern expands a regex pattern by generating permutations of any optional parameters
// and changing named parameters into their {openapi} equivalents.
func expandPattern(pattern string) []string {

	// This construct is added by GenericNameRegex and is much easier to remove now
	// than to compensate for in the other regexes.
	pattern = strings.Replace(pattern, `\w(([\w-.]+)?\w)?`, "", -1)

	// expand all optional regex elements into two paths. This approach is really only useful up to 2 optional
	// groups, but we probably don't want to deal with the exponential increase beyond that anyway.
	paths := []string{pattern}

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

type OASType struct {
	typ    string
	items  string
	format string
}

// convertType translates a FieldType into an OpenAPI type.
// In the case of arrays, a subtype is returns as well.
func convertType(t FieldType) OASType {
	ret := OASType{typ: "string"}

	switch t {
	case TypeString, TypeNameString:
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

var wsRe = regexp.MustCompile(`\s+`)

// cleanString prepares s for inclusion in the output YAML. This is currently just
// basic escaping, whitespace thinning, and wrapping in quotes.
func cleanString(s string) string {
	s = strings.TrimSpace(s)

	// TODO: no truncation for now.
	//if idx := strings.Index(s, "\n"); idx != -1 {
	//	s = s[0:idx] + "..."
	//}

	s = wsRe.ReplaceAllString(s, " ")
	//s = strings.Replace(s, `"`, `\"`, -1)

	//return fmt.Sprintf(`"%s"`, s)
	return s
}
