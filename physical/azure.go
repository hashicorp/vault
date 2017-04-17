package physical

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/Azure/azure-storage-go"
	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
)

// MaxBlobSize at this time
var MaxBlobSize = 1024 * 1024 * 4

// AzureBackend is a physical backend that stores data
// within an Azure blob container.
type AzureBackend struct {
	container  string
	client     storage.BlobStorageClient
	logger     log.Logger
	permitPool *PermitPool
}

// newAzureBackend constructs an Azure backend using a pre-existing
// bucket. Credentials can be provided to the backend, sourced
// from the environment, AWS credential files or by IAM role.
func newAzureBackend(conf map[string]string, logger log.Logger) (Backend, error) {

	container := os.Getenv("AZURE_BLOB_CONTAINER")
	if container == "" {
		container = conf["container"]
		if container == "" {
			return nil, fmt.Errorf("'container' must be set")
		}
	}

	accountName := os.Getenv("AZURE_ACCOUNT_NAME")
	if accountName == "" {
		accountName = conf["accountName"]
		if accountName == "" {
			return nil, fmt.Errorf("'accountName' must be set")
		}
	}

	accountKey := os.Getenv("AZURE_ACCOUNT_KEY")
	if accountKey == "" {
		accountKey = conf["accountKey"]
		if accountKey == "" {
			return nil, fmt.Errorf("'accountKey' must be set")
		}
	}

	client, err := storage.NewBasicClient(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure client: %v", err)
	}

	contObj := client.GetBlobService().GetContainerReference(container)
	created, err := contObj.CreateIfNotExists()
	if err != nil {
		return nil, fmt.Errorf("failed to upsert container: %v", err)
	}
	if created {
		err = contObj.SetPermissions(storage.ContainerPermissions{
			AccessType: storage.ContainerAccessTypePrivate,
		}, 0, "")
		if err != nil {
			return nil, fmt.Errorf("failed to set permissions on newly-created container: %v", err)
		}
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
		if logger.IsDebug() {
			logger.Debug("azure: max_parallel set", "max_parallel", maxParInt)
		}
	}

	a := &AzureBackend{
		container:  container,
		client:     client.GetBlobService(),
		logger:     logger,
		permitPool: NewPermitPool(maxParInt),
	}
	return a, nil
}

// Put is used to insert or update an entry
func (a *AzureBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"azure", "put"}, time.Now())

	if len(entry.Value) >= MaxBlobSize {
		return fmt.Errorf("Value is bigger than the current supported limit of 4MBytes")
	}

	blockID := base64.StdEncoding.EncodeToString([]byte("AAAA"))
	blocks := make([]storage.Block, 1)
	blocks[0] = storage.Block{ID: blockID, Status: storage.BlockStatusLatest}

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	err := a.client.PutBlock(a.container, entry.Key, blockID, entry.Value)

	err = a.client.PutBlockList(a.container, entry.Key, blocks)
	return err
}

// Get is used to fetch an entry
func (a *AzureBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"azure", "get"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	exists, _ := a.client.BlobExists(a.container, key)

	if !exists {
		return nil, nil
	}

	reader, err := a.client.GetBlob(a.container, key)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(reader)

	ent := &Entry{
		Key:   key,
		Value: data,
	}

	return ent, err
}

// Delete is used to permanently delete an entry
func (a *AzureBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"azure", "delete"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	_, err := a.client.DeleteBlobIfExists(a.container, key, nil)
	return err
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (a *AzureBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"azure", "list"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	contObj := a.client.GetContainerReference(a.container)
	list, err := contObj.ListBlobs(storage.ListBlobsParameters{Prefix: prefix})

	if err != nil {
		// Break early.
		return nil, err
	}

	keys := []string{}
	for _, blob := range list.Blobs {
		key := strings.TrimPrefix(blob.Name, prefix)
		if i := strings.Index(key, "/"); i == -1 {
			keys = append(keys, key)
		} else {
			keys = appendIfMissing(keys, key[:i+1])
		}
	}

	sort.Strings(keys)
	return keys, nil
}
