/*
Package gosnowflake is a pure Go Snowflake driver for the database/sql package.

Clients can use the database/sql package directly. For example:

	import (
		"database/sql"

		_ "github.com/snowflakedb/gosnowflake"
	)

	func main() {
		db, err := sql.Open("snowflake", "user:password@myaccount/mydb")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		...
	}

Connection String

Use the Open() function to create a database handle with connection parameters:

	db, err := sql.Open("snowflake", "<connection string>")

The Go Snowflake Driver supports the following connection syntaxes (or data source name (DSN) formats):

	* username[:password]@accountname/dbname/schemaname[?param1=value&...&paramN=valueN
	* username[:password]@accountname/dbname[?param1=value&...&paramN=valueN
	* username[:password]@hostname:port/dbname/schemaname?account=<your_account>[&param1=value&...&paramN=valueN]

where all parameters must be escaped or use `Config` and `DSN` to construct a DSN string.

The following example opens a database handle with the Snowflake account
myaccount where the username is jsmith, password is mypassword, database is
mydb, schema is testschema, and warehouse is mywh:

	db, err := sql.Open("snowflake", "jsmith:mypassword@myaccount/mydb/testschema?warehouse=mywh")

Connection Parameters

The connection string (DSN) can contain both connection parameters (described below) and session parameters
(https://docs.snowflake.com/en/sql-reference/parameters.html).

The following connection parameters are supported:

	* account <string>: Specifies the name of your Snowflake account, where string is the name
		assigned to your account by Snowflake. In the URL you received from
		Snowflake, your account name is the first segment in the domain (e.g.
		abc123 in https://abc123.snowflakecomputing.com). This parameter is
		optional if your account is specified after the @ character. If you
		are not on us-west-2 region or AWS deployment, then append the region
		after the account name, e.g. “<account>.<region>”. If you are not on
		AWS deployment, then append not only the region, but also the platform,
		e.g., “<account>.<region>.<platform>”. Account, region, and platform
		should be separated by a period (“.”), as shown above. If you are using
        a global url, then append connection group and "global",
        e.g., "account-<connection_group>.global". Account and connection group are
        separated by a dash ("-"), as shown above.

	* region <string>: DEPRECATED. You may specify a region, such as
		“eu-central-1”, with this parameter. However, since this parameter
		is deprecated, it is best to specify the region as part of the
		account parameter. For details, see the description of the account
		parameter.

	* database: Specifies the database to use by default in the client session
		(can be changed after login).

	* schema: Specifies the database schema to use by default in the client
		session (can be changed after login).

	* warehouse: Specifies the virtual warehouse to use by default for queries,
		loading, etc. in the client session (can be changed after login).

	* role: Specifies the role to use by default for accessing Snowflake
		objects in the client session (can be changed after login).

	* passcode: Specifies the passcode provided by Duo when using multi-factor authentication (MFA) for login.

	* passcodeInPassword: false by default. Set to true if the MFA passcode is
		embedded in the login password. Appends the MFA passcode to the end of the
		password.

	* loginTimeout: Specifies the timeout, in seconds, for login. The default
		is 60 seconds. The login request gives up after the timeout length if the
		HTTP response is success.

	* authenticator: Specifies the authenticator to use for authenticating user credentials:
		- To use the internal Snowflake authenticator, specify snowflake (Default).
		- To authenticate through Okta, specify https://<okta_account_name>.okta.com (URL prefix for Okta).
		- To authenticate using your IDP via a browser, specify externalbrowser.
		- To authenticate via OAuth, specify oauth and provide an OAuth Access Token (see the token parameter below).

	* application: Identifies your application to Snowflake Support.

	* insecureMode: false by default. Set to true to bypass the Online
		Certificate Status Protocol (OCSP) certificate revocation check.
		IMPORTANT: Change the default value for testing or emergency situations only.

	* token: a token that can be used to authenticate. Should be used in conjunction with the "oauth" authenticator.

	* client_session_keep_alive: Set to true have a heartbeat in the background every hour to keep the connection alive
		such that the connection session will never expire. Care should be taken in using this option as it opens up
		the access forever as long as the process is alive.

	* ocspFailOpen: true by default. Set to false to make OCSP check fail closed mode.

	* validateDefaultParameters: true by default. Set to false to disable checks on existence and privileges check for
								 Database, Schema, Warehouse and Role when setting up the connection

All other parameters are interpreted as session parameters (https://docs.snowflake.com/en/sql-reference/parameters.html).
For example, the TIMESTAMP_OUTPUT_FORMAT session parameter can be set by adding:

	...&TIMESTAMP_OUTPUT_FORMAT=MM-DD-YYYY...

A complete connection string looks similar to the following:

	my_user_name:my_password@ac123456/my_database/my_schema?my_warehouse=inventory_warehouse&role=my_user_role&DATE_OUTPUT_FORMAT=YYYY-MM-DD
                                                                ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ ^^^^^^^^^^^^^^^^^ ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
                                                                      connection                     connection           session
                                                                      parameter                      parameter            parameter

Session-level parameters can also be set by using the SQL command "ALTER SESSION"
(https://docs.snowflake.com/en/sql-reference/sql/alter-session.html).

Alternatively, use OpenWithConfig() function to create a database handle with the specified Config.

Proxy

The Go Snowflake Driver honors the environment variables HTTP_PROXY, HTTPS_PROXY and NO_PROXY for the forward proxy setting.

NO_PROXY specifies which hostname endings should be allowed to bypass the proxy server, e.g. :code:`no_proxy=.amazonaws.com` means that AWS S3 access does not need to go through the proxy.

NO_PROXY does not support wildcards. Each value specified should be one of the following:

    * The end of a hostname (or a complete hostname), for example: ".amazonaws.com" or "xy12345.snowflakecomputing.com".

    * An IP address, for example "192.196.1.15".

If more than one value is specified, values should be separated by commas, for example:

    no_proxy=localhost,.my_company.com,xy12345.snowflakecomputing.com,192.168.1.15,192.168.1.16


Logging

By default, the driver's builtin logger is exposing logrus's FieldLogger and default at INFO level.
Users can use SetLogger in driver.go to set a customized logger for gosnowflake package.

In order to enable debug logging for the driver, user could use SetLogLevel("debug") in SFLogger interface
as shown in demo code at cmd/logger.go. To redirect the logs SFlogger.SetOutput method could do the work.


Query request ID

Specific query request ID can be set in the context and will be passed through
in place of the default randomized request ID. For example:

	ctxWithID := WithRequestID(ctx, "app-request-id-...")
	rows, err := db.QueryContext(ctxWithID, query)

Canceling Query by CtrlC

From 0.5.0, a signal handling responsibility has moved to the applications. If you want to cancel a
query/command by Ctrl+C, add a os.Interrupt trap in context to execute methods that can take the context parameter,
e.g., QueryContext, ExecContext.

	// handle interrupt signal
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()
	... (connection)
	// execute a query
	rows, err := db.QueryContext(ctx, query)
	... (Ctrl+C to cancel the query)

See cmd/selectmany.go for the full example.

Supported Data Types

Queries return SQL column type information in the ColumnType type. The
DatabaseTypeName method returns strings representing Snowflake data types.
The following table shows those strings, the corresponding Snowflake data
type, and the corresponding Golang native data type. The columns are:

   1. The string representation of the data type.
   2. The SQL data type.
   3. The default Golang data type that is returned when you use snowflakeRows.Scan() to read data from
      JSON data format via an interface{}. (All returned values are JSON strings.)
   4. The standard Golang data type that is returned when you use snowflakeRows.Scan() to read data from
      JSON data format directly.
   5. Footnotes numbers.

This table shows the data types:

  =============================================================================================
                 |                                    | Default Go  | Supported  |
  String         |                                    | Data Type   | Go Data    |
  Representation | Snowflake Data Type                | for Scan()  | Types for  | Footnotes
                 |                                    | interface{} | Scan()     |
                 |                                    | (JSON)      | (JSON)     |
  ==============================================================================================
  BOOLEAN        | BOOLEAN                            | string      | bool       |
  TEXT           | VARCHAR/STRING                     | string      | string     |
  REAL           | REAL/DOUBLE                        | string      | float64    | [1]  [2]
  FIXED          | INTEGER that fits in int64         | string      | int64      | [1]  [2]
  FIXED          | NUMBER(P, S) where S > 0           | string      |            | [1]  [3]
  DATE           | DATE                               | string      | time.Time  |
  TIME           | TIME                               | string      | time.Time  |
  TIMESTAMP_LTZ  | TIMESTAMP_LTZ                      | string      | time.Time  |
  TIMESTAMP_NTZ  | TIMESTAMP_NTZ                      | string      | time.Time  |
  TIMESTAMP_TZ   | TIMESTAMP_TZ                       | string      | time.Time  |
  BINARY         | BINARY                             | string      | []byte     |
  ARRAY          | ARRAY                              | string      | string     |
  OBJECT         | OBJECT                             | string      | string     |
  VARIANT        | VARIANT                            | string      | string     |

Footnotes:

  [1] Converting from a higher precision data type to a lower precision data type via the snowflakeRows.Scan() method can lose low bits (lose precision), lose high bits (completely change the value), or result in error.

  [2] Attempting to convert from a higher precision data type to a lower precision data type via interface{} causes an error.

  [3] If the value in Snowflake is too large to fit into the corresponding Golang data type, then conversion can return either an int64 with the high bits truncated or an error.

Note: SQL NULL values are converted to Golang nil values, and vice-versa.

Binding Parameters to Array Variables For Batch Inserts

Version 1.3.9 (and later) of the Go Snowflake Driver supports the ability to bind an array variable to a parameter in an SQL
INSERT statement. You can use this technique to insert multiple rows in a single batch.

As an example, the following code inserts rows into a table that contains integer, float, boolean, and string columns. The example
binds arrays to the parameters in the INSERT statement.

	// Create a table containing an integer, float, boolean, and string column.
	_, err = db.Exec("create or replace table my_table(c1 int, c2 float, c3 boolean, c4 string)")
	...
	// Define the arrays containing the data to insert.
	intArray := []int{1, 2, 3}
	fltArray := []float64{0.1, 2.34, 5.678}
	boolArray := []bool{true, false, true}
	strArray := []string{"test1", "test2", "test3"}
	...
	// Insert the data from the arrays into the table.
	_, err = db.Exec("insert into my_table values (?, ?, ?, ?)", intArray, fltArray, boolArray, strArray)

Note: For alternative ways to load data into the Snowflake database (including bulk loading using the COPY command), see
Loading Data Into Snowflake (https://docs.snowflake.com/en/user-guide-data-load.html).

Binding a Parameter to a Time Type

Go's database/sql package supports the ability to bind a parameter in an SQL statement to a time.Time variable.
However, when the client binds data to send to the server, the driver cannot determine the correct Snowflake date/timestamp data
type to associate with the binding parameter. For example:

	dbt.mustExec("CREATE OR REPLACE TABLE tztest (id int, ntz, timestamp_ntz, ltz timestamp_ltz)")
	// ...
	stmt, err :=dbt.db.Prepare("INSERT INTO tztest(id,ntz,ltz) VALUES(1, ?, ?)")
	// ...
	tmValue time.Now()
	// ... Is tmValue a TIMESTAMP_NTZ or TIMESTAMP_LTZ?
	_, err = stmt.Exec(tmValue, tmValue)

To resolve this issue, a binding parameter flag is introduced that associates
any subsequent time.Time type to the DATE, TIME, TIMESTAMP_LTZ, TIMESTAMP_NTZ
or BINARY data type. The above example could be rewritten as follows:

	import (
		sf "github.com/snowflakedb/gosnowflake"
	)
	dbt.mustExec("CREATE OR REPLACE TABLE tztest (id int, ntz, timestamp_ntz, ltz timestamp_ltz)")
	// ...
	stmt, err :=dbt.db.Prepare("INSERT INTO tztest(id,ntz,ltz) VALUES(1, ?, ?)")
	// ...
	tmValue time.Now()
	// ...
	_, err = stmt.Exec(sf.DataTypeTimestampNtz, tmValue, sf.DataTypeTimestampLtz, tmValue)

Timestamps with Time Zones

The driver fetches TIMESTAMP_TZ (timestamp with time zone) data using the
offset-based Location types, which represent a collection of time offsets in
use in a geographical area, such as CET (Central European Time) or UTC
(Coordinated Universal Time). The offset-based Location data is generated and
cached when a Go Snowflake Driver application starts, and if the given offset
is not in the cache, it is generated dynamically.

Currently, Snowflake doesn't support the name-based Location types, e.g.,
"America/Los_Angeles".

For more information about Location types, see the Go documentation for https://golang.org/pkg/time/#Location.

Binary Data

Internally, this feature leverages the []byte data type. As a result, BINARY
data cannot be bound without the binding parameter flag. In the following
example, sf is an alias for the gosnowflake package:

	var b = []byte{0x01, 0x02, 0x03}
	_, err = stmt.Exec(sf.DataTypeBinary, b)

Maximum number of Result Set Chunk Downloader

The driver directly downloads a result set from the cloud storage if the size is large. It is
required to shift workloads from the Snowflake database to the clients for scale. The download takes place by goroutine
named "Chunk Downloader" asynchronously so that the driver can fetch the next result set while the application can
consume the current result set.

The application may change the number of result set chunk downloader if required. Note this doesn't help reduce
memory footprint by itself. Consider Custom JSON Decoder.

	import (
		sf "github.com/snowflakedb/gosnowflake"
	)
	sf.MaxChunkDownloadWorkers = 2


Experimental: Custom JSON Decoder for parsing Result Set

The application may have the driver use a custom JSON decoder that incrementally parses the result set as follows.

	import (
		sf "github.com/snowflakedb/gosnowflake"
	)
	sf.CustomJSONDecoderEnabled = true
	...

This option will reduce the memory footprint to half or even quarter, but it can significantly degrade the
performance depending on the environment. The test cases running on Travis Ubuntu box show five times less memory
footprint while four times slower. Be cautious when using the option.

JWT authentication

The Go Snowflake Driver supports JWT (JSON Web Token) authentication.

To enable this feature, construct the DSN with fields "authenticator=SNOWFLAKE_JWT&privateKey=<your_private_key>",
or using a Config structure specifying:

	config := &Config{
		...
		Authenticator: "SNOWFLAKE_JWT"
		PrivateKey:   "<your_private_key_struct in *rsa.PrivateKey type>"
	}

The <your_private_key> should be a base64 URL encoded PKCS8 rsa private key string. One way to encode a byte slice to URL
base 64 URL format is through the base64.URLEncoding.EncodeToString() function.

On the server side, you can alter the public key with the SQL command:

	ALTER USER <your_user_name> SET RSA_PUBLIC_KEY='<your_public_key>';

The <your_public_key> should be a base64 Standard encoded PKI public key string. One way to encode a byte slice to base
64 Standard format is through the base64.StdEncoding.EncodeToString() function.

To generate the valid key pair, you can execute the following commands in the shell:

	# generate 2048-bit pkcs8 encoded RSA private key
	openssl genpkey -algorithm RSA \
    	-pkeyopt rsa_keygen_bits:2048 \
    	-pkeyopt rsa_keygen_pubexp:65537 | \
  		openssl pkcs8 -topk8 -outform der > rsa-2048-private-key.p8

	# extract 2048-bit PKI encoded RSA public key from the private key
	openssl pkey -pubout -inform der -outform der \
    	-in rsa-2048-private-key.p8 \
    	-out rsa-2048-public-key.spki

Note: As of February 2020, Golang's official library does not support passcode-encrypted PKCS8 private key.
For security purposes, Snowflake highly recommends that you store the passcode-encrypted private key on the disk and
decrypt the key in your application using a library you trust.


Executing Multiple Statements in One Call

This feature is available in version 1.3.8 or later of the driver.

By default, Snowflake returns an error for queries issued with multiple statements.
This restriction helps protect against SQL Injection attacks (https://en.wikipedia.org/wiki/SQL_injection).

The multi-statement feature allows users skip this restriction and execute multiple SQL statements through a
single Golang function call. However, this opens up the possibility for SQL injection, so it should be used carefully.
The risk can be reduced by specifying the exact number of statements to be executed, which makes it more difficult to
inject a statement by appending it. More details are below.

The Go Snowflake Driver provides two functions that can execute multiple SQL statements in a single call:

	* db.QueryContext(): This function is used to execute queries, such as SELECT statements, that return a result set.
	* db.ExecContext(): This function is used to execute statements that don't return a result set (i.e. most DML and DDL statements).

To compose a multi-statement query, simply create a string that contains all the queries, separated by semicolons,
in the order in which the statements should be executed.


To protect against SQL Injection attacks while using the multi-statement feature, pass a Context that specifies
the number of statements in the string. For example:


	import (
		"context"
		"database/sql"
	)

	var multi_statement_query = "SELECT c1 FROM t1; SELECT c2 FROM t2"
	var number_of_statements = 2
	blank_context = context.Background()
	multi_statement_context, _ := WithMultiStatement(blank_context, number_of_statements)
	rows, err := db.QueryContext(multi_statement_context, multi_statement_query)


When multiple queries are executed by a single call to QueryContext(), multiple result sets are returned. After
you process the first result set, get the next result set (for the next SQL statement) by calling NextResultSet().

The following pseudo-code shows how to process multiple result sets:

	Execute the statement and get the result set(s):


		rows, err := db.QueryContext(ctx, multiStmtQuery)

	Retrieve the rows in the first query's result set:

		while rows.Next() {
			err = rows.Scan(&variable_1)
			if err != nil {
				t.Errorf("failed to scan: %#v", err)
			}
			...
		}

	Retrieve the remaining result sets and the rows in them:

		while rows.NextResultSet()  {

			while rows.Next() {
				...
			}

		}

The function db.execContext() returns a single result, which is the sum of the results of the individual statements.
For example, if your multi-statement query executed two UPDATE statements, each of which updated 10 rows,
then the result returned would be 20. Individual results for individual statements are not available.

The following code shows how to retrieve the result of a multi-statement query executed through db.ExecContext():

    Execute the SQL statements:

        res, err := db.ExecContext(ctx, multiStmtQuery)

    Get the summed result and store it in the variable named count:

        count, err := res.RowsAffected()


Note: Because a multi-statement ExecContext() returns a single value, you cannot detect offsetting errors.
For example, suppose you expected the return value to be 20 because you expected each UPDATE statement to
update 10 rows. If one UPDATE statement updated 15 rows and the other UPDATE statement updated only 5
rows, the total would still be 20. You would see no indication that the UPDATES had not functioned as
expected.


The ExecContext() function does not return an error if passed a query (e.g. a SELECT statement). However, it
still returns only a single value, not a result set, so using it to execute queries (or a mix of queries and non-query
statements) is impractical.

The QueryContext() function does not return an error if passed non-query statements (e.g. DML). The function
returns a result set for each statement, whether or not the statement is a query. For each non-query statement, the
result set contains a single row that contains a single column; the value is the number of rows changed by the
statement.

If you want to execute a mix of query and non-query statements (e.g. a mix of SELECT and DML statements) in a
multi-statement query, use QueryContext(). You can retrieve the result sets for the queries,
and you can retrieve or ignore the row counts for the non-query statements.


If any of the SQL statements fail to compile or execute, execution is aborted. Any previous statements that ran before are unaffected.

For example, if the statements below are run as one multi-statement query, the multi-statement query fails on the
third statement, and an exception is thrown.


	CREATE OR REPLACE TABLE test(n int);
	INSERT INTO TEST VALUES (1), (2);
	INSERT INTO TEST VALUES ('not_an_integer');  -- execution fails here
	INSERT INTO TEST VALUES (3);

If you then query the contents of the table named "test", the values 1 and 2 would be present.

When using the QueryContext() and ExecContext() functions, golang code can check for errors the usual way. For
example:

	rows, err := db.QueryContext(ctx, multiStmtQuery)
	if err != nil {
		Fatalf("failed to query multiple statements: %v", err)
	}

Preparing statements and using bind variables are also not supported for multi-statement queries.


Limitations

GET and PUT operations are unsupported.
*/
package gosnowflake
