package gocb

import (
	"encoding/json"
	"errors"

	gocbcore "github.com/couchbase/gocbcore/v10"
)

// Transcoder provides an interface for transforming Go values to and
// from raw bytes for storage and retreival from Couchbase data storage.
type Transcoder interface {
	// Decodes retrieved bytes into a Go type.
	Decode([]byte, uint32, interface{}) error

	// Encodes a Go type into bytes for storage.
	Encode(interface{}) ([]byte, uint32, error)
}

// JSONTranscoder implements the default transcoding behavior and applies JSON transcoding to all values.
//
// This will apply the following behavior to the value:
// binary ([]byte) -> error.
// default -> JSON value, JSON Flags.
type JSONTranscoder struct {
}

// NewJSONTranscoder returns a new JSONTranscoder.
func NewJSONTranscoder() *JSONTranscoder {
	return &JSONTranscoder{}
}

// Decode applies JSON transcoding behaviour to decode into a Go type.
func (t *JSONTranscoder) Decode(bytes []byte, flags uint32, out interface{}) error {
	valueType, compression := gocbcore.DecodeCommonFlags(flags)

	// Make sure compression is disabled
	if compression != gocbcore.NoCompression {
		return errors.New("unexpected value compression")
	}

	// Normal types of decoding
	if valueType == gocbcore.BinaryType {
		return errors.New("binary datatype is not supported by JSONTranscoder")
	} else if valueType == gocbcore.StringType {
		return errors.New("string datatype is not supported by JSONTranscoder")
	} else if valueType == gocbcore.JSONType {
		err := json.Unmarshal(bytes, &out)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("unexpected expectedFlags value")
}

// Encode applies JSON transcoding behaviour to encode a Go type.
func (t *JSONTranscoder) Encode(value interface{}) ([]byte, uint32, error) {
	var bytes []byte
	var flags uint32
	var err error

	switch typeValue := value.(type) {
	case []byte:
		return nil, 0, errors.New("binary data is not supported by JSONTranscoder")
	case *[]byte:
		return nil, 0, errors.New("binary data is not supported by JSONTranscoder")
	case json.RawMessage:
		bytes = typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	case *json.RawMessage:
		bytes = *typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	case *interface{}:
		return t.Encode(*typeValue)
	default:
		bytes, err = json.Marshal(value)
		if err != nil {
			return nil, 0, err
		}
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	}

	// No compression supported currently

	return bytes, flags, nil
}

// RawJSONTranscoder implements passthrough behavior of JSON data. This transcoder does not apply any serialization.
// It will forward data across the network without incurring unnecessary parsing costs.
//
// This will apply the following behavior to the value:
// binary ([]byte) -> JSON bytes, JSON expectedFlags.
// string -> JSON bytes, JSON expectedFlags.
// default -> error.
type RawJSONTranscoder struct {
}

// NewRawJSONTranscoder returns a new RawJSONTranscoder.
func NewRawJSONTranscoder() *RawJSONTranscoder {
	return &RawJSONTranscoder{}
}

// Decode applies raw JSON transcoding behaviour to decode into a Go type.
func (t *RawJSONTranscoder) Decode(bytes []byte, flags uint32, out interface{}) error {
	valueType, compression := gocbcore.DecodeCommonFlags(flags)

	// Make sure compression is disabled
	if compression != gocbcore.NoCompression {
		return errors.New("unexpected value compression")
	}

	// Normal types of decoding
	if valueType == gocbcore.BinaryType {
		return errors.New("binary datatype is not supported by RawJSONTranscoder")
	} else if valueType == gocbcore.StringType {
		return errors.New("string datatype is not supported by RawJSONTranscoder")
	} else if valueType == gocbcore.JSONType {
		switch typedOut := out.(type) {
		case *[]byte:
			*typedOut = bytes
			return nil
		case *string:
			*typedOut = string(bytes)
			return nil
		default:
			return errors.New("you must encode raw JSON data in a byte array or string")
		}
	}

	return errors.New("unexpected expectedFlags value")
}

// Encode applies raw JSON transcoding behaviour to encode a Go type.
func (t *RawJSONTranscoder) Encode(value interface{}) ([]byte, uint32, error) {
	var bytes []byte
	var flags uint32

	switch typeValue := value.(type) {
	case []byte:
		bytes = typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	case *[]byte:
		bytes = *typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	case string:
		bytes = []byte(typeValue)
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	case *string:
		bytes = []byte(*typeValue)
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	case json.RawMessage:
		bytes = typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	case *json.RawMessage:
		bytes = *typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	case *interface{}:
		return t.Encode(*typeValue)
	default:
		return nil, 0, makeInvalidArgumentsError("only binary and string data is supported by RawJSONTranscoder")
	}

	// No compression supported currently

	return bytes, flags, nil
}

// RawStringTranscoder implements passthrough behavior of raw string data. This transcoder does not apply any serialization.
//
// This will apply the following behavior to the value:
// string -> string bytes, string expectedFlags.
// default -> error.
type RawStringTranscoder struct {
}

// NewRawStringTranscoder returns a new RawStringTranscoder.
func NewRawStringTranscoder() *RawStringTranscoder {
	return &RawStringTranscoder{}
}

// Decode applies raw string transcoding behaviour to decode into a Go type.
func (t *RawStringTranscoder) Decode(bytes []byte, flags uint32, out interface{}) error {
	valueType, compression := gocbcore.DecodeCommonFlags(flags)

	// Make sure compression is disabled
	if compression != gocbcore.NoCompression {
		return errors.New("unexpected value compression")
	}

	// Normal types of decoding
	if valueType == gocbcore.BinaryType {
		return errors.New("only string datatype is supported by RawStringTranscoder")
	} else if valueType == gocbcore.StringType {
		switch typedOut := out.(type) {
		case *string:
			*typedOut = string(bytes)
			return nil
		case *interface{}:
			*typedOut = string(bytes)
			return nil
		default:
			return errors.New("you must encode a string in a string or interface")
		}
	} else if valueType == gocbcore.JSONType {
		return errors.New("only string datatype is supported by RawStringTranscoder")
	}

	return errors.New("unexpected expectedFlags value")
}

// Encode applies raw string transcoding behaviour to encode a Go type.
func (t *RawStringTranscoder) Encode(value interface{}) ([]byte, uint32, error) {
	var bytes []byte
	var flags uint32

	switch typeValue := value.(type) {
	case string:
		bytes = []byte(typeValue)
		flags = gocbcore.EncodeCommonFlags(gocbcore.StringType, gocbcore.NoCompression)
	case *string:
		bytes = []byte(*typeValue)
		flags = gocbcore.EncodeCommonFlags(gocbcore.StringType, gocbcore.NoCompression)
	case *interface{}:
		return t.Encode(*typeValue)
	default:
		return nil, 0, makeInvalidArgumentsError("only raw string data is supported by RawStringTranscoder")
	}

	// No compression supported currently

	return bytes, flags, nil
}

// RawBinaryTranscoder implements passthrough behavior of raw binary data. This transcoder does not apply any serialization.
//
// This will apply the following behavior to the value:
// binary ([]byte) -> binary bytes, binary expectedFlags.
// default -> error.
type RawBinaryTranscoder struct {
}

// NewRawBinaryTranscoder returns a new RawBinaryTranscoder.
func NewRawBinaryTranscoder() *RawBinaryTranscoder {
	return &RawBinaryTranscoder{}
}

// Decode applies raw binary transcoding behaviour to decode into a Go type.
func (t *RawBinaryTranscoder) Decode(bytes []byte, flags uint32, out interface{}) error {
	valueType, compression := gocbcore.DecodeCommonFlags(flags)

	// Make sure compression is disabled
	if compression != gocbcore.NoCompression {
		return errors.New("unexpected value compression")
	}

	// Normal types of decoding
	if valueType == gocbcore.BinaryType {
		switch typedOut := out.(type) {
		case *[]byte:
			*typedOut = bytes
			return nil
		case *interface{}:
			*typedOut = bytes
			return nil
		default:
			return errors.New("you must encode binary in a byte array or interface")
		}
	} else if valueType == gocbcore.StringType {
		return errors.New("only binary datatype is supported by RawBinaryTranscoder")
	} else if valueType == gocbcore.JSONType {
		return errors.New("only binary datatype is supported by RawBinaryTranscoder")
	}

	return errors.New("unexpected expectedFlags value")
}

// Encode applies raw binary transcoding behaviour to encode a Go type.
func (t *RawBinaryTranscoder) Encode(value interface{}) ([]byte, uint32, error) {
	var bytes []byte
	var flags uint32

	switch typeValue := value.(type) {
	case []byte:
		bytes = typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.BinaryType, gocbcore.NoCompression)
	case *[]byte:
		bytes = *typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.BinaryType, gocbcore.NoCompression)
	case *interface{}:
		return t.Encode(*typeValue)
	default:
		return nil, 0, makeInvalidArgumentsError("only raw binary data is supported by RawBinaryTranscoder")
	}

	// No compression supported currently

	return bytes, flags, nil
}

// LegacyTranscoder implements the behaviour for a backward-compatible transcoder. This transcoder implements
// behaviour matching that of gocb v1.
//
// This will apply the following behavior to the value:
// binary ([]byte) -> binary bytes, Binary expectedFlags.
// string -> string bytes, String expectedFlags.
// default -> JSON value, JSON expectedFlags.
type LegacyTranscoder struct {
}

// NewLegacyTranscoder returns a new LegacyTranscoder.
func NewLegacyTranscoder() *LegacyTranscoder {
	return &LegacyTranscoder{}
}

// Decode applies legacy transcoding behaviour to decode into a Go type.
func (t *LegacyTranscoder) Decode(bytes []byte, flags uint32, out interface{}) error {
	valueType, compression := gocbcore.DecodeCommonFlags(flags)

	// Make sure compression is disabled
	if compression != gocbcore.NoCompression {
		return errors.New("unexpected value compression")
	}

	// Normal types of decoding
	if valueType == gocbcore.BinaryType {
		switch typedOut := out.(type) {
		case *[]byte:
			*typedOut = bytes
			return nil
		case *interface{}:
			*typedOut = bytes
			return nil
		default:
			return errors.New("you must encode binary in a byte array or interface")
		}
	} else if valueType == gocbcore.StringType {
		switch typedOut := out.(type) {
		case *string:
			*typedOut = string(bytes)
			return nil
		case *interface{}:
			*typedOut = string(bytes)
			return nil
		default:
			return errors.New("you must encode a string in a string or interface")
		}
	} else if valueType == gocbcore.JSONType {
		err := json.Unmarshal(bytes, &out)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("unexpected expectedFlags value")
}

// Encode applies legacy transcoding behavior to encode a Go type.
func (t *LegacyTranscoder) Encode(value interface{}) ([]byte, uint32, error) {
	var bytes []byte
	var flags uint32
	var err error

	switch typeValue := value.(type) {
	case []byte:
		bytes = typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.BinaryType, gocbcore.NoCompression)
	case *[]byte:
		bytes = *typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.BinaryType, gocbcore.NoCompression)
	case string:
		bytes = []byte(typeValue)
		flags = gocbcore.EncodeCommonFlags(gocbcore.StringType, gocbcore.NoCompression)
	case *string:
		bytes = []byte(*typeValue)
		flags = gocbcore.EncodeCommonFlags(gocbcore.StringType, gocbcore.NoCompression)
	case json.RawMessage:
		bytes = typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	case *json.RawMessage:
		bytes = *typeValue
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	case *interface{}:
		return t.Encode(*typeValue)
	default:
		bytes, err = json.Marshal(value)
		if err != nil {
			return nil, 0, err
		}
		flags = gocbcore.EncodeCommonFlags(gocbcore.JSONType, gocbcore.NoCompression)
	}

	// No compression supported currently

	return bytes, flags, nil
}
