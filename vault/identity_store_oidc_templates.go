package vault

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/helper/identity"
)

var re = regexp.MustCompile(`"{{(\S+)}}"`)

const (
	text = iota
	handler
)

type parsedTemplate struct {
	chunks []chunk
}

type chunk interface {
	Render(*identity.Entity, []*identity.Group) (string, error)
}

type staticChunk struct {
	value string
}

func (sc *staticChunk) Render(*identity.Entity, []*identity.Group) (string, error) {
	return sc.value, nil
}

type dynamicChunk struct {
	matcher hFunc
}

func (dc *dynamicChunk) Render(entity *identity.Entity, groups []*identity.Group) (string, error) {
	return dc.matcher(entity, groups)
}

//func (t *parsedTemplate) addText(s string) {
//	c := chunk{
//		chunkType: text,
//		value:     s,
//	}
//	t.chunks = append(t.chunks, c)
//}
//
//func (t *parsedTemplate) addMatcher(f hFunc) {
//	t.chunks = append(t.chunks, chunk{
//		chunkType: handler,
//		matcher:   f,
//	})
//}

//func (t *parsedTemplate) addText(s string) {
//	c := chunk{
//		chunkType: text,
//		value:     s,
//	}
//	t.chunks = append(t.chunks, c)
//}
//
//func (t *parsedTemplate) addMatcher(f hFunc) {
//	t.chunks = append(t.chunks, chunk{
//		chunkType: handler,
//		matcher:   f,
//	})
//}

func (t *parsedTemplate) Render(entity *identity.Entity, groups []*identity.Group) (string, error) {
	var out strings.Builder

	for _, c := range t.chunks {
		result, err := c.Render(entity, groups)
		if err != nil {
			return "", err
		}
		out.WriteString(result)

		//switch c.chunkType {
		//case text:
		//	out.WriteString(c.value)
		//case handler:
		//	result, err := c.matcher(entity, groups)
		//	if err != nil {
		//		return "", err
		//	}
		//	out.WriteString(result)
		//}
	}

	return out.String(), nil
}

func CompileTemplate(tpl string) (parsedTemplate, error) {
	var out map[string]interface{}
	var pt parsedTemplate

	err := json.Unmarshal([]byte(tpl), &out)
	if err != nil {
		return pt, err
	}

	matches := re.FindAllStringSubmatchIndex(tpl, -1)

	i := 0
	for _, m := range matches {
		//pt.addText(tpl[i:m[0]])
		pt.chunks = append(pt.chunks, &staticChunk{tpl[i:m[0]]})

		v := tpl[m[2]:m[3]]

		var f hFunc
		for _, p := range patterns {
			m := p.pattern.FindStringSubmatch(v)
			if len(m) > 0 {
				f = func(entity *identity.Entity, groups []*identity.Group) (string, error) {
					return p.handler(entity, groups, m[1:])
				}
				break
			}
		}

		if f != nil {
			pt.chunks = append(pt.chunks, &dynamicChunk{f})
			//pt.addMatcher(f)
		} else {
			pt.chunks = append(pt.chunks, &staticChunk{v})
			//f := func(entity *identity.Entity, groups []*identity.Group) (string, error) {
			//	return staticTextHandler(nil, nil, []string{v})
			//}
			//pt.addMatcher(f)
		}

		i = m[1]
	}
	//pt.addText(tpl[i:])
	pt.chunks = append(pt.chunks, &staticChunk{tpl[i:]})

	return pt, nil
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
*/

type hFunc func(*identity.Entity, []*identity.Group) (string, error)
type handlerFunc func(*identity.Entity, []*identity.Group, []string) (string, error)

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
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) (string, error) {
			return quote(e.ID), nil
		},
	},
	{
		pattern: regexp.MustCompile(reFmt("identity.entity.name")),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) (string, error) {
			return quote(e.Name), nil
		},
	},
	{
		pattern: regexp.MustCompile(reFmt("identity.entity.metadata")),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) (string, error) {
			d, err := json.Marshal(e.Metadata)
			if err == nil {
				return string(d), nil
			}
			return `{}`, nil
		},
	},
	{
		pattern: regexp.MustCompile(reFmt(`identity.entity.metadata.<param>`)),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) (string, error) {
			return quote(e.Metadata[v[0]]), nil
		},
	},
	{
		pattern: regexp.MustCompile(reFmt(`identity.entity.aliases.<param>.metadata.<param>`)),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) (string, error) {
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
		pattern: regexp.MustCompile(reFmt(`identity.entity.group_names`)),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) (string, error) {
			return groupsToList(groups, "name"), nil
		},
	},
	{
		pattern: regexp.MustCompile(reFmt(`identity.entity.group_ids`)),
		handler: func(e *identity.Entity, groups []*identity.Group, v []string) (string, error) {
			return groupsToList(groups, "id"), nil
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
		result, err = p.handler(entity, groups, m[1:])
	}
	return result, err
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
