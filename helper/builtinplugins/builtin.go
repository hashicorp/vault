package builtinplugins

import (
	"github.com/hashicorp/vault/plugins/database/mssql"
	"github.com/hashicorp/vault/plugins/database/mysql"
	"github.com/hashicorp/vault/plugins/database/postgresql"
)

type BuiltinFactory func() (interface{}, error)

var plugins map[string]BuiltinFactory = map[string]BuiltinFactory{
	"mysql-database-plugin":      mysql.New,
	"postgresql-database-plugin": postgresql.New,
	"mssql-database-plugin":      mssql.New,
}

func Get(name string) (BuiltinFactory, bool) {
	f, ok := plugins[name]
	return f, ok
}

func Keys() []string {
	keys := make([]string, len(plugins))

	i := 0
	for k := range plugins {
		keys[i] = k
		i++
	}

	return keys
}
