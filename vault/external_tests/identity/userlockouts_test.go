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

	// read auth tune
	resp, err := active.Logical().Read("sys/auth/userpass/tune")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for reading auth tune")
	}
	_, err = standby.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}
	// login failure 1
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 2
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 3
	active.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 4
	// login should not fail as user lockout feature is disabled for this mount
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
		LockoutDuration:             "10s",
		LockoutCounterResetDuration: "10s",
	}
	err = client.Sys().TuneMount("auth/userpass", api.MountConfigInput{
		UserLockoutConfig: userlockoutConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// read auth tune
	resp, err := client.Logical().Read("sys/auth/userpass/tune")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for reading auth tune")
	}
	_, err = client.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}
	// login failure 1
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 2
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login with right credentials - successful login
	// userFailedLogin map that contains failed user login info reset
	_, err = client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}
	// login failure 1
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login failure 2 after successful login
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

func TestIdentityStore_LockoutDurationTest(t *testing.T) {
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

	// read auth tune
	resp, err := active.Logical().Read("sys/auth/userpass/tune")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for reading auth tune")
	}

	_, err = standby.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}
	// login failure 1
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 2
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 3
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
		t.Fatal(err.Error())
		t.Fatal("expected login to succeed as user is unlocked")
	}
}

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

	// read auth tune
	resp, err := active.Logical().Read("sys/auth/userpass/tune")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for reading auth tune")
	}
	_, err = standby.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}
	// login failure 1
	standby.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 2
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
