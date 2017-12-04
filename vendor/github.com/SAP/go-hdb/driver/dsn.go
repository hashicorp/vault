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

import (
	"net/url"
	"strconv"

	p "github.com/SAP/go-hdb/internal/protocol"
)

// DSN query parameters. For parameter client locale see http://help.sap.com/hana/SAP_HANA_SQL_Command_Network_Protocol_Reference_en.pdf.
const (
	DSNLocale  = "locale"  // Client locale as described in the protocol reference.
	DSNTimeout = "timeout" // Driver side connection timeout in seconds.
)
const (
	dsnBufferSize = "bufferSize"
	dsnFetchSize  = "fetchSize"
)

// DSN query default values.
const (
	DSNDefaultTimeout = 300 // Default value connection timeout (300 seconds = 5 minutes).
)
const (
	dsnDefaultFetchSize = 128
)

/*
DSN is here for the purposes of documentation only. A DSN string is an URL string with the following format

	hdb://<username>:<password>@<host address>:<port number>

and optional query parameters (see DSN query parameters and DSN query default values).

Example:

	hdb://myuser:mypassword@localhost:30015?timeout=60
*/
type DSN string

func parseDSN(dsn string) (*p.SessionPrm, error) {

	url, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	prm := &p.SessionPrm{Host: url.Host}

	if url.User != nil {
		prm.Username = url.User.Username()
		prm.Password, _ = url.User.Password()
	}

	values := url.Query()

	prm.BufferSize, _ = strconv.Atoi(values.Get(dsnBufferSize))

	prm.FetchSize, err = strconv.Atoi(values.Get(dsnFetchSize))
	if err != nil {
		prm.FetchSize = dsnDefaultFetchSize
	}
	prm.Timeout, err = strconv.Atoi(values.Get(DSNTimeout))
	if err != nil {
		prm.Timeout = DSNDefaultTimeout
	}

	prm.Locale = values.Get(DSNLocale)

	return prm, nil
}
