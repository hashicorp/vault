package healthcheck

import (
	"crypto/x509"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type CRLValidityPeriod struct {
	Enabled bool

	CRLExpiryPercentage      int
	DeltaCRLExpiryPercentage int

	UnsupportedVersion bool
	NoDeltas           bool

	CRLs      map[string]*x509.RevocationList
	DeltaCRLs map[string]*x509.RevocationList

	CRLConfig *PathFetch
}

func NewCRLValidityPeriodCheck() Check {
	return &CRLValidityPeriod{
		CRLs:      make(map[string]*x509.RevocationList),
		DeltaCRLs: make(map[string]*x509.RevocationList),
	}
}

func (h *CRLValidityPeriod) Name() string {
	return "crl_validity_period"
}

func (h *CRLValidityPeriod) IsEnabled() bool {
	return h.Enabled
}

func (h *CRLValidityPeriod) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"crl_expiry_pct_critical":       "95",
		"delta_crl_expiry_pct_critical": "95",
	}
}

func (h *CRLValidityPeriod) LoadConfig(config map[string]interface{}) error {
	value, err := parseutil.SafeParseIntRange(config["crl_expiry_pct_critical"], 1, 99)
	if err != nil {
		return fmt.Errorf("error parsing %v.crl_expiry_pct_critical=%v: %w", h.Name(), config["crl_expiry_pct_critical"], err)
	}
	h.CRLExpiryPercentage = int(value)

	value, err = parseutil.SafeParseIntRange(config["delta_crl_expiry_pct_critical"], 1, 99)
	if err != nil {
		return fmt.Errorf("error parsing %v.delta_crl_expiry_pct_critical=%v: %w", h.Name(), config["delta_crl_expiry_pct_critical"], err)
	}
	h.DeltaCRLExpiryPercentage = int(value)

	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *CRLValidityPeriod) FetchResources(e *Executor) error {
	exit, _, issuers, err := pkiFetchIssuers(e, func() {
		h.UnsupportedVersion = true
	})
	if exit {
		return err
	}

	for _, issuer := range issuers {
		exit, _, crl, err := pkiFetchIssuerCRL(e, issuer, false, func() {
			h.UnsupportedVersion = true
		})
		if exit {
			if err != nil {
				return err
			}
			continue
		}

		h.CRLs[issuer] = crl

		exit, _, delta, err := pkiFetchIssuerCRL(e, issuer, true, func() {
			h.NoDeltas = true
		})
		if exit {
			if err != nil {
				return err
			}
			continue
		}

		h.DeltaCRLs[issuer] = delta
	}

	// Check if the issuer is fetched yet.
	configRet, err := e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/config/crl")
	if err != nil {
		return err
	}

	h.CRLConfig = configRet

	return nil
}

func (h *CRLValidityPeriod) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/{{mount}}/issuers",
			Message:  "This health check requires Vault 1.11+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	now := time.Now()
	crlDisabled := false
	if h.CRLConfig != nil {
		if h.CRLConfig.IsSecretPermissionsError() {
			ret := Result{
				Status:   ResultInsufficientPermissions,
				Endpoint: "/{{mount}}/config/crl",
				Message:  "This prevents the health check from seeing if the CRL is disabled and dropping the severity of this check accordingly.",
			}

			if e.Client.Token() == "" {
				ret.Message = "No token available so unable read authenticated CRL configuration for this mount. " + ret.Message
			} else {
				ret.Message = "This token lacks so permission to read the CRL configuration for this mount. " + ret.Message
			}

			results = append(results, &ret)
		} else if h.CRLConfig.Secret != nil && h.CRLConfig.Secret.Data["disabled"] != nil {
			crlDisabled = h.CRLConfig.Secret.Data["disabled"].(bool)
		}
	}

	if h.NoDeltas && len(h.DeltaCRLs) == 0 {
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/{{mount}}/issuer/*/crl/delta",
			Message:  "This health check validates Delta CRLs on Vault 1.12+, but an earlier version of Vault was used. No results about delta CRL validity will be returned.",
		}
		results = append(results, &ret)
	}

	for name, crl := range h.CRLs {
		var ret Result
		ret.Status = ResultOK
		ret.Endpoint = "/{{mount}}/issuer/" + name + "/crl"
		ret.Message = fmt.Sprintf("CRL's validity (%v to %v) is OK.", crl.ThisUpdate.Format("2006-01-02"), crl.NextUpdate.Format("2006-01-02"))

		used := now.Sub(crl.ThisUpdate)
		total := crl.NextUpdate.Sub(crl.ThisUpdate)
		ratio := time.Duration((int64(total) * int64(h.CRLExpiryPercentage)) / int64(100))
		if used >= ratio {
			expWhen := crl.ThisUpdate.Add(ratio)
			ret.Status = ResultCritical
			ret.Message = fmt.Sprintf("CRL's validity is outside of suggested rotation window: CRL's next update is expected at %v, but expires within %v%% of validity window (starting on %v and ending on %v). It is suggested to rotate this CRL and start propagating it to hosts to avoid any issues caused by stale CRLs.", crl.NextUpdate.Format("2006-01-02"), h.CRLExpiryPercentage, crl.ThisUpdate.Format("2006-01-02"), expWhen.Format("2006-01-02"))

			if crlDisabled == true {
				ret.Status = ResultInformational
				ret.Message += " Because the CRL is disabled, this is less of a concern."
			}
		}

		results = append(results, &ret)
	}

	for name, crl := range h.DeltaCRLs {
		var ret Result
		ret.Status = ResultOK
		ret.Endpoint = "/{{mount}}/issuer/" + name + "/crl/delta"
		ret.Message = fmt.Sprintf("Delta CRL's validity (%v to %v) is OK.", crl.ThisUpdate.Format("2006-01-02"), crl.NextUpdate.Format("2006-01-02"))

		used := now.Sub(crl.ThisUpdate)
		total := crl.NextUpdate.Sub(crl.ThisUpdate)
		ratio := time.Duration((int64(total) * int64(h.DeltaCRLExpiryPercentage)) / int64(100))
		if used >= ratio {
			expWhen := crl.ThisUpdate.Add(ratio)
			ret.Status = ResultCritical
			ret.Message = fmt.Sprintf("Delta CRL's validity is outside of suggested rotation window: Delta CRL's next update is expected at %v, but expires within %v%% of validity window (starting on %v and ending on %v). It is suggested to rotate this Delta CRL and start propagating it to hosts to avoid any issues caused by stale CRLs.", crl.NextUpdate.Format("2006-01-02"), h.CRLExpiryPercentage, crl.ThisUpdate.Format("2006-01-02"), expWhen.Format("2006-01-02"))

			if crlDisabled == true {
				ret.Status = ResultInformational
				ret.Message += " Because the CRL is disabled, this is less of a concern."
			}
		}

		results = append(results, &ret)
	}

	return
}
