// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package parsing

import (
	"testing"

	"github.com/google/certificate-transparency-go/asn1"
	"github.com/stretchr/testify/require"
)

// TestAsn1UnmarshallNoTrailing tests the Asn1UnmarshallNoTrailing function returns
// errors as we expect if the input is not marshalled correctly or there is trailing
// data.
func TestAsn1UnmarshallNoTrailing(t *testing.T) {
	stringToMarshal := "a string"
	marshal, err := asn1.Marshal(stringToMarshal)
	require.NoError(t, err, "marshal failed")

	var myTestString string

	type args struct {
		b   []byte
		val any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"happy-path", args{marshal, &myTestString}, false},
		{"bad-marshalling", args{[]byte("incorrect"), &myTestString}, true},
		{"trailing-data", args{append(marshal, []byte("\n")...), &myTestString}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myTestString = ""
			err := Asn1UnmarshallNoTrailing(tt.args.b, tt.args.val)

			if (err != nil) != tt.wantErr {
				t.Errorf("Asn1UnmarshallNoTrailing() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if myTestString != stringToMarshal {
					t.Errorf("Asn1UnmarshallNoTrailing() = %v, want %v", myTestString, stringToMarshal)
				}
			}
		})
	}
}
