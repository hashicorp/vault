package framework

import (
	"sort"
	"strings"

	"github.com/hashicorp/vault/logical"
)

// PolicyMap is a specialization of PathMap that expects the values to
// be lists of policies. This assists in querying and loading policies
// from the PathMap.
type PolicyMap struct {
	*PathMap

	DefaultKey string
}

func (p *PolicyMap) Policies(s logical.Storage, names ...string) ([]string, error) {
	if p.DefaultKey != "" {
		newNames := make([]string, len(names)+1)
		newNames[0] = p.DefaultKey
		copy(newNames[1:], names)
		names = newNames
	}

	set := make(map[string]struct{})
	for _, name := range names {
		v, err := p.Get(s, name)
		if err != nil {
			return nil, err
		}

		for _, p := range strings.Split(v, ",") {
			if p = strings.TrimSpace(p); p != "" {
				set[p] = struct{}{}
			}
		}
	}

	list := make([]string, 0, len(set))
	for k, _ := range set {
		list = append(list, k)
	}
	sort.Strings(list)

	return list, nil
}
