package pkicli

import (
	"fmt"
	"strings"
)

type mapStringAny map[string]interface{}

type Params struct {
	params       mapStringAny
	missingKeys  []string
	invalidTypes []string
}

func newParams(m map[string]interface{}) *Params {
	p := make(map[string]interface{})
	for k, v := range m {
		p[k] = v
	}

	return &Params{params: p}
}

func (p *Params) clone() *Params {
	return newParams(p.params)
}

func (p *Params) put(k string, v interface{}) *Params {
	p.params[k] = v
	return p
}

func (p *Params) putDefault(k string, v interface{}) {
	if _, ok := p.params[k]; !ok {
		p.params[k] = v
	}
}

func (p *Params) privatePop(k string) (interface{}, bool) {
	v, ok := p.params[k]
	if !ok {
		return "", false
	}

	delete(p.params, k)
	return v, true
}

func (p *Params) pop(k string) string {
	if v, ok := p.privatePop(k); ok {
		if s, ok := v.(string); ok {
			return s
		}
		p.invalidTypes = append(p.invalidTypes, k)
		return ""
	}

	p.appendMissing(k)
	return ""
}

func (p *Params) appendMissing(k string) {
	p.missingKeys = append(p.missingKeys, k)
}

func (p *Params) popDefault(k, d string) string {
	if v, ok := p.privatePop(k); ok {
		if s, ok := v.(string); ok {
			return s
		}
		p.invalidTypes = append(p.invalidTypes, k)
		return ""
	}

	return d
}

func (p *Params) data() mapStringAny {
	data := mapStringAny{}
	for k, v := range p.params {
		if !strings.HasPrefix(k, "_") {
			data[k] = v
		}
	}
	return data
}

func (p *Params) error() error {
	if len(p.missingKeys) > 0 {
		return fmt.Errorf("missing parameters: %s", strings.Join(p.missingKeys, ","))
	}
	if len(p.invalidTypes) > 0 {
		return fmt.Errorf("invalid type for parameters: %s", strings.Join(p.invalidTypes, ","))
	}
	return nil
}

func (p *Params) hasAny(keys ...string) bool {
	for _, k := range keys {
		if _, ok := p.params[k]; ok {
			return true
		}
	}

	return false
}

func (p *Params) hasAll(keys ...string) bool {
	for _, k := range keys {
		if _, ok := p.params[k]; !ok {
			return false
		}
	}

	return true
}

func (p *Params) require(keys ...string) {
	for _, k := range keys {
		if !p.hasAny(k) {
			p.appendMissing(k)
		}
	}
}
