package mccockroachdb

import (
	log "github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPathPrivate(t *testing.T) {
	c := MultiClusterCockroachDB{logger: log.New(nil)}
	tests := []struct {
		give string
		want bool
	}{
		{
			give: "logical/712b4f75-32dd-20e6-ba37-cc76eeffb5d8/94f53cc0-3b94-d2ea-f5f4-c61496453876/policy/metadata",
			want: false,
		},
		{
			give: "sys/policy/default",
			want: false,
		},
		{
			give: "core/cluster/local/info",
			want: false,
		},
		{
			give: "counters/etc",
			want: true,
		},
		{
			give: "leader/etc",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			res := c.isPrivate(tt.give)
			assert.Equal(t, tt.want, res)
		})
	}
}
