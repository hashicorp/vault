package dbplugin_test

import (
	"errors"
	stdhttp "net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/plugins"
	"github.com/hashicorp/vault/vault"
	log "github.com/mgutz/logxi/v1"
)

type mockPlugin struct {
	users map[string][]string
}

func (m *mockPlugin) Type() (string, error) { return "mock", nil }
func (m *mockPlugin) CreateUser(statements dbplugin.Statements, usernameConf dbplugin.UsernameConfig, expiration time.Time) (username string, password string, err error) {
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
func (m *mockPlugin) RenewUser(statements dbplugin.Statements, username string, expiration time.Time) error {
	err := errors.New("err")
	if username == "" || expiration.IsZero() {
		return err
	}

	if _, ok := m.users[username]; !ok {
		return err
	}

	return nil
}
func (m *mockPlugin) RevokeUser(statements dbplugin.Statements, username string) error {
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
func (m *mockPlugin) Initialize(conf map[string]interface{}, _ bool) error {
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

func getCore(t *testing.T) ([]*vault.TestClusterCore, logical.SystemView) {
	coreConfig := &vault.CoreConfig{}

	handler1 := stdhttp.NewServeMux()
	handler2 := stdhttp.NewServeMux()
	handler3 := stdhttp.NewServeMux()

	// Chicken-and-egg: Handler needs a core. So we create handlers first, then
	// add routes chained to a Handler-created handler.
	cores := vault.TestCluster(t, []stdhttp.Handler{handler1, handler2, handler3}, coreConfig, false)
	handler1.Handle("/", http.Handler(cores[0].Core))
	handler2.Handle("/", http.Handler(cores[1].Core))
	handler3.Handle("/", http.Handler(cores[2].Core))

	core := cores[0]

	sys := vault.TestDynamicSystemView(core.Core)
	vault.TestAddTestPlugin(t, core.Core, "test-plugin", "TestPlugin_Main")

	return cores, sys
}

// This is not an actual test case, it's a helper function that will be executed
// by the go-plugin client via an exec call.
func TestPlugin_Main(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	plugin := &mockPlugin{
		users: make(map[string][]string),
	}

	args := []string{"--tls-skip-verify=true"}

	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(args)

	plugins.Serve(plugin, apiClientMeta.GetTLSConfig())
}

func TestPlugin_Initialize(t *testing.T) {
	cores, sys := getCore(t)
	for _, core := range cores {
		defer core.CloseListeners()
	}

	dbRaw, err := dbplugin.PluginFactory("test-plugin", sys, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connectionDetails := map[string]interface{}{
		"test": 1,
	}

	err = dbRaw.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = dbRaw.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPlugin_CreateUser(t *testing.T) {
	cores, sys := getCore(t)
	for _, core := range cores {
		defer core.CloseListeners()
	}

	db, err := dbplugin.PluginFactory("test-plugin", sys, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer db.Close()

	connectionDetails := map[string]interface{}{
		"test": 1,
	}

	err = db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConf := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	us, pw, err := db.CreateUser(dbplugin.Statements{}, usernameConf, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if us != "test" || pw != "test" {
		t.Fatal("expected username and password to be 'test'")
	}

	// try and save the same user again to verify it saved the first time, this
	// should return an error
	_, _, err = db.CreateUser(dbplugin.Statements{}, usernameConf, time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("expected an error, user wasn't created correctly")
	}
}

func TestPlugin_RenewUser(t *testing.T) {
	cores, sys := getCore(t)
	for _, core := range cores {
		defer core.CloseListeners()
	}

	db, err := dbplugin.PluginFactory("test-plugin", sys, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer db.Close()

	connectionDetails := map[string]interface{}{
		"test": 1,
	}
	err = db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConf := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	us, _, err := db.CreateUser(dbplugin.Statements{}, usernameConf, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = db.RenewUser(dbplugin.Statements{}, us, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPlugin_RevokeUser(t *testing.T) {
	cores, sys := getCore(t)
	for _, core := range cores {
		defer core.CloseListeners()
	}

	db, err := dbplugin.PluginFactory("test-plugin", sys, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer db.Close()

	connectionDetails := map[string]interface{}{
		"test": 1,
	}
	err = db.Initialize(connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	usernameConf := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	us, _, err := db.CreateUser(dbplugin.Statements{}, usernameConf, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test default revoke statememts
	err = db.RevokeUser(dbplugin.Statements{}, us)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Try adding the same username back so we can verify it was removed
	_, _, err = db.CreateUser(dbplugin.Statements{}, usernameConf, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
