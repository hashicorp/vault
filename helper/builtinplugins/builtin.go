package builtinplugins

import "github.com/hashicorp/vault-plugins/database/mysql"

var BuiltinPlugins = map[string]func() error{
	"mysql-database-plugin": mysql.Run,
	//	"postgres-database-plugin": postgres.Run,
}
