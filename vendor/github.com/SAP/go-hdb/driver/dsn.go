// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

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
