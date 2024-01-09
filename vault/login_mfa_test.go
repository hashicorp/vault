// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"strings"
	"testing"
)

func TestParseFactors(t *testing.T) {
	testcases := []struct {
		name                string
		invalidMFAHeaderVal []string
		expectedError       string
	}{
		{
			"two headers with passcode",
			[]string{"passcode", "foo"},
			"found multiple passcodes for the same MFA method",
		},
		{
			"single header with passcode=",
			[]string{"passcode="},
			"invalid passcode",
		},
		{
			"single invalid header",
			[]string{"foo="},
			"found an invalid MFA cred",
		},
		{
			"single header equal char",
			[]string{"=="},
			"found an invalid MFA cred",
		},
		{
			"two headers with passcode=",
			[]string{"passcode=foo", "foo"},
			"found multiple passcodes for the same MFA method",
		},
		{
			"two headers invalid name",
			[]string{"passcode=foo", "passcode=bar"},
			"found multiple passcodes for the same MFA method",
		},
		{
			"two headers, two invalid",
			[]string{"foo", "bar"},
			"found multiple passcodes for the same MFA method",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseMfaFactors(tc.invalidMFAHeaderVal)
			if err == nil {
				t.Fatal("nil error returned")
			}
			if !strings.Contains(err.Error(), tc.expectedError) {
				t.Fatalf("expected %s, got %v", tc.expectedError, err)
			}
		})
	}
}
