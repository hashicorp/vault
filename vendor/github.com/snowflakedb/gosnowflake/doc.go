/*
Package gosnowflake is a pure Go Snowflake driver for the database/sql package.

Clients can use the database/sql package directly. For example:

	import (
		"database/sql"

		_ "github.com/snowflakedb/gosnowflake"
	)

	func main() {
		db, err := sql.Open("snowflake", "user:password@my_organization-my_account/mydb")
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

	* username[:password]@<account_identifier>/dbname/schemaname[?param1=value&...&paramN=valueN]
	* username[:password]@<account_identifier>/dbname[?param1=value&...&paramN=valueN]
	* username[:password]@hostname:port/dbname/schemaname?account=<account_identifier>[&param1=value&...&paramN=valueN]

where all parameters must be escaped or use Config and DSN to construct a DSN string.

For information about account identifiers, see the Snowflake documentation
(https://docs.snowflake.com/en/user-guide/admin-account-identifier.html).

The following example opens a database handle with the Snowflake account
named "my_account" under the organization named "my_organization",
where the username is "jsmith", password is "mypassword", database is "mydb",
schema is "testschema", and warehouse is "mywh":

	db, err := sql.Open("snowflake", "jsmith:mypassword@my_organization-my_account/mydb/testschema?warehouse=mywh")


Connection Parameters

The connection string (DSN) can contain both connection parameters (described below) and session parameters
(https://docs.snowflake.com/en/sql-reference/parameters.html).

The following connection parameters are supported:

	* account <string>: Specifies your Snowflake account, where "<string>" is the account
		identifier assigned to your account by Snowflake.
		For information about account identifiers, see the Snowflake documentation
		(https://docs.snowflake.com/en/user-guide/admin-account-identifier.html).

		If you are using a global URL, then append the connection group and ".global"
		(e.g. "<account_identifier>-<connection_group>.global"). The account identifier and the
		connection group are separated by a dash ("-"), as shown above.

		This parameter is optional if your account identifier is specified after the "@" character
		in the connection string.

	* region <string>: DEPRECATED. You may specify a region, such as
		"eu-central-1", with this parameter. However, since this parameter
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

	* passcodeInPassword: false by default. Set to true if the MFA passcode is embedded
		in the login password. Appends the MFA passcode to the end of the password.

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

NO_PROXY specifies which hostname endings should be allowed to bypass the proxy server, e.g. no_proxy=.amazonaws.com means that Amazon S3 access does not need to go through the proxy.

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

A specific query request ID can be set in the context and will be passed through
in place of the default randomized request ID. For example:

	requestID := uuid.MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	ctxWithID := WithRequestID(ctx, requestID)
	rows, err := db.QueryContext(ctxWithID, query)

Canceling Query by CtrlC

From 0.5.0, a signal handling responsibility has moved to the applications. If you want to cancel a
query/command by Ctrl+C, add a os.Interrupt trap in context to execute methods that can take the context parameter
(e.g. QueryContext, ExecContext).

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

The Go Snowflake Driver now supports the Arrow data format for data transfers
between Snowflake and the Golang client. The Arrow data format avoids extra
conversions between binary and textual representations of the data. The Arrow
data format can improve performance and reduce memory consumption in clients.

Snowflake continues to support the JSON data format.

The data format is controlled by the session-level parameter
GO_QUERY_RESULT_FORMAT. To use JSON format, execute:

	ALTER SESSION SET GO_QUERY_RESULT_FORMAT = 'JSON';

The valid values for the parameter are:

	* ARROW (default)
	* JSON

If the user attempts to set the parameter to an invalid value, an error is
returned.

The parameter name and the parameter value are case-insensitive.

This parameter can be set only at the session level.

Usage notes:

	* The Arrow data format reduces rounding errors in floating point numbers. You might see slightly
	  different values for floating point numbers when using Arrow format than when using JSON format.
	  In order to take advantage of the increased precision, you must pass in the context.Context object
	  provided by the WithHigherPrecision function when querying.

	* Traditionally, the rows.Scan() method returned a string when a variable of types interface was passed
	  in. Turning on the flag ENABLE_HIGHER_PRECISION via WithHigherPrecision will return the natural,
	  expected data type as well.

	* For some numeric data types, the driver can retrieve larger values when using the Arrow format than
	  when using the JSON format. For example, using Arrow format allows the full range of SQL NUMERIC(38,0)
	  values to be retrieved, while using JSON format allows only values in the range supported by the
	  Golang int64 data type.

	  Users should ensure that Golang variables are declared using the appropriate data type for the full
	  range of values contained in the column. For an example, see below.

When using the Arrow format, the driver supports more Golang data types and
more ways to convert SQL values to those Golang data types. The table below
lists the supported Snowflake SQL data types and the corresponding Golang
data types. The columns are:

   1. The SQL data type.
   2. The default Golang data type that is returned when you use snowflakeRows.Scan() to read data from
      Arrow data format via an interface{}.
   3. The possible Golang data types that can be returned when you use snowflakeRows.Scan() to read data
      from Arrow data format directly.
   4. The default Golang data type that is returned when you use snowflakeRows.Scan() to read data from
      JSON data format via an interface{}. (All returned values are strings.)
   5. The standard Golang data type that is returned when you use snowflakeRows.Scan() to read data from
      JSON data format directly.

  Go Data Types for Scan()
  ===================================================================================================================
                         |                    ARROW                    |                    JSON
  ===================================================================================================================
  SQL Data Type          | Default Go Data Type   | Supported Go Data  | Default Go Data Type   | Supported Go Data
                         | for Scan() interface{} | Types for Scan()   | for Scan() interface{} | Types for Scan()
  ===================================================================================================================
    BOOLEAN              | bool                                        | string                 | bool
  -------------------------------------------------------------------------------------------------------------------
    VARCHAR              | string                                      | string
  -------------------------------------------------------------------------------------------------------------------
    DOUBLE               | float32, float64                  [1] , [2] | string                 | float32, float64
  -------------------------------------------------------------------------------------------------------------------
    INTEGER that         | int, int8, int16, int32, int64              | string                 | int, int8, int16,
    fits in int64        |                                   [1] , [2] |                        | int32, int64
  -------------------------------------------------------------------------------------------------------------------
    INTEGER that doesn't | int, int8, int16, int32, int64,  *big.Int   | string                 | error
    fit in int64         |                       [1] , [2] , [3] , [4] |
  -------------------------------------------------------------------------------------------------------------------
    NUMBER(P, S)         | float32, float64,  *big.Float               | string                 | float32, float64
    where S > 0          |                       [1] , [2] , [3] , [5] |
  -------------------------------------------------------------------------------------------------------------------
    DATE                 | time.Time                                   | string                 | time.Time
  -------------------------------------------------------------------------------------------------------------------
    TIME                 | time.Time                                   | string                 | time.Time
  -------------------------------------------------------------------------------------------------------------------
    TIMESTAMP_LTZ        | time.Time                                   | string                 | time.Time
  -------------------------------------------------------------------------------------------------------------------
    TIMESTAMP_NTZ        | time.Time                                   | string                 | time.Time
  -------------------------------------------------------------------------------------------------------------------
    TIMESTAMP_TZ         | time.Time                                   | string                 | time.Time
  -------------------------------------------------------------------------------------------------------------------
    BINARY               | []byte                                      | string                 | []byte
  -------------------------------------------------------------------------------------------------------------------
    ARRAY                | string                                      | string
  -------------------------------------------------------------------------------------------------------------------
    OBJECT               | string                                      | string
  -------------------------------------------------------------------------------------------------------------------
    VARIANT              | string                                      | string

  [1] Converting from a higher precision data type to a lower precision data type via the snowflakeRows.Scan()
      method can lose low bits (lose precision), lose high bits (completely change the value), or result in error.

  [2] Attempting to convert from a higher precision data type to a lower precision data type via interface{}
      causes an error.

  [3] Higher precision data types like *big.Int and *big.Float can be accessed by querying with a context
      returned by WithHigherPrecision().

  [4] You cannot directly Scan() into the alternative data types via snowflakeRows.Scan(), but can convert to
      those data types by using .Int64()/.String()/.Uint64() methods. For an example, see below.

  [5] You cannot directly Scan() into the alternative data types via snowflakeRows.Scan(), but can convert to
      those data types by using .Float32()/.String()/.Float64() methods. For an example, see below.

Note: SQL NULL values are converted to Golang nil values, and vice-versa.

The following example shows how to retrieve very large values using the math/big
package. This example retrieves a large INTEGER value to an interface and then
extracts a big.Int value from that interface. If the value fits into an int64,
then the code also copies the value to a variable of type int64. Note that a
context that enables higher precision must be passed in with the query.

       import "context"
       import "math/big"

       ...

       var my_interface interface{}
       var my_big_int_pointer *big.Int
       var my_int64 int64
       var rows snowflakeRows

       ...
       rows = db.QueryContext(WithHigherPrecision(context.Background), <query>)
       rows.Scan(&my_interface)
       my_big_int_pointer, ok = my_interface.(*big.Int)
       if my_big_int_pointer.IsInt64() {
           my_int64 = my_big_int_pointer.Int64()
       }

If the variable named "rows" is known to contain a big.Int, then you can use the following instead of scanning into an interface
and then converting to a big.Int:

       rows.Scan(&my_big_int_pointer)

If the variable named "rows" contains a big.Int, then each of the following fails:

       rows.Scan(&my_int64)

       my_int64, _ = my_interface.(int64)

Similar code and rules also apply to big.Float values.

If you are not sure what data type will be returned, you can use code similar to the following to check the data type
of the returned value:

    // Create variables into which you can scan the returned values.
    var i64 int64
    var bigIntPtr *big.Int

    for rows.Next() {
        // Get the data type info.
        column_types, err := rows.ColumnTypes()
        if err != nil {
            log.Fatalf("ERROR: ColumnTypes() failed. err: %v", err)
        }
        // The data type of the zeroeth column in the row.
        column_type := column_types[0].ScanType()
        // Choose the appropriate variable based on the data type.
        switch column_type {
            case reflect.TypeOf(i64):
                err = rows.Scan(&i64)
                fmt.Println("INFO: retrieved int64 value:")
                fmt.Println(i64)
            case reflect.TypeOf(bigIntPtr):
                err = rows.Scan(&bigIntPtr)
                fmt.Println("INFO: retrieved bigIntPtr value:")
                fmt.Println(bigIntPtr)
        }
    }


Binding Parameters

Binding allows a SQL statement to use a value that is stored in a Golang variable.

Without binding, a SQL statement specifies values by specifying literals inside the statement.
For example, the following statement uses the literal value ``42`` in an UPDATE statement:

	_, err = db.Exec("UPDATE table1 SET integer_column = 42 WHERE ID = 1000")

With binding, you can execute a SQL statement that uses a value that is inside a variable. For example:

	var my_integer_variable int = 42
	_, err = db.Exec("UPDATE table1 SET integer_column = ? WHERE ID = 1000", my_integer_variable)

The ``?`` inside the ``VALUES`` clause specifies that the SQL statement uses the value from a variable.

Binding data that involves time zones can require special handling. For details, see the section
titled "Timestamps with Time Zones".


Binding Parameters to Array Variables

Version 1.3.9 (and later) of the Go Snowflake Driver supports the ability to bind an array variable to a parameter in a SQL
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
	// Insert the data from the arrays and wrap in an Array() function into the table.
	_, err = db.Exec("insert into my_table values (?, ?, ?, ?)", Array(&intArray), Array(&fltArray), Array(&boolArray), Array(&strArray))

Note: For alternative ways to load data into the Snowflake database (including bulk loading using the COPY command), see
Loading Data into Snowflake (https://docs.snowflake.com/en/user-guide-data-load.html).


Batch Inserts and Binding Parameters

When you use array binding to insert a large number of values, the driver can
improve performance by streaming the data (without creating files on the local
machine) to a temporary stage for ingestion. The driver automatically does this
when the number of values exceeds a threshold (no changes are needed to user code).

In order for the driver to send the data to a temporary stage, the user must have the following privilege on the schema:

    CREATE STAGE

If the user does not have this privilege, the driver falls back to sending the data with the query to the Snowflake database.

In addition, the current database and schema for the session must be set. If these are not set,
the CREATE TEMPORARY STAGE command executed by the driver can fail with the following error:

	CREATE TEMPORARY STAGE SYSTEM$BIND file_format=(type=csv field_optionally_enclosed_by='"')
	Cannot perform CREATE STAGE. This session does not have a current schema. Call 'USE SCHEMA', or use a qualified name.

For alternative ways to load data into the Snowflake database (including bulk loading using the COPY command),
see Loading Data into Snowflake (https://docs.snowflake.com/en/user-guide-data-load.html).


Binding a Parameter to a Time Type

Go's database/sql package supports the ability to bind a parameter in a SQL statement to a time.Time variable.
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

Currently, Snowflake does not support the name-based Location types (e.g. "America/Los_Angeles").

For more information about Location types, see the Go documentation for https://golang.org/pkg/time/#Location.

Binary Data

Internally, this feature leverages the []byte data type. As a result, BINARY
data cannot be bound without the binding parameter flag. In the following
example, sf is an alias for the gosnowflake package:

	var b = []byte{0x01, 0x02, 0x03}
	_, err = stmt.Exec(sf.DataTypeBinary, b)

Maximum Number of Result Set Chunk Downloader

The driver directly downloads a result set from the cloud storage if the size is large. It is
required to shift workloads from the Snowflake database to the clients for scale. The download takes place by goroutine
named "Chunk Downloader" asynchronously so that the driver can fetch the next result set while the application can
consume the current result set.

The application may change the number of result set chunk downloader if required. Note this does not help reduce
memory footprint by itself. Consider Custom JSON Decoder.

	import (
		sf "github.com/snowflakedb/gosnowflake"
	)
	sf.MaxChunkDownloadWorkers = 2


Custom JSON Decoder for Parsing Result Set (Experimental)

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

The function db.ExecContext() returns a single result, which is the sum of the number of rows changed by each
individual statement. For example, if your multi-statement query executed two UPDATE statements, each of which
updated 10 rows, then the result returned would be 20. Individual row counts for individual statements are not
available.

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

Note: PUT statements are not supported for multi-statement queries.


If a SQL statement passed to ExecQuery() or QueryContext() fails to compile or execute, that statement is
aborted, and subsequent statements are not executed. Any statements prior to the aborted statement are unaffected.

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


Asynchronous Queries

The Go Snowflake Driver supports asynchronous execution of SQL statements.
Asynchronous execution allows you to start executing a statement and then
retrieve the result later without being blocked while waiting. While waiting
for the result of a SQL statement, you can perform other tasks, including
executing other SQL statements.

Most of the steps to execute an asynchronous query are the same as the
steps to execute a synchronous query. However, there is an additional step,
which is that you must call the WithAsyncMode() function to update
your Context object to specify that asynchronous mode is enabled.

In the code below, the call to "WithAsyncMode()" is specific
to asynchronous mode. The rest of the code is compatible with both
asynchronous mode and synchronous mode.

	...

	// Update your Context object to specify asynchronous mode:
	ctx := WithAsyncMode(context.Background())

	// Execute your query as usual by calling:
	rows, _ := db.QueryContext(ctx, query_string)

	// Retrieve the results as usual by calling:
	for rows.Next()  {
		err := rows.Scan(...)
		...
	}

The function db.QueryContext() returns an object of type snowflakeRows
regardless of whether the query is synchronous or asynchronous. However:

	* If the query is synchronous, then db.QueryContext() does not return until
		the query has finished and the result set has been loaded into the
		snowflakeRows object.
	* If the query is asynchronous, then db.QueryContext() returns a
		potentially incomplete snowflakeRows object that is filled in later
		in the background.

The call to the Next() function of snowflakeRows is always synchronous (i.e. blocking).
If the query has not yet completed and the snowflakeRows object (named "rows" in this
example) has not been filled in yet, then rows.Next() waits until the result set has been filled in.

More generally, calls to any Golang SQL API function implemented in snowflakeRows or
snowflakeResult are blocking calls, and wait if results are not yet available.
(Examples of other synchronous calls include: snowflakeRows.Err(), snowflakeRows.Columns(),
snowflakeRows.columnTypes(), snowflakeRows.Scan(), and snowflakeResult.RowsAffected().)

Because the example code above executes only one query and no other activity, there is
no significant difference in behavior between asynchronous and synchronous behavior.
The differences become significant if, for example, you want to perform some other
activity after the query starts and before it completes. The example code below starts
multiple queries, which run in the background, and then retrieves the results later.

This example uses small SELECT statements that do not retrieve enough data to require
asynchronous handling. However, the technique works for larger data sets, and for
situations where the programmer might want to do other work after starting the queries
and before retrieving the results.

	package gosnowflake

	import  (
		"context"
		"database/sql"
		"database/sql/driver"
		"fmt"
		"log"
		"os"
		sf "github.com/snowflakedb/gosnowflake"
    )

	...

	func DemonstrateAsyncMode(db *sql.DB) {
		// Enable asynchronous mode.
		ctx := WithAsyncMode(context.Background())
		// Establish connection
		conn, _ := db.Conn(ctx)

		// Unwrap connection
		err = conn.Raw(func(x interface{}) error {
			// Execute asynchronous query
			rows, _ := x.(driver.QueryerContext).QueryContext(ctx, "select 1", nil)
			defer rows.Close()

			// Retrieve and check results of the query after casting the result
			status := rows.(SnowflakeResult).GetStatus()
			if status == QueryStatusComplete {
				// do something
			} else if status == QueryStatusInProgress {
				// do something
			} else if status == QueryFailed {
				// do something
			}
			return nil
		})
	}


Support For PUT and GET

The Go Snowflake Driver supports the PUT and GET commands.

The PUT command copies a file from a local computer (the computer where the
Golang client is running) to a stage on the cloud platform. The GET command
copies data files from a stage on the cloud platform to a local computer.

See the following for information on the syntax and supported parameters:

  * PUT: https://docs.snowflake.com/en/sql-reference/sql/put.html
  * GET: https://docs.snowflake.com/en/sql-reference/sql/get.html

Using PUT

The following example shows how to run a PUT command by passing a string to the
db.Query() function:

  db.Query("PUT file://<local_file> <stage_identifier> <optional_parameters>")

"<local_file>" should include the file path as well as the name. Snowflake recommends
using an absolute path rather than a relative path. For example:

	db.Query("PUT file:///tmp/my_data_file @~ auto_compress=false overwrite=false")

Different client platforms (e.g. linux, Windows) have different path name
conventions. Ensure that you specify path names appropriately. This is
particularly important on Windows, which uses the backslash character as
both an escape character and as a separator in path names.

To send information from a stream (rather than a file) use code similar to the code below.
(The ReplaceAll() function is needed on Windows to handle backslashes in the path to the file.)

	fileStream, _ := os.OpenFile(fname, os.O_RDONLY, os.ModePerm)
	defer func() {
		if fileStream != nil {
			fileStream.Close()
		}
	} ()

	sql := "put 'file://%v' @%%%v auto_compress=true parallel=30"
	sqlText := fmt.Sprintf(sql,
		strings.ReplaceAll(fname, "\\", "\\\\"),
		tableName)
	dbt.mustExecContext(WithFileStream(context.Background(), fileStream),
		sqlText)

Note: PUT statements are not supported for multi-statement queries.

Using GET

The following example shows how to run a GET command by passing a string to the
db.Query() function:

  db.Query("GET file://<local_file> <stage_identifier> <optional_parameters>")

"<local_file>" should include the file path as well as the name. Snowflake recommends using
an absolute path rather than a relative path. For example:

  db.Query("GET file:///tmp/my_data_file @~ auto_compress=false overwrite=false")

*/
package gosnowflake
