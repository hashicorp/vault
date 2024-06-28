// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package physical

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

func TestTransactionalErrorInjector_TransactionLimits(t *testing.T) {
	tc := []struct {
		name        string
		be          Backend
		wantEntries int
		wantSize    int
	}{
		{
			name: "non-transactionlimits backend",
			be:   &TestTransactionalNonLimitBackend{},

			// Should return zeros to let the implementor choose defaults.
			wantEntries: 0,
			wantSize:    0,
		},
		{
			name: "transactionlimits backend",
			be: &TestTransactionalLimitBackend{
				MaxEntries: 123,
				MaxSize:    345,
			},

			// Should return underlying limits
			wantEntries: 123,
			wantSize:    345,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			logger := hclog.NewNullLogger()

			injector := NewTransactionalErrorInjector(tt.be, 0, logger)

			maxEntries, maxBytes := injector.TransactionLimits()

			require.Equal(t, tt.wantEntries, maxEntries)
			require.Equal(t, tt.wantSize, maxBytes)
		})
	}
}
