package impl

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ttacon/pouch"
)

func Test_Pouch(t *testing.T) {
	dbConn, err := sql.Open("mysql", "pouch:pouch@/pouch")
	if err != nil {
		t.Error(err)
		return
	}

	p := SQLPouch(dbConn)
	var f Food
	f.ID = 1
	err = p.Find(&f)

	if err != nil {
		t.Error("err should be nil, was: ", err)
	}

	if f.Name != "spinach" {
		t.Error("name should have been 'spinach', was: ", f.Name)
	}

	if f.Nil != nil {
		t.Error("Nil should have been nil, was: ", f.Nil)
	}
}

func Test_create(t *testing.T) {
	dbConn, err := sql.Open("mysql", "pouch:pouch@/pouch")
	if err != nil {
		t.Error(err)
		return
	}

	p := SQLPouch(dbConn)
	var f = Food{
		Name: "spinach",
	}
	err = p.Create(&f)

	if err != nil {
		t.Error("err should have been nil, was: ", err)
	}

	if f.ID != 3 {
		t.Error("ID should have been 3, was: ", f.ID)
	}
}

func pString(s string) *string {
	return &s
}

func Test_update(t *testing.T) {
	dbConn, err := sql.Open("mysql", "pouch:pouch@/pouch")
	if err != nil {
		t.Error(err)
		return
	}

	p := SQLPouch(dbConn)
	var f = Food{
		ID:  2,
		Nil: pString("YUMMY"),
	}
	err = p.Update(&f)

	if err != nil {
		t.Error("err should have been nil, was: ", err)
	}

	if f.Name != "" {
		t.Error("Name should have been empty/overridden, was: ", f.Name)
	}
}

func Test_delete(t *testing.T) {
	dbConn, err := sql.Open("mysql", "pouch:pouch@/pouch")
	if err != nil {
		t.Error(err)
		return
	}

	p := SQLPouch(dbConn)
	var f = Food{
		ID: 5,
	}
	err = p.Delete(&f)

	if err != nil {
		t.Error("err should have been nil, was: ", err)
	}
}

type Food struct {
	ID   int
	Name string
	Nil  *string
}

func (f *Food) IdentifiableFields() ([]string, []interface{}) {
	return []string{"ID"}, []interface{}{f.ID}
}

func (f *Food) GetFieldsFor([]string) []interface{} {
	return nil
}

func (f *Food) GetAllFields() ([]string, []interface{}) {
	return []string{"ID", "Name", "NullableField"}, []interface{}{
		&f.ID, &f.Name, &f.Nil,
	}
}

func (f *Food) Table() string {
	return "Food"
}

func (f *Food) SetIdentifier(i interface{}) error {
	id, _ := i.(int64)
	f.ID = int(id)
	return nil
}

func (f *Food) FieldsFor(cols []string) []interface{} {
	var vals = make([]interface{}, len(cols))
	for i, col := range cols {
		if col == "ID" {
			vals[i] = f.ID
		} else if col == "Name" {
			vals[i] = f.Name
		} else if col == "Nil" {
			vals[i] = f.Nil
		}
	}
	return vals
}

func (f *Food) FindableCopy() pouch.Findable {
	return &Food{}
}

func (f *Food) InsertableFields() ([]string, []interface{}) {
	var cols []string
	var vals []interface{}

	cols = append(cols, "Name")
	vals = append(vals, f.Name)

	if f.Nil != nil {
		cols = append(cols, "NullableField")
		vals = append(vals, *f.Nil)
	}
	return cols, vals
}
