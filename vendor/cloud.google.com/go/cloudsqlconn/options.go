// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudsqlconn

import (
	"context"
	"crypto/rsa"
	"net"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/cloudsqlconn/debug"
	"cloud.google.com/go/cloudsqlconn/errtype"
	"cloud.google.com/go/cloudsqlconn/instance"
	"cloud.google.com/go/cloudsqlconn/internal/cloudsql"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	apiopt "google.golang.org/api/option"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// An Option is an option for configuring a Dialer.
type Option func(d *dialerConfig)

type dialerConfig struct {
	rsaKey                 *rsa.PrivateKey
	sqladminOpts           []apiopt.ClientOption
	dialOpts               []DialOption
	dialFunc               func(ctx context.Context, network, addr string) (net.Conn, error)
	refreshTimeout         time.Duration
	useIAMAuthN            bool
	logger                 debug.ContextLogger
	lazyRefresh            bool
	iamLoginTokenSource    oauth2.TokenSource
	useragents             []string
	credentialsUniverse    string
	serviceUniverse        string
	setAdminAPIEndpoint    bool
	setUniverseDomain      bool
	setCredentials         bool
	setTokenSource         bool
	setIAMAuthNTokenSource bool
	resolver               instance.ConnectionNameResolver
	failoverPeriod         time.Duration
	// err tracks any dialer options that may have failed.
	err error
}

// WithOptions turns a list of Option's into a single Option.
func WithOptions(opts ...Option) Option {
	return func(d *dialerConfig) {
		for _, opt := range opts {
			opt(d)
		}
	}
}

// WithCredentialsFile returns an Option that specifies a service account
// or refresh token JSON credentials file to be used as the basis for
// authentication.
func WithCredentialsFile(filename string) Option {
	return func(d *dialerConfig) {
		b, err := os.ReadFile(filename)
		if err != nil {
			d.err = errtype.NewConfigError(err.Error(), "n/a")
			return
		}
		opt := WithCredentialsJSON(b)
		opt(d)
	}
}

// WithCredentialsJSON returns an Option that specifies a service account
// or refresh token JSON credentials to be used as the basis for authentication.
func WithCredentialsJSON(b []byte) Option {
	return func(d *dialerConfig) {
		c, err := google.CredentialsFromJSON(context.Background(), b, sqladmin.SqlserviceAdminScope)
		if err != nil {
			d.err = errtype.NewConfigError(err.Error(), "n/a")
			return
		}
		ud, err := c.GetUniverseDomain()
		if err != nil {
			d.err = errtype.NewConfigError(err.Error(), "n/a")
			return
		}
		d.credentialsUniverse = ud
		d.sqladminOpts = append(d.sqladminOpts, apiopt.WithCredentials(c))

		// Create another set of credentials scoped to login only
		scoped, err := google.CredentialsFromJSON(context.Background(), b, iamLoginScope)
		if err != nil {
			d.err = errtype.NewConfigError(err.Error(), "n/a")
			return
		}
		d.iamLoginTokenSource = scoped.TokenSource
		d.setCredentials = true
	}
}

// WithUserAgent returns an Option that sets the User-Agent.
func WithUserAgent(ua string) Option {
	return func(d *dialerConfig) {
		d.useragents = append(d.useragents, ua)
	}
}

// WithDefaultDialOptions returns an Option that specifies the default
// DialOptions used.
func WithDefaultDialOptions(opts ...DialOption) Option {
	return func(d *dialerConfig) {
		d.dialOpts = append(d.dialOpts, opts...)
	}
}

// WithTokenSource returns an Option that specifies an OAuth2 token source to be
// used as the basis for authentication.
//
// When Auth IAM AuthN is enabled, use WithIAMAuthNTokenSources to set the token
// source for login tokens separately from the API client token source.
// WithTokenSource should not be used with WithIAMAuthNTokenSources.
func WithTokenSource(s oauth2.TokenSource) Option {
	return func(d *dialerConfig) {
		d.setTokenSource = true
		d.setCredentials = true
		d.sqladminOpts = append(d.sqladminOpts, apiopt.WithTokenSource(s))
	}
}

// WithIAMAuthNTokenSources sets the oauth2.TokenSource for the API client and a
// second token source for IAM AuthN login tokens. The API client token source
// should have the following scopes:
//
//  1. https://www.googleapis.com/auth/sqlservice.admin, and
//  2. https://www.googleapis.com/auth/cloud-platform
//
// The IAM AuthN token source on the other hand should only have:
//
//  1. https://www.googleapis.com/auth/sqlservice.login.
//
// Prefer this option over WithTokenSource when using IAM AuthN which does not
// distinguish between the two token sources. WithIAMAuthNTokenSources should
// not be used with WithTokenSource.
func WithIAMAuthNTokenSources(apiTS, iamLoginTS oauth2.TokenSource) Option {
	return func(d *dialerConfig) {
		d.setIAMAuthNTokenSource = true
		d.setCredentials = true
		d.iamLoginTokenSource = iamLoginTS
		d.sqladminOpts = append(d.sqladminOpts, apiopt.WithTokenSource(apiTS))
	}
}

// WithRSAKey returns an Option that specifies a rsa.PrivateKey used to represent the client.
func WithRSAKey(k *rsa.PrivateKey) Option {
	return func(d *dialerConfig) {
		d.rsaKey = k
	}
}

// WithRefreshTimeout returns an Option that sets a timeout on refresh
// operations. Defaults to 60s.
func WithRefreshTimeout(t time.Duration) Option {
	return func(d *dialerConfig) {
		d.refreshTimeout = t
	}
}

// WithHTTPClient configures the underlying SQL Admin API client with the
// provided HTTP client. This option is generally unnecessary except for
// advanced use-cases.
func WithHTTPClient(client *http.Client) Option {
	return func(d *dialerConfig) {
		d.sqladminOpts = append(d.sqladminOpts, apiopt.WithHTTPClient(client))
	}
}

// WithAdminAPIEndpoint configures the underlying SQL Admin API client to use
// the provided URL.
func WithAdminAPIEndpoint(url string) Option {
	return func(d *dialerConfig) {
		d.sqladminOpts = append(d.sqladminOpts, apiopt.WithEndpoint(url))
		d.setAdminAPIEndpoint = true
		d.serviceUniverse = ""
	}
}

// WithUniverseDomain configures the underlying SQL Admin API client to use
// the provided universe domain. Enables Trusted Partner Cloud (TPC).
func WithUniverseDomain(ud string) Option {
	return func(d *dialerConfig) {
		d.sqladminOpts = append(d.sqladminOpts, apiopt.WithUniverseDomain(ud))
		d.serviceUniverse = ud
		d.setUniverseDomain = true
	}
}

// WithQuotaProject returns an Option that specifies the project used for quota and billing purposes.
func WithQuotaProject(p string) Option {
	return func(cfg *dialerConfig) {
		cfg.sqladminOpts = append(cfg.sqladminOpts, apiopt.WithQuotaProject(p))
	}
}

// WithDialFunc configures the function used to connect to the address on the
// named network. This option is generally unnecessary except for advanced
// use-cases. The function is used for all invocations of Dial. To configure
// a dial function per individual calls to dial, use WithOneOffDialFunc.
func WithDialFunc(dial func(ctx context.Context, network, addr string) (net.Conn, error)) Option {
	return func(d *dialerConfig) {
		d.dialFunc = dial
	}
}

// WithIAMAuthN enables automatic IAM Authentication. If no token source has
// been configured (such as with WithTokenSource, WithCredentialsFile, etc), the
// dialer will use the default token source as defined by
// https://pkg.go.dev/golang.org/x/oauth2/google#FindDefaultCredentialsWithParams.
//
// For documentation on automatic IAM Authentication, see
// https://cloud.google.com/sql/docs/postgres/authentication.
func WithIAMAuthN() Option {
	return func(d *dialerConfig) {
		d.useIAMAuthN = true
	}
}

// WithResolver replaces the default resolver with an alternate
// implementation to resolve the name in the database DSN to a Cloud SQL
// instance.
func WithResolver(r instance.ConnectionNameResolver) Option {
	return func(d *dialerConfig) {
		d.resolver = r
	}
}

// WithDNSResolver replaces the default resolver (which only resolves instance
// names) with the DNSResolver, which will attempt to first parse the instance
// name, and then will attempt to resolve the DNS TXT record to determine
// the instance name.
//
// First, add a record for your Cloud SQL instance to a **private** DNS server
// or a private Google Cloud DNS Zone used by your application.
//
// **Note:** You are strongly discouraged from adding DNS records for your
// Cloud SQL instances to a public DNS server. This would allow anyone on the
// internet to discover the Cloud SQL instance name.
//
// For example: suppose you wanted to use the domain name
// `prod-db.mycompany.example.com` to connect to your database instance
// `my-project:region:my-instance`. You would create the following DNS record:
//
//   - Record type: `TXT`
//   - Name: `prod-db.mycompany.example.com` – This is the domain name used by
//     the application
//   - Value: `my-project:region:my-instance` – This is the instance name
func WithDNSResolver() Option {
	return func(d *dialerConfig) {
		d.resolver = cloudsql.DNSResolver
	}
}

// WithFailoverPeriod will cause the connector to periodically check the SRV DNS
// records of instance configured using DNS names. By default, this is 30
// seconds. If this is set to 0, the connector will only check for domain name
// changes when establishing a new connection.
func WithFailoverPeriod(f time.Duration) Option {
	return func(d *dialerConfig) {
		d.failoverPeriod = f
	}
}

type debugLoggerWithoutContext struct {
	logger debug.Logger
}

// Debugf implements debug.ContextLogger.
func (d *debugLoggerWithoutContext) Debugf(_ context.Context, format string, args ...interface{}) {
	d.logger.Debugf(format, args...)
}

var _ debug.ContextLogger = new(debugLoggerWithoutContext)

// WithDebugLogger configures a debug lgoger for reporting on internal
// operations. By default the debug logger is disabled.
//
// Prefer WithContextDebugLogger instead
func WithDebugLogger(l debug.Logger) Option {
	return func(d *dialerConfig) {
		d.logger = &debugLoggerWithoutContext{l}
	}
}

// WithContextDebugLogger configures a debug logger for reporting on internal
// operations. By default the debug logger is disabled.
func WithContextDebugLogger(l debug.ContextLogger) Option {
	return func(d *dialerConfig) {
		d.logger = l
	}
}

// WithLazyRefresh configures the dialer to refresh certificates on an
// as-needed basis. If a certificate is expired when a connection request
// occurs, the Go Connector will block the attempt and refresh the certificate
// immediately. This option is useful when running the Go Connector in
// environments where the CPU may be throttled, thus preventing a background
// goroutine from running consistently (e.g., in Cloud Run the CPU is throttled
// outside of a request context causing the background refresh to fail).
func WithLazyRefresh() Option {
	return func(d *dialerConfig) {
		d.lazyRefresh = true
	}
}

// A DialOption is an option for configuring how a Dialer's Dial call is executed.
type DialOption func(d *dialConfig)

type dialConfig struct {
	dialFunc     func(ctx context.Context, network, addr string) (net.Conn, error)
	ipType       string
	tcpKeepAlive time.Duration
	useIAMAuthN  bool
}

// DialOptions turns a list of DialOption instances into an DialOption.
func DialOptions(opts ...DialOption) DialOption {
	return func(cfg *dialConfig) {
		for _, opt := range opts {
			opt(cfg)
		}
	}
}

// WithOneOffDialFunc configures the dial function on a one-off basis for an
// individual call to Dial. To configure a dial function across all invocations
// of Dial, use WithDialFunc.
func WithOneOffDialFunc(dial func(ctx context.Context, network, addr string) (net.Conn, error)) DialOption {
	return func(c *dialConfig) {
		c.dialFunc = dial
	}
}

// WithTCPKeepAlive returns a DialOption that specifies the tcp keep alive period for the connection returned by Dial.
func WithTCPKeepAlive(d time.Duration) DialOption {
	return func(cfg *dialConfig) {
		cfg.tcpKeepAlive = d
	}
}

// WithPublicIP returns a DialOption that specifies a public IP will be used to connect.
func WithPublicIP() DialOption {
	return func(cfg *dialConfig) {
		cfg.ipType = cloudsql.PublicIP
	}
}

// WithPrivateIP returns a DialOption that specifies a private IP (VPC) will be used to connect.
func WithPrivateIP() DialOption {
	return func(cfg *dialConfig) {
		cfg.ipType = cloudsql.PrivateIP
	}
}

// WithPSC returns a DialOption that specifies a PSC endpoint will be used to connect.
func WithPSC() DialOption {
	return func(cfg *dialConfig) {
		cfg.ipType = cloudsql.PSC
	}
}

// WithAutoIP returns a DialOption that selects the public IP if available and
// otherwise falls back to private IP. This option is present for backwards
// compatibility only and is not recommended for use in production.
func WithAutoIP() DialOption {
	return func(cfg *dialConfig) {
		cfg.ipType = cloudsql.AutoIP
	}
}

// WithDialIAMAuthN allows you to enable or disable IAM Authentication for this
// instance as described in the documentation for WithIAMAuthN. This value will
// override the Dialer-level configuration set with WithIAMAuthN.
//
// WARNING: This DialOption can cause a new Refresh operation to be triggered.
// Toggling this option on or off between Dials may cause increased API usage
// and/or delayed connection attempts.
func WithDialIAMAuthN(b bool) DialOption {
	return func(cfg *dialConfig) {
		cfg.useIAMAuthN = b
	}
}
