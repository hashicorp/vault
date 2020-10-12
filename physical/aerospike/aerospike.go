package aerospike

import (
	"context"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	aero "github.com/aerospike/aerospike-client-go"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
)

const (
	keyBin   = "keyBin"
	valueBin = "valueBin"
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
	namespace := conf["namespace"]
	set := conf["set"]

	policy := aero.NewClientPolicy()
	policy.User = conf["username"]
	policy.Password = conf["password"]

	port, err := strconv.Atoi(conf["port"])
	if err != nil {
		return nil, err
	}

	client, err := aero.NewClientWithPolicy(policy, conf["hostname"], port)
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

func (a *AerospikeBackend) key(userKey string) (*aero.Key, error) {
	return aero.NewKey(a.namespace, a.set, hash(userKey))
}

// Put is used to insert or update an entry.
func (a *AerospikeBackend) Put(ctx context.Context, entry *physical.Entry) error {
	aeroKey, err := a.key(entry.Key)
	if err != nil {
		return err
	}

	writePolicy := aero.NewWritePolicy(0, 0)
	writePolicy.RecordExistsAction = aero.REPLACE

	binMap := make(aero.BinMap, 2)
	binMap[keyBin] = entry.Key
	binMap[valueBin] = entry.Value

	return a.client.Put(writePolicy, aeroKey, binMap)
}

// Get is used to fetch an entry.
func (a *AerospikeBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	aeroKey, err := a.key(key)
	if err != nil {
		return nil, err
	}

	record, err := a.client.Get(nil, aeroKey)
	if err != nil {
		if err.Error() == "Key not found" {
			return nil, nil
		}
		return nil, err
	}

	return &physical.Entry{
		Key:   key,
		Value: record.Bins[valueBin].([]byte),
	}, nil
}

// Delete is used to permanently delete an entry.
func (a *AerospikeBackend) Delete(ctx context.Context, key string) error {
	aeroKey, err := a.key(key)
	if err != nil {
		return err
	}

	_, err = a.client.Delete(nil, aeroKey)
	return err
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (a *AerospikeBackend) List(ctx context.Context, prefix string) ([]string, error) {
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
			trimPrefix := strings.Replace(recordKey, prefix, "", 1)
			keys := strings.Split(trimPrefix, "/")
			if len(keys) == 1 {
				keyList = append(keyList, keys[0])
			} else {
				withSlash := keys[0] + "/"
				if !listContains(keyList, withSlash) {
					keyList = append(keyList, withSlash)
				}
			}
		}
	}

	return keyList, nil
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprint(h.Sum32())
}

func listContains(list []string, s string) bool {
	for _, i := range list {
		if i == s {
			return true
		}
	}
	return false
}
