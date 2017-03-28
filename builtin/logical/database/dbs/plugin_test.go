package dbs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	log "github.com/mgutz/logxi/v1"
)

type mockPlugin struct {
	users map[string][]string
	CredentialsProducer
}

func (m *mockPlugin) Type() string { return "mock" }
func (m *mockPlugin) CreateUser(statements Statements, username, password, expiration string) error {
	err := errors.New("err")
	if username == "" || password == "" || expiration == "" {
		return err
	}

	if _, ok := m.users[username]; ok {
		return err
	}

	m.users[username] = []string{password, expiration}

	return nil
}
func (m *mockPlugin) RenewUser(statements Statements, username, expiration string) error {
	err := errors.New("err")
	if username == "" || expiration == "" {
		return err
	}

	if _, ok := m.users[username]; !ok {
		return err
	}

	return nil
}
func (m *mockPlugin) RevokeUser(statements Statements, username string) error {
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
func (m *mockPlugin) Initialize(conf map[string]interface{}) error {
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

func getConf(t *testing.T) *DatabaseConfig {
	command := fmt.Sprintf("%s -test.run=TestPlugin_Main", os.Args[0])
	cmd := exec.Command(os.Args[0])
	hash := sha256.New()

	file, err := os.Open(cmd.Path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(hash, file)
	if err != nil {
		t.Fatal(err)
	}

	sum := hash.Sum(nil)

	conf := &DatabaseConfig{
		DatabaseType:   pluginTypeName,
		PluginCommand:  command,
		PluginChecksum: hex.EncodeToString(sum),
		ConnectionDetails: map[string]interface{}{
			"test": true,
		},
	}

	return conf
}

func getCore(t *testing.T) (*vault.Core, net.Listener, logical.SystemView) {
	core, _, _, ln := vault.TestCoreUnsealedWithListener(t)
	http.TestServerWithListener(t, ln, "", core)
	sys := vault.TestDynamicSystemView(core)

	return core, ln, sys
}

// This is not an actual test case, it's a helper function that will be executed
// by the go-plugin client via an exec call.
func TestPlugin_Main(t *testing.T) {
	if os.Getenv(pluginutil.PluginUnwrapTokenEnv) == "" {
		return
	}

	plugin := &mockPlugin{
		users:               make(map[string][]string),
		CredentialsProducer: &sqlCredentialsProducer{5, 50},
	}

	NewPluginServer(plugin)
}

func TestPlugin_Initialize(t *testing.T) {
	_, ln, sys := getCore(t)
	defer ln.Close()

	conf := getConf(t)
	dbRaw, err := PluginFactory(conf, sys, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = dbRaw.Initialize(conf.ConnectionDetails)
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

	conf := getConf(t)
	db, err := PluginFactory(conf, sys, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer db.Close()

	err = db.Initialize(conf.ConnectionDetails)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	username, err := db.GenerateUsername("test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err := db.GeneratePassword()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expiration, err := db.GenerateExpiration(time.Minute)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = db.CreateUser(Statements{}, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	// try and save the same user again to verify it saved the first time, this
	// should return an error
	err = db.CreateUser(Statements{}, username, password, expiration)
	if err == nil {
		t.Fatal("expected an error, user wasn't created correctly")
	}

	// Create one more user
	username, err = db.GenerateUsername("test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = db.CreateUser(Statements{}, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPlugin_RenewUser(t *testing.T) {
	_, ln, sys := getCore(t)
	defer ln.Close()

	conf := getConf(t)
	db, err := PluginFactory(conf, sys, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer db.Close()

	err = db.Initialize(conf.ConnectionDetails)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	username, err := db.GenerateUsername("test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err := db.GeneratePassword()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expiration, err := db.GenerateExpiration(time.Minute)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = db.CreateUser(Statements{}, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expiration, err = db.GenerateExpiration(time.Minute)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = db.RenewUser(Statements{}, username, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPlugin_RevokeUser(t *testing.T) {
	_, ln, sys := getCore(t)
	defer ln.Close()

	conf := getConf(t)
	db, err := PluginFactory(conf, sys, &log.NullLogger{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer db.Close()

	err = db.Initialize(conf.ConnectionDetails)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	username, err := db.GenerateUsername("test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	password, err := db.GeneratePassword()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expiration, err := db.GenerateExpiration(time.Minute)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = db.CreateUser(Statements{}, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test default revoke statememts
	err = db.RevokeUser(Statements{}, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Try adding the same username back so we can verify it was removed
	err = db.CreateUser(Statements{}, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	username, err = db.GenerateUsername("test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expiration, err = db.GenerateExpiration(time.Minute)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// try once more
	err = db.CreateUser(Statements{}, username, password, expiration)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = db.RevokeUser(Statements{}, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

}
