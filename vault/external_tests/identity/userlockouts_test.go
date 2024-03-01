// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	UserLockoutThresholdDefault = 5
)

// TestIdentityStore_DisableUserLockoutTest tests that user login will
// fail when supplied with wrong credentials. If the user is locked,
// it returns permission denied. Otherwise, it returns invalid user
// credentials error if the user lockout feature is disabled.
// It tests disabling the feature using env variable VAULT_DISABLE_USER_LOCKOUT
// and also using auth tune. Also, tests that env var has more precedence over
// settings in auth tune.
func TestIdentityStore_DisableUserLockoutTest(t *testing.T) {
	// reset to false before exiting
	defer os.Unsetenv("VAULT_DISABLE_USER_LOCKOUT")

	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	// standby client
	client := cluster.Cores[1].Client

	// enable userpass
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	// create a userpass user
	_, err = client.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	// get mount accessor for userpass mount
	secret, err := client.Logical().Read("sys/auth/userpass")
	if err != nil || secret == nil {
		t.Fatal(err)
	}
	mountAccessor := secret.Data["accessor"].(string)

	// variables for auth tune
	disableLockout := true
	enableLockout := false

	tests := []struct {
		name                        string
		setDisableUserLockoutEnvVar string
		// default is false
		setDisableLockoutAuthTune bool
		expectedUserLocked        bool
	}{
		{
			name:                        "Both unset, uses default behaviour i.e; user lockout feature enabled",
			setDisableUserLockoutEnvVar: "",
			setDisableLockoutAuthTune:   false,
			expectedUserLocked:          true,
		},
		{
			name:                        "User lockout feature is disabled using auth tune",
			setDisableUserLockoutEnvVar: "",
			setDisableLockoutAuthTune:   true,
			expectedUserLocked:          false,
		},
		{
			name:                        "User Lockout feature is disabled using env var VAULT_DISABLE_USER_LOCKOUT",
			setDisableUserLockoutEnvVar: "true",
			setDisableLockoutAuthTune:   false,
			expectedUserLocked:          false,
		},
		{
			name:                        "User lockout feature is enabled using env variable, disabled using auth tune",
			setDisableUserLockoutEnvVar: "false",
			setDisableLockoutAuthTune:   true,
			expectedUserLocked:          true,
		},
		{
			name:                        "User lockout feature is disabled using auth tune and env variable",
			setDisableUserLockoutEnvVar: "true",
			setDisableLockoutAuthTune:   true,
			expectedUserLocked:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setDisableUserLockoutEnvVar != "" {
				os.Setenv("VAULT_DISABLE_USER_LOCKOUT", tt.setDisableUserLockoutEnvVar)
			} else {
				os.Unsetenv("VAULT_DISABLE_USER_LOCKOUT")
			}

			var disableLockoutAuthTune *bool

			// default for disable lockout is false
			disableLockoutAuthTune = &enableLockout

			if tt.setDisableLockoutAuthTune == true {
				disableLockoutAuthTune = &disableLockout
			}

			// tune auth mount
			userlockoutConfig := &api.UserLockoutConfigInput{
				DisableLockout: disableLockoutAuthTune,
			}
			err := client.Sys().TuneMount("auth/userpass", api.MountConfigInput{
				UserLockoutConfig: userlockoutConfig,
			})
			if err != nil {
				t.Fatal(err)
			}

			// login for default lockout threshold times with wrong credentials
			for i := 0; i < UserLockoutThresholdDefault; i++ {
				_, err = client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
					"password": "wrongPassword",
				})
				if err == nil {
					t.Fatal("expected login to fail due to wrong credentials")
				}
				if !strings.Contains(err.Error(), "invalid username or password") {
					t.Fatal(err)
				}
			}

			// login to check if user locked
			_, err = client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
				"password": "wrongPassword",
			})
			if err == nil {
				t.Fatal("expected login to fail due to wrong credentials")
			}

			switch tt.expectedUserLocked {
			case true:
				if !strings.Contains(err.Error(), logical.ErrPermissionDenied.Error()) {
					t.Fatalf("expected user to get locked but got %v", err)
				}
				// user locked, unlock user to perform next test iteration
				if _, err = client.Logical().Write("sys/locked-users/"+mountAccessor+"/unlock/bsmith", nil); err != nil {
					t.Fatal(err)
				}

			default:
				if !strings.Contains(err.Error(), "invalid username or password") {
					t.Fatalf("expected user to be unlocked but locked, got  %v", err)
				}
			}
		})
	}
}
