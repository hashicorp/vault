package healthcheck

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type CAValidityPeriod struct {
	Enabled bool

	RootExpiries          map[ResultStatus]time.Duration
	IntermediateExpieries map[ResultStatus]time.Duration

	UnsupportedVersion bool

	Issuers map[string]*x509.Certificate
}

func NewCAValidityPeriodCheck() Check {
	return &CAValidityPeriod{
		RootExpiries:          make(map[ResultStatus]time.Duration, 3),
		IntermediateExpieries: make(map[ResultStatus]time.Duration, 3),
		Issuers:               make(map[string]*x509.Certificate),
	}
}

func (h *CAValidityPeriod) Name() string {
	return "ca_validity_period"
}

func (h *CAValidityPeriod) IsEnabled() bool {
	return h.Enabled
}

func (h *CAValidityPeriod) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"root_expiry_critical":              "180d",
		"intermediate_expiry_critical":      "30d",
		"root_expiry_warning":               "365d",
		"intermediate_expiry_warning":       "60d",
		"root_expiry_informational":         "730d",
		"intermediate_expiry_informational": "180d",
	}
}

func (h *CAValidityPeriod) LoadConfig(config map[string]interface{}) error {
	parameters := []string{
		"root_expiry_critical",
		"intermediate_expiry_critical",
		"root_expiry_warning",
		"intermediate_expiry_warning",
		"root_expiry_informational",
		"intermediate_expiry_informational",
	}
	for _, parameter := range parameters {
		name_split := strings.Split(parameter, "_")
		if len(name_split) != 3 || name_split[1] != "expiry" {
			return fmt.Errorf("bad parameter: %v / %v / %v", parameter, len(name_split), name_split[1])
		}

		status, present := NameResultStatusMap[name_split[2]]
		if !present {
			return fmt.Errorf("bad parameter: %v's type %v isn't in name map", parameter, name_split[2])
		}

		value_raw, present := config[parameter]
		if !present {
			return fmt.Errorf("parameter not present in config; Executor should've handled this for us: %v", parameter)
		}

		value, err := parseutil.ParseDurationSecond(value_raw)
		if err != nil {
			return fmt.Errorf("failed to parse parameter (%v=%v): %w", parameter, value_raw, err)
		}

		if name_split[0] == "root" {
			h.RootExpiries[status] = value
		} else if name_split[0] == "intermediate" {
			h.IntermediateExpieries[status] = value
		} else {
			return fmt.Errorf("bad parameter: %v's CA type isn't root/intermediate: %v", parameters, name_split[0])
		}
	}

	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *CAValidityPeriod) FetchResources(e *Executor) error {
	exit, _, issuers, err := pkiFetchIssuers(e, func() {
		h.UnsupportedVersion = true
	})
	if exit {
		return err
	}

	for _, issuer := range issuers {
		skip, _, cert, err := pkiFetchIssuer(e, issuer, func() {
			h.UnsupportedVersion = true
		})
		if skip {
			if err != nil {
				return err
			}
			continue
		}

		h.Issuers[issuer] = cert
	}

	return nil
}

func (h *CAValidityPeriod) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/{{mount}}/issuers",
			Message:  "This health check requires Vault 1.11+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	now := time.Now()

	for name, cert := range h.Issuers {
		var ret Result
		ret.Status = ResultOK
		ret.Endpoint = "/{{mount}}/issuer/" + name
		ret.Message = fmt.Sprintf("Issuer's validity (%v) is OK", cert.NotAfter.Format("2006-01-02"))

		hasSelfReference := bytes.Equal(cert.RawSubject, cert.RawIssuer)
		isSelfSigned := cert.CheckSignatureFrom(cert) == nil
		isRoot := hasSelfReference && isSelfSigned

		for _, criticality := range []ResultStatus{ResultCritical, ResultWarning, ResultInformational} {
			var d time.Duration
			if isRoot {
				d = h.RootExpiries[criticality]
			} else {
				d = h.IntermediateExpieries[criticality]
			}

			windowExpiry := now.Add(d)
			if cert.NotAfter.Before(windowExpiry) {
				ret.Status = criticality
				ret.Message = fmt.Sprintf("Issuer's validity is outside of the suggested rotation window: issuer is valid until %v but expires within %v (ending on %v). It is suggested to start rotating this issuer to new key material to avoid future downtime caused by this current issuer expiring.", cert.NotAfter.Format("2006-01-02"), FormatDuration(d), windowExpiry.Format("2006-01-02"))
				break
			}
		}

		results = append(results, &ret)
	}

	return
}
