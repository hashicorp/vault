package api

import (
	"fmt"
	"time"
)

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
	Id                     string              `hcl:"id,optional"`
	Name                   string              `hcl:"name,optional"`
	Type                   string              `hcl:"type,optional"`
	Command                string              `hcl:"command,optional"`
	Args                   []string            `hcl:"args,optional"`
	Path                   string              `hcl:"path,optional"`
	Protocol               string              `hcl:"protocol,optional"`
	PortLabel              string              `mapstructure:"port" hcl:"port,optional"`
	Expose                 bool                `hcl:"expose,optional"`
	AddressMode            string              `mapstructure:"address_mode" hcl:"address_mode,optional"`
	Interval               time.Duration       `hcl:"interval,optional"`
	Timeout                time.Duration       `hcl:"timeout,optional"`
	InitialStatus          string              `mapstructure:"initial_status" hcl:"initial_status,optional"`
	TLSSkipVerify          bool                `mapstructure:"tls_skip_verify" hcl:"tls_skip_verify,optional"`
	Header                 map[string][]string `hcl:"header,block"`
	Method                 string              `hcl:"method,optional"`
	CheckRestart           *CheckRestart       `mapstructure:"check_restart" hcl:"check_restart,block"`
	GRPCService            string              `mapstructure:"grpc_service" hcl:"grpc_service,optional"`
	GRPCUseTLS             bool                `mapstructure:"grpc_use_tls" hcl:"grpc_use_tls,optional"`
	TaskName               string              `mapstructure:"task" hcl:"task,optional"`
	SuccessBeforePassing   int                 `mapstructure:"success_before_passing" hcl:"success_before_passing,optional"`
	FailuresBeforeCritical int                 `mapstructure:"failures_before_critical" hcl:"failures_before_critical,optional"`
}

// Service represents a Consul service definition.
type Service struct {
	//FIXME Id is unused. Remove?
	Id                string            `hcl:"id,optional"`
	Name              string            `hcl:"name,optional"`
	Tags              []string          `hcl:"tags,optional"`
	CanaryTags        []string          `mapstructure:"canary_tags" hcl:"canary_tags,optional"`
	EnableTagOverride bool              `mapstructure:"enable_tag_override" hcl:"enable_tag_override,optional"`
	PortLabel         string            `mapstructure:"port" hcl:"port,optional"`
	AddressMode       string            `mapstructure:"address_mode" hcl:"address_mode,optional"`
	Checks            []ServiceCheck    `hcl:"check,block"`
	CheckRestart      *CheckRestart     `mapstructure:"check_restart" hcl:"check_restart,block"`
	Connect           *ConsulConnect    `hcl:"connect,block"`
	Meta              map[string]string `hcl:"meta,block"`
	CanaryMeta        map[string]string `hcl:"canary_meta,block"`
	TaskName          string            `mapstructure:"task" hcl:"task,optional"`
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

	s.Connect.Canonicalize()

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
	}
}

// ConsulConnect represents a Consul Connect jobspec stanza.
type ConsulConnect struct {
	Native         bool                  `hcl:"native,optional"`
	Gateway        *ConsulGateway        `hcl:"gateway,block"`
	SidecarService *ConsulSidecarService `mapstructure:"sidecar_service" hcl:"sidecar_service,block"`
	SidecarTask    *SidecarTask          `mapstructure:"sidecar_task" hcl:"sidecar_task,block"`
}

func (cc *ConsulConnect) Canonicalize() {
	if cc == nil {
		return
	}

	cc.SidecarService.Canonicalize()
	cc.SidecarTask.Canonicalize()
	cc.Gateway.Canonicalize()
}

// ConsulSidecarService represents a Consul Connect SidecarService jobspec
// stanza.
type ConsulSidecarService struct {
	Tags  []string     `hcl:"tags,optional"`
	Port  string       `hcl:"port,optional"`
	Proxy *ConsulProxy `hcl:"proxy,block"`
}

func (css *ConsulSidecarService) Canonicalize() {
	if css == nil {
		return
	}

	if len(css.Tags) == 0 {
		css.Tags = nil
	}

	css.Proxy.Canonicalize()
}

// SidecarTask represents a subset of Task fields that can be set to override
// the fields of the Task generated for the sidecar
type SidecarTask struct {
	Name          string                 `hcl:"name,optional"`
	Driver        string                 `hcl:"driver,optional"`
	User          string                 `hcl:"user,optional"`
	Config        map[string]interface{} `hcl:"config,block"`
	Env           map[string]string      `hcl:"env,block"`
	Resources     *Resources             `hcl:"resources,block"`
	Meta          map[string]string      `hcl:"meta,block"`
	KillTimeout   *time.Duration         `mapstructure:"kill_timeout" hcl:"kill_timeout,optional"`
	LogConfig     *LogConfig             `mapstructure:"logs" hcl:"logs,block"`
	ShutdownDelay *time.Duration         `mapstructure:"shutdown_delay" hcl:"shutdown_delay,optional"`
	KillSignal    string                 `mapstructure:"kill_signal" hcl:"kill_signal,optional"`
}

func (st *SidecarTask) Canonicalize() {
	if st == nil {
		return
	}

	if len(st.Config) == 0 {
		st.Config = nil
	}

	if len(st.Env) == 0 {
		st.Env = nil
	}

	if st.Resources == nil {
		st.Resources = DefaultResources()
	} else {
		st.Resources.Canonicalize()
	}

	if st.LogConfig == nil {
		st.LogConfig = DefaultLogConfig()
	} else {
		st.LogConfig.Canonicalize()
	}

	if len(st.Meta) == 0 {
		st.Meta = nil
	}

	if st.KillTimeout == nil {
		st.KillTimeout = timeToPtr(5 * time.Second)
	}

	if st.ShutdownDelay == nil {
		st.ShutdownDelay = timeToPtr(0)
	}
}

// ConsulProxy represents a Consul Connect sidecar proxy jobspec stanza.
type ConsulProxy struct {
	LocalServiceAddress string                 `mapstructure:"local_service_address" hcl:"local_service_address,optional"`
	LocalServicePort    int                    `mapstructure:"local_service_port" hcl:"local_service_port,optional"`
	ExposeConfig        *ConsulExposeConfig    `mapstructure:"expose" hcl:"expose,block"`
	Upstreams           []*ConsulUpstream      `hcl:"upstreams,block"`
	Config              map[string]interface{} `hcl:"config,block"`
}

func (cp *ConsulProxy) Canonicalize() {
	if cp == nil {
		return
	}

	cp.ExposeConfig.Canonicalize()

	if len(cp.Upstreams) == 0 {
		cp.Upstreams = nil
	}

	if len(cp.Config) == 0 {
		cp.Config = nil
	}
}

// ConsulUpstream represents a Consul Connect upstream jobspec stanza.
type ConsulUpstream struct {
	DestinationName string `mapstructure:"destination_name" hcl:"destination_name,optional"`
	LocalBindPort   int    `mapstructure:"local_bind_port" hcl:"local_bind_port,optional"`
	Datacenter      string `mapstructure:"datacenter" hcl:"datacenter,optional"`
}

type ConsulExposeConfig struct {
	Path []*ConsulExposePath `mapstructure:"path" hcl:"path,block"`
}

func (cec *ConsulExposeConfig) Canonicalize() {
	if cec == nil {
		return
	}

	if len(cec.Path) == 0 {
		cec.Path = nil
	}
}

type ConsulExposePath struct {
	Path          string `hcl:"path,optional"`
	Protocol      string `hcl:"protocol,optional"`
	LocalPathPort int    `mapstructure:"local_path_port" hcl:"local_path_port,optional"`
	ListenerPort  string `mapstructure:"listener_port" hcl:"listener_port,optional"`
}

// ConsulGateway is used to configure one of the Consul Connect Gateway types.
type ConsulGateway struct {
	// Proxy is used to configure the Envoy instance acting as the gateway.
	Proxy *ConsulGatewayProxy `hcl:"proxy,block"`

	// Ingress represents the Consul Configuration Entry for an Ingress Gateway.
	Ingress *ConsulIngressConfigEntry `hcl:"ingress,block"`

	// Terminating represents the Consul Configuration Entry for a Terminating Gateway.
	Terminating *ConsulTerminatingConfigEntry `hcl:"terminating,block"`

	// Mesh is not yet supported.
	// Mesh *ConsulMeshConfigEntry
}

func (g *ConsulGateway) Canonicalize() {
	if g == nil {
		return
	}
	g.Proxy.Canonicalize()
	g.Ingress.Canonicalize()
	g.Terminating.Canonicalize()
}

func (g *ConsulGateway) Copy() *ConsulGateway {
	if g == nil {
		return nil
	}

	return &ConsulGateway{
		Proxy:       g.Proxy.Copy(),
		Ingress:     g.Ingress.Copy(),
		Terminating: g.Terminating.Copy(),
	}
}

type ConsulGatewayBindAddress struct {
	Name    string `hcl:",label"`
	Address string `mapstructure:"address" hcl:"address,optional"`
	Port    int    `mapstructure:"port" hcl:"port,optional"`
}

var (
	// defaultGatewayConnectTimeout is the default amount of time connections to
	// upstreams are allowed before timing out.
	defaultGatewayConnectTimeout = 5 * time.Second
)

// ConsulGatewayProxy is used to tune parameters of the proxy instance acting as
// one of the forms of Connect gateways that Consul supports.
//
// https://www.consul.io/docs/connect/proxies/envoy#gateway-options
type ConsulGatewayProxy struct {
	ConnectTimeout                  *time.Duration                       `mapstructure:"connect_timeout" hcl:"connect_timeout,optional"`
	EnvoyGatewayBindTaggedAddresses bool                                 `mapstructure:"envoy_gateway_bind_tagged_addresses" hcl:"envoy_gateway_bind_tagged_addresses,optional"`
	EnvoyGatewayBindAddresses       map[string]*ConsulGatewayBindAddress `mapstructure:"envoy_gateway_bind_addresses" hcl:"envoy_gateway_bind_addresses,block"`
	EnvoyGatewayNoDefaultBind       bool                                 `mapstructure:"envoy_gateway_no_default_bind" hcl:"envoy_gateway_no_default_bind,optional"`
	EnvoyDNSDiscoveryType           string                               `mapstructure:"envoy_dns_discovery_type" hcl:"envoy_dns_discovery_type,optional"`
	Config                          map[string]interface{}               `hcl:"config,block"` // escape hatch envoy config
}

func (p *ConsulGatewayProxy) Canonicalize() {
	if p == nil {
		return
	}

	if p.ConnectTimeout == nil {
		// same as the default from consul
		p.ConnectTimeout = timeToPtr(defaultGatewayConnectTimeout)
	}

	if len(p.EnvoyGatewayBindAddresses) == 0 {
		p.EnvoyGatewayBindAddresses = nil
	}

	if len(p.Config) == 0 {
		p.Config = nil
	}
}

func (p *ConsulGatewayProxy) Copy() *ConsulGatewayProxy {
	if p == nil {
		return nil
	}

	var binds map[string]*ConsulGatewayBindAddress = nil
	if p.EnvoyGatewayBindAddresses != nil {
		binds = make(map[string]*ConsulGatewayBindAddress, len(p.EnvoyGatewayBindAddresses))
		for k, v := range p.EnvoyGatewayBindAddresses {
			binds[k] = v
		}
	}

	var config map[string]interface{} = nil
	if p.Config != nil {
		config = make(map[string]interface{}, len(p.Config))
		for k, v := range p.Config {
			config[k] = v
		}
	}

	return &ConsulGatewayProxy{
		ConnectTimeout:                  timeToPtr(*p.ConnectTimeout),
		EnvoyGatewayBindTaggedAddresses: p.EnvoyGatewayBindTaggedAddresses,
		EnvoyGatewayBindAddresses:       binds,
		EnvoyGatewayNoDefaultBind:       p.EnvoyGatewayNoDefaultBind,
		EnvoyDNSDiscoveryType:           p.EnvoyDNSDiscoveryType,
		Config:                          config,
	}
}

// ConsulGatewayTLSConfig is used to configure TLS for a gateway.
type ConsulGatewayTLSConfig struct {
	Enabled bool `hcl:"enabled,optional"`
}

func (tc *ConsulGatewayTLSConfig) Canonicalize() {
}

func (tc *ConsulGatewayTLSConfig) Copy() *ConsulGatewayTLSConfig {
	if tc == nil {
		return nil
	}

	return &ConsulGatewayTLSConfig{
		Enabled: tc.Enabled,
	}
}

// ConsulIngressService is used to configure a service fronted by the ingress gateway.
type ConsulIngressService struct {
	// Namespace is not yet supported.
	// Namespace string
	Name string `hcl:"name,optional"`

	Hosts []string `hcl:"hosts,optional"`
}

func (s *ConsulIngressService) Canonicalize() {
	if s == nil {
		return
	}

	if len(s.Hosts) == 0 {
		s.Hosts = nil
	}
}

func (s *ConsulIngressService) Copy() *ConsulIngressService {
	if s == nil {
		return nil
	}

	var hosts []string = nil
	if n := len(s.Hosts); n > 0 {
		hosts = make([]string, n)
		copy(hosts, s.Hosts)
	}

	return &ConsulIngressService{
		Name:  s.Name,
		Hosts: hosts,
	}
}

const (
	defaultIngressListenerProtocol = "tcp"
)

// ConsulIngressListener is used to configure a listener on a Consul Ingress
// Gateway.
type ConsulIngressListener struct {
	Port     int                     `hcl:"port,optional"`
	Protocol string                  `hcl:"protocol,optional"`
	Services []*ConsulIngressService `hcl:"service,block"`
}

func (l *ConsulIngressListener) Canonicalize() {
	if l == nil {
		return
	}

	if l.Protocol == "" {
		// same as default from consul
		l.Protocol = defaultIngressListenerProtocol
	}

	if len(l.Services) == 0 {
		l.Services = nil
	}
}

func (l *ConsulIngressListener) Copy() *ConsulIngressListener {
	if l == nil {
		return nil
	}

	var services []*ConsulIngressService = nil
	if n := len(l.Services); n > 0 {
		services = make([]*ConsulIngressService, n)
		for i := 0; i < n; i++ {
			services[i] = l.Services[i].Copy()
		}
	}

	return &ConsulIngressListener{
		Port:     l.Port,
		Protocol: l.Protocol,
		Services: services,
	}
}

// ConsulIngressConfigEntry represents the Consul Configuration Entry type for
// an Ingress Gateway.
//
// https://www.consul.io/docs/agent/config-entries/ingress-gateway#available-fields
type ConsulIngressConfigEntry struct {
	// Namespace is not yet supported.
	// Namespace string

	TLS       *ConsulGatewayTLSConfig  `hcl:"tls,block"`
	Listeners []*ConsulIngressListener `hcl:"listener,block"`
}

func (e *ConsulIngressConfigEntry) Canonicalize() {
	if e == nil {
		return
	}

	e.TLS.Canonicalize()

	if len(e.Listeners) == 0 {
		e.Listeners = nil
	}

	for _, listener := range e.Listeners {
		listener.Canonicalize()
	}
}

func (e *ConsulIngressConfigEntry) Copy() *ConsulIngressConfigEntry {
	if e == nil {
		return nil
	}

	var listeners []*ConsulIngressListener = nil
	if n := len(e.Listeners); n > 0 {
		listeners = make([]*ConsulIngressListener, n)
		for i := 0; i < n; i++ {
			listeners[i] = e.Listeners[i].Copy()
		}
	}

	return &ConsulIngressConfigEntry{
		TLS:       e.TLS.Copy(),
		Listeners: listeners,
	}
}

type ConsulLinkedService struct {
	Name     string `hcl:"name,optional"`
	CAFile   string `hcl:"ca_file,optional"`
	CertFile string `hcl:"cert_file,optional"`
	KeyFile  string `hcl:"key_file,optional"`
	SNI      string `hcl:"sni,optional"`
}

func (s *ConsulLinkedService) Canonicalize() {
	// nothing to do for now
}

func (s *ConsulLinkedService) Copy() *ConsulLinkedService {
	if s == nil {
		return nil
	}

	return &ConsulLinkedService{
		Name:     s.Name,
		CAFile:   s.CAFile,
		CertFile: s.CertFile,
		KeyFile:  s.KeyFile,
		SNI:      s.SNI,
	}
}

// ConsulTerminatingConfigEntry represents the Consul Configuration Entry type
// for a Terminating Gateway.
//
// https://www.consul.io/docs/agent/config-entries/terminating-gateway#available-fields
type ConsulTerminatingConfigEntry struct {
	// Namespace is not yet supported.
	// Namespace string

	Services []*ConsulLinkedService `hcl:"service,block"`
}

func (e *ConsulTerminatingConfigEntry) Canonicalize() {
	if e == nil {
		return
	}

	if len(e.Services) == 0 {
		e.Services = nil
	}

	for _, service := range e.Services {
		service.Canonicalize()
	}
}

func (e *ConsulTerminatingConfigEntry) Copy() *ConsulTerminatingConfigEntry {
	if e == nil {
		return nil
	}

	var services []*ConsulLinkedService = nil
	if n := len(e.Services); n > 0 {
		services = make([]*ConsulLinkedService, n)
		for i := 0; i < n; i++ {
			services[i] = e.Services[i].Copy()
		}
	}

	return &ConsulTerminatingConfigEntry{
		Services: services,
	}
}

// ConsulMeshConfigEntry is not yet supported.
// type ConsulMeshConfigEntry struct {
// }
