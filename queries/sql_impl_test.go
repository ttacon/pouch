package queries

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ttacon/pretty"
)

func Test_Pouch(t *testing.T) {
	dbConn, err := sql.Open("mysql", "trey:chips@/dbgen")
	if err != nil {
		t.Error(err)
		return
	}

	p := SQLPouch(dbConn)
	var f Food
	f.ID = 1
	err = p.Find(&f)
	fmt.Println("err: ", err)
	pretty.Println(f)
}

func Test_create(t *testing.T) {
	dbConn, err := sql.Open("mysql", "trey:chips@/dbgen")
	if err != nil {
		t.Error(err)
		return
	}

	p := SQLPouch(dbConn)
	var f = Food{
		Name: "spinach",
	}
	err = p.Create(&f)
	fmt.Println("err: ", err)
	pretty.Println(f)
}

func Test_update(t *testing.T) {
	dbConn, err := sql.Open("mysql", "trey:chips@/dbgen")
	if err != nil {
		t.Error(err)
		return
	}

	p := SQLPouch(dbConn)
	var f = Food{
		Nil: "YUMMY",
	}
	err = p.Update(&f)
	fmt.Println("err: ", err)
	pretty.Println(f)
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
func (f *Food) InsertableFields() ([]string, []interface{}) {
	var cols []string
	var vals []interface{}

	if len(f.Name) > 0 {
		cols = append(cols, "Name")
		vals = append(vals, f.Name)
	}

	if f.Nil != nil {
		cols = append(cols, "NullableField")
		vals = append(vals, *f.Nil)
	}
	return cols, vals
}
