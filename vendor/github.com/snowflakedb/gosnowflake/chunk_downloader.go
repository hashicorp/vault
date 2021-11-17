// Copyright (c) 2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/arrow/go/arrow/ipc"
	"github.com/apache/arrow/go/arrow/memory"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type chunkDownloader interface {
	totalUncompressedSize() (acc int64)
	hasNextResultSet() bool
	nextResultSet() error
	start() error
	next() (chunkRowType, error)
	reset()
	getChunkMetas() []execResponseChunk
	getQueryResultFormat() resultFormat
	getRowType() []execResponseRowType
	setNextChunkDownloader(downloader chunkDownloader)
	getNextChunkDownloader() chunkDownloader
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
	CurrentIndex       int
	ChunkHeader        map[string]string
	ChunkMetas         []execResponseChunk
	Chunks             map[int][]chunkRowType
	ChunksChan         chan int
	ChunksError        chan *chunkError
	ChunksErrorCounter int
	ChunksFinalErrors  []*chunkError
	ChunksMutex        *sync.Mutex
	DoneDownloadCond   *sync.Cond
	NextDownloader     chunkDownloader
	Qrmk               string
	QueryResultFormat  string
	RowSet             rowSetType
	FuncDownload       func(context.Context, *snowflakeChunkDownloader, int)
	FuncDownloadHelper func(context.Context, *snowflakeChunkDownloader, int) error
	FuncGet            func(context.Context, *snowflakeChunkDownloader, string, map[string]string, time.Duration) (*http.Response, error)
}

func (scd *snowflakeChunkDownloader) totalUncompressedSize() (acc int64) {
	for _, c := range scd.ChunkMetas {
		acc += c.UncompressedSize
	}
	return
}

func (scd *snowflakeChunkDownloader) hasNextResultSet() bool {
	if len(scd.ChunkMetas) == 0 && scd.NextDownloader == nil {
		return false // no extra chunk
	}
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

	if scd.getQueryResultFormat() == arrowFormat && scd.RowSet.RowSetBase64 != "" {
		// if the rowsetbase64 retrieved from the server is empty, move on to downloading chunks
		var err error
		firstArrowChunk := buildFirstArrowChunk(scd.RowSet.RowSetBase64)
		higherPrecision := higherPrecisionEnabled(scd.ctx)
		scd.CurrentChunk, err = firstArrowChunk.decodeArrowChunk(scd.RowSet.RowType, higherPrecision)
		scd.CurrentChunkSize = firstArrowChunk.rowCount
		if err != nil {
			return err
		}
	}

	// start downloading chunks if exists
	chunkMetaLen := len(scd.ChunkMetas)
	if chunkMetaLen > 0 {
		logger.Debugf("MaxChunkDownloadWorkers: %v", MaxChunkDownloadWorkers)
		logger.Debugf("chunks: %v, total bytes: %d", chunkMetaLen, scd.totalUncompressedSize())
		scd.ChunksMutex = &sync.Mutex{}
		scd.DoneDownloadCond = sync.NewCond(scd.ChunksMutex)
		scd.Chunks = make(map[int][]chunkRowType)
		scd.ChunksChan = make(chan int, chunkMetaLen)
		scd.ChunksError = make(chan *chunkError, MaxChunkDownloadWorkers)
		for i := 0; i < chunkMetaLen; i++ {
			logger.Debugf("add chunk to channel ChunksChan: %v", i+1)
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
		logger.Infof("schedule chunk: %v", nextIdx+1)
		go scd.FuncDownload(scd.ctx, scd, nextIdx)
	default:
		// no more download
		logger.Info("no more download")
	}
}

func (scd *snowflakeChunkDownloader) checkErrorRetry() (err error) {
	select {
	case errc := <-scd.ChunksError:
		if scd.ChunksErrorCounter < maxChunkDownloaderErrorCounter && errc.Error != context.Canceled {
			// add the index to the chunks channel so that the download will be retried.
			go scd.FuncDownload(scd.ctx, scd, errc.Index)
			scd.ChunksErrorCounter++
			logger.Warningf("chunk idx: %v, err: %v. retrying (%v/%v)...",
				errc.Index, errc.Error, scd.ChunksErrorCounter, maxChunkDownloaderErrorCounter)
		} else {
			scd.ChunksFinalErrors = append(scd.ChunksFinalErrors, errc)
			logger.Warningf("chunk idx: %v, err: %v. no further retry", errc.Index, errc.Error)
			return errc.Error
		}
	default:
		logger.Info("no error is detected.")
	}
	return nil
}

func (scd *snowflakeChunkDownloader) next() (chunkRowType, error) {
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
			logger.Debugf("waiting for chunk idx: %v/%v",
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
		logger.Debugf("ready: chunk %v", scd.CurrentChunkIndex+1)
		scd.CurrentChunk = scd.Chunks[scd.CurrentChunkIndex]
		scd.ChunksMutex.Unlock()
		scd.CurrentChunkSize = len(scd.CurrentChunk)

		// kick off the next download
		scd.schedule()
	}

	logger.Debugf("no more data")
	if len(scd.ChunkMetas) > 0 {
		close(scd.ChunksError)
		close(scd.ChunksChan)
	}
	return chunkRowType{}, io.EOF
}

func (scd *snowflakeChunkDownloader) reset() {
	scd.Chunks = nil // detach all chunks. No way to go backward without reinitialize it.
}

func (scd *snowflakeChunkDownloader) getChunkMetas() []execResponseChunk {
	return scd.ChunkMetas
}

func (scd *snowflakeChunkDownloader) getQueryResultFormat() resultFormat {
	if scd.QueryResultFormat == "json" {
		return jsonFormat
	}
	return arrowFormat
}

func (scd *snowflakeChunkDownloader) setNextChunkDownloader(nextDownloader chunkDownloader) {
	scd.NextDownloader = nextDownloader
}

func (scd *snowflakeChunkDownloader) getNextChunkDownloader() chunkDownloader {
	return scd.NextDownloader
}

func (scd *snowflakeChunkDownloader) getRowType() []execResponseRowType {
	return scd.RowSet.RowType
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
	logger.Infof("download start chunk: %v", idx+1)
	defer scd.DoneDownloadCond.Broadcast()

	if err := scd.FuncDownloadHelper(ctx, scd, idx); err != nil {
		logger.Errorf(
			"failed to extract HTTP response body. URL: %v, err: %v", scd.ChunkMetas[idx].URL, err)
		scd.ChunksError <- &chunkError{Index: idx, Error: err}
	} else if scd.ctx.Err() == context.Canceled || scd.ctx.Err() == context.DeadlineExceeded {
		scd.ChunksError <- &chunkError{Index: idx, Error: scd.ctx.Err()}
	}
}

func downloadChunkHelper(ctx context.Context, scd *snowflakeChunkDownloader, idx int) error {
	headers := make(map[string]string)
	if len(scd.ChunkHeader) > 0 {
		logger.Debug("chunk header is provided.")
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
	logger.Infof("response returned chunk: %v, resp: %v", idx+1, resp)
	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(bufStream)
		if err != nil {
			return err
		}
		logger.Infof("HTTP: %v, URL: %v, Body: %v", resp.StatusCode, scd.ChunkMetas[idx].URL, b)
		logger.Infof("Header: %v", resp.Header)
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
	if scd.getQueryResultFormat() != arrowFormat {
		var decRespd [][]*string
		if !CustomJSONDecoderEnabled {
			dec := json.NewDecoder(st)
			for {
				if err = dec.Decode(&decRespd); err == io.EOF {
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
		highPrec := higherPrecisionEnabled(scd.ctx)
		respd, err = arc.decodeArrowChunk(scd.RowSet.RowType, highPrec)
		if err != nil {
			return err
		}
	}
	logger.Debugf(
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

type streamChunkDownloader struct {
	ctx            context.Context
	id             int64
	fetcher        streamChunkFetcher
	readErr        error
	rowStream      chan []*string
	Total          int64
	ChunkMetas     []execResponseChunk
	NextDownloader chunkDownloader
	RowSet         rowSetType
}

func (scd *streamChunkDownloader) totalUncompressedSize() (acc int64) {
	return -1
}

func (scd *streamChunkDownloader) hasNextResultSet() bool {
	return scd.readErr == nil
}

func (scd *streamChunkDownloader) nextResultSet() error {
	return scd.readErr
}

func (scd *streamChunkDownloader) start() error {
	go func() {
		var readErr = io.EOF

		logger.WithContext(scd.ctx).Infof(
			"start downloading. downloader id: %v, %v/%v rows, %v chunks",
			scd.id, len(scd.RowSet.RowType), scd.Total, len(scd.ChunkMetas))
		t := time.Now()

		defer func() {
			if readErr == io.EOF {
				logger.WithContext(scd.ctx).Infof("downloading done. downloader id: %v", scd.id)
			} else {
				logger.WithContext(scd.ctx).Debugf("downloading error. downloader id: %v", scd.id)
			}
			scd.readErr = readErr
			close(scd.rowStream)

			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					readErr = err
				} else {
					readErr = fmt.Errorf("%v", r)
				}
			}
		}()

		logger.WithContext(scd.ctx).Infof("sending initial set of rows in %vms", time.Since(t).Microseconds())
		t = time.Now()
		for _, row := range scd.RowSet.JSON {
			scd.rowStream <- row
		}
		scd.RowSet.JSON = nil

		// Download and parse one chunk at a time. The fetcher will send each
		// parsed row to the row stream. When an error occurs, the fetcher will
		// stop writing to the row stream so we can stop processing immediately
		for i, chunk := range scd.ChunkMetas {
			logger.WithContext(scd.ctx).Infof("starting chunk fetch %d (%d rows)", i, chunk.RowCount)
			if err := scd.fetcher.fetch(chunk.URL, scd.rowStream); err != nil {
				logger.WithContext(scd.ctx).Debugf(
					"failed chunk fetch %d: %#v, downloader id: %v, %v/%v rows, %v chunks",
					i, err, scd.id, len(scd.RowSet.RowType), scd.Total, len(scd.ChunkMetas))
				readErr = fmt.Errorf("chunk fetch: %w", err)
				break
			}
			logger.WithContext(scd.ctx).Infof("fetched chunk %d (%d rows) in %vms", i, chunk.RowCount, time.Since(t).Microseconds())
			t = time.Now()
		}
	}()
	return nil
}

func (scd *streamChunkDownloader) next() (chunkRowType, error) {
	if row, ok := <-scd.rowStream; ok {
		return chunkRowType{RowSet: row}, nil
	}
	return chunkRowType{}, scd.readErr
}

func (scd *streamChunkDownloader) reset() {}

func (scd *streamChunkDownloader) getChunkMetas() []execResponseChunk {
	return scd.ChunkMetas
}

func (scd *streamChunkDownloader) getQueryResultFormat() resultFormat {
	return jsonFormat
}

func (scd *streamChunkDownloader) setNextChunkDownloader(nextDownloader chunkDownloader) {
	scd.NextDownloader = nextDownloader
}

func (scd *streamChunkDownloader) getNextChunkDownloader() chunkDownloader {
	return scd.NextDownloader
}

func (scd *streamChunkDownloader) getRowType() []execResponseRowType {
	return scd.RowSet.RowType
}

func useStreamDownloader(ctx context.Context) bool {
	val := ctx.Value(streamChunkDownload)
	if val == nil {
		return false
	}
	boolVal, ok := val.(bool)
	if !ok {
		return false
	}
	return boolVal
}

type streamChunkFetcher interface {
	fetch(url string, rows chan<- []*string) error
}

type httpStreamChunkFetcher struct {
	ctx      context.Context
	client   *http.Client
	clientIP net.IP
	headers  map[string]string
	qrmk     string
}

func newStreamChunkDownloader(
	ctx context.Context,
	fetcher streamChunkFetcher,
	total int64,
	rowType []execResponseRowType,
	firstRows [][]*string,
	chunks []execResponseChunk) *streamChunkDownloader {
	return &streamChunkDownloader{
		ctx:        ctx,
		id:         rand.Int63(),
		fetcher:    fetcher,
		readErr:    nil,
		rowStream:  make(chan []*string),
		Total:      total,
		ChunkMetas: chunks,
		RowSet:     rowSetType{RowType: rowType, JSON: firstRows},
	}
}

func (f *httpStreamChunkFetcher) fetch(URL string, rows chan<- []*string) error {
	if len(f.headers) == 0 {
		f.headers = map[string]string{
			headerSseCAlgorithm: headerSseCAes,
			headerSseCKey:       f.qrmk,
		}
	}

	fullURL, _ := url.Parse(URL)
	res, err := newRetryHTTP(context.Background(), f.client, http.NewRequest, fullURL, f.headers, 0).execute()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("status (%d): %s", res.StatusCode, string(b))
	}
	if err = copyChunkStream(res.Body, rows); err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

func copyChunkStream(body io.Reader, rows chan<- []*string) error {
	bufStream := bufio.NewReader(body)
	gzipMagic, err := bufStream.Peek(2)
	if err != nil {
		return err
	}
	var source io.Reader
	if gzipMagic[0] == 0x1f && gzipMagic[1] == 0x8b {
		// detect and decompress Gzip format data
		bufStream0, err := gzip.NewReader(bufStream)
		if err != nil {
			return err
		}
		defer bufStream0.Close()
		source = bufStream0
	} else {
		source = bufStream
	}
	r := io.MultiReader(strings.NewReader("["), source, strings.NewReader("]"))
	dec := json.NewDecoder(r)
	openToken := json.Delim('[')
	closeToken := json.Delim(']')
	for {
		if t, err := dec.Token(); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("delim open: %w", err)
		} else if t != openToken {
			return fmt.Errorf("delim open: got %T", t)
		}
		for dec.More() {
			var row []*string
			if err = dec.Decode(&row); err != nil {
				return fmt.Errorf("decode: %w", err)
			}
			rows <- row
		}
		if t, err := dec.Token(); err != nil {
			return fmt.Errorf("delim close: %w", err)
		} else if t != closeToken {
			return fmt.Errorf("delim close: got %T", t)
		}
	}
	return nil
}
