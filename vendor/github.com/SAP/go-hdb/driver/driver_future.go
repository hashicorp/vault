// +build future

/*
Copyright 2018 SAP SE

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
	"context"
	//"database/sql"
	"database/sql/driver"

	"github.com/SAP/go-hdb/driver/sqltrace"

	p "github.com/SAP/go-hdb/internal/protocol"
)

// database connection

func (c *conn) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	if c.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	done := make(chan struct{})
	go func() {
		var (
			qt             p.QueryType
			id             uint64
			prmFieldSet    *p.ParameterFieldSet
			resultFieldSet *p.ResultFieldSet
		)
		qt, id, prmFieldSet, resultFieldSet, err = c.session.Prepare(query)
		if err != nil {
			goto done
		}
		select {
		default:
		case <-ctx.Done():
			return
		}
		stmt, err = newStmt(qt, c.session, query, id, prmFieldSet, resultFieldSet)
	done:
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return stmt, err
	}
}

// QueryContext implements the database/sql/driver/QueryerContext interface.
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	if c.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	if len(args) != 0 {
		return nil, driver.ErrSkip //fast path not possible (prepare needed)
	}

	// direct execution of call procedure
	// - returns no parameter metadata (sps 82) but only field values
	// --> let's take the 'prepare way' for stored procedures
	//	if checkCallProcedure(query) {
	//		return nil, driver.ErrSkip
	//	}

	sqltrace.Traceln(query)

	done := make(chan struct{})
	go func() {
		var (
			id             uint64
			resultFieldSet *p.ResultFieldSet
			fieldValues    *p.FieldValues
			attributes     p.PartAttributes
		)
		id, resultFieldSet, fieldValues, attributes, err = c.session.QueryDirect(query)
		if err != nil {
			goto done
		}
		select {
		default:
		case <-ctx.Done():
			return
		}
		if id == 0 { // non select query
			rows = noResult
		} else {
			rows, err = newQueryResult(c.session, id, resultFieldSet, fieldValues, attributes)
		}
	done:
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return rows, err
	}
}

//statement

func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {

	if s.session.IsBad() {
		return nil, driver.ErrBadConn
	}

	done := make(chan struct{})
	go func() {
		rows, err = s.defaultQuery(ctx, args)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return rows, err
	}
}

func (s *stmt) defaultQuery(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {

	sqltrace.Tracef("%s %v", s.query, args)

	rid, values, attributes, err := s.session.Query(s.id, s.prmFieldSet, s.resultFieldSet, args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if rid == 0 { // non select query
		return noResult, nil
	}
	return newQueryResult(s.session, rid, s.resultFieldSet, values, attributes)
}
