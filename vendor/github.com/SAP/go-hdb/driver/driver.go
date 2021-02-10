// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"os"
	"sync/atomic"
)

// DriverVersion is the version number of the hdb driver.
const DriverVersion = "0.102.7"

// DriverName is the driver name to use with sql.Open for hdb databases.
const DriverName = "hdb"

// default application name.
var defaultApplicationName string

var hdbDriver = &Driver{}

func init() {
	defaultApplicationName, _ = os.Executable()
	sql.Register(DriverName, hdbDriver)
}

// driver

//  check if driver implements all required interfaces
var (
	_ driver.Driver        = (*Driver)(nil)
	_ driver.DriverContext = (*Driver)(nil)
)

// Stats contains driver statistics.
type Stats struct {
	OpenConnections  int // Number of open driver connections.
	OpenTransactions int // Number of open driver transactions.
	OpenStatements   int // Number of open driver database statements.
}

// Driver represents the go sql driver implementation for hdb.
type Driver struct {
	// Atomic access only. At top of struct to prevent mis-alignment on 32-bit platforms.
	numConn int64 // Number of connections.
	numTx   int64 // Number of transactions.
	numStmt int64 // Number of statements.
}

func (d *Driver) addConn(delta int64) { atomic.AddInt64(&d.numConn, delta) }
func (d *Driver) addTx(delta int64)   { atomic.AddInt64(&d.numTx, delta) }
func (d *Driver) addStmt(delta int64) { atomic.AddInt64(&d.numStmt, delta) }

// Open implements the driver.Driver interface.
func (d *Driver) Open(dsn string) (driver.Conn, error) {
	connector, err := NewDSNConnector(dsn)
	if err != nil {
		return nil, err
	}
	return connector.Connect(context.Background())
}

// OpenConnector implements the driver.DriverContext interface.
func (d *Driver) OpenConnector(dsn string) (driver.Connector, error) { return NewDSNConnector(dsn) }

// Name returns the driver name.
func (d *Driver) Name() string { return DriverName }

// Version returns the driver version.
func (d *Driver) Version() string { return DriverVersion }

// Stats returns driver statistics.
func (d *Driver) Stats() Stats {
	return Stats{
		OpenConnections:  int(atomic.LoadInt64(&d.numConn)),
		OpenTransactions: int(atomic.LoadInt64(&d.numTx)),
		OpenStatements:   int(atomic.LoadInt64(&d.numStmt)),
	}
}
