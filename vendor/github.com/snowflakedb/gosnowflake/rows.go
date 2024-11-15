// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"database/sql/driver"
	"io"
	"reflect"
	"strings"
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

// SnowflakeRows provides an API for methods exposed to the clients
type SnowflakeRows interface {
	GetQueryID() string
	GetStatus() queryStatus
	GetArrowBatches() ([]*ArrowBatch, error)
}

type snowflakeRows struct {
	sc                  *snowflakeConn
	ChunkDownloader     chunkDownloader
	tailChunkDownloader chunkDownloader
	queryID             string
	status              queryStatus
	err                 error
	errChannel          chan error
	location            *time.Location
	ctx                 context.Context
}

func (rows *snowflakeRows) getLocation() *time.Location {
	if rows.location == nil && rows.sc != nil && rows.sc.cfg != nil {
		rows.location = getCurrentLocation(rows.sc.cfg.Params)
	}
	return rows.location
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

func (rows *snowflakeRows) Close() (err error) {
	if err := rows.waitForAsyncQueryStatus(); err != nil {
		return err
	}
	logger.WithContext(rows.sc.ctx).Debugln("Rows.Close")
	return nil
}

// ColumnTypeDatabaseTypeName returns the database column name.
func (rows *snowflakeRows) ColumnTypeDatabaseTypeName(index int) string {
	if err := rows.waitForAsyncQueryStatus(); err != nil {
		return err.Error()
	}
	return strings.ToUpper(rows.ChunkDownloader.getRowType()[index].Type)
}

// ColumnTypeLength returns the length of the column
func (rows *snowflakeRows) ColumnTypeLength(index int) (length int64, ok bool) {
	if err := rows.waitForAsyncQueryStatus(); err != nil {
		return 0, false
	}
	if index < 0 || index > len(rows.ChunkDownloader.getRowType()) {
		return 0, false
	}
	switch rows.ChunkDownloader.getRowType()[index].Type {
	case "text", "variant", "object", "array", "binary":
		return rows.ChunkDownloader.getRowType()[index].Length, true
	}
	return 0, false
}

func (rows *snowflakeRows) ColumnTypeNullable(index int) (nullable, ok bool) {
	if err := rows.waitForAsyncQueryStatus(); err != nil {
		return false, false
	}
	if index < 0 || index > len(rows.ChunkDownloader.getRowType()) {
		return false, false
	}
	return rows.ChunkDownloader.getRowType()[index].Nullable, true
}

func (rows *snowflakeRows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	if err := rows.waitForAsyncQueryStatus(); err != nil {
		return 0, 0, false
	}
	rowType := rows.ChunkDownloader.getRowType()
	if index < 0 || index > len(rowType) {
		return 0, 0, false
	}
	switch rowType[index].Type {
	case "fixed":
		return rowType[index].Precision, rowType[index].Scale, true
	case "time":
		return rowType[index].Scale, 0, true
	case "timestamp":
		return rowType[index].Scale, 0, true
	}
	return 0, 0, false
}

func (rows *snowflakeRows) Columns() []string {
	if err := rows.waitForAsyncQueryStatus(); err != nil {
		return make([]string, 0)
	}
	logger.WithContext(rows.ctx).Debug("Rows.Columns")
	ret := make([]string, len(rows.ChunkDownloader.getRowType()))
	for i, n := 0, len(rows.ChunkDownloader.getRowType()); i < n; i++ {
		ret[i] = rows.ChunkDownloader.getRowType()[i].Name
	}
	return ret
}

func (rows *snowflakeRows) ColumnTypeScanType(index int) reflect.Type {
	if err := rows.waitForAsyncQueryStatus(); err != nil {
		return nil
	}
	return snowflakeTypeToGo(rows.ctx, getSnowflakeType(rows.ChunkDownloader.getRowType()[index].Type), rows.ChunkDownloader.getRowType()[index].Scale, rows.ChunkDownloader.getRowType()[index].Fields)
}

func (rows *snowflakeRows) GetQueryID() string {
	return rows.queryID
}

func (rows *snowflakeRows) GetStatus() queryStatus {
	return rows.status
}

// GetArrowBatches returns an array of ArrowBatch objects to retrieve data in arrow.Record format
func (rows *snowflakeRows) GetArrowBatches() ([]*ArrowBatch, error) {
	// Wait for all arrow batches before fetching.
	// Otherwise, a panic error "invalid memory address or nil pointer dereference" will be thrown.
	if err := rows.waitForAsyncQueryStatus(); err != nil {
		return nil, err
	}

	return rows.ChunkDownloader.getArrowBatches(), nil
}

func (rows *snowflakeRows) Next(dest []driver.Value) (err error) {
	if err = rows.waitForAsyncQueryStatus(); err != nil {
		return err
	}
	row, err := rows.ChunkDownloader.next()
	if err != nil {
		// includes io.EOF
		if err == io.EOF {
			rows.ChunkDownloader.reset()
		}
		return err
	}

	if rows.ChunkDownloader.getQueryResultFormat() == arrowFormat {
		for i, n := 0, len(row.ArrowRow); i < n; i++ {
			dest[i] = row.ArrowRow[i]
		}
	} else {
		for i, n := 0, len(row.RowSet); i < n; i++ {
			// could move to chunk downloader so that each go routine
			// can convert data
			err = stringToValue(rows.ctx, &dest[i], rows.ChunkDownloader.getRowType()[i], row.RowSet[i], rows.getLocation(), rows.sc.cfg.Params)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (rows *snowflakeRows) HasNextResultSet() bool {
	if err := rows.waitForAsyncQueryStatus(); err != nil {
		return false
	}
	return rows.ChunkDownloader.hasNextResultSet()
}

func (rows *snowflakeRows) NextResultSet() error {
	if err := rows.waitForAsyncQueryStatus(); err != nil {
		return err
	}
	if len(rows.ChunkDownloader.getChunkMetas()) == 0 {
		if rows.ChunkDownloader.getNextChunkDownloader() == nil {
			return io.EOF
		}
		rows.ChunkDownloader = rows.ChunkDownloader.getNextChunkDownloader()
		if err := rows.ChunkDownloader.start(); err != nil {
			return err
		}
	}
	return rows.ChunkDownloader.nextResultSet()
}

func (rows *snowflakeRows) waitForAsyncQueryStatus() error {
	// if async query, block until query is finished
	if rows.status == QueryStatusInProgress {
		err := <-rows.errChannel
		rows.status = QueryStatusComplete
		if err != nil {
			rows.status = QueryFailed
			rows.err = err
			return rows.err
		}
	} else if rows.status == QueryFailed {
		return rows.err
	}
	return nil
}

func (rows *snowflakeRows) addDownloader(newDL chunkDownloader) {
	if rows.ChunkDownloader == nil {
		rows.ChunkDownloader = newDL
		rows.tailChunkDownloader = newDL
		return
	}
	rows.tailChunkDownloader.setNextChunkDownloader(newDL)
	rows.tailChunkDownloader = newDL
}
