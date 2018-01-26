package framework

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/vault/logical"
)

// PolicyMap is a specialization of PathMap that expects the values to
// be lists of policies. This assists in querying and loading policies
// from the PathMap.
type PolicyMap struct {
	PathMap

	DefaultKey string
	PolicyKey  string
}

func (p *PolicyMap) Policies(ctx context.Context, s logical.Storage, names ...string) ([]string, error) {
	policyKey := "value"
	if p.PolicyKey != "" {
		policyKey = p.PolicyKey
	}

	if p.DefaultKey != "" {
		newNames := make([]string, len(names)+1)
		newNames[0] = p.DefaultKey
		copy(newNames[1:], names)
		names = newNames
	}

	set := make(map[string]struct{})
	for _, name := range names {
		v, err := p.Get(ctx, s, name)
		if err != nil {
			return nil, err
		}

		valuesRaw, ok := v[policyKey]
		if !ok {
			continue
		}

		values, ok := valuesRaw.(string)
		if !ok {
			continue
		}

		for _, p := range strings.Split(values, ",") {
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
