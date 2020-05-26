package gocb

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v9"
	"github.com/pkg/errors"
)

type jsonViewResponse struct {
	TotalRows uint64      `json:"total_rows,omitempty"`
	DebugInfo interface{} `json:"debug_info,omitempty"`
}

type jsonViewRow struct {
	ID    string          `json:"id"`
	Key   json.RawMessage `json:"key"`
	Value json.RawMessage `json:"value"`
}

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

type viewRowReader interface {
	NextRow() []byte
	Err() error
	MetaData() ([]byte, error)
	Close() error
}

// ViewResult implements an iterator interface which can be used to iterate over the rows of the query results.
type ViewResult struct {
	reader viewRowReader

	currentRow ViewRow
}

func newViewResult(reader viewRowReader) *ViewResult {
	return &ViewResult{
		reader: reader,
	}
}

// Next assigns the next result from the results into the value pointer, returning whether the read was successful.
func (r *ViewResult) Next() bool {
	rowBytes := r.reader.NextRow()
	if rowBytes == nil {
		return false
	}

	r.currentRow = ViewRow{}

	var rowData jsonViewRow
	if err := json.Unmarshal(rowBytes, &rowData); err == nil {
		r.currentRow.ID = rowData.ID
		r.currentRow.keyBytes = rowData.Key
		r.currentRow.valueBytes = rowData.Value
	}

	return true
}

// Row returns the contents of the current row.
func (r *ViewResult) Row() ViewRow {
	return r.currentRow
}

// Err returns any errors that have occurred on the stream
func (r *ViewResult) Err() error {
	return r.reader.Err()
}

// Close marks the results as closed, returning any errors that occurred during reading the results.
func (r *ViewResult) Close() error {
	return r.reader.Close()
}

// MetaData returns any meta-data that was available from this query.  Note that
// the meta-data will only be available once the object has been closed (either
// implicitly or explicitly).
func (r *ViewResult) MetaData() (*ViewMetaData, error) {
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
	if opts == nil {
		opts = &ViewOptions{}
	}

	span := b.tracer.StartSpan("ViewQuery", opts.parentSpan).
		SetTag("couchbase.service", "view")
	defer span.Finish()

	designDoc = b.maybePrefixDevDocument(opts.Namespace, designDoc)

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = b.timeoutsConfig.ViewTimeout
	}
	deadline := time.Now().Add(timeout)

	retryWrapper := b.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryWrapper = newRetryStrategyWrapper(opts.RetryStrategy)
	}

	urlValues, err := opts.toURLValues()
	if err != nil {
		return nil, errors.Wrap(err, "could not parse query options")
	}

	return b.execViewQuery(span.Context(), "_view", designDoc, viewName, *urlValues, deadline, retryWrapper)
}

func (b *Bucket) execViewQuery(
	span requestSpanContext,
	viewType, ddoc, viewName string,
	options url.Values,
	deadline time.Time,
	wrapper *retryStrategyWrapper,
) (*ViewResult, error) {
	cli := b.getCachedClient()
	provider, err := cli.getViewProvider()
	if err != nil {
		return nil, ViewError{
			InnerError:         wrapError(err, "failed to get query provider"),
			DesignDocumentName: ddoc,
			ViewName:           viewName,
		}
	}

	res, err := provider.ViewQuery(gocbcore.ViewQueryOptions{
		DesignDocumentName: ddoc,
		ViewType:           viewType,
		ViewName:           viewName,
		Options:            options,
		RetryStrategy:      wrapper,
		Deadline:           deadline,
		TraceContext:       span,
	})
	if err != nil {
		return nil, maybeEnhanceViewError(err)
	}

	return newViewResult(res), nil
}

func (b *Bucket) maybePrefixDevDocument(namespace DesignDocumentNamespace, ddoc string) string {
	designDoc := ddoc
	if namespace == DesignDocumentNamespaceProduction {
		designDoc = strings.TrimPrefix(ddoc, "dev_")
	} else {
		if !strings.HasPrefix(ddoc, "dev_") {
			designDoc = "dev_" + ddoc
		}
	}

	return designDoc
}
