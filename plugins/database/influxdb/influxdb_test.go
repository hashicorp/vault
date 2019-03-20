package influxdb

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/ory/dockertest"
)

const testInfluxRole = `CREATE USER "{{username}}" WITH PASSWORD '{{password}}';GRANT ALL ON "vault" TO "{{username}}";`

func prepareInfluxdbTestContainer(t *testing.T) (func(), string, int) {
	if os.Getenv("INFLUXDB_HOST") != "" {
		return func() {}, os.Getenv("INFLUXDB_HOST"), 0
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	ro := &dockertest.RunOptions{
		Repository: "influxdb",
		Tag:        "alpine",
		Env:        []string{"INFLUXDB_DB=vault", "INFLUXDB_ADMIN_USER=influx-root", "INFLUXDB_ADMIN_PASSWORD=influx-root", "INFLUXDB_HTTP_AUTH_ENABLED=true"},
	}
	resource, err := pool.RunWithOptions(ro)
	if err != nil {
		t.Fatalf("Could not start local influxdb docker container: %s", err)
	}

	cleanup := func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	port, _ := strconv.Atoi(resource.GetPort("8086/tcp"))
	address := "127.0.0.1"

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		cli, err := influx.NewHTTPClient(influx.HTTPConfig{
			Addr:     fmt.Sprintf("http://%s:%d", address, port),
			Username: "influx-root",
			Password: "influx-root",
		})
		if err != nil {
			return errwrap.Wrapf("Error creating InfluxDB Client: ", err)
		}
		defer cli.Close()
		_, _, err = cli.Ping(1)
		if err != nil {
			return errwrap.Wrapf("error checking cluster status: {{err}}", err)
		}
		return nil
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to influxdb docker container: %s", err)
	}
	return cleanup, address, port
}

func TestInfluxdb_Initialize(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	cleanup, address, port := prepareInfluxdbTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"host":     address,
		"port":     port,
		"username": "influx-root",
		"password": "influx-root",
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
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

	// test a string protocol
	connectionDetails = map[string]interface{}{
		"host":     address,
		"port":     strconv.Itoa(port),
		"username": "influx-root",
		"password": "influx-root",
	}

	_, err = db.Init(context.Background(), connectionDetails, true)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestInfluxdb_CreateUser(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	cleanup, address, port := prepareInfluxdbTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"host":     address,
		"port":     port,
		"username": "influx-root",
		"password": "influx-root",
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
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

	if err := testCredsExist(t, address, port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}
}

func TestMyInfluxdb_RenewUser(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	cleanup, address, port := prepareInfluxdbTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"host":     address,
		"port":     port,
		"username": "influx-root",
		"password": "influx-root",
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
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

	if err := testCredsExist(t, address, port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	err = db.RenewUser(context.Background(), statements, username, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestInfluxdb_RevokeUser(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	cleanup, address, port := prepareInfluxdbTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"host":     address,
		"port":     port,
		"username": "influx-root",
		"password": "influx-root",
	}

	db := new()
	_, err := db.Init(context.Background(), connectionDetails, true)
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

	if err = testCredsExist(t, address, port, username, password); err != nil {
		t.Fatalf("Could not connect with new credentials: %s", err)
	}

	// Test default revoke statements
	err = db.RevokeUser(context.Background(), statements, username)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err = testCredsExist(t, address, port, username, password); err == nil {
		t.Fatal("Credentials were not revoked")
	}
}
func TestInfluxdb_RotateRootCredentials(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}
	cleanup, address, port := prepareInfluxdbTestContainer(t)
	defer cleanup()

	connectionDetails := map[string]interface{}{
		"host":     address,
		"port":     port,
		"username": "influx-root",
		"password": "influx-root",
	}

	db := new()

	connProducer := db.influxdbConnectionProducer

	_, err := db.Init(context.Background(), connectionDetails, true)
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

func testCredsExist(t testing.TB, address string, port int, username, password string) error {
	cli, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%d", address, port),
		Username: username,
		Password: password,
	})
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
	if response.Error() != nil {
		return errwrap.Wrapf("error using the correct influx database: {{err}}", response.Error())
	}
	return nil
}
