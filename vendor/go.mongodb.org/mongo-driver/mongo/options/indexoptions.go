// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"
)

// CreateIndexesOptions represents options that can be used to configure IndexView.CreateOne and IndexView.CreateMany
// operations.
type CreateIndexesOptions struct {
	// The number of data-bearing members of a replica set, including the primary, that must complete the index builds
	// successfully before the primary marks the indexes as ready. This should either be a string or int32 value. The
	// semantics of the values are as follows:
	//
	// 1. String: specifies a tag. All members with that tag must complete the build.
	// 2. int: the number of members that must complete the build.
	// 3. "majority": A special value to indicate that more than half the nodes must complete the build.
	// 4. "votingMembers": A special value to indicate that all voting data-bearing nodes must complete.
	//
	// This option is only available on MongoDB versions >= 4.4. A client-side error will be returned if the option
	// is specified for MongoDB versions <= 4.2. The default value is nil, meaning that the server-side default will be
	// used. See dochub.mongodb.org/core/index-commit-quorum for more information.
	CommitQuorum interface{}

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	MaxTime *time.Duration
}

// CreateIndexes creates a new CreateIndexesOptions instance.
func CreateIndexes() *CreateIndexesOptions {
	return &CreateIndexesOptions{}
}

// SetMaxTime sets the value for the MaxTime field.
func (c *CreateIndexesOptions) SetMaxTime(d time.Duration) *CreateIndexesOptions {
	c.MaxTime = &d
	return c
}

// SetCommitQuorumInt sets the value for the CommitQuorum field as an int32.
func (c *CreateIndexesOptions) SetCommitQuorumInt(quorum int32) *CreateIndexesOptions {
	c.CommitQuorum = quorum
	return c
}

// SetCommitQuorumString sets the value for the CommitQuorum field as a string.
func (c *CreateIndexesOptions) SetCommitQuorumString(quorum string) *CreateIndexesOptions {
	c.CommitQuorum = quorum
	return c
}

// SetCommitQuorumMajority sets the value for the CommitQuorum to special "majority" value.
func (c *CreateIndexesOptions) SetCommitQuorumMajority() *CreateIndexesOptions {
	c.CommitQuorum = "majority"
	return c
}

// SetCommitQuorumVotingMembers sets the value for the CommitQuorum to special "votingMembers" value.
func (c *CreateIndexesOptions) SetCommitQuorumVotingMembers() *CreateIndexesOptions {
	c.CommitQuorum = "votingMembers"
	return c
}

// MergeCreateIndexesOptions combines the given CreateIndexesOptions into a single CreateIndexesOptions in a last one
// wins fashion.
func MergeCreateIndexesOptions(opts ...*CreateIndexesOptions) *CreateIndexesOptions {
	c := CreateIndexes()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.MaxTime != nil {
			c.MaxTime = opt.MaxTime
		}
		if opt.CommitQuorum != nil {
			c.CommitQuorum = opt.CommitQuorum
		}
	}

	return c
}

// DropIndexesOptions represents options that can be used to configure IndexView.DropOne and IndexView.DropAll
// operations.
type DropIndexesOptions struct {
	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	MaxTime *time.Duration
}

// DropIndexes creates a new DropIndexesOptions instance.
func DropIndexes() *DropIndexesOptions {
	return &DropIndexesOptions{}
}

// SetMaxTime sets the value for the MaxTime field.
func (d *DropIndexesOptions) SetMaxTime(duration time.Duration) *DropIndexesOptions {
	d.MaxTime = &duration
	return d
}

// MergeDropIndexesOptions combines the given DropIndexesOptions into a single DropIndexesOptions in a last-one-wins
// fashion.
func MergeDropIndexesOptions(opts ...*DropIndexesOptions) *DropIndexesOptions {
	c := DropIndexes()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.MaxTime != nil {
			c.MaxTime = opt.MaxTime
		}
	}

	return c
}

// ListIndexesOptions represents options that can be used to configure an IndexView.List operation.
type ListIndexesOptions struct {
	// The maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// The maximum amount of time that the query can run on the server. The default value is nil, meaning that there
	// is no time limit for query execution.
	MaxTime *time.Duration
}

// ListIndexes creates a new ListIndexesOptions instance.
func ListIndexes() *ListIndexesOptions {
	return &ListIndexesOptions{}
}

// SetBatchSize sets the value for the BatchSize field.
func (l *ListIndexesOptions) SetBatchSize(i int32) *ListIndexesOptions {
	l.BatchSize = &i
	return l
}

// SetMaxTime sets the value for the MaxTime field.
func (l *ListIndexesOptions) SetMaxTime(d time.Duration) *ListIndexesOptions {
	l.MaxTime = &d
	return l
}

// MergeListIndexesOptions combines the given ListIndexesOptions instances into a single *ListIndexesOptions in a
// last-one-wins fashion.
func MergeListIndexesOptions(opts ...*ListIndexesOptions) *ListIndexesOptions {
	c := ListIndexes()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.BatchSize != nil {
			c.BatchSize = opt.BatchSize
		}
		if opt.MaxTime != nil {
			c.MaxTime = opt.MaxTime
		}
	}

	return c
}

// IndexOptions represents options that can be used to configure a new index created through the IndexView.CreateOne
// or IndexView.CreateMany operations.
type IndexOptions struct {
	// If true, the index will be built in the background on the server and will not block other tasks. The default
	// value is false.
	//
	// Deprecated: This option has been deprecated in MongoDB version 4.2.
	Background *bool

	// The length of time, in seconds, for documents to remain in the collection. The default value is 0, which means
	// that documents will remain in the collection until they're explicitly deleted or the collection is dropped.
	ExpireAfterSeconds *int32

	// The name of the index. The default value is "[field1]_[direction1]_[field2]_[direction2]...". For example, an
	// index with the specification {name: 1, age: -1} will be named "name_1_age_-1".
	Name *string

	// If true, the index will only reference documents that contain the fields specified in the index. The default is
	// false.
	Sparse *bool

	// Specifies the storage engine to use for the index. The value must be a document in the form
	// {<storage engine name>: <options>}. The default value is nil, which means that the default storage engine
	// will be used. This option is only applicable for MongoDB versions >= 3.0 and is ignored for previous server
	// versions.
	StorageEngine interface{}

	// If true, the collection will not accept insertion or update of documents where the index key value matches an
	// existing value in the index. The default is false.
	Unique *bool

	// The index version number, either 0 or 1.
	Version *int32

	// The language that determines the list of stop words and the rules for the stemmer and tokenizer. This option
	// is only applicable for text indexes and is ignored for other index types. The default value is "english".
	DefaultLanguage *string

	// The name of the field in the collection's documents that contains the override language for the document. This
	// option is only applicable for text indexes and is ignored for other index types. The default value is the value
	// of the DefaultLanguage option.
	LanguageOverride *string

	// The index version number for a text index. See https://docs.mongodb.com/manual/core/index-text/#text-versions for
	// information about different version numbers.
	TextVersion *int32

	// A document that contains field and weight pairs. The weight is an integer ranging from 1 to 99,999, inclusive,
	// indicating the significance of the field relative to the other indexed fields in terms of the score. This option
	// is only applicable for text indexes and is ignored for other index types. The default value is nil, which means
	// that every field will have a weight of 1.
	Weights interface{}

	// The index version number for a 2D sphere index. See https://docs.mongodb.com/manual/core/2dsphere/#dsphere-v2 for
	// information about different version numbers.
	SphereVersion *int32

	// The precision of the stored geohash value of the location data. This option only applies to 2D indexes and is
	// ignored for other index types. The value must be between 1 and 32, inclusive. The default value is 26.
	Bits *int32

	// The upper inclusive boundary for longitude and latitude values. This option is only applicable to 2D indexes and
	// is ignored for other index types. The default value is 180.0.
	Max *float64

	// The lower inclusive boundary for longitude and latitude values. This option is only applicable to 2D indexes and
	// is ignored for other index types. The default value is -180.0.
	Min *float64

	// The number of units within which to group location values. Location values that are within BucketSize units of
	// each other will be grouped in the same bucket. This option is only applicable to geoHaystack indexes and is
	// ignored for other index types. The value must be greater than 0.
	BucketSize *int32

	// A document that defines which collection documents the index should reference. This option is only valid for
	// MongoDB versions >= 3.2 and is ignored for previous server versions.
	PartialFilterExpression interface{}

	// The collation to use for string comparisons for the index. This option is only valid for MongoDB versions >= 3.4.
	// For previous server versions, the driver will return an error if this option is used.
	Collation *Collation

	// A document that defines the wildcard projection for the index.
	WildcardProjection interface{}

	// If true, the index will exist on the target collection but will not be used by the query planner when executing
	// operations. This option is only valid for MongoDB versions >= 4.4. The default value is false.
	Hidden *bool
}

// Index creates a new IndexOptions instance.
func Index() *IndexOptions {
	return &IndexOptions{}
}

// SetBackground sets value for the Background field.
//
// Deprecated: This option has been deprecated in MongoDB version 4.2.
func (i *IndexOptions) SetBackground(background bool) *IndexOptions {
	i.Background = &background
	return i
}

// SetExpireAfterSeconds sets value for the ExpireAfterSeconds field.
func (i *IndexOptions) SetExpireAfterSeconds(seconds int32) *IndexOptions {
	i.ExpireAfterSeconds = &seconds
	return i
}

// SetName sets the value for the Name field.
func (i *IndexOptions) SetName(name string) *IndexOptions {
	i.Name = &name
	return i
}

// SetSparse sets the value of the Sparse field.
func (i *IndexOptions) SetSparse(sparse bool) *IndexOptions {
	i.Sparse = &sparse
	return i
}

// SetStorageEngine sets the value for the StorageEngine field.
func (i *IndexOptions) SetStorageEngine(engine interface{}) *IndexOptions {
	i.StorageEngine = engine
	return i
}

// SetUnique sets the value for the Unique field.
func (i *IndexOptions) SetUnique(unique bool) *IndexOptions {
	i.Unique = &unique
	return i
}

// SetVersion sets the value for the Version field.
func (i *IndexOptions) SetVersion(version int32) *IndexOptions {
	i.Version = &version
	return i
}

// SetDefaultLanguage sets the value for the DefaultLanguage field.
func (i *IndexOptions) SetDefaultLanguage(language string) *IndexOptions {
	i.DefaultLanguage = &language
	return i
}

// SetLanguageOverride sets the value of the LanguageOverride field.
func (i *IndexOptions) SetLanguageOverride(override string) *IndexOptions {
	i.LanguageOverride = &override
	return i
}

// SetTextVersion sets the value for the TextVersion field.
func (i *IndexOptions) SetTextVersion(version int32) *IndexOptions {
	i.TextVersion = &version
	return i
}

// SetWeights sets the value for the Weights field.
func (i *IndexOptions) SetWeights(weights interface{}) *IndexOptions {
	i.Weights = weights
	return i
}

// SetSphereVersion sets the value for the SphereVersion field.
func (i *IndexOptions) SetSphereVersion(version int32) *IndexOptions {
	i.SphereVersion = &version
	return i
}

// SetBits sets the value for the Bits field.
func (i *IndexOptions) SetBits(bits int32) *IndexOptions {
	i.Bits = &bits
	return i
}

// SetMax sets the value for the Max field.
func (i *IndexOptions) SetMax(max float64) *IndexOptions {
	i.Max = &max
	return i
}

// SetMin sets the value for the Min field.
func (i *IndexOptions) SetMin(min float64) *IndexOptions {
	i.Min = &min
	return i
}

// SetBucketSize sets the value for the BucketSize field
func (i *IndexOptions) SetBucketSize(bucketSize int32) *IndexOptions {
	i.BucketSize = &bucketSize
	return i
}

// SetPartialFilterExpression sets the value for the PartialFilterExpression field.
func (i *IndexOptions) SetPartialFilterExpression(expression interface{}) *IndexOptions {
	i.PartialFilterExpression = expression
	return i
}

// SetCollation sets the value for the Collation field.
func (i *IndexOptions) SetCollation(collation *Collation) *IndexOptions {
	i.Collation = collation
	return i
}

// SetWildcardProjection sets the value for the WildcardProjection field.
func (i *IndexOptions) SetWildcardProjection(wildcardProjection interface{}) *IndexOptions {
	i.WildcardProjection = wildcardProjection
	return i
}

// SetHidden sets the value for the Hidden field.
func (i *IndexOptions) SetHidden(hidden bool) *IndexOptions {
	i.Hidden = &hidden
	return i
}

// MergeIndexOptions combines the given IndexOptions into a single IndexOptions in a last-one-wins fashion.
func MergeIndexOptions(opts ...*IndexOptions) *IndexOptions {
	i := Index()

	for _, opt := range opts {
		if opt.Background != nil {
			i.Background = opt.Background
		}
		if opt.ExpireAfterSeconds != nil {
			i.ExpireAfterSeconds = opt.ExpireAfterSeconds
		}
		if opt.Name != nil {
			i.Name = opt.Name
		}
		if opt.Sparse != nil {
			i.Sparse = opt.Sparse
		}
		if opt.StorageEngine != nil {
			i.StorageEngine = opt.StorageEngine
		}
		if opt.Unique != nil {
			i.Unique = opt.Unique
		}
		if opt.Version != nil {
			i.Version = opt.Version
		}
		if opt.DefaultLanguage != nil {
			i.DefaultLanguage = opt.DefaultLanguage
		}
		if opt.LanguageOverride != nil {
			i.LanguageOverride = opt.LanguageOverride
		}
		if opt.TextVersion != nil {
			i.TextVersion = opt.TextVersion
		}
		if opt.Weights != nil {
			i.Weights = opt.Weights
		}
		if opt.SphereVersion != nil {
			i.SphereVersion = opt.SphereVersion
		}
		if opt.Bits != nil {
			i.Bits = opt.Bits
		}
		if opt.Max != nil {
			i.Max = opt.Max
		}
		if opt.Min != nil {
			i.Min = opt.Min
		}
		if opt.BucketSize != nil {
			i.BucketSize = opt.BucketSize
		}
		if opt.PartialFilterExpression != nil {
			i.PartialFilterExpression = opt.PartialFilterExpression
		}
		if opt.Collation != nil {
			i.Collation = opt.Collation
		}
		if opt.WildcardProjection != nil {
			i.WildcardProjection = opt.WildcardProjection
		}
		if opt.Hidden != nil {
			i.Hidden = opt.Hidden
		}
	}

	return i
}
