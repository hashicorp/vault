// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"errors"
	"fmt"
	"strconv"
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
	LockoutThreshold       uint64        `hcl:"-"`
	LockoutThresholdRaw    interface{}   `hcl:"lockout_threshold"`
	LockoutDuration        time.Duration `hcl:"-"`
	LockoutDurationRaw     interface{}   `hcl:"lockout_duration"`
	LockoutCounterReset    time.Duration `hcl:"-"`
	LockoutCounterResetRaw interface{}   `hcl:"lockout_counter_reset"`
	DisableLockout         bool          `hcl:"-"`
	DisableLockoutRaw      interface{}   `hcl:"disable_lockout"`
}

func GetSupportedUserLockoutsAuthMethods() []string {
	return []string{"userpass", "approle", "ldap"}
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
			// Supported auth methods for user lockout configuration: ldap, approle, userpass
			// "all" is used to apply the configuration to all supported auth methods
			switch userLockoutConfig.Type {
			case "all", "ldap", "approle", "userpass":
				result.found(userLockoutConfig.Type, userLockoutConfig.Type)
			default:
				return multierror.Prefix(fmt.Errorf("unsupported auth type %q", userLockoutConfig.Type), fmt.Sprintf("user_lockouts.%d:", i))
			}
		}

		// Lockout Parameters

		// Not setting raw entries to nil here as soon as they are parsed
		// as they are used to set the missing user lockout configuration values later.
		{
			if userLockoutConfig.LockoutThresholdRaw != nil {
				userLockoutThresholdString := fmt.Sprintf("%v", userLockoutConfig.LockoutThresholdRaw)
				if userLockoutConfig.LockoutThreshold, err = strconv.ParseUint(userLockoutThresholdString, 10, 64); err != nil {
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

	// Use raw entries to set values for user lockout configurations fields
	// that were not configured using config file.
	// The raw entries would mean that the entry was configured by the user using the config file.
	// If any of these fields are not configured using the config file (missing fields),
	// we set values for these fields with defaults
	// The issue with not being able to use non-raw entries is because of fields lockout threshold
	// and disable lockout. We cannot differentiate using non-raw entries if the user configured these fields
	// with values (0 and false) or if the user did not configure these values in config file at all.
	// The raw fields are set to nil after setting missing values in setNilValuesForRawUserLockoutFields function
	userLockoutsMap = setMissingUserLockoutValuesInMap(userLockoutsMap)
	for _, userLockoutValues := range userLockoutsMap {
		result.UserLockouts = append(result.UserLockouts, userLockoutValues)
	}
	return nil
}

// setUserLockoutValueAllInMap sets default user lockout values for key "all" (all auth methods)
// for user lockout fields that are not configured using config file
func setUserLockoutValueAllInMap(userLockoutAll *UserLockout) *UserLockout {
	if userLockoutAll.Type == "" {
		userLockoutAll.Type = "all"
	}
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
	return setNilValuesForRawUserLockoutFields(userLockoutAll)
}

// setMissingUserLockoutValuesInMap sets missing user lockout fields for auth methods
// with default values (from key "all") that are not configured using config file
func setMissingUserLockoutValuesInMap(userLockoutsMap map[string]*UserLockout) map[string]*UserLockout {
	// set values for "all" key with default values for "all" user lockout fields that are not configured
	// the "all" key values will be used as default values for other auth methods
	userLockoutAll, ok := userLockoutsMap["all"]
	switch ok {
	case true:
		userLockoutsMap["all"] = setUserLockoutValueAllInMap(userLockoutAll)
	default:
		userLockoutsMap["all"] = setUserLockoutValueAllInMap(&UserLockout{})
	}

	for _, userLockoutAuth := range userLockoutsMap {
		if userLockoutAuth.Type == "all" {
			continue
		}
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
		userLockoutAuth = setNilValuesForRawUserLockoutFields(userLockoutAuth)
		userLockoutsMap[userLockoutAuth.Type] = userLockoutAuth
	}
	return userLockoutsMap
}

// setNilValuesForRawUserLockoutFields sets nil values for user lockout Raw fields
func setNilValuesForRawUserLockoutFields(userLockout *UserLockout) *UserLockout {
	userLockout.LockoutThresholdRaw = nil
	userLockout.LockoutDurationRaw = nil
	userLockout.LockoutCounterResetRaw = nil
	userLockout.DisableLockoutRaw = nil
	return userLockout
}
