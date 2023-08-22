package vault

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

// transitoryGeneration is the Generation value used by SealWrappebValues for
// entries that need to be upgraded from pre Vault 1.15.
// Note that the value is 0, since that is the default value for BlobInfo.Generation.
const transitoryGeneration = 0

type SealWrappedValue struct {
	value wrapping.MultiWrapValue
}

// NewSealWrappedValue creates a new seal wrapped value. Note that this
// method will change to accept a slice of BlobInfos when multi-seal wrapping
// is added.
func NewSealWrappedValue(blobInfo *wrapping.BlobInfo) *SealWrappedValue {
	if blobInfo == nil {
		panic("cannot create a SealWrappedValue without a BlobInfo")
	}

	// TODO(victorr): see if we can use plaintext instead of blobInfo.Wrapped = false
	return &SealWrappedValue{
		value: wrapping.MultiWrapValue{
			Generation: 1, // TODO(SEALHA): Fill in Generation
			Slots:      []*wrapping.BlobInfo{blobInfo},
		},
	}
}

func newTransitorySealWrappedValue(blobInfo *wrapping.BlobInfo) *SealWrappedValue {
	ret := NewSealWrappedValue(blobInfo)
	ret.value.Generation = transitoryGeneration

	return ret
}

func (swv *SealWrappedValue) isTransitory() bool {
	return swv.value.Generation == transitoryGeneration
}

func (swv *SealWrappedValue) isPlaintext() bool {
	return !swv.isEncrypted()
}

// isEncrypted returns true if a BlobInfo has flag Wrapped set to true.
func (swv *SealWrappedValue) isEncrypted() bool {
	// Note that we set Wrapped == false for StoredBarrierKeysPath and recoveryKeyPath, so check
	// for the presence of KeyInfo as well as the Wrapped flag.

	blobInfo := swv.GetUniqueBlobInfo()
	return blobInfo.Wrapped || (blobInfo.KeyInfo != nil)
}

// GetUniqueBlobInfo returns the BlobInfo for the seal wrapped value. This
// method panics if there is more than one BlobInfo.
// TODO(SEALHA): Remove GetUniqueBlobInfo() once all callers can work with multiple encryptions.
func (swv *SealWrappedValue) GetUniqueBlobInfo() *wrapping.BlobInfo {
	if len(swv.value.Slots) > 1 {
		panic("expected exactly one BlobInfo but there are more than one of them")
	}

	return swv.value.Slots[0]
}

func (swv *SealWrappedValue) getPlaintextValue() ([]byte, error) {
	if swv.isEncrypted() {
		return nil, errors.New("cannot return plaintext value from a SealWrappedValue with encrypted data")
	}

	return swv.GetUniqueBlobInfo().Ciphertext, nil
}

var sealWrappedValueHeader = []byte("multiwrapvalue:1")

const sealWrappedValueHeaderLength = 16

func init() {
	// Check that the header is 16 bytes long.
	if len(sealWrappedValueHeader) != sealWrappedValueHeaderLength {
		panic(fmt.Sprintf("sealWrappedValueHeader must be %d bytes long, but it is %d", sealWrappedValueHeaderLength, len(sealWrappedValueHeader)))
	}
}

// Marshal a seal wrapped value. DO NOT USE DIRECTLY, use MarshalSealWrappedValue instead.
// The marshalled bytes consists of:
// a) A 16 byte header
// b) 4 bytes specifying the length of the remaining bytes
// c) the protobuf marshalling of the MultiWrapValue
func (swv *SealWrappedValue) marshal() ([]byte, error) {
	protoBytes, err := proto.Marshal(&swv.value)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	var appendErr error
	append := func(value []byte) {
		if appendErr == nil {
			_, appendErr = buf.Write(value)
		}
	}

	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(protoBytes)))

	append(sealWrappedValueHeader)
	append(lengthBytes)
	append(protoBytes)
	if appendErr != nil {
		return nil, appendErr
	}

	return buf.Bytes(), nil
}

// Unmarshal a seal wrapped value. DO NOT USE DIRECTLY, use UnmarshalSealWrappedValue instead.
func (swv *SealWrappedValue) unmarshal(value []byte) error {
	if len(value) < sealWrappedValueHeaderLength+4 {
		return errors.New("error unmarshalling SealWrappedValue, not enough bytes")
	}

	header := value[0:sealWrappedValueHeaderLength]
	lengthBytes := value[sealWrappedValueHeaderLength : sealWrappedValueHeaderLength+4]
	protoBytes := value[sealWrappedValueHeaderLength+4:]

	if bytes.Compare(sealWrappedValueHeader, header) != 0 {
		return errors.New("error unmarshalling SealWrappedValue, header mismatch")
	}
	length := binary.BigEndian.Uint32(lengthBytes)
	if int(length) != len(protoBytes) {
		return errors.New("error unmarshalling SealWrappedValue, length mismatch")
	}

	if err := proto.Unmarshal(protoBytes, &swv.value); err != nil {
		return err
	}

	return nil
}
