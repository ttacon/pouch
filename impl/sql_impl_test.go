package impl

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ttacon/pouch"
)

var (
	username    = flag.String("u", "pouch", "username for db")
	password    = flag.String("p", "pouch", "password for db")
	database    = flag.String("db", "pouch", "database to connect to")
	dbTable     = flag.String("dbt", "Food", "table to use as scratch for testing")
	forceCreate = flag.Bool("force-create", false, "create the data base if it doesn't exist")

	dbURI string
)

func init() {
	flag.Parse()

	dbURI = fmt.Sprintf("%s:%s@/%s", *username, *password, *database)
	dbConn, err := sql.Open("mysql", dbURI)
	if err != nil {
		panic("failed to connect to db, err: " + err.Error())
	}
	if *forceCreate {
		err = createTable(dbConn)
		if err != nil {
			panic("failed to force creation of table, err: " + err.Error())
		}
	}
	err = cleanDB(dbConn)
	if err != nil {
		panic("failed to clean db, err: " + err.Error())
	}
}

func cleanDB(db pouch.Executor) error {
	_, err := db.Exec(fmt.Sprintf(`
delete from %s
`, *dbTable))
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf(`
ALTER TABLE %s AUTO_INCREMENT = 1
`, *dbTable))
	return err
}

func createTable(db pouch.Executor) error {
	_, err := db.Exec(fmt.Sprintf(`
create table %s if not exists (
  ID int primary key auto_increment,
  Name varchar(64) not null,
  NullableString varchar(64)
) engine=InnoDB;`), *dbTable)
	return err
}

func Test_create(t *testing.T) {
	dbConn, err := sql.Open("mysql", dbURI)
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

	if f.ID != 1 {
		t.Error("ID should have been 1, was: ", f.ID)
	}
}

func Test_Pouch(t *testing.T) {
	dbConn, err := sql.Open("mysql", dbURI)
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

func pString(s string) *string {
	return &s
}

func Test_update(t *testing.T) {
	dbConn, err := sql.Open("mysql", dbURI)
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
	dbConn, err := sql.Open("mysql", dbURI)
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

func (f *Food) Merge(m pouch.Gettable) error {
	fNames, fVals := m.GetAllFields()
	if len(fNames) != len(fVals) {
		return errors.New("when merging foods, must have slices of same length")
	}

	for i, name := range fNames {
		if name == "ID" {
			f.ID = *fVals[i].(*int)
		} else if name == "Name" {
			f.Name = *fVals[i].(*string)
		} else if name == "Nil" {
			f.Nil = fVals[i].(*string)
		}
	}
	return nil
}

type mapFood struct {
	*Food
}

func (m *mapFood) Table() string {
	return "food:%v"
}
