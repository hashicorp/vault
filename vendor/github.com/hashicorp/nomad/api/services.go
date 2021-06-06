package api

import (
	"fmt"
	"time"
)

// CheckRestart describes if and when a task should be restarted based on
// failing health checks.
type CheckRestart struct {
	Limit          int            `mapstructure:"limit"`
	Grace          *time.Duration `mapstructure:"grace"`
	IgnoreWarnings bool           `mapstructure:"ignore_warnings"`
}

// Canonicalize CheckRestart fields if not nil.
func (c *CheckRestart) Canonicalize() {
	if c == nil {
		return
	}

	if c.Grace == nil {
		c.Grace = timeToPtr(1 * time.Second)
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

// ServiceCheck represents the consul health check that Nomad registers.
type ServiceCheck struct {
	//FIXME Id is unused. Remove?
	Id            string
	Name          string
	Type          string
	Command       string
	Args          []string
	Path          string
	Protocol      string
	PortLabel     string `mapstructure:"port"`
	AddressMode   string `mapstructure:"address_mode"`
	Interval      time.Duration
	Timeout       time.Duration
	InitialStatus string `mapstructure:"initial_status"`
	TLSSkipVerify bool   `mapstructure:"tls_skip_verify"`
	Header        map[string][]string
	Method        string
	CheckRestart  *CheckRestart `mapstructure:"check_restart"`
	GRPCService   string        `mapstructure:"grpc_service"`
	GRPCUseTLS    bool          `mapstructure:"grpc_use_tls"`
	TaskName      string        `mapstructure:"task"`
}

// Service represents a Consul service definition.
type Service struct {
	//FIXME Id is unused. Remove?
	Id           string
	Name         string
	Tags         []string
	CanaryTags   []string `mapstructure:"canary_tags"`
	PortLabel    string   `mapstructure:"port"`
	AddressMode  string   `mapstructure:"address_mode"`
	Checks       []ServiceCheck
	CheckRestart *CheckRestart `mapstructure:"check_restart"`
	Connect      *ConsulConnect
	Meta         map[string]string
}

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

	// Canonicalize CheckRestart on Checks and merge Service.CheckRestart
	// into each check.
	for i, check := range s.Checks {
		s.Checks[i].CheckRestart = s.CheckRestart.Merge(check.CheckRestart)
		s.Checks[i].CheckRestart.Canonicalize()
	}
}

// ConsulConnect represents a Consul Connect jobspec stanza.
type ConsulConnect struct {
	Native         bool
	SidecarService *ConsulSidecarService `mapstructure:"sidecar_service"`
	SidecarTask    *SidecarTask          `mapstructure:"sidecar_task"`
}

// ConsulSidecarService represents a Consul Connect SidecarService jobspec
// stanza.
type ConsulSidecarService struct {
	Tags  []string
	Port  string
	Proxy *ConsulProxy
}

// SidecarTask represents a subset of Task fields that can be set to override
// the fields of the Task generated for the sidecar
type SidecarTask struct {
	Name          string
	Driver        string
	User          string
	Config        map[string]interface{}
	Env           map[string]string
	Resources     *Resources
	Meta          map[string]string
	KillTimeout   *time.Duration `mapstructure:"kill_timeout"`
	LogConfig     *LogConfig     `mapstructure:"logs"`
	ShutdownDelay *time.Duration `mapstructure:"shutdown_delay"`
	KillSignal    string         `mapstructure:"kill_signal"`
}

// ConsulProxy represents a Consul Connect sidecar proxy jobspec stanza.
type ConsulProxy struct {
	LocalServiceAddress string `mapstructure:"local_service_address"`
	LocalServicePort    int    `mapstructure:"local_service_port"`
	Upstreams           []*ConsulUpstream
	Config              map[string]interface{}
}

// ConsulUpstream represents a Consul Connect upstream jobspec stanza.
type ConsulUpstream struct {
	DestinationName string `mapstructure:"destination_name"`
	LocalBindPort   int    `mapstructure:"local_bind_port"`
}
