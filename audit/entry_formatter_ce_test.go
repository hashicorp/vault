// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package audit

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

// TestEntryFormatter_excludeFields tests that we can exclude data based on the
// pre-configured conditions/fields of the EntryFormatter. It covers some scenarios
// where we expect errors due to invalid input, which is unlikely to happen in reality.
func TestEntryFormatter_excludeFields(t *testing.T) {
	// Create the formatter node.
	cfg, err := newFormatterConfig(&testHeaderFormatter{}, nil)
	require.NoError(t, err)
	ss := newStaticSalt(t)

	// We intentionally create the EntryFormatter manually, as we wouldn't be
	// able to set exclusions via NewEntryFormatter WithExclusions option.
	formatter := &entryFormatter{
		config: cfg,
		salter: ss,
		logger: hclog.NewNullLogger(),
		name:   "juan",
	}

	res, err := formatter.excludeFields(nil)
	require.Error(t, err)
	require.EqualError(t, err, "enterprise-only feature: audit exclusion")
	require.Nil(t, res)
}
