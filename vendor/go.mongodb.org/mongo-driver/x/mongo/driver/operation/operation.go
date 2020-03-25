package operation

//go:generate operationgen insert.toml operation insert.go
//go:generate operationgen find.toml operation find.go
//go:generate operationgen list_collections.toml operation list_collections.go
//go:generate operationgen createIndexes.toml operation createIndexes.go
//go:generate operationgen drop_collection.toml operation drop_collection.go
//go:generate operationgen distinct.toml operation distinct.go
//go:generate operationgen delete.toml operation delete.go
//go:generate operationgen drop_indexes.toml operation drop_indexes.go
//go:generate operationgen drop_database.toml operation drop_database.go
//go:generate operationgen commit_transaction.toml operation commit_transaction.go
//go:generate operationgen abort_transaction.toml operation abort_transaction.go
//go:generate operationgen count.toml operation count.go
//go:generate operationgen end_sessions.toml operation end_sessions.go
