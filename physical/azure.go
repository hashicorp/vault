package physical

import (
	"encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/storage"
	"io/ioutil"
	"sort"
	"strings"
)

// MaxBlobSize limits the value size per Blob to 4MB
var MaxBlobSize = 1024 * 1024 * 4

// AzureBackend is a backend storing data in Azure Blob Storage
type AzureBackend struct {
	container string
	client    storage.BlobStorageClient
}

func newAzureBackend(conf map[string]string) (Backend, error) {
	container, ok := conf["container"]

	if !ok {
		return nil, fmt.Errorf("Azure 'container' is required")
	}

	accountName, ok := conf["accountName"]

	if !ok {
		return nil, fmt.Errorf("An Azure 'accountName' is required")
	}

	accountKey, ok := conf["accountKey"]

	if !ok {
		return nil, fmt.Errorf("An Azure 'accountKey' is required")
	}

	client, err := storage.NewBasicClient(accountName, accountKey)

	if err != nil {
		return nil, fmt.Errorf("Failed to create Azure client: %s", err)
	}

	client.GetBlobService().CreateContainerIfNotExists(container, storage.ContainerAccessTypePrivate)

	backend := AzureBackend{container: container, client: client.GetBlobService()}
	return backend, nil
}

//Delete removes the blob {key} from the container
func (a AzureBackend) Delete(key string) error {
	_, err := a.client.DeleteBlobIfExists(a.container, key)
	return err
}

//Get returns the blob {key} from the container
func (a AzureBackend) Get(key string) (*Entry, error) {
	exists, _ := a.client.BlobExists(a.container, key)

	if !exists {
		return nil, nil
	}

	reader, err := a.client.GetBlob(a.container, key)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	ent := &Entry{
		Key:   key,
		Value: data,
	}

	return ent, nil
}

//Put updates the blob {key} in the container
func (a AzureBackend) Put(entry *Entry) error {
	if len(entry.Value) >= MaxBlobSize {
		return fmt.Errorf("Value is bigger than 4MB which is not supported at this time.")
	}

	// Create a 'dummy' blockID and a singleton array to commit the blob after upload
	blockID := base64.StdEncoding.EncodeToString([]byte("AAAA"))
	blocks := make([]storage.Block, 1)
	blocks[0] = storage.Block{ID: blockID, Status: storage.BlockStatusLatest}

	//Upload the data
	err := a.client.PutBlock(a.container, entry.Key, blockID, entry.Value)

	if err != nil {
		fmt.Println("Failed to upload blob", err)
		return err
	}

	//Commit the block written above
	err = a.client.PutBlockList(a.container, entry.Key, blocks)

	return err
}

//List returns all known blobs with {prefix}
func (a AzureBackend) List(prefix string) ([]string, error) {
	list, err := a.client.ListBlobs(a.container, storage.ListBlobsParameters{Prefix: prefix})

	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, blob := range list.Blobs {
		key := strings.TrimPrefix(blob.Name, prefix)

		if i := strings.Index(key, "/"); i == -1 {
			// Add objects only from the current 'folder'
			keys = append(keys, key)
		} else if i != -1 {
			// Add truncated 'folder' paths
			keys = appendIfMissing(keys, key[:i+1])
		}
	}
	sort.Strings(keys)
	return keys, nil
}
