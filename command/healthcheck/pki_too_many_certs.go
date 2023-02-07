package healthcheck

import (
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type TooManyCerts struct {
	Enabled            bool
	UnsupportedVersion bool

	CountCritical int
	CountWarning  int

	CertCounts int
}

func NewTooManyCertsCheck() Check {
	return &TooManyCerts{}
}

func (h *TooManyCerts) Name() string {
	return "too_many_certs"
}

func (h *TooManyCerts) IsEnabled() bool {
	return h.Enabled
}

func (h *TooManyCerts) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"count_critical": 250000,
		"count_warning":  50000,
	}
}

func (h *TooManyCerts) LoadConfig(config map[string]interface{}) error {
	value, err := parseutil.SafeParseIntRange(config["count_critical"], 1, 15000000)
	if err != nil {
		return fmt.Errorf("error parsing %v.count_critical: %w", h.Name(), err)
	}
	h.CountCritical = int(value)

	value, err = parseutil.SafeParseIntRange(config["count_warning"], 1, 15000000)
	if err != nil {
		return fmt.Errorf("error parsing %v.count_warning: %w", h.Name(), err)
	}
	h.CountWarning = int(value)

	h.Enabled, err = parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}

	return nil
}

func (h *TooManyCerts) FetchResources(e *Executor) error {
	exit, leavesRet, _, err := pkiFetchLeavesList(e, func() {
		h.UnsupportedVersion = true
	})
	if exit {
		return err
	}

	h.CertCounts = leavesRet.ParsedCache["count"].(int)

	return nil
}

func (h *TooManyCerts) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		// Shouldn't happen; /certs has been around forever.
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/{{mount}}/certs",
			Message:  "This health check requires Vault 1.11+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	ret := Result{
		Status:   ResultOK,
		Endpoint: "/{{mount}}/certs",
		Message:  "This mount has an OK number of stored certificates.",
	}

	baseMsg := "This PKI mount has %v outstanding stored certificates; consider using no_store=false on roles, running tidy operations periodically, and using shorter certificate lifetimes to reduce the storage pressure on this mount."
	if h.CertCounts >= h.CountCritical {
		ret.Status = ResultCritical
		ret.Message = fmt.Sprintf(baseMsg, h.CertCounts)
	} else if h.CertCounts >= h.CountWarning {
		ret.Status = ResultWarning
		ret.Message = fmt.Sprintf(baseMsg, h.CertCounts)
	}

	results = append(results, &ret)

	return
}
