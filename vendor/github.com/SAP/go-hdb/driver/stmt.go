package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"slices"
)

// check if statements implements all required interfaces.
var (
	_ driver.Stmt              = (*stmt)(nil)
	_ driver.StmtExecContext   = (*stmt)(nil)
	_ driver.StmtQueryContext  = (*stmt)(nil)
	_ driver.NamedValueChecker = (*stmt)(nil)
)

type stmt struct {
	session *session
	attrs   *connAttrs
	metrics *metrics
	query   string
	pr      *prepareResult
	// rows: stored procedures with table output parameters
	rows *sql.Rows
}

type totalRowsAffected int64

func (t *totalRowsAffected) add(r driver.Result) {
	if r == nil {
		return
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return
	}
	*t += totalRowsAffected(rows)
}

func newStmt(session *session, attrs *connAttrs, metrics *metrics, query string, pr *prepareResult) *stmt {
	metrics.msgCh <- gaugeMsg{idx: gaugeStmt, v: 1} // increment number of statements.
	return &stmt{session: session, attrs: attrs, metrics: metrics, query: query, pr: pr}
}

/*
NumInput differs dependent on statement (check is done in QueryContext and ExecContext):
- #args == #param (only in params):    query, exec, exec bulk (non control query)
- #args == #param (in and out params): exec call
- #args == 0:                          exec bulk (control query)
- #args == #input param:               query call.
*/
func (s *stmt) NumInput() int { return -1 }

func (s *stmt) Close() error {
	s.metrics.msgCh <- gaugeMsg{idx: gaugeStmt, v: -1} // decrement number of statements.

	if s.rows != nil {
		s.rows.Close()
	}

	if s.session.isBad() {
		return driver.ErrBadConn
	}
	return s.session.dropStatementID(context.Background(), s.pr.stmtID)
}

// CheckNamedValue implements NamedValueChecker interface.
func (s *stmt) CheckNamedValue(nv *driver.NamedValue) error {
	// conversion is happening as part of the exec, query call
	return nil
}

func (s *stmt) QueryContext(ctx context.Context, nvargs []driver.NamedValue) (driver.Rows, error) {
	if s.pr.isProcedureCall() {
		return nil, fmt.Errorf("invalid procedure call %s - please use Exec instead", s.query)
	}

	if err := s.session.preventSwitchUser(ctx); err != nil {
		return nil, err
	}
	return s.session.query(ctx, s.query, s.pr, nvargs)
}

func (s *stmt) ExecContext(ctx context.Context, nvargs []driver.NamedValue) (driver.Result, error) {
	if hookFn, ok := ctx.Value(connHookCtxKey).(connHookFn); ok {
		hookFn(choStmtExec)
	}

	var (
		result driver.Result
		err    error
	)

	if err := s.session.preventSwitchUser(ctx); err != nil {
		return nil, err
	}
	if s.pr.isProcedureCall() {
		result, s.rows, err = s.execCall(ctx, s.pr, nvargs)
	} else {
		result, err = s.execDefault(ctx, nvargs)
	}
	return result, err
}

func (s *stmt) execCall(ctx context.Context, pr *prepareResult, nvargs []driver.NamedValue) (driver.Result, *sql.Rows, error) {
	/*
		call without lob input parameters:
		--> callResult output parameter values are set after read call
		call with lob output parameters:
		--> callResult output parameter values are set after last lob input write
	*/

	cr, callArgs, numRow, err := s.session.execCall(ctx, s.query, pr, nvargs)
	if err != nil {
		return nil, nil, err
	}

	numOutArgs := len(callArgs.outArgs)
	// no output args -> done
	if numOutArgs == 0 {
		return driver.RowsAffected(numRow), nil, nil
	}

	numOutputField := len(cr.outFields)
	scanArgs := make([]any, numOutputField)
	for i := range numOutArgs {
		scanArgs[i] = callArgs.outArgs[i].Value.(sql.Out).Dest
	}
	// acccount for table output fields without call arguments.
	for i := numOutArgs; i < numOutputField; i++ {
		scanArgs[i] = new(sql.Rows)
	}

	// no table output parameters -> QueryRow
	if len(callArgs.outFields) == numOutArgs {
		if err := stdConnTracker.callDB().QueryRow("", cr).Scan(scanArgs...); err != nil {
			return nil, nil, err
		}
		return driver.RowsAffected(numRow), nil, nil
	}

	// table output parameters -> Query (needs to kept open)
	rows, err := stdConnTracker.callDB().Query("", cr)
	if err != nil {
		return nil, rows, err
	}
	if !rows.Next() {
		return nil, rows, rows.Err()
	}
	if err := rows.Scan(scanArgs...); err != nil {
		return nil, rows, err
	}
	return driver.RowsAffected(numRow), rows, nil
}

type execFn func(ctx context.Context, nvargs []driver.NamedValue) (driver.Result, error)

func (s *stmt) execDefault(ctx context.Context, nvargs []driver.NamedValue) (driver.Result, error) {
	numNVArg, numField := len(nvargs), s.pr.numField()

	if numNVArg == 0 {
		if numField != 0 {
			return nil, fmt.Errorf("invalid number of arguments %d - expected %d", numNVArg, numField)
		}
		return s.session.exec(ctx, s.query, s.pr, nvargs, 0)
	}
	if numNVArg == 1 {
		if execFn := s.detectExecFn(nvargs[0]); execFn != nil {
			return execFn(ctx, nvargs)
		}
	}
	if numNVArg == numField {
		return s.exec(ctx, s.pr, nvargs, 0)
	}
	if numNVArg%numField != 0 {
		return nil, fmt.Errorf("invalid number of arguments %d - multiple of %d expected", numNVArg, numField)
	}
	return s.execMany(ctx, nvargs)
}

// ErrEndOfRows is the error to be returned using a function based bulk exec to indicate
// the end of rows.
var ErrEndOfRows = errors.New("end of rows")

/*
Non 'atomic' (transactional) operation due to the split in packages (bulkSize),
execMany data might only be written partially to the database in case of hdb stmt errors.
*/
func (s *stmt) execFct(ctx context.Context, nvargs []driver.NamedValue) (driver.Result, error) {
	totalRowsAffected := totalRowsAffected(0)
	args := make([]driver.NamedValue, 0, s.pr.numField())
	scanArgs := make([]any, s.pr.numField())

	fct, ok := nvargs[0].Value.(func(args []any) error)
	if !ok {
		panic("invalid argument") // should never happen
	}

	done := false
	batch := 0
	for !done {
		args = args[:0]
		for range s.attrs._bulkSize {
			err := fct(scanArgs)
			if errors.Is(err, ErrEndOfRows) {
				done = true
				break
			}
			if err != nil {
				return driver.RowsAffected(totalRowsAffected), err
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
		}

		r, err := s.exec(ctx, s.pr, args, batch*s.attrs._bulkSize)
		totalRowsAffected.add(r)
		if err != nil {
			return driver.RowsAffected(totalRowsAffected), err
		}
		batch++
	}
	return driver.RowsAffected(totalRowsAffected), nil
}

/*
Non 'atomic' (transactional) operation due to the split in packages (bulkSize),
execMany data might only be written partially to the database in case of hdb stmt errors.
*/
func (s *stmt) execMany(ctx context.Context, nvargs []driver.NamedValue) (driver.Result, error) {
	bulkSize := s.attrs._bulkSize

	totalRowsAffected := totalRowsAffected(0)
	numField := s.pr.numField()
	numNVArg := len(nvargs)
	numRec := numNVArg / numField
	numBatch := numRec / bulkSize
	if numRec%bulkSize != 0 {
		numBatch++
	}

	for i := range numBatch {
		from := i * numField * bulkSize
		to := (i + 1) * numField * bulkSize
		if to > numNVArg {
			to = numNVArg
		}
		r, err := s.exec(ctx, s.pr, nvargs[from:to], i*bulkSize)
		totalRowsAffected.add(r)
		if err != nil {
			return driver.RowsAffected(totalRowsAffected), err
		}
	}
	return driver.RowsAffected(totalRowsAffected), nil
}

/*
exec executes a sql statement.

Bulk insert containing LOBs:
  - Precondition:
    .Sending more than one row with partial LOB data.
  - Observations:
    .In hdb version 1 and 2 'piecewise' LOB writing does work.
    .Same does not work in case of geo fields which are LOBs en,- decoded as well.
    .In hana version 4 'piecewise' LOB writing seems not to work anymore at all.
  - Server implementation (not documented):
    .'piecewise' LOB writing is only supported for the last row of a 'bulk insert'.
  - Current implementation:
    One server call in case of
    .'non bulk' execs or
    .'bulk' execs without LOBs
    else potential several server calls (split into packages).
  - Package invariant:
    .for all packages except the last one, the last row contains 'incomplete' LOB data ('piecewise' writing)
*/
func (s *stmt) exec(ctx context.Context, pr *prepareResult, nvargs []driver.NamedValue, ofs int) (driver.Result, error) {
	addLobDataRecs, err := convertExecArgs(pr.parameterFields, nvargs, s.attrs._cesu8Encoder(), s.attrs._lobChunkSize)
	if err != nil {
		return driver.ResultNoRows, err
	}

	// piecewise LOB handling
	numColumn := len(pr.parameterFields)
	totalRowsAffected := totalRowsAffected(0)
	from := 0
	for _, row := range addLobDataRecs {
		to := (row + 1) * numColumn

		r, err := s.session.exec(ctx, s.query, pr, nvargs[from:to], ofs)
		totalRowsAffected.add(r)
		if err != nil {
			return driver.RowsAffected(totalRowsAffected), err
		}
		from = to
	}
	return driver.RowsAffected(totalRowsAffected), nil
}
