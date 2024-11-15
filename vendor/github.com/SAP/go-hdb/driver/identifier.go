package driver

import (
	"strconv"

	"github.com/SAP/go-hdb/driver/internal/rand/alphanum"
)

// Identifier in hdb SQL statements like schema or table name.
type Identifier string

// RandomIdentifier returns a random Identifier prefixed by the prefix parameter.
// This function is used to generate database objects with random names for test and example code.
func RandomIdentifier(prefix string) Identifier {
	return Identifier(prefix + alphanum.ReadString(16))
}

func (i Identifier) isSimple() bool {
	// var reSimple = regexp.MustCompile("^[_A-Z][_#$A-Z0-9]*$")
	for i, r := range i {
		switch {
		case r == '_' || ('A' <= r && r <= 'Z'): // valid char
		case i != 0 && (r == '#' || r == '$' || ('0' <= r && r <= '9')): // valid char for non first char
		default:
			return false
		}
	}
	return true
}
func (i Identifier) String() string {
	if i.isSimple() {
		return string(i)
	}
	return strconv.Quote(string(i))
}
