package dependency

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

// Ensure implements
var (
	_ Dependency = (*ListPartitionsQuery)(nil)

	// ListPartitionsQuerySleepTime is the amount of time to sleep between
	// queries, since the endpoint does not support blocking queries.
	ListPartitionsQuerySleepTime = DefaultNonBlockingQuerySleepTime
)

// Partition is a partition in Consul.
type Partition struct {
	Name        string
	Description string
}

// ListPartitionsQuery is the representation of a requested partitions
// dependency from inside a template.
type ListPartitionsQuery struct {
	stopCh chan struct{}
}

func NewListPartitionsQuery() (*ListPartitionsQuery, error) {
	return &ListPartitionsQuery{
		stopCh: make(chan struct{}, 1),
	}, nil
}

func (c *ListPartitionsQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	opts = opts.Merge(&QueryOptions{})

	log.Printf("[TRACE] %s: GET %s", c, &url.URL{
		Path:     "/v1/partitions",
		RawQuery: opts.String(),
	})

	// This is certainly not elegant, but the partitions endpoint does not support
	// blocking queries, so we are going to "fake it until we make it". When we
	// first query, the LastIndex will be "0", meaning we should immediately
	// return data, but future calls will include a LastIndex. If we have a
	// LastIndex in the query metadata, sleep for 15 seconds before asking Consul
	// again.
	//
	// This is probably okay given the frequency in which partitions actually
	// change, but is technically not edge-triggering.
	if opts.WaitIndex != 0 {
		log.Printf("[TRACE] %s: long polling for %s", c, ListPartitionsQuerySleepTime)

		select {
		case <-c.stopCh:
			return nil, nil, ErrStopped
		case <-time.After(ListPartitionsQuerySleepTime):
		}
	}

	partitions, _, err := clients.Consul().Partitions().List(context.Background(), opts.ToConsulOpts())
	if err != nil {
		if strings.Contains(err.Error(), "Invalid URL path") {
			return nil, nil, fmt.Errorf("%s: Partitions are an enterprise feature: %w", c.String(), err)
		}

		return nil, nil, fmt.Errorf("%s: %w", c.String(), err)
	}

	log.Printf("[TRACE] %s: returned %d results", c, len(partitions))

	slices.SortFunc(partitions, func(i, j *api.Partition) int {
		return strings.Compare(i.Name, j.Name)
	})

	resp := []*Partition{}
	for _, partition := range partitions {
		if partition != nil {
			resp = append(resp, &Partition{
				Name:        partition.Name,
				Description: partition.Description,
			})
		}
	}

	// Use respWithMetadata which always increments LastIndex and results
	// in fetching new data for endpoints that don't support blocking queries
	return respWithMetadata(resp)
}

// CanShare returns if this dependency is shareable when consul-template is running in de-duplication mode.
func (c *ListPartitionsQuery) CanShare() bool {
	return true
}

func (c *ListPartitionsQuery) String() string {
	return "list.partitions"
}

func (c *ListPartitionsQuery) Stop() {
	close(c.stopCh)
}

func (c *ListPartitionsQuery) Type() Type {
	return TypeConsul
}
