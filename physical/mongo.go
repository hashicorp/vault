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

	logger.Printf("[DEBUG]: physical/mongo: newMongoBackend: (%v, %v, %v)", url, database, collection)
	return &MongoBackend{
		Url:        url,
		Database:   database,
		Collection: collection,
		logger:     logger,
	}, nil
}

func (b *MongoBackend) Delete(k string) error {
	b.logger.Printf("[DEBUG]: physical/mongo: Delete(%v)", k)
	session, err := b.activeSession()
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: could not establish mongo session: %v", err)
		return err
	}
	c := session.DB(b.Database).C(b.Collection)
	err = c.Remove(bson.M{"key": k})
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: Delete(%v): error in Remove() : %v", k, err)
		return err
	}
	return nil
}

func (b *MongoBackend) Get(k string) (*Entry, error) {
	b.logger.Printf("[DEBUG]: physical/mongo: Get(%v)", k)
	session, err := b.activeSession()
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: could not establish mongo session: %v", err)
		return nil, err
	}
	c := session.DB(b.Database).C(b.Collection)

	q := c.Find(bson.M{"key": k})
	n, err := q.Count()
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: Get(%v): error in Count() : %v", k, err)
		return nil, err
	}

	// not found requires us to return nil and not throw an error
	// make an exception
	if n <= 0 {
		b.logger.Printf("[DEBUG]: physical/mongo: Get(%v): not found", k)
		return nil, nil
	}
	var entry Entry
	err = q.One(&entry)
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: Get(%v): error in One(): %v", k, err)
		return nil, err
	}
	return &entry, nil
}

func (b *MongoBackend) Put(entry *Entry) error {
	b.logger.Printf("[DEBUG]: physical/mongo: Put(%v)", entry.Key)
	session, err := b.activeSession()
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: could not establish mongo session: %v", err)
		return err
	}
	c := session.DB(b.Database).C(b.Collection)
	_, err = c.Upsert(bson.M{"key": entry.Key}, entry)
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: Put(%v): error in Upsert(): %v", entry.Key, err)
		return err
	}
	return nil
}

func (b *MongoBackend) List(prefix string) ([]string, error) {
	b.logger.Printf("[DEBUG]: physical/mongo: List(%v)", prefix)
	session, err := b.activeSession()
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: could not establish mongo session: %v", err)
		return nil, err
	}
	c := session.DB(b.Database).C(b.Collection)

	var results []string

	p := strings.Replace(prefix, "/", "\\/", -1)
	regex := `^` + p + `[^/]*`

	iter := c.Find(bson.M{"key": bson.M{"$regex": bson.RegEx{Pattern: regex, Options: ""}}}).
		Select(bson.M{"key": 1}).
		Iter()

	var result Entry
	for iter.Next(&result) {
		b.logger.Printf("[DEBUG]: physical/mongo: List(%v): Next(%v)", prefix, result)
		results = append(results, strings.TrimPrefix(result.Key, prefix))
	}

	err = iter.Close()
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: List(%v): error in iter.Next(): %v", prefix, err)
		return nil, err
	}

	b.logger.Printf("[DEBUG]: physical/mongo: List(%v): %v", prefix, results)

	return results, nil
}

// func (b *MongoBackend) path(k string) (string, string) {
// 	path := filepath.Join(b.Path, k)
// 	key := filepath.Base(path)
// 	path = filepath.Dir(path)
// 	return path, "_" + key
// }
