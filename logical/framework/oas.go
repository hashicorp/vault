package framework

import (
	"regexp"

	"github.com/hashicorp/vault/version"
)

type Top struct {
	Swagger string     `json:"swagger"`
	Info    Info       `json:"info"`
	Paths   PathTopMap `json:"paths"`
}

func NewTop() Top {
	return Top{
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
	Summary     string            `json:"summary"`
	Description string            `json:"description"`
	Produces    []string          `json:"produces"`
	Parameters  []Parameter       `json:"parameters,omitempty"`
	Responses   map[int]Response2 `json:"responses"`
}

func NewMethodDetail() *MethodDetail {
	return &MethodDetail{
		Responses: make(map[int]Response2),
		Produces:  []string{"application/json"},
	}
}

type PathMethods struct {
	Get  *MethodDetail `json:"get,omitempty"`
	Post *MethodDetail `json:"post,omitempty"`
	Root bool          `json:"x-vault-root"`
}

type Parameter struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	In          string      `json:"in"`
	Type        string      `json:"type,omitempty"`
	Schema      interface{} `json:"schema,omitempty"`
	Required    bool        `json:"required,omitempty"`
}

type Response2 struct {
	Description string `json:"description"`
	Example     string `json:"example,omitempty"`
}

var StdRespOK2 = Response2{
	Description: "OK",
}

var StdRespNoContent2 = Response2{
	Description: "empty body",
}

type Schema struct {
	Type string `json:"type"`

	// TODO: should this really be interface{}?
	Properties map[string]interface{} `json:"properties"`
}

func NewSchema() *Schema {
	return &Schema{
		Type:       "object",
		Properties: make(map[string]interface{}),
	}
}

type Property2 struct {
	Type        string     `json:"type"`
	Description string     `json:"description,omitempty"`
	Items       *Property2 `json:"items,omitempty"`
	Format      string     `json:"format,omitempty"`
}

func pathFields(pattern string) []string {
	pathFieldsRe := regexp.MustCompile(`{(\w+)}`)

	r := pathFieldsRe.FindAllStringSubmatch(pattern, -1)
	ret := make([]string, 0, len(r))
	for _, t := range r {
		ret = append(ret, t[1])
	}
	return ret
}

//func procFrameworkPath2(p *Path) Top {
//	paths := procFrameworkPath(p)
//
//	ps := make(PathTopMap)
//
//	for _, path := range paths {
//		pm := PathMethods{
//			Get:  unpackMethod(path.Methods["GET"]),
//			Post: unpackMethod(path.Methods["POST"]),
//		}
//
//		ps[path.Pattern] = pm
//	}
//
//	t := Top{
//		Paths: ps,
//	}
//
//	return t
//}

//func unpackMethod(method *Method) *MethodDetail {
//	if method == nil {
//		return nil
//	}
//
//	parameters := make([]Parameter, 0)
//
//	for _, p := range method.PathFields {
//		parameter := Parameter{
//			In:          "path",
//			Name:        p.Name,
//			Description: p.Description,
//			Type:        p.Type,
//			Required:    true,
//		}
//
//		parameters = append(parameters, parameter)
//	}
//
//	if len(method.BodyFields) > 0 {
//		s := Schema{
//			Type:       "object",
//			Properties: make(map[string]interface{}),
//		}
//
//		for _, p := range method.BodyFields {
//			prop := Property2{
//				Description: p.Description,
//				Type:        p.Type,
//			}
//			if p.Type == "array" {
//				prop.Items = &Property2{
//					Type: p.SubType,
//				}
//			}
//
//			s.Properties[p.Name] = prop
//		}
//
//		p := Parameter{
//			In:     "body",
//			Name:   "body",
//			Schema: s,
//		}
//		parameters = append(parameters, p)
//	}
//
//	md := MethodDetail{
//		Summary:     method.Summary,
//		Description: method.Description,
//		Produces:    []string{"application/json"},
//		Parameters:  parameters,
//	}
//
//	return &md
//}
