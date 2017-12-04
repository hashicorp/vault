package builtinplugins

import (
	"github.com/hashicorp/vault/plugins/database/cassandra"
	"github.com/hashicorp/vault/plugins/database/hana"
	"github.com/hashicorp/vault/plugins/database/mongodb"
	"github.com/hashicorp/vault/plugins/database/mssql"
	"github.com/hashicorp/vault/plugins/database/mysql"
	"github.com/hashicorp/vault/plugins/database/postgresql"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
)

// BuiltinFactory is the func signature that should be returned by
// the plugin's New() func.
type BuiltinFactory func() (interface{}, error)

var plugins = map[string]BuiltinFactory{
	// These four plugins all use the same mysql implementation but with
	// different username settings passed by the constructor.
	"mysql-database-plugin":        mysql.New(mysql.MetadataLen, mysql.MetadataLen, mysql.UsernameLen),
	"mysql-aurora-database-plugin": mysql.New(credsutil.NoneLength, mysql.LegacyMetadataLen, mysql.LegacyUsernameLen),
	"mysql-rds-database-plugin":    mysql.New(credsutil.NoneLength, mysql.LegacyMetadataLen, mysql.LegacyUsernameLen),
	"mysql-legacy-database-plugin": mysql.New(credsutil.NoneLength, mysql.LegacyMetadataLen, mysql.LegacyUsernameLen),

	"postgresql-database-plugin": postgresql.New,
	"mssql-database-plugin":      mssql.New,
	"cassandra-database-plugin":  cassandra.New,
	"mongodb-database-plugin":    mongodb.New,
	"hana-database-plugin":       hana.New,
}

// Get returns the BuiltinFactory func for a particular backend plugin
// from the plugins map.
func Get(name string) (BuiltinFactory, bool) {
	f, ok := plugins[name]
	return f, ok
}

// Keys returns the list of plugin names that are considered builtin plugins.
func Keys() []string {
	keys := make([]string, len(plugins))

	i := 0
	for k := range plugins {
		keys[i] = k
		i++
	}

	return keys
}
