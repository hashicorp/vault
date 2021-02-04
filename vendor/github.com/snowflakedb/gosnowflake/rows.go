// Copyright (c) 2017-2019 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"bufio"
	"compress/gzip"
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/apache/arrow/go/arrow/ipc"
	"github.com/apache/arrow/go/arrow/memory"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	headerSseCAlgorithm = "x-amz-server-side-encryption-customer-algorithm"
	headerSseCKey       = "x-amz-server-side-encryption-customer-key"
	headerSseCAes       = "AES256"
)

var (
	// MaxChunkDownloadWorkers specifies the maximum number of goroutines used to download chunks
	MaxChunkDownloadWorkers = 10

	// CustomJSONDecoderEnabled has the chunk downloader use the custom JSON decoder to reduce memory footprint.
	CustomJSONDecoderEnabled = false
)

var (
	maxChunkDownloaderErrorCounter = 5
)

type snowflakeRows struct {
	sc              *snowflakeConn
	RowType         []execResponseRowType
	ChunkDownloader *snowflakeChunkDownloader
	queryID         string
}

func (rows *snowflakeRows) Close() (err error) {
	glog.V(2).Infoln("Rows.Close")
	return nil
}

type snowflakeValue interface{}

type chunkRowType struct {
	RowSet   []*string
	ArrowRow []snowflakeValue
}

type rowSetType struct {
	RowType      []execResponseRowType
	JSON         [][]*string
	RowSetBase64 string
}

type chunkError struct {
	Index int
	Error error
}

type snowflakeChunkDownloader struct {
	sc                 *snowflakeConn
	ctx                context.Context
	Total              int64
	TotalRowIndex      int64
	CellCount          int
	CurrentChunk       []chunkRowType
	CurrentChunkIndex  int
	CurrentChunkSize   int
	ChunksMutex        *sync.Mutex
	ChunkMetas         []execResponseChunk
	Chunks             map[int][]chunkRowType
	ChunksChan         chan int
	ChunksError        chan *chunkError
	ChunksErrorCounter int
	ChunksFinalErrors  []*chunkError
	Qrmk               string
	QueryResultFormat  string
	RowSet             rowSetType
	ChunkHeader        map[string]string
	CurrentIndex       int
	FuncDownload       func(context.Context, *snowflakeChunkDownloader, int)
	FuncDownloadHelper func(context.Context, *snowflakeChunkDownloader, int) error
	FuncGet            func(context.Context, *snowflakeChunkDownloader, string, map[string]string, time.Duration) (*http.Response, error)
	DoneDownloadCond   *sync.Cond
	NextDownloader     *snowflakeChunkDownloader
}

// ColumnTypeDatabaseTypeName returns the database column name.
func (rows *snowflakeRows) ColumnTypeDatabaseTypeName(index int) string {
	return strings.ToUpper(rows.RowType[index].Type)
}

// ColumnTypeLength returns the length of the column
func (rows *snowflakeRows) ColumnTypeLength(index int) (length int64, ok bool) {
	if index < 0 || index > len(rows.RowType) {
		return 0, false
	}
	switch rows.RowType[index].Type {
	case "text", "variant", "object", "array", "binary":
		return rows.RowType[index].Length, true
	}
	return 0, false
}

func (rows *snowflakeRows) ColumnTypeNullable(index int) (nullable, ok bool) {
	if index < 0 || index > len(rows.RowType) {
		return false, false
	}
	return rows.RowType[index].Nullable, true
}

func (rows *snowflakeRows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	if index < 0 || index > len(rows.RowType) {
		return 0, 0, false
	}
	switch rows.RowType[index].Type {
	case "fixed":
		return rows.RowType[index].Precision, rows.RowType[index].Scale, true
	case "time":
		return rows.RowType[index].Scale, 0, true
	case "timestamp":
		return rows.RowType[index].Scale, 0, true
	}
	return 0, 0, false
}

func (rows *snowflakeRows) Columns() []string {
	glog.V(3).Infoln("Rows.Columns")
	ret := make([]string, len(rows.RowType))
	for i, n := 0, len(rows.RowType); i < n; i++ {
		ret[i] = rows.RowType[i].Name
	}
	return ret
}

func (rows *snowflakeRows) ColumnTypeScanType(index int) reflect.Type {
	return snowflakeTypeToGo(rows.RowType[index].Type, rows.RowType[index].Scale)
}

func (rows *snowflakeRows) QueryID() string {
	return rows.queryID
}

func (rows *snowflakeRows) Next(dest []driver.Value) (err error) {
	row, err := rows.ChunkDownloader.Next()
	if err != nil {
		// includes io.EOF
		if err == io.EOF {
			rows.ChunkDownloader.Chunks = nil // detach all chunks. No way to go backward without reinitialize it.
		}
		return err
	}

	if rows.ChunkDownloader.QueryResultFormat == arrowFormat {
		for i, n := 0, len(row.ArrowRow); i < n; i++ {
			dest[i] = row.ArrowRow[i]
		}
	} else {
		for i, n := 0, len(row.RowSet); i < n; i++ {
			// could move to chunk downloader so that each go routine
			// can convert data
			err := stringToValue(&dest[i], rows.RowType[i], row.RowSet[i])
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (rows *snowflakeRows) HasNextResultSet() bool {
	if len(rows.ChunkDownloader.ChunkMetas) == 0 && rows.ChunkDownloader.NextDownloader == nil {
		return false // no extra chunk
	}
	return rows.ChunkDownloader.hasNextResultSet()
}

func (rows *snowflakeRows) NextResultSet() error {
	if len(rows.ChunkDownloader.ChunkMetas) == 0 {
		if rows.ChunkDownloader.NextDownloader == nil {
			return io.EOF
		}
		rows.ChunkDownloader = rows.ChunkDownloader.NextDownloader
		rows.ChunkDownloader.start()
	}
	return rows.ChunkDownloader.nextResultSet()
}

func (scd *snowflakeChunkDownloader) totalUncompressedSize() (acc int64) {
	for _, c := range scd.ChunkMetas {
		acc += c.UncompressedSize
	}
	return
}

func (scd *snowflakeChunkDownloader) hasNextResultSet() bool {
	// next result set exists if current chunk has remaining result sets or there is another downloader
	return scd.CurrentChunkIndex < len(scd.ChunkMetas) || scd.NextDownloader != nil
}

func (scd *snowflakeChunkDownloader) nextResultSet() error {
	// no error at all times as the next chunk/resultset is automatically read
	if scd.CurrentChunkIndex < len(scd.ChunkMetas) {
		return nil
	}
	return io.EOF
}

func (scd *snowflakeChunkDownloader) start() error {
	scd.CurrentChunkSize = len(scd.RowSet.JSON) // cache the size
	scd.CurrentIndex = -1                       // initial chunks idx
	scd.CurrentChunkIndex = -1                  // initial chunk

	scd.CurrentChunk = make([]chunkRowType, scd.CurrentChunkSize)
	populateJSONRowSet(scd.CurrentChunk, scd.RowSet.JSON)

	if scd.QueryResultFormat == arrowFormat && scd.RowSet.RowSetBase64 != "" {
		// if the rowsetbase64 retrieved from the server is empty, move on to downloading chunks
		var err error
		firstArrowChunk := buildFirstArrowChunk(scd.RowSet.RowSetBase64)
		scd.CurrentChunk, err = firstArrowChunk.decodeArrowChunk(scd.RowSet.RowType)
		scd.CurrentChunkSize = firstArrowChunk.rowCount
		if err != nil {
			return err
		}
	}

	// start downloading chunks if exists
	chunkMetaLen := len(scd.ChunkMetas)
	if chunkMetaLen > 0 {
		glog.V(2).Infof("MaxChunkDownloadWorkers: %v", MaxChunkDownloadWorkers)
		glog.V(2).Infof("chunks: %v, total bytes: %d", chunkMetaLen, scd.totalUncompressedSize())
		scd.ChunksMutex = &sync.Mutex{}
		scd.DoneDownloadCond = sync.NewCond(scd.ChunksMutex)
		scd.Chunks = make(map[int][]chunkRowType)
		scd.ChunksChan = make(chan int, chunkMetaLen)
		scd.ChunksError = make(chan *chunkError, MaxChunkDownloadWorkers)
		for i := 0; i < chunkMetaLen; i++ {
			glog.V(2).Infof("add chunk to channel ChunksChan: %v", i+1)
			scd.ChunksChan <- i
		}
		for i := 0; i < intMin(MaxChunkDownloadWorkers, chunkMetaLen); i++ {
			scd.schedule()
		}
	}
	return nil
}

func (scd *snowflakeChunkDownloader) schedule() {
	select {
	case nextIdx := <-scd.ChunksChan:
		glog.V(2).Infof("schedule chunk: %v", nextIdx+1)
		go scd.FuncDownload(scd.ctx, scd, nextIdx)
	default:
		// no more download
		glog.V(2).Info("no more download")
	}
}

func (scd *snowflakeChunkDownloader) checkErrorRetry() (err error) {
	select {
	case errc := <-scd.ChunksError:
		if scd.ChunksErrorCounter < maxChunkDownloaderErrorCounter && errc.Error != context.Canceled {
			// add the index to the chunks channel so that the download will be retried.
			go scd.FuncDownload(scd.ctx, scd, errc.Index)
			scd.ChunksErrorCounter++
			glog.V(2).Infof("chunk idx: %v, err: %v. retrying (%v/%v)...",
				errc.Index, errc.Error, scd.ChunksErrorCounter, maxChunkDownloaderErrorCounter)
		} else {
			scd.ChunksFinalErrors = append(scd.ChunksFinalErrors, errc)
			glog.V(2).Infof("chunk idx: %v, err: %v. no further retry", errc.Index, errc.Error)
			return errc.Error
		}
	default:
		glog.V(2).Info("no error is detected.")
	}
	return nil
}
func (scd *snowflakeChunkDownloader) Next() (chunkRowType, error) {
	for {
		scd.CurrentIndex++
		if scd.CurrentIndex < scd.CurrentChunkSize {
			return scd.CurrentChunk[scd.CurrentIndex], nil
		}
		scd.CurrentChunkIndex++ // next chunk
		scd.CurrentIndex = -1   // reset
		if scd.CurrentChunkIndex >= len(scd.ChunkMetas) {
			break
		}

		scd.ChunksMutex.Lock()
		if scd.CurrentChunkIndex > 1 {
			scd.Chunks[scd.CurrentChunkIndex-1] = nil // detach the previously used chunk
		}

		for scd.Chunks[scd.CurrentChunkIndex] == nil {
			glog.V(2).Infof("waiting for chunk idx: %v/%v",
				scd.CurrentChunkIndex+1, len(scd.ChunkMetas))

			err := scd.checkErrorRetry()
			if err != nil {
				scd.ChunksMutex.Unlock()
				return chunkRowType{}, err
			}

			// wait for chunk downloader goroutine to broadcast the event,
			// 1) one chunk download finishes or 2) an error occurs.
			scd.DoneDownloadCond.Wait()
		}
		glog.V(2).Infof("ready: chunk %v", scd.CurrentChunkIndex+1)
		scd.CurrentChunk = scd.Chunks[scd.CurrentChunkIndex]
		scd.ChunksMutex.Unlock()
		scd.CurrentChunkSize = len(scd.CurrentChunk)

		// kick off the next download
		scd.schedule()
	}

	glog.V(2).Infof("no more data")
	if len(scd.ChunkMetas) > 0 {
		close(scd.ChunksError)
		close(scd.ChunksChan)
	}
	return chunkRowType{}, io.EOF
}

func getChunk(
	ctx context.Context,
	scd *snowflakeChunkDownloader,
	fullURL string,
	headers map[string]string,
	timeout time.Duration) (
	*http.Response, error) {
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	return newRetryHTTP(ctx, scd.sc.rest.Client, http.NewRequest, u, headers, timeout).execute()
}

/* largeResultSetReader is a reader that wraps the large result set with leading and tailing brackets. */
type largeResultSetReader struct {
	status int
	body   io.Reader
}

func (r *largeResultSetReader) Read(p []byte) (n int, err error) {
	if r.status == 0 {
		p[0] = 0x5b // initial 0x5b ([)
		r.status = 1
		return 1, nil
	}
	if r.status == 1 {
		var len int
		len, err = r.body.Read(p)
		if err == io.EOF {
			r.status = 2
			return len, nil
		}
		if err != nil {
			return 0, err
		}
		return len, nil
	}
	if r.status == 2 {
		p[0] = 0x5d // tail 0x5d (])
		r.status = 3
		return 1, nil
	}
	// ensure no data and EOF
	return 0, io.EOF
}

func downloadChunk(ctx context.Context, scd *snowflakeChunkDownloader, idx int) {
	glog.V(2).Infof("download start chunk: %v", idx+1)
	defer scd.DoneDownloadCond.Broadcast()

	if err := scd.FuncDownloadHelper(ctx, scd, idx); err != nil {
		glog.V(1).Infof(
			"failed to extract HTTP response body. URL: %v, err: %v", scd.ChunkMetas[idx].URL, err)
		glog.Flush()
		scd.ChunksError <- &chunkError{Index: idx, Error: err}
	} else if scd.ctx.Err() == context.Canceled || scd.ctx.Err() == context.DeadlineExceeded {
		scd.ChunksError <- &chunkError{Index: idx, Error: scd.ctx.Err()}
	}
}

func downloadChunkHelper(ctx context.Context, scd *snowflakeChunkDownloader, idx int) error {
	headers := make(map[string]string)
	if len(scd.ChunkHeader) > 0 {
		glog.V(2).Info("chunk header is provided.")
		for k, v := range scd.ChunkHeader {
			headers[k] = v
		}
	} else {
		headers[headerSseCAlgorithm] = headerSseCAes
		headers[headerSseCKey] = scd.Qrmk
	}

	resp, err := scd.FuncGet(ctx, scd, scd.ChunkMetas[idx].URL, headers, scd.sc.rest.RequestTimeout)
	if err != nil {
		return err
	}
	bufStream := bufio.NewReader(resp.Body)
	defer resp.Body.Close()
	glog.V(2).Infof("response returned chunk: %v, resp: %v", idx+1, resp)
	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(bufStream)
		if err != nil {
			return err
		}
		glog.V(1).Infof("HTTP: %v, URL: %v, Body: %v", resp.StatusCode, scd.ChunkMetas[idx].URL, b)
		glog.V(1).Infof("Header: %v", resp.Header)
		glog.Flush()
		return &SnowflakeError{
			Number:      ErrFailedToGetChunk,
			SQLState:    SQLStateConnectionFailure,
			Message:     errMsgFailedToGetChunk,
			MessageArgs: []interface{}{idx},
		}
	}
	return decodeChunk(scd, idx, bufStream)
}

func decodeChunk(scd *snowflakeChunkDownloader, idx int, bufStream *bufio.Reader) (err error) {
	gzipMagic, err := bufStream.Peek(2)
	if err != nil {
		return err
	}
	start := time.Now()
	var source io.Reader
	if gzipMagic[0] == 0x1f && gzipMagic[1] == 0x8b {
		// detects and uncompresses Gzip format data
		bufStream0, err := gzip.NewReader(bufStream)
		if err != nil {
			return err
		}
		defer bufStream0.Close()
		source = bufStream0
	} else {
		source = bufStream
	}
	st := &largeResultSetReader{
		status: 0,
		body:   source,
	}
	var respd []chunkRowType
	if scd.QueryResultFormat != arrowFormat {
		var decRespd [][]*string
		if !CustomJSONDecoderEnabled {
			dec := json.NewDecoder(st)
			for {
				if err := dec.Decode(&decRespd); err == io.EOF {
					break
				} else if err != nil {
					return err
				}
			}
		} else {
			decRespd, err = decodeLargeChunk(st, scd.ChunkMetas[idx].RowCount, scd.CellCount)
			if err != nil {
				return err
			}
		}
		respd = make([]chunkRowType, len(decRespd))
		populateJSONRowSet(respd, decRespd)
	} else {
		ipcReader, err := ipc.NewReader(source)
		if err != nil {
			return err
		}
		arc := arrowResultChunk{
			*ipcReader,
			0,
			int(scd.totalUncompressedSize()),
			memory.NewGoAllocator(),
		}
		respd, err = arc.decodeArrowChunk(scd.RowSet.RowType)
		if err != nil {
			return err
		}
	}
	glog.V(2).Infof(
		"decoded %d rows w/ %d bytes in %s (chunk %v)",
		scd.ChunkMetas[idx].RowCount,
		scd.ChunkMetas[idx].UncompressedSize,
		time.Since(start), idx+1,
	)

	scd.ChunksMutex.Lock()
	defer scd.ChunksMutex.Unlock()
	scd.Chunks[idx] = respd
	return nil
}

func populateJSONRowSet(dst []chunkRowType, src [][]*string) {
	// populate string rowset from src to dst's chunkRowType struct's RowSet field
	for i, row := range src {
		dst[i].RowSet = row
	}
}
