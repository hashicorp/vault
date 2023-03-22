// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

// TestIdentityStore_UserLockoutTest tests that the user gets locked after
// more than 1 failed login request than the number specified for
// lockout threshold field in user lockout configuration. It also
// tests that the user gets unlocked after the duration specified
// for lockout duration field has passed
func TestIdentityStore_UserLockoutTest(t *testing.T) {
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
	active := cluster.Cores[0].Client
	standby := cluster.Cores[1].Client

	err := active.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	// tune auth mount
	userlockoutConfig := &api.UserLockoutConfigInput{
		LockoutThreshold:            "3",
		LockoutDuration:             "5s",
		LockoutCounterResetDuration: "5s",
	}
	err = active.Sys().TuneMount("auth/userpass", api.MountConfigInput{
		UserLockoutConfig: userlockoutConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// create a user for userpass
	_, err = standby.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	// login failure count 1
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login failure count 2
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login failure count 3
	active.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login : permission denied as user locked out
	_, err = standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err == nil {
		t.Fatal("expected login to fail as user locked out")
	}
	if !strings.Contains(err.Error(), logical.ErrPermissionDenied.Error()) {
		t.Fatalf("expected to see permission denied error as user locked out, got %v", err)
	}

	time.Sleep(5 * time.Second)

	// login with right password and wait for user to get unlocked
	_, err = standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal("expected login to succeed as user is unlocked")
	}
}

// TestIdentityStore_UserFailedLoginMapResetOnSuccess tests that
// the user lockout feature is reset for a user after one successfull attempt
// after multiple failed login attempts (within lockout threshold)
func TestIdentityStore_UserFailedLoginMapResetOnSuccess(t *testing.T) {
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

	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	// tune auth mount
	userlockoutConfig := &api.UserLockoutConfigInput{
		LockoutThreshold:            "3",
		LockoutDuration:             "5s",
		LockoutCounterResetDuration: "5s",
	}
	err = client.Sys().TuneMount("auth/userpass", api.MountConfigInput{
		UserLockoutConfig: userlockoutConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// create a user for userpass
	_, err = client.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	// login failure count 1
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login failure count 2
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login with right credentials - successful login
	// entry for this user is removed from userFailedLoginInfo map
	_, err = client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	// login failure count 3, is now count 1 after successful login
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login failure count 4, is now count 2 after successful login
	// error should not be permission denied as user not locked out
	_, err = client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	if err == nil {
		t.Fatal("expected login to fail due to wrong credentials")
	}
	if !strings.Contains(err.Error(), "invalid username or password") {
		t.Fatalf("expected to see invalid username or password error as user is not locked out, got %v", err)
	}
}

// TestIdentityStore_DisableUserLockoutTest tests that user login will
// fail when supplied with wrong credentials. If the user is locked,
// it returns permission denied. In this case, it returns invalid user
// credentials error as the user lockout feature is disabled and the
// user did not get locked after multiple failed login attempts
func TestIdentityStore_DisableUserLockoutTest(t *testing.T) {
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

	active := cluster.Cores[0].Client
	standby := cluster.Cores[1].Client

	err := active.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	// tune auth mount
	disableLockout := true
	userlockoutConfig := &api.UserLockoutConfigInput{
		LockoutThreshold: "3",
		DisableLockout:   &disableLockout,
	}
	err = active.Sys().TuneMount("auth/userpass", api.MountConfigInput{
		UserLockoutConfig: userlockoutConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// create a userpass user
	_, err = standby.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	// login failure count 1
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login failure count 2
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login failure count 3
	active.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login failure count 4
	_, err = standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	if err == nil {
		t.Fatal("expected login to fail due to wrong credentials")
	}
	if !strings.Contains(err.Error(), "invalid username or password") {
		t.Fatalf("expected to see invalid username or password error as user is not locked out, got %v", err)
	}
}

// TestIdentityStore_LockoutCounterResetTest tests that the user lockout counter
// for a user is reset after no failed login attempts for a duration
// as specified for lockout counter reset field in user lockout configuration
func TestIdentityStore_LockoutCounterResetTest(t *testing.T) {
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
	active := cluster.Cores[0].Client
	standby := cluster.Cores[1].Client

	err := active.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	// tune auth mount
	userlockoutConfig := &api.UserLockoutConfigInput{
		LockoutThreshold:            "3",
		LockoutCounterResetDuration: "5s",
	}
	err = active.Sys().TuneMount("auth/userpass", api.MountConfigInput{
		UserLockoutConfig: userlockoutConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// create a user for userpass
	_, err = standby.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}

	// login failure count 1
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure count 2
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// set sleep timer to reset login counter
	time.Sleep(5 * time.Second)

	// login failure 3, count should be reset, this will be treated as failed count 1
	active.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 4, this will be treated as failed count 2
	_, err = standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	if err == nil {
		t.Fatal("expected login to fail due to wrong credentials")
	}
	if !strings.Contains(err.Error(), "invalid username or password") {
		t.Fatalf("expected to see invalid username or password error as user is not locked out, got %v", err)
	}
}

// TestIdentityStore_UnlockUserTest tests the user is
// unlocked if locked  using
// sys/locked-users/[mount_accessor]/unlock/[alias-identifier]
func TestIdentityStore_UnlockUserTest(t *testing.T) {
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
	active := cluster.Cores[0].Client
	standby := cluster.Cores[1].Client

	// enable userpass auth method on path userpass
	if err := active.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	}); err != nil {
		t.Fatal(err)
	}

	// get mount accessor for userpass mount
	secret, err := standby.Logical().Read("sys/auth/userpass")
	if err != nil || secret == nil {
		t.Fatal(err)
	}
	mountAccessor := secret.Data["accessor"].(string)

	// tune auth mount
	userlockoutConfig := &api.UserLockoutConfigInput{
		LockoutThreshold: "2",
		LockoutDuration:  "5m",
	}
	if err = active.Sys().TuneMount("auth/userpass", api.MountConfigInput{
		UserLockoutConfig: userlockoutConfig,
	}); err != nil {
		t.Fatal(err)
	}

	// create a user for userpass
	if _, err = standby.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	}); err != nil {
		t.Fatal(err)
	}

	// create another user for userpass with a different case
	if _, err = standby.Logical().Write("auth/userpass/users/bSmith", map[string]interface{}{
		"password": "training",
	}); err != nil {
		t.Fatal(err)
	}

	// login failure count 1
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure count 2
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login : permission denied as user locked out
	if _, err = standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "training",
	}); err == nil {
		t.Fatal("expected login to fail as user locked out")
	}
	if !strings.Contains(err.Error(), logical.ErrPermissionDenied.Error()) {
		t.Fatalf("expected to see permission denied error as user locked out, got %v", err)
	}

	// unlock user
	if _, err = standby.Logical().Write("sys/locked-users/"+mountAccessor+"/unlock/bsmith", nil); err != nil {
		t.Fatal(err)
	}

	// login: should be successful as user unlocked
	if _, err = standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "training",
	}); err != nil {
		t.Fatal("expected login to succeed as user is unlocked")
	}

	// login failure count 1 for user bSmith
	standby.Logical().Write("auth/userpass/login/bSmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure count 2 for user bSmith
	standby.Logical().Write("auth/userpass/login/bSmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login : permission denied as user locked out for user bSmith
	if _, err = standby.Logical().Write("auth/userpass/login/bSmith", map[string]interface{}{
		"password": "training",
	}); err == nil {
		t.Fatal("expected login to fail as user locked out")
	}
	if !strings.Contains(err.Error(), logical.ErrPermissionDenied.Error()) {
		t.Fatalf("expected to see permission denied error as user locked out, got %v", err)
	}

	// unlock user bSmith
	if _, err = standby.Logical().Write("sys/locked-users/"+mountAccessor+"/unlock/bSmith", nil); err != nil {
		t.Fatal(err)
	}

	// login: should be successful as user bSmith unlocked
	if _, err = standby.Logical().Write("auth/userpass/login/bSmith", map[string]interface{}{
		"password": "training",
	}); err != nil {
		t.Fatal("expected login to succeed as user is unlocked")
	}

	// unlock unlocked user
	if _, err = active.Logical().Write("sys/locked-users/mountAccessor/unlock/bsmith", nil); err != nil {
		t.Fatal(err)
	}
}
