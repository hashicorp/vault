package dbplugin_test

import (
	"errors"
	"net"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	log "github.com/mgutz/logxi/v1"
)

type mockPlugin struct {
	users map[string][]string
}

func (m *mockPlugin) Type() (string, error) { return "mock", nil }
func (m *mockPlugin) CreateUser(statements dbplugin.Statements, usernamePrefix string, expiration time.Time) (username string, password string, err error) {
	err = errors.New("err")
	if usernamePrefix == "" || expiration.IsZero() {
		return "", "", err
	}

	if _, ok := m.users[usernamePrefix]; ok {
		return "", "", err
	}

	m.users[usernamePrefix] = []string{password}

	return usernamePrefix, "test", nil
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

func getCore(t *testing.T) (*vault.Core, net.Listener, logical.SystemView) {
	core, _, _, ln := vault.TestCoreUnsealedWithListener(t)
	http.TestServerWithListener(t, ln, "", core)
	sys := vault.TestDynamicSystemView(core)
	vault.TestAddTestPlugin(t, core, "test-plugin", "TestPlugin_Main")

	return core, ln, sys
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

	dbplugin.NewPluginServer(plugin)
}

func TestPlugin_Initialize(t *testing.T) {
	_, ln, sys := getCore(t)
	defer ln.Close()

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
	_, ln, sys := getCore(t)
	defer ln.Close()

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

	us, pw, err := db.CreateUser(dbplugin.Statements{}, "test", time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if us != "test" || pw != "test" {
		t.Fatal("expected username and password to be 'test'")
	}

	// try and save the same user again to verify it saved the first time, this
	// should return an error
	_, _, err = db.CreateUser(dbplugin.Statements{}, "test", time.Now().Add(time.Minute))
	if err == nil {
		t.Fatal("expected an error, user wasn't created correctly")
	}
}

func TestPlugin_RenewUser(t *testing.T) {
	_, ln, sys := getCore(t)
	defer ln.Close()

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

	us, _, err := db.CreateUser(dbplugin.Statements{}, "test", time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = db.RenewUser(dbplugin.Statements{}, us, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPlugin_RevokeUser(t *testing.T) {
	_, ln, sys := getCore(t)
	defer ln.Close()

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

	us, _, err := db.CreateUser(dbplugin.Statements{}, "test", time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test default revoke statememts
	err = db.RevokeUser(dbplugin.Statements{}, us)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Try adding the same username back so we can verify it was removed
	_, _, err = db.CreateUser(dbplugin.Statements{}, "test", time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
