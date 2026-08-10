// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/logical"
	uberatomic "go.uber.org/atomic"
)

const (
	// DefaultBillingRetentionMonths is the default number of months of billing data to retain.
	// This includes the current month plus previous months (e.g., 37 = current + 36 previous months).
	DefaultBillingRetentionMonths = 37

	// MinBillingRetentionMonths is the minimum allowed retention period (13 months = 1 year + current month)
	MinBillingRetentionMonths = 13

	// MaxBillingRetentionMonths is the maximum allowed retention period (72 months = 6 years)
	MaxBillingRetentionMonths = 72

	// DefaultAttributionRetentionMonths is the default number of months of attribution data to retain.
	DefaultAttributionRetentionMonths = 37

	BillingWriteInterval = 10 * time.Minute
	// pluginCountsSendTimeout is the timeout for sending plugin counts to the active node
	PluginCountsSendTimeout = 30 * time.Second
	// pluginCountsStandbyTime is how long to wait before sending plugin counts from a perf standby
	PluginCountsStandbyTime = 10 * time.Minute

	BillingSubPath                          = "billing/"
	BillingConfigPath                       = "config"
	ReplicatedPrefix                        = "replicated/"
	RoleHWMCountsHWM                        = "maxRoleCounts/"
	TotpHWMCountsHWM                        = "maxTotpCounts/"
	KvHWMCountsHWM                          = "maxKvCounts/"
	KmseHWMCountsHWM                        = "maxKmseCounts/"
	TransitDataProtectionCallCountsPrefix   = "transitDataProtectionCallCounts/"
	TransformDataProtectionCallCountsPrefix = "transformDataProtectionCallCounts/"
	GcpKmsDataProtectionCallCountsPrefix    = "gcpKmsDataProtectionCallCounts/"
	LocalPrefix                             = "local/"
	ThirdPartyPluginsPrefix                 = "thirdPartyPluginCounts/"
	KmipEnabledPrefix                       = "kmipEnabled/"
	PkiDurationAdjustedCountPrefix          = "normalizedCertsIssued/"
	SpiffeJwtNormalizedTokenUnits           = "spiffeJwtNormalizedTokenUnits/"
	MetricsLastUpdatedAtPrefix              = "metricsLastUpdatedAt/"
	SSHCertificateMetric                    = "ssh/normalized-certs-issued"
	SSHOTPMetric                            = "ssh/credential-count"
	OidcDurationAdjustedCountPrefix         = "oidcNormalizedTokenUnits/"
	ExternalCaDurationAdjustedCountPrefix   = "externalCaNormalizedCertsIssued/"

	AttributionMaxPrefix = "attribution/maximum/"

	// Role and managed key sub-types for storing attribution data of consumption billing metrics
	AWSDynamicRoles            = "aws_dynamic"
	AWSStaticRoles             = "aws_static"
	AzureDynamicRoles          = "azure_dynamic"
	AzureStaticRoles           = "azure_static"
	DatabaseDynamicRoles       = "database_dynamic"
	DatabaseStaticRoles        = "database_static"
	GCPRolesets                = "gcp_dynamic"
	GCPStaticAccounts          = "gcp_static"
	GCPImpersonatedAccounts    = "gcp_impersonated"
	LDAPDynamicRoles           = "ldap_dynamic"
	LDAPStaticRoles            = "ldap_static"
	LDAPLibrarySets            = "ldap_library_sets"
	OpenLDAPDynamicRoles       = "openldap_dynamic"
	OpenLDAPLibrarySets        = "openldap_library_sets"
	OpenLDAPStaticRoles        = "openldap_static"
	AlicloudDynamicRoles       = "alicloud_dynamic"
	RabbitMQDynamicRoles       = "rabbitmq_dynamic"
	ConsulDynamicRoles         = "consul_dynamic"
	NomadDynamicRoles          = "nomad_dynamic"
	KubernetesDynamicRoles     = "kubernetes_dynamic"
	MongoDBAtlasDynamicRoles   = "mongodb_atlas_dynamic"
	TerraformCloudDynamicRoles = "terraformcloud_dynamic"
	OSLocalAccountRoles        = "os_local_account_static"

	TotpKeys = "totp"
	KmseKeys = "kmse"
)

var BillingMonthStorageFormat = "%s%d/%02d/%s" // e.g replicated/2026/01/maxKvCounts/

type ConsumptionBilling struct {
	// BillingStorageLock controls access to the billing storage paths
	BillingStorageLock sync.RWMutex

	BillingConfig            BillingConfig
	DataProtectionCallCounts DataProtectionCallCounts
	Logger                   log.Logger

	// KmipSeenEnabledThisMonth tracks whether KMIP has been enabled during the current billing month.
	// This is used to avoid scanning all mounts every 10 minutes for KMIP billing detection.
	KmipSeenEnabledThisMonth atomic.Bool

	IdentityTokenUnits IdentityTokenUnits

	// ExternalCaCertUnits tracks duration-adjusted PKI external CA certificate units
	ExternalCaCertUnits *uberatomic.Float64
}

type BillingConfig struct {
	// DisableConsumptionBilling disables the periodic consumption billing metrics
	// worker that walks the entire KV store every 10 minutes. When set to true,
	// the billing worker is not registered during post-unseal, preventing the
	// memory and performance impact of the periodic KV store enumeration.
	// This can be set via the VAULT_DISABLE_CONSUMPTION_BILLING environment variable.
	DisableConsumptionBilling bool
	// For testing purposes. The cadence at which billing metrics are updated
	MetricsUpdateCadence time.Duration
	// For testing purposes. The cadence at which plugin counts are sent from perf standby to active
	PluginCountsSendCadence time.Duration
	// For testin purposes. TestOverrideClock holds a custom clock to modify time.Now, time.Ticker, time.Timer.
	// If nil, the default functions from the time package are used
	TestOverrideClock timeutil.Clock
}
