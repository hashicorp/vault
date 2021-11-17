package magic

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
)

var (
	// Lnk matches Microsoft lnk binary format.
	Lnk = prefix([]byte{0x4C, 0x00, 0x00, 0x00, 0x01, 0x14, 0x02, 0x00})
	// Wasm matches a web assembly File Format file.
	Wasm = prefix([]byte{0x00, 0x61, 0x73, 0x6D})
	// Exe matches a Windows/DOS executable file.
	Exe = prefix([]byte{0x4D, 0x5A})
	// Elf matches an Executable and Linkable Format file.
	Elf = prefix([]byte{0x7F, 0x45, 0x4C, 0x46})
	// Nes matches a Nintendo Entertainment system ROM file.
	Nes = prefix([]byte{0x4E, 0x45, 0x53, 0x1A})
	// TzIf matches a Time Zone Information Format (TZif) file.
	TzIf = prefix([]byte("TZif"))
)

// Java bytecode and Mach-O binaries share the same magic number.
// More info here https://github.com/threatstack/libmagic/blob/master/magic/Magdir/cafebabe
func classOrMachOFat(in []byte) bool {
	// There should be at least 8 bytes for both of them because the only way to
	// quickly distinguish them is by comparing byte at position 7
	if len(in) < 8 {
		return false
	}

	return bytes.HasPrefix(in, []byte{0xCA, 0xFE, 0xBA, 0xBE})
}

// Class matches a java class file.
func Class(raw []byte, limit uint32) bool {
	return classOrMachOFat(raw) && raw[7] > 30
}

// MachO matches Mach-O binaries format.
func MachO(raw []byte, limit uint32) bool {
	if classOrMachOFat(raw) && raw[7] < 20 {
		return true
	}

	if len(raw) < 4 {
		return false
	}

	be := binary.BigEndian.Uint32(raw)
	le := binary.LittleEndian.Uint32(raw)

	return be == macho.Magic32 ||
		le == macho.Magic32 ||
		be == macho.Magic64 ||
		le == macho.Magic64
}

// Swf matches an Adobe Flash swf file.
func Swf(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte("CWS")) ||
		bytes.HasPrefix(raw, []byte("FWS")) ||
		bytes.HasPrefix(raw, []byte("ZWS"))
}

// Dbf matches a dBase file.
// https://www.dbase.com/Knowledgebase/INT/db7_file_fmt.htm
func Dbf(raw []byte, limit uint32) bool {
	if len(raw) < 4 {
		return false
	}

	// 3rd and 4th bytes contain the last update month and day of month
	if !(0 < raw[2] && raw[2] < 13 && 0 < raw[3] && raw[3] < 32) {
		return false
	}

	// dbf type is dictated by the first byte
	dbfTypes := []byte{
		0x02, 0x03, 0x04, 0x05, 0x30, 0x31, 0x32, 0x42, 0x62, 0x7B, 0x82,
		0x83, 0x87, 0x8A, 0x8B, 0x8E, 0xB3, 0xCB, 0xE5, 0xF5, 0xF4, 0xFB,
	}
	for _, b := range dbfTypes {
		if raw[0] == b {
			return true
		}
	}

	return false
}

// ElfObj matches an object file.
func ElfObj(raw []byte, limit uint32) bool {
	return len(raw) > 17 && ((raw[16] == 0x01 && raw[17] == 0x00) ||
		(raw[16] == 0x00 && raw[17] == 0x01))
}

// ElfExe matches an executable file.
func ElfExe(raw []byte, limit uint32) bool {
	return len(raw) > 17 && ((raw[16] == 0x02 && raw[17] == 0x00) ||
		(raw[16] == 0x00 && raw[17] == 0x02))
}

// ElfLib matches a shared library file.
func ElfLib(raw []byte, limit uint32) bool {
	return len(raw) > 17 && ((raw[16] == 0x03 && raw[17] == 0x00) ||
		(raw[16] == 0x00 && raw[17] == 0x03))
}

// ElfDump matches a core dump file.
func ElfDump(raw []byte, limit uint32) bool {
	return len(raw) > 17 && ((raw[16] == 0x04 && raw[17] == 0x00) ||
		(raw[16] == 0x00 && raw[17] == 0x04))
}

// Dcm matches a DICOM medical format file.
func Dcm(raw []byte, limit uint32) bool {
	return len(raw) > 131 &&
		bytes.Equal(raw[128:132], []byte{0x44, 0x49, 0x43, 0x4D})
}

// Marc matches a MARC21 (MAchine-Readable Cataloging) file.
func Marc(raw []byte, limit uint32) bool {
	// File is at least 24 bytes ("leader" field size).
	if len(raw) < 24 {
		return false
	}

	// Fixed bytes at offset 20.
	if !bytes.Equal(raw[20:24], []byte("4500")) {
		return false
	}

	// First 5 bytes are ASCII digits.
	for i := 0; i < 5; i++ {
		if raw[i] < '0' || raw[i] > '9' {
			return false
		}
	}

	// Field terminator is present.
	return bytes.Contains(raw, []byte{0x1E})
}
