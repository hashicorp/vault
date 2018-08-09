/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

// DSN parameters. For parameter client locale see http://help.sap.com/hana/SAP_HANA_SQL_Command_Network_Protocol_Reference_en.pdf.
const (
	DSNLocale    = "locale"    // Client locale as described in the protocol reference.
	DSNTimeout   = "timeout"   // Driver side connection timeout in seconds.
	DSNFetchSize = "fetchSize" // Maximum number of fetched records from database by database/sql/driver/Rows.Next().
)

/*
DSN TLS parameters.
For more information please see https://golang.org/pkg/crypto/tls/#Config.
For more flexibility in TLS configuration please see driver.Connector.
*/
const (
	DSNTLSRootCAFile         = "TLSRootCAFile"         // Path,- filename to root certificate(s).
	DSNTLSServerName         = "TLSServerName"         // ServerName to verify the hostname.
	DSNTLSInsecureSkipVerify = "TLSInsecureSkipVerify" // Controls whether a client verifies the server's certificate chain and host name.
)

// DSN default values.
const (
	DefaultTimeout   = 300 // Default value connection timeout (300 seconds = 5 minutes).
	DefaultFetchSize = 128 // Default value fetchSize.
)

// DSN minimal values.
const (
	minTimeout   = 0 // Minimal timeout value.
	minFetchSize = 1 // Minimal fetchSize value.
)

/*
DSN is here for the purposes of documentation only. A DSN string is an URL string with the following format

	"hdb://<username>:<password>@<host address>:<port number>"

and optional query parameters (see DSN query parameters and DSN query default values).

Example:
	"hdb://myuser:mypassword@localhost:30015?timeout=60"

Examples TLS connection:
	"hdb://myuser:mypassword@localhost:39013?TLSRootCAFile=trust.pem"
	"hdb://myuser:mypassword@localhost:39013?TLSRootCAFile=trust.pem&TLSServerName=hostname"
	"hdb://myuser:mypassword@localhost:39013?TLSInsecureSkipVerify"
*/
type DSN string
