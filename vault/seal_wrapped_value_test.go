package vault

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/vault/seal"

	"github.com/golang/protobuf/proto"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/stretchr/testify/assert"
)

func TestSealWrappedValue_marshal_unmarshal(t *testing.T) {
	blobInfo := &wrapping.BlobInfo{
		Wrapped:    false,
		Ciphertext: []byte("plaintext, actually"),
	}
	oneBlobInfo := []*wrapping.BlobInfo{blobInfo}
	twoBlobInfos := []*wrapping.BlobInfo{blobInfo, blobInfo}

	tests := []struct {
		name    string
		value   *SealWrappedValue
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "a BlobInfo generation 0",
			value: &SealWrappedValue{
				value: seal.MultiWrapValue{
					Generation: 0,
					Slots:      oneBlobInfo,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "two BlobInfos generation 7",
			value: &SealWrappedValue{
				value: seal.MultiWrapValue{
					Generation: 7,
					Slots:      twoBlobInfos,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "two BlobInfos generation 42",
			value: &SealWrappedValue{
				value: seal.MultiWrapValue{
					Generation: 42,
					Slots:      twoBlobInfos,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.value.marshal()
			if !tt.wantErr(t, err, fmt.Sprintf("marshal()")) {
				return
			}

			unmarshalled := &SealWrappedValue{}
			assert.NoError(t, unmarshalled.unmarshal(got))
			assert.True(t, proto.Equal(&tt.value.value, &unmarshalled.value), "%v != %v", tt.value.value, unmarshalled.value)
		})
	}
}

func TestSealWrappedValue_unmarshalMultiWrapError_error_conditions(t *testing.T) {
	unmarshal := (&SealWrappedValue{}).unmarshal

	assert.EqualError(t, unmarshal([]byte{1, 2, 3}), "error unmarshalling SealWrappedValue, not enough bytes")

	swv := &SealWrappedValue{
		value: seal.MultiWrapValue{
			Generation: 0,
			Slots: []*wrapping.BlobInfo{
				{
					Wrapped:    false,
					Ciphertext: []byte("plaintext, actually"),
				},
			},
		},
	}

	bytes, err := swv.marshal()
	assert.NoError(t, err)
	assert.NoError(t, unmarshal(bytes))

	badHeader := []byte("oops")
	badHeader = append(badHeader, bytes[len(badHeader):]...)

	assert.EqualError(t, unmarshal(badHeader), "error unmarshalling SealWrappedValue, header mismatch")

	badLength := bytes[0:sealWrappedValueHeaderLength]
	badLength = append(badLength, 1, 2, 3, 4)
	badLength = append(badLength, bytes[sealWrappedValueHeaderLength+4:]...)

	assert.EqualError(t, unmarshal(badLength), "error unmarshalling SealWrappedValue, length mismatch")
}
