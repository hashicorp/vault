package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/mitchellh/reflectwalk"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/kalafut/q"
)

var re = regexp.MustCompile(`"{{(\S+)}}"`)

const (
	text = iota
	str
	obj
)

type parsedTemplate struct {
	chunks []chunk
}

type chunk struct {
	chunkType int

	value    string
	dflt     string
	variable string
}

func (t *parsedTemplate) addText(s string) {
	c := chunk{
		chunkType: text,
		value:     s,
	}
	t.chunks = append(t.chunks, c)
}

func (t *parsedTemplate) addStringParam(name string) {
	t.chunks = append(t.chunks, chunk{
		chunkType: str,
		variable:  name,
		dflt:      `""`,
	})
}

func (t *parsedTemplate) render(entity identity.Entity) string {
	var out strings.Builder

	for _, c := range t.chunks {
		switch c.chunkType {
		case text:
			out.WriteString(c.value)
		case str:
			if c.variable == "identity.entity.id" {
				out.WriteString(fmt.Sprintf(`"%s"`, entity.ID))
			} else {
				out.WriteString(c.dflt)
			}
		}
	}

	return out.String()
}

func ABC(tpl string, e identity.Entity) (map[string]interface{}, error) {
	var out map[string]interface{}
	var pt parsedTemplate

	var compact bytes.Buffer
	err := json.Compact(&compact, []byte(tpl))

	//err := jsonutil.DecodeJSON([]byte(tpl), &out)
	err = json.Unmarshal([]byte(tpl), &out)
	q.Q(out)

	matches := re.FindAllStringSubmatchIndex(tpl, -1)

	i := 0
	for _, m := range matches {
		pt.addText(tpl[i:m[0]])

		v := tpl[m[2]:m[3]]
		if v == "identity.entity.id" {
			pt.addStringParam("identity.entity.id")
		} else {
			pt.addText(v)
		}

		i = m[1]
	}
	pt.addText(tpl[i:])

	ent := identity.Entity{
		ID: "abc-123",
	}
	q.Q(pt.render(ent))

	//walkTest()

	return out, err
}

type walker struct{}

func (w walker) Map(m reflect.Value) error {
	q.Q(m)
	return nil
}

func (w walker) MapElem(m, k, v reflect.Value) error {
	if v.Kind() == reflect.String {
		q.Q(m, k.String(), v.String())
	}
	return nil
}

var _ reflectwalk.MapWalker = (*walker)(nil)

func walkTest() {
	exp := map[string]interface{}{
		"basic": float64(42),
	}
	w := walker{}

	q.Q(reflectwalk.Walk(exp, w))
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

type matcher struct {
	pattern *regexp.Regexp
	handler func(identity.Entity, []string) string
}

var patterns = []matcher{
	{
		pattern: regexp.MustCompile(`^identity\.entity\.id$`),
		handler: func(e identity.Entity, v []string) string {
			return e.ID
		},
	},
	{
		pattern: regexp.MustCompile(`^identity\.entity\.name$`),
		handler: func(e identity.Entity, v []string) string {
			return e.Name
		},
	},
	{
		pattern: regexp.MustCompile(`^identity\.entity\.metadata$`),
		handler: func(e identity.Entity, v []string) string {
			d, err := json.Marshal(e.Metadata)
			if err == nil {
				return string(d)
			}
			return ""
		},
	},
	{
		pattern: regexp.MustCompile(`^identity\.entity\.metadata\.(\S+)$`),
		handler: func(e identity.Entity, v []string) string {
			return e.Metadata[v[0]]
		},
	},
	{
		pattern: regexp.MustCompile(`^identity\.entity\.aliases\.(\S+)\.metadata\.(\S+)$`),
		handler: func(e identity.Entity, v []string) string {
			name, key := v[0], v[1]
			for _, alias := range e.Aliases {
				if alias.Name == name {
					return alias.Metadata[key]
				}
			}
			return ""
		},
	},
}

func classifyParameters(entity identity.Entity, s string) (result string, err error) {
	for _, p := range patterns {
		m := p.pattern.FindStringSubmatch(s)
		if len(m) > 0 {
			result = p.handler(entity, m[1:])
			break
		}
	}

	return result, err
}
