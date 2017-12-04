package azure

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	storage "github.com/Azure/azure-sdk-for-go/storage"
	log "github.com/mgutz/logxi/v1"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/physical"
)

// MaxBlobSize at this time
var MaxBlobSize = 1024 * 1024 * 4

// AzureBackend is a physical backend that stores data
// within an Azure blob container.
type AzureBackend struct {
	container  *storage.Container
	logger     log.Logger
	permitPool *physical.PermitPool
}

// NewAzureBackend constructs an Azure backend using a pre-existing
// bucket. Credentials can be provided to the backend, sourced
// from the environment, AWS credential files or by IAM role.
func NewAzureBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	name := os.Getenv("AZURE_BLOB_CONTAINER")
	if name == "" {
		name = conf["container"]
		if name == "" {
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
	client.HTTPClient = cleanhttp.DefaultPooledClient()

	blobClient := client.GetBlobService()
	container := blobClient.GetContainerReference(name)
	_, err = container.CreateIfNotExists(&storage.CreateContainerOptions{
		Access: storage.ContainerAccessTypePrivate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create %q container: %v", name, err)
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
		logger:     logger,
		permitPool: physical.NewPermitPool(maxParInt),
	}
	return a, nil
}

// Put is used to insert or update an entry
func (a *AzureBackend) Put(entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"azure", "put"}, time.Now())

	if len(entry.Value) >= MaxBlobSize {
		return fmt.Errorf("value is bigger than the current supported limit of 4MBytes")
	}

	blockID := base64.StdEncoding.EncodeToString([]byte("AAAA"))
	blocks := make([]storage.Block, 1)
	blocks[0] = storage.Block{ID: blockID, Status: storage.BlockStatusLatest}

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	blob := &storage.Blob{
		Container: a.container,
		Name:      entry.Key,
	}
	if err := blob.PutBlock(blockID, entry.Value, nil); err != nil {
		return err
	}

	return blob.PutBlockList(blocks, nil)
}

// Get is used to fetch an entry
func (a *AzureBackend) Get(key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"azure", "get"}, time.Now())

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	blob := &storage.Blob{
		Container: a.container,
		Name:      key,
	}
	exists, err := blob.Exists()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	reader, err := blob.Get(nil)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)

	ent := &physical.Entry{
		Key:   key,
		Value: data,
	}

	return ent, err
}

// Delete is used to permanently delete an entry
func (a *AzureBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"azure", "delete"}, time.Now())

	blob := &storage.Blob{
		Container: a.container,
		Name:      key,
	}

	a.permitPool.Acquire()
	defer a.permitPool.Release()

	_, err := blob.DeleteIfExists(nil)
	return err
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (a *AzureBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"azure", "list"}, time.Now())

	a.permitPool.Acquire()
	list, err := a.container.ListBlobs(storage.ListBlobsParameters{Prefix: prefix})
	if err != nil {
		// Break early.
		a.permitPool.Release()
		return nil, err
	}
	a.permitPool.Release()

	keys := []string{}
	for _, blob := range list.Blobs {
		key := strings.TrimPrefix(blob.Name, prefix)
		if i := strings.Index(key, "/"); i == -1 {
			keys = append(keys, key)
		} else {
			keys = strutil.AppendIfMissing(keys, key[:i+1])
		}
	}

	sort.Strings(keys)
	return keys, nil
}
