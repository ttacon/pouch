package impl

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/ttacon/pouch"
)

func TestDynamicPouch(t *testing.T) {
	Convey("given a dynamic pouch", t, func() {
		d := NewDynamicPouch(map[string]string{
			"dynamo:foo": "spinach",
			"dynamo:bar": "kale",
			"dynamo:baz": "mocha",
		})
		Convey("with no defined functions", func() {
			Convey("it should error out when trying to do anything", func() {
				err := d.Find(nil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "no Find function has been defined")
			})
		})

		Convey("with only find defined", func() {
			d.SetFind(func(f pouch.Findable, i interface{}) error {
				var m, ok = i.(map[string]string)
				if !ok {
					return errors.New("unexpected data backer")
				}
				_, vals := f.IdentifiableFields()
				id := fmt.Sprintf(f.Table(), vals...)

				v, ok := m[id]
				if !ok {
					return errors.New("failed to find desired entity")
				}
				fields := f.GetFieldsFor([]string{"col"})
				field := fields[0].(*string)
				*field = v
				return nil
			})

			t0 := &dynamicTestStruct{
				id: "foo",
			}

			err := d.Find(t0)
			So(err, ShouldBeNil)
			So(t0.field, ShouldEqual, "spinach")

			Convey("but if we try to create something, it won't work", func() {
				t0 := &dynamoTestWrapper{
					id:    "yolo",
					field: "yoloBuds",
				}
				err := d.Create(t0)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "no Create function has been defined")
			})
		})
		Convey("with only findAll defined", func() {
			t0 := &dynamicTestStruct{id: "foo"}
			t1 := &dynamicTestStruct{id: "bar"}
			t2 := &dynamicTestStruct{id: "baz"}
			fs := []pouch.Findable{t0, t1, t2}
			err := d.FindAll(fs)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no FindAll function has been defined")

			d.SetFindAll(func(fs []pouch.Findable, i interface{}) error {
				var m, ok = i.(map[string]string)
				if !ok {
					return errors.New("unexpected data backer")
				}

				for _, f := range fs {
					_, vals := f.IdentifiableFields()
					id := fmt.Sprintf(f.Table(), vals...)

					v, ok := m[id]
					if !ok {
						return errors.New("failed to find desired entity")
					}
					fields := f.GetFieldsFor([]string{"col"})
					field := fields[0].(*string)
					*field = v
				}
				return nil
			})

			err = d.FindAll(fs)
			So(err, ShouldBeNil)
			So(t0.field, ShouldEqual, "spinach")
			So(t1.field, ShouldEqual, "kale")
			So(t2.field, ShouldEqual, "mocha")

			Convey("but if we try to create something, it won't work", func() {
				t0 := &dynamoTestWrapper{
					id:    "yolo",
					field: "yoloBuds",
				}
				err := d.Create(t0)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "no Create function has been defined")
			})
		})
		Convey("with only create defined", func() {
			t0 := &dynamoTestWrapper{id: "olaf"}
			err := d.Create(t0)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no Create function has been defined")

			d.SetCreate(func(c pouch.Createable, i interface{}) error {
				var m, ok = i.(map[string]string)
				if !ok {
					return errors.New("unexpected data backer")
				}

				_, fields := c.InsertableFields()
				field, ok := fields[0].(string)
				if !ok {
					return errors.New("invalid insertable field type")
				}

				m[c.Table()] = field
				return nil
			})
			err = d.Create(t0)
			So(err, ShouldBeNil)
		})
		Convey("with only createAll defined", func() {
			t0 := &dynamoTestWrapper{id: "olaf"}
			t1 := &dynamoTestWrapper{id: "bar"}
			t2 := &dynamoTestWrapper{id: "baz"}
			cs := []pouch.Createable{t0, t1, t2}
			err := d.CreateAll(cs)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no CreateAll function has been defined")

			d.SetCreateAll(func(cs []pouch.Createable, i interface{}) error {
				var m, ok = i.(map[string]string)
				if !ok {
					return errors.New("unexpected data backer")
				}
				for _, c := range cs {
					_, fields := c.InsertableFields()
					field, ok := fields[0].(string)
					if !ok {
						return errors.New("invalid insertable field type")
					}

					m[c.Table()] = field
				}
				return nil
			})
			err = d.CreateAll(cs)
			So(err, ShouldBeNil)
		})
		Convey("with only update defined", func() {
			t0 := &dynamoTestWrapper{
				id:    "baz",
				field: "vanillaLatte",
			}
			err := d.Update(t0)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no Update function has been defined")

			d.SetUpdate(func(c pouch.Updateable, i interface{}) error {
				var m, ok = i.(map[string]string)
				if !ok {
					return errors.New("unexpected data backer")
				}

				_, fields := c.InsertableFields()
				field, ok := fields[0].(string)
				if !ok {
					return errors.New("invalid insertable field type")
				}

				if _, ok := m[c.Table()]; !ok {
					return errors.New("entity does not exist, refusing to update")
				}

				m[c.Table()] = field
				return nil
			})

			err = d.Update(t0)
			So(err, ShouldBeNil)
		})
		Convey("with only updateAll defined", func() {
			t0 := &dynamoTestWrapper{
				id:    "baz",
				field: "vanillaLatte",
			}
			t1 := &dynamoTestWrapper{
				id:    "bar",
				field: "blueBerries",
			}
			err := d.UpdateAll([]pouch.Updateable{t0, t1})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no UpdateAll function has been defined")

			d.SetUpdateAll(func(cs []pouch.Updateable, i interface{}) error {
				var m, ok = i.(map[string]string)
				if !ok {
					return errors.New("unexpected data backer")
				}
				for _, c := range cs {
					_, fields := c.InsertableFields()
					field, ok := fields[0].(string)
					if !ok {
						return errors.New("invalid insertable field type")
					}

					if _, ok := m[c.Table()]; !ok {
						return errors.New("entity does not exist, refusing to update")
					}

					m[c.Table()] = field
				}
				return nil
			})

			err = d.UpdateAll([]pouch.Updateable{t0, t1})
			So(err, ShouldBeNil)
		})
		Convey("with only dlete defined", func() {
			t0 := &dynamoTestWrapper{id: "baz"}

			err := d.Delete(t0)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no Delete function has been defined")

			d.SetDlete(func(d pouch.Deleteable, i interface{}) error {
				var m, ok = i.(map[string]string)
				if !ok {
					return errors.New("unexpected data backer")
				}

				if _, ok := m[d.Table()]; !ok {
					return errors.New("entity does not exist, can't delete")
				}

				delete(m, d.Table())
				return nil
			})

			err = d.Delete(t0)
			So(err, ShouldBeNil)
		})
		Convey("with only dleteAll defined", func() {
			t0 := &dynamoTestWrapper{id: "baz"}
			t1 := &dynamoTestWrapper{id: "bar"}

			err := d.DeleteAll([]pouch.Deleteable{t0, t1})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no DeleteAll function has been defined")

			d.SetDleteAll(func(ds []pouch.Deleteable, i interface{}) error {
				var m, ok = i.(map[string]string)
				if !ok {
					return errors.New("unexpected data backer")
				}
				for _, d := range ds {
					if _, ok := m[d.Table()]; !ok {
						return errors.New("entity does not exist, can't delete")
					}

					delete(m, d.Table())
				}
				return nil
			})

			err = d.DeleteAll([]pouch.Deleteable{t0, t1})
			So(err, ShouldBeNil)
		})
	})
}

type dynamicTestStruct struct {
	id    string
	field string
}

func (d *dynamicTestStruct) GetFieldsFor(cols []string) []interface{} {
	return []interface{}{&d.field}
}

func (d *dynamicTestStruct) GetAllFields() ([]string, []interface{}) {
	return nil, nil
}

func (d *dynamicTestStruct) Table() string {
	return "dynamo:%s"
}

func (d *dynamicTestStruct) FindableCopy() pouch.Findable {
	return nil
}

func (d *dynamicTestStruct) IdentifiableFields() ([]string, []interface{}) {
	return nil, []interface{}{d.id}
}

type dynamoTestWrapper dynamicTestStruct

func (d *dynamoTestWrapper) IdentifiableFields() ([]string, []interface{}) {
	return nil, []interface{}{d.id}
}

func (d *dynamoTestWrapper) FieldsFor(s []string) []interface{} {
	return nil
}

func (d *dynamoTestWrapper) InsertableFields() ([]string, []interface{}) {
	return []string{"field"}, []interface{}{d.field}
}

func (d *dynamoTestWrapper) SetIdentifier(i interface{}) error {
	return nil
}

func (d *dynamoTestWrapper) Table() string {
	return "dynamo:" + d.id
}
