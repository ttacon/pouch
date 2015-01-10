package pouch

import "database/sql"

type Pouch interface {
	Queryable
	Storage
}

type Storage interface {
	Find(Findable) error
	Create(Createable) error
	Update(Updateable) error
	Delete(Deleteable) error
}

type Query interface {
	Queryable
	Storage
}

type Createable interface {
	Insertable
	Tableable
}

type Gettable interface {
	GetFieldsFor([]string) []interface{}
	GetAllFields() ([]string, []interface{})
	Tableable
}

type Tableable interface {
	Table() string
}

type Findable interface {
	Identifiable
	Gettable
}

type Insertable interface {
	InsertableFields() ([]string, []interface{})
	SetIdentifier(interface{}) error
}

type Identifiable interface {
	IdentifiableFields() ([]string, []interface{})
}

type Updateable interface {
	Insertable
	Identifiable
	Tableable
}

type Deleteable interface {
	Identifiable
	Tableable
}

type Queryable interface {
	GroupBy(spec string) Query
	OrderBy(spec string) Query
	Where(frag string, val ...interface{}) Query
	Limit(lim int) Query
	Offset(off int) Query
}

type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
