package plugin

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/pb"
)

func newGRPCStorageClient(conn *grpc.ClientConn) *GRPCStorageClient {
	return &GRPCStorageClient{
		client: pb.NewStorageClient(conn),
	}
}

// GRPCStorageClient is an implementation of logical.Storage that communicates
// over RPC.
type GRPCStorageClient struct {
	client pb.StorageClient
}

func (s *GRPCStorageClient) List(ctx context.Context, prefix string) ([]string, error) {
	reply, err := s.client.List(ctx, &pb.StorageListArgs{
		Prefix: prefix,
	}, largeMsgGRPCCallOpts...)
	if err != nil {
		return []string{}, err
	}
	if reply.Err != "" {
		return reply.Keys, errors.New(reply.Err)
	}
	return reply.Keys, nil
}

func (s *GRPCStorageClient) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	reply, err := s.client.Get(ctx, &pb.StorageGetArgs{
		Key: key,
	}, largeMsgGRPCCallOpts...)
	if err != nil {
		return nil, err
	}
	if reply.Err != "" {
		return nil, errors.New(reply.Err)
	}
	return pb.ProtoStorageEntryToLogicalStorageEntry(reply.Entry), nil
}

func (s *GRPCStorageClient) Put(ctx context.Context, entry *logical.StorageEntry) error {
	reply, err := s.client.Put(ctx, &pb.StoragePutArgs{
		Entry: pb.LogicalStorageEntryToProtoStorageEntry(entry),
	}, largeMsgGRPCCallOpts...)
	if err != nil {
		return err
	}
	if reply.Err != "" {
		return errors.New(reply.Err)
	}
	return nil
}

func (s *GRPCStorageClient) Delete(ctx context.Context, key string) error {
	reply, err := s.client.Delete(ctx, &pb.StorageDeleteArgs{
		Key: key,
	})
	if err != nil {
		return err
	}
	if reply.Err != "" {
		return errors.New(reply.Err)
	}
	return nil
}

// StorageServer is a net/rpc compatible structure for serving
type GRPCStorageServer struct {
	impl logical.Storage
}

func (s *GRPCStorageServer) List(ctx context.Context, args *pb.StorageListArgs) (*pb.StorageListReply, error) {
	keys, err := s.impl.List(ctx, args.Prefix)
	return &pb.StorageListReply{
		Keys: keys,
		Err:  pb.ErrToString(err),
	}, nil
}

func (s *GRPCStorageServer) Get(ctx context.Context, args *pb.StorageGetArgs) (*pb.StorageGetReply, error) {
	storageEntry, err := s.impl.Get(ctx, args.Key)
	return &pb.StorageGetReply{
		Entry: pb.LogicalStorageEntryToProtoStorageEntry(storageEntry),
		Err:   pb.ErrToString(err),
	}, nil
}

func (s *GRPCStorageServer) Put(ctx context.Context, args *pb.StoragePutArgs) (*pb.StoragePutReply, error) {
	err := s.impl.Put(ctx, pb.ProtoStorageEntryToLogicalStorageEntry(args.Entry))
	return &pb.StoragePutReply{
		Err: pb.ErrToString(err),
	}, nil
}

func (s *GRPCStorageServer) Delete(ctx context.Context, args *pb.StorageDeleteArgs) (*pb.StorageDeleteReply, error) {
	err := s.impl.Delete(ctx, args.Key)
	return &pb.StorageDeleteReply{
		Err: pb.ErrToString(err),
	}, nil
}

// NOOPStorage is used to deny access to the storage interface while running a
// backend plugin in metadata mode.
type NOOPStorage struct{}

func (s *NOOPStorage) List(_ context.Context, prefix string) ([]string, error) {
	return []string{}, nil
}

func (s *NOOPStorage) Get(_ context.Context, key string) (*logical.StorageEntry, error) {
	return nil, nil
}

func (s *NOOPStorage) Put(_ context.Context, entry *logical.StorageEntry) error {
	return nil
}

func (s *NOOPStorage) Delete(_ context.Context, key string) error {
	return nil
}
