package dbutil

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyCreationStatement = errors.New("empty creation statements")
)

// Query templates a query for us.
func QueryHelper(tpl string, data map[string]string) string {
	for k, v := range data {
		tpl = strings.Replace(tpl, fmt.Sprintf("{{%s}}", k), v, -1)
	}

	return tpl
}
