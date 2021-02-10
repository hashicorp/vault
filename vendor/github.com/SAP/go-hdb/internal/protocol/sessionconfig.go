// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"github.com/SAP/go-hdb/internal/container/vermap"
)

// SessionConfig represents the session relevant driver connector options.
type SessionConfig struct {
	DriverVersion, DriverName           string
	ApplicationName, Username, Password string
	SessionVariables                    *vermap.VerMap
	Locale                              string
	FetchSize, LobChunkSize             int
	Dfv                                 int
	Legacy                              bool
}
