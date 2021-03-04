package mstypes

import (
	"bytes"
	"errors"

	"gopkg.in/jcmturner/rpc.v1/ndr"
)

// Compression format assigned numbers.
const (
	CompressionFormatNone       uint16 = 0
	CompressionFormatLZNT1      uint16 = 2
	CompressionFormatXPress     uint16 = 3
	CompressionFormatXPressHuff uint16 = 4
)

// ClaimsSourceTypeAD https://msdn.microsoft.com/en-us/library/hh553809.aspx
const ClaimsSourceTypeAD uint16 = 1

// Claim Type assigned numbers
const (
	ClaimTypeIDInt64    uint16 = 1
	ClaimTypeIDUInt64   uint16 = 2
	ClaimTypeIDString   uint16 = 3
	ClaimsTypeIDBoolean uint16 = 6
)

// ClaimsBlob implements https://msdn.microsoft.com/en-us/library/hh554119.aspx
type ClaimsBlob struct {
	Size        uint32
	EncodedBlob EncodedBlob
}

// EncodedBlob are the bytes of the encoded Claims
type EncodedBlob []byte

// Size returns the size of the bytes of the encoded Claims
func (b EncodedBlob) Size(c interface{}) int {
	cb := c.(ClaimsBlob)
	return int(cb.Size)
}

// ClaimsSetMetadata implements https://msdn.microsoft.com/en-us/library/hh554073.aspx
type ClaimsSetMetadata struct {
	ClaimsSetSize             uint32
	ClaimsSetBytes            []byte `ndr:"pointer,conformant"`
	CompressionFormat         uint16 // Enum see constants for options
	UncompressedClaimsSetSize uint32
	ReservedType              uint16
	ReservedFieldSize         uint32
	ReservedField             []byte `ndr:"pointer,conformant"`
}

// ClaimsSet reads the ClaimsSet type from the NDR encoded ClaimsSetBytes in the ClaimsSetMetadata
func (m *ClaimsSetMetadata) ClaimsSet() (c ClaimsSet, err error) {
	if len(m.ClaimsSetBytes) < 1 {
		err = errors.New("no bytes available for ClaimsSet")
		return
	}
	// TODO switch statement to decompress ClaimsSetBytes
	if m.CompressionFormat != CompressionFormatNone {
		err = errors.New("compressed ClaimsSet not currently supported")
		return
	}
	dec := ndr.NewDecoder(bytes.NewReader(m.ClaimsSetBytes))
	err = dec.Decode(&c)
	return
}

// ClaimsSet implements https://msdn.microsoft.com/en-us/library/hh554122.aspx
type ClaimsSet struct {
	ClaimsArrayCount  uint32
	ClaimsArrays      []ClaimsArray `ndr:"pointer,conformant"`
	ReservedType      uint16
	ReservedFieldSize uint32
	ReservedField     []byte `ndr:"pointer,conformant"`
}

// ClaimsArray implements https://msdn.microsoft.com/en-us/library/hh536458.aspx
type ClaimsArray struct {
	ClaimsSourceType uint16
	ClaimsCount      uint32
	ClaimEntries     []ClaimEntry `ndr:"pointer,conformant"`
}

// ClaimEntry is a NDR union that implements https://msdn.microsoft.com/en-us/library/hh536374.aspx
type ClaimEntry struct {
	ID         string           `ndr:"pointer,conformant,varying"`
	Type       uint16           `ndr:"unionTag"`
	TypeInt64  ClaimTypeInt64   `ndr:"unionField"`
	TypeUInt64 ClaimTypeUInt64  `ndr:"unionField"`
	TypeString ClaimTypeString  `ndr:"unionField"`
	TypeBool   ClaimTypeBoolean `ndr:"unionField"`
}

// SwitchFunc is the ClaimEntry union field selection function
func (u ClaimEntry) SwitchFunc(_ interface{}) string {
	switch u.Type {
	case ClaimTypeIDInt64:
		return "TypeInt64"
	case ClaimTypeIDUInt64:
		return "TypeUInt64"
	case ClaimTypeIDString:
		return "TypeString"
	case ClaimsTypeIDBoolean:
		return "TypeBool"
	}
	return ""
}

// ClaimTypeInt64 is a claim of type int64
type ClaimTypeInt64 struct {
	ValueCount uint32
	Value      []int64 `ndr:"pointer,conformant"`
}

// ClaimTypeUInt64 is a claim of type uint64
type ClaimTypeUInt64 struct {
	ValueCount uint32
	Value      []uint64 `ndr:"pointer,conformant"`
}

// ClaimTypeString is a claim of type string
type ClaimTypeString struct {
	ValueCount uint32
	Value      []LPWSTR `ndr:"pointer,conformant"`
}

// ClaimTypeBoolean is a claim of type bool
type ClaimTypeBoolean struct {
	ValueCount uint32
	Value      []bool `ndr:"pointer,conformant"`
}
