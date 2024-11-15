// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// SearchIndexesOptions represents options that can be used to configure a SearchIndexView.
type SearchIndexesOptions struct {
	Name *string
	Type *string
}

// SearchIndexes creates a new SearchIndexesOptions instance.
func SearchIndexes() *SearchIndexesOptions {
	return &SearchIndexesOptions{}
}

// SetName sets the value for the Name field.
func (sio *SearchIndexesOptions) SetName(name string) *SearchIndexesOptions {
	sio.Name = &name
	return sio
}

// SetType sets the value for the Type field.
func (sio *SearchIndexesOptions) SetType(typ string) *SearchIndexesOptions {
	sio.Type = &typ
	return sio
}

// CreateSearchIndexesOptions represents options that can be used to configure a SearchIndexView.CreateOne or
// SearchIndexView.CreateMany operation.
type CreateSearchIndexesOptions struct {
}

// ListSearchIndexesOptions represents options that can be used to configure a SearchIndexView.List operation.
type ListSearchIndexesOptions struct {
	AggregateOpts *AggregateOptions
}

// DropSearchIndexOptions represents options that can be used to configure a SearchIndexView.DropOne operation.
type DropSearchIndexOptions struct {
}

// UpdateSearchIndexOptions represents options that can be used to configure a SearchIndexView.UpdateOne operation.
type UpdateSearchIndexOptions struct {
}
