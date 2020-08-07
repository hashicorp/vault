// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// DefaultName is the default name for a GridFS bucket.
var DefaultName = "fs"

// DefaultChunkSize is the default size of each file chunk in bytes.
var DefaultChunkSize int32 = 255 * 1024 // 255 KiB

// DefaultRevision is the default revision number for a download by name operation.
var DefaultRevision int32 = -1

// BucketOptions represents all possible options to configure a GridFS bucket.
type BucketOptions struct {
	Name           *string                    // The bucket name. Defaults to "fs".
	ChunkSizeBytes *int32                     // The chunk size in bytes. Defaults to 255KB.
	WriteConcern   *writeconcern.WriteConcern // The write concern for the bucket. Defaults to the write concern of the database.
	ReadConcern    *readconcern.ReadConcern   // The read concern for the bucket. Defaults to the read concern of the database.
	ReadPreference *readpref.ReadPref         // The read preference for the bucket. Defaults to the read preference of the database.
}

// GridFSBucket creates a new *BucketOptions
func GridFSBucket() *BucketOptions {
	return &BucketOptions{
		Name:           &DefaultName,
		ChunkSizeBytes: &DefaultChunkSize,
	}
}

// SetName sets the name for the bucket. Defaults to "fs" if not set.
func (b *BucketOptions) SetName(name string) *BucketOptions {
	b.Name = &name
	return b
}

// SetChunkSizeBytes sets the chunk size in bytes for the bucket. Defaults to 255KB if not set.
func (b *BucketOptions) SetChunkSizeBytes(i int32) *BucketOptions {
	b.ChunkSizeBytes = &i
	return b
}

// SetWriteConcern sets the write concern for the bucket.
func (b *BucketOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *BucketOptions {
	b.WriteConcern = wc
	return b
}

// SetReadConcern sets the read concern for the bucket.
func (b *BucketOptions) SetReadConcern(rc *readconcern.ReadConcern) *BucketOptions {
	b.ReadConcern = rc
	return b
}

// SetReadPreference sets the read preference for the bucket.
func (b *BucketOptions) SetReadPreference(rp *readpref.ReadPref) *BucketOptions {
	b.ReadPreference = rp
	return b
}

// MergeBucketOptions combines the given *BucketOptions into a single *BucketOptions.
// If the name or chunk size is not set in any of the given *BucketOptions, the resulting *BucketOptions will have
// name "fs" and chunk size 255KB.
func MergeBucketOptions(opts ...*BucketOptions) *BucketOptions {
	b := GridFSBucket()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Name != nil {
			b.Name = opt.Name
		}
		if opt.ChunkSizeBytes != nil {
			b.ChunkSizeBytes = opt.ChunkSizeBytes
		}
		if opt.WriteConcern != nil {
			b.WriteConcern = opt.WriteConcern
		}
		if opt.ReadConcern != nil {
			b.ReadConcern = opt.ReadConcern
		}
		if opt.ReadPreference != nil {
			b.ReadPreference = opt.ReadPreference
		}
	}

	return b
}

// UploadOptions represents all possible options for a GridFS upload operation.  If a registry is nil, bson.DefaultRegistry
// will be used when converting the Metadata interface to BSON.
type UploadOptions struct {
	ChunkSizeBytes *int32              // Chunk size in bytes. Defaults to the chunk size of the bucket.
	Metadata       interface{}         // User data for the 'metadata' field of the files collection document.
	Registry       *bsoncodec.Registry // The registry to use for converting filters. Defaults to bson.DefaultRegistry.
}

// GridFSUpload creates a new *UploadOptions
func GridFSUpload() *UploadOptions {
	return &UploadOptions{Registry: bson.DefaultRegistry}
}

// SetChunkSizeBytes sets the chunk size in bytes for the upload. Defaults to 255KB if not set.
func (u *UploadOptions) SetChunkSizeBytes(i int32) *UploadOptions {
	u.ChunkSizeBytes = &i
	return u
}

// SetMetadata specfies the metadata for the upload.
func (u *UploadOptions) SetMetadata(doc interface{}) *UploadOptions {
	u.Metadata = doc
	return u
}

// MergeUploadOptions combines the given *UploadOptions into a single *UploadOptions.
// If the chunk size is not set in any of the given *UploadOptions, the resulting *UploadOptions will have chunk size
// 255KB.
func MergeUploadOptions(opts ...*UploadOptions) *UploadOptions {
	u := GridFSUpload()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ChunkSizeBytes != nil {
			u.ChunkSizeBytes = opt.ChunkSizeBytes
		}
		if opt.Metadata != nil {
			u.Metadata = opt.Metadata
		}
		if opt.Registry != nil {
			u.Registry = opt.Registry
		}
	}

	return u
}

// NameOptions represents all options that can be used for a GridFS download by name operation.
type NameOptions struct {
	Revision *int32 // Which revision (documents with the same filename and different uploadDate). Defaults to -1 (the most recent revision).
}

// GridFSName creates a new *NameOptions
func GridFSName() *NameOptions {
	return &NameOptions{}
}

// SetRevision specifies which revision of the file to retrieve. Defaults to -1.
// * Revision numbers are defined as follows:
// * 0 = the original stored file
// * 1 = the first revision
// * 2 = the second revision
// * etc…
// * -2 = the second most recent revision
// * -1 = the most recent revision
func (n *NameOptions) SetRevision(r int32) *NameOptions {
	n.Revision = &r
	return n
}

// MergeNameOptions combines the given *NameOptions into a single *NameOptions in a last one wins fashion.
func MergeNameOptions(opts ...*NameOptions) *NameOptions {
	n := GridFSName()
	n.Revision = &DefaultRevision

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Revision != nil {
			n.Revision = opt.Revision
		}
	}

	return n
}

// GridFSFindOptions represents all options for a GridFS find operation.
type GridFSFindOptions struct {
	BatchSize       *int32
	Limit           *int32
	MaxTime         *time.Duration
	NoCursorTimeout *bool
	Skip            *int32
	Sort            interface{}
}

// GridFSFind creates a new GridFSFindOptions instance.
func GridFSFind() *GridFSFindOptions {
	return &GridFSFindOptions{}
}

// SetBatchSize sets the number of documents to return in each batch.
func (f *GridFSFindOptions) SetBatchSize(i int32) *GridFSFindOptions {
	f.BatchSize = &i
	return f
}

// SetLimit specifies a limit on the number of results.
// A negative limit implies that only 1 batch should be returned.
func (f *GridFSFindOptions) SetLimit(i int32) *GridFSFindOptions {
	f.Limit = &i
	return f
}

// SetMaxTime specifies the max time to allow the query to run.
func (f *GridFSFindOptions) SetMaxTime(d time.Duration) *GridFSFindOptions {
	f.MaxTime = &d
	return f
}

// SetNoCursorTimeout specifies whether or not cursors should time out after a period of inactivity.
func (f *GridFSFindOptions) SetNoCursorTimeout(b bool) *GridFSFindOptions {
	f.NoCursorTimeout = &b
	return f
}

// SetSkip specifies the number of documents to skip before returning.
func (f *GridFSFindOptions) SetSkip(i int32) *GridFSFindOptions {
	f.Skip = &i
	return f
}

// SetSort specifies the order in which to return documents.
func (f *GridFSFindOptions) SetSort(sort interface{}) *GridFSFindOptions {
	f.Sort = sort
	return f
}

// MergeGridFSFindOptions combines the argued GridFSFindOptions into a single GridFSFindOptions in a last-one-wins fashion
func MergeGridFSFindOptions(opts ...*GridFSFindOptions) *GridFSFindOptions {
	fo := GridFSFind()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.BatchSize != nil {
			fo.BatchSize = opt.BatchSize
		}
		if opt.Limit != nil {
			fo.Limit = opt.Limit
		}
		if opt.MaxTime != nil {
			fo.MaxTime = opt.MaxTime
		}
		if opt.NoCursorTimeout != nil {
			fo.NoCursorTimeout = opt.NoCursorTimeout
		}
		if opt.Skip != nil {
			fo.Skip = opt.Skip
		}
		if opt.Sort != nil {
			fo.Sort = opt.Sort
		}
	}

	return fo
}
