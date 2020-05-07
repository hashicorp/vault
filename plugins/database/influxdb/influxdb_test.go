package influxdb

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	influx "github.com/influxdata/influxdb/client/v2"
)

const testInfluxRole = `CREATE USER "{{username}}" WITH PASSWORD '{{password}}';GRANT ALL ON "vault" TO "{{username}}";`

type Config struct {
	docker.ServiceURL
	Username string
	Password string
}

var _ docker.ServiceConfig = &Config{}

func (c *Config) apiConfig() influx.HTTPConfig {
	return influx.HTTPConfig{
		Addr:     c.URL().String(),
		Username: c.Username,
		Password: c.Password,
	}
}

func (c *Config) connectionParams() map[string]interface{} {
	pieces := strings.Split(c.Address(), ":")
	port, _ := strconv.Atoi(pieces[1])
	return map[string]interface{}{
		"host":     pieces[0],
		"port":     port,
		"username": c.Username,
		"password": c.Password,
	}
}

func prepareInfluxdbTestContainer(t *testing.T) (func(), *Config) {
	c := &Config{
		Username: "influx-root",
		Password: "influx-root",
	}
	if host := os.Getenv("INFLUXDB_HOST"); host != "" {
		c.ServiceURL = *docker.NewServiceURL(url.URL{Scheme: "http", Host: host})
		return func() {}, c
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo: "influxdb",
		ImageTag:  "alpine",
		Env: []string{
			"INFLUXDB_DB=vault",
			"INFLUXDB_ADMIN_USER=" + c.Username,
			"INFLUXDB_ADMIN_PASSWORD=" + c.Password,
			"INFLUXDB_HTTP_AUTH_ENABLED=true"},
		Ports: []string{"8086/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start docker InfluxDB: %s", err)
	}
	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		c.ServiceURL = *docker.NewServiceURL(url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", host, port),
		})
		cli, err := influx.NewHTTPClient(c.apiConfig())
		if err != nil {
			return nil, errwrap.Wrapf("error creating InfluxDB client: {{err}}", err)
		}
		defer cli.Close()
		_, _, err = cli.Ping(1)
		if err != nil {
			return nil, errwrap.Wrapf("error checking cluster status: {{err}}", err)
		}

		return c, nil
	})
	if err != nil {
		t.Fatalf("Could not start docker InfluxDB: %s", err)
	}

	return svc.Cleanup, svc.Config.(*Config)
}

func TestInfluxdb_Initialize(t *testing.T) {
	cleanup, config := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	_, err := db.Init(context.Background(), config.connectionParams(), true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !db.Initialized {
		t.Fatal("Database should be initialized")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	connectionParams := config.connectionParams()
	connectionParams["port"] = strconv.Itoa(connectionParams["port"].(int))
	_, err = db.Init(context.Background(), connectionParams, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestInfluxdb_CreateUser(t *testing.T) {
	cleanup, config := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	_, err := db.Init(context.Background(), config.connectionParams(), true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testInfluxRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	config.Username, config.Password = username, password
	if err := testCredsExist(t, config); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMyInfluxdb_RenewUser(t *testing.T) {
	cleanup, config := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	_, err := db.Init(context.Background(), config.connectionParams(), true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testInfluxRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	config.Username, config.Password = username, password
	if err = testCredsExist(t, config); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = db.RenewUser(context.Background(), statements, username, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

// TestInfluxdb_RevokeDeletedUser tests attempting to revoke a user that was
// deleted externally. Guards against a panic, see
// https://github.com/hashicorp/vault/issues/6734
func TestInfluxdb_RevokeDeletedUser(t *testing.T) {
	cleanup, config := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	_, err := db.Init(context.Background(), config.connectionParams(), true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testInfluxRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	config.Username, config.Password = username, password
	if err = testCredsExist(t, config); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// call cleanup to remove database
	cleanup()

	// attempt to revoke the user after database is gone
	err = db.RevokeUser(context.Background(), statements, username)
	if err == nil {
		t.Fatalf("Expected err, got nil")
	}
}

func TestInfluxdb_RevokeUser(t *testing.T) {
	cleanup, config := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()
	_, err := db.Init(context.Background(), config.connectionParams(), true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	statements := dbplugin.Statements{
		Creation: []string{testInfluxRole},
	}

	usernameConfig := dbplugin.UsernameConfig{
		DisplayName: "test",
		RoleName:    "test",
	}

	username, password, err := db.CreateUser(context.Background(), statements, usernameConfig, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	config.Username, config.Password = username, password
	if err = testCredsExist(t, config); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revoke statements
	err = db.RevokeUser(context.Background(), dbplugin.Statements{}, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, config); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}
func TestInfluxdb_RotateRootCredentials(t *testing.T) {
	cleanup, config := prepareInfluxdbTestContainer(t)
	defer cleanup()

	db := new()

	connProducer := db.influxdbConnectionProducer

	_, err := db.Init(context.Background(), config.connectionParams(), true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !connProducer.Initialized {
		t.Fatal("Database should be initialized")
	}

	newConf, err := db.RotateRootCredentials(context.Background(), nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if newConf["password"] == "influx-root" {
		t.Fatal("password was not updated")
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testCredsExist(t testing.TB, c *Config) error {
	cli, err := influx.NewHTTPClient(c.apiConfig())
	if err != nil {
		return errwrap.Wrapf("Error creating InfluxDB Client: ", err)
	}
	defer cli.Close()
	_, _, err = cli.Ping(1)
	if err != nil {
		return errwrap.Wrapf("error checking server ping: {{err}}", err)
	}
	q := influx.NewQuery("SHOW SERIES ON vault", "", "")
	response, err := cli.Query(q)
	if err != nil {
		return errwrap.Wrapf("error querying influxdb server: {{err}}", err)
	}
	if response != nil && response.Error() != nil {
		return errwrap.Wrapf("error using the correct influx database: {{err}}", response.Error())
	}
	return nil
}
