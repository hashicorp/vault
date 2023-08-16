package aerospike

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"time"

	aero "github.com/aerospike/aerospike-client-go/v5"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/physical"
)

const (
	keyBin   = "keyBin"
	valueBin = "valueBin"

	defaultNamespace = "test"

	defaultHostname = "127.0.0.1"
	defaultPort     = 3000

	keyNotFoundError = "Key not found"
)

// AerospikeBackend is a physical backend that stores data in Aerospike.
type AerospikeBackend struct {
	client    *aero.Client
	namespace string
	set       string
	logger    log.Logger
}

// Verify AerospikeBackend satisfies the correct interface.
var _ physical.Backend = (*AerospikeBackend)(nil)

// NewAerospikeBackend constructs an AerospikeBackend backend.
func NewAerospikeBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	namespace, ok := conf["namespace"]
	if !ok {
		namespace = defaultNamespace
	}
	set := conf["set"]

	policy, err := buildClientPolicy(conf)
	if err != nil {
		return nil, err
	}

	client, err := buildAerospikeClient(conf, policy)
	if err != nil {
		return nil, err
	}

	return &AerospikeBackend{
		client:    client,
		namespace: namespace,
		set:       set,
		logger:    logger,
	}, nil
}

func buildAerospikeClient(conf map[string]string, policy *aero.ClientPolicy) (*aero.Client, error) {
	hostListString, ok := conf["hostlist"]
	if !ok || hostListString == "" {
		hostname, ok := conf["hostname"]
		if !ok || hostname == "" {
			hostname = defaultHostname
		}

		portString, ok := conf["port"]
		if !ok || portString == "" {
			portString = strconv.Itoa(defaultPort)
		}

		port, err := strconv.Atoi(portString)
		if err != nil {
			return nil, err
		}

		return aero.NewClientWithPolicy(policy, hostname, port)
	}

	hostList, err := parseHostList(hostListString)
	if err != nil {
		return nil, err
	}

	return aero.NewClientWithPolicyAndHost(policy, hostList...)
}

func buildClientPolicy(conf map[string]string) (*aero.ClientPolicy, error) {
	policy := aero.NewClientPolicy()

	policy.User = conf["username"]
	policy.Password = conf["password"]

	authMode := aero.AuthModeInternal
	if mode, ok := conf["auth_mode"]; ok {
		switch strings.ToUpper(mode) {
		case "EXTERNAL":
			authMode = aero.AuthModeExternal
		case "INTERNAL":
			authMode = aero.AuthModeInternal
		default:
			return nil, fmt.Errorf("'auth_mode' must be one of {INTERNAL, EXTERNAL}")
		}
	}
	policy.AuthMode = authMode
	policy.ClusterName = conf["cluster_name"]

	if timeoutString, ok := conf["timeout"]; ok {
		timeout, err := strconv.Atoi(timeoutString)
		if err != nil {
			return nil, err
		}
		policy.Timeout = time.Duration(timeout) * time.Millisecond
	}

	if idleTimeoutString, ok := conf["idle_timeout"]; ok {
		idleTimeout, err := strconv.Atoi(idleTimeoutString)
		if err != nil {
			return nil, err
		}
		policy.IdleTimeout = time.Duration(idleTimeout) * time.Millisecond
	}

	return policy, nil
}

func (a *AerospikeBackend) key(userKey string) (*aero.Key, error) {
	return aero.NewKey(a.namespace, a.set, hash(userKey))
}

// Put is used to insert or update an entry.
func (a *AerospikeBackend) Put(_ context.Context, entry *physical.Entry) error {
	aeroKey, err := a.key(entry.Key)
	if err != nil {
		return err
	}

	// replace the Aerospike record if exists
	writePolicy := aero.NewWritePolicy(0, 0)
	writePolicy.RecordExistsAction = aero.REPLACE

	binMap := make(aero.BinMap, 2)
	binMap[keyBin] = entry.Key
	binMap[valueBin] = entry.Value

	return a.client.Put(writePolicy, aeroKey, binMap)
}

// Get is used to fetch an entry.
func (a *AerospikeBackend) Get(_ context.Context, key string) (*physical.Entry, error) {
	aeroKey, err := a.key(key)
	if err != nil {
		return nil, err
	}

	record, err := a.client.Get(nil, aeroKey)
	if err != nil {
		if strings.Contains(err.Error(), keyNotFoundError) {
			return nil, nil
		}
		return nil, err
	}

	value, ok := record.Bins[valueBin]
	if !ok {
		return nil, fmt.Errorf("Value bin was not found in the record")
	}

	return &physical.Entry{
		Key:   key,
		Value: value.([]byte),
	}, nil
}

// Delete is used to permanently delete an entry.
func (a *AerospikeBackend) Delete(_ context.Context, key string) error {
	aeroKey, err := a.key(key)
	if err != nil {
		return err
	}

	_, err = a.client.Delete(nil, aeroKey)
	return err
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (a *AerospikeBackend) List(_ context.Context, prefix string) ([]string, error) {
	recordSet, err := a.client.ScanAll(nil, a.namespace, a.set)
	if err != nil {
		return nil, err
	}

	var keyList []string
	for res := range recordSet.Results() {
		if res.Err != nil {
			return nil, res.Err
		}
		recordKey := res.Record.Bins[keyBin].(string)
		if strings.HasPrefix(recordKey, prefix) {
			trimPrefix := strings.TrimPrefix(recordKey, prefix)
			keys := strings.Split(trimPrefix, "/")
			if len(keys) == 1 {
				keyList = append(keyList, keys[0])
			} else {
				withSlash := keys[0] + "/"
				if !strutil.StrListContains(keyList, withSlash) {
					keyList = append(keyList, withSlash)
				}
			}
		}
	}

	return keyList, nil
}

func parseHostList(list string) ([]*aero.Host, error) {
	hosts := strings.Split(list, ",")
	var hostList []*aero.Host
	for _, host := range hosts {
		if host == "" {
			continue
		}
		split := strings.Split(host, ":")
		switch len(split) {
		case 1:
			hostList = append(hostList, aero.NewHost(split[0], defaultPort))
		case 2:
			port, err := strconv.Atoi(split[1])
			if err != nil {
				return nil, err
			}
			hostList = append(hostList, aero.NewHost(split[0], port))
		default:
			return nil, fmt.Errorf("Invalid 'hostlist' configuration")
		}
	}
	return hostList, nil
}

func hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash[:])
}
