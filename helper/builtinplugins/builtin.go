package builtinplugins

import (
	"github.com/hashicorp/vault-plugins/database/mysql"
	"github.com/hashicorp/vault-plugins/database/postgresql"
)

var BuiltinPlugins = map[string]func() error{
	"mysql-database-plugin":      mysql.Run,
	"postgresql-database-plugin": postgresql.Run,
}
