package gocb

import (
	"encoding/json"
	"errors"
)

// ViewMetaData provides access to the meta-data properties of a view query result.
type ViewMetaData struct {
	TotalRows uint64
	Debug     interface{}
}

func (meta *ViewMetaData) fromData(data jsonViewResponse) error {
	meta.TotalRows = data.TotalRows
	meta.Debug = data.DebugInfo

	return nil
}

// ViewRow represents a single row returned from a view query.
type ViewRow struct {
	ID         string
	keyBytes   []byte
	valueBytes []byte
}

// Key returns the key associated with this view row.
func (vr *ViewRow) Key(valuePtr interface{}) error {
	return json.Unmarshal(vr.keyBytes, valuePtr)
}

// Value returns the value associated with this view row.
func (vr *ViewRow) Value(valuePtr interface{}) error {
	return json.Unmarshal(vr.valueBytes, valuePtr)
}

// ViewResultRaw provides raw access to views data.
// VOLATILE: This API is subject to change at any time.
type ViewResultRaw struct {
	reader viewRowReader
}

// NextBytes returns the next row as bytes.
func (vrr *ViewResultRaw) NextBytes() []byte {
	return vrr.reader.NextRow()
}

// Err returns any errors that have occurred on the stream
func (vrr *ViewResultRaw) Err() error {
	err := vrr.reader.Err()
	if err != nil {
		return maybeEnhanceViewError(err)
	}

	return nil
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (vrr *ViewResultRaw) Close() error {
	err := vrr.reader.Close()
	if err != nil {
		return maybeEnhanceViewError(err)
	}

	return nil
}

// MetaData returns any meta-data that was available from this query as bytes.
func (vrr *ViewResultRaw) MetaData() ([]byte, error) {
	return vrr.reader.MetaData()
}

// ViewResult implements an iterator interface which can be used to iterate over the rows of the query results.
type ViewResult struct {
	reader viewRowReader

	currentRow ViewRow
	jsonErr    error
}

func newViewResult(reader viewRowReader) *ViewResult {
	return &ViewResult{
		reader: reader,
	}
}

// Raw returns a ViewResultRaw which can be used to access the raw byte data from view queries.
// Calling this function invalidates the underlying ViewResult which will no longer be able to be used.
// VOLATILE: This API is subject to change at any time.
func (r *ViewResult) Raw() *ViewResultRaw {
	vr := &ViewResultRaw{
		reader: r.reader,
	}

	r.reader = nil
	return vr
}

// Next assigns the next result from the results into the value pointer, returning whether the read was successful.
func (r *ViewResult) Next() bool {
	if r.reader == nil {
		return false
	}

	rowBytes := r.reader.NextRow()
	if rowBytes == nil {
		return false
	}

	r.currentRow = ViewRow{}

	var rowData jsonViewRow
	if err := json.Unmarshal(rowBytes, &rowData); err != nil {
		// This should never happen but if it does then lets store it in a best efforts basis and maybe the next
		// row will be ok. We can then return this from .Err().
		r.jsonErr = err
		return true
	}

	r.currentRow.ID = rowData.ID
	r.currentRow.keyBytes = rowData.Key
	r.currentRow.valueBytes = rowData.Value

	return true
}

// Row returns the contents of the current row.
func (r *ViewResult) Row() ViewRow {
	if r.reader == nil {
		return ViewRow{}
	}

	return r.currentRow
}

// Err returns any errors that have occurred on the stream
func (r *ViewResult) Err() error {
	if r.reader == nil {
		return errors.New("result object is no longer valid")
	}

	err := r.reader.Err()
	if err != nil {
		return maybeEnhanceViewError(err)
	}
	// This is an error from json unmarshal so no point in trying to enhance it.
	return r.jsonErr
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (r *ViewResult) Close() error {
	if r.reader == nil {
		return r.Err()
	}

	err := r.reader.Close()
	if err != nil {
		return maybeEnhanceViewError(err)
	}

	return nil
}

// MetaData returns any meta-data that was available from this query.  Note that
// the meta-data will only be available once the object has been closed (either
// implicitly or explicitly).
func (r *ViewResult) MetaData() (*ViewMetaData, error) {
	if r.reader == nil {
		return nil, r.Err()
	}

	metaDataBytes, err := r.reader.MetaData()
	if err != nil {
		return nil, err
	}

	var jsonResp jsonViewResponse
	err = json.Unmarshal(metaDataBytes, &jsonResp)
	if err != nil {
		return nil, err
	}

	var metaData ViewMetaData
	err = metaData.fromData(jsonResp)
	if err != nil {
		return nil, err
	}

	return &metaData, nil
}

// ViewQuery performs a view query and returns a list of rows or an error.
func (b *Bucket) ViewQuery(designDoc string, viewName string, opts *ViewOptions) (*ViewResult, error) {
	return autoOpControl(b.viewController(), func(provider viewProvider) (*ViewResult, error) {
		if opts == nil {
			opts = &ViewOptions{}
		}

		return provider.ViewQuery(designDoc, viewName, opts)
	})
}
