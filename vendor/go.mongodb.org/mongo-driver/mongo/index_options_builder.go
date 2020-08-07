// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

// IndexOptionsBuilder is deprecated and unused.  Use mongo/options.IndexOptions instead.
type IndexOptionsBuilder struct {
	document bson.D
}

// NewIndexOptionsBuilder is deprecated.
func NewIndexOptionsBuilder() *IndexOptionsBuilder {
	return &IndexOptionsBuilder{}
}

// Background is deprecated.
func (iob *IndexOptionsBuilder) Background(background bool) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"background", background})
	return iob
}

// ExpireAfterSeconds is deprecated.
func (iob *IndexOptionsBuilder) ExpireAfterSeconds(expireAfterSeconds int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"expireAfterSeconds", expireAfterSeconds})
	return iob
}

// Name is deprecated.
func (iob *IndexOptionsBuilder) Name(name string) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"name", name})
	return iob
}

// Sparse is deprecated.
func (iob *IndexOptionsBuilder) Sparse(sparse bool) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"sparse", sparse})
	return iob
}

// StorageEngine is deprecated.
func (iob *IndexOptionsBuilder) StorageEngine(storageEngine interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"storageEngine", storageEngine})
	return iob
}

// Unique is deprecated.
func (iob *IndexOptionsBuilder) Unique(unique bool) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"unique", unique})
	return iob
}

// Version is deprecated.
func (iob *IndexOptionsBuilder) Version(version int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"v", version})
	return iob
}

// DefaultLanguage is deprecated.
func (iob *IndexOptionsBuilder) DefaultLanguage(defaultLanguage string) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"default_language", defaultLanguage})
	return iob
}

// LanguageOverride is deprecated.
func (iob *IndexOptionsBuilder) LanguageOverride(languageOverride string) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"language_override", languageOverride})
	return iob
}

// TextVersion is deprecated.
func (iob *IndexOptionsBuilder) TextVersion(textVersion int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"textIndexVersion", textVersion})
	return iob
}

// Weights is deprecated.
func (iob *IndexOptionsBuilder) Weights(weights interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"weights", weights})
	return iob
}

// SphereVersion is deprecated.
func (iob *IndexOptionsBuilder) SphereVersion(sphereVersion int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"2dsphereIndexVersion", sphereVersion})
	return iob
}

// Bits is deprecated.
func (iob *IndexOptionsBuilder) Bits(bits int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"bits", bits})
	return iob
}

// Max is deprecated.
func (iob *IndexOptionsBuilder) Max(max float64) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"max", max})
	return iob
}

// Min is deprecated.
func (iob *IndexOptionsBuilder) Min(min float64) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"min", min})
	return iob
}

// BucketSize is deprecated.
func (iob *IndexOptionsBuilder) BucketSize(bucketSize int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"bucketSize", bucketSize})
	return iob
}

// PartialFilterExpression is deprecated.
func (iob *IndexOptionsBuilder) PartialFilterExpression(partialFilterExpression interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"partialFilterExpression", partialFilterExpression})
	return iob
}

// Collation is deprecated.
func (iob *IndexOptionsBuilder) Collation(collation interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"collation", collation})
	return iob
}

// Build is deprecated.
func (iob *IndexOptionsBuilder) Build() bson.D {
	return iob.document
}
