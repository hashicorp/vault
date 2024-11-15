package godo

import (
	"context"
	"fmt"
	"net/http"
)

const (
	cachePath             = "cache"
	dropletsPath          = "droplets"
	forwardingRulesPath   = "forwarding_rules"
	loadBalancersBasePath = "/v2/load_balancers"
)

const (
	// Load Balancer types
	LoadBalancerTypeGlobal          = "GLOBAL"
	LoadBalancerTypeRegional        = "REGIONAL"
	LoadBalancerTypeRegionalNetwork = "REGIONAL_NETWORK"

	// Load Balancer network types
	LoadBalancerNetworkTypeExternal = "EXTERNAL"
	LoadBalancerNetworkTypeInternal = "INTERNAL"
)

// LoadBalancersService is an interface for managing load balancers with the DigitalOcean API.
// See: https://docs.digitalocean.com/reference/api/api-reference/#tag/Load-Balancers
type LoadBalancersService interface {
	Get(context.Context, string) (*LoadBalancer, *Response, error)
	List(context.Context, *ListOptions) ([]LoadBalancer, *Response, error)
	Create(context.Context, *LoadBalancerRequest) (*LoadBalancer, *Response, error)
	Update(ctx context.Context, lbID string, lbr *LoadBalancerRequest) (*LoadBalancer, *Response, error)
	Delete(ctx context.Context, lbID string) (*Response, error)
	AddDroplets(ctx context.Context, lbID string, dropletIDs ...int) (*Response, error)
	RemoveDroplets(ctx context.Context, lbID string, dropletIDs ...int) (*Response, error)
	AddForwardingRules(ctx context.Context, lbID string, rules ...ForwardingRule) (*Response, error)
	RemoveForwardingRules(ctx context.Context, lbID string, rules ...ForwardingRule) (*Response, error)
	PurgeCache(ctx context.Context, lbID string) (*Response, error)
}

// LoadBalancer represents a DigitalOcean load balancer configuration.
// Tags can only be provided upon the creation of a Load Balancer.
type LoadBalancer struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	IP   string `json:"ip,omitempty"`
	// SizeSlug is mutually exclusive with SizeUnit. Only one should be specified
	SizeSlug string `json:"size,omitempty"`
	// SizeUnit is mutually exclusive with SizeSlug. Only one should be specified
	SizeUnit                     uint32           `json:"size_unit,omitempty"`
	Type                         string           `json:"type,omitempty"`
	Algorithm                    string           `json:"algorithm,omitempty"`
	Status                       string           `json:"status,omitempty"`
	Created                      string           `json:"created_at,omitempty"`
	ForwardingRules              []ForwardingRule `json:"forwarding_rules,omitempty"`
	HealthCheck                  *HealthCheck     `json:"health_check,omitempty"`
	StickySessions               *StickySessions  `json:"sticky_sessions,omitempty"`
	Region                       *Region          `json:"region,omitempty"`
	DropletIDs                   []int            `json:"droplet_ids,omitempty"`
	Tag                          string           `json:"tag,omitempty"`
	Tags                         []string         `json:"tags,omitempty"`
	RedirectHttpToHttps          bool             `json:"redirect_http_to_https,omitempty"`
	EnableProxyProtocol          bool             `json:"enable_proxy_protocol,omitempty"`
	EnableBackendKeepalive       bool             `json:"enable_backend_keepalive,omitempty"`
	VPCUUID                      string           `json:"vpc_uuid,omitempty"`
	DisableLetsEncryptDNSRecords *bool            `json:"disable_lets_encrypt_dns_records,omitempty"`
	ValidateOnly                 bool             `json:"validate_only,omitempty"`
	ProjectID                    string           `json:"project_id,omitempty"`
	HTTPIdleTimeoutSeconds       *uint64          `json:"http_idle_timeout_seconds,omitempty"`
	Firewall                     *LBFirewall      `json:"firewall,omitempty"`
	Domains                      []*LBDomain      `json:"domains,omitempty"`
	GLBSettings                  *GLBSettings     `json:"glb_settings,omitempty"`
	TargetLoadBalancerIDs        []string         `json:"target_load_balancer_ids,omitempty"`
	Network                      string           `json:"network,omitempty"`
}

// String creates a human-readable description of a LoadBalancer.
func (l LoadBalancer) String() string {
	return Stringify(l)
}

// URN returns the load balancer ID in a valid DO API URN form.
func (l LoadBalancer) URN() string {
	return ToURN("LoadBalancer", l.ID)
}

// AsRequest creates a LoadBalancerRequest that can be submitted to Update with the current values of the LoadBalancer.
// Modifying the returned LoadBalancerRequest will not modify the original LoadBalancer.
func (l LoadBalancer) AsRequest() *LoadBalancerRequest {
	r := LoadBalancerRequest{
		Name:                         l.Name,
		Algorithm:                    l.Algorithm,
		SizeSlug:                     l.SizeSlug,
		SizeUnit:                     l.SizeUnit,
		Type:                         l.Type,
		ForwardingRules:              append([]ForwardingRule(nil), l.ForwardingRules...),
		DropletIDs:                   append([]int(nil), l.DropletIDs...),
		Tag:                          l.Tag,
		RedirectHttpToHttps:          l.RedirectHttpToHttps,
		EnableProxyProtocol:          l.EnableProxyProtocol,
		EnableBackendKeepalive:       l.EnableBackendKeepalive,
		VPCUUID:                      l.VPCUUID,
		DisableLetsEncryptDNSRecords: l.DisableLetsEncryptDNSRecords,
		ValidateOnly:                 l.ValidateOnly,
		ProjectID:                    l.ProjectID,
		HTTPIdleTimeoutSeconds:       l.HTTPIdleTimeoutSeconds,
		TargetLoadBalancerIDs:        append([]string(nil), l.TargetLoadBalancerIDs...),
		Network:                      l.Network,
	}

	if l.DisableLetsEncryptDNSRecords != nil {
		*r.DisableLetsEncryptDNSRecords = *l.DisableLetsEncryptDNSRecords
	}

	if l.HealthCheck != nil {
		r.HealthCheck = &HealthCheck{}
		*r.HealthCheck = *l.HealthCheck
	}

	if l.StickySessions != nil {
		r.StickySessions = &StickySessions{}
		*r.StickySessions = *l.StickySessions
	}

	if l.Region != nil {
		r.Region = l.Region.Slug
	}

	if l.Firewall != nil {
		r.Firewall = l.Firewall.deepCopy()
	}

	for _, domain := range l.Domains {
		lbDomain := &LBDomain{}
		*lbDomain = *domain
		lbDomain.VerificationErrorReasons = append([]string(nil), domain.VerificationErrorReasons...)
		lbDomain.SSLValidationErrorReasons = append([]string(nil), domain.SSLValidationErrorReasons...)
		r.Domains = append(r.Domains, lbDomain)
	}

	if l.GLBSettings != nil {
		r.GLBSettings = l.GLBSettings.deepCopy()
	}

	return &r
}

// ForwardingRule represents load balancer forwarding rules.
type ForwardingRule struct {
	EntryProtocol  string `json:"entry_protocol,omitempty"`
	EntryPort      int    `json:"entry_port,omitempty"`
	TargetProtocol string `json:"target_protocol,omitempty"`
	TargetPort     int    `json:"target_port,omitempty"`
	CertificateID  string `json:"certificate_id,omitempty"`
	TlsPassthrough bool   `json:"tls_passthrough,omitempty"`
}

// String creates a human-readable description of a ForwardingRule.
func (f ForwardingRule) String() string {
	return Stringify(f)
}

// HealthCheck represents optional load balancer health check rules.
type HealthCheck struct {
	Protocol               string `json:"protocol,omitempty"`
	Port                   int    `json:"port,omitempty"`
	Path                   string `json:"path,omitempty"`
	CheckIntervalSeconds   int    `json:"check_interval_seconds,omitempty"`
	ResponseTimeoutSeconds int    `json:"response_timeout_seconds,omitempty"`
	HealthyThreshold       int    `json:"healthy_threshold,omitempty"`
	UnhealthyThreshold     int    `json:"unhealthy_threshold,omitempty"`
	ProxyProtocol          *bool  `json:"proxy_protocol,omitempty"`
}

// String creates a human-readable description of a HealthCheck.
func (h HealthCheck) String() string {
	return Stringify(h)
}

// StickySessions represents optional load balancer session affinity rules.
type StickySessions struct {
	Type             string `json:"type,omitempty"`
	CookieName       string `json:"cookie_name,omitempty"`
	CookieTtlSeconds int    `json:"cookie_ttl_seconds,omitempty"`
}

// String creates a human-readable description of a StickySessions instance.
func (s StickySessions) String() string {
	return Stringify(s)
}

// LBFirewall holds the allow and deny rules for a loadbalancer's firewall.
// Currently, allow and deny rules support cidrs and ips.
// Please use the helper methods (IPSourceFirewall/CIDRSourceFirewall) to format the allow/deny rules.
type LBFirewall struct {
	Allow []string `json:"allow,omitempty"`
	Deny  []string `json:"deny,omitempty"`
}

func (lbf *LBFirewall) deepCopy() *LBFirewall {
	return &LBFirewall{
		Allow: append([]string(nil), lbf.Allow...),
		Deny:  append([]string(nil), lbf.Deny...),
	}
}

// IPSourceFirewall takes an IP (string) and returns a formatted ip source firewall rule
func IPSourceFirewall(ip string) string { return fmt.Sprintf("ip:%s", ip) }

// CIDRSourceFirewall takes a CIDR notation IP address and prefix length string
// like "192.0.2.0/24" and returns a formatted cidr source firewall rule
func CIDRSourceFirewall(cidr string) string { return fmt.Sprintf("cidr:%s", cidr) }

// String creates a human-readable description of an LBFirewall instance.
func (f LBFirewall) String() string {
	return Stringify(f)
}

// LoadBalancerRequest represents the configuration to be applied to an existing or a new load balancer.
type LoadBalancerRequest struct {
	Name      string `json:"name,omitempty"`
	Algorithm string `json:"algorithm,omitempty"`
	Region    string `json:"region,omitempty"`
	// SizeSlug is mutually exclusive with SizeUnit. Only one should be specified
	SizeSlug string `json:"size,omitempty"`
	// SizeUnit is mutually exclusive with SizeSlug. Only one should be specified
	SizeUnit                     uint32           `json:"size_unit,omitempty"`
	Type                         string           `json:"type,omitempty"`
	ForwardingRules              []ForwardingRule `json:"forwarding_rules,omitempty"`
	HealthCheck                  *HealthCheck     `json:"health_check,omitempty"`
	StickySessions               *StickySessions  `json:"sticky_sessions,omitempty"`
	DropletIDs                   []int            `json:"droplet_ids,omitempty"`
	Tag                          string           `json:"tag,omitempty"`
	Tags                         []string         `json:"tags,omitempty"`
	RedirectHttpToHttps          bool             `json:"redirect_http_to_https,omitempty"`
	EnableProxyProtocol          bool             `json:"enable_proxy_protocol,omitempty"`
	EnableBackendKeepalive       bool             `json:"enable_backend_keepalive,omitempty"`
	VPCUUID                      string           `json:"vpc_uuid,omitempty"`
	DisableLetsEncryptDNSRecords *bool            `json:"disable_lets_encrypt_dns_records,omitempty"`
	ValidateOnly                 bool             `json:"validate_only,omitempty"`
	ProjectID                    string           `json:"project_id,omitempty"`
	HTTPIdleTimeoutSeconds       *uint64          `json:"http_idle_timeout_seconds,omitempty"`
	Firewall                     *LBFirewall      `json:"firewall,omitempty"`
	Domains                      []*LBDomain      `json:"domains,omitempty"`
	GLBSettings                  *GLBSettings     `json:"glb_settings,omitempty"`
	TargetLoadBalancerIDs        []string         `json:"target_load_balancer_ids,omitempty"`
	Network                      string           `json:"network,omitempty"`
}

// String creates a human-readable description of a LoadBalancerRequest.
func (l LoadBalancerRequest) String() string {
	return Stringify(l)
}

type forwardingRulesRequest struct {
	Rules []ForwardingRule `json:"forwarding_rules,omitempty"`
}

func (l forwardingRulesRequest) String() string {
	return Stringify(l)
}

type dropletIDsRequest struct {
	IDs []int `json:"droplet_ids,omitempty"`
}

func (l dropletIDsRequest) String() string {
	return Stringify(l)
}

// LBDomain defines domain names required to ingress traffic to a Global LB
type LBDomain struct {
	// Name defines the domain fqdn
	Name string `json:"name"`
	// IsManaged indicates if the domain is DO-managed
	IsManaged bool `json:"is_managed"`
	// CertificateID indicates ID of a TLS certificate
	CertificateID string `json:"certificate_id,omitempty"`
	// Status indicates the domain validation status
	Status string `json:"status,omitempty"`
	// VerificationErrorReasons indicates any domain verification errors
	VerificationErrorReasons []string `json:"verification_error_reasons,omitempty"`
	// SSLValidationErrorReasons indicates any domain SSL validation errors
	SSLValidationErrorReasons []string `json:"ssl_validation_error_reasons,omitempty"`
}

// String creates a human-readable description of a LBDomain
func (d LBDomain) String() string {
	return Stringify(d)
}

// GLBSettings define settings for configuring a Global LB
type GLBSettings struct {
	// TargetProtocol is the outgoing traffic protocol.
	TargetProtocol string `json:"target_protocol"`
	// EntryPort is the outgoing traffic port.
	TargetPort uint32 `json:"target_port"`
	// CDNSettings is the CDN configurations
	CDN *CDNSettings `json:"cdn"`
	// RegionPriorities embeds regional priority information for regional active-passive failover policy
	RegionPriorities map[string]uint32 `json:"region_priorities,omitempty"`
	// FailoverThreshold embeds failover threshold percentage for regional active-passive failover policy
	FailoverThreshold uint32 `json:"failover_threshold,omitempty"`
}

// String creates a human-readable description of a GLBSettings
func (s GLBSettings) String() string {
	return Stringify(s)
}

func (s GLBSettings) deepCopy() *GLBSettings {
	settings := &GLBSettings{
		TargetProtocol:    s.TargetProtocol,
		TargetPort:        s.TargetPort,
		RegionPriorities:  s.RegionPriorities,
		FailoverThreshold: s.FailoverThreshold,
	}
	if s.CDN != nil {
		settings.CDN = &CDNSettings{IsEnabled: s.CDN.IsEnabled}
	}
	return settings
}

// CDNSettings define CDN settings for a Global LB
type CDNSettings struct {
	// IsEnabled is the caching enabled flag
	IsEnabled bool `json:"is_enabled"`
}

// String creates a human-readable description of a CDNSettings
func (c CDNSettings) String() string {
	return Stringify(c)
}

type loadBalancersRoot struct {
	LoadBalancers []LoadBalancer `json:"load_balancers"`
	Links         *Links         `json:"links"`
	Meta          *Meta          `json:"meta"`
}

type loadBalancerRoot struct {
	LoadBalancer *LoadBalancer `json:"load_balancer"`
}

// LoadBalancersServiceOp handles communication with load balancer-related methods of the DigitalOcean API.
type LoadBalancersServiceOp struct {
	client *Client
}

var _ LoadBalancersService = &LoadBalancersServiceOp{}

// Get an existing load balancer by its identifier.
func (l *LoadBalancersServiceOp) Get(ctx context.Context, lbID string) (*LoadBalancer, *Response, error) {
	path := fmt.Sprintf("%s/%s", loadBalancersBasePath, lbID)

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadBalancerRoot)
	resp, err := l.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.LoadBalancer, resp, err
}

// List load balancers, with optional pagination.
func (l *LoadBalancersServiceOp) List(ctx context.Context, opt *ListOptions) ([]LoadBalancer, *Response, error) {
	path, err := addOptions(loadBalancersBasePath, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := l.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadBalancersRoot)
	resp, err := l.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.LoadBalancers, resp, err
}

// Create a new load balancer with a given configuration.
func (l *LoadBalancersServiceOp) Create(ctx context.Context, lbr *LoadBalancerRequest) (*LoadBalancer, *Response, error) {
	req, err := l.client.NewRequest(ctx, http.MethodPost, loadBalancersBasePath, lbr)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadBalancerRoot)
	resp, err := l.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.LoadBalancer, resp, err
}

// Update an existing load balancer with new configuration.
func (l *LoadBalancersServiceOp) Update(ctx context.Context, lbID string, lbr *LoadBalancerRequest) (*LoadBalancer, *Response, error) {
	path := fmt.Sprintf("%s/%s", loadBalancersBasePath, lbID)

	req, err := l.client.NewRequest(ctx, "PUT", path, lbr)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadBalancerRoot)
	resp, err := l.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.LoadBalancer, resp, err
}

// Delete a load balancer by its identifier.
func (l *LoadBalancersServiceOp) Delete(ctx context.Context, ldID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", loadBalancersBasePath, ldID)

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}

// AddDroplets adds droplets to a load balancer.
func (l *LoadBalancersServiceOp) AddDroplets(ctx context.Context, lbID string, dropletIDs ...int) (*Response, error) {
	path := fmt.Sprintf("%s/%s/%s", loadBalancersBasePath, lbID, dropletsPath)

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, &dropletIDsRequest{IDs: dropletIDs})
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}

// RemoveDroplets removes droplets from a load balancer.
func (l *LoadBalancersServiceOp) RemoveDroplets(ctx context.Context, lbID string, dropletIDs ...int) (*Response, error) {
	path := fmt.Sprintf("%s/%s/%s", loadBalancersBasePath, lbID, dropletsPath)

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, &dropletIDsRequest{IDs: dropletIDs})
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}

// AddForwardingRules adds forwarding rules to a load balancer.
func (l *LoadBalancersServiceOp) AddForwardingRules(ctx context.Context, lbID string, rules ...ForwardingRule) (*Response, error) {
	path := fmt.Sprintf("%s/%s/%s", loadBalancersBasePath, lbID, forwardingRulesPath)

	req, err := l.client.NewRequest(ctx, http.MethodPost, path, &forwardingRulesRequest{Rules: rules})
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}

// RemoveForwardingRules removes forwarding rules from a load balancer.
func (l *LoadBalancersServiceOp) RemoveForwardingRules(ctx context.Context, lbID string, rules ...ForwardingRule) (*Response, error) {
	path := fmt.Sprintf("%s/%s/%s", loadBalancersBasePath, lbID, forwardingRulesPath)

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, &forwardingRulesRequest{Rules: rules})
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}

// PurgeCache purges the CDN cache of a global load balancer by its identifier.
func (l *LoadBalancersServiceOp) PurgeCache(ctx context.Context, ldID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s/%s", loadBalancersBasePath, ldID, cachePath)

	req, err := l.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return l.client.Do(ctx, req, nil)
}
