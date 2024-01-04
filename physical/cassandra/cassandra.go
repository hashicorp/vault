// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cassandra

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/gocql/gocql"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/physical"
)

// CassandraBackend is a physical backend that stores data in Cassandra.
type CassandraBackend struct {
	sess  *gocql.Session
	table string

	logger log.Logger
}

// Verify CassandraBackend satisfies the correct interfaces
var _ physical.Backend = (*CassandraBackend)(nil)

// NewCassandraBackend constructs a Cassandra backend using a pre-existing
// keyspace and table.
func NewCassandraBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	splitArray := func(v string) []string {
		return strings.FieldsFunc(v, func(r rune) bool {
			return r == ','
		})
	}

	var (
		hosts        = splitArray(conf["hosts"])
		port         = 9042
		explicitPort = false
		keyspace     = conf["keyspace"]
		table        = conf["table"]
		consistency  = gocql.LocalQuorum
	)

	if len(hosts) == 0 {
		hosts = []string{"localhost"}
	}
	for i, hp := range hosts {
		h, ps, err := net.SplitHostPort(hp)
		if err != nil {
			continue
		}
		p, err := strconv.Atoi(ps)
		if err != nil {
			return nil, err
		}

		if explicitPort && p != port {
			return nil, fmt.Errorf("all hosts must have the same port")
		}
		hosts[i], port = h, p
		explicitPort = true
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
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.Consistency = consistency

	if retryCountStr, ok := conf["simple_retry_policy_retries"]; ok {
		retryCount, err := strconv.Atoi(retryCountStr)
		if err != nil || retryCount <= 0 {
			return nil, fmt.Errorf("'simple_retry_policy_retries' must be a positive integer")
		}
		cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: retryCount}
	}

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
			return nil, fmt.Errorf("authentication is not supported with protocol version < 2")
		}
		authenticator := gocql.PasswordAuthenticator{Username: username}
		if password, ok := conf["password"]; ok {
			authenticator.Password = password
		}
		cluster.Authenticator = authenticator
	}

	if initialConnectionTimeoutStr, ok := conf["initial_connection_timeout"]; ok {
		initialConnectionTimeout, err := strconv.Atoi(initialConnectionTimeoutStr)
		if err != nil || initialConnectionTimeout <= 0 {
			return nil, fmt.Errorf("'initial_connection_timeout' must be a positive integer")
		}
		cluster.ConnectTimeout = time.Duration(initialConnectionTimeout) * time.Second
	}

	if connTimeoutStr, ok := conf["connection_timeout"]; ok {
		connectionTimeout, err := strconv.Atoi(connTimeoutStr)
		if err != nil || connectionTimeout <= 0 {
			return nil, fmt.Errorf("'connection_timeout' must be a positive integer")
		}
		cluster.Timeout = time.Duration(connectionTimeout) * time.Second
	}

	if err := setupCassandraTLS(conf, cluster); err != nil {
		return nil, err
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
		logger: logger,
	}
	return impl, nil
}

func setupCassandraTLS(conf map[string]string, cluster *gocql.ClusterConfig) error {
	tlsOnStr, ok := conf["tls"]
	if !ok {
		return nil
	}

	tlsOn, err := strconv.Atoi(tlsOnStr)
	if err != nil {
		return fmt.Errorf("'tls' must be an integer (0 or 1)")
	}

	if tlsOn == 0 {
		return nil
	}

	tlsConfig := &tls.Config{}
	if pemBundlePath, ok := conf["pem_bundle_file"]; ok {
		pemBundleData, err := ioutil.ReadFile(pemBundlePath)
		if err != nil {
			return fmt.Errorf("error reading pem bundle from %q: %w", pemBundlePath, err)
		}
		pemBundle, err := certutil.ParsePEMBundle(string(pemBundleData))
		if err != nil {
			return fmt.Errorf("error parsing 'pem_bundle': %w", err)
		}
		tlsConfig, err = pemBundle.GetTLSConfig(certutil.TLSClient)
		if err != nil {
			return err
		}
	} else if pemJSONPath, ok := conf["pem_json_file"]; ok {
		pemJSONData, err := ioutil.ReadFile(pemJSONPath)
		if err != nil {
			return fmt.Errorf("error reading json bundle from %q: %w", pemJSONPath, err)
		}
		pemJSON, err := certutil.ParsePKIJSON([]byte(pemJSONData))
		if err != nil {
			return err
		}
		tlsConfig, err = pemJSON.GetTLSConfig(certutil.TLSClient)
		if err != nil {
			return err
		}
	}

	if tlsSkipVerifyStr, ok := conf["tls_skip_verify"]; ok {
		tlsSkipVerify, err := strconv.Atoi(tlsSkipVerifyStr)
		if err != nil {
			return fmt.Errorf("'tls_skip_verify' must be an integer (0 or 1)")
		}
		if tlsSkipVerify == 0 {
			tlsConfig.InsecureSkipVerify = false
		} else {
			tlsConfig.InsecureSkipVerify = true
		}
	}

	if tlsMinVersion, ok := conf["tls_min_version"]; ok {
		switch tlsMinVersion {
		case "tls10":
			tlsConfig.MinVersion = tls.VersionTLS10
		case "tls11":
			tlsConfig.MinVersion = tls.VersionTLS11
		case "tls12":
			tlsConfig.MinVersion = tls.VersionTLS12
		case "tls13":
			tlsConfig.MinVersion = tls.VersionTLS13
		default:
			return fmt.Errorf("'tls_min_version' must be one of `tls10`, `tls11`, `tls12` or `tls13`")
		}
	}

	cluster.SslOpts = &gocql.SslOptions{
		Config:                 tlsConfig,
		EnableHostVerification: !tlsConfig.InsecureSkipVerify,
	}
	return nil
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
	vals := append([]string{""}, physical.Prefixes(key)...)
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
func (c *CassandraBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"cassandra", "put"}, time.Now())

	// Execute inserts to each key prefix simultaneously
	stmt := fmt.Sprintf(`INSERT INTO "%s" (bucket, key, value) VALUES (?, ?, ?)`, c.table)
	buckets := c.buckets(entry.Key)
	results := make(chan error, len(buckets))
	for i, _bucket := range buckets {
		go func(i int, bucket string) {
			var value []byte
			if i == len(buckets)-1 {
				// Only store the full value if this is the leaf bucket where the entry will actually be read
				// otherwise this write is just to allow for list operations
				value = entry.Value
			}
			results <- c.sess.Query(stmt, bucket, entry.Key, value).Exec()
		}(i, _bucket)
	}
	for i := 0; i < len(buckets); i++ {
		if err := <-results; err != nil {
			return err
		}
	}
	return nil
}

// Get is used to fetch an entry
func (c *CassandraBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
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

	return &physical.Entry{
		Key:   key,
		Value: v,
	}, nil
}

// Delete is used to permanently delete an entry
func (c *CassandraBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"cassandra", "delete"}, time.Now())

	stmt := fmt.Sprintf(`DELETE FROM "%s" WHERE bucket = ? AND key = ?`, c.table)
	buckets := c.buckets(key)
	results := make(chan error, len(buckets))

	for _, bucket := range buckets {
		go func(bucket string) {
			results <- c.sess.Query(stmt, bucket, key).Exec()
		}(bucket)
	}

	for i := 0; i < len(buckets); i++ {
		if err := <-results; err != nil {
			return err
		}
	}
	return nil
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (c *CassandraBackend) List(ctx context.Context, prefix string) ([]string, error) {
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
