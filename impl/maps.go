package impl

import (
	"errors"
	"fmt"

	"github.com/ttacon/pouch"
)

////////// Map Pouch implementation //////////
type mapPouch struct {
	db     map[string]interface{}
	nextID int
}

func MapPouch(m map[string]interface{}) pouch.Pouch {
	return &mapPouch{
		db: m,
	}
}

func (s *mapPouch) GroupBy(spec string) pouch.Query {
	return &mapQuery{
		db:           s.db,
		groupBySpecs: []string{spec},
	}
}

func (s *mapPouch) OrderBy(spec string) pouch.Query {
	return &mapQuery{
		db:           s.db,
		orderBySpecs: []string{spec},
	}
}

func (s *mapPouch) Where(frag string, vals ...interface{}) pouch.Query {
	return &mapQuery{
		db: s.db,
		constraints: []constraintPair{
			constraintPair{
				frag: frag,
				vals: vals,
			},
		},
	}
}

func (s *mapPouch) Limit(lim int) pouch.Query {
	return &mapQuery{
		db:    s.db,
		limit: lim,
	}
}

func (s *mapPouch) Offset(off int) pouch.Query {
	return &mapQuery{
		db:     s.db,
		offset: off,
	}
}

func (s *mapPouch) Find(i pouch.Findable) error {
	return findMapEntry(s.db, i)
}

func (s *mapPouch) FindAll(fs []pouch.Findable) error {
	return findAllEntries(s.db, fs)
}

func (s *mapPouch) Create(i pouch.Createable) error {
	return createMapEntry(s.db, &s.nextID, i)
}

func (s *mapPouch) CreateAll(cs []pouch.Createable) error {
	return createAllMapEntries(s.db, &s.nextID, cs)
}

func (s *mapPouch) Update(u pouch.Updateable) error {
	return updateMapEntry(s.db, u)
}

func (s *mapPouch) UpdateAll(u []pouch.Updateable) error {
	return updateAllMapEntries(s.db, u)
}

func (s *mapPouch) Delete(i pouch.Deleteable) error {
	return deleteMapEntry(s.db, i)
}

func (s *mapPouch) DeleteAll(ds []pouch.Deleteable) error {
	return deleteAllMapEntries(s.db, ds)
}

////////// SQL pouch.Query implementation //////////

type mapQuery struct {
	db           map[string]interface{}
	groupBySpecs []string
	orderBySpecs []string
	constraints  []constraintPair
	limit        int
	offset       int
	nextID       int
}

func (s *mapQuery) Find(i pouch.Findable) error {
	return findMapEntry(s.db, i)
}

func (s *mapQuery) FindAll(fs []pouch.Findable) error {
	var temp = fs
	if s.limit > 0 && s.limit < len(fs) {
		temp = temp[0:s.limit]
	}
	return findAllEntries(s.db, temp)
}

func (s *mapQuery) Create(i pouch.Createable) error {
	// TODO(ttacon): this isn't go routine safe, s.nextID either needs a lock
	// of we pass it in and assume no errors occur
	return createMapEntry(s.db, &s.nextID, i)
}

func (s *mapQuery) CreateAll(cs []pouch.Createable) error {
	return createAllMapEntries(s.db, &s.nextID, cs)
}

func (s *mapQuery) Update(u pouch.Updateable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapQuery) UpdateAll(us []pouch.Updateable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapQuery) Delete(i pouch.Deleteable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapQuery) DeleteAll(ds []pouch.Deleteable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapQuery) GroupBy(spec string) pouch.Query {
	s.groupBySpecs = append(s.groupBySpecs, spec)
	return s
}

func (s *mapQuery) OrderBy(spec string) pouch.Query {
	s.orderBySpecs = append(s.orderBySpecs, spec)
	return s
}

func (s *mapQuery) Where(frag string, vals ...interface{}) pouch.Query {
	s.constraints = append(s.constraints, constraintPair{
		frag: frag,
		vals: vals,
	})
	return s
}

func (s *mapQuery) Limit(lim int) pouch.Query {
	s.limit = lim
	return s
}

func (s *mapQuery) Offset(off int) pouch.Query {
	s.offset = off
	return s
}

func (s *mapQuery) FindEntities(template pouch.Findable, res *[]pouch.Findable) error {
	// TODO(ttaco): do it
	return nil
}

////////// help fns //////////

func findMapEntry(m map[string]interface{}, i pouch.Findable) error {
	table := i.Table()
	if len(table) == 0 {
		return errors.New("can't find an entity that has no table/key format")
	}

	_, ifies := i.IdentifiableFields()
	key := fmt.Sprintf(table, ifies...)
	v, ok := m[key]
	if !ok {
		return errors.New("could not find entity: " + key)
	}

	g, ok := v.(pouch.Gettable)
	if !ok {
		return errors.New("cannot merge with non-mergeable entity")
	}

	return i.Merge(g)
}

func createMapEntry(m map[string]interface{}, nextID *int, c pouch.Createable) error {
	table := c.Table()
	if len(table) == 0 {
		return errors.New("can't find an entity that has no table/key format")
	}
	nxt := *nextID
	key := fmt.Sprintf(table, nxt)
	*nextID = nxt + 1
	m[key] = c

	return c.SetIdentifier(nxt)
}

func updateMapEntry(m map[string]interface{}, u pouch.Updateable) error {
	// TODO(ttacon): offer merges or just DESTROY the existing entity?
	table := u.Table()
	if len(table) == 0 {
		return errors.New("can't find an entity that has no table/key format")
	}

	_, ifies := u.IdentifiableFields()
	key := fmt.Sprintf(table, ifies...)
	_, ok := m[key]
	if !ok {
		return errors.New("could not find entity: " + key)
	}

	m[key] = u
	return nil
}

func deleteMapEntry(m map[string]interface{}, d pouch.Deleteable) error {
	table := d.Table()
	if len(table) == 0 {
		return errors.New("can't find an entity that has no table/key format")
	}

	_, ifies := d.IdentifiableFields()
	key := fmt.Sprintf(table, ifies...)
	delete(m, key)
	return nil
}

////////// bulk helpers //////////
func findAllEntries(m map[string]interface{}, fs []pouch.Findable) error {
	// TODO(ttacon): try all of them or return first error?
	// or return bulk error that is simply an "error"?
	// (I like the second option)
	var err error
	for _, f := range fs {
		terr := findMapEntry(m, f)
		if terr != nil {
			// write now overwrite if non-nil, will make this bulk soon
			err = terr
		}
	}
	return err
}

func createAllMapEntries(m map[string]interface{}, nextID *int, cs []pouch.Createable) error {
	for _, c := range cs {
		if err := createMapEntry(m, nextID, c); err != nil {
			// TODO(ttacon): or should we just keep trying?
			return err
		}
	}
	return nil
}

func updateAllMapEntries(m map[string]interface{}, us []pouch.Updateable) error {
	var err error
	for _, u := range us {
		if terr := updateMapEntry(m, u); terr != nil {
			err = terr
		}
	}
	return err
}

func deleteAllMapEntries(m map[string]interface{}, ds []pouch.Deleteable) error {
	var err error
	for _, d := range ds {
		if terr := deleteMapEntry(m, d); terr != nil {
			err = terr
		}
	}
	return err
}
