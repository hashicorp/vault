// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WriteModel is an interface implemented by models that can be used in a BulkWrite operation. Each WriteModel
// represents a write.
//
// This interface is implemented by InsertOneModel, DeleteOneModel, DeleteManyModel, ReplaceOneModel, UpdateOneModel,
// and UpdateManyModel. Custom implementations of this interface must not be used.
type WriteModel interface {
	writeModel()
}

// InsertOneModel is used to insert a single document in a BulkWrite operation.
type InsertOneModel struct {
	Document interface{}
}

// NewInsertOneModel creates a new InsertOneModel.
func NewInsertOneModel() *InsertOneModel {
	return &InsertOneModel{}
}

// SetDocument specifies the document to be inserted. The document cannot be nil. If it does not have an _id field when
// transformed into BSON, one will be added automatically to the marshalled document. The original document will not be
// modified.
func (iom *InsertOneModel) SetDocument(doc interface{}) *InsertOneModel {
	iom.Document = doc
	return iom
}

func (*InsertOneModel) writeModel() {}

// DeleteOneModel is used to delete at most one document in a BulkWriteOperation.
type DeleteOneModel struct {
	Filter    interface{}
	Collation *options.Collation
	Hint      interface{}
}

// NewDeleteOneModel creates a new DeleteOneModel.
func NewDeleteOneModel() *DeleteOneModel {
	return &DeleteOneModel{}
}

// SetFilter specifies a filter to use to select the document to delete. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (dom *DeleteOneModel) SetFilter(filter interface{}) *DeleteOneModel {
	dom.Filter = filter
	return dom
}

// SetCollation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (dom *DeleteOneModel) SetCollation(collation *options.Collation) *DeleteOneModel {
	dom.Collation = collation
	return dom
}

// SetHint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.4. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (dom *DeleteOneModel) SetHint(hint interface{}) *DeleteOneModel {
	dom.Hint = hint
	return dom
}

func (*DeleteOneModel) writeModel() {}

// DeleteManyModel is used to delete multiple documents in a BulkWrite operation.
type DeleteManyModel struct {
	Filter    interface{}
	Collation *options.Collation
	Hint      interface{}
}

// NewDeleteManyModel creates a new DeleteManyModel.
func NewDeleteManyModel() *DeleteManyModel {
	return &DeleteManyModel{}
}

// SetFilter specifies a filter to use to select documents to delete. The filter must be a document containing query
// operators. It cannot be nil.
func (dmm *DeleteManyModel) SetFilter(filter interface{}) *DeleteManyModel {
	dmm.Filter = filter
	return dmm
}

// SetCollation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (dmm *DeleteManyModel) SetCollation(collation *options.Collation) *DeleteManyModel {
	dmm.Collation = collation
	return dmm
}

// SetHint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.4. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (dmm *DeleteManyModel) SetHint(hint interface{}) *DeleteManyModel {
	dmm.Hint = hint
	return dmm
}

func (*DeleteManyModel) writeModel() {}

// ReplaceOneModel is used to replace at most one document in a BulkWrite operation.
type ReplaceOneModel struct {
	Collation   *options.Collation
	Upsert      *bool
	Filter      interface{}
	Replacement interface{}
	Hint        interface{}
}

// NewReplaceOneModel creates a new ReplaceOneModel.
func NewReplaceOneModel() *ReplaceOneModel {
	return &ReplaceOneModel{}
}

// SetHint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (rom *ReplaceOneModel) SetHint(hint interface{}) *ReplaceOneModel {
	rom.Hint = hint
	return rom
}

// SetFilter specifies a filter to use to select the document to replace. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (rom *ReplaceOneModel) SetFilter(filter interface{}) *ReplaceOneModel {
	rom.Filter = filter
	return rom
}

// SetReplacement specifies a document that will be used to replace the selected document. It cannot be nil and cannot
// contain any update operators (https://docs.mongodb.com/manual/reference/operator/update/).
func (rom *ReplaceOneModel) SetReplacement(rep interface{}) *ReplaceOneModel {
	rom.Replacement = rep
	return rom
}

// SetCollation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (rom *ReplaceOneModel) SetCollation(collation *options.Collation) *ReplaceOneModel {
	rom.Collation = collation
	return rom
}

// SetUpsert specifies whether or not the replacement document should be inserted if no document matching the filter is
// found. If an upsert is performed, the _id of the upserted document can be retrieved from the UpsertedIDs field of the
// BulkWriteResult.
func (rom *ReplaceOneModel) SetUpsert(upsert bool) *ReplaceOneModel {
	rom.Upsert = &upsert
	return rom
}

func (*ReplaceOneModel) writeModel() {}

// UpdateOneModel is used to update at most one document in a BulkWrite operation.
type UpdateOneModel struct {
	Collation    *options.Collation
	Upsert       *bool
	Filter       interface{}
	Update       interface{}
	ArrayFilters *options.ArrayFilters
	Hint         interface{}
}

// NewUpdateOneModel creates a new UpdateOneModel.
func NewUpdateOneModel() *UpdateOneModel {
	return &UpdateOneModel{}
}

// SetHint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (uom *UpdateOneModel) SetHint(hint interface{}) *UpdateOneModel {
	uom.Hint = hint
	return uom
}

// SetFilter specifies a filter to use to select the document to update. The filter must be a document containing query
// operators. It cannot be nil. If the filter matches multiple documents, one will be selected from the matching
// documents.
func (uom *UpdateOneModel) SetFilter(filter interface{}) *UpdateOneModel {
	uom.Filter = filter
	return uom
}

// SetUpdate specifies the modifications to be made to the selected document. The value must be a document containing
// update operators (https://docs.mongodb.com/manual/reference/operator/update/). It cannot be nil or empty.
func (uom *UpdateOneModel) SetUpdate(update interface{}) *UpdateOneModel {
	uom.Update = update
	return uom
}

// SetArrayFilters specifies a set of filters to determine which elements should be modified when updating an array
// field.
func (uom *UpdateOneModel) SetArrayFilters(filters options.ArrayFilters) *UpdateOneModel {
	uom.ArrayFilters = &filters
	return uom
}

// SetCollation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (uom *UpdateOneModel) SetCollation(collation *options.Collation) *UpdateOneModel {
	uom.Collation = collation
	return uom
}

// SetUpsert specifies whether or not a new document should be inserted if no document matching the filter is found. If
// an upsert is performed, the _id of the upserted document can be retrieved from the UpsertedIDs field of the
// BulkWriteResult.
func (uom *UpdateOneModel) SetUpsert(upsert bool) *UpdateOneModel {
	uom.Upsert = &upsert
	return uom
}

func (*UpdateOneModel) writeModel() {}

// UpdateManyModel is used to update multiple documents in a BulkWrite operation.
type UpdateManyModel struct {
	Collation    *options.Collation
	Upsert       *bool
	Filter       interface{}
	Update       interface{}
	ArrayFilters *options.ArrayFilters
	Hint         interface{}
}

// NewUpdateManyModel creates a new UpdateManyModel.
func NewUpdateManyModel() *UpdateManyModel {
	return &UpdateManyModel{}
}

// SetHint specifies the index to use for the operation. This should either be the index name as a string or the index
// specification as a document. This option is only valid for MongoDB versions >= 4.2. Server versions >= 3.4 will
// return an error if this option is specified. For server versions < 3.4, the driver will return a client-side error if
// this option is specified. The driver will return an error if this option is specified during an unacknowledged write
// operation. The driver will return an error if the hint parameter is a multi-key map. The default value is nil, which
// means that no hint will be sent.
func (umm *UpdateManyModel) SetHint(hint interface{}) *UpdateManyModel {
	umm.Hint = hint
	return umm
}

// SetFilter specifies a filter to use to select documents to update. The filter must be a document containing query
// operators. It cannot be nil.
func (umm *UpdateManyModel) SetFilter(filter interface{}) *UpdateManyModel {
	umm.Filter = filter
	return umm
}

// SetUpdate specifies the modifications to be made to the selected documents. The value must be a document containing
// update operators (https://docs.mongodb.com/manual/reference/operator/update/). It cannot be nil or empty.
func (umm *UpdateManyModel) SetUpdate(update interface{}) *UpdateManyModel {
	umm.Update = update
	return umm
}

// SetArrayFilters specifies a set of filters to determine which elements should be modified when updating an array
// field.
func (umm *UpdateManyModel) SetArrayFilters(filters options.ArrayFilters) *UpdateManyModel {
	umm.ArrayFilters = &filters
	return umm
}

// SetCollation specifies a collation to use for string comparisons. The default is nil, meaning no collation will be
// used.
func (umm *UpdateManyModel) SetCollation(collation *options.Collation) *UpdateManyModel {
	umm.Collation = collation
	return umm
}

// SetUpsert specifies whether or not a new document should be inserted if no document matching the filter is found. If
// an upsert is performed, the _id of the upserted document can be retrieved from the UpsertedIDs field of the
// BulkWriteResult.
func (umm *UpdateManyModel) SetUpsert(upsert bool) *UpdateManyModel {
	umm.Upsert = &upsert
	return umm
}

func (*UpdateManyModel) writeModel() {}
