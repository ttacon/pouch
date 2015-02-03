package impl

import "github.com/ttacon/pouch"

// NOTE(ttacon):
// The general idea is that users can pass in closures that interact with
// the backing data store.
type dynamicPouch struct {
	l      Logger
	backer interface{}

	// the secret aioli:
	find      func(pouch.Findable, interface{}) error
	findAll   func([]pouch.Findable, interface{}) error
	create    func(pouch.Createable, interface{}) error
	createAll func([]pouch.Createable, interface{}) error
	update    func(pouch.Updateable, interface{}) error
	updateAll func([]pouch.Updateable, interface{}) error
	dlete     func(pouch.Deleteable, interface{}) error
	dleteAll  func([]pouch.Deleteable, interface{}) error
}

type DynamicPouch interface {
	pouch.Pouch
	SetFind(func(pouch.Findable, interface{}) error)
	SetFindAll(func([]pouch.Findable, interface{}) error)
	SetCreate(func(pouch.Createable, interface{}) error)
	SetCreateAll(func([]pouch.Createable, interface{}) error)
	SetUpdate(func(pouch.Updateable, interface{}) error)
	SetUpdateAll(func([]pouch.Updateable, interface{}) error)
	SetDlete(func(pouch.Deleteable, interface{}) error)
	SetDleteAll(func([]pouch.Deleteable, interface{}) error)
}

type DynamicQuery interface {
	pouch.Query
	SetFind(func(pouch.Findable, interface{}) error)
	SetFindAll(func([]pouch.Findable, interface{}) error)
	SetCreate(func(pouch.Createable, interface{}) error)
	SetCreateAll(func([]pouch.Createable, interface{}) error)
	SetUpdate(func(pouch.Updateable, interface{}) error)
	SetUpdateAll(func([]pouch.Updateable, interface{}) error)
	SetDlete(func(pouch.Deleteable, interface{}) error)
	SetDleteAll(func([]pouch.Deleteable, interface{}) error)
}

func NewDynamicPouch(backer interface{}) pouch.Pouch {
	return &dynamicPouch{
		l:      defaultLogger(),
		backer: backer,
	}
}

func (s *dynamicPouch) GroupBy(spec string) pouch.Query {
	return &dynamicFilter{
		groupBySpecs: []string{spec},
		l:            s.l,
	}
}

func (s *dynamicPouch) OrderBy(spec string) pouch.Query {
	return &dynamicFilter{
		orderBySpecs: []string{spec},
		l:            s.l,
	}
}

func (s *dynamicPouch) Where(frag string, vals ...interface{}) pouch.Query {
	return &dynamicFilter{
		constraints: []constraintPair{
			constraintPair{
				frag: frag,
				vals: vals,
			},
		},
		l: s.l,
	}
}

func (s *dynamicPouch) Limit(lim int) pouch.Query {
	return &dynamicFilter{
		limit: lim,
		l:     s.l,
	}
}

func (s *dynamicPouch) Offset(off int) pouch.Query {
	return &dynamicFilter{
		offset: off,
		l:      s.l,
	}
}

func (s *dynamicPouch) Find(i pouch.Findable) error {
	return s.find(i, s.backer)
}

func (s *dynamicPouch) FindAll(fs []pouch.Findable) error {
	return s.findAll(fs, s.backer)
}

func (s *dynamicPouch) Create(i pouch.Createable) error {
	return s.create(i, s.backer)
}

func (s *dynamicPouch) CreateAll(cs []pouch.Createable) error {
	return s.createAll(cs, s.backer)
}

func (s *dynamicPouch) Update(u pouch.Updateable) error {
	return s.update(u, s.backer)
}

func (s *dynamicPouch) UpdateAll(u []pouch.Updateable) error {
	return s.updateAll(u, s.backer)
}

func (s *dynamicPouch) Delete(i pouch.Deleteable) error {
	return s.dlete(i, s.backer)
}

func (s *dynamicPouch) DeleteAll(ds []pouch.Deleteable) error {
	return s.dleteAll(ds, s.backer)
}

////////// SQL pouch.Query implementation //////////
type dynamicFilter struct {
	backer       interface{}
	groupBySpecs []string
	orderBySpecs []string
	constraints  []constraintPair
	limit        int
	offset       int
	l            Logger

	// the secret aioli:
	find      func(pouch.Findable, interface{}) error
	findAll   func([]pouch.Findable, interface{}) error
	create    func(pouch.Createable, interface{}) error
	createAll func([]pouch.Createable, interface{}) error
	update    func(pouch.Updateable, interface{}) error
	updateAll func([]pouch.Updateable, interface{}) error
	dlete     func(pouch.Deleteable, interface{}) error
	dleteAll  func([]pouch.Deleteable, interface{}) error
}

func (s *dynamicFilter) Find(i pouch.Findable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *dynamicFilter) FindAll(fs []pouch.Findable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *dynamicFilter) Create(i pouch.Createable) error {
	return s.create(i, s.backer)
}

func (s *dynamicFilter) CreateAll(cs []pouch.Createable) error {
	return s.createAll(cs, s.backer)
}

func (s *dynamicFilter) Update(u pouch.Updateable) error {
	return s.update(u, s.backer)
}

func (s *dynamicFilter) UpdateAll(us []pouch.Updateable) error {
	return s.updateAll(us, s.backer)
}

func (s *dynamicFilter) Delete(i pouch.Deleteable) error {
	return s.dlete(i, s.backer)
}

func (s *dynamicFilter) DeleteAll(ds []pouch.Deleteable) error {
	return s.dleteAll(ds, s.backer)
}

func (s *dynamicFilter) GroupBy(spec string) pouch.Query {
	s.groupBySpecs = append(s.groupBySpecs, spec)
	return s
}

func (s *dynamicFilter) OrderBy(spec string) pouch.Query {
	s.orderBySpecs = append(s.orderBySpecs, spec)
	return s
}

func (s *dynamicFilter) Where(frag string, vals ...interface{}) pouch.Query {
	s.constraints = append(s.constraints, constraintPair{
		frag: frag,
		vals: vals,
	})
	return s
}

func (s *dynamicFilter) Limit(lim int) pouch.Query {
	s.limit = lim
	return s
}

func (s *dynamicFilter) Offset(off int) pouch.Query {
	s.offset = off
	return s
}

func (s *dynamicFilter) FindEntities(template pouch.Findable, res *[]pouch.Findable) error {
	// TODO(ttacon): do it
	return nil
}

////////// actually making the aioli //////////

func (d *dynamicPouch) SetFind(fn func(pouch.Findable, interface{}) error)          { d.find = fn }
func (d *dynamicPouch) SetFindAll(fn func([]pouch.Findable, interface{}) error)     { d.findAll = fn }
func (d *dynamicPouch) SetCreate(fn func(pouch.Createable, interface{}) error)      { d.create = fn }
func (d *dynamicPouch) SetCreateAll(fn func([]pouch.Createable, interface{}) error) { d.createAll = fn }
func (d *dynamicPouch) SetUpdate(fn func(pouch.Updateable, interface{}) error)      { d.update = fn }
func (d *dynamicPouch) SetUpdateAll(fn func([]pouch.Updateable, interface{}) error) { d.updateAll = fn }
func (d *dynamicPouch) SetDlete(fn func(pouch.Deleteable, interface{}) error)       { d.dlete = fn }
func (d *dynamicPouch) SetDleteAll(fn func([]pouch.Deleteable, interface{}) error)  { d.dleteAll = fn }

func (d *dynamicFilter) SetFind(fn func(pouch.Findable, interface{}) error)      { d.find = fn }
func (d *dynamicFilter) SetFindAll(fn func([]pouch.Findable, interface{}) error) { d.findAll = fn }
func (d *dynamicFilter) SetCreate(fn func(pouch.Createable, interface{}) error)  { d.create = fn }
func (d *dynamicFilter) SetCreateAll(fn func([]pouch.Createable, interface{}) error) {
	d.createAll = fn
}
func (d *dynamicFilter) SetUpdate(fn func(pouch.Updateable, interface{}) error) { d.update = fn }
func (d *dynamicFilter) SetUpdateAll(fn func([]pouch.Updateable, interface{}) error) {
	d.updateAll = fn
}
func (d *dynamicFilter) SetDlete(fn func(pouch.Deleteable, interface{}) error)      { d.dlete = fn }
func (d *dynamicFilter) SetDleteAll(fn func([]pouch.Deleteable, interface{}) error) { d.dleteAll = fn }
