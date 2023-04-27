// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodb

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	paths "path"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers/certhelpers"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/ory/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestInit_clientTLS(t *testing.T) {
	t.Skip("Skipping this test because CircleCI can't mount the files we need without further investigation: " +
		"https://support.circleci.com/hc/en-us/articles/360007324514-How-can-I-mount-volumes-to-docker-containers-")

	// Set up temp directory so we can mount it to the docker container
	confDir := makeTempDir(t)
	defer os.RemoveAll(confDir)

	// Create certificates for Mongo authentication
	caCert := certhelpers.NewCert(t,
		certhelpers.CommonName("test certificate authority"),
		certhelpers.IsCA(true),
		certhelpers.SelfSign(),
	)
	serverCert := certhelpers.NewCert(t,
		certhelpers.CommonName("server"),
		certhelpers.DNS("localhost"),
		certhelpers.Parent(caCert),
	)
	clientCert := certhelpers.NewCert(t,
		certhelpers.CommonName("client"),
		certhelpers.DNS("client"),
		certhelpers.Parent(caCert),
	)

	writeFile(t, paths.Join(confDir, "ca.pem"), caCert.CombinedPEM(), 0o644)
	writeFile(t, paths.Join(confDir, "server.pem"), serverCert.CombinedPEM(), 0o644)
	writeFile(t, paths.Join(confDir, "client.pem"), clientCert.CombinedPEM(), 0o644)

	// //////////////////////////////////////////////////////
	// Set up Mongo config file
	rawConf := `
net:
   tls:
      mode: preferTLS
      certificateKeyFile: /etc/mongo/server.pem
      CAFile: /etc/mongo/ca.pem
      allowInvalidHostnames: true`

	writeFile(t, paths.Join(confDir, "mongod.conf"), []byte(rawConf), 0o644)

	// //////////////////////////////////////////////////////
	// Start Mongo container
	retURL, cleanup := startMongoWithTLS(t, "latest", confDir)
	defer cleanup()

	// //////////////////////////////////////////////////////
	// Set up x509 user
	mClient := connect(t, retURL)

	setUpX509User(t, mClient, clientCert)

	// //////////////////////////////////////////////////////
	// Test
	mongo := new()

	initReq := dbplugin.InitializeRequest{
		Config: map[string]interface{}{
			"connection_url":      retURL,
			"allowed_roles":       "*",
			"tls_certificate_key": clientCert.CombinedPEM(),
			"tls_ca":              caCert.Pem,
		},
		VerifyConnection: true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := mongo.Initialize(ctx, initReq)
	if err != nil {
		t.Fatalf("Unable to initialize mongo engine: %s", err)
	}

	// Initialization complete. The connection was established, but we need to ensure
	// that we're connected as the right user
	whoamiCmd := map[string]interface{}{
		"connectionStatus": 1,
	}

	client, err := mongo.Connection(ctx)
	if err != nil {
		t.Fatalf("Unable to make connection to Mongo: %s", err)
	}
	result := client.Database("test").RunCommand(ctx, whoamiCmd)
	if result.Err() != nil {
		t.Fatalf("Unable to connect to Mongo: %s", err)
	}

	expected := connStatus{
		AuthInfo: authInfo{
			AuthenticatedUsers: []user{
				{
					User: fmt.Sprintf("CN=%s", clientCert.Template.Subject.CommonName),
					DB:   "$external",
				},
			},
			AuthenticatedUserRoles: []role{
				{
					Role: "readWrite",
					DB:   "test",
				},
				{
					Role: "userAdminAnyDatabase",
					DB:   "admin",
				},
			},
		},
		Ok: 1,
	}
	// Sort the AuthenticatedUserRoles because Mongo doesn't return them in the same order every time
	// Thanks Mongo! /tableflip
	sort.Sort(expected.AuthInfo.AuthenticatedUserRoles)

	actual := connStatus{}
	err = result.Decode(&actual)
	if err != nil {
		t.Fatalf("Unable to decode connection status: %s", err)
	}

	sort.Sort(actual.AuthInfo.AuthenticatedUserRoles)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Actual:%#v\nExpected:\n%#v", actual, expected)
	}
}

func makeTempDir(t *testing.T) (confDir string) {
	confDir, err := ioutil.TempDir(".", "mongodb-test-data")
	if err != nil {
		t.Fatalf("Unable to make temp directory: %s", err)
	}
	// Convert the directory to an absolute path because docker needs it when mounting
	confDir, err = filepath.Abs(filepath.Clean(confDir))
	if err != nil {
		t.Fatalf("Unable to determine where temp directory is on absolute path: %s", err)
	}
	return confDir
}

func startMongoWithTLS(t *testing.T, version string, confDir string) (retURL string, cleanup func()) {
	if os.Getenv("MONGODB_URL") != "" {
		return os.Getenv("MONGODB_URL"), func() {}
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}
	pool.MaxWait = 30 * time.Second

	containerName := "mongo-unit-test"

	// Remove previously running container if it is still running because cleanup failed
	err = pool.RemoveContainerByName(containerName)
	if err != nil {
		t.Fatalf("Unable to remove old running containers: %s", err)
	}

	runOpts := &dockertest.RunOptions{
		Name:       containerName,
		Repository: "mongo",
		Tag:        version,
		Cmd:        []string{"mongod", "--config", "/etc/mongo/mongod.conf"},
		// Mount the directory from local filesystem into the container
		Mounts: []string{
			fmt.Sprintf("%s:/etc/mongo", confDir),
		},
	}

	resource, err := pool.RunWithOptions(runOpts)
	if err != nil {
		t.Fatalf("Could not start local mongo docker container: %s", err)
	}
	resource.Expire(30)

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	uri := url.URL{
		Scheme: "mongodb",
		Host:   fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp")),
	}
	retURL = uri.String()

	// exponential backoff-retry
	err = pool.Retry(func() error {
		var err error
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(retURL))
		if err = client.Disconnect(ctx); err != nil {
			t.Fatal()
		}
		return client.Ping(ctx, readpref.Primary())
	})
	if err != nil {
		cleanup()
		t.Fatalf("Could not connect to mongo docker container: %s", err)
	}

	return retURL, cleanup
}

func connect(t *testing.T, uri string) (client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("Unable to make connection to Mongo: %s", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Fatalf("Failed to ping Mongo server: %s", err)
	}

	return client
}

func setUpX509User(t *testing.T, client *mongo.Client, cert certhelpers.Certificate) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := fmt.Sprintf("CN=%s", cert.Template.Subject.CommonName)

	cmd := &createUserCommand{
		Username: username,
		Roles: []interface{}{
			mongodbRole{
				Role: "readWrite",
				DB:   "test",
			},
			mongodbRole{
				Role: "userAdminAnyDatabase",
				DB:   "admin",
			},
		},
	}

	result := client.Database("$external").RunCommand(ctx, cmd)
	err := result.Err()
	if err != nil {
		t.Fatalf("Failed to create x509 user in database: %s", err)
	}
}

type connStatus struct {
	AuthInfo authInfo `bson:"authInfo"`
	Ok       int      `bson:"ok"`
}

type authInfo struct {
	AuthenticatedUsers     []user `bson:"authenticatedUsers"`
	AuthenticatedUserRoles roles  `bson:"authenticatedUserRoles"`
}

type user struct {
	User string `bson:"user"`
	DB   string `bson:"db"`
}

type role struct {
	Role string `bson:"role"`
	DB   string `bson:"db"`
}

type roles []role

func (r roles) Len() int           { return len(r) }
func (r roles) Less(i, j int) bool { return r[i].Role < r[j].Role }
func (r roles) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// ////////////////////////////////////////////////////////////////////////////
// Writing to file
// ////////////////////////////////////////////////////////////////////////////
func writeFile(t *testing.T, filename string, data []byte, perms os.FileMode) {
	t.Helper()

	err := ioutil.WriteFile(filename, data, perms)
	if err != nil {
		t.Fatalf("Unable to write to file [%s]: %s", filename, err)
	}
}
