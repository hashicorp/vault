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

	Type       string
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

	XForwardedForAuthorizedAddrs        []*sockaddr.SockAddrMarshaler `hcl:"-"`
	XForwardedForAuthorizedAddrsRaw     interface{}                   `hcl:"x_forwarded_for_authorized_addrs,alias:XForwardedForAuthorizedAddrs"`
	XForwardedForHopSkips               int64                         `hcl:"-"`
	XForwardedForHopSkipsRaw            interface{}                   `hcl:"x_forwarded_for_hop_skips,alias:XForwardedForHopSkips"`
	XForwardedForRejectNotPresent       bool                          `hcl:"-"`
	XForwardedForRejectNotPresentRaw    interface{}                   `hcl:"x_forwarded_for_reject_not_present,alias:XForwardedForRejectNotPresent"`
	XForwardedForRejectNotAuthorized    bool                          `hcl:"-"`
	XForwardedForRejectNotAuthorizedRaw interface{}                   `hcl:"x_forwarded_for_reject_not_authorized,alias:XForwardedForRejectNotAuthorized"`

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

func ParseListeners(result *SharedConfig, list *ast.ObjectList) error {
	var err error
	result.Listeners = make([]*Listener, 0, len(list.Items))
	for i, item := range list.Items {
		var l Listener
		if err := hcl.DecodeObject(&l, item.Val); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("listeners.%d:", i))
		}
		if rendered, err := ParseSingleIPTemplate(l.Address); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("listeners.%d:", i))
		} else {
			l.Address = rendered
		}
		if rendered, err := ParseSingleIPTemplate(l.ClusterAddress); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("listeners.%d:", i))
		} else {
			l.ClusterAddress = rendered
		}

		// Hacky way, for now, to get the values we want for sanitizing
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("listeners.%d:", i))
		}
		l.RawConfig = m

		// Base values
		{
			switch {
			case l.Type != "":
			case len(item.Keys) == 1:
				l.Type = strings.ToLower(item.Keys[0].Token.Value().(string))
			default:
				return multierror.Prefix(errors.New("listener type must be specified"), fmt.Sprintf("listeners.%d:", i))
			}

			l.Type = strings.ToLower(l.Type)
			switch l.Type {
			case "tcp", "unix":
				result.found(l.Type, l.Type)
			default:
				return multierror.Prefix(fmt.Errorf("unsupported listener type %q", l.Type), fmt.Sprintf("listeners.%d:", i))
			}

			if l.PurposeRaw != nil {
				if l.Purpose, err = parseutil.ParseCommaStringSlice(l.PurposeRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("unable to parse 'purpose' in listener type %q: %w", l.Type, err), fmt.Sprintf("listeners.%d:", i))
				}
				for i, v := range l.Purpose {
					l.Purpose[i] = strings.ToLower(v)
				}

				l.PurposeRaw = nil
			}

			switch l.Role {
			case "default", "metrics_only", "":
				result.found(l.Type, l.Type)
			default:
				return multierror.Prefix(fmt.Errorf("unsupported listener role %q", l.Role), fmt.Sprintf("listeners.%d:", i))
			}
		}
		// Request Parameters
		{
			if l.MaxRequestSizeRaw != nil {
				if l.MaxRequestSize, err = parseutil.ParseInt(l.MaxRequestSizeRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing max_request_size: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.MaxRequestSizeRaw = nil
			}

			if l.MaxRequestDurationRaw != nil {
				if l.MaxRequestDuration, err = parseutil.ParseDurationSecond(l.MaxRequestDurationRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing max_request_duration: %w", err), fmt.Sprintf("listeners.%d", i))
				}
				if l.MaxRequestDuration < 0 {
					return multierror.Prefix(errors.New("max_request_duration cannot be negative"), fmt.Sprintf("listeners.%d", i))
				}

				l.MaxRequestDurationRaw = nil
			}

			if l.RequireRequestHeaderRaw != nil {
				if l.RequireRequestHeader, err = parseutil.ParseBool(l.RequireRequestHeaderRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for require_request_header: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.RequireRequestHeaderRaw = nil
			}
		}

		// TLS Parameters
		{
			if l.TLSDisableRaw != nil {
				if l.TLSDisable, err = parseutil.ParseBool(l.TLSDisableRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for tls_disable: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.TLSDisableRaw = nil
			}

			if l.TLSCipherSuitesRaw != "" {
				if l.TLSCipherSuites, err = tlsutil.ParseCiphers(l.TLSCipherSuitesRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for tls_cipher_suites: %w", err), fmt.Sprintf("listeners.%d", i))
				}
			}

			if l.TLSRequireAndVerifyClientCertRaw != nil {
				if l.TLSRequireAndVerifyClientCert, err = parseutil.ParseBool(l.TLSRequireAndVerifyClientCertRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for tls_require_and_verify_client_cert: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.TLSRequireAndVerifyClientCertRaw = nil
			}

			if l.TLSDisableClientCertsRaw != nil {
				if l.TLSDisableClientCerts, err = parseutil.ParseBool(l.TLSDisableClientCertsRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for tls_disable_client_certs: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.TLSDisableClientCertsRaw = nil
			}
		}

		// HTTP timeouts
		{
			if l.HTTPReadTimeoutRaw != nil {
				if l.HTTPReadTimeout, err = parseutil.ParseDurationSecond(l.HTTPReadTimeoutRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing http_read_timeout: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.HTTPReadTimeoutRaw = nil
			}

			if l.HTTPReadHeaderTimeoutRaw != nil {
				if l.HTTPReadHeaderTimeout, err = parseutil.ParseDurationSecond(l.HTTPReadHeaderTimeoutRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing http_read_header_timeout: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.HTTPReadHeaderTimeoutRaw = nil
			}

			if l.HTTPWriteTimeoutRaw != nil {
				if l.HTTPWriteTimeout, err = parseutil.ParseDurationSecond(l.HTTPWriteTimeoutRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing http_write_timeout: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.HTTPWriteTimeoutRaw = nil
			}

			if l.HTTPIdleTimeoutRaw != nil {
				if l.HTTPIdleTimeout, err = parseutil.ParseDurationSecond(l.HTTPIdleTimeoutRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing http_idle_timeout: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.HTTPIdleTimeoutRaw = nil
			}
		}

		// Proxy Protocol config
		{
			if l.ProxyProtocolAuthorizedAddrsRaw != nil {
				if l.ProxyProtocolAuthorizedAddrs, err = parseutil.ParseAddrs(l.ProxyProtocolAuthorizedAddrsRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing proxy_protocol_authorized_addrs: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				switch l.ProxyProtocolBehavior {
				case "allow_authorized", "deny_authorized":
					if len(l.ProxyProtocolAuthorizedAddrs) == 0 {
						return multierror.Prefix(errors.New("proxy_protocol_behavior set to allow or deny only authorized addresses but no proxy_protocol_authorized_addrs value"), fmt.Sprintf("listeners.%d", i))
					}
				}

				l.ProxyProtocolAuthorizedAddrsRaw = nil
			}
		}

		// X-Forwarded-For config
		{
			if l.XForwardedForAuthorizedAddrsRaw != nil {
				if l.XForwardedForAuthorizedAddrs, err = parseutil.ParseAddrs(l.XForwardedForAuthorizedAddrsRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing x_forwarded_for_authorized_addrs: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				for _, m := range l.XForwardedForAuthorizedAddrs {
					if _, ok := m.SockAddr.(sockaddr.UnixSock); ok {
						// X-Forwarded-For headers should only be used with proxies/load balancers, and these mechanisms
						// are TCP based, not unix socket based. If any of our parsed addresses type check as unix
						// sockets, that almost certainly means they are malformed, because the parsing logic in
						// the sockaddr library tries IPv4 first, then IPv6, then assumes it's a unix socket if it
						// has a / in it (which any IPv4 address in CIDR notation will). Therefore, error here, rather
						// than accepting it and letting Vault panic at run time.
						return multierror.Prefix(fmt.Errorf("error parsing x_forwarded_for_authorized_addrs: %v does not appear to be valid", m), fmt.Sprintf("listeners.%d", i))
					}
				}

				l.XForwardedForAuthorizedAddrsRaw = nil
			}

			if l.XForwardedForHopSkipsRaw != nil {
				if l.XForwardedForHopSkips, err = parseutil.ParseInt(l.XForwardedForHopSkipsRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing x_forwarded_for_hop_skips: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				if l.XForwardedForHopSkips < 0 {
					return multierror.Prefix(fmt.Errorf("x_forwarded_for_hop_skips cannot be negative but set to %d", l.XForwardedForHopSkips), fmt.Sprintf("listeners.%d", i))
				}

				l.XForwardedForHopSkipsRaw = nil
			}

			if l.XForwardedForRejectNotAuthorizedRaw != nil {
				if l.XForwardedForRejectNotAuthorized, err = parseutil.ParseBool(l.XForwardedForRejectNotAuthorizedRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for x_forwarded_for_reject_not_authorized: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.XForwardedForRejectNotAuthorizedRaw = nil
			}

			if l.XForwardedForRejectNotPresentRaw != nil {
				if l.XForwardedForRejectNotPresent, err = parseutil.ParseBool(l.XForwardedForRejectNotPresentRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for x_forwarded_for_reject_not_present: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.XForwardedForRejectNotPresentRaw = nil
			}
		}

		// Telemetry
		{
			if l.Telemetry.UnauthenticatedMetricsAccessRaw != nil {
				if l.Telemetry.UnauthenticatedMetricsAccess, err = parseutil.ParseBool(l.Telemetry.UnauthenticatedMetricsAccessRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for telemetry.unauthenticated_metrics_access: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.Telemetry.UnauthenticatedMetricsAccessRaw = nil
			}
		}

		// Profiling
		{
			if l.Profiling.UnauthenticatedPProfAccessRaw != nil {
				if l.Profiling.UnauthenticatedPProfAccess, err = parseutil.ParseBool(l.Profiling.UnauthenticatedPProfAccessRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for profiling.unauthenticated_pprof_access: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.Profiling.UnauthenticatedPProfAccessRaw = nil
			}
		}

		// InFlight Request logging
		{
			if l.InFlightRequestLogging.UnauthenticatedInFlightAccessRaw != nil {
				if l.InFlightRequestLogging.UnauthenticatedInFlightAccess, err = parseutil.ParseBool(l.InFlightRequestLogging.UnauthenticatedInFlightAccessRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for inflight_requests_logging.unauthenticated_in_flight_requests_access: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.InFlightRequestLogging.UnauthenticatedInFlightAccessRaw = ""
			}
		}

		// CORS
		{
			if l.CorsEnabledRaw != nil {
				if l.CorsEnabled, err = parseutil.ParseBool(l.CorsEnabledRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for cors_enabled: %w", err), fmt.Sprintf("listeners.%d", i))
				}

				l.CorsEnabledRaw = nil
			}

			if strutil.StrListContains(l.CorsAllowedOrigins, "*") && len(l.CorsAllowedOrigins) > 1 {
				return multierror.Prefix(errors.New("cors_allowed_origins must only contain a wildcard or only non-wildcard values"), fmt.Sprintf("listeners.%d", i))
			}

			if len(l.CorsAllowedHeadersRaw) > 0 {
				for _, header := range l.CorsAllowedHeadersRaw {
					l.CorsAllowedHeaders = append(l.CorsAllowedHeaders, textproto.CanonicalMIMEHeaderKey(header))
				}
			}
		}

		// HTTP Headers
		{
			// if CustomResponseHeadersRaw is nil, we still need to set the default headers
			customHeadersMap, err := ParseCustomResponseHeaders(l.CustomResponseHeadersRaw)
			if err != nil {
				return multierror.Prefix(fmt.Errorf("failed to parse custom_response_headers: %w", err), fmt.Sprintf("listeners.%d", i))
			}
			l.CustomResponseHeaders = customHeadersMap
			l.CustomResponseHeadersRaw = nil
		}

		result.Listeners = append(result.Listeners, &l)

		// Chroot Namespace
		{
			// If a valid ChrootNamespace value exists, then canonicalize the namespace value
			if l.ChrootNamespaceRaw != nil {
				if l.ChrootNamespace, err = parseutil.ParseString(l.ChrootNamespaceRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for chroot_namespace: %w", err), fmt.Sprintf("listeners.%d", i))
				} else {
					l.ChrootNamespace = namespace.Canonicalize(l.ChrootNamespace)
				}

				l.ChrootNamespaceRaw = nil
			}
		}
	}

	return nil
}

// ParseSingleIPTemplate is used as a helper function to parse out a single IP
// address from a config parameter.
// If the input doesn't appear to contain the 'template' format,
// it will return the specified input unchanged.
func ParseSingleIPTemplate(ipTmpl string) (string, error) {
	r := regexp.MustCompile("{{.*?}}")
	if !r.MatchString(ipTmpl) {
		return ipTmpl, nil
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
