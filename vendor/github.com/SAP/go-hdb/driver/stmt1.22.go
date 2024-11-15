//go:build !go1.23

package driver

import (
	"database/sql/driver"
)

func (s *stmt) detectExecFn(nvarg driver.NamedValue) execFn {
	if _, ok := nvarg.Value.(func(args []any) error); ok {
		return s.execFct
	}
	return nil
}
