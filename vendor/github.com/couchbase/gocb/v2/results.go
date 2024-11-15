package gocb

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/couchbase/gocbcore/v10/memd"
)

// Result is the base type for the return types of operations
type Result struct {
	cas Cas
}

// Cas returns the cas of the result.
func (d *Result) Cas() Cas {
	return d.cas
}

// GetResult is the return type of Get operations.
type GetResult struct {
	Result
	transcoder Transcoder
	flags      uint32
	contents   []byte
	expiryTime *time.Time
}

// Content assigns the value of the result into the valuePtr using default decoding.
func (d *GetResult) Content(valuePtr interface{}) error {
	return d.transcoder.Decode(d.contents, d.flags, valuePtr)
}

// Expiry returns the expiry value for the result if it available.  Note that a nil
// pointer indicates that the Expiry was not fetched, while a valid pointer to a zero
// Duration indicates that the document will never expire.
// Deprecated: Use ExpiryTime instead.
func (d *GetResult) Expiry() *time.Duration {
	if d.expiryTime == nil {
		return nil
	}

	var t time.Duration
	if !d.expiryTime.IsZero() {
		t = time.Until(*d.expiryTime)
	}

	return &t
}

// ExpiryTime returns the expiry time for the result if it available.
// This function will return a zero time if the value either was not fetched or the
// document does not have an expiry time.
func (d *GetResult) ExpiryTime() time.Time {
	if d.expiryTime == nil {
		return time.Time{}
	}

	return *d.expiryTime
}

func (d *GetResult) fromFullProjection(ops []LookupInSpec, result *LookupInResult, fields []string) error {
	if len(fields) == 0 {
		// This is a special case where user specified a full doc fetch with expiration.
		d.contents = result.contents[0].data
		return nil
	}

	if len(result.contents) != 1 {
		return makeInvalidArgumentsError("fromFullProjection should only be called with 1 subdoc result")
	}

	resultContent := result.contents[0]
	if resultContent.err != nil {
		return resultContent.err
	}

	var content map[string]interface{}
	err := json.Unmarshal(resultContent.data, &content)
	if err != nil {
		return err
	}

	newContent := make(map[string]interface{})
	for _, field := range fields {
		parts := d.pathParts(field)
		d.set(parts, newContent, content[field])
	}

	bytes, err := json.Marshal(newContent)
	if err != nil {
		return wrapError(err, "could not marshal result contents")
	}
	d.contents = bytes

	return nil
}

func (d *GetResult) fromSubDoc(ops []LookupInSpec, result *LookupInResult) error {
	content := make(map[string]interface{})

	for i, op := range ops {
		err := result.contents[i].err
		if err != nil {
			logDebugf("Omitting %s from result due to error: %s", op.path, err)
			continue
		}

		parts := d.pathParts(op.path)
		d.set(parts, content, result.contents[i].data)
	}

	bytes, err := json.Marshal(content)
	if err != nil {
		return wrapError(err, "could not marshal result contents")
	}
	d.contents = bytes

	return nil
}

type subdocPath struct {
	path    string
	isArray bool
}

func (d *GetResult) pathParts(pathStr string) []subdocPath {
	pathLen := len(pathStr)
	var elemIdx int
	var i int
	var paths []subdocPath

	for i < pathLen {
		ch := pathStr[i]
		i++

		if ch == '[' {
			// opening of an array
			isArr := false
			arrayStart := i

			for i < pathLen {
				arrCh := pathStr[i]
				if arrCh == ']' {
					isArr = true
					i++
					break
				} else if arrCh == '.' {
					i++
					break
				}
				i++
			}

			if isArr {
				paths = append(paths, subdocPath{path: pathStr[elemIdx : arrayStart-1], isArray: true})
			} else {
				paths = append(paths, subdocPath{path: pathStr[elemIdx:i], isArray: false})
			}
			elemIdx = i

			if i < pathLen && pathStr[i] == '.' {
				i++
				elemIdx = i
			}
		} else if ch == '.' {
			paths = append(paths, subdocPath{path: pathStr[elemIdx : i-1]})
			elemIdx = i
		}
	}

	if elemIdx != i {
		// this should only ever be an object as an array would have ended in [...]
		paths = append(paths, subdocPath{path: pathStr[elemIdx:i]})
	}

	return paths
}

func (d *GetResult) set(paths []subdocPath, content interface{}, value interface{}) interface{} {
	path := paths[0]
	if len(paths) == 1 {
		if path.isArray {
			arr := make([]interface{}, 0)
			arr = append(arr, value)
			if _, ok := content.(map[string]interface{}); ok {
				content.(map[string]interface{})[path.path] = arr
			} else if _, ok := content.([]interface{}); ok {
				content = append(content.([]interface{}), arr)
			} else {
				logErrorf("Projections encountered a non-array or object content assigning an array")
			}
		} else {
			if _, ok := content.([]interface{}); ok {
				elem := make(map[string]interface{})
				elem[path.path] = value
				content = append(content.([]interface{}), elem)
			} else {
				content.(map[string]interface{})[path.path] = value
			}
		}
		return content
	}

	if path.isArray {
		if _, ok := content.([]interface{}); ok {
			var m []interface{}
			content = append(content.([]interface{}), d.set(paths[1:], m, value))
			return content
		} else if cMap, ok := content.(map[string]interface{}); ok {
			cMap[path.path] = make([]interface{}, 0)
			cMap[path.path] = d.set(paths[1:], cMap[path.path], value)
			return content

		} else {
			logErrorf("Projections encountered a non-array or object content assigning an array")
		}
	} else {
		if arr, ok := content.([]interface{}); ok {
			m := make(map[string]interface{})
			m[path.path] = make(map[string]interface{})
			content = append(arr, m)
			d.set(paths[1:], m[path.path], value)
			return content
		}
		cMap, ok := content.(map[string]interface{})
		if !ok {
			// this isn't possible but the linter won't play nice without it
			logErrorf("Failed to assert projection content to a map")
		}
		cMap[path.path] = make(map[string]interface{})
		return d.set(paths[1:], cMap[path.path], value)
	}

	return content
}

// LookupInResult is the return type for LookupIn.
type LookupInResult struct {
	Result
	contents []lookupInPartial
}

type lookupInPartial struct {
	data json.RawMessage
	err  error
	op   memd.SubDocOpType
}

func (pr *lookupInPartial) as(valuePtr interface{}) error {
	if pr.err != nil {
		return pr.err
	}

	if valuePtr == nil {
		return nil
	}

	if valuePtr, ok := valuePtr.(*[]byte); ok {
		*valuePtr = pr.data
		return nil
	}

	return json.Unmarshal(pr.data, valuePtr)
}

func (pr *lookupInPartial) exists() bool {
	if pr.op == memd.SubDocOpExists {
		// For exists spec we need to try to parse the result as a bool
		// and return the bool value.
		var exists bool
		err := pr.as(&exists)
		if err != nil {
			return false
		}

		return exists
	}

	// For any other spec we can just check for the presence of an error.
	return pr.err == nil
}

// ContentAt retrieves the value of the operation by its index. The index is the position of
// the operation as it was added to the builder.
func (lir *LookupInResult) ContentAt(idx uint, valuePtr interface{}) error {
	if idx >= uint(len(lir.contents)) {
		return makeInvalidArgumentsError("invalid index")
	}
	return lir.contents[idx].as(valuePtr)
}

// Exists verifies that the item at idx exists.
func (lir *LookupInResult) Exists(idx uint) bool {
	if idx >= uint(len(lir.contents)) {
		return false
	}
	return lir.contents[idx].exists()
}

// LookupInAllReplicasResult represents the results of a LookupInAllReplicas operation.
type LookupInAllReplicasResult struct {
	res *coreReplicasResult
}

// Next fetches the next replica result.
func (r *LookupInAllReplicasResult) Next() *LookupInReplicaResult {
	res := r.res.Next()
	if res == nil {
		return nil
	}
	return res.(*LookupInReplicaResult)
}

// Close cancels all remaining get replica requests.
func (r *LookupInAllReplicasResult) Close() error {
	return r.res.Close()
}

// LookupInReplicaResult is the return type of LookupInReplica operations.
type LookupInReplicaResult struct {
	*LookupInResult
	isReplica bool
}

// IsReplica returns whether or not this result came from a replica server.
func (r *LookupInReplicaResult) IsReplica() bool {
	return r.isReplica
}

// ExistsResult is the return type of Exist operations.
type ExistsResult struct {
	Result
	docExists bool
}

// Exists returns whether or not the document exists.
func (d *ExistsResult) Exists() bool {
	return d.docExists
}

// MutationResult is the return type of any store related operations. It contains Cas and mutation tokens.
type MutationResult struct {
	Result
	mt *MutationToken
}

// MutationToken returns the mutation token belonging to an operation.
func (mr MutationResult) MutationToken() *MutationToken {
	return mr.mt
}

// MutateInResult is the return type of any mutate in related operations.
// It contains Cas, mutation tokens and any returned content.
type MutateInResult struct {
	MutationResult
	contents []mutateInPartial
}

type mutateInPartial struct {
	data json.RawMessage
}

func (pr *mutateInPartial) as(valuePtr interface{}) error {
	if valuePtr == nil {
		return nil
	}

	if valuePtr, ok := valuePtr.(*[]byte); ok {
		*valuePtr = pr.data
		return nil
	}

	return json.Unmarshal(pr.data, valuePtr)
}

// ContentAt retrieves the value of the operation by its index. The index is the position of
// the operation as it was added to the builder.
func (mir MutateInResult) ContentAt(idx uint, valuePtr interface{}) error {
	return mir.contents[idx].as(valuePtr)
}

// CounterResult is the return type of counter operations.
type CounterResult struct {
	MutationResult
	content uint64
}

// MutationToken returns the mutation token belonging to an operation.
func (mr CounterResult) MutationToken() *MutationToken {
	return mr.mt
}

// Cas returns the Cas value for a document following an operation.
func (mr CounterResult) Cas() Cas {
	return mr.cas
}

// Content returns the new value for the counter document.
func (mr CounterResult) Content() uint64 {
	return mr.content
}

// GetReplicaResult is the return type of GetReplica operations.
type GetReplicaResult struct {
	GetResult
	isReplica bool
}

// IsReplica returns whether or not this result came from a replica server.
func (r *GetReplicaResult) IsReplica() bool {
	return r.isReplica
}

// ScanResult is the return type of Scan operations.
type ScanResult struct {
	resultChan chan *ScanResultItem
	cancelFn   func(error)
	err        error
	errLock    sync.Mutex

	limit    uint64
	numItems uint64

	peeked unsafe.Pointer
}

func (sr *ScanResult) setErr(err error) {
	sr.errLock.Lock()
	sr.err = err
	sr.errLock.Unlock()
}

// Next returns the next item on the stream, if there are no items remaining then nil is returned.
func (sr *ScanResult) Next() *ScanResultItem {
	peeked := atomic.SwapPointer(&sr.peeked, nil)
	if peeked != nil {
		atomic.AddUint64(&sr.numItems, 1)
		return (*ScanResultItem)(peeked)
	}

	item, more := <-sr.resultChan
	if !more {
		return nil
	}

	// If we're doing a sampling scan then we need to only write data into the channel
	// if we haven't seen the number of items that the user requested. Otherwise
	// we need to cancel the streams
	numItems := atomic.AddUint64(&sr.numItems, 1)
	if sr.limit == 0 || numItems <= sr.limit {
		return item
	}

	sr.cancelFn(nil)
	return nil
}

// Err returns any errors that have occurred on the stream.
func (sr *ScanResult) Err() error {
	sr.errLock.Lock()
	err := sr.err
	sr.errLock.Unlock()

	return err
}

// Close cancels the stream, returning any errors that occurred during reading the results.
func (sr *ScanResult) Close() error {
	sr.errLock.Lock()
	err := sr.err
	sr.errLock.Unlock()

	if err != nil {
		return err
	}

	sr.cancelFn(ErrRequestCanceled)

	return nil
}

// ScanResultItem represents an item that is returning on the stream from a Scan operation.
type ScanResultItem struct {
	Result
	transcoder Transcoder
	id         string
	flags      uint32
	contents   []byte
	expiryTime time.Time
	keysOnly   bool
}

// IDOnly returns whether the scan generating this item was made with IDsOnly set.
func (sri *ScanResultItem) IDOnly() bool {
	return sri.keysOnly
}

// ID returns the id of the item.
func (sri *ScanResultItem) ID() string {
	return sri.id
}

// Cas returns the Cas of the item.
func (sri *ScanResultItem) Cas() Cas {
	return sri.cas
}

// Content assigns the value of the result into the valuePtr using default decoding.
// If IDsOnly was set on the ScanOptions then this will return an error.
func (sri *ScanResultItem) Content(valuePtr interface{}) error {
	if sri.keysOnly {
		return makeInvalidArgumentsError("scan was called with IDsOnly set to true, content can never be set")
	}
	return sri.transcoder.Decode(sri.contents, sri.flags, valuePtr)
}

// ExpiryTime returns the expiry time for the result if available.
// This function will return a zero time if the value either was not fetched or the
// document does not have an expiry time.
func (sri *ScanResultItem) ExpiryTime() time.Time {
	return sri.expiryTime
}
