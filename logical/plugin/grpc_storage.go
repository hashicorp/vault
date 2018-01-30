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
	})
	if err != nil {
		return reply.Keys, err
	}
	if reply.Err != "" {
		return reply.Keys, errors.New(reply.Err)
	}
	return reply.Keys, nil
}

func (s *GRPCStorageClient) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	reply, err := s.client.Get(ctx, &pb.StorageGetArgs{
		Key: key,
	})
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
	})
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
