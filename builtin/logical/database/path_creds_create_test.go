// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"testing"
)

// TestSanitizeDisplayName tests that display names are sanitized before use.
func TestSanitizeDisplayName(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"empty": {
			input: "",
			want:  "",
		},
		"plain alphanumeric is unchanged": {
			input: "alice42",
			want:  "alice42",
		},
		"hyphen and underscore are preserved": {
			input: "alice-bob_carol",
			want:  "alice-bob_carol",
		},
		"spaces become hyphens": {
			input: "alice bob",
			want:  "alice-bob",
		},
		"single quote becomes hyphen": {
			input: "o'brien",
			want:  "o-brien",
		},
		"semicolon and quotes are neutralized": {
			input: `a';DROP`,
			want:  "a--DROP",
		},
		"plpgsql injection payload is fully neutralized": {
			input: `x'; DO $$ BEGIN END $$;--`,
			want:  "x---DO----BEGIN-END------",
		},
		"newline and tab become hyphens": {
			input: "a\nb\tc",
			want:  "a-b-c",
		},
		"non-ascii is replaced": {
			input: "naïve",
			want:  "na-ve",
		},

		"mysql @@version becomes --version": {
			input: "@@version",
			want:  "--version",
		},

		"mysql backtick identifier": {
			input: "`users`",
			want:  "-users-",
		},

		"mysql comment injection": {
			input: "admin`-- ",
			want:  "admin----",
		},

		"mssql xp_cmdshell is unchanged": {
			input: "xp_cmdshell",
			want:  "xp_cmdshell",
		},

		"mssql square bracket delimiters": {
			input: "[dbo].[users]",
			want:  "-dbo---users-",
		},

		"mssql variable @@SERVERNAME": {
			input: "@@SERVERNAME",
			want:  "--SERVERNAME",
		},

		"300-char null bytes all become hyphens": {
			input: string(make([]byte, 300)), // 300 null bytes — all outside allowlist
			want: func() string {
				b := make([]byte, 300)
				for i := range b {
					b[i] = '-'
				}
				return string(b)
			}(),
		},
		"300-char alphanumeric is unchanged": {
			input: func() string {
				b := make([]byte, 300)
				for i := range b {
					b[i] = 'a'
				}
				return string(b)
			}(),
			want: func() string {
				b := make([]byte, 300)
				for i := range b {
					b[i] = 'a'
				}
				return string(b)
			}(),
		},

		"unicode fullwidth single quote ＇ is replaced": {
			input: "x\uff07 OR 1=1", // ＇ U+FF07
			want:  "x--OR-1-1",
		},
		"unicode fullwidth semicolon ； is replaced": {
			input: "x\uff1bDROP", // ； U+FF1B
			want:  "x-DROP",
		},
		"unicode fullwidth backtick ｀ is replaced": {
			input: "admin\uff60--", // ｀ U+FF40
			want:  "admin---",
		},

		"unicode cyrillic А (U+0410) is replaced": {
			input: "\u0410DMINxp_cmdshell", // looks like "ADMIN..." but Cyrillic А
			want:  "-DMINxp_cmdshell",
		},
		"unicode fake alice with cyrillic а and с": {
			input: "\u0430li\u0441e", // а=U+0430, с=U+0441
			want:  "-li-e",
		},

		"unicode zero-width space U+200B is replaced": {
			input: "abc\u200bdef",
			want:  "abc-def",
		},
		"unicode right-to-left override U+202E is replaced": {
			input: "abc\u202edef",
			want:  "abc-def",
		},
		"unicode null byte is replaced": {
			input: "abc\x00def",
			want:  "abc-def",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := sanitizeDisplayName(tc.input)
			if got != tc.want {
				t.Fatalf("sanitizeDisplayName(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
