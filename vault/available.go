package vault

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const (
	listAvailableAuditsKey = "audits"
	listAvailableAuthsKey  = "auths"
	listAvailableMountsKey = "mounts"
)

// listAvailable returns the list of installed type by name. This can list
// the installed audits devices, auth methods, secrets engines, etc.
func (c *Core) listAvailable(_ context.Context, typ string) ([]string, error) {
	var rv reflect.Value

	typ = strings.TrimSuffix(typ, "/")
	switch typ {
	case listAvailableAuditsKey:
		rv = reflect.ValueOf(c.auditBackends)
	case listAvailableAuthsKey:
		rv = reflect.ValueOf(c.credentialBackends)
	case listAvailableMountsKey:
		rv = reflect.ValueOf(c.logicalBackends)
	default:
		return nil, fmt.Errorf("unknown type: %s", typ)
	}

	// This should never happen
	if rv.Kind() != reflect.Map || rv.Type().Key().Kind() != reflect.String {
		return nil, fmt.Errorf("not a string-keyed map: %s", typ)
	}

	names := make([]string, len(rv.MapKeys()))
	for i, v := range rv.MapKeys() {
		names[i] = v.String()
	}
	sort.Strings(names)
	return names, nil
}
