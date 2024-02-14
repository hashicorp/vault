// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package healthcheck

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type EnableAutoTidy struct {
	Enabled            bool
	UnsupportedVersion bool

	IntervalDurationCritical time.Duration
	IntervalDurationWarning  time.Duration
	PauseDurationCritical    time.Duration
	PauseDurationWarning     time.Duration

	TidyConfig *PathFetch
}

func NewEnableAutoTidyCheck() Check {
	return &EnableAutoTidy{}
}

func (h *EnableAutoTidy) Name() string {
	return "enable_auto_tidy"
}

func (h *EnableAutoTidy) IsEnabled() bool {
	return h.Enabled
}

func (h *EnableAutoTidy) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"interval_duration_critical": "7d",
		"interval_duration_warning":  "2d",
		"pause_duration_critical":    "1s",
		"pause_duration_warning":     "200ms",
	}
}

func (h *EnableAutoTidy) fromConfig(config map[string]interface{}, param string) (time.Duration, error) {
	value, err := parseutil.ParseDurationSecond(config[param])
	if err != nil {
		return time.Duration(0), fmt.Errorf("failed to parse parameter %v.%v=%v: %w", h.Name(), param, config[param], err)
	}

	return value, nil
}

func (h *EnableAutoTidy) LoadConfig(config map[string]interface{}) error {
	var err error

	h.IntervalDurationCritical, err = h.fromConfig(config, "interval_duration_critical")
	if err != nil {
		return err
	}

	h.IntervalDurationWarning, err = h.fromConfig(config, "interval_duration_warning")
	if err != nil {
		return err
	}

	h.PauseDurationCritical, err = h.fromConfig(config, "pause_duration_critical")
	if err != nil {
		return err
	}

	h.PauseDurationWarning, err = h.fromConfig(config, "pause_duration_warning")
	if err != nil {
		return err
	}

	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *EnableAutoTidy) FetchResources(e *Executor) error {
	var err error
	h.TidyConfig, err = e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/config/auto-tidy")
	if err != nil {
		return err
	}

	if h.TidyConfig.IsUnsupportedPathError() {
		h.UnsupportedVersion = true
	}

	return nil
}

func (h *EnableAutoTidy) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/{{mount}}/config/auto-tidy",
			Message:  "This health check requires Vault 1.12+, but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	if h.TidyConfig == nil {
		return
	}

	if h.TidyConfig.IsSecretPermissionsError() {
		ret := Result{
			Status:   ResultInsufficientPermissions,
			Endpoint: "/{{mount}}/config/auto-tidy",
			Message:  "This prevents the health check from functioning at all, as it cannot .",
		}

		if e.Client.Token() == "" {
			ret.Message = "No token available so unable read authenticated auto-tidy configuration for this mount. " + ret.Message
		} else {
			ret.Message = "This token lacks permission to read the auto-tidy configuration for this mount. " + ret.Message
		}

		return []*Result{&ret}, nil
	}

	isEnabled := h.TidyConfig.Secret.Data["enabled"].(bool)
	intervalDuration, err := parseutil.ParseDurationSecond(h.TidyConfig.Secret.Data["interval_duration"])
	if err != nil {
		return nil, fmt.Errorf("error parsing API response from server for interval_duration: %w", err)
	}

	pauseDuration, err := parseutil.ParseDurationSecond(h.TidyConfig.Secret.Data["pause_duration"])
	if err != nil {
		return nil, fmt.Errorf("error parsing API response from server for pause_duration: %w", err)
	}

	if !isEnabled {
		ret := Result{
			Status:   ResultInformational,
			Endpoint: "/{{mount}}/config/auto-tidy",
			Message:  "Auto-tidy is currently disabled; consider enabling auto-tidy to execute tidy operations periodically. This helps the health and performance of a mount.",
		}
		results = append(results, &ret)
	} else {
		baseMsg := "Auto-tidy is configured with too long of a value for %v (%v); this could impact performance as tidies run too infrequently or take too long to execute."

		if intervalDuration >= h.IntervalDurationCritical {
			ret := Result{
				Status:   ResultCritical,
				Endpoint: "/{{mount}}/config/auto-tidy",
				Message:  fmt.Sprintf(baseMsg, "interval_duration", intervalDuration),
			}
			results = append(results, &ret)
		} else if intervalDuration >= h.IntervalDurationWarning {
			ret := Result{
				Status:   ResultWarning,
				Endpoint: "/{{mount}}/config/auto-tidy",
				Message:  fmt.Sprintf(baseMsg, "interval_duration", intervalDuration),
			}
			results = append(results, &ret)
		}

		if pauseDuration >= h.PauseDurationCritical {
			ret := Result{
				Status:   ResultCritical,
				Endpoint: "/{{mount}}/config/auto-tidy",
				Message:  fmt.Sprintf(baseMsg, "pause_duration", pauseDuration),
			}
			results = append(results, &ret)
		} else if pauseDuration >= h.PauseDurationWarning {
			ret := Result{
				Status:   ResultWarning,
				Endpoint: "/{{mount}}/config/auto-tidy",
				Message:  fmt.Sprintf(baseMsg, "pause_duration", pauseDuration),
			}
			results = append(results, &ret)
		}

		if len(results) == 0 {
			ret := Result{
				Status:   ResultOK,
				Endpoint: "/{{mount}}/config/auto-tidy",
				Message:  "Auto-tidy is enabled and configured appropriately.",
			}
			results = append(results, &ret)
		}
	}

	return
}
