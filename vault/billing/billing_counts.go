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
	// BillingRetentionMonths is the number of months of billing data to retain.
	// This includes the current month plus previous months (e.g., 37 = current + 36 previous months).
	BillingRetentionMonths = 37

	BillingSubPath                          = "billing/"
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

	BillingWriteInterval = 10 * time.Minute
	// pluginCountsSendTimeout is the timeout for sending plugin counts to the active node
	PluginCountsSendTimeout = 30 * time.Second
	// pluginCountsStandbyTime is how long to wait before sending plugin counts from a perf standby
	PluginCountsStandbyTime = 10 * time.Minute
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
	// For testing purposes. The cadence at which billing metrics are updated
	MetricsUpdateCadence time.Duration
	// For testing purposes. The cadence at which plugin counts are sent from perf standby to active
	PluginCountsSendCadence time.Duration
	// For testin purposes. TestOverrideClock holds a custom clock to modify time.Now, time.Ticker, time.Timer.
	// If nil, the default functions from the time package are used
	TestOverrideClock timeutil.Clock
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

type DataProtectionCallCounts struct {
	Transit   *atomic.Uint64 `json:"transit,omitempty"`
	Transform *atomic.Uint64 `json:"transform,omitempty"`
	GcpKms    *atomic.Uint64 `json:"gcpkms,omitempty"`
}

// IdentityTokenUnits tracks billing metrics for identity and authentication services
type IdentityTokenUnits struct {
	// OidcTokenDuration tracks the token duration units (seconds, not duration-adjusted) for billing purposes in memory.
	// This value is normalized before flushing to storage and is reset to 0 after flush in UpdateOidcDurationAdjustedCount.
	OidcTokenDuration *uberatomic.Float64 `json:"oidc,omitempty"`

	// SpiffeJwt stores duration-adjusted JWT token units as float64
	// We need to use the uberAtomic package to store atomic float64 values
	SpiffeJwt *uberatomic.Float64 `json:"spiffe_jwt,omitempty"`
}

var _ logical.ConsumptionBillingManager = (*ConsumptionBilling)(nil)

func (s *ConsumptionBilling) WriteBillingData(ctx context.Context, mountType string, data map[string]interface{}) error {
	if s == nil {
		return nil
	}

	switch mountType {
	case "transit":
		val, ok := data["count"].(uint64)
		if !ok {
			err := fmt.Errorf("invalid value type for transit")
			return err
		}

		s.DataProtectionCallCounts.Transit.Add(val)
	case "transform":
		val, ok := data["count"].(uint64)
		if !ok {
			err := fmt.Errorf("invalid value type for transform")
			return err
		}

		s.DataProtectionCallCounts.Transform.Add(val)
	case "spiffe":
		// SPIFFE JWT uses float64 for duration-adjusted units
		val, ok := data["units"].(float64)
		if !ok {
			err := fmt.Errorf("invalid value type for spiffe")
			return err
		}

		s.IdentityTokenUnits.SpiffeJwt.Add(val)
	case "gcpkms":
		val, ok := data["count"].(uint64)
		if !ok {
			err := fmt.Errorf("invalid value type for gcp kms")
			return err
		}

		s.DataProtectionCallCounts.GcpKms.Add(val)
	case "external-ca":
		// External CA uses float64 for duration-adjusted units
		val, ok := data["units"].(float64)
		if !ok {
			err := fmt.Errorf("invalid value type for external-ca")
			return err
		}

		s.ExternalCaCertUnits.Add(val)
	default:
		err := fmt.Errorf("unknown metric type: %s", mountType)
		return err
	}
	return nil
}
