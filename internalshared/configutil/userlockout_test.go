package configutil

import (
	"reflect"
	"testing"
	"time"
)

func TestParseUserLockout(t *testing.T) {
	t.Parallel()
	t.Run("Missing user lockout block in config file", func(t *testing.T) {
		t.Parallel()
		inputConfig := make(map[string]*UserLockout)
		expectedConfig := make(map[string]*UserLockout)
		expectedConfigall := &UserLockout{}
		expectedConfigall.Type = "all"
		expectedConfigall.LockoutThreshold = UserLockoutThresholdDefault
		expectedConfigall.LockoutDuration = UserLockoutDurationDefault
		expectedConfigall.LockoutCounterReset = UserLockoutCounterResetDefault
		expectedConfigall.DisableLockout = DisableUserLockoutDefault
		expectedConfig["all"] = expectedConfigall

		outputConfig := setMissingUserLockoutValuesInMap(inputConfig)
		if !reflect.DeepEqual(expectedConfig["all"], outputConfig["all"]) {
			t.Errorf("user lockout config: expected %#v\nactual %#v", expectedConfig["all"], outputConfig["all"])
		}
	})
	t.Run("setting default lockout counter reset and lockout duration for userpass in config ", func(t *testing.T) {
		t.Parallel()
		// input user lockout in config file
		inputConfig := make(map[string]*UserLockout)
		configAll := &UserLockout{}
		configAll.Type = "all"
		configAll.LockoutCounterReset = 20 * time.Minute
		configAll.LockoutCounterResetRaw = "1200000000000"
		inputConfig["all"] = configAll
		configUserpass := &UserLockout{}
		configUserpass.Type = "userpass"
		configUserpass.LockoutDuration = 10 * time.Minute
		configUserpass.LockoutDurationRaw = "600000000000"
		inputConfig["userpass"] = configUserpass

		expectedConfig := make(map[string]*UserLockout)
		expectedConfigall := &UserLockout{}
		expectedConfigUserpass := &UserLockout{}
		// expected default values
		expectedConfigall.Type = "all"
		expectedConfigall.LockoutThreshold = UserLockoutThresholdDefault
		expectedConfigall.LockoutDuration = UserLockoutDurationDefault
		expectedConfigall.LockoutCounterReset = 20 * time.Minute
		expectedConfigall.DisableLockout = DisableUserLockoutDefault
		// expected values for userpass
		expectedConfigUserpass.Type = "userpass"
		expectedConfigUserpass.LockoutThreshold = UserLockoutThresholdDefault
		expectedConfigUserpass.LockoutDuration = 10 * time.Minute
		expectedConfigUserpass.LockoutCounterReset = 20 * time.Minute
		expectedConfigUserpass.DisableLockout = DisableUserLockoutDefault
		expectedConfig["all"] = expectedConfigall
		expectedConfig["userpass"] = expectedConfigUserpass

		outputConfig := setMissingUserLockoutValuesInMap(inputConfig)
		if !reflect.DeepEqual(expectedConfig["all"], outputConfig["all"]) {
			t.Errorf("user lockout config: expected %#v\nactual %#v", expectedConfig["all"], outputConfig["all"])
		}
		if !reflect.DeepEqual(expectedConfig["userpass"], outputConfig["userpass"]) {
			t.Errorf("user lockout config: expected %#v\nactual %#v", expectedConfig["userpass"], outputConfig["userpass"])
		}
	})
}
