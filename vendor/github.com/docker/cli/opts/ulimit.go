package opts

import (
	"fmt"
	"sort"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-units"
)

// UlimitOpt defines a map of Ulimits
type UlimitOpt struct {
	values *map[string]*container.Ulimit
}

// NewUlimitOpt creates a new UlimitOpt. Ulimits are not validated.
func NewUlimitOpt(ref *map[string]*container.Ulimit) *UlimitOpt {
	// TODO(thaJeztah): why do we need a map with pointers here?
	if ref == nil {
		ref = &map[string]*container.Ulimit{}
	}
	return &UlimitOpt{ref}
}

// Set validates a Ulimit and sets its name as a key in UlimitOpt
func (o *UlimitOpt) Set(val string) error {
	// FIXME(thaJeztah): these functions also need to be moved over from go-units.
	l, err := units.ParseUlimit(val)
	if err != nil {
		return err
	}

	(*o.values)[l.Name] = l

	return nil
}

// String returns Ulimit values as a string. Values are sorted by name.
func (o *UlimitOpt) String() string {
	out := make([]string, 0, len(*o.values))
	for _, v := range *o.values {
		out = append(out, v.String())
	}
	sort.Strings(out)
	return fmt.Sprintf("%v", out)
}

// GetList returns a slice of pointers to Ulimits. Values are sorted by name.
func (o *UlimitOpt) GetList() []*container.Ulimit {
	ulimits := make([]*container.Ulimit, 0, len(*o.values))
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
