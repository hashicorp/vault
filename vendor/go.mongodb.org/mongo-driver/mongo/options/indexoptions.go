// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"time"
)

// CreateIndexesOptions represents all possible options for the CreateOne() and CreateMany() functions.
type CreateIndexesOptions struct {
	MaxTime *time.Duration // The maximum amount of time to allow the query to run.
}

// CreateIndexes creates a new CreateIndexesOptions instance.
func CreateIndexes() *CreateIndexesOptions {
	return &CreateIndexesOptions{}
}

// SetMaxTime specifies the maximum amount of time to allow the query to run.
func (c *CreateIndexesOptions) SetMaxTime(d time.Duration) *CreateIndexesOptions {
	c.MaxTime = &d
	return c
}

// MergeCreateIndexesOptions combines the given *CreateIndexesOptions into a single *CreateIndexesOptions in a last one
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
	}

	return c
}

// DropIndexesOptions represents all possible options for the DropIndexes() function.
type DropIndexesOptions struct {
	MaxTime *time.Duration
}

// DropIndexes creates a new DropIndexesOptions instance.
func DropIndexes() *DropIndexesOptions {
	return &DropIndexesOptions{}
}

// SetMaxTime specifies the maximum amount of time to allow the query to run.
func (d *DropIndexesOptions) SetMaxTime(duration time.Duration) *DropIndexesOptions {
	d.MaxTime = &duration
	return d
}

// MergeDropIndexesOptions combines the given *DropIndexesOptions into a single *DropIndexesOptions in a last one
// wins fashion.
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

// ListIndexesOptions represents all possible options for the ListIndexes() function.
type ListIndexesOptions struct {
	BatchSize *int32
	MaxTime   *time.Duration
}

// ListIndexes creates a new ListIndexesOptions instance.
func ListIndexes() *ListIndexesOptions {
	return &ListIndexesOptions{}
}

// SetBatchSize specifies the number of documents to return in every batch.
func (l *ListIndexesOptions) SetBatchSize(i int32) *ListIndexesOptions {
	l.BatchSize = &i
	return l
}

// SetMaxTime specifies the maximum amount of time to allow the query to run.
func (l *ListIndexesOptions) SetMaxTime(d time.Duration) *ListIndexesOptions {
	l.MaxTime = &d
	return l
}

// MergeListIndexesOptions combines the given *ListIndexesOptions into a single *ListIndexesOptions in a last one
// wins fashion.
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

// IndexOptions represents all possible options to configure a new index.
type IndexOptions struct {
	Background              *bool
	ExpireAfterSeconds      *int32
	Name                    *string
	Sparse                  *bool
	StorageEngine           interface{}
	Unique                  *bool
	Version                 *int32
	DefaultLanguage         *string
	LanguageOverride        *string
	TextVersion             *int32
	Weights                 interface{}
	SphereVersion           *int32
	Bits                    *int32
	Max                     *float64
	Min                     *float64
	BucketSize              *int32
	PartialFilterExpression interface{}
	Collation               *Collation
	WildcardProjection      interface{}
}

// Index creates a new *IndexOptions
func Index() *IndexOptions {
	return &IndexOptions{}
}

// SetBackground sets the background option. If true, the server will create the index in the background and not block
// other tasks
func (i *IndexOptions) SetBackground(background bool) *IndexOptions {
	i.Background = &background
	return i
}

// SetExpireAfterSeconds specifies the number of seconds for a document to remain in a collection.
func (i *IndexOptions) SetExpireAfterSeconds(seconds int32) *IndexOptions {
	i.ExpireAfterSeconds = &seconds
	return i
}

// SetName specifies a name for the index.
// If not set, a name will be generated in the format "[field]_[direction]".
// If multiple indexes are created for the same key pattern with different collations, a name must be provided to avoid
// ambiguity.
func (i *IndexOptions) SetName(name string) *IndexOptions {
	i.Name = &name
	return i
}

// SetSparse sets the sparse option.
// If true, the index will only reference documents with the specified field in the index.
func (i *IndexOptions) SetSparse(sparse bool) *IndexOptions {
	i.Sparse = &sparse
	return i
}

// SetStorageEngine specifies the storage engine to use.
// Valid for server versions >= 3.0
func (i *IndexOptions) SetStorageEngine(engine interface{}) *IndexOptions {
	i.StorageEngine = engine
	return i
}

// SetUnique forces the index to be unique.
func (i *IndexOptions) SetUnique(unique bool) *IndexOptions {
	i.Unique = &unique
	return i
}

// SetVersion specifies the index version number, either 0 or 1.
func (i *IndexOptions) SetVersion(version int32) *IndexOptions {
	i.Version = &version
	return i
}

// SetDefaultLanguage specifies the default language for text indexes.
// If not set, this will default to english.
func (i *IndexOptions) SetDefaultLanguage(language string) *IndexOptions {
	i.DefaultLanguage = &language
	return i
}

// SetLanguageOverride specifies the field in the document to override the language.
func (i *IndexOptions) SetLanguageOverride(override string) *IndexOptions {
	i.LanguageOverride = &override
	return i
}

// SetTextVersion specifies the text index version number.
// MongoDB version 2.4 can only support version 1.
// MongoDB versions 2.6 and higher can support versions 1 or 2.
func (i *IndexOptions) SetTextVersion(version int32) *IndexOptions {
	i.TextVersion = &version
	return i
}

// SetWeights specifies fields in the index and their corresponding weight values.
func (i *IndexOptions) SetWeights(weights interface{}) *IndexOptions {
	i.Weights = weights
	return i
}

// SetSphereVersion specifies the 2dsphere index version number.
// MongoDB version 2.4 can only support version 1.
// MongoDB versions 2.6 and higher can support versions 1 or 2.
func (i *IndexOptions) SetSphereVersion(version int32) *IndexOptions {
	i.SphereVersion = &version
	return i
}

// SetBits specifies the precision of the stored geo hash in the 2d index, from 1 to 32.
func (i *IndexOptions) SetBits(bits int32) *IndexOptions {
	i.Bits = &bits
	return i
}

// SetMax specifies the maximum boundary for latitude and longitude in the 2d index.
func (i *IndexOptions) SetMax(max float64) *IndexOptions {
	i.Max = &max
	return i
}

// SetMin specifies the minimum boundary for latitude and longitude in the 2d index.
func (i *IndexOptions) SetMin(min float64) *IndexOptions {
	i.Min = &min
	return i
}

// SetBucketSize specifies number of units within which to group the location values in a geo haystack index.
func (i *IndexOptions) SetBucketSize(bucketSize int32) *IndexOptions {
	i.BucketSize = &bucketSize
	return i
}

// SetPartialFilterExpression specifies a filter for use in a partial index. Only documents that match the filter
// expression are included in the index.
func (i *IndexOptions) SetPartialFilterExpression(expression interface{}) *IndexOptions {
	i.PartialFilterExpression = expression
	return i
}

// SetCollation specifies a Collation to use for the operation.
// Valid for server versions >= 3.4
func (i *IndexOptions) SetCollation(collation *Collation) *IndexOptions {
	i.Collation = collation
	return i
}

// SetWildcardProjection specifies a wildcard projection for a wildcard index.
func (i *IndexOptions) SetWildcardProjection(wildcardProjection interface{}) *IndexOptions {
	i.WildcardProjection = wildcardProjection
	return i
}

// MergeIndexOptions combines the given *IndexOptions into a single *IndexOptions in a last one wins fashion.
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
	}

	return i
}
