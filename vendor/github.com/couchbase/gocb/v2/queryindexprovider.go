package gocb

import "time"

type queryIndexProvider interface {
	CreatePrimaryIndex(c *Collection, bucketName string, opts *CreatePrimaryQueryIndexOptions) error
	CreateIndex(c *Collection, bucketName, indexName string, fields []string, opts *CreateQueryIndexOptions) error
	DropPrimaryIndex(c *Collection, bucketName string, opts *DropPrimaryQueryIndexOptions) error
	DropIndex(c *Collection, bucketName, indexName string, opts *DropQueryIndexOptions) error
	GetAllIndexes(c *Collection, bucketName string, opts *GetAllQueryIndexesOptions) ([]QueryIndex, error)
	BuildDeferredIndexes(c *Collection, bucketName string, opts *BuildDeferredQueryIndexOptions) ([]string, error)
	WatchIndexes(c *Collection, bucketName string, watchList []string, timeout time.Duration, opts *WatchQueryIndexOptions) error
}
