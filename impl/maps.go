package impl

import (
	"errors"
	"fmt"

	"github.com/ttacon/pouch"
)

////////// Map Pouch implementation //////////
type mapPouch struct {
	db map[string]interface{}
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
	// TODO(ttaco): do it
	return nil
}

func (s *mapPouch) CreateAll(cs []pouch.Createable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapPouch) Update(u pouch.Updateable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapPouch) UpdateAll(u []pouch.Updateable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapPouch) Delete(i pouch.Deleteable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapPouch) DeleteAll(ds []pouch.Deleteable) error {
	// TODO(ttaco): do it
	return nil
}

////////// SQL pouch.Query implementation //////////

type mapQuery struct {
	db           map[string]interface{}
	groupBySpecs []string
	orderBySpecs []string
	constraints  []constraintPair
	limit        int
	offset       int
}

func (s *mapQuery) Find(i pouch.Findable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapQuery) FindAll(fs []pouch.Findable) error {
	var temp = fs
	if s.limit > 0 && s.limit < len(fs) {
		temp = temp[0:s.limit]
	}
	return findAllEntries(s.db, temp)
}

func (s *mapQuery) Create(i pouch.Createable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapQuery) CreateAll(cs []pouch.Createable) error {
	// TODO(ttaco): do it
	return nil
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
