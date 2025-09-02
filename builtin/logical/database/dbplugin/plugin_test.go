// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package dbplugin_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

type mockPlugin struct {
	users map[string][]string
}

var _ dbplugin.Database = &mockPlugin{}

func (m *mockPlugin) Type() (string, error) { return "mock", nil }
func (m *mockPlugin) CreateUser(_ context.Context, statements dbplugin.Statements, usernameConf dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
	err = errors.New("err")
	if usernameConf.DisplayName == "" || expiration.IsZero() {
		return "", "", err
	}

	if _, ok := m.users[usernameConf.DisplayName]; ok {
		return "", "", err
	}

	m.users[usernameConf.DisplayName] = []string{password}

	return usernameConf.DisplayName, "test", nil
}

func (m *mockPlugin) RenewUser(_ context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	err := errors.New("err")
	if username == "" || expiration.IsZero() {
		return err
	}

	if _, ok := m.users[username]; !ok {
		return err
	}

	return nil
}

func (m *mockPlugin) RevokeUser(_ context.Context, statements dbplugin.Statements, username string) error {
	err := errors.New("err")
	if username == "" {
		return err
	}

	if _, ok := m.users[username]; !ok {
		return err
	}

	delete(m.users, username)
	return nil
}

func (m *mockPlugin) RotateRootCredentials(_ context.Context, statements []string) (map[string]interface{}, error) {
	return nil, nil
}

func (m *mockPlugin) Init(_ context.Context, conf map[string]interface{}, _ bool) (map[string]interface{}, error) {
	err := errors.New("err")
	if len(conf) != 1 {
		return nil, err
	}

	return conf, nil
}

func (m *mockPlugin) Initialize(_ context.Context, conf map[string]interface{}, _ bool) error {
	err := errors.New("err")
	if len(conf) != 1 {
		return err
	}

	return nil
}

func (m *mockPlugin) Close() error {
	m.users = nil
	return nil
}

func (m *mockPlugin) GenerateCredentials(ctx context.Context) (password string, err error) {
	return password, err
}

func (m *mockPlugin) SetCredentials(ctx context.Context, statements dbplugin.Statements, staticConfig dbplugin.StaticUserConfig) (username string, password string, err error) {
	return username, password, err
}

func getCluster(t *testing.T) (*vault.TestCluster, logical.SystemView) {
	t.Helper()
	pluginDir := corehelpers.MakeTestPluginDir(t)
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		PluginDirectory: pluginDir,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	cores := cluster.Cores

	sys := vault.TestDynamicSystemView(cores[0].Core, nil)
	vault.TestAddTestPlugin(t, cores[0].Core, "test-plugin", consts.PluginTypeDatabase, "", "TestPlugin_GRPC_Main", []string{})

	return cluster, sys
}

// This is not an actual test case, it's a helper function that will be executed
// by the go-plugin client via an exec call.
func TestPlugin_GRPC_Main(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" && os.Getenv(pluginutil.PluginMetadataModeEnv) != "true" {
		return
	}

	plugin := &mockPlugin{
		users: make(map[string][]string),
	}

	args := []string{"--tls-skip-verify=true"}

	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)

	dbplugin.Serve(plugin, api.VaultPluginTLSProvider(apiClientMeta.GetTLSConfig()))
}

func TestPlugin_Init(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	dbRaw, err := dbplugin.PluginFactoryVersion(namespace.RootContext(nil), "test-plugin", "", sys, log.NewNullLogger())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connectionDetails := map[string]interface{}{
		"test": 1,
	}

	_, err = dbRaw.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = dbRaw.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPlugin_CreateUser(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	db, err := dbplugin.PluginFactoryVersion(namespace.RootContext(nil), "test-plugin", "", sys, log.NewNullLogger())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer db.Close()

	connectionDetails := map[string]interface{}{
		"test": 1,
	}

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConf := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	us, pw, err := db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConf, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if us != "test" || pw != "test" {
		t.Fatal("expected username and password to be 'test'")
	}

	// try and save the same user again to verify it saved the first time, this
	// should return an error
	_, _, err = db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConf, time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("expected an error, user wasn't created correctly")
	}
}

func TestPlugin_RenewUser(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	db, err := dbplugin.PluginFactoryVersion(namespace.RootContext(nil), "test-plugin", "", sys, log.NewNullLogger())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer db.Close()

	connectionDetails := map[string]interface{}{
		"test": 1,
	}
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConf := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	us, _, err := db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConf, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = db.RenewUser(context.Background(), dbplugin.Statements{}, us, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPlugin_RevokeUser(t *testing.T) {
	cluster, sys := getCluster(t)
	defer cluster.Cleanup()

	db, err := dbplugin.PluginFactoryVersion(namespace.RootContext(nil), "test-plugin", "", sys, log.NewNullLogger())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer db.Close()

	connectionDetails := map[string]interface{}{
		"test": 1,
	}
	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConf := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	us, _, err := db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConf, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test default revoke statements
	err = db.RevokeUser(context.Background(), dbplugin.Statements{}, us)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Try adding the same username back so we can verify it was removed
	_, _, err = db.CreateUser(context.Background(), dbplugin.Statements{}, usernameConf, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
