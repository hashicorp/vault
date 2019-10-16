// +build ignore

/*
Input to cgo -godefs.
*/

package host

/*
#define KERNEL
#include <sys/types.h>
#include <sys/time.h>
#include <utmpx.h>
#include "freebsd_headers/utxdb.h"

enum {
	sizeofPtr = sizeof(void*),
};

*/
import "C"

// Machine characteristics; for internal use.

const (
	sizeofPtr      = C.sizeofPtr
	sizeofShort    = C.sizeof_short
	sizeofInt      = C.sizeof_int
	sizeofLong     = C.sizeof_long
	sizeofLongLong = C.sizeof_longlong
	sizeOfUtmpx    = C.sizeof_struct_futx
)

// Basic types

type (
	_C_short     C.short
	_C_int       C.int
	_C_long      C.long
	_C_long_long C.longlong
)

type Utmp C.struct_utmp // for FreeBSD 9.0 compatibility
type Utmpx C.struct_futx
