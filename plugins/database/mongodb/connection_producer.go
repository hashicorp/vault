package mongodb

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	"github.com/mitchellh/mapstructure"

	"gopkg.in/mgo.v2"
)

// mongoDBConnectionProducer implements ConnectionProducer and provides an
// interface for databases to make connections.
type mongoDBConnectionProducer struct {
	ConnectionURL string `json:"connection_url" structs:"connection_url" mapstructure:"connection_url"`

	Initialized bool
	Type        string
	session     *mgo.Session
	sync.Mutex
}

// Initialize parses connection configuration.
func (c *mongoDBConnectionProducer) Initialize(conf map[string]interface{}, verifyConnection bool) error {
	c.Lock()
	defer c.Unlock()

	err := mapstructure.WeakDecode(conf, c)
	if err != nil {
		return err
	}

	if len(c.ConnectionURL) == 0 {
		return fmt.Errorf("connection_url cannot be empty")
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	if verifyConnection {
		if _, err := c.Connection(); err != nil {
			return fmt.Errorf("error verifying connection: %s", err)
		}

		if err := c.session.Ping(); err != nil {
			return fmt.Errorf("error verifying connection: %s", err)
		}
	}

	return nil
}

// Connection creates a database connection.
func (c *mongoDBConnectionProducer) Connection() (interface{}, error) {
	if !c.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	if c.session != nil {
		return c.session, nil
	}

	dialInfo, err := parseMongoURL(c.ConnectionURL)
	if err != nil {
		return nil, err
	}

	c.session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}
	c.session.SetSyncTimeout(1 * time.Minute)
	c.session.SetSocketTimeout(1 * time.Minute)

	return nil, nil
}

// Close terminates the database connection.
func (c *mongoDBConnectionProducer) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.session != nil {
		c.session.Close()
	}

	c.session = nil

	return nil
}

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
