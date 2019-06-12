package identity

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// parsedTemplates is a sequence of chunks to be rendered in order
type CompiledTemplate struct {
	chunks []*chunk
}

// chunk holds a function that will render entity/group info into a string
// appropriate for the matching template parameter.
type chunk struct {
	renderer func(*Entity, []*Group) (string, error)
	str      string
}

// paramMatcher defines a regex to match and possibly capture text from a
// template parameter, not including the outer braces.
type paramMatcher struct {
	// regex to match against, with optional capture groups
	pattern *regexp.Regexp

	// handler to be called during runtime to process this parameter with the
	// given entity and groups. The returned string must be a valid JSON string,
	// array, or object. captures will contain any of the parameterized elements
	// that were captured during template parsing. e.g.
	//
	// identity.entity.aliases.aws_123.metadata.region
	//                         ^^^^^^^          ^^^^^^
	//                       captures[0]       captures[1]
	handler func(entity *Entity, groups []*Group, captures []string) (string, error)
}

// parameterRE is used to locate all potential template parameters (i.e. anything
// within {{...}} with no spaces in the string.
var parameterRE = regexp.MustCompile(`"{{(\S+)}}"`)

// patterns is the set of all supported template parameters, along with their handlers.
// The pattern is a regex, but definitions use a small helper to simplify the pattern
// and ensure consistency in how string are matches, that the regex is bounded at the
// beginning and end, etc.
var patterns = []paramMatcher{
	{
		pattern: regexp.MustCompile(regexify("identity.entity.id")),
		handler: func(e *Entity, groups []*Group, v []string) (string, error) {
			return quote(e.ID), nil
		},
	},
	{
		pattern: regexp.MustCompile(regexify("identity.entity.name")),
		handler: func(e *Entity, groups []*Group, v []string) (string, error) {
			return quote(e.Name), nil
		},
	},
	{
		pattern: regexp.MustCompile(regexify("identity.entity.metadata")),
		handler: func(e *Entity, groups []*Group, v []string) (string, error) {
			d, err := json.Marshal(e.Metadata)
			if err == nil {
				return string(d), nil
			}
			return `{}`, nil
		},
	},
	{
		pattern: regexp.MustCompile(regexify(`identity.entity.metadata.<param>`)),
		handler: func(e *Entity, groups []*Group, v []string) (string, error) {
			return quote(e.Metadata[v[0]]), nil
		},
	},
	{
		pattern: regexp.MustCompile(regexify(`identity.entity.aliases.<param>.metadata.<param>`)),
		handler: func(e *Entity, groups []*Group, v []string) (string, error) {
			name, key := v[0], v[1]
			for _, alias := range e.Aliases {
				if alias.Name == name {
					return quote(alias.Metadata[key]), nil
				}
			}
			return quote(""), nil
		},
	},
	{
		pattern: regexp.MustCompile(regexify(`identity.entity.group_names`)),
		handler: func(e *Entity, groups []*Group, v []string) (string, error) {
			return groupsToArray(groups, "name"), nil
		},
	},
	{
		pattern: regexp.MustCompile(regexify(`identity.entity.group_ids`)),
		handler: func(e *Entity, groups []*Group, v []string) (string, error) {
			return groupsToArray(groups, "id"), nil
		},
	},
}

// NewCompiledTemplate will process a template string into a CompiledTemplate
// object, appropriate for populating with entity and group information. Template
// definitions should be valid JSON. It is an error if either the template isn't
// JSON, or a given template parameter isn't supported.
func NewCompiledTemplate(template string) (*CompiledTemplate, error) {
	var pt CompiledTemplate
	var tmp map[string]interface{}

	// Even before being rendered, templates should be valid JSON. Check that
	// now so we can return a descriptive errors if necessary.
	err := json.Unmarshal([]byte(template), &tmp)
	if err != nil {
		return nil, err
	}

	// Find all possible parameters {{...something...}}. matches will be a list
	// of 4-element slices provides character indices of both the entire match
	// start and end (m[0], m[1]), and the element within braces (m[2], m[3]).
	matches := parameterRE.FindAllStringSubmatchIndex(template, -1)

	// idx will point to the current character offset in the template
	idx := 0

	for _, m := range matches {
		// Add a chunkRenderer of static text from out current pointer to the start of the
		// next match.
		pt.chunks = append(pt.chunks, &chunk{
			str: template[idx:m[0]],
		})

		param := template[m[2]:m[3]]

		// Search parameter pattern looking for a match. If one is found, create
		// create a dynamic chunkRenderer using the handler for that parameter, closed
		// over with the parameter string(s) found for this template.
		var c *chunk
		for _, p := range patterns {
			// Test for a match, retaining any captures. For example:
			//
			// identity.entity.aliases.<mount_accessor>.metadata.<key>
			//                              [1]                   [2]
			// |----------------------[0]-----------------------------|
			submatches := p.pattern.FindStringSubmatch(param)

			if len(submatches) > 0 {
				handler := p.handler
				f := func(entity *Entity, groups []*Group) (string, error) {
					return handler(entity, groups, submatches[1:])
				}
				c = &chunk{renderer: f}
				break
			}
		}

		// Failing to match, just output the original string, including braces
		if c == nil {
			return nil, fmt.Errorf("invalid template parameter %q", param)
			//c = &chunk{str: template[m[0]:m[1]]}
		}
		pt.chunks = append(pt.chunks, c)

		// Advance index to the end of the entire match
		idx = m[1]
	}

	// Add remainder of template string
	pt.chunks = append(pt.chunks, &chunk{
		str: template[idx:],
	})

	return &pt, nil
}

func (t *CompiledTemplate) Render(entity *Entity, groups []*Group) (string, error) {
	var out strings.Builder

	for _, c := range t.chunks {
		result, err := c.Render(entity, groups)
		if err != nil {
			return "", err
		}
		out.WriteString(result)
	}

	return out.String(), nil
}

// Render generates a string version from the entity and groups.
// This will invoke the wrapped handler that was matched during template parsing,
// or return a fixed string for static text.
func (dc *chunk) Render(entity *Entity, groups []*Group) (string, error) {
	if dc.renderer != nil {
		return dc.renderer(entity, groups)
	} else {
		return dc.str, nil
	}
}

// regexify creates a regex from a simpler, more readable pattern.
func regexify(s string) string {
	s = strings.ReplaceAll(s, ".", `\.`)

	// TODO: named parameters might be better than <param>
	s = strings.ReplaceAll(s, "<param>", `([^\s.]+)`)

	return "^" + s + "$"
}

// groupsToArray is a helper to extract either the ID or Name from
// a list of groups into a JSON array.
func groupsToArray(groups []*Group, element string) string {
	var out strings.Builder

	groupsLen := len(groups)

	out.WriteString("[")
	for i, g := range groups {
		var v string
		switch element {
		case "name":
			v = g.Name
		case "id":
			v = g.ID
		}
		out.WriteString(quote(v))
		if i < groupsLen-1 {
			out.WriteString(",")
		}
	}
	out.WriteString("]")

	return out.String()
}

func quote(s string) string {
	return fmt.Sprintf(`"%s"`, s)
}
