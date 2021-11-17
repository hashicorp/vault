package pgx

import (
	"io"

	"github.com/jackc/pgx/pgtype"
)

// LargeObjects is a structure used to access the large objects API. It is only
// valid within the transaction where it was created.
//
// For more details see: http://www.postgresql.org/docs/current/static/largeobjects.html
type LargeObjects struct {
	// Has64 is true if the server is capable of working with 64-bit numbers
	Has64 bool
	fp    *fastpath
}

const largeObjectFns = `select proname, oid from pg_catalog.pg_proc
where proname in (
'lo_open',
'lo_close',
'lo_create',
'lo_unlink',
'lo_lseek',
'lo_lseek64',
'lo_tell',
'lo_tell64',
'lo_truncate',
'lo_truncate64',
'loread',
'lowrite')
and pronamespace = (select oid from pg_catalog.pg_namespace where nspname = 'pg_catalog')`

// LargeObjects returns a LargeObjects instance for the transaction.
func (tx *Tx) LargeObjects() (*LargeObjects, error) {
	if tx.conn.fp == nil {
		tx.conn.fp = newFastpath(tx.conn)
	}
	if _, exists := tx.conn.fp.fns["lo_open"]; !exists {
		res, err := tx.Query(largeObjectFns)
		if err != nil {
			return nil, err
		}
		if err := tx.conn.fp.addFunctions(res); err != nil {
			return nil, err
		}
	}

	lo := &LargeObjects{fp: tx.conn.fp}
	_, lo.Has64 = lo.fp.fns["lo_lseek64"]

	return lo, nil
}

type LargeObjectMode int32

const (
	LargeObjectModeWrite LargeObjectMode = 0x20000
	LargeObjectModeRead  LargeObjectMode = 0x40000
)

// Create creates a new large object. If id is zero, the server assigns an
// unused OID.
func (o *LargeObjects) Create(id pgtype.OID) (pgtype.OID, error) {
	newOID, err := fpInt32(o.fp.CallFn("lo_create", []fpArg{fpIntArg(int32(id))}))
	return pgtype.OID(newOID), err
}

// Open opens an existing large object with the given mode.
func (o *LargeObjects) Open(oid pgtype.OID, mode LargeObjectMode) (*LargeObject, error) {
	fd, err := fpInt32(o.fp.CallFn("lo_open", []fpArg{fpIntArg(int32(oid)), fpIntArg(int32(mode))}))
	return &LargeObject{fd: fd, lo: o}, err
}

// Unlink removes a large object from the database.
func (o *LargeObjects) Unlink(oid pgtype.OID) error {
	_, err := o.fp.CallFn("lo_unlink", []fpArg{fpIntArg(int32(oid))})
	return err
}

// A LargeObject is a large object stored on the server. It is only valid within
// the transaction that it was initialized in. It implements these interfaces:
//
//    io.Writer
//    io.Reader
//    io.Seeker
//    io.Closer
type LargeObject struct {
	fd int32
	lo *LargeObjects
}

// Write writes p to the large object and returns the number of bytes written
// and an error if not all of p was written.
func (o *LargeObject) Write(p []byte) (int, error) {
	n, err := fpInt32(o.lo.fp.CallFn("lowrite", []fpArg{fpIntArg(o.fd), p}))
	return int(n), err
}

// Read reads up to len(p) bytes into p returning the number of bytes read.
func (o *LargeObject) Read(p []byte) (int, error) {
	res, err := o.lo.fp.CallFn("loread", []fpArg{fpIntArg(o.fd), fpIntArg(int32(len(p)))})
	if len(res) < len(p) {
		err = io.EOF
	}
	return copy(p, res), err
}

// Seek moves the current location pointer to the new location specified by offset.
func (o *LargeObject) Seek(offset int64, whence int) (n int64, err error) {
	if o.lo.Has64 {
		n, err = fpInt64(o.lo.fp.CallFn("lo_lseek64", []fpArg{fpIntArg(o.fd), fpInt64Arg(offset), fpIntArg(int32(whence))}))
	} else {
		var n32 int32
		n32, err = fpInt32(o.lo.fp.CallFn("lo_lseek", []fpArg{fpIntArg(o.fd), fpIntArg(int32(offset)), fpIntArg(int32(whence))}))
		n = int64(n32)
	}
	return
}

// Tell returns the current read or write location of the large object
// descriptor.
func (o *LargeObject) Tell() (n int64, err error) {
	if o.lo.Has64 {
		n, err = fpInt64(o.lo.fp.CallFn("lo_tell64", []fpArg{fpIntArg(o.fd)}))
	} else {
		var n32 int32
		n32, err = fpInt32(o.lo.fp.CallFn("lo_tell", []fpArg{fpIntArg(o.fd)}))
		n = int64(n32)
	}
	return
}

// Trunctes the large object to size.
func (o *LargeObject) Truncate(size int64) (err error) {
	if o.lo.Has64 {
		_, err = o.lo.fp.CallFn("lo_truncate64", []fpArg{fpIntArg(o.fd), fpInt64Arg(size)})
	} else {
		_, err = o.lo.fp.CallFn("lo_truncate", []fpArg{fpIntArg(o.fd), fpIntArg(int32(size))})
	}
	return
}

// Close closees the large object descriptor.
func (o *LargeObject) Close() error {
	_, err := o.lo.fp.CallFn("lo_close", []fpArg{fpIntArg(o.fd)})
	return err
}
