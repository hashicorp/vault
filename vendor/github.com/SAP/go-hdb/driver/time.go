// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"database/sql"
)

// NullTime represents an time.Time that may be null.
// Deprecated: Please use database/sql NullTime instead.
type NullTime = sql.NullTime
