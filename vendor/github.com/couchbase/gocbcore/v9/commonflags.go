package gocbcore

const (
	// Legacy flag format for JSON data.
	lfJSON = 0

	// Common flags mask
	cfMask = 0xFF000000
	// Common flags mask for data format
	cfFmtMask = 0x0F000000
	// Common flags mask for compression mode.
	cfCmprMask = 0xE0000000

	// Common flag format for sdk-private data.
	cfFmtPrivate = 1 << 24 // nolint: deadcode,varcheck,unused
	// Common flag format for JSON data.
	cfFmtJSON = 2 << 24
	// Common flag format for binary data.
	cfFmtBinary = 3 << 24
	// Common flag format for string data.
	cfFmtString = 4 << 24

	// Common flags compression for disabled compression.
	cfCmprNone = 0 << 29
)

// DataType represents the type of data for a value
type DataType uint32

// CompressionType indicates the type of compression for a value
type CompressionType uint32

const (
	// UnknownType indicates the values type is unknown.
	UnknownType = DataType(0)

	// JSONType indicates the value is JSON data.
	JSONType = DataType(1)

	// BinaryType indicates the value is binary data.
	BinaryType = DataType(2)

	// StringType indicates the value is string data.
	StringType = DataType(3)
)

const (
	// UnknownCompression indicates that the compression type is unknown.
	UnknownCompression = CompressionType(0)

	// NoCompression indicates that no compression is being used.
	NoCompression = CompressionType(1)
)

// EncodeCommonFlags encodes a data type and compression type into a flags
// value using the common flags specification.
func EncodeCommonFlags(valueType DataType, compression CompressionType) uint32 {
	var flags uint32

	switch valueType {
	case JSONType:
		flags |= cfFmtJSON
	case BinaryType:
		flags |= cfFmtBinary
	case StringType:
		flags |= cfFmtString
	case UnknownType:
		// flags |= ?
	}

	switch compression {
	case NoCompression:
		// flags |= 0
	case UnknownCompression:
		// flags |= ?
	}

	return flags
}

// DecodeCommonFlags decodes a flags value into a data type and compression type
// using the common flags specification.
func DecodeCommonFlags(flags uint32) (DataType, CompressionType) {
	// Check for legacy flags
	if flags&cfMask == 0 {
		// Legacy Flags
		if flags == lfJSON {
			// Legacy JSON
			flags = cfFmtJSON
		} else {
			return UnknownType, UnknownCompression
		}
	}

	valueType := UnknownType
	compression := UnknownCompression

	if flags&cfFmtMask == cfFmtBinary {
		valueType = BinaryType
	} else if flags&cfFmtMask == cfFmtString {
		valueType = StringType
	} else if flags&cfFmtMask == cfFmtJSON {
		valueType = JSONType
	}

	if flags&cfCmprMask == cfCmprNone {
		compression = NoCompression
	}

	return valueType, compression
}
