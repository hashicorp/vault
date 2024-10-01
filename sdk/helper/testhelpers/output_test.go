// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testhelpers

import (
	"fmt"
	"reflect"
	"testing"
)

func TestToMap(t *testing.T) {
	type s struct {
		A string            `json:"a"`
		B []byte            `json:"b"`
		C map[string]string `json:"c"`
		D string            `json:"-"`
	}
	type args struct {
		in s
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "basic",
			args:    args{s{A: "a", B: []byte("bytes"), C: map[string]string{"k": "v"}, D: "d"}},
			want:    "map[a:a b:277089d91c0bdf4f2e6862ba7e4a07605119431f5d13f726dd352b06f1b206a9 c:map[k:v]]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := ToMap(&tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := fmt.Sprintf("%s", m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
