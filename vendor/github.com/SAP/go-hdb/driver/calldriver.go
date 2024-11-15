package driver

import (
	"context"
	"database/sql/driver"
	"fmt"
)

// Boilerplate to define a minimal sql driver implementation.
// To be used for converting stored procedure output parameters
// including sql.Rows output table parameters to guarantee
// exactly the same conversion behavior as for sql.Query.
var (
	_ driver.Driver            = (*callDriver)(nil)
	_ driver.Connector         = (*callConnector)(nil)
	_ driver.Conn              = (*callConn)(nil)
	_ driver.NamedValueChecker = (*callConn)(nil)
	_ driver.QueryerContext    = (*callConn)(nil)
)

type callDriver struct{}

var (
	defCallDriver = &callDriver{}
	defCallConn   = &callConn{}
)

func (d *callDriver) Open(name string) (driver.Conn, error) { return defCallConn, nil }

type callConnector struct{}

func (c *callConnector) Connect(context.Context) (driver.Conn, error) { return defCallConn, nil }
func (c *callConnector) Driver() driver.Driver                        { return defCallDriver }

type callConn struct{}

func (c *callConn) Prepare(query string) (driver.Stmt, error)   { panic("not implemented") }
func (c *callConn) Close() error                                { return nil }
func (c *callConn) Begin() (driver.Tx, error)                   { panic("not implemented") }
func (c *callConn) CheckNamedValue(nv *driver.NamedValue) error { return nil }

// QueryContext is used to convert the stored procedure output parameters.
func (c *callConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid argument length %d - expected 1", len(args))
	}
	cr, ok := args[0].Value.(*callResult)
	if !ok {
		return nil, fmt.Errorf("invalid argument type %T", args[0])
	}
	return cr, nil
}
