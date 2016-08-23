package physical

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// this part was derived from and adjusted for our needs from:
//   builtin/logical/mongodb/util.go
//
// Unfortunately, mgo doesn't support the ssl parameter in its MongoDB URI parsing logic, so we have to handle that
// ourselves. See https://github.com/go-mgo/mgo/issues/84
func (b *MongoBackend) parseMongoURI(rawUri string) (*mgo.DialInfo, error) {
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

	uriSsl := false

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
				uriSsl = true
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

	// deal with TLS
	if uriSsl || b.tls {
		tlsConfig := tls.Config{}

		if b.tlsSkipVerify {
			tlsConfig.InsecureSkipVerify = true
		}

		if b.tlsCAFile != "" {
			caBytes, err := ioutil.ReadFile(b.tlsCAFile)
			if err != nil {
				return nil, errors.New("could not read CA data from '" + b.tlsCAFile + "'")
			}
			caPool := x509.NewCertPool()
			ok := caPool.AppendCertsFromPEM(caBytes)
			if !ok {
				b.logger.Printf("[WARN]: physical/mongo: could not parse CAs from '%v'. Are you sure they are PEM encoded?", b.tlsCAFile)
			}
			tlsConfig.RootCAs = caPool
		}

		if b.tlsKeyFile != "" && b.tlsCertFile != "" {
			cert, err := tls.LoadX509KeyPair(b.tlsCertFile, b.tlsKeyFile)
			if err != nil {
				return nil, errors.New("could not load cert and/or key from '" + b.tlsCertFile + " / '" + b.tlsKeyFile + "'")
			}

			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tlsConfig)
		}
	}

	return &info, nil
}

var keyIndex mgo.Index = mgo.Index{
	Key:    []string{"key"},
	Unique: true,
}

var lastCheckedInIndex mgo.Index = mgo.Index{
	Key:         []string{"lastCheckedIn"},
	ExpireAfter: 5 * time.Second,
}

// this part was derived from and adjusted for our needs from:
//   builtin/logical/mongodb/backend.go
//
func (b *MongoBackend) activeSession() (*mgo.Session, error) {
	b.l.Lock()
	defer b.l.Unlock()

	if b.session != nil {
		if err := b.session.Ping(); err == nil {
			return b.session, nil
		}
		b.session.Close()
	}

	// we need to establish a new session
	b.logger.Printf("[INFO]: physical/mongo: establishing new MongoDB session")
	dialInfo, err := b.parseMongoURI(b.Url)
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: could not parse MongoDB URI: %v", err)
		return nil, err
	}

	b.session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: could not establish connection to MongoDB: %v", err)
		return nil, err
	}
	b.session.SetSyncTimeout(1 * time.Minute)
	b.session.SetSocketTimeout(1 * time.Minute)

	// ensure indices
	b.logger.Printf("[INFO]: physical/mongo: ensuring MongoDB indices exist")
	c := b.session.DB(b.Database).C(b.Collection)
	err = c.EnsureIndex(keyIndex)
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: could not EnsureIndex() for 'key': %v", err)
		return nil, err
	}

	err = c.EnsureIndex(lastCheckedInIndex)
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: could not EnsureIndex() for 'lastCheckedIn': %v", err)
		return nil, err
	}

	return b.session, nil
}

// MongoBackend is a physical backend that stores data on disk
type MongoBackend struct {
	Url           string
	Database      string
	Collection    string
	l             sync.Mutex
	session       *mgo.Session
	logger        *log.Logger
	haEnabled     bool
	lockOnPrimary string
	tls           bool
	tlsSkipVerify bool
	tlsCertFile   string
	tlsKeyFile    string
	tlsCAFile     string
	lockValue     string
	leaderCh      chan struct{}
	isLeader      bool
}

// newMongoBackend constructs a MongoBackend using the given directory
func newMongoBackend(conf map[string]string, logger *log.Logger) (Backend, error) {
	url, ok := conf["url"]
	if !ok {
		url = "mongodb://127.0.0.1:27017/vault"
	}

	database, ok := conf["database"]
	if !ok {
		database = ""
	}

	collection, ok := conf["collection"]
	if !ok {
		collection = "vault"
	}

	tls := false
	_, ok = conf["tls"]
	if ok {
		tls = true
	}

	haEnabled := false
	_, ok = conf["ha_enabled"]
	if ok {
		haEnabled = true
	}

	lockOnPrimary, ok := conf["lock_on_primary"]
	if !ok {
		lockOnPrimary = ""
	}

	tlsSkipVerify := false
	_, ok = conf["tls_skip_verify"]
	if ok {
		tlsSkipVerify = true
	}

	tlsCAFile, ok := conf["tls_ca_file"]
	if !ok {
		tlsCAFile = ""
	}

	tlsCertFile, ok := conf["tls_cert_file"]
	if !ok {
		tlsCertFile = ""
	}

	tlsKeyFile, ok := conf["tls_key_file"]
	if !ok {
		tlsKeyFile = ""
	}

	// TODO: add TLS config options

	logger.Printf("[DEBUG]: physical/mongo: newMongoBackend: (%v, %v, %v)", url, database, collection)
	return &MongoBackend{
		Url:           url,
		Database:      database,
		Collection:    collection,
		logger:        logger,
		haEnabled:     haEnabled,
		lockOnPrimary: lockOnPrimary,
		tls:           tls,
		tlsSkipVerify: tlsSkipVerify,
		tlsCAFile:     tlsCAFile,
		tlsCertFile:   tlsCertFile,
		tlsKeyFile:    tlsKeyFile,
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

	// The prefix needs to get its slashes replaced with '\/' so that it can form
	// a proper regex
	p := strings.Replace(prefix, "/", "\\/", -1)
	regex := `^` + p

	iter := c.Find(bson.M{"key": bson.M{"$regex": bson.RegEx{Pattern: regex, Options: ""}}}).
		Select(bson.M{"key": 1}).
		Iter()

	var result Entry
	for iter.Next(&result) {
		b.logger.Printf("[DEBUG]: physical/mongo: List(%v): Next(%v)", prefix, result)
		// we remove the prefix from the result and add it to the return list
		key := strings.TrimPrefix(result.Key, prefix)
		if strings.ContainsAny(key, "/") {
			dirKey := strings.SplitAfter(key, "/")[0]
			inResults := false
			for _, a := range results {
				if a == dirKey {
					inResults = true
				}
			}
			if !inResults {
				//results = append(results, dirKey)
				//append([]string{"Prepend Item"}, data...)
				results = append([]string{dirKey}, results...)
			}
		} else {
			results = append(results, key)
		}
	}

	err = iter.Close()
	if err != nil {
		b.logger.Printf("[ERROR]: physical/mongo: List(%v): error in iter.Next(): %v", prefix, err)
		return nil, err
	}

	// sort
	sort.Strings(results)

	b.logger.Printf("[DEBUG]: physical/mongo: List(%v): %v", prefix, results)

	return results, nil
}

// This is configurable right now, as HA should be considered experimental
// The default is off
func (b *MongoBackend) HAEnabled() bool {
	b.logger.Printf("[DEBUG]: physical/mongo: HAEnabled(%v)", b.haEnabled)
	return b.haEnabled
}

// LockWith is used for mutual exclusion based on the given key.
func (b *MongoBackend) LockWith(key, value string) (Lock, error) {
	b.logger.Printf("[DEBUG]: physical/mongo: LockWith(%v, %v)", key, value)
	l := &MongoLock{
		b:     b,
		key:   key,
		value: value,
	}
	return l, nil
}

type MongoLock struct {
	b        *MongoBackend
	key      string
	value    string
	llCh     chan struct{}
	llChOpen bool
}

type LockEntry struct {
	Key           string
	Value         string
	LastCheckedIn bson.MongoTimestamp
}

// Lock is used to acquire the given lock
// The stopCh is optional and if closed should interrupt the lock
// acquisition attempt. The return struct should be closed when
// leadership is lost.
func (l *MongoLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	l.b.logger.Printf("[DEBUG]: physical/mongo: Lock(%v)", l.value)

	// the "read" lock needs special treatment
	if l.value == "read" {
		return l.b.leaderCh, nil
	}

	// get the primary and see if we even want to acquire the lock
	if l.b.lockOnPrimary != "" {
		currentPrimary, err := l.getPrimary()
		if err != nil {
			l.b.logger.Printf("[ERROR]: physical/mongo: could not get primary: %v", err)
			return nil, err
		}

		if l.b.lockOnPrimary != currentPrimary {
			return nil, errors.New("current MongoDB primary is not '" + l.b.lockOnPrimary + "', but '" + currentPrimary + "'")
		}
	}

	session, err := l.b.activeSession()
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: could not establish mongo session: %v", err)
		return nil, err
	}
	c := session.DB(l.b.Database).C(l.b.Collection)

	err = c.Insert(bson.M{"key": l.key, "value": l.value, "lastCheckedIn": bson.Now()})
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: Lock already acquired by other vault instance: %v", err)
		return nil, errors.New("Lock already acquired by other vault instance: " + err.Error())
	}

	// acquired :)
	l.b.lockValue = l.value
	l.b.isLeader = true
	l.b.leaderCh = make(chan struct{})

	// poll and update the lock as long as we are the leader
	abortCh1 := make(chan struct{})
	abortCh2 := make(chan struct{})
	l.llCh = make(chan struct{})
	l.llChOpen = true
	go func() {
	updateLoop:
		for {
			select {
			case <-abortCh1:
				l.b.logger.Printf("[DEBUG]: physical/mongo: abortCh1")
				close(abortCh1)
				break updateLoop
			case <-time.After(time.Second * 3):
				if l.b.isLeader {
					err = l.activePoller()
					if err != nil {
						l.b.logger.Printf("[ERROR]: physical/mongo: l.Lock(): activePoller() error, going to lose leadership: %v", err)
						if l.llChOpen {
							l.llCh <- struct{}{}
						}
					}
				}
			}
		}
	}()

	// poll and break if the primary changes
	if l.b.lockOnPrimary != "" {
		go func() {
		primaryLoop:
			for {
				l.b.logger.Printf("[DEBUG]: physical/mongo: primary monitor")
				select {
				case <-abortCh2:
					l.b.logger.Printf("[DEBUG]: physical/mongo: abortCh2")
					close(abortCh2)
					break primaryLoop
				case <-time.After(time.Second * 1):
					currentPrimary, err := l.getPrimary()
					if err != nil {
						l.b.logger.Printf("[WARN]: physical/mongo: could not get primary: %v", err)
						break
					}

					if l.b.lockOnPrimary != currentPrimary {
						l.b.logger.Printf("[INFO]: physical/mongo: MongoDB PRIMARY changed to '%v'. Dropping leadership.", currentPrimary)
						if l.llChOpen {
							l.llCh <- struct{}{}
						}
					}
				}
			}
		}()
	}

	// abort and call lose leadership if there is something in the llCh channel
	go func() {
		l.b.logger.Printf("[DEBUG]: physical/mongo: llCh before")
		if l.llChOpen {
			<-l.llCh
			l.llChOpen = false
			close(l.llCh)
			l.llCh = nil
			l.b.logger.Printf("[DEBUG]: physical/mongo: llCh after")

			//l.loseLeadership()
			l.b.logger.Printf("[DEBUG]: physical/mongo: loseLeadership() begin")
			l.b.lockValue = ""
			l.b.isLeader = false
			abortCh1 <- struct{}{}
			abortCh2 <- struct{}{}
			close(l.b.leaderCh)
			l.b.leaderCh = nil
			l.b.logger.Printf("[DEBUG]: physical/mongo: loseLeadership() end")
		}
	}()

	return l.b.leaderCh, nil
}

// func (l *MongoLock) loseLeadership() {
// 	l.b.logger.Printf("[DEBUG]: physical/mongo: loseLeadership() begin")
// 	l.b.lockValue = ""
// 	l.b.isLeader = false
// 	l.abortCh1 <- struct{}{}
// 	l.abortCh2 <- struct{}{}
// 	close(l.b.leaderCh)
// 	l.b.leaderCh = nil
// 	l.b.logger.Printf("[DEBUG]: physical/mongo: loseLeadership() end")
// }

// Unlock is used to release the lock
func (l *MongoLock) Unlock() error {
	l.b.logger.Printf("[DEBUG]: physical/mongo: Unlock()")
	if l.value == "read" {
		l.b.logger.Printf("[DEBUG]: physical/mongo: Unlock() for mysterious 'read' lock!")
		return nil
	}

	// lose leadership, stop go routines and close channels
	if l.b.isLeader && l.llChOpen {
		l.llCh <- struct{}{}
	}

	session, err := l.b.activeSession()
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: could not establish mongo session: %v", err)
		return err
	}
	c := session.DB(l.b.Database).C(l.b.Collection)

	err = c.Remove(bson.M{"key": l.key, "value": l.value})
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: l.Unlock() could not remove lock: %v", err)
		return err
	}

	return nil
}

func (l *MongoLock) activePoller() error {
	l.b.logger.Printf("[DEBUG]: physical/mongo: activePoller()")
	session, err := l.b.activeSession()
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: could not establish mongo session: %v", err)
		return err
	}
	c := session.DB(l.b.Database).C(l.b.Collection)

	err = c.Update(bson.M{"key": l.key, "value": l.value}, bson.M{"$currentDate": bson.M{"lastCheckedIn": bson.M{"$type": "timestamp"}}})
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: activePoller(): Update() failed: %v", err)
		return err
	}

	return nil
}

/*func (l *MongoLock) standbyPoller() error {
	session, err := l.b.activeSession()
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: could not establish mongo session: %v", err)
		return err
	}
	c := session.DB(l.b.Database).C(l.b.Collection)

	// remove leader if timestamp too old
  time.Now().S
	return nil
}*/

type isMasterResult struct {
	IsMaster       bool
	Secondary      bool
	Primary        string
	Hosts          []string
	Passives       []string
	Tags           bson.D
	Msg            string
	SetName        string `bson:"setName"`
	MaxWireVersion int    `bson:"maxWireVersion"`
}

func (l *MongoLock) getPrimary() (string, error) {
	var result isMasterResult
	session, err := l.b.activeSession()
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: could not establish mongo session: %v", err)
		return "", err
	}
	err = session.Run("ismaster", &result)
	if err != nil {
		l.b.logger.Printf("[WARN]: physical/mongo: primaryPoller(): %v", err)
		return "", err
	}

	return result.Primary, nil
}

// Returns the value of the lock and if it is held
func (l *MongoLock) Value() (bool, string, error) {
	l.b.logger.Printf("[DEBUG]: physical/mongo: Value(%v)", l.value)
	session, err := l.b.activeSession()
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: could not establish mongo session: %v", err)
		return false, "", err
	}
	c := session.DB(l.b.Database).C(l.b.Collection)

	q := c.Find(bson.M{"key": l.key})
	n, err := q.Count()
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: l.Value(): error in Count() : %v", err)
		return false, "", err
	}

	if n <= 0 {
		//return false, "", errors.New("nobody holds the lock yet")
		//if l.value == "read" {
		//	return true, "", nil
		//} else {
		//	return l.b.isLeader, "", nil
		//}
		return false, "", nil
	}

	var entry LockEntry
	err = q.One(&entry)
	if err != nil {
		l.b.logger.Printf("[ERROR]: physical/mongo: l.Value(): error in One(): %v", err)
		return false, "", err
	}

	//l.b.logger.Printf("[DEBUG]: physical/mongo: LockEntry: %v", entry)
	//l.b.logger.Printf("[DEBUG]: physical/mongo: entry.Value '%v', l.Value '%v'", entry.Value, l.b.lockValue)

	// the read lock
	// if l.value == "read" {
	// 	l.b.logger.Printf("[ERROR]: physical/mongo: l.Value(read): %v", entry.Value)
	// 	return true, entry.Value, nil
	// }

	// // our default lock
	// l.b.logger.Printf("[ERROR]: physical/mongo: l.Value(%v): %v", l.value, entry.Value)
	// return l.b.isLeader, entry.Value, nil
	return true, entry.Value, nil
}
