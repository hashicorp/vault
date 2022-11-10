package healthcheck

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type TidyLastRun struct {
	Enabled            bool
	UnsupportedVersion bool

	LastRunCritical time.Duration
	LastRunWarning  time.Duration

	TidyStatus *PathFetch
}

func NewTidyLastRunCheck() Check {
	return &TidyLastRun{}
}

func (h *TidyLastRun) Name() string {
	return "tidy_last_run"
}

func (h *TidyLastRun) IsEnabled() bool {
	return h.Enabled
}

func (h *TidyLastRun) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"last_run_critical": "7d",
		"last_run_warning":  "2d",
	}
}

func (h *TidyLastRun) LoadConfig(config map[string]interface{}) error {
	var err error
	h.LastRunCritical, err = parseutil.ParseDurationSecond(config["last_run_critical"])
	if err != nil {
		return fmt.Errorf("failed to parse parameter %v.%v=%v: %w", h.Name(), "last_run_critical", config["last_run_critical"], err)
	}

	h.LastRunWarning, err = parseutil.ParseDurationSecond(config["last_run_warning"])
	if err != nil {
		return fmt.Errorf("failed to parse parameter %v.%v=%v: %w", h.Name(), "last_run_warning", config["last_run_warning"], err)
	}

	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *TidyLastRun) FetchResources(e *Executor) error {
	var err error

	h.TidyStatus, err = e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/tidy-status")
	if err != nil {
		return fmt.Errorf("failed to fetch mount's tidy-status value: %v", err)
	}

	if h.TidyStatus.IsUnsupportedPathError() {
		h.UnsupportedVersion = true
	}

	return nil
}

func (h *TidyLastRun) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		// Shouldn't happen; roles have been around forever.
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/{{mount}}/tidy-status",
			Message:  "This health check requires Vault 1.10+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	baseMsg := "Tidy hasn't run in the last %v; this can point to problems with the mount's auto-tidy configuration or an external tidy executor; this can impact PKI's and Vault's performance if not run regularly."

	ret := Result{
		Status:   ResultOK,
		Endpoint: "/{{mount}}/tidy-status",
		Message:  "Tidy has run recently on this mount.",
	}

	if h.TidyStatus.Secret != nil && h.TidyStatus.Secret.Data != nil {
		when := h.TidyStatus.Secret.Data["time_finished"]
		if when == nil {
			ret.Status = ResultCritical
			ret.Message = "Tidy hasn't run since this mount was created; this can point to problems with the mount's auto-tidy configuration or an external tidy executor; this can impact PKI's and Vault's performance if not run regularly. It is suggested to enable auto-tidy on this mount."
		} else {
			now := time.Now()
			lastRunCritical := now.Add(-1 * h.LastRunCritical)
			lastRunWarning := now.Add(-1 * h.LastRunWarning)

			whenT := when.(*time.Time)

			if whenT.Before(lastRunCritical) {
				ret.Status = ResultCritical
				ret.Message = fmt.Sprintf(baseMsg, h.LastRunCritical)
			} else if whenT.Before(lastRunWarning) {
				ret.Status = ResultWarning
				ret.Message = fmt.Sprintf(baseMsg, h.LastRunWarning)
			}
		}
	}

	results = append(results, &ret)

	return
}
