package identity

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// chunkRenderer can convert and an entity+group list into JSON
type chunkRenderer interface {
	Render(*Entity, []*Group) (string, error)
}

// staticChunk holds a string that will be rendered without change
type staticChunk struct {
	str string
}

// Render generates a string version from the entity and groups.
// For staticChunk, this will just be a fixed string.
func (sc *staticChunk) Render(*Entity, []*Group) (string, error) {
	return sc.str, nil
}

// dynamicChunk holds a function that will render entity/group info into a string
// appropriate for the matching template parameter.
type dynamicChunk struct {
	renderer func(*Entity, []*Group) (string, error)
}

// Render generates a string version from the entity and groups.
// For dynamicChunk, this will invoked the wrapped handler that was
// matched during template parsing.
func (dc *dynamicChunk) Render(entity *Entity, groups []*Group) (string, error) {
	return dc.renderer(entity, groups)
}

// parsedTemplates is a sequence of chunks to be rendered in order
type CompiledTemplate struct {
	chunks []chunkRenderer
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
		pt.chunks = append(pt.chunks, &staticChunk{
			str: template[idx:m[0]],
		})

		param := template[m[2]:m[3]]

		// Search parameter pattern looking for a match. If one is found, create
		// create a dynamic chunkRenderer using the handler for that parameter, closed
		// over with the parameter string(s) found for this template.
		var c chunkRenderer
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
				c = &dynamicChunk{renderer: f}
				break
			}
		}

		// Failing to match, just output the original string, including braces
		if c == nil {
			c = &staticChunk{str: template[m[0]:m[1]]}
		}
		pt.chunks = append(pt.chunks, c)

		// Advance index to the end of the entire match
		idx = m[1]
	}

	// Add remainder of template string
	pt.chunks = append(pt.chunks, &staticChunk{
		str: template[idx:],
	})

	return &pt, nil
}

type paramMatcher struct {
	pattern *regexp.Regexp
	handler func(*Entity, []*Group, []string) (string, error)
}

var parameterRE = regexp.MustCompile(`"{{(\S+)}}"`)
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
