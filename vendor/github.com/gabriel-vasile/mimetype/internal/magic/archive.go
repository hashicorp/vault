package magic

import (
	"bytes"
	"encoding/binary"
)

var (
	// SevenZ matches a 7z archive.
	SevenZ = prefix([]byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C})
	// Gzip matches gzip files based on http://www.zlib.org/rfc-gzip.html#header-trailer.
	Gzip = prefix([]byte{0x1f, 0x8b})
	// Tar matches a (t)ape (ar)chive file.
	Tar = offset([]byte("ustar"), 257)
	// Fits matches an Flexible Image Transport System file.
	Fits = prefix([]byte{
		0x53, 0x49, 0x4D, 0x50, 0x4C, 0x45, 0x20, 0x20, 0x3D, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x54,
	})
	// Xar matches an eXtensible ARchive format file.
	Xar = prefix([]byte{0x78, 0x61, 0x72, 0x21})
	// Bz2 matches a bzip2 file.
	Bz2 = prefix([]byte{0x42, 0x5A, 0x68})
	// Ar matches an ar (Unix) archive file.
	Ar = prefix([]byte{0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E})
	// Deb matches a Debian package file.
	Deb = offset([]byte{
		0x64, 0x65, 0x62, 0x69, 0x61, 0x6E, 0x2D,
		0x62, 0x69, 0x6E, 0x61, 0x72, 0x79,
	}, 8)
	// Warc matches a Web ARChive file.
	Warc = prefix([]byte("WARC/"))
	// Cab matches a Cabinet archive file.
	Cab = prefix([]byte("MSCF"))
	// Xz matches an xz compressed stream based on https://tukaani.org/xz/xz-file-format.txt.
	Xz = prefix([]byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00})
	// Lzip matches an Lzip compressed file.
	Lzip = prefix([]byte{0x4c, 0x5a, 0x49, 0x50})
)

// Zstd matches a Zstandard archive file.
func Zstd(raw []byte, limit uint32) bool {
	return len(raw) >= 4 &&
		(0x22 <= raw[0] && raw[0] <= 0x28 || raw[0] == 0x1E) && // Different Zstandard versions.
		bytes.HasPrefix(raw[1:], []byte{0xB5, 0x2F, 0xFD})
}

// Rpm matches an RPM or Delta RPM package file.
func Rpm(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte{0xed, 0xab, 0xee, 0xdb}) ||
		bytes.HasPrefix(raw, []byte("drpm"))
}

// Cpio matches a cpio archive file.
func Cpio(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte("070707")) ||
		bytes.HasPrefix(raw, []byte("070701")) ||
		bytes.HasPrefix(raw, []byte("070702"))
}

// Rar matches a RAR archive file.
func Rar(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte("Rar!\x1A\x07\x00")) ||
		bytes.HasPrefix(raw, []byte("Rar!\x1A\x07\x01\x00"))
}

// Crx matches a Chrome extension file: a zip archive prepended by a package header.
func Crx(raw []byte, limit uint32) bool {
	const minHeaderLen = 16
	if len(raw) < minHeaderLen || !bytes.HasPrefix(raw, []byte("Cr24")) {
		return false
	}
	pubkeyLen := binary.LittleEndian.Uint32(raw[8:12])
	sigLen := binary.LittleEndian.Uint32(raw[12:16])
	zipOffset := minHeaderLen + pubkeyLen + sigLen
	if uint32(len(raw)) < zipOffset {
		return false
	}
	return Zip(raw[zipOffset:], limit)
}
