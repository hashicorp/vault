package pgx

import (
	"context"
	"errors"
	"io"
)

// LargeObjects is a structure used to access the large objects API. It is only valid within the transaction where it
// was created.
//
// For more details see: http://www.postgresql.org/docs/current/static/largeobjects.html
type LargeObjects struct {
	tx Tx
}

type LargeObjectMode int32

const (
	LargeObjectModeWrite LargeObjectMode = 0x20000
	LargeObjectModeRead  LargeObjectMode = 0x40000
)

// Create creates a new large object. If oid is zero, the server assigns an unused OID.
func (o *LargeObjects) Create(ctx context.Context, oid uint32) (uint32, error) {
	err := o.tx.QueryRow(ctx, "select lo_create($1)", oid).Scan(&oid)
	return oid, err
}

// Open opens an existing large object with the given mode. ctx will also be used for all operations on the opened large
// object.
func (o *LargeObjects) Open(ctx context.Context, oid uint32, mode LargeObjectMode) (*LargeObject, error) {
	var fd int32
	err := o.tx.QueryRow(ctx, "select lo_open($1, $2)", oid, mode).Scan(&fd)
	if err != nil {
		return nil, err
	}
	return &LargeObject{fd: fd, tx: o.tx, ctx: ctx}, nil
}

// Unlink removes a large object from the database.
func (o *LargeObjects) Unlink(ctx context.Context, oid uint32) error {
	var result int32
	err := o.tx.QueryRow(ctx, "select lo_unlink($1)", oid).Scan(&result)
	if err != nil {
		return err
	}

	if result != 1 {
		return errors.New("failed to remove large object")
	}

	return nil
}

// A LargeObject is a large object stored on the server. It is only valid within the transaction that it was initialized
// in. It uses the context it was initialized with for all operations. It implements these interfaces:
//
//	io.Writer
//	io.Reader
//	io.Seeker
//	io.Closer
type LargeObject struct {
	ctx context.Context
	tx  Tx
	fd  int32
}

// Write writes p to the large object and returns the number of bytes written and an error if not all of p was written.
func (o *LargeObject) Write(p []byte) (int, error) {
	var n int
	err := o.tx.QueryRow(o.ctx, "select lowrite($1, $2)", o.fd, p).Scan(&n)
	if err != nil {
		return n, err
	}

	if n < 0 {
		return 0, errors.New("failed to write to large object")
	}

	return n, nil
}

// Read reads up to len(p) bytes into p returning the number of bytes read.
func (o *LargeObject) Read(p []byte) (int, error) {
	var res []byte
	err := o.tx.QueryRow(o.ctx, "select loread($1, $2)", o.fd, len(p)).Scan(&res)
	copy(p, res)
	if err != nil {
		return len(res), err
	}

	if len(res) < len(p) {
		err = io.EOF
	}
	return len(res), err
}

// Seek moves the current location pointer to the new location specified by offset.
func (o *LargeObject) Seek(offset int64, whence int) (n int64, err error) {
	err = o.tx.QueryRow(o.ctx, "select lo_lseek64($1, $2, $3)", o.fd, offset, whence).Scan(&n)
	return n, err
}

// Tell returns the current read or write location of the large object descriptor.
func (o *LargeObject) Tell() (n int64, err error) {
	err = o.tx.QueryRow(o.ctx, "select lo_tell64($1)", o.fd).Scan(&n)
	return n, err
}

// Truncate the large object to size.
func (o *LargeObject) Truncate(size int64) (err error) {
	_, err = o.tx.Exec(o.ctx, "select lo_truncate64($1, $2)", o.fd, size)
	return err
}

// Close the large object descriptor.
func (o *LargeObject) Close() error {
	_, err := o.tx.Exec(o.ctx, "select lo_close($1)", o.fd)
	return err
}
