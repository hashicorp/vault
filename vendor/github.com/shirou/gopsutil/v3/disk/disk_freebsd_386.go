// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs types_freebsd.go

package disk

const (
	sizeofPtr        = 0x4
	sizeofShort      = 0x2
	sizeofInt        = 0x4
	sizeofLong       = 0x4
	sizeofLongLong   = 0x8
	sizeofLongDouble = 0x8

	devstat_NO_DATA = 0x00
	devstat_READ    = 0x01
	devstat_WRITE   = 0x02
	devstat_FREE    = 0x03
)

const (
	sizeOfdevstat = 0xf0
)

type (
	_C_short       int16
	_C_int         int32
	_C_long        int32
	_C_long_long   int64
	_C_long_double int64
)

type devstat struct {
	Sequence0     uint32
	Allocated     int32
	Start_count   uint32
	End_count     uint32
	Busy_from     bintime
	Dev_links     _Ctype_struct___0
	Device_number uint32
	Device_name   [16]int8
	Unit_number   int32
	Bytes         [4]uint64
	Operations    [4]uint64
	Duration      [4]bintime
	Busy_time     bintime
	Creation_time bintime
	Block_size    uint32
	Tag_types     [3]uint64
	Flags         uint32
	Device_type   uint32
	Priority      uint32
	Id            *byte
	Sequence1     uint32
}

type bintime struct {
	Sec  int32
	Frac uint64
}

type _Ctype_struct___0 struct {
	Empty uint32
}
