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
// can be found, created, updated and deleted. For specific
// information, with respect to error cases and implementation
// details, please see the implementation for your backing
// storage medium.
type Storage interface {
	// Find retrieves an entity from the backing storage medium.
	Find(Findable) error
	// FindAll takes a slice of findable entities (that contain
	// enough information to uniquely identify them in the backing
	// storage medium) and retrieves the rest of each Findable
	// entities information. It should be noted, that the Findable
	// entities to be retrieved to not have to be of the same
	// underlying Go type.
	FindAll([]Findable) error

	// Create stores a Creatable entity in the backing storage medium.
	Create(Createable) error
	// CreateAll stores all given Creatable entities in the backing
	// storage medium. It should be noted that, similarly to FindAll,
	// the underlying types do not need to be the same. HOWEVER,
	// this does mean that for certain backing storage media (i.e. SQL)
	// this function will have to run multiple queries. If you want a bulk
	// insert, please see (YET_TO_BE_IMPLEMENTED).
	CreateAll([]Createable) error

	// Update identifies the Updateable entity in the backing storage
	// media, and updates it the given information that Updateable
	// entity currently contains. Whether or not the provided
	// functionality is a full update or a partial update, is up to
	// underlying Pouch implementation.
	Update(Updateable) error
	// UpdateAll updates all of the provided entities in the backing
	// storage media. Like [Find/Create]All, the underlying types do not
	// have to be the same and as such this function is not meant to
	// guarantee that this is a single interaction for certain media
	// (i.e. SQL based systems). However, this can be (read should be)
	// specified to run inside a transaction.
	UpdateAll([]Updateable) error

	// Delete deletes the given entity from a backing storage media.
	Delete(Deleteable) error
	// DeleteAll removes all given entities from the backing storage media.
	DeleteAll([]Deleteable) error
}

// A Query is a direct gateway to interact with the storage and
// retrieval of entities which is backed by a storage system.
// A Query can also be modified to either add or remove criterions
// that are used to filter searches for entities in the backing
// storage system.
type Query interface {
	Queryable
	Storage

	// FindEntities retrieves all entities that satisfy the Query's
	// current criterions. This can thus also be used to retrieve
	// all entities of a given type, to implement pagination, etc.
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
// *sql.Tx to be used as Pouches.
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
