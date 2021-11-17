/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protocol

// TableResult is the package internal representation of a table like output parameter of a stored procedure.
type TableResult struct {
	id             uint64
	resultFieldSet *ResultFieldSet
	fieldValues    *FieldValues
	attrs          partAttributes
}

func newTableResult(s *Session, size int) *TableResult {
	return &TableResult{
		resultFieldSet: newResultFieldSet(size),
		fieldValues:    newFieldValues(),
	}
}

// ID returns the resultset id.
func (r *TableResult) ID() uint64 {
	return r.id
}

// FieldSet returns the field metadata of the table.
func (r *TableResult) FieldSet() *ResultFieldSet {
	return r.resultFieldSet
}

// FieldValues returns the field values (fetched resultset part) of the table.
func (r *TableResult) FieldValues() *FieldValues {
	return r.fieldValues
}

// Attrs returns the PartAttributes interface of the fetched resultset part.
func (r *TableResult) Attrs() PartAttributes {
	return r.attrs
}
