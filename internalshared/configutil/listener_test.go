// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSingleIPTemplate(t *testing.T) {
	type args struct {
		ipTmpl string
	}
	tests := []struct {
		name    string
		arg     string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "test https addr",
			arg:     "https://vaultproject.io:8200",
			want:    "https://vaultproject.io:8200",
			wantErr: assert.NoError,
		},
		{
			name:    "test invalid template func",
			arg:     "{{FooBar}}",
			want:    "",
			wantErr: assert.Error,
		},
		{
			name:    "test partial template",
			arg:     "{{FooBar",
			want:    "{{FooBar",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSingleIPTemplate(tt.arg)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseSingleIPTemplate(%v)", tt.arg)) {
				return
			}

			assert.Equalf(t, tt.want, got, "ParseSingleIPTemplate(%v)", tt.arg)
		})
	}
}
