package opts

import (
	"fmt"
	"sort"

	"github.com/docker/go-units"
)

// UlimitOpt defines a map of Ulimits
type UlimitOpt struct {
	values *map[string]*units.Ulimit
}

// NewUlimitOpt creates a new UlimitOpt. Ulimits are not validated.
func NewUlimitOpt(ref *map[string]*units.Ulimit) *UlimitOpt {
	if ref == nil {
		ref = &map[string]*units.Ulimit{}
	}
	return &UlimitOpt{ref}
}

// Set validates a Ulimit and sets its name as a key in UlimitOpt
func (o *UlimitOpt) Set(val string) error {
	l, err := units.ParseUlimit(val)
	if err != nil {
		return err
	}

	(*o.values)[l.Name] = l

	return nil
}

// String returns Ulimit values as a string. Values are sorted by name.
func (o *UlimitOpt) String() string {
	var out []string
	for _, v := range *o.values {
		out = append(out, v.String())
	}
	sort.Strings(out)
	return fmt.Sprintf("%v", out)
}

// GetList returns a slice of pointers to Ulimits. Values are sorted by name.
func (o *UlimitOpt) GetList() []*units.Ulimit {
	var ulimits []*units.Ulimit
	for _, v := range *o.values {
		ulimits = append(ulimits, v)
	}
	sort.SliceStable(ulimits, func(i, j int) bool {
		return ulimits[i].Name < ulimits[j].Name
	})
	return ulimits
}

// Type returns the option type
func (o *UlimitOpt) Type() string {
	return "ulimit"
}
