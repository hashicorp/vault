// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"errors"
	"fmt"
	"net/textproto"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/go-sockaddr/template"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/helper/namespace"
)

const (
	TCP  ListenerType = "tcp"
	Unix ListenerType = "unix"
)

// ListenerType represents the supported types of listener.
type ListenerType string

type ListenerTelemetry struct {
	UnusedKeys                      UnusedKeyMap `hcl:",unusedKeyPositions"`
	UnauthenticatedMetricsAccess    bool         `hcl:"-"`
	UnauthenticatedMetricsAccessRaw interface{}  `hcl:"unauthenticated_metrics_access,alias:UnauthenticatedMetricsAccess"`
}

type ListenerProfiling struct {
	UnusedKeys                    UnusedKeyMap `hcl:",unusedKeyPositions"`
	UnauthenticatedPProfAccess    bool         `hcl:"-"`
	UnauthenticatedPProfAccessRaw interface{}  `hcl:"unauthenticated_pprof_access,alias:UnauthenticatedPProfAccessRaw"`
}

type ListenerInFlightRequestLogging struct {
	UnusedKeys                       UnusedKeyMap `hcl:",unusedKeyPositions"`
	UnauthenticatedInFlightAccess    bool         `hcl:"-"`
	UnauthenticatedInFlightAccessRaw interface{}  `hcl:"unauthenticated_in_flight_requests_access,alias:unauthenticatedInFlightAccessRaw"`
}

// Listener is the listener configuration for the server.
type Listener struct {
	UnusedKeys UnusedKeyMap `hcl:",unusedKeyPositions"`
	RawConfig  map[string]interface{}

	Type       ListenerType
	Purpose    []string    `hcl:"-"`
	PurposeRaw interface{} `hcl:"purpose"`
	Role       string      `hcl:"role"`

	Address                 string        `hcl:"address"`
	ClusterAddress          string        `hcl:"cluster_address"`
	MaxRequestSize          int64         `hcl:"-"`
	MaxRequestSizeRaw       interface{}   `hcl:"max_request_size"`
	MaxRequestDuration      time.Duration `hcl:"-"`
	MaxRequestDurationRaw   interface{}   `hcl:"max_request_duration"`
	RequireRequestHeader    bool          `hcl:"-"`
	RequireRequestHeaderRaw interface{}   `hcl:"require_request_header"`

	TLSDisable                       bool        `hcl:"-"`
	TLSDisableRaw                    interface{} `hcl:"tls_disable"`
	TLSCertFile                      string      `hcl:"tls_cert_file"`
	TLSKeyFile                       string      `hcl:"tls_key_file"`
	TLSMinVersion                    string      `hcl:"tls_min_version"`
	TLSMaxVersion                    string      `hcl:"tls_max_version"`
	TLSCipherSuites                  []uint16    `hcl:"-"`
	TLSCipherSuitesRaw               string      `hcl:"tls_cipher_suites"`
	TLSRequireAndVerifyClientCert    bool        `hcl:"-"`
	TLSRequireAndVerifyClientCertRaw interface{} `hcl:"tls_require_and_verify_client_cert"`
	TLSClientCAFile                  string      `hcl:"tls_client_ca_file"`
	TLSDisableClientCerts            bool        `hcl:"-"`
	TLSDisableClientCertsRaw         interface{} `hcl:"tls_disable_client_certs"`

	HTTPReadTimeout          time.Duration `hcl:"-"`
	HTTPReadTimeoutRaw       interface{}   `hcl:"http_read_timeout"`
	HTTPReadHeaderTimeout    time.Duration `hcl:"-"`
	HTTPReadHeaderTimeoutRaw interface{}   `hcl:"http_read_header_timeout"`
	HTTPWriteTimeout         time.Duration `hcl:"-"`
	HTTPWriteTimeoutRaw      interface{}   `hcl:"http_write_timeout"`
	HTTPIdleTimeout          time.Duration `hcl:"-"`
	HTTPIdleTimeoutRaw       interface{}   `hcl:"http_idle_timeout"`

	ProxyProtocolBehavior           string                        `hcl:"proxy_protocol_behavior"`
	ProxyProtocolAuthorizedAddrs    []*sockaddr.SockAddrMarshaler `hcl:"-"`
	ProxyProtocolAuthorizedAddrsRaw interface{}                   `hcl:"proxy_protocol_authorized_addrs,alias:ProxyProtocolAuthorizedAddrs"`

	XForwardedForAuthorizedAddrs          []*sockaddr.SockAddrMarshaler `hcl:"-"`
	XForwardedForAuthorizedAddrsRaw       interface{}                   `hcl:"x_forwarded_for_authorized_addrs,alias:XForwardedForAuthorizedAddrs"`
	XForwardedForHopSkips                 int64                         `hcl:"-"`
	XForwardedForHopSkipsRaw              interface{}                   `hcl:"x_forwarded_for_hop_skips,alias:XForwardedForHopSkips"`
	XForwardedForRejectNotPresent         bool                          `hcl:"-"`
	XForwardedForRejectNotPresentRaw      interface{}                   `hcl:"x_forwarded_for_reject_not_present,alias:XForwardedForRejectNotPresent"`
	XForwardedForRejectNotAuthorized      bool                          `hcl:"-"`
	XForwardedForRejectNotAuthorizedRaw   interface{}                   `hcl:"x_forwarded_for_reject_not_authorized,alias:XForwardedForRejectNotAuthorized"`
	XForwardedForClientCertHeader         string                        `hcl:"x_forwarded_for_client_cert_header,alias:XForwardedForClientCertHeader"`
	XForwardedForClientCertHeaderDecoders string                        `hcl:"x_forwarded_for_client_cert_header_decoders,alias:XForwardedForClientCertHeaderDecoders"`

	SocketMode  string `hcl:"socket_mode"`
	SocketUser  string `hcl:"socket_user"`
	SocketGroup string `hcl:"socket_group"`

	AgentAPI *AgentAPI `hcl:"agent_api"`

	ProxyAPI *ProxyAPI `hcl:"proxy_api"`

	Telemetry              ListenerTelemetry              `hcl:"telemetry"`
	Profiling              ListenerProfiling              `hcl:"profiling"`
	InFlightRequestLogging ListenerInFlightRequestLogging `hcl:"inflight_requests_logging"`

	// RandomPort is used only for some testing purposes
	RandomPort bool `hcl:"-"`

	CorsEnabledRaw        interface{} `hcl:"cors_enabled"`
	CorsEnabled           bool        `hcl:"-"`
	CorsAllowedOrigins    []string    `hcl:"cors_allowed_origins"`
	CorsAllowedHeaders    []string    `hcl:"-"`
	CorsAllowedHeadersRaw []string    `hcl:"cors_allowed_headers,alias:cors_allowed_headers"`

	// Custom Http response headers
	CustomResponseHeaders    map[string]map[string]string `hcl:"-"`
	CustomResponseHeadersRaw interface{}                  `hcl:"custom_response_headers"`

	// ChrootNamespace will prepend the specified namespace to requests
	ChrootNamespaceRaw interface{} `hcl:"chroot_namespace"`
	ChrootNamespace    string      `hcl:"-"`

	// Per-listener redaction configuration
	RedactAddressesRaw   any  `hcl:"redact_addresses"`
	RedactAddresses      bool `hcl:"-"`
	RedactClusterNameRaw any  `hcl:"redact_cluster_name"`
	RedactClusterName    bool `hcl:"-"`
	RedactVersionRaw     any  `hcl:"redact_version"`
	RedactVersion        bool `hcl:"-"`

	// DisableReplicationStatusEndpoint disables the unauthenticated replication status endpoints
	DisableReplicationStatusEndpointsRaw interface{} `hcl:"disable_replication_status_endpoints"`
	DisableReplicationStatusEndpoints    bool        `hcl:"-"`

	// DisableRequestLimiter allows per-listener disabling of the Request Limiter.
	DisableRequestLimiterRaw any  `hcl:"disable_request_limiter"`
	DisableRequestLimiter    bool `hcl:"-"`
}

// AgentAPI allows users to select which parts of the Agent API they want enabled.
type AgentAPI struct {
	EnableQuit bool `hcl:"enable_quit"`
}

// ProxyAPI allows users to select which parts of the Vault Proxy API they want enabled.
type ProxyAPI struct {
	EnableQuit bool `hcl:"enable_quit"`
}

func (l *Listener) GoString() string {
	return fmt.Sprintf("*%#v", *l)
}

func (l *Listener) Validate(path string) []ConfigError {
	results := append(ValidateUnusedFields(l.UnusedKeys, path), ValidateUnusedFields(l.Telemetry.UnusedKeys, path)...)
	return append(results, ValidateUnusedFields(l.Profiling.UnusedKeys, path)...)
}

// ParseSingleIPTemplate is used as a helper function to parse out a single IP
// address from a config parameter.
// If the input doesn't appear to contain the 'template' format,
// it will return the specified input unchanged.
func ParseSingleIPTemplate(ipTmpl string) (string, error) {
	r := regexp.MustCompile("{{.*?}}")
	if !r.MatchString(ipTmpl) {
		return NormalizeAddr(ipTmpl), nil
	}

	out, err := template.Parse(ipTmpl)
	if err != nil {
		return "", fmt.Errorf("unable to parse address template %q: %v", ipTmpl, err)
	}

	ips := strings.Split(out, " ")
	switch len(ips) {
	case 0:
		return "", errors.New("no addresses found, please configure one")
	case 1:
		return strings.TrimSpace(ips[0]), nil
	default:
		return "", fmt.Errorf("multiple addresses found (%q), please configure one", out)
	}
}

// ParseListeners attempts to parse the AST list of objects into listeners.
func ParseListeners(list *ast.ObjectList) ([]*Listener, error) {
	listeners := make([]*Listener, len(list.Items))

	for i, item := range list.Items {
		l, err := parseListener(item)
		if err != nil {
			return nil, multierror.Prefix(err, fmt.Sprintf("listeners.%d:", i))
		}
		listeners[i] = l
	}

	return listeners, nil
}

// parseListener attempts to parse the AST object into a listener.
func parseListener(item *ast.ObjectItem) (*Listener, error) {
	var l *Listener
	var err error

	// Decode the current item
	if err = hcl.DecodeObject(&l, item.Val); err != nil {
		return nil, err
	}

	// Parse and update address if required.
	if l.Address, err = ParseSingleIPTemplate(l.Address); err != nil {
		return nil, err
	}

	// Parse and update cluster address if required.
	if l.ClusterAddress, err = ParseSingleIPTemplate(l.ClusterAddress); err != nil {
		return nil, err
	}

	// Get the values for sanitizing
	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return nil, err
	}
	l.RawConfig = m

	// Parse type, but supply a fallback if type wasn't set.
	var fallbackType string
	if len(item.Keys) == 1 {
		fallbackType = strings.ToLower(item.Keys[0].Token.Value().(string))
	}

	if err = l.parseType(fallbackType); err != nil {
		return nil, err
	}

	// Parse out each set off settings for the listener.
	for _, parser := range []func() error{
		l.parseRequestSettings,
		l.parseTLSSettings,
		l.parseHTTPTimeoutSettings,
		l.parseProxySettings,
		l.parseForwardedForSettings,
		l.parseTelemetrySettings,
		l.parseProfilingSettings,
		l.parseInFlightRequestSettings,
		l.parseCORSSettings,
		l.parseHTTPHeaderSettings,
		l.parseChrootNamespaceSettings,
		l.parseRedactionSettings,
		l.parseDisableReplicationStatusEndpointSettings,
		l.parseDisableRequestLimiter,
	} {
		err := parser()
		if err != nil {
			return nil, err
		}
	}

	return l, nil
}

// Normalize returns the lower case string version of a listener type.
func (t ListenerType) Normalize() ListenerType {
	return ListenerType(strings.ToLower(string(t)))
}

// String returns the string version of a listener type.
func (t ListenerType) String() string {
	return string(t.Normalize())
}

// parseAndClearBool parses a raw setting as a bool configuration parameter. If
// the raw value is successfully parsed, the parsedSetting argument is set to it
// and the rawSetting argument is cleared. Otherwise, the rawSetting argument is
// left unchanged and an error is returned.
func parseAndClearBool(rawSetting *interface{}, parsedSetting *bool) error {
	var err error

	if *rawSetting != nil {
		*parsedSetting, err = parseutil.ParseBool(*rawSetting)
		if err != nil {
			return err
		}

		*rawSetting = nil
	}

	return nil
}

// parseAndClearString parses a raw setting as a string configuration parameter.
// If the raw value is successfully parsed, the parsedSetting argument is set to
// it and the rawSetting argument is cleared. Otherwise, the rawSetting argument
// is left unchanged and an error is returned.
func parseAndClearString(rawSetting *interface{}, parsedSetting *string) error {
	var err error

	if *rawSetting != nil {
		*parsedSetting, err = parseutil.ParseString(*rawSetting)
		if err != nil {
			return err
		}

		*rawSetting = nil
	}

	return nil
}

// parseAndClearInt parses a raw setting as an integer configuration parameter.
// If the raw value is successfully parsed, the parsedSetting argument is set to
// it and the rawSetting argument is cleared. Otherwise, the rawSetting argument
// is left unchanged and an error is returned.
func parseAndClearInt(rawSetting *interface{}, parsedSetting *int64) error {
	var err error

	if *rawSetting != nil {
		*parsedSetting, err = parseutil.ParseInt(*rawSetting)
		if err != nil {
			return err
		}

		*rawSetting = nil
	}

	return nil
}

// parseAndClearDurationSecond parses a raw setting as a time duration
// configuration parameter. If the raw value is successfully parsed, the
// parsedSetting argument is set to it and the rawSetting argument is cleared.
// Otherwise, the rawSetting argument is left unchanged and an error is
// returned.
func parseAndClearDurationSecond(rawSetting *interface{}, parsedSetting *time.Duration) error {
	var err error

	if *rawSetting != nil {
		*parsedSetting, err = parseutil.ParseDurationSecond(*rawSetting)
		if err != nil {
			return err
		}

		*rawSetting = nil
	}

	return nil
}

// parseDisableReplicationStatusEndpointSettings attempts to parse the raw
// disable_replication_status_endpoints setting. The receiving Listener's
// DisableReplicationStatusEndpoints field will be set with the successfully
// parsed value.
func (l *Listener) parseDisableReplicationStatusEndpointSettings() error {
	if l.Type != TCP {
		return nil
	}

	if err := parseAndClearBool(&l.DisableReplicationStatusEndpointsRaw, &l.DisableReplicationStatusEndpoints); err != nil {
		return fmt.Errorf("invalid value for disable_replication_status_endpoints: %w", err)
	}

	return nil
}

// parseDisableRequestLimiter attempts to parse the raw disable_request_limiter
// setting. The receiving Listener's DisableRequestLimiter field will be set
// with the successfully parsed value or return an error
func (l *Listener) parseDisableRequestLimiter() error {
	if err := parseAndClearBool(&l.DisableRequestLimiterRaw, &l.DisableRequestLimiter); err != nil {
		return fmt.Errorf("invalid value for disable_request_limiter: %w", err)
	}

	return nil
}

// parseChrootNamespace attempts to parse the raw listener chroot namespace settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseChrootNamespaceSettings() error {
	var (
		err     error
		setting string
	)

	err = parseAndClearString(&l.ChrootNamespaceRaw, &setting)
	if err != nil {
		return fmt.Errorf("invalid value for chroot_namespace: %w", err)
	}

	l.ChrootNamespace = namespace.Canonicalize(setting)

	return nil
}

// parseType attempts to sanitize and validate the type set on the listener.
// If the listener has no type set, the fallback value will be used.
// The state of the listener will be modified.
func (l *Listener) parseType(fallback string) error {
	switch {
	case l.Type != "":
	case fallback != "":
	default:
		return errors.New("listener type must be specified")
	}

	// Use type if available, otherwise fall back.
	rawType := l.Type
	if rawType == "" {
		rawType = ListenerType(fallback)
	}

	parsedType := rawType.Normalize()

	// Sanity check the values
	switch parsedType {
	case TCP, Unix:
	default:
		return fmt.Errorf("unsupported listener type %q", parsedType)
	}

	l.Type = parsedType

	return nil
}

// parseRequestSettings attempts to parse the raw listener request settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseRequestSettings() error {
	if err := parseAndClearInt(&l.MaxRequestSizeRaw, &l.MaxRequestSize); err != nil {
		return fmt.Errorf("error parsing max_request_size: %w", err)
	}

	if l.MaxRequestDurationRaw != nil {
		maxRequestDuration, err := parseutil.ParseDurationSecond(l.MaxRequestDurationRaw)
		if err != nil {
			return fmt.Errorf("error parsing max_request_duration: %w", err)
		}

		if maxRequestDuration < 0 {
			return errors.New("max_request_duration cannot be negative")
		}

		l.MaxRequestDuration = maxRequestDuration
		l.MaxRequestDurationRaw = nil
	}

	if err := parseAndClearBool(&l.RequireRequestHeaderRaw, &l.RequireRequestHeader); err != nil {
		return fmt.Errorf("invalid value for require_request_header: %w", err)
	}

	if err := parseAndClearBool(&l.DisableRequestLimiterRaw, &l.DisableRequestLimiter); err != nil {
		return fmt.Errorf("invalid value for disable_request_limiter: %w", err)
	}

	return nil
}

// parseTLSSettings attempts to parse the raw listener TLS settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseTLSSettings() error {
	if err := parseAndClearBool(&l.TLSDisableRaw, &l.TLSDisable); err != nil {
		return fmt.Errorf("invalid value for tls_disable: %w", err)
	}

	if l.TLSCipherSuitesRaw != "" {
		tlsCipherSuites, err := tlsutil.ParseCiphers(l.TLSCipherSuitesRaw)
		if err != nil {
			return fmt.Errorf("invalid value for tls_cipher_suites: %w", err)
		}
		l.TLSCipherSuites = tlsCipherSuites
	}

	if err := parseAndClearBool(&l.TLSRequireAndVerifyClientCertRaw, &l.TLSRequireAndVerifyClientCert); err != nil {
		return fmt.Errorf("invalid value for tls_require_and_verify_client_cert: %w", err)
	}

	if err := parseAndClearBool(&l.TLSDisableClientCertsRaw, &l.TLSDisableClientCerts); err != nil {
		return fmt.Errorf("invalid value for tls_disable_client_certs: %w", err)
	}

	// Clear raw values after successful parsing.
	l.TLSCipherSuitesRaw = ""

	return nil
}

// parseHTTPHeaderSettings attempts to parse the raw listener HTTP header settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseHTTPHeaderSettings() error {
	// Custom response headers are only supported by TCP listeners.
	// Clear raw data and return early if it was something else.
	if l.Type != TCP {
		l.CustomResponseHeadersRaw = nil
		return nil
	}

	// if CustomResponseHeadersRaw is nil, we still need to set the default headers
	customHeadersMap, err := ParseCustomResponseHeaders(l.CustomResponseHeadersRaw)
	if err != nil {
		return fmt.Errorf("failed to parse custom_response_headers: %w", err)
	}

	l.CustomResponseHeaders = customHeadersMap
	l.CustomResponseHeadersRaw = nil

	return nil
}

// parseHTTPTimeoutSettings attempts to parse the raw listener HTTP timeout settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseHTTPTimeoutSettings() error {
	if err := parseAndClearDurationSecond(&l.HTTPReadTimeoutRaw, &l.HTTPReadTimeout); err != nil {
		return fmt.Errorf("error parsing http_read_timeout: %w", err)
	}

	if err := parseAndClearDurationSecond(&l.HTTPReadHeaderTimeoutRaw, &l.HTTPReadHeaderTimeout); err != nil {
		return fmt.Errorf("error parsing http_read_header_timeout: %w", err)
	}

	if err := parseAndClearDurationSecond(&l.HTTPWriteTimeoutRaw, &l.HTTPWriteTimeout); err != nil {
		return fmt.Errorf("error parsing http_write_timeout: %w", err)
	}

	if err := parseAndClearDurationSecond(&l.HTTPIdleTimeoutRaw, &l.HTTPIdleTimeout); err != nil {
		return fmt.Errorf("error parsing http_idle_timeout: %w", err)
	}

	return nil
}

// parseProxySettings attempts to parse the raw listener proxy settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseProxySettings() error {
	var err error

	if l.ProxyProtocolAuthorizedAddrsRaw != nil {
		l.ProxyProtocolAuthorizedAddrs, err = parseutil.ParseAddrs(l.ProxyProtocolAuthorizedAddrsRaw)
		if err != nil {
			return fmt.Errorf("error parsing proxy_protocol_authorized_addrs: %w", err)
		}
	}

	// Validation/sanity check on allowed settings for behavior.
	switch l.ProxyProtocolBehavior {
	case "allow_authorized", "deny_unauthorized", "use_always", "":
		// Ignore these cases, they're all valid values.
		// In the case of 'allow_authorized' and 'deny_unauthorized', we don't need
		// to check how many addresses we have in ProxyProtocolAuthorizedAddrs
		// as parseutil.ParseAddrs returns "one or more addresses" (or an error)
		// so we'd have returned earlier.
	default:
		return fmt.Errorf("unsupported value supplied for proxy_protocol_behavior: %q", l.ProxyProtocolBehavior)
	}

	// Clear raw values after successful parsing.
	l.ProxyProtocolAuthorizedAddrsRaw = nil

	return nil
}

// parseForwardedForSettings attempts to parse the raw listener x-forwarded-for settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseForwardedForSettings() error {
	var err error

	if l.XForwardedForAuthorizedAddrsRaw != nil {
		if l.XForwardedForAuthorizedAddrs, err = parseutil.ParseAddrs(l.XForwardedForAuthorizedAddrsRaw); err != nil {
			return fmt.Errorf("error parsing x_forwarded_for_authorized_addrs: %w", err)
		}
	}

	if l.XForwardedForHopSkipsRaw != nil {
		if l.XForwardedForHopSkips, err = parseutil.ParseInt(l.XForwardedForHopSkipsRaw); err != nil {
			return fmt.Errorf("error parsing x_forwarded_for_hop_skips: %w", err)
		}

		if l.XForwardedForHopSkips < 0 {
			return fmt.Errorf("x_forwarded_for_hop_skips cannot be negative but set to %d", l.XForwardedForHopSkips)
		}

		l.XForwardedForHopSkipsRaw = nil
	}

	if err := parseAndClearBool(&l.XForwardedForRejectNotAuthorizedRaw, &l.XForwardedForRejectNotAuthorized); err != nil {
		return fmt.Errorf("invalid value for x_forwarded_for_reject_not_authorized: %w", err)
	}

	if err := parseAndClearBool(&l.XForwardedForRejectNotPresentRaw, &l.XForwardedForRejectNotPresent); err != nil {
		return fmt.Errorf("invalid value for x_forwarded_for_reject_not_present: %w", err)
	}

	// Clear raw values after successful parsing.
	l.XForwardedForAuthorizedAddrsRaw = nil

	return nil
}

// parseTelemetrySettings attempts to parse the raw listener telemetry settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseTelemetrySettings() error {
	if err := parseAndClearBool(&l.Telemetry.UnauthenticatedMetricsAccessRaw, &l.Telemetry.UnauthenticatedMetricsAccess); err != nil {
		return fmt.Errorf("invalid value for telemetry.unauthenticated_metrics_access: %w", err)
	}

	return nil
}

// parseProfilingSettings attempts to parse the raw listener profiling settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseProfilingSettings() error {
	if err := parseAndClearBool(&l.Profiling.UnauthenticatedPProfAccessRaw, &l.Profiling.UnauthenticatedPProfAccess); err != nil {
		return fmt.Errorf("invalid value for profiling.unauthenticated_pprof_access: %w", err)
	}

	return nil
}

// parseInFlightRequestSettings attempts to parse the raw listener in-flight request logging settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseInFlightRequestSettings() error {
	if err := parseAndClearBool(&l.InFlightRequestLogging.UnauthenticatedInFlightAccessRaw, &l.InFlightRequestLogging.UnauthenticatedInFlightAccess); err != nil {
		return fmt.Errorf("invalid value for inflight_requests_logging.unauthenticated_in_flight_requests_access: %w", err)
	}

	return nil
}

// parseCORSSettings attempts to parse the raw listener CORS settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseCORSSettings() error {
	if err := parseAndClearBool(&l.CorsEnabledRaw, &l.CorsEnabled); err != nil {
		return fmt.Errorf("invalid value for cors_enabled: %w", err)
	}

	if strutil.StrListContains(l.CorsAllowedOrigins, "*") && len(l.CorsAllowedOrigins) > 1 {
		return errors.New("cors_allowed_origins must only contain a wildcard or only non-wildcard values")
	}

	if len(l.CorsAllowedHeadersRaw) > 0 {
		for _, header := range l.CorsAllowedHeadersRaw {
			l.CorsAllowedHeaders = append(l.CorsAllowedHeaders, textproto.CanonicalMIMEHeaderKey(header))
		}
	}

	l.CorsAllowedHeadersRaw = nil

	return nil
}

// parseRedactionSettings attempts to parse the raw listener redaction settings.
// The state of the listener will be modified, raw data will be cleared upon
// successful parsing.
func (l *Listener) parseRedactionSettings() error {
	// Redaction is only supported on TCP listeners.
	// Clear raw data and return early if it was something else.
	if l.Type != TCP {
		l.RedactAddressesRaw = nil
		l.RedactClusterNameRaw = nil
		l.RedactVersionRaw = nil

		return nil
	}

	var err error

	if l.RedactAddressesRaw != nil {
		if l.RedactAddresses, err = parseutil.ParseBool(l.RedactAddressesRaw); err != nil {
			return fmt.Errorf("invalid value for redact_addresses: %w", err)
		}
	}
	if l.RedactClusterNameRaw != nil {
		if l.RedactClusterName, err = parseutil.ParseBool(l.RedactClusterNameRaw); err != nil {
			return fmt.Errorf("invalid value for redact_cluster_name: %w", err)
		}
	}
	if l.RedactVersionRaw != nil {
		if l.RedactVersion, err = parseutil.ParseBool(l.RedactVersionRaw); err != nil {
			return fmt.Errorf("invalid value for redact_version: %w", err)
		}
	}

	l.RedactAddressesRaw = nil
	l.RedactClusterNameRaw = nil
	l.RedactVersionRaw = nil

	return nil
}
