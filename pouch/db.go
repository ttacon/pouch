package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ttacon/go-utils/db/sqlutil"
	"github.com/ttacon/pouch/pouch/defs"
)

func generateStructsFrom(db *sql.DB) ([]*defs.StructInfo, error) {
	util := sqlutil.New(db)
	tables, err := util.ShowTables("")
	if err != nil {
		return nil, err
	}

	var toGen = make([]*defs.StructInfo, len(tables))

	for i, table := range tables {
		columns, err := util.DescribeTable(table)
		if err != nil {
			return nil, err
		}

		toGen[i] = structFrom(table, columns)
		toGen[i].Table = table
	}

	return toGen, nil
}

func structFrom(name string, columns []sqlutil.ColumnInfo) *defs.StructInfo {
	var fields = make([]defs.FieldInfo, len(columns))
	var idField *string
	for i, column := range columns {
		fields[i] = defs.FieldInfo{
			Name:      column.Field,
			Column:    column.Field,
			IsPointer: column.Null == "YES",
			Type:      goType(column.Type),
		}
		if column.Key == "PRI" {
			// TODO(ttacon): make this not mysql specific
			fields[i].IsPrimaryKey = true
			if column.Extra == "auto_increment" {
				var f = column.Field
				idField = &f
			}
		}
	}
	s := &defs.StructInfo{
		Name:   name,
		Fields: fields,
	}
	if idField != nil {
		s.IDField = *idField
		s.HasAutoGenIDField = true
	}
	return s
}

func goType(typ string) string {
	// for now let's strip from the first '('
	firstParen := strings.Index(typ, "(")
	if firstParen > 0 {
		typ = typ[:firstParen]
	}
	switch typ {
	case "boolean":
		return "bool"
	case "tinyint":
		return "int8"
	case "tinyint unsigned":
		return "uint8"
	case "smallint":
		return "int16"
	case "smallint unsigned":
		return "uint16"
	case "int":
		return "int"
	case "int unsigned":
		return "uint"
	case "bigint":
		return "int64"
	case "bigint unsigned":
		return "uint64"
		// TODO(ttacon): how should we know about float32s?
	case "double":
		return "float64"
	case "mediumblob":
		return "[]uint8"
	case "datetime":
		return "time.Time"
	default:
		return "string"
	}
}

func createTablesFn(db *sql.DB, ts []*defs.StructInfo) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	fmt.Println(dbgenPrmpt, "creating tables:")
	for _, t := range ts {
		var createQuery string
		creatable, ok := interface{}(t).(defs.Creatable)
		if !ok {
			createQuery = reflCreateQuery(t)
		} else {
			createQuery = creatable.CreateQuery()
		}
		fmt.Print("\tcreating table ", t.Name, ": ")
		_, err := db.Exec(createQuery)
		if err != nil {
			fmt.Println(errorX)
			tx.Rollback()
			return err
		}
		fmt.Println(checkY)
	}
	return tx.Commit()
}

func reflCreateQuery(t *defs.StructInfo) string {
	query := "CREATE TABLE " + t.Name + " ("

	// loop through fields and add to query
	for i, field := range t.Fields {
		// query tag for other info like maxsize
		query += fmt.Sprintf("\n\t%s %s %s",
			field.Column, sqlType(field.Type), extraSQLInfo(field))
		if i < len(t.Fields)-1 {
			query += ","
		}
	}

	// what about engine, auto inc start charset?
	// put them on tableInfo?
	return query + "\n);"
}

func sqlType(t string) string {
	switch t {
	case "bool":
		return "boolean"
	case "int8":
		return "tinyint"
	case "uint8":
		return "tinyint unsigned"
	case "int16":
		return "smallint"
	case "uint16":
		return "smallint unsigned"
	case "int":
		return "int"
	case "uint":
		return "int unsigned"
	case "int64":
		return "bigint"
	case "uint64":
		return "bigint unsigned"
	case "float64":
		fallthrough
	case "float32":
		return "double"
	case "[]uint8":
		return "mediumblob"
	case "time.Time":
		return "datetime"
	default:

		return "varchar(255)"
	}
}

func extraSQLInfo(f defs.FieldInfo) string {
	var buf = bytes.NewBuffer(nil)
	if f.IsPrimaryKey {
		buf.WriteString("primary key")
	}

	if !f.IsPointer {
		buf.WriteString("not null")
	}

	return buf.String()
}
