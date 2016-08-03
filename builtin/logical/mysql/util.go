package mysql

import (
	"fmt"
	"strings"
)

// Query templates a query for us.
func Query(tpl string, data map[string]string) string {
	for k, v := range data {
		tpl = strings.Replace(tpl, fmt.Sprintf("{{%s}}", k), v, -1)
	}

	return tpl
}
