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

// DefaultChunkSize is the default size of each file chunk in bytes (255 KiB).
var DefaultChunkSize int32 = 255 * 1024

// DefaultRevision is the default revision number for a download by name operation.
var DefaultRevision int32 = -1

// BucketOptions represents options that can be used to configure GridFS bucket.
type BucketOptions struct {
	// The name of the bucket. The default value is "fs".
	Name *string

	// The number of bytes in each chunk in the bucket. The default value is 255 KiB.
	ChunkSizeBytes *int32

	// The write concern for the bucket. The default value is the write concern of the database from which the bucket
	// is created.
	WriteConcern *writeconcern.WriteConcern

	// The read concern for the bucket. The default value is the read concern of the database from which the bucket
	// is created.
	ReadConcern *readconcern.ReadConcern

	// The read preference for the bucket. The default value is the read preference of the database from which the
	// bucket is created.
	ReadPreference *readpref.ReadPref
}

// GridFSBucket creates a new BucketOptions instance.
func GridFSBucket() *BucketOptions {
	return &BucketOptions{
		Name:           &DefaultName,
		ChunkSizeBytes: &DefaultChunkSize,
	}
}

// SetName sets the value for the Name field.
func (b *BucketOptions) SetName(name string) *BucketOptions {
	b.Name = &name
	return b
}

// SetChunkSizeBytes sets the value for the ChunkSize field.
func (b *BucketOptions) SetChunkSizeBytes(i int32) *BucketOptions {
	b.ChunkSizeBytes = &i
	return b
}

// SetWriteConcern sets the value for the WriteConcern field.
func (b *BucketOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *BucketOptions {
	b.WriteConcern = wc
	return b
}

// SetReadConcern sets the value for the ReadConcern field.
func (b *BucketOptions) SetReadConcern(rc *readconcern.ReadConcern) *BucketOptions {
	b.ReadConcern = rc
	return b
}

// SetReadPreference sets the value for the ReadPreference field.
func (b *BucketOptions) SetReadPreference(rp *readpref.ReadPref) *BucketOptions {
	b.ReadPreference = rp
	return b
}

// MergeBucketOptions combines the given BucketOptions instances into a single BucketOptions in a last-one-wins fashion.
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

// UploadOptions represents options that can be used to configure a GridFS upload operation.
type UploadOptions struct {
	// The number of bytes in each chunk in the bucket. The default value is DefaultChunkSize (255 KiB).
	ChunkSizeBytes *int32

	// Additional application data that will be stored in the "metadata" field of the document in the files collection.
	// The default value is nil, which means that the document in the files collection will not contain a "metadata"
	// field.
	Metadata interface{}

	// The BSON registry to use for converting filters to BSON documents. The default value is bson.DefaultRegistry.
	Registry *bsoncodec.Registry
}

// GridFSUpload creates a new UploadOptions instance.
func GridFSUpload() *UploadOptions {
	return &UploadOptions{Registry: bson.DefaultRegistry}
}

// SetChunkSizeBytes sets the value for the ChunkSize field.
func (u *UploadOptions) SetChunkSizeBytes(i int32) *UploadOptions {
	u.ChunkSizeBytes = &i
	return u
}

// SetMetadata sets the value for the Metadata field.
func (u *UploadOptions) SetMetadata(doc interface{}) *UploadOptions {
	u.Metadata = doc
	return u
}

// MergeUploadOptions combines the given UploadOptions instances into a single UploadOptions in a last-one-wins fashion.
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

// NameOptions represents options that can be used to configure a GridFS DownloadByName operation.
type NameOptions struct {
	// Specifies the revision of the file to retrieve. Revision numbers are defined as follows:
	//
	// * 0 = the original stored file
	// * 1 = the first revision
	// * 2 = the second revision
	// * etc..
	// * -2 = the second most recent revision
	// * -1 = the most recent revision.
	//
	// The default value is -1
	Revision *int32
}

// GridFSName creates a new NameOptions instance.
func GridFSName() *NameOptions {
	return &NameOptions{}
}

// SetRevision sets the value for the Revision field.
func (n *NameOptions) SetRevision(r int32) *NameOptions {
	n.Revision = &r
	return n
}

// MergeNameOptions combines the given NameOptions instances into a single *NameOptions in a last-one-wins fashion.
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

// GridFSFindOptions represents options that can be used to configure a GridFS Find operation.
type GridFSFindOptions struct {
	// If true, the server can write temporary data to disk while executing the find operation. The default value
	// is false. This option is only valid for MongoDB versions >= 4.4. For previous server versions, the server will
	// return an error if this option is used.
	AllowDiskUse *bool

	// The maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// The maximum number of documents to return. The default value is 0, which means that all documents matching the
	// filter will be returned. A negative limit specifies that the resulting documents should be returned in a single
	// batch. The default value is 0.
	Limit *int32

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	MaxTime *time.Duration

	// If true, the cursor created by the operation will not timeout after a period of inactivity. The default value
	// is false.
	NoCursorTimeout *bool

	// The number of documents to skip before adding documents to the result. The default value is 0.
	Skip *int32

	// A document specifying the order in which documents should be returned.  The driver will return an error if the
	// sort parameter is a multi-key map.
	Sort interface{}
}

// GridFSFind creates a new GridFSFindOptions instance.
func GridFSFind() *GridFSFindOptions {
	return &GridFSFindOptions{}
}

// SetAllowDiskUse sets the value for the AllowDiskUse field.
func (f *GridFSFindOptions) SetAllowDiskUse(b bool) *GridFSFindOptions {
	f.AllowDiskUse = &b
	return f
}

// SetBatchSize sets the value for the BatchSize field.
func (f *GridFSFindOptions) SetBatchSize(i int32) *GridFSFindOptions {
	f.BatchSize = &i
	return f
}

// SetLimit sets the value for the Limit field.
func (f *GridFSFindOptions) SetLimit(i int32) *GridFSFindOptions {
	f.Limit = &i
	return f
}

// SetMaxTime sets the value for the MaxTime field.
func (f *GridFSFindOptions) SetMaxTime(d time.Duration) *GridFSFindOptions {
	f.MaxTime = &d
	return f
}

// SetNoCursorTimeout sets the value for the NoCursorTimeout field.
func (f *GridFSFindOptions) SetNoCursorTimeout(b bool) *GridFSFindOptions {
	f.NoCursorTimeout = &b
	return f
}

// SetSkip sets the value for the Skip field.
func (f *GridFSFindOptions) SetSkip(i int32) *GridFSFindOptions {
	f.Skip = &i
	return f
}

// SetSort sets the value for the Sort field.
func (f *GridFSFindOptions) SetSort(sort interface{}) *GridFSFindOptions {
	f.Sort = sort
	return f
}

// MergeGridFSFindOptions combines the given GridFSFindOptions instances into a single GridFSFindOptions in a
// last-one-wins fashion.
func MergeGridFSFindOptions(opts ...*GridFSFindOptions) *GridFSFindOptions {
	fo := GridFSFind()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.AllowDiskUse != nil {
			fo.AllowDiskUse = opt.AllowDiskUse
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
