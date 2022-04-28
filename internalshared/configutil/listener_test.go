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
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "test non-template addr",
			args:    args{"vaultproject.io"},
			want:    "vaultproject.io",
			wantErr: assert.NoError,
		},
		{
			name:    "test invalid template func",
			args:    args{"{{FooBar}}"},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name:    "test partial template",
			args:    args{"{{FooBar"},
			want:    "{{FooBar",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSingleIPTemplate(tt.args.ipTmpl)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseSingleIPTemplate(%v)", tt.args.ipTmpl)) {
				return
			}

			assert.Equalf(t, tt.want, got, "ParseSingleIPTemplate(%v)", tt.args.ipTmpl)
		})
	}
}
