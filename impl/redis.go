package impl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ttacon/pouch"
	"github.com/ttacon/pouch/impl/piface"
)

type redisPouch struct {
	db piface.RedisFace
}

func RedisPouch(db piface.RedisFace) pouch.Pouch {
	return &redisPouch{
		db: db,
	}
}

func (r *redisPouch) GroupBy(spec string) pouch.Query {
	// TODO(ttacon): do it
	return &redisQuery{
		db: r.db,
	}
}

func (r *redisPouch) OrderBy(spec string) pouch.Query {
	// TODO(ttacon): do it
	return &redisQuery{
		db: r.db,
	}
}

func (r *redisPouch) Where(frag string, vals ...interface{}) pouch.Query {
	// TODO(ttacon): do it
	return &redisQuery{
		db: r.db,
	}
}

func (r *redisPouch) Limit(lim int) pouch.Query {
	return &redisQuery{
		db:    r.db,
		limit: lim,
	}
}

func (r *redisPouch) Offset(off int) pouch.Query {
	return &redisQuery{
		db:     r.db,
		offset: off,
	}
}

func (r *redisPouch) Find(i pouch.Findable) error {
	var key string
	v, redisDecorated := i.(pouch.RedisDecorated)
	if !redisDecorated {
		key = i.Table()
		_, vals := i.IdentifiableFields()
		if len(vals) == 0 {
			return errors.New("must provide values to find entity by")
		}

		var gnr8d string = fmt.Sprintf(key, vals...)
		if gnr8d == key || strings.Contains(gnr8d, "%s") {
			return errors.New("invalid generated key: " + gnr8d)
		}
		key = gnr8d
	} else {
		key = v.KeyFormula()
	}

	if len(key) == 0 {
		return errors.New("must return an entity formula to be able to retrieve")
	}

	var keys, fields = i.GetAllFields()
	if len(keys) != len(fields) {
		return errors.New("invalid number of fields/keys returned")
	}

	foundFields, err := r.db.Hmget(key, keys[0], keys[1:]...)
	if err != nil {
		return err
	}

	if redisDecorated {
		for i, fieldKey := range keys {
			if err = v.SetFieldFromString(fieldKey, foundFields[i]); err != nil {
				return err
			}
		}
		return nil
	}

	for i, field := range fields {
		s, ok := field.(*string)
		if !ok {
			// TODO(ttacon): make error point out RedisDecorated
			return errors.New("can't store string in non-string type")
		}
		*s = foundFields[i]
	}

	return nil
}

var ErrUnimplemented = errors.New("pouch: not implemented")

func (r *redisPouch) FindAll(fs []pouch.Findable) error {
	// TODO(ttaco): do it
	return ErrUnimplemented
}

func (r *redisPouch) Create(i pouch.Createable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisPouch) CreateAll(cs []pouch.Createable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisPouch) Update(u pouch.Updateable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisPouch) UpdateAll(u []pouch.Updateable) error {
	// TODO(ttaco): do it
	return ErrUnimplemented
}

func (r *redisPouch) Delete(i pouch.Deleteable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisPouch) DeleteAll(ds []pouch.Deleteable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

////////// Redis pouch.Query implementation //////////

type redisQuery struct {
	db           piface.RedisFace
	groupBySpecs []string
	orderBySpecs []string
	limit        int
	offset       int
}

func (r *redisQuery) Find(i pouch.Findable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisQuery) FindAll(fs []pouch.Findable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisQuery) Create(i pouch.Createable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisQuery) CreateAll(cs []pouch.Createable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisQuery) Update(u pouch.Updateable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisQuery) UpdateAll(us []pouch.Updateable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisQuery) Delete(i pouch.Deleteable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisQuery) DeleteAll(ds []pouch.Deleteable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}

func (r *redisQuery) GroupBy(spec string) pouch.Query {
	// TODO(ttacon): do it
	return r
}

func (r *redisQuery) OrderBy(spec string) pouch.Query {
	// TODO(ttacon): do it
	return r
}

func (r *redisQuery) Where(frag string, vals ...interface{}) pouch.Query {
	// TODO(ttacon): do it
	return r
}

func (r *redisQuery) Limit(lim int) pouch.Query {
	// TODO(ttacon): do it
	return r
}

func (r *redisQuery) Offset(off int) pouch.Query {
	// TODO(ttacon): do it
	return r
}

func (r *redisQuery) FindEntities(template pouch.Findable, res *[]pouch.Findable) error {
	// TODO(ttacon): do it
	return ErrUnimplemented
}
