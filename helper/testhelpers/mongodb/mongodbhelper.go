package mongodb

import (
	"context"
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

	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"gopkg.in/mgo.v2"
)

// PrepareTestContainer calls PrepareTestContainerWithDatabase without a
// database name value, which results in configuring a database named "test"
func PrepareTestContainer(t *testing.T, version string) (cleanup func(), retURL string) {
	return PrepareTestContainerWithDatabase(t, version, "")
}

// PrepareTestContainerWithDatabase configures a test container with a given
// database name, to test non-test/admin database configurations
func PrepareTestContainerWithDatabase(t *testing.T, version, dbName string) (func(), string) {
	if os.Getenv("MONGODB_URL") != "" {
		return func() {}, os.Getenv("MONGODB_URL")
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo: "mongo",
		ImageTag:  version,
		Ports:     []string{"27017/tcp"},
	})
	if err != nil {
		t.Fatalf("could not start docker mongo: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		connURL := fmt.Sprintf("mongodb://%s:%d", host, port)
		if dbName != "" {
			connURL = fmt.Sprintf("%s/%s", connURL, dbName)
		}
		dialInfo, err := ParseMongoURL(connURL)
		if err != nil {
			return nil, err
		}

		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			return nil, err
		}
		defer session.Close()

		session.SetSyncTimeout(1 * time.Minute)
		session.SetSocketTimeout(1 * time.Minute)
		err = session.Ping()
		if err != nil {
			return nil, err
		}

		return docker.NewServiceURLParse(connURL)
	})

	if err != nil {
		t.Fatalf("could not start docker mongo: %s", err)
	}

	return svc.Cleanup, svc.Config.URL().String()
}

// ParseMongoURL will parse a connection string and return a configured dialer
func ParseMongoURL(rawURL string) (*mgo.DialInfo, error) {
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
