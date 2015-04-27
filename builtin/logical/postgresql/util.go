package postgresql

import "strings"

// SplitSQL is used to split a series of SQL statements
func SplitSQL(sql string) []string {
	parts := strings.Split(sql, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		clean := strings.TrimSpace(p)
		if len(clean) > 0 {
			out = append(out, clean)
		}
	}
	return out
}
