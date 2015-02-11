package pouchtest

import "testing"

func TestOrderedTestExecutorExec(t *testing.T) {
	e := NewOrderedTestExecutor().
		RegisterExec(SQLResult(1, 1))

	res, err := e.Exec("GIBBERISH")
	if err != nil {
		t.Errorf("Expected err to be nil, was: %v", err)
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		t.Errorf("Expected err from LastInsertId() to be nil, was: %v", err)
	}
	if lastInsertID != 1 {
		t.Errorf("Expected last inserted id to be 1, was: %d", lastInsertID)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		t.Errorf("Expected err from RowsAffected() to be nil, was: %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("Expected rows affected to be 1, was: %d", rowsAffected)
	}
}

func TestSQLRows(t *testing.T) {
	rows := NewSQLRows(
		[]string{"ID", "Name"},
		[]SQLRow{
			NewSQLRow([]interface{}{1, "Blossom"}),
			NewSQLRow([]interface{}{2, "Bubbles"}),
			NewSQLRow([]interface{}{3, "Buttercup"}),
		},
	)

	if err := rows.Err(); err != nil {
		t.Error("err should have been nil, was: ", err)
	}

	var (
		id   int
		name string
	)
	err := rows.Scan(&id, &name)
	if err == nil || err.Error() != "must call Next() before scan" {
		t.Error("should not have been nil, was: ", err)
	}

	// get Blossom
	if ok := rows.Next(); !ok {
		t.Error("rows.Next() should have been ok, it wasn't")
	}
	err = rows.Scan(&id, &name)
	if err != nil {
		t.Error("should have been nil, was: ", err)
	}
	if id != 1 {
		t.Error("id should have been 1, was: ", id)
	}
	if name != "Blossom" {
		t.Error("name should have been Blossom, was: ", name)
	}

	// get Bubbles
	if ok := rows.Next(); !ok {
		t.Error("rows.Next() should have been ok, it wasn't")
	}
	err = rows.Scan(&id, &name)
	if err != nil {
		t.Error("should have been nil, was: ", err)
	}
	if id != 2 {
		t.Error("id should have been 2, was: ", id)
	}
	if name != "Bubbles" {
		t.Error("name should have been Bubbles, was: ", name)
	}

	// get Buttercup
	if ok := rows.Next(); !ok {
		t.Error("rows.Next() should have been ok, it wasn't")
	}
	err = rows.Scan(&id, &name)
	if err != nil {
		t.Error("should have been nil, was: ", err)
	}
	if id != 3 {
		t.Error("id should have been 3, was: ", id)
	}
	if name != "Buttercup" {
		t.Error("name should have been Buttercup, was: ", name)
	}

	if ok := rows.Next(); ok {
		t.Error("rows.Next() should not have been ok, it was")
	}
}

func TestOrderedTestExecutorQuery(t *testing.T) {
	query := NewSQLRows(
		[]string{"ID", "Name"},
		[]SQLRow{
			NewSQLRow([]interface{}{1, "HelloThur!"}),
		},
	)
	e := NewOrderedTestExecutor().
		RegisterQuery(query)

	rows, err := e.Query("YOLO, IT DOESN'T MATTER!")
	if err != nil {
		t.Error("error should have been nil was: ", err)
	}
	if ok := rows.Next(); !ok {
		t.Error("rows.Next() should have been ok, it wasn't")
	}

}

func TestOrderedTestExecutorQueryRow(t *testing.T) {
	e := NewOrderedTestExecutor().
		RegisterQueryRow(NewSQLRow([]interface{}{1.0, []byte("HelloThur!")}))

	row := e.QueryRow("yolo")
	var (
		entropy float64
		data    []byte
	)
	err := row.Scan(&entropy, &data)
	if err != nil {
		t.Error("err should have been nil, was: ", err)
	}

	if entropy != 1.0 {
		// ruh-roh float equality
		t.Error("entropy should have been 1, was: ", entropy)
	}
	if string(data) != "HelloThur!" {
		t.Error("got unexpected data back: ", string(data))
	}
}
