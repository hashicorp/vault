// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// DefaultIndexOptions represents the default options for a collection to apply on new indexes. This type can be used
// when creating a new collection through the CreateCollectionOptions.SetDefaultIndexOptions method.
type DefaultIndexOptions struct {
	// Specifies the storage engine to use for the index. The value must be a document in the form
	// {<storage engine name>: <options>}. The default value is nil, which means that the default storage engine
	// will be used.
	StorageEngine interface{}
}

// DefaultIndex creates a new DefaultIndexOptions instance.
func DefaultIndex() *DefaultIndexOptions {
	return &DefaultIndexOptions{}
}

// SetStorageEngine sets the value for the StorageEngine field.
func (d *DefaultIndexOptions) SetStorageEngine(storageEngine interface{}) *DefaultIndexOptions {
	d.StorageEngine = storageEngine
	return d
}

// CreateCollectionOptions represents options that can be used to configure a CreateCollection operation.
type CreateCollectionOptions struct {
	// Specifies if the collection is capped (see https://docs.mongodb.com/manual/core/capped-collections/). If true,
	// the SizeInBytes option must also be specified. The default value is false.
	Capped *bool

	// Specifies the default collation for the new collection. This option is only valid for MongoDB versions >= 3.4.
	// For previous server versions, the driver will return an error if this option is used. The default value is nil.
	Collation *Collation

	// Specifies a default configuration for indexes on the collection. This option is only valid for MongoDB versions
	// >= 3.4. The default value is nil, meaning indexes will be configured using server defaults.
	DefaultIndexOptions *DefaultIndexOptions

	// Specifies the maximum number of documents allowed in a capped collection. The limit specified by the SizeInBytes
	// option takes precedence over this option. If a capped collection reaches its size limit, old documents will be
	// removed, regardless of the number of documents in the collection. The default value is 0, meaning the maximum
	// number of documents is unbounded.
	MaxDocuments *int64

	// Specifies the maximum size in bytes for a capped collection. The default value is 0.
	SizeInBytes *int64

	// Specifies the storage engine to use for the index. The value must be a document in the form
	// {<storage engine name>: <options>}. The default value is nil, which means that the default storage engine
	// will be used.
	StorageEngine interface{}

	// Specifies what should happen if a document being inserted does not pass validation. Valid values are "error" and
	// "warn". See https://docs.mongodb.com/manual/core/schema-validation/#accept-or-reject-invalid-documents for more
	// information. This option is only valid for MongoDB versions >= 3.2. The default value is "error".
	ValidationAction *string

	// Specifies how strictly the server applies validation rules to existing documents in the collection during update
	// operations. Valid values are "off", "strict", and "moderate". See
	// https://docs.mongodb.com/manual/core/schema-validation/#existing-documents for more information. This option is
	// only valid for MongoDB versions >= 3.2. The default value is "strict".
	ValidationLevel *string

	// A document specifying validation rules for the collection. See
	// https://docs.mongodb.com/manual/core/schema-validation/ for more information about schema validation. This option
	// is only valid for MongoDB versions >= 3.2. The default value is nil, meaning no validator will be used for the
	// collection.
	Validator interface{}
}

// CreateCollection creates a new CreateCollectionOptions instance.
func CreateCollection() *CreateCollectionOptions {
	return &CreateCollectionOptions{}
}

// SetCapped sets the value for the Capped field.
func (c *CreateCollectionOptions) SetCapped(capped bool) *CreateCollectionOptions {
	c.Capped = &capped
	return c
}

// SetCollation sets the value for the Collation field.
func (c *CreateCollectionOptions) SetCollation(collation *Collation) *CreateCollectionOptions {
	c.Collation = collation
	return c
}

// SetDefaultIndexOptions sets the value for the DefaultIndexOptions field.
func (c *CreateCollectionOptions) SetDefaultIndexOptions(opts *DefaultIndexOptions) *CreateCollectionOptions {
	c.DefaultIndexOptions = opts
	return c
}

// SetMaxDocuments sets the value for the MaxDocuments field.
func (c *CreateCollectionOptions) SetMaxDocuments(max int64) *CreateCollectionOptions {
	c.MaxDocuments = &max
	return c
}

// SetSizeInBytes sets the value for the SizeInBytes field.
func (c *CreateCollectionOptions) SetSizeInBytes(size int64) *CreateCollectionOptions {
	c.SizeInBytes = &size
	return c
}

// SetStorageEngine sets the value for the StorageEngine field.
func (c *CreateCollectionOptions) SetStorageEngine(storageEngine interface{}) *CreateCollectionOptions {
	c.StorageEngine = &storageEngine
	return c
}

// SetValidationAction sets the value for the ValidationAction field.
func (c *CreateCollectionOptions) SetValidationAction(action string) *CreateCollectionOptions {
	c.ValidationAction = &action
	return c
}

// SetValidationLevel sets the value for the ValidationLevel field.
func (c *CreateCollectionOptions) SetValidationLevel(level string) *CreateCollectionOptions {
	c.ValidationLevel = &level
	return c
}

// SetValidator sets the value for the Validator field.
func (c *CreateCollectionOptions) SetValidator(validator interface{}) *CreateCollectionOptions {
	c.Validator = validator
	return c
}

// MergeCreateCollectionOptions combines the given CreateCollectionOptions instances into a single
// CreateCollectionOptions in a last-one-wins fashion.
func MergeCreateCollectionOptions(opts ...*CreateCollectionOptions) *CreateCollectionOptions {
	cc := CreateCollection()

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.Capped != nil {
			cc.Capped = opt.Capped
		}
		if opt.Collation != nil {
			cc.Collation = opt.Collation
		}
		if opt.DefaultIndexOptions != nil {
			cc.DefaultIndexOptions = opt.DefaultIndexOptions
		}
		if opt.MaxDocuments != nil {
			cc.MaxDocuments = opt.MaxDocuments
		}
		if opt.SizeInBytes != nil {
			cc.SizeInBytes = opt.SizeInBytes
		}
		if opt.StorageEngine != nil {
			cc.StorageEngine = opt.StorageEngine
		}
		if opt.ValidationAction != nil {
			cc.ValidationAction = opt.ValidationAction
		}
		if opt.ValidationLevel != nil {
			cc.ValidationLevel = opt.ValidationLevel
		}
		if opt.Validator != nil {
			cc.Validator = opt.Validator
		}
	}

	return cc
}

// CreateViewOptions represents options that can be used to configure a CreateView operation.
type CreateViewOptions struct {
	// Specifies the default collation for the new collection. This option is only valid for MongoDB versions >= 3.4.
	// For previous server versions, the driver will return an error if this option is used. The default value is nil.
	Collation *Collation
}

// CreateView creates an new CreateViewOptions instance.
func CreateView() *CreateViewOptions {
	return &CreateViewOptions{}
}

// SetCollation sets the value for the Collation field.
func (c *CreateViewOptions) SetCollation(collation *Collation) *CreateViewOptions {
	c.Collation = collation
	return c
}

// MergeCreateViewOptions combines the given CreateViewOptions instances into a single CreateViewOptions in a
// last-one-wins fashion.
func MergeCreateViewOptions(opts ...*CreateViewOptions) *CreateViewOptions {
	cv := CreateView()

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.Collation != nil {
			cv.Collation = opt.Collation
		}
	}

	return cv
}
