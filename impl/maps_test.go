package impl

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/ttacon/pouch"
	"github.com/ttacon/pretty"
)

func TestMapPouch_Find(t *testing.T) {
	Convey("when using a map based pouch", t, func() {
		Convey("given a map of foods", func() {
			foods := make(map[string]interface{})
			foods["food:1"] = &mapFood{
				&Food{
					ID:   1,
					Name: "map based",
				},
			}
			p := MapPouch(foods)
			Convey("so we can find a food the normal pouch way", func() {
				m := &mapFood{
					&Food{
						ID: 1,
					},
				}
				err := p.Find(m)
				So(err, ShouldBeNil)
				So(m.Name, ShouldEqual, "map based")
			})
		})
	})
}

func TestMapPouch_FindEntities(t *testing.T) {
	Convey("when using a map based pouch", t, func() {
		Convey("given a map of foods", func() {
			foods := make(map[string]interface{})
			foods["food:1"] = &mapFood{
				&Food{
					ID:   1,
					Name: "map based",
				},
			}
			foods["food:2"] = &mapFood{
				&Food{
					ID:   2,
					Name: "kale",
				},
			}
			foods["food:3"] = &mapFood{
				&Food{
					ID:   3,
					Name: "spinach",
				},
			}
			p := MapPouch(foods)
			Convey("so we can find a bunch foods the normal pouch way", func() {
				m0 := &mapFood{
					&Food{
						ID: 1,
					},
				}
				m1 := &mapFood{
					&Food{
						ID: 2,
					},
				}
				m2 := &mapFood{
					&Food{
						ID: 3,
					},
				}
				var m []pouch.Findable = []pouch.Findable{m0, m1, m2}

				err := p.FindAll(m)
				So(err, ShouldBeNil)
				So(m0.Name, ShouldEqual, "map based")
				So(m1.Name, ShouldEqual, "kale")
				So(m2.Name, ShouldEqual, "spinach")
			})

			Convey("when a map pouch has a limit, it only queries for the first n entities", func() {
				m0 := &mapFood{&Food{ID: 1}}
				m1 := &mapFood{&Food{ID: 2}}
				m2 := &mapFood{&Food{ID: 3}}
				var m []pouch.Findable = []pouch.Findable{m0, m1, m2}
				err := p.Limit(2).FindAll(m)
				So(err, ShouldBeNil)
				So(m0.Name, ShouldEqual, "map based")
				So(m1.Name, ShouldEqual, "kale")
				So(m2.Name, ShouldEqual, "")
			})
		})
	})
}

func TestMapPouch(t *testing.T) {
	Convey("given a map pouch", t, func() {
		foods := make(map[string]interface{})
		p := MapPouch(foods)
		Convey("we can store entities in it", func() {
			f := &mapFood{&Food{Name: "squash"}}
			pretty.Println(p)
			err := p.Create(f)
			pretty.Println(p)
			So(err, ShouldBeNil)
			So(f.ID, ShouldEqual, 0)

			Convey("and if we add another, the IDs increase as expected", func() {
				f := &mapFood{&Food{Name: "yolo"}}
				pretty.Println(p)
				err := p.Create(f)
				pretty.Println(p)
				So(err, ShouldBeNil)
				So(f.ID, ShouldEqual, 1)
			})
			Convey("we can also update entities in it", func() {
				f2 := &mapFood{&Food{ID: f.ID}}
				err := p.Find(f2)
				So(err, ShouldBeNil)
				So(f2.Name, ShouldEqual, "squash")

				f2.Name = "super squash"
				err = p.Update(f2)
				So(err, ShouldBeNil)

				f2 = &mapFood{&Food{ID: f.ID}}
				err = p.Find(f2)
				So(err, ShouldBeNil)
				So(f2.Name, ShouldEqual, "super squash")
			})
			Convey("and we can also delete entities in the map pouch", func() {
				f2 := &mapFood{&Food{ID: f.ID}}
				err := p.Find(f2)
				So(err, ShouldBeNil)
				So(f2.Name, ShouldEqual, "squash")

				err = p.Delete(f2)
				So(err, ShouldBeNil)

				err = p.Find(f2)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestMapPouchBulkOperations(t *testing.T) {
	Convey("given a map pouch", t, func() {
		foods := make(map[string]interface{})
		p := MapPouch(foods)
		Convey("we can bulk store entities in it", func() {
			m0 := &mapFood{&Food{Name: "map based"}}
			m1 := &mapFood{&Food{Name: "kale"}}
			m2 := &mapFood{&Food{Name: "spinach"}}
			var m []pouch.Createable = []pouch.Createable{m0, m1, m2}

			err := p.CreateAll(m)
			So(err, ShouldBeNil)
			So(m0.ID, ShouldEqual, 0)
			So(m1.ID, ShouldEqual, 1)
			So(m2.ID, ShouldEqual, 2)

			Convey("we can bulk find entities", func() {
				m0 := &mapFood{&Food{ID: 0}}
				m1 := &mapFood{&Food{ID: 1}}
				m2 := &mapFood{&Food{ID: 2}}
				var m = []pouch.Findable{m0, m1, m2}

				err := p.FindAll(m)
				So(err, ShouldBeNil)
				So(m0.Name, ShouldEqual, "map based")
				So(m1.Name, ShouldEqual, "kale")
				So(m2.Name, ShouldEqual, "spinach")

				Convey("we can bulk update entities", func() {
					m0.Name = "map based2"
					m1.Name = "kale2"
					m2.Name = "spinach2"
					err = p.UpdateAll([]pouch.Updateable{m0, m1, m2})

					So(err, ShouldBeNil)
					So(m0.Name, ShouldEqual, "map based2")
					So(m1.Name, ShouldEqual, "kale2")
					So(m2.Name, ShouldEqual, "spinach2")
				})

				Convey("and we can bulk delete", func() {
					err = p.DeleteAll([]pouch.Deleteable{m0, m1, m2})
					So(err, ShouldBeNil)

					err = p.FindAll(m)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}
