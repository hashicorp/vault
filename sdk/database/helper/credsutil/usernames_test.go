// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package credsutil

import (
	"regexp"
	"testing"
)

func TestGenerateUsername(t *testing.T) {
	type testCase struct {
		displayName    string
		displayNameLen int

		roleName    string
		roleNameLen int

		usernameLen int
		separator   string
		caseOp      CaseOp

		regex string
	}
	tests := map[string]testCase{
		"all opts": {
			displayName:    "abcdefghijklmonpqrstuvwxyz",
			displayNameLen: 10,
			roleName:       "zyxwvutsrqpnomlkjihgfedcba",
			roleNameLen:    10,
			usernameLen:    45,
			separator:      ".",
			caseOp:         KeepCase,

			regex: "^v.abcdefghij.zyxwvutsrq.[a-zA-Z0-9]{20}.$",
		},
		"no separator": {
			displayName:    "abcdefghijklmonpqrstuvwxyz",
			displayNameLen: 10,
			roleName:       "zyxwvutsrqpnomlkjihgfedcba",
			roleNameLen:    10,
			usernameLen:    45,
			separator:      "",
			caseOp:         KeepCase,

			regex: "^vabcdefghijzyxwvutsrq[a-zA-Z0-9]{20}[0-9]{4}$",
		},
		"lowercase": {
			displayName:    "abcdefghijklmonpqrstuvwxyz",
			displayNameLen: 10,
			roleName:       "zyxwvutsrqpnomlkjihgfedcba",
			roleNameLen:    10,
			usernameLen:    45,
			separator:      "_",
			caseOp:         Lowercase,

			regex: "^v_abcdefghij_zyxwvutsrq_[a-z0-9]{20}_$",
		},
		"uppercase": {
			displayName:    "abcdefghijklmonpqrstuvwxyz",
			displayNameLen: 10,
			roleName:       "zyxwvutsrqpnomlkjihgfedcba",
			roleNameLen:    10,
			usernameLen:    45,
			separator:      "_",
			caseOp:         Uppercase,

			regex: "^V_ABCDEFGHIJ_ZYXWVUTSRQ_[A-Z0-9]{20}_$",
		},
		"short username": {
			displayName:    "abcdefghijklmonpqrstuvwxyz",
			displayNameLen: 5,
			roleName:       "zyxwvutsrqpnomlkjihgfedcba",
			roleNameLen:    5,
			usernameLen:    15,
			separator:      "_",
			caseOp:         KeepCase,

			regex: "^v_abcde_zyxwv_[a-zA-Z0-9]{1}$",
		},
		"long username": {
			displayName:    "abcdefghijklmonpqrstuvwxyz",
			displayNameLen: 0,
			roleName:       "zyxwvutsrqpnomlkjihgfedcba",
			roleNameLen:    0,
			usernameLen:    100,
			separator:      "_",
			caseOp:         KeepCase,

			regex: "^v_abcdefghijklmonpqrstuvwxyz_zyxwvutsrqpnomlkjihgfedcba_[a-zA-Z0-9]{20}_[0-9]{1,23}$",
		},
		"zero max length": {
			displayName:    "abcdefghijklmonpqrstuvwxyz",
			displayNameLen: 0,
			roleName:       "zyxwvutsrqpnomlkjihgfedcba",
			roleNameLen:    0,
			usernameLen:    0,
			separator:      "_",
			caseOp:         KeepCase,

			regex: "^v_abcdefghijklmonpqrstuvwxyz_zyxwvutsrqpnomlkjihgfedcba_[a-zA-Z0-9]{20}_[0-9]+$",
		},
		"no display name": {
			displayName:    "abcdefghijklmonpqrstuvwxyz",
			displayNameLen: NoneLength,
			roleName:       "zyxwvutsrqpnomlkjihgfedcba",
			roleNameLen:    15,
			usernameLen:    100,
			separator:      "_",
			caseOp:         KeepCase,

			regex: "^v_zyxwvutsrqpnoml_[a-zA-Z0-9]{20}_[0-9]+$",
		},
		"no role name": {
			displayName:    "abcdefghijklmonpqrstuvwxyz",
			displayNameLen: 15,
			roleName:       "zyxwvutsrqpnomlkjihgfedcba",
			roleNameLen:    NoneLength,
			usernameLen:    100,
			separator:      "_",
			caseOp:         KeepCase,

			regex: "^v_abcdefghijklmon_[a-zA-Z0-9]{20}_[0-9]+$",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			re := regexp.MustCompile(test.regex)

			username, err := GenerateUsername(
				DisplayName(test.displayName, test.displayNameLen),
				RoleName(test.roleName, test.roleNameLen),
				Separator(test.separator),
				MaxLength(test.usernameLen),
				Case(test.caseOp),
			)
			if err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !re.MatchString(username) {
				t.Fatalf("username %q does not match regex %q", username, test.regex)
			}
		})
	}
}
