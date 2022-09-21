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

const (
	UserLockoutThresholdDefault    = 5
	UserLockoutDurationDefault     = 15 * time.Minute
	UserLockoutCounterResetDefault = 15 * time.Minute
	DisableUserLockoutDefault      = false
)

type UserLockoutConfig struct {
	Type                   string
	LockoutThreshold       int64         `hcl:"-"`
	LockoutThresholdRaw    interface{}   `hcl:"lockout_threshold"`
	LockoutDuration        time.Duration `hcl:"-"`
	LockoutDurationRaw     interface{}   `hcl:"lockout_duration"`
	LockoutCounterReset    time.Duration `hcl:"-"`
	LockoutCounterResetRaw interface{}   `hcl:"lockout_counter_reset"`
	DisableLockout         bool          `hcl:"-"`
	DisableLockoutRaw      interface{}   `hcl:"disable_lockout"`
}

func ParseUserLockouts(result *SharedConfig, list *ast.ObjectList) error {
	var err error
	result.UserLockoutConfigs = make([]*UserLockoutConfig, 0, len(list.Items))
	userLockoutConfigsMap := make(map[string]*UserLockoutConfig)
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
			}

			if userLockoutConfig.LockoutDurationRaw != nil {
				if userLockoutConfig.LockoutDuration, err = parseutil.ParseDurationSecond(userLockoutConfig.LockoutDurationRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing lockout_duration: %w", err), fmt.Sprintf("user_lockouts.%d", i))
				}
				if userLockoutConfig.LockoutDuration < 0 {
					return multierror.Prefix(errors.New("lockout_duration cannot be negative"), fmt.Sprintf("user_lockouts.%d", i))
				}

			}

			if userLockoutConfig.LockoutCounterResetRaw != nil {
				if userLockoutConfig.LockoutCounterReset, err = parseutil.ParseDurationSecond(userLockoutConfig.LockoutCounterResetRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("error parsing lockout_counter_reset: %w", err), fmt.Sprintf("user_lockouts.%d", i))
				}
				if userLockoutConfig.LockoutCounterReset < 0 {
					return multierror.Prefix(errors.New("lockout_counter_reset cannot be negative"), fmt.Sprintf("user_lockouts.%d", i))
				}

			}
			if userLockoutConfig.DisableLockoutRaw != nil {
				if userLockoutConfig.DisableLockout, err = parseutil.ParseBool(userLockoutConfig.DisableLockoutRaw); err != nil {
					return multierror.Prefix(fmt.Errorf("invalid value for disable_lockout: %w", err), fmt.Sprintf("user_lockouts.%d", i))
				}
			}
		}
		userLockoutConfigsMap[userLockoutConfig.Type] = &userLockoutConfig
	}
	// userLockoutConfigsMap = SetDefaultUserLockoutValuesInMap(userLockoutConfigsMap)
	userLockoutConfigsMap = SetMissingUserLockoutValuesInMap(userLockoutConfigsMap)
	for _, userLockoutValues := range userLockoutConfigsMap {
		result.UserLockoutConfigs = append(result.UserLockoutConfigs, userLockoutValues)
	}
	return nil
}

// setDefaultUserLockoutValuesInMap sets user lockout default values for key "all" for user lockout fields thats are not configured
func setDefaultUserLockoutValuesInMap(userLockoutConfigsMap map[string]*UserLockoutConfig) map[string]*UserLockoutConfig {
	if userLockoutAll, ok := userLockoutConfigsMap["all"]; !ok {
		var tmpUserLockoutConfig UserLockoutConfig
		tmpUserLockoutConfig.LockoutThreshold = UserLockoutThresholdDefault
		tmpUserLockoutConfig.LockoutDuration = UserLockoutDurationDefault
		tmpUserLockoutConfig.LockoutCounterReset = UserLockoutCounterResetDefault
		tmpUserLockoutConfig.DisableLockout = DisableUserLockoutDefault
		userLockoutConfigsMap["all"] = &tmpUserLockoutConfig

	} else {
		if userLockoutAll.LockoutThresholdRaw == nil {
			userLockoutAll.LockoutThreshold = UserLockoutThresholdDefault
		}
		if userLockoutAll.LockoutDurationRaw == nil {
			userLockoutAll.LockoutDuration = UserLockoutDurationDefault
		}
		if userLockoutAll.LockoutCounterResetRaw == nil {
			userLockoutAll.LockoutCounterReset = UserLockoutCounterResetDefault
		}
		if userLockoutAll.DisableLockoutRaw == nil {
			userLockoutAll.DisableLockout = DisableUserLockoutDefault
		}
		userLockoutConfigsMap[userLockoutAll.Type] = userLockoutAll
	}
	return userLockoutConfigsMap
}

// setDefaultUserLockoutValuesInMap sets missing user lockout fields for other auth types with default values (from key "all")
func SetMissingUserLockoutValuesInMap(userLockoutConfigsMap map[string]*UserLockoutConfig) map[string]*UserLockoutConfig {
	userLockoutConfigsMap = setDefaultUserLockoutValuesInMap(userLockoutConfigsMap)
	for _, userLockoutAuth := range userLockoutConfigsMap {
		// set missing values
		if userLockoutAuth.LockoutThresholdRaw == nil {
			userLockoutAuth.LockoutThreshold = userLockoutConfigsMap["all"].LockoutThreshold
		}
		if userLockoutAuth.LockoutDurationRaw == nil {
			userLockoutAuth.LockoutDuration = userLockoutConfigsMap["all"].LockoutDuration
		}
		if userLockoutAuth.LockoutCounterResetRaw == nil {
			userLockoutAuth.LockoutCounterReset = userLockoutConfigsMap["all"].LockoutCounterReset
		}
		if userLockoutAuth.DisableLockoutRaw == nil {
			userLockoutAuth.DisableLockout = userLockoutConfigsMap["all"].DisableLockout
		}

		// set nil values to Raw fields
		if userLockoutAuth.LockoutThresholdRaw != nil {
			userLockoutAuth.LockoutThresholdRaw = nil
		}
		if userLockoutAuth.LockoutDurationRaw != nil {
			userLockoutAuth.LockoutDurationRaw = nil
		}
		if userLockoutAuth.LockoutCounterResetRaw != nil {
			userLockoutAuth.LockoutCounterResetRaw = nil
		}
		if userLockoutAuth.DisableLockoutRaw != nil {
			userLockoutAuth.DisableLockoutRaw = nil
		}

		userLockoutConfigsMap[userLockoutAuth.Type] = userLockoutAuth
	}

	return userLockoutConfigsMap
}
