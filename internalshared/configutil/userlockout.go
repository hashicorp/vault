package configutil

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type UserLockoutConfig struct{
	Type string
	LockoutThreshold int64 `hcl:"-"`
	LockoutThresholdRaw interface{} `hcl:"lockout_threshold"`
	LockoutDuration time.Duration `hcl:"-"`
	LockoutDurationRaw interface{} `hcl:"lockout_duration"`
	LockoutCounterReset time.Duration `hcl:"-"`
	LockoutCounterResetRaw interface{} `hcl:"lockout_counter_reset"`
	DisableLockout bool `hcl:"-"`
	DisableLockoutRaw interface{} `hcl:"disable_lockout"`
}

func ParseUserLockouts(result *SharedConfig, list *ast.ObjectList) error {
	var err error
	result.UserLockoutConfigs = make([]*UserLockoutConfig, 0, len(list.Items))
	for i, item := range list.Items {
		var userLockoutConfig UserLockoutConfig
		if err := hcl.DecodeObject(&userLockoutConfig, item.Val); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("userLockouts.%d:", i))
		}
	
		// Base values
		{
			switch {
			case userLockoutConfig.Type != "":
			case len(item.Keys) == 1:
				userLockoutConfig.Type = strings.ToLower(item.Keys[0].Token.Value().(string))
			default:
				return multierror.Prefix(errors.New("auth type for user lockout must be specified, if it applies to all auth methods specify \"all\" "), fmt.Sprintf("user_lockouts.%d:", i))
			}

			userLockoutConfig.Type = strings.ToLower(userLockoutConfig.Type)
			switch userLockoutConfig.Type {
			case "all", "ldap", "approle", "userpass":
				result.found(userLockoutConfig.Type, userLockoutConfig.Type)
			default:
				return multierror.Prefix(fmt.Errorf("unsupported auth type %q", userLockoutConfig.Type), fmt.Sprintf("user_lockouts.%d:", i))
			}
		}

		// Lockout Parameters
		{
			if userLockoutConfig.LockoutThresholdRaw != nil {
				if userLockoutConfig.LockoutThreshold, err = parseutil.ParseInt(userLockoutConfig.LockoutThresholdRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing lockout_threshold: %w", err), fmt.Sprintf("user_lockouts.%d", i))
				}

				userLockoutConfig.LockoutThresholdRaw = nil
			}

			if userLockoutConfig.LockoutDurationRaw != nil {
				if userLockoutConfig.LockoutDuration, err = parseutil.ParseDurationSecond(userLockoutConfig.LockoutDurationRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing lockout_duration: %w", err), fmt.Sprintf("user_lockouts.%d", i))
				}
				if userLockoutConfig.LockoutDuration < 0 {
					return multierror.Prefix(errors.New("lockout_duration cannot be negative"), fmt.Sprintf("user_lockouts.%d", i))
				}

				userLockoutConfig.LockoutDurationRaw = nil
			}

			if userLockoutConfig.LockoutCounterResetRaw != nil {
				if userLockoutConfig.LockoutCounterReset, err = parseutil.ParseDurationSecond(userLockoutConfig.LockoutCounterResetRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing lockout_counter_reset: %w", err), fmt.Sprintf("user_lockouts.%d", i))
				}
				if userLockoutConfig.LockoutCounterReset < 0 {
					return multierror.Prefix(errors.New("lockout_counter_reset cannot be negative"), fmt.Sprintf("user_lockouts.%d", i))
				}

				userLockoutConfig.LockoutCounterResetRaw = nil
			}

			if userLockoutConfig.DisableLockoutRaw != nil {
				if userLockoutConfig.DisableLockout, err = parseutil.ParseBool(userLockoutConfig.DisableLockoutRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for disable_lockout: %w", err), fmt.Sprintf("user_lockouts.%d", i))
				}

				userLockoutConfig.DisableLockoutRaw = nil
			}

		}

		result.UserLockoutConfigs = append(result.UserLockoutConfigs, &userLockoutConfig)
	}

	return nil
}
