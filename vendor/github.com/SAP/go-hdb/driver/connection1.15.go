// +build go1.15

// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"database/sql/driver"
)

//  check if conn implements all required interfaces
var (
	_ driver.Validator = (*Conn)(nil)
)

// IsValid implements the driver.Validator interface.
func (c *Conn) IsValid() bool {
	c.lock()
	defer c.unlock()

	return !c.dbConn.isBad()
}
