// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/enthelpers/activationflags"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/stretchr/testify/require"
)

// TestActivationFlags_Read tests the read operation for the activation flags.
func TestActivationFlags_Read(t *testing.T) {
	t.Run("given an initial state then read flags and expect all to be unactivated", func(t *testing.T) {
		core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{})

		resp, err := core.systemBackend.HandleRequest(
			context.Background(),
			&logical.Request{
				Operation: logical.ReadOperation,
				Path:      prefixActivationFlags,
				Storage:   core.systemBarrierView,
			},
		)

		require.NoError(t, err)
		require.Equal(t, resp.Data, map[string]interface{}{
			"activated":   []string{},
			"unactivated": activationflags.NewFeatureActivationFlags().ValidActivationFeatures(),
		})
	})
}

// TestActivationFlags_BadFeatureName tests a nonexistent feature name or a missing feature name
// in the activation-flags path API call.
func TestActivationFlags_BadFeatureName(t *testing.T) {
	core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{})

	tests := map[string]struct {
		featureName string
	}{
		"if no feature name is provided then expect unsupported path": {
			featureName: "",
		},
		"if an invalid feature name is provided then expect unsupported path": {
			featureName: "fake-feature",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := core.router.Route(
				namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace),
				&logical.Request{
					Operation: logical.UpdateOperation,
					Path:      fmt.Sprintf("sys/%s/%s/%s", prefixActivationFlags, tt.featureName, verbActivationFlagsActivate),
					Storage:   core.systemBarrierView,
				},
			)

			require.Error(t, err)
			require.Nil(t, resp)
			require.Equal(t, err, logical.ErrUnsupportedPath)
		})
	}
}

// TestSystemBackend_activationFlagsPaths tests that the expected paths are returned based on the
// HVD-related environment variables set and the administrative namespace.
func TestSystemBackend_activationFlagsPaths(t *testing.T) {
	printPaths := func(paths []*framework.Path) string {
		list := ""
		for _, p := range paths {
			list += fmt.Sprintf("Pattern: %s\n", (*p).Pattern)
		}
		return list
	}

	core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{})

	tests := map[string]struct {
		settingAdminNSPath  bool
		enableHVDTierEnvVar bool
		enableOnHVDEnvVar   bool
		wantedPathCount     int
	}{
		"if registering paths to the root namespace then expect 2 paths": {
			settingAdminNSPath:  false,
			enableHVDTierEnvVar: false,
			enableOnHVDEnvVar:   false,
			wantedPathCount:     2,
		},
		"if registering paths to the root namespace with EnvHVDTierSecretsSync set then expect 2 paths": {
			settingAdminNSPath:  false,
			enableHVDTierEnvVar: true,
			enableOnHVDEnvVar:   false,
			wantedPathCount:     2,
		},
		"if registering paths to the root namespace on HVD then expect 2 paths": {
			settingAdminNSPath: false,
			enableOnHVDEnvVar:  true,
			wantedPathCount:    3,
		},
		"if registering paths to the admin namespace without EnvHVDTierSecretsSync set then expect 1 path": {
			settingAdminNSPath:  true,
			enableHVDTierEnvVar: false,
			enableOnHVDEnvVar:   false,
			wantedPathCount:     1,
		},
		"if registering paths to the admin namespace with EnvHVDTierSecretsSync set then expect 1 path": {
			settingAdminNSPath:  true,
			enableHVDTierEnvVar: true,
			enableOnHVDEnvVar:   false,
			wantedPathCount:     1,
		},
		"if registering paths to the admin namespace with EnvHVDFlag set then expect 1 path": {
			settingAdminNSPath: true,
			enableOnHVDEnvVar:  false,
			wantedPathCount:    1,
		},
		"if registering paths to the admin namespace with EnvHVDTierSecretsSync and on HVD then expect 2 paths": {
			settingAdminNSPath:  true,
			enableHVDTierEnvVar: true,
			enableOnHVDEnvVar:   true,
			wantedPathCount:     3,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.enableHVDTierEnvVar {
				t.Setenv(EnvHVDTierSecretsSync, "enable")
			}
			if tt.enableOnHVDEnvVar {
				t.Setenv(activationflags.EnvHVDFlag, "set")
			}

			got := core.systemBackend.activationFlagsPaths(tt.settingAdminNSPath)
			gotLength := len(got)
			if gotLength != tt.wantedPathCount {
				t.Errorf("Received %d paths back from SystemBackend.activationFlagsPaths(), wanted %d paths.\nReceived the following path patterns back:\n%s", gotLength, tt.wantedPathCount, printPaths(got))
			}
		})
	}
}

// TestActivationFlags_Write tests the write operations for the activation flags, for both
// activate and deactivate, whether we're running in Enterprise or HVD, and whether we're on
// a plus tier or not when in HVD.
func TestActivationFlags_Write(t *testing.T) {
	testActivationFlags_Write_Activate(t)
	testActivationFlags_Write_Deactivate(t)
}

// TestActivationFlags_Invalidate: Tests State Change Logic for Flags
//
// This logic determines changes to flags, but it does NOT account for flags that have been deleted.
// As of this writing, flag removal is not supported for activation flags.
//
// Valid State Transitions:
// 1. Unset (new flag) -> Active
// 2. Active -> Inactive
// 3. Inactive -> Active
//
// Behavior notes:
//   - If a flag does not exist in-memory (`!ok`), it is treated as a new flag.
//     Nodes should only react to the new flag if its state is being set to "Active".
//   - If a flag exists in-memory, any change in its value (e.g., Active -> Inactive) is considered valid
//     and is marked as a state change.
func TestActivationFlags_Invalidate(t *testing.T) {
	testCluster := func(t *testing.T) (*bytes.Buffer, *Core, *Core) {
		logOut := new(bytes.Buffer)
		logger := log.New(&log.LoggerOptions{
			Mutex:  &sync.Mutex{},
			Level:  log.Trace,
			Output: logOut,
		})
		inm, err := inmem.NewTransactionalInmem(nil, logger)
		require.NoError(t, err)
		inmha, err := inmem.NewInmemHA(nil, logger)
		require.NoError(t, err)
		coreConfig := &CoreConfig{
			Logger:                    logger,
			Physical:                  inm,
			HAPhysical:                inmha.(physical.HABackend),
			DisablePerformanceStandby: false,
		}
		cluster := NewTestCluster(t, coreConfig, &TestClusterOptions{NumCores: 2})
		cluster.Start()
		cores := cluster.Cores
		active := cores[0].Core
		perfStandby := cores[1].Core
		TestWaitActive(t, cores[0].Core)
		TestWaitPerfStandby(t, cores[1].Core)
		return logOut, active, perfStandby
	}

	moveToInactiveState := func(t *testing.T, active *Core, perfStandby *Core) {
		err := active.systemBackend.Core.FeatureActivationFlags.Write(context.Background(), "newflag", false)
		require.NoError(t, err)
		corehelpers.RetryUntil(t, 5*time.Second, func() error {
			enabled := perfStandby.systemBackend.Core.FeatureActivationFlags.IsActivationFlagEnabled("newflag")
			if enabled {
				return errors.New("flag should NOT be enabled")
			}
			return nil
		})
	}

	moveToActiveState := func(t *testing.T, active *Core, perfStandby *Core) {
		err := active.systemBackend.Core.FeatureActivationFlags.Write(context.Background(), "newflag", true)
		require.NoError(t, err)
		corehelpers.RetryUntil(t, 5*time.Second, func() error {
			enabled := perfStandby.systemBackend.Core.FeatureActivationFlags.IsActivationFlagEnabled("newflag")
			if !enabled {
				return errors.New("flag should be enabled")
			}
			return nil
		})
	}

	tsts := map[string]struct {
		fromState       string
		toState         string
		validTransition bool
	}{
		"Invalid transition to inactive": {
			"unset",
			"inactive",
			false,
		},
		"Valid transition to active": {
			"unset",
			"active",
			true,
		},
		"Valid transition from inactive to active": {
			"inactive",
			"active",
			true,
		},
		"Valid transition from active to inactive": {
			"active",
			"inactive",
			true,
		},
	}

	for tn, tst := range tsts {
		t.Run(tn, func(t *testing.T) {
			logOut, active, perfStandby := testCluster(t)

			switch tst.fromState {
			case "inactive":
				moveToInactiveState(t, active, perfStandby)
			case "active":
				moveToActiveState(t, active, perfStandby)
			default:
				// noop
			}

			switch tst.toState {
			case "inactive":
				moveToInactiveState(t, active, perfStandby)
			case "active":
				moveToActiveState(t, active, perfStandby)
			default:
				// noop
			}

			corehelpers.RetryUntil(t, 5*time.Second, func() error {
				if tst.validTransition && !strings.Contains(logOut.String(), "on update activation flag called for activation flag") {
					return errors.New("flag should have been updated")
				}

				if !tst.validTransition && strings.Contains(logOut.String(), "on update activation flag called for activation flag") {
					return errors.New("flag should have NOT been updated")
				}

				return nil
			})

			switch tst.toState {
			case "inactive":
				corehelpers.RetryUntil(t, 5*time.Second, func() error {
					enabled := perfStandby.systemBackend.Core.FeatureActivationFlags.IsActivationFlagEnabled("newflag")
					if enabled {
						return errors.New("flag should NOT be enabled")
					}
					return nil
				})
			case "active":
				corehelpers.RetryUntil(t, 5*time.Second, func() error {
					enabled := perfStandby.systemBackend.Core.FeatureActivationFlags.IsActivationFlagEnabled("newflag")
					if !enabled {
						return errors.New("flag should be enabled")
					}
					return nil
				})
			default:
				// noop
			}
		})
	}
}

type activationFlagsTestCase struct {
	name          string
	featureName   string
	namespaceName string
	isHVD         bool
	isPlusTier    bool
	wantResp      *logical.Response
	wantErr       error
}

func newActivateTestCases(featureName string) []activationFlagsTestCase {
	newActivateTestCase := func(featureName string, namespaceName string, isHVD, isPlusTier bool, wantErr error) activationFlagsTestCase {
		tc := activationFlagsTestCase{
			name:          fmt.Sprintf("namespace=%q,isHVD=%t,isPlusTier=%t, if %s is activated then expect it in the activated list and removed from the unactivated list", namespaceName, isHVD, isPlusTier, featureName),
			featureName:   featureName,
			namespaceName: namespaceName,
			isHVD:         isHVD,
			isPlusTier:    isPlusTier,
			wantResp:      nil,
			wantErr:       wantErr,
		}

		if wantErr == nil {
			tc.wantResp = &logical.Response{
				Data: map[string]interface{}{
					"activated": []string{featureName},
					"unactivated": slices.DeleteFunc(
						activationflags.NewFeatureActivationFlags().ValidActivationFeatures(),
						func(s string) bool {
							return s == featureName
						},
					),
				},
			}
		}

		return tc
	}

	tests := []activationFlagsTestCase{
		newActivateTestCase(featureName, "root", false, false, nil),
		newActivateTestCase(featureName, "root", false, true, nil),
		newActivateTestCase(featureName, "root", true, false, nil),
		newActivateTestCase(featureName, "root", true, true, nil),
		newActivateTestCase(featureName, "my-namespace", false, false, logical.ErrUnsupportedPath),
		newActivateTestCase(featureName, "my-namespace", false, true, logical.ErrUnsupportedPath),
		newActivateTestCase(featureName, "my-namespace", true, false, logical.ErrUnsupportedPath),
		newActivateTestCase(featureName, "my-namespace", true, true, nil),
	}

	return tests
}

func newDeactivateTestCase(featureName string) []activationFlagsTestCase {
	newDeactivateTestCase := func(featureName string, namespaceName string, isHVD, isPlusTier bool, wantErr error) activationFlagsTestCase {
		tc := activationFlagsTestCase{
			name:          fmt.Sprintf("namespace=%q,isHVD=%t,isPlusTier=%t, if %s is deactivated then expect an empty activated list and it in the unactivated list", namespaceName, isHVD, isPlusTier, featureName),
			featureName:   featureName,
			namespaceName: namespaceName,
			isHVD:         isHVD,
			isPlusTier:    isPlusTier,
			wantResp:      nil,
			wantErr:       wantErr,
		}

		if wantErr == nil {
			tc.wantResp = &logical.Response{
				Data: map[string]interface{}{
					"activated":   []string{},
					"unactivated": activationflags.NewFeatureActivationFlags().ValidActivationFeatures(),
				},
			}
		}

		return tc
	}

	tests := []activationFlagsTestCase{
		// We don't need to test cases where the activate step fails here.
		newDeactivateTestCase(featureName, "root", false, false, logical.ErrUnsupportedPath),
		newDeactivateTestCase(featureName, "root", false, true, logical.ErrUnsupportedPath),
		newDeactivateTestCase(featureName, "root", true, false, nil),
		newDeactivateTestCase(featureName, "root", true, true, nil),
		newDeactivateTestCase(featureName, "my-namespace", true, true, nil),
	}

	return tests
}

func newActivateTestCall(t *testing.T, tc activationFlagsTestCase) (*Core, *namespace.Namespace, *logical.Response, error) {
	coreConf := &CoreConfig{}

	if tc.isPlusTier {
		t.Setenv(EnvHVDTierSecretsSync, "enable")
	}

	if tc.isHVD {
		t.Setenv(activationflags.EnvHVDFlag, "set")
		if tc.namespaceName != "root" && tc.namespaceName != "" {
			coreConf = &CoreConfig{AdministrativeNamespacePath: tc.namespaceName}
		}
	}

	core, _, _ := TestCoreUnsealedWithConfig(t, coreConf)

	ns := namespace.RootNamespace
	if tc.namespaceName != "root" && tc.namespaceName != "" {
		ns = givenNamespace(namespace.RootContext(nil), t, core, tc.namespaceName)
	}

	resp, err := core.router.Route(
		namespace.ContextWithNamespace(context.Background(), ns),
		&logical.Request{
			Operation: logical.UpdateOperation,
			Path:      fmt.Sprintf("sys/%s/%s/%s", prefixActivationFlags, tc.featureName, verbActivationFlagsActivate),
			Storage:   core.systemBarrierView,
		},
	)

	return core, ns, resp, err
}

func newDeactivateTestCall(core *Core, ns *namespace.Namespace, tc activationFlagsTestCase) (*logical.Response, error) {
	return core.router.Route(
		namespace.ContextWithNamespace(context.Background(), ns),
		&logical.Request{
			Operation: logical.UpdateOperation,
			Path:      fmt.Sprintf("sys/%s/%s/%s", prefixActivationFlags, tc.featureName, verbActivationFlagsDeactivate),
			Storage:   core.systemBarrierView,
		},
	)
}

func testActivationFlags_Write_Activate(t *testing.T) {
	// Add new activation flag features to test here as they are added.
	tests := []activationFlagsTestCase{}
	tests = append(tests, newActivateTestCases(activationflags.SecretsSync)...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, resp, err := newActivateTestCall(t, tt)
			require.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				require.Equal(t, tt.wantResp.Data, resp.Data)
			} else {
				require.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func testActivationFlags_Write_Deactivate(t *testing.T) {
	// Add new activation flag features to test here as they are added.
	tests := []activationFlagsTestCase{}
	tests = append(tests, newDeactivateTestCase(activationflags.SecretsSync)...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, ns, _, _ := newActivateTestCall(t, tt)

			resp, err := newDeactivateTestCall(core, ns, tt)
			require.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				require.Equal(t, tt.wantResp.Data, resp.Data)
			} else {
				require.Equal(t, tt.wantResp, resp)
			}
		})
	}
}
