package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/kalafut/q"
)

var re = regexp.MustCompile(`"{{(\S+)}}"`)

const (
	text = iota
	handler
)

type parsedTemplate struct {
	chunks []chunk
}

type chunk struct {
	chunkType int

	value    string
	dflt     string
	variable string
	matcher  hFunc
}

func (t *parsedTemplate) addText(s string) {
	c := chunk{
		chunkType: text,
		value:     s,
	}
	t.chunks = append(t.chunks, c)
}

func (t *parsedTemplate) addMatcher(f hFunc) {
	t.chunks = append(t.chunks, chunk{
		chunkType: handler,
		matcher:   f,
	})
}

func (t *parsedTemplate) render(entity *identity.Entity) string {
	var out strings.Builder

	for _, c := range t.chunks {
		switch c.chunkType {
		case text:
			out.WriteString(c.value)
		case handler:
			out.WriteString(c.matcher(entity))
		}
	}

	return out.String()
}

func ABC(tpl string, e *identity.Entity, groups []*identity.Group) (string, error) {
	var out map[string]interface{}
	var pt parsedTemplate

	var compact bytes.Buffer
	err := json.Compact(&compact, []byte(tpl))

	err = json.Unmarshal([]byte(tpl), &out)

	matches := re.FindAllStringSubmatchIndex(tpl, -1)

	i := 0
	for _, m := range matches {
		pt.addText(tpl[i:m[0]])

		v := tpl[m[2]:m[3]]
		q.Q()

		var f hFunc
		for _, p := range patterns {
			m := p.pattern.FindStringSubmatch(v)
			if len(m) > 0 {
				//q.Q(v, p.pattern)
				f = func(entity *identity.Entity) string {
					return p.handler(entity, groups, m[1:])
				}
				break
			}
		}

		if f != nil {
			pt.addMatcher(f)
		} else {
			pt.addText(v)
		}

		i = m[1]
	}
	pt.addText(tpl[i:])

	result := pt.render(e)
	q.Q(result)

	return result, err
}

/*
identity.entity.id
identity.entity.name
identity.entity.metadata
identity.entity.metadata.<key>
identity.entity.group_ids
identity.entity.group_names
identity.entity.aliases.<mount_accessor>
identity.entity.aliases.<mount_accessor>.name
identity.entity.aliases.<mount_accessor>.metadata
identity.entity.aliases.<mount_accessor>.metadata.<key>
identity.groups.<group id>.id
identity.groups.<group id>.name
identity.groups.<group id>.metadata
identity.groups.<group id>.metadata.<key>
*/

type hFunc func(*identity.Entity) string
type handlerFunc func(*identity.Entity, []*identity.Group, []string) string

type matcher struct {
	pattern *regexp.Regexp
	handler handlerFunc
}

// TODO: maybe get rid of this
func quote(s string) string {
	return fmt.Sprintf(`"%s"`, s)
}

// TODO: named parameters would be better than <param>
func reFmt(s string) string {
	s = strings.ReplaceAll(s, ".", `\.`)
	s = strings.ReplaceAll(s, "<param>", `([^\s.]+)`)
	return "^" + s + "$"
}

var patterns = []matcher{
	{
		pattern: regexp.MustCompile(reFmt("identity.entity.id")),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) string {
			return quote(e.ID)
		},
	},
	{
		pattern: regexp.MustCompile(reFmt("identity.entity.name")),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) string {
			return quote(e.Name)
		},
	},
	{
		pattern: regexp.MustCompile(reFmt("identity.entity.metadata")),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) string {
			d, err := json.Marshal(e.Metadata)
			if err == nil {
				return string(d)
			}
			return `{}`
		},
	},
	{
		pattern: regexp.MustCompile(reFmt(`identity.entity.metadata.<param>`)),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) string {
			return quote(e.Metadata[v[0]])
		},
	},
	{
		pattern: regexp.MustCompile(reFmt(`identity.entity.aliases.<param>.metadata.<param>`)),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) string {
			name, key := v[0], v[1]
			for _, alias := range e.Aliases {
				if alias.Name == name {
					return quote(alias.Metadata[key])
				}
			}
			return quote("")
		},
	},
	{
		pattern: regexp.MustCompile(reFmt(`identity.entity.group_names`)),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) string {
			return groupsToList(groups, "name")
		},
	},
	{
		pattern: regexp.MustCompile(reFmt(`identity.entity.group_ids`)),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) string {
			return groupsToList(groups, "id")
		},
	},
}

func groupsToList(groups []*identity.Group, element string) string {
	var out strings.Builder

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
		if i < len(groups)-1 {
			out.WriteString(",")
		}
	}
	out.WriteString("]")

	return out.String()
}

func classifyParameters(entity *identity.Entity, groups []*identity.Group, s string) (result string, err error) {
	p := findPattern(s)
	if p != nil {
		m := p.pattern.FindStringSubmatch(s)
		result = p.handler(entity, groups, m[1:])
	}
	return result, nil
}

func findPattern(s string) *matcher {
	for _, p := range patterns {
		m := p.pattern.FindStringSubmatch(s)
		if len(m) > 0 {
			return &p
		}
	}

	return nil
}
