package dependency

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	// Ensure implements
	_ Dependency = (*KVKeysQuery)(nil)

	// KVKeysQueryRe is the regular expression to use.
	KVKeysQueryRe = regexp.MustCompile(`\A` + prefixRe + dcRe + `\z`)
)

// KVKeysQuery queries the KV store for a single key.
type KVKeysQuery struct {
	stopCh chan struct{}

	dc     string
	prefix string
}

// NewKVKeysQuery parses a string into a dependency.
func NewKVKeysQuery(s string) (*KVKeysQuery, error) {
	if s != "" && !KVKeysQueryRe.MatchString(s) {
		return nil, fmt.Errorf("kv.keys: invalid format: %q", s)
	}

	m := regexpMatch(KVKeysQueryRe, s)
	return &KVKeysQuery{
		stopCh: make(chan struct{}, 1),
		dc:     m["dc"],
		prefix: m["prefix"],
	}, nil
}

// Fetch queries the Consul API defined by the given client.
func (d *KVKeysQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(&QueryOptions{
		Datacenter: d.dc,
	})

	log.Printf("[TRACE] %s: GET %s", d, &url.URL{
		Path:     "/v1/kv/" + d.prefix,
		RawQuery: opts.String(),
	})

	list, qm, err := clients.Consul().KV().Keys(d.prefix, "", opts.ToConsulOpts())
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	keys := make([]string, len(list))
	for i, v := range list {
		v = strings.TrimPrefix(v, d.prefix)
		v = strings.TrimLeft(v, "/")
		keys[i] = v
	}

	log.Printf("[TRACE] %s: returned %d results", d, len(list))

	rm := &ResponseMetadata{
		LastIndex:   qm.LastIndex,
		LastContact: qm.LastContact,
	}

	return keys, rm, nil
}

// CanShare returns a boolean if this dependency is shareable.
func (d *KVKeysQuery) CanShare() bool {
	return true
}

// String returns the human-friendly version of this dependency.
func (d *KVKeysQuery) String() string {
	prefix := d.prefix
	if d.dc != "" {
		prefix = prefix + "@" + d.dc
	}
	return fmt.Sprintf("kv.keys(%s)", prefix)
}

// Stop halts the dependency's fetch function.
func (d *KVKeysQuery) Stop() {
	close(d.stopCh)
}

// Type returns the type of this dependency.
func (d *KVKeysQuery) Type() Type {
	return TypeConsul
}
