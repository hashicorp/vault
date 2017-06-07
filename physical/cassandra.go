package physical

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/armon/go-metrics"
	"github.com/gocql/gocql"
)

// CassandraBackend is a physical backend that stores data in Cassandra.
type CassandraBackend struct {
	sess  *gocql.Session
	table string

	logger log.Logger
}

// newCassandraBackend constructs a Cassandra backend using a pre-existing
// keyspace and table.
func newCassandraBackend(conf map[string]string, logger log.Logger) (Backend, error) {
	splitArray := func(v string) []string {
		return strings.FieldsFunc(v, func(r rune) bool {
			return r == ','
		})
	}

	var (
		hosts       = splitArray(conf["hosts"])
		keyspace    = conf["keyspace"]
		table       = conf["table"]
		consistency = gocql.LocalQuorum
	)

	if len(hosts) == 0 {
		hosts = []string{"localhost"}
	}
	if keyspace == "" {
		keyspace = "vault"
	}
	if table == "" {
		table = "entries"
	}
	if cs, ok := conf["consistency"]; ok {
		switch cs {
		case "ANY":
			consistency = gocql.Any
		case "ONE":
			consistency = gocql.One
		case "TWO":
			consistency = gocql.Two
		case "THREE":
			consistency = gocql.Three
		case "QUORUM":
			consistency = gocql.Quorum
		case "ALL":
			consistency = gocql.All
		case "LOCAL_QUORUM":
			consistency = gocql.LocalQuorum
		case "EACH_QUORUM":
			consistency = gocql.EachQuorum
		case "LOCAL_ONE":
			consistency = gocql.LocalOne
		default:
			return nil, fmt.Errorf("'consistency' must be one of {ANY, ONE, TWO, THREE, QUORUM, ALL, LOCAL_QUORUM, EACH_QUORUM, LOCAL_ONE}")
		}
	}

	connectStart := time.Now()
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace

	cluster.ProtoVersion = 2
	if protoVersionStr, ok := conf["protocol_version"]; ok {
		protoVersion, err := strconv.Atoi(protoVersionStr)
		if err != nil {
			return nil, fmt.Errorf("'protocol_version' must be an integer")
		}
		cluster.ProtoVersion = protoVersion
	}

	if username, ok := conf["username"]; ok {
		if cluster.ProtoVersion < 2 {
			return nil, fmt.Errorf("Authentication is not supported with protocol version < 2")
		}
		authenticator := gocql.PasswordAuthenticator{Username: username}
		if password, ok := conf["password"]; ok {
			authenticator.Password = password
		}
		cluster.Authenticator = authenticator
	}

	if connTimeoutStr, ok := conf["connection_timeout"]; ok {
		connectionTimeout, err := strconv.Atoi(connTimeoutStr)
		if err != nil {
			return nil, fmt.Errorf("'connection_timeout' must be an integer")
		}
		cluster.Timeout = time.Duration(connectionTimeout) * time.Millisecond
	}

	sess, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	metrics.MeasureSince([]string{"cassandra", "connect"}, connectStart)
	sess.SetConsistency(consistency)

	impl := &CassandraBackend{
		sess:   sess,
		table:  table,
		logger: logger}
	return impl, nil
}

// bucketName sanitises a bucket name for Cassandra
func (c *CassandraBackend) bucketName(name string) string {
	if name == "" {
		name = "."
	}
	return strings.TrimRight(name, "/")
}

// bucket returns all the prefix buckets the key should be stored at
func (c *CassandraBackend) buckets(key string) []string {
	vals := append([]string{""}, prefixes(key)...)
	for i, v := range vals {
		vals[i] = c.bucketName(v)
	}
	return vals
}

// bucket returns the most specific bucket for the key
func (c *CassandraBackend) bucket(key string) string {
	bs := c.buckets(key)
	return bs[len(bs)-1]
}

// Put is used to insert or update an entry
func (c *CassandraBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"cassandra", "put"}, time.Now())

	stmt := fmt.Sprintf(`INSERT INTO "%s" (bucket, key, value) VALUES (?, ?, ?)`, c.table)
	batch := gocql.NewBatch(gocql.LoggedBatch)
	for _, bucket := range c.buckets(entry.Key) {
		batch.Entries = append(batch.Entries, gocql.BatchEntry{
			Stmt: stmt,
			Args: []interface{}{bucket, entry.Key, entry.Value}})
	}
	return c.sess.ExecuteBatch(batch)
}

// Get is used to fetch an entry
func (c *CassandraBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"cassandra", "get"}, time.Now())

	v := []byte(nil)
	stmt := fmt.Sprintf(`SELECT value FROM "%s" WHERE bucket = ? AND key = ? LIMIT 1`, c.table)
	q := c.sess.Query(stmt, c.bucket(key), key)
	if err := q.Scan(&v); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &Entry{
		Key:   key,
		Value: v,
	}, nil
}

// Delete is used to permanently delete an entry
func (c *CassandraBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"cassandra", "delete"}, time.Now())

	stmt := fmt.Sprintf(`DELETE FROM "%s" WHERE bucket = ? AND key = ?`, c.table)
	batch := gocql.NewBatch(gocql.LoggedBatch)
	for _, bucket := range c.buckets(key) {
		batch.Entries = append(batch.Entries, gocql.BatchEntry{
			Stmt: stmt,
			Args: []interface{}{bucket, key}})
	}
	return c.sess.ExecuteBatch(batch)
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (c *CassandraBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"cassandra", "list"}, time.Now())

	stmt := fmt.Sprintf(`SELECT key FROM "%s" WHERE bucket = ?`, c.table)
	q := c.sess.Query(stmt, c.bucketName(prefix))
	iter := q.Iter()
	k, keys := "", []string{}
	for iter.Scan(&k) {
		// Only return the next "component" (with a trailing slash if it has children)
		k = strings.TrimPrefix(k, prefix)
		if parts := strings.SplitN(k, "/", 2); len(parts) > 1 {
			k = parts[0] + "/"
		} else {
			k = parts[0]
		}

		// Deduplicate; this works because the keys are sorted
		if len(keys) > 0 && keys[len(keys)-1] == k {
			continue
		}
		keys = append(keys, k)
	}
	return keys, iter.Close()
}
