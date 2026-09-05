// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
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

	// MinAttributionRetentionMonths is the minimum allowed attribution retention period.
	// A value of 0 disables attribution storage entirely and wipes all existing attribution data.
	MinAttributionRetentionMonths = 0

	// MaxAttributionRetentionMonths is the maximum allowed attribution retention period (72 months = 6 years).
	MaxAttributionRetentionMonths = 72

	// AttributionConfigPath is the storage key for the attribution retention configuration.
	AttributionConfigPath = "attribution_config"

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

	// Mount type constants for consumption billing attribution.
	MountTypeTransit   = "transit"
	MountTypeTransform = "transform"
	MountTypeGcpKms    = "gcpkms"
	MountTypeSpiffe    = "spiffe"
	MountTypeOidc      = "oidc"
	MountTypeExCa      = "external-ca"
)

var BillingMonthStorageFormat = "%s%d/%02d/%s" // e.g replicated/2026/01/maxKvCounts/

type ConsumptionBilling struct {
	// BillingStorageLock controls access to the billing storage paths
	BillingStorageLock sync.RWMutex

	BillingConfig        BillingConfig
	SecretEngineCounts   SecretEngineCounts
	Logger               log.Logger
	GetParentNamespaceID func(string) string

	// KmipSeenEnabledThisMonth tracks whether KMIP has been enabled during the current billing month.
	// This is used to avoid scanning all mounts every 10 minutes for KMIP billing detection.
	KmipSeenEnabledThisMonth atomic.Bool
}

type BillingConfig struct {
	// For testing purposes. The cadence at which billing metrics are updated
	MetricsUpdateCadence time.Duration
	// For testing purposes. The cadence at which plugin counts are sent from perf standby to active
	PluginCountsSendCadence time.Duration
	// For testing purposes. TestOverrideClock holds a custom clock to modify time.Now, time.Ticker, time.Timer.
	// If nil, the default functions from the time package are used
	TestOverrideClock timeutil.Clock
	// OnMetricsSent is called in tests to observe the proto that would be sent
	// to the control hub. It is never set in production.
	OnMetricsSent func([]byte, error)
}

func GetMonthlyBillingMetricPath(localPrefix string, now time.Time, billingMetric string) string {
	// Normalize to avoid double slashes since our prefixes include trailing "/".
	// Example: localPrefix="replicated/", billingMetric="maxKvCounts/" =>
	// "replicated/2026/01/maxKvCounts/"
	year := now.Year()
	month := int(now.Month())
	return fmt.Sprintf(BillingMonthStorageFormat, localPrefix, year, month, billingMetric)
}

func GetMonthlyBillingPath(localPrefix string, now time.Time) string {
	return fmt.Sprintf(BillingMonthStorageFormat, localPrefix, now.Year(), int(now.Month()), "")
}

func GetAttributionMaxPath(localPathPrefix string, month time.Time, attributionMetricName string) string {
	return GetMonthlyBillingMetricPath(localPathPrefix, month, AttributionMaxPrefix+attributionMetricName)
}

// SecretEngineCounts holds in-memory billing counters for all secret engine metric types.
type SecretEngineCounts struct {
	// Integer-counted data-protection engines
	Transit   DataProtectionEngineCounts
	Transform DataProtectionEngineCounts
	GcpKms    DataProtectionEngineCounts

	// Float-counted credential/certificate metrics
	Oidc       CredentialUnits
	Spiffe     CredentialUnits
	ExternalCa CredentialUnits
}

// DataProtectionEngineCounts holds a uint64 monthly counter and mount attribution
// for integer-counted data-protection engines (Transit, Transform, GcpKms).
type DataProtectionEngineCounts struct {
	// MonthlyCount is the monthly consumption billing count for the engine.
	MonthlyCount *atomic.Uint64 `json:"monthlyCount,omitempty"`
	// AttributionTracker holds the in-memory mount/namespace attribution for this engine.
	AttributionTracker
}

// CredentialUnits holds a float64 units counter and mount attribution for
// float-counted metrics (Oidc, Spiffe, ExternalCa).
type CredentialUnits struct {
	// MonthlyUnits is the monthly consumption billing units for the engine.
	MonthlyUnits *uberatomic.Float64
	// AttributionTracker holds the in-memory mount/namespace attribution.
	AttributionTracker
}

// AttributionTracker holds an in-memory map of mount attribution entries and the
// lock that protects it. Embedding this struct gives any metric type a consistent
// pair of fields (MountAttribution / MountAttributionLock) and a shared
// AccumulateMountAttributions method, so attribution tracking can be added to new
// metric types without duplicating fields or logic.
type AttributionTracker struct {
	// MountAttribution contains mount/namespace breakdown keyed by mount accessor.
	MountAttribution map[string]logical.MountAttribution
	// MountAttributionLock protects access to MountAttribution.
	MountAttributionLock sync.RWMutex
}

var _ logical.ConsumptionBillingManager = (*ConsumptionBilling)(nil)

// AccumulateMountAttributions extracts mount/namespace fields from a WriteBillingData
// payload and upserts them into the MountAttribution map, incrementing the count
// for an already-seen accessor or inserting a new entry.
// The mountPath is stored without its namespace prefix so attribution paths are
// consistent regardless of namespace depth.
func (d *AttributionTracker) AccumulateMountAttributions(ctx context.Context, data map[string]interface{}, count float64, getParentNamespaceID func(string) string) error {
	// Resolve namespace info from context — plugins do not pass namespace fields.
	namespaceID := ""
	namespacePath := ""
	if ns, err := namespace.FromContext(ctx); err == nil && ns != nil {
		namespaceID = ns.ID
		namespacePath = ns.Path
	}
	parentNamespaceID := ""
	if getParentNamespaceID != nil {
		parentNamespaceID = getParentNamespaceID(namespacePath)
	}

	// Extract the rest of the info from the plugin data.
	mountPath, ok := data["mountPath"].(string)
	if !ok {
		return fmt.Errorf("invalid value type for mountPath")
	}
	// Strip the namespace prefix so we store only the mount-local path.
	// e.g. "ns1/transit/" with namespacePath "ns1/" becomes "transit/".
	if namespacePath != "" {
		mountPath = strings.TrimPrefix(mountPath, namespacePath)
	}
	mountAccessor, ok := data["mountAccessor"].(string)
	if !ok {
		return fmt.Errorf("invalid value type for mountAccessor")
	}
	if mountAccessor == "" {
		return fmt.Errorf("mountAccessor cannot be empty")
	}
	mountType, ok := data["mountType"].(string)
	if !ok {
		return fmt.Errorf("invalid value type for mountType")
	}
	backendAwareUUID, ok := data["backendAwareUUID"].(string)
	if !ok {
		return fmt.Errorf("invalid value type for backendAwareUUID")
	}
	mountRunningVersion, ok := data["mountRunningVersion"].(string)
	if !ok {
		return fmt.Errorf("invalid value type for mountRunningVersion")
	}
	isExternal, _ := data["isExternal"].(bool)
	d.MountAttributionLock.Lock()
	var prev float64
	if existing, exists := d.MountAttribution[mountAccessor]; exists {
		if f, ok2 := existing.Count.(float64); ok2 {
			prev = f
		}
	}
	// Always write the full entry from the current request so that any metadata
	// change (e.g. namespace move, plugin upgrade) is reflected immediately.
	// Only the accumulated count is carried over from the previous entry.
	d.MountAttribution[mountAccessor] = logical.MountAttribution{
		MountPath:           mountPath,
		MountAccessor:       mountAccessor,
		MountType:           mountType,
		NamespaceID:         namespaceID,
		NamespacePath:       namespacePath,
		ParentNamespaceID:   parentNamespaceID,
		BackendAwareUUID:    backendAwareUUID,
		Count:               prev + count,
		MountRunningVersion: mountRunningVersion,
		IsExternal:          isExternal,
	}
	d.MountAttributionLock.Unlock()

	return nil
}

func (s *ConsumptionBilling) WriteBillingData(ctx context.Context, mountType string, data map[string]interface{}) error {
	if s == nil {
		return nil
	}

	switch mountType {
	case MountTypeTransit:
		val, ok := data["count"].(uint64)
		if !ok {
			err := fmt.Errorf("invalid value type for transit")
			return err
		}

		s.SecretEngineCounts.Transit.MonthlyCount.Add(val)
		if err := s.SecretEngineCounts.Transit.AccumulateMountAttributions(ctx, data, float64(val), s.GetParentNamespaceID); err != nil {
			return err
		}
	case MountTypeTransform:
		val, ok := data["count"].(uint64)
		if !ok {
			err := fmt.Errorf("invalid value type for transform")
			return err
		}

		s.SecretEngineCounts.Transform.MonthlyCount.Add(val)
		if err := s.SecretEngineCounts.Transform.AccumulateMountAttributions(ctx, data, float64(val), s.GetParentNamespaceID); err != nil {
			return err
		}
	case MountTypeSpiffe:
		// SPIFFE JWT uses float64 for duration-adjusted units
		val, ok := data["units"].(float64)
		if !ok {
			err := fmt.Errorf("invalid value type for spiffe")
			return err
		}

		s.SecretEngineCounts.Spiffe.MonthlyUnits.Add(val)
		if err := s.SecretEngineCounts.Spiffe.AccumulateMountAttributions(ctx, data, float64(val), s.GetParentNamespaceID); err != nil {
			return err
		}
	case MountTypeGcpKms:
		val, ok := data["count"].(uint64)
		if !ok {
			err := fmt.Errorf("invalid value type for gcp kms")
			return err
		}

		s.SecretEngineCounts.GcpKms.MonthlyCount.Add(val)
		if err := s.SecretEngineCounts.GcpKms.AccumulateMountAttributions(ctx, data, float64(val), s.GetParentNamespaceID); err != nil {
			return err
		}
	case MountTypeExCa:
		// External CA uses float64 for duration-adjusted units
		val, ok := data["units"].(float64)
		if !ok {
			err := fmt.Errorf("invalid value type for external-ca")
			return err
		}

		s.SecretEngineCounts.ExternalCa.MonthlyUnits.Add(val)
		if err := s.SecretEngineCounts.ExternalCa.AccumulateMountAttributions(ctx, data, val, s.GetParentNamespaceID); err != nil {
			return err
		}
	default:
		err := fmt.Errorf("unknown metric type: %s", mountType)
		return err
	}
	return nil
}
