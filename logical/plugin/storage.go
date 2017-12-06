package plugin

import (
	"net/rpc"

	"github.com/hashicorp/vault/logical"
)

// StorageClient is an implementation of logical.Storage that communicates
// over RPC.
type StorageClient struct {
	client *rpc.Client
}

func (s *StorageClient) List(prefix string) ([]string, error) {
	var reply StorageListReply
	err := s.client.Call("Plugin.List", prefix, &reply)
	if err != nil {
		return reply.Keys, err
	}
	if reply.Error != nil {
		return reply.Keys, reply.Error
	}
	return reply.Keys, nil
}

func (s *StorageClient) Get(key string) (*logical.StorageEntry, error) {
	var reply StorageGetReply
	err := s.client.Call("Plugin.Get", key, &reply)
	if err != nil {
		return nil, err
	}
	if reply.Error != nil {
		return nil, reply.Error
	}
	return reply.StorageEntry, nil
}

func (s *StorageClient) Put(entry *logical.StorageEntry) error {
	var reply StoragePutReply
	err := s.client.Call("Plugin.Put", entry, &reply)
	if err != nil {
		return err
	}
	if reply.Error != nil {
		return reply.Error
	}
	return nil
}

func (s *StorageClient) Delete(key string) error {
	var reply StorageDeleteReply
	err := s.client.Call("Plugin.Delete", key, &reply)
	if err != nil {
		return err
	}
	if reply.Error != nil {
		return reply.Error
	}
	return nil
}

// StorageServer is a net/rpc compatible structure for serving
type StorageServer struct {
	impl logical.Storage
}

func (s *StorageServer) List(prefix string, reply *StorageListReply) error {
	keys, err := s.impl.List(prefix)
	*reply = StorageListReply{
		Keys:  keys,
		Error: wrapError(err),
	}
	return nil
}

func (s *StorageServer) Get(key string, reply *StorageGetReply) error {
	storageEntry, err := s.impl.Get(key)
	*reply = StorageGetReply{
		StorageEntry: storageEntry,
		Error:        wrapError(err),
	}
	return nil
}

func (s *StorageServer) Put(entry *logical.StorageEntry, reply *StoragePutReply) error {
	err := s.impl.Put(entry)
	*reply = StoragePutReply{
		Error: wrapError(err),
	}
	return nil
}

func (s *StorageServer) Delete(key string, reply *StorageDeleteReply) error {
	err := s.impl.Delete(key)
	*reply = StorageDeleteReply{
		Error: wrapError(err),
	}
	return nil
}

type StorageListReply struct {
	Keys  []string
	Error error
}

type StorageGetReply struct {
	StorageEntry *logical.StorageEntry
	Error        error
}

type StoragePutReply struct {
	Error error
}

type StorageDeleteReply struct {
	Error error
}

// NOOPStorage is used to deny access to the storage interface while running a
// backend plugin in metadata mode.
type NOOPStorage struct{}

func (s *NOOPStorage) List(prefix string) ([]string, error) {
	return []string{}, nil
}

func (s *NOOPStorage) Get(key string) (*logical.StorageEntry, error) {
	return nil, nil
}

func (s *NOOPStorage) Put(entry *logical.StorageEntry) error {
	return nil
}

func (s *NOOPStorage) Delete(key string) error {
	return nil
}
