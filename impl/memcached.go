package impl

import (
	"encoding/json"
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ttacon/pouch"
)

type MemcacheKeyable interface {
	MemcacheKey() string
}

type PouchMarshaler interface {
	PouchMarshal() ([]byte, error)
}

type PouchUnmarshaler interface {
	PouchUnmarshal([]byte) error
}

type memcachedPouch struct {
	c *memcache.Client
	l Logger
}

func NewMemcachePouch(c *memcache.Client) pouch.Pouch {
	return &memcachedPouch{
		c: c,
		l: defaultLogger(),
	}
}

func (s *memcachedPouch) GroupBy(spec string) pouch.Query {
	return &memcachedFilter{
		c:            s.c,
		groupBySpecs: []string{spec},
		l:            s.l,
	}
}

func (s *memcachedPouch) OrderBy(spec string) pouch.Query {
	return &memcachedFilter{
		c:            s.c,
		orderBySpecs: []string{spec},
		l:            s.l,
	}
}

func (s *memcachedPouch) Where(frag string, vals ...interface{}) pouch.Query {
	return &memcachedFilter{
		c: s.c,
		constraints: []constraintPair{
			constraintPair{
				frag: frag,
				vals: vals,
			},
		},
		l: s.l,
	}
}

func (s *memcachedPouch) Limit(lim int) pouch.Query {
	return &memcachedFilter{
		c:     s.c,
		limit: lim,
		l:     s.l,
	}
}

func (s *memcachedPouch) Offset(off int) pouch.Query {
	return &memcachedFilter{
		c:      s.c,
		offset: off,
		l:      s.l,
	}
}

func (s *memcachedPouch) Find(i pouch.Findable) error {
	return memcacheFind(s.c, i)
}

func (s *memcachedPouch) FindAll(fs []pouch.Findable) error {
	return memcacheFindAll(s.c, fs)
}

func (s *memcachedPouch) Create(i pouch.Createable) error {
	return memcacheCreate(s.c, i)
}

func (s *memcachedPouch) CreateAll(cs []pouch.Createable) error {
	return memcacheCreateAll(s.c, cs)
}

func (s *memcachedPouch) Update(u pouch.Updateable) error {
	return memcacheUpdate(s.c, u)
}

func (s *memcachedPouch) UpdateAll(us []pouch.Updateable) error {
	return memcacheUpdateAll(s.c, us)
}

func (s *memcachedPouch) Delete(i pouch.Deleteable) error {
	return memcacheDelete(s.c, i)
}

func (s *memcachedPouch) DeleteAll(ds []pouch.Deleteable) error {
	return memcacheDeleteAll(s.c, ds)
}

////////// SQL pouch.Query implementation //////////

type memcachedFilter struct {
	c            *memcache.Client
	groupBySpecs []string
	orderBySpecs []string
	constraints  []constraintPair
	limit        int
	offset       int
	l            Logger
}

func (s *memcachedFilter) Find(i pouch.Findable) error {
	return memcacheFind(s.c, i)
}

func (s *memcachedFilter) FindAll(fs []pouch.Findable) error {
	return memcacheFindAll(s.c, fs)
}

func (s *memcachedFilter) Create(i pouch.Createable) error {
	return memcacheCreate(s.c, i)
}

func (s *memcachedFilter) CreateAll(cs []pouch.Createable) error {
	return memcacheCreateAll(s.c, cs)
}

func (s *memcachedFilter) Update(u pouch.Updateable) error {
	return memcacheUpdate(s.c, u)
}

func (s *memcachedFilter) UpdateAll(us []pouch.Updateable) error {
	return memcacheUpdateAll(s.c, us)
}

func (s *memcachedFilter) Delete(i pouch.Deleteable) error {
	return memcacheDelete(s.c, i)
}

func (s *memcachedFilter) DeleteAll(ds []pouch.Deleteable) error {
	return memcacheDeleteAll(s.c, ds)
}

// NOTE(ttacon): so these don't really do anything
func (s *memcachedFilter) GroupBy(spec string) pouch.Query {
	s.groupBySpecs = append(s.groupBySpecs, spec)
	return s
}

func (s *memcachedFilter) OrderBy(spec string) pouch.Query {
	s.orderBySpecs = append(s.orderBySpecs, spec)
	return s
}

func (s *memcachedFilter) Where(frag string, vals ...interface{}) pouch.Query {
	s.constraints = append(s.constraints, constraintPair{
		frag: frag,
		vals: vals,
	})
	return s
}

func (s *memcachedFilter) Limit(lim int) pouch.Query {
	s.limit = lim
	return s
}

func (s *memcachedFilter) Offset(off int) pouch.Query {
	s.offset = off
	return s
}

func (s *memcachedFilter) FindEntities(template pouch.Findable, res *[]pouch.Findable) error {
	// TODO(ttacon): do it
	return nil
}

////////// memcache helpers //////////
func memcacheCreate(c *memcache.Client, i pouch.Createable) error {
	var key = i.Table()
	if mk, ok := i.(MemcacheKeyable); ok {
		key = mk.MemcacheKey()
	}

	// TODO(ttacon): use marshaller/unmarshaller interface if available
	var dMap = make(map[string]interface{})
	cols, fields := i.InsertableFields()
	if len(cols) != len(fields) || len(cols) == 0 {
		// TODO(ttacon): better error message
		return errors.New("invalid implementation for MemcachePouch")
	}
	for i, col := range cols {
		dMap[col] = fields[i]
	}
	dbytes, err := json.Marshal(dMap)
	if err != nil {
		return err
	}

	return c.Add(&memcache.Item{
		Key:   key,
		Value: dbytes,
	})
}

func memcacheCreateAll(c *memcache.Client, cs []pouch.Createable) error {
	// TODO(ttacon): if we fail halfway through, should we delete the old
	// ones we put in? what if those were updates?
	for _, i := range cs {
		if err := memcacheCreate(c, i); err != nil {
			return err
		}
	}
	return nil
}

func memcacheUpdate(c *memcache.Client, u pouch.Updateable) error {
	var key = u.Table()
	if mk, ok := u.(MemcacheKeyable); ok {
		key = mk.MemcacheKey()
	}

	// TODO(ttacon): use marshaller/unmarshaller interface if available
	var dMap = make(map[string]interface{})
	cols, fields := u.InsertableFields()
	if len(cols) != len(fields) || len(cols) == 0 {
		// TODO(ttacon): better error message
		return errors.New("invalid implementation for MemcachePouch")
	}
	for i, col := range cols {
		dMap[col] = fields[i]
	}
	dbytes, err := json.Marshal(dMap)
	if err != nil {
		return err
	}

	return c.Replace(&memcache.Item{
		Key:   key,
		Value: dbytes,
	})
}

func memcacheUpdateAll(c *memcache.Client, us []pouch.Updateable) error {
	// TODO(ttacon): same question here as in memcacheCreateAll
	// one possible solution is to make another memcache implementation
	// that maintains the state during bulk updates/creates?
	for _, u := range us {
		if err := memcacheUpdate(c, u); err != nil {
			return err
		}
	}
	return nil
}

func memcacheFind(c *memcache.Client, f pouch.Findable) error {
	var key = f.Table()
	if mk, ok := f.(MemcacheKeyable); ok {
		key = mk.MemcacheKey()
	}

	item, err := c.Get(key)
	if err != nil {
		return err
	}

	// TODO(ttacon): use unmarshaller interface if it exists
	var fields = make(map[string]interface{})
	err = json.Unmarshal(item.Value, &fields)
	if err != nil {
		return err
	}

	return f.SetFields(fields)
}

func memcacheFindAll(c *memcache.Client, fs []pouch.Findable) error {
	for _, f := range fs {
		if err := memcacheFind(c, f); err != nil {
			return err
		}
	}
	return nil
}

func memcacheDelete(c *memcache.Client, d pouch.Deleteable) error {
	var key = d.Table()
	if mk, ok := d.(MemcacheKeyable); ok {
		key = mk.MemcacheKey()
	}

	return c.Delete(key)
}

func memcacheDeleteAll(c *memcache.Client, ds []pouch.Deleteable) error {
	for _, d := range ds {
		if err := memcacheDelete(c, d); err != nil {
			return err
		}
	}
	return nil
}
