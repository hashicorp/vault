package pgx

import (
	"strconv"
)

// QueryArgs is a container for arguments to an SQL query. It is helpful when
// building SQL statements where the number of arguments is variable.
type QueryArgs []interface{}

var placeholders []string

func init() {
	placeholders = make([]string, 64)

	for i := 1; i < 64; i++ {
		placeholders[i] = "$" + strconv.Itoa(i)
	}
}

// Append adds a value to qa and returns the placeholder value for the
// argument. e.g. $1, $2, etc.
func (qa *QueryArgs) Append(v interface{}) string {
	*qa = append(*qa, v)
	if len(*qa) < len(placeholders) {
		return placeholders[len(*qa)]
	}
	return "$" + strconv.Itoa(len(*qa))
}
