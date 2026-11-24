// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestStatusColumnValues verifies the column_values JSON built for the initial
// status: an empty status yields an empty object so no columns are set, a
// status is written as a label on the "status" column, the label is passed
// through exactly as given (Monday.com labels are case-sensitive), and the
// result is always valid JSON that round-trips the original label.
func TestStatusColumnValues(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		status string
		want   string
	}{
		"empty status sends an empty object": {status: "", want: `{}`},
		"label is set on the status column":  {status: "Not started", want: `{"status":{"label":"Not started"}}`},
		"label case is preserved":            {status: "Not Started", want: `{"status":{"label":"Not Started"}}`},
		"surrounding whitespace is kept":     {status: " Done ", want: `{"status":{"label":" Done "}}`},
		"quotes are escaped":                 {status: `Say "hi"`, want: `{"status":{"label":"Say \"hi\""}}`},
		"backslashes are escaped":            {status: `a\b`, want: `{"status":{"label":"a\\b"}}`},
		"unicode is preserved":               {status: "Prêt", want: `{"status":{"label":"Prêt"}}`},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := statusColumnValues(tc.status)
			require.Equal(t, columnValuesJSON(tc.want), got)
			require.True(t, json.Valid([]byte(got)), "column_values must be valid JSON")

			if tc.status == "" {
				return
			}
			var decoded map[string]map[string]string
			require.NoError(t, json.Unmarshal([]byte(got), &decoded))
			require.Equal(t, tc.status, decoded["status"]["label"])
		})
	}
}

// TestColumnValuesJSON_GetGraphQLType verifies the type name hasura emits in
// the operation's variable declaration. It must be JSON, not the String it
// would otherwise infer from the underlying Go string.
func TestColumnValuesJSON_GetGraphQLType(t *testing.T) {
	t.Parallel()

	require.Equal(t, "JSON", columnValuesJSON("{}").GetGraphQLType())
}
