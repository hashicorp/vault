package mysql

import (
	"crypto/rand"
	"fmt"
	"strings"
)

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

// Query templates a query for us.
func Query(tpl string, data map[string]string) string {
	for k, v := range data {
		tpl = strings.Replace(tpl, fmt.Sprintf("{{%s}}", k), v, -1)
	}

	return tpl
}

// generateUUID is used to generate a random UUID
func generateUUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
}
