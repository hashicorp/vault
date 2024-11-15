package gocb

import (
	"encoding/json"
)

// TransactionQueryResult allows access to the results of a query.
type TransactionQueryResult struct {
	results  []json.RawMessage
	idx      int
	rowBytes json.RawMessage

	metadata *QueryMetaData
	endpoint string
}

func newTransactionQueryResult(results []json.RawMessage, meta *QueryMetaData, endpoint string) *TransactionQueryResult {
	return &TransactionQueryResult{
		results:  results,
		metadata: meta,
		endpoint: endpoint,
	}
}

// Next assigns the next result from the results into the value pointer, returning whether the read was successful.
func (r *TransactionQueryResult) Next() bool {
	if r.idx >= len(r.results) {
		return false
	}

	r.rowBytes = r.results[r.idx]
	r.idx++

	return true
}

// Row returns the contents of the current row
func (r *TransactionQueryResult) Row(valuePtr interface{}) error {
	if r.rowBytes == nil {
		return ErrNoResult
	}

	if bytesPtr, ok := valuePtr.(*json.RawMessage); ok {
		*bytesPtr = r.rowBytes
		return nil
	}

	return json.Unmarshal(r.rowBytes, valuePtr)
}

// One assigns the first value from the results into the value pointer.
func (r *TransactionQueryResult) One(valuePtr interface{}) error {
	// Prime the row
	if !r.Next() {
		return ErrNoResult
	}

	err := r.Row(valuePtr)
	if err != nil {
		return err
	}

	return nil
}

// MetaData returns any meta-data that was available from this query.  Note that
// the meta-data will only be available once the object has been closed (either
// implicitly or explicitly).
func (r *TransactionQueryResult) MetaData() (*QueryMetaData, error) {
	return r.metadata, nil
}
