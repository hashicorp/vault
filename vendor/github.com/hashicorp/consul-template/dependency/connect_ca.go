package dependency

import (
	"log"
	"net/url"

	"github.com/pkg/errors"
)

var (
	// Ensure implements
	_ Dependency = (*ConnectCAQuery)(nil)
)

type ConnectCAQuery struct {
	stopCh chan struct{}
}

func NewConnectCAQuery() *ConnectCAQuery {
	return &ConnectCAQuery{
		stopCh: make(chan struct{}, 1),
	}
}

func (d *ConnectCAQuery) Fetch(clients *ClientSet, opts *QueryOptions) (
	interface{}, *ResponseMetadata, error,
) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(nil)
	log.Printf("[TRACE] %s: GET %s", d, &url.URL{
		Path:     "/v1/agent/connect/ca/roots",
		RawQuery: opts.String(),
	})

	certs, md, err := clients.Consul().Agent().ConnectCARoots(
		opts.ToConsulOpts())
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	log.Printf("[TRACE] %s: returned %d results", d, len(certs.Roots))
	log.Printf("[TRACE] %s: %#v ", d, md)

	rm := &ResponseMetadata{
		LastIndex:   md.LastIndex,
		LastContact: md.LastContact,
		Block:       true,
	}

	return certs.Roots, rm, nil
}

func (d *ConnectCAQuery) Stop() {
	close(d.stopCh)
}

func (d *ConnectCAQuery) CanShare() bool {
	return false
}

func (d *ConnectCAQuery) Type() Type {
	return TypeConsul
}

func (d *ConnectCAQuery) String() string {
	return "connect.caroots"
}
