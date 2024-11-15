// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cloudsqlconn provides functions for authorizing and encrypting
// connections. These functions can be used with a database driver to
// connect to a Cloud SQL instance.
//
// The instance connection name for a Cloud SQL instance is always in the
// format "project:region:instance".
//
// # Creating a Dialer
//
// To start working with this package, create a Dialer. There are two ways of
// creating a Dialer, which one you use depends on your database driver.
//
// # Postgres
//
// Postgres users have the option of using the [database/sql] interface or using [pgx] directly.
//
// To use a dialer with [pgx], we recommend using connection pooling with
// [pgxpool]. To create the dialer use the NewDialer func.
//
//	import (
//	    "context"
//	    "net"
//
//	    "cloud.google.com/go/cloudsqlconn"
//	    "github.com/jackc/pgx/v4/pgxpool"
//	)
//
//	func connect() {
//	    // Configure the driver to connect to the database
//	    dsn := "user=myuser password=mypass dbname=mydb sslmode=disable"
//	    config, err := pgxpool.ParseConfig(dsn)
//	    if err != nil {
//		    // handle error
//	    }
//
//	    // Create a new dialer with any options
//	    d, err := cloudsqlconn.NewDialer(context.Background())
//	    if err != nil {
//		    // handle error
//	    }
//
//	    // Tell the driver to use the Cloud SQL Go Connector to create connections
//	    config.ConnConfig.DialFunc = func(ctx context.Context, _ string, instance string) (net.Conn, error) {
//		    return d.Dial(ctx, "project:region:instance")
//	    }
//
//	    // Interact with the driver directly as you normally would
//	    conn, err := pgxpool.ConnectConfig(context.Background(), config)
//	    if err != nil {
//		    // handle error
//	    }
//
//	    // call cleanup when you're done with the database connection
//	    cleanup := func() error { return d.Close() }
//	    // ... etc
//	}
//
// To use [database/sql], call pgxv4.RegisterDriver with any necessary Dialer
// configuration.
//
// Note: the connection string must use the keyword/value format
// with host set to the instance connection name. The returned cleanup func
// will stop the dialer's background refresh goroutine and so should only be
// called when you're done with the Dialer.
//
//	import (
//	    "database/sql"
//
//	    "cloud.google.com/go/cloudsqlconn"
//	    "cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
//	)
//
//	func connect() {
//	    // adjust options as needed
//	    cleanup, err := pgxv4.RegisterDriver("cloudsql-postgres", cloudsqlconn.WithIAMAuthN())
//	    if err != nil {
//	    	// ... handle error
//	    }
//	    // call cleanup when you're done with the database connection
//	    defer cleanup()
//
//	    db, err := sql.Open(
//	        "cloudsql-postgres",
//	        "host=project:region:instance user=myuser password=mypass dbname=mydb sslmode=disable",
//	    )
//	    // ... etc
//	}
//
// # MySQL
//
// MySQL users should use [database/sql]. Use mysql.RegisterDriver with any
// necessary Dialer configuration.
//
// Note: The returned cleanup func will stop the dialer's background refresh
// goroutine and should only be called when you're done with the Dialer.
//
//	import (
//	    "database/sql"
//
//	    "cloud.google.com/go/cloudsqlconn"
//	    "cloud.google.com/go/cloudsqlconn/mysql/mysql"
//	)
//
//	func connect() {
//	    // adjust options as needed
//	    cleanup, err := mysql.RegisterDriver("cloudsql-mysql", cloudsqlconn.WithIAMAuthN())
//	    if err != nil {
//	        // ... handle error
//	    }
//	    // call cleanup when you're done with the database connection
//	    defer cleanup()
//
//	    db, err := sql.Open(
//	        "cloudsql-mysql",
//	        "myuser:mypass@cloudsql-mysql(project:region:instance)/mydb",
//	    )
//	    // ... etc
//	}
//
// # SQL Server
//
// SQL Server users should use [database/sql]. Use mssql.RegisterDriver with any
// necessary Dialer configuration.
//
// Note: The returned cleanup func will stop the dialer's background refresh
// goroutine and should only be called when you're done with the Dialer.
//
//	import (
//	    "database/sql"
//
//	    "cloud.google.com/go/cloudsqlconn"
//	    "cloud.google.com/go/cloudsqlconn/sqlserver/mssql"
//	)
//
//	func connect() {
//	    cleanup, err := mssql.RegisterDriver("cloudsql-sqlserver")
//	    if err != nil {
//	        // ... handle error
//	    }
//	    // call cleanup when you're done with the database connection
//	    defer cleanup()
//
//	    db, err := sql.Open(
//	        "cloudsql-sqlserver",
//	        "sqlserver://user:password@localhost?database=mydb&cloudsql=project:region:instance",
//	    )
//	    // ... etc
//	}
//
// [database/sql]: https://pkg.go.dev/database/sql
// [pgx]: https://github.com/jackc/pgx
// [pgxpool]: https://pkg.go.dev/github.com/jackc/pgx/v4/pgxpool
package cloudsqlconn
