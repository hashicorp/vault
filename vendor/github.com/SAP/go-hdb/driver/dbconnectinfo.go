package driver

import (
	"fmt"
)

// DBConnectInfo represents the connection information attributes returned by hdb.
type DBConnectInfo struct {
	DatabaseName string
	Host         string
	Port         int
	IsConnected  bool
}

func (ci *DBConnectInfo) String() string {
	return fmt.Sprintf("Database Name: %s Host: %s Port: %d connected: %t", ci.DatabaseName, ci.Host, ci.Port, ci.IsConnected)
}
