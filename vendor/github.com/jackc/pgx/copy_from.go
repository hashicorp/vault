package pgx

import (
	"bytes"
	"fmt"
	"io"

	"github.com/jackc/pgx/pgio"
	"github.com/jackc/pgx/pgproto3"
	"github.com/pkg/errors"
)

// CopyFromRows returns a CopyFromSource interface over the provided rows slice
// making it usable by *Conn.CopyFrom.
func CopyFromRows(rows [][]interface{}) CopyFromSource {
	return &copyFromRows{rows: rows, idx: -1}
}

type copyFromRows struct {
	rows [][]interface{}
	idx  int
}

func (ctr *copyFromRows) Next() bool {
	ctr.idx++
	return ctr.idx < len(ctr.rows)
}

func (ctr *copyFromRows) Values() ([]interface{}, error) {
	return ctr.rows[ctr.idx], nil
}

func (ctr *copyFromRows) Err() error {
	return nil
}

// CopyFromSource is the interface used by *Conn.CopyFrom as the source for copy data.
type CopyFromSource interface {
	// Next returns true if there is another row and makes the next row data
	// available to Values(). When there are no more rows available or an error
	// has occurred it returns false.
	Next() bool

	// Values returns the values for the current row.
	Values() ([]interface{}, error)

	// Err returns any error that has been encountered by the CopyFromSource. If
	// this is not nil *Conn.CopyFrom will abort the copy.
	Err() error
}

type copyFrom struct {
	conn          *Conn
	tableName     Identifier
	columnNames   []string
	rowSrc        CopyFromSource
	readerErrChan chan error
}

func (ct *copyFrom) readUntilReadyForQuery() {
	for {
		msg, err := ct.conn.rxMsg()
		if err != nil {
			ct.readerErrChan <- err
			close(ct.readerErrChan)
			return
		}

		switch msg := msg.(type) {
		case *pgproto3.ReadyForQuery:
			ct.conn.rxReadyForQuery(msg)
			close(ct.readerErrChan)
			return
		case *pgproto3.CommandComplete:
		case *pgproto3.ErrorResponse:
			ct.readerErrChan <- ct.conn.rxErrorResponse(msg)
		default:
			err = ct.conn.processContextFreeMsg(msg)
			if err != nil {
				ct.readerErrChan <- ct.conn.processContextFreeMsg(msg)
			}
		}
	}
}

func (ct *copyFrom) waitForReaderDone() error {
	var err error
	for err = range ct.readerErrChan {
	}
	return err
}

func (ct *copyFrom) run() (int, error) {
	quotedTableName := ct.tableName.Sanitize()
	cbuf := &bytes.Buffer{}
	for i, cn := range ct.columnNames {
		if i != 0 {
			cbuf.WriteString(", ")
		}
		cbuf.WriteString(quoteIdentifier(cn))
	}
	quotedColumnNames := cbuf.String()

	ps, err := ct.conn.Prepare("", fmt.Sprintf("select %s from %s", quotedColumnNames, quotedTableName))
	if err != nil {
		return 0, err
	}

	err = ct.conn.sendSimpleQuery(fmt.Sprintf("copy %s ( %s ) from stdin binary;", quotedTableName, quotedColumnNames))
	if err != nil {
		return 0, err
	}

	err = ct.conn.readUntilCopyInResponse()
	if err != nil {
		return 0, err
	}

	panicked := true

	go ct.readUntilReadyForQuery()
	defer ct.waitForReaderDone()
	defer func() {
		if panicked {
			ct.conn.die(errors.New("panic while in copy from"))
		}
	}()

	buf := ct.conn.wbuf
	buf = append(buf, copyData)
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)

	buf = append(buf, "PGCOPY\n\377\r\n\000"...)
	buf = pgio.AppendInt32(buf, 0)
	buf = pgio.AppendInt32(buf, 0)

	var sentCount int

	moreRows := true
	for moreRows {
		select {
		case err = <-ct.readerErrChan:
			panicked = false
			return 0, err
		default:
		}

		var addedRows int
		var err error
		moreRows, buf, addedRows, err = ct.buildCopyBuf(buf, ps)
		if err != nil {
			panicked = false
			ct.cancelCopyIn()
			return 0, err
		}
		sentCount += addedRows
		pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

		_, err = ct.conn.conn.Write(buf)
		if err != nil {
			panicked = false
			ct.conn.die(err)
			return 0, err
		}

		// Directly manipulate wbuf to reset to reuse the same buffer
		buf = buf[0:5]

	}

	if ct.rowSrc.Err() != nil {
		panicked = false
		ct.cancelCopyIn()
		return 0, ct.rowSrc.Err()
	}

	buf = pgio.AppendInt16(buf, -1) // terminate the copy stream
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	buf = append(buf, copyDone)
	buf = pgio.AppendInt32(buf, 4)

	_, err = ct.conn.conn.Write(buf)
	if err != nil {
		panicked = false
		ct.conn.die(err)
		return 0, err
	}

	err = ct.waitForReaderDone()
	if err != nil {
		panicked = false
		return 0, err
	}

	panicked = false
	return sentCount, nil
}

func (ct *copyFrom) buildCopyBuf(buf []byte, ps *PreparedStatement) (bool, []byte, int, error) {
	var rowCount int

	for ct.rowSrc.Next() {
		values, err := ct.rowSrc.Values()
		if err != nil {
			return false, nil, 0, err
		}
		if len(values) != len(ct.columnNames) {
			return false, nil, 0, errors.Errorf("expected %d values, got %d values", len(ct.columnNames), len(values))
		}

		buf = pgio.AppendInt16(buf, int16(len(ct.columnNames)))
		for i, val := range values {
			buf, err = encodePreparedStatementArgument(ct.conn.ConnInfo, buf, ps.FieldDescriptions[i].DataType, val)
			if err != nil {
				return false, nil, 0, err
			}
		}

		rowCount++

		if len(buf) > 65536 {
			return true, buf, rowCount, nil
		}
	}

	return false, buf, rowCount, nil
}

func (c *Conn) readUntilCopyInResponse() error {
	for {
		msg, err := c.rxMsg()
		if err != nil {
			return err
		}

		switch msg := msg.(type) {
		case *pgproto3.CopyInResponse:
			return nil
		default:
			err = c.processContextFreeMsg(msg)
			if err != nil {
				return err
			}
		}
	}
}

func (ct *copyFrom) cancelCopyIn() error {
	buf := ct.conn.wbuf
	buf = append(buf, copyFail)
	sp := len(buf)
	buf = pgio.AppendInt32(buf, -1)
	buf = append(buf, "client error: abort"...)
	buf = append(buf, 0)
	pgio.SetInt32(buf[sp:], int32(len(buf[sp:])))

	_, err := ct.conn.conn.Write(buf)
	if err != nil {
		ct.conn.die(err)
		return err
	}

	return nil
}

// CopyFrom uses the PostgreSQL copy protocol to perform bulk data insertion.
// It returns the number of rows copied and an error.
//
// CopyFrom requires all values use the binary format. Almost all types
// implemented by pgx use the binary format by default. Types implementing
// Encoder can only be used if they encode to the binary format.
func (c *Conn) CopyFrom(tableName Identifier, columnNames []string, rowSrc CopyFromSource) (int, error) {
	ct := &copyFrom{
		conn:          c,
		tableName:     tableName,
		columnNames:   columnNames,
		rowSrc:        rowSrc,
		readerErrChan: make(chan error),
	}

	return ct.run()
}

// CopyFromReader uses the PostgreSQL textual format of the copy protocol
func (c *Conn) CopyFromReader(r io.Reader, sql string) (CommandTag, error) {
	if err := c.sendSimpleQuery(sql); err != nil {
		return "", err
	}

	if err := c.readUntilCopyInResponse(); err != nil {
		return "", err
	}
	buf := c.wbuf

	buf = append(buf, copyData)
	sp := len(buf)
	for {
		n, err := r.Read(buf[5:cap(buf)])
		if err == io.EOF && n == 0 {
			break
		}
		buf = buf[0 : n+5]
		pgio.SetInt32(buf[sp:], int32(n+4))

		if _, err := c.conn.Write(buf); err != nil {
			return "", err
		}
	}

	buf = buf[:0]
	buf = append(buf, copyDone)
	buf = pgio.AppendInt32(buf, 4)

	if _, err := c.conn.Write(buf); err != nil {
		return "", err
	}

	for {
		msg, err := c.rxMsg()
		if err != nil {
			return "", err
		}

		switch msg := msg.(type) {
		case *pgproto3.ReadyForQuery:
			c.rxReadyForQuery(msg)
			return "", err
		case *pgproto3.CommandComplete:
			return CommandTag(msg.CommandTag), nil
		case *pgproto3.ErrorResponse:
			return "", c.rxErrorResponse(msg)
		default:
			return "", c.processContextFreeMsg(msg)
		}
	}
}
