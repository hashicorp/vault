// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"net/url"
	"time"
)

// ServiceRegistration is an instance of a single allocation advertising itself
// as a named service with a specific address. Each registration is constructed
// from the job specification Service block. Whether the service is registered
// within Nomad, and therefore generates a ServiceRegistration is controlled by
// the Service.Provider parameter.
type ServiceRegistration struct {

	// ID is the unique identifier for this registration. It currently follows
	// the Consul service registration format to provide consistency between
	// the two solutions.
	ID string

	// ServiceName is the human friendly identifier for this service
	// registration.
	ServiceName string

	// Namespace represents the namespace within which this service is
	// registered.
	Namespace string

	// NodeID is Node.ID on which this service registration is currently
	// running.
	NodeID string

	// Datacenter is the DC identifier of the node as identified by
	// Node.Datacenter.
	Datacenter string

	// JobID is Job.ID and represents the job which contained the service block
	// which resulted in this service registration.
	JobID string

	// AllocID is Allocation.ID and represents the allocation within which this
	// service is running.
	AllocID string

	// Tags are determined from either Service.Tags or Service.CanaryTags and
	// help identify this service. Tags can also be used to perform lookups of
	// services depending on their state and role.
	Tags []string

	// Address is the IP address of this service registration. This information
	// comes from the client and is not guaranteed to be routable; this depends
	// on cluster network topology.
	Address string

	// Port is the port number on which this service registration is bound. It
	// is determined by a combination of factors on the client.
	Port int

	CreateIndex uint64
	ModifyIndex uint64
}

// ServiceRegistrationListStub represents all service registrations held within a
// single namespace.
type ServiceRegistrationListStub struct {

	// Namespace details the namespace in which these services have been
	// registered.
	Namespace string

	// Services is a list of services found within the namespace.
	Services []*ServiceRegistrationStub
}

// ServiceRegistrationStub is the stub object describing an individual
// namespaced service. The object is built in a manner which would allow us to
// add additional fields in the future, if we wanted.
type ServiceRegistrationStub struct {

	// ServiceName is the human friendly name for this service as specified
	// within Service.Name.
	ServiceName string

	// Tags is a list of unique tags found for this service. The list is
	// de-duplicated automatically by Nomad.
	Tags []string
}

// Services is used to query the service endpoints.
type Services struct {
	client *Client
}

// Services returns a new handle on the services endpoints.
func (c *Client) Services() *Services {
	return &Services{client: c}
}

// List can be used to list all service registrations currently stored within
// the target namespace. It returns a stub response object.
func (s *Services) List(q *QueryOptions) ([]*ServiceRegistrationListStub, *QueryMeta, error) {
	var resp []*ServiceRegistrationListStub
	qm, err := s.client.query("/v1/services", &resp, q)
	if err != nil {
		return nil, qm, err
	}
	return resp, qm, nil
}

// Get is used to return a list of service registrations whose name matches the
// specified parameter.
func (s *Services) Get(serviceName string, q *QueryOptions) ([]*ServiceRegistration, *QueryMeta, error) {
	var resp []*ServiceRegistration
	qm, err := s.client.query("/v1/service/"+url.PathEscape(serviceName), &resp, q)
	if err != nil {
		return nil, qm, err
	}
	return resp, qm, nil
}

// Delete can be used to delete an individual service registration as defined
// by its service name and service ID.
func (s *Services) Delete(serviceName, serviceID string, q *WriteOptions) (*WriteMeta, error) {
	path := fmt.Sprintf("/v1/service/%s/%s", url.PathEscape(serviceName), url.PathEscape(serviceID))
	wm, err := s.client.delete(path, nil, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// CheckRestart describes if and when a task should be restarted based on
// failing health checks.
type CheckRestart struct {
	Limit          int            `mapstructure:"limit" hcl:"limit,optional"`
	Grace          *time.Duration `mapstructure:"grace" hcl:"grace,optional"`
	IgnoreWarnings bool           `mapstructure:"ignore_warnings" hcl:"ignore_warnings,optional"`
}

// Canonicalize CheckRestart fields if not nil.
func (c *CheckRestart) Canonicalize() {
	if c == nil {
		return
	}

	if c.Grace == nil {
		c.Grace = pointerOf(1 * time.Second)
	}
}

// Copy returns a copy of CheckRestart or nil if unset.
func (c *CheckRestart) Copy() *CheckRestart {
	if c == nil {
		return nil
	}

	nc := new(CheckRestart)
	nc.Limit = c.Limit
	if c.Grace != nil {
		g := *c.Grace
		nc.Grace = &g
	}
	nc.IgnoreWarnings = c.IgnoreWarnings
	return nc
}

// Merge values from other CheckRestart over default values on this
// CheckRestart and return merged copy.
func (c *CheckRestart) Merge(o *CheckRestart) *CheckRestart {
	if c == nil {
		// Just return other
		return o
	}

	nc := c.Copy()

	if o == nil {
		// Nothing to merge
		return nc
	}

	if o.Limit > 0 {
		nc.Limit = o.Limit
	}

	if o.Grace != nil {
		nc.Grace = o.Grace
	}

	if o.IgnoreWarnings {
		nc.IgnoreWarnings = o.IgnoreWarnings
	}

	return nc
}

// ServiceCheck represents a Nomad job-submitters view of a Consul service health check.
type ServiceCheck struct {
	Name                   string              `hcl:"name,optional"`
	Type                   string              `hcl:"type,optional"`
	Command                string              `hcl:"command,optional"`
	Args                   []string            `hcl:"args,optional"`
	Path                   string              `hcl:"path,optional"`
	Protocol               string              `hcl:"protocol,optional"`
	PortLabel              string              `mapstructure:"port" hcl:"port,optional"`
	Expose                 bool                `hcl:"expose,optional"`
	AddressMode            string              `mapstructure:"address_mode" hcl:"address_mode,optional"`
	Advertise              string              `hcl:"advertise,optional"`
	Interval               time.Duration       `hcl:"interval,optional"`
	Timeout                time.Duration       `hcl:"timeout,optional"`
	InitialStatus          string              `mapstructure:"initial_status" hcl:"initial_status,optional"`
	Notes                  string              `hcl:"notes,optional"`
	TLSServerName          string              `mapstructure:"tls_server_name" hcl:"tls_server_name,optional"`
	TLSSkipVerify          bool                `mapstructure:"tls_skip_verify" hcl:"tls_skip_verify,optional"`
	Header                 map[string][]string `hcl:"header,block"`
	Method                 string              `hcl:"method,optional"`
	CheckRestart           *CheckRestart       `mapstructure:"check_restart" hcl:"check_restart,block"`
	GRPCService            string              `mapstructure:"grpc_service" hcl:"grpc_service,optional"`
	GRPCUseTLS             bool                `mapstructure:"grpc_use_tls" hcl:"grpc_use_tls,optional"`
	TaskName               string              `mapstructure:"task" hcl:"task,optional"`
	SuccessBeforePassing   int                 `mapstructure:"success_before_passing" hcl:"success_before_passing,optional"`
	FailuresBeforeCritical int                 `mapstructure:"failures_before_critical" hcl:"failures_before_critical,optional"`
	FailuresBeforeWarning  int                 `mapstructure:"failures_before_warning" hcl:"failures_before_warning,optional"`
	Body                   string              `hcl:"body,optional"`
	OnUpdate               string              `mapstructure:"on_update" hcl:"on_update,optional"`
}

// Service represents a Nomad job-submitters view of a Consul or Nomad service.
type Service struct {
	Name              string            `hcl:"name,optional"`
	Tags              []string          `hcl:"tags,optional"`
	CanaryTags        []string          `mapstructure:"canary_tags" hcl:"canary_tags,optional"`
	EnableTagOverride bool              `mapstructure:"enable_tag_override" hcl:"enable_tag_override,optional"`
	PortLabel         string            `mapstructure:"port" hcl:"port,optional"`
	AddressMode       string            `mapstructure:"address_mode" hcl:"address_mode,optional"`
	Address           string            `hcl:"address,optional"`
	Checks            []ServiceCheck    `hcl:"check,block"`
	CheckRestart      *CheckRestart     `mapstructure:"check_restart" hcl:"check_restart,block"`
	Connect           *ConsulConnect    `hcl:"connect,block"`
	Meta              map[string]string `hcl:"meta,block"`
	CanaryMeta        map[string]string `hcl:"canary_meta,block"`
	TaggedAddresses   map[string]string `hcl:"tagged_addresses,block"`
	TaskName          string            `mapstructure:"task" hcl:"task,optional"`
	OnUpdate          string            `mapstructure:"on_update" hcl:"on_update,optional"`
	Identity          *WorkloadIdentity `hcl:"identity,block"`
	Weights           *ServiceWeights   `mapstructure:"weights" hcl:"weights,block"`

	// Provider defines which backend system provides the service registration,
	// either "consul" (default) or "nomad".
	Provider string `hcl:"provider,optional"`

	// Cluster is valid only for Nomad Enterprise with provider: consul
	Cluster string `hcl:"cluster,optional"`
}

const (
	OnUpdateRequireHealthy = "require_healthy"
	OnUpdateIgnoreWarn     = "ignore_warnings"
	OnUpdateIgnore         = "ignore"

	// ServiceProviderConsul is the default provider for services when no
	// parameter is set.
	ServiceProviderConsul = "consul"
)

// Canonicalize the Service by ensuring its name and address mode are set. Task
// will be nil for group services.
func (s *Service) Canonicalize(t *Task, tg *TaskGroup, job *Job) {
	if s.Name == "" {
		if t != nil {
			s.Name = fmt.Sprintf("%s-%s-%s", *job.Name, *tg.Name, t.Name)
		} else {
			s.Name = fmt.Sprintf("%s-%s", *job.Name, *tg.Name)
		}
	}

	// Default to AddressModeAuto
	if s.AddressMode == "" {
		s.AddressMode = "auto"
	}

	// Default to OnUpdateRequireHealthy
	if s.OnUpdate == "" {
		s.OnUpdate = OnUpdateRequireHealthy
	}

	// Default the service provider.
	if s.Provider == "" {
		s.Provider = ServiceProviderConsul
	}
	if s.Cluster == "" {
		s.Cluster = "default"
	}

	if len(s.Meta) == 0 {
		s.Meta = nil
	}

	if len(s.CanaryMeta) == 0 {
		s.CanaryMeta = nil
	}

	if len(s.TaggedAddresses) == 0 {
		s.TaggedAddresses = nil
	}

	s.Connect.Canonicalize()
	s.Weights.Canonicalize()

	// Canonicalize CheckRestart on Checks and merge Service.CheckRestart
	// into each check.
	for i, check := range s.Checks {
		s.Checks[i].CheckRestart = s.CheckRestart.Merge(check.CheckRestart)
		s.Checks[i].CheckRestart.Canonicalize()

		if s.Checks[i].SuccessBeforePassing < 0 {
			s.Checks[i].SuccessBeforePassing = 0
		}

		if s.Checks[i].FailuresBeforeCritical < 0 {
			s.Checks[i].FailuresBeforeCritical = 0
		}

		if s.Checks[i].FailuresBeforeWarning < 0 {
			s.Checks[i].FailuresBeforeWarning = 0
		}

		// Inhert Service
		if s.Checks[i].OnUpdate == "" {
			s.Checks[i].OnUpdate = s.OnUpdate
		}
	}
}

// ServiceWeights is the jobspec block which configures how a service instance
// is weighted in a DNS SRV request based on the service's health status.
type ServiceWeights struct {
	Passing int `hcl:"passing,optional"`
	Warning int `hcl:"warning,optional"`
}

func (weights *ServiceWeights) Canonicalize() {
	if weights == nil {
		return
	}

	if weights.Passing <= 0 {
		weights.Passing = 1
	}
	if weights.Warning <= 0 {
		weights.Warning = 1
	}
}
