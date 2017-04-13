package builtinplugins

import (
	"github.com/hashicorp/vault/plugins/database/mysql"
	"github.com/hashicorp/vault/plugins/database/postgresql"
)

var BuiltinPlugins *builtinPlugins = &builtinPlugins{
	plugins: map[string]func() error{
		"mysql-database-plugin":      mysql.Run,
		"postgresql-database-plugin": postgresql.Run,
	},
}

// The list of builtin plugins should not be changed by any other package, so we
// store them in an unexported variable in this unexported struct.
type builtinPlugins struct {
	plugins map[string]func() error
}

func (b *builtinPlugins) Get(name string) (func() error, bool) {
	f, ok := b.plugins[name]
	return f, ok
}

func (b *builtinPlugins) Keys() []string {
	keys := make([]string, len(b.plugins))

	i := 0
	for k := range b.plugins {
		keys[i] = k
		i++
	}

	return keys
}
