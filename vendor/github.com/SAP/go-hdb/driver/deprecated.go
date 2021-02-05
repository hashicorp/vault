// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"database/sql/driver"
)

// deprecated driver interface methods

// Prepare implements the driver.Conn interface.
func (*Conn) Prepare(query string) (driver.Stmt, error) { panic("deprecated") }

// Begin implements the driver.Conn interface.
func (*Conn) Begin() (driver.Tx, error) { panic("deprecated") }

// Exec implements the driver.Execer interface.
func (*Conn) Exec(query string, args []driver.Value) (driver.Result, error) { panic("deprecated") }

// Query implements the driver.Queryer interface.
func (*Conn) Query(query string, args []driver.Value) (driver.Rows, error) { panic("deprecated") }

func (*stmt) Exec(args []driver.Value) (driver.Result, error)             { panic("deprecated") }
func (*stmt) Query(args []driver.Value) (rows driver.Rows, err error)     { panic("deprecated") }
func (*callStmt) Exec(args []driver.Value) (driver.Result, error)         { panic("deprecated") }
func (*callStmt) Query(args []driver.Value) (rows driver.Rows, err error) { panic("deprecated") }

// replaced driver interface methods
// sql.Stmt.ColumnConverter --> replaced by CheckNamedValue
