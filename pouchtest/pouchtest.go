package pouchtest

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type TestExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (SQLRows, error)
	QueryRow(query string, args ...interface{}) SQLRow

	RegisterExec(id string, s sql.Result)
	RegisterQuery(id string, s SQLRows)
	RegisterQueryRow(id string, s SQLRow)
}

type SQLRow interface {
	SQLScannable
}

type SQLScannable interface {
	Scan(dest ...interface{}) error
}

type SQLRows interface {
	Close() error
	Columns() ([]string, error)
	Err() error
	Next() error
	SQLScannable
}

type testExecutor struct {
	results map[string]sql.Result
	rows    map[string]SQLRows
	row     map[string]SQLRow
}

type sqlResult struct {
	lastInsertId, rowsAffected     int64
	lastInsertErr, rowsAffectedErr error
}

func (s sqlResult) LastInsertId() (int64, error) {
	return s.lastInsertId, s.lastInsertErr
}

func (s sqlResult) RowsAffected() (int64, error) {
	return s.rowsAffected, s.rowsAffectedErr
}

func (t *testExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {
	resultID := fmt.Sprintf(strings.Replace(query, "?", "%v", -1), args...)
	if r, ok := t.results[resultID]; ok {
		// TODO(ttacon): allow errors to be specified to be returned instead
		return r, nil
	}
	return nil, errors.New("invalid query/args combo, not registered with TestExecutor")
}

func (t *testExecutor) Query(query string, args ...interface{}) (SQLRows, error) {
	queryID := fmt.Sprintf(strings.Replace(query, "?", "%v", -1), args...)
	if q, ok := t.rows[queryID]; ok {
		// TODO(ttacon): allow errors to be specified to be returned instead
		return q, nil
	}
	return nil, errors.New("invalid query/args combo, not registered with TestExecutor")
}

func (t *testExecutor) QueryRow(query string, args ...interface{}) SQLRow {
	queryRowID := fmt.Sprintf(strings.Replace(query, "?", "%v", -1), args...)
	if q, ok := t.row[queryRowID]; ok {
		// TODO(ttacon): allow errors to be specified to be returned instead
		return q
	}
	return nil
}

func (t *testExecutor) RegisterExec(id string, s sql.Result) { t.results[id] = s }
func (t *testExecutor) RegisterQuery(id string, s SQLRows)   { t.rows[id] = s }
func (t *testExecutor) RegisterQueryRow(id string, s SQLRow) { t.row[id] = s }
