package physical

import (
	"encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/storage"
	"io/ioutil"
	"sort"
	"strings"
)

// AzureHABackend is a backend storing data in Azure Blob Storage
type AzureHABackend struct {
	container string
	client    storage.BlobStorageClient
	zkBackend *ZookeeperBackend
}

func newAzureHABackend(conf map[string]string) (Backend, error) {
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

	genBackend, zkerr := newZookeeperBackend(conf)

	if zkerr != nil {
		return nil, zkerr
	}

	zkBackend, ok := genBackend.(*ZookeeperBackend)

	backend := &AzureHABackend{
		container: container,
		client:    client.GetBlobService(),
		zkBackend: zkBackend,
	}
	return backend, nil
}

//Delete removes the blob {key} from the container
func (a *AzureHABackend) Delete(key string) error {
	_, err := a.client.DeleteBlobIfExists(a.container, key)
	return err
}

//Get returns the blob {key} from the container
func (a *AzureHABackend) Get(key string) (*Entry, error) {
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
func (a *AzureHABackend) Put(entry *Entry) error {
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
func (a *AzureHABackend) List(prefix string) ([]string, error) {
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

// LockWith is used for mutual exclusion based on the given key.
func (a *AzureHABackend) LockWith(key, value string) (Lock, error) {
	l := &ZookeeperHALock{
		in:    a.zkBackend,
		key:   key,
		value: value,
	}
	return l, nil
}
