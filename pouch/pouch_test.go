package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/ttacon/go-utils/db/sqlutil"
) // not a fan of these style imports

func Test_dbInfoProvided(t *testing.T) {
	Convey("When checking if we need to contact a database or not", t, func() {
		Convey("Given db login info", func() {
			pieces := []string{"", "pouch", "pouch", ""}
			Convey("we should know we're in db retrieve mode", func() {
				So(dbInfoProvided(pieces...), ShouldEqual, true)
			})
		})

		Convey("Not given any db login info", func() {
			pieces := []string{}
			Convey("we should know we're not supposed to touch a db", func() {
				So(dbInfoProvided(pieces...), ShouldEqual, false)
			})
		})
	})
}

func Test_structFrom(t *testing.T) {
	deliciousness := sqlutil.ColumnInfo{
		Field: "Deliciousness",
		Null:  "YES",
		Type:  "int",
	}
	auto_id := sqlutil.ColumnInfo{
		Field: "ID",
		Null:  "NO",
		Type:  "int",
		Key:   "PRI",
		Extra: "auto_increment",
	}
	id := sqlutil.ColumnInfo{
		Field: "ID",
		Null:  "NO",
		Type:  "int",
		Key:   "PRI",
		Extra: "auto_increment",
	}
	name := sqlutil.ColumnInfo{
		Field: "Name",
		Null:  "NO",
		Type:  "varchar(64)",
	}
	Convey("When generating struct info from db information", t, FailureContinues, func() {
		Convey("A table named 'Food'", func() {
			tableName := "Food"
			Convey("With only one column", func() {
				columns := []sqlutil.ColumnInfo{deliciousness}
				Convey("Should have the correct structure", func() {
					st := structFrom(tableName, columns)
					So(st.Name, ShouldEqual, tableName)
					So(st.IDField, ShouldBeEmpty)
					So(len(st.Fields), ShouldEqual, 1)
					So(st.Fields[0].Name, ShouldEqual, "Deliciousness")
					So(st.Fields[0].Column, ShouldEqual, "Deliciousness")
					So(st.Fields[0].IsPointer, ShouldBeTrue)
					So(st.Fields[0].Type, ShouldEqual, "int")
				})
			})
			Convey("With two columns, one an id", func() {
				columns := []sqlutil.ColumnInfo{id, deliciousness}
				Convey("Should have the correct structure", func() {
					_ = structFrom(tableName, columns)
					// TODO(ttacon): do it
				})
			})
			Convey("With two columns, one an 'auto_increment'ed id", func() {
				columns := []sqlutil.ColumnInfo{auto_id, deliciousness}
				Convey("Should have the correct structure", func() {
					_ = structFrom(tableName, columns)
					// TODO(ttacon): do it
				})
			})
			Convey("With a bunch of columns", func() {
				columns := []sqlutil.ColumnInfo{auto_id, name, deliciousness}
				Convey("Should have the correct structure", func() {
					_ = structFrom(tableName, columns)
					// TODO(ttacon): do it
				})
			})
		})
	})
}

func Test_loadTemplates(t *testing.T) {
	Convey("Loading templates should be errorless", t, func() {
		So(loadTemplates(), ShouldBeNil)
	})
}

func Test_goType(t *testing.T) {
	tests := [][]string{
		[]string{"boolean", "bool"},
		[]string{"tinyint", "int8"},
		[]string{"tinyint unsigned", "uint8"},
		[]string{"smallint", "int16"},
		[]string{"smallint unsigned", "uint16"},
		[]string{"int", "int"},
		[]string{"int unsigned", "uint"},
		[]string{"bigint", "int64"},
		[]string{"bigint unsigned", "uint64"},
		[]string{"double", "float64"},
		[]string{"mediumblob", "[]uint8"},
		[]string{"datetime", "time.Time"}, // TODO(ttacon): should this be *time.Time?
		[]string{"tinyblob", "string"},
		[]string{"varchar(5)", "string"},
		[]string{"varchar(255)", "string"},
	}

	Convey("When identifying the corresponding go type", t, func() {
		for _, test := range tests {
			Convey(test[0]+" -> "+test[1], func() {
				So(goType(test[0]), ShouldEqual, test[1])
			})
		}
	})
}
