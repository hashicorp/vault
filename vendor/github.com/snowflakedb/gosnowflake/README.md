## Support

For official support and urgent, production-impacting issues, please [contact Snowflake Support](https://community.snowflake.com/s/article/How-To-Submit-a-Support-Case-in-Snowflake-Lodge).

# Go Snowflake Driver

<a href="https://codecov.io/github/snowflakedb/gosnowflake?branch=master">
    <img alt="Coverage" src="https://codecov.io/github/snowflakedb/gosnowflake/coverage.svg?branch=master">
</a>
<a href="https://github.com/snowflakedb/gosnowflake/actions?query=workflow%3A%22Build+and+Test%22">
    <img src="https://github.com/snowflakedb/gosnowflake/workflows/Build%20and%20Test/badge.svg?branch=master">
</a>
<a href="http://www.apache.org/licenses/LICENSE-2.0.txt">
    <img src="http://img.shields.io/:license-Apache%202-brightgreen.svg">
</a>
<a href="https://goreportcard.com/report/github.com/snowflakedb/gosnowflake">
    <img src="https://goreportcard.com/badge/github.com/snowflakedb/gosnowflake">
</a>

This topic provides instructions for installing, running, and modifying the Go Snowflake Driver. The driver supports Go's [database/sql](https://golang.org/pkg/database/sql/) package.

# Prerequisites

The following software packages are required to use the Go Snowflake Driver.

## Go

The latest driver requires the [Go language](https://golang.org/) 1.20 or higher. The supported operating systems are Linux, Mac OS, and Windows, but you may run the driver on other platforms if the Go language works correctly on those platforms.


# Installation

If you don't have a project initialized, set it up.

```sh
go mod init example.com/snowflake
```

Get Gosnowflake source code, if not installed.

```sh
go get -u github.com/snowflakedb/gosnowflake
```

# Docs

For detailed documentation and basic usage examples, please see the documentation at
[godoc.org](https://godoc.org/github.com/snowflakedb/gosnowflake/).

## Note

This driver currently does not support GCP regional endpoints. Please ensure that any workloads using through this driver do not require support for regional endpoints on GCP. If you have questions about this, please contact Snowflake Support.

# Sample Programs

Snowflake provides a set of sample programs to test with. Set the environment variable ``$GOPATH`` to the top directory of your workspace, e.g., ``~/go`` and make certain to
include ``$GOPATH/bin`` in the environment variable ``$PATH``. Run the ``make`` command to build all sample programs.

```
make install
```

In the following example, the program ``select1.go`` is built and installed in ``$GOPATH/bin`` and can be run from the command line:

```
SNOWFLAKE_TEST_ACCOUNT=<your_account> \
SNOWFLAKE_TEST_USER=<your_user> \
SNOWFLAKE_TEST_PASSWORD=<your_password> \
select1
Congrats! You have successfully run SELECT 1 with Snowflake DB!
```

# Development

The developer notes are hosted with the source code on [GitHub](https://github.com/snowflakedb/gosnowflake).

## Testing Code


Set the Snowflake connection info in ``parameters.json``:

```
{
    "testconnection": {
        "SNOWFLAKE_TEST_USER":      "<your_user>",
        "SNOWFLAKE_TEST_PASSWORD":  "<your_password>",
        "SNOWFLAKE_TEST_ACCOUNT":   "<your_account>",
        "SNOWFLAKE_TEST_WAREHOUSE": "<your_warehouse>",
        "SNOWFLAKE_TEST_DATABASE":  "<your_database>",
        "SNOWFLAKE_TEST_SCHEMA":    "<your_schema>",
        "SNOWFLAKE_TEST_ROLE":      "<your_role>"
    }
}
```

Install [jq](https://stedolan.github.io/jq) so that the parameters can get parsed correctly, and run ``make test`` in your Go development environment:

```
make test
```

## customizing Logging Tags

If you would like to ensure that certain tags are always present in the logs, `RegisterClientLogContextHook` can be used in your init function. See example below.
```
import "github.com/snowflakedb/gosnowflake"

func init() {
    // each time the logger is used, the logs will contain a REQUEST_ID field with requestID the value extracted 
    // from the context
	gosnowflake.RegisterClientLogContextHook("REQUEST_ID", func(ctx context.Context) interface{} {
		return requestIdFromContext(ctx)
	})
}
```

## Setting Log Level
If you want to change the log level, `SetLogLevel` can be used in your init function like this:
```
import "github.com/snowflakedb/gosnowflake"

func init() {
    // The following line changes the log level to debug
	_ = gosnowflake.GetLogger().SetLogLevel("debug")
}
```
The following is a list of options you can pass in to set the level from least to most verbose: 
- `"OFF"`
- `"error"`
- `"warn"`
- `"print"`
- `"trace"`
- `"debug"`
- `"info"`


## Capturing Code Coverage

Configure your testing environment as described above and run ``make cov``. The coverage percentage will be printed on the console when the testing completes.

```
make cov
```

For more detailed analysis, results are printed to ``coverage.txt`` in the project directory.

To read the coverage report, run:

```
go tool cover -html=coverage.txt
```

## Submitting Pull Requests

You may use your preferred editor to edit the driver code. Make certain to run ``make fmt lint`` before submitting any pull request to Snowflake. This command formats your source code according to the standard Go style and detects any coding style issues.

## Runaway `dbus-daemon` processes on certain OS
This only affects certain Linux distributions, one of them is confirmed to be RHEL. Due to a bug in one of the dependencies (`keyring`),
on the affected OS, each invocation of a program depending on gosnowflake (or any other program depending on the same `keyring`),
will generate a new instance of `dbus-daemon` fork which can, due to not being cleaned up, eventually fill the fd limits.

Until we replace the offending dependency with one which doesn't have the bug, a workaround needs to be applied, which can be:
* cleaning up the runaway processes periodically
* setting envvar `DBUS_SESSION_BUS_ADDRESS=$XDG_RUNTIME_DIR/bus` (if that socket exists, or create it) or even `DBUS_SESSION_BUS_ADDRESS=/dev/null`

Details in [issue 773](https://github.com/snowflakedb/gosnowflake/issues/773)
