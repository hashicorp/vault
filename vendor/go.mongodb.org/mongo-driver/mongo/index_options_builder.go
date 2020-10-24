// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

// IndexOptionsBuilder specifies options for a new index.
//
// Deprecated: Use the IndexOptions type in the mongo/options package instead.
type IndexOptionsBuilder struct {
	document bson.D
}

// NewIndexOptionsBuilder creates a new IndexOptionsBuilder.
//
// Deprecated: Use the Index function in mongo/options instead.
func NewIndexOptionsBuilder() *IndexOptionsBuilder {
	return &IndexOptionsBuilder{}
}

// Background specifies a value for the background option.
//
// Deprecated: Use the IndexOptions.SetBackground function in mongo/options instead.
func (iob *IndexOptionsBuilder) Background(background bool) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"background", background})
	return iob
}

// ExpireAfterSeconds specifies a value for the expireAfterSeconds option.
//
// Deprecated: Use the IndexOptions.SetExpireAfterSeconds function in mongo/options instead.
func (iob *IndexOptionsBuilder) ExpireAfterSeconds(expireAfterSeconds int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"expireAfterSeconds", expireAfterSeconds})
	return iob
}

// Name specifies a value for the name option.
//
// Deprecated: Use the IndexOptions.SetName function in mongo/options instead.
func (iob *IndexOptionsBuilder) Name(name string) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"name", name})
	return iob
}

// Sparse specifies a value for the sparse option.
//
// Deprecated: Use the IndexOptions.SetSparse function in mongo/options instead.
func (iob *IndexOptionsBuilder) Sparse(sparse bool) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"sparse", sparse})
	return iob
}

// StorageEngine specifies a value for the storageEngine option.
//
// Deprecated: Use the IndexOptions.SetStorageEngine function in mongo/options instead.
func (iob *IndexOptionsBuilder) StorageEngine(storageEngine interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"storageEngine", storageEngine})
	return iob
}

// Unique specifies a value for the unique option.
//
// Deprecated: Use the IndexOptions.SetUnique function in mongo/options instead.
func (iob *IndexOptionsBuilder) Unique(unique bool) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"unique", unique})
	return iob
}

// Version specifies a value for the version option.
//
// Deprecated: Use the IndexOptions.SetVersion function in mongo/options instead.
func (iob *IndexOptionsBuilder) Version(version int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"v", version})
	return iob
}

// DefaultLanguage specifies a value for the default_language option.
//
// Deprecated: Use the IndexOptions.SetDefaultLanguage function in mongo/options instead.
func (iob *IndexOptionsBuilder) DefaultLanguage(defaultLanguage string) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"default_language", defaultLanguage})
	return iob
}

// LanguageOverride specifies a value for the language_override option.
//
// Deprecated: Use the IndexOptions.SetLanguageOverride function in mongo/options instead.
func (iob *IndexOptionsBuilder) LanguageOverride(languageOverride string) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"language_override", languageOverride})
	return iob
}

// TextVersion specifies a value for the textIndexVersion option.
//
// Deprecated: Use the IndexOptions.SetTextVersion function in mongo/options instead.
func (iob *IndexOptionsBuilder) TextVersion(textVersion int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"textIndexVersion", textVersion})
	return iob
}

// Weights specifies a value for the weights option.
//
// Deprecated: Use the IndexOptions.SetWeights function in mongo/options instead.
func (iob *IndexOptionsBuilder) Weights(weights interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"weights", weights})
	return iob
}

// SphereVersion specifies a value for the 2dsphereIndexVersion option.
//
// Deprecated: Use the IndexOptions.SetSphereVersion function in mongo/options instead.
func (iob *IndexOptionsBuilder) SphereVersion(sphereVersion int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"2dsphereIndexVersion", sphereVersion})
	return iob
}

// Bits specifies a value for the bits option.
//
// Deprecated: Use the IndexOptions.SetBits function in mongo/options instead.
func (iob *IndexOptionsBuilder) Bits(bits int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"bits", bits})
	return iob
}

// Max specifies a value for the max option.
//
// Deprecated: Use the IndexOptions.SetMax function in mongo/options instead.
func (iob *IndexOptionsBuilder) Max(max float64) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"max", max})
	return iob
}

// Min specifies a value for the min option.
//
// Deprecated: Use the IndexOptions.SetMin function in mongo/options instead.
func (iob *IndexOptionsBuilder) Min(min float64) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"min", min})
	return iob
}

// BucketSize specifies a value for the bucketSize option.
//
// Deprecated: Use the IndexOptions.SetBucketSize function in mongo/options instead.
func (iob *IndexOptionsBuilder) BucketSize(bucketSize int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"bucketSize", bucketSize})
	return iob
}

// PartialFilterExpression specifies a value for the partialFilterExpression option.
//
// Deprecated: Use the IndexOptions.SetPartialFilterExpression function in mongo/options instead.
func (iob *IndexOptionsBuilder) PartialFilterExpression(partialFilterExpression interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"partialFilterExpression", partialFilterExpression})
	return iob
}

// Collation specifies a value for the collation option.
//
// Deprecated: Use the IndexOptions.SetCollation function in mongo/options instead.
func (iob *IndexOptionsBuilder) Collation(collation interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"collation", collation})
	return iob
}

// Build finishes constructing an the builder.
//
// Deprecated: Use the IndexOptions type in the mongo/options package instead.
func (iob *IndexOptionsBuilder) Build() bson.D {
	return iob.document
}
