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

type UserLockout struct {
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
	result.UserLockouts = make([]*UserLockout, 0, len(list.Items))
	userLockoutsMap := make(map[string]*UserLockout)
	for i, item := range list.Items {
		var userLockoutConfig UserLockout
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
		userLockoutsMap[userLockoutConfig.Type] = &userLockoutConfig
	}

	userLockoutsMap = setMissingUserLockoutValuesInMap(userLockoutsMap)
	for _, userLockoutValues := range userLockoutsMap {
		result.UserLockouts = append(result.UserLockouts, userLockoutValues)
	}
	return nil
}

// setDefaultUserLockoutValuesInMap sets default user lockout values for key "all" (all auth methods) for user lockout fields that are not configured using config file
func setDefaultUserLockoutValuesInMap(userLockoutsMap map[string]*UserLockout) map[string]*UserLockout {
	if userLockoutAll, ok := userLockoutsMap["all"]; !ok {
		var tmpUserLockoutConfig UserLockout
		tmpUserLockoutConfig.Type = "all"
		tmpUserLockoutConfig.LockoutThreshold = UserLockoutThresholdDefault
		tmpUserLockoutConfig.LockoutDuration = UserLockoutDurationDefault
		tmpUserLockoutConfig.LockoutCounterReset = UserLockoutCounterResetDefault
		tmpUserLockoutConfig.DisableLockout = DisableUserLockoutDefault
		userLockoutsMap["all"] = &tmpUserLockoutConfig

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
		userLockoutsMap[userLockoutAll.Type] = userLockoutAll
	}
	return userLockoutsMap
}

// setDefaultUserLockoutValuesInMap sets missing user lockout fields for auth methods with default values (from key "all") that are not configured using config file
func setMissingUserLockoutValuesInMap(userLockoutsMap map[string]*UserLockout) map[string]*UserLockout {
	userLockoutsMap = setDefaultUserLockoutValuesInMap(userLockoutsMap)
	for _, userLockoutAuth := range userLockoutsMap {
		// set missing values
		if userLockoutAuth.LockoutThresholdRaw == nil {
			userLockoutAuth.LockoutThreshold = userLockoutsMap["all"].LockoutThreshold
		}
		if userLockoutAuth.LockoutDurationRaw == nil {
			userLockoutAuth.LockoutDuration = userLockoutsMap["all"].LockoutDuration
		}
		if userLockoutAuth.LockoutCounterResetRaw == nil {
			userLockoutAuth.LockoutCounterReset = userLockoutsMap["all"].LockoutCounterReset
		}
		if userLockoutAuth.DisableLockoutRaw == nil {
			userLockoutAuth.DisableLockout = userLockoutsMap["all"].DisableLockout
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

		userLockoutsMap[userLockoutAuth.Type] = userLockoutAuth
	}

	return userLockoutsMap
}
