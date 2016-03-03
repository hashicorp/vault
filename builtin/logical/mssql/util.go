package mssql

import (
	"fmt"
	"strings"
)

// BuildDsn creates a DSN with some default and overriden values
func BuildDsn(dsn string) string {
	dsnParts := make(map[string]string)
	parts := strings.Split(dsn, ";")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		lst := strings.SplitN(part, "=", 2)
		name := strings.TrimSpace(strings.ToLower(lst[0]))
		if len(name) == 0 {
			continue
		}
		var value string
		if len(lst) > 1 {
			value = strings.TrimSpace(lst[1])
		}
		dsnParts[name] = value
	}

	// Default app name to vault
	if _, exists := dsnParts["app name"]; !exists {
		dsnParts["app name"] = "vault"
	}

	var newDsn string
	for k, v := range dsnParts {
		newDsn = newDsn + k + "=" + v + ";"
	}
	return newDsn
}

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
