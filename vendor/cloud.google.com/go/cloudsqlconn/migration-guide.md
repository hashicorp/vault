# Migrating from Cloud SQL Proxy v1 to the Go Connector

The Go Connector supports an improved version of the drivers available in v1.
Unlike V1, the Go Connectors drivers support:

1. Configuring a driver with all supported Go connector options
1. Configuring multiple drivers per engine type using distinct registered driver
   names
1. Support for SQL Server
1. (Postgres only) Configuring a connection using pgx directly (see README for
   details).

Below are examples of the Cloud SQL Proxy invocation vs the new Go connector
invocation.

## MySQL

### Cloud SQL Proxy

``` golang
import (
	"database/sql"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
)

func connectMySQL() *sql.DB {
	cfg := mysql.Cfg("project:region:instance", "user", "password")
	cfg.DBName = "DB_1"
	cfg.ParseTime = true

	db, err := mysql.DialCfg(cfg)
	if err != nil {
		// handle error as necessary
	}
	return db
}
```

### Cloud SQL Go Connector

``` golang
import (
	"database/sql"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/mysql/mysql"
)

func connectMySQL() *sql.DB {
	// Register a driver using whatever name you like.
	cleanup, err := mysql.RegisterDriver(
		"cloudsql-mysql",
		// any desired options go here, for example:
		cloudsqlconn.WithCredentialsFile("key.json"),
	)
	if err != nil {
		// handle error as necessary
	}
	// call cleanup to close the underylying driver when you're done with the
	// db.
	defer cleanup()

	db, err := sql.Open(
		"cloudsql-mysql", // matches the name registered above
		"myuser:mypass@cloudsql-mysql(project:region:instance)/mydb",
	)
	if err != nil {
		// handle error as necessary
	}
	return db
}
```

## Postgres

### Cloud SQL Proxy

``` golang
import (
	"database/sql"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)
func connectPostgres() *sql.DB {
	db, err := sql.Open(
		"cloudsqlpostgres",
		"host=project:region:instance user=postgres dbname=postgres password=password sslmode=disable",
	)
	if err != nil {
		// handle error as necessary
	}
	return db
}
```

### Cloud SQL Go Connector

``` golang
import (
	"database/sql"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
)

func connectPostgres() *sql.DB {
	// Register a driver using whatever name you like.
	cleanup, err := pgxv4.RegisterDriver(
		"cloudsql-postgres",
		// any desired options go here, for example:
		cloudsqlconn.WithCredentialsFile("key.json"),
		cloudsqlconn.WithIAMAuthN(),
	)
	if err != nil {
		// handle error as necessary
	}
	// call cleanup to close the underylying driver when you're done with the
	// db.
	defer cleanup()
	db, err := sql.Open(
		"cloudsql-postgres", // matches the name registered above
		"host=project:region:instance user=postgres password=password dbname=postgres sslmode=disable",
	)
	if err != nil {
		// handle error as necessary
	}
	return db
}
```

## SQL Server

### Cloud SQL Proxy

The Cloud SQL Proxy does not support SQL Server as a driver.

### Cloud SQL Go Connector

``` golang
import (
	"database/sql"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/sqlserver/mssql"
)

func connectSQLServer() *sql.DB {
	// Register a driver using whatever name you like.
	cleanup, err := mssql.RegisterDriver(
		"cloudsql-sqlserver",
		// any desired options go here, for example:
		cloudsqlconn.WithCredentialsFile("key.json"),
	)
	if err != nil {
		// handle error as necessary
	}
	// call cleanup when you're done with the database connection
	defer cleanup()

	db, err := sql.Open(
		"cloudsql-sqlserver", // matches the name registered above
		"sqlserver://user:password@localhost?database=mydb&cloudsql=project:region:instance",
	)
	if err != nil {
		// handle error as necessary
	}
	return db
}
```
