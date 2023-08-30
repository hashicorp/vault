package vault

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/vault/seal"

	"github.com/golang/protobuf/proto"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/stretchr/testify/assert"
)

func TestMarshalSealWrappedValue(t *testing.T) {
	isBlobInfo := func(bytes []byte) bool {
		err := proto.Unmarshal(bytes, &wrapping.BlobInfo{})
		return err == nil
	}

	isSealWrapValue := func(bytes []byte) bool {
		if isBlobInfo(bytes) {
			return false
		}
		err := (&SealWrappedValue{}).unmarshal(bytes)
		return err == nil
	}

	blobInfo := &wrapping.BlobInfo{
		Wrapped:    false,
		Ciphertext: []byte("plaintext, actually"),
	}
	oneBlobInfo := []*wrapping.BlobInfo{blobInfo}
	twoBlobInfos := []*wrapping.BlobInfo{blobInfo, blobInfo}

	wantBlobInfo := true
	wantMultiWrappedValue := !wantBlobInfo

	tests := []struct {
		name    string
		value   *SealWrappedValue
		want    bool
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
			want:    wantBlobInfo,
			wantErr: assert.NoError,
		},
		{
			name: "two BlobInfos generation 0",
			value: &SealWrappedValue{
				value: seal.MultiWrapValue{
					Generation: 0,
					Slots:      twoBlobInfos,
				},
			},
			want:    wantMultiWrappedValue,
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
			want:    wantMultiWrappedValue,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalSealWrappedValue(tt.value)
			if !tt.wantErr(t, err, fmt.Sprintf("MarshalSealWrappedValue(%v)", tt.value)) {
				return
			}
			if tt.want == wantBlobInfo {
				assertTrue(t, isBlobInfo(got), "expecting bytes to be a marshalled BlobInfo")
			} else {
				assertTrue(t, isSealWrapValue(got), "expecting bytes to be a marshalled SealWrappedValue")
			}

			unmarshalled, err := UnmarshalSealWrappedValue(got)
			assert.NoError(t, err)
			assert.True(t, proto.Equal(&tt.value.value, &unmarshalled.value), "%v != %v", tt.value.value, unmarshalled.value)
		})
	}
}
