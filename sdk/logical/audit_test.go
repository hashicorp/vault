// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLogInput_BexprDatum ensures that we can transform a LogInput
// into a LogInputBexpr to be used in audit filtering.
func TestLogInput_BexprDatum(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Request            *Request
		Namespace          string
		ExpectedPath       string
		ExpectedMountPoint string
		ExpectedMountType  string
		ExpectedNamespace  string
		ExpectedOperation  string
	}{
		"nil-no-namespace": {
			Request:            nil,
			Namespace:          "",
			ExpectedPath:       "",
			ExpectedMountPoint: "",
			ExpectedMountType:  "",
			ExpectedNamespace:  "",
			ExpectedOperation:  "",
		},
		"nil-namespace": {
			Request:            nil,
			Namespace:          "juan",
			ExpectedPath:       "",
			ExpectedMountPoint: "",
			ExpectedMountType:  "",
			ExpectedNamespace:  "juan",
			ExpectedOperation:  "",
		},
		"happy-path": {
			Request: &Request{
				MountPoint: "IAmAMountPoint",
				MountType:  "IAmAMountType",
				Operation:  CreateOperation,
				Path:       "IAmAPath",
			},
			Namespace:          "juan",
			ExpectedPath:       "IAmAPath",
			ExpectedMountPoint: "IAmAMountPoint",
			ExpectedMountType:  "IAmAMountType",
			ExpectedNamespace:  "juan",
			ExpectedOperation:  "create",
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l := &LogInput{Request: tc.Request}

			d := l.BexprDatum(tc.Namespace)

			require.Equal(t, tc.ExpectedPath, d.Path)
			require.Equal(t, tc.ExpectedMountPoint, d.MountPoint)
			require.Equal(t, tc.ExpectedMountType, d.MountType)
			require.Equal(t, tc.ExpectedNamespace, d.Namespace)
			require.Equal(t, tc.ExpectedOperation, d.Operation)
		})
	}
}
