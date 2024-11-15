package dependency

import (
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"

	capi "github.com/hashicorp/consul/api"
)

const exportedServicesEndpointLabel = "list.exportedServices"

// Ensure implements
var _ Dependency = (*ListExportedServicesQuery)(nil)

// ListExportedServicesQuery is the representation of a requested exported services
// dependency from inside a template.
type ListExportedServicesQuery struct {
	stopCh    chan struct{}
	partition string
}

type ExportedService struct {
	// Name of the service
	Service string

	// Partition of the service
	Partition string

	// Namespace of the service
	Namespace string

	// Consumers is a list of downstream consumers of the service.
	Consumers ResolvedConsumers
}

type ResolvedConsumers struct {
	Peers      []string
	Partitions []string
}

func fromConsulExportedService(svc capi.ResolvedExportedService) ExportedService {
	exportedService := ExportedService{
		Service: svc.Service,
		Consumers: ResolvedConsumers{
			Partitions: []string{},
			Peers:      []string{},
		},
	}

	if len(svc.Consumers.Partitions) > 0 {
		exportedService.Consumers.Partitions = slices.Clone(svc.Consumers.Partitions)
	}

	if len(svc.Consumers.Peers) > 0 {
		exportedService.Consumers.Peers = slices.Clone(svc.Consumers.Peers)
	}

	return exportedService
}

// NewListExportedServicesQuery parses a string of the format @dc.
func NewListExportedServicesQuery(s string) (*ListExportedServicesQuery, error) {
	return &ListExportedServicesQuery{
		stopCh:    make(chan struct{}),
		partition: s,
	}, nil
}

func (c *ListExportedServicesQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-c.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(&QueryOptions{
		ConsulPartition: c.partition,
	})

	log.Printf("[TRACE] %s: GET %s", c, &url.URL{
		Path:     "/v1/exported-services",
		RawQuery: opts.String(),
	})

	consulExportedServices, qm, err := clients.Consul().ExportedServices(opts.ToConsulOpts())
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", c.String(), err)
	}

	exportedServices := make([]ExportedService, 0, len(consulExportedServices))
	for _, exportedService := range consulExportedServices {
		exportedServices = append(exportedServices, fromConsulExportedService(exportedService))
	}

	log.Printf("[TRACE] %s: returned %d results", c, len(exportedServices))

	slices.SortStableFunc(exportedServices, func(i, j ExportedService) int {
		return strings.Compare(i.Service, j.Service)
	})

	rm := &ResponseMetadata{
		LastContact: qm.LastContact,
		LastIndex:   qm.LastIndex,
	}

	return exportedServices, rm, nil
}

// CanShare returns if this dependency is shareable when consul-template is running in de-duplication mode.
func (c *ListExportedServicesQuery) CanShare() bool {
	return true
}

func (c *ListExportedServicesQuery) String() string {
	return fmt.Sprintf("%s(%s)", exportedServicesEndpointLabel, c.partition)
}

func (c *ListExportedServicesQuery) Stop() {
	close(c.stopCh)
}

func (c *ListExportedServicesQuery) Type() Type {
	return TypeConsul
}
