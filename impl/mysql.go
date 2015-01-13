package impl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ttacon/builder"
	"github.com/ttacon/pouch"
)

////////// SQL Pouch implementation //////////
type sqlPouch struct {
	db pouch.Executor
}

func SQLPouch(db pouch.Executor) pouch.Pouch {
	return &sqlPouch{
		db: db,
	}
}

func (s *sqlPouch) GroupBy(spec string) pouch.Query {
	return &sqlQuery{
		db:           s.db,
		groupBySpecs: []string{spec},
	}
}

func (s *sqlPouch) OrderBy(spec string) pouch.Query {
	return &sqlQuery{
		db:           s.db,
		orderBySpecs: []string{spec},
	}
}

func (s *sqlPouch) Where(frag string, vals ...interface{}) pouch.Query {
	return &sqlQuery{
		db: s.db,
		constraints: []constraintPair{
			constraintPair{
				frag: frag,
				vals: vals,
			},
		},
	}
}

func (s *sqlPouch) Limit(lim int) pouch.Query {
	return &sqlQuery{
		db:    s.db,
		limit: lim,
	}
}

func (s *sqlPouch) Offset(off int) pouch.Query {
	return &sqlQuery{
		db:     s.db,
		offset: off,
	}
}

func (s *sqlPouch) Find(i pouch.Findable) error {
	return findEntity(s.db, i, "", nil)
}

func (s *sqlPouch) FindAll(fs []pouch.Findable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *sqlPouch) Create(i pouch.Createable) error {
	return createEntity(s.db, i, "")
}

func (s *sqlPouch) CreateAll(cs []pouch.Createable) error {
	return createAll(s.db, cs)
}

func (s *sqlPouch) Update(u pouch.Updateable) error {
	return updateEntity(s.db, u, "")
}

func (s *sqlPouch) UpdateAll(u []pouch.Updateable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *sqlPouch) Delete(i pouch.Deleteable) error {
	return deleteEntity(s.db, i, "")
}

func (s *sqlPouch) DeleteAll(ds []pouch.Deleteable) error {
	return deleteAll(s.db, ds)
}

// TODO(ttacon): reuse these as we add other dialects
func findEntity(db pouch.Executor, i pouch.Findable, rest string, ps []interface{}) error {
	cols, fields := i.GetAllFields()
	if len(cols) == 0 || len(fields) == 0 {
		return errors.New("must provide columns to select from")
	}

	table := i.Table()
	if len(table) == 0 {
		return errors.New("entity is not known to map to any table")
	}

	var query = builder.NewBuilderString("select ")
	for i, col := range cols {
		if i > 0 {
			query.WriteString(",\n  ")
		}
		query.WriteString(col)
	}
	query.WriteString("\nfrom " + table + "\n")

	ids, vals := i.IdentifiableFields()
	if len(ids) == 0 || len(vals) == 0 {
		return errors.New("no identifying information for entity")
	}

	ps = append(vals, ps...)
	query.WriteString("where ")
	for i, id := range ids {
		if i > 0 && i < len(ids)-1 {
			query.WriteString(" AND ")
		}
		query.WriteString(id + " = ? ")
	}

	if len(rest) > 0 {
		query.WriteString(" AND " + rest)
	}

	row := db.QueryRow(query.String(), ps...)
	return row.Scan(fields...)
}

func createEntity(db pouch.Executor, i pouch.Createable, rest string) error {
	var cols, vals = i.InsertableFields()
	if len(cols) == 0 || len(vals) == 0 {
		return errors.New("cannot insert empty entity")
	}
	if len(cols) != len(vals) {
		return errors.New("[inserting], there cannot be more columns than values")
	}

	placeholders := "?"
	if len(vals) > 1 {
		placeholders += strings.Repeat(", ?", len(vals)-1)
	}

	table := i.Table()
	if len(table) == 0 {
		return errors.New("this entity is not known to be associated with any table")
	}

	var query = builder.NewBuilderString("insert into " + table)
	query.WriteString("(\n" + strings.Join(cols, ", ") + "\n) values ")
	query.WriteString("(\n" + placeholders + "\n)")
	res, err := db.Exec(query.String(), vals...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	return i.SetIdentifier(id)
}

func updateEntity(db pouch.Executor, u pouch.Updateable, rest string) error {
	var cols, vals = u.InsertableFields()
	if len(cols) == 0 || len(vals) == 0 {
		return errors.New("cannot insert empty entity")
	}
	if len(cols) != len(vals) {
		return errors.New("[inserting], there cannot be more columns than values")
	}

	table := u.Table()
	if len(table) == 0 {
		return errors.New("this entity is not known to be associated with any table")
	}

	ids, idVals := u.IdentifiableFields()
	if len(ids) == 0 || len(idVals) == 0 {
		return errors.New("no identifying information for entity")
	}

	var query = builder.NewBuilderString("update " + table + "\nset ")
	for i, col := range cols {
		if i > 0 && i < len(cols) {
			query.WriteString(", ")
		}
		query.WriteString(col + " = ?")
	}

	vals = append(vals, idVals...)
	query.WriteString("\nwhere ")
	for i, id := range ids {
		if i > 0 && i < len(ids)-1 {
			query.WriteString(" AND ")
		}
		query.WriteString(id + " = ? ")
	}

	_, err := db.Exec(query.String(), vals...)
	return err
}

func deleteEntity(db pouch.Executor, d pouch.Deleteable, rest string) error {
	table := d.Table()
	if len(table) == 0 {
		return errors.New("this entity is not known to be associated with any table")
	}

	ids, idVals := d.IdentifiableFields()
	if len(ids) == 0 || len(idVals) == 0 {
		return errors.New("no identifying information for entity")
	}

	var query = builder.NewBuilderString("delete\nfrom " + table + "\nwhere ")
	for i, id := range ids {
		if i > 0 && i < len(ids) {
			query.WriteString(", ")
		}
		query.WriteString(id + " = ?")
	}

	_, err := db.Exec(query.String(), idVals...)
	return err
}

////////// SQL pouch.Query implementation //////////

type sqlQuery struct {
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

func (s *sqlQuery) Find(i pouch.Findable) error {
	rest, vals := buildConstraints(s)
	return findEntity(s.db, i, rest, vals)
}

func (s *sqlQuery) FindAll(fs []pouch.Findable) error {
	rest, ps := buildConstraints(s)
	return findAll(s.db, fs, rest, ps)
}

func (s *sqlQuery) Create(i pouch.Createable) error {
	rest, _ := buildConstraints(s)
	return createEntity(s.db, i, rest)
}

func (s *sqlQuery) CreateAll(cs []pouch.Createable) error {
	return createAll(s.db, cs)
}

func (s *sqlQuery) Update(u pouch.Updateable) error {
	rest, _ := buildConstraints(s)
	return updateEntity(s.db, u, rest)
}

func (s *sqlQuery) UpdateAll(us []pouch.Updateable) error {
	// TODO(ttaco): do it
	return nil
}

func (s *sqlQuery) Delete(i pouch.Deleteable) error {
	rest, _ := buildConstraints(s)
	return deleteEntity(s.db, i, rest)
}

func (s *sqlQuery) DeleteAll(ds []pouch.Deleteable) error {
	return deleteAll(s.db, ds)
}

func (s *sqlQuery) GroupBy(spec string) pouch.Query {
	s.groupBySpecs = append(s.groupBySpecs, spec)
	return s
}

func (s *sqlQuery) OrderBy(spec string) pouch.Query {
	s.orderBySpecs = append(s.orderBySpecs, spec)
	return s
}

func (s *sqlQuery) Where(frag string, vals ...interface{}) pouch.Query {
	s.constraints = append(s.constraints, constraintPair{
		frag: frag,
		vals: vals,
	})
	return s
}

func (s *sqlQuery) Limit(lim int) pouch.Query {
	s.limit = lim
	return s
}

func (s *sqlQuery) Offset(off int) pouch.Query {
	s.offset = off
	return s
}

func (s *sqlQuery) FindEntities(template pouch.Findable, res *[]pouch.Findable) error {
	rest, ps := buildConstraints(s)
	return findEntities(s.db, template, res, rest, ps)
}

//TODO(ttacon): add HAVING
func buildConstraints(s *sqlQuery) (string, []interface{}) {
	var constraints = builder.NewBuilder(nil)
	var vals []interface{}
	for i, constraint := range s.constraints {
		if i > 0 && i < len(s.constraints)-1 {
			// TODO(ttacon): decent way to specify AND vs OR
			constraints.WriteString(" AND ")
		}
		constraints.WriteString(constraint.frag)
		constraints.WriteString("\n")
		vals = append(vals, constraint.vals...)
	}

	for i, group := range s.groupBySpecs {
		if i == 0 {
			constraints.WriteString("group by ")
		}
		if i > 0 && i < len(s.groupBySpecs)-1 {
			constraints.WriteString(", ")
		}
		constraints.WriteString(group)
		constraints.WriteString("\n")
	}

	for i, order := range s.orderBySpecs {
		if i == 0 {
			constraints.WriteString("order by ")
		}
		if i > 0 && i < len(s.orderBySpecs)-1 {
			constraints.WriteString(", ")
		}
		constraints.WriteString(order)
		constraints.WriteString("\n")
	}

	if s.limit > 0 {
		if s.offset > 0 {
			constraints.WriteString(fmt.Sprintf("limit %d, %d", s.offset, s.limit))
		} else {
			constraints.WriteString(fmt.Sprintf("limit %d", s.limit))
		}
	}
	return constraints.String(), vals
}

////////// *All functions //////////
func findAll(db pouch.Executor, fs []pouch.Findable, rest string, ps []interface{}) error {
	// it assumes fs is full of entities who know their identifying info
	if len(fs) == 0 {
		// how's this for a cryptic error lol
		return errors.New("cannot find non-existent entities")
	}

	for _, i := range fs {
		table := i.Table()
		if len(table) == 0 {
			return errors.New("entity is not known to map to any table")
		}

		cols, fields := i.GetAllFields()
		if len(cols) == 0 || len(fields) == 0 {
			return errors.New("must provide columns to select from")
		}

		var query = builder.NewBuilderString("select ")
		for i, col := range cols {
			if i > 0 {
				query.WriteString(",\n  ")
			}
			query.WriteString(col)
		}
		query.WriteString("\nfrom " + table + "\n")

		ids, vals := i.IdentifiableFields()
		if len(ids) == 0 || len(vals) == 0 {
			return errors.New("no identifying information for entity")
		}

		ps = append(vals, ps...)
		query.WriteString("where ")
		for i, id := range ids {
			if i > 0 && i < len(ids)-1 {
				query.WriteString(" AND ")
			}
			query.WriteString(id + " = ? ")
		}

		if len(rest) > 0 {
			query.WriteString(" AND " + rest)
		}

		row := db.QueryRow(query.String(), ps...)

		if err := row.Scan(fields...); err != nil {
			return err
		}
	}
	return nil
}

func createAll(db pouch.Executor, cs []pouch.Createable) error {
	if len(cs) == 0 {
		return errors.New("no entities to insert (empty slice)")
	}

	// TODO(ttacon): move the validation logic out of the loop
	// (prevalidate everything, will help)

	for _, i := range cs {
		var cols, vals = i.InsertableFields()
		if len(cols) == 0 || len(vals) == 0 {
			return errors.New("cannot insert empty entity")
		}
		if len(cols) != len(vals) {
			return errors.New("[inserting], there cannot be more columns than values")
		}

		placeholders := "?"
		if len(vals) > 1 {
			placeholders += strings.Repeat(", ?", len(vals)-1)
		}

		table := i.Table()
		if len(table) == 0 {
			return errors.New("this entity is not known to be associated with any table")
		}

		var query = builder.NewBuilderString("insert into " + table)
		query.WriteString("(\n" + strings.Join(cols, ", ") + "\n) values ")
		query.WriteString("(\n" + placeholders + "\n)")
		res, err := db.Exec(query.String(), vals...)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		err = i.SetIdentifier(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteAll(db pouch.Executor, ds []pouch.Deleteable) error {
	if len(ds) == 0 {
		return errors.New("[deleteAll] no entities to delete")
	}

	for _, d := range ds {
		table := d.Table()
		if len(table) == 0 {
			return errors.New("this entity is not known to be associated with any table")
		}

		ids, idVals := d.IdentifiableFields()
		if len(ids) == 0 || len(idVals) == 0 {
			return errors.New("no identifying information for entity")
		}

		var query = builder.NewBuilderString("delete\nfrom " + table + "\nwhere ")
		for i, id := range ids {
			if i > 0 && i < len(ids) {
				query.WriteString(", ")
			}
			query.WriteString(id + " = ?")
		}

		_, err := db.Exec(query.String(), idVals...)
		if err != nil {
			return err
		}
	}
	return nil
}

////////// findEntities //////////
func findEntities(
	db pouch.Executor,
	example pouch.Findable,
	fs *[]pouch.Findable,
	rest string,
	ps []interface{}) error {

	table := example.Table()
	if len(table) == 0 {
		return errors.New("entity is not known to map to any table")
	}

	cols, _ := example.GetAllFields()
	if len(cols) == 0 {
		return errors.New("must provide columns to select from")
	}

	var query = builder.NewBuilderString("select ")
	for i, col := range cols {
		if i > 0 {
			query.WriteString(",\n  ")
		}
		query.WriteString(col)
	}
	query.WriteString("\nfrom " + table + "\n")

	if len(rest) > 0 {
		query.WriteString("WHERE " + rest)
	}

	fmt.Println("query: ", query.String())
	rows, err := db.Query(query.String(), ps...)
	if err != nil {
		return err
	}

	for rows.Next() {
		cop := example.FindableCopy()
		fields := cop.GetFieldsFor(cols)
		err = rows.Scan(fields...)
		if err != nil {
			return err
		}
		*fs = append(*fs, cop)
	}
	return rows.Err()
}
