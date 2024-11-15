//go:build go1.23

package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"iter"
	"slices"
)

func (s *stmt) detectExecFn(nvarg driver.NamedValue) execFn {
	switch nvarg.Value.(type) {
	case func(args []any) error:
		return s.execFct
	case iter.Seq[[]any]:
		return s.execSeq
	default:
		return nil
	}
}

func (s *stmt) execSeq(ctx context.Context, nvargs []driver.NamedValue) (driver.Result, error) {
	totalRowsAffected := totalRowsAffected(0)
	args := make([]driver.NamedValue, 0, s.pr.numField())

	seq, ok := nvargs[0].Value.(iter.Seq[[]any])
	if !ok {
		panic("invalid argument") // should never happen
	}

	batch, n := 0, 0
	for scanArgs := range seq {
		if len(scanArgs) != s.pr.numField() {
			return driver.RowsAffected(totalRowsAffected), fmt.Errorf("invalid number of args %d - expected %d", len(scanArgs), s.pr.numField())
		}

		args = slices.Grow(args, len(scanArgs))
		for i, scanArg := range scanArgs {
			nv := driver.NamedValue{Ordinal: i + 1}
			if t, ok := scanArg.(sql.NamedArg); ok {
				nv.Name = t.Name
				nv.Value = t.Value
			} else {
				nv.Name = ""
				nv.Value = scanArg
			}
			args = append(args, nv)
		}

		n++
		if n >= s.attrs._bulkSize {
			r, err := s.exec(ctx, s.pr, args, batch*s.attrs._bulkSize)
			totalRowsAffected.add(r)
			if err != nil {
				return driver.RowsAffected(totalRowsAffected), err
			}
			args = args[:0]
			batch++
		}
	}

	if n > 0 {
		r, err := s.exec(ctx, s.pr, args, batch*s.attrs._bulkSize)
		totalRowsAffected.add(r)
		if err != nil {
			return driver.RowsAffected(totalRowsAffected), err
		}
	}

	return driver.RowsAffected(totalRowsAffected), nil
}
