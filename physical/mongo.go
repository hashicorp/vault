package physical

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Unfortunately, mgo doesn't support the ssl parameter in its MongoDB URI parsing logic, so we have to handle that
// ourselves. See https://github.com/go-mgo/mgo/issues/84
func parseMongoURI(rawUri string) (*mgo.DialInfo, error) {
	uri, err := url.Parse(rawUri)
	if err != nil {
		return nil, err
	}

	info := mgo.DialInfo{
		Addrs:    strings.Split(uri.Host, ","),
		Database: strings.TrimPrefix(uri.Path, "/"),
		Timeout:  10 * time.Second,
	}

	if uri.User != nil {
		info.Username = uri.User.Username()
		info.Password, _ = uri.User.Password()
	}

	query := uri.Query()
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

func (b *MongoBackend) activeSession() (*mgo.Session, error) {
	b.l.Lock()
	defer b.l.Unlock()

	if b.session != nil {
		if err := b.session.Ping(); err == nil {
			return b.session, nil
		}
		b.session.Close()
	}

	dialInfo, err := parseMongoURI(b.Url)
	if err != nil {
		return nil, err
	}

	b.session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}
	b.session.SetSyncTimeout(1 * time.Minute)
	b.session.SetSocketTimeout(1 * time.Minute)

	// handled by mgo
	//if b.Database == "" {
	//	b.Database = dialInfo.Database
	//}

	return b.session, nil
}

// MongoBackend is a physical backend that stores data on disk
type MongoBackend struct {
	Url        string
	Database   string
	Collection string
	l          sync.Mutex
	session    *mgo.Session
	logger     *log.Logger
}

// newMongoBackend constructs a MongoBackend using the given directory
func newMongoBackend(conf map[string]string, logger *log.Logger) (Backend, error) {
	url, ok := conf["url"]
	if !ok {
		return nil, fmt.Errorf("'url' must be set")
	}

	database, ok := conf["database"]
	if !ok {
		database = ""
	}

	collection, ok := conf["collection"]
	if !ok {
		collection = "vault"
	}

	return &MongoBackend{
		Url:        url,
		Database:   database,
		Collection: collection,
		logger:     logger,
	}, nil
}

func (b *MongoBackend) Delete(k string) error {
	session, err := b.activeSession()
	if err != nil {
		return err
	}
	c := session.DB(b.Database).C(b.Collection)
	return c.Remove(bson.M{"Key": k})
}

func (b *MongoBackend) Get(k string) (*Entry, error) {
	session, err := b.activeSession()
	if err != nil {
		return nil, err
	}
	c := session.DB(b.Database).C(b.Collection)

	var entry Entry
	err = c.Find(bson.M{"Key": k}).One(&entry)
	return &entry, err
}

func (b *MongoBackend) Put(entry *Entry) error {
	session, err := b.activeSession()
	if err != nil {
		return err
	}
	c := session.DB(b.Database).C(b.Collection)
	_, err = c.Upsert(bson.M{"Key": entry.Key}, entry)
	return err
}

func (b *MongoBackend) List(prefix string) ([]string, error) {
	session, err := b.activeSession()
	if err != nil {
		return nil, err
	}
	c := session.DB(b.Database).C(b.Collection)
	var results []string
	err = c.Find(bson.M{"Key": prefix}).Select(bson.M{"Key": 1}).All(&results)
	return results, err
}

// func (b *MongoBackend) path(k string) (string, string) {
// 	path := filepath.Join(b.Path, k)
// 	key := filepath.Base(path)
// 	path = filepath.Dir(path)
// 	return path, "_" + key
// }
