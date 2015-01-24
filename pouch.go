package pouch

import "database/sql"

// A Pouch is anything which can act as a backing Storage and which
// we can query for entities.
type Pouch interface {
	Queryable
	Storage
}

// Storage is the interface implemented by anything which can
// be interacted with to store or retrieve entities which
// can be found, created, updated and deleted.
type Storage interface {
	FindAll([]Findable) error
	Find(Findable) error

	Create(Createable) error
	CreateAll([]Createable) error

	Update(Updateable) error
	UpdateAll([]Updateable) error

	Delete(Deleteable) error
	DeleteAll([]Deleteable) error
}

// A Query is a direct gateway to interact with the storage and
// retrieval of entities which is backed by a Storage system.
// A Query can also be modified to either add or remove criterions
// that are used to filter searches for entities in the backing
// Storage system.
type Query interface {
	Queryable
	Storage
	FindEntities(Findable, *[]Findable) error
}

// A Creatable entity is one that knows where it is meant to be
// explicitly stored and knows what parts of it need to be stored.
type Createable interface {
	Insertable
	Tableable
}

// A Gettable entity is one that knows what needs to be retrieved
// from an underlying Storage system in order to rebuild itself.
// It also knows where it should be retrieved from inside a Storage
// system.
type Gettable interface {
	GetFieldsFor([]string) []interface{}
	GetAllFields() ([]string, []interface{})
	Tableable
}

type Mergeable interface {
	Merge(Gettable) error
}

// A Tableable entity is one which knows where in a Storage system
// it is supposed to be stored and retrieved from.
type Tableable interface {
	Table() string
}

// A Findable entity is one that is Gettable from a Storage system
// and that knows how to identify itself in the underlying Storage.
type Findable interface {
	Identifiable
	Gettable
	FindableCopy() Findable
	Mergeable
}

// An Insertable entity is one which knows what data from itself
// needs to be stored in an underlying storage system.
type Insertable interface {
	FieldsFor([]string) []interface{}
	InsertableFields() ([]string, []interface{})
	SetIdentifier(interface{}) error
}

// An Identifiable entity is one which knows how to find itself
// in a Storage system.
type Identifiable interface {
	IdentifiableFields() ([]string, []interface{})
}

// An Updateable entity is one which knows how to update itself
// in a Storage system.
type Updateable interface {
	Insertable
	Identifiable
	Tableable
}

// A Deleteable entity is one which knows how to delete only itself
// from a Storage system.
type Deleteable interface {
	Identifiable
	Tableable
}

// Anything that is Queryable knows how to filter queries for itself.
type Queryable interface {
	GroupBy(spec string) Query
	OrderBy(spec string) Query
	Where(frag string, val ...interface{}) Query
	Limit(lim int) Query
	Offset(off int) Query
}

// Executor is a convenience wrapper that allows both *sql.DB and
// *sql.DB to be as Pouches.
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
