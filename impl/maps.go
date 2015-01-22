package impl

import "github.com/ttacon/pouch"

////////// SQL Pouch implementation //////////
type mapPouch struct {
	db map[string]interface{}
}

func MapPouch() pouch.Pouch {
	return &mapPouch{
		db: make(map[string]interface{}),
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
	// TODO(ttaco): do it
	return nil
}

func (s *mapPouch) FindAll(fs []pouch.Findable) error {
	// TODO(ttaco): do it
	return nil
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
	db           pouch.Executor
	groupBySpecs []string
	orderBySpecs []string
	constraints  []constraintPair
	limit        int
	offset       int
}

type constraintPair struct {
	frag string
	vals []interface{}
}

func (s *mapQuery) Find(i pouch.Findable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *mapQuery) FindAll(fs []pouch.Findable) error {
	// TODO(ttaco): do it
	return nil
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
