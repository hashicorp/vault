package mongodb

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"gopkg.in/mgo.v2"
)

// PrepareTestContainer calls PrepareTestContainerWithDatabase without a
// database name value, which results in configuring a database named "test"
func PrepareTestContainer(t *testing.T, version string) (cleanup func(), retURL string) {
	return PrepareTestContainerWithDatabase(t, version, "")
}

// PrepareTestContainerWithDatabase configures a test container with a given
// database name, to test non-test/admin database configurations
func PrepareTestContainerWithDatabase(t *testing.T, version, dbName string) (cleanup func(), retURL string) {
	if os.Getenv("MONGODB_URL") != "" {
		return func() {}, os.Getenv("MONGODB_URL")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", version, []string{})
	if err != nil {
		t.Fatalf("Could not start local mongo docker container: %s", err)
	}

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	retURL = fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))
	if dbName != "" {
		retURL = fmt.Sprintf("%s/%s", retURL, dbName)
	}

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		dialInfo, err := parseMongoURL(retURL)
		if err != nil {
			return err
		}

		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			return err
		}
		defer session.Close()
		session.SetSyncTimeout(1 * time.Minute)
		session.SetSocketTimeout(1 * time.Minute)
		return session.Ping()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to mongo docker container: %s", err)
	}

	return
}

// parseMongoURL will parse a connection string and return a configured dialer
func parseMongoURL(rawURL string) (*mgo.DialInfo, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	info := mgo.DialInfo{
		Addrs:    strings.Split(url.Host, ","),
		Database: strings.TrimPrefix(url.Path, "/"),
		Timeout:  10 * time.Second,
	}

	if url.User != nil {
		info.Username = url.User.Username()
		info.Password, _ = url.User.Password()
	}

	query := url.Query()
	for key, values := range query {
		var value string
		if len(values) > 0 {
			value = values[0]
		}

		switch key {
		case "authSource":
			info.Source = value
		case "authMechanism":
			info.Mechanism = value
		case "gssapiServiceName":
			info.Service = value
		case "replicaSet":
			info.ReplicaSetName = value
		case "maxPoolSize":
			poolLimit, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("bad value for maxPoolSize: " + value)
			}
			info.PoolLimit = poolLimit
		case "ssl":
			// Unfortunately, mgo doesn't support the ssl parameter in its MongoDB URI parsing logic, so we have to handle that
			// ourselves. See https://github.com/go-mgo/mgo/issues/84
			ssl, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.New("bad value for ssl: " + value)
			}
			if ssl {
				info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
					return tls.Dial("tcp", addr.String(), &tls.Config{})
				}
			}
		case "connect":
			if value == "direct" {
				info.Direct = true
				break
			}
			if value == "replicaSet" {
				break
			}
			fallthrough
		default:
			return nil, errors.New("unsupported connection URL option: " + key + "=" + value)
		}
	}

	return &info, nil
}
